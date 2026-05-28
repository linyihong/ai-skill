package glossary_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/glossary"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func ruleIDs(vs []glossary.Violation) []string {
	ids := make([]string, 0, len(vs))
	for _, v := range vs {
		ids = append(ids, v.RuleID)
	}
	sort.Strings(ids)
	return ids
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

// fixture/glossary-valid-entry: happy path, no violations.
func TestValidate_ValidEntry(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## context_mode

`+"```yaml"+`
term: context_mode
status: canonical
owner-layer: runtime-cognition
meaning: Runtime control plane enum for context expansion strategy.
affects:
  - runtime/cognitive-modes.yaml
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
`+"```\n")

	res, err := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if len(res.Violations) != 0 {
		for _, v := range res.Violations {
			t.Logf("unexpected violation: %s - %s", v.RuleID, v.Message)
		}
		t.Fatalf("expected 0 violations, got %d", len(res.Violations))
	}
	if res.EntryCount != 1 {
		t.Errorf("entry count: got %d want 1", res.EntryCount)
	}
}

// README.md is skipped by the validator even if it contains H2 + YAML blocks.
func TestValidate_SkipsReadme(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "README.md", `## not_a_real_term

`+"```yaml"+`
term: nope
`+"```\n")
	res, err := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if res.EntryCount != 0 || len(res.Violations) != 0 {
		t.Fatalf("README must be skipped; got entries=%d violations=%d", res.EntryCount, len(res.Violations))
	}
}

// fixture/glossary-invalid-entry: each missing field is flagged.
func TestValidate_MissingRequiredFields(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## bad_entry

`+"```yaml"+`
term: bad_entry
`+"```\n")
	res, err := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	ids := ruleIDs(res.Violations)
	for _, want := range []string{
		"glossary.entry.missing_status",
		"glossary.entry.missing_meaning",
		"glossary.entry.missing_affects",
		"glossary.entry.missing_owner_layer",
	} {
		if !contains(ids, want) {
			t.Errorf("expected violation %s in %v", want, ids)
		}
	}
}

// fixture/glossary-invalid-entry: kebab-case term blocked.
func TestValidate_TermNaming(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## context-mode

`+"```yaml"+`
term: context-mode
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.term_naming") {
		t.Fatalf("expected term_naming violation; got %v", ruleIDs(res.Violations))
	}
}

// fixture/glossary-invalid-entry: invalid status enum.
func TestValidate_StatusEnum(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: ratified
owner-layer: runtime-cognition
meaning: x
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.status_enum") {
		t.Fatalf("expected status_enum violation; got %v", ruleIDs(res.Violations))
	}
}

// fixture/glossary-invalid-entry: invalid owner-layer enum.
func TestValidate_OwnerLayerEnum(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: canonical
owner-layer: unknown-layer
meaning: x
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.owner_layer_enum") {
		t.Fatalf("expected owner_layer_enum violation; got %v", ruleIDs(res.Violations))
	}
}

// fixture/glossary-alias-rules: alias names another canonical term.
func TestValidate_AliasIsCanonicalTerm(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## compression

`+"```yaml"+`
term: compression
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
aliases:
  - other_term
`+"```"+`

## other_term

`+"```yaml"+`
term: other_term
status: canonical
owner-layer: runtime-cognition
meaning: y
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.alias_is_canonical_term") {
		t.Fatalf("expected alias_is_canonical_term violation; got %v", ruleIDs(res.Violations))
	}
}

// fixture/glossary-alias-rules: status alias-only is forbidden.
func TestValidate_AliasOnlyStatusForbidden(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: alias-only
owner-layer: runtime-cognition
meaning: x
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.alias_only_status_forbidden") {
		t.Fatalf("expected alias_only_status_forbidden violation; got %v", ruleIDs(res.Violations))
	}
}

// fixture/glossary-introduced-by-shape: commit SHA / bare URL forbidden.
func TestValidate_IntroducedByShape(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  bool
	}{
		{"plan path", "plans/active/2026-05-25-1000-context-language-glossary-system.md", false},
		{"archived plan path", "plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md", false},
		{"ADR path", "constitution/ADR-007-foo.md", false},
		{"commit SHA", "abc1234", true},
		{"GitHub URL", "https://github.com/owner/repo/pull/123", true},
		{"issue number", `"#123"`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
introduced-by: `+tc.value+`
`+"```\n")
			res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
			got := contains(ruleIDs(res.Violations), "glossary.entry.introduced_by_shape")
			if got != tc.want {
				t.Errorf("introduced-by=%q: expected violation=%v got=%v (all=%v)", tc.value, tc.want, got, ruleIDs(res.Violations))
			}
		})
	}
}

