// Package audit implements the `ai-skill runtime audit` subcommand.
//
// The audit walks four runtime surfaces and classifies each entry into one of
// four buckets: auto-detected (pulled by a discovery signal), consumed (used
// by a commit-msg validator / hook), intentionally manual (explicit
// manual_activation annotation), or orphan (no consumer).
//
// Plan: plans/active/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md
// Phase: 2 (Inventory Tool)
package audit

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Classification buckets for routes / surfaces / scenarios.
const (
	ClassAutoDetected = "auto-detected"
	ClassConsumed     = "consumed"
	ClassManual       = "intentionally-manual"
	ClassOrphan       = "orphan"
)

// Inventory is the full audit snapshot.
type Inventory struct {
	Repo      string         `json:"repo"`
	Routes    []RouteEntry   `json:"routes"`
	Surfaces  []SurfaceEntry `json:"surfaces"`
	Scenarios []ScenarioEntry `json:"scenarios"`
	Summary   InventorySummary `json:"summary"`
	Warnings  []string       `json:"warnings,omitempty"`
}

// InventorySummary counts entries per classification per category.
type InventorySummary struct {
	RouteCounts    map[string]int `json:"route_counts"`
	SurfaceCounts  map[string]int `json:"surface_counts"`
	ScenarioCounts map[string]int `json:"scenario_counts"`
	OrphanTotal    int            `json:"orphan_total"`
}

// RouteEntry classifies a routing-registry record.
type RouteEntry struct {
	ID             string `json:"id"`
	Classification string `json:"classification"`
	Evidence       string `json:"evidence"`
}

// SurfaceEntry classifies a runtime.db generated_surfaces row.
type SurfaceEntry struct {
	TargetKey      string `json:"target_key"`
	SourcePath     string `json:"source_path"`
	Classification string `json:"classification"`
	Evidence       string `json:"evidence"`
}

// ScenarioEntry classifies a validation/scenarios/failure-derived/*.yaml file.
type ScenarioEntry struct {
	ID             string `json:"id"`
	Path           string `json:"path"`
	Classification string `json:"classification"`
	Evidence       string `json:"evidence"`
}

// Options controls audit behavior.
type Options struct {
	Repo string
}

// Build runs the full audit and returns the classified inventory.
func Build(opts Options) (*Inventory, error) {
	repo := opts.Repo
	if repo == "" {
		repo = "."
	}
	abs, err := filepath.Abs(repo)
	if err != nil {
		return nil, fmt.Errorf("resolve repo path: %w", err)
	}

	inv := &Inventory{
		Repo: abs,
		Summary: InventorySummary{
			RouteCounts:    map[string]int{},
			SurfaceCounts:  map[string]int{},
			ScenarioCounts: map[string]int{},
		},
	}

	signalDescriptions, manualAnnotations, err := readSignalsAndManual(abs)
	if err != nil {
		return nil, err
	}

	hooksSource, err := readSourceFiles(filepath.Join(abs, "scripts", "ai-skill-cli", "internal", "app"))
	if err != nil {
		return nil, err
	}
	auditSource, _ := readSourceFiles(filepath.Join(abs, "scripts", "ai-skill-cli", "internal", "audit"))
	cliSource := hooksSource + "\n" + auditSource

	routes, err := readRouteIDs(filepath.Join(abs, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		return nil, err
	}
	for _, id := range routes {
		entry := classifyRoute(id, signalDescriptions, manualAnnotations, cliSource)
		inv.Routes = append(inv.Routes, entry)
		inv.Summary.RouteCounts[entry.Classification]++
	}

	surfaces, err := readGeneratedSurfaces(filepath.Join(abs, "runtime", "runtime.db"))
	if err != nil {
		return nil, err
	}
	for _, s := range surfaces {
		entry := classifySurface(s, cliSource)
		inv.Surfaces = append(inv.Surfaces, entry)
		inv.Summary.SurfaceCounts[entry.Classification]++
	}

	scenarios, err := readScenarios(filepath.Join(abs, "validation", "scenarios"))
	if err != nil {
		return nil, err
	}
	for _, sc := range scenarios {
		entry := classifyScenario(sc, cliSource, abs)
		inv.Scenarios = append(inv.Scenarios, entry)
		inv.Summary.ScenarioCounts[entry.Classification]++
	}

	inv.Summary.OrphanTotal =
		inv.Summary.RouteCounts[ClassOrphan] +
			inv.Summary.SurfaceCounts[ClassOrphan] +
			inv.Summary.ScenarioCounts[ClassOrphan]

	sort.Slice(inv.Routes, func(i, j int) bool { return inv.Routes[i].ID < inv.Routes[j].ID })
	sort.Slice(inv.Surfaces, func(i, j int) bool { return inv.Surfaces[i].TargetKey < inv.Surfaces[j].TargetKey })
	sort.Slice(inv.Scenarios, func(i, j int) bool { return inv.Scenarios[i].ID < inv.Scenarios[j].ID })

	return inv, nil
}

// ── Readers ─────────────────────────────────────────────────────────────────

type routingRegistry struct {
	Records []struct {
		ID                 string `yaml:"id"`
		ManualActivation   *struct {
			Reason string `yaml:"reason"`
		} `yaml:"manual_activation"`
	} `yaml:"records"`
}

// readRouteIDs returns route ids and tracks which have manual_activation set
// via the side-channel `manualAnnotations` returned from readSignalsAndManual.
func readRouteIDs(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read routing-registry: %w", err)
	}
	var reg routingRegistry
	if err := yaml.Unmarshal(raw, &reg); err != nil {
		return nil, fmt.Errorf("parse routing-registry: %w", err)
	}
	ids := make([]string, 0, len(reg.Records))
	for _, r := range reg.Records {
		if r.ID != "" {
			ids = append(ids, r.ID)
		}
	}
	return ids, nil
}

