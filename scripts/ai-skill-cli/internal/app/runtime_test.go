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
	if len(result.PlannedActions) != 6 {
		t.Fatalf("expected six planned validators, got %#v", result.PlannedActions)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime validate dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeValidateDefaultNativeDoesNotNeedRuby(t *testing.T) {
	repo := repoRootForTest(t)
	ensureRuntimeIndexForRepoTest(t, repo)
	t.Setenv("PATH", emptyPathDir(t))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "validate", "--repo", repo, "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Mode != "native" {
		t.Fatalf("expected native mode, got %#v", result.Mode)
	}
	if !hasCheckStatus(result.Checks, "knowledge_runtime_native", "ok") {
		t.Fatalf("expected native knowledge runtime validation, got %#v", result.Checks)
	}
}

func TestNativeRoutingRegistryValidationRequiresSourceOfTruthGate(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "README.md"), "# Test Repo\n")
	registry := runtimeRoutingRegistry{Records: []runtimeRouteRecord{
		{
			ID:               "route.test.missing-gate",
			TaskIntent:       "test missing gate",
			PrimarySource:    "README.md",
			RankingReason:    "README is the test primary source.",
			ValidationSignal: "test validation signal",
			Metadata: runtimeRouteMetadata{
				Priority:           "P1",
				Confidence:         "high",
				CompatibilityState: "test-active",
			},
			Model: runtimeRouteModel{
				Profile:          "small",
				CompressionLevel: "summary-first",
			},
		},
	}}
	if err := nativeRoutingRegistryValidation(repo, registry); err == nil || !strings.Contains(err.Error(), "missing source_of_truth_gate") {
		t.Fatalf("expected missing source_of_truth_gate failure, got %v", err)
	}
}

func TestNativeRoutingRegistryValidationRequiresWorkflowActivation(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "workflow", "demo", "execution-flow.md"), "# Demo workflow\n")
	registry := runtimeRoutingRegistry{Records: []runtimeRouteRecord{
		{
			ID:                "route.workflow.demo",
			TaskIntent:        "demo workflow",
			PrimarySource:     "workflow/demo/execution-flow.md",
			SourceOfTruthGate: "demo-workflow-active",
			RankingReason:     "Workflow demo primary source.",
			ValidationSignal:  "Demo workflow route validated.",
			Metadata: runtimeRouteMetadata{
				Priority:           "P2",
				Confidence:         "high",
				CompatibilityState: "new-layer-promoted",
			},
			Model: runtimeRouteModel{
				Profile:          "large",
				CompressionLevel: "source-backed",
			},
		},
	}}
	if err := nativeRoutingRegistryValidation(repo, registry); err == nil || !strings.Contains(err.Error(), "workflow route missing activation_triggers") {
		t.Fatalf("expected missing workflow activation failure, got %v", err)
	}
}

func TestRuntimeRefreshDryRunPlansNativeActions(t *testing.T) {
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
		t.Fatalf("expected six planned native refresh actions, got %#v", result.PlannedActions)
	}
	if !strings.Contains(result.PlannedActions[0], "write native refresh report") {
		t.Fatalf("expected first planned action to write native report, got %#v", result.PlannedActions)
	}
	if len(result.Mutations) != 0 {
		t.Fatalf("runtime refresh dry-run must not mutate, got %#v", result.Mutations)
	}
}

func TestRuntimeRefreshDefaultNativeDoesNotNeedRuby(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	t.Setenv("PATH", emptyPathDir(t))
	writeRuntimeNativeReportSourceFixture(t, repo)
	writeFile(t, filepath.Join(repo, "README.md"), "# Test Repo\n")
	writeFile(t, filepath.Join(repo, "CORE_BOOTSTRAP.md"), "# Core\n")
	writeFile(t, filepath.Join(repo, "workflow", "test.md"), "# Workflow Test\n")
	writeFile(t, filepath.Join(repo, "knowledge", "summaries", "README.md"), "# knowledge summaries\n")
	writeFile(t, filepath.Join(repo, "skills", "demo", "feedback_history", "lesson.md"), "# Feedback\n\n### Feedback Lesson\n\n#### One-line Summary\nfeedback route lesson\n")

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
	if result.Mode != "native_refresh" {
		t.Fatalf("expected native_refresh mode, got %#v", result.Mode)
	}
	if hasCheckStatus(result.Checks, "ruby", "ok") || hasCheckStatus(result.Checks, "sqlite3", "ok") {
		t.Fatalf("native refresh should not check Ruby/sqlite3, got %#v", result.Checks)
	}
	for _, name := range []string{"knowledge_runtime_report", "model_context_report", "model_checklists", "runtime_sqlite_index", "runtime_index_native", "knowledge_runtime_native"} {
		if !hasCheckStatus(result.Checks, name, "ok") {
			t.Fatalf("expected ok check for %s, got %#v", name, result.Checks)
		}
	}
}

