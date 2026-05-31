package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// EnforcementRegistryLintError is one finding from Phase 3 lint.
//
// Format follows Phase 7 scenario contract:
//
//	LINT ERROR [<type>]
//	  <key>: <value>
//	  ...
//	  message: <one-line explanation + remediation>
type EnforcementRegistryLintError struct {
	Type    string
	Fields  []EnforcementRegistryLintField
	Message string
}

type EnforcementRegistryLintField struct {
	Key   string
	Value string
}

func (e EnforcementRegistryLintError) Format() string {
	var b strings.Builder
	fmt.Fprintf(&b, "LINT ERROR [%s]\n", e.Type)
	for _, f := range e.Fields {
		fmt.Fprintf(&b, "  %s: %s\n", f.Key, f.Value)
	}
	fmt.Fprintf(&b, "  message: %s", e.Message)
	return b.String()
}

// registrySnapshot is a parsed view of enforcement-registry.yaml used by lint.
type registrySnapshot struct {
	EnforcementMode struct {
		OrphanRule     string `yaml:"orphan_rule"`
		OrphanExecutor string `yaml:"orphan_executor"`
	} `yaml:"enforcement_mode"`
	ExecutorKindSpec struct {
		BindingRequiredFor []string `yaml:"binding_required_for"`
	} `yaml:"executor_kind_spec"`
	InternalHelperAllowlist struct {
		Symbols []string `yaml:"symbols"`
	} `yaml:"internal_helper_allowlist"`
	RuleClasses []registryRuleClass `yaml:"rule_classes"`
}

type registryRuleClass struct {
	ID                string             `yaml:"id"`
	Coverage          string             `yaml:"coverage"`
	SourceFiles       []string           `yaml:"source_files"`
	Executors         []registryExecutor `yaml:"executors"`
	ExecutorsPlanned  []registryExecutor `yaml:"executors_planned"`
	Rationale         string             `yaml:"rationale"`
	SunsetDecision    *registrySunset    `yaml:"sunset_decision"`
	ChildPlan         string             `yaml:"child_plan"`
	TargetPromotion   string             `yaml:"target_promotion"`
	ReplacedBy        string             `yaml:"replaced_by"`
	RemovalDate       string             `yaml:"removal_date"`
	ObjectiveImpossBc string             `yaml:"objective_validation_impossible_because"`
	ResearchQuestions []string           `yaml:"research_questions"`
	UnblockTimeline   string             `yaml:"estimated_unblock_timeline"`
}

type registryExecutor struct {
	File         string `yaml:"file"`
	Symbol       string `yaml:"symbol"`
	ExecutorKind string `yaml:"executor_kind"`
	HookPhase    string `yaml:"hook_phase"`
	BlockOrWarn  string `yaml:"block_or_warn"`
}

type registrySunset struct {
	RevisitWhen     string `yaml:"revisit_when"`
	SuccessCriteria string `yaml:"success_criteria"`
	RevisitOwner    string `yaml:"revisit_owner"`
}

// LintEnforcementRegistry runs the Phase 3 compile-time lint against the
// registry at <repo>/enforcement/enforcement-registry.yaml. Returns the
// list of findings (empty slice on PASS) plus a load error if the registry
// itself is missing or malformed.
func LintEnforcementRegistry(repo string) ([]EnforcementRegistryLintError, error) {
	reg, regPath, err := loadRegistrySnapshot(repo)
	if err != nil {
		return nil, err
	}
	var errs []EnforcementRegistryLintError
	errs = append(errs, lintOrphanRules(repo, reg)...)
	errs = append(errs, lintMissingExecutorSymbols(repo, reg)...)
	errs = append(errs, lintBehavioralIncompleteSunset(reg)...)
	errs = append(errs, lintDeprecatedDisposal(reg)...)
	errs = append(errs, lintOrphanExecutors(repo, reg)...)
	_ = regPath
	return errs, nil
}

