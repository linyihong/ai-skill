package app

import (
	"reflect"
	"testing"
)

// testRegistry builds a small in-memory registry covering the cases the
// detector must handle, independent of the real routing-registry.yaml.
func testRegistry() runtimeRoutingRegistry {
	return runtimeRoutingRegistry{Records: []runtimeRouteRecord{
		{
			// two-phase auto-detect route
			ID:             "route.analysis.web",
			RouteType:      "analysis",
			ActivationMode: "auto-detect",
			ActivationTriggers: runtimeRouteTriggers{
				ActivationAnyOf: &runtimeActivationAnyOf{
					UserSignals:    []string{"web scraping", "爬蟲"},
					ContextSignals: []string{"**/analysis/web/**", "**/*.html"},
				},
				ReinforcementAnyOf: &runtimeReinforcementAnyOf{
					ArtifactSignals: []string{"anti-bot", "selector|xpath"},
				},
			},
		},
		{
			// another two-phase auto-detect route (for multi-hit)
			ID:             "route.intelligence.architectural-fit",
			RouteType:      "intelligence",
			ActivationMode: "auto-detect",
			ActivationTriggers: runtimeRouteTriggers{
				ActivationAnyOf: &runtimeActivationAnyOf{
					UserSignals: []string{"DDD", "bounded context"},
				},
			},
		},
		{
			// legacy flat route (no activation_any_of); must normalize
			ID:             "route.workflow.apk-analysis",
			RouteType:      "workflow",
			ActivationMode: "auto-detect",
			ActivationTriggers: runtimeRouteTriggers{
				UserSignals:     []string{"Frida", "抓包"},
				FileChangeGlobs: []string{"**/scripts/frida/**"},
			},
		},
		{
			// advisory route — reinforcement-capable, weak signals
			ID:             "route.intelligence.engineering.heuristics",
			RouteType:      "intelligence",
			ActivationMode: "advisory",
			ActivationTriggers: runtimeRouteTriggers{
				ActivationAnyOf: &runtimeActivationAnyOf{
					UserSignals: []string{"magic bytes"},
				},
			},
		},
		{
			// on-demand route — must NEVER flow through the detector even if it
			// somehow carries signals
			ID:             "route.intelligence.engineering.agent-architecture",
			RouteType:      "intelligence",
			ActivationMode: "on-demand",
			ActivationTriggers: runtimeRouteTriggers{
				ActivationAnyOf: &runtimeActivationAnyOf{
					UserSignals: []string{"agent architecture"},
				},
			},
		},
		{
			// always-on route, no triggers — ignored
			ID:        "route.runtime.phase-machine",
			RouteType: "runtime_core",
		},
	}}
}

func detectedIDs(routes []DetectedRoute) []string {
	ids := make([]string, 0, len(routes))
	for _, r := range routes {
		ids = append(ids, r.RouteID)
	}
	return ids
}

func findRoute(routes []DetectedRoute, id string) (DetectedRoute, bool) {
	for _, r := range routes {
		if r.RouteID == id {
			return r, true
		}
	}
	return DetectedRoute{}, false
}

func TestDetectWorkflows_SingleHit_UserSignal(t *testing.T) {
	got := DetectWorkflows(testRegistry(),
		[]DetectorMessage{{Role: "user", Text: "幫我做 web scraping 抓一個網站"}}, nil)
	if ids := detectedIDs(got); !reflect.DeepEqual(ids, []string{"route.analysis.web"}) {
		t.Fatalf("expected only route.analysis.web, got %v", ids)
	}
	r, _ := findRoute(got, "route.analysis.web")
	if !r.Activated {
		t.Fatalf("expected Activated=true for user-signal hit")
	}
	if !reflect.DeepEqual(r.UserSignalHits, []string{"web scraping"}) {
		t.Fatalf("unexpected user hits: %v", r.UserSignalHits)
	}
}

