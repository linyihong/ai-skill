package app

// `ai-skill plans tree` — Phase 3 of plan-tree-hierarchy-governance.
//
// Pure read-only introspection: walks plans/active + plans/archived,
// parses frontmatter (reusing Phase 2 scanAllPlanFrontmatter), and
// renders the tree built from `parent` pointers. Output formats: text
// (ASCII), JSON (nested), markdown (bullet list).
//
// CLI is a visualization companion to the Phase 2 commit-msg validators
// (validatePlanTree*); they own enforcement, this owns observation. The
// renderers are registered in enforcement-registry.yaml's
// internal_helper_allowlist (not as rule_class executors).
//
// See: plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/03-cli-tree-subcommand.md

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/planvalidate"
)

// runPlansValidate is Phase 2.4: a THIN CLI consumer of the planvalidate engine.
//
// It is transport only — it owns no validation logic. It loads the plan set via
// the shared loader (normalizedPlansFromRoot, fixed to <root>/plans/active and
// <root>/plans/archived — scope (A)), calls the engine entrypoint
// planvalidate.Validate, and renders the findings as text or JSON. Exit code is
// the manual consumer's transport mapping: non-zero iff a blocking finding exists.
//
// Explicit non-goals (do NOT add here): custom plans dir, schema dialect
// handling, loader plugins, external path conventions, filtering/policy. Those
// are Q8 / Phase 3, not the consumer surface.
func runPlansValidate(args []string, stdout io.Writer, stderr io.Writer) int {
	opts := struct{ root, format string }{}
	fs := newFlagSet("plans validate", stderr)
	fs.StringVar(&opts.root, "root", ".", "repository root (default: current directory)")
	fs.StringVar(&opts.format, "format", "text", "render format: text | json")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	switch opts.format {
	case "text", "json":
	default:
		_, _ = fmt.Fprintf(stderr, "invalid --format %q (want: text | json)\n", opts.format)
		return ExitInvalidUsage
	}

	models, compat := normalizedPlansFromRoot(opts.root)
	findings := planvalidate.Validate(planvalidate.ValidationContext{
		Root:          opts.root,
		ExecutionMode: planvalidate.ModeManual,
	}, models)
	findings = append(findings, compat...) // compat-layer rejects (e.g. unsupported schema_version)

	blocking := 0
	for _, f := range findings {
		if f.Blocking {
			blocking++
		}
	}

	if opts.format == "json" {
		type jsonFinding struct {
			RuleID   string `json:"rule_id"`
			Message  string `json:"message"`
			Blocking bool   `json:"blocking"`
		}
		payload := struct {
			Root     string        `json:"root"`
			Plans    int           `json:"plans"`
			Blocking int           `json:"blocking"`
			Findings []jsonFinding `json:"findings"`
		}{Root: opts.root, Plans: len(models), Blocking: blocking, Findings: []jsonFinding{}}
		for _, f := range findings {
			payload.Findings = append(payload.Findings, jsonFinding{f.RuleID, f.Message, f.Blocking})
		}
		b, _ := json.MarshalIndent(payload, "", "  ")
		_, _ = fmt.Fprintln(stdout, string(b))
	} else {
		for _, f := range findings {
			tag := "warn"
			if f.Blocking {
				tag = "BLOCK"
			}
			_, _ = fmt.Fprintf(stdout, "[%s] %s: %s\n", tag, f.RuleID, f.Message)
		}
		_, _ = fmt.Fprintf(stdout, "plans=%d findings=%d blocking=%d\n", len(models), len(findings), blocking)
	}

	if blocking > 0 {
		return ExitValidationFailed
	}
	return ExitSuccess
}

func runPlans(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill plans <tree|validate> [flags]")
		return ExitInvalidUsage
	}
	switch args[0] {
	case "tree":
		return runPlansTree(args[1:], stdout, stderr)
	case "validate":
		return runPlansValidate(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		_, _ = fmt.Fprintln(stdout, "usage: ai-skill plans <tree|validate> [flags]")
		_, _ = fmt.Fprintln(stdout, "")
		_, _ = fmt.Fprintln(stdout, "subcommands:")
		_, _ = fmt.Fprintln(stdout, "  tree      render plan tree built from frontmatter parent pointers")
		_, _ = fmt.Fprintln(stdout, "  validate  run the plan_profile.core engine over <root>/plans (thin consumer)")
		return ExitSuccess
	default:
		_, _ = fmt.Fprintf(stderr, "unknown plans subcommand: %s\n", args[0])
		return ExitInvalidUsage
	}
}

type plansTreeOptions struct {
	root           string
	state          string // active | archived | all
	format         string // text | json | markdown
	includeOrphans bool
}

