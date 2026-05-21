package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type runtimeOptions struct {
	command       string
	repoPath      string
	dryRun        bool
	assertSource  string
	assertKeyword string
	jsonOutput    bool
	plainOutput   bool
}

type runtimeValidator struct {
	name string
	path string
}

func runRuntime(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill runtime <validate|refresh|compile> [flags]")
		return ExitInvalidUsage
	}
	opts := runtimeOptions{command: args[0]}
	if opts.command != "validate" && opts.command != "refresh" && opts.command != "compile" {
		_, _ = fmt.Fprintf(stderr, "unsupported runtime command: %s\n", opts.command)
		return ExitInvalidUsage
	}

	fs := newFlagSet("runtime "+opts.command, stderr)
	fs.StringVar(&opts.repoPath, "repo", ".", "Ai-skill repository path")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview runtime wrapper scripts without executing")
	fs.StringVar(&opts.assertSource, "assert-source", "", "source path expected in generated surfaces")
	fs.StringVar(&opts.assertKeyword, "assert-keyword", "", "keyword expected in generated surfaces")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildRuntimeResult(opts)
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

func buildRuntimeResult(opts runtimeOptions) Result {
	switch opts.command {
	case "compile":
		return buildRuntimeCompileResult(opts)
	case "refresh":
		return buildRuntimeRefreshResult(opts)
	default:
		return buildRuntimeValidateResult(opts)
	}
}

func buildRuntimeValidateResult(opts runtimeOptions) Result {
	result := Result{
		Command:        "runtime validate",
		Mode:           "wrapper",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	if opts.dryRun {
		result.Mode = "dry_run"
	}

	repo, repoCheck := resolveExistingDir("repo", opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message, Remediation: "Pass --repo with the Ai-skill repository root."}
		return result
	}

	validators := runtimeValidators(repo)
	for _, validator := range validators {
		result.PlannedActions = append(result.PlannedActions, "run ruby validator: "+validator.path)
		if _, err := os.Stat(validator.path); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "missing_validator", Message: validator.path, Remediation: "Run from a complete Ai-skill checkout."}
			result.Checks = append(result.Checks, Check{Name: validator.name, Status: "missing", Message: validator.path})
			return result
		}
	}

	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "wrapper_mode", Status: "ok", Message: "dry-run only; validators not executed"})
		return result
	}

	ruby, rubyCheck := requiredExecutable("ruby", []string{"--version"}, "Install Ruby to use wrapper-mode runtime validators until native Go validators replace them.")
	result.Checks = append(result.Checks, rubyCheck)
	if rubyCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_ruby", Message: "Ruby is required for runtime validate wrapper mode.", Remediation: rubyCheck.Remediation}
		return result
	}

	sqlite, sqliteCheck := requiredExecutable("sqlite3", []string{"--version"}, "Install sqlite3 CLI for wrapper-mode runtime DB and runtime index validators until native Go validators replace them.")
	result.Checks = append(result.Checks, sqliteCheck)
	if sqliteCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_sqlite3", Message: "sqlite3 CLI is required for runtime validate wrapper mode.", Remediation: sqliteCheck.Remediation}
		return result
	}
	_ = sqlite

	for _, validator := range validators {
		check := runRuntimeValidator(repo, ruby, validator)
		result.Checks = append(result.Checks, check)
		if check.Status != "ok" {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "runtime_validation_failed", Message: check.Message, Remediation: "Inspect validator output and fix the runtime source or generated surface."}
			return result
		}
	}

	return result
}

func buildRuntimeRefreshResult(opts runtimeOptions) Result {
	result := Result{
		Command:        "runtime refresh",
		Mode:           "wrapper",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	if opts.dryRun {
		result.Mode = "dry_run"
	}

	repo, repoCheck := resolveExistingDir("repo", opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message, Remediation: "Pass --repo with the Ai-skill repository root."}
		return result
	}

	for _, script := range runtimeRefreshScripts(repo) {
		result.PlannedActions = append(result.PlannedActions, "run ruby refresh step: "+script)
		if _, err := os.Stat(script); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "missing_runtime_script", Message: script, Remediation: "Run from a complete Ai-skill checkout."}
			result.Checks = append(result.Checks, Check{Name: filepath.Base(script), Status: "missing", Message: script})
			return result
		}
	}

	wrapper := filepath.Join(repo, "scripts", "refresh-knowledge-runtime.rb")
	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "wrapper_mode", Status: "ok", Message: "dry-run only; refresh wrapper not executed"})
		return result
	}

	ruby, rubyCheck := requiredExecutable("ruby", []string{"--version"}, "Install Ruby to use wrapper-mode runtime refresh until native Go refresh replaces it.")
	result.Checks = append(result.Checks, rubyCheck)
	if rubyCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_ruby", Message: "Ruby is required for runtime refresh wrapper mode.", Remediation: rubyCheck.Remediation}
		return result
	}

	sqlite, sqliteCheck := requiredExecutable("sqlite3", []string{"--version"}, "Install sqlite3 CLI for wrapper-mode runtime refresh until native Go SQLite refresh replaces it.")
	result.Checks = append(result.Checks, sqliteCheck)
	if sqliteCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_sqlite3", Message: "sqlite3 CLI is required for runtime refresh wrapper mode.", Remediation: sqliteCheck.Remediation}
		return result
	}
	_ = sqlite

	git, gitCheck := requiredExecutable("git", []string{"--version"}, "Install Git because wrapper-mode runtime refresh validates generated SQLite index git-ignore boundaries.")
	result.Checks = append(result.Checks, gitCheck)
	if gitCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_git", Message: "Git is required for runtime refresh wrapper mode.", Remediation: gitCheck.Remediation}
		return result
	}
	_ = git

	check := runRuntimeScript(repo, ruby, "runtime_refresh", wrapper)
	result.Checks = append(result.Checks, check)
	if check.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_refresh_failed", Message: check.Message, Remediation: "Inspect refresh output and fix the failing generator or validator."}
		return result
	}

	return result
}

