package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func repoRootForLint(t *testing.T) string {
	t.Helper()
	// internal/app -> ../../../..  is repo root
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "..", "..", "..", ".."))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "enforcement", "enforcement-registry.yaml")); err != nil {
		t.Skipf("registry not found at %s: %v", root, err)
	}
	return root
}

// TestEnforcementRegistryLintDryRun is the Phase 3 step-5 dry-run probe.
// It does NOT fail the build on findings — it just prints them so the
// agent / human can decide backfill strategy before wiring lint into
// `ai-skill runtime compile`.
//
// To run:
//
//	go test ./scripts/ai-skill-cli/internal/app -run TestEnforcementRegistryLintDryRun -v
func TestEnforcementRegistryLintDryRun(t *testing.T) {
	root := repoRootForLint(t)
	errs, err := LintEnforcementRegistry(root)
	if err != nil {
		t.Fatalf("lint: %v", err)
	}
	if len(errs) == 0 {
		t.Log("PASS: zero lint findings")
		return
	}
	// Group by type for legibility.
	byType := map[string][]EnforcementRegistryLintError{}
	for _, e := range errs {
		byType[e.Type] = append(byType[e.Type], e)
	}
	types := make([]string, 0, len(byType))
	for k := range byType {
		types = append(types, k)
	}
	sort.Strings(types)
	var b strings.Builder
	fmt.Fprintf(&b, "\n=== Enforcement Registry Lint Dry-Run ===\nTotal findings: %d\n\n", len(errs))
	for _, typ := range types {
		fmt.Fprintf(&b, "── %s (%d) ──\n", typ, len(byType[typ]))
		for _, e := range byType[typ] {
			fmt.Fprintln(&b, e.Format())
			fmt.Fprintln(&b)
		}
	}
	t.Log(b.String())
}

// Fixture-based fail/pass tests per Phase 7 scenarios.

