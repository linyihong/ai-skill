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
	}, []derivedMatchToken{
		{
			Token:                "SecretProject",
			CanonicalToken:       "SecretProject",
			EntityName:           "Secret Project",
			Kind:                 "codename",
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
	// Finding must name the protected entity, not just the token (Phase 1D).
	if !strings.Contains(got, `entity "Secret Project"`) {
		t.Fatalf("finding should name the entity, got:\n%s", got)
	}
}

func TestSanitizationScanAllowsPrivateTokenInProjectLocalLayer(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{
		"plans/":        true,
		".agent-goals/": false,
	}, []derivedMatchToken{
		{
			Token:                "SecretProject",
			CanonicalToken:       "SecretProject",
			EntityName:           "Secret Project",
			Kind:                 "codename",
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

func TestSanitizationIncidentScoreWarnsWithoutBlocking(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{"plans/": true}, nil)
	seedSanitizationPatterns(t, repo)

	rel := "plans/active/incident.md"
	writeFile(t, filepath.Join(repo, rel), "Review `docs/20260606-private-flow.md` because user said \"please inspect this broken project flow\".\n")
	stageAll(t, repo, rel)

	if got := validateSanitizationStagedContent(repo, []string{rel}); got != "" {
		t.Fatalf("incident score must not block, got:\n%s", got)
	}
	warning := warnSanitizationIncidentScore(repo, []string{rel})
	if !strings.Contains(warning, "incident_score=") {
		t.Fatalf("expected incident-score warning, got:\n%s", warning)
	}
}

func TestSanitizationIncidentScoreIgnoresLowScoreAndArchivedPaths(t *testing.T) {
	repo := initTempGitRepo(t)
	seedSanitizationRuntimeDB(t, repo, map[string]bool{
		"plans/":          true,
		"plans/archived/": false,
	}, nil)
	seedSanitizationPatterns(t, repo)

	lowRel := "plans/active/route.md"
	writeFile(t, filepath.Join(repo, lowRel), "route.workflow.travel-planning is reusable framework terminology.\n")
	archivedRel := "plans/archived/incident.md"
	writeFile(t, filepath.Join(repo, archivedRel), "Review `docs/20260606-private-flow.md` because user said \"please inspect this broken project flow\".\n")
	stageAll(t, repo, lowRel, archivedRel)

	if got := warnSanitizationIncidentScore(repo, []string{lowRel, archivedRel}); got != "" {
		t.Fatalf("low-score and archived content should not warn, got:\n%s", got)
	}
}

func seedSanitizationRuntimeDB(t *testing.T, repo string, topology map[string]bool, tokens []derivedMatchToken) {
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
	if _, err := db.Exec(`CREATE TABLE derived_match_tokens (token TEXT NOT NULL, canonical_token TEXT NOT NULL, entity_name TEXT NOT NULL, kind TEXT NOT NULL, owning_project_id TEXT NOT NULL, source_metadata_path TEXT NOT NULL, suggested_placeholder TEXT NOT NULL)`); err != nil {
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
		if _, err := db.Exec(`INSERT INTO derived_match_tokens (token, canonical_token, entity_name, kind, owning_project_id, source_metadata_path, suggested_placeholder) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			token.Token, token.CanonicalToken, token.EntityName, token.Kind, token.OwningProjectID, token.SourceMetadataPath, token.SuggestedPlaceholder,
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
	config := map[string]any{
		"incident_score": map[string]any{
			"enabled":                 true,
			"warn_if_total_score_gte": 7,
			"signals": map[string]any{
				"filename_pattern": map[string]any{
					"weight": 5,
					"patterns": []map[string]any{{
						"regex": "`?docs/[0-9]{8}-[^`<\\s]+\\.md`?",
					}},
				},
				"quoted_user_text": map[string]any{"weight": 5, "min_runes": 6},
				"domain_noun_cluster": map[string]any{
					"weight": 1, "min_dash_terms": 3,
				},
			},
		},
	}
	if _, err := db.Exec(`INSERT INTO sanitization_patterns (category, content) VALUES ('__config__', ?)`, runtimeJSON(config)); err != nil {
		t.Fatal(err)
	}
}