func loadRegistrySnapshot(repo string) (*registrySnapshot, string, error) {
	p := filepath.Join(repo, "enforcement", "enforcement-registry.yaml")
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, p, fmt.Errorf("read enforcement-registry.yaml: %w", err)
	}
	var snap registrySnapshot
	if err := yaml.Unmarshal(data, &snap); err != nil {
		return nil, p, fmt.Errorf("parse enforcement-registry.yaml: %w", err)
	}
	return &snap, p, nil
}

// ─────────────────────────────────────────────────────────────────────
// orphan_rule: yaml under enforcement/runtime/governance with a top-level
// `id:` field whose path is not referenced by any rule_class.source_files.
// ─────────────────────────────────────────────────────────────────────

var topLevelIDPattern = regexp.MustCompile(`(?m)^id:\s*([^\s#]+)`)

func lintOrphanRules(repo string, reg *registrySnapshot) []EnforcementRegistryLintError {
	bound := map[string]bool{}
	for _, rc := range reg.RuleClasses {
		for _, sf := range rc.SourceFiles {
			bound[normalizeSourcePath(sf)] = true
		}
	}
	var errs []EnforcementRegistryLintError
	roots := []string{"enforcement", "runtime", "governance"}
	for _, root := range roots {
		base := filepath.Join(repo, root)
		_ = filepath.Walk(base, func(p string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			if info == nil || info.IsDir() {
				return nil
			}
			name := info.Name()
			if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
				return nil
			}
			rel, err := filepath.Rel(repo, p)
			if err != nil {
				return nil
			}
			rel = filepath.ToSlash(rel)
			// Skip the registry itself and projection mirrors.
			if rel == "enforcement/enforcement-registry.yaml" {
				return nil
			}
			content, err := os.ReadFile(p)
			if err != nil {
				return nil
			}
			id := extractTopLevelID(string(content))
			if id == "" {
				return nil
			}
			if bound[rel] {
				return nil
			}
			errs = append(errs, EnforcementRegistryLintError{
				Type: "orphan_rule",
				Fields: []EnforcementRegistryLintField{
					{"file", rel},
					{"declared_id", id},
				},
				Message: "rule yaml declares id but no enforcement-registry rule_class binds it. Add a rule_class entry to enforcement/enforcement-registry.yaml or mark coverage=deprecated with replaced_by or removal_date.",
			})
			return nil
		})
	}
	sortLintErrors(errs)
	return errs
}

func extractTopLevelID(content string) string {
	// Strip yaml front-matter shebang/leading comments, just match first id: line.
	m := topLevelIDPattern.FindStringSubmatch(content)
	if len(m) < 2 {
		return ""
	}
	return strings.Trim(m[1], `"'`)
}

func normalizeSourcePath(p string) string {
	// Strip anchor (`#section`) and trim leading slashes.
	if idx := strings.Index(p, "#"); idx >= 0 {
		p = p[:idx]
	}
	return strings.TrimLeft(filepath.ToSlash(p), "/")
}

// ─────────────────────────────────────────────────────────────────────
// missing_executor_symbol: coverage=mechanical executors whose symbol is
// not present in the declared file (Go grep). Symbols in
// internal_helper_allowlist are exempt. runtime_state_machine_phase kind
// is skipped (symbol resolves against runtime.db phases, not Go AST).
// pending_implementation uses executors_planned and is NOT lint-checked.
// ─────────────────────────────────────────────────────────────────────

