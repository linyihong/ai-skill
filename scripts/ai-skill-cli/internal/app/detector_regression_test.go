package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// repoRootForDetector resolves the Ai-skill repo root from the test's working
// directory (internal/app -> ../../../.. ). Skips when the real registry is not
// present (e.g. isolated checkout) so the unit suite still passes.
func repoRootForDetector(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "..", "..", "..", ".."))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "knowledge", "runtime", "routing-registry.yaml")); err != nil {
		t.Skipf("real routing-registry not found at %s: %v", root, err)
	}
	return root
}

// TestDetectorTravelPlanningRegression is the mechanical regression for the
// 2026-05-31 failure (validation/scenarios/runtime/
// workflow-detector-travel-planning-regression-v1.yaml): a travel-planning task
// MUST activate route.workflow.travel-planning against the REAL registry.
func TestDetectorTravelPlanningRegression(t *testing.T) {
	root := repoRootForDetector(t)
	registry, err := readRuntimeRoutingRegistry(filepath.Join(root, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	ctx := BuildRuntimeContext(registry,
		[]DetectorMessage{{Role: "user", Text: "幫我規劃下個月去京都的五天旅遊行程，預算有限"}},
		nil, time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC))

	if ctx.Status != StatusDetected {
		t.Fatalf("regression: travel task must be detected, got status=%q (the 2026-05-31 bug was no-match)", ctx.Status)
	}
	if ctx.ActiveRoute != "route.workflow.travel-planning" {
		t.Fatalf("regression: expected route.workflow.travel-planning, got active=%q detected=%v", ctx.ActiveRoute, ctx.DetectedRoutes)
	}
	if ctx.Conflict {
		t.Fatalf("regression: travel task should be a clean single match, got conflict=true (detected=%v)", ctx.DetectedRoutes)
	}
}

// TestDetectorDeterministicMatchAgainstRealRegistry backs
// workflow-detector-deterministic-match-v1.yaml.
func TestDetectorDeterministicMatchAgainstRealRegistry(t *testing.T) {
	root := repoRootForDetector(t)
	registry, err := readRuntimeRoutingRegistry(filepath.Join(root, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	msg := []DetectorMessage{{Role: "user", Text: "幫我規劃一個四天三夜的東京旅遊行程"}}
	a := BuildRuntimeContext(registry, msg, nil, time.Now().UTC())
	b := BuildRuntimeContext(registry, msg, nil, time.Now().UTC())
	if a.ActiveRoute != "route.workflow.travel-planning" {
		t.Fatalf("expected travel-planning, got %q", a.ActiveRoute)
	}
	// deterministic: same input → same active route + detected set
	if a.ActiveRoute != b.ActiveRoute || len(a.DetectedRoutes) != len(b.DetectedRoutes) {
		t.Fatalf("non-deterministic: a=%+v b=%+v", a, b)
	}
}
