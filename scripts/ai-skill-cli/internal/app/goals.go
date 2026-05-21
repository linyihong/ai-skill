package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/pathutil"
)

type goalsOptions struct {
	command     string
	projectPath string
	dryRun      bool
	jsonOutput  bool
	plainOutput bool
	id          string
	title       string
	source      string
	priority    string
	next        string
	criteria    string
	status      string
	note        string
	owner       string
	parent      string
	reason      string
	missing     string
	decision    string
	strengthen  string
	parallel    string
	validated   bool
	superseded  bool
	plans       stringSliceFlag
	todos       stringSliceFlag
}

type stringSliceFlag []string

func (s *stringSliceFlag) String() string { return strings.Join(*s, ",") }
func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func runGoals(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill goals <status|init|start|update|split|pause|complete|cleanup> [flags]")
		return ExitInvalidUsage
	}

	opts := goalsOptions{command: args[0]}
	fs := newFlagSet("goals "+opts.command, stderr)
	fs.StringVar(&opts.projectPath, "project", ".", "project root")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview planned changes without writing")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	fs.StringVar(&opts.id, "id", "", "goal id")
	fs.StringVar(&opts.title, "title", "", "goal title")
	fs.StringVar(&opts.source, "source", "", "source request")
	fs.StringVar(&opts.priority, "priority", "P1", "goal priority")
	fs.StringVar(&opts.next, "next", "", "next action")
	fs.StringVar(&opts.criteria, "criteria", "", "completion criteria")
	fs.StringVar(&opts.status, "status", "", "goal status")
	fs.StringVar(&opts.note, "note", "", "progress note")
	fs.StringVar(&opts.owner, "owner", "", "goal owner")
	fs.StringVar(&opts.parent, "parent", "", "parent goal id")
	fs.StringVar(&opts.reason, "reason", "", "pause reason")
	fs.StringVar(&opts.missing, "missing", "", "missing work")
	fs.StringVar(&opts.decision, "decision", "", "decision needed")
	fs.StringVar(&opts.strengthen, "strengthen", "", "needs strengthening")
	fs.StringVar(&opts.parallel, "parallelization", "parallelizable", "parallelization mode")
	fs.BoolVar(&opts.validated, "validated", false, "confirm validation before deleting completed goal")
	fs.BoolVar(&opts.superseded, "superseded", false, "mark pause as superseded")
	fs.Var(&opts.plans, "plan", "planning reference")
	fs.Var(&opts.todos, "todo", "todo reference")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}

	result := buildGoalsResult(opts)
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

func buildGoalsResult(opts goalsOptions) Result {
	result := Result{
		Command:        "goals " + opts.command,
		Mode:           "check",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}

	project, projectCheck := resolveTargetProject(opts.projectPath)
	result.Checks = append(result.Checks, projectCheck)
	if projectCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_project", Message: projectCheck.Message, Remediation: "Pass --project with an existing project directory."}
		return result
	}

	switch opts.command {
	case "status":
		return buildGoalsStatusResult(result, project)
	case "init":
		return buildGoalsInitResult(result, project, opts.dryRun)
	case "start":
		return buildGoalsStartResult(result, project, opts)
	case "update":
		return buildGoalsUpdateResult(result, project, opts)
	case "split":
		return buildGoalsSplitResult(result, project, opts)
	case "pause":
		return buildGoalsPauseResult(result, project, opts)
	case "complete":
		return buildGoalsCompleteResult(result, project, opts)
	case "cleanup":
		return buildGoalsCleanupResult(result, project)
	default:
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_goals_command", Message: "unsupported goals command: " + opts.command, Remediation: "Use status, init, start, update, split, pause, complete, or cleanup."}
		return result
	}
}

func buildGoalsStatusResult(result Result, project string) Result {
	root := filepath.Join(project, ".agent-goals")
	goalsDir := filepath.Join(root, "goals")
	locksDir := filepath.Join(root, "locks")
	if _, err := os.Stat(root); err != nil {
		result.Checks = append(result.Checks, Check{Name: "goal_ledger", Status: "missing", Message: ".agent-goals does not exist"})
		return result
	}

	result.Checks = append(result.Checks,
		pathExistenceCheck("goal_ledger", root),
		pathExistenceCheck("goals_dir", goalsDir),
		pathExistenceCheck("locks_dir", locksDir),
	)
	goalFiles := listMarkdownFiles(goalsDir)
	lockFiles := listDirEntries(locksDir)
	result.Checks = append(result.Checks,
		Check{Name: "goal_files", Status: "ok", Message: fmt.Sprintf("%d files", len(goalFiles))},
		Check{Name: "locks", Status: "ok", Message: fmt.Sprintf("%d entries", len(lockFiles))},
	)
	return result
}