func writeRegistryFixture(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "enforcement"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "runtime"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "governance"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "scripts", "ai-skill-cli", "internal", "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "enforcement", "enforcement-registry.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestLintOrphanRule_Detected(t *testing.T) {
	dir := writeRegistryFixture(t, `
schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for: []
internal_helper_allowlist:
  symbols: []
rule_classes: []
`)
	// Create an orphan rule yaml.
	orphan := "id: enforcement.test-orphan-rule\nstatus: active\n"
	if err := os.WriteFile(filepath.Join(dir, "enforcement", "__test_orphan__.yaml"), []byte(orphan), 0o644); err != nil {
		t.Fatal(err)
	}
	errs, err := LintEnforcementRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range errs {
		if e.Type == "orphan_rule" {
			for _, f := range e.Fields {
				if f.Key == "declared_id" && f.Value == "enforcement.test-orphan-rule" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Fatalf("expected orphan_rule with declared_id enforcement.test-orphan-rule; got %d errs: %+v", len(errs), errs)
	}
}

func TestLintOrphanRule_BoundPasses(t *testing.T) {
	dir := writeRegistryFixture(t, `
schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for: []
internal_helper_allowlist:
  symbols: []
rule_classes:
  - id: bound
    coverage: behavioral_only
    source_files:
      - enforcement/__test_bound__.yaml
    rationale: "fixture"
    sunset_decision:
      revisit_when: "next quarter"
      success_criteria: "test"
`)
	bound := "id: enforcement.test-bound\nstatus: active\n"
	if err := os.WriteFile(filepath.Join(dir, "enforcement", "__test_bound__.yaml"), []byte(bound), 0o644); err != nil {
		t.Fatal(err)
	}
	errs, err := LintEnforcementRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range errs {
		if e.Type == "orphan_rule" {
			t.Fatalf("did not expect orphan_rule for bound yaml: %s", e.Format())
		}
	}
}

func TestLintMissingExecutorSymbol_Detected(t *testing.T) {
	dir := writeRegistryFixture(t, `
schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for:
    - commit_msg_validator
internal_helper_allowlist:
  symbols: []
rule_classes:
  - id: __test_missing_executor__
    coverage: mechanical
    source_files: []
    executors:
      - file: scripts/ai-skill-cli/internal/app/hooks.go
        symbol: validateNonexistentFakeSymbol
        executor_kind: commit_msg_validator
        block_or_warn: block
    rationale: "fixture"
`)
	// Empty hooks.go so the symbol cannot resolve.
	if err := os.WriteFile(filepath.Join(dir, "scripts", "ai-skill-cli", "internal", "app", "hooks.go"), []byte("package app\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	errs, err := LintEnforcementRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range errs {
		if e.Type == "missing_executor_symbol" {
			hasSym := false
			hasCls := false
			for _, f := range e.Fields {
				if f.Key == "expected_symbol" && f.Value == "validateNonexistentFakeSymbol" {
					hasSym = true
				}
				if f.Key == "rule_class" && f.Value == "__test_missing_executor__" {
					hasCls = true
				}
			}
			if hasSym && hasCls {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("expected missing_executor_symbol; got: %+v", errs)
	}
}

func TestLintMissingExecutorSymbol_PresentPasses(t *testing.T) {
	dir := writeRegistryFixture(t, `
schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for:
    - commit_msg_validator
internal_helper_allowlist:
  symbols: []
rule_classes:
  - id: __test_present__
    coverage: mechanical
    source_files: []
    executors:
      - file: scripts/ai-skill-cli/internal/app/hooks.go
        symbol: validateReal
        executor_kind: commit_msg_validator
        block_or_warn: block
    rationale: "fixture"
`)
	src := "package app\n\nfunc validateReal() string { return \"\" }\n"
	if err := os.WriteFile(filepath.Join(dir, "scripts", "ai-skill-cli", "internal", "app", "hooks.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	errs, err := LintEnforcementRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range errs {
		if e.Type == "missing_executor_symbol" {
			t.Fatalf("did not expect missing_executor_symbol: %s", e.Format())
		}
	}
}

func TestLintBehavioralIncompleteSunset_AllBadCases(t *testing.T) {
	cases := []struct {
		name    string
		yaml    string
		missing string
	}{
		{"no_sunset", `
rule_classes:
  - id: bad_no_sunset
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
`, "sunset_decision"},
		{"missing_revisit", `
rule_classes:
  - id: bad_missing_revisit
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      success_criteria: "x"
`, "sunset_decision.revisit_when"},
		{"missing_criteria", `
rule_classes:
  - id: bad_missing_criteria
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
`, "sunset_decision.success_criteria"},
		{"empty_strings", `
rule_classes:
  - id: bad_empty
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: ""
      success_criteria: ""
`, "sunset_decision.revisit_when"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			header := `schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for: []
internal_helper_allowlist:
  symbols: []
`
			dir := writeRegistryFixture(t, header+tc.yaml)
			errs, err := LintEnforcementRegistry(dir)
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, e := range errs {
				if e.Type == "behavioral_only_incomplete_sunset" {
					for _, f := range e.Fields {
						if f.Key == "missing_field" && f.Value == tc.missing {
							found = true
						}
					}
				}
			}
			if !found {
				t.Fatalf("%s: expected behavioral_only_incomplete_sunset(%s); got %+v", tc.name, tc.missing, errs)
			}
		})
	}
}

func TestLintBehavioralIncompleteSunset_GoodPasses(t *testing.T) {
	dir := writeRegistryFixture(t, `
schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for: []
internal_helper_allowlist:
  symbols: []
rule_classes:
  - id: good
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "fixture"
    sunset_decision:
      revisit_when: "Phase X lands"
      success_criteria: "Y achieves Z"
      revisit_owner: "maintainer"
`)
	errs, err := LintEnforcementRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range errs {
		if e.Type == "behavioral_only_incomplete_sunset" {
			t.Fatalf("did not expect behavioral_only_incomplete_sunset for good fixture: %s", e.Format())
		}
	}
}

func TestLintDeprecatedDisposal(t *testing.T) {
	header := `schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for: []
internal_helper_allowlist:
  symbols: []
`
	t.Run("missing_both", func(t *testing.T) {
		dir := writeRegistryFixture(t, header+`
rule_classes:
  - id: dep_missing
    coverage: deprecated
    source_files: []
    rationale: "x"
`)
		errs, _ := LintEnforcementRegistry(dir)
		found := false
		for _, e := range errs {
			if e.Type == "deprecated_missing_disposal" {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected deprecated_missing_disposal; got %+v", errs)
		}
	})
	t.Run("past_removal_date", func(t *testing.T) {
		dir := writeRegistryFixture(t, header+`
rule_classes:
  - id: dep_past
    coverage: deprecated
    source_files: []
    rationale: "x"
    removal_date: "2020-01-01"
`)
		errs, _ := LintEnforcementRegistry(dir)
		found := false
		for _, e := range errs {
			if e.Type == "deprecated_past_removal_date" {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected deprecated_past_removal_date; got %+v", errs)
		}
	})
	t.Run("future_date_passes", func(t *testing.T) {
		dir := writeRegistryFixture(t, header+`
rule_classes:
  - id: dep_ok
    coverage: deprecated
    source_files: []
    rationale: "x"
    removal_date: "2099-12-31"
`)
		errs, _ := LintEnforcementRegistry(dir)
		for _, e := range errs {
			if e.Type == "deprecated_missing_disposal" || e.Type == "deprecated_past_removal_date" {
				t.Fatalf("unexpected: %s", e.Format())
			}
		}
	})
}

// ─────────────────────────────────────────────────────────────────────
// Schema patch v2 tests — 9 new lint checks
// ─────────────────────────────────────────────────────────────────────

const v2Header = `schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for: []
internal_helper_allowlist:
  symbols: []
`

func hasFinding(errs []EnforcementRegistryLintError, typ string, sev string) bool {
	for _, e := range errs {
		if e.Type == typ {
			s := e.Severity
			if s == "" {
				s = SeverityFail
			}
			if s == sev {
				return true
			}
		}
	}
	return false
}

func TestLintBehavioralMissingRationale(t *testing.T) {
	// FAIL: rationale empty
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: bad
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    sunset_decision:
      revisit_when: "X"
      success_criteria: "20% coverage"
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "behavioral_only_missing_rationale", SeverityFail) {
		t.Fatalf("expected behavioral_only_missing_rationale FAIL; got %+v", errs)
	}
}

func TestLintBehavioralRecommendedFields_WarnOnly(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: bare
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "fixture"
    sunset_decision:
      revisit_when: "X"
      success_criteria: "20% coverage"
`)
	errs, _ := LintEnforcementRegistry(dir)
	for _, expected := range []string{
		"behavioral_only_missing_revisit_owner",
		"behavioral_only_missing_last_reviewed_at",
		"behavioral_only_missing_last_review_summary",
		"behavioral_only_missing_depends_on_rule_classes",
	} {
		if !hasFinding(errs, expected, SeverityWarn) {
			t.Errorf("expected %s WARNING; not found in %+v", expected, errs)
		}
		// Also ensure not FAIL:
		if hasFinding(errs, expected, SeverityFail) {
			t.Errorf("expected %s as WARNING not FAIL", expected)
		}
	}
}

func TestLintBehavioralReviewAge(t *testing.T) {
	// > 24 months → FAIL
	oldDate := time.Now().UTC().AddDate(-3, 0, 0).Format("2006-01-02")
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: stale
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "10% improvement"
      last_reviewed_at: "`+oldDate+`"
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "behavioral_only_review_age", SeverityFail) {
		t.Fatalf("expected behavioral_only_review_age FAIL for >24m old; got %+v", errs)
	}
}

func TestLintBehavioralReviewAge_MissingSkipped(t *testing.T) {
	// U1 mutual exclusion: missing last_reviewed_at → no review_age finding
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: fresh
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "20% coverage"
`)
	errs, _ := LintEnforcementRegistry(dir)
	for _, e := range errs {
		if e.Type == "behavioral_only_review_age" {
			t.Fatalf("review_age should be skipped when last_reviewed_at missing (U1 mutual exclusion): %s", e.Format())
		}
	}
}

func TestLintBehavioralVagueSuccessCriteria(t *testing.T) {
	for _, bad := range []string{"TBD", "future", "eventually", "TODO"} {
		t.Run(bad, func(t *testing.T) {
			dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: vague
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "improve `+bad+` someday"
`)
			errs, _ := LintEnforcementRegistry(dir)
			if !hasFinding(errs, "behavioral_only_vague_success_criteria", SeverityFail) {
				t.Errorf("expected vague_success_criteria FAIL for %q; got %+v", bad, errs)
			}
		})
	}
}

func TestLintBehavioralMissingMeasurableSignal(t *testing.T) {
	// "quality improves" — no digit, no %, no framework noun → WARNING
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: nosignal
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "quality improves significantly"
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "behavioral_only_missing_measurable_signal", SeverityWarn) {
		t.Errorf("expected missing_measurable_signal WARNING; got %+v", errs)
	}
}

func TestLintBehavioralMissingMeasurableSignal_PercentPasses(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: hassignal
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "coverage >= 80%"
`)
	errs, _ := LintEnforcementRegistry(dir)
	for _, e := range errs {
		if e.Type == "behavioral_only_missing_measurable_signal" {
			t.Fatalf("expected pass with %% token; got %s", e.Format())
		}
	}
}

func TestLintBehavioralRevisitChain_ChainOfBehavioralFails(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: upstream
    coverage: behavioral_only
    source_files: [enforcement/up.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "30% coverage"
  - id: downstream
    coverage: behavioral_only
    source_files: [enforcement/down.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "20% coverage"
      depends_on_rule_classes: [upstream]
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "behavioral_only_revisit_chain", SeverityFail) {
		t.Fatalf("expected behavioral_only_revisit_chain FAIL for behavioral→behavioral; got %+v", errs)
	}
}

func TestLintBehavioralRevisitChain_UnknownReference(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: dangler
    coverage: behavioral_only
    source_files: [enforcement/foo.md]
    rationale: "x"
    sunset_decision:
      revisit_when: "x"
      success_criteria: "20% coverage"
      depends_on_rule_classes: [does_not_exist]
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "behavioral_only_revisit_chain", SeverityFail) {
		t.Fatalf("expected revisit_chain FAIL for unknown ref; got %+v", errs)
	}
}

func TestLintUpstreamChainResolution_CycleDetected(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: a
    coverage: mechanical
    source_files: []
    executors: []
    rationale: "x"
    upstream_classes: [b]
  - id: b
    coverage: mechanical
    source_files: []
    executors: []
    rationale: "x"
    upstream_classes: [a]
`)
	errs, _ := LintEnforcementRegistry(dir)
	found := false
	for _, e := range errs {
		if e.Type == "upstream_chain_resolution" && e.IsFail() {
			for _, f := range e.Fields {
				if f.Key == "cycle" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Fatalf("expected upstream_chain_resolution cycle FAIL; got %+v", errs)
	}
}

func TestLintUpstreamChainResolution_UnknownReference(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: a
    coverage: mechanical
    source_files: []
    executors: []
    rationale: "x"
    upstream_classes: [ghost]
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "upstream_chain_resolution", SeverityFail) {
		t.Fatalf("expected upstream_chain_resolution FAIL for unknown ref; got %+v", errs)
	}
}

func TestLintClassSizeReviewThreshold(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
governance_thresholds:
  source_files_review_threshold: 3
rule_classes:
  - id: big
    coverage: mechanical
    source_files: [a.md, b.md, c.md, d.md, e.md]
    executors: []
    rationale: "x"
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "class_size_review_threshold", SeverityWarn) {
		t.Fatalf("expected class_size_review_threshold WARNING; got %+v", errs)
	}
}

func TestLintClassSizeReviewThreshold_WithExemptionStillEmitsWarning(t *testing.T) {
	// B6: exemption rationale does NOT suppress warning; it annotates it.
	dir := writeRegistryFixture(t, v2Header+`
governance_thresholds:
  source_files_review_threshold: 3
rule_classes:
  - id: big_but_justified
    coverage: mechanical
    source_files: [a.md, b.md, c.md, d.md, e.md]
    executors: []
    rationale: "x"
    size_review_exemption_rationale: "Cross-cutting writeback index; cohesive class."
`)
	errs, _ := LintEnforcementRegistry(dir)
	found := false
	for _, e := range errs {
		if e.Type == "class_size_review_threshold" {
			found = true
			gotRationale := false
			for _, f := range e.Fields {
				if f.Key == "exemption_rationale" {
					gotRationale = true
				}
			}
			if !gotRationale {
				t.Fatalf("expected exemption_rationale field on warning; got %s", e.Format())
			}
		}
	}
	if !found {
		t.Fatalf("warning should still emit even with exemption; got %+v", errs)
	}
}

func TestLintBaselineSnapshotGovernance_StrictModeSkips(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
bootstrap_mode: strict
rule_classes: []
`)
	errs, _ := LintEnforcementRegistry(dir)
	for _, e := range errs {
		if e.Type == "baseline_snapshot_missing_governance" {
			t.Fatalf("strict mode should skip baseline lint: %s", e.Format())
		}
	}
}

func TestLintBaselineSnapshotGovernance_MissingOwner(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
bootstrap_mode: baseline_snapshot_v1
baseline_snapshot:
  baseline_created_at: "2026-06-01"
  baseline_burndown_target_date: "2027-06-01"
  entries: []
rule_classes: []
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "baseline_snapshot_missing_governance", SeverityFail) {
		t.Fatalf("expected baseline_snapshot_missing_governance FAIL for missing owner; got %+v", errs)
	}
}

func TestLintBaselineSnapshotGovernance_EmptyEntrySummary(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
bootstrap_mode: baseline_snapshot_v1
baseline_snapshot:
  baseline_owner: "maintainer"
  entries:
    - finding_type: orphan_rule
      identifier: foo.yaml
      baseline_review_summary: ""
      accepted_at: "2026-06-01"
      burndown_owner: "maintainer"
rule_classes: []
`)
	errs, _ := LintEnforcementRegistry(dir)
	if !hasFinding(errs, "baseline_snapshot_missing_governance", SeverityFail) {
		t.Fatalf("expected baseline_snapshot_missing_governance FAIL for empty summary; got %+v", errs)
	}
}

func TestLintPendingImplementationChildPlanValidity(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: pending_bad
    coverage: pending_implementation
    source_files: []
    child_plan: plans/active/does-not-exist.md
    target_promotion: "later"
`)
	errs, _ := LintEnforcementRegistry(dir)
	found := false
	for _, e := range errs {
		if e.Type == "pending_implementation_child_plan_validity" && e.IsFail() {
			for _, f := range e.Fields {
				if f.Key == "violation" && f.Value == "(a) path_resolves" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Fatalf("expected child_plan_validity (a) FAIL for nonexistent path; got %+v", errs)
	}
}

func TestLintPendingImplementationChildPlanValidity_AnchorStripped(t *testing.T) {
	// B5: path with #anchor must be normalized before resolution
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: pending_ok
    coverage: pending_implementation
    source_files: []
    child_plan: "plans/active/stub.md#phase-3"
    target_promotion: "later"
`)
	// Create stub plan satisfying b/c/d
	if err := os.MkdirAll(filepath.Join(dir, "plans", "active"), 0o755); err != nil {
		t.Fatal(err)
	}
	stub := `# Stub Plan
Owner: maintainer

## Phase 0
preflight

## Validation Plan
- [ ] tested
`
	if err := os.WriteFile(filepath.Join(dir, "plans", "active", "stub.md"), []byte(stub), 0o644); err != nil {
		t.Fatal(err)
	}
	errs, _ := LintEnforcementRegistry(dir)
	for _, e := range errs {
		if e.Type == "pending_implementation_child_plan_validity" {
			t.Fatalf("expected anchor-stripped path to resolve and all rules to pass; got %s", e.Format())
		}
	}
}

// Step 7 wire: buildEnforcementRegistryLintCheck severity-aware behavior.

func TestEnforcementRegistryLintCheck_FailBlocks(t *testing.T) {
	// A mechanical class with a nonexistent executor symbol → FAIL → blocks.
	// Note: binding_required_for must include the kind, else missing-symbol
	// lint skips it (mirrors the real registry's executor_kind_spec).
	header := `schema_version: 2
id: enforcement.enforcement-registry
enforcement_mode:
  orphan_rule: fail
  orphan_executor: fail
executor_kind_spec:
  binding_required_for:
    - commit_msg_validator
internal_helper_allowlist:
  symbols: []
`
	dir := writeRegistryFixture(t, header+`
rule_classes:
  - id: broken
    coverage: mechanical
    source_files: []
    rationale: "x"
    executors:
      - file: scripts/ai-skill-cli/internal/app/hooks.go
        symbol: thisSymbolDoesNotExistAnywhere
        executor_kind: commit_msg_validator
        block_or_warn: block
`)
	// Empty hooks.go so symbol cannot resolve.
	if err := os.WriteFile(filepath.Join(dir, "scripts", "ai-skill-cli", "internal", "app", "hooks.go"), []byte("package app\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	check, blocks := buildEnforcementRegistryLintCheck(dir)
	if !blocks {
		t.Fatalf("expected FAIL to block; check=%+v", check)
	}
	if check.Status != "failed" {
		t.Fatalf("expected status failed; got %q", check.Status)
	}
	if !strings.Contains(check.Message, "FAIL:") {
		t.Fatalf("expected summary with FAIL count; got %q", check.Message)
	}
}

func TestEnforcementRegistryLintCheck_WarningDoesNotBlock(t *testing.T) {
	// class_size warning only → does not block, status warning.
	dir := writeRegistryFixture(t, v2Header+`
governance_thresholds:
  source_files_review_threshold: 2
rule_classes:
  - id: big
    coverage: mechanical
    source_files: [a.md, b.md, c.md]
    rationale: "x"
    executors: []
`)
	check, blocks := buildEnforcementRegistryLintCheck(dir)
	if blocks {
		t.Fatalf("warning must not block; check=%+v", check)
	}
	if check.Status != "warning" {
		t.Fatalf("expected status warning; got %q (%s)", check.Status, check.Message)
	}
	if !strings.Contains(check.Message, "Compile PASSED") {
		t.Fatalf("expected PASSED summary; got %q", check.Message)
	}
}

func TestEnforcementRegistryLintCheck_CleanIsOK(t *testing.T) {
	dir := writeRegistryFixture(t, v2Header+`
rule_classes: []
`)
	check, blocks := buildEnforcementRegistryLintCheck(dir)
	if blocks {
		t.Fatalf("clean registry must not block; %+v", check)
	}
	if check.Status != "ok" {
		t.Fatalf("expected status ok; got %q (%s)", check.Status, check.Message)
	}
}

func TestLintPendingImplementationChildPlanValidity_BCD_Warning(t *testing.T) {
	// path resolves (a OK) but body missing Phase 0 / owner / acceptance → WARNINGs
	dir := writeRegistryFixture(t, v2Header+`
rule_classes:
  - id: pending_thin
    coverage: pending_implementation
    source_files: []
    child_plan: plans/active/thin.md
    target_promotion: "later"
`)
	if err := os.MkdirAll(filepath.Join(dir, "plans", "active"), 0o755); err != nil {
		t.Fatal(err)
	}
	thin := "# Thin\n\nJust a title and one line.\n"
	if err := os.WriteFile(filepath.Join(dir, "plans", "active", "thin.md"), []byte(thin), 0o644); err != nil {
		t.Fatal(err)
	}
	errs, _ := LintEnforcementRegistry(dir)
	violations := map[string]bool{}
	for _, e := range errs {
		if e.Type == "pending_implementation_child_plan_validity" && e.Severity == SeverityWarn {
			for _, f := range e.Fields {
				if f.Key == "violation" {
					violations[f.Value] = true
				}
			}
		}
	}
	for _, want := range []string{"(b) phase_0_heading", "(c) owner_present", "(d) acceptance_section"} {
		if !violations[want] {
			t.Errorf("expected WARNING violation %q; got violations=%v", want, violations)
		}
	}
}
