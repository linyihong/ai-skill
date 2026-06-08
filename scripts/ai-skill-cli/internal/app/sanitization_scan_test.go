package app

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizationScanBlocksPrivateTokenInSharedLayer(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{
		"plans/":        true,
		".agent-goals/": false,
	}, []derivedForbiddenToken{
		{
			Token:                "SecretProject",
			CanonicalToken:       "SecretProject",
			OwningProjectID:      "secret-project",
			SourceMetadataPath:   ".agent-goals/demo/.ai-skill-project.yaml",
			SuggestedPlaceholder: "<SECRET_PROJECT>",
		},
	})

	rel := "plans/active/example.md"
	writeFile(t, filepath.Join(repo, rel), "# Example\n\nSecretProject incident detail\n")
	stageAll(t, repo, rel)

	got := validateSanitizationStagedContent(repo, []string{rel})
	if got == "" {
		t.Fatal("expected sanitization finding for shared-layer private token")
	}
	if !strings.Contains(got, `plans/active/example.md:3 contains "SecretProject"`) {
		t.Fatalf("unexpected finding:\n%s", got)
	}
}

func TestSanitizationScanAllowsPrivateTokenInProjectLocalLayer(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{
		"plans/":        true,
		".agent-goals/": false,
	}, []derivedForbiddenToken{
		{
			Token:                "SecretProject",
			CanonicalToken:       "SecretProject",
			OwningProjectID:      "secret-project",
			SourceMetadataPath:   ".agent-goals/demo/.ai-skill-project.yaml",
			SuggestedPlaceholder: "<SECRET_PROJECT>",
		},
	})

	rel := ".agent-goals/demo/incident.md"
	writeFile(t, filepath.Join(repo, rel), "# Incident\n\nSecretProject local evidence\n")
	stageAll(t, repo, rel)

	if got := validateSanitizationStagedContent(repo, []string{rel}); got != "" {
		t.Fatalf("project-local content should not be blocked, got:\n%s", got)
	}
}

func TestSanitizationScanBootstrapSafeWithNoPrivateTokens(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{
		"plans/": true,
	}, nil)

	rel := "plans/active/new-framework-concept.md"
	writeFile(t, filepath.Join(repo, rel), "# New Framework Concept\n\nActivationBridgeV2 is reusable terminology.\n")
	stageAll(t, repo, rel)

	if got := validateSanitizationStagedContent(repo, []string{rel}); got != "" {
		t.Fatalf("new framework concepts must pass when no project metadata declares private tokens, got:\n%s", got)
	}
}

func TestCompileDerivedForbiddenTokensProjectsPrivateMetadata(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, ".agent-goals", "demo", ".ai-skill-project.yaml"), `project:
  id: secret-project
  visibility: private
  private_tokens:
    - SecretProject
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
	if err := compileDerivedForbiddenTokens(repo, db); err != nil {
		t.Fatal(err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM derived_forbidden_tokens WHERE token IN ('SecretProject', 'secret-project', 'SECRET_PROJECT')`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3 token variants, got %d", count)
	}
}

func seedSanitizationRuntimeDB(t *testing.T, repo string, topology map[string]bool, tokens []derivedForbiddenToken) {
	t.Helper()
	dbPath := filepath.Join(repo, "runtime", "runtime.db")
	mkParent(t, dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE repository_topology (subtree TEXT PRIMARY KEY, content TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE derived_forbidden_tokens (token TEXT NOT NULL, canonical_token TEXT NOT NULL, owning_project_id TEXT NOT NULL, source_metadata_path TEXT NOT NULL, suggested_placeholder TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	for subtree, shared := range topology {
		row := map[string]any{"subtree": subtree, "shared": shared}
		if _, err := db.Exec(`INSERT INTO repository_topology (subtree, content) VALUES (?, ?)`, subtree, runtimeJSON(row)); err != nil {
			t.Fatal(err)
		}
	}
	for _, token := range tokens {
		if _, err := db.Exec(`INSERT INTO derived_forbidden_tokens (token, canonical_token, owning_project_id, source_metadata_path, suggested_placeholder) VALUES (?, ?, ?, ?, ?)`,
			token.Token, token.CanonicalToken, token.OwningProjectID, token.SourceMetadataPath, token.SuggestedPlaceholder,
		); err != nil {
			t.Fatal(err)
		}
	}
}
