package app

import (
	"reflect"
	"testing"
	"time"
)

var rcClock = time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

func rcMsg(text string) []DetectorMessage { return []DetectorMessage{{Role: "user", Text: text}} }

func hasRoute(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}

func TestBuildRuntimeContext_SingleActivation(t *testing.T) {
	ctx := BuildRuntimeContext(testRegistry(), rcMsg("幫我做 web scraping"), nil, rcClock)
	if ctx.Status != StatusDetected {
		t.Fatalf("status = %q, want detected", ctx.Status)
	}
	if ctx.ActiveRoute != "route.analysis.web" {
		t.Fatalf("active = %q", ctx.ActiveRoute)
	}
	if ctx.Conflict {
		t.Fatal("unexpected conflict")
	}
	if ctx.EffectiveMode != ModeAutoDetect {
		t.Fatalf("mode = %q", ctx.EffectiveMode)
	}
	if !ctx.Substantive {
		t.Fatal("expected substantive")
	}
}

func TestBuildRuntimeContext_ConflictNoAutoPick(t *testing.T) {
	ctx := BuildRuntimeContext(testRegistry(), rcMsg("用 Frida 抓這個 APK 並評估 DDD 架構"), nil, rcClock)
	if !ctx.Conflict {
		t.Fatal("expected conflict")
	}
	if ctx.ActiveRoute != "" {
		t.Fatalf("conflict must NOT auto-pick, got %q", ctx.ActiveRoute)
	}
	if ctx.Status != StatusDetected {
		t.Fatalf("status = %q", ctx.Status)
	}
	want := []string{"route.intelligence.architectural-fit", "route.workflow.apk-analysis"}
	if !reflect.DeepEqual(ctx.DetectedRoutes, want) {
		t.Fatalf("detected = %v, want %v", ctx.DetectedRoutes, want)
	}
}

func TestBuildRuntimeContext_NoMatchOnGreeting(t *testing.T) {
	ctx := BuildRuntimeContext(testRegistry(), rcMsg("hi 早安"), nil, rcClock)
	if ctx.Status != StatusNoMatch {
		t.Fatalf("status = %q, want no-match", ctx.Status)
	}
	if ctx.Substantive {
		t.Fatal("greeting must not be substantive")
	}
	if ctx.ActiveRoute != "" {
		t.Fatalf("active = %q", ctx.ActiveRoute)
	}
}

func TestBuildRuntimeContext_ManualLock(t *testing.T) {
	// "鎖定" + a single route's signal ("web scraping") → manual-lock
	ctx := BuildRuntimeContext(testRegistry(),
		rcMsg("這個專案之後都用 web scraping，鎖定"), nil, rcClock)
	if ctx.Status != StatusLocked {
		t.Fatalf("status = %q, want locked", ctx.Status)
	}
	if ctx.EffectiveMode != ModeManualLock {
		t.Fatalf("mode = %q, want manual-lock", ctx.EffectiveMode)
	}
	if ctx.ActiveRoute != "route.analysis.web" {
		t.Fatalf("active = %q", ctx.ActiveRoute)
	}
}

func TestBuildRuntimeContext_ManualLockAmbiguousDoesNotLock(t *testing.T) {
	// lock sentinel but TWO routes' signals present → cannot resolve → no lock,
	// falls through to normal detection (which is then a conflict).
	ctx := BuildRuntimeContext(testRegistry(),
		rcMsg("鎖定：web scraping 跟 DDD 都要"), nil, rcClock)
	if ctx.Status == StatusLocked {
		t.Fatal("ambiguous lock must NOT lock")
	}
}

func TestBuildRuntimeContext_ManualUnlockRestoresAutoDetect(t *testing.T) {
	// lock earlier, unlock later (unlock wins) → auto-detection resumes
	transcript := []DetectorMessage{
		{Role: "user", Text: "鎖定 web scraping"},
		{Role: "user", Text: "回到自動偵測，現在我要 DDD 架構評估"},
	}
	ctx := BuildRuntimeContext(testRegistry(), transcript, nil, rcClock)
	if ctx.Status == StatusLocked {
		t.Fatal("unlock should clear manual-lock")
	}
	if ctx.ActiveRoute != "route.intelligence.architectural-fit" {
		t.Fatalf("after unlock expected auto-detect architectural-fit, got %q", ctx.ActiveRoute)
	}
}