func buildGoalsInitResult(result Result, project string, dryRun bool) Result {
	result.Mode = "dry_run"
	root := filepath.Join(project, ".agent-goals")
	result.PlannedActions = []string{
		"create directory: " + filepath.Join(root, "goals"),
		"create directory: " + filepath.Join(root, "locks"),
		"create file: " + filepath.Join(root, "README.md"),
		"ensure git exclude contains .agent-goals/",
	}
	if !dryRun {
		result.Mode = "write"
		if err := ensureGoalsLedger(project); err != nil {
			return goalsError(result, ExitGeneralFailure, "goals_init_failed", err.Error())
		}
		result.Mutations = append(result.Mutations, result.PlannedActions...)
	}
	return result
}

func buildGoalsStartResult(result Result, project string, opts goalsOptions) Result {
	result.Mode = "write"
	if opts.id == "" || opts.title == "" || opts.source == "" {
		return goalsError(result, ExitInvalidUsage, "missing_required_goal_fields", "start requires --id, --title, and --source")
	}
	if err := validateGoalID(opts.id); err != nil {
		return goalsError(result, ExitInvalidUsage, "invalid_goal_id", err.Error())
	}
	if err := ensureGoalsLedger(project); err != nil {
		return goalsError(result, ExitGeneralFailure, "goals_init_failed", err.Error())
	}
	release, err := acquireGoalLock(project, opts.id)
	if err != nil {
		return goalsError(result, ExitUnsafeRepoState, "active_goal_lock", err.Error())
	}
	defer release()
	file := goalPath(project, opts.id)
	now := nowUTC()
	owner := defaultGoalOwner()
	content := goalContent(goalContentOptions{
		ID: opts.id, Title: opts.title, Source: opts.source, Priority: defaultString(opts.priority, "P1"),
		Owner: owner, Parallel: defaultString(opts.parallel, "parallelizable"), Created: now, Updated: now,
		Project: project, Next: defaultString(opts.next, "Define the next concrete action."),
		Criteria: defaultString(opts.criteria, "Define observable completion criteria."),
		Plans:    opts.plans, Todos: opts.todos,
	})
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		return goalsError(result, ExitGeneralFailure, "write_failed", err.Error())
	}
	if err := refreshGoalsIndex(project); err != nil {
		return goalsError(result, ExitGeneralFailure, "index_refresh_failed", err.Error())
	}
	result.Mutations = append(result.Mutations, "wrote goal: "+file, "refreshed goal index")
	return result
}

func buildGoalsUpdateResult(result Result, project string, opts goalsOptions) Result {
	result.Mode = "write"
	if opts.id == "" {
		return goalsError(result, ExitInvalidUsage, "missing_goal_id", "update requires --id")
	}
	release, file, fail := lockExistingGoal(result, project, opts.id)
	if fail != nil {
		return *fail
	}
	defer release()
	textBytes, err := os.ReadFile(file)
	if err != nil {
		return goalsError(result, ExitGeneralFailure, "read_failed", err.Error())
	}
	text := string(textBytes)
	now := nowUTC()
	updates := map[string]string{"updated": now}
	if opts.status != "" {
		updates["status"] = opts.status
	}
	if opts.owner != "" {
		updates["owner"] = opts.owner
	}
	if opts.parallel != "" {
		updates["parallelization"] = opts.parallel
	}
	text = setFrontmatterValues(text, updates)
	if opts.next != "" {
		text = replaceSection(text, "Next Action", opts.next)
	}
	for _, plan := range opts.plans {
		text = addPlanningLink(text, "plan", plan, "linked during update")
	}
	for _, todo := range opts.todos {
		text = addPlanningLink(text, "todo", todo, "linked during update")
	}
	text = setOpenWork(text, "Missing work", opts.missing)
	text = setOpenWork(text, "Decision needed", opts.decision)
	text = setOpenWork(text, "Needs strengthening", opts.strengthen)
	text = appendProgress(text, defaultString(opts.note, "Goal updated."))
	if err := os.WriteFile(file, []byte(text), 0o644); err != nil {
		return goalsError(result, ExitGeneralFailure, "write_failed", err.Error())
	}
	if err := refreshGoalsIndex(project); err != nil {
		return goalsError(result, ExitGeneralFailure, "index_refresh_failed", err.Error())
	}
	result.Mutations = append(result.Mutations, "updated goal: "+file, "refreshed goal index")
	return result
}

