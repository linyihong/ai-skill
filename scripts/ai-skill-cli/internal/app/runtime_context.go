package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// runtime_context.go implements Phase 4.0 of the Workflow Activation Engine:
// the in-memory RuntimeContext derived from the current task's transcript.
//
// Persistence model (important): the PreToolUse hook runs as a fresh process
// per tool call, so there is no live in-memory state shared across calls. Since
// DetectWorkflows is deterministic, RuntimeContext is REBUILT from the
// transcript on each invocation and yields the same result — no store is
// required for correctness. This is exactly why SQLite persistence is deferred
// (Phase 4.1, conditional): nothing in-task needs it. A store only becomes
// necessary for cross-session replay / analytics / multi-agent handoff.
//
// Lifecycle rules (plan Phase 4.0):
//   - substantive(msg): vocabulary-based, NOT length-based.
//   - explicit pivot sentinel → re-detect over post-pivot turns only.
//   - manual-lock sentinel → EffectiveMode=manual-lock, detector yields.
//   - manual-unlock sentinel → restore auto-detection.
//   - NO implicit keyword-drift invalidation (would mis-fire on drill-down).

// ActivationMode is the 5-value activation mode enum (registry
// §activation_mode_spec). manual-lock is runtime-assigned only.
type ActivationMode string

const (
	ModeAlwaysOn   ActivationMode = "always-on"
	ModeAutoDetect ActivationMode = "auto-detect"
	ModeOnDemand   ActivationMode = "on-demand"
	ModeAdvisory   ActivationMode = "advisory"
	ModeManualLock ActivationMode = "manual-lock"
)

// CanActivate reports whether a route in this mode may take activation
// ownership of the task — i.e. become RuntimeContext.ActiveRoute via the
// detector. This mirrors routing-registry.yaml §activation_mode_spec
// capability_matrix `can_activate` and is the INVARIANT behind active_route:
//
//	active_route != "" (detector path)  ⇒  EffectiveMode.CanActivate() == true
//
// Only `auto-detect` activates via the detector. `advisory` is reinforce/
// suggestion ONLY (can_activate=false) — it must never standalone-lock. The
// detector's DetectedRoute.Activated means "an activation signal MATCHED"
// (a detector fact), which is NOT the same as "ACTIVATABLE" (a mode policy);
// conflating the two is the contract violation this guards. `on-demand`
// activates only by explicit user invocation (never through the detector;
// such routes are filtered out upstream), `always-on` is always loaded, and
// `manual-lock` is a user-assigned lock handled on a separate path.
func (m ActivationMode) CanActivate() bool {
	return m == ModeAutoDetect
}

// RuntimeStatus is the RuntimeContext lifecycle status.
type RuntimeStatus string

const (
	StatusNoMatch  RuntimeStatus = "no-match"
	StatusDetected RuntimeStatus = "detected"
	StatusLocked   RuntimeStatus = "locked" // user manual-lock
)

// DetectionSig records which signal axes fired for the active/primary route.
type DetectionSig struct {
	UserSignalHits    []string `json:"user_signal_hits,omitempty"`
	ContextSignalHits []string `json:"context_signal_hits,omitempty"`
	ArtifactReinforce []string `json:"artifact_reinforce,omitempty"` // Phase 2; not an activation axis
}

// RuntimeContext is the in-memory workflow activation state for one task.
type RuntimeContext struct {
	ActiveRoute      string         `json:"active_route"`     // single locked route, "" when none / conflict
	DetectedRoutes   []string       `json:"detected_routes"`  // ALL routes whose signals matched (incl. advisory suggestions)
	CandidateRoutes  []string       `json:"candidate_routes"` // subset of DetectedRoutes whose mode CanActivate() (auto-detect)
	DetectionSource  DetectionSig   `json:"detection_source"` // axes that fired for ActiveRoute
	ActivatedAt      time.Time      `json:"activated_at,omitempty"`
	LastReinforcedAt time.Time      `json:"last_reinforced_at,omitempty"`
	Status           RuntimeStatus  `json:"status"`
	EffectiveMode    ActivationMode `json:"effective_mode,omitempty"`
	// Conflict is true when >1 route activated; ActiveRoute stays "" and the
	// caller routes to workflow/workflow-routing.md Stage 2 (no auto-pick).
	Conflict bool `json:"conflict"`
	// Substantive reflects whether the latest considered user turn carried task
	// vocabulary; detection is only meaningful on a substantive turn.
	Substantive bool `json:"substantive"`
}

