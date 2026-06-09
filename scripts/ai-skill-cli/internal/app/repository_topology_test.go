package app

import (
	"path/filepath"
	"strings"
	"testing"
)

// repository_topology_test.go — Phase 1B canonical loader tests.
//
// Spec reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1B + runtime/repository-topology-migration.md.
//
// Plan-mandated cases (4):
//   1. v1 read
//   2. v2 read
//   3. v2 write round-trip
//   4. reject missing owner/purpose
//
// Regression / edge cases follow.

func TestLoadRepositoryTopology_V1Read(t *testing.T) {
	body := []byte(`
schema_version: 1
status: active
owner_layer: runtime

runtime_projection:
  enabled: true
  target_key: runtime.repository_topology.config

shared_layer_classification:
  - subtree: plans/
    shared: true
  - subtree: workflow/
    shared: true
  - subtree: scripts/
    shared: false

expected_consumers:
  - sanitization
  - workflow_activation

invariants:
  - "Topology is a source-of-truth surface, not a detector heuristic."
`)
	file, err := ParseRepositoryTopology(body)
	if err != nil {
		t.Fatalf("expected v1 schema to parse; got: %v", err)
	}
	if file.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion=1, got %d", file.SchemaVersion)
	}
	if len(file.Subtrees) != 3 {
		t.Fatalf("expected 3 subtrees, got %d", len(file.Subtrees))
	}
	plans := file.Subtrees[0]
	if plans.Path != "plans/" || !plans.SharedLayer {
		t.Errorf("unexpected first subtree: %+v", plans)
	}
	// v1 has no owner/purpose; ensure they are empty after read.
	if plans.Owner != "" || plans.Purpose != "" {
		t.Errorf("expected v1 read to leave owner/purpose empty; got owner=%q purpose=%q", plans.Owner, plans.Purpose)
	}
	scripts := file.Subtrees[2]
	if scripts.Path != "scripts/" || scripts.SharedLayer {
		t.Errorf("unexpected scripts subtree: %+v", scripts)
	}
	// v1 expected_consumers is intentionally dropped on normalize (frozen
	// governance decision: consumer list is derived, not maintained).
	if file.ConsumerTracking != nil {
		t.Errorf("expected v1 read to leave ConsumerTracking nil; got %+v", file.ConsumerTracking)
	}
	if file.RuntimeProjection == nil || !file.RuntimeProjection.Enabled {
		t.Errorf("expected runtime_projection.enabled to be carried through")
	}
}

func TestLoadRepositoryTopology_V2Read(t *testing.T) {
	body := []byte(`
schema_version: 2
status: active
owner_layer: runtime

runtime_projection:
  enabled: true
  target_key: runtime.repository_topology.config

consumer_tracking:
  strategy: code_reference
  rationale: |
    Manual consumer lists go stale. This block is intentionally permanent.

subtrees:
  - path: plans/
    shared_layer: true
    owner: framework-maintainer
    purpose: "Plan tracking; referenced by enforcement-registry child_plan"
  - path: scripts/
    shared_layer: false
    owner: tooling-maintainer
    purpose: "CLI / runtime implementation; not consumed as reusable knowledge"
`)
	file, err := ParseRepositoryTopology(body)
	if err != nil {
		t.Fatalf("expected v2 schema to parse; got: %v", err)
	}
	if file.SchemaVersion != 2 {
		t.Errorf("expected SchemaVersion=2, got %d", file.SchemaVersion)
	}
	if len(file.Subtrees) != 2 {
		t.Fatalf("expected 2 subtrees, got %d", len(file.Subtrees))
	}
	plans := file.Subtrees[0]
	if plans.Path != "plans/" || plans.Owner != "framework-maintainer" || plans.Purpose == "" {
		t.Errorf("unexpected plans subtree: %+v", plans)
	}
	if file.ConsumerTracking == nil || file.ConsumerTracking.Strategy != "code_reference" {
		t.Errorf("expected ConsumerTracking.Strategy=code_reference; got %+v", file.ConsumerTracking)
	}
}