func buildGoalsSplitResult(result Result, project string, opts goalsOptions) Result {
	if opts.parent == "" || opts.id == "" || opts.title == "" {
		return goalsError(result, ExitInvalidUsage, "missing_required_goal_fields", "split requires --parent, --id, and --title")
	}
	parent := goalPath(project, opts.parent)
	if _, err := os.Stat(parent); err != nil {
		return goalsError(result, ExitValidationFailed, "parent_goal_not_found", err.Error())
	}
	opts.source = "Child goal of " + opts.parent
	if opts.criteria == "" {
		opts.criteria = "Complete child goal " + opts.title + "."
	}
	if opts.priority == "" {
		opts.priority = "P2"
	}
	result = buildGoalsStartResult(result, project, opts)
	if result.ExitCode != ExitSuccess {
		return result
	}
	if bytes, err := os.ReadFile(parent); err == nil {
		_ = os.WriteFile(parent, []byte(appendProgress(string(bytes), "Split child goal "+opts.id+": "+opts.title)), 0o644)
		_ = refreshGoalsIndex(project)
	}
	return result
}

func buildGoalsPauseResult(result Result, project string, opts goalsOptions) Result {
	if opts.id == "" {
		return goalsError(result, ExitInvalidUsage, "missing_goal_id", "pause requires --id")
	}
	opts.status = "paused"
	if opts.superseded {
		opts.status = "superseded"
	}
	opts.note = defaultString(opts.reason, "Goal "+opts.status+".")
	opts.next = "Resume only when this goal becomes the highest priority again."
	return buildGoalsUpdateResult(result, project, opts)
}

func buildGoalsCompleteResult(result Result, project string, opts goalsOptions) Result {
	result.Mode = "write"
	if opts.id == "" {
		return goalsError(result, ExitInvalidUsage, "missing_goal_id", "complete requires --id")
	}
	if !opts.validated {
		opts.status = "needs-validation"
		opts.note = defaultString(opts.note, "Completion requested without validation; keeping goal.")
		return buildGoalsUpdateResult(result, project, opts)
	}
	release, file, fail := lockExistingGoal(result, project, opts.id)
	if fail != nil {
		return *fail
	}
	defer release()
	if err := os.Remove(file); err != nil {
		return goalsError(result, ExitGeneralFailure, "delete_failed", err.Error())
	}
	if err := refreshGoalsIndex(project); err != nil {
		return goalsError(result, ExitGeneralFailure, "index_refresh_failed", err.Error())
	}
	result.Mutations = append(result.Mutations, "deleted completed goal: "+opts.id, "refreshed goal index")
	return result
}

func buildGoalsCleanupResult(result Result, project string) Result {
	result.Mode = "write"
	if err := ensureGoalsLedger(project); err != nil {
		return goalsError(result, ExitGeneralFailure, "goals_init_failed", err.Error())
	}
	removed, err := cleanupGoalLocks(project)
	if err != nil {
		return goalsError(result, ExitGeneralFailure, "cleanup_failed", err.Error())
	}
	if err := refreshGoalsIndex(project); err != nil {
		return goalsError(result, ExitGeneralFailure, "index_refresh_failed", err.Error())
	}
	result.Mutations = append(result.Mutations, fmt.Sprintf("removed stale locks: %d", removed), "refreshed goal index")
	return result
}

func goalsError(result Result, code int, errCode string, message string) Result {
	result.Status = "blocked"
	result.ExitCode = code
	result.Error = &CommandError{Code: errCode, Message: message}
	return result
}

func pathExistenceCheck(name string, path string) Check {
	normalized, err := pathutil.NormalizeForReport(path)
	if err != nil {
		normalized = path
	}
	if _, err := os.Stat(path); err != nil {
		return Check{Name: name, Status: "missing", Message: normalized}
	}
	return Check{Name: name, Status: "ok", Message: normalized}
}

func listMarkdownFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	files := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	return files
}

