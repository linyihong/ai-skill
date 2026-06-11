package app

import "testing"

// failure_authority_test.go — proves the Authority Classifier honours the
// contract across MULTIPLE subject kinds (not a path-welded special-case) and
// that the Standing precedence + safe defaults hold.
// Spec: governance/lifecycle/governance-pattern-library-draft.md
//       §"Authority Classification Contract".

func TestClassifyFailureAuthority_MetadataFile(t *testing.T) {
	cases := []struct {
		name    string
		subject AuthoritySubject
		want    AuthorityLevel
	}{
		{
			name:    "shared layer has standing to block",
			subject: AuthoritySubject{Kind: SubjectMetadataFile, Path: "enforcement/x/.ai-skill-project.yaml", SharedLayer: SharedTrue},
			want:    Authoritative,
		},
		{
			name:    "explicitly non-shared subtree cannot block",
			subject: AuthoritySubject{Kind: SubjectMetadataFile, Path: ".agent-goals/demo/.ai-skill-project.yaml", SharedLayer: SharedFalse},
			want:    NonAuthoritative,
		},
		{
			name:    "project-local owner cannot block",
			subject: AuthoritySubject{Kind: SubjectMetadataFile, Path: ".agent-goals/demo/.ai-skill-project.yaml", Owner: "project-local"},
			want:    NonAuthoritative,
		},
		{
			name:    "topology miss keeps standing (fail-safe toward protection)",
			subject: AuthoritySubject{Kind: SubjectMetadataFile, Path: "somewhere/.ai-skill-project.yaml", SharedLayer: SharedUnknown},
			want:    Authoritative,
		},
		{
			name:    "demotion wins over shared signal (project-local + shared:true)",
			subject: AuthoritySubject{Kind: SubjectMetadataFile, Owner: "project-local", SharedLayer: SharedTrue},
			want:    NonAuthoritative,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ClassifyFailureAuthority(tc.subject); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestClassifyFailureAuthority_RuntimeIndexRow(t *testing.T) {
	if got := ClassifyFailureAuthority(AuthoritySubject{Kind: SubjectRuntimeIndexRow, Registered: true}); got != Authoritative {
		t.Fatalf("registered runtime-index row must be authoritative, got %s", got)
	}
	if got := ClassifyFailureAuthority(AuthoritySubject{Kind: SubjectRuntimeIndexRow, Registered: false}); got != NonAuthoritative {
		t.Fatalf("unregistered runtime-index row must not block, got %s", got)
	}
}

func TestClassifyFailureAuthority_DiscoveryProviderNeverBlocks(t *testing.T) {
	// Even with attributes that look authoritative for other kinds, a discovery
	// provider is advisory-by-contract and must never have standing to halt.
	s := AuthoritySubject{Kind: SubjectDiscoveryProvider, Registered: true, SharedLayer: SharedTrue}
	if got := ClassifyFailureAuthority(s); got != NonAuthoritative {
		t.Fatalf("discovery provider must be advisory-only, got %s", got)
	}
}

func TestClassifyFailureAuthority_UnknownKindMustEarnStanding(t *testing.T) {
	// An unmodelled kind (e.g. generated-surface) does not get to block by
	// default. Granting a new blocking kind must be an explicit code change.
	for _, kind := range []SubjectKind{SubjectGeneratedSurface, SubjectKind("brand-new-kind"), SubjectKind("")} {
		if got := ClassifyFailureAuthority(AuthoritySubject{Kind: kind, SharedLayer: SharedTrue, Registered: true}); got != NonAuthoritative {
			t.Fatalf("unmodelled kind %q must default to non-authoritative, got %s", kind, got)
		}
	}
}

func TestAuthorityLevelString(t *testing.T) {
	if Authoritative.String() != "authoritative" || NonAuthoritative.String() != "non-authoritative" {
		t.Fatalf("unexpected String() values: %s / %s", Authoritative, NonAuthoritative)
	}
}
