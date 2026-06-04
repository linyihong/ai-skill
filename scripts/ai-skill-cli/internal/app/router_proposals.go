package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// router_proposals.go implements `ai-skill router proposals` (Workflow
// Activation Engine Phase 6.1): the occurrence-tracking store + promotion state
// machine for route-candidate proposals produced on detector miss.
//
// The store (runtime/router/route-candidate-proposals.yaml) is a DATA file, not
// a projected runtime contract — it is read/written here, never compiled into
// runtime.db.

const (
	proposalStatusAccumulating   = "accumulating"
	proposalStatusReadyForReview = "ready_for_review"
	proposalStatusPromoted       = "promoted"
	proposalStatusRejected       = "rejected"
	proposalStatusStale          = "stale"
)

// promotion thresholds (plan Phase 6.1 promotion rules)
const (
	proposalReadyOccurrenceThreshold = 5
	proposalReadyRecencyDays         = 30
	proposalStaleAgeDays             = 60
	proposalStalePruneAgeDays        = 90
)

// RouteProposal is one pending route-candidate.
type RouteProposal struct {
	CandidateID          string   `yaml:"candidate_id"`
	FirstSeen            string   `yaml:"first_seen"`
	LastSeen             string   `yaml:"last_seen"`
	OccurrenceCount      int      `yaml:"occurrence_count"`
	DetectedCapabilities []string `yaml:"detected_capabilities,omitempty"`
	SuggestedUserSignals []string `yaml:"suggested_user_signals,omitempty"`
	SuggestedRouteType   string   `yaml:"suggested_route_type,omitempty"`
	Status               string   `yaml:"status"`
	RejectedReason       string   `yaml:"rejected_reason,omitempty"`
}

type routeProposalStore struct {
	SchemaVersion int             `yaml:"schema_version"`
	Proposals     []RouteProposal `yaml:"proposals"`
}

func routeProposalsPath(repo string) string {
	return filepath.Join(repo, "runtime", "router", "route-candidate-proposals.yaml")
}

func loadProposalStore(path string) (routeProposalStore, error) {
	store := routeProposalStore{SchemaVersion: 1}
	body, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil // empty store
		}
		return store, err
	}
	if err := yaml.Unmarshal(body, &store); err != nil {
		return store, err
	}
	if store.SchemaVersion == 0 {
		store.SchemaVersion = 1
	}
	return store, nil
}

func saveProposalStore(path string, store routeProposalStore) error {
	if store.SchemaVersion == 0 {
		store.SchemaVersion = 1
	}
	if store.Proposals == nil {
		store.Proposals = []RouteProposal{}
	}
	out, err := yaml.Marshal(store)
	if err != nil {
		return err
	}
	header := "# route-candidate-proposals.yaml — Workflow Activation Engine Phase 6.1\n" +
		"# DATA STORE, NOT a projected runtime contract. Managed by\n" +
		"# `ai-skill router proposals`. [skip-runtime-yaml-projection]\n\n"
	return os.WriteFile(path, append([]byte(header), out...), 0o644)
}

// findProposal returns the index of candidateID, or -1.
func findProposal(store *routeProposalStore, candidateID string) int {
	for i := range store.Proposals {
		if store.Proposals[i].CandidateID == candidateID {
			return i
		}
	}
	return -1
}

