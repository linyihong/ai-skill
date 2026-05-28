package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// RenderMarkdown writes a human-readable markdown report.
func RenderMarkdown(w io.Writer, inv *Inventory) error {
	var sb strings.Builder
	sb.WriteString("# Ai-skill Runtime Audit Report\n\n")
	sb.WriteString(fmt.Sprintf("Repo: `%s`\n\n", inv.Repo))

	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Category | auto-detected | consumed | intentionally-manual | orphan |\n")
	sb.WriteString("|---|---:|---:|---:|---:|\n")
	sb.WriteString(fmt.Sprintf("| Routes | %d | %d | %d | %d |\n",
		inv.Summary.RouteCounts[ClassAutoDetected],
		inv.Summary.RouteCounts[ClassConsumed],
		inv.Summary.RouteCounts[ClassManual],
		inv.Summary.RouteCounts[ClassOrphan]))
	sb.WriteString(fmt.Sprintf("| Surfaces | %d | %d | %d | %d |\n",
		inv.Summary.SurfaceCounts[ClassAutoDetected],
		inv.Summary.SurfaceCounts[ClassConsumed],
		inv.Summary.SurfaceCounts[ClassManual],
		inv.Summary.SurfaceCounts[ClassOrphan]))
	sb.WriteString(fmt.Sprintf("| Scenarios | %d | %d | %d | %d |\n",
		inv.Summary.ScenarioCounts[ClassAutoDetected],
		inv.Summary.ScenarioCounts[ClassConsumed],
		inv.Summary.ScenarioCounts[ClassManual],
		inv.Summary.ScenarioCounts[ClassOrphan]))
	sb.WriteString(fmt.Sprintf("\n**Orphan total**: %d\n\n", inv.Summary.OrphanTotal))

	sb.WriteString("## Routes\n\n")
	sb.WriteString("| ID | Classification | Evidence |\n|---|---|---|\n")
	for _, r := range inv.Routes {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", r.ID, r.Classification, escapePipe(r.Evidence)))
	}

	sb.WriteString("\n## Generated surfaces\n\n")
	sb.WriteString("| target_key | source_path | Classification | Evidence |\n|---|---|---|---|\n")
	for _, s := range inv.Surfaces {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", s.TargetKey, escapePipe(s.SourcePath), s.Classification, escapePipe(s.Evidence)))
	}

	sb.WriteString("\n## Validation scenarios\n\n")
	sb.WriteString("| ID | Path | Classification | Evidence |\n|---|---|---|---|\n")
	for _, sc := range inv.Scenarios {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", sc.ID, escapePipe(sc.Path), sc.Classification, escapePipe(sc.Evidence)))
	}

	if len(inv.Warnings) > 0 {
		sb.WriteString("\n## Warnings\n\n")
		for _, w := range inv.Warnings {
			sb.WriteString("- ")
			sb.WriteString(w)
			sb.WriteString("\n")
		}
	}

	_, err := io.WriteString(w, sb.String())
	return err
}

// RenderJSON writes a machine-readable JSON report.
func RenderJSON(w io.Writer, inv *Inventory) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(inv)
}

func escapePipe(s string) string {
	return strings.ReplaceAll(s, "|", `\|`)
}
