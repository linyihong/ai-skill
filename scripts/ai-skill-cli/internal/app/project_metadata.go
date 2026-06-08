package app

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// project_metadata.go — CANONICAL parser for <PROJECT_ROOT>/.ai-skill-project.yaml.
//
// Schema canonical: metadata/project/ai-skill-project-schema.yaml (v1).
// Plan reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1A.
//
// SCOPE DISCIPLINE — Phase 1A (this file) defines the canonical parser
// only. It is NOT wired into any consumer. The legacy reader in
// sanitization_scan.go::readProjectMetadata (flat `private_tokens` /
// `private_entities []string`) remains the live consumer of project
// metadata until Phase 1C lands the new projection rule. See
// metadata/project/migration-notes.md for the full migration trajectory.
//
// Reviewer constraint: this file MUST NOT import from or reference
// sanitization_scan.go symbols. Phase 1A scope explicitly forbids
// touching the legacy reader; even a one-line cross-file dependency
// would blur the 1A / 1C / 1D phase boundaries that the plan rewrite
// of 2026-06-08 established.

// ProjectMetadataFile is the root structure of <PROJECT_ROOT>/.ai-skill-project.yaml.
type ProjectMetadataFile struct {
	Project ProjectMetadata `yaml:"project"`

	// Unknown is captured for forward-compat: future schema versions may
	// add top-level siblings to `project:`. The parser tolerates unknown
	// top-level fields but does not consume them.
	Unknown map[string]any `yaml:",inline"`
}

// ProjectMetadata is the per-project declaration body.
type ProjectMetadata struct {
	ID              string           `yaml:"id"`
	Visibility      string           `yaml:"visibility"`
	PrivateEntities []ProjectEntity  `yaml:"private_entities"`

	// Unknown captures forward-compat fields at the project level.
	Unknown map[string]any `yaml:",inline"`
}

// ProjectEntity is a single private entity declaration.
//
// Governance layer fields: Name, Kind. These flow into audit / findings /
// debug output and identify the entity in human-readable terms.
//
// Execution layer fields: MatchTokens, CaseVariants. These describe what
// the scanner compares against; case-variant expansion is performed by
// Phase 1C projection rule, NOT by this parser.
type ProjectEntity struct {
	Name         string       `yaml:"name"`
	Kind         string       `yaml:"kind"`
	MatchTokens  []string     `yaml:"match_tokens"`
	CaseVariants CaseVariants `yaml:"case_variants"`

	Unknown map[string]any `yaml:",inline"`
}

// CaseVariants is either the literal string "auto" or an explicit list of
// case-variant strings. The schema field accepts both shapes; this type
// captures which shape was declared so Phase 1C can branch correctly.
type CaseVariants struct {
	// Mode is "auto" (default) or "explicit". "auto" means Phase 1C
	// projection rule expands match_tokens via standard case rules.
	// "explicit" means the projection rule uses Variants verbatim.
	Mode string

	// Variants is populated only when Mode == "explicit".
	Variants []string
}

// UnmarshalYAML accepts either the string "auto" or a list of strings.
// Empty / missing field defaults to {Mode: "auto"}.
func (c *CaseVariants) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		s := strings.TrimSpace(value.Value)
		if s == "" || s == "auto" {
			c.Mode = "auto"
			return nil
		}
		return fmt.Errorf("case_variants scalar must be 'auto' (got %q)", s)
	case yaml.SequenceNode:
		c.Mode = "explicit"
		c.Variants = nil
		for _, item := range value.Content {
			if item.Kind != yaml.ScalarNode {
				return fmt.Errorf("case_variants list entries must be strings")
			}
			s := strings.TrimSpace(item.Value)
			if s == "" {
				return fmt.Errorf("case_variants list entries must be non-empty strings")
			}
			c.Variants = append(c.Variants, s)
		}
		return nil
	case 0:
		// missing field
		c.Mode = "auto"
		return nil
	}
	return fmt.Errorf("case_variants must be the string 'auto' or a list of strings")
}

// ProjectMetadataValidationError aggregates rule failures from a single
// parse. Each entry references the validation_rules id declared in
// metadata/project/ai-skill-project-schema.yaml.
type ProjectMetadataValidationError struct {
	Failures []ProjectMetadataValidationFailure
}

// ProjectMetadataValidationFailure is one rule failure occurrence.
type ProjectMetadataValidationFailure struct {
	RuleID  string
	Path    string // dotted path into the file, e.g. "project.private_entities[0].kind"
	Message string
}

func (e *ProjectMetadataValidationError) Error() string {
	if len(e.Failures) == 0 {
		return "project metadata validation: no failures (this is a bug)"
	}
	lines := []string{fmt.Sprintf("project metadata validation: %d failure(s)", len(e.Failures))}
	for _, f := range e.Failures {
		lines = append(lines, fmt.Sprintf("  - [%s] %s: %s", f.RuleID, f.Path, f.Message))
	}
	return strings.Join(lines, "\n")
}

var projectIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)
var allowedVisibilities = map[string]bool{"private": true, "public": true}
var allowedKinds = map[string]bool{
	"codename": true, "client": true, "product": true, "individual": true, "other": true,
}

// LoadProjectMetadata reads and validates a .ai-skill-project.yaml file
// against the v1 canonical schema. Returns the parsed structure on
// success, or a *ProjectMetadataValidationError describing all rule
// failures on validation failure (the parser collects all failures
// rather than short-circuiting on the first; this gives the user a
// complete fix-list in one pass).
//
// I/O errors (file not found, malformed YAML) are returned as plain
// errors, not validation errors — they are categorically different.
func LoadProjectMetadata(path string) (*ProjectMetadataFile, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read project metadata file %s: %w", path, err)
	}
	return ParseProjectMetadata(body)
}

// ParseProjectMetadata is LoadProjectMetadata's byte-slice form, useful
// for tests that synthesize fixtures in-memory.
func ParseProjectMetadata(body []byte) (*ProjectMetadataFile, error) {
	var file ProjectMetadataFile
	if err := yaml.Unmarshal(body, &file); err != nil {
		return nil, fmt.Errorf("parse project metadata yaml: %w", err)
	}
	if verr := validateProjectMetadata(&file); verr != nil {
		return &file, verr
	}
	return &file, nil
}

func validateProjectMetadata(file *ProjectMetadataFile) error {
	var failures []ProjectMetadataValidationFailure

	p := file.Project
	if strings.TrimSpace(p.ID) == "" {
		failures = append(failures, ProjectMetadataValidationFailure{
			RuleID:  "project.id.required",
			Path:    "project.id",
			Message: "project.id must be present and non-empty",
		})
	} else if !projectIDPattern.MatchString(p.ID) {
		failures = append(failures, ProjectMetadataValidationFailure{
			RuleID:  "project.id.kebab",
			Path:    "project.id",
			Message: fmt.Sprintf("project.id must be kebab-case (got %q)", p.ID),
		})
	}

	if !allowedVisibilities[strings.TrimSpace(p.Visibility)] {
		failures = append(failures, ProjectMetadataValidationFailure{
			RuleID:  "project.visibility.enum",
			Path:    "project.visibility",
			Message: fmt.Sprintf("project.visibility must be 'private' or 'public' (got %q)", p.Visibility),
		})
	}

	for i, entity := range p.PrivateEntities {
		entityPath := fmt.Sprintf("project.private_entities[%d]", i)
		if strings.TrimSpace(entity.Name) == "" {
			failures = append(failures, ProjectMetadataValidationFailure{
				RuleID:  "entity.name.required",
				Path:    entityPath + ".name",
				Message: "entity name must be present and non-empty",
			})
		}
		if !allowedKinds[strings.TrimSpace(entity.Kind)] {
			failures = append(failures, ProjectMetadataValidationFailure{
				RuleID:  "entity.kind.required",
				Path:    entityPath + ".kind",
				Message: fmt.Sprintf("entity kind must be one of [codename, client, product, individual, other] (got %q)", entity.Kind),
			})
		}
		nonEmptyTokenCount := 0
		for _, tok := range entity.MatchTokens {
			if strings.TrimSpace(tok) != "" {
				nonEmptyTokenCount++
			}
		}
		if nonEmptyTokenCount == 0 {
			failures = append(failures, ProjectMetadataValidationFailure{
				RuleID:  "entity.match_tokens.required",
				Path:    entityPath + ".match_tokens",
				Message: "entity match_tokens must contain at least one non-empty string",
			})
		}
		// case_variants shape was validated by the UnmarshalYAML
		// implementation; no further check needed here. Anything that
		// reached this point is either {Mode: "auto"} or
		// {Mode: "explicit", Variants: [...]} with non-empty entries.
		switch entity.CaseVariants.Mode {
		case "", "auto", "explicit":
			// ok
		default:
			failures = append(failures, ProjectMetadataValidationFailure{
				RuleID:  "entity.case_variants.shape",
				Path:    entityPath + ".case_variants",
				Message: fmt.Sprintf("internal: unexpected CaseVariants.Mode %q", entity.CaseVariants.Mode),
			})
		}
	}

	if len(failures) == 0 {
		return nil
	}
	return &ProjectMetadataValidationError{Failures: failures}
}

// IsValidationError reports whether err is a *ProjectMetadataValidationError.
// Helpful for callers that want to branch on "structural error" vs
// "validation error".
func IsValidationError(err error) bool {
	var v *ProjectMetadataValidationError
	return errors.As(err, &v)
}
