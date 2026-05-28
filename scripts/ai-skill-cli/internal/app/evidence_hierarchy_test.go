package app

import (
	"strings"
	"testing"
)

func TestValidateEvidenceHierarchyOptOut(t *testing.T) {
	body := "feat(x): Phase 1 complete\n\n[skip-evidence-hierarchy]\n"
	got := validateEvidenceHierarchy(body, []string{"scripts/ai-skill-cli/internal/app/x.go"}, "")
	if got != "" {
		t.Errorf("opt-out should bypass, got %q", got)
	}
}

func TestValidateEvidenceHierarchyNoSuccessClaim(t *testing.T) {
	got := validateEvidenceHierarchy("feat(x): add stuff", []string{"scripts/ai-skill-cli/internal/app/x.go"}, "")
	if got != "" {
		t.Errorf("non-completion body should bypass, got %q", got)
	}
}

func TestValidateEvidenceHierarchyDocOnlyBypass(t *testing.T) {
	got := validateEvidenceHierarchy("docs: Phase 1 complete", []string{"README.md"}, "")
	if got != "" {
		t.Errorf("doc-only staging should bypass, got %q", got)
	}
}

func TestValidateEvidenceHierarchyBlocksUnsupportedCompletionClaim(t *testing.T) {
	got := validateEvidenceHierarchy("feat(x): Phase 1 完成", []string{"scripts/ai-skill-cli/internal/app/x.go"}, "")
	if !strings.Contains(got, "evidence-hierarchy") {
		t.Errorf("expected evidence-hierarchy violation for bare completion claim, got %q", got)
	}
}

func TestValidateEvidenceHierarchyAcceptsTestPassEvidence(t *testing.T) {
	body := "feat(x): done — tests pass and fixture green"
	got := validateEvidenceHierarchy(body, []string{"scripts/ai-skill-cli/internal/app/x.go"}, "")
	if got != "" {
		t.Errorf("test-pass evidence should satisfy validator, got %q", got)
	}
}

func TestValidateEvidenceHierarchyAcceptsAuditEvidence(t *testing.T) {
	body := "feat(x): completed — `ai-skill runtime audit` shows orphan count down from 242 to 230"
	got := validateEvidenceHierarchy(body, []string{"runtime/foo.yaml"}, "")
	if got != "" {
		t.Errorf("audit evidence should satisfy validator, got %q", got)
	}
}

func TestValidateEvidenceHierarchyAcceptsScenarioCitation(t *testing.T) {
	body := "feat(x): ✅ scenario evidence-hierarchy-v1 covers this path"
	got := validateEvidenceHierarchy(body, []string{"validation/scenarios/foo.yaml"}, "")
	if got != "" {
		t.Errorf("scenario citation should satisfy validator, got %q", got)
	}
}

func TestValidateEvidenceHierarchyAcceptsCommitHashEvidence(t *testing.T) {
	body := "feat(x): done — based on commit abc1234 which already validated this surface"
	got := validateEvidenceHierarchy(body, []string{"governance/foo.yaml"}, "")
	if got != "" {
		t.Errorf("commit-hash reference should satisfy validator, got %q", got)
	}
}