func runPlansTree(args []string, stdout io.Writer, stderr io.Writer) int {
	opts := plansTreeOptions{}
	fs := newFlagSet("plans tree", stderr)
	fs.StringVar(&opts.root, "root", ".", "repository root (default: current directory)")
	fs.StringVar(&opts.state, "state", "all", "filter by plan location: active | archived | all")
	fs.StringVar(&opts.format, "format", "text", "render format: text | json | markdown")
	fs.BoolVar(&opts.includeOrphans, "include-orphans", false, "include sub/spike plans whose parent does not resolve")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	switch opts.state {
	case "active", "archived", "all":
	default:
		_, _ = fmt.Fprintf(stderr, "invalid --state %q (want: active | archived | all)\n", opts.state)
		return ExitInvalidUsage
	}
	switch opts.format {
	case "text", "json", "markdown":
	default:
		_, _ = fmt.Fprintf(stderr, "invalid --format %q (want: text | json | markdown)\n", opts.format)
		return ExitInvalidUsage
	}

	all := scanAllPlanFrontmatter(opts.root)
	nodes := buildPlanTree(all, opts.state, opts.includeOrphans)

	var out string
	switch opts.format {
	case "text":
		out = renderPlanTreeText(nodes, opts.root)
	case "json":
		out = renderPlanTreeJSON(nodes, opts.root)
	case "markdown":
		out = renderPlanTreeMarkdown(nodes, opts.root)
	}
	_, _ = fmt.Fprint(stdout, out)
	return ExitSuccess
}

// ---------------------------------------------------------------------------
// Tree building
// ---------------------------------------------------------------------------

// planTreeNode is one node in the rendered tree.
type planTreeNode struct {
	ID             string          `json:"id"`
	PlanKind       string          `json:"plan_kind,omitempty"`
	Status         string          `json:"status,omitempty"`
	Path           string          `json:"path"`
	Location       string          `json:"location"` // active | archived
	Progress       string          `json:"progress,omitempty"`
	BlockerCount   int             `json:"blocker_count"`
	ArchiveReady   bool            `json:"archive_ready,omitempty"`
	Required       bool            `json:"required_for_completion,omitempty"`
	Children       []*planTreeNode `json:"children,omitempty"`
	IsOrphan       bool            `json:"is_orphan,omitempty"`
	UnresolvedRef  string          `json:"unresolved_parent,omitempty"`
}

// buildPlanTree converts a flat frontmatter list into rooted trees.
//
//   - state filter: drop plans not in the requested set (active|archived|all)
//   - children indexed by parent.id; nodes with parent pointing nowhere are
//     emitted as orphans iff includeOrphans (else dropped)
//   - main plans (or anything with parent==nil/"") form tree roots
//   - cycle protection: any node visited twice during DFS is skipped
//     (and a warning is encoded in UnresolvedRef = "cycle:<id>")
func buildPlanTree(all []PlanFrontmatter, state string, includeOrphans bool) []*planTreeNode {
	// Apply state filter and build basic node list.
	byID := map[string]*planTreeNode{}
	var nodes []*planTreeNode
	for _, p := range all {
		if !p.HasFrontmatter || p.ID == "" {
			continue
		}
		loc := planLocation(p.Path)
		if state != "all" && loc != state {
			continue
		}
		n := &planTreeNode{
			ID:       p.ID,
			PlanKind: p.PlanKind,
			Status:   p.Status,
			Path:     p.Path,
			Location: loc,
			Required: p.RequiredForCompletion != nil && *p.RequiredForCompletion,
		}
		// Acceptance-criteria progress (best-effort; ignore errors).
		if res, err := ScanCheckboxesInFile(p.Path); err == nil {
			total := len(res.CheckedLines) + len(res.UncheckedLines)
			if total > 0 {
				n.Progress = fmt.Sprintf("%d/%d", len(res.CheckedLines), total)
			}
		}
		byID[p.ID] = n
		nodes = append(nodes, n)
	}

	// Link children to parents (rebuild the parent pointer from the raw list).
	parentOf := map[string]string{}
	for _, p := range all {
		if p.HasFrontmatter && p.ID != "" {
			parentOf[p.ID] = strings.TrimSpace(p.Parent)
		}
	}
	var roots []*planTreeNode
	var orphans []*planTreeNode
	for _, n := range nodes {
		parentID := parentOf[n.ID]
		if parentID == "" {
			// root candidate
			roots = append(roots, n)
			continue
		}
		parent, ok := byID[parentID]
		if !ok {
			n.IsOrphan = true
			n.UnresolvedRef = parentID
			if includeOrphans {
				orphans = append(orphans, n)
			}
			continue
		}
		parent.Children = append(parent.Children, n)
	}

	// Sort children deterministically, then compute aggregate fields.
	for _, n := range nodes {
		sort.Slice(n.Children, func(i, j int) bool { return n.Children[i].ID < n.Children[j].ID })
	}
	sort.Slice(roots, func(i, j int) bool { return roots[i].ID < roots[j].ID })
	sort.Slice(orphans, func(i, j int) bool { return orphans[i].ID < orphans[j].ID })

	// Compute blocker_count + archive_ready per node (post-link).
	for _, n := range nodes {
		var blockers int
		allReqCompleted := true
		for _, c := range n.Children {
			if c.Required && c.Status != "completed" {
				blockers++
				allReqCompleted = false
			}
		}
		n.BlockerCount = blockers
		if n.PlanKind == "main" {
			n.ArchiveReady = allReqCompleted
		}
	}

	// Cycle detection (defensive — frontmatter parents are id-based so a cycle
	// requires deliberate self-reference, but check anyway).
	visited := map[string]bool{}
	var visit func(n *planTreeNode, stack map[string]bool) bool
	visit = func(n *planTreeNode, stack map[string]bool) bool {
		if stack[n.ID] {
			n.UnresolvedRef = "cycle:" + n.ID
			return true
		}
		if visited[n.ID] {
			return false
		}
		visited[n.ID] = true
		stack[n.ID] = true
		for _, c := range n.Children {
			if visit(c, stack) {
				return true
			}
		}
		delete(stack, n.ID)
		return false
	}
	for _, r := range roots {
		visit(r, map[string]bool{})
	}

	if includeOrphans {
		return append(roots, orphans...)
	}
	return roots
}