func TestDetectWorkflows_CaseInsensitive(t *testing.T) {
	// "ddd" lower-case in transcript should match signal "DDD"
	got := DetectWorkflows(testRegistry(),
		[]DetectorMessage{{Role: "user", Text: "should I adopt ddd here?"}}, nil)
	if ids := detectedIDs(got); !reflect.DeepEqual(ids, []string{"route.intelligence.architectural-fit"}) {
		t.Fatalf("expected architectural-fit, got %v", ids)
	}
}

func TestDetectWorkflows_MultiHit_Conflict(t *testing.T) {
	got := DetectWorkflows(testRegistry(),
		[]DetectorMessage{{Role: "user", Text: "對這個 APK 用 Frida 抓包，順便評估 DDD 架構"}}, nil)
	ids := detectedIDs(got)
	want := []string{"route.intelligence.architectural-fit", "route.workflow.apk-analysis"}
	if !reflect.DeepEqual(ids, want) {
		t.Fatalf("expected multi-hit %v (sorted), got %v", want, ids)
	}
}

func TestDetectWorkflows_NoMatch(t *testing.T) {
	got := DetectWorkflows(testRegistry(),
		[]DetectorMessage{{Role: "user", Text: "hi 早安，今天天氣如何"}}, nil)
	if len(got) != 0 {
		t.Fatalf("expected no match, got %v", detectedIDs(got))
	}
}

func TestDetectWorkflows_LegacyFlatNormalization(t *testing.T) {
	// legacy flat route hit via user_signals (no activation_any_of present)
	got := DetectWorkflows(testRegistry(),
		[]DetectorMessage{{Role: "user", Text: "用 Frida hook 這個 app"}}, nil)
	if ids := detectedIDs(got); !reflect.DeepEqual(ids, []string{"route.workflow.apk-analysis"}) {
		t.Fatalf("expected apk-analysis via legacy flat, got %v", ids)
	}
}

func TestDetectWorkflows_LegacyFlatContextGlob(t *testing.T) {
	// legacy file_change_globs normalized to context_signals (pre-Read path)
	got := DetectWorkflows(testRegistry(), nil,
		[]DetectorFile{{Path: "repo/scripts/frida/hook_login.js"}})
	if ids := detectedIDs(got); !reflect.DeepEqual(ids, []string{"route.workflow.apk-analysis"}) {
		t.Fatalf("expected apk-analysis via legacy glob, got %v", ids)
	}
	r, _ := findRoute(got, "route.workflow.apk-analysis")
	if !r.Activated || len(r.ContextSignalHits) != 1 {
		t.Fatalf("expected context-signal activation, got %+v", r)
	}
}

func TestDetectWorkflows_ContextGlobTwoPhase(t *testing.T) {
	got := DetectWorkflows(testRegistry(), nil,
		[]DetectorFile{{Path: "project/analysis/web/scraper.md"}})
	if ids := detectedIDs(got); !reflect.DeepEqual(ids, []string{"route.analysis.web"}) {
		t.Fatalf("expected analysis.web via context glob, got %v", ids)
	}
}

func TestDetectWorkflows_ReinforcementOnly_IsLateDetected(t *testing.T) {
	// No user/context signal — only artifact content matches. Route must appear
	// but Activated=false (late-detected), per the two-phase invariant.
	got := DetectWorkflows(testRegistry(), nil,
		[]DetectorFile{{Path: "notes.txt", Content: "the site uses an anti-bot gateway"}})
	r, ok := findRoute(got, "route.analysis.web")
	if !ok {
		t.Fatalf("expected analysis.web surfaced via reinforcement, got %v", detectedIDs(got))
	}
	if r.Activated {
		t.Fatalf("reinforcement-only hit must NOT activate (late-detected), got Activated=true")
	}
	if len(r.ArtifactReinforce) == 0 {
		t.Fatalf("expected artifact reinforcement recorded, got %+v", r)
	}
}

func TestDetectWorkflows_ReinforcementRegexAlternation(t *testing.T) {
	got := DetectWorkflows(testRegistry(), nil,
		[]DetectorFile{{Path: "page.md", Content: "use a CSS selector to grab the node"}})
	r, ok := findRoute(got, "route.analysis.web")
	if !ok || r.Activated {
		t.Fatalf("expected reinforce-only via 'selector|xpath' regex, got %+v / %v", r, detectedIDs(got))
	}
}

