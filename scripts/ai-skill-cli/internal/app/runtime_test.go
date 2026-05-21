package app

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
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
	if len(result.PlannedActions) != 6 {
		t.Fatalf("expected six planned refresh steps, got %#v", result.PlannedActions)
	}
	if !strings.Contains(result.PlannedActions[0], "generate-model-context-report.rb --write") {
		t.Fatalf("expected first planned step to write model context report, got %#v", result.PlannedActions)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime refresh dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeRefreshExecutesOrderedSteps(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	requireExecutableForTest(t, "git")
	writeRuntimeRefreshRecorderScripts(t, repo, "")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	for _, step := range runtimeRefreshSteps(repo) {
		if !hasCheckStatus(result.Checks, step.name, "ok") {
			t.Fatalf("expected ok check for %s, got %#v", step.name, result.Checks)
		}
	}

	log := readTestFile(t, filepath.Join(repo, "refresh.log"))
	expected := strings.Join([]string{
		"model_context_report --write",
		"model_checklists --write",
		"knowledge_runtime_report --write",
		"runtime_sqlite_index",
		"runtime_sqlite_index_validation",
		"knowledge_runtime_validation",
	}, "\n") + "\n"
	if log != expected {
		t.Fatalf("unexpected refresh order:\n%s", log)
	}
}

func TestRuntimeRefreshStopsOnFirstFailedStep(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	requireExecutableForTest(t, "git")
	writeRuntimeRefreshRecorderScripts(t, repo, "knowledge_runtime_report")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Status != "blocked" || result.Error == nil || result.Error.Code != "runtime_refresh_failed" {
		t.Fatalf("expected runtime_refresh_failed block, got %#v", result)
	}
	if !hasCheckStatus(result.Checks, "knowledge_runtime_report", "failed") {
		t.Fatalf("expected failed check for knowledge_runtime_report, got %#v", result.Checks)
	}
	if hasCheckStatus(result.Checks, "runtime_sqlite_index", "ok") {
		t.Fatalf("runtime_sqlite_index should not run after failed report step: %#v", result.Checks)
	}

	log := readTestFile(t, filepath.Join(repo, "refresh.log"))
	expected := strings.Join([]string{
		"model_context_report --write",
		"model_checklists --write",
		"knowledge_runtime_report --write",
	}, "\n") + "\n"
	if log != expected {
		t.Fatalf("unexpected refresh order before failure:\n%s", log)
	}
}

func TestRuntimeRefreshNativeReportsWritesGoReportsThenRunsRemainingRubySteps(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	requireExecutableForTest(t, "git")
	writeRuntimeRefreshRecorderScripts(t, repo, "")
	writeRuntimeNativeReportSourceFixture(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--native-reports", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Mode != "wrapper_native_reports" {
		t.Fatalf("expected wrapper_native_reports mode, got %#v", result.Mode)
	}
	for _, name := range []string{"knowledge_runtime_report", "model_context_report", "model_checklists", "runtime_sqlite_index", "runtime_sqlite_index_validation", "knowledge_runtime_validation"} {
		if !hasCheckStatus(result.Checks, name, "ok") {
			t.Fatalf("expected ok check for %s, got %#v", name, result.Checks)
		}
	}
	if len(result.Mutations) != 3 {
		t.Fatalf("expected three native report mutations, got %#v", result.Mutations)
	}

	log := readTestFile(t, filepath.Join(repo, "refresh.log"))
	expected := strings.Join([]string{
		"runtime_sqlite_index",
		"runtime_sqlite_index_validation",
		"knowledge_runtime_validation",
	}, "\n") + "\n"
	if log != expected {
		t.Fatalf("unexpected remaining Ruby step order:\n%s", log)
	}

	expectedRuntimeReport, err := buildNativeKnowledgeRuntimeReport(repo)
	if err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, filepath.Join(repo, "knowledge", "runtime", "runtime-report.md")); got != expectedRuntimeReport {
		t.Fatalf("native runtime report content mismatch")
	}
}