func TestRuntimeCompileDryRunPlansCompiler(t *testing.T) {
	repo := fakeRuntimeRepo(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "compile", "--repo", repo, "--dry-run", "--assert-source", "runtime/runtime.db", "--assert-keyword", "phase", "--json"}, &stdout, &stderr)
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
	if !hasCheckStatus(result.Checks, "runtime_compile_native", "ok") {
		t.Fatalf("expected runtime_compile_native ok, got %#v", result.Checks)
	}
}

func TestRuntimeCompileNativeCompilerWritesRuntimeDBWithoutRuby(t *testing.T) {
	repo := repoRootForTest(t)
	t.Setenv("PATH", emptyPathDir(t))
	outputDB := filepath.Join(t.TempDir(), "runtime-native.db")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"runtime", "compile", "--repo", repo, "--db", outputDB, "--json"}, &stdout, &stderr)
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
	if result.Results[0].ID != "phase-machine" || result.Results[0].SourcePath != "runtime/runtime.db" {
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

	runtimeReport, err := buildNativeKnowledgeRuntimeReport(repo)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(runtimeReport, "# Knowledge Runtime Report") || !strings.Contains(runtimeReport, "`route.bootstrap.ai-skill`") {
		t.Fatalf("runtime report missing golden anchors")
	}

	modelReport, err := buildNativeModelContextReport(repo)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(modelReport, "# Model Context Report") || !strings.Contains(modelReport, "## Profile View") {
		t.Fatalf("model context report missing golden anchors")
	}

	modelChecklists, err := buildNativeModelChecklists(repo)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(modelChecklists, "# Model Checklists") || !strings.Contains(modelChecklists, "## Profile Checklists") {
		t.Fatalf("model checklists missing golden anchors")
	}

	temp := t.TempDir()
	indexPath := filepath.Join(temp, "runtime-index.sqlite")
	if err := buildNativeRuntimeSQLiteIndex(repo, indexPath); err != nil {
		t.Fatal(err)
	}
	assertSQLiteCountAtLeast(t, indexPath, "atoms", 60)
	assertSQLiteCountAtLeast(t, indexPath, "sources", 50)
	assertSQLiteScalar(t, indexPath, "SELECT COUNT(*) FROM fts WHERE fts MATCH '\"runtime\"'", "nonzero")

	runtimeDBPath := filepath.Join(temp, "runtime.db")
	check := buildNativeRuntimeDBFromSources(repo, runtimeDBPath)
	if check.Status != "ok" {
		t.Fatalf("native compiler failed: %#v", check)
	}
	assertSQLiteCountAtLeast(t, runtimeDBPath, "generated_surfaces", 1)
	assertSQLiteScalar(t, runtimeDBPath, "SELECT COUNT(*) FROM decision_recording WHERE section = '__config__'", "nonzero")
	assertSQLiteScalar(t, runtimeDBPath, "SELECT COUNT(*) FROM generated_surfaces WHERE source_path = 'plans/active/*.md'", "nonzero")
	assertSQLiteScalar(t, runtimeDBPath, "SELECT COUNT(*) FROM compiler_metadata WHERE key = 'compiler_version'", "nonzero")
}

func TestRuntimeCompilerGoSnapshotHarnessIsStable(t *testing.T) {
	repo := repoRootForTest(t)
	temp := t.TempDir()
	firstDB := filepath.Join(temp, "runtime-a.db")
	secondDB := filepath.Join(temp, "runtime-b.db")

	if check := buildNativeRuntimeDBFromSources(repo, firstDB); check.Status != "ok" {
		t.Fatalf("native compiler first run failed: %#v", check)
	}
	if check := buildNativeRuntimeDBFromSources(repo, secondDB); check.Status != "ok" {
		t.Fatalf("native compiler second run failed: %#v", check)
	}

	if got, want := runtimeCompilerSnapshot(t, firstDB), runtimeCompilerSnapshot(t, secondDB); !reflect.DeepEqual(got, want) {
		t.Fatalf("Go compiler snapshots differ: %s", firstRowDiff(got, want))
	}
}

