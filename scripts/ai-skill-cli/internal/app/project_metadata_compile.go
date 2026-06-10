package app

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// project_metadata_compile.go — Phase 1C₂ projection rule for
// <PROJECT_ROOT>/.ai-skill-project.yaml.
//
// Plan reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1C₂.
//
// This file wires the Phase 1A canonical parser (`LoadProjectMetadata`)
// into the compile pipeline, projecting each project's `private_entities`
// declarations into TWO tables that mirror the governance/execution layer
// split the plan's Q8 resolution established:
//
//   - derived_private_entities (governance layer): one row per declared
//     entity, carrying the human-readable identity (name + kind) the
//     scanner's findings and audit output refer to.
//   - derived_match_tokens (execution layer): one row per case-variant of
//     every entity's match_tokens — the literal surface the Phase 1D
//     scanner compares staged content against. Case-variant expansion
//     happens HERE (projection), never in the parser (Phase 1A is a pure
//     reader).
//
// SCOPE DISCIPLINE — Phase 1C₂ is ADDITIVE. It does NOT touch:
//   - sanitization_scan.go (the legacy reader/scanner; Phase 1D migrates
//     the scanner's query from derived_forbidden_tokens to
//     derived_match_tokens and retires the legacy table)
//   - the legacy compileDerivedForbiddenTokens projection, which keeps
//     populating derived_forbidden_tokens until Phase 1D
//
// The legacy and new projections walk the same .ai-skill-project.yaml
// files. Source-file traceability registration is left to the legacy
// projection (insertProjectMetadataSourceFile) to avoid a
// runtime_source_files PRIMARY KEY (source_path) collision; .ai-skill-
// project.yaml files are discovered project files, not runtime canonical
// documents, so no assertRuntimeCanonicalDocumentsProjected invariant
// applies to them.

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

// compileProjectMetadataDerived walks every .ai-skill-project.yaml in the
// repo, parses it with the Phase 1A canonical parser, and projects the
// governance + execution rows.
//
// Transition tolerance: a file that fails the v1 schema validation is
// SKIPPED rather than aborting the whole compile. During the migration
// window, legacy project files may use the flat `private_tokens` shape
// that the new parser does not recognise; those files are still covered by
// the legacy derived_forbidden_tokens projection, so skipping them here is
// safe and additive. I/O errors (unreadable file, malformed YAML) are
// still propagated, since they are categorically different from a
// schema-shape mismatch.
func compileProjectMetadataDerived(repo string, db *sql.DB) error {
	metadataFiles, err := discoverProjectMetadataFiles(repo)
	if err != nil {
		return err
	}
	for _, rel := range metadataFiles {
		metadata, err := LoadProjectMetadata(filepath.Join(repo, filepath.FromSlash(rel)))
		if err != nil {
			if IsValidationError(err) {
				// Covered by the legacy projection during transition.
				continue
			}
			return fmt.Errorf("read %s: %w", rel, err)
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
	}
	return nil
}
