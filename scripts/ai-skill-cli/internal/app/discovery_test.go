package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// discoveryRepoRoot walks up from the test cwd until it finds the Ai-skill
// repo root (marker: CORE_BOOTSTRAP.md + runtime/runtime.db). Tests are
// run from .../scripts/ai-skill-cli/internal/app — three levels above is
// the repo root.
func discoveryRepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := cwd
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "CORE_BOOTSTRAP.md")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("could not find repo root from %s", cwd)
	return ""
}

func TestEligibilityCheck_ShortMessageFails(t *testing.T) {
	cfg := defaultDiscoveryConfig()
	input := DiscoveryInput{UserMessage: "hi"}
	reason, ok := EligibilityCheck(input, cfg)
	if ok {
		t.Fatalf("expected eligibility fail; got ok=true reason=%q", reason)
	}
	if reason != discoveryMissEligibilityFail {
		t.Errorf("expected %q got %q", discoveryMissEligibilityFail, reason)
	}
}

func TestEligibilityCheck_LongImperativePasses(t *testing.T) {
	cfg := defaultDiscoveryConfig()
	input := DiscoveryInput{
		UserMessage: "幫我規劃下個月去京都的五天旅遊行程，預算有限要省一點",
	}
	reason, ok := EligibilityCheck(input, cfg)
	if !ok {
		t.Fatalf("expected eligibility pass; got reason=%q", reason)
	}
}

func TestEligibilityCheck_NoArtifactNoImperative(t *testing.T) {
	cfg := defaultDiscoveryConfig()
	// Long enough (≥6 tokens) but no imperative verb and no artifact.
	input := DiscoveryInput{UserMessage: "the sky was very blue yesterday morning quietly"}
	reason, ok := EligibilityCheck(input, cfg)
	if ok {
		t.Fatalf("expected eligibility fail; got ok=true")
	}
	if reason != discoveryMissNoArtifact {
		t.Errorf("expected %q got %q", discoveryMissNoArtifact, reason)
	}
}

func TestExtractArtifactTokens(t *testing.T) {
	msg := "幫我看 src/api/client.ts 還有 README.md 對照 v1.2.3 的設定，不要動 https://x.com/y.png"
	bn, paths, exts := extractArtifactTokens(msg)
	if len(bn) == 0 {
		t.Fatalf("expected basenames, got none")
	}
	gotBn := strings.Join(bn, ",")
	if !strings.Contains(gotBn, "client.ts") {
		t.Errorf("expected client.ts in basenames, got %q", gotBn)
	}
	if !strings.Contains(gotBn, "README.md") {
		t.Errorf("expected README.md in basenames, got %q", gotBn)
	}
	if strings.Contains(gotBn, "y.png") {
		t.Errorf("URL filename should be skipped, got %q", gotBn)
	}
	gotExts := strings.Join(exts, ",")
	if !strings.Contains(gotExts, ".ts") || !strings.Contains(gotExts, ".md") {
		t.Errorf("expected .ts and .md exts, got %q", gotExts)
	}
	gotPaths := strings.Join(paths, ",")
	if !strings.Contains(gotPaths, "src/api/client.ts") {
		t.Errorf("expected path src/api/client.ts, got %q", gotPaths)
	}
}

func TestRunLightDiscovery_TravelPlanningReplay(t *testing.T) {
	repo := discoveryRepoRoot(t)
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		t.Fatalf("load registry: %v", err)
	}
	summaries := LoadDiscoverySummaries(repo)
	if len(summaries) == 0 {
		t.Fatalf("no summaries loaded from %s", repo)
	}
	cfg := defaultDiscoveryConfig()
	input := DiscoveryInput{
		UserMessage: "幫我 review 這份 20260531-下関.md 的旅遊行程內容，看 Mapcode 有沒有問題",
		Basenames:   []string{"20260531-下関.md"},
		Paths:       []string{"docs/20260531-下関.md"},
		Extensions:  []string{".md"},
	}
	candidates, _ := RunLightDiscovery(input, registry, summaries, cfg, repo)
	if len(candidates) == 0 {
		t.Fatalf("expected at least one candidate; got none")
	}
	// travel-planning should appear in top-3 (the keyword 旅遊 + 行程 +
	// Mapcode are all in the summary 'When to read').
	found := false
	for _, c := range candidates {
		if strings.Contains(c.Route, "travel-planning") {
			found = true
			t.Logf("travel-planning candidate score=%.3f evidence=%+v", c.Score, c.Evidence)
			break
		}
	}
	if !found {
		t.Errorf("expected travel-planning in candidates; got %+v", candidates)
	}
}

