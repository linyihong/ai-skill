package app

import (
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
)

type hooksOptions struct {
	command     string
	repoPath    string
	dryRun      bool
	force       bool
	jsonOutput  bool
	plainOutput bool
	positional  []string
}

func runHooks(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill hooks <install|run> [flags]")
		return ExitInvalidUsage
	}
	opts := hooksOptions{command: args[0]}
	if opts.command != "install" && opts.command != "run" {
		_, _ = fmt.Fprintf(stderr, "unsupported hooks command: %s\n", opts.command)
		return ExitInvalidUsage
	}

	name := "hooks " + opts.command
	if opts.command == "run" {
		if len(args) < 2 {
			_, _ = fmt.Fprintln(stderr, "usage: ai-skill hooks run <pre-commit|commit-msg|post-commit|pre-push|session-start|pre-tool-use|post-tool-use|user-prompt-submit|stop> [flags]")
			return ExitInvalidUsage
		}
		opts.command = "run " + args[1]
		args = append([]string{"run"}, args[2:]...)
	}
	fs := newFlagSet(name, stderr)
	fs.StringVar(&opts.repoPath, "repo", ".", "repository path")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview hook installation without writing")
	fs.BoolVar(&opts.force, "force", false, "allow overwriting existing hook targets")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	opts.positional = fs.Args()
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	// Claude Code hooks: bypass git/repo checks and the Result machinery.
	// They write raw JSON or plain text to stdout/stderr and return exit codes directly.
	// Project directory comes from CLAUDE_PROJECT_DIR env var (set by Claude Code).
	if strings.HasPrefix(opts.command, "run ") {
		projectDir := os.Getenv("CLAUDE_PROJECT_DIR")
		if projectDir == "" {
			projectDir = opts.repoPath
			if projectDir == "." {
				if cwd, err := os.Getwd(); err == nil {
					projectDir = cwd
				}
			}
		}
		switch opts.command {
		case "run session-start":
			return runSessionStartHook(projectDir, stdout, stderr)
		case "run pre-tool-use":
			return runPreToolUseHook(projectDir, stdout, stderr)
		case "run post-tool-use":
			return runPostToolUseHook(projectDir, stdout, stderr)
		case "run user-prompt-submit":
			return runUserPromptSubmitHook(projectDir, stdout, stderr)
		case "run stop":
			return runStopHook(projectDir, stdout, stderr)
		}
		// Fall through to buildHooksRunResult for git hooks (pre-commit, commit-msg, etc.)
	}

	var result Result
	if strings.HasPrefix(opts.command, "run ") {
		result = buildHooksRunResult(opts)
	} else {
		result = buildHooksInstallResult(opts)
	}
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

func buildHooksInstallResult(opts hooksOptions) Result {
	result := Result{
		Command:        "hooks install",
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

	gitCheck := checkGit()
	result.Checks = append(result.Checks, gitCheck)
	if gitCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{
			Code:        "missing_git",
			Message:     "Git is required for hooks install but was not found in PATH.",
			Remediation: "Install Git and ensure the git executable is available in PATH.",
		}
		return result
	}

	root, repoCheck := closeLoopRepoRoot(opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message, Remediation: "Pass --repo with a valid Git working tree."}
		return result
	}

	operationCheck := closeLoopUnsafeStateCheck(root)
	if operationCheck.Status != "ok" {
		operationCheck.Status = "warning"
		operationCheck.Remediation = "Hook installation planning is allowed, but commit or push must not be triggered while Git has an active operation."
		result.Checks = append(result.Checks, operationCheck)
	} else {
		result.Checks = append(result.Checks, Check{Name: "git_operation", Status: "ok", Message: "no merge, rebase, or cherry-pick state detected"})
	}

	targetDir, targetCheck := gitHooksTargetDir(root)
	result.Checks = append(result.Checks, targetCheck)
	if targetCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "missing_hooks_target", Message: targetCheck.Message, Remediation: "Use a normal Git working tree with a .git/hooks directory."}
		return result
	}

	hooks := []string{"pre-commit", "commit-msg", "post-commit", "pre-push"}
	for _, hook := range hooks {
		result.PlannedActions = append(result.PlannedActions, fmt.Sprintf("install Go hook adapter: %s", filepath.Join(targetDir, hook)))
	}
	conflicts := existingHookTargets(targetDir, hooks)
	if len(conflicts) > 0 && !opts.force {
		result.Checks = append(result.Checks, Check{Name: "hook_conflicts", Status: "failed", Message: strings.Join(conflicts, ", ")})
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "target_exists", Message: "hook targets already exist: " + strings.Join(conflicts, ", "), Remediation: "Pass --force only after reviewing the planned overwrite list."}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "hook_conflicts", Status: "ok", Message: "no blocking hook target conflicts"})
	if !opts.dryRun {
		for _, hook := range hooks {
			target := filepath.Join(targetDir, hook)
			if err := os.WriteFile(target, []byte(hookAdapterContent(hook)), 0o755); err != nil {
				result.Status = "blocked"
				result.ExitCode = ExitGeneralFailure
				result.Error = &CommandError{Code: "hook_install_failed", Message: err.Error()}
				return result
			}
			result.Mutations = append(result.Mutations, "installed hook adapter: "+target)
		}
	}
	return result
}

func buildHooksRunResult(opts hooksOptions) Result {
	result := Result{
		Command:        "hooks " + opts.command,
		Mode:           "write",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	gitCheck := checkGit()
	result.Checks = append(result.Checks, gitCheck)
	if gitCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_git", Message: "Git is required for hook runner but was not found in PATH."}
		return result
	}
	root, repoCheck := closeLoopRepoRoot(opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message}
		return result
	}
	switch opts.command {
	case "run pre-commit":
		return runPreCommitHook(result, root)
	case "run commit-msg":
		return runCommitMsgHook(result, root, opts.positional)
	case "run post-commit":
		if os.Getenv("AI_SKILL_SYNC_CURSOR_BUNDLE") == "1" {
			result.Checks = append(result.Checks, Check{Name: "cursor_bundle_sync", Status: "skipped", Message: "Go sync-cursor-bundle write mode is not enabled"})
		} else {
			result.Checks = append(result.Checks, Check{Name: "cursor_bundle_sync", Status: "skipped", Message: "reference-only default"})
		}
		return result
	case "run pre-push":
		return runPrePushHook(result, root)
	default:
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_hook", Message: "unsupported hook runner: " + opts.command}
		return result
	}
}

func runPreCommitHook(result Result, root string) Result {
	staged, err := gitLines(root, "diff", "--cached", "--name-only")
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitGeneralFailure
		result.Error = &CommandError{Code: "staged_lookup_failed", Message: err.Error()}
		return result
	}
	if hasRuntimeSourceChange(staged) {
		var stdout strings.Builder
		var stderr strings.Builder
		code := Run([]string{"runtime", "compile", "--repo", root, "--json"}, &stdout, &stderr)
		if code != ExitSuccess {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "runtime_compile_failed", Message: compactSummary(stdout.String() + stderr.String())}
			return result
		}
		if _, err := exec.Command("git", "-C", root, "add", filepath.Join(root, "runtime", "runtime.db")).CombinedOutput(); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "runtime_db_stage_failed", Message: err.Error()}
			return result
		}
		result.Mutations = append(result.Mutations, "compiled and staged runtime/runtime.db")
	}
	if hasRuntimeDBChange(staged) || hasKnowledgeValidationChange(staged) {
		var stdout strings.Builder
		var stderr strings.Builder
		code := Run([]string{"runtime", "validate", "--repo", root, "--json"}, &stdout, &stderr)
		if code != ExitSuccess {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "runtime_validate_failed", Message: compactSummary(stdout.String() + stderr.String())}
			return result
		}
		result.Checks = append(result.Checks, Check{Name: "runtime_validation", Status: "ok", Message: "staged runtime/knowledge/validation changes passed"})
	}
	// Gate: new shell scripts must be implemented in Go instead (cross-platform policy).
	// Detects only newly Added (.sh) files; modifications to existing scripts are allowed
	// while they await migration (mark with [skip-go-migration] to suppress).
	if msg := validateNoNewShellScripts(root, staged); msg != "" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{
			Code:        "new_shell_script_forbidden",
			Message:     msg,
			Remediation: "Implement the hook/script logic in Go (scripts/ai-skill-cli/internal/app/hooks.go or a new subcommand). See plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md §platform-governance and enforcement/failure-patterns/shell-script-added-without-go-migration.md. To suppress for a transitional .sh wrapper, add '[skip-go-migration]' as a standalone line in the commit message body.",
		}
		return result
	}

	if len(result.Mutations) == 0 {
		result.Checks = append(result.Checks, Check{Name: "pre_commit", Status: "ok", Message: "no runtime or knowledge hook action required"})
	}
	return result
}

// validateNoNewShellScripts returns a non-empty error message if any newly Added
// .sh files appear in the staged set, enforcing the cross-platform Go-first policy.
// Modifications to existing .sh files are allowed (they are pending migration).
// Opt-out: '[skip-go-migration]' standalone line in commit message body.
func validateNoNewShellScripts(root string, staged []string) string {
	// Check for opt-out marker in commit message (COMMIT_EDITMSG)
	msgPath := filepath.Join(root, ".git", "COMMIT_EDITMSG")
	if data, err := os.ReadFile(msgPath); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.TrimSpace(line) == "[skip-go-migration]" {
				return ""
			}
		}
	}

	// Get only Added (new) files from staged set
	added, err := gitLines(root, "diff", "--cached", "--diff-filter=A", "--name-only")
	if err != nil {
		return "" // fail-open: don't block on git error
	}
	var newShells []string
	for _, f := range added {
		if strings.HasSuffix(f, ".sh") {
			newShells = append(newShells, f)
		}
	}
	if len(newShells) == 0 {
		return ""
	}
	return "new shell script(s) staged: " + strings.Join(newShells, ", ") +
		" — cross-platform policy requires Go implementation instead of .sh"
}

// ---------------------------------------------------------------------------
// Claude Code hook helpers
// ---------------------------------------------------------------------------

// appendLog appends a line to a diagnostic log file (best-effort, no-op on error).
func appendLog(path, msg string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = fmt.Fprintln(f, msg)
}

// readFileSafe returns file contents or a "(missing: <path>)" placeholder.
func readFileSafe(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "(missing: " + path + ")"
	}
	return string(data)
}

func resolveClaudeAiSkillRepo(projectDir string) string {
	candidates := []string{}
	if env := strings.TrimSpace(os.Getenv("AI_SKILL_REPO")); env != "" {
		candidates = append(candidates, env)
	}
	candidates = append(candidates, projectDir)
	if cwd, err := os.Getwd(); err == nil {
		for {
			candidates = append(candidates, cwd)
			parent := filepath.Dir(cwd)
			if parent == cwd {
				break
			}
			cwd = parent
		}
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(candidate, "CORE_BOOTSTRAP.md")); err == nil {
			if _, err := os.Stat(filepath.Join(candidate, "runtime", "core-bootstrap.yaml")); err == nil {
				return candidate
			}
		}
	}
	return projectDir
}

