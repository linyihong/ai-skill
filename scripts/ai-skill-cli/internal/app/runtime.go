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
	command     string
	repoPath    string
	dryRun      bool
	jsonOutput  bool
	plainOutput bool
}

type runtimeValidator struct {
	name string
	path string
}

func runRuntime(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill runtime validate [flags]")
		return ExitInvalidUsage
	}
	opts := runtimeOptions{command: args[0]}
	if opts.command != "validate" {
		_, _ = fmt.Fprintf(stderr, "unsupported runtime command: %s\n", opts.command)
		return ExitInvalidUsage
	}

	fs := newFlagSet("runtime validate", stderr)
	fs.StringVar(&opts.repoPath, "repo", ".", "Ai-skill repository path")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview runtime validators without executing")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildRuntimeValidateResult(opts)
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

func runtimeValidators(repo string) []runtimeValidator {
	return []runtimeValidator{
		{name: "knowledge_runtime", path: filepath.Join(repo, "scripts", "validate-knowledge-runtime.rb")},
		{name: "runtime_db", path: filepath.Join(repo, "scripts", "validate-runtime-db.rb")},
		{name: "runtime_sqlite_index", path: filepath.Join(repo, "scripts", "validate-runtime-sqlite-index.rb")},
	}
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
	cmd := exec.Command(ruby, validator.path)
	cmd.Dir = repo
	cmd.Env = runtimeWrapperEnv(os.Environ())
	output, err := cmd.CombinedOutput()
	message := strings.TrimSpace(string(output))
	if err != nil {
		if message == "" {
			message = err.Error()
		}
		return Check{Name: validator.name, Status: "failed", Message: message}
	}
	if message == "" {
		message = "ok"
	}
	return Check{Name: validator.name, Status: "ok", Message: message}
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
