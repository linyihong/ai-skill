package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type aiSkillProjectMetadataFile struct {
	Project aiSkillProjectMetadata `yaml:"project"`
}

type aiSkillProjectMetadata struct {
	ID              string   `yaml:"id"`
	Visibility      string   `yaml:"visibility"`
	PrivateTokens   []string `yaml:"private_tokens"`
	PrivateEntities []string `yaml:"private_entities"`
}

type repositoryTopologyRow struct {
	Subtree string `json:"subtree"`
	Shared  bool   `json:"shared"`
}

type derivedForbiddenToken struct {
	Token                string
	CanonicalToken       string
	OwningProjectID      string
	SourceMetadataPath   string
	SuggestedPlaceholder string
}

const repositoryTopologyRuntimeTargetKey = "runtime.repository_topology.config"

func compileDerivedForbiddenTokens(repo string, db *sql.DB) error {
	metadataFiles, err := discoverProjectMetadataFiles(repo)
	if err != nil {
		return err
	}
	for _, rel := range metadataFiles {
		metadata, err := readProjectMetadata(filepath.Join(repo, filepath.FromSlash(rel)))
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}
		if strings.ToLower(strings.TrimSpace(metadata.Project.Visibility)) != "private" {
			continue
		}
		projectID := strings.TrimSpace(metadata.Project.ID)
		if projectID == "" {
			projectID = strings.TrimSuffix(filepath.Base(filepath.Dir(rel)), string(filepath.Separator))
		}
		placeholder := projectPlaceholder(projectID)
		for _, canonical := range append(metadata.Project.PrivateTokens, metadata.Project.PrivateEntities...) {
			canonical = strings.TrimSpace(canonical)
			if canonical == "" {
				continue
			}
			for _, token := range sanitizationTokenVariants(canonical) {
				if _, err := db.Exec(`INSERT OR REPLACE INTO derived_forbidden_tokens (token, canonical_token, owning_project_id, source_metadata_path, suggested_placeholder) VALUES (?, ?, ?, ?, ?)`,
					token, canonical, projectID, rel, placeholder,
				); err != nil {
					return err
				}
			}
		}
		if err := insertProjectMetadataSourceFile(db, rel); err != nil {
			return err
		}
	}
	return nil
}

func insertProjectMetadataSourceFile(db *sql.DB, rel string) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO runtime_source_files (source_path, source_kind, target_table, compile_rule, compiled_at, compiler_version, status) VALUES (?, 'project_metadata', 'derived_forbidden_tokens', 'project_metadata_private_tokens', datetime('now'), ?, 'synced')`, rel, goRuntimeCompilerVersion)
	return err
}

func validateSanitizationStagedContent(root string, staged []string) string {
	if root == "" || len(staged) == 0 {
		return ""
	}
	topology, tokens, err := loadSanitizationRuntimeData(filepath.Join(root, "runtime", "runtime.db"))
	if err != nil || len(topology) == 0 || len(tokens) == 0 {
		return ""
	}

	resolver := newStagedBlobResolver(root)
	var findings []string
	stagedSorted := append([]string(nil), staged...)
	sort.Strings(stagedSorted)
	for _, rel := range stagedSorted {
		rel = filepath.ToSlash(rel)
		if !sanitizationPathIsShared(rel, topology) {
			continue
		}
		content, err := resolver.Read(rel)
		if err != nil {
			continue
		}
		findings = append(findings, sanitizationFindingsForContent(rel, string(content), tokens)...)
	}
	if len(findings) == 0 {
		return ""
	}
	return "sanitization-scan:\n  forbidden private token(s) in staged shared-layer content:\n    - " +
		strings.Join(findings, "\n    - ") +
		"\n  Remediation: replace project-specific details with placeholders such as `<PROJECT_ROOT>` or the suggested placeholder; keep incident evidence in project-local docs."
}

func discoverProjectMetadataFiles(repo string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(repo, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", "vendor":
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != ".ai-skill-project.yaml" {
			return nil
		}
		rel, err := filepath.Rel(repo, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(files)
	return files, err
}

func readProjectMetadata(path string) (aiSkillProjectMetadataFile, error) {
	var metadata aiSkillProjectMetadataFile
	content, err := os.ReadFile(path)
	if err != nil {
		return metadata, err
	}
	if err := yaml.Unmarshal(content, &metadata); err != nil {
		return metadata, err
	}
	return metadata, nil
}

func loadSanitizationRuntimeData(dbPath string) (map[string]bool, []derivedForbiddenToken, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	if ok, err := generatedSurfaceTargetExists(db, repositoryTopologyRuntimeTargetKey); err == nil && !ok {
		return nil, nil, fmt.Errorf("required runtime projection missing: %s", repositoryTopologyRuntimeTargetKey)
	}
	topology, err := loadRepositoryTopology(db)
	if err != nil {
		return nil, nil, err
	}
	tokens, err := loadDerivedForbiddenTokens(db)
	if err != nil {
		return nil, nil, err
	}
	return topology, tokens, nil
}

func generatedSurfaceTargetExists(db *sql.DB, targetKey string) (bool, error) {
	var tableCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'generated_surfaces'`).Scan(&tableCount); err != nil {
		return false, err
	}
	if tableCount == 0 {
		return true, nil
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM generated_surfaces WHERE target_key = ?`, targetKey).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func loadRepositoryTopology(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT subtree, content FROM repository_topology`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]bool{}
	for rows.Next() {
		var subtree string
		var raw string
		if err := rows.Scan(&subtree, &raw); err != nil {
			return nil, err
		}
		var row repositoryTopologyRow
		_ = json.Unmarshal([]byte(raw), &row)
		if row.Subtree == "" {
			row.Subtree = subtree
		}
		normalized := filepath.ToSlash(strings.TrimSpace(row.Subtree))
		if normalized != "" && !strings.HasSuffix(normalized, "/") {
			normalized += "/"
		}
		result[normalized] = row.Shared
	}
	return result, rows.Err()
}

