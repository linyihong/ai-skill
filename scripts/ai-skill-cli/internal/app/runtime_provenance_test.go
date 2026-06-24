package app

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

// provenanceRepo writes a minimal routing registry declaring route.feedback.history
// so feedbackCanonicalSink derives the sink ("feedback/history/") from the
// registry rather than a hardcoded literal.
func provenanceRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"),
		"records:\n"+
			"  - id: route.feedback.history\n"+
			"    route_type: feedback\n"+
			"    primary_source: feedback/history/README.md\n")
	return repo
}

// provenanceIndex builds a minimal SQLite index whose atoms table holds the given
// (type, source_path) rows, and returns its path.
func provenanceIndex(t *testing.T, repo string, rows [][2]string) string {
	t.Helper()
	path := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE atoms (id TEXT PRIMARY KEY, source_path TEXT, type TEXT)`); err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		if _, err := db.Exec(`INSERT INTO atoms (id, source_path, type) VALUES (?, ?, ?)`,
			"atom-"+string(rune('a'+i)), row[1], row[0]); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

// TestNativeRuntimeIndexProvenanceCheck locks exactly four cases (ADR-004
// Completion Plan Phase 0/P0-A): provenance, not token; route presence must not
// count as lesson presence.
func TestNativeRuntimeIndexProvenanceCheck(t *testing.T) {
	cases := []struct {
		name       string
		rows       [][2]string // (type, source_path)
		wantStatus string
	}{
		{
			// no provenance atoms -> check fails (refresh continues; wiring is non-gating).
			name:       "no_provenance_atoms",
			rows:       [][2]string{{"route", "feedback/history/README.md"}, {"reference", "runtime/runtime.db"}},
			wantStatus: "failed",
		},
		{
			// route atom only -> check fails (anti-mask: route token must not pass).
			name:       "route_atom_only",
			rows:       [][2]string{{"route", "feedback/history/README.md"}},
			wantStatus: "failed",
		},
		{
			// valid lesson atom -> check passes.
			name:       "valid_lesson_atom",
			rows:       [][2]string{{"feedback-pattern", "feedback/history/apk-analysis/common/lesson.md"}},
			wantStatus: "ok",
		},
		{
			// sink mismatch (old-world path) -> check fails (drift resistance).
			name:       "sink_mismatch",
			rows:       [][2]string{{"feedback-pattern", "skills/demo/feedback_history/lesson.md"}},
			wantStatus: "failed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := provenanceRepo(t)
			path := provenanceIndex(t, repo, tc.rows)
			check := nativeRuntimeIndexProvenanceCheck(repo, path)
			if check.Name != "runtime_index_feedback_provenance" {
				t.Fatalf("unexpected check name %q", check.Name)
			}
			if check.Status != tc.wantStatus {
				t.Fatalf("case %s: want status %q, got %q (msg=%s)", tc.name, tc.wantStatus, check.Status, check.Message)
			}
		})
	}
}
