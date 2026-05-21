package app

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type runtimeOptions struct {
	command        string
	repoPath       string
	dryRun         bool
	nativeReports  bool
	nativeIndex    bool
	nativeCompiler bool
	assertSource   string
	assertKeyword  string
	keyword        string
	dbPath         string
	layer          string
	queryType      string
	statusFilter   string
	source         string
	target         string
	graphQuery     bool
	limit          int
	jsonOutput     bool
	plainOutput    bool
	positionals    []string
}

type runtimeNativeReportTarget struct {
	name  string
	path  string
	build func(string) (string, error)
}

type runtimeRoutingRegistry struct {
	Records []runtimeRouteRecord `yaml:"records"`
}

type runtimeRouteRecord struct {
	ID                   string               `yaml:"id"`
	TaskIntent           string               `yaml:"task_intent"`
	PrimarySource        string               `yaml:"primary_source"`
	RequiredDependencies []string             `yaml:"required_dependencies"`
	RankingReason        string               `yaml:"ranking_reason"`
	ValidationSignal     string               `yaml:"validation_signal"`
	Metadata             runtimeRouteMetadata `yaml:"metadata"`
	Model                runtimeRouteModel    `yaml:"model"`
}

type runtimeRouteMetadata struct {
	Priority           string `yaml:"priority"`
	Confidence         string `yaml:"confidence"`
	ContextCost        any    `yaml:"context_cost"`
	CompatibilityState string `yaml:"compatibility_state"`
}

type runtimeRouteModel struct {
	Profile          string `yaml:"profile"`
	CompressionLevel string `yaml:"compression_level"`
	Reason           string `yaml:"reason"`
}

type runtimeRefreshPolicy struct {
	Status         string   `yaml:"status"`
	DecisionValues []string `yaml:"decision_values"`
}

type runtimeSummaryRecord struct {
	File      string
	AtomID    string
	Lifecycle string
	Summary   string
}

type runtimeGraphRecord struct {
	File      string
	ID        string
	Source    string
	Status    string
	EdgeCount int
}

type runtimeIndexAtom struct {
	ID               string
	SourcePath       string
	Layer            string
	Type             string
	Status           string
	Priority         string
	Confidence       string
	ContextCost      string
	Tags             string
	Domains          string
	Title            string
	Summary          string
	WhenToRead       string
	ValidationSignal string
}

type runtimeIndexEdge struct {
	GraphID    string
	SourcePath string
	EdgeType   string
	TargetPath string
	Reason     string
	Validation string
}

func runRuntime(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill runtime <validate|refresh|compile|query> [flags]")
		return ExitInvalidUsage
	}
	opts := runtimeOptions{command: args[0], limit: 8}
	if opts.command != "validate" && opts.command != "refresh" && opts.command != "compile" && opts.command != "query" {
		_, _ = fmt.Fprintf(stderr, "unsupported runtime command: %s\n", opts.command)
		return ExitInvalidUsage
	}

	fs := newFlagSet("runtime "+opts.command, stderr)
	fs.StringVar(&opts.repoPath, "repo", ".", "Ai-skill repository path")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview runtime wrapper scripts without executing")
	fs.BoolVar(&opts.nativeReports, "native-reports", false, "opt in to Go-native Markdown report generation during runtime refresh")
	fs.BoolVar(&opts.nativeIndex, "native-index", false, "opt in to Go-native SQLite index generation during runtime refresh")
	fs.BoolVar(&opts.nativeCompiler, "native-compiler", false, "deprecated no-op; runtime compile is Go-native")
	fs.StringVar(&opts.assertSource, "assert-source", "", "source path expected in generated surfaces")
	fs.StringVar(&opts.assertKeyword, "assert-keyword", "", "keyword expected in generated surfaces")
	fs.StringVar(&opts.keyword, "keyword", "", "keyword or phrase to query in the runtime index")
	fs.StringVar(&opts.dbPath, "db", "", "SQLite runtime index path")
	fs.StringVar(&opts.layer, "layer", "", "filter runtime query results by layer")
	fs.StringVar(&opts.queryType, "type", "", "filter runtime query results by type")
	fs.StringVar(&opts.statusFilter, "status", "", "filter runtime query results by status")
	fs.StringVar(&opts.source, "source", "", "filter runtime graph query results by source")
	fs.StringVar(&opts.target, "target", "", "filter runtime graph query results by target")
	fs.BoolVar(&opts.graphQuery, "graph", false, "query knowledge graph YAML edges instead of the SQLite runtime index")
	fs.IntVar(&opts.limit, "limit", 8, "maximum runtime query results")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	opts.positionals = fs.Args()
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
	case "query":
		return buildRuntimeQueryResult(opts)
	case "refresh":
		return buildRuntimeRefreshResult(opts)
	default:
		return buildRuntimeValidateResult(opts)
	}
}

func buildRuntimeValidateResult(opts runtimeOptions) Result {
	result := Result{
		Command:        "runtime validate",
		Mode:           "native",
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

	result.PlannedActions = append(result.PlannedActions, "run native runtime DB validation")
	result.PlannedActions = append(result.PlannedActions, "run native runtime SQLite index validation")
	result.PlannedActions = append(result.PlannedActions, "run native knowledge runtime validation")

	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "native_mode", Status: "ok", Message: "dry-run only; validators not executed"})
		return result
	}

	nativeDBCheck := nativeRuntimeDBValidation(filepath.Join(repo, "runtime", "runtime.db"))
	result.Checks = append(result.Checks, nativeDBCheck)
	if nativeDBCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_db_native_failed", Message: nativeDBCheck.Message, Remediation: "Fix runtime source files or run runtime compile."}
		return result
	}

	nativeIndexCheck := nativeRuntimeIndexValidation(repo, filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite"))
	result.Checks = append(result.Checks, nativeIndexCheck)
	if nativeIndexCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_index_native_failed", Message: nativeIndexCheck.Message, Remediation: "Fix knowledge/runtime/sqlite/runtime-index.sqlite or run runtime refresh."}
		return result
	}

	knowledgeRuntimeCheck := nativeKnowledgeRuntimeValidation(repo)
	result.Checks = append(result.Checks, knowledgeRuntimeCheck)
	if knowledgeRuntimeCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "knowledge_runtime_native_failed", Message: knowledgeRuntimeCheck.Message, Remediation: "Fix knowledge runtime sources or generated reports."}
		return result
	}

	return result
}

func buildRuntimeRefreshResult(opts runtimeOptions) Result {
	opts.nativeReports = true
	opts.nativeIndex = true
	result := Result{
		Command:        "runtime refresh",
		Mode:           "native_refresh",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	if opts.dryRun {
		result.Mode = "dry_run"
	}
	if (opts.nativeReports || opts.nativeIndex) && !opts.dryRun {
		result.Mode = runtimeRefreshMode(opts)
	}

	repo, repoCheck := resolveExistingDir("repo", opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message, Remediation: "Pass --repo with the Ai-skill repository root."}
		return result
	}
	nativeTargets := []runtimeNativeReportTarget{}
	nativeTargets = runtimeNativeReportTargets(repo)
	for _, target := range nativeTargets {
		result.PlannedActions = append(result.PlannedActions, "write native refresh report: "+target.path)
	}
	nativeIndexPath := ""
	nativeIndexPath = filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	result.PlannedActions = append(result.PlannedActions, "write native runtime SQLite index: "+nativeIndexPath)
	result.PlannedActions = append(result.PlannedActions, "run native runtime SQLite index validation")
	result.PlannedActions = append(result.PlannedActions, "run native knowledge runtime validation")

	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "refresh_mode", Status: "ok", Message: "dry-run only; refresh steps not executed"})
		return result
	}

	for _, target := range nativeTargets {
		check := writeNativeRuntimeReport(repo, target)
		result.Checks = append(result.Checks, check)
		if check.Status != "ok" {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "runtime_refresh_failed", Message: "native report failed: " + target.name + ": " + check.Message, Remediation: "Inspect native report output and fix the runtime source or generator parity."}
			return result
		}
		result.Mutations = append(result.Mutations, target.path)
	}
	check := writeNativeRuntimeSQLiteIndex(repo, nativeIndexPath)
	result.Checks = append(result.Checks, check)
	if check.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_refresh_failed", Message: "native SQLite index failed: " + check.Message, Remediation: "Inspect native index output."}
		return result
	}
	result.Mutations = append(result.Mutations, nativeIndexPath)

	for _, check := range []Check{
		nativeRuntimeIndexValidation(repo, nativeIndexPath),
		nativeKnowledgeRuntimeValidation(repo),
	} {
		result.Checks = append(result.Checks, check)
		if check.Status != "ok" {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "runtime_refresh_failed", Message: check.Message, Remediation: "Inspect native refresh output and source files."}
			return result
		}
	}

	return result
}

