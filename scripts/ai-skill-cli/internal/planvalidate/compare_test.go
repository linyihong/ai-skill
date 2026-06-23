package planvalidate

import "testing"

func f(rule string) Finding { return Finding{RuleID: rule, Blocking: true} }

// Gate C.2: equality is by RuleID+Blocking, not Message — differing messages on
// the same rule are Same, not divergence.
func TestCompare_IgnoresMessage(t *testing.T) {
	legacy := []Finding{{RuleID: "plan_tree.unique_id", Message: "legacy wording", Blocking: true}}
	engine := []Finding{{RuleID: "plan_tree.unique_id", Message: "engine wording", Blocking: true}}
	d := Compare(legacy, engine, CompareHints{})
	if len(d.Same) != 1 || !d.Converged() {
		t.Fatalf("message-only difference should be Same/converged, got %+v", d)
	}
}

// Gate C.3: genuine gaps land in Missing / Extra.
func TestCompare_GenuineGaps(t *testing.T) {
	legacy := []Finding{f("plan_tree.frontmatter")}
	engine := []Finding{f("plan_tree.parent_reference")}
	d := Compare(legacy, engine, CompareHints{})
	if len(d.Missing) != 1 || d.Missing[0] != "plan_tree.frontmatter" {
		t.Fatalf("expected frontmatter Missing, got %+v", d)
	}
	if len(d.Extra) != 1 || d.Extra[0] != "plan_tree.parent_reference" {
		t.Fatalf("expected parent_reference Extra, got %+v", d)
	}
	if d.Converged() {
		t.Fatalf("genuine gaps must not converge")
	}
}

// Gate C.3: an engine-only finding explained by opt-out is Transport, not Extra
// (engine is policy-free; legacy suppressed it via opt-out).
func TestCompare_OptOutBecomesTransport(t *testing.T) {
	engine := []Finding{f("plan_tree.frontmatter")}
	d := Compare(nil, engine, CompareHints{OptedOut: map[string]bool{"plan_tree.frontmatter": true}})
	if len(d.Transport) != 1 || len(d.Extra) != 0 {
		t.Fatalf("opt-out engine-only finding should be Transport not Extra, got %+v", d)
	}
	if !d.Converged() {
		t.Fatalf("transport difference is benign and should converge, got %+v", d)
	}
}

// Gate C.3: a difference explained by execution context is Context, not a gap.
func TestCompare_ContextBucket(t *testing.T) {
	legacy := []Finding{f("plan_tree.archive_order")}
	d := Compare(legacy, nil, CompareHints{ContextSensitive: map[string]bool{"plan_tree.archive_order": true}})
	if len(d.Context) != 1 || len(d.Missing) != 0 {
		t.Fatalf("context-sensitive difference should be Context not Missing, got %+v", d)
	}
	if !d.Converged() {
		t.Fatalf("context difference is benign and should converge, got %+v", d)
	}
}
