package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRuntimeValidateDryRunPlansValidators(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "validate", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "runtime validate" || result.Mode != "dry_run" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.PlannedActions) != 3 {
		t.Fatalf("expected three planned validators, got %#v", result.PlannedActions)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime validate dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeValidateBlocksMissingRubyBeforeWrapper(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "validate", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_ruby" {
		t.Fatalf("expected missing_ruby, got %#v", result.Error)
	}
}

func TestRuntimeRefreshDryRunPlansWrapperCommands(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "runtime refresh" || result.Mode != "dry_run" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.PlannedActions) != 7 {
		t.Fatalf("expected seven planned refresh scripts, got %#v", result.PlannedActions)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime refresh dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeRefreshBlocksMissingRubyBeforeWrapper(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_ruby" {
		t.Fatalf("expected missing_ruby, got %#v", result.Error)
	}
}

func TestRuntimeCompileDryRunPlansCompiler(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "compile", "--repo", repo, "--dry-run", "--assert-source", "runtime/compiler/embedded_data.rb", "--assert-keyword", "phase", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "runtime compile" || result.Mode != "dry_run" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime compile dry-run must not mutate, got %#v", result.Mutations)
	}
	if !hasCheckStatus(result.Checks, "wrapper_mode", "ok") {
		t.Fatalf("expected wrapper_mode ok, got %#v", result.Checks)
	}
}

func TestRuntimeCompileBlocksMissingRubyBeforeWrapper(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "compile", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitMissingDependency {
		t.Fatalf("expected missing dependency, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_ruby" {
		t.Fatalf("expected missing_ruby, got %#v", result.Error)
	}
}

func TestRuntimeValidateBlocksMissingValidator(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "scripts", "validate-knowledge-runtime.rb"), "# ok\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "validate", "--repo", repo, "--dry-run", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s", code, stderr.String())
	}
}

func TestRuntimeUnsupportedSubcommandReturnsInvalidUsage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "migrate"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unsupported runtime command") {
		t.Fatalf("expected unsupported message, got %q", stderr.String())
	}
}

func TestRuntimeWrapperEnvForcesUTF8(t *testing.T) {
	env := runtimeWrapperEnv([]string{"PATH=/bin", "LANG=C"})
	if !containsEnv(env, "LANG=C.UTF-8") {
		t.Fatalf("expected LANG override, got %#v", env)
	}
	if !containsEnv(env, "LC_ALL=C.UTF-8") {
		t.Fatalf("expected LC_ALL override, got %#v", env)
	}
}

func TestNativeRuntimeDBValidationAcceptsValidFixture(t *testing.T) {
	path := createNativeRuntimeDBFixture(t)
	check := nativeRuntimeDBValidation(path)
	if check.Status != "ok" {
		t.Fatalf("expected native validation ok, got %#v", check)
	}
}

func TestNativeRuntimeDBValidationBlocksMissingTable(t *testing.T) {
	path := createNativeRuntimeDBFixture(t)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("DROP TABLE gates"); err != nil {
		t.Fatal(err)
	}

	check := nativeRuntimeDBValidation(path)
	if check.Status != "failed" || !strings.Contains(check.Message, "missing required table: gates") {
		t.Fatalf("expected missing table failure, got %#v", check)
	}
}

func TestNativeRuntimeDBValidationBlocksInvalidJSON(t *testing.T) {
	path := createNativeRuntimeDBFixture(t)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("UPDATE phases SET entry_conditions = ? WHERE id = ?", "{bad", "phases-0"); err != nil {
		t.Fatal(err)
	}

	check := nativeRuntimeDBValidation(path)
	if check.Status != "failed" || !strings.Contains(check.Message, "phases.entry_conditions") {
		t.Fatalf("expected invalid JSON failure, got %#v", check)
	}
}

func TestNativeRuntimeDBValidationReportsStaleCompilerMetadata(t *testing.T) {
	path := createNativeRuntimeDBFixture(t)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("UPDATE compiler_metadata SET value = ? WHERE key = 'compiled_at'", "2026-05-19T00:00:00Z"); err != nil {
		t.Fatal(err)
	}

	originalNow := nativeRuntimeNow
	nativeRuntimeNow = func() time.Time { return time.Date(2026, 5, 21, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nativeRuntimeNow = originalNow })

	check := nativeRuntimeDBValidation(path)
	if check.Status != "ok" {
		t.Fatalf("stale metadata should warn without failing, got %#v", check)
	}
	if !strings.Contains(check.Message, "warning: runtime.db is 60.0 hours old") {
		t.Fatalf("expected stale warning, got %#v", check)
	}
}

func fakeRuntimeRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	for _, name := range []string{
		"generate-model-context-report.rb",
		"generate-model-checklists.rb",
		"generate-knowledge-runtime-report.rb",
		"generate-runtime-sqlite-index.rb",
		"refresh-knowledge-runtime.rb",
		"validate-knowledge-runtime.rb",
		"validate-runtime-db.rb",
		"validate-runtime-sqlite-index.rb",
	} {
		writeFile(t, filepath.Join(repo, "scripts", name), "#!/usr/bin/env ruby\nputs 'ok'\n")
	}
	writeFile(t, filepath.Join(repo, "runtime", "compiler", "compiler-engine.rb"), "#!/usr/bin/env ruby\nputs 'compiled'\n")
	copyFile(t, createNativeRuntimeDBFixture(t), filepath.Join(repo, "runtime", "runtime.db"))
	return repo
}

func createNativeRuntimeDBFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "runtime.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for _, table := range nativeRuntimeRequiredTables {
		if _, err := db.Exec(nativeRuntimeCreateTableSQL(table)); err != nil {
			t.Fatalf("create %s: %v", table, err)
		}
	}
	for table, minimum := range nativeRuntimeMinimumRows {
		for i := 0; i < minimum; i++ {
			if _, err := db.Exec(nativeRuntimeInsertSQL(table, i)); err != nil {
				t.Fatalf("insert %s: %v", table, err)
			}
		}
	}
	if _, err := db.Exec("INSERT OR REPLACE INTO compiler_metadata (key, value) VALUES ('compiler_version', 'test'), ('compiled_at', '2026-05-21T00:00:00Z')"); err != nil {
		t.Fatal(err)
	}
	return path
}

func nativeRuntimeCreateTableSQL(table string) string {
	if table == "compiler_metadata" {
		return "CREATE TABLE compiler_metadata (key TEXT PRIMARY KEY, value TEXT)"
	}
	columns := []string{"id TEXT PRIMARY KEY"}
	for _, column := range nativeRuntimeJSONColumns[table] {
		columns = append(columns, column+" TEXT")
	}
	return fmt.Sprintf("CREATE TABLE %s (%s)", table, strings.Join(columns, ", "))
}

func nativeRuntimeInsertSQL(table string, index int) string {
	if table == "compiler_metadata" {
		return fmt.Sprintf("INSERT OR REPLACE INTO compiler_metadata (key, value) VALUES ('key-%d', 'value-%d')", index, index)
	}
	columns := []string{"id"}
	values := []string{fmt.Sprintf("'%s-%d'", table, index)}
	for _, column := range nativeRuntimeJSONColumns[table] {
		columns = append(columns, column)
		values = append(values, "'[]'")
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(values, ", "))
}

func copyFile(t *testing.T, source string, target string) {
	t.Helper()
	content, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func containsEnv(env []string, item string) bool {
	for _, value := range env {
		if value == item {
			return true
		}
	}
	return false
}