func buildRuntimeCompileResult(opts runtimeOptions) Result {
	result := Result{
		Command:        "runtime compile",
		Mode:           "native_compiler",
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

	outputDB := runtimeCompileDBPath(repo, opts.dbPath)
	result.PlannedActions = append(result.PlannedActions, "compile runtime YAML and prose sources to runtime.db: "+outputDB)
	if opts.dryRun {
		result.PlannedActions = append(result.PlannedActions, "validate Go-native compiler output: "+outputDB)
	}
	if opts.assertSource != "" || opts.assertKeyword != "" {
		result.PlannedActions = append(result.PlannedActions, "assert generated surface: source="+opts.assertSource+" keyword="+opts.assertKeyword)
	}

	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "runtime_compile_native", Status: "ok", Message: "dry-run only; compiler not executed"})
		return result
	}
	check := buildNativeRuntimeDBFromSources(repo, outputDB)
	result.Checks = append(result.Checks, check)
	if check.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_compile_failed", Message: check.Message, Remediation: "Inspect runtime source YAML/prose and fix the compiler input."}
		return result
	}
	result.Mutations = append(result.Mutations, outputDB)

	return result
}

func runtimeCompileDBPath(repo string, dbPath string) string {
	if strings.TrimSpace(dbPath) == "" {
		return filepath.Join(repo, "runtime", "runtime.db")
	}
	if filepath.IsAbs(dbPath) {
		return dbPath
	}
	return filepath.Join(repo, dbPath)
}

func buildRuntimeQueryResult(opts runtimeOptions) Result {
	result := Result{
		Command:        "runtime query",
		Mode:           "native",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
		Results:        []QueryResult{},
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

	if opts.limit < 1 {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_limit", Message: "--limit must be greater than zero.", Remediation: "Pass --limit with a positive integer."}
		return result
	}
	if opts.graphQuery {
		return buildRuntimeGraphQueryResult(result, repo, opts)
	}

	keyword := strings.TrimSpace(opts.keyword)
	if keyword == "" && len(opts.positionals) > 0 {
		keyword = strings.TrimSpace(opts.positionals[0])
	}
	if keyword == "" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "missing_keyword", Message: "runtime query requires --keyword or a positional query.", Remediation: "Pass a search term, for example: ai-skill runtime query --keyword phase."}
		return result
	}

	dbPath := runtimeQueryDBPath(repo, opts.dbPath)
	result.Checks = append(result.Checks, Check{Name: "runtime_index", Status: "ok", Message: dbPath})
	result.PlannedActions = append(result.PlannedActions, "query native runtime index: "+dbPath)
	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "query_mode", Status: "ok", Message: "dry-run only; runtime index not queried"})
		return result
	}
	if _, err := os.Stat(dbPath); err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "missing_runtime_index", Message: dbPath, Remediation: "Run ai-skill runtime refresh after runtime index generation is available."}
		result.Checks[len(result.Checks)-1].Status = "missing"
		return result
	}

	queryResults, err := runRuntimeIndexQuery(dbPath, keyword, opts)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_query_failed", Message: err.Error(), Remediation: "Regenerate and validate the runtime SQLite index."}
		result.Checks = append(result.Checks, Check{Name: "runtime_query", Status: "failed", Message: err.Error()})
		return result
	}
	result.Results = queryResults
	result.Checks = append(result.Checks, Check{Name: "runtime_query", Status: "ok", Message: fmt.Sprintf("%d result(s)", len(queryResults))})
	return result
}

func runtimeRefreshMode(opts runtimeOptions) string {
	return "native_refresh"
}

func runtimeNativeReportTargets(repo string) []runtimeNativeReportTarget {
	return []runtimeNativeReportTarget{
		{name: "knowledge_runtime_report", path: filepath.Join(repo, "knowledge", "runtime", "runtime-report.md"), build: buildNativeKnowledgeRuntimeReport},
		{name: "model_context_report", path: filepath.Join(repo, "knowledge", "runtime", "model-context-report.md"), build: buildNativeModelContextReport},
		{name: "model_checklists", path: filepath.Join(repo, "knowledge", "runtime", "model-checklists.md"), build: buildNativeModelChecklists},
	}
}

func writeNativeRuntimeReport(repo string, target runtimeNativeReportTarget) Check {
	content, err := target.build(repo)
	if err != nil {
		return Check{Name: target.name, Status: "failed", Message: err.Error()}
	}
	if err := os.MkdirAll(filepath.Dir(target.path), 0o755); err != nil {
		return Check{Name: target.name, Status: "failed", Message: err.Error()}
	}
	if err := os.WriteFile(target.path, []byte(content), 0o644); err != nil {
		return Check{Name: target.name, Status: "failed", Message: err.Error()}
	}
	return Check{Name: target.name, Status: "ok", Message: target.path}
}

func writeNativeRuntimeSQLiteIndex(repo string, path string) Check {
	if err := buildNativeRuntimeSQLiteIndex(repo, path); err != nil {
		return Check{Name: "runtime_sqlite_index", Status: "failed", Message: err.Error()}
	}
	return Check{Name: "runtime_sqlite_index", Status: "ok", Message: path}
}

