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
	case "doctor":
		return runDoctor(args[1:], stdout, stderr)
	case "init-project":
		return runInitProject(args[1:], stdout, stderr)
	case "goals":
		return runGoals(args[1:], stdout, stderr)
	case "close-loop":
		return runCloseLoop(args[1:], stdout, stderr)
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
	_, _ = fmt.Fprintln(w, "  doctor    check local runtime and repository readiness")
	_, _ = fmt.Fprintln(w, "  init-project    plan project-local AI tool bootstrap files")
	_, _ = fmt.Fprintln(w, "  goals    inspect or plan project-local goal ledger changes")
	_, _ = fmt.Fprintln(w, "  close-loop    inspect repository close-loop readiness")
}

func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	return fs
}
