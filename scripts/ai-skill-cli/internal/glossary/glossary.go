// Package glossary implements the parser and validator for
// knowledge/glossary/*.md entries. Schema spec source-of-truth is
// knowledge/glossary/README.md.
package glossary

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Entry is the YAML body of a single glossary term entry.
type Entry struct {
	Term         string     `yaml:"term"`
	Status       string     `yaml:"status"`
	Meaning      string     `yaml:"meaning"`
	Affects      []string   `yaml:"affects"`
	OwnerLayer   string     `yaml:"owner-layer"`
	Aliases      []string   `yaml:"aliases,omitempty"`
	AntiMeaning  string     `yaml:"anti-meaning,omitempty"`
	Excludes     []string   `yaml:"excludes,omitempty"`
	RelatedTerms []Relation `yaml:"related-terms,omitempty"`
	IntroducedBy string     `yaml:"introduced-by,omitempty"`
	DeprecatedBy string     `yaml:"deprecated-by,omitempty"`
}

// Relation is a single directed relation between two glossary terms.
type Relation struct {
	Type   string `yaml:"type"`
	Target string `yaml:"target"`
}

// ParsedEntry pairs a YAML Entry with its source location.
type ParsedEntry struct {
	Entry
	SourceFile  string
	Heading     string
	HeadingLine int
	YAMLLine    int
}

// Violation is a single schema check failure.
type Violation struct {
	File        string `json:"file,omitempty"`
	Term        string `json:"term,omitempty"`
	RuleID      string `json:"rule_id"`
	Message     string `json:"message"`
	Remediation string `json:"remediation,omitempty"`
}

// ValidateOptions controls validate behavior.
type ValidateOptions struct {
	GlossaryDir string
}

// ValidateResult aggregates counts and violations.
type ValidateResult struct {
	EntryCount    int
	AliasCount    int
	RelationCount int
	Violations    []Violation
}

// Allowed enums. Source: knowledge/glossary/README.md.
var (
	AllowedStatuses = []string{
		"canonical", "candidate", "deprecated", "superseded",
		"alias-only", "experimental", "project-local",
	}
	AllowedOwnerLayers = []string{
		"runtime-cognition", "semantic-routing", "workflow-orchestration",
		"validation-governance", "memory-replay", "runtime-projection",
		"architecture-contracts", "ecosystem-adaptation", "runtime-economics",
	}
	AllowedRelationTypes = []string{
		"alias_of", "related_to", "conflicts_with",
		"owned_by", "used_by", "deprecated_by", "replaced_by",
		"derived_from", "aggregates",
	}
	SymmetricRelations = map[string]bool{
		"related_to":     true,
		"conflicts_with": true,
	}
)

var (
	snakeCaseRe    = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	h2HeadingRe    = regexp.MustCompile(`^## (.+)$`)
	yamlOpenRe     = regexp.MustCompile("^```yaml\\s*$")
	yamlCloseRe    = regexp.MustCompile("^```\\s*$")
	introducedByRe = regexp.MustCompile(`^(plans/[A-Za-z0-9._/\-]+\.md|constitution/ADR-[A-Za-z0-9_\-]+\.md)$`)
)

// LoadDir walks dir and parses every *.md file except README.md.
func LoadDir(dir string) ([]ParsedEntry, []Violation, error) {
	var allEntries []ParsedEntry
	var allViolations []Violation
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("glossary path is not a directory: %s", dir)
	}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		if filepath.Base(path) == "README.md" {
			return nil
		}
		entries, vs, perr := ParseFile(path)
		if perr != nil {
			return perr
		}
		allEntries = append(allEntries, entries...)
		allViolations = append(allViolations, vs...)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return allEntries, allViolations, nil
}

// ParseFile reads path and parses glossary entries.
func ParseFile(path string) ([]ParsedEntry, []Violation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	entries, vs := ParseBytes(path, data)
	return entries, vs, nil
}

