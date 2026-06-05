package app

import (
	"regexp"
	"sort"
	"strings"
)

// detector.go implements the deterministic Workflow Activation Engine detector
// (plan plans/archived/2026-05-31-1900-workflow-activation-engine.md Phase 3).
//
// Design invariants (do NOT change without updating the plan + ADR-006):
//   - Deterministic: a route is detected iff at least one of its activation
//     signals matches. NO weighted scoring, NO thresholds, NO ranking.
//   - Two-phase: activation_any_of signals (pre-Read user/context) can ACTIVATE
//     a route; reinforcement_any_of signals (post-Read artifact content) only
//     REINFORCE — they never activate a route on their own. A reinforcement-only
//     hit is surfaced as a "late-detected" candidate (Activated=false) for
//     trigger-coverage analysis, but the detector does not lock on it.
//   - Backward compatible: legacy flat triggers (top-level user_signals /
//     file_change_globs) are normalized into the two-phase model.
//
// This file is a pure function over an already-parsed registry; it performs no
// file IO. Hook integration (PreToolUse dedupe via in-memory RuntimeContext)
// lands with Phase 4.0.

// DetectorMessage is one transcript turn fed to the detector. Only Text is
// matched; Role is retained for callers that want to restrict matching to user
// turns (see DetectWorkflows docs).
type DetectorMessage struct {
	Role string
	Text string
}

// DetectorFile is an open / referenced file. Path is matched against
// context_signals (pre-Read globs); Content is matched against
// reinforcement artifact_signals (post-Read). Content may be empty when only
// the path is known (pre-Read).
type DetectorFile struct {
	Path    string
	Content string
}

// DetectedRoute is one route the detector matched.
type DetectedRoute struct {
	RouteID       string
	EffectiveMode string // resolved activation_mode (explicit or derived from route_type)
	// Activated is true when at least one activation_any_of signal (user or
	// context) matched — i.e. the route may lock as active_route. When false but
	// the route still appears here, only reinforcement matched (late-detected).
	Activated         bool
	UserSignalHits    []string
	ContextSignalHits []string
	ArtifactReinforce []string
}

// routeTypeDefaultMode mirrors routing-registry.yaml §route_type_spec.enum.
// Used to resolve a route's effective activation_mode when it is not declared
// explicitly. must-declare types (analysis, intelligence) have no safe default
// here; callers should rely on the explicit activation_mode for those, and we
// fall back to "auto-detect" only so a mis-annotated record still participates
// rather than silently vanishing (the must_declare lint guards authoring).
var routeTypeDefaultMode = map[string]string{
	"bootstrap":     "always-on",
	"runtime_core":  "always-on",
	"runtime_doc":   "on-demand",
	"workflow":      "auto-detect",
	"analysis":      "auto-detect", // must-declare; explicit mode expected
	"intelligence":  "auto-detect", // must-declare; explicit mode expected
	"governance":    "on-demand",
	"constitution":  "on-demand",
	"architecture":  "on-demand",
	"feedback":      "advisory",
	"metadata":      "on-demand",
	"ai_tools":      "on-demand",
	"models":        "advisory",
	"memory":        "on-demand",
	"validation":    "on-demand",
	"anti_patterns": "advisory",
}

// effectiveActivationMode resolves the mode the detector should treat a route
// as having: the explicit activation_mode wins, else the route_type default.
func effectiveActivationMode(rec runtimeRouteRecord) string {
	if m := strings.TrimSpace(rec.ActivationMode); m != "" {
		return m
	}
	if m, ok := routeTypeDefaultMode[strings.TrimSpace(rec.RouteType)]; ok {
		return m
	}
	return ""
}

// detectorParticipatingModes are the modes whose routes the detector evaluates.
// auto-detect can activate; advisory can reinforce / be suggested. always-on,
// on-demand and manual-lock never flow through the detector.
func detectorModeParticipates(mode string) bool {
	switch mode {
	case "auto-detect", "advisory":
		return true
	default:
		return false
	}
}

// normalizedTriggers is the two-phase view after folding legacy flat fields.
type normalizedTriggers struct {
	userSignals     []string
	contextSignals  []string
	artifactSignals []string
}

// normalizeRouteTriggers folds the legacy flat form into the two-phase model,
// per routing-registry.yaml §activation_triggers_spec.backward_compat:
//   - top-level user_signals      -> activation_any_of.user_signals
//   - top-level file_change_globs  -> activation_any_of.context_signals
//   - top-level artifact_signals   -> (none today; reserved)
//
// Explicit two-phase fields are merged with the legacy fields so a record using
// either or both forms is handled. Results are de-duplicated, order-stable.
func normalizeRouteTriggers(t runtimeRouteTriggers) normalizedTriggers {
	var n normalizedTriggers
	n.userSignals = append(n.userSignals, t.UserSignals...)
	n.contextSignals = append(n.contextSignals, t.FileChangeGlobs...)
	if t.ActivationAnyOf != nil {
		n.userSignals = append(n.userSignals, t.ActivationAnyOf.UserSignals...)
		n.contextSignals = append(n.contextSignals, t.ActivationAnyOf.ContextSignals...)
	}
	if t.ReinforcementAnyOf != nil {
		n.artifactSignals = append(n.artifactSignals, t.ReinforcementAnyOf.ArtifactSignals...)
	}
	n.userSignals = dedupeStable(n.userSignals)
	n.contextSignals = dedupeStable(n.contextSignals)
	n.artifactSignals = dedupeStable(n.artifactSignals)
	return n
}