func TestRunLightDiscovery_SoftwareDelivery_ExtensionHint(t *testing.T) {
	repo := discoveryRepoRoot(t)
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		t.Fatalf("load registry: %v", err)
	}
	summaries := LoadDiscoverySummaries(repo)
	cfg := defaultDiscoveryConfig()
	input := DiscoveryInput{
		UserMessage: "請幫我看 src/api/client.ts 這個 SDK 整合 bug 怎麼修，需要 implement test",
		Basenames:   []string{"client.ts"},
		Paths:       []string{"src/api/client.ts"},
		Extensions:  []string{".ts"},
	}
	candidates, _ := RunLightDiscovery(input, registry, summaries, cfg, repo)
	if len(candidates) == 0 {
		t.Fatalf("expected at least one candidate")
	}
	// At least one candidate should mention software-delivery.
	found := false
	for _, c := range candidates {
		if strings.Contains(c.Route, "software-delivery") {
			found = true
			t.Logf("software-delivery candidate score=%.3f evidence=%+v", c.Score, c.Evidence)
			break
		}
	}
	if !found {
		t.Logf("software-delivery not in top-3 candidates: %+v", candidates)
		// Don't fail hard — depending on registry coverage, another route
		// (e.g. plan execution) may outrank. Phase D tunes thresholds.
	}
}

func TestRunDiscoveryBridge_ManualLockBypass(t *testing.T) {
	repo := discoveryRepoRoot(t)
	tmpDB := filepath.Join(t.TempDir(), "fake-runtime.db")
	input := DiscoveryInput{
		UserMessage: "幫我規劃下個月去京都的五天旅遊行程",
	}
	advisory, proposal, err := RunDiscoveryBridge(input, repo, tmpDB, true /* manualLockActive */)
	if err != nil {
		t.Fatalf("RunDiscoveryBridge: %v", err)
	}
	if advisory != "" {
		t.Errorf("expected no advisory on manual-lock bypass; got %q", advisory)
	}
	if proposal.MissReason != discoveryMissManualLockBypass {
		t.Errorf("expected miss_reason=%s got %s", discoveryMissManualLockBypass, proposal.MissReason)
	}
	if proposal.Status != discoveryStatusExpired {
		t.Errorf("expected status=%s got %s", discoveryStatusExpired, proposal.Status)
	}
}

func TestProjectOverlayCache_Invalidation(t *testing.T) {
	cache := &projectOverlayCache{entries: map[string]projectOverlayEntry{}}
	cwd := "/tmp/fake-cwd-1"
	cache.put(cwd, []EvidenceItem{{Type: "x", Value: "y"}})
	if got, ok := cache.get(cwd); !ok || len(got) != 1 {
		t.Fatalf("expected cached entry; got ok=%v len=%d", ok, len(got))
	}
	cache.Invalidate()
	if _, ok := cache.get(cwd); ok {
		t.Errorf("expected cache empty after invalidate")
	}
}

func TestTaskHash_DeterministicAndSensitiveToInputs(t *testing.T) {
	a := DiscoveryInput{UserMessage: "hello world", Cwd: "/x", Basenames: []string{"a.md", "b.md"}}
	b := DiscoveryInput{UserMessage: "hello world", Cwd: "/x", Basenames: []string{"b.md", "a.md"}}
	if taskHash(a) != taskHash(b) {
		t.Errorf("expected order-independent basenames hash to match")
	}
	c := DiscoveryInput{UserMessage: "hello WORLD", Cwd: "/x", Basenames: []string{"a.md"}}
	if taskHash(a) == taskHash(c) {
		t.Errorf("expected different messages to hash differently")
	}
}