// md5Short returns the first 12 hex chars of the MD5 hash of s.
func md5Short(s string) string {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", h)[:12]
}

// claudeFileExists reports whether path exists on disk.
func claudeFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// extractAssistantTexts scans a JSONL transcript file and returns all assistant
// message text bodies in document order.  Each element is the joined text of
// one assistant turn.
func extractAssistantTexts(transcriptPath string) []string {
	f, err := os.Open(transcriptPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var results []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 2*1024*1024), 2*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		roleField := entry["type"]
		if roleField == nil {
			roleField = entry["role"]
		}
		var role string
		if roleField != nil {
			_ = json.Unmarshal(roleField, &role)
		}
		if role != "assistant" {
			continue
		}
		var chunks []string
		if msgRaw, ok := entry["message"]; ok {
			var msg map[string]json.RawMessage
			if err := json.Unmarshal(msgRaw, &msg); err == nil {
				if cRaw, ok := msg["content"]; ok {
					var s string
					if err := json.Unmarshal(cRaw, &s); err == nil {
						chunks = append(chunks, s)
					} else {
						var items []json.RawMessage
						if err := json.Unmarshal(cRaw, &items); err == nil {
							for _, item := range items {
								var m map[string]json.RawMessage
								if err := json.Unmarshal(item, &m); err == nil {
									if tRaw, ok := m["text"]; ok {
										var t string
										if err := json.Unmarshal(tRaw, &t); err == nil {
											chunks = append(chunks, t)
										}
									}
								} else {
									var s string
									if err := json.Unmarshal(item, &s); err == nil {
										chunks = append(chunks, s)
									}
								}
							}
						}
					}
				}
			}
		} else if cRaw, ok := entry["content"]; ok {
			var s string
			if err := json.Unmarshal(cRaw, &s); err == nil {
				chunks = append(chunks, s)
			}
		}
		if len(chunks) > 0 {
			results = append(results, strings.Join(chunks, "\n"))
		}
	}
	return results
}

type gitRepoStatusReport struct {
	Root  string
	Rel   string
	Lines []string
}

func collectDirtyGitRepoReports(projectDir string) []gitRepoStatusReport {
	seen := map[string]bool{}
	var repos []string
	if root, err := hookGitOutput(projectDir, "rev-parse", "--show-toplevel"); err == nil {
		root = strings.TrimSpace(root)
		if root != "" {
			repos = append(repos, root)
			seen[root] = true
		}
	}
	_ = filepath.WalkDir(projectDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			switch name {
			case ".git":
				root := filepath.Dir(path)
				if !seen[root] {
					repos = append(repos, root)
					seen[root] = true
				}
				return filepath.SkipDir
			case "node_modules", ".cache", "tmp", "vendor":
				return filepath.SkipDir
			}
			return nil
		}
		if name == ".git" {
			root := filepath.Dir(path)
			if !seen[root] {
				repos = append(repos, root)
				seen[root] = true
			}
		}
		return nil
	})
	sort.Strings(repos)

	var reports []gitRepoStatusReport
	for _, repo := range repos {
		out, err := hookGitOutput(repo, "status", "--short", "--branch")
		if err != nil {
			continue
		}
		lines := nonEmptyLines(out)
		if !gitStatusNeedsReport(lines) {
			continue
		}
		rel, err := filepath.Rel(projectDir, repo)
		if err != nil || rel == "." {
			rel = filepath.Base(repo)
		}
		reports = append(reports, gitRepoStatusReport{Root: repo, Rel: filepath.ToSlash(rel), Lines: lines})
	}
	return reports
}

func nonEmptyLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func gitStatusNeedsReport(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	if len(lines) > 1 {
		return true
	}
	branch := lines[0]
	return strings.Contains(branch, "ahead") ||
		strings.Contains(branch, "behind") ||
		strings.Contains(branch, "gone") ||
		strings.Contains(branch, "diverged")
}

func formatDirtyGitRepoReport(projectDir string) string {
	reports := collectDirtyGitRepoReports(projectDir)
	if len(reports) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("[ai-skill nested git report]\n")
	b.WriteString("Dirty Git repositories were detected under the project root. The final response MUST include a combined `### Project Git Report` section. ")
	b.WriteString("If one repo changed, report that repo; if multiple repos changed, merge them into one section with one bullet per repo. ")
	b.WriteString("Do not claim a clean close-loop until every listed repo is handled or explicitly marked as pre-existing/unrelated.\n\n")
	for _, report := range reports {
		b.WriteString("- ")
		b.WriteString(report.Rel)
		b.WriteString("\n")
		limit := len(report.Lines)
		if limit > 12 {
			limit = 12
		}
		for i := 0; i < limit; i++ {
			b.WriteString("  ")
			b.WriteString(report.Lines[i])
			b.WriteByte('\n')
		}
		if len(report.Lines) > limit {
			b.WriteString("  ... ")
			b.WriteString(itoa(len(report.Lines) - limit))
			b.WriteString(" more status lines\n")
		}
	}
	return b.String()
}

func hookGitOutput(root string, args ...string) (string, error) {
	output, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// ---------------------------------------------------------------------------
// Claude Code hooks (SessionStart / PreToolUse / PostToolUse / UserPromptSubmit / Stop)
// Cross-platform Go replacement for .claude/hooks/*.sh per
// plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md.
// ---------------------------------------------------------------------------

// runSessionStartHook implements the Claude Code SessionStart hook.
// Queries runtime.db for phase/obligation/gate counts, reads 4 bootstrap files,
// emits hookSpecificOutput JSON, and writes a TTL flag for runPreToolUseHook.
func runSessionStartHook(projectDir string, stdout io.Writer, stderr io.Writer) int {
	const logFile = "/tmp/ai-skill-sessionstart-hook.log"
	ts := time.Now().Format("2006-01-02T15:04:05")
	appendLog(logFile, "=== "+ts+" SessionStart hook fired (Go) ===")

	aiSkillRepo := resolveClaudeAiSkillRepo(projectDir)
	phase, obligCount, gateCount := "unknown", "?", "?"
	dbPath := filepath.Join(aiSkillRepo, "runtime", "runtime.db")
	if db, err := sql.Open("sqlite", dbPath); err == nil {
		defer db.Close()
		if row := db.QueryRow("SELECT phase_id FROM phase_machine LIMIT 1"); row != nil {
			_ = row.Scan(&phase)
		}
		if row := db.QueryRow("SELECT COUNT(*) FROM obligations"); row != nil {
			_ = row.Scan(&obligCount)
		}
		if row := db.QueryRow("SELECT COUNT(*) FROM gates"); row != nil {
			_ = row.Scan(&gateCount)
		}
	}
	const activePerTurn = "obligation.cognitive.mode_report, obligation.finality.close_loop_check"

	coreBootstrap := readFileSafe(filepath.Join(aiSkillRepo, "CORE_BOOTSTRAP.md"))
	ruleWeight := readFileSafe(filepath.Join(aiSkillRepo, "enforcement", "rule-weight.md"))
	dependency := readFileSafe(filepath.Join(aiSkillRepo, "enforcement", "dependency-reading.md"))
	goalLedger := readFileSafe(filepath.Join(aiSkillRepo, "enforcement", "conversation-goal-ledger.md"))

	context := fmt.Sprintf(
		"[ai-skill SessionStart] Bootstrap auto-loaded. The agent does NOT need to read these files again "+
			"— they are already in context. Your first user-facing response MUST begin with this Bootstrap "+
			"Receipt (verbatim), then proceed to answer the user:\n\n"+
			"Bootstrap: rules=✓ phase=%s obligations=%s gates=%s\n"+
			"Active per-turn obligations: %s\n\n"+
			"Final response MUST also end with a Cognitive Mode 報告 block (compact form is fine for trivial "+
			"tasks). Per-turn enforcement: see runtime/core-bootstrap.yaml §per_turn_obligations.\n\n"+
			"--- CORE_BOOTSTRAP.md (companion) ---\n%s\n\n"+
			"--- enforcement/rule-weight.md ---\n%s\n\n"+
			"--- enforcement/dependency-reading.md ---\n%s\n\n"+
			"--- enforcement/conversation-goal-ledger.md ---\n%s",
		phase, obligCount, gateCount, activePerTurn,
		coreBootstrap, ruleWeight, dependency, goalLedger,
	)

	output := map[string]interface{}{
		"hookSpecificOutput": map[string]interface{}{
			"hookEventName":     "SessionStart",
			"additionalContext": context,
		},
	}
	if err := json.NewEncoder(stdout).Encode(output); err != nil {
		_, _ = fmt.Fprintln(stderr, "encode error:", err)
		return ExitGeneralFailure
	}

	projectHash := md5Short(projectDir)
	flagFile := "/tmp/ai-skill-sessionstart-" + projectHash + ".flag"
	_ = os.WriteFile(flagFile, []byte(fmt.Sprintf("%d", time.Now().Unix())), 0o644)
	appendLog(logFile, "wrote sessionstart flag: "+flagFile)
	appendLog(logFile, fmt.Sprintf("phase=%s obligations=%s gates=%s", phase, obligCount, gateCount))
	return ExitSuccess
}

// runPreToolUseHook implements the Claude Code PreToolUse hook.
// Blocks non-Read tool calls until "Bootstrap:" is found in an assistant message.
// Uses cache file + SessionStart TTL flag to avoid redundant transcript scans.
func runPreToolUseHook(projectDir string, stdout io.Writer, stderr io.Writer) int {
	const logFile = "/tmp/ai-skill-bootstrap-hook.log"
	ts := time.Now().Format("2006-01-02T15:04:05")
	appendLog(logFile, "=== "+ts+" PreToolUse hook fired (Go) ===")

	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "ALLOW_BAD_INPUT:", err)
		return ExitSuccess
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw, &payload); err != nil {
		_, _ = fmt.Fprintln(stderr, "ALLOW_BAD_INPUT:", err)
		return ExitSuccess
	}

	var toolName, transcriptPath string
	if v, ok := payload["tool_name"]; ok {
		_ = json.Unmarshal(v, &toolName)
	}
	if v, ok := payload["transcript_path"]; ok {
		_ = json.Unmarshal(v, &transcriptPath)
	}
	appendLog(logFile, fmt.Sprintf("DIAG tool=%q transcript=%q", toolName, transcriptPath))
	_, _ = fmt.Fprintf(stderr, "DIAG tool=%q transcript=%q\n", toolName, transcriptPath)

	if toolName == "Read" {
		_, _ = fmt.Fprintln(stderr, "ALLOW_READ_TOOL:", toolName)
		return ExitSuccess
	}

	if transcriptPath == "" || !claudeFileExists(transcriptPath) {
		_, _ = fmt.Fprintln(stderr, "ALLOW_NO_TRANSCRIPT:", transcriptPath)
		return ExitSuccess
	}

	cacheKey := md5Short(transcriptPath)
	cacheFile := "/tmp/ai-skill-bootstrap-" + cacheKey + ".done"
	if claudeFileExists(cacheFile) {
		_, _ = fmt.Fprintln(stderr, "ALLOW_CACHED")
		return ExitSuccess
	}

	if projectDir != "" {
		projectHash := md5Short(projectDir)
		flagFile := "/tmp/ai-skill-sessionstart-" + projectHash + ".flag"
		if data, err := os.ReadFile(flagFile); err == nil {
			var flagTs int64
			if _, err2 := fmt.Sscanf(strings.TrimSpace(string(data)), "%d", &flagTs); err2 == nil {
				if time.Now().Unix()-flagTs < 120 {
					_ = os.WriteFile(cacheFile, []byte{}, 0o644)
					_, _ = fmt.Fprintln(stderr, "ALLOW_SESSIONSTART_FLAG")
					appendLog(logFile, "exit_code: 0 (sessionstart flag)")
					return ExitSuccess
				}
			}
		}
	}

	for _, text := range extractAssistantTexts(transcriptPath) {
		if strings.Contains(text, "Bootstrap:") {
			_ = os.WriteFile(cacheFile, []byte{}, 0o644)
			_, _ = fmt.Fprintln(stderr, "ALLOW_RECEIPT_FOUND")
			appendLog(logFile, "exit_code: 0 (receipt found)")
			return ExitSuccess
		}
	}

	_, _ = fmt.Fprintln(stderr, "BLOCK_NO_RECEIPT")
	appendLog(logFile, "exit_code: 2 (block no receipt)")
	_, _ = fmt.Fprint(stderr, `[ai-skill PreToolUse hook] Bootstrap Receipt missing.

Before calling any tool other than Read, you MUST:
1. Read CORE_BOOTSTRAP.md
2. Query runtime/runtime.db for phase / obligations / gates
3. Read the 3 required_reads: enforcement/rule-weight.md, enforcement/dependency-reading.md, enforcement/conversation-goal-ledger.md
4. Output the Bootstrap Receipt in your first user-facing response:
   Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>
   Active per-turn obligations: <obligation ids>

Only Read tool calls are allowed before the Receipt is emitted.
`)
	return ExitValidationFailed
}

