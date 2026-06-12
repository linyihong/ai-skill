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

// SeverityFail is a P0 blocking finding; compile must fail.
// SeverityWarn is a non-blocking finding surfaced to the maintainer.
const (
	SeverityFail = "FAIL"
	SeverityWarn = "WARNING"
)

// EnforcementRegistryLintError is one finding from Phase 3 lint.
//
// Format follows Phase 7 scenario contract:
//
//	LINT [FAIL|WARNING] [<type>]
//	  <key>: <value>
//	  ...
//	  message: <one-line explanation + remediation>
type EnforcementRegistryLintError struct {
	Type     string
	Severity string // SeverityFail | SeverityWarn ; empty defaults to FAIL for backward compat
	Fields   []EnforcementRegistryLintField
	Message  string
}

type EnforcementRegistryLintField struct {
	Key   string
	Value string
}

// IsFail returns true if this finding should block compile. Empty Severity
// is treated as FAIL for backward compat with pre-v2 lint output.
func (e EnforcementRegistryLintError) IsFail() bool {
	return e.Severity == "" || e.Severity == SeverityFail
}

func (e EnforcementRegistryLintError) Format() string {
	sev := e.Severity
	if sev == "" {
		sev = SeverityFail
	}
	var b strings.Builder
	fmt.Fprintf(&b, "LINT %s [%s]\n", sev, e.Type)
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
	// Schema patch v2 A2: bootstrap_mode is single source of truth.
	BootstrapMode    string                    `yaml:"bootstrap_mode"`
	BaselineSnapshot registryBaselineSnapshot  `yaml:"baseline_snapshot"`
	GovernanceThresholds registryGovernanceThresholds `yaml:"governance_thresholds"`
	ExecutorKindSpec struct {
		BindingRequiredFor []string `yaml:"binding_required_for"`
	} `yaml:"executor_kind_spec"`
	InternalHelperAllowlist struct {
		Symbols []string `yaml:"symbols"`
	} `yaml:"internal_helper_allowlist"`
	RuleClasses []registryRuleClass `yaml:"rule_classes"`
}

type registryGovernanceThresholds struct {
	SourceFilesReviewThreshold int `yaml:"source_files_review_threshold"`
}

type registryBaselineSnapshot struct {
	BaselineCreatedAt          string                          `yaml:"baseline_created_at"`
	BaselineBurndownTargetDate string                          `yaml:"baseline_burndown_target_date"`
	BaselineOwner              string                          `yaml:"baseline_owner"`
	Entries                    []registryBaselineSnapshotEntry `yaml:"entries"`
}

type registryBaselineSnapshotEntry struct {
	FindingType           string `yaml:"finding_type"`
	Identifier            string `yaml:"identifier"`
	BaselineReviewSummary string `yaml:"baseline_review_summary"`
	AcceptedAt            string `yaml:"accepted_at"`
	BurndownOwner         string `yaml:"burndown_owner"`
}

type registryRuleClass struct {
	ID                          string             `yaml:"id"`
	Coverage                    string             `yaml:"coverage"`
	SourceFiles                 []string           `yaml:"source_files"`
	Executors                   []registryExecutor `yaml:"executors"`
	ExecutorsPlanned            []registryExecutor `yaml:"executors_planned"`
	Rationale                   string             `yaml:"rationale"`
	SunsetDecision              *registrySunset    `yaml:"sunset_decision"`
	ChildPlan                   string             `yaml:"child_plan"`
	TargetPromotion             string             `yaml:"target_promotion"`
	ReplacedBy                  string             `yaml:"replaced_by"`
	RemovalDate                 string             `yaml:"removal_date"`
	ObjectiveImpossBc           string             `yaml:"objective_validation_impossible_because"`
	ResearchQuestions           []string           `yaml:"research_questions"`
	UnblockTimeline             string             `yaml:"estimated_unblock_timeline"`
	// Schema patch v2 additions:
	UpstreamClasses             []string           `yaml:"upstream_classes"`
	SizeReviewExemptionRationale string            `yaml:"size_review_exemption_rationale"`
	// Phase 4.5 — registry self-governance. Required on rule_class
	// entries whose coverage value is being demoted (e.g. mechanical →
	// behavioral_only). Validator validateEnforcementRegistryTransition
	// blocks demotion commits that lack this field. Format:
	// `constitution/ADR-<NNN>-<slug>.md` (file must exist under <repo>).
	AdrReference string `yaml:"adr_reference"`
	DemotionRationale string `yaml:"demotion_rationale"`
	// F19 (validation_scenario_governance): coverage_evidence binds a
	// rule_class to the validation scenarios that prove its mechanical
	// coverage. Parsed here (shared struct) and consumed by scenario_lint.go
	// (LintValidationScenarios).
	CoverageEvidence *registryCoverageEvidence `yaml:"coverage_evidence"`
}

