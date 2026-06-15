package app

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// discovery.go — Workflow Activation Discovery Bridge (Phase A: Light
// Discovery). Plan: plans/active/2026-06-06-1700-workflow-activation-
// discovery-bridge.md.
//
// SCOPE — Phase A only. Phase B (Deep Discovery via PostToolUse:Read
// hijack) deferred to next session.
//
// CRITICAL INVARIANT — scoring is for advisory RANKING, not gating. This
// module MUST NOT be wired into any activation decision path. Activation
// remains the responsibility of detector.go against routing-registry's
// deterministic activation_triggers. See plan §Decision Rationale
// §Non-Goals封印 ("Discovery proposals MUST NOT satisfy activation_triggers").
//
// SOURCE-OF-TRUTH GUARDRAIL — the project overlay scanner produces signal
// FACTS only (cwd-matches-overlay-path, overlay-declares-route-hint).
// Candidate routes are produced by the scoring stage from
// routing-registry.yaml + knowledge/summaries/*.md, never directly from
// overlay metadata.

const (
	discoveryScoringVersion = "light-v1"

	// Status enum (paired with runtime/discovery-bridge.yaml.proposal_statuses).
	discoveryStatusAwaitingPhaseB = "awaiting_phase_b"
	discoveryStatusAdvised        = "advised"
	discoveryStatusDismissed      = "dismissed"
	discoveryStatusRejected       = "rejected"
	discoveryStatusExpired        = "expired"

	// miss_reason enum (paired with runtime/discovery-bridge.yaml.miss_reasons).
	discoveryMissNoArtifact        = "no_artifact_reference"
	discoveryMissInsufficientSig   = "insufficient_signal"
	discoveryMissBelowThreshold    = "confidence_below_threshold"
	discoveryMissCostBudget        = "cost_budget_exceeded"
	discoveryMissManualLockBypass  = "manual_lock_bypass"
	discoveryMissEligibilityFail   = "eligibility_gate_fail"
)

// DiscoveryConfig mirrors runtime/discovery-bridge.yaml. Loaded from the
// projected generated_surfaces[runtime.discovery.config] surface, falling
// back to defaults if the surface is unavailable (hook running outside
// the Ai-skill repo).
type DiscoveryConfig struct {
	PhaseALight struct {
		Threshold              float64            `yaml:"threshold"`
		BudgetP95Ms            int                `yaml:"budget_p95_ms"`
		TopNCandidates         int                `yaml:"top_n_candidates"`
		EvidencePerCandidate   int                `yaml:"evidence_per_candidate_cap"`
		Weights                map[string]float64 `yaml:"weights"`
		// DormantFeatures maps a feature name → reason for being inactive. A
		// dormant feature has a declared (reserved) weight but no live producer.
		// Per the Dormant-Feature Invariant it is excluded from the score
		// denominator, contributes no score, and emits no evidence — see
		// activeWeights + scoreRoute. Keeps the threshold calibrated against
		// only features that actually produce a signal.
		DormantFeatures        map[string]string  `yaml:"dormant_features"`
		ExtensionRouteHints    map[string]string  `yaml:"extension_route_hints"`
	} `yaml:"phase_a_light"`
	Eligibility struct {
		MinUserMsgTokens             int  `yaml:"min_user_msg_tokens"`
		RequireArtifactOrImperative  bool `yaml:"require_artifact_or_imperative"`
		RecentProposalDedupHours     int  `yaml:"recent_proposal_dedup_hours"`
	} `yaml:"eligibility"`
	ProposalStore struct {
		TTLHours         int  `yaml:"ttl_hours"`
		OnWriteEviction  bool `yaml:"on_write_eviction"`
	} `yaml:"proposal_store"`
	Advisory struct {
		MaxTokensPerInject       int  `yaml:"max_tokens_per_inject"`
		CumulativeCapPerSession  int  `yaml:"cumulative_cap_per_session"`
		NonBlocking              bool `yaml:"non_blocking"`
	} `yaml:"advisory"`
}

