package app

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// project_metadata_compile.go — projection rule for
// <PROJECT_ROOT>/.ai-skill-project.yaml.
//
// Plan reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1C₂ (introduced the tables) + §Phase 1D
// (this file became the SOLE projection; legacy derived_forbidden_tokens
// retired and the scanner migrated to derived_match_tokens).
//
// Each project's `private_entities` declarations project into two tables
// mirroring the governance/execution layer split (plan Q8):
//
//   - derived_private_entities (governance layer): one row per declared
//     entity, carrying the human-readable identity (name + kind) findings
//     and audit output refer to.
//   - derived_match_tokens (execution layer): one row per case-variant of
//     every entity's match_tokens — the literal surface the scanner
//     (sanitization_scan.go::loadDerivedMatchTokens) compares staged
//     content against. Case-variant expansion happens HERE, never in the
//     parser (project_metadata.go is a pure reader).
//
// SHAPE-AWARE TOLERANCE (Phase 1D, plan §"Phase 1D — Shape-Aware Skip
// Remediation"): legacy flat-shape files (`private_tokens:` or a
// `private_entities:` list of scalars) are a DEPRECATED schema with no
// projection anymore — they are tolerated with an stderr warning and
// skipped. A file that IS new-schema shape but fails validation is a
// genuine misconfiguration and HARD-FAILS the compile, so a typo can never
// silently drop a project's protected tokens (the silent-skip gap 1C₂'s
// self-review flagged).

// derivedPrivateEntityRow is one governance-layer row.
type derivedPrivateEntityRow struct {
	EntityName           string
	Kind                 string
	OwningProjectID      string
	SourceMetadataPath   string
	SuggestedPlaceholder string
}

// derivedMatchTokenRow is one execution-layer row. CanonicalToken records
// the declared match_token this variant expanded from (or the variant
// itself for explicitly-enumerated case_variants), so governance debug can
// trace a matched token back to its declaration.
type derivedMatchTokenRow struct {
	Token                string
	CanonicalToken       string
	EntityName           string
	Kind                 string
	OwningProjectID      string
	SourceMetadataPath   string
	SuggestedPlaceholder string
}

// projectMetadataDerivedRows is the pure (DB-free) projection core: given a
// parsed metadata file and its repo-relative path, it returns the
// governance + execution rows. Keeping this pure makes case-variant
// expansion and cross-entity token collision behaviour unit-testable
// without a database.
//
// Only `visibility: private` projects contribute rows; public projects
// declare no protected surface. An entity with no non-empty match_tokens
// contributes a governance row but no token rows (the Phase 1A validator
// already rejects such entities, so this is a defensive no-op).
func projectMetadataDerivedRows(metadata *ProjectMetadataFile, rel string) ([]derivedPrivateEntityRow, []derivedMatchTokenRow) {
	if metadata == nil {
		return nil, nil
	}
	if strings.ToLower(strings.TrimSpace(metadata.Project.Visibility)) != "private" {
		return nil, nil
	}
	projectID := strings.TrimSpace(metadata.Project.ID)
	placeholder := projectPlaceholder(projectID)

	var entities []derivedPrivateEntityRow
	var tokens []derivedMatchTokenRow

	for _, entity := range metadata.Project.PrivateEntities {
		name := strings.TrimSpace(entity.Name)
		kind := strings.TrimSpace(entity.Kind)
		entities = append(entities, derivedPrivateEntityRow{
			EntityName:           name,
			Kind:                 kind,
			OwningProjectID:      projectID,
			SourceMetadataPath:   rel,
			SuggestedPlaceholder: placeholder,
		})

		for token, canonical := range expandEntityMatchTokens(entity) {
			tokens = append(tokens, derivedMatchTokenRow{
				Token:                token,
				CanonicalToken:       canonical,
				EntityName:           name,
				Kind:                 kind,
				OwningProjectID:      projectID,
				SourceMetadataPath:   rel,
				SuggestedPlaceholder: placeholder,
			})
		}
	}

	// Deterministic ordering for stable inserts and test assertions.
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].EntityName < entities[j].EntityName
	})
	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].EntityName != tokens[j].EntityName {
			return tokens[i].EntityName < tokens[j].EntityName
		}
		return tokens[i].Token < tokens[j].Token
	})
	return entities, tokens
}

// expandEntityMatchTokens returns the full set of execution-layer tokens
// for one entity, mapping each variant to the canonical declared form it
// traces to.
//
//   - Every match_token literal always maps to itself.
//   - case_variants: auto (the default) adds the standard case variants of
//     each match_token (via the shared sanitizationTokenVariants helper),
//     each mapping back to the match_token it derived from.
//   - case_variants: [explicit, list] adds those forms verbatim (entity-
//     level, not per-token), each mapping to itself; auto-derivation is
//     suppressed because the author fully enumerated the surface.
//
// The first writer of a given token wins the canonical mapping, so a
// match_token literal is never overwritten by a later auto-variant.
func expandEntityMatchTokens(entity ProjectEntity) map[string]string {
	out := map[string]string{}
	add := func(tok, canonical string) {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			return
		}
		if _, exists := out[tok]; !exists {
			out[tok] = strings.TrimSpace(canonical)
		}
	}

	explicit := entity.CaseVariants.Mode == "explicit"
	for _, mt := range entity.MatchTokens {
		mt = strings.TrimSpace(mt)
		if mt == "" {
			continue
		}
		add(mt, mt)
		if !explicit {
			for _, variant := range sanitizationTokenVariants(mt) {
				add(variant, mt)
			}
		}
	}
	if explicit {
		for _, variant := range entity.CaseVariants.Variants {
			add(variant, variant)
		}
	}
	return out
}

