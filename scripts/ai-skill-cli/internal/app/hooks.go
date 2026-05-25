package app

import (
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
			_, _ = fmt.Fprintln(stderr, "usage: ai-skill hooks run <pre-commit|commit-msg|post-commit|pre-push> [flags]")
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
	if len(result.Mutations) == 0 {
		result.Checks = append(result.Checks, Check{Name: "pre_commit", Status: "ok", Message: "no runtime or knowledge hook action required"})
	}
	return result
}

// runCommitMsgHook enforces Phase 4 behavioral wiring of
// gate.execution.cognitive_mode_resolved. Commit message body must contain
// the literal '### Cognitive Mode 報告' block (template defined in
// models/cognitive-modes/README.md). Mechanical commits may opt out via
// '[skip-cognitive-mode]' in the body. Merge commits auto-skip.
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

	// Primary path: Cognitive Mode 報告 block present → run Phase 3 behavioral validators.
	// Checked BEFORE opt-out marker to avoid false positives when commit body
	// documents/quotes the opt-out token (e.g. "Opt-out via '[skip-cognitive-mode]'").
	if strings.Contains(text, "### Cognitive Mode 報告") {
		modes := parseCognitiveModeBlock(text)
		staged, _ := gitLines(root, "diff", "--cached", "--name-only")

		// Phase 3.1-B / 3.3-B / 3.4-B behavioral validators
		var violations []string
		if v := validateExecutionModeFloors(modes, staged); v != "" {
			violations = append(violations, v)
		}
		if v := validateGovernanceModeConsistency(modes, staged, text); v != "" {
			violations = append(violations, v)
		}
		if v := validateMemoryModeSubdir(modes, staged); v != "" {
			violations = append(violations, v)
		}
		if v := validatePlanStatusSync(text, staged); v != "" {
			violations = append(violations, v)
		}
		if v := validateTokenBudget(modes, text); v != "" {
			violations = append(violations, v)
		}
		if v := validateAdaptiveTriggers(modes, text); v != "" {
			violations = append(violations, v)
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
		Message:     "Commit message body must include '### Cognitive Mode 報告' block (4-dim execution/context/governance/memory resolution).",
		Remediation: "Add the block per models/cognitive-modes/README.md template, or add a standalone '[skip-cognitive-mode]' trailer line for mechanical commits.",
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
		case "execution_mode", "context_mode", "governance_mode", "memory_mode":
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
//   1. ≥1 completion vocabulary word
//   2. ≥1 "Phase <num>" / "phase <num>" reference
//   3. ≥1 plans/active/<f>.md path reference
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
		"FAST|INDEX_ONLY|LIGHT|NONE":                       1000,
		"NORMAL|SUMMARY_FIRST|STANDARD|EPISODIC":           5000,
		"DEEP|SOURCE_BACKED|STRICT|DECISION_REPLAY":        20000,
		"FORENSIC|GRAPH_ASSISTED|STRICT|FAILURE_REPLAY":    50000,
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
