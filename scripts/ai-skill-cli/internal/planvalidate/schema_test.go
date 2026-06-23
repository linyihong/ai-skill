package planvalidate

import (
	"reflect"
	"strings"
	"testing"
)

// Phase 2.1 / Gate B acceptance tests.

// Gate B.3 — bidirectional fixtures: the absent-version plan (today's existing
// plans) and the explicit-current-version plan must normalize to the SAME model.
// This proves the compatibility layer absorbs the first real version boundary
// (absent == baseline) rather than leaking it downstream.
func TestNormalize_AbsentAndExplicitVersionProduceSameModel(t *testing.T) {
	fields := map[string]string{
		"id":                      "sub-x",
		"plan_kind":               "sub",
		"status":                  "draft",
		"parent":                  "main-x",
		"required_for_completion": "true",
		"sub_plan_reason":         "r",
	}
	legacy := RawPlan{Path: "plans/active/x.md", Location: "active", SchemaVersion: "", Fields: fields}
	current := RawPlan{Path: "plans/active/x.md", Location: "active", SchemaVersion: "1", Fields: fields}

	lm, err := Normalize(legacy)
	if err != nil {
		t.Fatalf("legacy normalize: unexpected error: %v", err)
	}
	cm, err := Normalize(current)
	if err != nil {
		t.Fatalf("current normalize: unexpected error: %v", err)
	}
	if !reflect.DeepEqual(lm, cm) {
		t.Fatalf("absent vs explicit version produced different models:\n legacy=%+v\ncurrent=%+v", lm, cm)
	}
}

// Gate B.1 — the engine-facing model must not carry a schema version. Mechanize
// the "grep version -> 0 hit in the engine's input" acceptance via reflection so
// it cannot regress when fields are added later.
func TestNormalizedPlanModel_HasNoVersionField(t *testing.T) {
	tp := reflect.TypeOf(NormalizedPlanModel{})
	for i := 0; i < tp.NumField(); i++ {
		name := strings.ToLower(tp.Field(i).Name)
		if strings.Contains(name, "version") || strings.Contains(name, "schema") {
			t.Fatalf("NormalizedPlanModel must not expose a version/schema field, found %q", tp.Field(i).Name)
		}
	}
}

// Gate B.2 support — an unsupported version is rejected in the compat layer, so
// no downstream validator ever has to branch on a version.
func TestNormalize_UnsupportedVersionRejected(t *testing.T) {
	_, err := Normalize(RawPlan{SchemaVersion: "99", Fields: map[string]string{"id": "x"}})
	if err == nil {
		t.Fatalf("expected error for unsupported schema_version, got nil")
	}
}

// required_for_completion is three-state: absent (nil), true, false. Confirm the
// compat layer preserves the distinction for the engine.
func TestNormalize_RequiredForCompletionTriState(t *testing.T) {
	absent, _ := Normalize(RawPlan{Fields: map[string]string{"id": "x"}})
	if absent.RequiredForCompletion != nil {
		t.Fatalf("absent required_for_completion should be nil")
	}
	tru, _ := Normalize(RawPlan{Fields: map[string]string{"required_for_completion": "true"}})
	if tru.RequiredForCompletion == nil || !*tru.RequiredForCompletion {
		t.Fatalf("required_for_completion=true should normalize to *true")
	}
	fal, _ := Normalize(RawPlan{Fields: map[string]string{"required_for_completion": "false"}})
	if fal.RequiredForCompletion == nil || *fal.RequiredForCompletion {
		t.Fatalf("required_for_completion=false should normalize to *false")
	}
}
