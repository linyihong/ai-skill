package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
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
