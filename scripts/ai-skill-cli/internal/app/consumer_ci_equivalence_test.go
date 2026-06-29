package app

import (
	"reflect"
	"strings"
	"testing"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/planvalidate"
)

// Phase 3.3b — hook ↔ CI consumer equivalence via the Canonical Observation
// Record. CI is modeled as an engine-backed adapter (F.1): it builds a
// ValidationContext, calls the SAME engine entrypoint, applies its own policy,
// and renders — it owns no validation/schema logic. So the CI COR is exactly the
// engine observation over the normalized checkout (normalizedPlansFromRoot reads
// the full working tree = the CI checkout; the snapshot normalization lives in
// the adapter/loader, never the engine — F.2).
//
// Equivalence is judged on COR (CI COR == hook COR), never on exit code (exit MAY
// differ).

// ciObservation is the CI consumer: engine-backed, full-checkout snapshot.
func ciObservation(text, root string) ObservationRecord { return engineObservation(text, root) }

func TestConsumerEquivalence_HookVsCI_Normalized(t *testing.T) {
	cases := []struct {
		name, text string
		root       func(*testing.T) string
		staged     []string
	}{
		{"valid", "feat: x\n", validTree, []string{"plans/active/01.md"}},
		{"violation", "feat: x\n", violationTree, []string{"plans/active/01.md"}},
		{"opt-out", "feat: x\n\n[skip-plan-tree-parent-reference]\n", violationTree, []string{"plans/active/01.md"}},
	}
	for _, c := range cases {
		root := c.root(t)
		hook := legacyObservation(c.text, c.staged, root)
		ci := ciObservation(c.text, root)
		if !reflect.DeepEqual(hook, ci) {
			t.Fatalf("%s: hook != normalized-CI COR\n hook=%+v\n ci=%+v", c.name, hook, ci)
		}
	}
}

// F.3 asymmetric proof: equivalence is a contract, not a both-empty coincidence.
// On a checkout, a CI that naively reused the hook's staged semantics sees an
// EMPTY staged set and would miss the violation (raw CI transport != hook). The
// CI adapter that NORMALIZES the snapshot (full checkout) reproduces the hook COR.
func TestConsumerEquivalence_HookVsCI_AsymmetricProof(t *testing.T) {
	root := violationTree(t)
	text := "feat: x\n"

	hook := legacyObservation(text, []string{"plans/active/01.md"}, root)
	rawCI := legacyObservation(text, []string{}, root)   // naive: staged empty on checkout
	ciNorm := ciObservation(text, root)                  // adapter-normalized full checkout

	if hook.Findings["plan_tree.parent_reference"] != true {
		t.Fatalf("precondition: hook should see the violation")
	}
	if reflect.DeepEqual(hook.Findings, rawCI.Findings) {
		t.Fatalf("asymmetry expected: raw CI (staged-empty) should differ from hook, both=%+v", hook.Findings)
	}
	if !reflect.DeepEqual(hook.Findings, ciNorm.Findings) {
		t.Fatalf("normalized CI must equal hook: hook=%+v ciNorm=%+v", hook.Findings, ciNorm.Findings)
	}
}

// F.2 guard: the engine input contract must not gain snapshot-origin awareness,
// or the equivalence matrix ("Input snapshot MAY differ") would be undermined.
func TestValidationContext_NoSnapshotOriginField(t *testing.T) {
	tp := reflect.TypeOf(planvalidate.ValidationContext{})
	for i := 0; i < tp.NumField(); i++ {
		name := strings.ToLower(tp.Field(i).Name)
		if strings.Contains(name, "snapshot") || strings.Contains(name, "origin") || strings.Contains(name, "staged") {
			t.Fatalf("ValidationContext field %q leaks snapshot origin into the engine (F.2)", tp.Field(i).Name)
		}
	}
}
