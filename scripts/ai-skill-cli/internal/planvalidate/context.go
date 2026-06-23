// Package planvalidate defines the reusable, consumer-agnostic contract for
// plan validation.
//
// Phase 2.0 (Gate A) intentionally ships the ValidationContext type only — no
// engine, no methods, no commit parsing, no schema handling — so the contract
// can be proven feedable from the commit-msg hook, CI, and manual CLI consumers
// before any validation logic is extracted out of hooks.go.
//
// See plans/active/2026-06-22-1009-plans-system-portability-and-delivery-integration/
// 01-external-repo-plan-system-shared-binary.md (Phase 2.0 / Gate A).
package planvalidate

// ExecutionMode identifies which consumer is driving a validation run.
//
// It exists because some plan validators are execution-context-coupled (e.g.
// staged-blob vs working-tree reads, and opt-out transport). The governance
// rule is "engine receives effective policy / consumer resolves policy source":
// each consumer sets ExecutionMode, and the engine never infers it.
type ExecutionMode string

const (
	// ModeCommit is set by the git commit-msg hook: the staged set and the
	// commit message (the opt-out / justification transport) are available.
	ModeCommit ExecutionMode = "commit"
	// ModeCI is set by a CI run: there is no commit message, so opt-out must
	// arrive via configuration rather than a message trailer.
	ModeCI ExecutionMode = "ci"
	// ModeManual is set by a manual CLI / ad-hoc invocation: opt-out arrives
	// via a flag.
	ModeManual ExecutionMode = "manual"
)

// ChangedSet is the set of repo-relative paths a run should consider changed.
//
// It is a named slice (not a bare []string) so the contract can evolve — e.g.
// gain per-path status metadata — without breaking call sites.
type ChangedSet []string

// ValidationMetadata carries evolvable, consumer-supplied context that is not
// yet a first-class field (e.g. "HEAD", working-tree state).
//
// Keeping it a string map preserves room for the still-unresolved context
// minimal-set (Open Question Q1) without locking the struct shape prematurely.
type ValidationMetadata map[string]string

// ValidationContext is the consumer-agnostic input contract for plan validation.
//
// Gate A constraint: this type carries data only. It deliberately defines no
// methods, no parsing, and no validation logic, and it omits any schema version
// — Gate B keeps frontmatter-version resolution inside the schema compatibility
// layer, so the engine consumes a normalized model and never sees a version.
type ValidationContext struct {
	// Root is the repository root the run operates against; plans/active and
	// plans/archived live under it. Cross-repo reuse is via Root, mirroring
	// `ai-skill plans tree --root`.
	Root string
	// ChangedSet is the changed/staged paths in scope for this run. A nil or
	// empty set is valid (e.g. a manual run that scans all plans).
	ChangedSet ChangedSet
	// ExecutionMode is the driving consumer (commit | ci | manual).
	ExecutionMode ExecutionMode
	// Metadata is evolvable consumer-supplied context (see Q1).
	Metadata ValidationMetadata
}
