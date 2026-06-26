package app

// Plan-tree validators (Phase 2 of 2026-06-02-1200-plan-tree-hierarchy-governance).
//
// Five commit-msg validators that mechanically enforce the plan tree frontmatter
// schema established in 01-frontmatter-schema.md:
//
//   - validatePlanTreeFrontmatter        block    sub/spike missing required fields
//   - validatePlanTreeArchiveOrder       block    main archive blocked by pending required child
//   - validatePlanTreeParentReference    block    parent: <id> must resolve to a real plan
//   - validatePlanTreeUniqueID           block    no two plans may share frontmatter id
//   - validatePlanTreeFolderConvention   warning  folder shape advisory (depth/_plan.md/NN- prefix)
//
// All validators are pre-existing-plan friendly: a plan file without YAML
// frontmatter is silently skipped (the Phase 4 migration sub-plan handles
// retro-fitting old plans). This avoids a one-shot break of the existing
// repository when these validators land.
//
// Files inside any path segment named "fixtures" are excluded from cross-plan
// scans (uniqueness / parent-existence indexes), so example/testdata files
// shipped alongside the schema docs don't collide with real plans.
//
// See: plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/02-validator-implementation.md

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// PlanFrontmatter captures the minimal-governance frontmatter fields used by
// the five plan-tree validators. Unknown fields are tolerated and ignored.
type PlanFrontmatter struct {
	Path                  string // repo-relative path (e.g. "plans/active/foo.md")
	HasFrontmatter        bool   // true if a YAML frontmatter block was found
	ID                    string
	PlanKind              string // "main" | "sub" | "spike" | ""
	Status                string // "draft" | "in-progress" | "completed" | ""
	Parent                string
	HasParentField        bool // distinguish parent: null from missing parent field
	RequiredForCompletion *bool
	HasReasonField        bool
	SubPlanReason         string // raw trimmed value; empty string = block
	SchemaVersion         string // declared schema_version (quotes stripped); "" = absent
}