func buildNativeRuntimeSQLiteIndex(repo string, path string) error {
	records, err := runtimeIndexRecords(repo)
	if err != nil {
		return err
	}
	edges, err := runtimeIndexEdges(repo)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(`
PRAGMA journal_mode=OFF;
PRAGMA synchronous=OFF;
CREATE TABLE atoms (id TEXT PRIMARY KEY, source_path TEXT, layer TEXT, type TEXT, status TEXT, priority TEXT, confidence TEXT, context_cost TEXT, tags TEXT, domains TEXT, title TEXT, summary TEXT, when_to_read TEXT, validation_signal TEXT);
CREATE TABLE sources (source_path TEXT PRIMARY KEY, layer TEXT, title TEXT, checksum TEXT, bytes INTEGER);
CREATE TABLE edges (graph_id TEXT, source_path TEXT, edge_type TEXT, target_path TEXT, reason TEXT, validation TEXT);
CREATE VIRTUAL TABLE fts USING fts5(id, source_path, title, summary, tags, when_to_read, validation_signal);
`); err != nil {
		return err
	}
	for _, record := range records {
		if _, err := db.Exec(`INSERT INTO atoms VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			record.ID, record.SourcePath, record.Layer, record.Type, record.Status, record.Priority, record.Confidence, record.ContextCost, record.Tags, record.Domains, record.Title, record.Summary, record.WhenToRead, record.ValidationSignal); err != nil {
			return err
		}
		if _, err := db.Exec(`INSERT INTO fts VALUES (?, ?, ?, ?, ?, ?, ?)`,
			record.ID, record.SourcePath, record.Title, record.Summary, record.Tags, record.WhenToRead, record.ValidationSignal); err != nil {
			return err
		}
	}
	sources, err := runtimeIndexSources(repo, records)
	if err != nil {
		return err
	}
	for _, source := range sources {
		if _, err := db.Exec(`INSERT INTO sources VALUES (?, ?, ?, ?, ?)`, source.sourcePath, source.layer, source.title, source.checksum, source.bytes); err != nil {
			return err
		}
	}
	for _, edge := range edges {
		if _, err := db.Exec(`INSERT INTO edges VALUES (?, ?, ?, ?, ?, ?)`, edge.GraphID, edge.SourcePath, edge.EdgeType, edge.TargetPath, edge.Reason, edge.Validation); err != nil {
			return err
		}
	}
	return nil
}

func runtimeIndexRecords(repo string) ([]runtimeIndexAtom, error) {
	summaries, err := runtimeIndexSummaryRecords(repo)
	if err != nil {
		return nil, err
	}
	routes, err := runtimeIndexRouteRecords(repo)
	if err != nil {
		return nil, err
	}
	feedback, err := runtimeIndexFeedbackRecords(repo)
	if err != nil {
		return nil, err
	}
	records := append([]runtimeIndexAtom{}, summaries...)
	records = append(records, routes...)
	records = append(records, feedback...)
	return records, nil
}

func runtimeIndexSummaryRecords(repo string) ([]runtimeIndexAtom, error) {
	paths, err := filepath.Glob(filepath.Join(repo, "knowledge", "summaries", "*.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	records := []runtimeIndexAtom{}
	for _, path := range paths {
		if filepath.Base(path) == "README.md" {
			continue
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fields := parseRuntimeSummaryTable(string(content))
		sourcePath := runtimeRepoRelativeLink(repo, path, runtimeLinksFromCell(fields["Source path"]))
		if sourcePath == "" {
			relative, err := filepath.Rel(repo, path)
			if err != nil {
				return nil, err
			}
			sourcePath = filepath.ToSlash(relative)
		}
		lifecycle := runtimeStripMarkup(fields["Lifecycle"])
		confidence := "medium"
		if lifecycle == "validated" {
			confidence = "high"
		}
		records = append(records, runtimeIndexAtom{
			ID: runtimeStripMarkup(fields["Atom ID"]), SourcePath: sourcePath, Layer: runtimeLayerFor(sourcePath),
			Type: "summary", Status: lifecycle, Priority: "P2", Confidence: confidence, ContextCost: "low",
			Tags: "summary,atom", Title: runtimeStripMarkup(fields["Atom ID"]), Summary: fields["Summary"],
			WhenToRead: fields["When to read"], ValidationSignal: fields["Validation signal"],
		})
	}
	return records, nil
}

func runtimeIndexRouteRecords(repo string) ([]runtimeIndexAtom, error) {
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		return nil, err
	}
	records := []runtimeIndexAtom{}
	for _, record := range registry.Records {
		records = append(records, runtimeIndexAtom{
			ID: record.ID, SourcePath: record.PrimarySource, Layer: runtimeLayerFor(record.PrimarySource), Type: "route",
			Status: record.Metadata.CompatibilityState, Priority: record.Metadata.Priority, Confidence: record.Metadata.Confidence,
			ContextCost: rubyLikeString(record.Metadata.ContextCost), Tags: strings.Join(nonEmptyStrings("route", record.Model.Profile, record.Model.CompressionLevel), ","),
			Title: record.TaskIntent, Summary: record.RankingReason, WhenToRead: record.TaskIntent, ValidationSignal: record.ValidationSignal,
		})
	}
	return records, nil
}

func runtimeIndexFeedbackRecords(repo string) ([]runtimeIndexAtom, error) {
	roots, err := filepath.Glob(filepath.Join(repo, "skills", "*", "feedback_history"))
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for _, root := range roots {
		if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}
			paths = append(paths, path)
			return nil
		}); err != nil {
			return nil, err
		}
	}
	sort.Strings(paths)
	records := []runtimeIndexAtom{}
	for _, path := range paths {
		if filepath.Base(path) == "README.md" {
			continue
		}
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		content := string(contentBytes)
		relative, err := filepath.Rel(repo, path)
		if err != nil {
			return nil, err
		}
		relative = filepath.ToSlash(relative)
		parts := strings.Split(relative, "/")
		skill := ""
		if len(parts) > 1 {
			skill = parts[1]
		}
		category := ""
		if len(parts) > 4 {
			category = strings.Join(parts[3:len(parts)-1], "/")
		}
		status := regexp.MustCompile(`(?m)^Status:\s*([^\n]+)`).FindStringSubmatch(content)
		statusValue := ""
		if len(status) > 1 {
			statusValue = strings.TrimSpace(status[1])
		}
		if statusValue == "" {
			statusValue = "candidate"
		}
		confidence := "medium"
		if statusValue == "promoted" || statusValue == "validated" {
			confidence = "high"
		}
		title := firstRuntimeLineWithPrefix(content, "### ")
		title = strings.TrimSpace(strings.TrimPrefix(title, "### "))
		if title == "" {
			title = strings.TrimSuffix(filepath.Base(relative), ".md")
		}
		records = append(records, runtimeIndexAtom{
			ID: "feedback." + skill + "." + strings.TrimSuffix(filepath.Base(relative), ".md"), SourcePath: relative, Layer: "skills", Type: "feedback-pattern",
			Status: statusValue, Priority: "P2", Confidence: confidence, ContextCost: "medium", Tags: strings.Join(nonEmptyStrings(skill, "feedback", category), ","),
			Domains: skill, Title: title, Summary: firstRuntimeHeadingAfter(content, "#### One-line Summary"),
			WhenToRead: "Feedback lesson lookup for " + skill + ".", ValidationSignal: "Open canonical feedback lesson at " + relative + ".",
		})
	}
	return records, nil
}

func runtimeIndexEdges(repo string) ([]runtimeIndexEdge, error) {
	paths, err := filepath.Glob(filepath.Join(repo, "knowledge", "graphs", "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	edges := []runtimeIndexEdge{}
	for _, path := range paths {
		graph, err := readKnowledgeGraphFile(path)
		if err != nil {
			return nil, err
		}
		for _, edge := range graph.Edges {
			edges = append(edges, runtimeIndexEdge{GraphID: graph.ID, SourcePath: graph.Source, EdgeType: edge.Type, TargetPath: edge.Target, Reason: edge.Reason, Validation: edge.Validation})
		}
	}
	return edges, nil
}

type runtimeIndexSource struct {
	sourcePath string
	layer      string
	title      string
	checksum   string
	bytes      int
}

func runtimeIndexSources(repo string, records []runtimeIndexAtom) ([]runtimeIndexSource, error) {
	seen := map[string]bool{}
	sourcePaths := []string{}
	for _, record := range records {
		if record.SourcePath == "" || seen[record.SourcePath] {
			continue
		}
		seen[record.SourcePath] = true
		sourcePaths = append(sourcePaths, record.SourcePath)
	}
	sort.Strings(sourcePaths)
	sources := []runtimeIndexSource{}
	for _, sourcePath := range sourcePaths {
		path := filepath.Join(repo, filepath.FromSlash(sourcePath))
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		title := filepath.Base(sourcePath)
		checksum := ""
		size := 0
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			title = runtimeTitleFromMarkdown(string(content))
			sum := sha256.Sum256(content)
			checksum = fmt.Sprintf("%x", sum)
			size = len(content)
		}
		sources = append(sources, runtimeIndexSource{sourcePath: sourcePath, layer: runtimeLayerFor(sourcePath), title: title, checksum: checksum, bytes: size})
	}
	return sources, nil
}

func runtimeQueryDBPath(repo string, dbPath string) string {
	if strings.TrimSpace(dbPath) == "" {
		return filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	}
	if filepath.IsAbs(dbPath) {
		return filepath.Clean(dbPath)
	}
	return filepath.Join(repo, dbPath)
}

func runRuntimeIndexQuery(dbPath string, keyword string, opts runtimeOptions) ([]QueryResult, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `SELECT bm25(fts) AS rank,
       atoms.id,
       atoms.source_path,
       atoms.layer,
       atoms.type,
       atoms.status,
       atoms.priority,
       atoms.confidence,
       atoms.context_cost,
       atoms.summary
FROM fts
JOIN atoms ON atoms.id = fts.id
WHERE fts MATCH ?`
	args := []any{runtimeFTSMatchLiteral(keyword)}
	if opts.layer != "" {
		sqlText += " AND atoms.layer = ?"
		args = append(args, opts.layer)
	}
	if opts.queryType != "" {
		sqlText += " AND atoms.type = ?"
		args = append(args, opts.queryType)
	}
	if opts.statusFilter != "" {
		sqlText += " AND atoms.status = ?"
		args = append(args, opts.statusFilter)
	}
	sqlText += ` ORDER BY rank,
       CASE atoms.priority WHEN 'P0' THEN 0 WHEN 'P1' THEN 1 WHEN 'P2' THEN 2 ELSE 3 END,
       CASE atoms.confidence WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END,
       CASE atoms.context_cost WHEN 'low' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END,
       atoms.id
LIMIT ?`
	args = append(args, opts.limit)

	rows, err := db.Query(sqlText, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []QueryResult{}
	for rows.Next() {
		var item QueryResult
		if err := rows.Scan(&item.Rank, &item.ID, &item.SourcePath, &item.Layer, &item.Type, &item.Status, &item.Priority, &item.Confidence, &item.ContextCost, &item.Summary); err != nil {
			return nil, err
		}
		item.MatchReason = "fts keyword match: " + keyword
		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func runtimeFTSMatchLiteral(keyword string) string {
	escaped := strings.ReplaceAll(keyword, `"`, `""`)
	return `"` + escaped + `"`
}

func buildRuntimeGraphQueryResult(result Result, repo string, opts runtimeOptions) Result {
	keyword := strings.TrimSpace(opts.keyword)
	if keyword == "" && len(opts.positionals) > 0 {
		keyword = strings.TrimSpace(opts.positionals[0])
	}
	if keyword == "" && opts.source == "" && opts.target == "" && opts.queryType == "" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "missing_graph_filter", Message: "runtime graph query requires --source, --target, --type, --keyword, or a positional query.", Remediation: "Pass at least one graph filter."}
		return result
	}

	graphDir := filepath.Join(repo, "knowledge", "graphs")
	result.Checks = append(result.Checks, Check{Name: "knowledge_graphs", Status: "ok", Message: graphDir})
	result.PlannedActions = append(result.PlannedActions, "query native knowledge graph edges: "+graphDir)
	if opts.dryRun {
		result.Checks = append(result.Checks, Check{Name: "query_mode", Status: "ok", Message: "dry-run only; knowledge graph files not queried"})
		return result
	}

	results, err := runKnowledgeGraphQuery(repo, graphDir, keyword, opts)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "runtime_graph_query_failed", Message: err.Error(), Remediation: "Inspect knowledge/graphs YAML files and query filters."}
		result.Checks = append(result.Checks, Check{Name: "knowledge_graph_query", Status: "failed", Message: err.Error()})
		return result
	}
	result.Results = results
	result.Checks = append(result.Checks, Check{Name: "knowledge_graph_query", Status: "ok", Message: fmt.Sprintf("%d result(s)", len(results))})
	return result
}