func TestBuildRuntimeContext_ExplicitPivotReDetectsPostPivot(t *testing.T) {
	// First turn = apk; pivot turn switches to web. Only post-pivot considered.
	transcript := []DetectorMessage{
		{Role: "user", Text: "用 Frida 抓 APK"},
		{Role: "user", Text: "換任務，現在做 web scraping"},
	}
	ctx := BuildRuntimeContext(testRegistry(), transcript, nil, rcClock)
	if ctx.ActiveRoute != "route.analysis.web" {
		t.Fatalf("post-pivot active = %q, want analysis.web", ctx.ActiveRoute)
	}
	// apk must NOT linger from the pre-pivot turn
	for _, r := range ctx.DetectedRoutes {
		if r == "route.workflow.apk-analysis" {
			t.Fatal("pre-pivot route must not survive explicit pivot")
		}
	}
}

func TestBuildRuntimeContext_NoImplicitDriftInvalidation(t *testing.T) {
	// Drill-down: the workflow keyword appears once, then several sub-questions
	// with NO further trigger keyword. Detection still holds from the early turn
	// because the whole (post-pivot) transcript is matched — no implicit drift.
	transcript := []DetectorMessage{
		{Role: "user", Text: "幫我做 web scraping"},
		{Role: "user", Text: "這個欄位怎麼解析"},
		{Role: "user", Text: "那分頁呢"},
		{Role: "user", Text: "速度太慢怎麼辦"},
	}
	ctx := BuildRuntimeContext(testRegistry(), transcript, nil, rcClock)
	if ctx.ActiveRoute != "route.analysis.web" {
		t.Fatalf("drill-down should NOT lose the route, got %q", ctx.ActiveRoute)
	}
}

func TestBuildRuntimeContext_ReinforcementOnlyStaysNoMatch(t *testing.T) {
	// Only artifact content matches → late-detected; not an activation.
	ctx := BuildRuntimeContext(testRegistry(), nil,
		[]DetectorFile{{Path: "n.md", Content: "anti-bot gateway present"}}, rcClock)
	if ctx.Status != StatusNoMatch {
		t.Fatalf("reinforcement-only must be no-match (not activated), got %q", ctx.Status)
	}
	// but the candidate is surfaced in DetectedRoutes for coverage analysis
	found := false
	for _, r := range ctx.DetectedRoutes {
		if r == "route.analysis.web" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected late-detected candidate in DetectedRoutes")
	}
}

// TestActivationContractMatrix pins the 3-case can_activate matrix so the
// advisory contract violation can never regress:
//
//	Case 1  auto-detect only        -> active_route set
//	Case 2  advisory only           -> active_route nil, advisory in DetectedRoutes only
//	Case 3  auto-detect + advisory  -> active_route = auto-detect; advisory suggestion only
func TestActivationContractMatrix(t *testing.T) {
	reg := testRegistry()

	// Case 1 — auto-detect hit locks.
	c1 := BuildRuntimeContext(reg, rcMsg("幫我做 web scraping"), nil, rcClock)
	if c1.ActiveRoute != "route.analysis.web" || c1.EffectiveMode != ModeAutoDetect {
		t.Fatalf("Case1: auto-detect must lock, got active=%q mode=%q", c1.ActiveRoute, c1.EffectiveMode)
	}

	// Case 2 — advisory only never locks.
	c2 := BuildRuntimeContext(reg, rcMsg("看一下 magic bytes"), nil, rcClock)
	if c2.ActiveRoute != "" {
		t.Fatalf("Case2: advisory-only must NOT lock, got %q", c2.ActiveRoute)
	}
	if !hasRoute(c2.DetectedRoutes, "route.intelligence.engineering.heuristics") || len(c2.CandidateRoutes) != 0 {
		t.Fatalf("Case2: advisory must be a detected suggestion, not a candidate; detected=%v candidates=%v", c2.DetectedRoutes, c2.CandidateRoutes)
	}

	// Case 3 — auto-detect wins over co-occurring advisory.
	c3 := BuildRuntimeContext(reg, rcMsg("評估 DDD 架構，順便看 magic bytes"), nil, rcClock)
	if c3.ActiveRoute != "route.intelligence.architectural-fit" {
		t.Fatalf("Case3: auto-detect must win, got %q", c3.ActiveRoute)
	}
	if !hasRoute(c3.DetectedRoutes, "route.intelligence.engineering.heuristics") {
		t.Fatalf("Case3: advisory must remain a detected suggestion, got %v", c3.DetectedRoutes)
	}
}

