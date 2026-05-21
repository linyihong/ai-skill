package app

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
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

	result.Checks = append(result.Checks, pathCheck())

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
		result.Checks = append(result.Checks, hooksPathCheck())
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
		result.Checks = append(result.Checks, rubyCheck())
		result.Checks = append(result.Checks, pythonCheck())
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

func pathCheck() Check {
	pathValue := os.Getenv("PATH")
	if pathValue == "" {
		return Check{
			Name:        "path",
			Status:      "missing",
			Message:     "PATH is empty",
			Remediation: "Set PATH so external dependencies such as Git can be discovered.",
		}
	}

	summary := pathutil.SummarizePathList(pathValue)
	if summary.EmptyEntries > 0 {
		return Check{
			Name:        "path",
			Status:      "warning",
			Message:     fmt.Sprintf("%d PATH entries, %d empty", summary.Entries, summary.EmptyEntries),
			Remediation: "Remove empty PATH entries to avoid platform-specific command lookup surprises.",
		}
	}

	return Check{Name: "path", Status: "ok", Message: fmt.Sprintf("%d entries", summary.Entries)}
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
	normalized, err := pathutil.NormalizeForReport(strings.TrimSpace(string(output)))
	if err != nil {
		return Check{Name: "repo_root", Status: "failed", Message: err.Error()}
	}
	return Check{Name: "repo_root", Status: "ok", Message: normalized}
}

func hooksPathCheck() Check {
	output, err := exec.Command("git", "config", "--get", "core.hooksPath").Output()
	if err != nil {
		return Check{Name: "hooks_path", Status: "unset", Message: "core.hooksPath is not configured"}
	}
	value := strings.TrimSpace(string(output))
	if value == "" {
		return Check{Name: "hooks_path", Status: "unset", Message: "core.hooksPath is empty"}
	}
	normalized, err := pathutil.NormalizeForReport(value)
	if err != nil {
		return Check{Name: "hooks_path", Status: "failed", Message: err.Error()}
	}
	return Check{Name: "hooks_path", Status: "ok", Message: normalized}
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

	fileCheck := nativeSQLiteFileBackedProof()
	if fileCheck.Status != "ok" {
		return fileCheck
	}
	return Check{Name: "native_sqlite", Status: "ok", Message: "modernc.org/sqlite in-memory and file-backed write/query/integrity checks succeeded"}
}

func nativeSQLiteFileBackedProof() Check {
	file, err := os.CreateTemp("", "ai-skill-sqlite-proof-*.db")
	if err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	path := file.Name()
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	defer os.Remove(path)

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE proof (id INTEGER PRIMARY KEY, label TEXT NOT NULL)"); err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	if _, err := db.Exec("INSERT INTO proof (label) VALUES (?)", "native"); err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}

	var label string
	if err := db.QueryRow("SELECT label FROM proof WHERE id = 1").Scan(&label); err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	if label != "native" {
		return Check{Name: "native_sqlite", Status: "failed", Message: "unexpected file-backed SQLite result"}
	}

	var integrity string
	if err := db.QueryRow("PRAGMA integrity_check").Scan(&integrity); err != nil {
		return Check{Name: "native_sqlite", Status: "failed", Message: err.Error()}
	}
	if integrity != "ok" {
		return Check{Name: "native_sqlite", Status: "failed", Message: integrity}
	}
	return Check{Name: "native_sqlite", Status: "ok"}
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
	normalized, err := pathutil.NormalizeForReport(path)
	if err != nil {
		return Check{Name: "runtime_db", Status: "failed", Message: err.Error()}
	}
	return Check{Name: "runtime_db", Status: "ok", Message: normalized}
}

func rubyCheck() Check {
	return executableVersionCheck("ruby", []string{"--version"}, "Ruby is only required for wrapper-mode runtime compiler / validators.")
}

func pythonCheck() Check {
	if check := executableVersionCheck("python3", []string{"--version"}, "Python is only required for wrapper-mode helpers."); check.Status == "ok" {
		check.Name = "python"
		return check
	}
	check := executableVersionCheck("python", []string{"--version"}, "Python is only required for wrapper-mode helpers.")
	check.Name = "python"
	return check
}

func executableVersionCheck(name string, args []string, remediation string) Check {
	path, err := exec.LookPath(name)
	if err != nil {
		return Check{Name: name, Status: "missing_optional", Message: name + " not found in PATH", Remediation: remediation}
	}
	output, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return Check{Name: name, Status: "failed", Message: strconv.Quote(strings.TrimSpace(string(output)))}
	}
	return Check{Name: name, Status: "ok", Message: strings.TrimSpace(string(output))}
}