type knowledgeGraphFile struct {
	ID     string               `yaml:"id"`
	Source string               `yaml:"source"`
	Status string               `yaml:"status"`
	Edges  []knowledgeGraphEdge `yaml:"edges"`
}

type knowledgeGraphEdge struct {
	Type       string `yaml:"type"`
	Target     string `yaml:"target"`
	Reason     string `yaml:"reason"`
	Validation string `yaml:"validation"`
}

func runKnowledgeGraphQuery(repo string, graphDir string, keyword string, opts runtimeOptions) ([]QueryResult, error) {
	entries, err := os.ReadDir(graphDir)
	if err != nil {
		return nil, err
	}
	results := []QueryResult{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		path := filepath.Join(graphDir, entry.Name())
		graph, err := readKnowledgeGraphFile(path)
		if err != nil {
			return nil, err
		}
		graphFile, err := filepath.Rel(repo, path)
		if err != nil {
			return nil, err
		}
		for _, edge := range graph.Edges {
			item := QueryResult{
				ID:          graph.ID,
				SourcePath:  graph.Source,
				Type:        edge.Type,
				Status:      graph.Status,
				GraphID:     graph.ID,
				GraphSource: graph.Source,
				EdgeType:    edge.Type,
				Target:      edge.Target,
				Reason:      edge.Reason,
				Validation:  edge.Validation,
				GraphFile:   filepath.ToSlash(graphFile),
				MatchReason: "knowledge graph filter match",
			}
			if !knowledgeGraphResultMatches(item, keyword, opts) {
				continue
			}
			results = append(results, item)
			if len(results) >= opts.limit {
				return results, nil
			}
		}
	}
	return results, nil
}

func readKnowledgeGraphFile(path string) (knowledgeGraphFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return knowledgeGraphFile{}, err
	}
	var graph knowledgeGraphFile
	if err := yaml.Unmarshal(content, &graph); err != nil {
		return knowledgeGraphFile{}, err
	}
	return graph, nil
}

func knowledgeGraphResultMatches(item QueryResult, keyword string, opts runtimeOptions) bool {
	if !containsFold(item.GraphSource, opts.source) {
		return false
	}
	if !containsFold(item.Target, opts.target) {
		return false
	}
	if !containsFold(item.EdgeType, opts.queryType) {
		return false
	}
	if keyword != "" {
		haystack := strings.Join([]string{item.GraphID, item.GraphSource, item.EdgeType, item.Target, item.Reason, item.Validation}, " ")
		if !containsFold(haystack, keyword) {
			return false
		}
	}
	return true
}

func containsFold(value string, needle string) bool {
	if needle == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}

func nativeRuntimeIndexValidation(repo string, path string) Check {
	if _, err := os.Stat(path); err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: "runtime index not found: " + path}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: err.Error()}
	}
	defer db.Close()

	if err := nativeIntegrityCheck(db); err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: err.Error()}
	}
	if err := nativeRuntimeIndexTablesCheck(db); err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: err.Error()}
	}
	counts, err := nativeRuntimeIndexCountsCheck(db)
	if err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: err.Error()}
	}
	if err := nativeRuntimeIndexSourceReferencesCheck(db); err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: err.Error()}
	}
	if err := nativeRuntimeIndexChecksumsCheck(repo, db); err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: err.Error()}
	}
	if err := nativeRuntimeIndexFTSCheck(db); err != nil {
		return Check{Name: "runtime_index_native", Status: "failed", Message: err.Error()}
	}
	return Check{Name: "runtime_index_native", Status: "ok", Message: fmt.Sprintf("SQLite runtime index checks passed: atoms=%d sources=%d edges=%d fts=%d", counts["atoms"], counts["sources"], counts["edges"], counts["fts"])}
}

