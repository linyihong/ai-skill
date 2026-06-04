package app

import (
	"reflect"
	"testing"
	"time"
)

var rcClock = time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

func rcMsg(text string) []DetectorMessage { return []DetectorMessage{{Role: "user", Text: text}} }

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
