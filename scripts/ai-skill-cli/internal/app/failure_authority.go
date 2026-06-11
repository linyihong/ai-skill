package app

// failure_authority.go — Authority Classifier: the FIRST implementation of the
// Authority Classification Contract.
//
// Spec: governance/lifecycle/governance-pattern-library-draft.md
//       §"Authority Classification Contract (the missing layer)".
//
// PRINCIPLE — Standing (Validity ≠ Authority). Many things can be wrong; not
// every wrong thing has standing to halt the pipeline. This classifier answers
// the *standing* question (may this subject block?) — deliberately BEFORE, and
// independently of, the *validity* question (is this subject correct?).
//
// SUBJECT, NOT PATH. Authority is resolved over an abstract subject, never over
// a raw path. Only the metadata-file kind is path-shaped; runtime-index-row and
// discovery-provider derive standing from entirely different signals. Welding
// authority := path-classification would silently misjudge the non-file kinds —
// which is exactly the coupling this contract exists to prevent.
//
// OBSERVATION-STAGE. The Failure Authority family is pre-gate (N<5,
// governance-pattern-library-draft.md). This file is an *executor input*: it
// contributes the family's 4th sample but does NOT promote the pattern, and no
// caller is wired yet (the Finding A executor is the next step).

// AuthorityLevel is the standing verdict for a subject.
type AuthorityLevel int

const (
	// NonAuthoritative — may emit a warning, but MUST NOT block compile/commit.
	NonAuthoritative AuthorityLevel = iota
	// Authoritative — has standing to fail the process.
	Authoritative
)

func (l AuthorityLevel) String() string {
	switch l {
	case Authoritative:
		return "authoritative"
	case NonAuthoritative:
		return "non-authoritative"
	default:
		return "unknown"
	}
}

// SubjectKind enumerates the source kinds whose standing the contract resolves.
// Each kind derives authority from kind-specific signals, not from a shared
// path heuristic.
type SubjectKind string

const (
	SubjectMetadataFile      SubjectKind = "metadata-file"      // .ai-skill-project.yaml etc.; standing from topology
	SubjectRuntimeIndexRow   SubjectKind = "runtime-index-row"  // standing from registration (sources row)
	SubjectDiscoveryProvider SubjectKind = "discovery-provider" // advisory-by-contract; never blocks
	SubjectGeneratedSurface  SubjectKind = "generated-surface"  // not yet modelled (no sample) — see default case
)

// SharedLayerState is tri-state on purpose. A topology lookup can return
// shared, not-shared, OR no-match. The distinction is a standing decision: on a
// topology MISS the safe default is to KEEP authority (fail-safe toward
// protection), never to silently demote a possibly-shared file to a warning and
// let a real leak through. A bare bool cannot express "unknown" and would
// default that case to the dangerous direction.
type SharedLayerState int

const (
	SharedUnknown SharedLayerState = iota // topology miss / not populated
	SharedTrue                            // declared shared_layer: true
	SharedFalse                           // declared shared_layer: false
)

// ownerProjectLocal is the topology owner value that denotes a project-local,
// non-authoritative subtree (e.g. .agent-goals/).
const ownerProjectLocal = "project-local"

// AuthoritySubject is the language-neutral input the contract classifies.
// Executors construct one of these from their own world (a file + topology, a
// runtime-index row, a discovery provider) and ask for a standing verdict.
type AuthoritySubject struct {
	Kind        SubjectKind
	Path        string           // when path-bound (metadata-file)
	Owner       string           // topology owner, e.g. "project-local"
	SharedLayer SharedLayerState // metadata-file: from repository-topology
	Registered  bool             // runtime-index-row: has a sources row
	Advisory    bool             // discovery-provider: advisory-by-contract
}

// ClassifyFailureAuthority is the single contract surface every executor calls.
// It returns whether the subject has standing to block the pipeline.
//
// Resolution precedence (per kind), with non-authoritative signals acting as
// demotions that win over authoritative ones:
//
//	metadata-file:
//	  owner == project-local            -> NonAuthoritative (demotion wins)
//	  shared_layer == false             -> NonAuthoritative
//	  shared_layer == true | unknown    -> Authoritative   (fail-safe on miss)
//	runtime-index-row:
//	  registered                        -> Authoritative
//	  not registered                    -> NonAuthoritative
//	discovery-provider:
//	  always                            -> NonAuthoritative (advisory-by-contract)
//	unknown / not-yet-modelled kind:
//	  always                            -> NonAuthoritative (must EARN standing;
//	                                       adding a blocking kind requires an
//	                                       explicit case here, by design)
func ClassifyFailureAuthority(s AuthoritySubject) AuthorityLevel {
	switch s.Kind {
	case SubjectMetadataFile:
		return classifyMetadataFileAuthority(s)
	case SubjectRuntimeIndexRow:
		if s.Registered {
			return Authoritative
		}
		return NonAuthoritative
	case SubjectDiscoveryProvider:
		// Advisory-by-contract: a discovery failure must never halt the system.
		return NonAuthoritative
	default:
		// Unknown or not-yet-modelled kind (e.g. generated-surface). A subject
		// the contract has not explicitly granted standing does not get to
		// block. This default is intentional and tested: extending the
		// contract with a NEW blocking kind requires adding a case above, so
		// that grant is always a conscious decision rather than an accident.
		return NonAuthoritative
	}
}

// classifyMetadataFileAuthority resolves standing for a metadata-file subject
// from its topology classification. Demotion signals (project-local owner,
// explicitly non-shared subtree) win; otherwise the file keeps standing,
// including on a topology miss (SharedUnknown) to avoid downgrading a
// possibly-shared file and letting a real leak pass as a mere warning.
func classifyMetadataFileAuthority(s AuthoritySubject) AuthorityLevel {
	if s.Owner == ownerProjectLocal {
		return NonAuthoritative
	}
	if s.SharedLayer == SharedFalse {
		return NonAuthoritative
	}
	return Authoritative
}
