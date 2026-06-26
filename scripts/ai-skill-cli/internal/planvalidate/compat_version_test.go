package planvalidate

import (
	"reflect"
	"testing"
)

// Phase 3.2 Compatibility Slice acceptance — subject=artifact (schema_version),
// axis=schema, exactly one upgrade. Three cases per the 3.2 contract.

func findingSet(fs []Finding) map[cmpKey]bool {
	m := map[cmpKey]bool{}
	for _, f := range fs {
		m[cmpKey{f.RuleID, f.Blocking}] = true
	}
	return m
}

// modelsWithVersion builds the same plan set (a violation tree) declaring a given
// schema_version, then normalizes. Because Normalize absorbs the version, the
// resulting models — and thus findings — must be invariant to the version.
func modelsWithVersion(t *testing.T, ver string) []NormalizedPlanModel {
	t.Helper()
	raws := []RawPlan{
		{Path: "plans/active/_plan.md", SchemaVersion: ver, Fields: map[string]string{"id": "m", "plan_kind": "main", "parent": "null"}},
		{Path: "plans/active/01.md", SchemaVersion: ver, Fields: map[string]string{"id": "s", "plan_kind": "sub", "parent": "ghost", "required_for_completion": "true", "sub_plan_reason": "x"}},
	}
	var ms []NormalizedPlanModel
	for _, r := range raws {
		m, err := Normalize(r)
		if err != nil {
			t.Fatalf("normalize ver %q: unexpected error %v", ver, err)
		}
		ms = append(ms, m)
	}
	return ms
}

// supported upgrade (schema_version 1 -> 2) must preserve findings (Structural =
// same RuleID set, Behavioral = same Blocking). Single axis x single subject.
func TestCompat_SupportedVersionUpgradePreservesFindings(t *testing.T) {
	v1 := Validate(ValidationContext{}, modelsWithVersion(t, "1"))
	v2 := Validate(ValidationContext{}, modelsWithVersion(t, "2"))
	if len(v1) == 0 {
		t.Fatalf("expected the violation tree to produce findings to compare")
	}
	if !reflect.DeepEqual(findingSet(v1), findingSet(v2)) {
		t.Fatalf("schema_version 1->2 changed findings: v1=%v v2=%v", v1, v2)
	}
}

// unsupported combination must reject deterministically AND diagnosably: same
// rejection stage + same reason class every run (classify on type, not string).
func TestCompat_UnsupportedRejectDeterministicDiagnosable(t *testing.T) {
	for _, bad := range []string{"99", "abc", "99"} { // "99" twice = determinism
		_, err := Normalize(RawPlan{SchemaVersion: bad, Fields: map[string]string{"id": "x"}})
		ce, ok := err.(*CompatError)
		if !ok {
			t.Fatalf("ver %q: expected *CompatError, got %T (%v)", bad, err, err)
		}
		if ce.Stage != "normalize" || ce.ReasonClass != "unsupported_schema_version" {
			t.Fatalf("ver %q: diagnostic drift stage=%q reason=%q", bad, ce.Stage, ce.ReasonClass)
		}
	}
}

// no-version-change baseline: running the (now version-set-aware) machinery
// without changing any version must leave findings unchanged — negative proof
// that the upgrade machinery itself does not alter results.
func TestCompat_NoChangeBaselineStable(t *testing.T) {
	a := Validate(ValidationContext{}, modelsWithVersion(t, "1"))
	b := Validate(ValidationContext{}, modelsWithVersion(t, "1"))
	if !reflect.DeepEqual(findingSet(a), findingSet(b)) {
		t.Fatalf("no-change baseline drifted across runs: a=%v b=%v", a, b)
	}
}
