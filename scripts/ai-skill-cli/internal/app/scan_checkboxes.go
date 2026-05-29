package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// CheckboxScanResult holds the result of scanning a file or string for
// Markdown task-list checkboxes.
type CheckboxScanResult struct {
	File          string   `json:"file,omitempty"`
	UncheckedLines []string `json:"unchecked_lines"`
	CheckedLines   []string `json:"checked_lines"`
}

// HasUnchecked reports whether any unchecked "- [ ]" lines were found.
func (r CheckboxScanResult) HasUnchecked() bool { return len(r.UncheckedLines) > 0 }

// ScanCheckboxes scans arbitrary string content for Markdown checkbox lines.
// It is a pure function with no shell dependency — equivalent to
// grep -nE '^\s*- \[ \]' for unchecked and grep -nE '^\s*- \[x\]' for checked.
func ScanCheckboxes(content string) CheckboxScanResult {
	var unchecked, checked []string
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [ ]") {
			unchecked = append(unchecked, trimmed)
		} else if strings.HasPrefix(trimmed, "- [x]") || strings.HasPrefix(trimmed, "- [X]") {
			checked = append(checked, trimmed)
		}
	}
	return CheckboxScanResult{UncheckedLines: unchecked, CheckedLines: checked}
}

// ScanCheckboxesInFile reads the file at path and calls ScanCheckboxes.
func ScanCheckboxesInFile(path string) (CheckboxScanResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return CheckboxScanResult{}, err
	}
	result := ScanCheckboxes(string(data))
	result.File = path
	return result, nil
}

// runScanCheckboxes implements the "scan-checkboxes" CLI subcommand.
// Usage: ai-skill scan-checkboxes <file> [--format plain|json] [--exit-code]
func runScanCheckboxes(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := newFlagSet("scan-checkboxes", stderr)
	format := fs.String("format", "plain", "output format: plain or json")
	exitCode := fs.Bool("exit-code", false, "exit 1 when unchecked items are found")

	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if fs.NArg() != 1 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill scan-checkboxes <file> [--format plain|json] [--exit-code]")
		return ExitInvalidUsage
	}
	if *format != "plain" && *format != "json" {
		_, _ = fmt.Fprintf(stderr, "--format must be plain or json, got %q\n", *format)
		return ExitInvalidUsage
	}

	filePath := fs.Arg(0)
	result, err := ScanCheckboxesInFile(filePath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "scan-checkboxes: %v\n", err)
		return ExitGeneralFailure
	}

	if *format == "json" {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			_, _ = fmt.Fprintf(stderr, "scan-checkboxes: encode json: %v\n", err)
			return ExitGeneralFailure
		}
	} else {
		_, _ = fmt.Fprintf(stdout, "file: %s\n", result.File)
		_, _ = fmt.Fprintf(stdout, "unchecked: %d\n", len(result.UncheckedLines))
		_, _ = fmt.Fprintf(stdout, "checked:   %d\n", len(result.CheckedLines))
		if len(result.UncheckedLines) > 0 {
			_, _ = fmt.Fprintln(stdout, "\nunchecked items:")
			for _, line := range result.UncheckedLines {
				_, _ = fmt.Fprintf(stdout, "  %s\n", line)
			}
		}
	}

	if *exitCode && result.HasUnchecked() {
		return ExitGeneralFailure
	}
	return ExitSuccess
}