func buildRuntimeCompileResult(opts runtimeOptions) Result {
	result := Result{
		Command:        "runtime compile",
		Mode:           "wrapper",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	if opts.dryRun {
		result.Mode = "dry_run"
	}

	repo, repoCheck := resolveExistingDir("repo", opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message, Remediation: "Pass --repo with the Ai-skill repository root."}
		return result
	}

	compiler := filepath.Join(repo, "runtime", "compiler", "compiler-engine.rb")
	result.PlannedActions = append(result.PlannedActions, "run ruby compiler: "+compiler)
	if opts.dryRun {
		result.PlannedActions = append(result.PlannedActions, "run ruby compiler check: "+compiler+" --diff")
	}
	if opts.assertSource != "" || opts.assertKeyword != "" {
		result.PlannedActions = append(result.PlannedActions, "assert generated surface: source="+opts.assertSource+" keyword="+opts.assertKeyword)
	}
	if _, err := os.Stat(compiler); err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "missing_runtime_compiler", Message: compiler, Remediation: "Run from a complete Ai-skill checkout."}
		result.Checks = append(result.Checks, Check{Name: "runtime_compiler", Status: "missing", Message: compiler})
		return result
	}

	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "wrapper_mode", Status: "ok", Message: "dry-run only; compiler not executed"})
		return result
	}

	ruby, rubyCheck := requiredExecutable("ruby", []string{"--version"}, "Install Ruby to use wrapper-mode runtime compiler until native Go compiler replaces it.")
	result.Checks = append(result.Checks, rubyCheck)
	if rubyCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_ruby", Message: "Ruby is required for runtime compile wrapper mode.", Remediation: rubyCheck.Remediation}
		return result
	}

	sqlite, sqliteCheck := requiredExecutable("sqlite3", []string{"--version"}, "Install sqlite3 CLI for wrapper-mode runtime compiler until native Go SQLite compile replaces it.")
	result.Checks = append(result.Checks, sqliteCheck)
	if sqliteCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{Code: "missing_sqlite3", Message: "sqlite3 CLI is required for runtime compile wrapper mode.", Remediation: sqliteCheck.Remediation}
		return result
	}
	_ = sqlite

	check := runRuntimeScript(repo, ruby, "runtime_compile", compiler)
	result.Checks = append(result.Checks, check)
	if check.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_compile_failed", Message: check.Message, Remediation: "Inspect compiler output and fix the runtime source or compiler input."}
		return result
	}

	return result
}

func runtimeValidators(repo string) []runtimeValidator {
	return []runtimeValidator{
		{name: "knowledge_runtime", path: filepath.Join(repo, "scripts", "validate-knowledge-runtime.rb")},
		{name: "runtime_db", path: filepath.Join(repo, "scripts", "validate-runtime-db.rb")},
		{name: "runtime_sqlite_index", path: filepath.Join(repo, "scripts", "validate-runtime-sqlite-index.rb")},
	}
}

func runtimeRefreshScripts(repo string) []string {
	names := []string{
		"generate-model-context-report.rb",
		"generate-model-checklists.rb",
		"generate-knowledge-runtime-report.rb",
		"generate-runtime-sqlite-index.rb",
		"validate-runtime-sqlite-index.rb",
		"validate-knowledge-runtime.rb",
		"refresh-knowledge-runtime.rb",
	}
	scripts := make([]string, 0, len(names))
	for _, name := range names {
		scripts = append(scripts, filepath.Join(repo, "scripts", name))
	}
	return scripts
}

func requiredExecutable(name string, args []string, remediation string) (string, Check) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", Check{Name: name, Status: "missing", Message: name + " not found in PATH", Remediation: remediation}
	}
	output, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return path, Check{Name: name, Status: "failed", Message: strings.TrimSpace(string(output)), Remediation: remediation}
	}
	return path, Check{Name: name, Status: "ok", Message: strings.TrimSpace(string(output))}
}

func runRuntimeValidator(repo string, ruby string, validator runtimeValidator) Check {
	return runRuntimeScript(repo, ruby, validator.name, validator.path)
}

func runRuntimeScript(repo string, ruby string, name string, path string) Check {
	cmd := exec.Command(ruby, path)
	cmd.Dir = repo
	cmd.Env = runtimeWrapperEnv(os.Environ())
	output, err := cmd.CombinedOutput()
	message := strings.TrimSpace(string(output))
	if err != nil {
		if message == "" {
			message = err.Error()
		}
		return Check{Name: name, Status: "failed", Message: message}
	}
	if message == "" {
		message = "ok"
	}
	return Check{Name: name, Status: "ok", Message: message}
}

func runtimeWrapperEnv(base []string) []string {
	env := append([]string{}, base...)
	env = upsertEnv(env, "LANG", "C.UTF-8")
	env = upsertEnv(env, "LC_ALL", "C.UTF-8")
	return env
}

func upsertEnv(env []string, key string, value string) []string {
	prefix := key + "="
	for i, existing := range env {
		if strings.HasPrefix(existing, prefix) {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}