// recordProposalOccurrence adds a new proposal or bumps an existing one. This is
// the write path that the detector-miss → Discovery feedback loop calls (the
// hot-hook auto-call is deferred until Discovery graph traversal can supply
// detected_capabilities; until then this is driven via `router proposals
// record` / tests). Terminal-status proposals (promoted/rejected) are not
// re-opened by a bump.
func recordProposalOccurrence(store *routeProposalStore, candidateID string, signals, capabilities []string, routeType string, now time.Time) {
	ts := now.UTC().Format(time.RFC3339)
	if i := findProposal(store, candidateID); i >= 0 {
		p := &store.Proposals[i]
		if p.Status == proposalStatusPromoted || p.Status == proposalStatusRejected {
			return
		}
		p.OccurrenceCount++
		p.LastSeen = ts
		p.SuggestedUserSignals = dedupeStable(append(p.SuggestedUserSignals, signals...))
		p.DetectedCapabilities = dedupeStable(append(p.DetectedCapabilities, capabilities...))
		if p.Status == proposalStatusStale {
			p.Status = proposalStatusAccumulating // re-activated by a fresh hit
		}
		return
	}
	store.Proposals = append(store.Proposals, RouteProposal{
		CandidateID:          candidateID,
		FirstSeen:            ts,
		LastSeen:             ts,
		OccurrenceCount:      1,
		DetectedCapabilities: dedupeStable(capabilities),
		SuggestedUserSignals: dedupeStable(signals),
		SuggestedRouteType:   routeType,
		Status:               proposalStatusAccumulating,
	})
}

// applyProposalLifecycle runs the promotion state machine + stale pruning.
// Returns (transitions, pruned).
func applyProposalLifecycle(store *routeProposalStore, now time.Time) (transitions int, pruned int) {
	kept := store.Proposals[:0]
	for i := range store.Proposals {
		p := store.Proposals[i]
		ageDays := proposalAgeDays(p.LastSeen, now)
		if p.Status == proposalStatusAccumulating {
			switch {
			case p.OccurrenceCount >= proposalReadyOccurrenceThreshold && ageDays <= proposalReadyRecencyDays:
				p.Status = proposalStatusReadyForReview
				transitions++
			case p.OccurrenceCount < proposalReadyOccurrenceThreshold && ageDays > proposalStaleAgeDays:
				p.Status = proposalStatusStale
				transitions++
			}
		}
		// hard-prune stale entries that have been stale a long time
		if p.Status == proposalStatusStale && ageDays > proposalStalePruneAgeDays {
			pruned++
			continue
		}
		kept = append(kept, p)
	}
	store.Proposals = kept
	return transitions, pruned
}

// proposalAgeDays returns whole days between lastSeen (RFC3339) and now. A
// malformed/empty timestamp is treated as very old so it ages out rather than
// lingering forever.
func proposalAgeDays(lastSeen string, now time.Time) int {
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(lastSeen))
	if err != nil {
		return proposalStalePruneAgeDays + 1
	}
	return int(now.UTC().Sub(t.UTC()).Hours() / 24)
}