func TestWriteRepositoryTopology_V2RoundTrip(t *testing.T) {
	original := &RepositoryTopologyFile{
		SchemaVersion: 2,
		Status:        "active",
		OwnerLayer:    "runtime",
		RuntimeProjection: &RuntimeProjectionConfig{
			Enabled:   true,
			TargetKey: "runtime.repository_topology.config",
		},
		ConsumerTracking: &ConsumerTracking{
			Strategy:  "code_reference",
			Rationale: "Manual consumer lists go stale.",
		},
		Subtrees: []Subtree{
			{Path: "workflow/", SharedLayer: true, Owner: "framework-maintainer", Purpose: "Cross-skill workflow"},
			{Path: "plans/", SharedLayer: true, Owner: "framework-maintainer", Purpose: "Plan tracking"},
			{Path: "scripts/", SharedLayer: false, Owner: "tooling-maintainer", Purpose: "CLI implementation"},
		},
		Invariants: []string{"Topology is a source-of-truth surface."},
	}

	path := filepath.Join(t.TempDir(), "repository-topology.yaml")
	if err := WriteRepositoryTopology(path, original); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	roundtrip, err := LoadRepositoryTopology(path)
	if err != nil {
		t.Fatalf("read-back failed: %v", err)
	}
	if roundtrip.SchemaVersion != 2 {
		t.Errorf("expected round-tripped SchemaVersion=2, got %d", roundtrip.SchemaVersion)
	}
	if len(roundtrip.Subtrees) != 3 {
		t.Fatalf("expected 3 subtrees after round-trip, got %d", len(roundtrip.Subtrees))
	}
	// Writer sorts subtrees by path for deterministic output. Verify the
	// sorted order is observable in the round-tripped file.
	wantPaths := []string{"plans/", "scripts/", "workflow/"}
	for i, want := range wantPaths {
		if roundtrip.Subtrees[i].Path != want {
			t.Errorf("subtree[%d]: expected sorted path %q, got %q", i, want, roundtrip.Subtrees[i].Path)
		}
	}
	// Verify a full subtree to ensure owner/purpose survive.
	plans := roundtrip.Subtrees[0]
	if plans.Owner != "framework-maintainer" || plans.Purpose != "Plan tracking" {
		t.Errorf("plans subtree round-trip lost fields: %+v", plans)
	}
	if roundtrip.ConsumerTracking == nil || roundtrip.ConsumerTracking.Strategy != "code_reference" {
		t.Errorf("ConsumerTracking lost in round-trip: %+v", roundtrip.ConsumerTracking)
	}
}

