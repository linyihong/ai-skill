// Command evidencecandidate is the Evidence Candidate Scanner v0 — an ASSEMBLER.
//
// Design: plans/active/2026-06-16-1131-evidence-candidate-system.md §Phase 1C.
//
// It is Go-first (toolchain placement under scripts/ai-skill-cli/cmd/) but is a
// manual utility: discoverable, NOT routable. Three concepts kept separate:
// implementation language (Go) / toolchain placement (cmd/) / authority (none).
//
// It is a stateless assembler: given an artifact + EXPLICIT criteria_hits
// (annotated outside the scanner) + the rule registry, it assembles a
// well-formed candidate and persists it to the gitignored inbox.
//
//	contract:
//	  input:      artifact + criteria_hits + matched_plans (+ criteria_source.actor)
//	  output:     candidate (stdout line + inbox/<id>.json)
//	  authority:  none
//	  side_effect: inbox only
//	does:     schema validate / pointer resolve / dedupe / invariant check / persist inbox
//	does NOT: infer / match / classify / score / rank / accept / expire
//
// HARD GUARDS (must hold so this never grows back into runtime):
//   - Guard 1 (removable): nothing may depend on this binary — no route.* /
//     runtime.db / build pipeline / commit hook / generated surface. Deleting it
//     must leave Phase 1A/1B intact. NOT registered as an `ai-skill` dispatch
//     target (discoverable != routable).
//   - Guard 2 (output = artifact, not state): side effects are stdout +
//     inbox/<id>.json ONLY. It MUST NOT mutate runtime.db, the registry, any
//     plan, or memory. accept/discard/expire happen elsewhere (human).
//   - Guard 3 (exit code != maturity): exit 0 = assembled, 1 = invalid input.
//     No other codes — never encode "accepted"/"matured" in the exit status.
//
// Invariants: source.artifact must reference an original artifact (not another
// candidate); criteria_hits must originate outside the scanner
// (criteria_source.actor present and not the scanner itself). No confidence is
// produced (Q1 frozen). Output ordering is undefined (emit only).
//
// Usage (from the module dir scripts/ai-skill-cli/):
//
//	go run ./cmd/evidencecandidate -base ../../governance/evidence-candidates < input.json
//
// Stateless: same input -> same candidate id (content hash) -> idempotent write.
package main

import (
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	scannerActors = map[string]bool{"scanner": true, "scanner-v0": true, "self": true}
	candidateIDRe = regexp.MustCompile(`^C-[0-9a-fA-F]{6,}$`)
)

type source struct {
	Repo     string `json:"repo"`
	Artifact string `json:"artifact"`
	Commit   string `json:"commit"`
}

type criteriaSource struct {
	Actor string `json:"actor"`
}

type input struct {
	Source         source         `json:"source"`
	MatchedPlans   []string       `json:"matched_plans"`
	CriteriaHits   []string       `json:"criteria_hits"`
	CriteriaSource criteriaSource `json:"criteria_source"`
}

type candidate struct {
	ID             string         `json:"id"`
	Source         source         `json:"source"`
	MatchedPlans   []string       `json:"matched_plans"`
	CriteriaHits   []string       `json:"criteria_hits"`
	CriteriaSource criteriaSource `json:"criteria_source"`
	Status         string         `json:"status"`
}

// pointerStatus reads a registry pointer's declared status (e.g. "resolved" /
// "section_pending"). Pointer files are pointer-only YAML.
func pointerStatus(path string) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var p struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal(body, &p); err != nil {
		return "", err
	}
	return p.Status, nil
}

