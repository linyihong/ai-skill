package main

import (
	"os"
	"path/filepath"
	"testing"
)

// newBase creates a temp evidence-candidates dir with registry pointers for the
// given plans, so assemble's pointer-resolve can run without the real repo.
func newBase(t *testing.T, plans ...string) string {
	t.Helper()
	base := t.TempDir()
	reg := filepath.Join(base, "evidence-rules")
	if err := os.MkdirAll(reg, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, p := range plans {
		if err := os.WriteFile(filepath.Join(reg, p+".pointer.yaml"), []byte("plan_ref: x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return base
}

func validInput() input {
	in := input{
		Source:       source{Repo: "current_repo", Artifact: "plans/active/x.md", Commit: "abc"},
		MatchedPlans: []string{"economics"},
		CriteriaHits: []string{"owner_ambiguity"},
	}
	in.CriteriaSource.Actor = "human"
	return in
}

func TestAssembleValid(t *testing.T) {
	c, err := assemble(validInput(), newBase(t, "economics"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Status != "create" {
		t.Errorf("status = %q, want create (scanner must not set accept/discard/expire)", c.Status)
	}
	if c.ID == "" {
		t.Error("id not assigned")
	}
}

func TestDeterministicID(t *testing.T) {
	a := validInput()
	b := validInput()
	b.MatchedPlans = []string{"economics"} // same content, would-be different order is sorted
	if deterministicID(a) != deterministicID(b) {
		t.Error("same content must yield same id (idempotency)")
	}
	c := validInput()
	c.Source.Artifact = "plans/active/y.md"
	if deterministicID(a) == deterministicID(c) {
		t.Error("different artifact must yield different id")
	}
}

func TestInvariantSourceIsCandidate(t *testing.T) {
	in := validInput()
	in.Source.Artifact = "C-deadbeef"
	if _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: candidate must not reference another candidate")
	}
}

func TestInvariantCriteriaFromScanner(t *testing.T) {
	in := validInput()
	in.CriteriaSource.Actor = "scanner-v0"
	if _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: criteria_hits must originate outside scanner")
	}
}

func TestSchemaMissingActor(t *testing.T) {
	in := validInput()
	in.CriteriaSource.Actor = ""
	if _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: criteria_source.actor required")
	}
}

func TestSchemaEmptyCriteria(t *testing.T) {
	in := validInput()
	in.CriteriaHits = nil
	if _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: criteria_hits must be non-empty")
	}
}

func TestPointerResolveMissing(t *testing.T) {
	in := validInput()
	in.MatchedPlans = []string{"nonexistent-plan"}
	if _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: no registry pointer for matched_plan")
	}
}
