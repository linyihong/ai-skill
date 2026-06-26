package app

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/planvalidate"
)

// Phase 2.4 acceptance. The CLI integration test calls the binary end-to-end and
// compares against the ENGINE ENTRYPOINT (planvalidate.Validate) only — it does
// not reach into validator internals. This guards the "CLI = transport surface"
// contract: the CLI must project engine findings, never compute its own.

type cliFinding struct {
	RuleID   string `json:"rule_id"`
	Message  string `json:"message"`
	Blocking bool   `json:"blocking"`
}
type cliPayload struct {
	Plans    int          `json:"plans"`
	Blocking int          `json:"blocking"`
	Findings []cliFinding `json:"findings"`
}

func runValidateJSON(t *testing.T, root string) (cliPayload, int) {
	t.Helper()
	var out, errb bytes.Buffer
	code := Run([]string{"plans", "validate", "--root", root, "--format", "json"}, &out, &errb)
	var p cliPayload
	if err := json.Unmarshal(out.Bytes(), &p); err != nil {
		t.Fatalf("CLI json parse failed: %v\noutput=%s", err, out.String())
	}
	return p, code
}

// CLI output == engine entrypoint projection (same root, same findings).
func TestPlansValidateCLI_MatchesEngineEntrypoint(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: null\n---")
	makePlan(t, tmp, "plans/active/01-x.md",
		"---\nid: s\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: ghost\nrequired_for_completion: true\nsub_plan_reason: x\n---")

	cli, _ := runValidateJSON(t, tmp)

	// Engine entrypoint, computed independently in the test.
	models, compat := normalizedPlansFromRoot(tmp)
	engine := planvalidate.Validate(
		planvalidate.ValidationContext{Root: tmp, ExecutionMode: planvalidate.ModeManual},
		models)
	engine = append(engine, compat...)

	if len(cli.Findings) != len(engine) {
		t.Fatalf("CLI findings=%d != engine findings=%d", len(cli.Findings), len(engine))
	}
	got := map[string]bool{}
	for _, f := range cli.Findings {
		got[f.RuleID] = true
	}
	for _, f := range engine {
		if !got[f.RuleID] {
			t.Fatalf("engine finding %q missing from CLI projection", f.RuleID)
		}
	}
}

// Phase 3.2 end-to-end (loader -> compat layer -> consumer): a supported
// schema_version flows and validates clean; an unsupported one is a deterministic,
// diagnosable blocking reject surfaced by the CLI (not silently degraded).
func TestPlansValidateCLI_SchemaVersionEndToEnd(t *testing.T) {
	// supported version "2" (quoted, as real frontmatter) -> clean, exit 0
	sup := t.TempDir()
	makePlan(t, sup, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-25\nparent: null\nschema_version: \"2\"\n---")
	if cli, code := runValidateJSON(t, sup); code != ExitSuccess || len(cli.Findings) != 0 {
		t.Fatalf("supported schema_version 2 should be clean exit 0, got code=%d findings=%d", code, len(cli.Findings))
	}
	// unsupported version "99" -> blocking compat reject, exit 30
	uns := t.TempDir()
	makePlan(t, uns, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-25\nparent: null\nschema_version: \"99\"\n---")
	cli, code := runValidateJSON(t, uns)
	if code != ExitValidationFailed {
		t.Fatalf("unsupported schema_version should exit ExitValidationFailed, got %d", code)
	}
	found := false
	for _, f := range cli.Findings {
		if f.RuleID == "compat.unsupported_schema_version" && f.Blocking {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected blocking compat.unsupported_schema_version finding, got %+v", cli.Findings)
	}
}

// Valid canonical tree -> CLI clean, exit 0.
func TestPlansValidateCLI_ValidTreeExitZero(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: null\n---")
	makePlan(t, tmp, "plans/active/01-x.md",
		"---\nid: s\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: m\nrequired_for_completion: true\nsub_plan_reason: x\n---")
	cli, code := runValidateJSON(t, tmp)
	if code != ExitSuccess || len(cli.Findings) != 0 {
		t.Fatalf("valid tree should be clean exit 0, got code=%d findings=%d", code, len(cli.Findings))
	}
}

// Violation -> CLI non-zero exit, AND three-way equivalence: the legacy hook
// validator fires the same rule the CLI (engine) reports on the same tree.
func TestPlansValidateCLI_ViolationThreeWayEquivalence(t *testing.T) {
	tmp := t.TempDir()
	makePlan(t, tmp, "plans/active/_plan.md",
		"---\nid: m\nplan_kind: main\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: null\n---")
	makePlan(t, tmp, "plans/active/01-x.md",
		"---\nid: s\nplan_kind: sub\nstatus: draft\nowner: t\ncreated: 2026-06-24\nparent: ghost\nrequired_for_completion: true\nsub_plan_reason: x\n---")

	cli, code := runValidateJSON(t, tmp)
	if code != ExitValidationFailed {
		t.Fatalf("violation should exit ExitValidationFailed, got %d", code)
	}
	cliHasParentRef := false
	for _, f := range cli.Findings {
		if f.RuleID == "plan_tree.parent_reference" {
			cliHasParentRef = true
		}
	}
	if !cliHasParentRef {
		t.Fatalf("CLI should report parent_reference for unresolved parent")
	}
	// Legacy hook validator (the other consumer) must agree on the same tree.
	legacy := validatePlanTreeParentReference("commit\n", []string{"plans/active/01-x.md"}, tmp)
	if legacy == "" {
		t.Fatalf("legacy validator should also fire parent_reference (three-way divergence)")
	}
}
