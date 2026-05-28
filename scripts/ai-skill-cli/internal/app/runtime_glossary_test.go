package app

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const glossaryProjectionFixture = "## context_mode\n\n" +
	"```yaml\n" +
	"term: context_mode\n" +
	"status: canonical\n" +
	"owner-layer: runtime-cognition\n" +
	"meaning: Test entry for context_mode.\n" +
	"affects:\n" +
	"  - runtime/cognitive-modes.yaml\n" +
	"  - workflow/foo.md\n" +
	"  - plans/active/x.md\n" +
	"aliases:\n" +
	"  - ctx_mode\n" +
	"related-terms:\n" +
	"  - { type: related_to, target: execution_mode }\n" +
	"introduced-by: plans/active/2026-05-25-1000-context-language-glossary-system.md\n" +
	"```\n\n" +
	"## execution_mode\n\n" +
	"```yaml\n" +
	"term: execution_mode\n" +
	"status: canonical\n" +
	"owner-layer: runtime-cognition\n" +
	"meaning: Test entry for execution_mode.\n" +
	"affects:\n" +
	"  - runtime/cognitive-modes.yaml\n" +
	"related-terms:\n" +
	"  - { type: related_to, target: context_mode }\n" +
	"```\n"

func setupGlossaryProjection(t *testing.T) *sql.DB {
	t.Helper()
	repo := t.TempDir()
	glossaryDir := filepath.Join(repo, "knowledge", "glossary")
	if err := os.MkdirAll(glossaryDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(glossaryDir, "ai-skill.md"), []byte(glossaryProjectionFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := populateGlossaryProjection(db, repo); err != nil {
		t.Fatalf("populate: %v", err)
	}
	return db
}

func TestPopulateGlossaryProjection_TermCount(t *testing.T) {
	db := setupGlossaryProjection(t)
	var n int
	if err := db.QueryRow("SELECT count(*) FROM glossary_terms").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("glossary_terms count: got %d want 2", n)
	}
}

func TestPopulateGlossaryProjection_OwnerLookup(t *testing.T) {
	db := setupGlossaryProjection(t)
	var owner, status, canonical, aliases string
	if err := db.QueryRow(
		"SELECT owner_layer, status, canonical_source, aliases FROM glossary_terms WHERE term=?",
		"context_mode",
	).Scan(&owner, &status, &canonical, &aliases); err != nil {
		t.Fatal(err)
	}
	if owner != "runtime-cognition" {
		t.Errorf("owner: got %q want runtime-cognition", owner)
	}
	if status != "canonical" {
		t.Errorf("status: got %q want canonical", status)
	}
	if !strings.HasSuffix(canonical, "knowledge/glossary/ai-skill.md") {
		t.Errorf("canonical_source: got %q want suffix knowledge/glossary/ai-skill.md", canonical)
	}
	if !strings.Contains(aliases, "ctx_mode") {
		t.Errorf("aliases: got %q want to contain ctx_mode", aliases)
	}
}

func TestPopulateGlossaryProjection_ReverseRelationQuery(t *testing.T) {
	db := setupGlossaryProjection(t)
	// 反向查找：誰 related_to context_mode？
	var src string
	if err := db.QueryRow(
		"SELECT source_term FROM glossary_relations WHERE target_term=? AND relation_type=?",
		"context_mode", "related_to",
	).Scan(&src); err != nil {
		t.Fatal(err)
	}
	if src != "execution_mode" {
		t.Errorf("reverse relation lookup: got %q want execution_mode", src)
	}
}

func TestPopulateGlossaryProjection_UsageClassification(t *testing.T) {
	db := setupGlossaryProjection(t)
	rows, err := db.Query(
		"SELECT source_file, source_type FROM glossary_usage WHERE term=? ORDER BY source_file",
		"context_mode",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var file, srcType string
		if err := rows.Scan(&file, &srcType); err != nil {
			t.Fatal(err)
		}
		got[file] = srcType
	}
	want := map[string]string{
		"plans/active/x.md":            "plan",
		"runtime/cognitive-modes.yaml": "runtime",
		"workflow/foo.md":              "workflow",
	}
	for f, w := range want {
		if got[f] != w {
			t.Errorf("usage %q: got source_type %q want %q", f, got[f], w)
		}
	}
}

func TestPopulateGlossaryProjection_DriftQuery_ConflictingOwner(t *testing.T) {
	db := setupGlossaryProjection(t)
	// near-duplicate fork drift: same owner_layer with multiple canonical entries
	var owner string
	var n int
	if err := db.QueryRow(
		"SELECT owner_layer, count(*) FROM glossary_terms WHERE status='canonical' GROUP BY owner_layer ORDER BY count(*) DESC LIMIT 1",
	).Scan(&owner, &n); err != nil {
		t.Fatal(err)
	}
	if owner != "runtime-cognition" || n != 2 {
		t.Errorf("drift aggregate: got owner=%q n=%d want runtime-cognition / 2", owner, n)
	}
}

func TestClassifyGlossaryUsageSourceType(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"workflow/foo/bar.md", "workflow"},
		{"validation/scenarios/baz.yaml", "validation"},
		{"runtime/runtime.db", "runtime"},
		{"knowledge/glossary/ai-skill.md", "knowledge"},
		{"constitution/ADR-007.md", "adr"},
		{"plans/active/foo.md", "plan"},
		{"memory/project/README.md", "memory"},
		{"unknown/path.md", "knowledge"},
		{"./workflow/foo.md", "workflow"},
	}
	for _, tc := range cases {
		if got := classifyGlossaryUsageSourceType(tc.path); got != tc.want {
			t.Errorf("classify %q: got %q want %q", tc.path, got, tc.want)
		}
	}
}