func lintMissingExecutorSymbols(repo string, reg *registrySnapshot) []EnforcementRegistryLintError {
	allow := map[string]bool{}
	for _, s := range reg.InternalHelperAllowlist.Symbols {
		allow[s] = true
	}
	requiredKinds := map[string]bool{}
	for _, k := range reg.ExecutorKindSpec.BindingRequiredFor {
		requiredKinds[k] = true
	}
	// File cache: symbols-defined-in-file
	fileSymbols := map[string]map[string]bool{}
	getSyms := func(rel string) map[string]bool {
		if syms, ok := fileSymbols[rel]; ok {
			return syms
		}
		full := filepath.Join(repo, filepath.FromSlash(rel))
		data, err := os.ReadFile(full)
		syms := map[string]bool{}
		if err == nil {
			for _, m := range goFuncDeclPattern.FindAllStringSubmatch(string(data), -1) {
				if len(m) >= 2 {
					syms[m[1]] = true
				}
			}
		}
		fileSymbols[rel] = syms
		return syms
	}
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "mechanical" {
			continue
		}
		for _, ex := range rc.Executors {
			if ex.ExecutorKind == "runtime_state_machine_phase" {
				continue
			}
			if ex.ExecutorKind != "" && !requiredKinds[ex.ExecutorKind] {
				continue
			}
			if ex.Symbol == "" {
				errs = append(errs, EnforcementRegistryLintError{
					Type: "missing_executor_symbol",
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"expected_symbol", "(empty)"},
						{"file", ex.File},
						{"executor_kind", ex.ExecutorKind},
					},
					Message: "rule_class declares mechanical coverage but executor.symbol is empty. Add the symbol or change coverage to pending_implementation.",
				})
				continue
			}
			if allow[ex.Symbol] {
				continue
			}
			if !strings.HasSuffix(ex.File, ".go") {
				// Non-Go executors (e.g. runtime.db references) — skip symbol check.
				continue
			}
			syms := getSyms(ex.File)
			if syms[ex.Symbol] {
				continue
			}
			errs = append(errs, EnforcementRegistryLintError{
				Type: "missing_executor_symbol",
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"expected_symbol", ex.Symbol},
					{"file", ex.File},
					{"executor_kind", ex.ExecutorKind},
				},
				Message: "rule_class declares mechanical coverage but executor symbol not found. Either implement the symbol, change coverage to pending_implementation with child_plan, or move symbol to internal_helper_allowlist if it is genuinely helper-only.",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}

var goFuncDeclPattern = regexp.MustCompile(`(?m)^func(?:\s*\([^)]*\))?\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`)

// ─────────────────────────────────────────────────────────────────────
// behavioral_only_incomplete_sunset: Q2 dual-required lock on
// sunset_decision.revisit_when + sunset_decision.success_criteria.
// ─────────────────────────────────────────────────────────────────────

func lintBehavioralIncompleteSunset(reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "behavioral_only" {
			continue
		}
		if rc.SunsetDecision == nil {
			errs = append(errs, EnforcementRegistryLintError{
				Type: "behavioral_only_incomplete_sunset",
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"missing_field", "sunset_decision"},
				},
				Message: "behavioral_only requires both revisit_when (event trigger) AND success_criteria (objective test). Single-field schema lets sunset decay into 'have criteria but no one checks' — more dangerous than no criteria with checks. Add the missing field or change coverage to not_mechanizable + impossibility rationale.",
			})
			continue
		}
		if strings.TrimSpace(rc.SunsetDecision.RevisitWhen) == "" {
			errs = append(errs, EnforcementRegistryLintError{
				Type: "behavioral_only_incomplete_sunset",
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"missing_field", "sunset_decision.revisit_when"},
				},
				Message: "behavioral_only requires both revisit_when (event trigger) AND success_criteria (objective test). Add revisit_when or change coverage to not_mechanizable.",
			})
		}
		if strings.TrimSpace(rc.SunsetDecision.SuccessCriteria) == "" {
			errs = append(errs, EnforcementRegistryLintError{
				Type: "behavioral_only_incomplete_sunset",
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"missing_field", "sunset_decision.success_criteria"},
				},
				Message: "behavioral_only requires both revisit_when (event trigger) AND success_criteria (objective test). Add success_criteria or change coverage to not_mechanizable.",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// deprecated_missing_disposal / deprecated_past_removal_date:
// require replaced_by OR removal_date; removal_date must be future ISO-8601.
// ─────────────────────────────────────────────────────────────────────

func lintDeprecatedDisposal(reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	today := time.Now().UTC().Format("2006-01-02")
	activeIDs := map[string]bool{}
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "deprecated" {
			activeIDs[rc.ID] = true
		}
	}
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "deprecated" {
			continue
		}
		hasReplaced := strings.TrimSpace(rc.ReplacedBy) != ""
		hasRemoval := strings.TrimSpace(rc.RemovalDate) != ""
		if !hasReplaced && !hasRemoval {
			errs = append(errs, EnforcementRegistryLintError{
				Type: "deprecated_missing_disposal",
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
				},
				Message: "deprecated rule_class requires either replaced_by (pointing to an active mechanical class) or removal_date (future ISO-8601). Add one or remove the entry.",
			})
			continue
		}
		if hasReplaced && !activeIDs[rc.ReplacedBy] {
			errs = append(errs, EnforcementRegistryLintError{
				Type: "deprecated_missing_disposal",
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"replaced_by", rc.ReplacedBy},
				},
				Message: "deprecated.replaced_by must resolve to an active rule_class id.",
			})
		}
		if hasRemoval {
			if _, err := time.Parse("2006-01-02", rc.RemovalDate); err != nil {
				errs = append(errs, EnforcementRegistryLintError{
					Type: "deprecated_missing_disposal",
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"removal_date", rc.RemovalDate},
					},
					Message: "deprecated.removal_date must be ISO-8601 YYYY-MM-DD.",
				})
				continue
			}
			if rc.RemovalDate < today {
				errs = append(errs, EnforcementRegistryLintError{
					Type: "deprecated_past_removal_date",
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"removal_date", rc.RemovalDate},
						{"today", today},
					},
					Message: "deprecated rule_class is past removal_date but still present. Actually remove the entry, or extend removal_date with rationale.",
				})
			}
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// orphan_executor: hooks.go top-level functions matching dispatcher /
// validator naming that are neither bound by any rule_class.executors[]
// nor in internal_helper_allowlist.
// ─────────────────────────────────────────────────────────────────────