func nativeKnowledgeRuntimeValidation(repo string) Check {
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		return Check{Name: "knowledge_runtime_native", Status: "failed", Message: err.Error()}
	}
	if len(registry.Records) == 0 {
		return Check{Name: "knowledge_runtime_native", Status: "failed", Message: "routing registry has no records"}
	}
	summaries, err := runtimeSummaryRecords(repo)
	if err != nil {
		return Check{Name: "knowledge_runtime_native", Status: "failed", Message: err.Error()}
	}
	if len(summaries) == 0 {
		return Check{Name: "knowledge_runtime_native", Status: "failed", Message: "knowledge summaries are empty"}
	}
	graphs, err := runtimeGraphRecords(repo)
	if err != nil {
		return Check{Name: "knowledge_runtime_native", Status: "failed", Message: err.Error()}
	}
	if len(graphs) == 0 {
		return Check{Name: "knowledge_runtime_native", Status: "failed", Message: "knowledge graphs are empty"}
	}
	reportAnchors := map[string]string{
		filepath.Join(repo, "knowledge", "runtime", "runtime-report.md"):       "# Knowledge Runtime Report",
		filepath.Join(repo, "knowledge", "runtime", "model-context-report.md"): "# Model Context Report",
		filepath.Join(repo, "knowledge", "runtime", "model-checklists.md"):     "# Model Checklists",
		filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"):   "records:",
		filepath.Join(repo, "knowledge", "runtime", "refresh-policy.yaml"):     "decision_values:",
		filepath.Join(repo, "knowledge", "summaries", "README.md"):             "knowledge",
	}
	for path, anchor := range reportAnchors {
		content, err := os.ReadFile(path)
		if err != nil {
			return Check{Name: "knowledge_runtime_native", Status: "failed", Message: err.Error()}
		}
		if !strings.Contains(string(content), anchor) {
			return Check{Name: "knowledge_runtime_native", Status: "failed", Message: filepath.ToSlash(path) + " missing anchor " + anchor}
		}
	}
	return Check{Name: "knowledge_runtime_native", Status: "ok", Message: fmt.Sprintf("Knowledge runtime checks passed: registry_records=%d summaries=%d graphs=%d", len(registry.Records), len(summaries), len(graphs))}
}

func nativeRuntimeIndexTablesCheck(db *sql.DB) error {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type IN ('table', 'virtual')")
	if err != nil {
		return err
	}
	defer rows.Close()

	tables := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables[name] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, table := range []string{"atoms", "sources", "edges", "fts"} {
		if !tables[table] {
			return fmt.Errorf("missing table: %s", table)
		}
	}
	return nil
}

func nativeRuntimeIndexCountsCheck(db *sql.DB) (map[string]int, error) {
	counts := map[string]int{}
	for _, table := range []string{"atoms", "sources", "edges", "fts"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			return nil, fmt.Errorf("%s count failed: %w", table, err)
		}
		counts[table] = count
	}
	if counts["atoms"] == 0 {
		return nil, fmt.Errorf("atoms table is empty")
	}
	if counts["sources"] == 0 {
		return nil, fmt.Errorf("sources table is empty")
	}
	if counts["fts"] != counts["atoms"] {
		return nil, fmt.Errorf("fts count does not match atoms count")
	}
	return counts, nil
}

func nativeRuntimeIndexSourceReferencesCheck(db *sql.DB) error {
	rows, err := db.Query("SELECT source_path FROM atoms WHERE source_path NOT IN (SELECT source_path FROM sources) LIMIT 10")
	if err != nil {
		return err
	}
	defer rows.Close()

	missing := []string{}
	for rows.Next() {
		var sourcePath string
		if err := rows.Scan(&sourcePath); err != nil {
			return err
		}
		missing = append(missing, sourcePath)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(missing) > 0 {
		return fmt.Errorf("atoms reference missing sources: %s", strings.Join(missing, ", "))
	}
	return nil
}

func nativeRuntimeIndexChecksumsCheck(repo string, db *sql.DB) error {
	rows, err := db.Query("SELECT source_path, checksum FROM sources WHERE checksum IS NOT NULL AND checksum != ''")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sourcePath, checksum string
		if err := rows.Scan(&sourcePath, &checksum); err != nil {
			return err
		}
		path := filepath.Join(repo, sourcePath)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("source path missing on disk: %s", sourcePath)
		}
		current := fmt.Sprintf("%x", sha256.Sum256(content))
		if current != checksum {
			return fmt.Errorf("stale checksum for %s", sourcePath)
		}
	}
	return rows.Err()
}

func nativeRuntimeIndexFTSCheck(db *sql.DB) error {
	var feedbackHits int
	if err := db.QueryRow("SELECT COUNT(*) FROM fts WHERE fts MATCH ?", runtimeFTSMatchLiteral("feedback")).Scan(&feedbackHits); err != nil {
		return err
	}
	if feedbackHits == 0 {
		return fmt.Errorf("expected feedback FTS hits")
	}
	var routeHits int
	if err := db.QueryRow("SELECT COUNT(*) FROM fts WHERE fts MATCH ?", runtimeFTSMatchLiteral("route")).Scan(&routeHits); err != nil {
		return err
	}
	if routeHits == 0 {
		return fmt.Errorf("expected route FTS hits")
	}
	var rankedRoute string
	err := db.QueryRow(`SELECT atoms.id FROM fts
JOIN atoms ON atoms.id = fts.id
WHERE fts MATCH ?
ORDER BY bm25(fts), CASE atoms.priority WHEN 'P0' THEN 0 WHEN 'P1' THEN 1 WHEN 'P2' THEN 2 ELSE 3 END
LIMIT 1`, runtimeFTSMatchLiteral("feedback")).Scan(&rankedRoute)
	if err != nil || rankedRoute == "" {
		return fmt.Errorf("expected ranked query result")
	}
	return nil
}

func nativeRuntimeIndexGitIgnoreCheck(repo string, path string, git string) Check {
	rel, err := filepath.Rel(repo, path)
	if err != nil {
		return Check{Name: "runtime_index_git_ignore", Status: "failed", Message: err.Error()}
	}
	rel = filepath.ToSlash(rel)
	output, err := exec.Command(git, "-C", repo, "check-ignore", rel).CombinedOutput()
	message := strings.TrimSpace(string(output))
	if err != nil || message == "" {
		if message == "" {
			message = "generated DB is not ignored by git: " + rel
		}
		return Check{Name: "runtime_index_git_ignore", Status: "failed", Message: message}
	}
	return Check{Name: "runtime_index_git_ignore", Status: "ok", Message: rel + " is ignored by git"}
}

var nativeRuntimeRequiredTables = []string{
	"phases", "phase_transitions", "obligations", "gates",
	"transaction_states", "transaction_transitions", "transaction_rules", "transaction_templates",
	"activation_rules", "core_bootstrap_rules",
	"discovery_checkpoints", "discovery_search_strategy",
	"generated_surfaces", "compiler_metadata",
	"runtime_budget", "context_ttl_policy", "circuit_breaker", "context_pollution",
	"context_health_score", "intelligence_routing", "obligation_ledger",
	"language_policy", "output_rules", "governance_gates", "blocking_gates",
	"phase_machine", "pipeline_context_flow", "guard_chain", "relevance_engine",
	"session_lifecycle", "prompt_artifact_templates", "prompt_composition_rules",
	"recovery_strategies", "state_repair", "obligation_rebuild", "phase_reconciliation",
	"execution_queue", "priority_scheduler", "activation_rules_mirror",
	"transaction_templates_ext", "distributed_locks", "multi_agent_coordination",
	"async_job_lifecycle", "capability_checkpoints",
	"compiler_rules",
}

var nativeRuntimeMinimumRows = map[string]int{
	"phases": 8, "obligations": 15, "gates": 15, "activation_rules": 10,
	"core_bootstrap_rules": 2, "discovery_checkpoints": 3, "compiler_metadata": 2,
	"runtime_budget": 1, "context_ttl_policy": 1, "circuit_breaker": 1,
	"context_pollution": 1, "context_health_score": 1, "intelligence_routing": 1,
	"obligation_ledger": 1, "language_policy": 1, "output_rules": 1,
	"governance_gates": 1, "blocking_gates": 1, "phase_machine": 1,
	"pipeline_context_flow": 1, "guard_chain": 1, "relevance_engine": 1,
	"session_lifecycle": 1, "prompt_artifact_templates": 1, "prompt_composition_rules": 1,
	"recovery_strategies": 1, "state_repair": 1, "obligation_rebuild": 1,
	"phase_reconciliation": 1, "execution_queue": 1, "priority_scheduler": 1,
	"activation_rules_mirror": 1, "transaction_templates_ext": 1,
	"distributed_locks": 1, "multi_agent_coordination": 1, "async_job_lifecycle": 1,
	"capability_checkpoints": 1, "compiler_rules": 1,
}