func TestDetectWorkflows_AdvisoryActivatesButIsAdvisoryMode(t *testing.T) {
	got := DetectWorkflows(testRegistry(),
		[]DetectorMessage{{Role: "user", Text: "看一下 magic bytes 怎麼判斷"}}, nil)
	r, ok := findRoute(got, "route.intelligence.engineering.heuristics")
	if !ok {
		t.Fatalf("expected heuristics advisory route, got %v", detectedIDs(got))
	}
	if r.EffectiveMode != "advisory" {
		t.Fatalf("expected EffectiveMode=advisory, got %q", r.EffectiveMode)
	}
}

func TestDetectWorkflows_OnDemandNeverParticipates(t *testing.T) {
	// agent-architecture is on-demand and must not be detected even on a direct
	// keyword match (user deliberately chose on-demand to avoid auto-firing).
	got := DetectWorkflows(testRegistry(),
		[]DetectorMessage{{Role: "user", Text: "discuss agent architecture for this system"}}, nil)
	if _, ok := findRoute(got, "route.intelligence.engineering.agent-architecture"); ok {
		t.Fatalf("on-demand route must never flow through detector, got %v", detectedIDs(got))
	}
}

func TestEffectiveActivationMode_DerivesFromType(t *testing.T) {
	if m := effectiveActivationMode(runtimeRouteRecord{RouteType: "workflow"}); m != "auto-detect" {
		t.Fatalf("workflow default = %q, want auto-detect", m)
	}
	// explicit mode overrides type default
	if m := effectiveActivationMode(runtimeRouteRecord{RouteType: "intelligence", ActivationMode: "on-demand"}); m != "on-demand" {
		t.Fatalf("explicit override = %q, want on-demand", m)
	}
}

func TestGlobToRegexp_DoubleStarMatchesZeroSegments(t *testing.T) {
	re := globToRegexp("**/*.apk")
	if re == nil {
		t.Fatal("nil regexp")
	}
	for _, p := range []string{"app.apk", "a/b/c/app.apk", "deep/nested/path/x.apk"} {
		if !re.MatchString(p) {
			t.Fatalf("expected %q to match **/*.apk", p)
		}
	}
	for _, p := range []string{"app.apkx", "app.txt"} {
		if re.MatchString(p) {
			t.Fatalf("did not expect %q to match **/*.apk", p)
		}
	}
}

func TestGlobToRegexp_SingleStarDoesNotCrossSegment(t *testing.T) {
	re := globToRegexp("scripts/*.js")
	if re.MatchString("scripts/sub/hook.js") {
		t.Fatal("single * must not cross a path separator")
	}
	if !re.MatchString("scripts/hook.js") {
		t.Fatal("expected scripts/hook.js to match")
	}
}

func TestNormalizeRouteTriggers_MergesBothForms(t *testing.T) {
	n := normalizeRouteTriggers(runtimeRouteTriggers{
		UserSignals:     []string{"legacy-user"},
		FileChangeGlobs: []string{"**/legacy/**"},
		ActivationAnyOf: &runtimeActivationAnyOf{
			UserSignals:    []string{"new-user", "legacy-user"}, // dup folded
			ContextSignals: []string{"**/new/**"},
		},
		ReinforcementAnyOf: &runtimeReinforcementAnyOf{ArtifactSignals: []string{"art"}},
	})
	if !reflect.DeepEqual(n.userSignals, []string{"legacy-user", "new-user"}) {
		t.Fatalf("user merge/dedupe wrong: %v", n.userSignals)
	}
	if !reflect.DeepEqual(n.contextSignals, []string{"**/legacy/**", "**/new/**"}) {
		t.Fatalf("context merge wrong: %v", n.contextSignals)
	}
	if !reflect.DeepEqual(n.artifactSignals, []string{"art"}) {
		t.Fatalf("artifact wrong: %v", n.artifactSignals)
	}
}
