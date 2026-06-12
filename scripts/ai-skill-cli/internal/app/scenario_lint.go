package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LintValidationScenarios is the F19 (validation_scenario_governance)
// mechanical executor. It enforces that the validation scenarios bound to a
// rule_class via coverage_evidence are real, well-formed governance evidence:
//
//   1.1 structural lint of referenced scenarios (id/given/when/then = FAIL;
//       domain + validation.detection_command = WARNING)
//   1.2 coverage_evidence.{validation_scenarios,regression_scenarios} paths
//       must resolve to existing files (FAIL: dangling_coverage_ref)
//   1.3 coverage_target_pct governance floor on mechanical classes
//       (< 50 FAIL / 50-79 WARNING / >= 80 OK)
//   1.4 each regression_scenarios entry should link back to a real failure
//       pattern (WARNING: regression_unlinked_pattern)
//
// Scope decision (child plan Q1/Q2): structural lint applies ONLY to
// registry-referenced scenarios. The broader 207-scenario corpus contains
// many legitimately looser legacy formats; blanket FAIL would break compile
// and audit non-evidence files. detection_command is NOT a universal invariant
// (routing scenarios use then.validation[]), so it is WARNING-only.
//
// Findings reuse EnforcementRegistryLintError so the runtime compile pipeline
// can render them through the same severity-aware summary.
func LintValidationScenarios(repo string) ([]EnforcementRegistryLintError, error) {
	reg, _, err := loadRegistrySnapshot(repo)
	if err != nil {
		return nil, err
	}
	var errs []EnforcementRegistryLintError
	errs = append(errs, lintCoverageEvidencePathExistence(repo, reg)...)   // 1.2
	errs = append(errs, lintReferencedScenarioStructure(repo, reg)...)     // 1.1
	errs = append(errs, lintCoverageTargetFloor(reg)...)                   // 1.3
	errs = append(errs, lintRegressionPatternLinkage(repo, reg)...)        // 1.4
	sortLintErrors(errs)
	return errs, nil
}

// scenarioDoc captures only the fields the structural lint inspects. Each
// field is a yaml.Node so an absent key (Kind == 0) is distinguishable from a
// present-but-empty one.
type scenarioDoc struct {
	ID         yaml.Node `yaml:"id"`
	Domain     yaml.Node `yaml:"domain"`
	Given      yaml.Node `yaml:"given"`
	When       yaml.Node `yaml:"when"`
	Then       yaml.Node `yaml:"then"`
	Validation struct {
		DetectionCommand yaml.Node `yaml:"detection_command"`
	} `yaml:"validation"`
}

func nodePresent(n yaml.Node) bool { return n.Kind != 0 }

func nodeScalarNonEmpty(n yaml.Node) bool {
	return n.Kind != 0 && strings.TrimSpace(n.Value) != ""
}

// referencedScenarioPaths returns the de-duplicated set of scenario paths bound
// by any rule_class coverage_evidence (validation_scenarios + regression_scenarios),
// each tagged with the rule_class that referenced it for finding attribution.
type scenarioRef struct {
	rawPath   string
	pathOnly  string
	ruleClass string
	isRegress bool
}

func collectReferencedScenarios(reg *registrySnapshot) []scenarioRef {
	var refs []scenarioRef
	for _, rc := range reg.RuleClasses {
		if rc.CoverageEvidence == nil {
			continue
		}
		add := func(p string, regress bool) {
			refs = append(refs, scenarioRef{
				rawPath:   p,
				pathOnly:  stripAnchor(p),
				ruleClass: rc.ID,
				isRegress: regress,
			})
		}
		for _, p := range rc.CoverageEvidence.ValidationScenarios {
			add(p, false)
		}
		for _, p := range rc.CoverageEvidence.RegressionScenarios {
			add(p, true)
		}
	}
	return refs
}

func stripAnchor(p string) string {
	if idx := strings.Index(p, "#"); idx >= 0 {
		p = p[:idx]
	}
	return strings.TrimSpace(filepath.ToSlash(p))
}

// ─────────────────────────────────────────────────────────────────────
// 1.2 dangling_coverage_ref (FAIL): coverage_evidence scenario paths must
// resolve to existing files.
// ─────────────────────────────────────────────────────────────────────