var nativeRuntimeJSONColumns = map[string][]string{
	"phases":                []string{"entry_conditions", "allowed_actions", "forbidden_actions", "blocking_gates", "obligations", "phase_transition_triggers"},
	"obligations":           []string{"verification", "depends_on", "linked_gates"},
	"transaction_states":    []string{"entry_conditions", "allowed_actions", "forbidden_actions", "blocking_gates"},
	"activation_rules":      []string{"activation_when"},
	"discovery_checkpoints": []string{"discovery_targets"},
}

func nativeRuntimeDBValidation(path string) Check {
	info, err := os.Stat(path)
	if err != nil {
		return Check{Name: "runtime_db_native", Status: "failed", Message: "runtime.db not found: " + path}
	}
	if info.IsDir() {
		return Check{Name: "runtime_db_native", Status: "failed", Message: "runtime.db path is a directory: " + path}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return Check{Name: "runtime_db_native", Status: "failed", Message: err.Error()}
	}
	defer db.Close()

	if err := nativeIntegrityCheck(db); err != nil {
		return Check{Name: "runtime_db_native", Status: "failed", Message: err.Error()}
	}
	if err := nativeRequiredTablesCheck(db); err != nil {
		return Check{Name: "runtime_db_native", Status: "failed", Message: err.Error()}
	}
	if err := nativeMinimumRowsCheck(db); err != nil {
		return Check{Name: "runtime_db_native", Status: "failed", Message: err.Error()}
	}
	if err := nativeJSONColumnsCheck(db); err != nil {
		return Check{Name: "runtime_db_native", Status: "failed", Message: err.Error()}
	}
	warning, err := nativeCompilerMetadataCheck(db)
	if err != nil {
		return Check{Name: "runtime_db_native", Status: "failed", Message: err.Error()}
	}
	message := "Go native runtime.db integrity, schema, row count, JSON, and compiler metadata checks passed"
	if warning != "" {
		message += "; warning: " + warning
	}
	return Check{Name: "runtime_db_native", Status: "ok", Message: message}
}

func nativeIntegrityCheck(db *sql.DB) error {
	var result string
	if err := db.QueryRow("PRAGMA integrity_check").Scan(&result); err != nil {
		return fmt.Errorf("integrity_check query failed: %w", err)
	}
	if result != "ok" {
		return fmt.Errorf("integrity_check failed: %s", result)
	}
	return nil
}

func nativeRequiredTablesCheck(db *sql.DB) error {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type IN ('table', 'virtual')")
	if err != nil {
		return fmt.Errorf("table lookup failed: %w", err)
	}
	defer rows.Close()

	tables := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables[name] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, table := range nativeRuntimeRequiredTables {
		if !tables[table] {
			return fmt.Errorf("missing required table: %s", table)
		}
	}
	return nil
}

func nativeMinimumRowsCheck(db *sql.DB) error {
	for table, minimum := range nativeRuntimeMinimumRows {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			return fmt.Errorf("%s count failed: %w", table, err)
		}
		if count < minimum {
			return fmt.Errorf("%s: %d rows, expected at least %d", table, count, minimum)
		}
	}
	return nil
}

func nativeJSONColumnsCheck(db *sql.DB) error {
	for table, columns := range nativeRuntimeJSONColumns {
		rows, err := db.Query("SELECT " + strings.Join(columns, ", ") + " FROM " + table + " LIMIT 5")
		if err != nil {
			return fmt.Errorf("%s JSON query failed: %w", table, err)
		}
		values := make([]sql.NullString, len(columns))
		scanTargets := make([]any, len(columns))
		for i := range values {
			scanTargets[i] = &values[i]
		}
		rowIndex := 0
		for rows.Next() {
			if err := rows.Scan(scanTargets...); err != nil {
				rows.Close()
				return err
			}
			for i, value := range values {
				text := strings.TrimSpace(value.String)
				if !value.Valid || text == "" || text == "[]" || text == "{}" {
					continue
				}
				var decoded any
				if err := json.Unmarshal([]byte(text), &decoded); err != nil {
					rows.Close()
					return fmt.Errorf("%s.%s row %d invalid JSON", table, columns[i], rowIndex)
				}
			}
			rowIndex++
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}
		rows.Close()
	}
	return nil
}

var nativeRuntimeNow = time.Now

func nativeCompilerMetadataCheck(db *sql.DB) (string, error) {
	rows, err := db.Query("SELECT key, value FROM compiler_metadata")
	if err != nil {
		return "", fmt.Errorf("compiler_metadata query failed: %w", err)
	}
	defer rows.Close()

	metadata := map[string]string{}
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return "", err
		}
		metadata[key] = value
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if metadata["compiler_version"] == "" {
		return "", fmt.Errorf("compiler_metadata missing compiler_version")
	}
	if metadata["compiled_at"] == "" {
		return "", fmt.Errorf("compiler_metadata missing compiled_at")
	}
	compiledAt, err := time.Parse(time.RFC3339, metadata["compiled_at"])
	if err != nil {
		return "", fmt.Errorf("compiler_metadata compiled_at invalid: %w", err)
	}
	age := nativeRuntimeNow().Sub(compiledAt)
	if age > 24*time.Hour {
		return fmt.Sprintf("runtime.db is %.1f hours old (compiled at %s)", age.Hours(), metadata["compiled_at"]), nil
	}
	return "", nil
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

func buildNativeModelContextReport(repo string) (string, error) {
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		return "", err
	}
	records := registry.Records
	lines := []string{
		"# Model Context Report",
		"",
		"本檔由 `ai-skill runtime refresh` 產生，依 `knowledge/runtime/routing-registry.yaml` 的 model 欄位整理 model-aware context loading view。",
		"",
		"## Source Surfaces",
		"",
		"| Surface | Path | Purpose |",
		"| --- | --- | --- |",
		"| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 提供每條 route 的 model profile 與 compression level。 |",
		"| Model profiles | [`../../models/profiles/README.md`](../../models/profiles/README.md) | 定義 `small`、`large`、`specialized` 的讀取深度與 guardrails。 |",
		"| Compression strategy | [`../../models/compression/README.md`](../../models/compression/README.md) | 定義 `summary-first`、`source-backed`、`graph-assisted` 等壓縮層級。 |",
		"",
		"## Profile View",
		"",
	}

	for _, profile := range sortedRouteGroupKeys(records, func(record runtimeRouteRecord) string {
		return record.Model.Profile
	}) {
		displayProfile := profile
		if displayProfile == "" {
			displayProfile = "unspecified"
		}
		lines = append(lines,
			"### `"+runtimeMDEscape(displayProfile)+"`",
			"",
			"| Route | Primary source | Compression | Reason |",
			"| --- | --- | --- | --- |",
		)
		for _, record := range records {
			if record.Model.Profile != profile {
				continue
			}
			lines = append(lines, "| `"+runtimeMDEscape(record.ID)+"` | `"+runtimeMDEscape(record.PrimarySource)+"` | `"+runtimeMDEscape(record.Model.CompressionLevel)+"` | "+runtimeMDEscape(record.Model.Reason)+" |")
		}
		lines = append(lines, "")
	}

	lines = append(lines,
		"## Compression View",
		"",
		"| Compression level | Routes | Escalation note |",
		"| --- | --- | --- |",
	)
	for _, level := range sortedRouteGroupKeys(records, func(record runtimeRouteRecord) string {
		return record.Model.CompressionLevel
	}) {
		routeIDs := []string{}
		for _, record := range records {
			if record.Model.CompressionLevel == level {
				routeIDs = append(routeIDs, "`"+runtimeMDEscape(record.ID)+"`")
			}
		}
		lines = append(lines, "| `"+runtimeMDEscape(level)+"` | "+strings.Join(routeIDs, ", ")+" | "+runtimeCompressionEscalationNote(level)+" |")
	}

	lines = append(lines,
		"",
		"## Agent Output Shape",
		"",
		"使用本 report 決定 model-aware loading 時，回報：",
		"",
		"```text",
		"Profile:",
		"Compression level:",
		"Primary source:",
		"Summaries used:",
		"Required full sources:",
		"Deferred sources:",
		"Escalation trigger:",
		"Validation signal:",
		"```",
		"",
		"## Validation",
		"",
		"- 產生前應先確認 `routing-registry.yaml` 可通過 `ai-skill runtime validate`。",
		"- 產生後應重新執行 `ai-skill runtime validate`，檢查本 report links。",
		"- 本報告是 generated view，不取代 `models/profiles/README.md`、`models/compression/README.md` 或 routing registry。",
		"",
	)
	return strings.Join(lines, "\n"), nil
}