func TestWriteRepositoryTopology_RejectsMissingOwner(t *testing.T) {
	file := &RepositoryTopologyFile{
		SchemaVersion: 2,
		Subtrees: []Subtree{
			{Path: "plans/", SharedLayer: true, Owner: "", Purpose: "x"},
		},
	}
	path := filepath.Join(t.TempDir(), "repository-topology.yaml")
	err := WriteRepositoryTopology(path, file)
	if err == nil {
		t.Fatal("expected write to fail when owner is missing")
	}
	if !IsRepositoryTopologyValidationError(err) {
		t.Errorf("expected validation error type; got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "v2.subtrees.owner.required") {
		t.Errorf("expected owner.required rule citation; got: %s", err.Error())
	}
}

func TestWriteRepositoryTopology_RejectsMissingPurpose(t *testing.T) {
	file := &RepositoryTopologyFile{
		SchemaVersion: 2,
		Subtrees: []Subtree{
			{Path: "plans/", SharedLayer: true, Owner: "x", Purpose: ""},
		},
	}
	path := filepath.Join(t.TempDir(), "repository-topology.yaml")
	err := WriteRepositoryTopology(path, file)
	if err == nil {
		t.Fatal("expected write to fail when purpose is missing")
	}
	if !strings.Contains(err.Error(), "v2.subtrees.purpose.required") {
		t.Errorf("expected purpose.required rule citation; got: %s", err.Error())
	}
}

func TestWriteRepositoryTopology_RejectsMissingPath(t *testing.T) {
	file := &RepositoryTopologyFile{
		SchemaVersion: 2,
		Subtrees: []Subtree{
			{Path: "", SharedLayer: true, Owner: "x", Purpose: "y"},
		},
	}
	path := filepath.Join(t.TempDir(), "repository-topology.yaml")
	err := WriteRepositoryTopology(path, file)
	if err == nil {
		t.Fatal("expected write to fail when path is missing")
	}
	if !strings.Contains(err.Error(), "v2.subtrees.path.required") {
		t.Errorf("expected path.required rule citation; got: %s", err.Error())
	}
}

func TestWriteRepositoryTopology_RejectsV1ReadDirectly(t *testing.T) {
	// A file loaded from v1 has empty owner/purpose; writing it as-is
	// must fail. Phase 1B explicitly does NOT auto-fabricate values
	// during migration — callers must populate before writing.
	body := []byte(`
schema_version: 1
shared_layer_classification:
  - subtree: plans/
    shared: true
`)
	file, err := ParseRepositoryTopology(body)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}
	path := filepath.Join(t.TempDir(), "repository-topology.yaml")
	if err := WriteRepositoryTopology(path, file); err == nil {
		t.Fatal("expected write to fail because v1 read leaves owner/purpose empty")
	}
}

func TestLoadRepositoryTopology_LiveFileParses(t *testing.T) {
	// Regression guard: the live runtime/repository-topology.yaml must
	// continue to parse via the canonical loader. After Phase 1C₁
	// (2026-06-09) the live file was upgraded from v1 to v2; this test
	// now asserts v2. The test was renamed from LiveV1FileParses to
	// LiveFileParses so future schema bumps can update the version
	// assertion without renaming the test.
	repo := discoveryRepoRoot(t)
	path := filepath.Join(repo, "runtime", "repository-topology.yaml")
	file, err := LoadRepositoryTopology(path)
	if err != nil {
		t.Fatalf("live file must parse cleanly; got: %v", err)
	}
	if file.SchemaVersion != 2 {
		t.Errorf("expected live file to be v2 (Phase 1C₁ migrated it); got SchemaVersion=%d", file.SchemaVersion)
	}
	if len(file.Subtrees) == 0 {
		t.Error("expected at least one subtree in live topology")
	}
	// v2 contract: every subtree must declare owner + purpose.
	for _, s := range file.Subtrees {
		if s.Owner == "" || s.Purpose == "" {
			t.Errorf("live v2 subtree %q missing owner/purpose: %+v", s.Path, s)
		}
	}
	// v2 contract: consumer_tracking must be present and frozen.
	if file.ConsumerTracking == nil {
		t.Error("expected ConsumerTracking on live v2 file")
	} else if file.ConsumerTracking.Strategy != "code_reference" {
		t.Errorf("expected ConsumerTracking.Strategy=code_reference; got %q", file.ConsumerTracking.Strategy)
	}
}

func TestLoadRepositoryTopology_SchemaVersionInferenceV2(t *testing.T) {
	body := []byte(`
subtrees:
  - path: plans/
    shared_layer: true
    owner: x
    purpose: y
`)
	file, err := ParseRepositoryTopology(body)
	if err != nil {
		t.Fatalf("inference parse failed: %v", err)
	}
	if file.SchemaVersion != 2 {
		t.Errorf("expected SchemaVersion=2 inferred from `subtrees:` field; got %d", file.SchemaVersion)
	}
}