func listDirEntries(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

type goalContentOptions struct {
	ID       string
	Title    string
	Source   string
	Priority string
	Owner    string
	Parallel string
	Created  string
	Updated  string
	Project  string
	Next     string
	Criteria string
	Plans    []string
	Todos    []string
}

func ledgerRoot(project string) string { return filepath.Join(project, ".agent-goals") }
func goalsDir(project string) string   { return filepath.Join(ledgerRoot(project), "goals") }
func locksDir(project string) string   { return filepath.Join(ledgerRoot(project), "locks") }
func goalPath(project, id string) string {
	return filepath.Join(goalsDir(project), id+".md")
}

func ensureGoalsLedger(project string) error {
	if err := os.MkdirAll(goalsDir(project), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(locksDir(project), 0o755); err != nil {
		return err
	}
	if err := ensureGoalGitExclude(project); err != nil {
		return err
	}
	return refreshGoalsIndex(project)
}

func ensureGoalGitExclude(project string) error {
	gitDir := filepath.Join(project, ".git")
	if info, err := os.Stat(gitDir); err != nil || !info.IsDir() {
		return nil
	}
	exclude := filepath.Join(gitDir, "info", "exclude")
	if err := os.MkdirAll(filepath.Dir(exclude), 0o755); err != nil {
		return err
	}
	var existing string
	if bytes, err := os.ReadFile(exclude); err == nil {
		existing = string(bytes)
	}
	if strings.Contains(existing, ".agent-goals/") {
		return nil
	}
	return os.WriteFile(exclude, []byte(strings.TrimRight(existing, "\n")+"\n.agent-goals/\n"), 0o644)
}

func refreshGoalsIndex(project string) error {
	root := ledgerRoot(project)
	if err := os.MkdirAll(goalsDir(project), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(locksDir(project), 0o755); err != nil {
		return err
	}
	files := listMarkdownFiles(goalsDir(project))
	lines := []string{
		"# Agent Goals",
		"",
		"Temporary conversation goal ledger for humans and agents. Do not commit this directory.",
		"Use this table to see unfinished work, needed decisions, priority, and what needs strengthening.",
		"Delete goal files only after completion criteria and validation are satisfied.",
		"",
		"## Goal Table",
		"",
		"| Priority | Status | Mode | Owner | Lock | Goal | Open Work / Decisions | Planning / Todo Links | Next Action | Updated |",
		"| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |",
	}
	if len(files) == 0 {
		lines = append(lines, "| - | - | - | - | unlocked | No active goals | - | - | - | - |")
	} else {
		for _, name := range files {
			path := filepath.Join(goalsDir(project), name)
			textBytes, _ := os.ReadFile(path)
			text := string(textBytes)
			meta := parseFrontmatter(text)
			title := firstHeading(text, strings.TrimSuffix(name, ".md"))
			id := defaultString(meta["id"], strings.TrimSuffix(name, ".md"))
			row := []string{
				meta["priority"],
				meta["status"],
				meta["parallelization"],
				meta["owner"],
				goalLockState(project, id),
				fmt.Sprintf("[%s](goals/%s)", title, name),
				openWorkSummary(text),
				planningLinksSummary(text),
				compactSummary(sectionText(text, "Next Action", "Completion Criteria")),
				meta["updated"],
			}
			lines = append(lines, "| "+strings.Join(row, " | ")+" |")
		}
	}
	lines = append(lines,
		"",
		"## Notes",
		"",
		"- Goal details live in `goals/`.",
		"- Locks live in `locks/` and are temporary.",
		"- This file is an index generated by `ai-skill goals`.",
		"- If `Lock` shows another active owner for overlapping work, stop and ask before editing.",
	)
	return os.WriteFile(filepath.Join(root, "README.md"), []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func acquireGoalLock(project, id string) (func(), error) {
	if err := validateGoalID(id); err != nil {
		return nil, err
	}
	lock := filepath.Join(locksDir(project), id+".lock")
	if info, err := os.Stat(lock); err == nil && info.IsDir() {
		age := time.Since(info.ModTime())
		ttl := goalLockTTL()
		if age < ttl {
			ownerBytes, _ := os.ReadFile(filepath.Join(lock, "owner"))
			return nil, fmt.Errorf("active goal lock detected: %s owner=%s ageSeconds=%d", id, strings.TrimSpace(string(ownerBytes)), int(age.Seconds()))
		}
		if err := os.RemoveAll(lock); err != nil {
			return nil, err
		}
	}
	if err := os.MkdirAll(locksDir(project), 0o755); err != nil {
		return nil, err
	}
	if err := os.Mkdir(lock, 0o755); err != nil {
		return nil, err
	}
	_ = os.WriteFile(filepath.Join(lock, "pid"), []byte(strconv.Itoa(os.Getpid())+"\n"), 0o644)
	_ = os.WriteFile(filepath.Join(lock, "owner"), []byte(defaultGoalOwner()+"\n"), 0o644)
	_ = os.WriteFile(filepath.Join(lock, "startedAt"), []byte(nowUTC()+"\n"), 0o644)
	return func() { _ = os.RemoveAll(lock) }, nil
}

func lockExistingGoal(result Result, project, id string) (func(), string, *Result) {
	if err := validateGoalID(id); err != nil {
		fail := goalsError(result, ExitInvalidUsage, "invalid_goal_id", err.Error())
		return nil, "", &fail
	}
	file := goalPath(project, id)
	if _, err := os.Stat(file); err != nil {
		fail := goalsError(result, ExitValidationFailed, "goal_not_found", err.Error())
		return nil, "", &fail
	}
	release, err := acquireGoalLock(project, id)
	if err != nil {
		fail := goalsError(result, ExitUnsafeRepoState, "active_goal_lock", err.Error())
		return nil, "", &fail
	}
	return release, file, nil
}

func cleanupGoalLocks(project string) (int, error) {
	entries, err := os.ReadDir(locksDir(project))
	if err != nil {
		return 0, nil
	}
	removed := 0
	ttl := goalLockTTL()
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasSuffix(entry.Name(), ".lock") {
			continue
		}
		path := filepath.Join(locksDir(project), entry.Name())
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if time.Since(info.ModTime()) >= ttl {
			if err := os.RemoveAll(path); err != nil {
				return removed, err
			}
			removed++
		}
	}
	return removed, nil
}

func goalLockTTL() time.Duration {
	raw := os.Getenv("AGENT_GOALS_LOCK_TTL_SECONDS")
	if raw == "" {
		raw = "1800"
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		seconds = 1800
	}
	return time.Duration(seconds) * time.Second
}

func validateGoalID(id string) error {
	if id == "" {
		return fmt.Errorf("goal id is required")
	}
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			continue
		}
		return fmt.Errorf("invalid id %q; use letters, numbers, dot, underscore, or dash", id)
	}
	return nil
}

func goalContent(opts goalContentOptions) string {
	var b strings.Builder
	fmt.Fprintf(&b, "---\nid: %s\npriority: %s\nstatus: active\nowner: %s\nparallelization: %s\ncreated: %s\nupdated: %s\nproject: %s\n---\n\n", opts.ID, opts.Priority, opts.Owner, opts.Parallel, opts.Created, opts.Updated, opts.Project)
	fmt.Fprintf(&b, "# %s\n\n## Source Request\n%s\n\n## Scope\n- In:\n- Out:\n- Affected paths/repos:\n\n## Subgoals\n- [ ] %s\n\n", opts.Title, opts.Source, opts.Title)
	b.WriteString("## Planning / Todo Links\n| Type | Reference | Status / Note |\n| --- | --- | --- |\n")
	for _, plan := range opts.Plans {
		fmt.Fprintf(&b, "| plan | %s | linked at creation |\n", plan)
	}
	for _, todo := range opts.Todos {
		fmt.Fprintf(&b, "| todo | %s | linked at creation |\n", todo)
	}
	fmt.Fprintf(&b, "\n## Owner / Lock Decision\n- Owner: %s\n- Lock: unlocked between helper operations\n- Parallelization: %s\n\n", opts.Owner, opts.Parallel)
	b.WriteString("## Open Work / Decisions\n- Missing work: none\n- Decision needed: none\n- Needs strengthening: none\n\n## Dependencies\n- none\n\n")
	fmt.Fprintf(&b, "## Progress\n- %s: Goal created.\n\n## Next Action\n%s\n\n## Completion Criteria\n- [ ] %s\n\n## Validation\n- not yet validated\n\n## Handoff Notes\n- none\n", nowUTC(), opts.Next, opts.Criteria)
	return b.String()
}

func nowUTC() string { return time.Now().UTC().Format("2006-01-02T15:04:05Z") }

func defaultGoalOwner() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "unknown"
	}
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		user = "unknown"
	}
	return user + "@" + host
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func parseFrontmatter(text string) map[string]string {
	meta := map[string]string{}
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return meta
	}
	for _, line := range lines[1:] {
		if line == "---" {
			break
		}
		if key, value, ok := strings.Cut(line, ": "); ok {
			meta[key] = value
		}
	}
	return meta
}