// runPostToolUseHook implements the Claude Code PostToolUse hook.
// Injects a Bootstrap Receipt reminder when the receipt is not yet present.
// Always exits 0 (PostToolUse cannot reliably block without breaking tool results).
func runPostToolUseHook(projectDir string, stdout io.Writer, stderr io.Writer) int {
	raw, _ := io.ReadAll(os.Stdin)
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ExitSuccess
	}

	var transcriptPath string
	if v, ok := payload["transcript_path"]; ok {
		_ = json.Unmarshal(v, &transcriptPath)
	}

	if transcriptPath != "" {
		cacheKey := md5Short(transcriptPath)
		cacheFile := "/tmp/ai-skill-bootstrap-" + cacheKey + ".done"
		if claudeFileExists(cacheFile) {
			_, _ = fmt.Fprintln(stderr, "CACHED_DONE")
			return ExitSuccess
		}
		if claudeFileExists(transcriptPath) {
			for _, text := range extractAssistantTexts(transcriptPath) {
				if strings.Contains(text, "Bootstrap:") {
					_ = os.WriteFile(cacheFile, []byte{}, 0o644)
					_, _ = fmt.Fprintln(stderr, "RECEIPT_FOUND")
					return ExitSuccess
				}
			}
		}
	}

	reminder := "[ai-skill PostToolUse] Bootstrap Receipt not yet emitted. " +
		"Before writing your next response, you MUST:\n" +
		"1. Read CORE_BOOTSTRAP.md\n" +
		"2. Query runtime/runtime.db (phase / obligations / gates)\n" +
		"3. Read enforcement/rule-weight.md, enforcement/dependency-reading.md, enforcement/conversation-goal-ledger.md\n" +
		"4. Output Bootstrap Receipt as the first line of your response:\n" +
		"   Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>\n" +
		"   Active per-turn obligations: <obligation ids>"
	output := map[string]interface{}{
		"hookSpecificOutput": map[string]interface{}{
			"hookEventName":     "PostToolUse",
			"additionalContext": reminder,
		},
	}
	if err := json.NewEncoder(stdout).Encode(output); err == nil {
		_, _ = fmt.Fprintln(stderr, "INJECTED_REMINDER")
	}
	return ExitSuccess
}

// runUserPromptSubmitHook implements the Claude Code UserPromptSubmit hook.
// Injects per-turn obligation reminder + CORE_BOOTSTRAP.md as additionalContext.
func runUserPromptSubmitHook(projectDir string, stdout io.Writer, stderr io.Writer) int {
	const logFile = "/tmp/ai-skill-prompt-hook.log"
	appendLog(logFile, time.Now().Format("2006-01-02T15:04:05")+" UserPromptSubmit fired (Go)")

	aiSkillRepo := resolveClaudeAiSkillRepo(projectDir)
	bootstrap := readFileSafe(filepath.Join(aiSkillRepo, "CORE_BOOTSTRAP.md"))
	gitReport := formatDirtyGitRepoReport(projectDir)
	combined := "[ai-skill per-turn obligation] Final response MUST end with ### Cognitive Mode 報告 block. " +
		"First-turn ALSO outputs Bootstrap Receipt. Canonical spec: runtime/core-bootstrap.yaml.\n\n---\n" +
		bootstrap
	if gitReport != "" {
		combined += "\n\n---\n" + gitReport
	}

	output := map[string]interface{}{
		"hookSpecificOutput": map[string]interface{}{
			"hookEventName":     "UserPromptSubmit",
			"additionalContext": combined,
		},
	}
	if err := json.NewEncoder(stdout).Encode(output); err != nil {
		_, _ = fmt.Fprintln(stderr, "encode error:", err)
	}
	return ExitSuccess
}

// runStopHook implements the Claude Code Stop hook.
// Blocks stop if the last assistant message lacks a Cognitive Mode block
// (either compact "Cognitive: X..." or full "### Cognitive Mode 報告").
func runStopHook(projectDir string, stdout io.Writer, stderr io.Writer) int {
	const logFile = "/tmp/ai-skill-stop-hook.log"
	ts := time.Now().Format("2006-01-02T15:04:05")
	appendLog(logFile, "=== "+ts+" Stop hook fired (Go) ===")

	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "ALLOW_BAD_INPUT:", err)
		return ExitSuccess
	}
	appendLog(logFile, "input_json: "+string(raw))

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw, &payload); err != nil {
		_, _ = fmt.Fprintln(stderr, "ALLOW_BAD_INPUT:", err)
		return ExitSuccess
	}

	if v, ok := payload["stop_hook_active"]; ok {
		var active bool
		if err := json.Unmarshal(v, &active); err == nil && active {
			_, _ = fmt.Fprintln(stderr, "ALLOW_LOOP_GUARD")
			return ExitSuccess
		}
	}

	var transcriptPath string
	if v, ok := payload["transcript_path"]; ok {
		_ = json.Unmarshal(v, &transcriptPath)
	}
	if transcriptPath == "" || !claudeFileExists(transcriptPath) {
		_, _ = fmt.Fprintln(stderr, "ALLOW_NO_TRANSCRIPT: path="+transcriptPath)
		return ExitSuccess
	}

	texts := extractAssistantTexts(transcriptPath)
	if len(texts) == 0 {
		_, _ = fmt.Fprintln(stderr, "ALLOW_NO_ASSISTANT_MSG")
		return ExitSuccess
	}
	lastText := texts[len(texts)-1]

	tail := lastText
	if len(tail) > 200 {
		tail = tail[len(tail)-200:]
	}
	diagMsg := fmt.Sprintf("DIAG last_msg_len=%d tail=%q", len(lastText), tail)
	appendLog(logFile, diagMsg)
	_, _ = fmt.Fprintln(stderr, diagMsg)

	cogRe := regexp.MustCompile(`(### Cognitive Mode 報告|(?:^|\n)Cognitive: [A-Z])`)
	if cogRe.MatchString(lastText) {
		if formatDirtyGitRepoReport(projectDir) != "" {
			gitReportRe := regexp.MustCompile(`(?m)^### (Project Git Report|Git Repo Report|Git Repository Report)\b`)
			if !gitReportRe.MatchString(lastText) {
				_, _ = fmt.Fprintln(stderr, "BLOCK_MISSING_PROJECT_GIT_REPORT")
				appendLog(logFile, "exit_code: 2 (block missing project git report)")
				_, _ = fmt.Fprint(stderr, "[ai-skill Stop hook] Dirty Git repositories were detected under the project root, but your final response did not include `### Project Git Report`.\n\n"+
					"If one nested Git repo changed, report that repo. If multiple nested Git repos changed, merge them into one `### Project Git Report` section with one bullet per repo and clearly distinguish current-task changes from pre-existing/unrelated dirty state.\n")
				return ExitValidationFailed
			}
		}
		_, _ = fmt.Fprintln(stderr, "ALLOW_BLOCK_PRESENT")
		appendLog(logFile, "exit_code: 0")
		return ExitSuccess
	}

	_, _ = fmt.Fprintln(stderr, "BLOCK_MISSING")
	appendLog(logFile, "exit_code: 2")
	_, _ = fmt.Fprint(stderr, "[ai-skill Stop hook] Missing obligation: your final response did not include the `### Cognitive Mode 報告` block.\n\n"+
		"Per runtime/core-bootstrap.yaml §per_turn_obligations[obligation.cognitive.mode_report], every final user-facing "+
		"response MUST end with a Cognitive Mode block (compact 1-line for trivial all-default tasks: "+
		"`Cognitive: <e>·<c>·<g>·<m> / V:<v> / Cost:<cost> / Sig:<signal>`; full 6-row markdown table otherwise).\n\n"+
		"Please append the block to your response now, then stop again. Canonical format spec: runtime/core-bootstrap.yaml. "+
		"Query active obligations: `ai-skill runtime obligations`.\n")
	return ExitValidationFailed
}

// cognitiveV2Defaults are the 6 default dim values for the v2 compact form.
// Compact form is only valid when all 6 dims match these defaults exactly.
var cognitiveV2Defaults = map[string]string{
	"execution_mode":  "NORMAL",
	"context_mode":    "SUMMARY_FIRST",
	"governance_mode": "STANDARD",
	"memory_mode":     "NONE",
	"validation_mode": "CHECKLIST",
	"cognitive_cost":  "LOW",
}

// parseCompactCognitiveLine parses a v2 compact Cognitive Contract line:
//
//	"Cognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:<signal>"
//
// Returns a map with the 6 dim values, or nil if the line is not a valid compact form.
func parseCompactCognitiveLine(line string) map[string]string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "Cognitive: ") {
		return nil
	}
	rest := strings.TrimPrefix(line, "Cognitive: ")
	segments := strings.Split(rest, " / ")
	if len(segments) < 3 {
		return nil
	}
	dims := strings.Split(segments[0], "·")
	if len(dims) < 4 {
		return nil
	}
	modes := map[string]string{
		"execution_mode":  dims[0],
		"context_mode":    dims[1],
		"governance_mode": dims[2],
		"memory_mode":     dims[3],
	}
	for _, seg := range segments[1:] {
		switch {
		case strings.HasPrefix(seg, "V:"):
			modes["validation_mode"] = strings.TrimPrefix(seg, "V:")
		case strings.HasPrefix(seg, "Cost:"):
			modes["cognitive_cost"] = strings.TrimPrefix(seg, "Cost:")
		case strings.HasPrefix(seg, "Sig:"):
			modes["activation_signal"] = strings.TrimPrefix(seg, "Sig:")
		}
	}
	return modes
}