func lintCoverageEvidencePathExistence(repo string, reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, ref := range collectReferencedScenarios(reg) {
		full := filepath.Join(repo, filepath.FromSlash(ref.pathOnly))
		if _, err := os.Stat(full); err != nil {
			kind := "validation_scenarios"
			if ref.isRegress {
				kind = "regression_scenarios"
			}
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "dangling_coverage_ref",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", ref.ruleClass},
					{"coverage_evidence_field", kind},
					{"scenario_path", ref.rawPath},
				},
				Message: "coverage_evidence references a validation scenario that does not exist. Create the scenario at that path, or repoint coverage_evidence to an existing scenario. coverage_evidence must be real evidence, not an aspirational filename.",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// 1.1 referenced_scenario_structure: core BDD skeleton (id/given/when/then)
// FAIL; quality fields (domain / validation.detection_command) WARNING.
// Only runs on referenced, existing scenario files.
// ─────────────────────────────────────────────────────────────────────

func lintReferencedScenarioStructure(repo string, reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	seen := map[string]bool{}
	for _, ref := range collectReferencedScenarios(reg) {
		if seen[ref.pathOnly] {
			continue
		}
		seen[ref.pathOnly] = true
		full := filepath.Join(repo, filepath.FromSlash(ref.pathOnly))
		data, err := os.ReadFile(full)
		if err != nil {
			continue // existence handled by 1.2
		}
		var doc scenarioDoc
		if uerr := yaml.Unmarshal(data, &doc); uerr != nil {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "referenced_scenario_structure",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", ref.ruleClass},
					{"scenario_path", ref.pathOnly},
					{"parse_error", uerr.Error()},
				},
				Message: "referenced scenario yaml does not parse. Fix the YAML syntax.",
			})
			continue
		}
		fail := func(field string) {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "referenced_scenario_structure",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", ref.ruleClass},
					{"scenario_path", ref.pathOnly},
					{"missing_field", field},
				},
				Message: "referenced scenario missing required structural field. id/given/when/then are the BDD skeleton invariant; a scenario bound as coverage_evidence must declare them.",
			})
		}
		warn := func(field, why string) {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "referenced_scenario_structure",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", ref.ruleClass},
					{"scenario_path", ref.pathOnly},
					{"missing_field", field},
				},
				Message: why,
			})
		}
		if !nodeScalarNonEmpty(doc.ID) {
			fail("id")
		}
		if !nodePresent(doc.Given) {
			fail("given")
		}
		if !nodePresent(doc.When) {
			fail("when")
		}
		if !nodePresent(doc.Then) {
			fail("then")
		}
		if !nodeScalarNonEmpty(doc.Domain) {
			warn("domain", "referenced scenario missing domain (quality signal, not blocker). Add domain to classify the scenario.")
		}
		if !nodeScalarNonEmpty(doc.Validation.DetectionCommand) {
			warn("validation.detection_command", "referenced scenario has no validation.detection_command (quality signal). Routing-style scenarios legitimately use then.validation[] instead; add a detection_command if this scenario should be auto-verifiable.")
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// 1.3 coverage_target_below_floor: mechanical classes declaring a
// coverage_target_pct must clear the governance floor (< 50 FAIL,
// 50-79 WARNING). A mechanical class with coverage_evidence but no
// coverage_target_pct gets a WARNING to declare one.
// ─────────────────────────────────────────────────────────────────────

func lintCoverageTargetFloor(reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, rc := range reg.RuleClasses {
		if rc.Coverage != "mechanical" || rc.CoverageEvidence == nil {
			continue
		}
		pct := rc.CoverageEvidence.CoverageTargetPct
		if pct == nil {
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "coverage_target_below_floor",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"coverage_target_pct", "(absent)"},
				},
				Message: "mechanical rule_class with coverage_evidence should declare coverage_target_pct so the governance floor (>= 80 OK / 50-79 WARNING / < 50 FAIL) can be enforced.",
			})
			continue
		}
		switch {
		case *pct < 50:
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "coverage_target_below_floor",
				Severity: SeverityFail,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"coverage_target_pct", fmt.Sprintf("%d", *pct)},
					{"floor", "50"},
				},
				Message: "mechanical rule_class declares coverage_target_pct below the 50% FAIL floor. A class claiming mechanical coverage must target at least 50% (recommended >= 80%). Raise the target or reconsider the coverage classification.",
			})
		case *pct < 80:
			errs = append(errs, EnforcementRegistryLintError{
				Type:     "coverage_target_below_floor",
				Severity: SeverityWarn,
				Fields: []EnforcementRegistryLintField{
					{"rule_class", rc.ID},
					{"coverage_target_pct", fmt.Sprintf("%d", *pct)},
					{"recommended_floor", "80"},
				},
				Message: "mechanical rule_class coverage_target_pct is in the 50-79 WARNING band; consider raising toward the >= 80 recommended floor.",
			})
		}
	}
	sortLintErrors(errs)
	return errs
}

// ─────────────────────────────────────────────────────────────────────
// 1.4 regression_unlinked_pattern (WARNING): each regression scenario should
// trace back to a real failure pattern, via a failure_source: block, a
// reference to enforcement/failure-patterns/<x>.md, or living under
// validation/scenarios/failure-derived/.
// ─────────────────────────────────────────────────────────────────────

func lintRegressionPatternLinkage(repo string, reg *registrySnapshot) []EnforcementRegistryLintError {
	var errs []EnforcementRegistryLintError
	for _, ref := range collectReferencedScenarios(reg) {
		if !ref.isRegress {
			continue
		}
		full := filepath.Join(repo, filepath.FromSlash(ref.pathOnly))
		data, err := os.ReadFile(full)
		if err != nil {
			continue // existence handled by 1.2
		}
		body := string(data)
		linked := strings.Contains(ref.pathOnly, "validation/scenarios/failure-derived/") ||
			containsFailureSourceKey(body) ||
			strings.Contains(body, "enforcement/failure-patterns/")
		if linked {
			continue
		}
		errs = append(errs, EnforcementRegistryLintError{
			Type:     "regression_unlinked_pattern",
			Severity: SeverityWarn,
			Fields: []EnforcementRegistryLintField{
				{"rule_class", ref.ruleClass},
				{"scenario_path", ref.pathOnly},
			},
			Message: "regression scenario does not link back to a failure pattern. Add a failure_source: block, reference enforcement/failure-patterns/<x>.md, or place it under validation/scenarios/failure-derived/ so the regression's origin is traceable.",
		})
	}
	sortLintErrors(errs)
	return errs
}

func containsFailureSourceKey(body string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "failure_source:") {
			return true
		}
	}
	return false
}
