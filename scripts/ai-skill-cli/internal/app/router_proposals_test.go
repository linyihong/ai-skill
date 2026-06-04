package app

import (
	"testing"
	"time"
)

var rpNow = time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

func daysAgo(n int) string {
	return rpNow.Add(-time.Duration(n) * 24 * time.Hour).Format(time.RFC3339)
}

func TestRecordProposalOccurrence_AddThenBump(t *testing.T) {
	var store routeProposalStore
	recordProposalOccurrence(&store, "governance-audit", []string{"audit"}, []string{"governance/x"}, "workflow", rpNow)
	if len(store.Proposals) != 1 || store.Proposals[0].OccurrenceCount != 1 {
		t.Fatalf("expected 1 new proposal count=1, got %+v", store.Proposals)
	}
	if store.Proposals[0].Status != proposalStatusAccumulating {
		t.Fatalf("new proposal must be accumulating, got %s", store.Proposals[0].Status)
	}
	recordProposalOccurrence(&store, "governance-audit", []string{"compliance"}, nil, "workflow", rpNow)
	if len(store.Proposals) != 1 || store.Proposals[0].OccurrenceCount != 2 {
		t.Fatalf("bump should keep one entry count=2, got %+v", store.Proposals)
	}
	// signals merged + deduped
	if len(store.Proposals[0].SuggestedUserSignals) != 2 {
		t.Fatalf("signals merge wrong: %v", store.Proposals[0].SuggestedUserSignals)
	}
}

func TestRecordProposalOccurrence_TerminalNotReopened(t *testing.T) {
	store := routeProposalStore{Proposals: []RouteProposal{
		{CandidateID: "x", OccurrenceCount: 9, Status: proposalStatusPromoted, LastSeen: daysAgo(1)},
	}}
	recordProposalOccurrence(&store, "x", []string{"y"}, nil, "workflow", rpNow)
	if store.Proposals[0].OccurrenceCount != 9 || store.Proposals[0].Status != proposalStatusPromoted {
		t.Fatalf("promoted proposal must not be re-opened, got %+v", store.Proposals[0])
	}
}

func TestApplyLifecycle_AccumulatingToReady(t *testing.T) {
	store := routeProposalStore{Proposals: []RouteProposal{
		{CandidateID: "ready", OccurrenceCount: 5, Status: proposalStatusAccumulating, LastSeen: daysAgo(10)},
	}}
	tr, _ := applyProposalLifecycle(&store, rpNow)
	if tr != 1 || store.Proposals[0].Status != proposalStatusReadyForReview {
		t.Fatalf("expected ready_for_review, got %+v (tr=%d)", store.Proposals[0], tr)
	}
}

func TestApplyLifecycle_BelowThresholdStaysAccumulating(t *testing.T) {
	store := routeProposalStore{Proposals: []RouteProposal{
		{CandidateID: "young", OccurrenceCount: 4, Status: proposalStatusAccumulating, LastSeen: daysAgo(10)},
	}}
	applyProposalLifecycle(&store, rpNow)
	if store.Proposals[0].Status != proposalStatusAccumulating {
		t.Fatalf("occurrence<5 & recent must stay accumulating, got %s", store.Proposals[0].Status)
	}
}

func TestApplyLifecycle_RecentHighCountNotStaleEvenIfManyOccurrences(t *testing.T) {
	// occurrence>=5 but last_seen old (40d): not ready (recency fails), and not
	// stale (stale requires occurrence<5). Stays accumulating.
	store := routeProposalStore{Proposals: []RouteProposal{
		{CandidateID: "oldbusy", OccurrenceCount: 6, Status: proposalStatusAccumulating, LastSeen: daysAgo(40)},
	}}
	applyProposalLifecycle(&store, rpNow)
	if store.Proposals[0].Status != proposalStatusAccumulating {
		t.Fatalf("stale recency-only with high count should stay accumulating, got %s", store.Proposals[0].Status)
	}
}

func TestApplyLifecycle_AccumulatingToStale(t *testing.T) {
	store := routeProposalStore{Proposals: []RouteProposal{
		{CandidateID: "old", OccurrenceCount: 2, Status: proposalStatusAccumulating, LastSeen: daysAgo(70)},
	}}
	tr, _ := applyProposalLifecycle(&store, rpNow)
	if tr != 1 || store.Proposals[0].Status != proposalStatusStale {
		t.Fatalf("expected stale, got %+v", store.Proposals[0])
	}
}

func TestApplyLifecycle_PrunesVeryOldStale(t *testing.T) {
	store := routeProposalStore{Proposals: []RouteProposal{
		{CandidateID: "ancient", OccurrenceCount: 1, Status: proposalStatusStale, LastSeen: daysAgo(100)},
		{CandidateID: "keep", OccurrenceCount: 5, Status: proposalStatusReadyForReview, LastSeen: daysAgo(5)},
	}}
	_, pruned := applyProposalLifecycle(&store, rpNow)
	if pruned != 1 {
		t.Fatalf("expected 1 pruned, got %d", pruned)
	}
	if len(store.Proposals) != 1 || store.Proposals[0].CandidateID != "keep" {
		t.Fatalf("only 'keep' should remain, got %+v", store.Proposals)
	}
}

func TestApplyLifecycle_StaleReactivatedByRecord(t *testing.T) {
	store := routeProposalStore{Proposals: []RouteProposal{
		{CandidateID: "revived", OccurrenceCount: 2, Status: proposalStatusStale, LastSeen: daysAgo(70)},
	}}
	recordProposalOccurrence(&store, "revived", []string{"z"}, nil, "workflow", rpNow)
	if store.Proposals[0].Status != proposalStatusAccumulating {
		t.Fatalf("a fresh hit must re-activate a stale proposal, got %s", store.Proposals[0].Status)
	}
	if store.Proposals[0].OccurrenceCount != 3 {
		t.Fatalf("count should bump to 3, got %d", store.Proposals[0].OccurrenceCount)
	}
}

func TestProposalStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/p.yaml"
	store := routeProposalStore{SchemaVersion: 1, Proposals: []RouteProposal{
		{CandidateID: "rt", FirstSeen: daysAgo(3), LastSeen: daysAgo(1), OccurrenceCount: 2, Status: proposalStatusAccumulating, SuggestedUserSignals: []string{"a"}},
	}}
	if err := saveProposalStore(path, store); err != nil {
		t.Fatal(err)
	}
	got, err := loadProposalStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Proposals) != 1 || got.Proposals[0].CandidateID != "rt" || got.Proposals[0].OccurrenceCount != 2 {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
}

func TestLoadProposalStore_MissingIsEmpty(t *testing.T) {
	got, err := loadProposalStore(t.TempDir() + "/does-not-exist.yaml")
	if err != nil {
		t.Fatalf("missing file must not error: %v", err)
	}
	if len(got.Proposals) != 0 {
		t.Fatalf("expected empty store, got %+v", got)
	}
}