func setFrontmatterValues(text string, updates map[string]string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return text
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return text
	}
	for key, value := range updates {
		prefix := key + ": "
		found := false
		for i := 1; i < end; i++ {
			if strings.HasPrefix(lines[i], prefix) {
				lines[i] = prefix + value
				found = true
				break
			}
		}
		if !found {
			lines = append(lines[:end], append([]string{prefix + value}, lines[end:]...)...)
			end++
		}
	}
	return strings.Join(lines, "\n")
}

func replaceSection(text, section, value string) string {
	marker := "\n## " + section + "\n"
	start := strings.Index(text, marker)
	if start == -1 {
		return text
	}
	start += len(marker)
	end := strings.Index(text[start:], "\n## ")
	if end == -1 {
		return text[:start] + value + "\n"
	}
	return text[:start] + value + "\n" + text[start+end:]
}

func addPlanningLink(text, kind, reference, note string) string {
	if reference == "" || strings.Contains(text, "| "+kind+" | "+reference+" |") {
		return text
	}
	marker := "\n## Planning / Todo Links\n"
	idx := strings.Index(text, marker)
	if idx == -1 {
		return text
	}
	start := idx + len(marker)
	next := strings.Index(text[start:], "\n## ")
	if next == -1 {
		return text
	}
	insert := start + next
	return text[:insert] + fmt.Sprintf("| %s | %s | %s |\n", kind, reference, note) + text[insert:]
}