// validateCognitiveContractFormat is the v2 entry-point name used in
// per_commit_obligations (core-bootstrap.yaml). The enforcement logic lives in
// runCommitMsgHook which accepts both compact and full form.
const validateCognitiveContractFormat = "obligation.commit.cognitive_mode_block"

// runCommitMsgHook enforces Phase 4 behavioral wiring of
// gate.execution.cognitive_mode_resolved. Commit message body must contain
// either a v2 compact Cognitive line or a '### Cognitive Mode 報告' full table
// (template defined in models/cognitive-modes/README.md v2). Mechanical commits
// may opt out via '[skip-cognitive-mode]' in the body. Merge commits auto-skip.
func runCommitMsgHook(result Result, root string, positional []string) Result {
	if len(positional) == 0 {
		// Hook called without message file path; cannot enforce, fail open with warning.
		result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "warning", Message: "no commit message file passed; check skipped"})
		return result
	}
	msgPath := positional[0]
	if !filepath.IsAbs(msgPath) {
		msgPath = filepath.Join(root, msgPath)
	}
	body, err := os.ReadFile(msgPath)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitGeneralFailure
		result.Error = &CommandError{Code: "commit_msg_read_failed", Message: err.Error()}
		return result
	}
	text := string(body)

	// Auto-skip merge commits (git auto-generated header)
	if strings.HasPrefix(strings.TrimLeft(text, " \t\n"), "Merge ") {
		result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "skipped", Message: "merge commit auto-skipped"})
		return result
	}

	// v2 compact form path: "Cognitive: <e>·<c>·<g>·<m> / V:<v> / Cost:<cost> / Sig:<sig>"
	// Valid only when all 6 dims are at their default values. Non-default dims require full form.
	for _, line := range strings.Split(text, "\n") {
		if compactModes := parseCompactCognitiveLine(line); compactModes != nil {
			var nonDefault []string
			for dim, val := range compactModes {
				if dim == "cognitive_cost" {
					continue // derived — not validated here; Phase 3 adds cost class check
				}
				if def, ok := cognitiveV2Defaults[dim]; ok && val != def {
					nonDefault = append(nonDefault, dim+"="+val)
				}
			}
			if len(nonDefault) > 0 {
				sort.Strings(nonDefault)
				result.Status = "blocked"
				result.ExitCode = ExitValidationFailed
				result.Error = &CommandError{
					Code:        "cognitive_compact_non_default",
					Message:     "Compact form used but non-default dims detected: " + strings.Join(nonDefault, ", ") + ". Compact form is only valid when all 6 dims are at default values.",
					Remediation: "Replace the compact line with a full ### Cognitive Mode 報告 table (6-dim v2 format) per models/cognitive-modes/README.md.",
				}
				return result
			}
			ctx := commitMsgCtx{text: text, staged: nil, root: root, modes: compactModes}
			var violations []string
			if v := validateCognitiveCost(ctx.modes); v != "" {
				violations = append(violations, v)
			}
			if v := validateActivationSignals(ctx); v != "" {
				violations = append(violations, v)
			}
			if len(violations) > 0 {
				result.Status = "blocked"
				result.ExitCode = ExitValidationFailed
				result.Error = &CommandError{
					Code:        "cognitive_compact_violations",
					Message:     "Declared compact Cognitive Contract conflicts with v2 validation:\n  - " + strings.Join(violations, "\n  - "),
					Remediation: "Use a known discovery signal and the derived cognitive_cost for the declared execution/context tuple.",
				}
				return result
			}
			if v := validateCapabilitySnippet(ctx.modes, ctx.text); v != "" {
				result.Status = "blocked"
				result.ExitCode = ExitValidationFailed
				result.Error = &CommandError{Code: "cognitive_compact_capability_violation", Message: v}
				return result
			}
			result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "ok", Message: "Cognitive Contract v2 compact form present (all dims at default)"})
			return result
		}
	}

	// Primary path: full form '### Cognitive Mode 報告' block → run Phase 3 behavioral validators.
	// Checked BEFORE opt-out marker to avoid false positives when commit body
	// documents/quotes the opt-out token (e.g. "Opt-out via '[skip-cognitive-mode]'").
	if strings.Contains(text, "### Cognitive Mode 報告") {
		modes := parseCognitiveModeBlock(text)
		staged, _ := gitLines(root, "diff", "--cached", "--name-only")

		// Phase 6 dispatcher: read per_commit_obligations order from
		// generated_surfaces[runtime.core_bootstrap.contract] and dispatch
		// validators by id via registry. Fallback to hardcoded order if
		// runtime.db is unavailable or not yet projected (e.g. fresh clone
		// before first `runtime compile`).
		ctx := commitMsgCtx{text: text, staged: staged, root: root, modes: modes}
		v2Violations := []string{}
		if v := validateCognitiveCost(ctx.modes); v != "" {
			v2Violations = append(v2Violations, v)
		}
		if v := validateActivationSignals(ctx); v != "" {
			v2Violations = append(v2Violations, v)
		}
		if v := validateCapabilitySnippet(ctx.modes, ctx.text); v != "" {
			v2Violations = append(v2Violations, v)
		}
		if len(v2Violations) > 0 {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{
				Code:        "cognitive_contract_v2_violations",
				Message:     "Declared Cognitive Contract v2 block conflicts with validation:\n  - " + strings.Join(v2Violations, "\n  - "),
				Remediation: "Use known activation_reason signals, derived cognitive_cost, and Capability summary for high-risk modes.",
			}
			return result
		}
		order := readPerCommitObligationsOrder(root)
		if len(order) == 0 {
			order = defaultCommitMsgDispatchOrder
		}
		var violations []string
		for _, id := range order {
			validator, ok := commitMsgValidatorRegistry[id]
			if !ok {
				// Obligation declared in YAML but no Go validator registered.
				// Skip silently (allows YAML to declare future-planned
				// obligations without breaking hook).
				continue
			}
			if v := validator(ctx); v != "" {
				violations = append(violations, v)
			}
		}
		if len(violations) > 0 {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{
				Code:        "cognitive_mode_violations",
				Message:     "Declared Cognitive Mode block conflicts with commit content:\n  - " + strings.Join(violations, "\n  - "),
				Remediation: "Update the Cognitive Mode block to match what the commit actually does, or split the commit. See runtime/cognitive-modes-*.yaml for contract details.",
			}
			return result
		}
		result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "ok", Message: "Cognitive Mode 報告 present + Phase 3 validators passed"})
		return result
	}

	// Fallback path: opt-out marker on its own line (require leading whitespace or BOL
	// to reduce false positives from prose mentions). Mechanical commits should
	// place the marker as a standalone trailer line.
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "[skip-cognitive-mode]" {
			result.Checks = append(result.Checks, Check{Name: "cognitive_mode_block", Status: "skipped", Message: "[skip-cognitive-mode] opt-out marker present on its own line"})
			return result
		}
	}

	result.Status = "blocked"
	result.ExitCode = ExitValidationFailed
	result.Error = &CommandError{
		Code:        "cognitive_mode_block_missing",
		Message:     "Commit message body must include a Cognitive Contract v2 block: compact single-line form (all-default dims) or full '### Cognitive Mode 報告' table (6-dim, non-default or high-risk).",
		Remediation: "Add compact form 'Cognitive: NORMAL·SUMMARY_FIRST·STANDARD·NONE / V:CHECKLIST / Cost:LOW / Sig:<signal>' for trivial commits, or full table per models/cognitive-modes/README.md v2 template. For mechanical commits, add a standalone '[skip-cognitive-mode]' trailer line.",
	}
	return result
}

// parseCognitiveModeBlock extracts the 4-dim mode resolution from a Cognitive
// Mode 報告 markdown table. Returns a map with keys execution_mode, context_mode,
// governance_mode, memory_mode (uppercase values). Missing/malformed rows
// produce empty strings; callers must handle empty values defensively.
func parseCognitiveModeBlock(text string) map[string]string {
	modes := map[string]string{}
	lines := strings.Split(text, "\n")
	inBlock := false
	for _, line := range lines {
		if strings.Contains(line, "### Cognitive Mode 報告") {
			inBlock = true
			continue
		}
		if !inBlock {
			continue
		}
		// Stop at next header or blank line after we've started parsing rows
		if strings.HasPrefix(strings.TrimSpace(line), "## ") {
			break
		}
		// Markdown table row: | dim_name | VALUE | reason |
		if !strings.HasPrefix(strings.TrimSpace(line), "|") {
			continue
		}
		parts := strings.Split(line, "|")
		// parts: "" | dim | value | reason | ""  → len >= 4
		if len(parts) < 4 {
			continue
		}
		dim := strings.TrimSpace(parts[1])
		val := strings.TrimSpace(parts[2])
		// Skip header and separator rows
		if dim == "" || dim == "維度" || strings.HasPrefix(dim, "---") {
			continue
		}
		if val == "" || strings.HasPrefix(val, "---") {
			continue
		}
		switch dim {
		case "execution_mode", "context_mode", "governance_mode", "memory_mode",
			"validation_mode", "cognitive_cost":
			modes[dim] = val
		}
	}
	return modes
}

// validateExecutionModeFloors implements Phase 3.1-B: enforce execution_mode
// floor requirements per runtime/cognitive-modes-phase-integration.yaml.
// Returns empty string when valid, otherwise a violation message.
func validateExecutionModeFloors(modes map[string]string, staged []string) string {
	exec := modes["execution_mode"]
	gov := modes["governance_mode"]
	ctx := modes["context_mode"]
	mem := modes["memory_mode"]

	// FAST cannot touch governance/, enforcement/, or runtime/ (auto-escalation rule)
	if exec == "FAST" {
		for _, f := range staged {
			if strings.HasPrefix(f, "governance/") || strings.HasPrefix(f, "enforcement/") || strings.HasPrefix(f, "runtime/") {
				return "execution_mode=FAST forbidden when staged files touch governance/, enforcement/, or runtime/ (auto-escalation rule per cognitive-modes-phase-integration.yaml). File: " + f
			}
		}
	}

	// DEEP / FORENSIC / RECOVERY require governance_mode ≥ STRICT
	if exec == "DEEP" || exec == "FORENSIC" || exec == "RECOVERY" {
		if gov != "STRICT" && gov != "LOCKDOWN" {
			return "execution_mode=" + exec + " requires governance_mode ≥ STRICT (declared: " + gov + ")"
		}
	}

	// DEEP requires context_mode ≥ SOURCE_BACKED
	if exec == "DEEP" && ctx != "SOURCE_BACKED" && ctx != "GRAPH_ASSISTED" {
		return "execution_mode=DEEP requires context_mode ≥ SOURCE_BACKED (declared: " + ctx + ")"
	}
	// FORENSIC requires context_mode = GRAPH_ASSISTED
	if exec == "FORENSIC" && ctx != "GRAPH_ASSISTED" {
		return "execution_mode=FORENSIC requires context_mode=GRAPH_ASSISTED (declared: " + ctx + ")"
	}
	// RECOVERY requires context_mode = CHECKLIST_FIRST and memory_mode = FAILURE_REPLAY
	if exec == "RECOVERY" {
		if ctx != "CHECKLIST_FIRST" {
			return "execution_mode=RECOVERY requires context_mode=CHECKLIST_FIRST (declared: " + ctx + ")"
		}
		if mem != "FAILURE_REPLAY" {
			return "execution_mode=RECOVERY requires memory_mode=FAILURE_REPLAY (declared: " + mem + ")"
		}
	}

	return ""
}

