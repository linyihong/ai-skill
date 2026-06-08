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

func TestSanitizationScanBlocksGenericPatternsInSharedLayer(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{name: "email", content: "contact test@example.com\n", want: "email/email_address"},
		{name: "phone", content: "call +1 555 123 4567\n", want: "phone/international_phone"},
		{name: "os path", content: "open /Users/alice/project\n", want: "os_absolute_path/macos_user_path"},
		{name: "credential", content: "token sk-abcdefghijklmnopqrstuvwxyz123456\n", want: "credential_pattern/openai_style_secret"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := initTempGitRepo(t)
			seedSanitizationRuntimeDB(t, repo, map[string]bool{"plans/": true}, nil)
			seedSanitizationPatterns(t, repo)

			rel := "plans/active/example.md"
			writeFile(t, filepath.Join(repo, rel), tc.content)
			stageAll(t, repo, rel)

			got := validateSanitizationStagedContent(repo, []string{rel})
			if !strings.Contains(got, tc.want) {
				t.Fatalf("expected %s finding, got:\n%s", tc.want, got)
			}
		})
	}
}

func TestSanitizationScanAllowsPlaceholdersAndPatternConfigSelfReference(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{
		"plans/":   true,
		"runtime/": true,
	}, nil)
	seedSanitizationPatterns(t, repo)

	placeholderRel := "plans/active/placeholders.md"
	writeFile(t, filepath.Join(repo, placeholderRel), "Use <PROJECT_ROOT> instead of a local path.\n")
	patternRel := "runtime/sanitization-patterns.yaml"
	writeFile(t, filepath.Join(repo, patternRel), "regex: '/Users/[^/\\\\s]+'\n")
	stageAll(t, repo, placeholderRel, patternRel)

	if got := validateSanitizationStagedContent(repo, []string{placeholderRel, patternRel}); got != "" {
		t.Fatalf("placeholders and pattern config self-reference should pass, got:\n%s", got)
	}
}

func TestSanitizationScanAllowsPatternDocumentationLine(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{"plans/": true}, nil)
	seedSanitizationPatterns(t, repo)

	rel := "plans/active/pattern-doc.md"
	writeFile(t, filepath.Join(repo, rel), "- OS absolute path pattern: `/Users/[^/\\s]+` / `/home/[^/\\s]+`\n")
	stageAll(t, repo, rel)

	if got := validateSanitizationStagedContent(repo, []string{rel}); got != "" {
		t.Fatalf("pattern documentation should not self-trigger, got:\n%s", got)
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
	if _, err := db.Exec(`CREATE TABLE sanitization_patterns (category TEXT PRIMARY KEY, content TEXT NOT NULL)`); err != nil {
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

func seedSanitizationPatterns(t *testing.T, repo string) {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(repo, "runtime", "runtime.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	families := []map[string]any{
		{
			"category": "email",
			"patterns": []map[string]any{{
				"id": "email_address", "regex": `[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`, "suggestion": "replace email",
			}},
		},
		{
			"category": "phone",
			"patterns": []map[string]any{{
				"id": "international_phone", "regex": `\+[0-9][0-9 .()-]{7,}[0-9]`, "suggestion": "replace phone",
			}},
		},
		{
			"category": "os_absolute_path",
			"patterns": []map[string]any{{
				"id": "macos_user_path", "regex": `/Users/[^/\s]+`, "suggestion": "replace path",
			}},
		},
		{
			"category": "credential_pattern",
			"patterns": []map[string]any{{
				"id": "openai_style_secret", "regex": `sk-[A-Za-z0-9_-]{20,}`, "suggestion": "remove secret",
			}},
		},
	}
	for _, family := range families {
		if _, err := db.Exec(`INSERT INTO sanitization_patterns (category, content) VALUES (?, ?)`, family["category"], runtimeJSON(family)); err != nil {
			t.Fatal(err)
		}
	}
}
