package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
)

type initProjectOptions struct {
	projectPath string
	tools       string
	dryRun      bool
	force       bool
	jsonOutput  bool
	plainOutput bool
}

type plannedFile struct {
	tool        string
	path        string
	description string
}

func runInitProject(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := newFlagSet("init-project", stderr)
	opts := initProjectOptions{}
	fs.StringVar(&opts.projectPath, "project", "", "target project directory")
	fs.StringVar(&opts.tools, "tools", "roo,cursor,claude", "comma-separated tools: roo,cursor,claude")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview planned files without writing")
	fs.BoolVar(&opts.force, "force", false, "allow overwriting existing files")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")

	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildInitProjectResult(opts)
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

func buildInitProjectResult(opts initProjectOptions) Result {
	result := Result{
		Command:        "init-project",
		Mode:           "dry_run",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}

	if !opts.dryRun {
		result.Status = "blocked"
		result.ExitCode = ExitPartialCloseBlocked
		result.Error = &CommandError{
			Code:        "write_mode_not_implemented",
			Message:     "init-project currently supports dry-run planning only.",
			Remediation: "Run with --dry-run until write mode has fixture-backed parity with init-new-project.sh.",
		}
	}

	target, targetCheck := resolveTargetProject(opts.projectPath)
	result.Checks = append(result.Checks, targetCheck)
	if targetCheck.Status != "ok" && result.ExitCode == ExitSuccess {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{
			Code:        "invalid_project",
			Message:     targetCheck.Message,
			Remediation: "Pass --project with an existing project directory.",
		}
	}

	tools, toolsCheck := parseInitTools(opts.tools)
	result.Checks = append(result.Checks, toolsCheck)
	if toolsCheck.Status == "failed" && result.ExitCode == ExitSuccess {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{
			Code:        "invalid_tools",
			Message:     toolsCheck.Message,
			Remediation: "Use one or more supported tools: roo,cursor,claude.",
		}
	}

	if target == "" || len(tools) == 0 {
		return result
	}

	files := initProjectPlannedFiles(target, tools)
	conflicts := []string{}
	for _, file := range files {
		result.PlannedActions = append(result.PlannedActions, fmt.Sprintf("create %s: %s", file.description, file.path))
		if _, err := os.Stat(file.path); err == nil {
			conflicts = append(conflicts, file.path)
		}
	}
	if len(conflicts) > 0 && !opts.force && result.ExitCode == ExitSuccess {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{
			Code:        "target_exists",
			Message:     "target files already exist: " + strings.Join(conflicts, ", "),
			Remediation: "Pass --force only after reviewing the planned overwrite list.",
		}
		result.Checks = append(result.Checks, Check{Name: "conflicts", Status: "failed", Message: fmt.Sprintf("%d existing target files", len(conflicts))})
	} else {
		result.Checks = append(result.Checks, Check{Name: "conflicts", Status: "ok", Message: "no blocking file conflicts"})
	}

	return result
}

func resolveTargetProject(projectPath string) (string, Check) {
	if strings.TrimSpace(projectPath) == "" {
		return "", Check{Name: "project", Status: "failed", Message: "--project is required"}
	}

	abs, err := filepath.Abs(projectPath)
	if err != nil {
		return "", Check{Name: "project", Status: "failed", Message: err.Error()}
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", Check{Name: "project", Status: "failed", Message: err.Error()}
	}
	if !info.IsDir() {
		return "", Check{Name: "project", Status: "failed", Message: "project path is not a directory"}
	}

	normalized, err := pathutil.NormalizeForReport(abs)
	if err != nil {
		return "", Check{Name: "project", Status: "failed", Message: err.Error()}
	}
	return abs, Check{Name: "project", Status: "ok", Message: normalized}
}

func parseInitTools(value string) ([]string, Check) {
	supported := map[string]bool{"roo": true, "cursor": true, "claude": true}
	seen := map[string]bool{}
	ignored := []string{}
	for _, part := range strings.Split(value, ",") {
		tool := strings.ToLower(strings.TrimSpace(part))
		if tool == "" {
			continue
		}
		if !supported[tool] {
			ignored = append(ignored, tool)
			continue
		}
		seen[tool] = true
	}

	tools := make([]string, 0, len(seen))
	for tool := range seen {
		tools = append(tools, tool)
	}
	sort.Strings(tools)
	if len(tools) == 0 {
		return nil, Check{Name: "tools", Status: "failed", Message: "no supported tools selected"}
	}
	if len(ignored) > 0 {
		return tools, Check{Name: "tools", Status: "warning", Message: "ignored unsupported tools: " + strings.Join(ignored, ", ")}
	}
	return tools, Check{Name: "tools", Status: "ok", Message: strings.Join(tools, ",")}
}

func initProjectPlannedFiles(target string, tools []string) []plannedFile {
	files := []plannedFile{}
	for _, tool := range tools {
		switch tool {
		case "roo":
			files = append(files, plannedFile{tool: tool, path: filepath.Join(target, ".roomodes"), description: "Roo Code settings"})
		case "cursor":
			files = append(files,
				plannedFile{tool: tool, path: filepath.Join(target, ".cursor", "rules", "ai-skill-bootstrap.mdc"), description: "Cursor rule"},
				plannedFile{tool: tool, path: filepath.Join(target, ".cursor", "hooks.json"), description: "Cursor hooks"},
			)
		case "claude":
			files = append(files, plannedFile{tool: tool, path: filepath.Join(target, "CLAUDE.md"), description: "Claude Code settings"})
		}
	}
	files = append(files, plannedFile{tool: "common", path: filepath.Join(target, ".agent-goals", "README.md"), description: "agent goals ledger"})
	return files
}
