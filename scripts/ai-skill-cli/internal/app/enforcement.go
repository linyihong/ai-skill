package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// runEnforcement dispatches `ai-skill enforcement <subcommand>`. Phase 4 of
// plans/active/2026-05-31-2100-mechanical-enforcement-registry.md exposes
// the existing LintEnforcementRegistry engine and a coverage aggregation
// as public CLI surfaces.
func runEnforcement(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill enforcement <lint|coverage> [flags]")
		return ExitInvalidUsage
	}
	switch args[0] {
	case "lint":
		return runEnforcementLint(args[1:], stdout, stderr)
	case "coverage":
		return runEnforcementCoverage(args[1:], stdout, stderr)
	case "transition-check":
		return runEnforcementTransitionCheck(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		_, _ = fmt.Fprintln(stdout, "usage: ai-skill enforcement <lint|coverage|transition-check> [flags]")
		_, _ = fmt.Fprintln(stdout, "")
		_, _ = fmt.Fprintln(stdout, "subcommands:")
		_, _ = fmt.Fprintln(stdout, "  lint                run the enforcement-registry lint engine and report findings")
		_, _ = fmt.Fprintln(stdout, "  coverage            aggregate rule_class coverage status + verification + runtime observation")
		_, _ = fmt.Fprintln(stdout, "  transition-check    Phase 4.5 R1/R2/R3 — detect rule_class coverage transitions in staged registry diff and block missing-ADR demotions / unverified promotions")
		return ExitSuccess
	default:
		_, _ = fmt.Fprintf(stderr, "unknown enforcement subcommand: %s\n", args[0])
		return ExitInvalidUsage
	}
}

// ─────────────────────────────────────────────────────────────────────
// `ai-skill enforcement lint` — thin CLI wrapper around
// LintEnforcementRegistry. Reuses Phase 3 engine; adds arg parsing,
// severity-grouped output, optional --check filter, and assertion-mode
// --expect-finding for cross-platform scenario detection_commands.
// ─────────────────────────────────────────────────────────────────────

type enforcementLintOptions struct {
	repo           string
	registry       string
	check          string
	expectFinding  string
	expectSeverity string
	jsonOutput     bool
	plainOutput    bool
}

func runEnforcementLint(args []string, stdout io.Writer, stderr io.Writer) int {
	opts := enforcementLintOptions{}
	fs := newFlagSet("enforcement lint", stderr)
	fs.StringVar(&opts.repo, "repo", ".", "Ai-skill repository root used to resolve rule yaml + executor files")
	fs.StringVar(&opts.registry, "registry", "", "override path to enforcement-registry.yaml (default: <repo>/enforcement/enforcement-registry.yaml)")
	fs.StringVar(&opts.check, "check", "", "filter findings by lint type (substring match on Type field)")
	fs.StringVar(&opts.expectFinding, "expect-finding", "", "assertion mode: exit 0 if at least one finding's Type or any field value contains this substring, exit 30 otherwise")
	fs.StringVar(&opts.expectSeverity, "expect-severity", "", "assertion mode: narrow --expect-finding match to findings of this severity (FAIL or WARNING)")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable text output (default)")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	root, err := resolveEnforcementRepo(opts.repo)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "resolve repo: %v\n", err)
		return ExitInvalidUsage
	}

	errs, lintErr := lintWithOverride(root, opts.registry)
	if lintErr != nil {
		_, _ = fmt.Fprintf(stderr, "load enforcement-registry: %v\n", lintErr)
		return ExitValidationFailed
	}

	filtered := errs
	if strings.TrimSpace(opts.check) != "" {
		filtered = filterFindings(errs, opts.check)
	}

	// Assertion mode: exit 0 if any matching finding present, 30 otherwise.
	if strings.TrimSpace(opts.expectFinding) != "" {
		matched := assertionMatch(filtered, opts.expectFinding, opts.expectSeverity)
		result := buildLintAssertionResult(opts, root, filtered, matched)
		if opts.jsonOutput {
			if err := writeJSON(stdout, result); err != nil {
				_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
				return ExitGeneralFailure
			}
		} else {
			if err := writePlain(stdout, result); err != nil {
				_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
				return ExitGeneralFailure
			}
		}
		return result.ExitCode
	}

	// Regular mode: render all findings, exit 30 if any FAIL, 0 otherwise.
	result := buildLintReportResult(opts, root, filtered)
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

