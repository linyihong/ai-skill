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
	tool             string
	path             string
	description      string
	private          bool
	preserveExisting bool
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
			if file.preserveExisting {
				continue
			}
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
		if file.preserveExisting {
			if _, err := os.Stat(file.path); err == nil {
				result.Mutations = append(result.Mutations, fmt.Sprintf("preserved existing %s: %s", file.description, file.path))
				continue
			}
		}
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
				plannedFile{tool: tool, path: filepath.Join(target, ".copilot", "README.md"), description: "GitHub Copilot guided startup README"},
				plannedFile{tool: tool, path: filepath.Join(target, ".copilot", "bootstrap-prompt.md"), description: "GitHub Copilot bootstrap prompt"},
				plannedFile{tool: tool, path: filepath.Join(target, ".copilot", "start-copilot.sh"), description: "GitHub Copilot Unix startup wrapper"},
			)
		}
	}
	files = append(files,
		plannedFile{tool: "common", path: filepath.Join(target, ".agent-goals", "README.md"), description: "agent goals ledger"},
		plannedFile{tool: "common", path: filepath.Join(target, ".ai-skill", ".gitignore"), description: "Ai-skill local config ignore rules"},
		plannedFile{tool: "common", path: filepath.Join(target, ".ai-skill", "project", "README.md"), description: "Ai-skill project overlay index", preserveExisting: true},
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
	if filepath.Ext(path) == ".sh" {
		mode = 0o755
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
		if strings.HasSuffix(file.path, filepath.Join(".copilot", "README.md")) {
			return initProjectCopilotReadmeContent(repo)
		}
		if strings.HasSuffix(file.path, filepath.Join(".copilot", "bootstrap-prompt.md")) {
			return initProjectCopilotBootstrapPromptContent(repo)
		}
		if strings.HasSuffix(file.path, filepath.Join(".copilot", "start-copilot.sh")) {
			return initProjectCopilotStartScriptContent(repo)
		}
		return initProjectCopilotInstructionsContent(repo)
	case "common":
		if strings.HasSuffix(file.path, filepath.Join(".ai-skill", ".gitignore")) {
			return initProjectLocalGitignoreContent()
		}
		if strings.HasSuffix(file.path, filepath.Join(".ai-skill", "project", "README.md")) {
			return initProjectOverlayReadmeContent(repo)
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
	return fmt.Sprintf(`本檔是 thin tool-entry pointer，不是 canonical 規則本體。

Resolve `+"`AI_SKILL_REPO`"+`：優先用環境變數；否則讀 project-local `+"`.ai-skill/local.env`"+`。不要把本機絕對路徑寫入可 commit 檔案。

**IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
1. 讀 `+"`%s/CORE_BOOTSTRAP.md`"+`
2. 讀 `+"`%s/runtime/core-bootstrap.yaml`"+`
3. 依 canonical contract 執行 Bootstrap Receipt、final Cognitive Mode 報告與 close-loop obligations。
4. 若 `+"`<PROJECT_ROOT>/.ai-skill/project/README.md`"+` 存在，bootstrap 後讀取 project overlay index。

Summary 的「Resume directly」不豁免 bootstrap；resume / continuation session 同樣必須 dereference canonical files。

修改本檔前先讀 `+"`%s/runtime/bootstrap-entry-points.yaml`"+`。新 obligation 加到 runtime/core-bootstrap.yaml；工具差異放 `+"`%s/ai-tools/agent/<tool>.md`"+`。
`, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder)
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
	// stop: when AI_SKILL_REPO is unresolvable the close-out hook fails CLOSED
	// (exit 2 + failClosed) — a missing repo must not let a session end without
	// the Cognitive Mode / Project Git Report check.
	stopCommand := "sh -c 'ROOT=\"${CURSOR_PROJECT_DIR:-$(pwd)}\"; if [ -z \"${AI_SKILL_REPO:-}\" ] && [ -f \"$ROOT/.ai-skill/local.env\" ]; then . \"$ROOT/.ai-skill/local.env\"; fi; if [ -z \"${AI_SKILL_REPO:-}\" ]; then echo \"AI_SKILL_REPO is not set; Ai-skill Cursor close-out hook cannot validate Cognitive Mode\" >&2; exit 2; fi; case \"$(uname -s 2>/dev/null | tr A-Z a-z)\" in darwin) os=darwin ;; linux) os=linux ;; mingw*|msys*|cygwin*) os=windows ;; *) os=unknown ;; esac; arch=\"$(uname -m 2>/dev/null || echo unknown)\"; case \"$arch\" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; esac; suffix=\"\"; [ \"$os\" = \"windows\" ] && suffix=\".exe\"; exec \"$AI_SKILL_REPO/scripts/ai-skill-cli/bin/ai-skill-$os-$arch$suffix\" hooks run stop --repo \"$ROOT\"'"
	// preToolUse: a blockable Cursor event. When AI_SKILL_REPO is unresolvable the
	// gate fails OPEN (exit 0) — a missing repo must NOT block every tool call.
	// The binary itself emits {"permission":"deny"} on a real gate hit; allow is
	// empty stdout + exit 0.
	preToolUseCommand := "sh -c 'ROOT=\"${CURSOR_PROJECT_DIR:-$(pwd)}\"; if [ -z \"${AI_SKILL_REPO:-}\" ] && [ -f \"$ROOT/.ai-skill/local.env\" ]; then . \"$ROOT/.ai-skill/local.env\"; fi; if [ -z \"${AI_SKILL_REPO:-}\" ]; then echo \"AI_SKILL_REPO is not set; skipping Ai-skill Cursor pre-tool-use gate\" >&2; exit 0; fi; case \"$(uname -s 2>/dev/null | tr A-Z a-z)\" in darwin) os=darwin ;; linux) os=linux ;; mingw*|msys*|cygwin*) os=windows ;; *) os=unknown ;; esac; arch=\"$(uname -m 2>/dev/null || echo unknown)\"; case \"$arch\" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; esac; suffix=\"\"; [ \"$os\" = \"windows\" ] && suffix=\".exe\"; exec \"$AI_SKILL_REPO/scripts/ai-skill-cli/bin/ai-skill-$os-$arch$suffix\" hooks run pre-tool-use --repo \"$ROOT\"'"
	settings := map[string]any{
		"version": 1,
		"hooks": map[string]any{
			"sessionStart": []map[string]any{{
				"type":    "prompt",
				"prompt":  "This project uses Ai-skill. Before answering, load <AI_SKILL_REPO>/CORE_BOOTSTRAP.md and <AI_SKILL_REPO>/runtime/core-bootstrap.yaml, emit the Bootstrap Receipt, and include Cognitive Mode only in the final response / session close-out. If dirty root or nested Git repositories are detected at close-out, the final response must include a combined ### Project Git Report.",
				"timeout": 10,
			}},
			"preToolUse": []map[string]any{{
				"command": preToolUseCommand,
				"timeout": 10,
			}},
			"stop": []map[string]any{{
				"command":    stopCommand,
				"timeout":    10,
				"failClosed": true,
				"loop_limit": 2,
			}},
		},
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func initProjectClaudeContent(repo string) (string, error) {
	return fmt.Sprintf("# Claude Code Auto-Bootstrap\n\n%s", initProjectBootstrapText(repo)), nil
}

func initProjectClaudeSettingsContent(repo string) (string, error) {
	command := func(event string) string {
		return fmt.Sprintf("sh -c 'ROOT=\"${CLAUDE_PROJECT_DIR:-$(pwd)}\"; if [ -z \"${AI_SKILL_REPO:-}\" ] && [ -f \"$ROOT/.ai-skill/local.env\" ]; then . \"$ROOT/.ai-skill/local.env\"; fi; if [ -z \"${AI_SKILL_REPO:-}\" ]; then for candidate in \"$HOME/Documents/Ai-skill\" \"$HOME/Ai-skill\" \"$PWD/../Ai-skill\"; do if [ -d \"$candidate/scripts/ai-skill-cli/bin\" ]; then AI_SKILL_REPO=\"$candidate\"; break; fi; done; fi; if [ -z \"${AI_SKILL_REPO:-}\" ]; then echo \"AI_SKILL_REPO is not set; skipping Ai-skill hook %s\" >&2; exit 0; fi; case \"$(uname -s 2>/dev/null | tr A-Z a-z)\" in darwin) os=darwin ;; linux) os=linux ;; mingw*|msys*|cygwin*) os=windows ;; *) os=unknown ;; esac; arch=\"$(uname -m 2>/dev/null || echo unknown)\"; case \"$arch\" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; esac; suffix=\"\"; [ \"$os\" = \"windows\" ] && suffix=\".exe\"; exec \"$AI_SKILL_REPO/scripts/ai-skill-cli/bin/ai-skill-$os-$arch$suffix\" hooks run %s --repo \"$ROOT\"'", event, event)
	}
	settings := map[string]any{
		"description": "Claude Code project-local Ai-skill hooks. Commands execute the canonical Ai-skill repo-local Go binary and use CLAUDE_PROJECT_DIR as the target project root for final Cognitive Mode and nested Git reports.",
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

本檔是 thin generic agent entry，不是 canonical 規則本體。適用 Codex、Cursor partial、Aider、Cline 等 `+"`AGENTS.md`"+`-aware tools。

Resolve `+"`AI_SKILL_REPO`"+`：優先用環境變數；否則讀 project-local `+"`.ai-skill/local.env`"+`。不要把本機絕對路徑寫入可 commit 檔案。

**IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
1. 讀 %s
2. 讀 %s
3. 讀 %s 選擇 active tool adapter
4. 若 <PROJECT_ROOT>/.ai-skill/project/README.md 存在，bootstrap 後讀取 project overlay index。
5. 依 canonical contract 執行 Bootstrap Receipt、final Cognitive Mode 報告與 close-loop obligations。

Summary 的「Resume directly」不豁免 bootstrap。修改本檔前先讀 %s；工具差異放 `+"`ai-tools/agent/<tool>.md`"+`。
`, aiSkillRepoPlaceholder+"/CORE_BOOTSTRAP.md",
		aiSkillRepoPlaceholder+"/runtime/core-bootstrap.yaml",
		aiSkillRepoPlaceholder+"/ai-tools/README.md",
		aiSkillRepoPlaceholder+"/runtime/bootstrap-entry-points.yaml"), nil
}

func initProjectCopilotInstructionsContent(repo string) (string, error) {
	return fmt.Sprintf(`# GitHub Copilot Bootstrap Entry

本檔為 Copilot project-wide custom instructions 的 thin pointer。不要在此複製 bootstrap obligations、格式、enum、examples、goal ledger、close-loop 或 runtime phase 細節。

## Mandatory Startup

在回覆任何使用者請求前，必須先完成下方 Required Reads 並依 canonical bootstrap contract 回報 Bootstrap Receipt。不得因為請求看似只是列檔、簡單查詢、read-only、說明原因或已由對話摘要提供 context 而跳過 bootstrap。

## Local Resolution

本專案由 `+"`ai-skill init-project`"+` 建立 `+"`.ai-skill/local.env`"+`。若 `+"`AI_SKILL_REPO`"+` 尚未存在，先讀該檔以 resolve `+"`%s`"+`，但不要把本機絕對路徑寫入此檔。

## Required Reads

1. `+"`%s/CORE_BOOTSTRAP.md`"+`
2. `+"`%s/runtime/core-bootstrap.yaml`"+`
3. `+"`%s/ai-tools/agent/copilot.md`"+`

依 canonical bootstrap contract 執行 required reads、Bootstrap Receipt、final close-out Cognitive Mode reporting 與 close-loop checks。若 Copilot 功能無法強制執行某項 gate，回報限制，並讓 repository hooks / CI / `+"`ai-skill runtime validate`"+` 作為 enforcement boundary。
`, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder), nil
}

func initProjectCopilotScopedInstructionsContent(repo string) (string, error) {
	return fmt.Sprintf(`---
applyTo: "**"
---

# Ai-skill Copilot Scoped Pointer

This scoped instruction file is a thin pointer for GitHub Copilot. Do not copy shared rules, bootstrap formats, runtime phase details, goal ledger procedures, or close-loop checklists here.

Before answering any user request in this repository, read these files and complete the canonical Bootstrap Receipt. Do not treat simple file listings, read-only queries, explanation-only requests, or resumed context as exemptions:

1. `+"`%s/CORE_BOOTSTRAP.md`"+`
2. `+"`%s/runtime/core-bootstrap.yaml`"+`
3. `+"`%s/ai-tools/agent/copilot.md`"+`

If Copilot cannot enforce a required runtime gate directly, report the limitation and rely on repository hooks, CI, and `+"`ai-skill runtime validate`"+` for enforcement.
`, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder), nil
}

func initProjectCopilotReadmeContent(repo string) (string, error) {
	return fmt.Sprintf(`# GitHub Copilot Guided Startup

Copilot custom instructions are advisory. For a new session, start with the generated bootstrap prompt in [`+"`bootstrap-prompt.md`"+`](bootstrap-prompt.md), or run:

`+"```bash"+`
.copilot/start-copilot.sh
`+"```"+`

The script is a temporary bootstrap wrapper. It only resolves `+"`AI_SKILL_REPO`"+` / `+"`.ai-skill/local.env`"+` and calls the repo-local `+"`ai-skill copilot start`"+` command. Deletion condition: remove this wrapper after Copilot supports a native automatic bootstrap-on-session setting, or after the project adopts a non-shell launcher with equivalent behavior.

Hard enforcement still lives in repository hooks, CI, and `+"`ai-skill runtime validate`"+`; this startup package only reduces new-session bootstrap misses.

Canonical sources:

1. `+"`%s/CORE_BOOTSTRAP.md`"+`
2. `+"`%s/runtime/core-bootstrap.yaml`"+`
3. `+"`%s/ai-tools/agent/copilot.md`"+`
`, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder, aiSkillRepoPlaceholder), nil
}

func initProjectCopilotBootstrapPromptContent(repo string) (string, error) {
	return copilotBootstrapPrompt("<PROJECT_ROOT>", aiSkillRepoPlaceholder) + "\n", nil
}

func initProjectCopilotStartScriptContent(repo string) (string, error) {
	return `#!/bin/sh
set -eu

# Temporary bootstrap wrapper.
# Deletion condition: remove after Copilot supports native automatic bootstrap-on-session
# or after this project adopts a non-shell launcher with equivalent behavior.
# This wrapper only locates and invokes the repo-local ai-skill binary.

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
if [ -z "${AI_SKILL_REPO:-}" ] && [ -f "$ROOT/.ai-skill/local.env" ]; then
  . "$ROOT/.ai-skill/local.env"
fi
if [ -z "${AI_SKILL_REPO:-}" ]; then
  echo "AI_SKILL_REPO is not set; configure it or run ai-skill init-project again." >&2
  exit 1
fi

case "$(uname -s 2>/dev/null | tr A-Z a-z)" in
  darwin) os=darwin ;;
  linux) os=linux ;;
  mingw*|msys*|cygwin*) os=windows ;;
  *) os=unknown ;;
esac
arch="$(uname -m 2>/dev/null || echo unknown)"
case "$arch" in
  arm64|aarch64) arch=arm64 ;;
  x86_64|amd64) arch=amd64 ;;
esac
suffix=""
[ "$os" = "windows" ] && suffix=".exe"

exec "$AI_SKILL_REPO/scripts/ai-skill-cli/bin/ai-skill-$os-$arch$suffix" copilot start --project "$ROOT" "$@"
`, nil
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
	return `local.env
`, nil
}

func initProjectOverlayReadmeContent(repo string) (string, error) {
	return `# Project Ai-skill Overlay

This directory is reserved for project-local Ai-skill overlay rules. Put project-specific rules in ` + "`rules/`" + ` and link to them from this README.

Tool-specific bootstrap entries should load this index instead of duplicating project rule bodies. Do not add separate project-rule ` + "`.mdc`" + ` files under ` + "`.cursor/rules/`" + `; keep Cursor on ` + "`.cursor/rules/ai-skill-bootstrap.mdc`" + ` so validation exercises the shared overlay.
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