type discoveryDoc struct {
	Signals []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Pattern     string `yaml:"pattern"`
	} `yaml:"signals"`
}

// readSignalsAndManual extracts (1) the combined signal description text used
// to detect auto-detected routes via string mention, and (2) the set of route
// ids carrying an explicit manual_activation annotation.
func readSignalsAndManual(repo string) (signalText string, manualSet map[string]string, err error) {
	manualSet = map[string]string{}

	discoveryPath := filepath.Join(repo, "runtime", "cognitive-modes-discovery.yaml")
	raw, derr := os.ReadFile(discoveryPath)
	if derr == nil {
		var doc discoveryDoc
		if uerr := yaml.Unmarshal(raw, &doc); uerr == nil {
			var sb strings.Builder
			for _, s := range doc.Signals {
				sb.WriteString(s.Description)
				sb.WriteString("\n")
				sb.WriteString(s.Pattern)
				sb.WriteString("\n")
			}
			signalText = sb.String()
		}
	}

	regRaw, rerr := os.ReadFile(filepath.Join(repo, "knowledge", "runtime", "routing-registry.yaml"))
	if rerr == nil {
		var reg routingRegistry
		if err := yaml.Unmarshal(regRaw, &reg); err == nil {
			for _, r := range reg.Records {
				if r.ManualActivation != nil && r.ManualActivation.Reason != "" {
					manualSet[r.ID] = r.ManualActivation.Reason
				}
			}
		}
	}
	return signalText, manualSet, nil
}

func readSourceFiles(dir string) (string, error) {
	var sb strings.Builder
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".go") {
			return nil
		}
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return nil
		}
		sb.Write(b)
		sb.WriteString("\n")
		return nil
	})
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

type rawSurface struct {
	TargetKey  string
	SourcePath string
}

func readGeneratedSurfaces(dbPath string) ([]rawSurface, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("runtime.db missing: %w", err)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open runtime.db: %w", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT target_key, source_path FROM generated_surfaces")
	if err != nil {
		return nil, fmt.Errorf("query generated_surfaces: %w", err)
	}
	defer rows.Close()
	var out []rawSurface
	for rows.Next() {
		var s rawSurface
		if err := rows.Scan(&s.TargetKey, &s.SourcePath); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

type rawScenario struct {
	ID   string
	Path string
}

var scenarioIDRe = regexp.MustCompile(`(?m)^id:\s*(.+?)\s*$`)

func readScenarios(scenariosDir string) ([]rawScenario, error) {
	var out []rawScenario
	err := filepath.Walk(scenariosDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".yaml") && !strings.HasSuffix(p, ".yml") {
			return nil
		}
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return nil
		}
		m := scenarioIDRe.FindSubmatch(b)
		if len(m) < 2 {
			return nil
		}
		out = append(out, rawScenario{ID: string(m[1]), Path: p})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ── Classifiers ─────────────────────────────────────────────────────────────

func classifyRoute(id, signalText string, manual map[string]string, cliSource string) RouteEntry {
	if reason, ok := manual[id]; ok {
		return RouteEntry{ID: id, Classification: ClassManual, Evidence: "manual_activation reason: " + reason}
	}
	if signalText != "" && strings.Contains(signalText, id) {
		return RouteEntry{ID: id, Classification: ClassAutoDetected, Evidence: "discovery signal description references " + id}
	}
	if cliSource != "" && strings.Contains(cliSource, id) {
		return RouteEntry{ID: id, Classification: ClassConsumed, Evidence: "Go source under scripts/ai-skill-cli/ references " + id}
	}
	return RouteEntry{ID: id, Classification: ClassOrphan, Evidence: "no signal / Go consumer / manual_activation"}
}

func classifySurface(s rawSurface, cliSource string) SurfaceEntry {
	if cliSource != "" && strings.Contains(cliSource, s.TargetKey) {
		return SurfaceEntry{TargetKey: s.TargetKey, SourcePath: s.SourcePath, Classification: ClassConsumed, Evidence: "Go source references target_key"}
	}
	return SurfaceEntry{TargetKey: s.TargetKey, SourcePath: s.SourcePath, Classification: ClassOrphan, Evidence: "no Go consumer references target_key"}
}

func classifyScenario(s rawScenario, cliSource, repo string) ScenarioEntry {
	rel, _ := filepath.Rel(repo, s.Path)
	if cliSource != "" && strings.Contains(cliSource, s.ID) {
		return ScenarioEntry{ID: s.ID, Path: rel, Classification: ClassConsumed, Evidence: "Go source references scenario id"}
	}
	return ScenarioEntry{ID: s.ID, Path: rel, Classification: ClassOrphan, Evidence: "no Go consumer references scenario id"}
}
