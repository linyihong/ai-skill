package planvalidate

import "sort"

// Phase 2.3 (Gate C): shadow divergence accounting.
//
// Compare contrasts the legacy validators' findings with the engine's findings
// so a shadow run can record parity WITHOUT affecting any exit code (Gate C.1 is
// the consumer's responsibility — Compare is pure and side-effect-free).
//
// Equality is by RuleID + Blocking, never by Message (Gate C.2): message wording
// is allowed to differ between legacy and engine without counting as divergence.
//
// Divergence is bucketed (Gate C.3) rather than a flat "mismatch", because some
// differences are EXPECTED and benign:
//   - Transport: the engine is policy-free, so it emits a finding the legacy
//     path suppressed via opt-out (commit message / config / flag). Driven by
//     CompareHints.OptedOut.
//   - Context: a difference attributable to execution context (staged-blob vs
//     working-tree, ExecutionMode). Driven by CompareHints.ContextSensitive.
// Only Missing / Extra are genuine, unexplained divergences that must converge
// before any later phase switches the hook over to the engine.

// Divergence buckets the comparison of two finding sets, keyed by RuleID.
type Divergence struct {
	Same      []string // fired in both (RuleID+Blocking match)
	Missing   []string // legacy fired, engine did not — genuine gap
	Extra     []string // engine fired, legacy did not — genuine gap
	Transport []string // explained by opt-out policy (consumer-resolved)
	Context   []string // explained by execution context (staged/worktree, mode)
}

// CompareHints lets the consumer (which knows policy + execution context, things
// the engine deliberately does not) reclassify otherwise-genuine gaps into the
// benign Transport / Context buckets.
type CompareHints struct {
	OptedOut         map[string]bool // RuleID -> opted out for this run
	ContextSensitive map[string]bool // RuleID -> differs due to execution context
}

// key identifies a finding for equality purposes: RuleID + Blocking (Gate C.2).
type cmpKey struct {
	ruleID   string
	blocking bool
}

func presence(fs []Finding) map[cmpKey]bool {
	m := map[cmpKey]bool{}
	for _, f := range fs {
		m[cmpKey{f.RuleID, f.Blocking}] = true
	}
	return m
}

// Compare buckets legacy vs engine findings. Pure: no IO, no exit code.
func Compare(legacy, engine []Finding, hints CompareHints) Divergence {
	l := presence(legacy)
	e := presence(engine)

	union := map[cmpKey]bool{}
	for k := range l {
		union[k] = true
	}
	for k := range e {
		union[k] = true
	}

	var d Divergence
	for k := range union {
		inL, inE := l[k], e[k]
		switch {
		case inL && inE:
			d.Same = append(d.Same, k.ruleID)
		case hints.ContextSensitive[k.ruleID]:
			d.Context = append(d.Context, k.ruleID)
		case hints.OptedOut[k.ruleID]:
			d.Transport = append(d.Transport, k.ruleID)
		case inL && !inE:
			d.Missing = append(d.Missing, k.ruleID)
		default: // inE && !inL
			d.Extra = append(d.Extra, k.ruleID)
		}
	}
	sort.Strings(d.Same)
	sort.Strings(d.Missing)
	sort.Strings(d.Extra)
	sort.Strings(d.Transport)
	sort.Strings(d.Context)
	return d
}

// Converged reports whether there are no genuine (unexplained) gaps. Transport
// and Context differences are expected and do not block convergence.
func (d Divergence) Converged() bool {
	return len(d.Missing) == 0 && len(d.Extra) == 0
}