func buildNativeModelChecklists(repo string) (string, error) {
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		return "", err
	}
	records := registry.Records
	lines := []string{
		"# Model Checklists",
		"",
		"本檔由 `ai-skill runtime refresh` 產生，將 routing registry 中的 model profile / compression level 轉成 agent 可直接使用的 context-loading checklist。",
		"",
		"## Source Surfaces",
		"",
		"| Surface | Path | Purpose |",
		"| --- | --- | --- |",
		"| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | 提供 route、primary source、dependencies、model profile 與 compression level。 |",
		"| Model profiles | [`../../models/profiles/README.md`](../../models/profiles/README.md) | 定義 profile guardrails。 |",
		"| Compression strategy | [`../../models/compression/README.md`](../../models/compression/README.md) | 定義 escalation rules。 |",
		"",
		"## Profile Checklists",
		"",
	}

	for _, profile := range sortedRouteGroupKeys(records, func(record runtimeRouteRecord) string {
		return record.Model.Profile
	}) {
		displayProfile := profile
		if displayProfile == "" {
			displayProfile = "unspecified"
		}
		lines = append(lines,
			"### `"+runtimeMDEscape(displayProfile)+"`",
			"",
			"Guardrails:",
			"",
		)
		for _, item := range runtimeProfileGuardrails(profile) {
			lines = append(lines, "- "+item)
		}
		lines = append(lines,
			"",
			"| Route | Checklist |",
			"| --- | --- |",
		)
		for _, record := range records {
			if record.Model.Profile != profile {
				continue
			}
			dependencies := []string{}
			for _, path := range record.RequiredDependencies {
				dependencies = append(dependencies, "`"+path+"`")
			}
			checklist := strings.Join([]string{
				"Primary: `" + record.PrimarySource + "`",
				"Compression: `" + record.Model.CompressionLevel + "`",
				"Required: " + strings.Join(dependencies, "<br>"),
				"Validation: " + runtimeMDEscape(record.ValidationSignal),
			}, "<br>")
			lines = append(lines, "| `"+runtimeMDEscape(record.ID)+"` | "+checklist+" |")
		}
		lines = append(lines, "")
	}

	lines = append(lines,
		"## Escalation Checklist",
		"",
		"- Summary / registry 與 source-of-truth 可能不一致時，讀全文。",
		"- 任務需要修改、commit、push、readback 或 promotion 時，升級到 `source-backed`。",
		"- 涉及 safety、secrets、authorization、source/mirror 或 destructive actions 時，升級到 full source 和 enforcement rules。",
		"- Routing registry 指向 candidate path，但 old entrypoint 仍 active 時，保留 old entrypoint gate。",
		"- Validation signal 不足以支持結論時，停止並讀 required dependencies。",
		"",
		"## Validation",
		"",
		"- 產生前應先確認 `routing-registry.yaml` 可通過 `ai-skill runtime validate`。",
		"- 產生後應重新執行 `ai-skill runtime validate`，檢查本 report links。",
		"- 本檔是 generated view，不取代 model source docs 或 routing registry。",
		"",
	)
	return strings.Join(lines, "\n"), nil
}

func buildNativeKnowledgeRuntimeReport(repo string) (string, error) {
	registry, err := readRuntimeRoutingRegistry(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		return "", err
	}
	summaries, err := runtimeSummaryRecords(repo)
	if err != nil {
		return "", err
	}
	graphs, err := runtimeGraphRecords(repo)
	if err != nil {
		return "", err
	}
	policy, err := readRuntimeRefreshPolicy(filepath.Join(repo, "knowledge", "runtime", "refresh-policy.yaml"))
	if err != nil {
		return "", err
	}
	status := policy.Status
	if status == "" {
		status = "unknown"
	}

	lines := []string{
		"# Knowledge Runtime Report",
		"",
		"本檔由 `ai-skill runtime refresh` 產生，彙整 runtime registry、summaries、graphs 與 refresh policy 的目前狀態。",
		"",
		"## Source Surfaces",
		"",
		"| Surface | Path | Count / Status |",
		"| --- | --- | --- |",
		fmt.Sprintf("| Routing registry | [`routing-registry.yaml`](routing-registry.yaml) | %d records |", len(registry.Records)),
		"| Refresh policy | [`refresh-policy.yaml`](refresh-policy.yaml) | " + status + " |",
		"| Model context report | [`model-context-report.md`](model-context-report.md) | generated view |",
		"| Model checklists | [`model-checklists.md`](model-checklists.md) | generated view |",
		"| SQLite runtime index | [`sqlite/`](sqlite/) | generated lookup cache prototype |",
		fmt.Sprintf("| Summaries | [`../summaries/`](../summaries/) | %d files |", len(summaries)),
		fmt.Sprintf("| Graph records | [`../graphs/`](../graphs/) | %d files |", len(graphs)),
		"",
		"## Routing Records",
		"",
		"| ID | Primary source | Model | Compression | Validation signal |",
		"| --- | --- | --- | --- | --- |",
	}
	for _, record := range registry.Records {
		lines = append(lines, "| `"+runtimeMDEscape(record.ID)+"` | `"+runtimeMDEscape(record.PrimarySource)+"` | `"+runtimeMDEscape(record.Model.Profile)+"` | `"+runtimeMDEscape(record.Model.CompressionLevel)+"` | "+runtimeMDEscape(record.ValidationSignal)+" |")
	}
	lines = append(lines,
		"",
		"## Summary Records",
		"",
		"| Atom ID | Lifecycle | File | Summary |",
		"| --- | --- | --- | --- |",
	)
	for _, summary := range summaries {
		base := filepath.Base(summary.File)
		lines = append(lines, "| `"+runtimeMDEscape(summary.AtomID)+"` | `"+runtimeMDEscape(summary.Lifecycle)+"` | [`"+runtimeMDEscape(base)+"`](../summaries/"+base+") | "+runtimeMDEscape(summary.Summary)+" |")
	}
	lines = append(lines,
		"",
		"## Graph Records",
		"",
		"| ID | Source | Status | Edges | File |",
		"| --- | --- | --- | --- | --- |",
	)
	for _, graph := range graphs {
		base := filepath.Base(graph.File)
		lines = append(lines, "| `"+runtimeMDEscape(graph.ID)+"` | `"+runtimeMDEscape(graph.Source)+"` | `"+runtimeMDEscape(graph.Status)+"` | "+fmt.Sprintf("%d", graph.EdgeCount)+" | [`"+runtimeMDEscape(base)+"`](../graphs/"+base+") |")
	}
	lines = append(lines,
		"",
		"## Refresh Decisions",
		"",
		"| Decision value | Meaning |",
		"| --- | --- |",
	)
	for _, decision := range policy.DecisionValues {
		lines = append(lines, "| `"+runtimeMDEscape(decision)+"` | 由 `refresh-policy.yaml` 定義的 generated surface decision。 |")
	}
	lines = append(lines,
		"",
		"## Validation",
		"",
		"- 產生前應先執行 `ai-skill runtime validate`。",
		"- 產生後應執行 Markdown link check、lints、close-loop dry run、commit / push / readback。",
		"- 本報告是 generated view，不取代 `routing-registry.yaml`、`refresh-policy.yaml`、summary 或 graph source files。",
		"",
	)
	return strings.Join(lines, "\n"), nil
}

