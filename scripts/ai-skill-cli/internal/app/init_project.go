package app

import (
	"encoding/json"
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
	private     bool
}

func runInitProject(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := newFlagSet("init-project", stderr)
	opts := initProjectOptions{}
	fs.StringVar(&opts.projectPath, "project", "", "target project directory")
	fs.StringVar(&opts.tools, "tools", "roo,cursor,claude,gemini,codex,copilot", "comma-separated tools: roo,cursor,claude,gemini,codex,copilot")
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
			Remediation: "Use one or more supported tools: roo,cursor,claude,gemini,codex,copilot.",
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
	supported := map[string]bool{"roo": true, "cursor": true, "claude": true, "gemini": true, "codex": true, "copilot": true}
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
			files = append(files,
				plannedFile{tool: tool, path: filepath.Join(target, "CLAUDE.md"), description: "Claude Code settings"},
				plannedFile{tool: tool, path: filepath.Join(target, ".claude", "settings.json"), description: "Claude Code hooks"},
			)
		case "gemini":
			files = append(files, plannedFile{tool: tool, path: filepath.Join(target, "GEMINI.md"), description: "Gemini CLI settings"})
		case "codex":
			files = append(files, plannedFile{tool: tool, path: filepath.Join(target, "AGENTS.md"), description: "Generic agent entry (AGENTS.md — Codex / Aider / Cline / other AGENTS.md-aware)"})
		case "copilot":
			files = append(files,
				plannedFile{tool: tool, path: filepath.Join(target, ".github", "copilot-instructions.md"), description: "GitHub Copilot project instructions"},
				plannedFile{tool: tool, path: filepath.Join(target, ".github", "instructions", "ai-skill.instructions.md"), description: "GitHub Copilot scoped instructions"},
			)
		}
	}
	files = append(files,
		plannedFile{tool: "common", path: filepath.Join(target, ".agent-goals", "README.md"), description: "agent goals ledger"},
		plannedFile{tool: "common", path: filepath.Join(target, ".ai-skill", ".gitignore"), description: "Ai-skill local config ignore rules"},
		plannedFile{tool: "common", path: filepath.Join(target, ".ai-skill", "local.env"), description: "Ai-skill local environment", private: true},
	)
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
	mode := os.FileMode(0o644)
	if filepath.Base(path) == "local.env" {
		mode = 0o600
	}
	return os.WriteFile(path, content, mode)
}

func initProjectFileContent(file plannedFile, repo string) (string, error) {
	switch file.tool {
	case "roo":
		return initProjectRooContent(repo)
	case "cursor":
		if strings.HasSuffix(file.path, "hooks.json") {
			return initProjectCursorHooksContent()
		}
		return initProjectCursorRuleContent(repo)
	case "claude":
		if strings.HasSuffix(file.path, filepath.Join(".claude", "settings.json")) {
			return initProjectClaudeSettingsContent(repo)
		}
		return initProjectClaudeContent(repo)
	case "gemini":
		return initProjectGeminiContent(repo)
	case "codex":
		return initProjectCodexContent(repo)
	case "copilot":
		if strings.HasSuffix(file.path, filepath.Join(".github", "instructions", "ai-skill.instructions.md")) {
			return initProjectCopilotScopedInstructionsContent(repo)
		}
		return initProjectCopilotInstructionsContent(repo)
	case "common":
		if strings.HasSuffix(file.path, filepath.Join(".ai-skill", ".gitignore")) {
			return initProjectLocalGitignoreContent()
		}
		if strings.HasSuffix(file.path, filepath.Join(".ai-skill", "local.env")) {
			return initProjectLocalEnvContent(repo)
		}
		return initProjectGoalsReadmeContent(repo)
	default:
		return "", fmt.Errorf("unsupported init-project template: %s", file.tool)
	}
}

const aiSkillRepoPlaceholder = "<AI_SKILL_REPO>"