// ParseBytes parses raw Markdown content into entries.
func ParseBytes(file string, data []byte) ([]ParsedEntry, []Violation) {
	lines := strings.Split(string(data), "\n")
	var entries []ParsedEntry
	var violations []Violation

	i := 0
	for i < len(lines) {
		line := strings.TrimRight(lines[i], "\r")
		m := h2HeadingRe.FindStringSubmatch(line)
		if m == nil {
			i++
			continue
		}
		heading := strings.TrimSpace(m[1])
		headingLine := i + 1

		// Find ```yaml after optional blank lines.
		j := i + 1
		for j < len(lines) && strings.TrimSpace(strings.TrimRight(lines[j], "\r")) == "" {
			j++
		}
		if j >= len(lines) || !yamlOpenRe.MatchString(strings.TrimRight(lines[j], "\r")) {
			violations = append(violations, Violation{
				File:    file,
				Term:    heading,
				RuleID:  "glossary.entry.yaml_block_missing",
				Message: fmt.Sprintf("H2 heading '%s' at %s:%d is not followed by a ```yaml code block", heading, file, headingLine),
				Remediation: "Add a ```yaml ... ``` block immediately after the H2 heading.",
			})
			i++
			continue
		}
		yamlStart := j + 1
		k := yamlStart
		for k < len(lines) && !yamlCloseRe.MatchString(strings.TrimRight(lines[k], "\r")) {
			k++
		}
		if k >= len(lines) {
			violations = append(violations, Violation{
				File:    file,
				Term:    heading,
				RuleID:  "glossary.entry.yaml_block_unterminated",
				Message: fmt.Sprintf("YAML block for '%s' starting at %s:%d is not terminated", heading, file, j+1),
			})
			i++
			continue
		}
		body := strings.Join(lines[yamlStart:k], "\n")
		var entry Entry
		if err := yaml.Unmarshal([]byte(body), &entry); err != nil {
			violations = append(violations, Violation{
				File:    file,
				Term:    heading,
				RuleID:  "glossary.entry.yaml_parse_error",
				Message: fmt.Sprintf("YAML parse error for '%s' at %s:%d: %v", heading, file, j+1, err),
			})
			i = k + 1
			continue
		}
		entries = append(entries, ParsedEntry{
			Entry:       entry,
			SourceFile:  file,
			Heading:     heading,
			HeadingLine: headingLine,
			YAMLLine:    j + 1,
		})
		i = k + 1
	}
	return entries, violations
}

// Validate runs all schema checks against entries discovered in opts.GlossaryDir.
func Validate(opts ValidateOptions) (ValidateResult, error) {
	entries, parseViolations, err := LoadDir(opts.GlossaryDir)
	if err != nil {
		return ValidateResult{}, err
	}
	violations := append([]Violation{}, parseViolations...)

	termIndex := make(map[string]*ParsedEntry, len(entries))
	for i := range entries {
		termIndex[entries[i].Term] = &entries[i]
	}

	aliasCount, relationCount := 0, 0
	for i := range entries {
		e := &entries[i]
		violations = append(violations, checkHeadingMatchesTerm(e)...)
		violations = append(violations, checkRequiredFields(e)...)
		violations = append(violations, checkStatus(e)...)
		violations = append(violations, checkOwnerLayer(e)...)
		violations = append(violations, checkTermNaming(e)...)
		violations = append(violations, checkRelationTypes(e)...)
		violations = append(violations, checkRelationTargets(e, termIndex)...)
		violations = append(violations, checkAliasRules(e, termIndex)...)
		violations = append(violations, checkIntroducedBy(e)...)
		violations = append(violations, checkDeprecatedBy(e)...)
		violations = append(violations, checkExcludes(e, termIndex)...)
		aliasCount += len(e.Aliases)
		relationCount += len(e.RelatedTerms)
	}
	violations = append(violations, checkAliasCycles(entries)...)
	violations = append(violations, checkSymmetricRelations(entries, termIndex)...)

	sort.SliceStable(violations, func(i, j int) bool {
		if violations[i].File != violations[j].File {
			return violations[i].File < violations[j].File
		}
		if violations[i].Term != violations[j].Term {
			return violations[i].Term < violations[j].Term
		}
		return violations[i].RuleID < violations[j].RuleID
	})

	return ValidateResult{
		EntryCount:    len(entries),
		AliasCount:    aliasCount,
		RelationCount: relationCount,
		Violations:    violations,
	}, nil
}

func checkHeadingMatchesTerm(e *ParsedEntry) []Violation {
	if e.Term == "" {
		return nil
	}
	if e.Heading != e.Term {
		return []Violation{{
			File:        e.SourceFile,
			Term:        e.Term,
			RuleID:      "glossary.entry.heading_term_mismatch",
			Message:     fmt.Sprintf("H2 heading '%s' at %s:%d does not match YAML term '%s'", e.Heading, e.SourceFile, e.HeadingLine, e.Term),
			Remediation: fmt.Sprintf("Change H2 to '## %s' or change YAML term to '%s'", e.Term, e.Heading),
		}}
	}
	return nil
}