// DetectWorkflows runs the deterministic detector over a parsed registry.
//
// transcript: recent turns (caller decides the window). All Role values are
// matched; pass only user turns if activation should ignore assistant text.
// openFiles: referenced files. Path drives context_signal (pre-Read) matching;
// Content drives artifact reinforcement (post-Read) matching.
//
// Returns every participating route with at least one signal hit. Order is
// stable (sorted by RouteID) so callers and tests are deterministic. Conflict
// resolution (len > 1) and miss handling (len == 0) are the caller's job
// (PreToolUse hook → workflow/workflow-routing.md Stage 2 / no-match log).
func DetectWorkflows(registry runtimeRoutingRegistry, transcript []DetectorMessage, openFiles []DetectorFile) []DetectedRoute {
	var sb strings.Builder
	for _, m := range transcript {
		sb.WriteString(m.Text)
		sb.WriteString("\n")
	}
	transcriptLower := strings.ToLower(sb.String())

	var detected []DetectedRoute
	for _, rec := range registry.Records {
		mode := effectiveActivationMode(rec)
		if !detectorModeParticipates(mode) {
			continue
		}
		n := normalizeRouteTriggers(rec.ActivationTriggers)
		if len(n.userSignals) == 0 && len(n.contextSignals) == 0 && len(n.artifactSignals) == 0 {
			continue
		}

		userHits := matchSubstrings(transcriptLower, n.userSignals)
		ctxHits := matchContextGlobs(openFiles, n.contextSignals)
		artHits := matchArtifactRegexes(openFiles, n.artifactSignals)

		activated := len(userHits) > 0 || len(ctxHits) > 0
		if !activated && len(artHits) == 0 {
			continue
		}
		detected = append(detected, DetectedRoute{
			RouteID:           rec.ID,
			EffectiveMode:     mode,
			Activated:         activated,
			UserSignalHits:    userHits,
			ContextSignalHits: ctxHits,
			ArtifactReinforce: artHits,
		})
	}
	sort.Slice(detected, func(i, j int) bool { return detected[i].RouteID < detected[j].RouteID })
	return detected
}

// matchSubstrings returns the signals that appear (case-insensitive literal
// substring) in the already-lowercased haystack.
func matchSubstrings(haystackLower string, signals []string) []string {
	var hits []string
	for _, s := range signals {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if strings.Contains(haystackLower, strings.ToLower(s)) {
			hits = append(hits, s)
		}
	}
	return hits
}

// matchContextGlobs returns the glob patterns that match at least one open file
// path. Globs support ** (any depth incl. zero segments) and * (within a path
// segment). Matching is case-insensitive.
func matchContextGlobs(files []DetectorFile, globs []string) []string {
	var hits []string
	for _, g := range globs {
		g = strings.TrimSpace(g)
		if g == "" {
			continue
		}
		re := globToRegexp(g)
		if re == nil {
			continue
		}
		for _, f := range files {
			if re.MatchString(strings.ToLower(f.Path)) {
				hits = append(hits, g)
				break
			}
		}
	}
	return hits
}

// matchArtifactRegexes returns the artifact signals (treated as
// case-insensitive regular expressions; literal fallback on compile error) that
// match the content of at least one open file. These reinforce only.
func matchArtifactRegexes(files []DetectorFile, signals []string) []string {
	var hits []string
	for _, s := range signals {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		re, err := regexp.Compile("(?is)" + s)
		matched := false
		for _, f := range files {
			if f.Content == "" {
				continue
			}
			if err != nil {
				if strings.Contains(strings.ToLower(f.Content), strings.ToLower(s)) {
					matched = true
					break
				}
				continue
			}
			if re.MatchString(f.Content) {
				matched = true
				break
			}
		}
		if matched {
			hits = append(hits, s)
		}
	}
	return hits
}

// globToRegexp converts a (lowercased) glob with ** and * into an anchored
// case-insensitive regexp. Returns nil on compile failure.
func globToRegexp(glob string) *regexp.Regexp {
	glob = strings.ToLower(glob)
	var b strings.Builder
	b.WriteString("(?i)^")
	i := 0
	for i < len(glob) {
		c := glob[i]
		switch c {
		case '*':
			if i+1 < len(glob) && glob[i+1] == '*' {
				// ** matches any sequence including path separators
				b.WriteString(".*")
				i += 2
				// swallow a trailing slash after ** so "**/x" also matches "x"
				if i < len(glob) && glob[i] == '/' {
					i++
				}
				continue
			}
			// single * matches within a segment (no slash)
			b.WriteString("[^/]*")
			i++
		case '.', '+', '(', ')', '|', '^', '$', '[', ']', '{', '}', '\\', '?':
			b.WriteByte('\\')
			b.WriteByte(c)
			i++
		default:
			b.WriteByte(c)
			i++
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil
	}
	return re
}

// dedupeStable removes duplicates while preserving first-seen order.
func dedupeStable(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