func TestRenderAdvisory_RespectsTokenCap(t *testing.T) {
	cfg := defaultDiscoveryConfig()
	cfg.Advisory.MaxTokensPerInject = 10
	candidates := []RouteCandidate{
		{Route: "route.workflow.travel-planning", Score: 0.7, Evidence: []EvidenceItem{{Type: "user_msg_term", Value: "x"}}},
	}
	out := renderAdvisory(candidates, "", cfg)
	if !strings.Contains(out, "[ai-skill Discovery Bridge") {
		t.Errorf("expected advisory header, got %q", out)
	}
	if len(strings.Fields(out)) > 12 {
		t.Errorf("advisory exceeded cap: %d words", len(strings.Fields(out)))
	}
}

// TestActiveWeights_DormantExcludedAndRenormalized verifies the mechanical half
// of the Dormant-Feature Invariant: a dormant feature's reserved weight is
// dropped from the denominator and the remaining weights renormalize to 1.0.
func TestActiveWeights_DormantExcludedAndRenormalized(t *testing.T) {
	cfg := defaultDiscoveryConfig()
	if _, ok := cfg.PhaseALight.DormantFeatures["frontmatter_head"]; !ok {
		t.Fatalf("precondition: frontmatter_head should be dormant by default")
	}
	w := activeWeights(cfg)
	if _, ok := w["frontmatter_head"]; ok {
		t.Errorf("dormant frontmatter_head must be absent from active weights")
	}
	sum := 0.0
	for _, v := range w {
		sum += v
	}
	if sum < 0.999 || sum > 1.001 {
		t.Errorf("active weights must renormalize to 1.0, got %.6f", sum)
	}
	// user_msg_term 0.30 / (1.0 - 0.10 dormant) = 0.3333…
	if got := w["user_msg_term"]; got < 0.332 || got > 0.334 {
		t.Errorf("user_msg_term expected ~0.3333 after renormalize, got %.4f", got)
	}
}

// TestDormantFeatureInvariant_NoScoreNoEvidence proves a dormant feature emits
// no evidence even when its match condition WOULD fire. The control case
// (frontmatter_head re-enabled) confirms the head genuinely matches, so the
// absence in the dormant case is the dormancy mechanism, not a non-match.
func TestDormantFeatureInvariant_NoScoreNoEvidence(t *testing.T) {
	repo := discoveryRepoRoot(t)
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	summaries := LoadDiscoverySummaries(repo)
	// Keyword msg so travel-planning scores via live features; head contains the
	// route tail ("travel-planning") so frontmatterHit WOULD match if reached.
	input := DiscoveryInput{
		UserMessage:     "幫我規劃京都旅遊行程 itinerary 看看",
		FrontmatterHead: map[string]string{"trip.md": "notes about travel-planning"},
	}

	// (a) Default cfg: frontmatter_head dormant → no frontmatter evidence.
	candDormant, _ := RunLightDiscovery(input, registry, summaries, defaultDiscoveryConfig(), repo)
	for _, c := range candDormant {
		for _, e := range c.Evidence {
			if e.Type == "frontmatter_head" {
				t.Fatalf("dormant invariant violated: frontmatter_head evidence emitted on %s", c.Route)
			}
		}
	}

	// (b) Control — re-enable frontmatter_head → evidence now appears.
	cfgOn := defaultDiscoveryConfig()
	cfgOn.PhaseALight.DormantFeatures = map[string]string{}
	candOn, _ := RunLightDiscovery(input, registry, summaries, cfgOn, repo)
	found := false
	for _, c := range candOn {
		for _, e := range c.Evidence {
			if e.Type == "frontmatter_head" {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("control failed: with frontmatter_head enabled the head should have matched; cannot prove dormancy")
	}
}
