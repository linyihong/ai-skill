package app

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// project_metadata_test.go — Phase 1A canonical parser tests.
//
// Spec reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1A — "schema validation test: parser 接受合法
// schema、拒絕缺 kind / 缺 match_tokens / case_variants 既非 auto 也非
// list 等變形".
//
// These tests exercise the v1 canonical parser ONLY. They MUST NOT
// reference the legacy reader (sanitization_scan.go::readProjectMetadata)
// or its types. Phase 1A scope discipline: parser surface is independent.

func TestLoadProjectMetadata_ValidMinimal(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities: []
`)
	file, err := ParseProjectMetadata(body)
	if err != nil {
		t.Fatalf("expected valid schema to parse; got error: %v", err)
	}
	if file.Project.ID != "example-project" {
		t.Errorf("expected project.id=example-project, got %q", file.Project.ID)
	}
	if file.Project.Visibility != "private" {
		t.Errorf("expected visibility=private, got %q", file.Project.Visibility)
	}
	if len(file.Project.PrivateEntities) != 0 {
		t.Errorf("expected 0 entities, got %d", len(file.Project.PrivateEntities))
	}
}

func TestLoadProjectMetadata_ValidWithEntities(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: ProjectAtlas
      kind: codename
      match_tokens:
        - Atlas
        - ProjectAtlas
      case_variants: auto
    - name: AcmeCorporation
      kind: client
      match_tokens:
        - Acme
        - AcmeCorp
      case_variants:
        - Acme
        - ACME
        - AcmeCorp
        - acme-corp
`)
	file, err := ParseProjectMetadata(body)
	if err != nil {
		t.Fatalf("expected valid schema to parse; got error: %v", err)
	}
	if len(file.Project.PrivateEntities) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(file.Project.PrivateEntities))
	}
	atlas := file.Project.PrivateEntities[0]
	if atlas.Name != "ProjectAtlas" || atlas.Kind != "codename" {
		t.Errorf("unexpected first entity: %+v", atlas)
	}
	if atlas.CaseVariants.Mode != "auto" {
		t.Errorf("expected atlas case_variants mode=auto, got %q", atlas.CaseVariants.Mode)
	}
	acme := file.Project.PrivateEntities[1]
	if acme.CaseVariants.Mode != "explicit" {
		t.Errorf("expected acme case_variants mode=explicit, got %q", acme.CaseVariants.Mode)
	}
	if len(acme.CaseVariants.Variants) != 4 {
		t.Errorf("expected 4 explicit acme variants, got %d", len(acme.CaseVariants.Variants))
	}
}