func TestNativeRuntimeCompilerBuildsFromSources(t *testing.T) {
	repo := repoRootForTest(t)
	temp := t.TempDir()
	nativeDB := filepath.Join(temp, "runtime-native.db")

	check := buildNativeRuntimeDBFromSources(repo, nativeDB)
	if check.Status != "ok" {
		t.Fatalf("native compiler failed: %#v", check)
	}

	assertSQLiteCountAtLeast(t, nativeDB, "generated_surfaces", 1)
	assertSQLiteCountAtLeast(t, nativeDB, "compiler_rules", 1)
	assertSQLiteCountAtLeast(t, nativeDB, "decision_recording", 1)
	assertSQLiteScalar(t, nativeDB, "SELECT COUNT(*) FROM runtime_budget WHERE model_name = 'default_budget'", "nonzero")
	assertSQLiteScalar(t, nativeDB, "SELECT COUNT(*) FROM runtime_budget WHERE model_name = 'layer:bootstrap'", "nonzero")
	assertSQLiteScalar(t, nativeDB, "SELECT COUNT(*) FROM runtime_budget WHERE model_name LIKE 'on_warning:%'", "nonzero")
	assertSQLiteScalar(t, nativeDB, "SELECT COUNT(*) FROM runtime_budget WHERE model_name LIKE 'on_hard_stop:%'", "nonzero")
	assertRuntimeCanonicalDocumentsProjected(t, nativeDB)
	assertSQLiteScalar(t, nativeDB, "SELECT COUNT(*) FROM compiler_metadata WHERE value = '2.0.0'", "nonzero")
	if check := nativeExecutableContractsValidation(repo, nativeDB); check.Status != "ok" {
		t.Fatalf("executable contract validation failed: %#v", check)
	}
}

