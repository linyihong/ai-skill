package app

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
)

// seedCanonicalDocumentRow inserts the minimum runtime_config_documents
// row needed by insertRuntimeSourceFile's checksum lookup. In real
// compile this row is populated by the earlier insertRuntimeConfigDocuments
// step; tests that invoke compileRepositoryTopology in isolation seed
// it manually so they don't have to spin up the full compile pipeline.
func seedCanonicalDocumentRow(t *testing.T, db *sql.DB, logicalID string) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT OR REPLACE INTO runtime_config_documents (logical_id, owner_layer, status, schema_version, content_json, checksum, updated_at) VALUES (?, 'runtime', 'active', '2.0', '{}', 'seed-checksum', datetime('now'))`,
		logicalID,
	); err != nil {
		t.Fatalf("seed runtime_config_documents row: %v", err)
	}
}

// repository_topology_compile_test.go — Phase 1C₁ projection rule tests.
//
// Spec reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1C + runtime/repository-topology-migration.md.
//
// The critical contract Phase 1C₁ guards:
//   - The new projection function writes JSON content with BOTH legacy
//     v1 keys (`subtree`, `shared`) AND v2 keys (`path`, `shared_layer`,
//     `owner`, `purpose`)
//   - The legacy scanner's loadRepositoryTopology() (in sanitization_scan.go)
//     keeps reading the table unchanged
// Phase 1D will retire the v1 keys; until then, breaking the dual-shape
// JSON contract breaks the live scanner.

func TestCompileRepositoryTopology_WritesBackwardCompatJSON(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "runtime", "repository-topology.yaml"), `
schema_version: 2
status: active
owner_layer: runtime

consumer_tracking:
  strategy: code_reference
  rationale: "test rationale"

subtrees:
  - path: plans/
    shared_layer: true
    owner: framework-maintainer
    purpose: "Plan tracking"
  - path: scripts/
    shared_layer: false
    owner: tooling-maintainer
    purpose: "CLI implementation"
`)

	dbPath := filepath.Join(repo, "runtime.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := createGoRuntimeSchema(db); err != nil {
		t.Fatal(err)
	}
	seedCanonicalDocumentRow(t, db, "runtime/repository-topology.yaml")
	if err := compileRepositoryTopology(repo, db); err != nil {
		t.Fatalf("compileRepositoryTopology: %v", err)
	}

	// Two subtrees + one __config__ row = 3 rows total.
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM repository_topology`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3 rows (2 subtrees + __config__), got %d", count)
	}

	// Inspect the plans/ row: JSON must carry BOTH v1 keys and v2 keys.
	var raw string
	if err := db.QueryRow(`SELECT content FROM repository_topology WHERE subtree = 'plans/'`).Scan(&raw); err != nil {
		t.Fatalf("query plans/ row: %v", err)
	}
	var content map[string]any
	if err := json.Unmarshal([]byte(raw), &content); err != nil {
		t.Fatalf("unmarshal content: %v", err)
	}
	// v1-compatible keys (Phase 1D will retire these)
	if got, want := content["subtree"], "plans/"; got != want {
		t.Errorf("v1 key 'subtree' = %v; want %q", got, want)
	}
	if got, want := content["shared"], true; got != want {
		t.Errorf("v1 key 'shared' = %v; want %v", got, want)
	}
	// v2 canonical keys
	if got, want := content["path"], "plans/"; got != want {
		t.Errorf("v2 key 'path' = %v; want %q", got, want)
	}
	if got, want := content["shared_layer"], true; got != want {
		t.Errorf("v2 key 'shared_layer' = %v; want %v", got, want)
	}
	if got, want := content["owner"], "framework-maintainer"; got != want {
		t.Errorf("v2 key 'owner' = %v; want %q", got, want)
	}
	if got, want := content["purpose"], "Plan tracking"; got != want {
		t.Errorf("v2 key 'purpose' = %v; want %q", got, want)
	}

	// Inspect scripts/ row for shared:false correctness.
	if err := db.QueryRow(`SELECT content FROM repository_topology WHERE subtree = 'scripts/'`).Scan(&raw); err != nil {
		t.Fatalf("query scripts/ row: %v", err)
	}
	if err := json.Unmarshal([]byte(raw), &content); err != nil {
		t.Fatalf("unmarshal scripts content: %v", err)
	}
	if got := content["shared"]; got != false {
		t.Errorf("scripts/ v1 key 'shared' = %v; want false", got)
	}
	if got := content["shared_layer"]; got != false {
		t.Errorf("scripts/ v2 key 'shared_layer' = %v; want false", got)
	}

	// __config__ row holds the full v2 document for runtime introspection.
	if err := db.QueryRow(`SELECT content FROM repository_topology WHERE subtree = '__config__'`).Scan(&raw); err != nil {
		t.Fatalf("query __config__ row: %v", err)
	}
	var configContent map[string]any
	if err := json.Unmarshal([]byte(raw), &configContent); err != nil {
		t.Fatalf("unmarshal __config__: %v", err)
	}
	if got := configContent["schema_version"]; got != float64(2) {
		t.Errorf("__config__ schema_version = %v; want 2", got)
	}
	if configContent["consumer_tracking"] == nil {
		t.Error("__config__ missing consumer_tracking block")
	}
}

// TestCompileRepositoryTopology_LegacyScannerCompat ensures that the
// legacy loadRepositoryTopology() function in sanitization_scan.go
// can still read the data Phase 1C₁'s projection writes. This is the
// hard contract that justifies the dual-shape JSON content: breaking
// it would break the live scanner's topology query path.
//
// This test is the cross-file regression guard the Phase 1A/1B
// discipline established: legacy consumer code stays untouched until
// Phase 1D, and to keep it untouched we must preserve the input shape
// it reads.
func TestCompileRepositoryTopology_LegacyScannerCompat(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "runtime", "repository-topology.yaml"), `
schema_version: 2
subtrees:
  - path: workflow/
    shared_layer: true
    owner: framework-maintainer
    purpose: "Workflows"
  - path: .agent-goals/
    shared_layer: false
    owner: project-local
    purpose: "Per-project goal ledger"
`)

	dbPath := filepath.Join(repo, "runtime.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := createGoRuntimeSchema(db); err != nil {
		t.Fatal(err)
	}
	seedCanonicalDocumentRow(t, db, "runtime/repository-topology.yaml")
	if err := compileRepositoryTopology(repo, db); err != nil {
		t.Fatalf("compileRepositoryTopology: %v", err)
	}

	// Call the legacy scanner's loader (in sanitization_scan.go). It must
	// produce the same map[path]→shared result the live scanner relies on.
	got, err := loadRepositoryTopology(db)
	if err != nil {
		t.Fatalf("loadRepositoryTopology (legacy reader): %v", err)
	}
	if shared, ok := got["workflow/"]; !ok || !shared {
		t.Errorf("legacy reader missing workflow/=true; got map=%+v", got)
	}
	if shared, ok := got[".agent-goals/"]; !ok || shared {
		t.Errorf("legacy reader missing .agent-goals/=false; got map=%+v", got)
	}
}