func initProjectBootstrapText(repo string) string {
	return fmt.Sprintf(`本專案使用 Ai-skill 知識庫。請在本機設定環境變數 `+"`AI_SKILL_REPO`"+` 指向 Ai-skill repository；文件中以 `+"`%s`"+` 表示該路徑。

## 本機設定

不要把本機 Ai-skill repo 絕對路徑 commit 進專案。啟用順序：

1. 優先讀本機環境變數 `+"`AI_SKILL_REPO`"+`
2. 如果環境變數不存在，讀 `+"`.ai-skill/local.env`"+`（由 `+"`ai-skill init-project`"+` 產生，已被 `+"`.ai-skill/.gitignore`"+` 排除）
3. Claude hook 會再嘗試 $HOME/Documents/Ai-skill 等 portable fallback

每台機器可用自己的環境變數設定：

    # macOS / Linux / WSL
    export AI_SKILL_REPO=/path/to/ai-skill

    # Windows PowerShell
    [Environment]::SetEnvironmentVariable("AI_SKILL_REPO", "C:\path\to\Ai-skill", "User")

建議把 macOS / Linux / WSL 設定放在 shell profile（例如 ~/.zshrc、~/.bashrc），Windows 則用 User environment variable。不要提交 .env、.ai-skill/local.env 或任何含個人路徑的檔案。

## Pointer 展開規則

本檔不是完整規則，只是 bootstrap pointer。Agent 必須在處理任務前完成：

1. 讀 .ai-skill/local.env（若 AI_SKILL_REPO 尚未存在）並 resolve <AI_SKILL_REPO>
2. 讀 <AI_SKILL_REPO>/CORE_BOOTSTRAP.md
3. 讀 <AI_SKILL_REPO>/runtime/core-bootstrap.yaml
4. 依 canonical bootstrap contract 執行 Bootstrap Receipt、per-turn Cognitive Mode 報告與 close-loop obligations

不能只停在本檔摘要；必須 dereference pointer 到 canonical files。

> **IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
> 本檔為 **thin tool-entry pointer**。所有 obligation / format / enum / example 的 canonical 來源在 `+"`%s/runtime/core-bootstrap.yaml`"+`；`+"`%s/CORE_BOOTSTRAP.md`"+` 是 human-readable companion。Session 啟動 first turn 必須讀取 companion 並依 canonical contract 執行 Bootstrap Receipt、Cognitive Mode 報告與 close-loop obligations。
> Summary 的「Resume directly」是對話 framing，**不豁免** runtime / governance bootstrap。Resume / continuation session 同樣須走完 bootstrap。

## 啟動序列

1. 讀 `+"`%s/CORE_BOOTSTRAP.md`"+` — bootstrap companion 與進入點
2. 載入 `+"`%s/runtime/core-bootstrap.yaml`"+` — canonical bootstrap contract（required reads / per-session / per-turn obligations）

## 修改本檔的規則

本檔是 entry pointer，不是 canonical content。修改前先讀 `+"`%s/runtime/bootstrap-entry-points.yaml`"+`。新 obligation 加到 runtime/core-bootstrap.yaml 並 refresh runtime surfaces；工具差異放 ai-tools/agent/<tool>.md。不要把 obligation 全文複製進本檔。
`, aiSkillRepoPlaceholder,
		aiSkillRepoPlaceholder,
		aiSkillRepoPlaceholder,
		aiSkillRepoPlaceholder,
		aiSkillRepoPlaceholder,
		aiSkillRepoPlaceholder)
}