// lintWithOverride lets `--registry <path>` swap in a synthetic registry
// while still resolving source files / hooks.go from <root>. Implemented by
// temporarily writing the override file into a symlink-free shadow path is
// brittle on Windows; simpler: directly read the override yaml, then call
// the standard lint with a per-invocation override via a temp file copy.
// The Phase 3 engine reads the registry at <repo>/enforcement/enforcement-registry.yaml,
// so when override is set we treat <override-dir> as a synthetic repo only if
// it contains a fixture layout — otherwise we copy the override into a temp
// repo whose enforcement/enforcement-registry.yaml points at the override.
func lintWithOverride(repo, registry string) ([]EnforcementRegistryLintError, error) {
	if strings.TrimSpace(registry) == "" {
		return LintEnforcementRegistry(repo)
	}
	abs, err := filepath.Abs(registry)
	if err != nil {
		return nil, fmt.Errorf("resolve --registry: %w", err)
	}
	if _, err := os.Stat(abs); err != nil {
		return nil, fmt.Errorf("--registry not readable: %w", err)
	}
	// Build a shadow repo whose enforcement/enforcement-registry.yaml is the
	// override. Other dirs (enforcement/runtime/governance for orphan walk,
	// scripts/ai-skill-cli/internal/app/hooks.go for orphan executors) are
	// hard-linked back to the original repo via a sibling overlay directory.
	shadow, err := os.MkdirTemp("", "ai-skill-enforcement-lint-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(shadow)
	if err := os.MkdirAll(filepath.Join(shadow, "enforcement"), 0o755); err != nil {
		return nil, err
	}
	overrideBytes, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}
	dstReg := filepath.Join(shadow, "enforcement", "enforcement-registry.yaml")
	if err := os.WriteFile(dstReg, overrideBytes, 0o644); err != nil {
		return nil, err
	}
	// Mirror the original repo's scanned trees into the shadow so the orphan
	// walks see them. Use junction-free copies of yaml-only content (small).
	// hooks.go is needed for orphan_executor; we copy it too.
	mirrorPaths := []string{"enforcement", "runtime", "governance"}
	for _, p := range mirrorPaths {
		src := filepath.Join(repo, p)
		dst := filepath.Join(shadow, p)
		_ = filepath.Walk(src, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil || info == nil {
				return nil
			}
			rel, err := filepath.Rel(src, path)
			if err != nil {
				return nil
			}
			out := filepath.Join(dst, rel)
			if info.IsDir() {
				return os.MkdirAll(out, 0o755)
			}
			// Skip the registry file itself; we already wrote the override.
			if filepath.ToSlash(filepath.Join(p, rel)) == "enforcement/enforcement-registry.yaml" {
				return nil
			}
			if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			_ = os.MkdirAll(filepath.Dir(out), 0o755)
			return os.WriteFile(out, data, 0o644)
		})
	}
	hooksRel := "scripts/ai-skill-cli/internal/app/hooks.go"
	hooksSrc := filepath.Join(repo, filepath.FromSlash(hooksRel))
	if data, err := os.ReadFile(hooksSrc); err == nil {
		hooksDst := filepath.Join(shadow, filepath.FromSlash(hooksRel))
		_ = os.MkdirAll(filepath.Dir(hooksDst), 0o755)
		_ = os.WriteFile(hooksDst, data, 0o644)
	}
	return LintEnforcementRegistry(shadow)
}

func filterFindings(errs []EnforcementRegistryLintError, substr string) []EnforcementRegistryLintError {
	s := strings.ToLower(substr)
	out := errs[:0:0]
	for _, e := range errs {
		if strings.Contains(strings.ToLower(e.Type), s) {
			out = append(out, e)
		}
	}
	return out
}

func assertionMatch(errs []EnforcementRegistryLintError, expect, severity string) []EnforcementRegistryLintError {
	want := strings.ToLower(strings.TrimSpace(expect))
	sev := strings.ToUpper(strings.TrimSpace(severity))
	var matched []EnforcementRegistryLintError
	for _, e := range errs {
		if sev != "" {
			actual := e.Severity
			if actual == "" {
				actual = SeverityFail
			}
			if actual != sev {
				continue
			}
		}
		hit := strings.Contains(strings.ToLower(e.Type), want)
		if !hit {
			hit = strings.Contains(strings.ToLower(e.Message), want)
		}
		if !hit {
			for _, f := range e.Fields {
				if strings.Contains(strings.ToLower(f.Value), want) || strings.Contains(strings.ToLower(f.Key), want) {
					hit = true
					break
				}
			}
		}
		if hit {
			matched = append(matched, e)
		}
	}
	return matched
}

