package app

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateRuntimeTriggerWiringOptOut(t *testing.T) {
	body := "feat(x): touch routing-registry\n\n[skip-runtime-trigger-wiring]\n"
	got := validateRuntimeTriggerWiring(body, []string{"knowledge/runtime/routing-registry.yaml"}, "")
	if got != "" {
		t.Errorf("opt-out should bypass, got %q", got)
	}
}

func TestValidateRuntimeTriggerWiringNoTrigger(t *testing.T) {
	got := validateRuntimeTriggerWiring("feat(x): unrelated", []string{"README.md"}, "")
	if got != "" {
		t.Errorf("non-runtime staging should bypass, got %q", got)
	}
}

func TestValidateRuntimeTriggerWiringBlocksOrphanRoute(t *testing.T) {
	repo := initTempGitRepo(t)
	rel := "knowledge/runtime/routing-registry.yaml"
	abs := filepath.Join(repo, rel)
	writeFile(t, abs, "records:\n  - id: route.existing\n    primary_source: foo.md\n")
	runGit(t, repo, "add", rel)
	runGit(t, repo, "commit", "-m", "init")
	writeFile(t, abs, "records:\n  - id: route.existing\n    primary_source: foo.md\n  - id: route.totally.unwired\n    primary_source: bar.md\n")
	runGit(t, repo, "add", rel)

	got := validateRuntimeTriggerWiring("feat(x): add route", []string{rel}, repo)
	if !strings.Contains(got, "route.totally.unwired") {
		t.Errorf("expected violation citing new orphan route, got %q", got)
	}
}

func TestValidateRuntimeTriggerWiringAcceptsManualAnnotation(t *testing.T) {
	repo := initTempGitRepo(t)
	rel := "knowledge/runtime/routing-registry.yaml"
	abs := filepath.Join(repo, rel)
	writeFile(t, abs, "records:\n  - id: route.existing\n    primary_source: foo.md\n")
	runGit(t, repo, "add", rel)
	runGit(t, repo, "commit", "-m", "init")
	writeFile(t, abs, "records:\n  - id: route.existing\n    primary_source: foo.md\n  - id: route.manual.thing\n    manual_activation:\n      reason: workflow_discovery\n    primary_source: bar.md\n")
	runGit(t, repo, "add", rel)

	got := validateRuntimeTriggerWiring("feat(x): add manual route", []string{rel}, repo)
	if got != "" {
		t.Errorf("manual_activation annotation should pass, got %q", got)
	}
}

func TestValidateRuntimeTriggerWiringAcceptsDiscoverySignalReference(t *testing.T) {
	repo := initTempGitRepo(t)
	regRel := "knowledge/runtime/routing-registry.yaml"
	discRel := "runtime/cognitive-modes-discovery.yaml"
	writeFile(t, filepath.Join(repo, regRel), "records:\n  - id: route.existing\n    primary_source: foo.md\n")
	writeFile(t, filepath.Join(repo, discRel), "signals:\n  - name: existing\n    description: nothing\n")
	runGit(t, repo, "add", regRel, discRel)
	runGit(t, repo, "commit", "-m", "init")
	// Add new route AND a discovery signal referencing it
	writeFile(t, filepath.Join(repo, regRel), "records:\n  - id: route.existing\n    primary_source: foo.md\n  - id: route.signal.linked\n    primary_source: bar.md\n")
	writeFile(t, filepath.Join(repo, discRel), "signals:\n  - name: existing\n    description: nothing\n  - name: new_one\n    description: loads route.signal.linked when relevant\n")
	runGit(t, repo, "add", regRel, discRel)

	got := validateRuntimeTriggerWiring("feat(x): add signal-linked route", []string{regRel, discRel}, repo)
	if got != "" {
		t.Errorf("signal mention should satisfy wiring, got %q", got)
	}
}

func TestValidateRuntimeTriggerWiringBlocksOrphanTargetKey(t *testing.T) {
	repo := initTempGitRepo(t)
	rel := "runtime/foo.yaml"
	abs := filepath.Join(repo, rel)
	writeFile(t, abs, "id: foo\n")
	runGit(t, repo, "add", rel)
	runGit(t, repo, "commit", "-m", "init")
	writeFile(t, abs, "id: foo\nruntime_projection:\n  enabled: true\n  target_key: runtime.lonely.surface\n")
	runGit(t, repo, "add", rel)

	got := validateRuntimeTriggerWiring("feat(x): add target_key", []string{rel}, repo)
	if !strings.Contains(got, "runtime.lonely.surface") {
		t.Errorf("expected violation citing new orphan target_key, got %q", got)
	}
}

func TestStagedAddedRouteIDsIgnoresContextLines(t *testing.T) {
	repo := initTempGitRepo(t)
	rel := "knowledge/runtime/routing-registry.yaml"
	abs := filepath.Join(repo, rel)
	writeFile(t, abs, "records:\n  - id: route.one\n  - id: route.two\n")
	runGit(t, repo, "add", rel)
	runGit(t, repo, "commit", "-m", "init")
	writeFile(t, abs, "records:\n  - id: route.one\n  - id: route.two\n  - id: route.three\n")
	runGit(t, repo, "add", rel)

	ids := stagedAddedRouteIDs(repo, rel)
	if len(ids) != 1 || ids[0] != "route.three" {
		t.Errorf("expected only route.three to be detected as added, got %v", ids)
	}
}

func TestStagedAddedTargetKeysHandlesMultiple(t *testing.T) {
	repo := initTempGitRepo(t)
	rel := "runtime/foo.yaml"
	abs := filepath.Join(repo, rel)
	writeFile(t, abs, "x: 1\n")
	runGit(t, repo, "add", rel)
	runGit(t, repo, "commit", "-m", "init")
	writeFile(t, abs, "x: 1\nruntime_projection:\n  target_key: runtime.alpha\nother:\n  target_key: runtime.beta\n")
	runGit(t, repo, "add", rel)

	keys := stagedAddedTargetKeys(repo, rel)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys added, got %v", keys)
	}
	if keys[0] != "runtime.alpha" || keys[1] != "runtime.beta" {
		t.Errorf("unexpected target keys: %v", keys)
	}
}