func checkRequiredFields(e *ParsedEntry) []Violation {
	var vs []Violation
	name := e.Term
	if name == "" {
		name = e.Heading
	}
	add := func(rule, field string) {
		vs = append(vs, Violation{
			File:    e.SourceFile,
			Term:    name,
			RuleID:  rule,
			Message: fmt.Sprintf("Entry '%s' at %s:%d missing required field '%s'", name, e.SourceFile, e.HeadingLine, field),
		})
	}
	if e.Term == "" {
		add("glossary.entry.missing_term", "term")
	}
	if e.Status == "" {
		add("glossary.entry.missing_status", "status")
	}
	if e.Meaning == "" {
		add("glossary.entry.missing_meaning", "meaning")
	}
	if len(e.Affects) == 0 {
		add("glossary.entry.missing_affects", "affects (non-empty array)")
	}
	if e.OwnerLayer == "" {
		add("glossary.entry.missing_owner_layer", "owner-layer")
	}
	return vs
}

func checkStatus(e *ParsedEntry) []Violation {
	if e.Status == "" {
		return nil
	}
	for _, ok := range AllowedStatuses {
		if e.Status == ok {
			return nil
		}
	}
	return []Violation{{
		File:        e.SourceFile,
		Term:        e.Term,
		RuleID:      "glossary.entry.status_enum",
		Message:     fmt.Sprintf("Entry '%s' has invalid status '%s'", e.Term, e.Status),
		Remediation: "Use one of: " + strings.Join(AllowedStatuses, ", "),
	}}
}

func checkOwnerLayer(e *ParsedEntry) []Violation {
	if e.OwnerLayer == "" {
		return nil
	}
	for _, ok := range AllowedOwnerLayers {
		if e.OwnerLayer == ok {
			return nil
		}
	}
	return []Violation{{
		File:        e.SourceFile,
		Term:        e.Term,
		RuleID:      "glossary.entry.owner_layer_enum",
		Message:     fmt.Sprintf("Entry '%s' has invalid owner-layer '%s'", e.Term, e.OwnerLayer),
		Remediation: "Use one of: " + strings.Join(AllowedOwnerLayers, ", "),
	}}
}

func checkTermNaming(e *ParsedEntry) []Violation {
	if e.Term == "" {
		return nil
	}
	if !snakeCaseRe.MatchString(e.Term) {
		return []Violation{{
			File:        e.SourceFile,
			Term:        e.Term,
			RuleID:      "glossary.entry.term_naming",
			Message:     fmt.Sprintf("Term '%s' is not snake_case (must match /^[a-z][a-z0-9_]*$/)", e.Term),
			Remediation: "Use snake_case: lowercase, underscores, digits (digits not as first char).",
		}}
	}
	return nil
}

func checkRelationTypes(e *ParsedEntry) []Violation {
	var vs []Violation
	for _, rel := range e.RelatedTerms {
		if rel.Type == "" {
			vs = append(vs, Violation{
				File:    e.SourceFile,
				Term:    e.Term,
				RuleID:  "glossary.entry.relation_missing_type",
				Message: fmt.Sprintf("Entry '%s' has related-terms item missing 'type' (target=%q)", e.Term, rel.Target),
			})
			continue
		}
		ok := false
		for _, allowed := range AllowedRelationTypes {
			if rel.Type == allowed {
				ok = true
				break
			}
		}
		if !ok {
			vs = append(vs, Violation{
				File:        e.SourceFile,
				Term:        e.Term,
				RuleID:      "glossary.entry.relation_type_enum",
				Message:     fmt.Sprintf("Entry '%s' has invalid relation type '%s' (target=%s)", e.Term, rel.Type, rel.Target),
				Remediation: "Use one of: " + strings.Join(AllowedRelationTypes, ", "),
			})
		}
		if rel.Target == "" {
			vs = append(vs, Violation{
				File:    e.SourceFile,
				Term:    e.Term,
				RuleID:  "glossary.entry.relation_missing_target",
				Message: fmt.Sprintf("Entry '%s' has relation type '%s' missing 'target'", e.Term, rel.Type),
			})
		}
	}
	return vs
}

func checkRelationTargets(e *ParsedEntry, index map[string]*ParsedEntry) []Violation {
	var vs []Violation
	for _, rel := range e.RelatedTerms {
		if rel.Target == "" {
			continue
		}
		if _, ok := index[rel.Target]; !ok {
			vs = append(vs, Violation{
				File:    e.SourceFile,
				Term:    e.Term,
				RuleID:  "glossary.entry.relation_target_unknown",
				Message: fmt.Sprintf("Entry '%s' has relation '%s' -> '%s' but target term is not defined in the glossary", e.Term, rel.Type, rel.Target),
			})
		}
	}
	return vs
}