// validateGovernanceModeConsistency implements Phase 3.3-B: enforce that the
// declared governance_mode matches the sensitivity of staged files, and that
// LOCKDOWN commits carry an [approved-by: ...] trailer.
func validateGovernanceModeConsistency(modes map[string]string, staged []string, text string) string {
	gov := modes["governance_mode"]
	if gov == "" {
		return "governance_mode missing from Cognitive Mode block"
	}

	// LIGHT: forbidden when staged files touch governance-critical paths
	if gov == "LIGHT" {
		for _, f := range staged {
			if strings.HasPrefix(f, "governance/") || strings.HasPrefix(f, "enforcement/") || strings.HasPrefix(f, "runtime/") || strings.HasPrefix(f, "validation/") {
				return "governance_mode=LIGHT forbidden when staged files include governance-critical paths (governance/, enforcement/, runtime/, validation/). File: " + f
			}
		}
	}

	// LOCKDOWN: require explicit [approved-by: <name>] trailer line
	if gov == "LOCKDOWN" {
		hasApproval := false
		for _, line := range strings.Split(text, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "[approved-by:") && strings.HasSuffix(trimmed, "]") {
				hasApproval = true
				break
			}
		}
		if !hasApproval {
			return "governance_mode=LOCKDOWN requires an [approved-by: <name>] trailer line in the commit body"
		}
	}

	return ""
}

// validateMemoryModeSubdir implements Phase 3.4-B: enforce that staged memory/
// files are inside the subdir corresponding to the declared memory_mode.
// NONE → no memory/ files; EPISODIC → memory/episodic/; DECISION_REPLAY →
// memory/decision/; FAILURE_REPLAY → memory/failure/; PROJECT_CONTEXT →
// memory/project/.
func validateMemoryModeSubdir(modes map[string]string, staged []string) string {
	mem := modes["memory_mode"]
	allowedPrefix := ""
	switch mem {
	case "NONE":
		allowedPrefix = "" // any memory/ touch is a violation
	case "EPISODIC":
		allowedPrefix = "memory/episodic/"
	case "DECISION_REPLAY":
		allowedPrefix = "memory/decision/"
	case "FAILURE_REPLAY":
		allowedPrefix = "memory/failure/"
	case "PROJECT_CONTEXT":
		allowedPrefix = "memory/project/"
	default:
		// Unknown mode value — let block-presence validator handle absence; tolerate unrecognized strings here
		return ""
	}

	// memory/README.md and memory/retrieval-governance/ are layer-level docs and
	// are not subject to the per-mode subdir restriction (they describe the layer
	// itself, not memory content).
	isLayerDoc := func(f string) bool {
		return f == "memory/README.md" ||
			strings.HasPrefix(f, "memory/retrieval-governance/")
	}

	for _, f := range staged {
		if !strings.HasPrefix(f, "memory/") {
			continue
		}
		if isLayerDoc(f) {
			continue
		}
		if mem == "NONE" {
			return "memory_mode=NONE but commit touches " + f + " (NONE forbids all memory/ writes per cognitive-modes-memory-integration.yaml)"
		}
		if !strings.HasPrefix(f, allowedPrefix) {
			return "memory_mode=" + mem + " allows only " + allowedPrefix + " but commit touches " + f
		}
	}
	return ""
}

// validatePlanStatusSync implements runtime/plan-status-sync-enforcement.yaml:
// when a commit body claims phase/milestone completion AND references an active
// plan file by path, that plan file MUST be in the staged set.
//
// Trigger composition (all three required to fire):
//  1. ≥1 completion vocabulary word
//  2. ≥1 "Phase <num>" / "phase <num>" reference
//  3. ≥1 plans/active/<f>.md path reference
//
// Opt-out: standalone "[skip-plan-status-sync]" trailer line.
var (
	planPathRE        = regexp.MustCompile(`plans/active/[^\s)"\]]+\.md`)
	phaseMentionRE    = regexp.MustCompile(`(?i)Phase\s+\d+(?:\.\d+)?(?:[\.-][A-Za-z]+)?`)
	completionPhrases = []string{
		"complete", "completed", "completes", "done", "finish", "finished",
		"完成", "結案", "結束", "✅",
	}
)

func validatePlanStatusSync(text string, staged []string) string {
	// Opt-out marker on its own line
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-plan-status-sync]" {
			return ""
		}
	}

	// Trigger 1: completion vocabulary
	hasCompletion := false
	lowered := strings.ToLower(text)
	for _, phrase := range completionPhrases {
		if strings.Contains(lowered, strings.ToLower(phrase)) {
			hasCompletion = true
			break
		}
	}
	if !hasCompletion {
		return ""
	}

	// Trigger 2: Phase N mention
	if !phaseMentionRE.MatchString(text) {
		return ""
	}

	// Trigger 3: plans/active/*.md reference
	planRefs := planPathRE.FindAllString(text, -1)
	if len(planRefs) == 0 {
		return ""
	}

	// Trigger fired. Each referenced plan must be in staged set.
	stagedSet := make(map[string]bool, len(staged))
	for _, s := range staged {
		stagedSet[s] = true
	}
	var missing []string
	seen := map[string]bool{}
	for _, ref := range planRefs {
		// Normalize: strip trailing markdown link garbage
		clean := strings.TrimRight(ref, "),]\"")
		if seen[clean] {
			continue
		}
		seen[clean] = true
		if !stagedSet[clean] {
			missing = append(missing, clean)
		}
	}
	if len(missing) == 0 {
		return ""
	}
	return "plan-status-sync: commit body claims phase completion and references " +
		strings.Join(missing, ", ") +
		" but the plan file is not in the staged set. Update the plan's Phase section in the same commit (runtime/plan-status-sync-enforcement.yaml). Use a [skip-plan-status-sync] trailer for retrospective references."
}

// validateTokenBudget implements runtime/cognitive-modes-token-budget.yaml:
// when a commit body declares a Token Estimate trailer AND the declared
// cognitive mode combination has a known budget, the estimate must not
// exceed max_tokens.
//
// Trigger: body contains "Token Estimate: <n>" (case-insensitive) AND
// Cognitive Mode block declared execution_mode.
//
// Budget table is hardcoded here to keep the hook self-contained;
// canonical source is runtime/cognitive-modes-token-budget.yaml. If the
// YAML budgets change, this function must be updated in sync.
//
// Opt-out: standalone "[skip-token-budget]" trailer line.
var tokenEstimateRE = regexp.MustCompile(`(?i)Token\s+Estimate:\s*(\d+)`)

func validateTokenBudget(modes map[string]string, text string) string {
	// Opt-out marker on its own line
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-token-budget]" {
			return ""
		}
	}

	match := tokenEstimateRE.FindStringSubmatch(text)
	if len(match) < 2 {
		// No declared estimate → no-op (validator is opt-in for this turn)
		return ""
	}
	estimate := 0
	for _, c := range match[1] {
		estimate = estimate*10 + int(c-'0')
	}

	exec := modes["execution_mode"]
	ctx := modes["context_mode"]
	gov := modes["governance_mode"]
	mem := modes["memory_mode"]

	// Exact tuple budgets per runtime/cognitive-modes-token-budget.yaml §budgets
	exactBudgets := map[string]int{
		"FAST|INDEX_ONLY|LIGHT|NONE":                    1000,
		"NORMAL|SUMMARY_FIRST|STANDARD|EPISODIC":        5000,
		"DEEP|SOURCE_BACKED|STRICT|DECISION_REPLAY":     20000,
		"FORENSIC|GRAPH_ASSISTED|STRICT|FAILURE_REPLAY": 50000,
	}
	// Default budget by execution_mode (when exact tuple not found)
	execDefaults := map[string]int{
		"FAST":     1000,
		"NORMAL":   5000,
		"DEEP":     20000,
		"FORENSIC": 50000,
		"RECOVERY": 50000,
	}

	key := exec + "|" + ctx + "|" + gov + "|" + mem
	budget, ok := exactBudgets[key]
	if !ok {
		budget, ok = execDefaults[exec]
		if !ok {
			// Unknown execution_mode → no enforcement
			return ""
		}
	}

	if estimate > budget {
		return "token_budget: declared Token Estimate=" + itoa(estimate) +
			" exceeds budget=" + itoa(budget) +
			" for mode tuple (execution_mode=" + exec +
			", context_mode=" + ctx +
			", governance_mode=" + gov +
			", memory_mode=" + mem +
			"). Downgrade context_mode (GRAPH_ASSISTED → SOURCE_BACKED → CHECKLIST_FIRST → SUMMARY_FIRST → INDEX_ONLY) or split the work. Use [skip-token-budget] only if exceptional."
	}
	return ""
}

// deriveCognitiveCost implements runtime/cognitive-modes-cost-class.yaml:
// cognitive_cost is derived from execution_mode × context_mode, not claimed by
// the agent.
func deriveCognitiveCost(executionMode, contextMode string) string {
	switch executionMode {
	case "FAST":
		if contextMode == "INDEX_ONLY" {
			return "LOW"
		}
		return "MEDIUM"
	case "NORMAL":
		if contextMode == "INDEX_ONLY" || contextMode == "SUMMARY_FIRST" {
			return "LOW"
		}
		if contextMode == "CHECKLIST_FIRST" || contextMode == "SOURCE_BACKED" || contextMode == "GRAPH_ASSISTED" {
			return "MEDIUM"
		}
	case "DEEP":
		return "HIGH"
	case "FORENSIC", "RECOVERY":
		return "VERY_HIGH"
	}
	return ""
}

func validateCognitiveCost(modes map[string]string) string {
	declared := modes["cognitive_cost"]
	if declared == "" {
		return "cognitive_cost missing from Cognitive Contract v2 block"
	}
	derived := deriveCognitiveCost(modes["execution_mode"], modes["context_mode"])
	if derived == "" {
		return "cognitive_cost: cannot derive cost for execution_mode=" + modes["execution_mode"] + " context_mode=" + modes["context_mode"]
	}
	if declared != derived {
		return "cognitive_cost mismatch: declared=" + declared + " derived=" + derived + " for execution_mode=" + modes["execution_mode"] + " context_mode=" + modes["context_mode"]
	}
	return ""
}

