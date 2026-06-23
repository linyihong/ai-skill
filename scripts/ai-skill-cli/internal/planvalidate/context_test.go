package planvalidate

import "testing"

// Phase 2.0 / Gate A acceptance: the ValidationContext contract must be
// constructible by all three intended consumers (commit-msg hook, CI, manual
// CLI) without any engine, parsing, or schema dependency.
//
// These are construction-only assertions — there is no validation logic to
// exercise yet. If a consumer cannot express its inputs through this contract,
// Gate A fails and the engine contract is not yet mature enough to proceed.

func TestValidationContext_HookConsumerCanConstruct(t *testing.T) {
	ctx := ValidationContext{
		Root:          "/repo",
		ChangedSet:    ChangedSet{"plans/active/x.md"},
		ExecutionMode: ModeCommit,
		// In commit mode the opt-out / justification transport is the commit
		// message; the consumer (not the engine) records where it resolved it.
		Metadata: ValidationMetadata{"opt_out_source": "commit-message"},
	}
	if ctx.ExecutionMode != ModeCommit {
		t.Fatalf("hook consumer: got mode %q, want %q", ctx.ExecutionMode, ModeCommit)
	}
	if len(ctx.ChangedSet) != 1 {
		t.Fatalf("hook consumer: changed set not preserved, got %d entries", len(ctx.ChangedSet))
	}
}

func TestValidationContext_CIConsumerCanConstruct(t *testing.T) {
	ctx := ValidationContext{
		Root:          "/repo",
		ChangedSet:    ChangedSet{"plans/active/x.md"},
		ExecutionMode: ModeCI,
		// CI has no commit message, so opt-out arrives via config.
		Metadata: ValidationMetadata{"opt_out_source": "config"},
	}
	if ctx.ExecutionMode != ModeCI {
		t.Fatalf("ci consumer: got mode %q, want %q", ctx.ExecutionMode, ModeCI)
	}
}

func TestValidationContext_ManualConsumerCanConstruct(t *testing.T) {
	ctx := ValidationContext{
		Root: "/repo",
		// A manual run may scan all plans rather than a changed subset.
		ChangedSet:    nil,
		ExecutionMode: ModeManual,
		Metadata:      nil,
	}
	if ctx.ExecutionMode != ModeManual {
		t.Fatalf("manual consumer: got mode %q, want %q", ctx.ExecutionMode, ModeManual)
	}
	if ctx.ChangedSet != nil {
		t.Fatalf("manual consumer: expected nil changed set to be valid")
	}
}
