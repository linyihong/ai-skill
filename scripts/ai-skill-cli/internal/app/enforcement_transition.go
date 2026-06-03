package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Phase 4.5 of plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md
// — Registry Self-Governance.
//
// This file implements the engine that powers both:
//   1. `ai-skill enforcement transition-check` CLI subcommand (public surface,
//      used by validation scenarios for cross-platform fixture-based testing)
//   2. `validateEnforcementRegistryTransition` commit-msg hook validator (Go
//      registry entry obligation.commit.enforcement_registry_transition;
//      blocks commits that violate R1/R2/R3 self-governance lint rules)
//
// The two surfaces share the same `checkRegistryTransitions` core so a
// scenario passing the CLI is genuine evidence the commit-msg path
// behaves identically.

// ─────────────────────────────────────────────────────────────────────
// Violation model
// ─────────────────────────────────────────────────────────────────────

type registryTransitionViolation struct {
	Code      string // R1_*, R2_*, R3_*
	RuleClass string // empty for R1 commit-msg-level violations
	From      string // empty for R1
	To        string // empty for R1
	Detail    string // human readable
}

func (v registryTransitionViolation) String() string {
	parts := []string{"[" + v.Code + "]"}
	if v.RuleClass != "" {
		parts = append(parts, "rule_class="+v.RuleClass)
	}
	if v.From != "" || v.To != "" {
		parts = append(parts, fmt.Sprintf("transition=%s→%s", v.From, v.To))
	}
	parts = append(parts, v.Detail)
	return strings.Join(parts, " ")
}

// transitionInput bundles the three artifacts every R1/R2/R3 check needs.
type transitionInput struct {
	repoRoot  string // for ADR file resolution + executor symbol grep
	oldSnap   *registrySnapshot
	newSnap   *registrySnapshot
	commitMsg string
}

// Demotion / promotion tables (Phase 4.5 contract — see
// registry-transition-{demotion-without-adr,promotion-verification-gap}-v1.yaml
// phase_4_5_validator_contract sections).
var demotionTable = map[string]map[string]bool{
	"mechanical": {
		"behavioral_only":        true,
		"not_mechanizable":       true,
		"research_required":      true,
		"pending_implementation": true,
		// mechanical → deprecated is end-of-life, not demotion (handled
		// by existing Phase 3 deprecated_disposal lint).
	},
	"pending_implementation": {
		"behavioral_only":   true,
		"research_required": true,
	},
	"behavioral_only": {
		"not_mechanizable": true,
	},
}

func isDemotion(from, to string) bool {
	if m, ok := demotionTable[from]; ok {
		return m[to]
	}
	return false
}

