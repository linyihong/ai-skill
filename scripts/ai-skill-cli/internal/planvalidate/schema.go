package planvalidate

import "fmt"

// Phase 2.1 (Gate B): the schema compatibility layer.
//
// This layer is the SINGLE place a plan frontmatter schema version is observed
// and resolved. It maps a version-bearing RawPlan to a version-free
// NormalizedPlanModel. Gate B constraints:
//
//   - B.1: the output is NormalizedPlanModel, never a schema/version object —
//     NormalizedPlanModel has no version field (enforced by a reflection test).
//   - B.2: all compatibility decisions live in Normalize; downstream validators
//     must never branch on a version (e.g. `if plan.Version >= 2`).
//   - B.3: absent and explicit-current versions must normalize to the SAME
//     model (bidirectional fixtures in schema_test.go).
//
// See plans/active/2026-06-22-1009-plans-system-portability-and-delivery-integration/
// 01-external-repo-plan-system-shared-binary.md (Phase 2.1 / Gate B).

// currentSchemaVersion is the plan frontmatter schema version this build
// normalizes to. An absent schema_version is treated as this baseline, so
// existing plans (which carry no schema_version field) normalize identically to
// plans that declare it explicitly. This is the first real compatibility
// boundary and the anchor for plan_schema versioning (Open Question Q3).
// currentSchemaVersion is the version an absent schema_version resolves to.
const currentSchemaVersion = "1"

// supportedSchemaVersions is the accepted set. Phase 3.2 (Compatibility Slice,
// subject=artifact): the set is intentionally EXTENSIBLE — adding a version here
// must not change findings for existing plans. Versions "1" and "2" currently
// normalize identically (no shape divergence yet); "2" is reserved to prove the
// acceptance set extends without semantic drift, NOT to encode a real migration.
// When a real shape change lands, Normalize gains per-version mapping; the
// engine-facing NormalizedPlanModel and its consumers stay unchanged (Gate B).
var supportedSchemaVersions = map[string]bool{"1": true, "2": true}

// CompatError is a typed, classified rejection from the compatibility layer so a
// reject is DIAGNOSABLE (Phase 3.2 patch): the same unsupported case must reject
// at the same Stage with the same ReasonClass on every run (no diagnostic drift).
// Consumers classify on Stage/ReasonClass, never on the message string.
type CompatError struct {
	Stage       string // where rejection happened, e.g. "normalize"
	ReasonClass string // stable class, e.g. "unsupported_schema_version"
	Detail      string // human detail (not for classification)
}

func (e *CompatError) Error() string {
	return e.Stage + "/" + e.ReasonClass + ": " + e.Detail
}

// RawPlan is the un-normalized frontmatter a schema loader extracts from a plan
// file, before any version resolution. It is the ONLY type in this package that
// carries a schema version; everything downstream consumes NormalizedPlanModel.
type RawPlan struct {
	Path     string
	Location string // active | archived
	// SchemaVersion is the declared frontmatter schema version. An empty string
	// means the field was absent and is resolved to currentSchemaVersion.
	SchemaVersion string
	// Fields holds the raw scalar frontmatter values (id, plan_kind, status, …)
	// exactly as parsed, before normalization.
	Fields map[string]string
}

// NormalizedPlanModel is the stable, version-free plan representation the engine
// will consume. It deliberately has NO schema version field: Gate B requires the
// engine never to see a version. Introducing a new schema version means
// extending Normalize, never changing this struct or its consumers.
type NormalizedPlanModel struct {
	Path                  string
	Location              string
	ID                    string
	PlanKind              string
	Status                string
	Parent                string
	RequiredForCompletion *bool
	SubPlanReason         string
}

// Normalize resolves a RawPlan's schema version and produces a version-free
// NormalizedPlanModel. It is the only function that reads SchemaVersion.
//
// Note (deferred, see plan): a future RawPlan carrying a deprecated-but-tolerated
// field will need to surface a migration warning. That requires distinguishing
// "missing field" from "deprecated field" and returning warnings alongside the
// model (a CompatibilityResult{ Model, Warnings }) so each consumer (hook / CI /
// CLI) decides what to do. Not built now — Normalize returns (model, error) and
// the warning channel is reserved for when a second schema version lands.
// normalizeNullScalar maps YAML null idioms to the empty string. A main plan
// declares `parent: null`; without this the engine would treat the literal
// string "null" as an unresolved parent id and emit a false positive. This is a
// compat-layer concern (Gate B): the engine consumes a clean model and must
// never see YAML idioms. Surfaced by the Vidoe-Test external plan tree
// (2026-06-24); same family as the quoted-scalar requirement (Q3).
func normalizeNullScalar(v string) string {
	switch v {
	case "null", "Null", "NULL", "~":
		return ""
	}
	return v
}

func Normalize(raw RawPlan) (NormalizedPlanModel, error) {
	version := raw.SchemaVersion
	if version == "" {
		version = currentSchemaVersion // absent == baseline
	}
	if !supportedSchemaVersions[version] {
		// Deterministic, diagnosable reject: same stage + reason class every run.
		return NormalizedPlanModel{}, &CompatError{
			Stage:       "normalize",
			ReasonClass: "unsupported_schema_version",
			Detail:      fmt.Sprintf("schema_version %q not in supported set", version),
		}
	}

	m := NormalizedPlanModel{
		Path:          raw.Path,
		Location:      raw.Location,
		ID:            raw.Fields["id"],
		PlanKind:      raw.Fields["plan_kind"],
		Status:        raw.Fields["status"],
		Parent:        normalizeNullScalar(raw.Fields["parent"]),
		SubPlanReason: raw.Fields["sub_plan_reason"],
	}
	if rfc, ok := raw.Fields["required_for_completion"]; ok {
		b := rfc == "true"
		m.RequiredForCompletion = &b
	}
	return m, nil
}
