package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
)

type goalsOptions struct {
	command     string
	projectPath string
	dryRun      bool
	jsonOutput  bool
	plainOutput bool
}

func runGoals(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill goals <status|init> [flags]")
		return ExitInvalidUsage
	}

	opts := goalsOptions{command: args[0]}
	fs := newFlagSet("goals "+opts.command, stderr)
	fs.StringVar(&opts.projectPath, "project", ".", "project root")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview planned changes without writing")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildGoalsResult(opts)
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

func buildGoalsResult(opts goalsOptions) Result {
	result := Result{
		Command:        "goals " + opts.command,
		Mode:           "check",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}

	project, projectCheck := resolveTargetProject(opts.projectPath)
	result.Checks = append(result.Checks, projectCheck)
	if projectCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_project", Message: projectCheck.Message, Remediation: "Pass --project with an existing project directory."}
		return result
	}

	switch opts.command {
	case "status":
		return buildGoalsStatusResult(result, project)
	case "init":
		return buildGoalsInitResult(result, project, opts.dryRun)
	default:
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_goals_command", Message: "unsupported goals command: " + opts.command, Remediation: "Use status or init for the current Phase 2 slice."}
		return result
	}
}

func buildGoalsStatusResult(result Result, project string) Result {
	root := filepath.Join(project, ".agent-goals")
	goalsDir := filepath.Join(root, "goals")
	locksDir := filepath.Join(root, "locks")
	if _, err := os.Stat(root); err != nil {
		result.Checks = append(result.Checks, Check{Name: "goal_ledger", Status: "missing", Message: ".agent-goals does not exist"})
		return result
	}

	result.Checks = append(result.Checks,
		pathExistenceCheck("goal_ledger", root),
		pathExistenceCheck("goals_dir", goalsDir),
		pathExistenceCheck("locks_dir", locksDir),
	)
	goalFiles := listMarkdownFiles(goalsDir)
	lockFiles := listDirEntries(locksDir)
	result.Checks = append(result.Checks,
		Check{Name: "goal_files", Status: "ok", Message: fmt.Sprintf("%d files", len(goalFiles))},
		Check{Name: "locks", Status: "ok", Message: fmt.Sprintf("%d entries", len(lockFiles))},
	)
	return result
}

func buildGoalsInitResult(result Result, project string, dryRun bool) Result {
	result.Mode = "dry_run"
	root := filepath.Join(project, ".agent-goals")
	result.PlannedActions = []string{
		"create directory: " + filepath.Join(root, "goals"),
		"create directory: " + filepath.Join(root, "locks"),
		"create file: " + filepath.Join(root, "README.md"),
		"ensure git exclude contains .agent-goals/",
	}
	if !dryRun {
		result.Status = "blocked"
		result.ExitCode = ExitPartialCloseBlocked
		result.Error = &CommandError{
			Code:        "write_mode_not_implemented",
			Message:     "goals init currently supports dry-run planning only.",
			Remediation: "Run with --dry-run until write mode has fixture-backed parity with agent-goals.sh.",
		}
	}
	return result
}

func pathExistenceCheck(name string, path string) Check {
	normalized, err := pathutil.NormalizeForReport(path)
	if err != nil {
		normalized = path
	}
	if _, err := os.Stat(path); err != nil {
		return Check{Name: name, Status: "missing", Message: normalized}
	}
	return Check{Name: name, Status: "ok", Message: normalized}
}

func listMarkdownFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	files := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	return files
}

func listDirEntries(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}