func setOpenWork(text, label, value string) string {
	if value == "" {
		return text
	}
	prefix := "- " + label + ":"
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = prefix + " " + value
			return strings.Join(lines, "\n")
		}
	}
	return text
}

func appendProgress(text, note string) string {
	marker := "\n## Progress\n"
	idx := strings.Index(text, marker)
	if idx == -1 {
		return text + "\n## Progress\n- " + nowUTC() + ": " + note + "\n"
	}
	start := idx + len(marker)
	next := strings.Index(text[start:], "\n## ")
	if next == -1 {
		return text + "- " + nowUTC() + ": " + note + "\n"
	}
	insert := start + next
	return text[:insert] + "- " + nowUTC() + ": " + note + "\n" + text[insert:]
}

func sectionText(text, section, nextSection string) string {
	marker := "\n## " + section + "\n"
	idx := strings.Index(text, marker)
	if idx == -1 {
		return ""
	}
	start := idx + len(marker)
	endMarker := "\n## " + nextSection + "\n"
	end := strings.Index(text[start:], endMarker)
	if end == -1 {
		return strings.TrimSpace(text[start:])
	}
	return strings.TrimSpace(text[start : start+end])
}

func compactSummary(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if value == "" {
		return "-"
	}
	if len(value) > 100 {
		return value[:97] + "..."
	}
	return value
}

func firstHeading(text, fallback string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return fallback
}

func openWorkSummary(text string) string {
	items := []string{}
	for _, line := range strings.Split(sectionText(text, "Open Work / Decisions", "Dependencies"), "\n") {
		for _, prefix := range []string{"- Missing work:", "- Decision needed:", "- Needs strengthening:"} {
			if strings.HasPrefix(line, prefix) {
				value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
				if value != "" && value != "none" {
					items = append(items, strings.TrimPrefix(strings.TrimSuffix(prefix, ":"), "- ")+": "+value)
				}
			}
		}
	}
	if len(items) == 0 {
		return "-"
	}
	return strings.Join(items, "<br>")
}

func planningLinksSummary(text string) string {
	section := sectionText(text, "Planning / Todo Links", "Owner / Lock Decision")
	items := []string{}
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "|") || strings.Contains(line, "---") || strings.Contains(line, "Reference") {
			continue
		}
		parts := strings.Split(strings.Trim(line, "|"), "|")
		if len(parts) >= 2 {
			items = append(items, strings.TrimSpace(parts[0])+": "+strings.TrimSpace(parts[1]))
		}
	}
	if len(items) == 0 {
		return "-"
	}
	return strings.Join(items, "<br>")
}

func goalLockState(project, id string) string {
	lock := filepath.Join(locksDir(project), id+".lock")
	if _, err := os.Stat(lock); err != nil {
		return "unlocked"
	}
	ownerBytes, _ := os.ReadFile(filepath.Join(lock, "owner"))
	pidBytes, _ := os.ReadFile(filepath.Join(lock, "pid"))
	owner := defaultString(strings.TrimSpace(string(ownerBytes)), "unknown")
	pid := strings.TrimSpace(string(pidBytes))
	if pid != "" {
		return "locked by " + owner + " pid=" + pid
	}
	return "locked by " + owner
}
