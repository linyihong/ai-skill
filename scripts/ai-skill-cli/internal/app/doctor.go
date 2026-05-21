package app

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	_ "modernc.org/sqlite"
)

type doctorOptions struct {
	jsonOutput   bool
	plainOutput  bool
	requireGit   bool
	requireWrite bool
	checkRuntime bool
}

func runDoctor(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := newFlagSet("doctor", stderr)
	opts := doctorOptions{}
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	fs.BoolVar(&opts.requireGit, "require-git", false, "block when git is not available")
	fs.BoolVar(&opts.requireWrite, "require-write", false, "check current directory write permission")
	fs.BoolVar(&opts.checkRuntime, "check-runtime", false, "check runtime database presence")

	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildDoctorResult(opts)
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

func buildDoctorResult(opts doctorOptions) Result {
	result := Result{
		Command:        "doctor",
		Mode:           "check",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}

	platform := platformCheck()
	result.Checks = append(result.Checks, platform)
	if platform.Status == "unsupported" {
		result.Status = "blocked"
		result.ExitCode = ExitUnsupportedPlatform
		result.Error = &CommandError{
			Code:        "unsupported_platform",
			Message:     "This platform is not supported for native ai-skill execution.",
			Remediation: platform.Remediation,
		}
	}

	gitCheck := checkGit()
	result.Checks = append(result.Checks, gitCheck)
	if opts.requireGit && gitCheck.Status != "ok" && result.ExitCode == ExitSuccess {
		result.Status = "blocked"
		result.ExitCode = ExitMissingDependency
		result.Error = &CommandError{
			Code:        "missing_git",
			Message:     "Git is required for this command but was not found in PATH.",
			Remediation: "Install Git and ensure the git executable is available in PATH.",
		}
	}

	if gitCheck.Status == "ok" {
		result.Checks = append(result.Checks, repoRootCheck())
	}

	if opts.requireWrite {
		writeCheck := checkWritePermission(".")
		result.Checks = append(result.Checks, writeCheck)
		if writeCheck.Status == "failed" && result.ExitCode == ExitSuccess {
			result.Status = "blocked"
			result.ExitCode = ExitPermissionDenied
			result.Error = &CommandError{
				Code:        "permission_denied",
				Message:     "Current directory is not writable.",
				Remediation: "Choose a writable working directory or adjust filesystem permissions.",
			}
		}
	}

	if opts.checkRuntime {
		result.Checks = append(result.Checks, nativeSQLiteCheck())
		result.Checks = append(result.Checks, runtimeDBCheck())
	}

	return result
}

func platformCheck() Check {
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		return Check{Name: "platform", Status: "ok", Message: runtime.GOOS}
	case "ios":
		return Check{
			Name:        "platform",
			Status:      "unsupported",
			Message:     "iOS is not a native arbitrary binary target.",
			Remediation: "Use an app-contained runtime, Browser/WASM, or SSH remote runner.",
		}
	default:
		return Check{Name: "platform", Status: "unsupported", Message: runtime.GOOS}
	}
}

func checkGit() Check {
	path, err := exec.LookPath("git")
	if err != nil {
		return Check{
			Name:        "git",
			Status:      "missing",
			Message:     "git executable not found in PATH",
			Remediation: "Install Git and ensure the git executable is available in PATH.",
		}
	}

	version, err := exec.Command(path, "--version").Output()
	if err != nil {
		return Check{Name: "git", Status: "failed", Message: err.Error()}
	}

	return Check{Name: "git", Status: "ok", Message: strings.TrimSpace(string(version))}
}

func repoRootCheck() Check {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return Check{Name: "repo_root", Status: "failed", Message: "not inside a Git work tree"}
	}
	return Check{Name: "repo_root", Status: "ok", Message: strings.TrimSpace(string(output))}
}

func checkWritePermission(dir string) Check {
	file, err := os.CreateTemp(dir, ".ai-skill-write-check-*")
	if err != nil {
		return Check{Name: "write_permission", Status: "failed", Message: err.Error()}
	}
	name := file.Name()
	closeErr := file.Close()
	removeErr := os.Remove(name)
	if closeErr != nil {
		return Check{Name: "write_permission", Status: "failed", Message: closeErr.Error()}
	}
	if removeErr != nil {
		return Check{Name: "write_permission", Status: "failed", Message: removeErr.Error()}
	}
	return Check{Name: "write_permission", Status: "ok", Message: dir}
}

func runtimeDBCheck() Check {
	candidates := []string{
		filepath.Join("runtime", "runtime.db"),
		filepath.Join("..", "..", "runtime", "runtime.db"),
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return runtimeDBIntegrityCheck(candidate)
		}
	}
	return Check{Name: "runtime_db", Status: "missing", Message: "runtime/runtime.db not found from current working directory"}
}

func nativeSQLiteCheck() Check {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	defer db.Close()

	var value int
	if err := db.QueryRow("SELECT 1").Scan(&value); err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	if value != 1 {
		return Check{Name: "native_sqlite", Status: "failed", Message: "unexpected SQLite result"}
	}
	return Check{Name: "native_sqlite", Status: "ok", Message: "modernc.org/sqlite in-memory query succeeded"}
}

func runtimeDBIntegrityCheck(path string) Check {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return Check{Name: "runtime_db", Status: "failed", Message: err.Error()}
	}
	defer db.Close()

	var result string
	if err := db.QueryRow("PRAGMA integrity_check").Scan(&result); err != nil {
		return Check{Name: "runtime_db", Status: "failed", Message: err.Error()}
	}
	if result != "ok" {
		return Check{Name: "runtime_db", Status: "failed", Message: result}
	}
	return Check{Name: "runtime_db", Status: "ok", Message: path}
}