func initProjectRooContent(repo string) (string, error) {
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
`, text, text), nil
}

func initProjectCursorRuleContent(repo string) (string, error) {
	return fmt.Sprintf(`---
description: Ai-skill 知識庫啟動流程
globs:
alwaysApply: true
---

# Ai-skill Bootstrap

%s`, initProjectBootstrapText(repo)), nil
}

func initProjectCursorHooksContent() (string, error) {
	return `{
  "sessionStart": [
    {
      "description": "提醒載入 Ai-skill 知識庫",
      "command": "echo '提示：本專案使用 Ai-skill 知識庫，請確認已載入 CORE_BOOTSTRAP.md'"
    }
  ]
}
`, nil
}

func initProjectClaudeContent(repo string) (string, error) {
	return fmt.Sprintf("# Claude Code Auto-Bootstrap\n\n%s", initProjectBootstrapText(repo)), nil
}

func initProjectClaudeSettingsContent(repo string) (string, error) {
	command := func(event string) string {
		return fmt.Sprintf("sh -c 'ROOT=\"${CLAUDE_PROJECT_DIR:-$(pwd)}\"; if [ -z \"${AI_SKILL_REPO:-}\" ] && [ -f \"$ROOT/.ai-skill/local.env\" ]; then . \"$ROOT/.ai-skill/local.env\"; fi; if [ -z \"${AI_SKILL_REPO:-}\" ]; then for candidate in \"$HOME/Documents/Ai-skill\" \"$HOME/Ai-skill\" \"$PWD/../Ai-skill\"; do if [ -d \"$candidate/scripts/ai-skill-cli/bin\" ]; then AI_SKILL_REPO=\"$candidate\"; break; fi; done; fi; if [ -z \"${AI_SKILL_REPO:-}\" ]; then echo \"AI_SKILL_REPO is not set; skipping Ai-skill hook %s\" >&2; exit 0; fi; case \"$(uname -s 2>/dev/null | tr A-Z a-z)\" in darwin) os=darwin ;; linux) os=linux ;; mingw*|msys*|cygwin*) os=windows ;; *) os=unknown ;; esac; arch=\"$(uname -m 2>/dev/null || echo unknown)\"; case \"$arch\" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; esac; suffix=\"\"; [ \"$os\" = \"windows\" ] && suffix=\".exe\"; exec \"$AI_SKILL_REPO/scripts/ai-skill-cli/bin/ai-skill-$os-$arch$suffix\" hooks run %s --repo \"$AI_SKILL_REPO\"'", event, event)
	}
	settings := map[string]any{
		"description": "Claude Code project-local Ai-skill hooks. Commands execute the canonical Ai-skill repo-local Go binary and use CLAUDE_PROJECT_DIR as the target project root for nested Git reports.",
		"hooks": map[string]any{
			"SessionStart": []map[string]any{{
				"matcher": "startup|resume|clear",
				"hooks":   []map[string]any{{"type": "command", "command": command("session-start"), "timeout": 30}},
			}},
			"UserPromptSubmit": []map[string]any{{
				"matcher": "",
				"hooks":   []map[string]any{{"type": "command", "command": command("user-prompt-submit"), "timeout": 30}},
			}},
			"PreToolUse": []map[string]any{{
				"matcher": "",
				"hooks":   []map[string]any{{"type": "command", "command": command("pre-tool-use"), "timeout": 10}},
			}},
			"PostToolUse": []map[string]any{{
				"matcher": "",
				"hooks":   []map[string]any{{"type": "command", "command": command("post-tool-use"), "timeout": 10}},
			}},
			"Stop": []map[string]any{{
				"matcher": "",
				"hooks":   []map[string]any{{"type": "command", "command": command("stop"), "timeout": 10}},
			}},
		},
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func initProjectGeminiContent(repo string) (string, error) {
	return fmt.Sprintf("# Gemini CLI Auto-Bootstrap\n\n%s", initProjectBootstrapText(repo)), nil
}

func initProjectCodexContent(repo string) (string, error) {
	return fmt.Sprintf(`# AGENTS.md — Generic Agent Bootstrap Entry

本檔為 thin generic agent entry。適用 Codex、Cursor partial、Aider、Cline 等遵循 `+"`AGENTS.md`"+` 慣例的 AI agent。Canonical obligations 在 Ai-skill repo 的 `+"`CORE_BOOTSTRAP.md`"+` + `+"`runtime/core-bootstrap.yaml`"+`。

## 本機設定

不要把本機 Ai-skill repo 絕對路徑 commit 進專案。啟用順序：

1. 優先讀本機環境變數 `+"`AI_SKILL_REPO`"+`
2. 如果環境變數不存在，讀 `+"`.ai-skill/local.env`"+`（由 `+"`ai-skill init-project`"+` 產生，已被 `+"`.ai-skill/.gitignore`"+` 排除）
3. Claude hook 會再嘗試 $HOME/Documents/Ai-skill 等 portable fallback

每台機器可用自己的 `+"`AI_SKILL_REPO`"+` 環境變數設定：

    # macOS / Linux / WSL
    export AI_SKILL_REPO=/path/to/ai-skill

    # Windows PowerShell
    [Environment]::SetEnvironmentVariable("AI_SKILL_REPO", "C:\path\to\Ai-skill", "User")

文件中以 `+"`<AI_SKILL_REPO>`"+` 表示該路徑。建議把 macOS / Linux / WSL 設定放在 shell profile，Windows 則用 User environment variable。不要提交 .env、.ai-skill/local.env 或任何含個人路徑的檔案。

## Pointer 展開規則

本檔不是完整規則，只是 bootstrap pointer。Agent 必須在處理任務前完成：

1. 讀 .ai-skill/local.env（若 AI_SKILL_REPO 尚未存在）並 resolve <AI_SKILL_REPO>
2. 讀 <AI_SKILL_REPO>/CORE_BOOTSTRAP.md
3. 讀 <AI_SKILL_REPO>/runtime/core-bootstrap.yaml
4. 依 canonical bootstrap contract 執行 Bootstrap Receipt、per-turn Cognitive Mode 報告與 close-loop obligations

不能只停在本檔摘要；必須 dereference pointer 到 canonical files。

## 啟動序列

1. %s — 必讀規則 + 所有 obligations
2. %s — OS layout
3. %s — routing hub（選擇對應 AI 工具的 adapter）
4. %s — runtime canonical（phase / obligations / gates）

## 修改規則

不在本檔加單一工具規則 — 用 routing hub 指向的對應 `+"`ai-tools/agent/<tool>.md`"+`。本檔保持 thin。
`, aiSkillRepoPlaceholder+"/CORE_BOOTSTRAP.md",
		aiSkillRepoPlaceholder+"/README.md",
		aiSkillRepoPlaceholder+"/ai-tools/README.md",
		aiSkillRepoPlaceholder+"/runtime/runtime.db"), nil
}

func initProjectCopilotInstructionsContent(repo string) (string, error) {
	return fmt.Sprintf(`# GitHub Copilot Bootstrap Entry

本檔為 Copilot project-wide custom instructions 的 thin pointer。不要在此複製 bootstrap obligations、格式、enum、examples、goal ledger、close-loop 或 runtime phase 細節。

## Local Resolution

本專案由 `+"`ai-skill init-project`"+` 建立 `+"`.ai-skill/local.env`"+`。若 `+"`AI_SKILL_REPO`"+` 尚未存在，先讀該檔以 resolve `+"`%s`"+`，但不要把本機絕對路徑寫入此檔。

## Required Reads

1. `+"`%s/CORE_BOOTSTRAP.md`"+`
2. `+"`%s/runtime/core-bootstrap.yaml`"+`
3. `+"`%s/ai-tools/agent/copilot.md`"+`

依 canonical bootstrap contract 執行 required reads、Bootstrap Receipt、per-turn Cognitive Mode reporting 與 close-loop checks。若 Copilot 功能無法強制執行某項 gate，回報限制，並讓 repository hooks / CI / `+"`ai-skill runtime validate`"+` 作為 enforcement boundary。
`, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder), nil
}

func initProjectCopilotScopedInstructionsContent(repo string) (string, error) {
	return fmt.Sprintf(`---
applyTo: "**"
---

# Ai-skill Copilot Scoped Pointer

This scoped instruction file is a thin pointer for GitHub Copilot. Do not copy shared rules, bootstrap formats, runtime phase details, goal ledger procedures, or close-loop checklists here.

Before acting on this repository, read:

1. `+"`%s/CORE_BOOTSTRAP.md`"+`
2. `+"`%s/runtime/core-bootstrap.yaml`"+`
3. `+"`%s/ai-tools/agent/copilot.md`"+`

If Copilot cannot enforce a required runtime gate directly, report the limitation and rely on repository hooks, CI, and `+"`ai-skill runtime validate`"+` for enforcement.
`, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder), nil
}

func initProjectGoalsReadmeContent(repo string) (string, error) {
	return fmt.Sprintf(`# Agent Goals

本目錄由 Ai-skill 對話目標帳本管理：%s/enforcement/conversation-goal-ledger.md
用 `+"`ai-skill goals`"+` 操作。

## 目前目標

（尚無 active goal）
`, aiSkillRepoPlaceholder), nil
}

func initProjectLocalGitignoreContent() (string, error) {
	return `*
!.gitignore
`, nil
}

func initProjectLocalEnvContent(repo string) (string, error) {
	return fmt.Sprintf(`# Local Ai-skill bootstrap configuration.
# This file contains a machine-local path and is ignored by .ai-skill/.gitignore.
# Do not commit this file.
export AI_SKILL_REPO=%s
`, shellSingleQuote(repo)), nil
}

func shellSingleQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