func loadDerivedForbiddenTokens(db *sql.DB) ([]derivedForbiddenToken, error) {
	rows, err := db.Query(`SELECT token, canonical_token, owning_project_id, source_metadata_path, suggested_placeholder FROM derived_forbidden_tokens ORDER BY token, owning_project_id, source_metadata_path`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tokens []derivedForbiddenToken
	for rows.Next() {
		var token derivedForbiddenToken
		if err := rows.Scan(&token.Token, &token.CanonicalToken, &token.OwningProjectID, &token.SourceMetadataPath, &token.SuggestedPlaceholder); err != nil {
			return nil, err
		}
		if strings.TrimSpace(token.Token) != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens, rows.Err()
}

func sanitizationPathIsShared(rel string, topology map[string]bool) bool {
	bestLen := -1
	shared := false
	for subtree, value := range topology {
		if subtree == "" {
			continue
		}
		if strings.HasPrefix(rel, subtree) && len(subtree) > bestLen {
			bestLen = len(subtree)
			shared = value
		}
	}
	return bestLen >= 0 && shared
}

func sanitizationFindingsForContent(rel string, content string, tokens []derivedForbiddenToken) []string {
	var findings []string
	lines := strings.Split(content, "\n")
	for lineIndex, line := range lines {
		for _, token := range tokens {
			if token.Token == "" || !strings.Contains(line, token.Token) {
				continue
			}
			findings = append(findings, fmt.Sprintf("%s:%d contains %q from %s; use %s", rel, lineIndex+1, token.Token, token.OwningProjectID, token.SuggestedPlaceholder))
		}
	}
	return findings
}

func sanitizationTokenVariants(token string) []string {
	parts := sanitizationTokenParts(token)
	variants := map[string]bool{token: true}
	if len(parts) > 0 {
		camel := ""
		for _, part := range parts {
			camel += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
		variants[camel] = true
		variants[strings.Join(lowerStrings(parts), "-")] = true
		variants[strings.ToUpper(strings.Join(parts, "_"))] = true
	}
	result := make([]string, 0, len(variants))
	for variant := range variants {
		variant = strings.TrimSpace(variant)
		if len([]rune(variant)) >= 3 {
			result = append(result, variant)
		}
	}
	sort.Strings(result)
	return result
}

func sanitizationTokenParts(token string) []string {
	var parts []string
	var current []rune
	flush := func() {
		if len(current) == 0 {
			return
		}
		parts = append(parts, strings.ToLower(string(current)))
		current = nil
	}
	for _, r := range token {
		if r == '-' || r == '_' || unicode.IsSpace(r) {
			flush()
			continue
		}
		if len(current) > 0 && unicode.IsUpper(r) && (unicode.IsLower(current[len(current)-1]) || unicode.IsDigit(current[len(current)-1])) {
			flush()
		}
		current = append(current, r)
	}
	flush()
	return parts
}

func lowerStrings(values []string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = strings.ToLower(value)
	}
	return result
}

func projectPlaceholder(projectID string) string {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return "<PRIVATE_PROJECT>"
	}
	normalized := regexp.MustCompile(`[^A-Za-z0-9]+`).ReplaceAllString(projectID, "_")
	normalized = strings.Trim(normalized, "_")
	if normalized == "" {
		return "<PRIVATE_PROJECT>"
	}
	return "<" + strings.ToUpper(normalized) + ">"
}