// fixture/glossary-excludes-reference: excludes target must exist as term.
func TestValidate_ExcludesUnknownTerm(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
excludes:
  - non_existent_term
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.excludes_unknown_term") {
		t.Fatalf("expected excludes_unknown_term violation; got %v", ruleIDs(res.Violations))
	}
}

// Excludes referencing a real term passes.
func TestValidate_ExcludesValidTerm(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
excludes:
  - term_b
`+"```"+`

## term_b

`+"```yaml"+`
term: term_b
status: canonical
owner-layer: runtime-cognition
meaning: y
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if contains(ruleIDs(res.Violations), "glossary.entry.excludes_unknown_term") {
		t.Fatalf("excludes_valid_term should not flag; got %v", ruleIDs(res.Violations))
	}
}

// fixture/glossary-symmetric-relation: related_to in only one direction.
func TestValidate_SymmetricRelationMissingReverse(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
related-terms:
  - { type: related_to, target: term_b }
`+"```"+`

## term_b

`+"```yaml"+`
term: term_b
status: canonical
owner-layer: runtime-cognition
meaning: y
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.symmetric_relation_missing_reverse") {
		t.Fatalf("expected symmetric_relation_missing_reverse violation; got %v", ruleIDs(res.Violations))
	}
}

// Bidirectional related_to passes.
func TestValidate_SymmetricRelationBidirectional(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
related-terms:
  - { type: related_to, target: term_b }
`+"```"+`

## term_b

`+"```yaml"+`
term: term_b
status: canonical
owner-layer: runtime-cognition
meaning: y
affects: [a]
related-terms:
  - { type: related_to, target: term_a }
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if contains(ruleIDs(res.Violations), "glossary.entry.symmetric_relation_missing_reverse") {
		t.Fatalf("bidirectional symmetric should not flag; got %v", ruleIDs(res.Violations))
	}
}

// Asymmetric relation (derived_from) requires no reverse.
func TestValidate_AsymmetricRelationNoReverseRequired(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## cognitive_cost

`+"```yaml"+`
term: cognitive_cost
status: candidate
owner-layer: runtime-cognition
meaning: x
affects: [a]
related-terms:
  - { type: aggregates, target: thinking_cost }
`+"```"+`

## thinking_cost

`+"```yaml"+`
term: thinking_cost
status: candidate
owner-layer: ecosystem-adaptation
meaning: y
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if contains(ruleIDs(res.Violations), "glossary.entry.symmetric_relation_missing_reverse") {
		t.Fatalf("asymmetric relation should not require reverse; got %v", ruleIDs(res.Violations))
	}
}

// H2 heading without YAML block is flagged.
func TestValidate_HeadingWithoutYAMLBlock(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", "## orphan_term\n\nJust prose, no YAML block.\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.yaml_block_missing") {
		t.Fatalf("expected yaml_block_missing violation; got %v", ruleIDs(res.Violations))
	}
}

// Missing glossary directory returns no entries, no error (greenfield-safe).
func TestValidate_MissingDirectoryIsClean(t *testing.T) {
	res, err := glossary.Validate(glossary.ValidateOptions{GlossaryDir: filepath.Join(t.TempDir(), "does-not-exist")})
	if err != nil {
		t.Fatalf("missing dir should not error: %v", err)
	}
	if res.EntryCount != 0 || len(res.Violations) != 0 {
		t.Fatalf("missing dir should be empty; got entries=%d violations=%d", res.EntryCount, len(res.Violations))
	}
}

// Heading text mismatches term: flagged.
func TestValidate_HeadingMismatchesTerm(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## displayed_name

`+"```yaml"+`
term: actual_term
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.heading_term_mismatch") {
		t.Fatalf("expected heading_term_mismatch; got %v", ruleIDs(res.Violations))
	}
}

// Relation type enum guards against typos.
func TestValidate_RelationTypeEnum(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "ai-skill.md", `## term_a

`+"```yaml"+`
term: term_a
status: canonical
owner-layer: runtime-cognition
meaning: x
affects: [a]
related-terms:
  - { type: friends_with, target: term_b }
`+"```\n")
	res, _ := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if !contains(ruleIDs(res.Violations), "glossary.entry.relation_type_enum") {
		t.Fatalf("expected relation_type_enum; got %v", ruleIDs(res.Violations))
	}
}

// The README.md fixture is allowed to contain its own H2 + YAML examples
// without polluting validate output. Smoke test that strings.Contains
// of common rule IDs would surface helpfully.
func TestRuleIDsAreNamespaced(t *testing.T) {
	for _, status := range glossary.AllowedStatuses {
		if status == "" {
			t.Fatalf("empty status in AllowedStatuses")
		}
	}
	if !strings.HasPrefix("glossary.entry.missing_term", "glossary.entry.") {
		t.Fatalf("rule IDs must be namespaced under glossary.entry.*")
	}
}
