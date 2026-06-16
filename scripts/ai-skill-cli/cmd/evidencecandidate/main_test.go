package main

import (
	"os"
	"path/filepath"
	"testing"
)

// newBase creates a temp evidence-candidates dir with RESOLVED registry pointers
// for the given plans, so assemble's status-aware resolve can run.
func newBase(t *testing.T, plans ...string) string {
	t.Helper()
	return newBaseWithStatus(t, "resolved", plans...)
}

func newBaseWithStatus(t *testing.T, status string, plans ...string) string {
	t.Helper()
	base := t.TempDir()
	reg := filepath.Join(base, "evidence-rules")
	if err := os.MkdirAll(reg, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, p := range plans {
		body := "plan_ref: x\nstatus: " + status + "\n"
		if err := os.WriteFile(filepath.Join(reg, p+".pointer.yaml"), []byte(body), 0o644); err != nil {
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
	c, pending, err := assemble(validInput(), newBase(t, "economics"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("resolved pointer should not be pending, got %v", pending)
	}
	if c.Status != "create" {
		t.Errorf("status = %q, want create (scanner must not set accept/discard/expire)", c.Status)
	}
	if c.ID == "" {
		t.Error("id not assigned")
	}
}

// index != consumable: a section_pending pointer must NOT be consumable —
// no error, but no candidate emitted (returned in pending).
func TestSectionPendingNotConsumable(t *testing.T) {
	base := newBaseWithStatus(t, "section_pending", "economics")
	c, pending, err := assemble(validInput(), base)
	if err != nil {
		t.Fatalf("section_pending must not be a hard error (exit 0), got %v", err)
	}
	if len(pending) == 0 {
		t.Fatal("section_pending pointer must be reported as pending (not consumable)")
	}
	if c.ID != "" {
		t.Errorf("no candidate must be assembled for a non-resolved pointer, got id %q", c.ID)
	}
}

func TestDeterministicID(t *testing.T) {
	a := validInput()
	b := validInput()
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
	if _, _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: candidate must not reference another candidate")
	}
}

func TestInvariantCriteriaFromScanner(t *testing.T) {
	in := validInput()
	in.CriteriaSource.Actor = "scanner-v0"
	if _, _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: criteria_hits must originate outside scanner")
	}
}

func TestSchemaMissingActor(t *testing.T) {
	in := validInput()
	in.CriteriaSource.Actor = ""
	if _, _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: criteria_source.actor required")
	}
}

func TestSchemaEmptyCriteria(t *testing.T) {
	in := validInput()
	in.CriteriaHits = nil
	if _, _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: criteria_hits must be non-empty")
	}
}

func TestPointerResolveMissing(t *testing.T) {
	in := validInput()
	in.MatchedPlans = []string{"nonexistent-plan"}
	if _, _, err := assemble(in, newBase(t, "economics")); err == nil {
		t.Error("expected reject: no registry pointer for matched_plan")
	}
}
