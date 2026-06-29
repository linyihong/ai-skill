package planvalidate

import (
	"strings"
	"testing"
)

func TestExtractFencedYAMLBlocks(t *testing.T) {
	md := "intro\n```yaml\nsteps:\n  - id: a\n```\n"
	blocks := ExtractFencedYAMLBlocks(md)
	if len(blocks) != 1 || !strings.Contains(blocks[0], "id: a") {
		t.Fatalf("unexpected blocks: %#v", blocks)
	}
}

func TestDetectIllegalIntentTransitions_structureToFeatureWithoutEquivalence(t *testing.T) {
	steps := []ImplementationStep{
		{ID: "prep-01", Intent: "structure", EquivalenceRequired: true},
		{ID: "feat-01", Intent: "feature"},
	}
	fs := DetectIllegalIntentTransitions(steps)
	if len(fs) != 1 || fs[0].RuleID != "implementation.intent.illegal_transition" {
		t.Fatalf("expected illegal transition, got %#v", fs)
	}
	if fs[0].Blocking {
		t.Fatal("dogfood advisory findings must not block")
	}
}

func TestDetectIllegalIntentTransitions_happyPath(t *testing.T) {
	steps := []ImplementationStep{
		{ID: "prep-01", Intent: "structure", EquivalenceRequired: true, EquivalenceEvidence: "TestParse"},
		{ID: "feat-01", Intent: "feature"},
	}
	if fs := DetectIllegalIntentTransitions(steps); len(fs) != 0 {
		t.Fatalf("expected no findings, got %#v", fs)
	}
}

func TestDetectIllegalIntentTransitions_featureToStructureWithoutReason(t *testing.T) {
	steps := []ImplementationStep{
		{ID: "feat-01", Intent: "feature"},
		{ID: "prep-02", Intent: "structure"},
	}
	fs := DetectIllegalIntentTransitions(steps)
	if len(fs) != 1 {
		t.Fatalf("expected illegal transition, got %#v", fs)
	}
}

func TestAdvisoryValidateImplementationIntent_forceExitDogfoodFixture(t *testing.T) {
	md := `# Implementation Plan: advisory illegal-transition scan

` + "```yaml\n" + `change_kind: feature
blocked_by_structure: true
execution_mode: preparatory_refactoring
steps:
  - id: prep-01
    intent: structure
    behavior_change:
      allowed: false
    checkpoint:
      observable_equivalence:
        required: true
        evidence: TestExtractFencedYAMLBlocks
  - id: prep-02
    intent: structure
    behavior_change:
      allowed: false
    checkpoint:
      observable_equivalence:
        required: true
        evidence: TestImplementationPlanParserDropped
  - id: feat-01
    intent: feature
    behavior_change:
      allowed: true
` + "```\n"

	fs, err := AdvisoryValidateImplementationIntent(md)
	if err != nil {
		t.Fatal(err)
	}
	// prep-02 had evidence but feature follows — legal transition from prep-02
	if len(fs) != 0 {
		t.Fatalf("force_exit shrunk path should pass transitions, got %#v", fs)
	}
}

func TestAdvisoryValidateImplementationIntent_fakeEquivalenceLens(t *testing.T) {
	steps := []ImplementationStep{
		{ID: "prep-01", Intent: "structure", EquivalenceRequired: true, EquivalenceEvidence: "manual note only"},
		{ID: "feat-01", Intent: "feature"},
	}
	// Parser would accept non-empty evidence; reviewer must catch fake equivalence manually.
	if fs := DetectIllegalIntentTransitions(steps); len(fs) != 0 {
		t.Fatalf("non-empty evidence satisfies mechanical check: %#v", fs)
	}
}