// assemble validates the input, enforces invariants, and does a STATUS-AWARE
// pointer resolve. "index != consumable": a pointer's mere existence does not
// make a plan a candidate target. A matched_plan whose pointer is missing is a
// hard reject (error). A pointer whose status != "resolved" (e.g.
// section_pending) is NOT consumable — the candidate is not emitted and the
// plan id is returned in `pending` (caller WARNs, exit 0, no error). No
// inference, matching, ranking, or lifecycle transition (Guard 2/3).
func assemble(in input, base string) (candidate, []string, error) {
	var c candidate
	if in.Source.Repo == "" {
		return c, nil, fmt.Errorf("source.repo required")
	}
	if in.Source.Artifact == "" {
		return c, nil, fmt.Errorf("source.artifact required")
	}
	if len(in.MatchedPlans) == 0 {
		return c, nil, fmt.Errorf("matched_plans must be non-empty")
	}
	if len(in.CriteriaHits) == 0 {
		return c, nil, fmt.Errorf("criteria_hits must be non-empty")
	}
	if in.CriteriaSource.Actor == "" {
		return c, nil, fmt.Errorf("criteria_source.actor required (criteria_hits MUST originate outside scanner)")
	}
	// invariant: criteria_hits originate outside the scanner
	if scannerActors[in.CriteriaSource.Actor] {
		return c, nil, fmt.Errorf("criteria_source.actor=%q is the scanner itself; criteria_hits must come from outside (human / matcher-v2=Phase 2)", in.CriteriaSource.Actor)
	}
	// invariant: source must be an original artifact, not another candidate
	if candidateIDRe.MatchString(in.Source.Artifact) || strings.Contains(in.Source.Artifact, "evidence-candidates/inbox") {
		return c, nil, fmt.Errorf("source.artifact %q looks like a candidate; candidate MUST NOT reference another candidate", in.Source.Artifact)
	}
	// status-aware pointer resolve
	var pending []string
	for _, plan := range in.MatchedPlans {
		ptr := filepath.Join(base, "evidence-rules", plan+".pointer.yaml")
		if _, err := os.Stat(ptr); err != nil {
			return c, nil, fmt.Errorf("no registry pointer for matched_plan %q (%s)", plan, ptr)
		}
		st, err := pointerStatus(ptr)
		if err != nil {
			return c, nil, fmt.Errorf("cannot read pointer status for %q: %v", plan, err)
		}
		if st != "resolved" {
			pending = append(pending, fmt.Sprintf("%s(status=%s)", plan, st))
		}
	}
	if len(pending) > 0 {
		// section_pending (or any non-resolved) pointer is not a candidate target.
		return c, pending, nil
	}
	c.ID = deterministicID(in)
	c.Source = in.Source
	c.MatchedPlans = in.MatchedPlans // order preserved as given; scanner does not rank
	c.CriteriaHits = in.CriteriaHits
	c.CriteriaSource = in.CriteriaSource
	c.Status = "create" // scanner never sets accept/discard/expire
	return c, pending, nil
}

// deterministicID hashes the order-independent content so the same artifact +
// annotations always yield the same candidate id (dedupe / idempotency).
func deterministicID(in input) string {
	mp := append([]string(nil), in.MatchedPlans...)
	ch := append([]string(nil), in.CriteriaHits...)
	sort.Strings(mp)
	sort.Strings(ch)
	basis, _ := json.Marshal(map[string]any{
		"source":        map[string]string{"repo": in.Source.Repo, "artifact": in.Source.Artifact, "commit": in.Source.Commit},
		"matched_plans": mp,
		"criteria_hits": ch,
	})
	sum := sha1.Sum(basis)
	return fmt.Sprintf("C-%x", sum[:4])
}

func main() {
	base := flag.String("base", "governance/evidence-candidates", "path to the evidence-candidates directory")
	flag.Parse()

	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "REJECT: cannot read stdin: %v\n", err)
		os.Exit(1)
	}
	var in input
	if err := json.Unmarshal(raw, &in); err != nil {
		fmt.Fprintf(os.Stderr, "REJECT: input is not valid JSON: %v\n", err)
		os.Exit(1)
	}
	c, pending, err := assemble(in, *base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "REJECT: %v\n", err)
		os.Exit(1)
	}
	if len(pending) > 0 {
		// index != consumable: pointer exists but is not resolved -> not a
		// candidate target. Not an error; warn and exit 0 without emitting.
		fmt.Fprintf(os.Stderr, "WARN: not emitted; matched_plan pointer(s) not resolved (index != consumable): %s\n", strings.Join(pending, ", "))
		return
	}
	inbox := filepath.Join(*base, "inbox")
	if err := os.MkdirAll(inbox, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "REJECT: cannot create inbox: %v\n", err)
		os.Exit(1)
	}
	out := filepath.Join(inbox, c.ID+".json")
	if _, err := os.Stat(out); err == nil {
		fmt.Printf("IDEMPOTENT: %s already in inbox (no duplicate written)\n", c.ID)
		return
	}
	body, _ := json.MarshalIndent(c, "", "  ")
	if err := os.WriteFile(out, append(body, '\n'), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "REJECT: cannot write candidate: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("EMIT: %s -> inbox/%s.json\n", c.ID, c.ID)
	// emit-only; ordering undefined; no ranking/scoring/confidence.
}