// registryCoverageEvidence is the coverage_evidence block on a rule_class.
// CoverageTargetPct is a pointer so the lint can distinguish "absent" (nil)
// from an explicit 0.
type registryCoverageEvidence struct {
	ExpectedInstanceCount int      `yaml:"expected_instance_count"`
	ValidationScenarios   []string `yaml:"validation_scenarios"`
	RegressionScenarios   []string `yaml:"regression_scenarios"`
	CoverageTargetPct     *int     `yaml:"coverage_target_pct"`
}

type registryExecutor struct {
	File         string `yaml:"file"`
	Symbol       string `yaml:"symbol"`
	ExecutorKind string `yaml:"executor_kind"`
	HookPhase    string `yaml:"hook_phase"`
	BlockOrWarn  string `yaml:"block_or_warn"`
}

type registrySunset struct {
	RevisitWhen       string `yaml:"revisit_when"`
	SuccessCriteria   string `yaml:"success_criteria"`
	RevisitOwner      string `yaml:"revisit_owner"`
	LastReviewedAt    string `yaml:"last_reviewed_at"`
	LastReviewSummary string `yaml:"last_review_summary"`
	// Pointer so the lint can distinguish "field absent" (nil → warning)
	// from "explicitly declared empty" ([] → conscious no-dependency, OK).
	DependsOnRuleClasses *[]string `yaml:"depends_on_rule_classes"`
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
	// Existing v1 lints (orphan / missing executor / deprecated):
	errs = append(errs, lintOrphanRules(repo, reg)...)
	errs = append(errs, lintMissingExecutorSymbols(repo, reg)...)
	errs = append(errs, lintBehavioralIncompleteSunset(reg)...)
	errs = append(errs, lintDeprecatedDisposal(reg)...)
	errs = append(errs, lintOrphanExecutors(repo, reg)...)
	// Schema patch v2 additions — behavioral_only family:
	errs = append(errs, lintBehavioralRecommendedFields(reg)...)
	errs = append(errs, lintBehavioralReviewAge(reg)...)
	errs = append(errs, lintBehavioralVagueSuccessCriteria(reg)...)
	errs = append(errs, lintBehavioralMissingMeasurableSignal(reg)...)
	errs = append(errs, lintBehavioralRevisitChain(reg)...)
	// Schema patch v2 — compile_time_lint_rules R1-R4:
	errs = append(errs, lintUpstreamChainResolution(reg)...)
	errs = append(errs, lintClassSizeReviewThreshold(reg)...)
	errs = append(errs, lintBaselineSnapshotGovernance(reg)...)
	errs = append(errs, lintPendingImplementationChildPlanValidity(repo, reg)...)
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
		// Round-4 T2: rationale is the 3rd hard-required field for behavioral_only.
		if strings.TrimSpace(rc.Rationale) == "" {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "behavioral_only_missing_rationale",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"missing_field", "rationale"},
				},
				Message: "behavioral_only requires rationale (one of 3 hard required: rationale + sunset_decision.revisit_when + sunset_decision.success_criteria). Add a one-paragraph rationale explaining why this class is intentionally not mechanized.",
			})
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

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — behavioral_only recommended fields (WARNING)
// Round-4 T2 + Round-3 + Round-5 U1: 4 recommended fields downgraded
// from required to recommended; lint emits WARNING (not FAIL).
// ─────────────────────────────────────────────────────────────────────

func lintBehavioralRecommendedFields(reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "behavioral_only" || rc.SunsetDecision == nil {
			continue
		}
		s := rc.SunsetDecision
		check := func(field, value string) {
			if strings.TrimSpace(value) == "" {
				errs = append(errs, EnforcementRegistryLintError{
					Type:     "behavioral_only_missing_" + field,
					Severity: SeverityWarn,
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"missing_field", "sunset_decision." + field},
					},
					Message: fmt.Sprintf("behavioral_only recommends sunset_decision.%s (governance signal, not blocker). Add the field to surface in coverage report.", field),
				})
			}
		}
		check("revisit_owner", s.RevisitOwner)
		check("last_reviewed_at", s.LastReviewedAt)
		check("last_review_summary", s.LastReviewSummary)
		if s.DependsOnRuleClasses == nil {
			// nil = field absent. Explicitly empty [] is a conscious
			// declaration of "no rule_class dependency" and is accepted.
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "behavioral_only_missing_depends_on_rule_classes",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"missing_field", "sunset_decision.depends_on_rule_classes"},
				},
				Message: "behavioral_only recommends sunset_decision.depends_on_rule_classes (structured chain; replaces NLP parse of revisit_when). Declare [] for no dependency, or list class ids if sunset is triggered by another rule_class state change.",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — behavioral_only_review_age (FAIL >24m)
