package app

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestRuntimeValidateDryRunPlansValidators(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "validate", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "runtime validate" || result.Mode != "dry_run" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.PlannedActions) != 3 {
		t.Fatalf("expected three planned validators, got %#v", result.PlannedActions)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime validate dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeValidateBlocksMissingRubyBeforeWrapper(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "validate", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_ruby" {
		t.Fatalf("expected missing_ruby, got %#v", result.Error)
	}
}

func TestRuntimeRefreshDryRunPlansWrapperCommands(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "runtime refresh" || result.Mode != "dry_run" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.PlannedActions) != 7 {
		t.Fatalf("expected seven planned refresh scripts, got %#v", result.PlannedActions)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime refresh dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeRefreshBlocksMissingRubyBeforeWrapper(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_ruby" {
		t.Fatalf("expected missing_ruby, got %#v", result.Error)
	}
}

func TestRuntimeValidateBlocksMissingValidator(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "scripts", "validate-knowledge-runtime.rb"), "# ok\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "validate", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s", code, stderr.String())
	}
}

func TestRuntimeUnsupportedSubcommandReturnsInvalidUsage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "compile"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unsupported runtime command") {
		t.Fatalf("expected unsupported message, got %q", stderr.String())
	}
}

func TestRuntimeWrapperEnvForcesUTF8(t *testing.T) {
	env := runtimeWrapperEnv([]string{"PATH=/bin", "LANG=C"})
	if !containsEnv(env, "LANG=C.UTF-8") {
		t.Fatalf("expected LANG override, got %#v", env)
	}
	if !containsEnv(env, "LC_ALL=C.UTF-8") {
		t.Fatalf("expected LC_ALL override, got %#v", env)
	}
}

func fakeRuntimeRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	for _, name := range []string{
		"generate-model-context-report.rb",
		"generate-model-checklists.rb",
		"generate-knowledge-runtime-report.rb",
		"generate-runtime-sqlite-index.rb",
		"refresh-knowledge-runtime.rb",
		"validate-knowledge-runtime.rb",
		"validate-runtime-db.rb",
		"validate-runtime-sqlite-index.rb",
	} {
		writeFile(t, filepath.Join(repo, "scripts", name), "#!/usr/bin/env ruby\nputs 'ok'\n")
	}
	return repo
}

func containsEnv(env []string, item string) bool {
	for _, value := range env {
		if value == item {
			return true
		}
	}
	return false
}