func readRuntimeRoutingRegistry(path string) (runtimeRoutingRegistry, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return runtimeRoutingRegistry{}, err
	}
	var registry runtimeRoutingRegistry
	if err := yaml.Unmarshal(content, &registry); err != nil {
		return runtimeRoutingRegistry{}, err
	}
	return registry, nil
}

func readRuntimeRefreshPolicy(path string) (runtimeRefreshPolicy, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return runtimeRefreshPolicy{}, err
	}
	var policy runtimeRefreshPolicy
	if err := yaml.Unmarshal(content, &policy); err != nil {
		return runtimeRefreshPolicy{}, err
	}
	return policy, nil
}

func runtimeSummaryRecords(repo string) ([]runtimeSummaryRecord, error) {
	paths, err := filepath.Glob(filepath.Join(repo, "knowledge", "summaries", "*.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	records := []runtimeSummaryRecord{}
	for _, path := range paths {
		if filepath.Base(path) == "README.md" {
			continue
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fields := parseRuntimeSummaryTable(string(content))
		relative, err := filepath.Rel(repo, path)
		if err != nil {
			return nil, err
		}
		records = append(records, runtimeSummaryRecord{
			File:      filepath.ToSlash(relative),
			AtomID:    strings.ReplaceAll(fields["Atom ID"], "`", ""),
			Lifecycle: strings.ReplaceAll(fields["Lifecycle"], "`", ""),
			Summary:   fields["Summary"],
		})
	}
	return records, nil
}

func parseRuntimeSummaryTable(content string) map[string]string {
	fields := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(strings.TrimSpace(line), "|")
		trimmed := []string{}
		for _, cell := range cells {
			trimmed = append(trimmed, strings.TrimSpace(cell))
		}
		if len(trimmed) < 3 {
			continue
		}
		key := trimmed[1]
		if key == "欄位" || runtimeMarkdownSeparator(key) {
			continue
		}
		fields[key] = trimmed[2]
	}
	return fields
}

func runtimeMarkdownSeparator(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char != '-' {
			return false
		}
	}
	return true
}

func runtimeGraphRecords(repo string) ([]runtimeGraphRecord, error) {
	paths, err := filepath.Glob(filepath.Join(repo, "knowledge", "graphs", "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	records := []runtimeGraphRecord{}
	for _, path := range paths {
		graph, err := readKnowledgeGraphFile(path)
		if err != nil {
			return nil, err
		}
		relative, err := filepath.Rel(repo, path)
		if err != nil {
			return nil, err
		}
		records = append(records, runtimeGraphRecord{
			File:      filepath.ToSlash(relative),
			ID:        graph.ID,
			Source:    graph.Source,
			Status:    graph.Status,
			EdgeCount: len(graph.Edges),
		})
	}
	return records, nil
}

func sortedRouteGroupKeys(records []runtimeRouteRecord, group func(runtimeRouteRecord) string) []string {
	seen := map[string]bool{}
	keys := []string{}
	for _, record := range records {
		key := group(record)
		if seen[key] {
			continue
		}
		seen[key] = true
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func runtimeMDEscape(value string) string {
	value = strings.ReplaceAll(value, "|", `\|`)
	return strings.ReplaceAll(value, "\n", " ")
}

func runtimeCompressionEscalationNote(level string) string {
	switch level {
	case "summary-first":
		return "適合先用 registry / summary 判斷 relevance；修改 source 時升級。"
	case "source-backed":
		return "需要 primary source 與 required dependencies；適合 writeback、migration 或 domain work。"
	case "graph-assisted":
		return "需要 graph records 輔助 dependency / conflict / promotion reasoning。"
	default:
		return "依 `models/compression/README.md` 的 escalation rules 判斷。"
	}
}

func runtimeProfileGuardrails(profile string) []string {
	switch profile {
	case "small":
		return []string{
			"先讀 index、registry、summary 或 generated lookup。",
			"不可跳過 required bootstrap、source-of-truth gate 或 validation signal。",
			"需要修改 canonical source、遇到 conflict、缺 validation signal 時升級。",
		}
	case "large":
		return []string{
			"讀 primary source、required dependencies 與 task-relevant related sources。",
			"回報 deferred sources 與 validation signal。",
			"涉及 writeback、promotion、enforcement rules 或 migration 時保持 source-backed。",
		}
	case "specialized":
		return []string{
			"先讀 routing registry 與 primary source，再讀 domain workflow / technique / adapter。",
			"不得讓工具能力覆蓋 enforcement rules、authorization 或 source-of-truth。",
			"保留 domain-specific validation 與 project evidence boundary。",
		}
	default:
		return []string{
			"先確認 registry record 的 model profile。",
			"依 `models/profiles/README.md` 與 `models/compression/README.md` 選讀取深度。",
		}
	}
}

func runtimeLinksFromCell(cell string) []string {
	matches := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`).FindAllStringSubmatch(cell, -1)
	links := []string{}
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		target := strings.SplitN(match[1], "#", 2)[0]
		if target != "" {
			links = append(links, target)
		}
	}
	return links
}

func runtimeRepoRelativeLink(repo string, basePath string, links []string) string {
	if len(links) == 0 || links[0] == "" {
		return ""
	}
	target := links[0]
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return target
	}
	resolved := filepath.Clean(filepath.Join(filepath.Dir(basePath), filepath.FromSlash(target)))
	relative, err := filepath.Rel(repo, resolved)
	if err != nil {
		return ""
	}
	return filepath.ToSlash(relative)
}

func runtimeStripMarkup(value string) string {
	value = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).ReplaceAllString(value, "$1")
	return strings.ReplaceAll(value, "`", "")
}

func runtimeLayerFor(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func rubyLikeString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	default:
		return rubyInspect(typed)
	}
}

func rubyInspect(value any) string {
	switch typed := value.(type) {
	case nil:
		return "nil"
	case string:
		return `"` + strings.ReplaceAll(typed, `"`, `\"`) + `"`
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case int, int64, float64:
		return fmt.Sprint(typed)
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			parts = append(parts, rubyInspect(item))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		keys := orderedRuntimeMetadataKeys(typed)
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, rubyInspect(key)+"=>"+rubyInspect(typed[key]))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case map[any]any:
		normalized := map[string]any{}
		for key, item := range typed {
			normalized[fmt.Sprint(key)] = item
		}
		return rubyInspect(normalized)
	default:
		return fmt.Sprint(typed)
	}
}

func orderedRuntimeMetadataKeys(values map[string]any) []string {
	preferred := []string{
		"estimated_tokens",
		"load_strategy",
		"cacheable",
		"read_frequency",
		"invalidation_triggers",
		"ttl",
		"provider_cache",
		"provider_cache_candidate",
		"prefix_stability",
		"cache_position",
		"churn_risk",
		"notes",
	}
	keys := []string{}
	seen := map[string]bool{}
	for _, key := range preferred {
		if _, ok := values[key]; ok {
			keys = append(keys, key)
			seen[key] = true
		}
	}
	remaining := []string{}
	for key := range values {
		if !seen[key] {
			remaining = append(remaining, key)
		}
	}
	sort.Strings(remaining)
	return append(keys, remaining...)
}

func nonEmptyStrings(values ...string) []string {
	result := []string{}
	for _, value := range values {
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func firstRuntimeLineWithPrefix(content string, prefix string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}

func firstRuntimeHeadingAfter(content string, marker string) string {
	lines := strings.Split(content, "\n")
	for index, line := range lines {
		if strings.TrimSpace(line) != marker {
			continue
		}
		for _, candidate := range lines[index+1:] {
			value := strings.TrimSpace(candidate)
			if value == "" || strings.HasPrefix(value, "#") {
				continue
			}
			return value
		}
	}
	return ""
}

func runtimeTitleFromMarkdown(content string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(regexp.MustCompile(`^#+\s*`).ReplaceAllString(line, ""))
		}
	}
	return ""
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