// TestAdvisoryNeverLocksActiveRoute is the contract guard for
// activation_mode_spec `advisory.can_activate: false`. An advisory-only hit
// must surface as a suggestion (DetectedRoutes) but NEVER lock ActiveRoute —
// otherwise it would falsely trigger the Phase 5 primary_source gate.
func TestAdvisoryNeverLocksActiveRoute(t *testing.T) {
	// "magic bytes" hits only route.intelligence.engineering.heuristics (advisory)
	ctx := BuildRuntimeContext(testRegistry(), rcMsg("看一下 magic bytes 怎麼判斷"), nil, rcClock)
	if ctx.ActiveRoute != "" {
		t.Fatalf("advisory route must NOT lock active_route, got %q (mode=%s)", ctx.ActiveRoute, ctx.EffectiveMode)
	}
	if ctx.Status != StatusNoMatch {
		t.Fatalf("advisory-only must be no-match (no lock), got %q", ctx.Status)
	}
	// still surfaced as a suggestion, but not a candidate
	if !hasRoute(ctx.DetectedRoutes, "route.intelligence.engineering.heuristics") {
		t.Fatalf("advisory route should appear in DetectedRoutes as suggestion, got %v", ctx.DetectedRoutes)
	}
	if len(ctx.CandidateRoutes) != 0 {
		t.Fatalf("advisory route must NOT be a candidate, got %v", ctx.CandidateRoutes)
	}
}

// TestAutoDetectPlusAdvisory_AdvisoryIsSuggestionOnly: when both an auto-detect
// and an advisory route match, the auto-detect one locks and the advisory one
// is suggestion-only.
func TestAutoDetectPlusAdvisory_AdvisoryIsSuggestionOnly(t *testing.T) {
	// "DDD" → architectural-fit (auto-detect); "magic bytes" → heuristics (advisory)
	ctx := BuildRuntimeContext(testRegistry(), rcMsg("評估 DDD 架構，順便看 magic bytes"), nil, rcClock)
	if ctx.ActiveRoute != "route.intelligence.architectural-fit" {
		t.Fatalf("auto-detect route must win the lock, got %q", ctx.ActiveRoute)
	}
	if ctx.Conflict {
		t.Fatal("one auto-detect + one advisory is NOT a conflict (only 1 candidate)")
	}
	if !reflect.DeepEqual(ctx.CandidateRoutes, []string{"route.intelligence.architectural-fit"}) {
		t.Fatalf("candidates must be the auto-detect route only, got %v", ctx.CandidateRoutes)
	}
	if !hasRoute(ctx.DetectedRoutes, "route.intelligence.engineering.heuristics") {
		t.Fatalf("advisory route should still be a detected suggestion, got %v", ctx.DetectedRoutes)
	}
}

func TestIsSubstantive(t *testing.T) {
	reg := testRegistry()
	if !IsSubstantive(reg, "幫我看一下") { // action verb 幫我
		t.Fatal("action verb should be substantive")
	}
	if !IsSubstantive(reg, "web scraping 怎麼做") { // domain noun from registry
		t.Fatal("domain noun should be substantive")
	}
	if IsSubstantive(reg, "hi 早安") {
		t.Fatal("greeting should not be substantive")
	}
}