func buildLintAssertionResult(opts enforcementLintOptions, repo string, all []EnforcementRegistryLintError, matched []EnforcementRegistryLintError) Result {
	result := Result{
		Command:        "enforcement lint",
		Mode:           "assert",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	result.Checks = append(result.Checks, Check{Name: "repo_root", Status: "ok", Message: repo})
	if opts.registry != "" {
		result.Checks = append(result.Checks, Check{Name: "registry_override", Status: "ok", Message: opts.registry})
	}
	if opts.check != "" {
		result.Checks = append(result.Checks, Check{Name: "check_filter", Status: "ok", Message: opts.check})
	}
	result.Checks = append(result.Checks, Check{Name: "expect_finding", Status: "ok", Message: opts.expectFinding})
	result.Checks = append(result.Checks, Check{Name: "findings_total", Status: "ok", Message: fmt.Sprintf("%d", len(all))})
	if len(matched) > 0 {
		result.Checks = append(result.Checks, Check{Name: "assertion", Status: "ok", Message: fmt.Sprintf("matched %d finding(s)", len(matched))})
		return result
	}
	result.Status = "failed"
	result.ExitCode = ExitValidationFailed
	result.Checks = append(result.Checks, Check{Name: "assertion", Status: "failed", Message: "no finding matched --expect-finding substring"})
	result.Error = &CommandError{
		Code:        "enforcement_lint_assertion_unmet",
		Message:     fmt.Sprintf("--expect-finding %q matched 0 findings (of %d total)", opts.expectFinding, len(all)),
		Remediation: "Adjust --expect-finding substring, run without --expect-finding to see full lint output, or fix the underlying rule to produce the expected finding.",
	}
	return result
}

func buildLintReportResult(opts enforcementLintOptions, repo string, errs []EnforcementRegistryLintError) Result {
	result := Result{
		Command:        "enforcement lint",
		Mode:           "report",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	result.Checks = append(result.Checks, Check{Name: "repo_root", Status: "ok", Message: repo})
	if opts.registry != "" {
		result.Checks = append(result.Checks, Check{Name: "registry_override", Status: "ok", Message: opts.registry})
	}
	if opts.check != "" {
		result.Checks = append(result.Checks, Check{Name: "check_filter", Status: "ok", Message: opts.check})
	}
	var fails, warns []EnforcementRegistryLintError
	for _, e := range errs {
		if e.IsFail() {
			fails = append(fails, e)
		} else {
			warns = append(warns, e)
		}
	}
	result.Checks = append(result.Checks, Check{Name: "findings_fail", Status: "ok", Message: fmt.Sprintf("%d", len(fails))})
	result.Checks = append(result.Checks, Check{Name: "findings_warn", Status: "ok", Message: fmt.Sprintf("%d", len(warns))})
	// Inline each finding as a Check so writePlain emits the canonical block.
	for _, e := range errs {
		status := "warning"
		if e.IsFail() {
			status = "failed"
		}
		result.Checks = append(result.Checks, Check{
			Name:    "lint." + e.Type,
			Status:  status,
			Message: lintCheckMessage(e),
		})
	}
	if len(fails) > 0 {
		result.Status = "failed"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{
			Code:        "enforcement_lint_failed",
			Message:     fmt.Sprintf("%d FAIL finding(s) in enforcement-registry.yaml", len(fails)),
			Remediation: "Resolve FAIL findings (see lint.* checks above). WARNING findings are governance signals and do not block.",
		}
	}
	return result
}

func lintCheckMessage(e EnforcementRegistryLintError) string {
	var b strings.Builder
	for i, f := range e.Fields {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s=%s", f.Key, f.Value)
	}
	if b.Len() == 0 {
		return e.Message
	}
	return b.String() + " — " + e.Message
}

func resolveEnforcementRepo(repoArg string) (string, error) {
	candidate := strings.TrimSpace(repoArg)
	if candidate == "" {
		candidate = "."
	}
	abs, err := filepath.Abs(candidate)
	if err != nil {
		return "", err
	}
	if isAiSkillRuntimeRepo(abs) {
		return abs, nil
	}
	if env := strings.TrimSpace(os.Getenv("AI_SKILL_REPO")); env != "" {
		if isAiSkillRuntimeRepo(env) {
			return env, nil
		}
	}
	// Soft fallback: even when the path doesn't look like Ai-skill canonical
	// repo (e.g. a fixture root), still accept it so scenarios with
	// --registry overrides can run.
	return abs, nil
}

// ─────────────────────────────────────────────────────────────────────
// `ai-skill enforcement coverage` — Phase 4 primary deliverable.
// Aggregates registry into 6 coverage buckets, computes verification
// level per class, surfaces runtime_observed gaps (or null + alert when
// runtime.executor_observations is not yet wired), and renders the
// report as text / json / markdown.
// ─────────────────────────────────────────────────────────────────────

type enforcementCoverageOptions struct {
	repo       string
	registry   string
	format     string
	diff       string
	detail     bool
	selfCheck  bool
	jsonStdout bool
}

func runEnforcementCoverage(args []string, stdout io.Writer, stderr io.Writer) int {
	opts := enforcementCoverageOptions{format: "text"}
	fs := newFlagSet("enforcement coverage", stderr)
	fs.StringVar(&opts.repo, "repo", ".", "Ai-skill repository root")
	fs.StringVar(&opts.registry, "registry", "", "override path to enforcement-registry.yaml")
	fs.StringVar(&opts.format, "format", "text", "output format: text | json | markdown")
	fs.StringVar(&opts.diff, "diff", "", "compare current registry to git revision (e.g. origin/main); prints added/removed/changed classes")
	fs.BoolVar(&opts.detail, "detail", false, "include per-class detail in text/markdown output (default: summary only)")
	fs.BoolVar(&opts.selfCheck, "self-check", false, "render all 3 formats internally and validate schema; cross-platform replacement for jq/grep scenarios")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	switch opts.format {
	case "text", "json", "markdown":
	default:
		_, _ = fmt.Fprintf(stderr, "unknown --format: %s (want text|json|markdown)\n", opts.format)
		return ExitInvalidUsage
	}

	root, err := resolveEnforcementRepo(opts.repo)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "resolve repo: %v\n", err)
		return ExitInvalidUsage
	}

	registryPath := opts.registry
	if strings.TrimSpace(registryPath) == "" {
		registryPath = filepath.Join(root, "enforcement", "enforcement-registry.yaml")
	}
	snap, err := loadRegistrySnapshotFromPath(registryPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "load registry: %v\n", err)
		return ExitValidationFailed
	}
	report := buildCoverageReport(root, snap)

	if opts.diff != "" {
		diffSection, derr := buildCoverageDiff(root, snap, opts.diff)
		if derr != nil {
			report.Alerts = append(report.Alerts, coverageAlert{
				Kind:    "diff_unavailable",
				Message: fmt.Sprintf("--diff %s could not be computed: %v", opts.diff, derr),
			})
		} else {
			report.Diff = diffSection
		}
	}

	if opts.selfCheck {
		code := runCoverageSelfCheck(report, stdout, stderr)
		return code
	}

	switch opts.format {
	case "text":
		_, _ = fmt.Fprint(stdout, renderCoverageText(report, opts.detail))
	case "markdown":
		_, _ = fmt.Fprint(stdout, renderCoverageMarkdown(report, opts.detail))
	case "json":
		out, _ := json.MarshalIndent(renderCoverageJSON(report), "", "  ")
		_, _ = fmt.Fprintln(stdout, string(out))
	}
	return ExitSuccess
}

func loadRegistrySnapshotFromPath(path string) (*registrySnapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var snap registrySnapshot
	if err := yaml.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &snap, nil
}

// ─────────────────────────────────────────────────────────────────────
// Coverage report data model
// ─────────────────────────────────────────────────────────────────────

type coverageReport struct {
	GeneratedAt          time.Time
	TotalRuleClasses     int
	Buckets              map[string]int      // 6-value coverage enum → count
	Classes              []coverageClassRow
	Alerts               []coverageAlert
	RuntimeMetricsWired  bool
	ObservationWindowDays int
	Diff                 *coverageDiff
}