func TestLoadProjectMetadata_RejectsMissingID(t *testing.T) {
	body := []byte(`
project:
  visibility: private
  private_entities: []
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for missing project.id")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected ProjectMetadataValidationError, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "project.id.required") {
		t.Errorf("expected error to cite project.id.required rule; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsBadIDFormat(t *testing.T) {
	body := []byte(`
project:
  id: Bad_ID_With_Underscores
  visibility: private
  private_entities: []
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for non-kebab id")
	}
	if !strings.Contains(err.Error(), "project.id.kebab") {
		t.Errorf("expected project.id.kebab rule citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsMissingVisibility(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  private_entities: []
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for missing visibility")
	}
	if !strings.Contains(err.Error(), "project.visibility.enum") {
		t.Errorf("expected project.visibility.enum rule citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsBadVisibility(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: secret
  private_entities: []
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for invalid visibility enum value")
	}
	if !strings.Contains(err.Error(), "project.visibility.enum") {
		t.Errorf("expected visibility enum citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsEntityMissingKind(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: ProjectAtlas
      match_tokens:
        - Atlas
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for entity missing kind")
	}
	if !strings.Contains(err.Error(), "entity.kind.required") {
		t.Errorf("expected entity.kind.required citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsEntityBadKind(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: ProjectAtlas
      kind: super-secret
      match_tokens:
        - Atlas
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for invalid kind enum")
	}
	if !strings.Contains(err.Error(), "entity.kind.required") {
		t.Errorf("expected entity.kind.required citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsEntityMissingMatchTokens(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: ProjectAtlas
      kind: codename
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for entity missing match_tokens")
	}
	if !strings.Contains(err.Error(), "entity.match_tokens.required") {
		t.Errorf("expected entity.match_tokens.required citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsEntityEmptyMatchTokens(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: ProjectAtlas
      kind: codename
      match_tokens:
        - "   "
        - ""
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for entity with only whitespace match_tokens")
	}
	if !strings.Contains(err.Error(), "entity.match_tokens.required") {
		t.Errorf("expected entity.match_tokens.required citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsEntityMissingName(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - kind: codename
      match_tokens:
        - Atlas
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure for entity missing name")
	}
	if !strings.Contains(err.Error(), "entity.name.required") {
		t.Errorf("expected entity.name.required citation; got:\n%s", err.Error())
	}
}

func TestLoadProjectMetadata_RejectsBadCaseVariantsScalar(t *testing.T) {
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: ProjectAtlas
      kind: codename
      match_tokens:
        - Atlas
      case_variants: notauto
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected parse failure for case_variants scalar != 'auto'")
	}
	// case_variants shape errors surface at unmarshal time as plain
	// errors (not validation errors), because they make the file
	// fundamentally malformed. IsValidationError should report false.
	if IsValidationError(err) {
		t.Errorf("expected raw parse error for case_variants malformed scalar, got ValidationError: %v", err)
	}
	if !strings.Contains(err.Error(), "case_variants") {
		t.Errorf("expected error to mention case_variants; got: %v", err)
	}
}

func TestLoadProjectMetadata_AggregatesMultipleFailures(t *testing.T) {
	body := []byte(`
project:
  id: ""
  visibility: secret
  private_entities:
    - kind: badkind
      match_tokens: []
`)
	_, err := ParseProjectMetadata(body)
	if err == nil {
		t.Fatal("expected validation failure")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected ProjectMetadataValidationError, got %T: %v", err, err)
	}
	// Should report all four failures: project.id.required, project.visibility.enum,
	// entity.name.required, entity.kind.required, entity.match_tokens.required.
	var v *ProjectMetadataValidationError
	if !errors.As(err, &v) {
		t.Fatalf("errors.As failed: %v", err)
	}
	if len(v.Failures) < 4 {
		t.Errorf("expected at least 4 aggregated failures, got %d:\n%s", len(v.Failures), err.Error())
	}
}

func TestLoadProjectMetadata_IOErrorDistinctFromValidation(t *testing.T) {
	_, err := LoadProjectMetadata(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("expected I/O error for missing file")
	}
	if IsValidationError(err) {
		t.Errorf("expected I/O error not to surface as ValidationError, got: %v", err)
	}
}

func TestLoadProjectMetadata_PublicVisibilityAllowed(t *testing.T) {
	body := []byte(`
project:
  id: public-project
  visibility: public
  private_entities: []
`)
	_, err := ParseProjectMetadata(body)
	if err != nil {
		t.Errorf("expected public visibility to be valid; got: %v", err)
	}
}

func TestLoadProjectMetadata_AllKindsAccepted(t *testing.T) {
	for _, kind := range []string{"codename", "client", "product", "individual", "other"} {
		body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: TestEntity
      kind: ` + kind + `
      match_tokens:
        - TestEntity
`)
		if _, err := ParseProjectMetadata(body); err != nil {
			t.Errorf("expected kind=%s to be accepted; got: %v", kind, err)
		}
	}
}

// TestLoadProjectMetadata_SingleCharIDAllowed locks in the Phase 1A
// review decision (2026-06-08 follow-up): single-character project ids
// (e.g. "a", "0") are accepted by both the canonical schema and the
// parser. The schema pattern is the canonical SOT for this contract;
// the parser regex MUST match it. If you tighten the regex to require
// ≥2 chars, also update the schema pattern in
// metadata/project/ai-skill-project-schema.yaml — they are paired.
func TestLoadProjectMetadata_SingleCharIDAllowed(t *testing.T) {
	body := []byte(`
project:
  id: a
  visibility: private
  private_entities: []
`)
	file, err := ParseProjectMetadata(body)
	if err != nil {
		t.Fatalf("expected single-char id to be accepted; got: %v", err)
	}
	if file.Project.ID != "a" {
		t.Errorf("expected project.id=a, got %q", file.Project.ID)
	}
	// Digit form too.
	body2 := []byte(`
project:
  id: "0"
  visibility: private
  private_entities: []
`)
	if _, err := ParseProjectMetadata(body2); err != nil {
		t.Errorf("expected single-digit id to be accepted; got: %v", err)
	}
}

// TestLoadProjectMetadata_UnknownFieldTolerance locks in the Phase 1A
// review decision (2026-06-08 follow-up): unknown fields are silently
// tolerated for forward-compat. The schema declares this contract; this
// test ensures the parser actually honours it. If a future phase wants
// diagnostics for typo'd / unrecognized fields, design it as a separate
// metadata-diagnostics surface — do NOT change the parser behaviour
// here without also revising the schema doc and removing this test.
func TestLoadProjectMetadata_UnknownFieldTolerance(t *testing.T) {
	// Top-level unknown sibling of `project:`.
	body := []byte(`
project:
  id: example-project
  visibility: private
  private_entities: []
unknown_top_level: anything
also_unknown:
  nested: ok
`)
	if _, err := ParseProjectMetadata(body); err != nil {
		t.Errorf("expected top-level unknown fields to be tolerated; got: %v", err)
	}

	// Unknown field within `project:`.
	body = []byte(`
project:
  id: example-project
  visibility: private
  private_entities: []
  future_field: someday
  another_future_one:
    structured: true
`)
	if _, err := ParseProjectMetadata(body); err != nil {
		t.Errorf("expected project-level unknown fields to be tolerated; got: %v", err)
	}

	// Unknown field within an entity.
	body = []byte(`
project:
  id: example-project
  visibility: private
  private_entities:
    - name: TestEntity
      kind: codename
      match_tokens:
        - TestToken
      severity: high          # future field (Phase 1C+)
      expires_at: 2099-01-01  # future field (Phase 1C+)
`)
	if _, err := ParseProjectMetadata(body); err != nil {
		t.Errorf("expected entity-level unknown fields to be tolerated; got: %v", err)
	}
}

// TestLoadProjectMetadata_ExampleFileParsesCleanly ensures the example
// fixture documented in metadata/project/example-ai-skill-project.yaml
// stays valid. Doubles as a regression for the documented schema.
func TestLoadProjectMetadata_ExampleFileParsesCleanly(t *testing.T) {
	repo := discoveryRepoRoot(t)
	path := filepath.Join(repo, "metadata", "project", "example-ai-skill-project.yaml")
	file, err := LoadProjectMetadata(path)
	if err != nil {
		t.Fatalf("example file must parse cleanly; got: %v", err)
	}
	if file.Project.ID == "" {
		t.Error("example file should have a project.id")
	}
	if len(file.Project.PrivateEntities) == 0 {
		t.Error("example file should declare at least one entity for documentation value")
	}
	// Spot-check the first entity has the structured shape (NOT a flat string).
	first := file.Project.PrivateEntities[0]
	if first.Name == "" || first.Kind == "" || len(first.MatchTokens) == 0 {
		t.Errorf("example file first entity should have name + kind + match_tokens; got: %+v", first)
	}
}