// runRouter dispatches `ai-skill router <subcommand>`.
func runRouter(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "proposals" {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill router proposals <list|record|promote|reject|gc> [flags]")
		return ExitInvalidUsage
	}
	rest := args[1:]
	if len(rest) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill router proposals <list|record|promote|reject|gc> [flags]")
		return ExitInvalidUsage
	}
	action := rest[0]

	// Extract an optional leading positional (candidate_id) BEFORE flag parsing:
	// Go's flag package stops at the first non-flag token, so flags placed after
	// a positional would be silently ignored. Pulling the positional out first
	// lets `record <id> --signals ...` and `reject <id> --reason ...` work.
	remaining := rest[1:]
	candidateID := ""
	if len(remaining) > 0 && !strings.HasPrefix(remaining[0], "-") {
		candidateID = remaining[0]
		remaining = remaining[1:]
	}

	fs := newFlagSet("router proposals "+action, stderr)
	var repo, statusFilter, reason, signals, capabilities, routeType string
	var jsonOut, plainOut bool
	fs.StringVar(&repo, "repo", ".", "Ai-skill repository path")
	fs.StringVar(&statusFilter, "status", "", "filter by status (list)")
	fs.StringVar(&reason, "reason", "", "rejection reason (reject)")
	fs.StringVar(&signals, "signals", "", "comma-separated suggested user_signals (record)")
	fs.StringVar(&capabilities, "capabilities", "", "comma-separated detected_capabilities (record)")
	fs.StringVar(&routeType, "route-type", "workflow", "suggested route_type (record)")
	fs.BoolVar(&jsonOut, "json", false, "machine-readable JSON output")
	fs.BoolVar(&plainOut, "plain", false, "human-readable output")
	if err := fs.Parse(remaining); err != nil {
		return ExitInvalidUsage
	}

	root, repoCheck := resolveRuntimeObligationsRepo(repo)
	result := Result{Command: "router proposals " + action, Mode: "native", Status: "success", ExitCode: ExitSuccess, Checks: []Check{repoCheck}, PlannedActions: []string{}, Mutations: []string{}}
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message}
		return emitRouterResult(result, jsonOut, stdout, stderr)
	}
	path := routeProposalsPath(root)
	store, err := loadProposalStore(path)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "proposal_store_unreadable", Message: err.Error()}
		return emitRouterResult(result, jsonOut, stdout, stderr)
	}
	now := time.Now().UTC()

	switch action {
	case "list":
		count := 0
		for _, p := range store.Proposals {
			if statusFilter != "" && p.Status != statusFilter {
				continue
			}
			count++
			result.Checks = append(result.Checks, Check{
				Name:    p.CandidateID,
				Status:  "ok",
				Message: fmt.Sprintf("status=%s occurrences=%d last_seen=%s signals=[%s]", p.Status, p.OccurrenceCount, p.LastSeen, strings.Join(p.SuggestedUserSignals, ", ")),
			})
		}
		result.Checks = append(result.Checks, Check{Name: "total", Status: "ok", Message: fmt.Sprintf("%d proposal(s)", count)})
	case "record":
		if candidateID == "" {
			result.Status = "blocked"
			result.ExitCode = ExitInvalidUsage
			result.Error = &CommandError{Code: "missing_candidate_id", Message: "usage: router proposals record <candidate_id> --signals a,b [--capabilities ...] [--route-type ...]"}
			break
		}
		recordProposalOccurrence(&store, candidateID, splitCSV(signals), splitCSV(capabilities), routeType, now)
		if err := saveProposalStore(path, store); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "write_failed", Message: err.Error()}
			break
		}
		result.Mutations = append(result.Mutations, path)
		result.Checks = append(result.Checks, Check{Name: candidateID, Status: "ok", Message: "recorded occurrence"})
	case "promote", "reject":
		if candidateID == "" {
			result.Status = "blocked"
			result.ExitCode = ExitInvalidUsage
			result.Error = &CommandError{Code: "missing_candidate_id", Message: "usage: router proposals " + action + " <candidate_id>"}
			break
		}
		idx := findProposal(&store, candidateID)
		if idx < 0 {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "not_found", Message: "no proposal: " + candidateID}
			break
		}
		if action == "promote" {
			store.Proposals[idx].Status = proposalStatusPromoted
			result.PlannedActions = append(result.PlannedActions,
				"add route to knowledge/runtime/routing-registry.yaml (human/governance review) with suggested user_signals; promotion does NOT auto-edit the canonical registry")
		} else {
			store.Proposals[idx].Status = proposalStatusRejected
			store.Proposals[idx].RejectedReason = reason
		}
		if err := saveProposalStore(path, store); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "write_failed", Message: err.Error()}
			break
		}
		result.Mutations = append(result.Mutations, path)
		result.Checks = append(result.Checks, Check{Name: candidateID, Status: "ok", Message: "status=" + store.Proposals[idx].Status})
	case "gc":
		transitions, pruned := applyProposalLifecycle(&store, now)
		if err := saveProposalStore(path, store); err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitGeneralFailure
			result.Error = &CommandError{Code: "write_failed", Message: err.Error()}
			break
		}
		result.Mutations = append(result.Mutations, path)
		result.Checks = append(result.Checks,
			Check{Name: "transitions", Status: "ok", Message: fmt.Sprintf("%d", transitions)},
			Check{Name: "pruned", Status: "ok", Message: fmt.Sprintf("%d", pruned)},
		)
	default:
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "unknown_action", Message: "unknown action: " + action}
	}

	return emitRouterResult(result, jsonOut, stdout, stderr)
}

func emitRouterResult(result Result, jsonOut bool, stdout, stderr io.Writer) int {
	if jsonOut {
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

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	sort.Strings(out)
	return out
}