func planLocation(path string) string {
	if strings.HasPrefix(path, "plans/active/") {
		return "active"
	}
	if strings.HasPrefix(path, "plans/archived/") {
		return "archived"
	}
	return ""
}

// ---------------------------------------------------------------------------
// Renderers
// ---------------------------------------------------------------------------

func renderPlanTreeText(roots []*planTreeNode, root string) string {
	var b strings.Builder
	if root == "" {
		root = "."
	}
	fmt.Fprintf(&b, "Plan Tree (root=%s, %d top-level node(s))\n", root, len(roots))
	for i, r := range roots {
		isLast := i == len(roots)-1
		renderNodeText(&b, r, "", isLast)
	}
	return b.String()
}

func renderNodeText(b *strings.Builder, n *planTreeNode, prefix string, isLast bool) {
	branch := "├── "
	cont := "│   "
	if isLast {
		branch = "└── "
		cont = "    "
	}
	line := formatNodeLine(n)
	fmt.Fprintf(b, "%s%s%s\n", prefix, branch, line)
	for i, c := range n.Children {
		renderNodeText(b, c, prefix+cont, i == len(n.Children)-1)
	}
}

func formatNodeLine(n *planTreeNode) string {
	var parts []string
	parts = append(parts, n.ID)
	if n.PlanKind != "" {
		parts = append(parts, "["+n.PlanKind+"]")
	}
	if n.Status != "" {
		parts = append(parts, "status="+n.Status)
	}
	if n.Progress != "" {
		parts = append(parts, "progress="+n.Progress)
	}
	if n.BlockerCount > 0 {
		parts = append(parts, fmt.Sprintf("blockers=%d", n.BlockerCount))
	}
	if n.PlanKind == "main" && n.ArchiveReady {
		parts = append(parts, "archive_ready=✓")
	}
	if n.Required {
		parts = append(parts, "required=true")
	}
	if n.IsOrphan {
		parts = append(parts, "ORPHAN(parent="+n.UnresolvedRef+")")
	} else if strings.HasPrefix(n.UnresolvedRef, "cycle:") {
		parts = append(parts, "CYCLE")
	}
	if n.Location != "" {
		parts = append(parts, "loc="+n.Location)
	}
	return strings.Join(parts, " ")
}

func renderPlanTreeJSON(roots []*planTreeNode, root string) string {
	payload := map[string]any{
		"root":         root,
		"top_level":    len(roots),
		"plan_kind":    "tree",
		"trees":        roots,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\":%q}\n", err.Error())
	}
	return string(data) + "\n"
}

func renderPlanTreeMarkdown(roots []*planTreeNode, root string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Plan Tree (root=%s)\n\n", root)
	for _, r := range roots {
		renderNodeMarkdown(&b, r, 0)
	}
	return b.String()
}

func renderNodeMarkdown(b *strings.Builder, n *planTreeNode, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(b, "%s- `%s` (%s) %s\n", indent, n.ID, statusOrUnknown(n.Status), markdownExtras(n))
	for _, c := range n.Children {
		renderNodeMarkdown(b, c, depth+1)
	}
}

func statusOrUnknown(s string) string {
	if s == "" {
		return "no-status"
	}
	return s
}

func markdownExtras(n *planTreeNode) string {
	var parts []string
	if n.Progress != "" {
		parts = append(parts, "progress=`"+n.Progress+"`")
	}
	if n.BlockerCount > 0 {
		parts = append(parts, fmt.Sprintf("blockers=`%d`", n.BlockerCount))
	}
	if n.PlanKind == "main" && n.ArchiveReady {
		parts = append(parts, "archive_ready=`✓`")
	}
	if n.IsOrphan {
		parts = append(parts, "**ORPHAN** (parent=`"+n.UnresolvedRef+"`)")
	}
	if n.Location != "" {
		parts = append(parts, "loc=`"+n.Location+"`")
	}
	return strings.Join(parts, " ")
}
