package app

import (
	"os"
	"path/filepath"
	"testing"
)

// writeGateRegistry writes a minimal routing-registry.yaml into a fake Ai-skill
// repo so workflowPrimarySourceGate can resolve routes + primary_source.
func writeGateRegistry(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "knowledge", "runtime"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := `registry_version: knowledge-routing/v2
records:
  - id: route.analysis.web
    route_type: analysis
    activation_mode: auto-detect
    task_intent: web scraping
    activation_triggers:
      activation_any_of:
        user_signals:
          - web scraping
    primary_source: analysis/web/README.md
  - id: route.intelligence.architectural-fit
    route_type: intelligence
    activation_mode: auto-detect
    task_intent: architecture
    activation_triggers:
      activation_any_of:
        user_signals:
          - DDD
    primary_source: intelligence/engineering/architecture/architectural-fit/README.md
`
	if err := os.WriteFile(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"), []byte(registry), 0o644); err != nil {
		t.Fatal(err)
	}
	return repo
}

// transcript line helpers
func userLine(text string) string {
	return `{"type":"user","message":{"role":"user","content":` + jsonQuote(text) + `}}`
}

// assistantReadLine emits an assistant turn containing a Read tool_use of path.
func assistantReadLine(path string) string {
	return `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"Read","input":{"file_path":` + jsonQuote(path) + `}}]}}`
}

func jsonQuote(s string) string {
	b := []byte{'"'}
	for _, r := range s {
		switch r {
		case '"', '\\':
			b = append(b, '\\', byte(r))
		default:
			b = append(b, []byte(string(r))...)
		}
	}
	b = append(b, '"')
	return string(b)
}

func writeTranscript(t *testing.T, lines ...string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "transcript.jsonl")
	data := ""
	for _, l := range lines {
		data += l + "\n"
	}
	if err := os.WriteFile(f, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestWorkflowGate_LockedAndUnread_Blocks(t *testing.T) {
	repo := writeGateRegistry(t)
	tr := writeTranscript(t, userLine("幫我做 web scraping 抓網站"))
	block, route, ps := workflowPrimarySourceGate(tr, repo)
	if !block {
		t.Fatalf("expected block (locked, primary_source unread); route=%q ps=%q", route, ps)
	}
	if route != "route.analysis.web" || ps != "analysis/web/README.md" {
		t.Fatalf("unexpected route/ps: %q %q", route, ps)
	}
}

func TestWorkflowGate_LockedAndRead_Allows(t *testing.T) {
	repo := writeGateRegistry(t)
	tr := writeTranscript(t,
		userLine("幫我做 web scraping 抓網站"),
		assistantReadLine("/abs/path/analysis/web/README.md"), // suffix match
	)
	block, _, _ := workflowPrimarySourceGate(tr, repo)
	if block {
		t.Fatal("expected allow once primary_source Read")
	}
}

func TestWorkflowGate_Conflict_Allows(t *testing.T) {
	repo := writeGateRegistry(t)
	// both routes activate → conflict → ActiveRoute empty → never block
	tr := writeTranscript(t, userLine("做 web scraping 並評估 DDD 架構"))
	block, _, _ := workflowPrimarySourceGate(tr, repo)
	if block {
		t.Fatal("conflict (>1 route) must never block")
	}
}

func TestWorkflowGate_Miss_Allows(t *testing.T) {
	repo := writeGateRegistry(t)
	tr := writeTranscript(t, userLine("hi 早安"))
	block, _, _ := workflowPrimarySourceGate(tr, repo)
	if block {
		t.Fatal("detector miss must never block")
	}
}

func TestWorkflowGate_FailOpen_NoRepo(t *testing.T) {
	tr := writeTranscript(t, userLine("幫我做 web scraping"))
	if block, _, _ := workflowPrimarySourceGate(tr, ""); block {
		t.Fatal("must fail open when repo unresolvable")
	}
	if block, _, _ := workflowPrimarySourceGate(tr, t.TempDir()); block {
		t.Fatal("must fail open when registry missing")
	}
}

func TestWorkflowGate_NoTranscript_Allows(t *testing.T) {
	repo := writeGateRegistry(t)
	if block, _, _ := workflowPrimarySourceGate("", repo); block {
		t.Fatal("no transcript path must fail open")
	}
}