// defaultDiscoveryConfig returns the fallback config used when the
// projected surface is unavailable. Values must stay aligned with
// runtime/discovery-bridge.yaml.
func defaultDiscoveryConfig() DiscoveryConfig {
	var c DiscoveryConfig
	c.PhaseALight.Threshold = 0.5
	c.PhaseALight.BudgetP95Ms = 30
	c.PhaseALight.TopNCandidates = 3
	c.PhaseALight.EvidencePerCandidate = 5
	c.PhaseALight.Weights = map[string]float64{
		"user_msg_term":    0.30,
		"summary_match":    0.25,
		"basename_term":    0.15,
		"path_segment":     0.10,
		"extension_hint":   0.05,
		"frontmatter_head": 0.10, // reserved weight; dormant until Phase B wires a producer
		"cwd_overlay":      0.05,
	}
	// frontmatter_head has a scorer branch + reserved weight (0.10) but NO
	// producer — buildDiscoveryInputFromTranscript never populates
	// DiscoveryInput.FrontmatterHead, so the branch could never fire in prod.
	// Left un-marked it was a calibration bug: the reserved 0.10 sat in the
	// denominator, capping the effective max score below the declared 1.0 and
	// silently miscalibrating the threshold. Marking it dormant excludes it
	// from the denominator (see activeWeights). Phase B re-enables it by
	// deleting this entry once a producer is connected.
	c.PhaseALight.DormantFeatures = map[string]string{
		"frontmatter_head": "producer_not_connected",
	}
	c.PhaseALight.ExtensionRouteHints = map[string]string{
		".py": "software-delivery", ".ts": "software-delivery", ".tsx": "software-delivery",
		".js": "software-delivery", ".jsx": "software-delivery", ".go": "software-delivery",
		".rb": "software-delivery", ".java": "software-delivery", ".kt": "software-delivery",
		".swift": "software-delivery", ".rs": "software-delivery",
	}
	c.Eligibility.MinUserMsgTokens = 6
	c.Eligibility.RequireArtifactOrImperative = true
	c.Eligibility.RecentProposalDedupHours = 24
	c.ProposalStore.TTLHours = 24
	c.ProposalStore.OnWriteEviction = true
	c.Advisory.MaxTokensPerInject = 200
	c.Advisory.CumulativeCapPerSession = 1000
	c.Advisory.NonBlocking = true
	return c
}

// LoadDiscoveryConfig reads the projected runtime.discovery.config surface
// from runtime.db. Falls back to defaultDiscoveryConfig if the surface or
// database is unavailable. The discovery bridge is fail-open by design —
// config miss must never block tool calls.
func LoadDiscoveryConfig(runtimeDB string) DiscoveryConfig {
	cfg := defaultDiscoveryConfig()
	if runtimeDB == "" {
		return cfg
	}
	db, err := sql.Open("sqlite3", runtimeDB)
	if err != nil {
		return cfg
	}
	defer db.Close()
	var raw string
	err = db.QueryRow(
		"SELECT data FROM generated_surfaces WHERE target_key='runtime.discovery.config' LIMIT 1",
	).Scan(&raw)
	if err != nil || raw == "" {
		return cfg
	}
	// generated_surfaces stores the projected document as JSON; unmarshal
	// into a flexible map and then re-marshal to YAML before parsing into
	// DiscoveryConfig — this keeps the parser tolerant of the projection
	// shape (the compiler may wrap the YAML body in metadata).
	var probe map[string]any
	if jerr := json.Unmarshal([]byte(raw), &probe); jerr == nil {
		body, merr := yaml.Marshal(probe)
		if merr == nil {
			var parsed DiscoveryConfig
			if yerr := yaml.Unmarshal(body, &parsed); yerr == nil {
				mergeDiscoveryConfig(&cfg, &parsed)
			}
		}
	} else if yerr := yaml.Unmarshal([]byte(raw), &cfg); yerr == nil {
		_ = yerr // direct YAML parse succeeded; cfg already populated
	}
	return cfg
}

// mergeDiscoveryConfig overlays non-zero parsed values onto cfg so a
// partial projection cannot silently zero out a default.
func mergeDiscoveryConfig(cfg, parsed *DiscoveryConfig) {
	if parsed.PhaseALight.Threshold > 0 {
		cfg.PhaseALight.Threshold = parsed.PhaseALight.Threshold
	}
	if parsed.PhaseALight.BudgetP95Ms > 0 {
		cfg.PhaseALight.BudgetP95Ms = parsed.PhaseALight.BudgetP95Ms
	}
	if parsed.PhaseALight.TopNCandidates > 0 {
		cfg.PhaseALight.TopNCandidates = parsed.PhaseALight.TopNCandidates
	}
	if parsed.PhaseALight.EvidencePerCandidate > 0 {
		cfg.PhaseALight.EvidencePerCandidate = parsed.PhaseALight.EvidencePerCandidate
	}
	if len(parsed.PhaseALight.Weights) > 0 {
		cfg.PhaseALight.Weights = parsed.PhaseALight.Weights
	}
	if len(parsed.PhaseALight.ExtensionRouteHints) > 0 {
		cfg.PhaseALight.ExtensionRouteHints = parsed.PhaseALight.ExtensionRouteHints
	}
	if parsed.Eligibility.MinUserMsgTokens > 0 {
		cfg.Eligibility.MinUserMsgTokens = parsed.Eligibility.MinUserMsgTokens
	}
	if parsed.Eligibility.RecentProposalDedupHours > 0 {
		cfg.Eligibility.RecentProposalDedupHours = parsed.Eligibility.RecentProposalDedupHours
	}
	if parsed.ProposalStore.TTLHours > 0 {
		cfg.ProposalStore.TTLHours = parsed.ProposalStore.TTLHours
	}
	if parsed.Advisory.MaxTokensPerInject > 0 {
		cfg.Advisory.MaxTokensPerInject = parsed.Advisory.MaxTokensPerInject
	}
	if parsed.Advisory.CumulativeCapPerSession > 0 {
		cfg.Advisory.CumulativeCapPerSession = parsed.Advisory.CumulativeCapPerSession
	}
}