func isPromotionToMechanical(from, to string) bool {
	if to != "mechanical" {
		return false
	}
	switch from {
	case "(new)", "pending_implementation", "research_required", "behavioral_only":
		return true
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────
// Core engine
// ─────────────────────────────────────────────────────────────────────

const transitionOptOutMarker = "[skip-registry-transition]"

func checkRegistryTransitions(in transitionInput) []registryTransitionViolation {
	if strings.Contains(in.commitMsg, transitionOptOutMarker) {
		return nil
	}
	if in.oldSnap == nil || in.newSnap == nil {
		return nil
	}

	// Build id → coverage maps; missing in old means "(new)".
	oldByID := map[string]string{}
	for _, rc := range in.oldSnap.RuleClasses {
		oldByID[rc.ID] = rc.Coverage
	}
	newByID := map[string]*registryRuleClass{}
	for i := range in.newSnap.RuleClasses {
		rc := &in.newSnap.RuleClasses[i]
		newByID[rc.ID] = rc
	}

	// Detect transitions: any new class whose coverage differs from
	// its old counterpart (or which is new entirely).
	type transition struct {
		id   string
		from string
		to   string
		rc   *registryRuleClass
	}
	var transitions []transition
	for id, rc := range newByID {
		from, existed := oldByID[id]
		if !existed {
			from = "(new)"
		}
		if from == rc.Coverage {
			continue
		}
		transitions = append(transitions, transition{id: id, from: from, to: rc.Coverage, rc: rc})
	}

	var violations []registryTransitionViolation

	// R1 — trailer + rationale gate. Fires once if any transition is
	// present and the commit message is missing either marker. The R1
	// gate does NOT block R2/R3 from firing; all violations surface
	// together so the maintainer sees the full picture in one shot.
	if len(transitions) > 0 {
		if !containsTrailer(in.commitMsg, "[registry-status-change]") {
			violations = append(violations, registryTransitionViolation{
				Code:   "R1_missing_trailer",
				Detail: "commit body must include [registry-status-change] trailer when staged diff changes rule_class coverage",
			})
		}
		if !containsRationaleLine(in.commitMsg) {
			violations = append(violations, registryTransitionViolation{
				Code:   "R1_missing_rationale",
				Detail: "commit body must include a `rationale: <text>` line explaining the status change",
			})
		}
	}

	// R2 / R3 — per transition.
	for _, t := range transitions {
		if isDemotion(t.from, t.to) {
			violations = append(violations, checkR2DemotionADR(in.repoRoot, t.id, t.from, t.to, t.rc)...)
		}
		if isPromotionToMechanical(t.from, t.to) {
			violations = append(violations, checkR3PromotionExecutor(in.repoRoot, t.id, t.from, t.to, t.rc, in.newSnap)...)
		}
	}

	return violations
}

// containsTrailer checks for a stand-alone trailer token on its own line
// (mirrors the existing skip-marker convention in hooks.go).
func containsTrailer(commitMsg, trailer string) bool {
	scan := strings.Split(commitMsg, "\n")
	for _, line := range scan {
		if strings.TrimSpace(line) == trailer {
			return true
		}
	}
	return false
}

var rationaleLinePattern = regexp.MustCompile(`(?im)^\s*rationale\s*:\s*\S`)

func containsRationaleLine(commitMsg string) bool {
	return rationaleLinePattern.MatchString(commitMsg)
}

func checkR2DemotionADR(repo, classID, from, to string, rc *registryRuleClass) []registryTransitionViolation {
	adr := strings.TrimSpace(rc.AdrReference)
	if adr == "" {
		return []registryTransitionViolation{{
			Code: "R2_demotion_missing_adr", RuleClass: classID, From: from, To: to,
			Detail: "demotion requires adr_reference field on the rule_class pointing to constitution/ADR-*.md",
		}}
	}
	if !strings.HasPrefix(adr, "constitution/ADR-") || !strings.HasSuffix(adr, ".md") {
		return []registryTransitionViolation{{
			Code: "R2_demotion_invalid_adr_format", RuleClass: classID, From: from, To: to,
			Detail: fmt.Sprintf("adr_reference %q must match constitution/ADR-*.md", adr),
		}}
	}
	full := filepath.Join(repo, filepath.FromSlash(adr))
	if _, err := os.Stat(full); err != nil {
		return []registryTransitionViolation{{
			Code: "R2_demotion_adr_unresolved", RuleClass: classID, From: from, To: to,
			Detail: fmt.Sprintf("adr_reference %q does not resolve to an existing file under <repo>", adr),
		}}
	}
	return nil
}

func checkR3PromotionExecutor(repo, classID, from, to string, rc *registryRuleClass, newSnap *registrySnapshot) []registryTransitionViolation {
	// Re-use the Phase 3 missing_executor_symbol engine, but restricted
	// to this single class within the new registry. We synthesize a
	// snapshot containing just this class so the existing
	// lintMissingExecutorSymbols logic surfaces the gap.
	scoped := *newSnap // shallow copy
	scoped.RuleClasses = []registryRuleClass{*rc}
	errs := lintMissingExecutorSymbols(repo, &scoped)
	if len(errs) == 0 {
		return nil
	}
	var out []registryTransitionViolation
	for _, e := range errs {
		// Pull the expected_symbol + file out of the lint finding.
		var sym, file string
		for _, f := range e.Fields {
			switch f.Key {
			case "expected_symbol":
				sym = f.Value
			case "file":
				file = f.Value
			}
		}
		out = append(out, registryTransitionViolation{
			Code: "R3_promotion_missing_executor", RuleClass: classID, From: from, To: to,
			Detail: fmt.Sprintf("promotion to mechanical requires symbol_exists; symbol %q not found in %s", sym, file),
		})
	}
	return out
}

// ─────────────────────────────────────────────────────────────────────
// `ai-skill enforcement transition-check` CLI subcommand
// ─────────────────────────────────────────────────────────────────────

type enforcementTransitionOptions struct {
	repo            string
	old             string
	newPath         string
	commitMsgFile   string
	commitMsgInline string
	expectViolation string
	jsonOutput      bool
	plainOutput     bool
}

func runEnforcementTransitionCheck(args []string, stdout io.Writer, stderr io.Writer) int {
	opts := enforcementTransitionOptions{}
	fs := newFlagSet("enforcement transition-check", stderr)
	fs.StringVar(&opts.repo, "repo", ".", "repo root (used to resolve ADR + executor file paths)")
	fs.StringVar(&opts.old, "old", "", "path to old (HEAD) enforcement-registry.yaml")
	fs.StringVar(&opts.newPath, "new", "", "path to new (staged) enforcement-registry.yaml")
	fs.StringVar(&opts.commitMsgFile, "commit-msg-file", "", "path to commit message text file")
	fs.StringVar(&opts.commitMsgInline, "commit-msg", "", "inline commit message string (alternative to --commit-msg-file)")
	fs.StringVar(&opts.expectViolation, "expect-violation", "", "assertion mode: exit 0 if any violation code contains this substring, exit 30 otherwise")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write plain text output (default)")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}
	if strings.TrimSpace(opts.old) == "" || strings.TrimSpace(opts.newPath) == "" {
		_, _ = fmt.Fprintln(stderr, "--old and --new are required")
		return ExitInvalidUsage
	}

	root, err := resolveEnforcementRepo(opts.repo)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "resolve repo: %v\n", err)
		return ExitInvalidUsage
	}
	oldSnap, err := loadRegistrySnapshotFromPath(opts.old)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "load --old: %v\n", err)
		return ExitValidationFailed
	}
	newSnap, err := loadRegistrySnapshotFromPath(opts.newPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "load --new: %v\n", err)
		return ExitValidationFailed
	}
	msg := opts.commitMsgInline
	if strings.TrimSpace(opts.commitMsgFile) != "" {
		data, err := os.ReadFile(opts.commitMsgFile)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "read --commit-msg-file: %v\n", err)
			return ExitInvalidUsage
		}
		msg = string(data)
	}

	violations := checkRegistryTransitions(transitionInput{
		repoRoot:  root,
		oldSnap:   oldSnap,
		newSnap:   newSnap,
		commitMsg: msg,
	})

	result := buildTransitionCheckResult(opts, violations)
	if opts.jsonOutput {
		if err := writeJSON(stdout, result); err != nil {
			_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
			return ExitGeneralFailure
		}
		return result.ExitCode
	}
	if err := writePlain(stdout, result); err != nil {
		_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
		return ExitGeneralFailure
	}
	return result.ExitCode
}