type coverageClassRow struct {
	ID                 string
	Coverage           string
	VerificationLevel  string
	RuntimeObservedPct *float64 // nil = not wired yet
	ChildPlan          string
	ReplacedBy         string
	RemovalDate        string
	RevisitWhen        string
	ResearchQuestions  []string
	UpstreamClasses    []string
	Alerts             []string
}

type coverageAlert struct {
	Kind     string `json:"kind"`
	RuleClass string `json:"rule_class,omitempty"`
	Message  string `json:"message"`
}

type coverageDiff struct {
	BaseRef string
	Added   []string
	Removed []string
	Changed []coverageDiffChange
}

type coverageDiffChange struct {
	ID       string
	FromCov  string
	ToCov    string
}

// The canonical 6-value bucket order. Buckets that don't appear in the
// registry still render with count 0 to satisfy the JSON contract.
var coverageBucketOrder = []string{
	"mechanical",
	"behavioral_only",
	"not_mechanizable",
	"pending_implementation",
	"research_required",
	"deprecated",
}

func buildCoverageReport(repo string, snap *registrySnapshot) *coverageReport {
	report := &coverageReport{
		GeneratedAt:          time.Now().UTC(),
		Buckets:              map[string]int{},
		Classes:              []coverageClassRow{},
		Alerts:               []coverageAlert{},
		ObservationWindowDays: 30,
	}
	for _, b := range coverageBucketOrder {
		report.Buckets[b] = 0
	}
	// Detect runtime metrics wiring: runtime.db has executor_observations table.
	report.RuntimeMetricsWired = detectExecutorObservationsTable(filepath.Join(repo, "runtime", "runtime.db"))
	if !report.RuntimeMetricsWired {
		report.Alerts = append(report.Alerts, coverageAlert{
			Kind:    "runtime_observations_not_wired",
			Message: "runtime.executor_observations table absent; runtime_observed_pct reported as null until Phase 5 wires the write path.",
		})
	}

	for _, rc := range snap.RuleClasses {
		report.TotalRuleClasses++
		bucket := rc.Coverage
		if _, known := report.Buckets[bucket]; !known {
			// Unknown coverage value — count under bucket but also alert.
			report.Buckets[bucket] = 0
			report.Alerts = append(report.Alerts, coverageAlert{
				Kind:      "unknown_coverage_value",
				RuleClass: rc.ID,
				Message:   fmt.Sprintf("rule_class %s has coverage=%q which is not in the 6-value enum", rc.ID, rc.Coverage),
			})
		}
		report.Buckets[bucket]++

		row := coverageClassRow{
			ID:                rc.ID,
			Coverage:          rc.Coverage,
			VerificationLevel: classifyVerification(repo, rc),
			ChildPlan:         rc.ChildPlan,
			ReplacedBy:        rc.ReplacedBy,
			RemovalDate:       rc.RemovalDate,
			ResearchQuestions: append([]string(nil), rc.ResearchQuestions...),
			UpstreamClasses:   append([]string(nil), rc.UpstreamClasses...),
		}
		if rc.SunsetDecision != nil {
			row.RevisitWhen = rc.SunsetDecision.RevisitWhen
		}
		// Runtime observed: null when metrics not wired (per user decision
		// for Phase 4). When wired (Phase 5), this would aggregate from
		// runtime.executor_observations.
		if report.RuntimeMetricsWired {
			pct := observedPctForClass(repo, rc.ID, report.ObservationWindowDays)
			row.RuntimeObservedPct = &pct
		}
		// Per-class alerts.
		today := time.Now().UTC().Format("2006-01-02")
		// Phase 4.5 R4 — deprecated past removal_date by 30+ days
		// (governance alert, NOT compile fail per registry self-governance
		// design: R4-R5 surface here for human review, do not block).
		if rc.Coverage == "deprecated" && rc.RemovalDate != "" {
			if rc.RemovalDate < today {
				row.Alerts = append(row.Alerts, "past_removal_date")
				report.Alerts = append(report.Alerts, coverageAlert{
					Kind:      "deprecated_past_removal_date",
					RuleClass: rc.ID,
					Message:   fmt.Sprintf("removal_date %s elapsed (today %s); actually remove or extend.", rc.RemovalDate, today),
				})
			}
			if pastByDays(rc.RemovalDate, 30) {
				row.Alerts = append(row.Alerts, "R4_governance_overdue")
				report.Alerts = append(report.Alerts, coverageAlert{
					Kind:      "R4_deprecated_governance_overdue",
					RuleClass: rc.ID,
					Message:   fmt.Sprintf("deprecated past removal_date by ≥ 30 days (removal_date=%s); governance decision required: actually remove or extend with new rationale.", rc.RemovalDate),
				})
			}
		}
		// Phase 4.5 R5 — research_required past estimated_unblock_timeline.
		if rc.Coverage == "research_required" && strings.TrimSpace(rc.UnblockTimeline) != "" {
			if isCalendarPast(rc.UnblockTimeline, today) {
				row.Alerts = append(row.Alerts, "R5_research_timeout")
				report.Alerts = append(report.Alerts, coverageAlert{
					Kind:      "R5_research_required_timeout",
					RuleClass: rc.ID,
					Message:   fmt.Sprintf("research_required past estimated_unblock_timeline=%s (today %s); promote to pending_implementation or demote with rationale.", rc.UnblockTimeline, today),
				})
			}
		}
		if rc.Coverage == "mechanical" && row.VerificationLevel == "symbol_only" {
			row.Alerts = append(row.Alerts, "scenarios_missing")
		}
		report.Classes = append(report.Classes, row)
	}
	sort.Slice(report.Classes, func(i, j int) bool { return report.Classes[i].ID < report.Classes[j].ID })
	return report
}