func checkAliasRules(e *ParsedEntry, index map[string]*ParsedEntry) []Violation {
	var vs []Violation
	if e.Status == "alias-only" {
		vs = append(vs, Violation{
			File:        e.SourceFile,
			Term:        e.Term,
			RuleID:      "glossary.entry.alias_only_status_forbidden",
			Message:     fmt.Sprintf("Entry '%s' uses forbidden status 'alias-only'; new entries must declare aliases via the 'aliases:' field on the canonical entry instead", e.Term),
			Remediation: "Move this term into the canonical entry's aliases: array, or set status to candidate/canonical/experimental.",
		})
	}
	for _, a := range e.Aliases {
		if other, exists := index[a]; exists && other.Term != e.Term {
			vs = append(vs, Violation{
				File:        e.SourceFile,
				Term:        e.Term,
				RuleID:      "glossary.entry.alias_is_canonical_term",
				Message:     fmt.Sprintf("Entry '%s' lists alias '%s' which is itself a canonical entry (defined at %s:%d)", e.Term, a, other.SourceFile, other.HeadingLine),
				Remediation: fmt.Sprintf("Either remove '%s' from aliases or merge the two entries.", a),
			})
		}
	}
	return vs
}

func checkAliasCycles(entries []ParsedEntry) []Violation {
	var vs []Violation
	aliasMap := make(map[string][]string, len(entries))
	for _, e := range entries {
		aliasMap[e.Term] = e.Aliases
	}
	for term, aliases := range aliasMap {
		for _, a := range aliases {
			back, ok := aliasMap[a]
			if !ok {
				continue
			}
			for _, b := range back {
				if b == term {
					vs = append(vs, Violation{
						Term:    term,
						RuleID:  "glossary.entry.alias_cycle",
						Message: fmt.Sprintf("Alias cycle detected: '%s' aliases '%s' which aliases '%s'", term, a, term),
					})
				}
			}
		}
	}
	return vs
}

func checkIntroducedBy(e *ParsedEntry) []Violation {
	if e.IntroducedBy == "" {
		return nil
	}
	if !introducedByRe.MatchString(e.IntroducedBy) {
		return []Violation{{
			File:        e.SourceFile,
			Term:        e.Term,
			RuleID:      "glossary.entry.introduced_by_shape",
			Message:     fmt.Sprintf("Entry '%s' has invalid 'introduced-by' value '%s'", e.Term, e.IntroducedBy),
			Remediation: "Must be 'plans/<path>.md' or 'constitution/ADR-XXX.md'. Commit SHAs / issue numbers / PR URLs are forbidden.",
		}}
	}
	return nil
}

func checkDeprecatedBy(e *ParsedEntry) []Violation {
	if e.DeprecatedBy == "" {
		return nil
	}
	if !introducedByRe.MatchString(e.DeprecatedBy) {
		return []Violation{{
			File:        e.SourceFile,
			Term:        e.Term,
			RuleID:      "glossary.entry.deprecated_by_shape",
			Message:     fmt.Sprintf("Entry '%s' has invalid 'deprecated-by' value '%s'", e.Term, e.DeprecatedBy),
			Remediation: "Must be 'plans/<path>.md' or 'constitution/ADR-XXX.md'.",
		}}
	}
	return nil
}

func checkExcludes(e *ParsedEntry, index map[string]*ParsedEntry) []Violation {
	var vs []Violation
	for _, target := range e.Excludes {
		if _, ok := index[target]; !ok {
			vs = append(vs, Violation{
				File:        e.SourceFile,
				Term:        e.Term,
				RuleID:      "glossary.entry.excludes_unknown_term",
				Message:     fmt.Sprintf("Entry '%s' lists 'excludes: %s' but no such glossary term exists", e.Term, target),
				Remediation: fmt.Sprintf("Either add a glossary entry for '%s' or remove it from excludes.", target),
			})
		}
	}
	return vs
}

func checkSymmetricRelations(entries []ParsedEntry, index map[string]*ParsedEntry) []Violation {
	var vs []Violation
	for i := range entries {
		e := &entries[i]
		for _, rel := range e.RelatedTerms {
			if !SymmetricRelations[rel.Type] {
				continue
			}
			target, ok := index[rel.Target]
			if !ok {
				continue // already flagged by checkRelationTargets
			}
			if target.Term == e.Term {
				continue // self-reference; no symmetry required
			}
			found := false
			for _, rev := range target.RelatedTerms {
				if rev.Type == rel.Type && rev.Target == e.Term {
					found = true
					break
				}
			}
			if !found {
				vs = append(vs, Violation{
					File:        e.SourceFile,
					Term:        e.Term,
					RuleID:      "glossary.entry.symmetric_relation_missing_reverse",
					Message:     fmt.Sprintf("Entry '%s' has symmetric relation '%s' -> '%s', but '%s' does not have matching '%s' -> '%s'", e.Term, rel.Type, rel.Target, rel.Target, rel.Type, e.Term),
					Remediation: fmt.Sprintf("Add `{ type: %s, target: %s }` to entry '%s'.related-terms.", rel.Type, e.Term, rel.Target),
				})
			}
		}
	}
	return vs
}