// DiscoveryInput is the cheap pre-Read signal envelope for Light Discovery.
// All fields are extracted from the PreToolUse hook payload + the transcript
// + the Ai-skill repo on disk — no new artifact Reads are issued.
type DiscoveryInput struct {
	UserMessage   string   // most recent user message text
	Basenames     []string // artifact basenames referenced in user_msg or recent Reads
	Paths         []string // full relative paths (Ai-skill or downstream project local)
	Extensions    []string // file extensions associated with referenced artifacts
	FrontmatterHead map[string]string // path → first ≤200B of file content (markdown frontmatter target)
	Cwd           string   // current working directory at hook time
}

// EvidenceItem is one specific signal contribution to a route candidate's
// score. Sanitization: no raw private paths — only basename or last 2 path
// segments are emitted. See plan §Phase A.2 "Per-candidate evidence_set".
type EvidenceItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// RouteCandidate is one ranked Discovery output row.
type RouteCandidate struct {
	Route    string         `json:"route"`     // route.* id from routing-registry
	Score    float64        `json:"score"`
	Evidence []EvidenceItem `json:"evidence"`
}

// DiscoveryProposal is the runtime.db row shape.
type DiscoveryProposal struct {
	ID                    int64
	TaskHash              string
	RouteCandidates       []RouteCandidate
	SignalSnapshot        map[string]any
	ScoringVersion        string
	CurrentBestConfidence float64
	Status                string
	MissReason            string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	ExpiresAt             time.Time
}

// projectOverlayCache is a per-process in-memory cache of overlay signal
// facts keyed by cwd. TTL is 30 minutes to bound staleness.
type projectOverlayCache struct {
	mu      sync.Mutex
	entries map[string]projectOverlayEntry
}

type projectOverlayEntry struct {
	facts     []EvidenceItem
	cachedAt  time.Time
}

const projectOverlayCacheTTL = 30 * time.Minute

var globalProjectOverlayCache = &projectOverlayCache{entries: map[string]projectOverlayEntry{}}

func (c *projectOverlayCache) get(cwd string) ([]EvidenceItem, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[cwd]
	if !ok {
		return nil, false
	}
	if time.Since(entry.cachedAt) > projectOverlayCacheTTL {
		delete(c.entries, cwd)
		return nil, false
	}
	return entry.facts, true
}

func (c *projectOverlayCache) put(cwd string, facts []EvidenceItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[cwd] = projectOverlayEntry{facts: facts, cachedAt: time.Now()}
}

// Invalidate clears the cache. Exported for tests.
func (c *projectOverlayCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = map[string]projectOverlayEntry{}
}

// taskHash hashes the substantive task fingerprint so repeat calls within
// the dedup window can be short-circuited. Hash is sha256-12 of
// trim(user_msg) + cwd + sorted basenames.
func taskHash(input DiscoveryInput) string {
	h := sha256.New()
	h.Write([]byte(strings.TrimSpace(input.UserMessage)))
	h.Write([]byte("|cwd:"))
	h.Write([]byte(input.Cwd))
	bn := append([]string{}, input.Basenames...)
	sort.Strings(bn)
	for _, b := range bn {
		h.Write([]byte("|"))
		h.Write([]byte(b))
	}
	return hex.EncodeToString(h.Sum(nil))[:24]
}

// EligibilityCheck returns ("", true) if Discovery should run; or a
// miss_reason and false if it should short-circuit.
func EligibilityCheck(input DiscoveryInput, cfg DiscoveryConfig) (string, bool) {
	tokens := tokenize(input.UserMessage)
	if len(tokens) < cfg.Eligibility.MinUserMsgTokens {
		return discoveryMissEligibilityFail, false
	}
	if cfg.Eligibility.RequireArtifactOrImperative {
		hasArtifact := len(input.Basenames) > 0 || len(input.Paths) > 0
		hasImperative := containsImperative(input.UserMessage)
		if !hasArtifact && !hasImperative {
			return discoveryMissNoArtifact, false
		}
	}
	return "", true
}

var imperativeRE = regexp.MustCompile(`(?i)\b(review|fix|implement|build|create|add|remove|delete|refactor|test|debug|analyze|plan|design|update|sync|merge|deploy|release|check|verify|check|inspect|investigate|trace|repro)\b|[請幫做改修檢評審分析計設加移除]`)

func containsImperative(msg string) bool {
	return imperativeRE.MatchString(msg)
}

// tokenRE matches alphanumeric runs OR a single CJK character. Per-char
// CJK tokenization reflects that each Han/Kana character carries roughly
// one semantic unit, which keeps min_user_msg_tokens consistent across
// Latin- and CJK-dominant messages.
var tokenRE = regexp.MustCompile(`[A-Za-z0-9_\-]+|\p{Han}|\p{Hiragana}|\p{Katakana}`)

func tokenize(s string) []string {
	if s == "" {
		return nil
	}
	out := tokenRE.FindAllString(s, -1)
	for i, t := range out {
		out[i] = strings.ToLower(t)
	}
	return out
}