func TestRuntimeRefreshNativeReportsAndIndexSkipsRubyGenerators(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	requireExecutableForTest(t, "git")
	writeRuntimeRefreshRecorderScripts(t, repo, "")
	writeRuntimeNativeReportSourceFixture(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "refresh", "--repo", repo, "--native-reports", "--native-index", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Mode != "wrapper_native_reports_index" {
		t.Fatalf("expected wrapper_native_reports_index mode, got %#v", result.Mode)
	}
	for _, name := range []string{"knowledge_runtime_report", "model_context_report", "model_checklists", "runtime_sqlite_index", "runtime_sqlite_index_validation", "knowledge_runtime_validation"} {
		if !hasCheckStatus(result.Checks, name, "ok") {
			t.Fatalf("expected ok check for %s, got %#v", name, result.Checks)
		}
	}
	if len(result.Mutations) != 4 {
		t.Fatalf("expected four native mutations, got %#v", result.Mutations)
	}

	log := readTestFile(t, filepath.Join(repo, "refresh.log"))
	expected := strings.Join([]string{
		"runtime_sqlite_index_validation",
		"knowledge_runtime_validation",
	}, "\n") + "\n"
	if log != expected {
		t.Fatalf("unexpected remaining Ruby step order:\n%s", log)
	}
	assertSQLiteCountAtLeast(t, filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite"), "atoms", 1)
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

func TestRuntimeCompileNativeCompilerWritesSnapshotWithoutRuby(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	t.Setenv("PATH", emptyPathDir(t))
	outputDB := filepath.Join(t.TempDir(), "runtime-native.db")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "compile", "--repo", repo, "--native-compiler", "--db", outputDB, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Mode != "native_compiler" {
		t.Fatalf("expected native_compiler mode, got %#v", result.Mode)
	}
	if !hasCheckStatus(result.Checks, "runtime_compile_native", "ok") {
		t.Fatalf("expected native compile check ok, got %#v", result.Checks)
	}
	if len(result.Mutations) != 1 || result.Mutations[0] != outputDB {
		t.Fatalf("expected output DB mutation, got %#v", result.Mutations)
	}
	assertSQLiteCountAtLeast(t, outputDB, "phases", 1)
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

func TestRuntimeQueryReturnsRankedResults(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "query", "--repo", repo, "--keyword", "phase", "--limit", "1", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "runtime query" || result.Mode != "native" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected one result, got %#v", result.Results)
	}
	if result.Results[0].ID != "phase-machine" || result.Results[0].SourcePath != "runtime/compiler/embedded_data.rb" {
		t.Fatalf("unexpected top result: %#v", result.Results[0])
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime query must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeQueryFiltersByLayerTypeStatus(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "query", "--keyword", "phase", "--repo", repo, "--layer", "workflow", "--type", "guide", "--status", "candidate", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if len(result.Results) != 1 || result.Results[0].ID != "workflow-phase-guide" {
		t.Fatalf("expected filtered workflow result, got %#v", result.Results)
	}
}

func TestRuntimeQueryEmptyResultSucceeds(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "query", "--repo", repo, "--keyword", "not-present", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if len(result.Results) != 0 {
		t.Fatalf("expected empty results, got %#v", result.Results)
	}
}

func TestRuntimeQueryBlocksMissingIndex(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	if err := os.Remove(filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "query", "--repo", repo, "--keyword", "phase", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_runtime_index" {
		t.Fatalf("expected missing runtime index, got %#v", result.Error)
	}
}

func TestRuntimeGraphQueryFiltersEdges(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	createKnowledgeGraphFixture(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "query", "--graph", "--repo", repo, "--source", "workflow/software-delivery", "--target", "artifact-gates", "--type", "related_to", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "runtime query" || result.Mode != "native" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected one graph result, got %#v", result.Results)
	}
	got := result.Results[0]
	if got.GraphID != "graph.test" || got.EdgeType != "related_to" || got.Target != "workflow/software-delivery/artifact-gates.md" {
		t.Fatalf("unexpected graph result: %#v", got)
	}
	if got.GraphFile != "knowledge/graphs/test-graph.yaml" {
		t.Fatalf("expected graph file path, got %#v", got)
	}
}

func TestRuntimeGraphQueryEmptyResultSucceeds(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	createKnowledgeGraphFixture(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "query", "--graph", "--repo", repo, "--keyword", "not-present", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if len(result.Results) != 0 {
		t.Fatalf("expected empty graph results, got %#v", result.Results)
	}
}

func TestRuntimeGraphQueryRequiresFilter(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	createKnowledgeGraphFixture(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "query", "--graph", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitInvalidUsage {
		t.Fatalf("expected invalid usage, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_graph_filter" {
		t.Fatalf("expected missing graph filter, got %#v", result.Error)
	}
}

func TestRuntimeGoldenFixtureCoversGeneratedSurfaces(t *testing.T) {
	repo := repoRootForTest(t)
	ruby := requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")

	runtimeReport := runRubyScript(t, repo, ruby, "scripts/generate-knowledge-runtime-report.rb")
	if !strings.Contains(runtimeReport, "# Knowledge Runtime Report") || !strings.Contains(runtimeReport, "`route.bootstrap.ai-skill`") {
		t.Fatalf("runtime report missing golden anchors")
	}

	modelReport := runRubyScript(t, repo, ruby, "scripts/generate-model-context-report.rb")
	if !strings.Contains(modelReport, "# Model Context Report") || !strings.Contains(modelReport, "## Profile View") {
		t.Fatalf("model context report missing golden anchors")
	}

	modelChecklists := runRubyScript(t, repo, ruby, "scripts/generate-model-checklists.rb")
	if !strings.Contains(modelChecklists, "# Model Checklists") || !strings.Contains(modelChecklists, "## Profile Checklists") {
		t.Fatalf("model checklists missing golden anchors")
	}

	temp := t.TempDir()
	indexPath := filepath.Join(temp, "runtime-index.sqlite")
	indexRel, err := filepath.Rel(repo, indexPath)
	if err != nil {
		t.Fatal(err)
	}
	runRubyScript(t, repo, ruby, "scripts/generate-runtime-sqlite-index.rb", "--output", filepath.ToSlash(indexRel))
	assertSQLiteCountAtLeast(t, indexPath, "atoms", 60)
	assertSQLiteCountAtLeast(t, indexPath, "sources", 50)
	assertSQLiteScalar(t, indexPath, "SELECT COUNT(*) FROM fts WHERE fts MATCH '\"runtime\"'", "nonzero")

	runtimeDBPath := filepath.Join(temp, "runtime.db")
	runRubyScript(t, repo, ruby, "runtime/compiler/compiler-engine.rb", "--db", runtimeDBPath)
	assertSQLiteCountAtLeast(t, runtimeDBPath, "generated_surfaces", 1)
	assertSQLiteScalar(t, runtimeDBPath, "SELECT COUNT(*) FROM generated_surfaces WHERE source_path = 'plans/active/*.md'", "nonzero")
	assertSQLiteScalar(t, runtimeDBPath, "SELECT COUNT(*) FROM compiler_metadata WHERE key = 'compiler_version'", "nonzero")
}

func TestRuntimeCompilerRubySnapshotHarnessIsStable(t *testing.T) {
	repo := repoRootForTest(t)
	ruby := requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	temp := t.TempDir()
	firstDB := filepath.Join(temp, "runtime-a.db")
	secondDB := filepath.Join(temp, "runtime-b.db")

	runRubyScript(t, repo, ruby, "runtime/compiler/compiler-engine.rb", "--db", firstDB)
	runRubyScript(t, repo, ruby, "runtime/compiler/compiler-engine.rb", "--db", secondDB)

	if got, want := runtimeCompilerSnapshot(t, firstDB), runtimeCompilerSnapshot(t, secondDB); !reflect.DeepEqual(got, want) {
		t.Fatalf("Ruby compiler snapshots differ: %s", firstRowDiff(got, want))
	}
}

func TestNativeRuntimeCompilerSnapshotMatchesRubyCompiler(t *testing.T) {
	repo := repoRootForTest(t)
	ruby := requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	temp := t.TempDir()
	rubyDB := filepath.Join(temp, "runtime-ruby.db")
	nativeDB := filepath.Join(temp, "runtime-native.db")

	runRubyScript(t, repo, ruby, "runtime/compiler/compiler-engine.rb", "--db", rubyDB)
	check := buildNativeRuntimeDBSnapshot(repo, nativeDB)
	if check.Status != "ok" {
		t.Fatalf("native compiler snapshot failed: %#v", check)
	}

	if got, want := runtimeCompilerGeneratedSurfaceSnapshot(t, nativeDB), runtimeCompilerGeneratedSurfaceSnapshot(t, rubyDB); !reflect.DeepEqual(got, want) {
		t.Fatalf("native compiler snapshot differs from Ruby: %s", firstRowDiff(got, want))
	}
}

func TestNativeModelContextReportMatchesRubyGenerator(t *testing.T) {
	repo := repoRootForTest(t)
	ruby := requireExecutableForTest(t, "ruby")

	rubyOutput := runRubyScriptStdout(t, repo, ruby, "scripts/generate-model-context-report.rb")
	goOutput, err := buildNativeModelContextReport(repo)
	if err != nil {
		t.Fatal(err)
	}
	if goOutput != rubyOutput {
		t.Fatalf("Go model context report does not match Ruby output: %s", firstStringDiff(goOutput, rubyOutput))
	}
}

func TestNativeModelChecklistsMatchesRubyGenerator(t *testing.T) {
	repo := repoRootForTest(t)
	ruby := requireExecutableForTest(t, "ruby")

	rubyOutput := runRubyScriptStdout(t, repo, ruby, "scripts/generate-model-checklists.rb")
	goOutput, err := buildNativeModelChecklists(repo)
	if err != nil {
		t.Fatal(err)
	}
	if goOutput != rubyOutput {
		t.Fatalf("Go model checklists report does not match Ruby output: %s", firstStringDiff(goOutput, rubyOutput))
	}
}

func TestNativeKnowledgeRuntimeReportMatchesRubyGenerator(t *testing.T) {
	repo := repoRootForTest(t)
	ruby := requireExecutableForTest(t, "ruby")

	rubyOutput := runRubyScriptStdout(t, repo, ruby, "scripts/generate-knowledge-runtime-report.rb")
	goOutput, err := buildNativeKnowledgeRuntimeReport(repo)
	if err != nil {
		t.Fatal(err)
	}
	if goOutput != rubyOutput {
		t.Fatalf("Go knowledge runtime report does not match Ruby output: %s", firstStringDiff(goOutput, rubyOutput))
	}
}

func TestNativeRuntimeSQLiteIndexMatchesRubyInvariants(t *testing.T) {
	repo := repoRootForTest(t)
	ruby := requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	temp := t.TempDir()
	rubyPath := filepath.Join(temp, "ruby-runtime-index.sqlite")
	goPath := filepath.Join(temp, "go-runtime-index.sqlite")
	rubyRel, err := filepath.Rel(repo, rubyPath)
	if err != nil {
		t.Fatal(err)
	}

	runRubyScript(t, repo, ruby, "scripts/generate-runtime-sqlite-index.rb", "--output", filepath.ToSlash(rubyRel))
	if err := buildNativeRuntimeSQLiteIndex(repo, goPath); err != nil {
		t.Fatal(err)
	}

	for _, table := range []string{"atoms", "sources", "edges", "fts"} {
		rubyCount := sqliteCount(t, rubyPath, table)
		goCount := sqliteCount(t, goPath, table)
		if goCount != rubyCount {
			t.Fatalf("%s count mismatch: go=%d ruby=%d", table, goCount, rubyCount)
		}
	}
	if got, want := sqliteSourceChecksums(t, goPath), sqliteSourceChecksums(t, rubyPath); !reflect.DeepEqual(got, want) {
		t.Fatalf("source checksums mismatch")
	}
	for _, table := range []string{"atoms", "sources", "edges", "fts"} {
		if got, want := sqliteRows(t, goPath, table), sqliteRows(t, rubyPath, table); !reflect.DeepEqual(got, want) {
			t.Fatalf("%s row-level mismatch: %s", table, firstRowDiff(got, want))
		}
	}
	for _, keyword := range []string{"runtime", "feedback", "route"} {
		query := "SELECT COUNT(*) FROM fts WHERE fts MATCH " + sqliteQuote(runtimeFTSMatchLiteral(keyword))
		if got, want := sqliteScalarInt(t, goPath, query), sqliteScalarInt(t, rubyPath, query); got != want {
			t.Fatalf("FTS hit mismatch for %q: go=%d ruby=%d", keyword, got, want)
		}
	}
}

func TestNativeRuntimeSQLiteIndexMatchesRubyRecursiveFeedback(t *testing.T) {
	ruby := requireExecutableForTest(t, "ruby")
	requireExecutableForTest(t, "sqlite3")
	sourceRepo := repoRootForTest(t)
	repo := t.TempDir()
	copyFile(t, filepath.Join(sourceRepo, "scripts", "generate-runtime-sqlite-index.rb"), filepath.Join(repo, "scripts", "generate-runtime-sqlite-index.rb"))
	writeFile(t, filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"), "records: []\n")
	writeFile(t, filepath.Join(repo, "skills", "demo", "feedback_history", "nested", "lesson.md"), `# Feedback Lesson

### Recursive Feedback Title

Status: promoted

#### One-line Summary
Recursive feedback summary.
`)
	temp := t.TempDir()
	rubyPath := filepath.Join(temp, "ruby-feedback.sqlite")
	goPath := filepath.Join(temp, "go-feedback.sqlite")
	rubyRel, err := filepath.Rel(repo, rubyPath)
	if err != nil {
		t.Fatal(err)
	}

	runRubyScript(t, repo, ruby, "scripts/generate-runtime-sqlite-index.rb", "--output", filepath.ToSlash(rubyRel))
	if err := buildNativeRuntimeSQLiteIndex(repo, goPath); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"atoms", "sources", "edges", "fts"} {
		if got, want := sqliteRows(t, goPath, table), sqliteRows(t, rubyPath, table); !reflect.DeepEqual(got, want) {
			t.Fatalf("%s recursive feedback rows mismatch: %s", table, firstRowDiff(got, want))
		}
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

func TestNativeRuntimeIndexValidationAcceptsValidFixture(t *testing.T) {
	repo := t.TempDir()
	path := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	createRuntimeIndexFixture(t, path)

	check := nativeRuntimeIndexValidation(repo, path)
	if check.Status != "ok" {
		t.Fatalf("expected native index validation ok, got %#v", check)
	}
}

func TestNativeRuntimeIndexValidationBlocksMissingTable(t *testing.T) {
	repo := t.TempDir()
	path := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	createRuntimeIndexFixture(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("DROP TABLE edges"); err != nil {
		t.Fatal(err)
	}

	check := nativeRuntimeIndexValidation(repo, path)
	if check.Status != "failed" || !strings.Contains(check.Message, "missing table: edges") {
		t.Fatalf("expected missing table failure, got %#v", check)
	}
}

func TestNativeRuntimeIndexValidationBlocksStaleChecksum(t *testing.T) {
	repo := t.TempDir()
	path := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	createRuntimeIndexFixture(t, path)
	writeFile(t, filepath.Join(repo, "runtime", "compiler", "embedded_data.rb"), "changed\n")

	check := nativeRuntimeIndexValidation(repo, path)
	if check.Status != "failed" || !strings.Contains(check.Message, "stale checksum") {
		t.Fatalf("expected stale checksum failure, got %#v", check)
	}
}

func TestNativeRuntimeIndexValidationBlocksFTSCountMismatch(t *testing.T) {
	repo := t.TempDir()
	path := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	createRuntimeIndexFixture(t, path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("DELETE FROM fts WHERE id = 'workflow-phase-guide'"); err != nil {
		t.Fatal(err)
	}

	check := nativeRuntimeIndexValidation(repo, path)
	if check.Status != "failed" || !strings.Contains(check.Message, "fts count does not match atoms count") {
		t.Fatalf("expected FTS count failure, got %#v", check)
	}
}

func TestNativeRuntimeIndexGitIgnoreCheckAcceptsIgnoredIndex(t *testing.T) {
	repo := initTempGitRepo(t)
	writeFile(t, filepath.Join(repo, ".gitignore"), "knowledge/runtime/sqlite/\n")
	runGit(t, repo, "add", ".gitignore")
	runGit(t, repo, "commit", "-m", "ignore runtime index")
	path := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	createRuntimeIndexFixture(t, path)

	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is required for git-ignore boundary test")
	}
	check := nativeRuntimeIndexGitIgnoreCheck(repo, path, git)
	if check.Status != "ok" {
		t.Fatalf("expected git-ignore check ok, got %#v", check)
	}
}

func TestNativeRuntimeIndexGitIgnoreCheckBlocksTrackedBoundary(t *testing.T) {
	repo := initTempGitRepo(t)
	path := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	createRuntimeIndexFixture(t, path)

	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is required for git-ignore boundary test")
	}
	check := nativeRuntimeIndexGitIgnoreCheck(repo, path, git)
	if check.Status != "failed" || !strings.Contains(check.Message, "generated DB is not ignored by git") {
		t.Fatalf("expected git-ignore boundary failure, got %#v", check)
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
	createRuntimeIndexFixture(t, filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite"))
	return repo
}

func writeRuntimeRefreshRecorderScripts(t *testing.T, repo string, failStep string) {
	t.Helper()
	for _, step := range runtimeRefreshSteps(repo) {
		failLine := ""
		if step.name == failStep {
			failLine = "exit 7\n"
		}
		writeFile(t, step.path, fmt.Sprintf(`#!/usr/bin/env ruby
line = ([%q] + ARGV).join(" ").strip
File.open(File.join(Dir.pwd, "refresh.log"), "a") { |file| file.puts(line) }
puts "#{line} ok"
%s`, step.name, failLine))
	}
}

func writeRuntimeNativeReportSourceFixture(t *testing.T, repo string) {
	t.Helper()
	writeFile(t, filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"), `records:
  - id: route.test.small
    primary_source: README.md
    required_dependencies:
      - CORE_BOOTSTRAP.md
      - README.md
    validation_signal: small route validated
    model:
      profile: small
      compression_level: summary-first
      reason: small reason
  - id: route.test.large
    primary_source: workflow/test.md
    required_dependencies:
      - workflow/test.md
    validation_signal: large route validated
    model:
      profile: large
      compression_level: source-backed
      reason: large reason
`)
	writeFile(t, filepath.Join(repo, "knowledge", "runtime", "refresh-policy.yaml"), `status: candidate
decision_values:
  - refresh_now
  - no_update_needed
`)
	writeFile(t, filepath.Join(repo, "knowledge", "summaries", "test-summary.md"), "# test.summary\n\n| 欄位 | 值 |\n| --- | --- |\n| Atom ID | `test.summary` |\n| Lifecycle | `validated` |\n| Summary | Test summary text. |\n")
	writeFile(t, filepath.Join(repo, "knowledge", "graphs", "test-graph.yaml"), `id: graph.test
source: README.md
status: candidate
edges:
  - type: depends_on
    target: CORE_BOOTSTRAP.md
`)
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func firstStringDiff(got string, want string) string {
	limit := len(got)
	if len(want) < limit {
		limit = len(want)
	}
	for i := 0; i < limit; i++ {
		if got[i] != want[i] {
			return fmt.Sprintf("first diff at byte %d: got %q want %q", i, diffWindow(got, i), diffWindow(want, i))
		}
	}
	return fmt.Sprintf("length mismatch: got %d bytes, want %d bytes", len(got), len(want))
}

func diffWindow(value string, index int) string {
	start := index - 80
	if start < 0 {
		start = 0
	}
	end := index + 80
	if end > len(value) {
		end = len(value)
	}
	return value[start:end]
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

func createRuntimeIndexFixture(t *testing.T, path string) {
	t.Helper()
	repo := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(path))))
	phaseSource := "phase machine runtime source feedback route\n"
	workflowSource := "phase workflow guide route feedback\n"
	writeFile(t, filepath.Join(repo, "runtime", "compiler", "embedded_data.rb"), phaseSource)
	writeFile(t, filepath.Join(repo, "workflow", "software-delivery", "execution-flow.md"), workflowSource)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`
CREATE TABLE atoms (
  id TEXT PRIMARY KEY,
  source_path TEXT,
  layer TEXT,
  type TEXT,
  status TEXT,
  priority TEXT,
  confidence TEXT,
  context_cost TEXT,
  summary TEXT
);
CREATE TABLE sources (
  source_path TEXT PRIMARY KEY,
  checksum TEXT
);
CREATE TABLE edges (
  source_id TEXT,
  target_id TEXT,
  type TEXT
);
CREATE VIRTUAL TABLE fts USING fts5(id UNINDEXED, content);
INSERT INTO atoms VALUES
  ('phase-machine', 'runtime/compiler/embedded_data.rb', 'runtime', 'reference', 'validated', 'P0', 'high', 'low', 'Phase machine runtime source.'),
  ('workflow-phase-guide', 'workflow/software-delivery/execution-flow.md', 'workflow', 'guide', 'candidate', 'P2', 'medium', 'medium', 'Workflow phase guide.');
INSERT INTO sources VALUES
  ('runtime/compiler/embedded_data.rb', ?),
  ('workflow/software-delivery/execution-flow.md', ?);
INSERT INTO edges VALUES
  ('phase-machine', 'workflow-phase-guide', 'relates_to');
INSERT INTO fts VALUES
  ('phase-machine', 'phase phase phase machine runtime source feedback route'),
  ('workflow-phase-guide', 'phase workflow guide feedback route');
`, testChecksum(phaseSource), testChecksum(workflowSource)); err != nil {
		t.Fatal(err)
	}
}

func testChecksum(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}

func createKnowledgeGraphFixture(t *testing.T, repo string) {
	t.Helper()
	writeFile(t, filepath.Join(repo, "knowledge", "graphs", "test-graph.yaml"), `id: graph.test
source: workflow/software-delivery/README.md
status: candidate
edges:
  - type: related_to
    target: workflow/software-delivery/artifact-gates.md
    reason: Artifact gates define delivery outputs.
    validation: Fixture validates graph query filters.
  - type: depends_on
    target: analysis/development-guidance/README.md
    reason: Analysis guidance supports workflow decisions.
    validation: Fixture validates empty query behavior.
`)
}

func repoRootForTest(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "scripts", "generate-knowledge-runtime-report.rb")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate Ai-skill repo root")
		}
		dir = parent
	}
}

func requireExecutableForTest(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s is required for golden fixture integration test", name)
	}
	return path
}

func runRubyScript(t *testing.T, repo string, ruby string, script string, args ...string) string {
	t.Helper()
	cmd := exec.Command(ruby, append([]string{script}, args...)...)
	cmd.Dir = repo
	cmd.Env = runtimeWrapperEnv(os.Environ())
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ruby %s failed: %v\n%s", script, err, string(output))
	}
	return string(output)
}

func runRubyScriptStdout(t *testing.T, repo string, ruby string, script string, args ...string) string {
	t.Helper()
	cmd := exec.Command(ruby, append([]string{script}, args...)...)
	cmd.Dir = repo
	cmd.Env = runtimeWrapperEnv(os.Environ())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("ruby %s failed: %v\n%s", script, err, stderr.String())
	}
	return string(output)
}

func assertSQLiteCountAtLeast(t *testing.T, path string, table string, minimum int) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count < minimum {
		t.Fatalf("%s row count = %d, expected at least %d", table, count, minimum)
	}
}

func assertSQLiteScalar(t *testing.T, path string, query string, expectation string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if expectation == "nonzero" && count == 0 {
		t.Fatalf("query returned zero: %s", query)
	}
}

func sqliteCount(t *testing.T, path string, table string) int {
	t.Helper()
	return sqliteScalarInt(t, path, "SELECT COUNT(*) FROM "+table)
}

func sqliteScalarInt(t *testing.T, path string, query string) int {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		t.Fatal(err)
	}
	return count
}

func sqliteSourceChecksums(t *testing.T, path string) map[string]string {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT source_path, checksum FROM sources ORDER BY source_path")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	checksums := map[string]string{}
	for rows.Next() {
		var sourcePath string
		var checksum string
		if err := rows.Scan(&sourcePath, &checksum); err != nil {
			t.Fatal(err)
		}
		checksums[sourcePath] = checksum
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return checksums
}

func sqliteRows(t *testing.T, path string, table string) []string {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	result := []string{}
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		scan := make([]any, len(columns))
		for index := range values {
			scan[index] = &values[index]
		}
		if err := rows.Scan(scan...); err != nil {
			t.Fatal(err)
		}
		cells := make([]string, len(values))
		for index, value := range values {
			if value.Valid {
				cells[index] = value.String
			} else {
				cells[index] = "<NULL>"
			}
		}
		result = append(result, strings.Join(cells, "\x1f"))
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	sort.Strings(result)
	return result
}

func runtimeCompilerSnapshot(t *testing.T, path string) []string {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	snapshot := []string{}
	for _, table := range nativeRuntimeRequiredTables {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		snapshot = append(snapshot, fmt.Sprintf("count:%s=%d", table, count))
	}
	rows, err := db.Query(`SELECT source_path, target_key, compile_rule, status, data FROM generated_surfaces ORDER BY source_path, target_key`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var sourcePath string
		var targetKey string
		var compileRule string
		var status string
		var data string
		if err := rows.Scan(&sourcePath, &targetKey, &compileRule, &status, &data); err != nil {
			t.Fatal(err)
		}
		snapshot = append(snapshot, strings.Join([]string{"surface", sourcePath, targetKey, compileRule, status, data}, "\x1f"))
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"compiler_version", "schema_version"} {
		var value string
		if err := db.QueryRow("SELECT value FROM compiler_metadata WHERE key = ?", key).Scan(&value); err != nil {
			t.Fatalf("metadata %s: %v", key, err)
		}
		snapshot = append(snapshot, "metadata:"+key+"="+value)
	}
	sort.Strings(snapshot)
	return snapshot
}

func runtimeCompilerGeneratedSurfaceSnapshot(t *testing.T, path string) []string {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	snapshot := []string{}
	for _, table := range []string{"generated_surfaces", "compiler_metadata"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		snapshot = append(snapshot, fmt.Sprintf("count:%s=%d", table, count))
	}
	rows, err := db.Query(`SELECT source_path, target_key, compile_rule, status, data FROM generated_surfaces ORDER BY source_path, target_key`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var sourcePath string
		var targetKey string
		var compileRule string
		var status string
		var data string
		if err := rows.Scan(&sourcePath, &targetKey, &compileRule, &status, &data); err != nil {
			t.Fatal(err)
		}
		snapshot = append(snapshot, strings.Join([]string{"surface", sourcePath, targetKey, compileRule, status, data}, "\x1f"))
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"compiler_version", "schema_version"} {
		var value string
		if err := db.QueryRow("SELECT value FROM compiler_metadata WHERE key = ?", key).Scan(&value); err != nil {
			t.Fatalf("metadata %s: %v", key, err)
		}
		snapshot = append(snapshot, "metadata:"+key+"="+value)
	}
	sort.Strings(snapshot)
	return snapshot
}

func firstRowDiff(got []string, want []string) string {
	limit := len(got)
	if len(want) < limit {
		limit = len(want)
	}
	for index := 0; index < limit; index++ {
		if got[index] != want[index] {
			return fmt.Sprintf("first row diff at %d: got=%q want=%q", index, got[index], want[index])
		}
	}
	return fmt.Sprintf("row count mismatch: got=%d want=%d", len(got), len(want))
}

func sqliteQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func containsEnv(env []string, item string) bool {
	for _, value := range env {
		if value == item {
			return true
		}
	}
	return false
}