func TestRuntimeValidateChecksCanonicalRuntimeDocuments(t *testing.T) {
	repo := fakeRuntimeRepo(t)
	db, err := sql.Open("sqlite", filepath.Join(repo, "runtime", "runtime.db"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("UPDATE runtime_config_documents SET content_json = 'not-json' WHERE logical_id = 'runtime-doc-0'"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	result := buildRuntimeValidateResult(runtimeOptions{repoPath: repo})
	if result.ExitCode != ExitValidationFailed {
		t.Fatalf("expected canonical runtime document failure, got %#v", result)
	}
	if result.Error == nil || result.Error.Code != "runtime_db_native_failed" {
		t.Fatalf("expected runtime_db_native_failed, got %#v", result.Error)
	}
}

func TestRuntimeCompilerCanonicalDocumentsAreProjected(t *testing.T) {
	repo := repoRootForTest(t)
	nativeDB := filepath.Join(t.TempDir(), "runtime-native.db")
	check := buildNativeRuntimeDBFromSources(repo, nativeDB)
	if check.Status != "ok" {
		t.Fatalf("native compiler failed: %#v", check)
	}
	assertRuntimeCanonicalDocumentsProjected(t, nativeDB)
}

func TestNativeRuntimeSQLiteIndexHasStableInvariants(t *testing.T) {
	repo := repoRootForTest(t)
	temp := t.TempDir()
	goPath := filepath.Join(temp, "go-runtime-index.sqlite")

	if err := buildNativeRuntimeSQLiteIndex(repo, goPath); err != nil {
		t.Fatal(err)
	}

	assertSQLiteCountAtLeast(t, goPath, "atoms", 60)
	assertSQLiteCountAtLeast(t, goPath, "sources", 50)
	assertSQLiteCountAtLeast(t, goPath, "edges", 1)
	for _, keyword := range []string{"runtime", "feedback", "route"} {
		query := "SELECT COUNT(*) FROM fts WHERE fts MATCH " + sqliteQuote(runtimeFTSMatchLiteral(keyword))
		if got := sqliteScalarInt(t, goPath, query); got == 0 {
			t.Fatalf("expected FTS hits for %q", keyword)
		}
	}
}

func TestNativeRuntimeSQLiteIndexIncludesRecursiveFeedback(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"), "records: []\n")
	writeFile(t, filepath.Join(repo, "skills", "demo", "feedback_history", "nested", "lesson.md"), `# Feedback Lesson

### Recursive Feedback Title

Status: promoted

#### One-line Summary
Recursive feedback summary.
`)
	temp := t.TempDir()
	goPath := filepath.Join(temp, "go-feedback.sqlite")

	if err := buildNativeRuntimeSQLiteIndex(repo, goPath); err != nil {
		t.Fatal(err)
	}
	assertSQLiteScalar(t, goPath, "SELECT COUNT(*) FROM atoms WHERE id = 'feedback.demo.lesson' AND source_path = 'skills/demo/feedback_history/nested/lesson.md'", "nonzero")
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
	writeFile(t, filepath.Join(repo, "runtime", "runtime.db"), "changed\n")

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
	copyFile(t, createNativeRuntimeDBFixture(t), filepath.Join(repo, "runtime", "runtime.db"))
	createRuntimeIndexFixture(t, filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite"))
	return repo
}

func writeRuntimeNativeReportSourceFixture(t *testing.T, repo string) {
	t.Helper()
	writeFile(t, filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"), `records:
  - id: route.test.small
    task_intent: test small route
    primary_source: README.md
    required_dependencies:
      - CORE_BOOTSTRAP.md
      - README.md
    source_of_truth_gate: test-small-active
    ranking_reason: README is the primary source for the small test route.
    validation_signal: small route validated
    metadata:
      priority: P1
      confidence: high
      compatibility_state: test-active
    model:
      profile: small
      compression_level: summary-first
      reason: small reason
  - id: route.test.large
    task_intent: test large route
    primary_source: workflow/test.md
    required_dependencies:
      - workflow/test.md
    source_of_truth_gate: test-large-active
    ranking_reason: workflow/test.md is the primary source for the large test route.
    validation_signal: large route validated
    metadata:
      priority: P2
      confidence: medium
      compatibility_state: test-active
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
	if table == "runtime_config_documents" {
		return "CREATE TABLE runtime_config_documents (logical_id TEXT PRIMARY KEY, content_json TEXT NOT NULL, checksum TEXT NOT NULL)"
	}
	if table == "runtime_config_projections" {
		return "CREATE TABLE runtime_config_projections (id TEXT PRIMARY KEY, logical_id TEXT, target_table TEXT, row_key TEXT, checksum TEXT)"
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
	if table == "runtime_config_documents" {
		return fmt.Sprintf("INSERT INTO runtime_config_documents (logical_id, content_json, checksum) VALUES ('runtime-doc-%d', '{}', 'checksum-%d')", index, index)
	}
	if table == "runtime_config_projections" {
		return fmt.Sprintf("INSERT INTO runtime_config_projections (id, logical_id, target_table, row_key, checksum) VALUES ('projection-%d', 'runtime-doc-%d', 'runtime_table', '__config__', 'checksum-%d')", index, index, index)
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
	workflowSource := "phase workflow guide route feedback\n"
	if _, err := os.Stat(filepath.Join(repo, "runtime", "runtime.db")); os.IsNotExist(err) {
		writeFile(t, filepath.Join(repo, "runtime", "runtime.db"), "runtime db fixture\n")
	} else if err != nil {
		t.Fatal(err)
	}
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
  ('phase-machine', 'runtime/runtime.db', 'runtime', 'reference', 'validated', 'P0', 'high', 'low', 'Phase machine runtime source.'),
  ('workflow-phase-guide', 'workflow/software-delivery/execution-flow.md', 'workflow', 'guide', 'candidate', 'P2', 'medium', 'medium', 'Workflow phase guide.');
INSERT INTO sources VALUES
  ('runtime/runtime.db', ?),
  ('workflow/software-delivery/execution-flow.md', ?);
INSERT INTO edges VALUES
  ('phase-machine', 'workflow-phase-guide', 'relates_to');
INSERT INTO fts VALUES
  ('phase-machine', 'phase phase phase machine runtime source feedback route'),
  ('workflow-phase-guide', 'phase workflow guide feedback route');
`, testChecksumFile(t, filepath.Join(repo, "runtime", "runtime.db")), testChecksum(workflowSource)); err != nil {
		t.Fatal(err)
	}
}

func testChecksum(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}

func testChecksumFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(content))
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
		if _, err := os.Stat(filepath.Join(dir, "CORE_BOOTSTRAP.md")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate Ai-skill repo root")
		}
		dir = parent
	}
}

func ensureRuntimeIndexForRepoTest(t *testing.T, repo string) {
	t.Helper()
	indexPath := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	if _, err := os.Stat(indexPath); err == nil {
		return
	} else if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if err := buildNativeRuntimeSQLiteIndex(repo, indexPath); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(indexPath)
	})
}

func assertRuntimeCanonicalDocumentsProjected(t *testing.T, dbPath string) {
	t.Helper()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	missing := []string{}
	for _, rel := range runtimeCanonicalDocumentPaths() {
		var docCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM runtime_config_documents WHERE logical_id = ?", rel).Scan(&docCount); err != nil {
			t.Fatal(err)
		}
		var projectionCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM runtime_config_projections WHERE logical_id = ?", rel).Scan(&projectionCount); err != nil {
			t.Fatal(err)
		}
		var sourceCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM runtime_source_files WHERE source_path = ? AND source_kind = 'db'", rel).Scan(&sourceCount); err != nil {
			t.Fatal(err)
		}
		if docCount == 0 || projectionCount == 0 || sourceCount == 0 {
			missing = append(missing, rel)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("runtime canonical documents missing projections: %s", strings.Join(missing, ", "))
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
