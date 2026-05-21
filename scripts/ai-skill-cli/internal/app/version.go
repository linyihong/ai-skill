package app

import (
	"fmt"
	"io"
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func runVersion(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := newFlagSet("version", stderr)
	jsonOutput := fs.Bool("json", false, "write machine-readable JSON output")
	plainOutput := fs.Bool("plain", false, "write human-readable output")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if *jsonOutput && *plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := Result{
		Command:  "version",
		Mode:     "info",
		Status:   "success",
		ExitCode: ExitSuccess,
		Checks: []Check{
			{Name: "version", Status: "ok", Message: Version},
			{Name: "commit", Status: "ok", Message: Commit},
			{Name: "date", Status: "ok", Message: Date},
		},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	if *jsonOutput {
		if err := writeJSON(stdout, result); err != nil {
			_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
			return ExitGeneralFailure
		}
		return result.ExitCode
	}
	if err := writePlain(stdout, result); err != nil {
		_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
		return ExitGeneralFailure
	}
	return result.ExitCode
}