// classifyVerification returns one of:
//   full / symbol_only / planned / behavioral / not_applicable
// Heuristic, not a strict mapping: full = mechanical + at least one
// validation scenario file exists for the class; symbol_only = mechanical
// without scenario coverage; planned = pending_implementation; behavioral
// = behavioral_only; not_applicable = not_mechanizable / deprecated.
func classifyVerification(repo string, rc registryRuleClass) string {
	switch rc.Coverage {
	case "not_mechanizable":
		return "not_applicable"
	case "deprecated":
		return "not_applicable"
	case "behavioral_only":
		return "behavioral"
	case "pending_implementation", "research_required":
		return "planned"
	case "mechanical":
		if scenarioExistsForClass(repo, rc.ID) {
			return "full"
		}
		return "symbol_only"
	}
	return "unknown"
}

func scenarioExistsForClass(repo, id string) bool {
	// Cheap heuristic: walk validation/scenarios/ and look for filenames
	// that contain the class id (snake-or-kebab tolerant).
	root := filepath.Join(repo, "validation", "scenarios")
	if _, err := os.Stat(root); err != nil {
		return false
	}
	want := strings.ReplaceAll(id, "_", "-")
	found := false
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		base := strings.ToLower(filepath.Base(p))
		if !strings.HasSuffix(base, ".yaml") && !strings.HasSuffix(base, ".yml") {
			return nil
		}
		if strings.Contains(base, want) || strings.Contains(base, id) {
			found = true
			return filepath.SkipDir
		}
		return nil
	})
	return found
}

func detectExecutorObservationsTable(dbPath string) bool {
	if _, err := os.Stat(dbPath); err != nil {
		return false
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return false
	}
	defer db.Close()
	var name string
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='executor_observations'")
	if err := row.Scan(&name); err != nil {
		return false
	}
	return name == "executor_observations"
}

func observedPctForClass(repo, classID string, windowDays int) float64 {
	dbPath := filepath.Join(repo, "runtime", "runtime.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0
	}
	defer db.Close()
	cutoff := time.Now().UTC().AddDate(0, 0, -windowDays).Format(time.RFC3339)
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM executor_observations WHERE rule_class_id = ? AND fired_at >= ?", classID, cutoff)
	if err := row.Scan(&count); err != nil {
		return 0
	}
	if count > 0 {
		return 100.0
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────
// --diff <ref>
// ─────────────────────────────────────────────────────────────────────

func buildCoverageDiff(repo string, current *registrySnapshot, ref string) (*coverageDiff, error) {
	cmd := exec.Command("git", "-C", repo, "show", ref+":enforcement/enforcement-registry.yaml")
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git show %s:enforcement/enforcement-registry.yaml: %w", ref, err)
	}
	var base registrySnapshot
	if err := yaml.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("parse base registry: %w", err)
	}
	baseByID := map[string]string{}
	for _, rc := range base.RuleClasses {
		baseByID[rc.ID] = rc.Coverage
	}
	currByID := map[string]string{}
	for _, rc := range current.RuleClasses {
		currByID[rc.ID] = rc.Coverage
	}
	diff := &coverageDiff{BaseRef: ref}
	for id, cov := range currByID {
		if _, ok := baseByID[id]; !ok {
			diff.Added = append(diff.Added, id)
			continue
		}
		if baseByID[id] != cov {
			diff.Changed = append(diff.Changed, coverageDiffChange{ID: id, FromCov: baseByID[id], ToCov: cov})
		}
	}
	for id := range baseByID {
		if _, ok := currByID[id]; !ok {
			diff.Removed = append(diff.Removed, id)
		}
	}
	sort.Strings(diff.Added)
	sort.Strings(diff.Removed)
	sort.Slice(diff.Changed, func(i, j int) bool { return diff.Changed[i].ID < diff.Changed[j].ID })
	return diff, nil
}

// ─────────────────────────────────────────────────────────────────────
// Rendering: text / markdown / json
// ─────────────────────────────────────────────────────────────────────