func TestLoadRepositoryTopology_SchemaVersionInferenceV1(t *testing.T) {
	body := []byte(`
shared_layer_classification:
  - subtree: plans/
    shared: true
`)
	file, err := ParseRepositoryTopology(body)
	if err != nil {
		t.Fatalf("inference parse failed: %v", err)
	}
	if file.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion=1 inferred from `shared_layer_classification:` field; got %d", file.SchemaVersion)
	}
}

func TestLoadRepositoryTopology_UnknownFieldTolerance(t *testing.T) {
	// Forward-compat: future schema versions may add fields. The loader
	// captures unknowns into Unknown maps but never rejects them.
	body := []byte(`
schema_version: 2
subtrees:
  - path: plans/
    shared_layer: true
    owner: x
    purpose: y
    future_severity: high
    future_expires_at: 2099-01-01
v3_future_field:
  nested: true
`)
	if _, err := ParseRepositoryTopology(body); err != nil {
		t.Errorf("expected unknown fields to be tolerated; got: %v", err)
	}
}

func TestLoadRepositoryTopology_EmptySubtreesOk(t *testing.T) {
	// A brand-new project may declare topology with no subtrees yet.
	body := []byte(`
schema_version: 2
status: active
`)
	file, err := ParseRepositoryTopology(body)
	if err != nil {
		t.Errorf("empty subtrees should be valid on read; got: %v", err)
	}
	if file == nil || len(file.Subtrees) != 0 {
		t.Errorf("expected empty subtrees slice; got %+v", file)
	}
}

func TestLoadRepositoryTopology_RejectsUnsupportedSchemaVersion(t *testing.T) {
	body := []byte(`
schema_version: 99
subtrees: []
`)
	_, err := ParseRepositoryTopology(body)
	if err == nil {
		t.Fatal("expected unsupported schema version to surface as error")
	}
	if !strings.Contains(err.Error(), "schema_version.unsupported") {
		t.Errorf("expected schema_version.unsupported rule citation; got: %s", err.Error())
	}
}