func parseActivationSignals(text string, modes map[string]string) []string {
	if sig := strings.TrimSpace(modes["activation_signal"]); sig != "" {
		return []string{sig}
	}
	lines := strings.Split(text, "\n")
	inBlock := false
	var signals []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "activation_reason:" {
			inBlock = true
			continue
		}
		if !inBlock {
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "##") || strings.HasPrefix(trimmed, "Capability summary:") {
			break
		}
		if strings.HasPrefix(trimmed, "- ") {
			sig := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if idx := strings.Index(sig, "#"); idx >= 0 {
				sig = strings.TrimSpace(sig[:idx])
			}
			if sig != "" {
				signals = append(signals, sig)
			}
		}
	}
	return signals
}

func readKnownCognitiveSignals(root string) map[string]bool {
	known := map[string]bool{}
	dbPath := filepath.Join(root, "runtime", "runtime.db")
	if db, err := sql.Open("sqlite", dbPath); err == nil {
		defer db.Close()
		var raw string
		if err := db.QueryRow("SELECT data FROM generated_surfaces WHERE target_key='runtime.cognitive_modes.discovery' LIMIT 1").Scan(&raw); err == nil {
			var doc map[string]any
			if json.Unmarshal([]byte(raw), &doc) == nil {
				if signals, ok := doc["signals"].([]any); ok {
					for _, item := range signals {
						if m, ok := item.(map[string]any); ok {
							if name, ok := m["name"].(string); ok && name != "" {
								known[name] = true
							}
						}
					}
				}
			}
		}
	}
	if len(known) > 0 {
		return known
	}
	// Fallback for fresh clones before runtime.db is compiled. The canonical
	// signal vocabulary remains runtime/cognitive-modes-discovery.yaml.
	body, err := os.ReadFile(filepath.Join(root, "runtime", "cognitive-modes-discovery.yaml"))
	if err != nil {
		return known
	}
	nameRE := regexp.MustCompile(`^\s*-\s+name:\s*([A-Za-z0-9_]+)\s*$`)
	for _, line := range strings.Split(string(body), "\n") {
		m := nameRE.FindStringSubmatch(line)
		if len(m) == 2 {
			known[m[1]] = true
		}
	}
	return known
}

func validateActivationSignals(ctx commitMsgCtx) string {
	signals := parseActivationSignals(ctx.text, ctx.modes)
	if len(signals) == 0 {
		return "activation_reason missing: Cognitive Contract v2 requires at least one discovery signal"
	}
	known := readKnownCognitiveSignals(ctx.root)
	if len(known) == 0 {
		return "activation_reason: known discovery signal list unavailable from runtime generated surface or YAML source"
	}
	var unknown []string
	for _, sig := range signals {
		if !known[sig] {
			unknown = append(unknown, sig)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return "activation_reason contains unknown discovery signal(s): " + strings.Join(unknown, ", ")
	}
	return ""
}

func validateCapabilitySnippet(modes map[string]string, text string) string {
	exec := modes["execution_mode"]
	gov := modes["governance_mode"]
	highRisk := exec == "DEEP" || exec == "FORENSIC" || exec == "RECOVERY" || gov == "STRICT" || gov == "LOCKDOWN"
	if !highRisk {
		return ""
	}
	if !strings.Contains(text, "Capability summary:") {
		return "capability snippet missing: high-risk Cognitive Contract modes require a Capability summary section"
	}
	return ""
}

// itoa avoids importing strconv solely for one call site.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	buf := [20]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// validateAdaptiveTriggers implements runtime/cognitive-modes-adaptive.yaml:
// 3 commit-msg-detectable adaptive triggers.
//
//  1. contradiction_risk — body contains contradiction-class keyword AND
//     references ≥2 distinct sources (plans/, constitution/, decisions/).
//     Requires governance_mode≥STRICT AND context_mode in {SOURCE_BACKED, GRAPH_ASSISTED}.
//
//  2. repeated_failure — body references ≥2 failure patterns OR uses
//     revert/hotfix/retry vocabulary ≥2 times. Requires
//     execution_mode=RECOVERY AND memory_mode=FAILURE_REPLAY.
//
//  3. budget_near_ceiling — Token Estimate ≥ 80% of mode tuple budget.
//     Warning level (not blocking).
//
// Opt-out: standalone "[skip-adaptive]" trailer line.
var (
	contradictionKeywordsRE = regexp.MustCompile(`(?i)contradict\w*|conflict\w*|mismatch\w*|discrepan\w*|衝突|矛盾|不一致`)
	failurePatternRefRE     = regexp.MustCompile(`enforcement/failure-patterns/[\w-]+\.md`)
	revertHotfixRE          = regexp.MustCompile(`(?i)\b(revert|hotfix|retry)\b`)
	sourceClassRE           = regexp.MustCompile(`(plans/|constitution/|decisions/)[^\s)"\]]+`)
)

func validateAdaptiveTriggers(modes map[string]string, text string) string {
	// Opt-out marker on its own line
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-adaptive]" {
			return ""
		}
	}

	exec := modes["execution_mode"]
	ctx := modes["context_mode"]
	gov := modes["governance_mode"]
	mem := modes["memory_mode"]

	var violations []string

	// Trigger 1: contradiction_risk
	if contradictionKeywordsRE.MatchString(text) {
		// Count distinct source references
		refs := sourceClassRE.FindAllString(text, -1)
		distinct := map[string]bool{}
		for _, r := range refs {
			distinct[r] = true
		}
		if len(distinct) >= 2 {
			govOK := gov == "STRICT" || gov == "LOCKDOWN"
			ctxOK := ctx == "SOURCE_BACKED" || ctx == "GRAPH_ASSISTED"
			if !govOK || !ctxOK {
				violations = append(violations,
					"adaptive: contradiction risk detected (cross-source keywords + ≥2 distinct source refs) but governance_mode="+
						gov+" / context_mode="+ctx+
						" below required floor (governance_mode≥STRICT, context_mode in {SOURCE_BACKED, GRAPH_ASSISTED}). Upgrade modes per runtime/cognitive-modes-adaptive.yaml §contradiction_risk.")
			}
		}
	}

	// Trigger 2: repeated_failure
	failureRefs := len(failurePatternRefRE.FindAllString(text, -1))
	revertHits := len(revertHotfixRE.FindAllString(text, -1))
	if failureRefs >= 2 || revertHits >= 2 {
		if exec != "RECOVERY" || mem != "FAILURE_REPLAY" {
			violations = append(violations,
				"adaptive: repeated failure signal (failure-pattern refs="+itoa(failureRefs)+
					", revert/hotfix/retry hits="+itoa(revertHits)+
					") requires execution_mode=RECOVERY and memory_mode=FAILURE_REPLAY (declared: execution_mode="+
					exec+", memory_mode="+mem+
					"). Adjust mode tuple per runtime/cognitive-modes-adaptive.yaml §repeated_failure.")
		}
	}

	// Trigger 3: budget_near_ceiling — warning only; emitted as a violation
	// line so it surfaces in commit-msg output, but agents should treat it
	// as advisory until adaptive infrastructure exists. Conservative: only
	// fire when an explicit Token Estimate is declared.
	match := tokenEstimateRE.FindStringSubmatch(text)
	if len(match) >= 2 {
		estimate := 0
		for _, c := range match[1] {
			estimate = estimate*10 + int(c-'0')
		}
		// Reuse the same budget table as validateTokenBudget
		exactBudgets := map[string]int{
			"FAST|INDEX_ONLY|LIGHT|NONE":                    1000,
			"NORMAL|SUMMARY_FIRST|STANDARD|EPISODIC":        5000,
			"DEEP|SOURCE_BACKED|STRICT|DECISION_REPLAY":     20000,
			"FORENSIC|GRAPH_ASSISTED|STRICT|FAILURE_REPLAY": 50000,
		}
		execDefaults := map[string]int{
			"FAST": 1000, "NORMAL": 5000, "DEEP": 20000, "FORENSIC": 50000, "RECOVERY": 50000,
		}
		key := exec + "|" + ctx + "|" + gov + "|" + mem
		budget, ok := exactBudgets[key]
		if !ok {
			budget, ok = execDefaults[exec]
		}
		if ok && estimate >= (budget*80/100) && estimate <= budget {
			violations = append(violations,
				"adaptive[warning]: Token Estimate="+itoa(estimate)+
					" is ≥80% of budget="+itoa(budget)+
					"; consider downgrading context_mode one step along the downgrade_path (GRAPH_ASSISTED → SOURCE_BACKED → CHECKLIST_FIRST → SUMMARY_FIRST → INDEX_ONLY) OR split the work. Suppress this notice with [skip-adaptive].")
		}
	}

	if len(violations) == 0 {
		return ""
	}
	return strings.Join(violations, "\n  - ")
}

// validateBootstrapEntryThinness implements runtime/bootstrap-entry-points.yaml:
// when a commit stages an AI-tool entry file (repo-root CLAUDE.md,
// .cursor/rules/ai-skill-bootstrap.mdc, .roomodes), enforce that the file
// remains a thin pointer to canonical sources.
//
// Checks:
//   - Line count ≤ 30
//   - No mode enum substrings (FAST/NORMAL/DEEP/FORENSIC/RECOVERY etc.)
//   - No Bootstrap Receipt format example with the literal phase pattern
//   - No full Cognitive Mode 報告 markdown table (the "| 維度 | 值 | 理由 |" header)
//
// Opt-out: standalone "[skip-bootstrap-thinness]" trailer line.
var bootstrapEntryPaths = []string{
	"CLAUDE.md",
	".cursor/rules/ai-skill-bootstrap.mdc",
	".roomodes",
	// AGENTS.md is generic agent entry (Codex, Cursor partial, Aider, Cline,
	// other AGENTS.md-aware tools); thinness applies equally.
	"AGENTS.md",
}

func validateBootstrapEntryThinness(text string, staged []string, root string) string {
	// Opt-out marker on its own line
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-bootstrap-thinness]" {
			return ""
		}
	}

	// Identify staged entry files
	stagedEntries := []string{}
	for _, s := range staged {
		for _, p := range bootstrapEntryPaths {
			if s == p {
				stagedEntries = append(stagedEntries, s)
				break
			}
		}
	}
	if len(stagedEntries) == 0 {
		return ""
	}

	// Forbidden substrings (each substring is enough to fail)
	forbiddenSubs := []string{
		"FAST/NORMAL/DEEP/FORENSIC/RECOVERY",
		"INDEX_ONLY/SUMMARY_FIRST/CHECKLIST_FIRST/SOURCE_BACKED/GRAPH_ASSISTED",
		"LIGHT/STANDARD/STRICT/LOCKDOWN",
		"NONE/EPISODIC/DECISION_REPLAY/FAILURE_REPLAY/PROJECT_CONTEXT",
		"Bootstrap: rules=✓ phase=phase.bootstrap obligations=",
		"| 維度 | 值 | 理由 |",
	}

	const maxLines = 30
	var violations []string
	for _, path := range stagedEntries {
		fullPath := path
		if !filepath.IsAbs(fullPath) {
			fullPath = filepath.Join(root, fullPath)
		}
		body, err := os.ReadFile(fullPath)
		if err != nil {
			// File might be staged for deletion; ignore
			continue
		}
		content := string(body)
		// Line count check (count lines including trailing-newline content)
		lineCount := 1
		for _, b := range content {
			if b == '\n' {
				lineCount++
			}
		}
		// Strip count if file ends with newline
		if strings.HasSuffix(content, "\n") {
			lineCount--
		}
		if lineCount > maxLines {
			violations = append(violations,
				"bootstrap-entry-thinness: "+path+" is "+itoa(lineCount)+
					" lines (max "+itoa(maxLines)+"); move obligation content to CORE_BOOTSTRAP.md or ai-tools/agent/<tool>.md per runtime/bootstrap-entry-points.yaml.")
			continue
		}
		// Forbidden substring check
		for _, sub := range forbiddenSubs {
			if strings.Contains(content, sub) {
				violations = append(violations,
					"bootstrap-entry-thinness: "+path+" contains canonical content fragment '"+
						sub+"'; this belongs in CORE_BOOTSTRAP.md / models/cognitive-modes/, not in tool entries.")
				break // one violation per file is enough
			}
		}
	}
	if len(violations) == 0 {
		return ""
	}
	return strings.Join(violations, "\n  - ")
}

