package app

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestSyncCursorBundleDryRunPlansFakeTargetWithoutWriting(t *testing.T) {
	repo := fakeSyncRepo(t)
	target := filepath.Join(t.TempDir(), ".cursor")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"sync-cursor-bundle", "--repo", repo, "--target", target, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "sync-cursor-bundle" {
		t.Fatalf("unexpected command: %q", result.Command)
	}
	if !hasCheckStatus(result.Checks, "target_boundary", "ok") {
		t.Fatalf("expected target boundary ok, got %#v", result.Checks)
	}
	if !hasCheckMessage(result.Checks, "mirror_strategy", "copy-fallback") {
		t.Fatalf("expected copy-fallback strategy, got %#v", result.Checks)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("dry-run must not mutate, got %#v", result.Mutations)
	}
	if pathExists(filepath.Join(target, "bundles")) {
		t.Fatal("dry-run wrote target bundles directory")
	}
}

func TestSyncCursorBundleRequiresExplicitTarget(t *testing.T) {
	repo := fakeSyncRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"sync-cursor-bundle", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage, got %d; stderr=%s", code, stderr.String())
	}
}

func TestSyncCursorBundleBlocksTargetInsideRepo(t *testing.T) {
	repo := fakeSyncRepo(t)
	target := filepath.Join(repo, ".cursor")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"sync-cursor-bundle", "--repo", repo, "--target", target, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "target_inside_repo" {
		t.Fatalf("expected target_inside_repo, got %#v", result.Error)
	}
}

func TestSyncCursorBundleBlocksMissingSource(t *testing.T) {
	repo := t.TempDir()
	target := filepath.Join(t.TempDir(), ".cursor")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"sync-cursor-bundle", "--repo", repo, "--target", target, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s", code, stderr.String())
	}
}

func TestSyncCursorBundlePlansSkillMirrors(t *testing.T) {
	repo := fakeSyncRepo(t)
	writeFile(t, filepath.Join(repo, "skills", "demo", "SKILL.md"), "# Demo\n")
	writeFile(t, filepath.Join(repo, "skills", "_template", "SKILL.md"), "# Template\n")
	target := filepath.Join(t.TempDir(), ".cursor")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"sync-cursor-bundle", "--repo", repo, "--target", target, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckMessage(result.Checks, "skills", "1 syncable skills") {
		t.Fatalf("expected one syncable skill, got %#v", result.Checks)
	}
}

func TestSyncCursorBundleWriteModeBlockedUntilParity(t *testing.T) {
	repo := fakeSyncRepo(t)
	target := filepath.Join(t.TempDir(), ".cursor")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"sync-cursor-bundle", "--repo", repo, "--target", target, "--json"}, &stdout, &stderr)
	if code != ExitPartialCloseBlocked {
		t.Fatalf("expected write mode blocked, got %d; stderr=%s", code, stderr.String())
	}
	if pathExists(filepath.Join(target, "bundles")) {
		t.Fatal("write-blocked mode wrote target bundles directory")
	}
}

func fakeSyncRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "enforcement", "README.md"), "# Enforcement\n")
	return repo
}
