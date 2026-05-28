package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestClassifyRouteManualWins(t *testing.T) {
	manual := map[string]string{"route.example": "workflow_discovery"}
	entry := classifyRoute("route.example", "route.example is in signal text", manual, "validate(\"route.example\")")
	if entry.Classification != ClassManual {
		t.Fatalf("manual_activation should win, got %q", entry.Classification)
	}
	if !strings.Contains(entry.Evidence, "workflow_discovery") {
		t.Fatalf("evidence must cite manual reason, got %q", entry.Evidence)
	}
}

func TestClassifyRouteAutoDetectedBeatsConsumed(t *testing.T) {
	entry := classifyRoute("route.foo", "signal references route.foo", nil, "validate(\"route.foo\")")
	if entry.Classification != ClassAutoDetected {
		t.Fatalf("auto-detected via signal should outrank consumed, got %q", entry.Classification)
	}
}

func TestClassifyRouteConsumed(t *testing.T) {
	entry := classifyRoute("route.bar", "unrelated signal text", nil, "validateBar(\"route.bar\")")
	if entry.Classification != ClassConsumed {
		t.Fatalf("consumed should fire when Go source references id, got %q", entry.Classification)
	}
}

func TestClassifyRouteOrphan(t *testing.T) {
	entry := classifyRoute("route.lonely", "", nil, "")
	if entry.Classification != ClassOrphan {
		t.Fatalf("expected orphan, got %q", entry.Classification)
	}
}

func TestClassifySurfaceConsumed(t *testing.T) {
	s := rawSurface{TargetKey: "runtime.example.key", SourcePath: "runtime/x.yaml"}
	entry := classifySurface(s, `Query("runtime.example.key")`)
	if entry.Classification != ClassConsumed {
		t.Fatalf("expected consumed, got %q", entry.Classification)
	}
}

func TestClassifySurfaceOrphan(t *testing.T) {
	s := rawSurface{TargetKey: "runtime.lonely.key", SourcePath: "runtime/x.yaml"}
	entry := classifySurface(s, "")
	if entry.Classification != ClassOrphan {
		t.Fatalf("expected orphan, got %q", entry.Classification)
	}
}

func TestClassifyScenarioConsumed(t *testing.T) {
	sc := rawScenario{ID: "orphan-routing-entry-v1", Path: "/tmp/x.yaml"}
	entry := classifyScenario(sc, `const id = "orphan-routing-entry-v1"`, "/tmp")
	if entry.Classification != ClassConsumed {
		t.Fatalf("expected consumed, got %q", entry.Classification)
	}
}

func TestClassifyScenarioOrphan(t *testing.T) {
	sc := rawScenario{ID: "lonely-scenario-v1", Path: "/tmp/x.yaml"}
	entry := classifyScenario(sc, "no references here", "/tmp")
	if entry.Classification != ClassOrphan {
		t.Fatalf("expected orphan, got %q", entry.Classification)
	}
}

func TestRenderMarkdownTablesAndSummary(t *testing.T) {
	inv := &Inventory{
		Repo: "/repo",
		Routes: []RouteEntry{
			{ID: "route.a", Classification: ClassAutoDetected, Evidence: "signal cites route.a"},
			{ID: "route.b", Classification: ClassOrphan, Evidence: "no consumer"},
		},
		Surfaces: []SurfaceEntry{
			{TargetKey: "runtime.x", SourcePath: "runtime/x.yaml", Classification: ClassConsumed, Evidence: "Go cites it"},
		},
		Scenarios: []ScenarioEntry{
			{ID: "scenario-x-v1", Path: "validation/scenarios/x.yaml", Classification: ClassOrphan, Evidence: "no Go ref"},
		},
		Summary: InventorySummary{
			RouteCounts:    map[string]int{ClassAutoDetected: 1, ClassOrphan: 1},
			SurfaceCounts:  map[string]int{ClassConsumed: 1},
			ScenarioCounts: map[string]int{ClassOrphan: 1},
			OrphanTotal:    2,
		},
	}
	var buf bytes.Buffer
	if err := RenderMarkdown(&buf, inv); err != nil {
		t.Fatalf("render: %v", err)
	}
	out := buf.String()
	for _, must := range []string{
		"# Ai-skill Runtime Audit Report",
		"Orphan total**: 2",
		"route.a",
		"route.b",
		"runtime.x",
		"scenario-x-v1",
		"## Summary",
		"## Routes",
		"## Generated surfaces",
		"## Validation scenarios",
	} {
		if !strings.Contains(out, must) {
			t.Errorf("markdown missing %q\nfull output:\n%s", must, out)
		}
	}
}

func TestRenderJSONRoundTrip(t *testing.T) {
	inv := &Inventory{
		Repo:    "/repo",
		Routes:  []RouteEntry{{ID: "route.a", Classification: ClassOrphan, Evidence: "x"}},
		Summary: InventorySummary{RouteCounts: map[string]int{ClassOrphan: 1}, OrphanTotal: 1},
	}
	var buf bytes.Buffer
	if err := RenderJSON(&buf, inv); err != nil {
		t.Fatalf("render: %v", err)
	}
	var decoded Inventory
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.Repo != "/repo" || len(decoded.Routes) != 1 || decoded.Summary.OrphanTotal != 1 {
		t.Fatalf("round trip mismatch: %#v", decoded)
	}
}

func TestEscapePipeInEvidence(t *testing.T) {
	inv := &Inventory{
		Repo: "/repo",
		Routes: []RouteEntry{
			{ID: "route.a", Classification: ClassOrphan, Evidence: "evidence | with | pipes"},
		},
		Summary: InventorySummary{RouteCounts: map[string]int{ClassOrphan: 1}},
	}
	var buf bytes.Buffer
	if err := RenderMarkdown(&buf, inv); err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(buf.String(), `evidence \| with \| pipes`) {
		t.Fatalf("pipes in evidence not escaped:\n%s", buf.String())
	}
}