func lintOrphanExecutors(repo string, reg *registrySnapshot) []EnforcementRegistryLintError {
	hooksRel := "scripts/ai-skill-cli/internal/app/hooks.go"
	full := filepath.Join(repo, filepath.FromSlash(hooksRel))
	data, err := os.ReadFile(full)
	if err != nil {
		return nil
	}
	bound := map[string]bool{}
	for _, rc := range reg.RuleClasses {
		for _, ex := range rc.Executors {
			if ex.Symbol != "" {
				bound[ex.Symbol] = true
			}
		}
		for _, ex := range rc.ExecutorsPlanned {
			if ex.Symbol != "" {
				bound[ex.Symbol] = true
			}
		}
	}
	allow := map[string]bool{}
	for _, s := range reg.InternalHelperAllowlist.Symbols {
		allow[s] = true
	}
	var errs []EnforcementRegistryLintError
	for _, m := range goFuncDeclPattern.FindAllStringSubmatch(string(data), -1) {
		if len(m) < 2 {
			continue
		}
		sym := m[1]
		// Only flag names that look like dispatchers or validators —
		// helpers (`appendLog`, `parseGitHubRemote`) without `run`/`validate`
		// prefix don't even reach the candidate set, mirroring Q4 intent
		// that internal_helper is the default and exemption is the rule.
		if !isExecutorCandidate(sym) {
			continue
		}
		if bound[sym] || allow[sym] {
			continue
		}
		errs = append(errs, EnforcementRegistryLintError{
			Type: "orphan_executor",
			Fields: []EnforcementRegistryLintField{
				{"file", hooksRel},
				{"symbol", sym},
			},
			Message: "executor-looking symbol (run* / validate*) in hooks.go is not bound by any rule_class.executors[] and not listed in internal_helper_allowlist. Add a binding, or move it to internal_helper_allowlist if it is genuinely helper-only.",
		})
	}
	sortLintErrors(errs)
	return errs
}

func isExecutorCandidate(sym string) bool {
	return strings.HasPrefix(sym, "run") || strings.HasPrefix(sym, "validate")
}

func sortLintErrors(errs []EnforcementRegistryLintError) {
	sort.SliceStable(errs, func(i, j int) bool {
		if errs[i].Type != errs[j].Type {
			return errs[i].Type < errs[j].Type
		}
		return errs[i].Format() < errs[j].Format()
	})
}