// validateCLIDocSync enforces runtime/cli-modification-policy.yaml
// gate.cli.command_contract_synced: when staged Go files under
// scripts/ai-skill-cli/internal/app/ contain newly added subcommand
// dispatch (`case "run X":` or `case "X":` for runtime subcommands) OR
// new public `runXxxHook` / `buildRuntimeXxxResult` function, the
// command-contract.md must also be staged.
//
// Heuristic: scan staged Go files for `case "run ` / new func patterns
// via `git diff --cached`. Conservative — false negatives preferred
// over false positives.
//
// Opt-out: standalone `[skip-cli-doc-sync]` trailer line.
func validateCLIDocSync(text string, staged []string, root string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-cli-doc-sync]" {
			return ""
		}
	}
	cliSourceStaged := false
	docStaged := false
	for _, s := range staged {
		if strings.HasPrefix(s, "scripts/ai-skill-cli/internal/app/") && strings.HasSuffix(s, ".go") {
			cliSourceStaged = true
		}
		if s == "scripts/ai-skill-cli/docs/command-contract.md" {
			docStaged = true
		}
	}
	if !cliSourceStaged || docStaged {
		return ""
	}
	// CLI Go file staged but command-contract.md not staged. Check
	// git diff for newly added subcommand dispatch or hook handler
	// patterns. If none, skip (pure refactor).
	cmd := exec.Command("git", "-C", root, "diff", "--cached", "--", "scripts/ai-skill-cli/internal/app/")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	diff := string(out)
	patterns := []string{
		`+\tcase "run `,
		`+\tcase "obligations"`,
		`+func runCommitMsgHook`,
		`+func runPrePushHook`,
		`+func runPreCommitHook`,
		`+func buildRuntimeObligationsResult`,
	}
	for _, p := range patterns {
		if strings.Contains(diff, p) {
			return "cli-doc-sync: CLI source change adds subcommand dispatch / hook handler but scripts/ai-skill-cli/docs/command-contract.md is not staged. Per runtime/cli-modification-policy.yaml gate.cli.command_contract_synced. Use [skip-cli-doc-sync] for non-contract-affecting refactors."
		}
	}
	// Also flag generic new `case "` addition + `runXxxHook` function names.
	if regexp.MustCompile(`(?m)^\+func run[A-Z][A-Za-z]+Hook\b`).MatchString(diff) {
		return "cli-doc-sync: CLI source change adds new runXxxHook function but command-contract.md not staged. See runtime/cli-modification-policy.yaml."
	}
	return ""
}

// validateGlossaryRetroOwn enforces runtime/cli-modification-policy.yaml
// gate.glossary.retro_own_required: when staged diff touches framework
// cognitive vocabulary surface (runtime/cognitive-modes*.yaml,
// runtime/economics/**, ecosystem/**), knowledge/glossary/ai-skill.md
// MUST also be staged. Forces agents to retro-own newly introduced
// framework terms in the canonical glossary instead of letting them
// drift as subsystem-local vocabulary.
//
// Opt-out: standalone `[skip-glossary-retro-own]` trailer line. The
// commit message body should explain why this particular framework
// surface change does not introduce new vocabulary (e.g., typo fix,
// refactor of existing field, comment-only change).
//
// Source: plans/active/2026-05-25-1000-context-language-glossary-system.md
// Phase 6.
func validateGlossaryRetroOwn(text string, staged []string, root string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-glossary-retro-own]" {
			return ""
		}
	}
	frameworkSurfaceStaged := false
	glossaryStaged := false
	for _, s := range staged {
		if strings.HasPrefix(s, "runtime/cognitive-modes") && strings.HasSuffix(s, ".yaml") {
			frameworkSurfaceStaged = true
		}
		if strings.HasPrefix(s, "runtime/economics/") {
			frameworkSurfaceStaged = true
		}
		if strings.HasPrefix(s, "ecosystem/") {
			frameworkSurfaceStaged = true
		}
		if s == "knowledge/glossary/ai-skill.md" {
			glossaryStaged = true
		}
	}
	if !frameworkSurfaceStaged || glossaryStaged {
		return ""
	}
	return "glossary-retro-own: staged change touches framework cognitive vocabulary surface (runtime/cognitive-modes*.yaml, runtime/economics/, ecosystem/) but knowledge/glossary/ai-skill.md is not staged. Per plans/active/2026-05-25-1000-context-language-glossary-system.md Phase 6 and runtime/cli-modification-policy.yaml gate.glossary.retro_own_required, new framework terms must retro-own a canonical glossary entry. Use [skip-glossary-retro-own] (standalone trailer line) if this change is a typo / refactor / comment-only edit and introduces no new term."
}

// validateRuntimeYamlProjects enforces the rule "every runtime/*.yaml
// must declare runtime_projection.enabled: true AND target_key". Plans
// that intentionally defer projection must include §Deferred Runtime
// Projection in plan AND use [skip-runtime-yaml-projection] opt-out.
//
// Opt-out: standalone `[skip-runtime-yaml-projection]` trailer line.
func validateRuntimeYamlProjects(text string, staged []string, root string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-runtime-yaml-projection]" {
			return ""
		}
	}
	var violations []string
	for _, s := range staged {
		if !strings.HasPrefix(s, "runtime/") || !strings.HasSuffix(s, ".yaml") {
			continue
		}
		full := s
		if !filepath.IsAbs(full) {
			full = filepath.Join(root, s)
		}
		body, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		content := string(body)
		hasProjection := strings.Contains(content, "runtime_projection:") &&
			(strings.Contains(content, "enabled: true") || strings.Contains(content, "enabled:true"))
		hasTargetKey := strings.Contains(content, "target_key:")
		if !hasProjection || !hasTargetKey {
			violations = append(violations,
				"runtime-yaml-projects: "+s+" missing runtime_projection.enabled:true or target_key. Default rule: runtime/*.yaml must project to runtime.db. If intentional deferral, declare §Deferred Runtime Projection in plan + use [skip-runtime-yaml-projection].")
		}
	}
	if len(violations) == 0 {
		return ""
	}
	return strings.Join(violations, "\n  - ")
}

// validateMarkdownYamlSync enforces sibling-pair markdown/YAML
// synchronization: if a staged .md file has a sibling .yaml in the same
// directory (same path stem), the .yaml must also be staged in the same
// commit. Cross-path companion pairs (e.g. plans/README.md ↔
// governance/lifecycle/system-upgrade-governance.yaml) require explicit
// mapping; not yet covered by this validator.
//
// Opt-out: standalone `[skip-markdown-yaml-sync]` trailer line.
func validateMarkdownYamlSync(text string, staged []string, root string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "[skip-markdown-yaml-sync]" {
			return ""
		}
	}
	stagedSet := make(map[string]bool, len(staged))
	for _, s := range staged {
		stagedSet[s] = true
	}
	var violations []string
	for _, s := range staged {
		if !strings.HasSuffix(s, ".md") {
			continue
		}
		stem := strings.TrimSuffix(s, ".md")
		sibling := stem + ".yaml"
		siblingFull := sibling
		if !filepath.IsAbs(siblingFull) {
			siblingFull = filepath.Join(root, sibling)
		}
		if _, err := os.Stat(siblingFull); err != nil {
			// Sibling YAML does not exist on disk; no companion enforcement.
			continue
		}
		if !stagedSet[sibling] {
			violations = append(violations,
				"markdown-yaml-sync: "+s+" staged but sibling companion "+sibling+" not staged. Canonical .md edits typically need YAML companion update. If markdown-only change is intentional, use [skip-markdown-yaml-sync].")
		}
	}
	if len(violations) == 0 {
		return ""
	}
	return strings.Join(violations, "\n  - ")
}

// commitMsgCtx is the uniform context for per-obligation validators
// dispatched from per_commit_obligations in
// generated_surfaces[runtime.core_bootstrap.contract]. Phase 6 of
// bootstrap-contract-yaml-migration.
type commitMsgCtx struct {
	text   string
	staged []string
	root   string
	modes  map[string]string
}

// commitMsgValidatorRegistry maps obligation IDs (matching
// per_commit_obligations[].id in runtime/core-bootstrap.yaml) to
// adapter closures that call the actual validator with the right
// arguments. Obligations declared in YAML but not registered here
// are silently skipped (allows YAML to declare future-planned
// obligations without breaking the hook). Conversely, validators
// registered here that lack a YAML obligation will not fire unless
// they appear in defaultCommitMsgDispatchOrder fallback.
//
// obligation.commit.cognitive_mode_block is the GATE (block-presence
// check); it is handled inline before dispatch and is NOT in this
// registry.
var commitMsgValidatorRegistry = map[string]func(commitMsgCtx) string{
	"obligation.commit.execution_mode_floors": func(c commitMsgCtx) string {
		return validateExecutionModeFloors(c.modes, c.staged)
	},
	"obligation.commit.governance_mode_consistency": func(c commitMsgCtx) string {
		return validateGovernanceModeConsistency(c.modes, c.staged, c.text)
	},
	"obligation.commit.memory_mode_subdir": func(c commitMsgCtx) string {
		return validateMemoryModeSubdir(c.modes, c.staged)
	},
	"obligation.commit.cognitive_cost": func(c commitMsgCtx) string {
		return validateCognitiveCost(c.modes)
	},
	"obligation.commit.activation_signals": func(c commitMsgCtx) string {
		return validateActivationSignals(c)
	},
	"obligation.commit.capability_snippet": func(c commitMsgCtx) string {
		return validateCapabilitySnippet(c.modes, c.text)
	},
	"obligation.commit.plan_status_sync": func(c commitMsgCtx) string {
		return validatePlanStatusSync(c.text, c.staged)
	},
	"obligation.commit.token_budget": func(c commitMsgCtx) string {
		return validateTokenBudget(c.modes, c.text)
	},
	"obligation.commit.adaptive_triggers": func(c commitMsgCtx) string {
		return validateAdaptiveTriggers(c.modes, c.text)
	},
	"obligation.commit.bootstrap_entry_thinness": func(c commitMsgCtx) string {
		return validateBootstrapEntryThinness(c.text, c.staged, c.root)
	},
	"obligation.commit.cli_doc_sync": func(c commitMsgCtx) string {
		return validateCLIDocSync(c.text, c.staged, c.root)
	},
	"obligation.commit.runtime_yaml_projects": func(c commitMsgCtx) string {
		return validateRuntimeYamlProjects(c.text, c.staged, c.root)
	},
	"obligation.commit.markdown_yaml_sync": func(c commitMsgCtx) string {
		return validateMarkdownYamlSync(c.text, c.staged, c.root)
	},
	"obligation.commit.glossary_retro_own": func(c commitMsgCtx) string {
		return validateGlossaryRetroOwn(c.text, c.staged, c.root)
	},
}

