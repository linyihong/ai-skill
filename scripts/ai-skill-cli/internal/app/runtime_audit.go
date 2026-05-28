package app

import (
	"fmt"
	"io"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/audit"
)

// nativeRuntimeAuditWarning runs the audit on the repo and reports the orphan
// totals as a non-blocking Check. The check always reports status "ok" so it
// never fails `ai-skill runtime validate`; orphan counts surface in the
// Message field. Use `ai-skill runtime audit` for the detailed report.
func nativeRuntimeAuditWarning(repo string) Check {
	inv, err := audit.Build(audit.Options{Repo: repo})
	if err != nil {
		return Check{
			Name:    "runtime_audit_warning",
			Status:  "ok",
			Message: fmt.Sprintf("audit skipped: %v (run `ai-skill runtime audit` for details)", err),
		}
	}
	msg := fmt.Sprintf(
		"orphan_total=%d (routes=%d, surfaces=%d, scenarios=%d); see `ai-skill runtime audit` for details",
		inv.Summary.OrphanTotal,
		inv.Summary.RouteCounts[audit.ClassOrphan],
		inv.Summary.SurfaceCounts[audit.ClassOrphan],
		inv.Summary.ScenarioCounts[audit.ClassOrphan],
	)
	return Check{Name: "runtime_audit_warning", Status: "ok", Message: msg}
}

// runRuntimeAudit handles the `ai-skill runtime audit` subcommand.
//
// Default output is human-readable markdown. Pass --json to get the full
// inventory in JSON for CI / tool consumption. The audit walks the routing
// registry, runtime.db generated_surfaces, and validation/scenarios/ then
// classifies each entry into auto-detected / consumed / intentionally-manual
// / orphan.
//
// Plan: plans/active/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md
// Phase: 2 (Inventory Tool, Graduation #1)
func runRuntimeAudit(opts runtimeOptions, stdout io.Writer, stderr io.Writer) int {
	root, repoCheck := closeLoopRepoRoot(opts.repoPath)
	if repoCheck.Status != "ok" {
		_, _ = fmt.Fprintln(stderr, repoCheck.Message)
		return ExitInvalidUsage
	}

	inv, err := audit.Build(audit.Options{Repo: root})
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "audit failed: %v\n", err)
		return ExitGeneralFailure
	}

	if opts.jsonOutput {
		if err := audit.RenderJSON(stdout, inv); err != nil {
			_, _ = fmt.Fprintf(stderr, "write json: %v\n", err)
			return ExitGeneralFailure
		}
		return ExitSuccess
	}

	if err := audit.RenderMarkdown(stdout, inv); err != nil {
		_, _ = fmt.Fprintf(stderr, "write markdown: %v\n", err)
		return ExitGeneralFailure
	}
	return ExitSuccess
}