// explicit lifecycle sentinels (deterministic substring match, case-insensitive)
var pivotSentinels = []string{"換任務", "現在我要", "換個話題", "new task", "switch to", "改做", "改成做"}
var lockSentinels = []string{"鎖定", "用這個 workflow", "跟我做", "之後都用", "manual lock", "lock workflow", "lock to"}
var unlockSentinels = []string{"回到自動偵測", "unlock", "解鎖", "取消鎖定", "auto-detect again"}

// defaultActionVerbs is the small fixed action-verb set for substantive()
// (domain_nouns are aggregated from the registry at runtime).
var defaultActionVerbs = []string{
	"幫我", "規劃", "寫", "做", "找", "比較", "設計", "評估", "檢查", "修", "分析", "實作", "處理", "建立",
	"plan", "write", "build", "find", "compare", "design", "evaluate", "check", "fix", "analyze", "implement",
}

// aggregateDomainNouns collects every activation user_signal across all
// participating routes — the registry IS the domain-noun vocabulary, so it
// stays in sync automatically (plan Phase 4.0 substantive() definition).
func aggregateDomainNouns(registry runtimeRoutingRegistry) []string {
	seen := map[string]bool{}
	var out []string
	for _, rec := range registry.Records {
		if !detectorModeParticipates(effectiveActivationMode(rec)) {
			continue
		}
		n := normalizeRouteTriggers(rec.ActivationTriggers)
		for _, s := range n.userSignals {
			s = strings.TrimSpace(s)
			if s == "" || seen[strings.ToLower(s)] {
				continue
			}
			seen[strings.ToLower(s)] = true
			out = append(out, s)
		}
	}
	return out
}