// defaultCommitMsgDispatchOrder is the fallback order if
// runtime.core_bootstrap.contract is not available (fresh clone /
// pre-compile). Mirrors per_commit_obligations[] order in
// runtime/core-bootstrap.yaml, excluding cognitive_mode_block (gate).
var defaultCommitMsgDispatchOrder = []string{
	"obligation.commit.execution_mode_floors",
	"obligation.commit.governance_mode_consistency",
	"obligation.commit.memory_mode_subdir",
	"obligation.commit.cognitive_cost",
	"obligation.commit.activation_signals",
	"obligation.commit.capability_snippet",
	"obligation.commit.plan_status_sync",
	"obligation.commit.token_budget",
	"obligation.commit.adaptive_triggers",
	"obligation.commit.bootstrap_entry_thinness",
	"obligation.commit.cli_doc_sync",
	"obligation.commit.runtime_yaml_projects",
	"obligation.commit.markdown_yaml_sync",
	"obligation.commit.glossary_retro_own",
}

// readPerCommitObligationsOrder reads the per_commit_obligations id
// list from generated_surfaces[runtime.core_bootstrap.contract], in
// the order declared in runtime/core-bootstrap.yaml. The gate obligation
// (cognitive_mode_block) is filtered out; only post-gate validators
// are returned. Returns nil if any step fails (runtime.db missing,
// contract not projected, JSON malformed) — caller should fall back
// to defaultCommitMsgDispatchOrder.
func readPerCommitObligationsOrder(root string) []string {
	dbPath := filepath.Join(root, "runtime", "runtime.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil
	}
	defer db.Close()
	var raw string
	err = db.QueryRow("SELECT data FROM generated_surfaces WHERE target_key='runtime.core_bootstrap.contract' LIMIT 1").Scan(&raw)
	if err != nil {
		return nil
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil
	}
	arr, _ := doc["per_commit_obligations"].([]any)
	ids := make([]string, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id, _ := m["id"].(string)
		if id == "" || id == "obligation.commit.cognitive_mode_block" {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func runPrePushHook(result Result, root string) Result {
	changed, upstream, err := cliCIPrePushPaths(root)
	if err != nil {
		result.Checks = append(result.Checks, Check{Name: "cli_ci_scope", Status: "warning", Message: err.Error(), Remediation: "Running Go tests conservatively because changed paths could not be resolved."})
		changed = []string{"scripts/ai-skill-cli/"}
	} else {
		result.Checks = append(result.Checks, Check{Name: "cli_ci_scope", Status: "ok", Message: "compared against " + upstream})
	}
	if !hasCLICIPreflightChange(changed) {
		result.Checks = append(result.Checks, Check{Name: "cli_ci_preflight", Status: "skipped", Message: "no CLI, hook, or workflow changes since upstream"})
		return result
	}
	result.Checks = append(result.Checks, githubWorkflowHistoryCheck(root))
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = filepath.Join(root, "scripts", "ai-skill-cli")
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{
			Code:        "cli_go_test_failed",
			Message:     compactSummary(string(output)),
			Remediation: "Run `cd scripts/ai-skill-cli && go test ./...`; if repo-local binaries are stale, rebuild `bin/` with `go run ./cmd/releasebuild --stable-names --version \"repo-$(git rev-parse --short HEAD)\" --commit \"$(git rev-parse --short HEAD)\" --dist bin`.",
		}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "cli_ci_preflight", Status: "ok", Message: "go test ./... passed"})
	return result
}

type githubWorkflowRunsResponse struct {
	WorkflowRuns []struct {
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
		HTMLURL    string `json:"html_url"`
		HeadSHA    string `json:"head_sha"`
	} `json:"workflow_runs"`
}

func githubWorkflowHistoryCheck(root string) Check {
	remote, err := exec.Command("git", "-C", root, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: "remote.origin.url unavailable"}
	}
	owner, repo, ok := parseGitHubRemote(strings.TrimSpace(string(remote)))
	if !ok {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: "origin is not a GitHub remote"}
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/ai-skill-cli.yml/runs?per_page=1&status=completed", owner, repo)
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: err.Error()}
	}
	req.Header.Set("User-Agent", "ai-skill-cli-pre-push")
	resp, err := client.Do(req)
	if err != nil {
		return Check{Name: "github_workflow_history", Status: "warning", Message: "could not query latest completed GitHub workflow: " + err.Error(), Remediation: "Continuing with local Go preflight; check GitHub Actions manually if needed."}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Check{Name: "github_workflow_history", Status: "warning", Message: fmt.Sprintf("GitHub workflow query returned HTTP %d", resp.StatusCode), Remediation: "Continuing with local Go preflight; check GitHub Actions manually if needed."}
	}
	var runs githubWorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&runs); err != nil {
		return Check{Name: "github_workflow_history", Status: "warning", Message: "could not decode GitHub workflow response: " + err.Error(), Remediation: "Continuing with local Go preflight; check GitHub Actions manually if needed."}
	}
	if len(runs.WorkflowRuns) == 0 {
		return Check{Name: "github_workflow_history", Status: "skipped", Message: "no completed ai-skill CLI workflow runs found"}
	}
	run := runs.WorkflowRuns[0]
	message := fmt.Sprintf("latest completed run %s at %s", run.Conclusion, run.HTMLURL)
	if run.Conclusion != "success" {
		return Check{Name: "github_workflow_history", Status: "warning", Message: message, Remediation: "Previous completed GitHub workflow is red; local go test must pass before pushing a fix."}
	}
	return Check{Name: "github_workflow_history", Status: "ok", Message: message}
}

func parseGitHubRemote(remote string) (string, string, bool) {
	remote = strings.TrimSuffix(strings.TrimSpace(remote), ".git")
	switch {
	case strings.HasPrefix(remote, "https://github.com/"):
		parts := strings.Split(strings.TrimPrefix(remote, "https://github.com/"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], true
		}
	case strings.HasPrefix(remote, "git@github.com:"):
		parts := strings.Split(strings.TrimPrefix(remote, "git@github.com:"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], true
		}
	case strings.HasPrefix(remote, "ssh://git@github.com/"):
		parts := strings.Split(strings.TrimPrefix(remote, "ssh://git@github.com/"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], true
		}
	}
	return "", "", false
}

func cliCIPrePushPaths(root string) ([]string, string, error) {
	upstreamOutput, err := exec.Command("git", "-C", root, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}").Output()
	if err != nil {
		changed, fallbackErr := gitLines(root, "diff", "--name-only", "HEAD^...HEAD")
		if fallbackErr != nil {
			return []string{}, "HEAD", nil
		}
		return changed, "HEAD^...HEAD", nil
	}
	upstream := strings.TrimSpace(string(upstreamOutput))
	changed, err := gitLines(root, "diff", "--name-only", upstream+"...HEAD")
	if err != nil {
		return nil, upstream, err
	}
	return changed, upstream, nil
}

func hasCLICIPreflightChange(paths []string) bool {
	for _, path := range paths {
		if strings.HasPrefix(path, "scripts/ai-skill-cli/") ||
			strings.HasPrefix(path, "scripts/git-hooks/") ||
			path == ".github/workflows/ai-skill-cli.yml" {
			return true
		}
	}
	return false
}

func gitLines(root string, args ...string) ([]string, error) {
	output, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines, nil
}

func hasRuntimeSourceChange(paths []string) bool {
	return false
}

func hasRuntimeDBChange(paths []string) bool {
	for _, path := range paths {
		if path == "runtime/runtime.db" {
			return true
		}
	}
	return false
}

func hasKnowledgeValidationChange(paths []string) bool {
	for _, path := range paths {
		if strings.HasPrefix(path, "knowledge/") || strings.HasPrefix(path, "validation/") || strings.HasPrefix(path, "scripts/validate-") {
			return true
		}
	}
	return false
}

func hookAdapterContent(hook string) string {
	return fmt.Sprintf(`#!/usr/bin/env sh
set -eu
ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
BIN="${AI_SKILL_CLI:-$ROOT/scripts/ai-skill-cli/bin/ai-skill}"
if [ ! -x "$BIN" ]; then
  case "$(uname -s 2>/dev/null | tr '[:upper:]' '[:lower:]')" in
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
  BIN="$ROOT/scripts/ai-skill-cli/bin/ai-skill-$os-$arch$suffix"
fi
exec "$BIN" hooks run %s --repo "$ROOT" "$@"
`, hook)
}

func hookSourceCheck(sourceDir string) Check {
	normalized, err := pathutil.NormalizeForReport(sourceDir)
	if err != nil {
		normalized = sourceDir
	}
	info, err := os.Stat(sourceDir)
	if err != nil {
		return Check{Name: "hook_source", Status: "missing", Message: normalized}
	}
	if !info.IsDir() {
		return Check{Name: "hook_source", Status: "failed", Message: normalized + " is not a directory"}
	}
	hooks := listHookFiles(sourceDir)
	if len(hooks) == 0 {
		return Check{Name: "hook_source", Status: "failed", Message: normalized + " contains no hook files"}
	}
	return Check{Name: "hook_source", Status: "ok", Message: fmt.Sprintf("%s (%d hooks)", normalized, len(hooks))}
}

func gitHooksTargetDir(root string) (string, Check) {
	gitDirOutput, err := exec.Command("git", "-C", root, "rev-parse", "--git-dir").Output()
	if err != nil {
		return "", Check{Name: "hook_target", Status: "failed", Message: "cannot resolve git dir"}
	}
	gitDir := strings.TrimSpace(string(gitDirOutput))
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(root, gitDir)
	}
	targetDir := filepath.Join(gitDir, "hooks")
	normalized, err := pathutil.NormalizeForReport(targetDir)
	if err != nil {
		normalized = targetDir
	}
	info, err := os.Stat(targetDir)
	if err != nil {
		return "", Check{Name: "hook_target", Status: "missing", Message: normalized}
	}
	if !info.IsDir() {
		return "", Check{Name: "hook_target", Status: "failed", Message: normalized + " is not a directory"}
	}
	return targetDir, Check{Name: "hook_target", Status: "ok", Message: normalized}
}

func listHookFiles(sourceDir string) []string {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil
	}
	hooks := []string{}
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			hooks = append(hooks, entry.Name())
		}
	}
	sort.Strings(hooks)
	return hooks
}

func existingHookTargets(targetDir string, hooks []string) []string {
	conflicts := []string{}
	for _, hook := range hooks {
		target := filepath.Join(targetDir, hook)
		if _, err := os.Stat(target); err == nil {
			conflicts = append(conflicts, target)
		}
	}
	sort.Strings(conflicts)
	return conflicts
}