func buildTransitionCheckResult(opts enforcementTransitionOptions, violations []registryTransitionViolation) Result {
	result := Result{
		Command:        "enforcement transition-check",
		Mode:           "report",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	result.Checks = append(result.Checks, Check{Name: "old_registry", Status: "ok", Message: opts.old})
	result.Checks = append(result.Checks, Check{Name: "new_registry", Status: "ok", Message: opts.newPath})
	result.Checks = append(result.Checks, Check{Name: "violations_total", Status: "ok", Message: fmt.Sprintf("%d", len(violations))})
	for _, v := range violations {
		result.Checks = append(result.Checks, Check{
			Name:    "violation." + v.Code,
			Status:  "failed",
			Message: v.String(),
		})
	}
	// Assertion mode.
	if strings.TrimSpace(opts.expectViolation) != "" {
		result.Mode = "assert"
		want := strings.ToLower(strings.TrimSpace(opts.expectViolation))
		matched := 0
		for _, v := range violations {
			if strings.Contains(strings.ToLower(v.Code), want) || strings.Contains(strings.ToLower(v.Detail), want) {
				matched++
			}
		}
		if matched > 0 {
			result.Checks = append(result.Checks, Check{Name: "assertion", Status: "ok", Message: fmt.Sprintf("matched %d violation(s)", matched)})
			return result
		}
		result.Status = "failed"
		result.ExitCode = ExitValidationFailed
		result.Checks = append(result.Checks, Check{Name: "assertion", Status: "failed", Message: "no violation matched --expect-violation substring"})
		result.Error = &CommandError{
			Code:    "enforcement_transition_assertion_unmet",
			Message: fmt.Sprintf("--expect-violation %q matched 0 violations (of %d total)", opts.expectViolation, len(violations)),
		}
		return result
	}
	// Regular mode.
	if len(violations) > 0 {
		result.Status = "failed"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{
			Code:        "enforcement_registry_transition_blocked",
			Message:     fmt.Sprintf("%d transition violation(s) — see violation.* checks above", len(violations)),
			Remediation: "Add [registry-status-change] trailer + rationale: line (R1); attach adr_reference for demotions (R2); ensure executor symbol exists in declared file before promoting to mechanical (R3). Opt-out via [skip-registry-transition] in commit body.",
		}
	}
	return result
}

// ─────────────────────────────────────────────────────────────────────
// commit-msg hook validator: validateEnforcementRegistryTransition
// ─────────────────────────────────────────────────────────────────────

// validateEnforcementRuleRegistrySync is the Phase 5 commit-msg validator
// (obligation.commit.enforcement_rule_registry_sync — 21st per_commit
// validator). It blocks commits that stage an enforcement/*.yaml or
// enforcement/*.md rule file without simultaneously registering it in
// enforcement-registry.yaml (either by an existing source_files binding
// or by staging the registry itself with a new rule_class entry).
//
// Relationship with Phase 3 compile-time orphan_rule lint: dual-gate.
// orphan_rule still fails compile-time on the entire registry; this
// commit-msg validator catches the same drift earlier with a tighter
// per-staged-file scope, so the failure surfaces at commit not at the
// next `ai-skill runtime compile` (which may be hours later).
//
// Opt-out: [skip-enforcement-registry-sync] trailer in commit body.
//
// Scope choice: only enforcement/ subtree for now. runtime/ and
// governance/ rule yamls are out of scope because they are typically
// edited alongside their owning module (their orphan_rule check at
// compile time is the primary gate). Expanding scope is a Phase 6+
// extension if the failure pattern surfaces there.
func validateEnforcementRuleRegistrySync(text string, staged []string, root string) string {
	if strings.Contains(text, "[skip-enforcement-registry-sync]") {
		return ""
	}
	registryRel := "enforcement/enforcement-registry.yaml"
	companionRel := "enforcement/enforcement-registry.md"

	registryStaged := false
	var candidates []string
	for _, p := range staged {
		rel := filepath.ToSlash(strings.TrimSpace(p))
		if rel == "" {
			continue
		}
		if rel == registryRel {
			registryStaged = true
			continue
		}
		if rel == companionRel {
			continue
		}
		if !strings.HasPrefix(rel, "enforcement/") {
			continue
		}
		if !strings.HasSuffix(rel, ".yaml") && !strings.HasSuffix(rel, ".yml") {
			continue
		}
		candidates = append(candidates, rel)
	}
	if len(candidates) == 0 {
		return ""
	}

	// Load current registry to check existing bindings. If registry is
	// missing the rule yaml binding AND the registry itself is staged,
	// assume the developer is adding the binding in the same commit
	// (trust + leave compile-time orphan_rule lint as belt).
	regPath := filepath.Join(root, filepath.FromSlash(registryRel))
	snap, err := loadRegistrySnapshotFromPath(regPath)
	if err != nil {
		// No registry to compare against — let compile-time lint handle it.
		return ""
	}
	bound := map[string]bool{}
	for _, rc := range snap.RuleClasses {
		for _, sf := range rc.SourceFiles {
			bound[normalizeSourcePath(sf)] = true
		}
	}

	var unbound []string
	for _, rel := range candidates {
		// Only flag files that actually declare a top-level `id:` — README
		// or notes files under enforcement/ without an id are not rule yamls.
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			continue
		}
		if extractTopLevelID(string(content)) == "" {
			continue
		}
		if bound[rel] {
			continue
		}
		unbound = append(unbound, rel)
	}
	if len(unbound) == 0 {
		return ""
	}
	if registryStaged {
		// Developer is editing the registry too — trust they are binding
		// the new files. Compile-time orphan_rule will catch any drift.
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "enforcement_rule_registry_sync: %d staged enforcement rule yaml(s) declare top-level id: but are not bound by enforcement-registry.yaml, AND enforcement-registry.yaml is not staged. Either:\n", len(unbound))
	for _, p := range unbound {
		fmt.Fprintf(&b, "  - %s\n", p)
	}
	b.WriteString("  Add a rule_class entry to enforcement/enforcement-registry.yaml (source_files: [<path>]) in this same commit, or use [skip-enforcement-registry-sync] opt-out.")
	return strings.TrimRight(b.String(), "\n")
}

