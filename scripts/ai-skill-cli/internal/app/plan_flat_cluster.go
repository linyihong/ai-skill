package app

// Flat multi-file plan cluster detection and folderize migration.
//
// When a main plan at plans/{active|archived}/<slug>.md grows companion
// files sharing the same timestamp-slug prefix (e.g. <slug>-dogfood-evidence.md),
// the UI convention is a folder with _plan.md + NN-<suffix>.md — same shape as
// plan-tree sub-plan folders but companions need not be sub-plans.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var planFlatMainBasenameRE = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-\d{4}-[a-z0-9-]+$`)

type flatPlanCompanion struct {
	RelPath string
	Suffix  string
}

type flatPlanCluster struct {
	Base       string
	Location   string
	MainRel    string
	Companions []flatPlanCompanion
}

func (c flatPlanCluster) memberCount() int {
	return 1 + len(c.Companions)
}

func (c flatPlanCluster) folderRel() string {
	return filepath.ToSlash(filepath.Join("plans", c.Location, c.Base))
}

func scanFlatPlanClusters(root string) []flatPlanCluster {
	var out []flatPlanCluster
	for _, loc := range []string{"active", "archived"} {
		baseDir := filepath.Join(root, "plans", loc)
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			continue
		}
		var basenames []string
		set := map[string]bool{}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
				continue
			}
			if pathContainsFixturesSegment(e.Name()) {
				continue
			}
			base := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
			if !planFlatMainBasenameRE.MatchString(base) {
				continue
			}
			basenames = append(basenames, base)
			set[base] = true
		}
		sort.Strings(basenames)
		seen := map[string]bool{}
		for _, mainBase := range basenames {
			if seen[mainBase] {
				continue
			}
			var companions []flatPlanCompanion
			prefix := mainBase + "-"
			for other := range set {
				if other == mainBase {
					continue
				}
				if strings.HasPrefix(other, prefix) {
					companions = append(companions, flatPlanCompanion{
						RelPath: filepath.ToSlash(filepath.Join("plans", loc, other+".md")),
						Suffix:  strings.TrimPrefix(other, prefix),
					})
				}
			}
			if len(companions) == 0 {
				continue
			}
			sort.Slice(companions, func(i, j int) bool {
				return companions[i].Suffix < companions[j].Suffix
			})
			seen[mainBase] = true
			out = append(out, flatPlanCluster{
				Base:       mainBase,
				Location:   loc,
				MainRel:    filepath.ToSlash(filepath.Join("plans", loc, mainBase+".md")),
				Companions: companions,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].MainRel < out[j].MainRel
	})
	return out
}

func companionTargetName(index int, suffix string) string {
	return fmt.Sprintf("%02d-%s.md", index, suffix)
}

func rewriteClusterLinks(content string, c flatPlanCluster) string {
	mainBase := c.Base + ".md"
	content = strings.ReplaceAll(content, mainBase, "_plan.md")
	for i, comp := range c.Companions {
		oldName := c.Base + "-" + comp.Suffix + ".md"
		content = strings.ReplaceAll(content, oldName, companionTargetName(i+1, comp.Suffix))
	}
	return content
}

func flatClusterWarningsForStaged(staged []string, root string) []string {
	if len(staged) == 0 {
		return nil
	}
	stagedSet := map[string]bool{}
	for _, s := range staged {
		stagedSet[filepath.ToSlash(s)] = true
	}
	var warnings []string
	for _, c := range scanFlatPlanClusters(root) {
		touches := false
		for _, p := range append([]string{c.MainRel}, companionPaths(c)...) {
			if stagedSet[p] {
				touches = true
				break
			}
		}
		if !touches {
			continue
		}
		warnings = append(warnings, fmt.Sprintf(
			"flat multi-file cluster %s (%d files) — move to %s/_plan.md + NN-<suffix>.md (`ai-skill plans folderize --cluster %s --dry-run`)",
			c.Base, c.memberCount(), c.folderRel(), c.Base,
		))
	}
	return warnings
}

func companionPaths(c flatPlanCluster) []string {
	var out []string
	for _, comp := range c.Companions {
		out = append(out, comp.RelPath)
	}
	return out
}
