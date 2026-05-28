package app

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopilotStartOutputsBootstrapPrompt(t *testing.T) {
	project := t.TempDir()
	writeFile(t, filepath.Join(project, ".ai-skill", "local.env"), "export AI_SKILL_REPO='/tmp/ai-skill'\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"copilot", "start", "--project", project, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "copilot start" || result.Mode != "guided_start" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("copilot start must not mutate files, got %#v", result.Mutations)
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected one generated prompt, got %#v", result.Results)
	}
	prompt := result.Results[0].Summary
	for _, expected := range []string{
		"<AI_SKILL_REPO>/CORE_BOOTSTRAP.md",
		"<AI_SKILL_REPO>/runtime/core-bootstrap.yaml",
		"<AI_SKILL_REPO>/ai-tools/agent/copilot.md",
		"Bootstrap Receipt",
		"Before answering any user request",
		"simple file listings",
		"ai-skill runtime validate",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt missing %q: %s", expected, prompt)
		}
	}
}

func TestCopilotStartBlocksMissingProject(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"copilot", "start", "--project", filepath.Join(t.TempDir(), "missing"), "--json"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage, got %d; stderr=%s", code, stderr.String())
	}
}