// IsSubstantive reports whether a message carries task intent, by vocabulary
// (domain noun OR action verb) — NOT by character count. An 8-char Chinese
// message can be a full task; a 20-char greeting is not.
func IsSubstantive(registry runtimeRoutingRegistry, message string) bool {
	lower := strings.ToLower(message)
	for _, n := range aggregateDomainNouns(registry) {
		if strings.Contains(lower, strings.ToLower(n)) {
			return true
		}
	}
	for _, v := range defaultActionVerbs {
		if strings.Contains(lower, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

// lastSentinelKind returns which of lock/unlock/pivot appeared most recently in
// the transcript (later turns win), or "" if none. Used to resolve lifecycle.
func lastSentinelKind(transcript []DetectorMessage) (kind string, atIndex int) {
	kind = ""
	atIndex = -1
	for i, m := range transcript {
		lower := strings.ToLower(m.Text)
		if containsAny(lower, unlockSentinels) {
			kind, atIndex = "unlock", i
		}
		if containsAny(lower, lockSentinels) {
			kind, atIndex = "lock", i
		}
		if containsAny(lower, pivotSentinels) {
			kind, atIndex = "pivot", i
		}
	}
	return kind, atIndex
}

func containsAny(haystackLower string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(haystackLower, strings.ToLower(n)) {
			return true
		}
	}
	return false
}

// BuildRuntimeContext derives the workflow activation state from a transcript.
// `now` is injected for deterministic timestamps in tests.
func BuildRuntimeContext(registry runtimeRoutingRegistry, transcript []DetectorMessage, openFiles []DetectorFile, now time.Time) RuntimeContext {
	ctx := RuntimeContext{Status: StatusNoMatch, DetectedRoutes: []string{}, CandidateRoutes: []string{}}

	// substantive flag = latest user turn carries task vocabulary
	for i := len(transcript) - 1; i >= 0; i-- {
		if transcript[i].Role != "" && transcript[i].Role != "user" {
			continue
		}
		ctx.Substantive = IsSubstantive(registry, transcript[i].Text)
		break
	}

	// explicit pivot: re-detect over post-pivot turns only
	kind, idx := lastSentinelKind(transcript)
	considered := transcript
	if kind == "pivot" && idx >= 0 {
		considered = transcript[idx:]
	}

	// manual-lock: user explicitly pinned a workflow. Resolve the route by
	// matching the lock turn's text against participating routes' user_signals;
	// lock only when exactly one route matches (deterministic, no guessing).
	if kind == "lock" && idx >= 0 {
		if route, src, ok := resolveManualLock(registry, transcript[idx].Text); ok {
			ctx.ActiveRoute = route
			ctx.EffectiveMode = ModeManualLock
			ctx.Status = StatusLocked
			ctx.DetectedRoutes = []string{route}
			ctx.CandidateRoutes = []string{route} // user-assigned lock is the active candidate
			ctx.DetectionSource = src
			ctx.ActivatedAt = now
			return ctx
		}
	}
	// manual-unlock: fall through to normal auto-detection (no lock applied).

	detected := DetectWorkflows(registry, considered, openFiles)
	// Three layers (per activation_mode_spec):
	//   DetectedRoutes  — every route whose signals matched (incl. advisory).
	//   candidates      — the ACTIVATABLE subset: Activated AND mode.CanActivate().
	//                     advisory routes match but are suggestion-only, never here.
	//   ActiveRoute     — the single selected candidate (none when 0 / conflict).
	var candidates []DetectedRoute
	for _, d := range detected {
		ctx.DetectedRoutes = append(ctx.DetectedRoutes, d.RouteID)
		if d.Activated && ActivationMode(d.EffectiveMode).CanActivate() {
			candidates = append(candidates, d)
			ctx.CandidateRoutes = append(ctx.CandidateRoutes, d.RouteID)
		}
		if len(d.ArtifactReinforce) > 0 {
			ctx.LastReinforcedAt = now
		}
	}
	sort.Strings(ctx.DetectedRoutes)
	sort.Strings(ctx.CandidateRoutes)

	switch len(candidates) {
	case 0:
		// No activatable route. May still carry advisory suggestions in
		// DetectedRoutes, but nothing locks — Status stays no-match.
		ctx.Status = StatusNoMatch
	case 1:
		ctx.Status = StatusDetected
		ctx.ActiveRoute = candidates[0].RouteID
		ctx.EffectiveMode = ActivationMode(candidates[0].EffectiveMode)
		ctx.DetectionSource = DetectionSig{
			UserSignalHits:    candidates[0].UserSignalHits,
			ContextSignalHits: candidates[0].ContextSignalHits,
			ArtifactReinforce: candidates[0].ArtifactReinforce,
		}
		ctx.ActivatedAt = now
	default:
		// conflict: never auto-pick; caller → workflow-routing.md Stage 2
		ctx.Status = StatusDetected
		ctx.Conflict = true
		ctx.ActivatedAt = now
	}

	// Invariant guard (active_route ⇒ CanActivate). ActiveRoute is only ever
	// assigned from `candidates`, which are CanActivate() by construction, so a
	// violation here is a programming error, not bad input. Fail safe: drop the
	// illegal lock rather than enforce a false primary_source gate downstream.
	if ctx.ActiveRoute != "" && ctx.EffectiveMode != ModeManualLock && !ctx.EffectiveMode.CanActivate() {
		ctx.ActiveRoute = ""
		ctx.EffectiveMode = ""
		ctx.Status = StatusNoMatch
		ctx.DetectionSource = DetectionSig{}
	}
	return ctx
}

// resolveManualLock matches a lock-turn's text against participating routes'
// user_signals and returns the single matching route (else ok=false).
func resolveManualLock(registry runtimeRoutingRegistry, lockText string) (string, DetectionSig, bool) {
	lower := strings.ToLower(lockText)
	var matchID string
	var src DetectionSig
	count := 0
	for _, rec := range registry.Records {
		if !detectorModeParticipates(effectiveActivationMode(rec)) {
			continue
		}
		n := normalizeRouteTriggers(rec.ActivationTriggers)
		hits := matchSubstrings(lower, n.userSignals)
		if len(hits) > 0 {
			count++
			matchID = rec.ID
			src = DetectionSig{UserSignalHits: hits}
		}
	}
	if count == 1 {
		return matchID, src, true
	}
	return "", DetectionSig{}, false
}

// extractTranscriptMessages reads a JSONL transcript and returns user +
// assistant turns as DetectorMessage in document order. Mirrors the role/
// content shapes handled by extractAssistantTexts (string content or content
// arrays with {text}). Non-user/assistant rows are skipped.
func extractTranscriptMessages(path string) []DetectorMessage {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var msgs []DetectorMessage
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 2*1024*1024), 2*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		roleField := entry["type"]
		if roleField == nil {
			roleField = entry["role"]
		}
		var role string
		if roleField != nil {
			_ = json.Unmarshal(roleField, &role)
		}
		if role != "user" && role != "assistant" {
			continue
		}
		text := extractEntryText(entry)
		if text != "" {
			msgs = append(msgs, DetectorMessage{Role: role, Text: text})
		}
	}
	return msgs
}

