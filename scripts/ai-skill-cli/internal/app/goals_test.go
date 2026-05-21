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

func TestGoalsInitWriteModeCreatesLedgerAndGitExclude(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".git", "info"), 0o755); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"goals", "init", "--project", project, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected write success, got %d; stderr=%s", code, stderr.String())
	}
	if !pathExists(filepath.Join(project, ".agent-goals", "goals")) {
		t.Fatal("goals dir was not created")
	}
	exclude, err := os.ReadFile(filepath.Join(project, ".git", "info", "exclude"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(exclude, []byte(".agent-goals/")) {
		t.Fatalf("git exclude missing .agent-goals/: %s", string(exclude))
	}
}

func TestGoalsLifecycleWriteMode(t *testing.T) {
	project := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"goals", "start", "--project", project, "--id", "demo", "--title", "Demo", "--source", "user request", "--plan", "plan.md", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("start failed: %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	goal := filepath.Join(project, ".agent-goals", "goals", "demo.md")
	if !pathExists(goal) {
		t.Fatal("goal file was not created")
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"goals", "update", "--project", project, "--id", "demo", "--next", "Ship it", "--missing", "none", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("update failed: %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	content, err := os.ReadFile(goal)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(content, []byte("Ship it")) {
		t.Fatalf("goal update missing next action: %s", string(content))
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"goals", "complete", "--project", project, "--id", "demo", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("complete without validation failed: %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !pathExists(goal) {
		t.Fatal("goal should remain without --validated")
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"goals", "complete", "--project", project, "--id", "demo", "--validated", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("complete validated failed: %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if pathExists(goal) {
		t.Fatal("validated complete should delete goal")
	}
}

func TestGoalsActiveLockBlocksWrite(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".agent-goals", "goals"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(project, ".agent-goals", "locks", "demo.lock"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".agent-goals", "goals", "demo.md"), []byte("---\nid: demo\n---\n# Demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"goals", "update", "--project", project, "--id", "demo", "--note", "blocked", "--json"}, &stdout, &stderr)
	if code != ExitUnsafeRepoState {
		t.Fatalf("expected active lock block, got %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
}
