package app

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/glossary"
)

type glossaryOptions struct {
	command     string
	repoPath    string
	glossaryDir string
	jsonOutput  bool
	plainOutput bool
}

func runGlossary(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill glossary <validate> [flags]")
		return ExitInvalidUsage
	}
	opts := glossaryOptions{command: args[0]}
	if opts.command != "validate" {
		_, _ = fmt.Fprintf(stderr, "unsupported glossary command: %s\n", opts.command)
		return ExitInvalidUsage
	}
	fs := newFlagSet("glossary "+opts.command, stderr)
	fs.StringVar(&opts.repoPath, "repo", ".", "Ai-skill repository path")
	fs.StringVar(&opts.glossaryDir, "glossary", "", "glossary directory (default <repo>/knowledge/glossary)")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildGlossaryValidateResult(opts)
	if opts.jsonOutput {
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

func buildGlossaryValidateResult(opts glossaryOptions) Result {
	dir := opts.glossaryDir
	if dir == "" {
		dir = filepath.Join(opts.repoPath, "knowledge", "glossary")
	}
	result := Result{
		Command:  "glossary validate",
		Mode:     "check",
		Status:   "success",
		ExitCode: ExitSuccess,
	}
	res, err := glossary.Validate(glossary.ValidateOptions{GlossaryDir: dir})
	if err != nil {
		result.Status = "failed"
		result.ExitCode = ExitGeneralFailure
		result.Error = &CommandError{Code: "general_failure", Message: fmt.Sprintf("validate: %v", err)}
		return result
	}
	summary := fmt.Sprintf("entries=%d aliases=%d relations=%d violations=%d",
		res.EntryCount, res.AliasCount, res.RelationCount, len(res.Violations))
	result.Checks = append(result.Checks, Check{
		Name:    "glossary_discovery",
		Status:  "ok",
		Message: fmt.Sprintf("glossary dir: %s", dir),
	})
	if len(res.Violations) == 0 {
		result.Checks = append(result.Checks, Check{
			Name:    "glossary_schema",
			Status:  "ok",
			Message: summary,
		})
		return result
	}
	result.Status = "failed"
	result.ExitCode = ExitValidationFailed
	result.Error = &CommandError{Code: "validation_failed", Message: summary}
	for _, v := range res.Violations {
		msg := v.Message
		if v.File != "" && !strings.Contains(msg, v.File) {
			msg = fmt.Sprintf("%s: %s", v.File, msg)
		}
		result.Checks = append(result.Checks, Check{
			Name:        v.RuleID,
			Status:      "error",
			Message:     msg,
			Remediation: v.Remediation,
		})
	}
	return result
}
