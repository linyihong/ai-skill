package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/planvalidate"
)

// Phase 3.3c — directional replacement proof (closes the Q1 evidence package).
// At least one REAL transport replacement (not a mock) must leave engine +
// remaining consumers + observation unchanged.

// corFromCLI extracts a COR from the manual-via-CLI adapter (a real adapter over
// the engine, distinct from the direct engineObservation adapter).
func corFromCLI(t *testing.T, root string) ObservationRecord {
	t.Helper()
	var out, errb bytes.Buffer
	Run([]string{"plans", "validate", "--root", root, "--format", "json"}, &out, &errb)
	var p cliPayload
	if err := json.Unmarshal(out.Bytes(), &p); err != nil {
		t.Fatalf("cli json: %v\n%s", err, out.String())
	}
	rec := newObservationRecord()
	for _, f := range p.Findings {
		rec.Findings[f.RuleID] = f.Blocking
	}
	return rec
}

// REAL replacement: swap the manual consumer's adapter (CLI <-> direct engine
// call). Both go through the same engine entrypoint; the observation must be
// preserved (R.3) including applicability.
func TestDirectionalReplacement_ManualAdapterSwap(t *testing.T) {
	root := violationTree(t)
	viaCLI := corFromCLI(t, root)
	viaDirect := engineObservation("", root) // no commit text -> no opt-out, matching CLI
	if !reflect.DeepEqual(viaCLI, viaDirect) {
		t.Fatalf("manual adapter swap changed observation:\n cli=%+v\n direct=%+v", viaCLI, viaDirect)
	}
	models, _ := normalizedPlansFromRoot(root)
	if !ruleApplicability(models)["plan_tree.parent_reference"] {
		t.Fatalf("parent_reference must be applicable on the violation tree (R.3)")
	}
}

// R.1 — removing the manual consumer activates no fallback: the remaining
// consumers (hook, CI) construct identical CORs whether or not manual is invoked.
func TestDirectionalReplacement_RemovalIndependence(t *testing.T) {
	root := violationTree(t)
	hook1 := legacyObservation("feat: x\n", []string{"plans/active/01.md"}, root)
	ci1 := ciObservation("feat: x\n", root)
	_ = corFromCLI(t, root) // invoke (then "remove") the manual consumer
	hook2 := legacyObservation("feat: x\n", []string{"plans/active/01.md"}, root)
	ci2 := ciObservation("feat: x\n", root)
	if !reflect.DeepEqual(hook1, hook2) || !reflect.DeepEqual(ci1, ci2) {
		t.Fatalf("manual presence/absence changed hook/CI COR — fallback/shared-state leak (R.1)")
	}
}

// R.3 anti-cheat — observation preservation is not set-equality: a rule that
// PASSED (applicable) and a rule that is SILENTLY INAPPLICABLE both yield zero
// findings but differ in applicability, so they are NOT observation-equal.
func TestDirectionalReplacement_ApplicabilityNotSetEquality(t *testing.T) {
	pass := validTree(t) // sub parent=m resolves -> parent_reference applicable + passes
	inap := t.TempDir()  // main only -> parent_reference inapplicable
	makePlan(t, inap, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-25\nparent: null\n---")

	passModels, _ := normalizedPlansFromRoot(pass)
	inapModels, _ := normalizedPlansFromRoot(inap)
	passAp := ruleApplicability(passModels)
	inapAp := ruleApplicability(inapModels)

	if !passAp["plan_tree.parent_reference"] {
		t.Fatalf("valid tree: parent_reference should be applicable")
	}
	if inapAp["plan_tree.parent_reference"] {
		t.Fatalf("main-only tree: parent_reference should be inapplicable")
	}
	// Both have empty parent_reference findings, but applicability differs.
	if reflect.DeepEqual(passAp, inapAp) {
		t.Fatalf("applicability must distinguish passed-vs-inapplicable")
	}
}

// R.2 — engine may evolve internally under the SAME contract without forcing
// adapter changes. Fingerprint the consumer-facing contract surface
// (ValidationContext + Finding); it must stay constant.
func contractFingerprint() string {
	var parts []string
	for _, tp := range []reflect.Type{
		reflect.TypeOf(planvalidate.ValidationContext{}),
		reflect.TypeOf(planvalidate.Finding{}),
	} {
		for i := 0; i < tp.NumField(); i++ {
			parts = append(parts, fmt.Sprintf("%s.%s:%s", tp.Name(), tp.Field(i).Name, tp.Field(i).Type.String()))
		}
	}
	sort.Strings(parts)
	out := ""
	for _, p := range parts {
		out += p + "|"
	}
	return out
}

func TestDirectionalReplacement_ContractStable(t *testing.T) {
	const golden = "Finding.Blocking:bool|Finding.Message:string|Finding.RuleID:string|ValidationContext.ChangedSet:planvalidate.ChangedSet|ValidationContext.ExecutionMode:planvalidate.ExecutionMode|ValidationContext.Metadata:planvalidate.ValidationMetadata|ValidationContext.Root:string|"
	if fp := contractFingerprint(); fp != golden {
		t.Fatalf("consumer contract surface changed -> adapters would need updating (R.2).\n got: %s\nwant: %s", fp, golden)
	}
}
