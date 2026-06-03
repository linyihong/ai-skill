package app

import (
	"flag"
	"fmt"
	"io"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill <command> [flags]")
		return ExitInvalidUsage
	}

	switch args[0] {
	case "version":
		return runVersion(args[1:], stdout, stderr)
	case "doctor":
		return runDoctor(args[1:], stdout, stderr)
	case "init-project":
		return runInitProject(args[1:], stdout, stderr)
	case "goals":
		return runGoals(args[1:], stdout, stderr)
	case "close-loop":
		return runCloseLoop(args[1:], stdout, stderr)
	case "hooks":
		return runHooks(args[1:], stdout, stderr)
	case "sync-cursor-bundle":
		return runSyncCursorBundle(args[1:], stdout, stderr)
	case "runtime":
		return runRuntime(args[1:], stdout, stderr)
	case "roo":
		return runRoo(args[1:], stdout, stderr)
	case "copilot":
		return runCopilot(args[1:], stdout, stderr)
	case "glossary":
		return runGlossary(args[1:], stdout, stderr)
	case "scan-checkboxes":
		return runScanCheckboxes(args[1:], stdout, stderr)
	case "enforcement":
		return runEnforcement(args[1:], stdout, stderr)
	case "plans":
		return runPlans(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		printUsage(stdout)
		return ExitSuccess
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		return ExitInvalidUsage
	}
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "usage: ai-skill <command> [flags]")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "commands:")
	_, _ = fmt.Fprintln(w, "  version    print build version metadata")
	_, _ = fmt.Fprintln(w, "  doctor    check local runtime and repository readiness")
	_, _ = fmt.Fprintln(w, "  init-project    plan project-local AI tool bootstrap files")
	_, _ = fmt.Fprintln(w, "  goals    inspect or plan project-local goal ledger changes")
	_, _ = fmt.Fprintln(w, "  close-loop    inspect repository close-loop readiness")
	_, _ = fmt.Fprintln(w, "  hooks    inspect or plan Git hook installation")
	_, _ = fmt.Fprintln(w, "  sync-cursor-bundle    plan Cursor bundle mirror sync")
	_, _ = fmt.Fprintln(w, "  runtime    wrap or inspect runtime tooling")
	_, _ = fmt.Fprintln(w, "  roo    manage guarded Roo Code settings")
	_, _ = fmt.Fprintln(w, "  copilot    generate GitHub Copilot guided bootstrap prompts")
	_, _ = fmt.Fprintln(w, "  glossary    validate knowledge/glossary/ entries")
	_, _ = fmt.Fprintln(w, "  scan-checkboxes    scan a Markdown file for unchecked/checked task-list items")
	_, _ = fmt.Fprintln(w, "  enforcement    run enforcement-registry lint or coverage report (Phase 4)")
	_, _ = fmt.Fprintln(w, "  plans    render plan-tree hierarchy (Phase 3 of plan-tree-hierarchy-governance)")
}

func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	return fs
}