// Round-5 U1 mutual exclusion: only fires when last_reviewed_at present.
// ─────────────────────────────────────────────────────────────────────

func lintBehavioralReviewAge(reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	now := time.Now().UTC()
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "behavioral_only" || rc.SunsetDecision == nil {
			continue
		}
		raw := strings.TrimSpace(rc.SunsetDecision.LastReviewedAt)
		if raw == "" {
			continue // U1: skip — missing_last_reviewed_at warning handles this
		}
		ts, err := time.Parse("2006-01-02", raw)
		if err != nil {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "behavioral_only_review_age",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"last_reviewed_at", raw},
				},
				Message: "last_reviewed_at must be ISO-8601 YYYY-MM-DD.",
			})
			continue
		}
		months := int(now.Sub(ts).Hours() / 24 / 30)
		if months > 24 {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "behavioral_only_review_age",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"last_reviewed_at", raw},
					{"months_since_review", fmt.Sprintf("%d", months)},
				},
				Message: "behavioral_only review age > 24 months. Re-review the class (verify revisit_when / success_criteria still apply, update last_reviewed_at + last_review_summary), promote to mechanical, or demote to not_mechanizable.",
			})
		} else if months > 12 {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "behavioral_only_review_age",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"last_reviewed_at", raw},
					{"months_since_review", fmt.Sprintf("%d", months)},
				},
				Message: "behavioral_only review age > 12 months; consider re-reviewing before age exceeds 24 months (FAIL threshold).",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — behavioral_only_vague_success_criteria (FAIL blacklist)
// ─────────────────────────────────────────────────────────────────────

var vagueBlacklist = []string{"TBD", "未來", "future", "eventually", "TODO"}

func lintBehavioralVagueSuccessCriteria(reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "behavioral_only" || rc.SunsetDecision == nil {
			continue
		}
		sc := rc.SunsetDecision.SuccessCriteria
		scLower := strings.ToLower(sc)
		for _, bad := range vagueBlacklist {
			if strings.Contains(scLower, strings.ToLower(bad)) {
				errs = append(errs, EnforcementRegistryLintError{
					Type:     "behavioral_only_vague_success_criteria",
					Severity: SeverityFail,
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"blacklist_match", bad},
					},
					Message: "success_criteria contains vague token; rewrite as concrete observable condition (event / count / state transition).",
				})
				break
			}
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — behavioral_only_missing_measurable_signal (WARNING whitelist)
// C9: TOKEN-LEVEL heuristic only.
// ─────────────────────────────────────────────────────────────────────

var measurableSignalPattern = regexp.MustCompile(`\d+|%|rule_class|executor|lint|coverage|validator|hook|scenario|gate`)

func lintBehavioralMissingMeasurableSignal(reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "behavioral_only" || rc.SunsetDecision == nil {
			continue
		}
		sc := rc.SunsetDecision.SuccessCriteria
		if strings.TrimSpace(sc) == "" {
			continue // FAIL already covered by behavioral_only_incomplete_sunset
		}
		if measurableSignalPattern.MatchString(sc) {
			continue
		}
		errs = append(errs, EnforcementRegistryLintError{
			Type:     "behavioral_only_missing_measurable_signal",
			Severity: SeverityWarn,
			Fields: []EnforcementRegistryLintField{
				{"rule_class", rc.ID},
			},
			Message: "success_criteria contains no measurable token (no digit / % / framework noun like rule_class/executor/lint/coverage/validator/hook/scenario/gate). Token-level heuristic — pass does NOT imply measurable, fail does NOT imply unmeasurable. Treat as advisory.",
		})
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — behavioral_only_revisit_chain (FAIL)
// Round-3 S1: lint reads ONLY depends_on_rule_classes (structured).
// Round-4 prevention: behavioral_only depending on behavioral_only forms
// decay chain (no one ever reviews).
// ─────────────────────────────────────────────────────────────────────

func lintBehavioralRevisitChain(reg *registrySnapshot) []EnforcementRegistryLintError {
	classByID := map[string]*registryRuleClass{}
	for i := range reg.RuleClasses {
		classByID[reg.RuleClasses[i].ID] = &reg.RuleClasses[i]
	}
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "behavioral_only" || rc.SunsetDecision == nil {
			continue
		}
		if rc.SunsetDecision.DependsOnRuleClasses == nil {
			continue
		}
		for _, depID := range *rc.SunsetDecision.DependsOnRuleClasses {
			dep, ok := classByID[depID]
			if !ok {
				errs = append(errs, EnforcementRegistryLintError{
					Type:     "behavioral_only_revisit_chain",
					Severity: SeverityFail,
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"depends_on", depID},
					},
					Message: "sunset_decision.depends_on_rule_classes references unknown rule_class. Either fix the id or remove the reference.",
				})
				continue
			}
			if dep.Coverage == "behavioral_only" {
				errs = append(errs, EnforcementRegistryLintError{
					Type:     "behavioral_only_revisit_chain",
					Severity: SeverityFail,
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"depends_on", depID},
						{"chain_type", "behavioral_only → behavioral_only"},
					},
					Message: "behavioral_only depending on another behavioral_only forms decay chain (no one ever reviews). Either depend on a mechanical/pending class, or merge the two classes.",
				})
			}
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — compile_time_lint_rules.R1: upstream_chain_resolution
// reference existence + cycle detection.
// ─────────────────────────────────────────────────────────────────────

