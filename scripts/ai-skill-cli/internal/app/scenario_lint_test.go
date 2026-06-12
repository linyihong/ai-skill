package app

import (
	"os"
	"path/filepath"
	"testing"
)

// writeScenarioLintRepo builds a temp repo with a registry yaml plus any
// scenario files, then returns the repo root for LintValidationScenarios.
func writeScenarioLintRepo(t *testing.T, registryYAML string, scenarios map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "enforcement"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "enforcement", "enforcement-registry.yaml"), []byte(registryYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	for rel, content := range scenarios {
		full := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func findingsByType(errs []EnforcementRegistryLintError, typ string) []EnforcementRegistryLintError {
	var out []EnforcementRegistryLintError
	for _, e := range errs {
		if e.Type == typ {
			out = append(out, e)
		}
	}
	return out
}

func countSeverity(errs []EnforcementRegistryLintError, typ, sev string) int {
	n := 0
	for _, e := range findingsByType(errs, typ) {
		if e.Severity == sev {
			n++
		}
	}
	return n
}

const wellFormedScenario = `id: sample-v1
domain: enforcement
type: failure-recovery
priority: P1
failure_source:
  pattern: some-pattern
given:
  a: b
when:
  action: x
then:
  expected_route:
    - ok
validation:
  detection_command: |
    echo PASS
`

// 1.2 dangling_coverage_ref
func TestScenarioLintDanglingRef(t *testing.T) {
	reg := `rule_classes:
  - id: demo
    coverage: mechanical
    coverage_evidence:
      coverage_target_pct: 100
      validation_scenarios:
        - validation/scenarios/x/missing.yaml
`
	// FAIL case: referenced file does not exist.
	root := writeScenarioLintRepo(t, reg, nil)
	errs, err := LintValidationScenarios(root)
	if err != nil {
		t.Fatalf("lint: %v", err)
	}
	if got := countSeverity(errs, "dangling_coverage_ref", SeverityFail); got != 1 {
		t.Fatalf("dangling FAIL: want 1, got %d (%v)", got, errs)
	}

	// PASS case: referenced file exists.
	root2 := writeScenarioLintRepo(t, reg, map[string]string{
		"validation/scenarios/x/missing.yaml": wellFormedScenario,
	})
	errs2, _ := LintValidationScenarios(root2)
	if got := len(findingsByType(errs2, "dangling_coverage_ref")); got != 0 {
		t.Fatalf("dangling PASS: want 0, got %d (%v)", got, errs2)
	}
}

// 1.1 referenced_scenario_structure
func TestScenarioLintStructure(t *testing.T) {
	reg := `rule_classes:
  - id: demo
    coverage: mechanical
    coverage_evidence:
      coverage_target_pct: 100
      validation_scenarios:
        - validation/scenarios/x/s.yaml
`
	// FAIL: missing `then` (skeleton).
	noThen := `id: s
domain: enforcement
given: {a: b}
when: {action: x}
validation:
  detection_command: "echo PASS"
`
	root := writeScenarioLintRepo(t, reg, map[string]string{"validation/scenarios/x/s.yaml": noThen})
	errs, _ := LintValidationScenarios(root)
	if got := countSeverity(errs, "referenced_scenario_structure", SeverityFail); got < 1 {
		t.Fatalf("structure FAIL(then): want >=1, got %d (%v)", got, errs)
	}

	// WARNING: routing scenario missing detection_command + domain only.
	noDetect := `id: s
given: {a: b}
when: {action: x}
then:
  validation:
    - assert something
`
	root2 := writeScenarioLintRepo(t, reg, map[string]string{"validation/scenarios/x/s.yaml": noDetect})
	errs2, _ := LintValidationScenarios(root2)
	if got := countSeverity(errs2, "referenced_scenario_structure", SeverityFail); got != 0 {
		t.Fatalf("structure: routing scenario should not FAIL on missing detection_command, got %d (%v)", got, errs2)
	}
	if got := countSeverity(errs2, "referenced_scenario_structure", SeverityWarn); got < 2 {
		t.Fatalf("structure WARN(domain+detection_command): want >=2, got %d (%v)", got, errs2)
	}

	// PASS: fully-formed scenario.
	root3 := writeScenarioLintRepo(t, reg, map[string]string{"validation/scenarios/x/s.yaml": wellFormedScenario})
	errs3, _ := LintValidationScenarios(root3)
	if got := len(findingsByType(errs3, "referenced_scenario_structure")); got != 0 {
		t.Fatalf("structure PASS: want 0, got %d (%v)", got, errs3)
	}
}

// 1.3 coverage_target_below_floor
func TestScenarioLintCoverageFloor(t *testing.T) {
	mk := func(pct string) string {
		return `rule_classes:
  - id: demo
    coverage: mechanical
    coverage_evidence:
` + pct + `      validation_scenarios:
        - validation/scenarios/x/s.yaml
`
	}
	scen := map[string]string{"validation/scenarios/x/s.yaml": wellFormedScenario}

	// FAIL: < 50.
	root := writeScenarioLintRepo(t, mk("      coverage_target_pct: 30\n"), scen)
	errs, _ := LintValidationScenarios(root)
	if got := countSeverity(errs, "coverage_target_below_floor", SeverityFail); got != 1 {
		t.Fatalf("floor FAIL(30): want 1, got %d (%v)", got, errs)
	}

	// WARNING: 50-79.
	root2 := writeScenarioLintRepo(t, mk("      coverage_target_pct: 60\n"), scen)
	errs2, _ := LintValidationScenarios(root2)
	if got := countSeverity(errs2, "coverage_target_below_floor", SeverityWarn); got != 1 {
		t.Fatalf("floor WARN(60): want 1, got %d (%v)", got, errs2)
	}

	// OK: >= 80.
	root3 := writeScenarioLintRepo(t, mk("      coverage_target_pct: 90\n"), scen)
	errs3, _ := LintValidationScenarios(root3)
	if got := len(findingsByType(errs3, "coverage_target_below_floor")); got != 0 {
		t.Fatalf("floor OK(90): want 0, got %d (%v)", got, errs3)
	}

	// WARNING: absent target on mechanical class with coverage_evidence.
	root4 := writeScenarioLintRepo(t, mk(""), scen)
	errs4, _ := LintValidationScenarios(root4)
	if got := countSeverity(errs4, "coverage_target_below_floor", SeverityWarn); got != 1 {
		t.Fatalf("floor WARN(absent): want 1, got %d (%v)", got, errs4)
	}
}

// 1.4 regression_unlinked_pattern
func TestScenarioLintRegressionLinkage(t *testing.T) {
	reg := `rule_classes:
  - id: demo
    coverage: mechanical
    coverage_evidence:
      coverage_target_pct: 100
      regression_scenarios:
        - validation/scenarios/x/r.yaml
`
	// WARNING: regression scenario with no failure linkage.
	unlinked := `id: r
domain: enforcement
given: {a: b}
when: {action: x}
then: {expected_route: [ok]}
validation:
  detection_command: "echo PASS"
`
	root := writeScenarioLintRepo(t, reg, map[string]string{"validation/scenarios/x/r.yaml": unlinked})
	errs, _ := LintValidationScenarios(root)
	if got := countSeverity(errs, "regression_unlinked_pattern", SeverityWarn); got != 1 {
		t.Fatalf("regression WARN: want 1, got %d (%v)", got, errs)
	}

	// PASS: regression scenario with failure_source block.
	root2 := writeScenarioLintRepo(t, reg, map[string]string{"validation/scenarios/x/r.yaml": wellFormedScenario})
	errs2, _ := LintValidationScenarios(root2)
	if got := len(findingsByType(errs2, "regression_unlinked_pattern")); got != 0 {
		t.Fatalf("regression PASS: want 0, got %d (%v)", got, errs2)
	}
}

// Real-repo guard: the live registry + referenced scenarios must lint clean of
// FAIL findings (WARNINGs allowed) so `runtime compile` stays green.
func TestScenarioLintRealRepoNoFail(t *testing.T) {
	root := repoRootForLint(t)
	errs, err := LintValidationScenarios(root)
	if err != nil {
		t.Fatalf("lint: %v", err)
	}
	for _, e := range errs {
		if e.IsFail() {
			t.Errorf("unexpected FAIL finding: %s", e.Format())
		}
	}
}
