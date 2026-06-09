package app

import (
	"database/sql"
	"fmt"
	"path/filepath"
)

// repository_topology_compile.go — Phase 1C₁ projection rule for
// runtime/repository-topology.yaml.
//
// This file is the **single point of wiring** between the Phase 1B canonical
// loader (`LoadRepositoryTopology`) and the `runtime.db.repository_topology`
// projected table consumed by the legacy scanner.
// Plan reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1C.
//
// Phase 1B established the loader; Phase 1C₁ (this commit) wires it into
// the compile pipeline. Phase 1D will then migrate the legacy scanner to
// read v2 fields directly. Until 1D lands, this function writes JSON
// content with BOTH v1 keys (`subtree`, `shared`) AND v2 keys (`path`,
// `shared_layer`, `owner`, `purpose`) so the legacy
// `sanitization_scan.go::repositoryTopologyRow` JSON unmarshalling
// continues to work unchanged.
//
// SCOPE DISCIPLINE — Phase 1C₁ touches:
//   - runtime/repository-topology.yaml (in-place v1 → v2 upgrade)
//   - runtime_compiler.go (removes line 339 tuple, adds call to the
//     function in this file)
//   - this file (new)
//   - repository_topology_test.go (Live*FileParses v1 → v2)
// Phase 1C₁ does NOT touch sanitization_scan.go — backward-compat JSON
// content is the contract that lets the legacy scanner continue working.

// compileRepositoryTopology reads runtime/repository-topology.yaml via the
// Phase 1B canonical loader and projects the entries into the
// runtime.db.repository_topology table.
//
// Each subtree becomes one row keyed by the path; the content column
// holds a JSON object containing BOTH v1-compatible keys and v2 keys so
// the projection is backward-compatible with the legacy scanner during
// the Phase 1D transition window.
//
// A `__config__` row holds the full canonical document for runtime
// introspection (mirrors the convention used by other tuple-driven
// projections).
func compileRepositoryTopology(repo string, db *sql.DB) error {
	path := filepath.Join(repo, "runtime", "repository-topology.yaml")
	file, err := LoadRepositoryTopology(path)
	if err != nil {
		return fmt.Errorf("compile repository topology: %w", err)
	}

	for _, s := range file.Subtrees {
		content := repositoryTopologyRowContent(s)
		if _, err := db.Exec(
			`INSERT OR REPLACE INTO repository_topology (subtree, content) VALUES (?, ?)`,
			s.Path, runtimeJSON(content),
		); err != nil {
			return fmt.Errorf("insert repository_topology row %q: %w", s.Path, err)
		}
	}

	configContent := map[string]any{
		"schema_version":    file.SchemaVersion,
		"status":            file.Status,
		"owner_layer":       file.OwnerLayer,
		"consumer_tracking": file.ConsumerTracking,
		"subtrees":          file.Subtrees,
		"invariants":        file.Invariants,
	}
	if _, err := db.Exec(
		`INSERT OR REPLACE INTO repository_topology (subtree, content) VALUES ('__config__', ?)`,
		runtimeJSON(configContent),
	); err != nil {
		return fmt.Errorf("insert repository_topology __config__: %w", err)
	}

	// Register the source-file entry + runtime_config_projections row.
	// The assertRuntimeCanonicalDocumentsProjected invariant
	// (runtime_test.go) requires every path in
	// runtimeCanonicalDocumentPaths() to have BOTH a runtime_source_files
	// row AND a runtime_config_projections row keyed by document
	// checksum. The insertRuntimeSourceFile helper handles both, but it
	// requires the runtime_config_documents row to already exist (so
	// it can read the checksum). In real compile this is guaranteed by
	// the earlier insertRuntimeConfigDocuments step; unit tests that
	// invoke this function in isolation must seed the document row
	// themselves (see TestCompileRepositoryTopology_* helpers).
	if err := insertRuntimeSourceFile(db, "runtime/repository-topology.yaml", "repository_topology", "repository_topology_v2"); err != nil {
		return fmt.Errorf("register runtime_source_files for repository_topology: %w", err)
	}
	return nil
}

// repositoryTopologyRowContent builds the JSON-shaped content payload for
// a single repository_topology row. The shape carries BOTH v1 keys
// (`subtree`, `shared`) AND v2 keys (`path`, `shared_layer`, `owner`,
// `purpose`) so the legacy scanner's JSON struct
// (`sanitization_scan.go::repositoryTopologyRow{Subtree, Shared}`)
// continues to populate correctly while v2 consumers can read the
// governance fields directly.
//
// Exported (lower-case but referenced in tests within this package).
// The doc comment captures the contract: removing either of the v1 keys
// requires Phase 1D's legacy-scanner retirement first.
func repositoryTopologyRowContent(s Subtree) map[string]any {
	return map[string]any{
		// v1-compatible keys (kept for sanitization_scan.go until Phase 1D)
		"subtree": s.Path,
		"shared":  s.SharedLayer,
		// v2 canonical keys
		"path":         s.Path,
		"shared_layer": s.SharedLayer,
		"owner":        s.Owner,
		"purpose":      s.Purpose,
	}
}