// uniqueLower returns the unique lowercase tokens of s, preserving order.
func uniqueLowerTokens(s string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, t := range tokenize(s) {
		if seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// RunLightDiscovery executes Phase A signal extraction + scoring. It does
// NOT write to runtime.db; the caller (PreToolUse hook adapter) is
// responsible for persistence via WriteDiscoveryProposal.
func RunLightDiscovery(input DiscoveryInput, registry runtimeRoutingRegistry, summaries []discoverySummary, cfg DiscoveryConfig, repoRoot string) ([]RouteCandidate, map[string]any) {
	// Signal snapshot — recorded verbatim so future scoring_version
	// re-runs over the same input remain reproducible.
	snapshot := map[string]any{
		"user_message_tokens": uniqueLowerTokens(input.UserMessage),
		"basenames":           input.Basenames,
		"paths":               input.Paths,
		"extensions":          input.Extensions,
		"cwd":                 input.Cwd,
	}

	overlayFacts := scanProjectOverlayFacts(input.Cwd)
	if len(overlayFacts) > 0 {
		snapshot["overlay_facts"] = overlayFacts
	}

	// Telemetry half of the Dormant-Feature Invariant: record which features
	// were inactive (reserved weight, no producer) and the normalized active
	// weight set actually used, so Phase D telemetry can group by feature
	// activity and never mistakes a dormant feature for a zero-scoring one.
	if len(cfg.PhaseALight.DormantFeatures) > 0 {
		snapshot["dormant_features"] = cfg.PhaseALight.DormantFeatures
	}
	snapshot["active_weights"] = activeWeights(cfg)

	userTokens := map[string]bool{}
	for _, t := range tokenize(input.UserMessage) {
		userTokens[t] = true
	}

	basenameTokens := map[string]bool{}
	for _, b := range input.Basenames {
		for _, t := range tokenize(b) {
			basenameTokens[t] = true
		}
	}

	pathSegments := map[string]bool{}
	for _, p := range input.Paths {
		for _, seg := range strings.FieldsFunc(filepath.ToSlash(p), func(r rune) bool {
			return r == '/' || r == '\\'
		}) {
			for _, t := range tokenize(seg) {
				pathSegments[t] = true
			}
		}
	}

	summaryByRoute := map[string]discoverySummary{}
	for _, s := range summaries {
		summaryByRoute[s.AtomID] = s
	}

	candidates := []RouteCandidate{}
	for _, rec := range registry.Records {
		if rec.ID == "" || rec.PrimarySource == "" {
			continue
		}
		score, evidence := scoreRoute(rec, summaryByRoute, userTokens, basenameTokens, pathSegments, input, cfg, overlayFacts)
		if score <= 0 {
			continue
		}
		if len(evidence) > cfg.PhaseALight.EvidencePerCandidate {
			evidence = evidence[:cfg.PhaseALight.EvidencePerCandidate]
		}
		candidates = append(candidates, RouteCandidate{
			Route:    rec.ID,
			Score:    score,
			Evidence: evidence,
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})
	if len(candidates) > cfg.PhaseALight.TopNCandidates {
		candidates = candidates[:cfg.PhaseALight.TopNCandidates]
	}
	return candidates, snapshot
}

// scoreRoute computes the weighted score for one route. Returns
// (score, evidence). Caller discoveryTruncates evidence to cap.
func scoreRoute(rec runtimeRouteRecord, summaryByRoute map[string]discoverySummary,
	userTokens, basenameTokens, pathSegments map[string]bool,
	input DiscoveryInput, cfg DiscoveryConfig, overlayFacts []EvidenceItem) (float64, []EvidenceItem) {

	w := activeWeights(cfg)

	// Route keyword set: task_intent + activation user_signals + summary
	// keywords. Lower-cased and tokenized.
	keywordTokens := map[string]bool{}
	addText := func(s string) {
		for _, t := range tokenize(s) {
			keywordTokens[t] = true
		}
	}
	addText(rec.TaskIntent)
	for _, sig := range rec.ActivationTriggers.UserSignals {
		addText(sig)
	}
	if rec.ActivationTriggers.ActivationAnyOf != nil {
		for _, sig := range rec.ActivationTriggers.ActivationAnyOf.UserSignals {
			addText(sig)
		}
	}
	addText(rec.ID)
	// Route ID like route.workflow.travel-planning → add 'travel', 'planning' tokens.
	for _, frag := range strings.Split(rec.ID, ".") {
		addText(strings.ReplaceAll(frag, "-", " "))
	}

	summary, hasSummary := summaryByRoute[deriveAtomIDFromRoute(rec)]
	if !hasSummary {
		// Try fallback: any summary whose AtomID contains the route's tail.
		tail := routeTail(rec.ID)
		if tail != "" {
			for atomID, s := range summaryByRoute {
				if strings.Contains(atomID, tail) {
					summary = s
					hasSummary = true
					break
				}
			}
		}
	}
	if hasSummary {
		addText(summary.Summary)
		addText(summary.WhenToRead)
	}

	score := 0.0
	evidence := []EvidenceItem{}

	// Every branch is guarded by presence in the active-weights map. A dormant
	// feature is absent from w (activeWeights drops it), so its branch is
	// skipped entirely — no score contribution AND no evidence emitted. This is
	// the Dormant-Feature Invariant: a dormant feature must not appear in the
	// denominator, the score, or the evidence set.

	// user_msg_term — fraction of route keywords present in user msg.
	if wv, ok := w["user_msg_term"]; ok {
		if hits := intersectCount(userTokens, keywordTokens); hits > 0 {
			match := normalizeMatch(hits, len(keywordTokens))
			score += wv * match
			evidence = append(evidence, EvidenceItem{Type: "user_msg_term", Value: fmt.Sprintf("%d/%d", hits, len(keywordTokens))})
		}
	}

	// summary_match — substring presence of summary text in user msg /
	// basenames. Binary full match when at least one notable phrase appears.
	if wv, ok := w["summary_match"]; ok {
		if hasSummary && summarySubstringHit(input, summary) {
			score += wv * 1.0
			evidence = append(evidence, EvidenceItem{Type: "summary_match", Value: discoveryTruncate(summary.Summary, 60)})
		}
	}

	if wv, ok := w["basename_term"]; ok {
		if hits := intersectCount(basenameTokens, keywordTokens); hits > 0 {
			match := normalizeMatch(hits, len(keywordTokens))
			score += wv * match
			var sampleBn string
			if len(input.Basenames) > 0 {
				sampleBn = input.Basenames[0]
			}
			evidence = append(evidence, EvidenceItem{Type: "basename_term", Value: sampleBn})
		}
	}

	if wv, ok := w["path_segment"]; ok {
		if hits := intersectCount(pathSegments, keywordTokens); hits > 0 {
			match := normalizeMatch(hits, len(keywordTokens))
			score += wv * match
			if len(input.Paths) > 0 {
				evidence = append(evidence, EvidenceItem{Type: "path_segment", Value: lastTwoSegments(input.Paths[0])})
			}
		}
	}

	if wv, ok := w["extension_hint"]; ok {
		if hint := extensionHit(rec, input.Extensions, cfg); hint != "" {
			score += wv * 1.0
			evidence = append(evidence, EvidenceItem{Type: "extension_hint", Value: hint})
		}
	}

	// frontmatter_head is dormant (no producer) → absent from w → this whole
	// branch is skipped: frontmatterHit is never even called, so no phantom
	// evidence can be emitted.
	if wv, ok := w["frontmatter_head"]; ok {
		if val := frontmatterHit(rec, input.FrontmatterHead); val != "" {
			score += wv * 1.0
			evidence = append(evidence, EvidenceItem{Type: "frontmatter_head", Value: val})
		}
	}

	if wv, ok := w["cwd_overlay"]; ok {
		if val := overlayHit(rec, overlayFacts); val != "" {
			score += wv * 1.0
			evidence = append(evidence, EvidenceItem{Type: "cwd_overlay", Value: val})
		}
	}

	return score, evidence
}

// activeWeights returns the scoring weights with dormant features removed and
// the remaining (active) weights renormalized to sum to 1.0. Dormant features
// (cfg.PhaseALight.DormantFeatures) carry a reserved weight that is excluded
// from the denominator. This is the mechanical half of the Dormant-Feature
// Invariant: a feature with a reserved weight but no live producer must not sit
// in the denominator (which would cap the effective max score below 1.0 and
// miscalibrate the threshold). Re-enabling a feature = delete it from
// DormantFeatures; the threshold philosophy is unchanged.
func activeWeights(cfg DiscoveryConfig) map[string]float64 {
	raw := cfg.PhaseALight.Weights
	if raw == nil {
		raw = defaultDiscoveryConfig().PhaseALight.Weights
	}
	dormant := cfg.PhaseALight.DormantFeatures
	sum := 0.0
	for k, v := range raw {
		if _, off := dormant[k]; off {
			continue
		}
		sum += v
	}
	out := map[string]float64{}
	if sum <= 0 {
		return out
	}
	for k, v := range raw {
		if _, off := dormant[k]; off {
			continue
		}
		out[k] = v / sum
	}
	return out
}

func intersectCount(a, b map[string]bool) int {
	n := 0
	if len(a) > len(b) {
		a, b = b, a
	}
	for k := range a {
		if b[k] {
			n++
		}
	}
	return n
}

func normalizeMatch(hits, total int) float64 {
	if total <= 0 {
		return 0
	}
	// Square-root scale: saturates quickly so a single strong hit on a
	// small keyword set still contributes meaningfully.
	ratio := float64(hits) / float64(total)
	if ratio > 1.0 {
		ratio = 1.0
	}
	return ratio
}

func summarySubstringHit(input DiscoveryInput, s discoverySummary) bool {
	low := strings.ToLower(input.UserMessage)
	if s.Summary != "" {
		for _, phrase := range significantPhrases(s.Summary) {
			if strings.Contains(low, phrase) {
				return true
			}
		}
	}
	if s.WhenToRead != "" {
		for _, phrase := range significantPhrases(s.WhenToRead) {
			if strings.Contains(low, phrase) {
				return true
			}
		}
	}
	for _, bn := range input.Basenames {
		bnLow := strings.ToLower(bn)
		for _, phrase := range significantPhrases(s.Summary) {
			if strings.Contains(bnLow, phrase) {
				return true
			}
		}
	}
	return false
}

// phraseRE picks multi-char tokens suitable as substring probes. CJK runs
// of ≥2 chars are kept (high semantic density), Latin runs of ≥4 chars are
// kept. Single Han chars and short Latin words like "the" are excluded —
// they false-match too easily.
var phraseRE = regexp.MustCompile(`[A-Za-z0-9_\-]{4,}|[\p{Han}\p{Hiragana}\p{Katakana}]{2,}`)

// significantPhrases extracts probe phrases from the summary text. Used by
// summary_match for substring presence checks against the user message.
func significantPhrases(s string) []string {
	if s == "" {
		return nil
	}
	matches := phraseRE.FindAllString(s, -1)
	out := []string{}
	seen := map[string]bool{}
	for _, t := range matches {
		t = strings.ToLower(t)
		if seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

func extensionHit(rec runtimeRouteRecord, exts []string, cfg DiscoveryConfig) string {
	if len(exts) == 0 || len(cfg.PhaseALight.ExtensionRouteHints) == 0 {
		return ""
	}
	for _, ext := range exts {
		hint, ok := cfg.PhaseALight.ExtensionRouteHints[strings.ToLower(ext)]
		if !ok {
			continue
		}
		if strings.Contains(rec.PrimarySource, hint) || strings.Contains(rec.ID, hint) {
			return ext + "→" + hint
		}
	}
	return ""
}

func frontmatterHit(rec runtimeRouteRecord, headByPath map[string]string) string {
	if len(headByPath) == 0 {
		return ""
	}
	tail := routeTail(rec.ID)
	if tail == "" {
		return ""
	}
	for path, head := range headByPath {
		low := strings.ToLower(head)
		if strings.Contains(low, tail) {
			return filepath.Base(path) + ":" + tail
		}
	}
	return ""
}

func overlayHit(rec runtimeRouteRecord, facts []EvidenceItem) string {
	tail := routeTail(rec.ID)
	if tail == "" {
		return ""
	}
	for _, f := range facts {
		if strings.Contains(strings.ToLower(f.Value), tail) {
			return f.Value
		}
	}
	return ""
}

func routeTail(routeID string) string {
	parts := strings.Split(routeID, ".")
	if len(parts) == 0 {
		return ""
	}
	return strings.ToLower(parts[len(parts)-1])
}

func deriveAtomIDFromRoute(rec runtimeRouteRecord) string {
	// routing-registry route IDs look like `route.workflow.travel-planning`;
	// summaries use `skill.travel-planning` / `governance.<...>` /
	// `runtime.<...>`. Try the last segment as a fallback key.
	parts := strings.Split(rec.ID, ".")
	if len(parts) == 0 {
		return rec.ID
	}
	return parts[len(parts)-1]
}

func discoveryTruncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func lastTwoSegments(p string) string {
	p = filepath.ToSlash(p)
	parts := strings.Split(p, "/")
	if len(parts) <= 2 {
		return p
	}
	return strings.Join(parts[len(parts)-2:], "/")
}

// discoverySummary is the subset of knowledge/summaries/*.md fields used
// by Light Discovery. Defined locally to avoid touching runtime.go's
// runtimeSummaryRecord (which currently omits WhenToRead).
type discoverySummary struct {
	File       string
	AtomID     string
	Summary    string
	WhenToRead string
}

// LoadDiscoverySummaries reads knowledge/summaries/*.md and returns the
// fields Discovery cares about.
func LoadDiscoverySummaries(repo string) []discoverySummary {
	if repo == "" {
		return nil
	}
	paths, err := filepath.Glob(filepath.Join(repo, "knowledge", "summaries", "*.md"))
	if err != nil {
		return nil
	}
	sort.Strings(paths)
	out := []discoverySummary{}
	for _, p := range paths {
		if filepath.Base(p) == "README.md" {
			continue
		}
		content, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		fields := parseRuntimeSummaryTable(string(content))
		rel, _ := filepath.Rel(repo, p)
		out = append(out, discoverySummary{
			File:       filepath.ToSlash(rel),
			AtomID:     strings.ReplaceAll(fields["Atom ID"], "`", ""),
			Summary:    fields["Summary"],
			WhenToRead: fields["When to read"],
		})
	}
	return out
}

// scanProjectOverlayFacts looks for project-local rules under
// <cwd>/.ai-skill/project/rules/*.md and extracts signal FACTS (NOT
// candidate routes) from frontmatter + title + first paragraph. Empty /
// missing overlay returns nil — Discovery is defensive about downstream
// project layouts.
//
// SOURCE-OF-TRUTH GUARDRAIL — this function MUST NOT emit route IDs.
// Facts are opaque tokens consumed by scoreRoute via overlayHit().
func scanProjectOverlayFacts(cwd string) []EvidenceItem {
	if cwd == "" {
		return nil
	}
	if cached, ok := globalProjectOverlayCache.get(cwd); ok {
		return cached
	}
	dir := filepath.Join(cwd, ".ai-skill", "project", "rules")
	entries, err := os.ReadDir(dir)
	if err != nil {
		globalProjectOverlayCache.put(cwd, nil)
		return nil
	}
	facts := []EvidenceItem{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		body, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		head := string(body)
		if len(head) > 2000 {
			head = head[:2000]
		}
		// Facts: frontmatter `tags:` / `kind:` / `domain:` lines; first
		// `# heading`. We do NOT extract `route:` because that would
		// cross the source-of-truth guardrail.
		for _, line := range strings.Split(head, "\n") {
			trim := strings.TrimSpace(line)
			low := strings.ToLower(trim)
			switch {
			case strings.HasPrefix(low, "tags:"), strings.HasPrefix(low, "kind:"), strings.HasPrefix(low, "domain:"):
				facts = append(facts, EvidenceItem{Type: "overlay_frontmatter", Value: trim})
			case strings.HasPrefix(trim, "# "):
				facts = append(facts, EvidenceItem{Type: "overlay_heading", Value: strings.TrimPrefix(trim, "# ")})
			}
		}
	}
	globalProjectOverlayCache.put(cwd, facts)
	return facts
}

// FindRecentProposal returns the most recent non-expired proposal for the
// given task_hash within the dedup window, or zero-value if none.
func FindRecentProposal(runtimeDB, taskHash string, dedupHours int) (DiscoveryProposal, bool) {
	if runtimeDB == "" || taskHash == "" {
		return DiscoveryProposal{}, false
	}
	db, err := sql.Open("sqlite3", runtimeDB)
	if err != nil {
		return DiscoveryProposal{}, false
	}
	defer db.Close()
	cutoff := time.Now().UTC().Add(-time.Duration(dedupHours) * time.Hour).Format(time.RFC3339)
	row := db.QueryRow(`SELECT id, task_hash, route_candidates_json, signal_snapshot_json,
		scoring_version, current_best_confidence, status, COALESCE(miss_reason,''),
		created_at, updated_at, expires_at
		FROM discovery_proposals
		WHERE task_hash=? AND updated_at >= ? AND status != ?
		ORDER BY updated_at DESC LIMIT 1`, taskHash, cutoff, discoveryStatusExpired)
	p := DiscoveryProposal{}
	var candJSON, sigJSON string
	var created, updated, expires string
	if err := row.Scan(&p.ID, &p.TaskHash, &candJSON, &sigJSON, &p.ScoringVersion,
		&p.CurrentBestConfidence, &p.Status, &p.MissReason,
		&created, &updated, &expires); err != nil {
		return DiscoveryProposal{}, false
	}
	_ = json.Unmarshal([]byte(candJSON), &p.RouteCandidates)
	_ = json.Unmarshal([]byte(sigJSON), &p.SignalSnapshot)
	p.CreatedAt, _ = time.Parse(time.RFC3339, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	p.ExpiresAt, _ = time.Parse(time.RFC3339, expires)
	return p, true
}

// WriteDiscoveryProposal inserts a new proposal row. If on-write eviction
// is enabled, expired rows for the same task_hash are DELETEd first.
func WriteDiscoveryProposal(runtimeDB string, p *DiscoveryProposal, cfg DiscoveryConfig) error {
	if runtimeDB == "" {
		return fmt.Errorf("runtime db path empty")
	}
	db, err := sql.Open("sqlite3", runtimeDB)
	if err != nil {
		return err
	}
	defer db.Close()
	if cfg.ProposalStore.OnWriteEviction {
		now := time.Now().UTC().Format(time.RFC3339)
		_, _ = db.Exec(`DELETE FROM discovery_proposals WHERE task_hash=? AND expires_at < ?`, p.TaskHash, now)
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = p.CreatedAt
	}
	if p.ExpiresAt.IsZero() {
		ttl := cfg.ProposalStore.TTLHours
		if ttl <= 0 {
			ttl = 24
		}
		p.ExpiresAt = p.CreatedAt.Add(time.Duration(ttl) * time.Hour)
	}
	if p.ScoringVersion == "" {
		p.ScoringVersion = discoveryScoringVersion
	}
	candJSON, err := json.Marshal(p.RouteCandidates)
	if err != nil {
		return err
	}
	sigJSON, err := json.Marshal(p.SignalSnapshot)
	if err != nil {
		return err
	}
	res, err := db.Exec(`INSERT INTO discovery_proposals
		(task_hash, route_candidates_json, signal_snapshot_json, scoring_version,
		 current_best_confidence, status, miss_reason, created_at, updated_at, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.TaskHash, string(candJSON), string(sigJSON), p.ScoringVersion,
		p.CurrentBestConfidence, p.Status, nullableString(p.MissReason),
		p.CreatedAt.Format(time.RFC3339), p.UpdatedAt.Format(time.RFC3339), p.ExpiresAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	p.ID = id
	return nil
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// RunDiscoveryBridge is the high-level entry called from the PreToolUse
// hook adapter when the detector misses. It returns the advisory text
// (possibly empty) plus the resulting proposal status, and persists the
// proposal row when applicable.
//
// Returns (advisoryText, proposal, error). advisoryText is "" when no
// advisory should be injected (status != advised). Discovery Bridge is
// fail-open: errors do not propagate to the hook caller as blocks.
func RunDiscoveryBridge(input DiscoveryInput, repoRoot, runtimeDB string, manualLockActive bool) (string, DiscoveryProposal, error) {
	cfg := LoadDiscoveryConfig(runtimeDB)
	p := DiscoveryProposal{
		TaskHash:       taskHash(input),
		ScoringVersion: discoveryScoringVersion,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	p.ExpiresAt = p.CreatedAt.Add(time.Duration(cfg.ProposalStore.TTLHours) * time.Hour)

	if manualLockActive {
		p.Status = discoveryStatusExpired
		p.MissReason = discoveryMissManualLockBypass
		return "", p, nil
	}

	// Dedup: if we already have a recent proposal for this task_hash, skip
	// regeneration. Caller can still inject the cached advisory text.
	if existing, ok := FindRecentProposal(runtimeDB, p.TaskHash, cfg.Eligibility.RecentProposalDedupHours); ok {
		if existing.Status == discoveryStatusAdvised {
			return renderAdvisory(existing.RouteCandidates, repoRoot, cfg), existing, nil
		}
		return "", existing, nil
	}

	if reason, ok := EligibilityCheck(input, cfg); !ok {
		p.Status = discoveryStatusExpired
		p.MissReason = reason
		_ = WriteDiscoveryProposal(runtimeDB, &p, cfg)
		return "", p, nil
	}

	registry, err := readRuntimeRoutingRegistry(filepath.Join(repoRoot, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		// fail-open: leave proposal unrecorded, no advisory
		return "", p, nil
	}
	summaries := LoadDiscoverySummaries(repoRoot)

	candidates, snapshot := RunLightDiscovery(input, registry, summaries, cfg, repoRoot)
	p.RouteCandidates = candidates
	p.SignalSnapshot = snapshot
	if len(candidates) == 0 {
		p.Status = discoveryStatusExpired
		p.MissReason = discoveryMissInsufficientSig
		_ = WriteDiscoveryProposal(runtimeDB, &p, cfg)
		return "", p, nil
	}
	p.CurrentBestConfidence = candidates[0].Score
	if p.CurrentBestConfidence >= cfg.PhaseALight.Threshold {
		p.Status = discoveryStatusAdvised
	} else {
		p.Status = discoveryStatusAwaitingPhaseB
		p.MissReason = discoveryMissBelowThreshold
	}
	if err := WriteDiscoveryProposal(runtimeDB, &p, cfg); err != nil {
		// Persistence failure is fail-open — still return advisory if eligible.
		// Caller can rely on the in-memory proposal.
	}
	if p.Status != discoveryStatusAdvised {
		return "", p, nil
	}
	return renderAdvisory(candidates, repoRoot, cfg), p, nil
}

// renderAdvisory produces the PreToolUse advisory text shown to the agent.
// Capped to cfg.Advisory.MaxTokensPerInject (roughly word-count, not
// tokenizer-perfect).
func renderAdvisory(candidates []RouteCandidate, repoRoot string, cfg DiscoveryConfig) string {
	if len(candidates) == 0 {
		return ""
	}
	registry, _ := readRuntimeRoutingRegistry(filepath.Join(repoRoot, "knowledge", "runtime", "routing-registry.yaml"))
	primaryByRoute := map[string]string{}
	for _, rec := range registry.Records {
		primaryByRoute[rec.ID] = rec.PrimarySource
	}
	var b strings.Builder
	b.WriteString("[ai-skill Discovery Bridge — advisory, non-blocking]\n")
	b.WriteString("Detector did not lock a route, but Light Discovery suggests these workflows may apply.\n")
	b.WriteString("Reading the primary_source listed below is OPTIONAL — Discovery never gates tool calls.\n\n")
	for i, c := range candidates {
		ps := primaryByRoute[c.Route]
		b.WriteString(fmt.Sprintf("  %d. %s  (score=%.2f)\n", i+1, c.Route, c.Score))
		if ps != "" {
			b.WriteString("     primary_source: " + ps + "\n")
		}
		if len(c.Evidence) > 0 {
			ev := []string{}
			for _, e := range c.Evidence {
				ev = append(ev, e.Type)
			}
			b.WriteString("     evidence: " + strings.Join(ev, ", ") + "\n")
		}
	}
	// Cap by approximate word count (rough proxy for tokens).
	out := b.String()
	if cfg.Advisory.MaxTokensPerInject > 0 {
		words := strings.Fields(out)
		if len(words) > cfg.Advisory.MaxTokensPerInject {
			out = strings.Join(words[:cfg.Advisory.MaxTokensPerInject], " ") + " …[discoveryTruncated]"
		}
	}
	return out
}