func TestLoadRepositoryTopology_RejectsV2SubtreeMissingPath(t *testing.T) {
	body := []byte(`
schema_version: 2
subtrees:
  - owner: x
    purpose: y
    shared_layer: true
`)
	_, err := ParseRepositoryTopology(body)
	if err == nil {
		t.Fatal("expected missing path to surface as validation error")
	}
	if !IsRepositoryTopologyValidationError(err) {
		t.Errorf("expected RepositoryTopologyValidationError; got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "v2.subtrees.path.required") {
		t.Errorf("expected path.required rule citation; got: %s", err.Error())
	}
}

// TestLoadRepositoryTopology_RejectsV1SubtreeMissingPath (gap G1 from
// 2026-06-09 review): verify the v1 read path also catches a subtree
// entry with no `subtree:` field. The rule
// `v1.shared_layer_classification.subtree.required` was defined in the
// loader but had no test coverage before this commit.
func TestLoadRepositoryTopology_RejectsV1SubtreeMissingPath(t *testing.T) {
	body := []byte(`
schema_version: 1
shared_layer_classification:
  - shared: true
    # subtree: intentionally missing
`)
	_, err := ParseRepositoryTopology(body)
	if err == nil {
		t.Fatal("expected v1 subtree missing path to surface as validation error")
	}
	if !IsRepositoryTopologyValidationError(err) {
		t.Errorf("expected RepositoryTopologyValidationError; got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "v1.shared_layer_classification.subtree.required") {
		t.Errorf("expected v1 subtree.required rule citation; got: %s", err.Error())
	}
}

// TestLoadRepositoryTopology_V1SubtreeUnknownFieldTolerance (gap G2
// from 2026-06-09 review): v1 subtree entries may carry forward-compat
// fields that v1 didn't define. The rawSubtreeV1.Unknown map captures
// them; this test asserts the capture path does not reject. Locks in
// forward-compat for the v1 subtree level (top-level + v2 subtree are
// covered by TestLoadRepositoryTopology_UnknownFieldTolerance).
func TestLoadRepositoryTopology_V1SubtreeUnknownFieldTolerance(t *testing.T) {
	body := []byte(`
schema_version: 1
shared_layer_classification:
  - subtree: plans/
    shared: true
    future_owner_hint: someone   # not part of v1 schema; must be tolerated
    future_purpose_hint: tracking
`)
	if _, err := ParseRepositoryTopology(body); err != nil {
		t.Errorf("expected v1 subtree unknown fields to be tolerated; got: %v", err)
	}
}

// TestLoadRepositoryTopology_IOErrorDistinctFromValidation (gap G5
// from 2026-06-09 review; sibling of project_metadata_test.go's
// TestLoadProjectMetadata_IOErrorDistinctFromValidation): a file-not-
// found error must surface as a plain I/O error, NOT as a
// RepositoryTopologyValidationError. Callers branch on the type to
// distinguish "file is missing" from "file is malformed schema".
func TestLoadRepositoryTopology_IOErrorDistinctFromValidation(t *testing.T) {
	_, err := LoadRepositoryTopology(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("expected I/O error for missing file")
	}
	if IsRepositoryTopologyValidationError(err) {
		t.Errorf("expected I/O error not to surface as ValidationError, got: %v", err)
	}
}

// TestLoadRepositoryTopology_ExplicitV1IgnoresV2Fields (gap G4 from
// 2026-06-09 review): when a YAML declares `schema_version: 1` but
// also contains v2-only fields (`subtrees:`, `consumer_tracking:`),
// the v2 fields are SILENTLY DROPPED. This is intentional — the author
// signalled v1 intent via the explicit version. Partial-migration files
// (where both shapes co-exist) would otherwise produce undefined
// behaviour.
//
// This test locks the behaviour so a future "convenience" refactor
// (e.g. "auto-upgrade when v2 fields are present") cannot land
// silently. Such a refactor would require an explicit schema_version
// bump in this test or a v3 schema design.
//
// See runtime/repository-topology-migration.md §Schema version
// precedence (mixed-shape files) for the prose rationale and the
// Phase 1C migration discipline this protects.
func TestLoadRepositoryTopology_ExplicitV1IgnoresV2Fields(t *testing.T) {
	body := []byte(`
schema_version: 1

# v1 fields — these are honoured
shared_layer_classification:
  - subtree: plans/
    shared: true

# v2 fields — these MUST be silently ignored
subtrees:
  - path: workflow/
    shared_layer: true
    owner: ghost-owner
    purpose: "should not appear in normalized output"
consumer_tracking:
  strategy: code_reference
  rationale: "should not appear in normalized output"
`)
	file, err := ParseRepositoryTopology(body)
	if err != nil {
		t.Fatalf("mixed-shape v1 file should parse; got: %v", err)
	}
	if file.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion=1 (explicit takes precedence); got %d", file.SchemaVersion)
	}
	// v1 subtree honoured.
	if len(file.Subtrees) != 1 {
		t.Fatalf("expected exactly 1 subtree (from v1 path); got %d: %+v", len(file.Subtrees), file.Subtrees)
	}
	if file.Subtrees[0].Path != "plans/" {
		t.Errorf("expected v1 subtree plans/; got %q", file.Subtrees[0].Path)
	}
	// v2 subtree silently dropped — workflow/ must NOT appear.
	for _, s := range file.Subtrees {
		if s.Path == "workflow/" {
			t.Errorf("v2 subtree leaked into v1 normalized output: %+v", s)
		}
		if s.Owner == "ghost-owner" {
			t.Errorf("v2 owner field leaked into v1 normalized output")
		}
	}
	// v2 consumer_tracking silently dropped.
	if file.ConsumerTracking != nil {
		t.Errorf("v2 consumer_tracking leaked into v1 normalized output: %+v", file.ConsumerTracking)
	}
}
