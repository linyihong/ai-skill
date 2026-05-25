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
	fs.StringVar(&opts.tools, "tools", "roo,cursor,claude,codex", "comma-separated tools: roo,cursor,claude,codex")
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
		result.Mode = "write"
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
			Remediation: "Use one or more supported tools: roo,cursor,claude,codex.",
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
	if result.ExitCode != ExitSuccess || opts.dryRun {
		return result
	}

	repo, repoCheck := resolveInitProjectAiSkillRepo()
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{
			Code:        "ai_skill_repo_unresolved",
			Message:     repoCheck.Message,
			Remediation: "Run ai-skill init-project from the Ai-skill repository or set the working directory to a checked-out repo.",
		}
		return result
	}
	for _, file := range files {
		content, err := initProjectFileContent(file, repo)
		if err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "template_error", Message: err.Error()}
			return result
		}
		if err := writeInitProjectFile(file.path, []byte(content), opts.force); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "write_failed", Message: err.Error()}
			return result
		}
		result.Mutations = append(result.Mutations, fmt.Sprintf("wrote %s: %s", file.description, file.path))
	}

	return result
}

func resolveInitProjectAiSkillRepo() (string, Check) {
	wd, err := os.Getwd()
	if err != nil {
		return "", Check{Name: "ai_skill_repo", Status: "failed", Message: err.Error()}
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "CORE_BOOTSTRAP.md")); err == nil {
			if _, err := os.Stat(filepath.Join(wd, "scripts", "ai-skill-cli")); err == nil {
				normalized, normErr := pathutil.NormalizeForReport(wd)
				if normErr != nil {
					normalized = wd
				}
				return wd, Check{Name: "ai_skill_repo", Status: "ok", Message: normalized}
			}
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", Check{Name: "ai_skill_repo", Status: "failed", Message: "could not locate Ai-skill repository root"}
		}
		wd = parent
	}
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
	supported := map[string]bool{"roo": true, "cursor": true, "claude": true, "codex": true}
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
		case "codex":
			files = append(files, plannedFile{tool: tool, path: filepath.Join(target, "AGENTS.md"), description: "Codex settings"})
		}
	}
	files = append(files, plannedFile{tool: "common", path: filepath.Join(target, ".agent-goals", "README.md"), description: "agent goals ledger"})
	return files
}

func writeInitProjectFile(path string, content []byte, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("target exists: %s", path)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func initProjectFileContent(file plannedFile, repo string) (string, error) {
	switch file.tool {
	case "roo":
		return initProjectRooContent(repo), nil
	case "cursor":
		if strings.HasSuffix(file.path, "hooks.json") {
			return initProjectCursorHooksContent(), nil
		}
		return initProjectCursorRuleContent(repo), nil
	case "claude":
		return initProjectClaudeContent(repo), nil
	case "codex":
		return initProjectCodexContent(repo), nil
	case "common":
		return initProjectGoalsReadmeContent(repo), nil
	default:
		return "", fmt.Errorf("unsupported init-project template: %s", file.tool)
	}
}

func initProjectBootstrapText(repo string) string {
	return fmt.Sprintf(`本專案使用 Ai-skill 知識庫：%s

> **IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
> 本檔為 **thin tool-entry pointer**。所有 obligation / format / enum / example 的 canonical 來源在 %s。Session 啟動 first turn 必須讀 CORE_BOOTSTRAP.md 並遵守其中**所有** obligations（含 Bootstrap Receipt、Cognitive Mode 報告 per-turn block、語言強制、knowledge-update-flow 等）。
> Summary 的「Resume directly」是對話 framing，**不豁免** runtime / governance bootstrap。Resume / continuation session 同樣須走完 bootstrap。

## 啟動序列

1. 讀 %s — 必讀規則 + Bootstrap Receipt + Cognitive Mode 報告 + 全部 per-session / per-turn obligations
2. 讀 %s — OS layout
3. 查 %s — 目前 phase / obligations / gates
4. 依任務 intent 查詢 routing-registry：%s

## 修改本檔的規則

本檔是 entry pointer，不是 canonical content。修改前先讀 %s。新 obligation 加到 CORE_BOOTSTRAP.md（cross-tool）或 ai-tools/agent/<tool>.md（tool-specific）— 不加到本檔。Commit-msg hook 會擋下違反 thinness 的修改。
`, repo,
		filepath.Join(repo, "CORE_BOOTSTRAP.md"),
		filepath.Join(repo, "CORE_BOOTSTRAP.md"),
		filepath.Join(repo, "README.md"),
		filepath.Join(repo, "runtime", "runtime.db"),
		filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"),
		filepath.Join(repo, "runtime", "bootstrap-entry-points.yaml"))
}

func initProjectRooContent(repo string) string {
	text := strings.ReplaceAll(initProjectBootstrapText(repo), "\n", "\\n")
	return fmt.Sprintf(`{
  "customModes": [
    {
      "slug": "code",
      "name": "Code",
      "roleDefinition": "You are a highly skilled software engineer.",
      "customInstructions": "%s",
      "groups": ["read", "edit", "command", "mcp"]
    },
    {
      "slug": "architect",
      "name": "Architect",
      "roleDefinition": "You are an expert software architect.",
      "customInstructions": "%s",
      "groups": ["read"]
    }
  ]
}
`, text, text)
}

func initProjectCursorRuleContent(repo string) string {
	return fmt.Sprintf(`---
description: Ai-skill 知識庫啟動流程
globs:
alwaysApply: true
---

# Ai-skill Bootstrap

%s`, initProjectBootstrapText(repo))
}

func initProjectCursorHooksContent() string {
	return `{
  "sessionStart": [
    {
      "description": "提醒載入 Ai-skill 知識庫",
      "command": "echo '提示：本專案使用 Ai-skill 知識庫，請確認已載入 CORE_BOOTSTRAP.md'"
    }
  ]
}
`
}

func initProjectClaudeContent(repo string) string {
	return fmt.Sprintf("# Claude Code Auto-Bootstrap\n\n%s", initProjectBootstrapText(repo))
}

func initProjectCodexContent(repo string) string {
	return fmt.Sprintf(`# Codex Adapter Bootstrap

本文件是 project-level Codex 自動載入入口，只負責指向 Ai-skill canonical source，不保存獨立規則副本。

啟動時，Codex 必須依序讀取：

1. %s
2. %s
3. %s
4. %s

若需要更新 Codex adapter 規則，請修改 Ai-skill repo 的 `+"`ai-tools/agent/codex.md`"+`，不要在本專案複製一套平行規則。
`, filepath.Join(repo, "CORE_BOOTSTRAP.md"),
		filepath.Join(repo, "README.md"),
		filepath.Join(repo, "ai-tools", "agent", "codex.md"),
		filepath.Join(repo, "runtime", "runtime.db"))
}

func initProjectGoalsReadmeContent(repo string) string {
	return fmt.Sprintf(`# Agent Goals

本目錄由 Ai-skill 對話目標帳本管理：%s
用 `+"`ai-skill goals`"+` 操作。

## 目前目標

（尚無 active goal）
`, filepath.Join(repo, "enforcement", "conversation-goal-ledger.md"))
}