func renderCoverageText(r *coverageReport, detail bool) string {
	var b strings.Builder
	date := r.GeneratedAt.Format("2006-01-02")
	fmt.Fprintf(&b, "Enforcement Coverage Report (%s)\n", date)
	// Plan §Phase 4 sample uses U+2550 box-drawing. We mirror it.
	b.WriteString(strings.Repeat("═", 39))
	b.WriteString("\n")
	fmt.Fprintf(&b, "Total Rule Classes: %d\n\n", r.TotalRuleClasses)

	for _, name := range coverageBucketOrder {
		count := r.Buckets[name]
		pct := 0
		if r.TotalRuleClasses > 0 {
			pct = (count*100 + r.TotalRuleClasses/2) / r.TotalRuleClasses
		}
		label := fmt.Sprintf("%-23s", coverageBucketLabel(name))
		fmt.Fprintf(&b, "  %s %3d  (%2d%%)\n", label, count, pct)
	}
	b.WriteString("\n")

	if detail {
		fmt.Fprintf(&b, "Per-class detail:\n\n")
		fmt.Fprintf(&b, "  %-32s %-22s %-18s %s\n", "Rule Class", "Status", "Verification", fmt.Sprintf("Runtime (%dd)", r.ObservationWindowDays))
		b.WriteString("  ")
		b.WriteString(strings.Repeat("─", 90))
		b.WriteString("\n")
		for _, c := range r.Classes {
			runtimeCol := "—"
			if c.RuntimeObservedPct != nil {
				runtimeCol = fmt.Sprintf("%.0f%%", *c.RuntimeObservedPct)
			} else {
				runtimeCol = "n/a (not wired)"
			}
			fmt.Fprintf(&b, "  %-32s %-22s %-18s %s\n", truncate(c.ID, 32), truncate(c.Coverage, 22), truncate(c.VerificationLevel, 18), runtimeCol)
		}
		b.WriteString("\n")
	}

	emitSection := func(title string, want string, withExtra func(coverageClassRow) string) {
		var rows []coverageClassRow
		for _, c := range r.Classes {
			if c.Coverage == want {
				rows = append(rows, c)
			}
		}
		if len(rows) == 0 {
			return
		}
		fmt.Fprintf(&b, "%s:\n", title)
		for _, c := range rows {
			extra := ""
			if withExtra != nil {
				extra = withExtra(c)
			}
			if extra == "" {
				fmt.Fprintf(&b, "  %s\n", c.ID)
			} else {
				fmt.Fprintf(&b, "  %-32s — %s\n", c.ID, extra)
			}
		}
		b.WriteString("\n")
	}
	emitSection("Pending implementation (active child plans)", "pending_implementation", func(c coverageClassRow) string {
		if c.ChildPlan == "" {
			return "child_plan TBD"
		}
		return c.ChildPlan
	})
	emitSection("Research required (no clear mechanization path)", "research_required", func(c coverageClassRow) string {
		if len(c.ResearchQuestions) == 0 {
			return ""
		}
		return strings.Join(c.ResearchQuestions, "; ")
	})
	emitSection("Behavioral_only awaiting sunset review", "behavioral_only", func(c coverageClassRow) string {
		return c.RevisitWhen
	})
	emitSection("Not_mechanizable (closed, will not appear in review queue)", "not_mechanizable", nil)
	emitSection("Deprecated (awaiting removal)", "deprecated", func(c coverageClassRow) string {
		if c.ReplacedBy != "" {
			return "replaced_by: " + c.ReplacedBy
		}
		if c.RemovalDate != "" {
			return "removal_date: " + c.RemovalDate
		}
		return ""
	})

	if len(r.Alerts) > 0 {
		gov, other := partitionGovernanceAlerts(r.Alerts)
		if len(gov) > 0 {
			b.WriteString("⚠ Governance Alerts (Phase 4.5 R4/R5):\n")
			for _, a := range gov {
				if a.RuleClass != "" {
					fmt.Fprintf(&b, "  [%s] %s: %s\n", a.Kind, a.RuleClass, a.Message)
				} else {
					fmt.Fprintf(&b, "  [%s] %s\n", a.Kind, a.Message)
				}
			}
			b.WriteString("\n")
		}
		if len(other) > 0 {
			b.WriteString("⚠ Alerts:\n")
			for _, a := range other {
				if a.RuleClass != "" {
					fmt.Fprintf(&b, "  [%s] %s: %s\n", a.Kind, a.RuleClass, a.Message)
				} else {
					fmt.Fprintf(&b, "  [%s] %s\n", a.Kind, a.Message)
				}
			}
			b.WriteString("\n")
		}
	}
	if r.Diff != nil {
		fmt.Fprintf(&b, "Changes vs %s:\n", r.Diff.BaseRef)
		for _, id := range r.Diff.Added {
			fmt.Fprintf(&b, "  + %s\n", id)
		}
		for _, id := range r.Diff.Removed {
			fmt.Fprintf(&b, "  - %s\n", id)
		}
		for _, c := range r.Diff.Changed {
			fmt.Fprintf(&b, "  ~ %s: %s → %s\n", c.ID, c.FromCov, c.ToCov)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func coverageBucketLabel(name string) string {
	switch name {
	case "mechanical":
		return "Mechanical:"
	case "behavioral_only":
		return "Behavioral only:"
	case "not_mechanizable":
		return "Not mechanizable:"
	case "pending_implementation":
		return "Pending impl:"
	case "research_required":
		return "Research required:"
	case "deprecated":
		return "Deprecated:"
	}
	return name + ":"
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

func renderCoverageMarkdown(r *coverageReport, detail bool) string {
	var b strings.Builder
	b.WriteString("# Enforcement Coverage Report\n\n")
	fmt.Fprintf(&b, "_Generated: %s UTC — Total Rule Classes: %d_\n\n", r.GeneratedAt.Format(time.RFC3339), r.TotalRuleClasses)
	b.WriteString("## Summary\n\n")
	b.WriteString("| Bucket | Count | % |\n")
	b.WriteString("| --- | ---: | ---: |\n")
	for _, name := range coverageBucketOrder {
		count := r.Buckets[name]
		pct := 0
		if r.TotalRuleClasses > 0 {
			pct = (count*100 + r.TotalRuleClasses/2) / r.TotalRuleClasses
		}
		fmt.Fprintf(&b, "| %s | %d | %d%% |\n", name, count, pct)
	}
	b.WriteString("\n")

	if detail {
		b.WriteString("## Per-class detail\n\n")
		b.WriteString("| Rule Class | Coverage | Verification | Runtime (window) |\n")
		b.WriteString("| --- | --- | --- | --- |\n")
		for _, c := range r.Classes {
			runtimeCol := "n/a (not wired)"
			if c.RuntimeObservedPct != nil {
				runtimeCol = fmt.Sprintf("%.0f%%", *c.RuntimeObservedPct)
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", c.ID, c.Coverage, c.VerificationLevel, runtimeCol)
		}
		b.WriteString("\n")
	}

	if len(r.Alerts) > 0 {
		gov, other := partitionGovernanceAlerts(r.Alerts)
		if len(gov) > 0 {
			b.WriteString("## Governance Alerts (Phase 4.5 R4/R5)\n\n")
			for _, a := range gov {
				if a.RuleClass != "" {
					fmt.Fprintf(&b, "- **%s** (`%s`): %s\n", a.Kind, a.RuleClass, a.Message)
				} else {
					fmt.Fprintf(&b, "- **%s**: %s\n", a.Kind, a.Message)
				}
			}
			b.WriteString("\n")
		}
		if len(other) > 0 {
			b.WriteString("## Alerts\n\n")
			for _, a := range other {
				if a.RuleClass != "" {
					fmt.Fprintf(&b, "- **%s** (`%s`): %s\n", a.Kind, a.RuleClass, a.Message)
				} else {
					fmt.Fprintf(&b, "- **%s**: %s\n", a.Kind, a.Message)
				}
			}
			b.WriteString("\n")
		}
	}
	if r.Diff != nil {
		fmt.Fprintf(&b, "## Changes vs %s\n\n", r.Diff.BaseRef)
		if len(r.Diff.Added) > 0 {
			b.WriteString("- **Added**: " + strings.Join(r.Diff.Added, ", ") + "\n")
		}
		if len(r.Diff.Removed) > 0 {
			b.WriteString("- **Removed**: " + strings.Join(r.Diff.Removed, ", ") + "\n")
		}
		for _, c := range r.Diff.Changed {
			fmt.Fprintf(&b, "- **Changed**: %s: %s → %s\n", c.ID, c.FromCov, c.ToCov)
		}
		b.WriteString("\n")
	}
	return b.String()
}

type coverageJSONOut struct {
	SchemaVersion     int                 `json:"schema_version"`
	GeneratedAt       string              `json:"generated_at"`
	TotalRuleClasses  int                 `json:"total_rule_classes"`
	ObservationWindow int                 `json:"observation_window_days"`
	Buckets           map[string]int      `json:"buckets"`
	PerClass          []coverageJSONClass `json:"per_class"`
	Alerts            []coverageAlert     `json:"alerts"`
	Diff              *coverageJSONDiff   `json:"diff,omitempty"`
}

type coverageJSONClass struct {
	ID                  string   `json:"id"`
	Coverage            string   `json:"coverage"`
	Verification        string   `json:"verification"`
	RuntimeObservedPct  *float64 `json:"runtime_observed_pct"`
	Alerts              []string `json:"alerts"`
	ChildPlan           string   `json:"child_plan,omitempty"`
	ReplacedBy          string   `json:"replaced_by,omitempty"`
	RemovalDate         string   `json:"removal_date,omitempty"`
	RevisitWhen         string   `json:"revisit_when,omitempty"`
	ResearchQuestions   []string `json:"research_questions,omitempty"`
	UpstreamClasses     []string `json:"upstream_classes,omitempty"`
}

type coverageJSONDiff struct {
	BaseRef string                  `json:"base_ref"`
	Added   []string                `json:"added"`
	Removed []string                `json:"removed"`
	Changed []coverageJSONDiffChange `json:"changed"`
}

type coverageJSONDiffChange struct {
	ID      string `json:"id"`
	FromCov string `json:"from_coverage"`
	ToCov   string `json:"to_coverage"`
}

func renderCoverageJSON(r *coverageReport) coverageJSONOut {
	out := coverageJSONOut{
		SchemaVersion:     1,
		GeneratedAt:       r.GeneratedAt.Format(time.RFC3339),
		TotalRuleClasses:  r.TotalRuleClasses,
		ObservationWindow: r.ObservationWindowDays,
		Buckets:           map[string]int{},
		PerClass:          []coverageJSONClass{},
		Alerts:            append([]coverageAlert(nil), r.Alerts...),
	}
	for _, b := range coverageBucketOrder {
		out.Buckets[b] = r.Buckets[b]
	}
	for _, c := range r.Classes {
		out.PerClass = append(out.PerClass, coverageJSONClass{
			ID:                 c.ID,
			Coverage:           c.Coverage,
			Verification:       c.VerificationLevel,
			RuntimeObservedPct: c.RuntimeObservedPct,
			Alerts:             append([]string(nil), c.Alerts...),
			ChildPlan:          c.ChildPlan,
			ReplacedBy:         c.ReplacedBy,
			RemovalDate:        c.RemovalDate,
			RevisitWhen:        c.RevisitWhen,
			ResearchQuestions:  c.ResearchQuestions,
			UpstreamClasses:    c.UpstreamClasses,
		})
	}
	if r.Diff != nil {
		d := &coverageJSONDiff{BaseRef: r.Diff.BaseRef, Added: r.Diff.Added, Removed: r.Diff.Removed}
		for _, ch := range r.Diff.Changed {
			d.Changed = append(d.Changed, coverageJSONDiffChange{ID: ch.ID, FromCov: ch.FromCov, ToCov: ch.ToCov})
		}
		out.Diff = d
	}
	return out
}

// runCoverageSelfCheck renders text / json / markdown internally and
// validates each against the scenario contract. Cross-platform replacement
// for `grep | jq -e | head -1` shell pipelines.
func runCoverageSelfCheck(r *coverageReport, stdout io.Writer, stderr io.Writer) int {
	checks := []selfCheckItem{}
	textOut := renderCoverageText(r, false)
	mdOut := renderCoverageMarkdown(r, false)
	jsonOut := renderCoverageJSON(r)

	// text checks
	textLines := strings.Split(textOut, "\n")
	if len(textLines) > 0 && strings.HasPrefix(textLines[0], "Enforcement Coverage Report (") && strings.HasSuffix(textLines[0], ")") {
		checks = append(checks, selfCheckItem{Name: "text_header", OK: true, Message: textLines[0]})
	} else {
		checks = append(checks, selfCheckItem{Name: "text_header", OK: false, Message: "first line does not match 'Enforcement Coverage Report (YYYY-MM-DD)'"})
	}
	if len(textLines) > 1 && strings.HasPrefix(textLines[1], "═") {
		checks = append(checks, selfCheckItem{Name: "text_underline", OK: true, Message: "U+2550 underline present"})
	} else {
		checks = append(checks, selfCheckItem{Name: "text_underline", OK: false, Message: "missing U+2550 underline on line 2"})
	}
	if strings.Contains(textOut, "Total Rule Classes:") {
		checks = append(checks, selfCheckItem{Name: "text_total_line", OK: true, Message: "present"})
	} else {
		checks = append(checks, selfCheckItem{Name: "text_total_line", OK: false, Message: "missing 'Total Rule Classes:' line"})
	}

	// markdown checks
	if strings.HasPrefix(mdOut, "# Enforcement Coverage Report") {
		checks = append(checks, selfCheckItem{Name: "markdown_h1", OK: true, Message: "h1 present"})
	} else {
		checks = append(checks, selfCheckItem{Name: "markdown_h1", OK: false, Message: "missing '# Enforcement Coverage Report' h1"})
	}
	if strings.Contains(mdOut, "## Summary") {
		checks = append(checks, selfCheckItem{Name: "markdown_summary_h2", OK: true, Message: "h2 present"})
	} else {
		checks = append(checks, selfCheckItem{Name: "markdown_summary_h2", OK: false, Message: "missing '## Summary' h2"})
	}

	// json checks
	if jsonOut.SchemaVersion == 1 {
		checks = append(checks, selfCheckItem{Name: "json_schema_version", OK: true, Message: "schema_version=1"})
	} else {
		checks = append(checks, selfCheckItem{Name: "json_schema_version", OK: false, Message: fmt.Sprintf("schema_version=%d, want 1", jsonOut.SchemaVersion)})
	}
	if len(jsonOut.Buckets) == len(coverageBucketOrder) {
		checks = append(checks, selfCheckItem{Name: "json_buckets_complete", OK: true, Message: fmt.Sprintf("%d buckets", len(jsonOut.Buckets))})
	} else {
		checks = append(checks, selfCheckItem{Name: "json_buckets_complete", OK: false, Message: fmt.Sprintf("buckets count=%d, want %d (6-value enum)", len(jsonOut.Buckets), len(coverageBucketOrder))})
	}
	if jsonOut.PerClass != nil {
		checks = append(checks, selfCheckItem{Name: "json_per_class_array", OK: true, Message: fmt.Sprintf("%d entries", len(jsonOut.PerClass))})
	} else {
		checks = append(checks, selfCheckItem{Name: "json_per_class_array", OK: false, Message: "per_class is null, want array"})
	}
	// Snake_case key check via raw marshal.
	raw, _ := json.Marshal(jsonOut)
	if !strings.Contains(string(raw), "schemaVersion") && strings.Contains(string(raw), "schema_version") {
		checks = append(checks, selfCheckItem{Name: "json_snake_case", OK: true, Message: "snake_case keys"})
	} else {
		checks = append(checks, selfCheckItem{Name: "json_snake_case", OK: false, Message: "non-snake_case key detected in JSON output"})
	}

	failed := 0
	for _, c := range checks {
		status := "ok"
		if !c.OK {
			status = "failed"
			failed++
		}
		_, _ = fmt.Fprintf(stdout, "- self_check.%s: %s - %s\n", c.Name, status, c.Message)
	}
	_, _ = fmt.Fprintf(stdout, "self-check summary: %d/%d passed\n", len(checks)-failed, len(checks))
	if failed > 0 {
		_, _ = fmt.Fprintf(stderr, "self-check FAILED (%d/%d)\n", failed, len(checks))
		return ExitValidationFailed
	}
	return ExitSuccess
}

type selfCheckItem struct {
	Name    string
	OK      bool
	Message string
}

// partitionGovernanceAlerts splits the alert list into Phase 4.5 R4/R5
// governance alerts vs everything else, preserving order within each
// bucket. R4/R5 alerts are surfaced under a dedicated `Governance Alerts`
// header so they are not lost in the generic alerts noise.
func partitionGovernanceAlerts(alerts []coverageAlert) (governance, other []coverageAlert) {
	for _, a := range alerts {
		if strings.HasPrefix(a.Kind, "R4_") || strings.HasPrefix(a.Kind, "R5_") {
			governance = append(governance, a)
		} else {
			other = append(other, a)
		}
	}
	return
}

// pastByDays returns true when isoDate (YYYY-MM-DD) is more than `days`
// days before today. Used by Phase 4.5 R4 (deprecated past removal_date
// by 30+ days). Returns false for malformed input — alerts only fire on
// clearly-overdue dates, not on parse errors.
func pastByDays(isoDate string, days int) bool {
	ts, err := time.Parse("2006-01-02", strings.TrimSpace(isoDate))
	if err != nil {
		return false
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	return ts.Before(cutoff)
}

// isCalendarPast accepts loose timeline strings like `2026-Q3`, `2026-12`,
// or ISO `2026-12-15`. Returns true when the timeline is unambiguously in
// the past relative to today (YYYY-MM-DD). Conservative: ambiguous inputs
// return false (do not fire R5 alert on uncertain dates).
func isCalendarPast(timeline, today string) bool {
	t := strings.TrimSpace(timeline)
	// ISO YYYY-MM-DD
	if _, err := time.Parse("2006-01-02", t); err == nil {
		return t < today
	}
	// YYYY-MM
	if _, err := time.Parse("2006-01", t); err == nil {
		return t+"-31" < today
	}
	// YYYY-Q[1-4]
	if len(t) == 7 && t[4:6] == "-Q" {
		yearPart := t[:4]
		quarter := t[6]
		if quarter >= '1' && quarter <= '4' {
			// End-of-quarter month.
			monthByQuarter := map[byte]string{'1': "03-31", '2': "06-30", '3': "09-30", '4': "12-31"}
			end := yearPart + "-" + monthByQuarter[quarter]
			return end < today
		}
	}
	return false
}