// validateEnforcementRegistryTransition is the commit-msg validator that
// enforces Phase 4.5 R1/R2/R3 at commit time. Returns empty string on
// pass; non-empty error description blocks the commit.
//
// Triggers only when enforcement/enforcement-registry.yaml is staged.
// Reads HEAD's version via `git show HEAD:enforcement/enforcement-registry.yaml`
// and the staged version via the working-tree file (already updated by the
// developer before commit-msg fires).
//
// Opt-out: include [skip-registry-transition] in the commit body.
func validateEnforcementRegistryTransition(text string, staged []string, root string) string {
	registryRel := "enforcement/enforcement-registry.yaml"
	stagedHit := false
	for _, p := range staged {
		if filepath.ToSlash(strings.TrimSpace(p)) == registryRel {
			stagedHit = true
			break
		}
	}
	if !stagedHit {
		return ""
	}
	if strings.Contains(text, transitionOptOutMarker) {
		return ""
	}
	// Read the previous (HEAD) version. If there is no HEAD (initial
	// commit), treat old as empty registry — every rule_class will be
	// "(new)" which only triggers R1 + (none of the demotion paths).
	oldData, err := exec.Command("git", "-C", root, "show", "HEAD:"+registryRel).Output()
	if err != nil {
		oldData = []byte("schema_version: 2\nrule_classes: []\n")
	}
	newData, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(registryRel)))
	if err != nil {
		return fmt.Sprintf("enforcement_registry_transition: cannot read staged %s: %v", registryRel, err)
	}
	var oldSnap, newSnap registrySnapshot
	if err := yaml.Unmarshal(oldData, &oldSnap); err != nil {
		return fmt.Sprintf("enforcement_registry_transition: cannot parse HEAD %s: %v", registryRel, err)
	}
	if err := yaml.Unmarshal(newData, &newSnap); err != nil {
		return fmt.Sprintf("enforcement_registry_transition: cannot parse staged %s: %v", registryRel, err)
	}
	violations := checkRegistryTransitions(transitionInput{
		repoRoot:  root,
		oldSnap:   &oldSnap,
		newSnap:   &newSnap,
		commitMsg: text,
	})
	if len(violations) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "enforcement_registry_transition: %d violation(s) — fix or use [skip-registry-transition] opt-out:\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&b, "  - %s\n", v.String())
	}
	return strings.TrimRight(b.String(), "\n")
}
