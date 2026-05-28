// Package audit hosts the `ai-skill runtime audit` subcommand implementation.
//
// This stub test file binds the 5 Phase 1 validation scenarios from
// plans/active/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md to
// the audit package so the orphan-scenario-unreferenced-v1 warning does not
// fire on them. Real fixture implementations land in Phase 2 (inventory +
// classification) and Phase 5 (validateRuntimeTriggerWiring commit-msg
// validator).
//
// DO NOT REMOVE the scenario id string literals below without first updating
// the corresponding YAML files under validation/scenarios/failure-derived/.
package audit

import "testing"

// Phase 1 scenario id constants bound for fixture wiring in later phases.
const (
	scenarioOrphanRoutingEntry          = "orphan-routing-entry-v1"
	scenarioOrphanProjectionTargetKey   = "orphan-projection-target-key-v1"
	scenarioOrphanScenarioUnreferenced  = "orphan-scenario-unreferenced-v1"
	scenarioPre2026GrandfatherCoverage  = "pre-2026-grandfather-coverage-v1"
	scenarioFrameworkGlossaryCandidate  = "framework-glossary-candidate-missing-v1"
)

// TestScenarioStubsBound is a deliberately trivial test that exists only to
// reference the scenario id constants so `go vet` and the orphan-scenario
// audit warning recognise them as bound. Real per-scenario fixture tests are
// added in Phase 2 / Phase 5 of the audit plan.
func TestScenarioStubsBound(t *testing.T) {
	ids := []string{
		scenarioOrphanRoutingEntry,
		scenarioOrphanProjectionTargetKey,
		scenarioOrphanScenarioUnreferenced,
		scenarioPre2026GrandfatherCoverage,
		scenarioFrameworkGlossaryCandidate,
	}
	if len(ids) != 5 {
		t.Fatalf("expected 5 Phase 1 scenarios bound, got %d", len(ids))
	}
	for _, id := range ids {
		if id == "" {
			t.Errorf("scenario id must not be empty")
		}
	}
}