// extractEntryText pulls the text body from a transcript entry whose content is
// either a string or an array of {text} / string items.
func extractEntryText(entry map[string]json.RawMessage) string {
	var chunks []string
	raw, ok := entry["message"]
	if ok {
		var msg map[string]json.RawMessage
		if json.Unmarshal(raw, &msg) == nil {
			if cRaw, ok := msg["content"]; ok {
				chunks = append(chunks, decodeContentChunks(cRaw)...)
			}
		}
	} else if cRaw, ok := entry["content"]; ok {
		chunks = append(chunks, decodeContentChunks(cRaw)...)
	}
	return strings.Join(chunks, "\n")
}

func decodeContentChunks(cRaw json.RawMessage) []string {
	var s string
	if json.Unmarshal(cRaw, &s) == nil {
		return []string{s}
	}
	var items []json.RawMessage
	if json.Unmarshal(cRaw, &items) != nil {
		return nil
	}
	var chunks []string
	for _, item := range items {
		var m map[string]json.RawMessage
		if json.Unmarshal(item, &m) == nil {
			if tRaw, ok := m["text"]; ok {
				var t string
				if json.Unmarshal(tRaw, &t) == nil {
					chunks = append(chunks, t)
				}
			}
			continue
		}
		var str string
		if json.Unmarshal(item, &str) == nil {
			chunks = append(chunks, str)
		}
	}
	return chunks
}

// buildRuntimeWorkflowContextResult implements `ai-skill runtime
// workflow-context --transcript <path>`: it rebuilds the in-memory
// RuntimeContext from a transcript and dumps it (Phase 4.0 deliverable). This
// is a read-only inspection command — no mutations, no persistence.
func buildRuntimeWorkflowContextResult(opts runtimeOptions) Result {
	result := Result{
		Command:        "runtime workflow-context",
		Mode:           "native",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{},
		Mutations:      []string{},
	}
	root, repoCheck := resolveRuntimeObligationsRepo(opts.repoPath)
	result.Checks = append(result.Checks, repoCheck)
	if repoCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitInvalidUsage
		result.Error = &CommandError{Code: "invalid_repo", Message: repoCheck.Message}
		return result
	}
	registry, err := readRuntimeRoutingRegistry(filepath.Join(root, "knowledge", "runtime", "routing-registry.yaml"))
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "registry_unavailable", Message: err.Error()}
		return result
	}

	var transcript []DetectorMessage
	if opts.transcriptPath != "" {
		transcript = extractTranscriptMessages(opts.transcriptPath)
		result.Checks = append(result.Checks, Check{
			Name: "transcript", Status: "ok",
			Message: fmt.Sprintf("%d turns from %s", len(transcript), opts.transcriptPath),
		})
	} else {
		result.Checks = append(result.Checks, Check{
			Name: "transcript", Status: "skipped",
			Message: "no --transcript provided; empty context",
		})
	}

	rc := BuildRuntimeContext(registry, transcript, nil, time.Now().UTC())
	result.Checks = append(result.Checks,
		Check{Name: "status", Status: "ok", Message: string(rc.Status)},
		Check{Name: "active_route", Status: "ok", Message: orNone(rc.ActiveRoute)},
		Check{Name: "effective_mode", Status: "ok", Message: orNone(string(rc.EffectiveMode))},
		Check{Name: "conflict", Status: "ok", Message: fmt.Sprintf("%t", rc.Conflict)},
		Check{Name: "substantive", Status: "ok", Message: fmt.Sprintf("%t", rc.Substantive)},
		Check{Name: "detected_routes", Status: "ok", Message: orNone(strings.Join(rc.DetectedRoutes, ", "))},
		Check{Name: "candidate_routes", Status: "ok", Message: orNone(strings.Join(rc.CandidateRoutes, ", "))},
	)
	if rc.Conflict {
		result.PlannedActions = append(result.PlannedActions,
			"multiple routes activated — resolve via workflow/workflow-routing.md Stage 2 (no auto-pick)")
	}
	return result
}

func orNone(s string) string {
	if strings.TrimSpace(s) == "" {
		return "(none)"
	}
	return s
}