// projectMetadataShape classifies a .ai-skill-project.yaml by the schema
// generation it declares, so compileProjectMetadataDerived can apply
// shape-aware tolerance.
type projectMetadataShape int

const (
	// shapeNewSchema: `private_entities:` is a list of mappings (objects),
	// or the file declares no private surface at all. Parsed strictly.
	shapeNewSchema projectMetadataShape = iota
	// shapeLegacyFlat: `private_tokens:` present, or `private_entities:` is
	// a list of scalars. Deprecated; tolerated-and-skipped.
	shapeLegacyFlat
)

// classifyProjectMetadataShape inspects the raw YAML structure (without the
// strict v1 parser) to decide whether a file is new-schema or legacy-flat.
// A YAML parse error is returned to the caller (categorically an I/O-class
// failure, not a tolerated skip).
func classifyProjectMetadataShape(body []byte) (projectMetadataShape, error) {
	var probe struct {
		Project struct {
			PrivateTokens   yaml.Node `yaml:"private_tokens"`
			PrivateEntities yaml.Node `yaml:"private_entities"`
		} `yaml:"project"`
	}
	if err := yaml.Unmarshal(body, &probe); err != nil {
		return shapeNewSchema, fmt.Errorf("parse project metadata yaml: %w", err)
	}
	if pt := probe.Project.PrivateTokens; pt.Kind == yaml.SequenceNode && len(pt.Content) > 0 {
		return shapeLegacyFlat, nil
	}
	if pe := probe.Project.PrivateEntities; pe.Kind == yaml.SequenceNode && len(pe.Content) > 0 {
		// Inspect the first element: scalar → legacy flat list of strings;
		// mapping → new-schema entity objects.
		if pe.Content[0].Kind == yaml.ScalarNode {
			return shapeLegacyFlat, nil
		}
	}
	return shapeNewSchema, nil
}

// compileProjectMetadataDerived walks every .ai-skill-project.yaml in the
// repo and projects the governance + execution rows. This is the SOLE
// project-metadata projection since Phase 1D retired the legacy
// derived_forbidden_tokens path.
//
// Shape-aware tolerance:
//   - legacy flat-shape file → stderr warning + skip (deprecated schema,
//     no projection target anymore)
//   - new-schema file that fails validation → HARD ERROR (a typo must not
//     silently drop protected tokens)
//   - I/O errors (unreadable file, malformed YAML) → propagated
func compileProjectMetadataDerived(repo string, db *sql.DB) error {
	metadataFiles, err := discoverProjectMetadataFiles(repo)
	if err != nil {
		return err
	}
	for _, rel := range metadataFiles {
		body, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(rel)))
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}
		shape, err := classifyProjectMetadataShape(body)
		if err != nil {
			return fmt.Errorf("classify %s: %w", rel, err)
		}
		if shape == shapeLegacyFlat {
			fmt.Fprintf(os.Stderr,
				"sanitization: skipping deprecated flat-shape project metadata %s; migrate private_tokens / scalar private_entities to the private_entities object schema (metadata/project/ai-skill-project-schema.yaml)\n",
				rel)
			continue
		}

		metadata, err := ParseProjectMetadata(body)
		if err != nil {
			if IsValidationError(err) {
				// New-schema shape but malformed: fail loudly. Silent skip
				// here is exactly the protection gap this remediation closes.
				return fmt.Errorf("invalid project metadata %s: %w", rel, err)
			}
			return fmt.Errorf("parse %s: %w", rel, err)
		}

		entities, tokens := projectMetadataDerivedRows(metadata, rel)
		for _, e := range entities {
			if _, err := db.Exec(
				`INSERT OR REPLACE INTO derived_private_entities (entity_name, kind, owning_project_id, source_metadata_path, suggested_placeholder) VALUES (?, ?, ?, ?, ?)`,
				e.EntityName, e.Kind, e.OwningProjectID, e.SourceMetadataPath, e.SuggestedPlaceholder,
			); err != nil {
				return fmt.Errorf("insert derived_private_entities row (%s/%s): %w", rel, e.EntityName, err)
			}
		}
		for _, tok := range tokens {
			if _, err := db.Exec(
				`INSERT OR REPLACE INTO derived_match_tokens (token, canonical_token, entity_name, kind, owning_project_id, source_metadata_path, suggested_placeholder) VALUES (?, ?, ?, ?, ?, ?, ?)`,
				tok.Token, tok.CanonicalToken, tok.EntityName, tok.Kind, tok.OwningProjectID, tok.SourceMetadataPath, tok.SuggestedPlaceholder,
			); err != nil {
				return fmt.Errorf("insert derived_match_tokens row (%s/%s/%s): %w", rel, tok.EntityName, tok.Token, err)
			}
		}
		// Traceability: register the source file now that the legacy
		// projection (which previously owned this registration) is retired.
		if err := insertProjectMetadataDerivedSourceFile(db, rel); err != nil {
			return err
		}
	}
	return nil
}

// insertProjectMetadataDerivedSourceFile records the .ai-skill-project.yaml
// source file in runtime_source_files, targeting the execution-layer table.
func insertProjectMetadataDerivedSourceFile(db *sql.DB, rel string) error {
	_, err := db.Exec(
		`INSERT OR REPLACE INTO runtime_source_files (source_path, source_kind, target_table, compile_rule, compiled_at, compiler_version, status) VALUES (?, 'project_metadata', 'derived_match_tokens', 'project_metadata_private_entities', datetime('now'), ?, 'synced')`,
		rel, goRuntimeCompilerVersion,
	)
	return err
}
