package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGoalsStatusReportsMissingLedger(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"goals", "status", "--project", project, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckStatus(result.Checks, "goal_ledger", "missing") {
		t.Fatalf("expected missing ledger check, got %#v", result.Checks)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("status must not mutate, got %#v", result.Mutations)
	}
}

func TestGoalsStatusReportsExistingLedger(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".agent-goals", "goals"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(project, ".agent-goals", "locks"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".agent-goals", "goals", "demo.md"), []byte("# Demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"goals", "status", "--project", project, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckStatus(result.Checks, "goal_files", "ok") {
		t.Fatalf("expected goal_files ok check, got %#v", result.Checks)
	}
}

func TestGoalsInitDryRunPlansWithoutWriting(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"goals", "init", "--project", project, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if len(result.PlannedActions) != 4 {
		t.Fatalf("expected 4 planned actions, got %#v", result.PlannedActions)
	}
	if pathExists(filepath.Join(project, ".agent-goals")) {
		t.Fatal("dry-run wrote .agent-goals")
	}
}

func TestGoalsInitWriteModeBlockedUntilParity(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"goals", "init", "--project", project, "--json"}, &stdout, &stderr)
	if code != ExitPartialCloseBlocked {
		t.Fatalf("expected blocked write mode, got %d; stderr=%s", code, stderr.String())
	}
}