func lintUpstreamChainResolution(reg *registrySnapshot) []EnforcementRegistryLintError {
	classIDs := map[string]bool{}
	upstream := map[string][]string{}
	for _, rc := range reg.RuleClasses {
		classIDs[rc.ID] = true
		upstream[rc.ID] = append([]string(nil), rc.UpstreamClasses...)
	}
	var errs []EnforcementRegistryLintError
	// 1. Unresolved references
	for _, rc := range reg.RuleClasses {
		for _, up := range rc.UpstreamClasses {
			if !classIDs[up] {
				errs = append(errs, EnforcementRegistryLintError{
					Type:     "upstream_chain_resolution",
					Severity: SeverityFail,
					Fields: []EnforcementRegistryLintField{
						{"rule_class", rc.ID},
						{"upstream_class", up},
					},
					Message: "upstream_classes references unknown rule_class id. Fix the id or remove the reference.",
				})
			}
		}
	}
	// 2. Cycle detection (DFS)
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := map[string]int{}
	var dfs func(id string, path []string) bool
	dfs = func(id string, path []string) bool {
		color[id] = gray
		path = append(path, id)
		for _, up := range upstream[id] {
			if !classIDs[up] {
				continue
			}
			if color[up] == gray {
				errs = append(errs, EnforcementRegistryLintError{
					Type:     "upstream_chain_resolution",
					Severity: SeverityFail,
					Fields: []EnforcementRegistryLintField{
						{"cycle", strings.Join(append(path, up), " → ")},
					},
					Message: "upstream_classes forms a cycle. Promotion chains must be acyclic. Restructure the dependency to break the cycle.",
				})
				return true
			}
			if color[up] == white {
				if dfs(up, path) {
					return true
				}
			}
		}
		color[id] = black
		return false
	}
	for id := range classIDs {
		if color[id] == white {
			dfs(id, nil)
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — compile_time_lint_rules.R2: class_size_review_threshold
// WARNING when source_files > threshold; suppressible via
// size_review_exemption_rationale (warning still emitted but includes rationale).
// ─────────────────────────────────────────────────────────────────────

func lintClassSizeReviewThreshold(reg *registrySnapshot) []EnforcementRegistryLintError {
	threshold := reg.GovernanceThresholds.SourceFilesReviewThreshold
	if threshold == 0 {
		return nil // not configured
	}
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		count := len(rc.SourceFiles)
		if count <= threshold {
			continue
		}
		fields := []EnforcementRegistryLintField{
			{"rule_class", rc.ID},
			{"source_files_count", fmt.Sprintf("%d", count)},
			{"threshold", fmt.Sprintf("%d", threshold)},
		}
		msg := "rule_class.source_files exceeds review threshold. Consider whether the class should be split, or add size_review_exemption_rationale documenting why this class is legitimately large."
		if rationale := strings.TrimSpace(rc.SizeReviewExemptionRationale); rationale != "" {
			fields = append(fields, EnforcementRegistryLintField{"exemption_rationale", rationale})
			msg = "rule_class.source_files exceeds review threshold; maintainer rationale recorded above (warning preserved by design, not suppressed)."
		}
		errs = append(errs, EnforcementRegistryLintError{
			Type:     "class_size_review_threshold",
			Severity: SeverityWarn,
			Fields:   fields,
			Message:  msg,
		})
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — compile_time_lint_rules.R3: baseline_snapshot_missing_governance
// FAIL when bootstrap_mode = baseline_snapshot_v1 AND (owner missing OR
// any entry has empty review_summary).
// ─────────────────────────────────────────────────────────────────────

func lintBaselineSnapshotGovernance(reg *registrySnapshot) []EnforcementRegistryLintError {
	if reg.BootstrapMode != "baseline_snapshot_v1" {
		return nil
	}
	var errs []EnforcementRegistryLintError
	if strings.TrimSpace(reg.BaselineSnapshot.BaselineOwner) == "" {
		errs = append(errs, EnforcementRegistryLintError{
			Type:     "baseline_snapshot_missing_governance",
			Severity: SeverityFail,
			Fields: []EnforcementRegistryLintField{
				{"missing_field", "baseline_snapshot.baseline_owner"},
			},
			Message: "bootstrap_mode=baseline_snapshot_v1 requires baseline_snapshot.baseline_owner (responsible party for burndown). Add the field.",
		})
	}
	for i, entry := range reg.BaselineSnapshot.Entries {
		s := strings.TrimSpace(entry.BaselineReviewSummary)
		if s == "" {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "baseline_snapshot_missing_governance",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"entry_index", fmt.Sprintf("%d", i)},
					{"identifier", entry.Identifier},
					{"missing_field", "baseline_review_summary"},
				},
				Message: "each baseline_snapshot.entries[] requires baseline_review_summary (>= 20 chars; what maintainer found when accepting into baseline).",
			})
			continue
		}
		if len(s) < 20 {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "baseline_snapshot_missing_governance",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"entry_index", fmt.Sprintf("%d", i)},
					{"identifier", entry.Identifier},
					{"summary_length", fmt.Sprintf("%d", len(s))},
				},
				Message: "baseline_review_summary too short (< 20 chars). Expand to describe what was reviewed when accepting this finding into baseline.",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 — compile_time_lint_rules.R4: pending_implementation_child_plan_validity
// (a) path resolves → FAIL ; (b)(c)(d) → WARNING.
// B5: path normalized via path.split('#')[0] before resolution.
// ─────────────────────────────────────────────────────────────────────

var (
	// Phase 0 may live at any heading depth (e.g. `### Phase 0` nested
	// under a `## Phase Plan` section). Rule (b) only asks that a Phase 0
	// outline exists, not that it is an h2.
	phase0Pattern     = regexp.MustCompile(`(?m)^#{2,}\s+Phase\s+0\b`)
	ownerPattern      = regexp.MustCompile(`(?mi)^(?:owner\s*:|.*\bOwner\s*:)\s*\S`)
	// Validation Plan / Acceptance also accepted at any heading depth.
	acceptancePattern = regexp.MustCompile(`(?m)^#{2,}\s+(Validation Plan|Acceptance)\b`)
)

func lintPendingImplementationChildPlanValidity(repo string, reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "pending_implementation" {
			continue
		}
		raw := strings.TrimSpace(rc.ChildPlan)
		if raw == "" {
			// pending_implementation.requires = [child_plan, target_promotion];
			// missing child_plan caught by future required-field lint. Here we
			// only validate when present.
			continue
		}
		// B5 anchor strip
		pathOnly := raw
		if idx := strings.Index(pathOnly, "#"); idx >= 0 {
			pathOnly = pathOnly[:idx]
		}
		full := filepath.Join(repo, filepath.FromSlash(pathOnly))
		content, err := os.ReadFile(full)
		if err != nil {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "pending_implementation_child_plan_validity",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"child_plan", raw},
					{"resolved_path", pathOnly},
					{"violation", "(a) path_resolves"},
				},
				Message: "child_plan path does not resolve to existing plans/active/*.md file. Fix the path or remove the rule_class.",
			})
			continue
		}
		body := string(content)
		if !phase0Pattern.MatchString(body) {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "pending_implementation_child_plan_validity",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"child_plan", raw},
					{"violation", "(b) phase_0_heading"},
				},
				Message: "child_plan body missing `## Phase 0` heading. Stub plans should at minimum outline Phase 0 preflight.",
			})
		}
		if !ownerPattern.MatchString(body) {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "pending_implementation_child_plan_validity",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"child_plan", raw},
					{"violation", "(c) owner_present"},
				},
				Message: "child_plan missing owner declaration (frontmatter `owner:` or body `Owner:`).",
			})
		}
		if !acceptancePattern.MatchString(body) {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "pending_implementation_child_plan_validity",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"child_plan", raw},
					{"violation", "(d) acceptance_section"},
				},
				Message: "child_plan missing `## Validation Plan` or `## Acceptance` section. Stub plans should declare completion criteria up front.",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}