var (
	planTreeFrontmatterDelim = []byte("---")
	planTreeYAMLKeyValueRE   = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*:\s*(.*)$`)
	planTreeFolderNameRE     = regexp.MustCompile(`^\d{2}-`)
)

// parsePlanFrontmatterFromBytes extracts the frontmatter fields from a plan
// markdown body. Returns a zero PlanFrontmatter with HasFrontmatter=false if
// the file does not start with a `---` frontmatter block.
func parsePlanFrontmatterFromBytes(path string, data []byte) PlanFrontmatter {
	pf := PlanFrontmatter{Path: path}
	text := string(data)
	// Tolerate UTF-8 BOM (some Windows editors prepend it) and leading whitespace.
	text = strings.TrimLeft(text, "\ufeff \t\r\n")
	if !strings.HasPrefix(text, "---") {
		return pf
	}
	idx := strings.Index(text, "---")
	if idx < 0 {
		return pf
	}
	rest := text[idx+3:]
	// Find the closing "---" on a line by itself.
	lines := strings.Split(rest, "\n")
	var body []string
	closed := false
	for i, line := range lines {
		if i == 0 && strings.TrimSpace(line) == "" {
			continue
		}
		if strings.TrimSpace(line) == "---" {
			closed = true
			break
		}
		body = append(body, line)
	}
	if !closed {
		return pf
	}
	pf.HasFrontmatter = true

	// Walk lines; tolerate folded scalars (`>` / `|`) for sub_plan_reason.
	var current string
	var foldingKey string
	var folded []string
	flushFolded := func() {
		if foldingKey == "" {
			return
		}
		val := strings.TrimSpace(strings.Join(folded, " "))
		assignField(&pf, foldingKey, val)
		foldingKey = ""
		folded = nil
	}
	for _, raw := range body {
		line := strings.TrimRight(raw, "\r")
		if foldingKey != "" {
			// Continuation if indented; otherwise flush and re-parse.
			if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
				folded = append(folded, strings.TrimSpace(line))
				continue
			}
			flushFolded()
		}
		m := planTreeYAMLKeyValueRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		key := m[1]
		val := strings.TrimSpace(m[2])
		current = key
		_ = current
		if val == ">" || val == ">-" || val == "|" || val == "|-" {
			foldingKey = key
			folded = nil
			continue
		}
		// Strip surrounding quotes.
		if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") && len(val) >= 2) ||
			(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") && len(val) >= 2) {
			val = val[1 : len(val)-1]
		}
		assignField(&pf, key, val)
	}
	flushFolded()
	return pf
}

func assignField(pf *PlanFrontmatter, key, val string) {
	switch key {
	case "id":
		pf.ID = val
	case "plan_kind":
		pf.PlanKind = val
	case "status":
		pf.Status = val
	case "parent":
		pf.HasParentField = true
		if val == "null" || val == "~" || val == "" {
			pf.Parent = ""
		} else {
			pf.Parent = val
		}
	case "required_for_completion":
		pf.HasReasonField = pf.HasReasonField // keep
		b := strings.ToLower(strings.TrimSpace(val)) == "true"
		f := strings.ToLower(strings.TrimSpace(val)) == "false"
		if b {
			t := true
			pf.RequiredForCompletion = &t
		} else if f {
			t := false
			pf.RequiredForCompletion = &t
		}
	case "sub_plan_reason":
		pf.HasReasonField = true
		pf.SubPlanReason = strings.TrimSpace(val)
	case "schema_version":
		// Quotes already stripped above ("1" -> 1), satisfying the Q3 loader
		// requirement; carried into RawPlan.SchemaVersion for the compat layer.
		pf.SchemaVersion = strings.TrimSpace(val)
	}
}

// parsePlanFrontmatterFile reads a path and parses its frontmatter.
func parsePlanFrontmatterFile(absPath string) (PlanFrontmatter, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return PlanFrontmatter{}, err
	}
	rel := absPath
	// Best-effort rel path: walk-up logic happens at call site.
	pf := parsePlanFrontmatterFromBytes(rel, data)
	return pf, nil
}

// scanAllPlanFrontmatter walks both plans/active and plans/archived under root
// and returns every parsed frontmatter. Files under any "fixtures" segment are
// excluded (those are documentation testdata, not real plans). Files without
// frontmatter are returned with HasFrontmatter=false so callers can choose to
// skip them or include them in coverage stats.
func scanAllPlanFrontmatter(root string) []PlanFrontmatter {
	var out []PlanFrontmatter
	for _, sub := range []string{"plans/active", "plans/archived"} {
		base := sub
		if root != "" {
			base = filepath.Join(root, sub)
		}
		_ = filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			rel = filepath.ToSlash(rel)
			if pathContainsFixturesSegment(rel) {
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil
			}
			pf := parsePlanFrontmatterFromBytes(rel, data)
			pf.Path = rel
			out = append(out, pf)
			return nil
		})
	}
	return out
}

func pathContainsFixturesSegment(rel string) bool {
	for _, seg := range strings.Split(rel, "/") {
		if seg == "fixtures" {
			return true
		}
	}
	return false
}

// stagedPlanPaths returns the subset of staged paths under plans/active or
// plans/archived that look like plan markdown.
func stagedPlanPaths(staged []string) []string {
	var out []string
	for _, s := range staged {
		s = filepath.ToSlash(s)
		if (strings.HasPrefix(s, "plans/active/") || strings.HasPrefix(s, "plans/archived/")) &&
			strings.HasSuffix(strings.ToLower(s), ".md") &&
			!pathContainsFixturesSegment(s) {
			out = append(out, s)
		}
	}
	return out
}

// readStagedPlan parses the staged plan path against the working tree (post-
// stage, pre-commit reflects what will be committed).
func readStagedPlan(root, rel string) (PlanFrontmatter, bool) {
	abs := rel
	if root != "" {
		abs = filepath.Join(root, rel)
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return PlanFrontmatter{}, false
	}
	pf := parsePlanFrontmatterFromBytes(rel, data)
	pf.Path = rel
	return pf, true
}

// ---------------------------------------------------------------------------
// Validator 1: validatePlanTreeFrontmatter (block)
//
// Sub or spike plans must declare: parent (non-empty), sub_plan_reason
// (non-empty string), required_for_completion (bool). Main plans require
// parent (may be null) but the schema permits an absent parent field on main;
// we accept either explicit `parent: null` or no parent field.
// ---------------------------------------------------------------------------
func validatePlanTreeFrontmatter(text string, staged []string, root string) string {
	if hasOptOutTrailer(text, "[skip-plan-tree-frontmatter]") {
		return ""
	}
	var violations []string
	for _, rel := range stagedPlanPaths(staged) {
		pf, ok := readStagedPlan(root, rel)
		if !ok || !pf.HasFrontmatter {
			continue
		}
		kind := pf.PlanKind
		if kind == "" {
			// Untagged plan — backward-compatible skip.
			continue
		}
		if kind != "sub" && kind != "spike" {
			continue
		}
		var missing []string
		if !pf.HasParentField || strings.TrimSpace(pf.Parent) == "" {
			missing = append(missing, "parent")
		}
		if !pf.HasReasonField || pf.SubPlanReason == "" {
			missing = append(missing, "sub_plan_reason (non-empty)")
		}
		if pf.RequiredForCompletion == nil {
			missing = append(missing, "required_for_completion")
		}
		if len(missing) > 0 {
			violations = append(violations, fmt.Sprintf("%s missing: %s", rel, strings.Join(missing, ", ")))
		}
	}
	if len(violations) == 0 {
		return ""
	}
	return "plan-tree-frontmatter: sub/spike plan(s) missing required frontmatter fields:\n    - " +
		strings.Join(violations, "\n    - ") +
		"\n  Add `parent: <main-id>`, `sub_plan_reason: <non-empty>` and `required_for_completion: true|false` " +
		"(see plans/active/2026-06-02-1200-plan-tree-hierarchy-governance/01-frontmatter-schema.md)" +
		"\n  Opt-out (emergency only): standalone `[skip-plan-tree-frontmatter]` trailer."
}

// ---------------------------------------------------------------------------
// Validator 2: validatePlanTreeArchiveOrder (block)
//
// When a main plan is being archived (its _plan.md or top-level .md moved
// into plans/archived/), every sub-plan declaring parent == <main>.id with
// required_for_completion: true must be status: completed (location-agnostic:
// still-active OR already-archived both qualify as long as status==completed).
// ---------------------------------------------------------------------------
func validatePlanTreeArchiveOrder(text string, staged []string, root string) string {
	if hasOptOutTrailer(text, "[skip-plan-tree-archive-order]") {
		return ""
	}
	archivedMains := []PlanFrontmatter{}
	for _, rel := range stagedPlanPaths(staged) {
		if !strings.HasPrefix(rel, "plans/archived/") {
			continue
		}
		pf, ok := readStagedPlan(root, rel)
		if !ok || !pf.HasFrontmatter {
			continue
		}
		if pf.PlanKind == "main" {
			archivedMains = append(archivedMains, pf)
		}
	}
	if len(archivedMains) == 0 {
		return ""
	}
	all := scanAllPlanFrontmatter(root)
	var violations []string
	for _, main := range archivedMains {
		if main.ID == "" {
			continue
		}
		for _, p := range all {
			if !p.HasFrontmatter {
				continue
			}
			if p.Parent != main.ID {
				continue
			}
			if p.RequiredForCompletion == nil || !*p.RequiredForCompletion {
				continue
			}
			if p.Status == "completed" {
				continue
			}
			violations = append(violations,
				fmt.Sprintf("main %s (%s) blocked by required sub %s (status=%s)",
					main.ID, main.Path, p.Path, displayStatus(p.Status)))
		}
	}
	if len(violations) == 0 {
		return ""
	}
	return "plan-tree-archive-order: cannot archive main plan(s) with unfinished required sub-plans:\n    - " +
		strings.Join(violations, "\n    - ") +
		"\n  Complete the required sub-plan(s) first or flip required_for_completion: false with rationale." +
		"\n  Opt-out (emergency only): standalone `[skip-plan-tree-archive-order]` trailer."
}

func displayStatus(s string) string {
	if s == "" {
		return "<missing>"
	}
	return s
}

// ---------------------------------------------------------------------------
// Validator 3: validatePlanTreeParentReference (block)
//
// Every sub/spike plan in the staged set whose parent field is non-empty must
// reference an id that exists somewhere in the repository (active or archived,
// excluding fixtures/). Prevents dangling parent pointers.
// ---------------------------------------------------------------------------
func validatePlanTreeParentReference(text string, staged []string, root string) string {
	if hasOptOutTrailer(text, "[skip-plan-tree-parent-reference]") {
		return ""
	}
	var stagedSubs []PlanFrontmatter
	for _, rel := range stagedPlanPaths(staged) {
		pf, ok := readStagedPlan(root, rel)
		if !ok || !pf.HasFrontmatter {
			continue
		}
		if pf.PlanKind != "sub" && pf.PlanKind != "spike" {
			continue
		}
		if strings.TrimSpace(pf.Parent) == "" {
			continue
		}
		stagedSubs = append(stagedSubs, pf)
	}
	if len(stagedSubs) == 0 {
		return ""
	}
	all := scanAllPlanFrontmatter(root)
	known := map[string]bool{}
	for _, p := range all {
		if p.HasFrontmatter && p.ID != "" {
			known[p.ID] = true
		}
	}
	var violations []string
	for _, p := range stagedSubs {
		if !known[p.Parent] {
			violations = append(violations,
				fmt.Sprintf("%s references parent: %q which does not resolve to any plan id", p.Path, p.Parent))
		}
	}
	if len(violations) == 0 {
		return ""
	}
	return "plan-tree-parent-reference: dangling parent pointer(s) detected:\n    - " +
		strings.Join(violations, "\n    - ") +
		"\n  Either fix the parent id or create the referenced main plan first." +
		"\n  Opt-out (emergency only): standalone `[skip-plan-tree-parent-reference]` trailer."
}

// ---------------------------------------------------------------------------
// Validator 4: validatePlanTreeUniqueID (block)
//
// No two plans (across active + archived, excluding fixtures/) may share an
// `id:` frontmatter value. Fires when staged plans introduce or modify the id.
// ---------------------------------------------------------------------------
func validatePlanTreeUniqueID(text string, staged []string, root string) string {
	if hasOptOutTrailer(text, "[skip-plan-tree-unique-id]") {
		return ""
	}
	all := scanAllPlanFrontmatter(root)
	byID := map[string][]string{}
	for _, p := range all {
		if !p.HasFrontmatter || p.ID == "" {
			continue
		}
		byID[p.ID] = append(byID[p.ID], p.Path)
	}
	stagedSet := map[string]bool{}
	for _, s := range staged {
		stagedSet[filepath.ToSlash(s)] = true
	}
	// Only surface duplicates whose duplicate set touches a staged file —
	// avoids re-litigating pre-existing repo state on unrelated commits.
	var ids []string
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var violations []string
	for _, id := range ids {
		paths := byID[id]
		if len(paths) < 2 {
			continue
		}
		touchesStage := false
		for _, p := range paths {
			if stagedSet[p] {
				touchesStage = true
				break
			}
		}
		if !touchesStage {
			continue
		}
		sort.Strings(paths)
		violations = append(violations, fmt.Sprintf("id %q appears in: %s", id, strings.Join(paths, ", ")))
	}
	if len(violations) == 0 {
		return ""
	}
	return "plan-tree-unique-id: duplicate frontmatter id(s) detected:\n    - " +
		strings.Join(violations, "\n    - ") +
		"\n  Plan ids must be globally unique across active + archived." +
		"\n  Opt-out (emergency only): standalone `[skip-plan-tree-unique-id]` trailer."
}

// ---------------------------------------------------------------------------
// Validator 5: validatePlanTreeFolderConvention (warning)
//
// Warning-only advisory checks against the UI convention:
//   - A sub-plan folder (one that contains any sub-plan files) should contain a _plan.md.
//   - Files inside a plan folder should match `^\d{2}-` (NN- prefix) or be `_plan.md`.
//   - Path depth under plans/active or plans/archived should be < 3 levels.
// Returns warnings prefixed with `plan-tree-folder-convention (warning):` —
// hooks.go renders warnings without blocking.
// ---------------------------------------------------------------------------
func validatePlanTreeFolderConvention(text string, staged []string, root string) string {
	if hasOptOutTrailer(text, "[skip-plan-tree-folder-convention]") {
		return ""
	}
	var warnings []string
	for _, rel := range stagedPlanPaths(staged) {
		segs := strings.Split(rel, "/")
		// segs[0]=plans, segs[1]=active|archived, segs[2..]=...
		if len(segs) < 3 {
			continue
		}
		base := segs[len(segs)-1]
		// Depth check: levels under plans/<active|archived>/.
		depth := len(segs) - 2
		if depth >= 3 {
			warnings = append(warnings, fmt.Sprintf("%s: nested depth %d (recommend < 3, consider splitting into independent main plan)", rel, depth))
		}
		// Filename convention (only for files inside a folder, not top-level
		// active/archived siblings).
		if depth >= 2 {
			if base != "_plan.md" && !planTreeFolderNameRE.MatchString(base) {
				warnings = append(warnings, fmt.Sprintf("%s: filename should be `_plan.md` or `NN-<slug>.md`", rel))
			}
		}
	}
	if len(warnings) == 0 {
		return ""
	}
	return "plan-tree-folder-convention (warning): UI convention advisories (non-blocking):\n    - " +
		strings.Join(warnings, "\n    - ") +
		"\n  Folder shape is a recommendation — frontmatter `parent` is the source of truth." +
		"\n  Opt-out: standalone `[skip-plan-tree-folder-convention]` trailer."
}

// hasOptOutTrailer returns true if the commit message body contains the given
// trailer on a line by itself (case-sensitive, whitespace-trimmed).
func hasOptOutTrailer(text, trailer string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == trailer {
			return true
		}
	}
	return false
}
