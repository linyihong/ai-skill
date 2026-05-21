package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDoctorRequireGitMissingBlocksWithStableJSON(t *testing.T) {
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"doctor", "--require-git", "--json"}, &stdout, &stderr)

	if code != ExitMissingDependency {
		t.Fatalf("expected exit %d, got %d; stderr=%s", ExitMissingDependency, code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v\n%s", err, stdout.String())
	}

	if result.Status != "blocked" {
		t.Fatalf("expected blocked status, got %q", result.Status)
	}
	if result.Error == nil || result.Error.Code != "missing_git" {
		t.Fatalf("expected missing_git error, got %#v", result.Error)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("doctor must not mutate files, got %#v", result.Mutations)
	}
}

func TestDoctorWithoutRequireGitReportsMissingButSucceeds(t *testing.T) {
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"doctor", "--json"}, &stdout, &stderr)

	if code != ExitSuccess {
		t.Fatalf("expected success exit, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Status != "success" {
		t.Fatalf("expected success status, got %q", result.Status)
	}
	if !hasCheckStatus(result.Checks, "git", "missing") {
		t.Fatalf("expected git missing check, got %#v", result.Checks)
	}
}

func TestUnknownCommandReturnsInvalidUsage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"nope"}, &stdout, &stderr)

	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("expected unknown command message, got %q", stderr.String())
	}
}

func TestDoctorPlainOutput(t *testing.T) {
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"doctor", "--require-git"}, &stdout, &stderr)

	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d", code)
	}
	if !strings.Contains(stdout.String(), "missing") {
		t.Fatalf("expected plain output to mention missing check, got %q", stdout.String())
	}
}

func emptyPathDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if runtime.GOOS == "windows" {
		// Windows may search the current directory before PATH. Run tests from an
		// empty temp directory so a repository-local executable cannot be found.
		previous, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = os.Chdir(previous)
		})
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Clean(dir)
}

func hasCheckStatus(checks []Check, name string, status string) bool {
	for _, check := range checks {
		if check.Name == name && check.Status == status {
			return true
		}
	}
	return false
}
