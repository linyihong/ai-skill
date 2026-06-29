package planvalidate

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ImplementationStep is a minimal step model for advisory intent-transition checks.
type ImplementationStep struct {
	ID                      string
	Intent                  string
	EquivalenceRequired     bool
	EquivalenceEvidence     string
	ExplicitReopenReason    string
}

type implementationStepsYAML struct {
	Steps []implementationStepYAML `yaml:"steps"`
}

type implementationStepYAML struct {
	ID                   string `yaml:"id"`
	Intent               string `yaml:"intent"`
	ExplicitReopenReason string `yaml:"explicit_reopen_reason"`
	Checkpoint           struct {
		ObservableEquivalence struct {
			Required bool   `yaml:"required"`
			Evidence string `yaml:"evidence"`
			Passed   bool   `yaml:"passed"`
		} `yaml:"observable_equivalence"`
	} `yaml:"checkpoint"`
}

// ParseImplementationStepsFromMarkdown scans ```yaml blocks for a top-level steps list.
func ParseImplementationStepsFromMarkdown(markdown string) ([]ImplementationStep, error) {
	for _, block := range ExtractFencedYAMLBlocks(markdown) {
		steps, ok, err := parseImplementationStepsYAML(block)
		if err != nil {
			return nil, err
		}
		if ok {
			return steps, nil
		}
	}
	return nil, nil
}

func parseImplementationStepsYAML(block string) ([]ImplementationStep, bool, error) {
	var doc implementationStepsYAML
	if err := yaml.Unmarshal([]byte(block), &doc); err != nil {
		return nil, false, err
	}
	if len(doc.Steps) == 0 {
		return nil, false, nil
	}
	out := make([]ImplementationStep, 0, len(doc.Steps))
	for _, step := range doc.Steps {
		if strings.TrimSpace(step.Intent) == "" {
			continue
		}
		out = append(out, ImplementationStep{
			ID:                   step.ID,
			Intent:               strings.TrimSpace(step.Intent),
			EquivalenceRequired:  step.Checkpoint.ObservableEquivalence.Required,
			EquivalenceEvidence:  strings.TrimSpace(step.Checkpoint.ObservableEquivalence.Evidence),
			ExplicitReopenReason: strings.TrimSpace(step.ExplicitReopenReason),
		})
		if step.Checkpoint.ObservableEquivalence.Passed {
			out[len(out)-1].EquivalenceEvidence = out[len(out)-1].EquivalenceEvidence + "|passed"
		}
	}
	if len(out) == 0 {
		return nil, false, nil
	}
	return out, true, nil
}

func structureEquivalenceSatisfied(step ImplementationStep) bool {
	if !step.EquivalenceRequired {
		return true
	}
	return step.EquivalenceEvidence != ""
}

// DetectIllegalIntentTransitions applies the Intent Transition Rule from
// workflow/software-delivery/implementation/execution-modes.md.
// Findings are advisory (Blocking=false) until dogfood matures.
func DetectIllegalIntentTransitions(steps []ImplementationStep) []Finding {
	var out []Finding
	for i := 1; i < len(steps); i++ {
		prev := steps[i-1]
		curr := steps[i]
		switch {
		case prev.Intent == "structure" && curr.Intent == "feature":
			if !structureEquivalenceSatisfied(prev) {
				out = append(out, Finding{
					RuleID:   "implementation.intent.illegal_transition",
					Message:  fmt.Sprintf("step %q structure → %q feature without observable_equivalence_passed", prev.ID, curr.ID),
					Blocking: false,
				})
			}
		case prev.Intent == "feature" && curr.Intent == "structure":
			if curr.ExplicitReopenReason == "" {
				out = append(out, Finding{
					RuleID:   "implementation.intent.illegal_transition",
					Message:  fmt.Sprintf("step %q feature → %q structure without explicit_reopen_reason", prev.ID, curr.ID),
					Blocking: false,
				})
			}
		}
	}
	return out
}

// AdvisoryValidateImplementationIntent parses steps from plan markdown and reports illegal transitions.
func AdvisoryValidateImplementationIntent(markdown string) ([]Finding, error) {
	steps, err := ParseImplementationStepsFromMarkdown(markdown)
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return nil, nil
	}
	return DetectIllegalIntentTransitions(steps), nil
}
