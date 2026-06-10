package app

import (
	"database/sql"
	"path/filepath"
	"sort"
	"testing"
)

// project_metadata_compile_test.go — Phase 1C₂ projection rule tests.
//
// Spec reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1C₂.
//
// The contracts Phase 1C₂ guards:
//   - case-variant expansion happens in projection, not in the parser
//   - explicit case_variants suppress auto-derivation and are taken verbatim
//   - cross-entity token collisions keep entity identity (governance debug)
//   - public projects contribute no rows (only private declare protected surface)
//   - the additive projection populates derived_private_entities +
//     derived_match_tokens without disturbing the legacy table

func tokenSet(rows []derivedMatchTokenRow) map[string]string {
	out := map[string]string{}
	for _, r := range rows {
		out[r.Token] = r.CanonicalToken
	}
	return out
}

func TestProjectMetadataDerivedRows_AutoCaseVariants(t *testing.T) {
	meta := &ProjectMetadataFile{Project: ProjectMetadata{
		ID:         "demo-project",
		Visibility: "private",
		PrivateEntities: []ProjectEntity{
			{
				Name:         "Project Foo",
				Kind:         "codename",
				MatchTokens:  []string{"ProjectFoo"},
				CaseVariants: CaseVariants{Mode: "auto"},
			},
		},
	}}

	entities, tokens := projectMetadataDerivedRows(meta, ".ai-skill-project.yaml")

	if len(entities) != 1 {
		t.Fatalf("expected 1 governance row, got %d", len(entities))
	}
	if entities[0].EntityName != "Project Foo" || entities[0].Kind != "codename" {
		t.Errorf("governance row identity wrong: %+v", entities[0])
	}
	if entities[0].SuggestedPlaceholder != "<DEMO_PROJECT>" {
		t.Errorf("placeholder = %q; want <DEMO_PROJECT>", entities[0].SuggestedPlaceholder)
	}

	set := tokenSet(tokens)
	// The literal must always be present and self-canonical.
	if canon, ok := set["ProjectFoo"]; !ok || canon != "ProjectFoo" {
		t.Errorf("literal ProjectFoo missing/wrong canonical: %v / %q", ok, canon)
	}
	// Auto expansion must produce kebab + screaming-snake variants, each
	// tracing back to the declared match_token.
	for _, want := range []string{"project-foo", "PROJECT_FOO"} {
		if canon, ok := set[want]; !ok {
			t.Errorf("auto variant %q missing", want)
		} else if canon != "ProjectFoo" {
			t.Errorf("auto variant %q canonical = %q; want ProjectFoo", want, canon)
		}
	}
	// Every token row must carry the entity identity.
	for _, r := range tokens {
		if r.EntityName != "Project Foo" {
			t.Errorf("token %q lost entity identity: %q", r.Token, r.EntityName)
		}
	}
}

func TestProjectMetadataDerivedRows_ExplicitSuppressesAuto(t *testing.T) {
	meta := &ProjectMetadataFile{Project: ProjectMetadata{
		ID:         "demo",
		Visibility: "private",
		PrivateEntities: []ProjectEntity{
			{
				Name:         "Acme",
				Kind:         "client",
				MatchTokens:  []string{"Acme"},
				CaseVariants: CaseVariants{Mode: "explicit", Variants: []string{"AcmeCorp", "acme-inc"}},
			},
		},
	}}

	_, tokens := projectMetadataDerivedRows(meta, ".ai-skill-project.yaml")
	set := tokenSet(tokens)

	// Literal + explicit variants present.
	for _, want := range []string{"Acme", "AcmeCorp", "acme-inc"} {
		if _, ok := set[want]; !ok {
			t.Errorf("explicit-mode token %q missing", want)
		}
	}
	// Auto-derivation must be suppressed: the screaming-snake form of "Acme"
	// (ACME) must NOT appear unless explicitly listed.
	if _, ok := set["ACME"]; ok {
		t.Errorf("explicit mode must suppress auto variant ACME, but it was emitted")
	}
}

func TestProjectMetadataDerivedRows_CrossEntityCollisionKeepsIdentity(t *testing.T) {
	meta := &ProjectMetadataFile{Project: ProjectMetadata{
		ID:         "demo",
		Visibility: "private",
		PrivateEntities: []ProjectEntity{
			{Name: "Alpha", Kind: "codename", MatchTokens: []string{"shared-token"}, CaseVariants: CaseVariants{Mode: "explicit", Variants: []string{"shared-token"}}},
			{Name: "Beta", Kind: "product", MatchTokens: []string{"shared-token"}, CaseVariants: CaseVariants{Mode: "explicit", Variants: []string{"shared-token"}}},
		},
	}}

	_, tokens := projectMetadataDerivedRows(meta, ".ai-skill-project.yaml")

	var owners []string
	for _, r := range tokens {
		if r.Token == "shared-token" {
			owners = append(owners, r.EntityName)
		}
	}
	sort.Strings(owners)
	if len(owners) != 2 || owners[0] != "Alpha" || owners[1] != "Beta" {
		t.Errorf("colliding token must yield one row per entity with distinct identity; got owners %v", owners)
	}
}

func TestProjectMetadataDerivedRows_PublicProjectContributesNothing(t *testing.T) {
	meta := &ProjectMetadataFile{Project: ProjectMetadata{
		ID:         "open-source",
		Visibility: "public",
		PrivateEntities: []ProjectEntity{
			{Name: "X", Kind: "codename", MatchTokens: []string{"X-secret"}, CaseVariants: CaseVariants{Mode: "auto"}},
		},
	}}
	entities, tokens := projectMetadataDerivedRows(meta, ".ai-skill-project.yaml")
	if len(entities) != 0 || len(tokens) != 0 {
		t.Errorf("public project must contribute no rows; got %d entities, %d tokens", len(entities), len(tokens))
	}
}

func TestCompileProjectMetadataDerived_PopulatesBothTables(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "acme-app", ".ai-skill-project.yaml"), `
project:
  id: acme-app
  visibility: private
  private_entities:
    - name: Acme Customer
      kind: client
      match_tokens:
        - "AcmeCustomer"
      case_variants: auto
`)
	// A public project file must be ignored entirely.
	writeFile(t, filepath.Join(repo, "oss-lib", ".ai-skill-project.yaml"), `
project:
  id: oss-lib
  visibility: public
  private_entities: []
`)

	dbPath := filepath.Join(repo, "runtime.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := createGoRuntimeSchema(db); err != nil {
		t.Fatal(err)
	}

	if err := compileProjectMetadataDerived(repo, db); err != nil {
		t.Fatalf("compileProjectMetadataDerived: %v", err)
	}

	var entityCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM derived_private_entities`).Scan(&entityCount); err != nil {
		t.Fatal(err)
	}
	if entityCount != 1 {
		t.Fatalf("expected 1 governance row (public skipped), got %d", entityCount)
	}

	var name, kind, project, placeholder string
	if err := db.QueryRow(`SELECT entity_name, kind, owning_project_id, suggested_placeholder FROM derived_private_entities`).Scan(&name, &kind, &project, &placeholder); err != nil {
		t.Fatal(err)
	}
	if name != "Acme Customer" || kind != "client" || project != "acme-app" || placeholder != "<ACME_APP>" {
		t.Errorf("governance row wrong: name=%q kind=%q project=%q placeholder=%q", name, kind, project, placeholder)
	}

	// Execution layer must hold the literal + auto variants.
	for _, want := range []string{"AcmeCustomer", "acme-customer", "ACME_CUSTOMER"} {
		var n int
		if err := db.QueryRow(`SELECT COUNT(*) FROM derived_match_tokens WHERE token = ? AND entity_name = 'Acme Customer'`, want).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n != 1 {
			t.Errorf("expected match token %q present once, got %d", want, n)
		}
	}

	// Bootstrap safety: a brand-new framework concept that no project
	// declares must not appear in derived_match_tokens.
	var leaked int
	if err := db.QueryRow(`SELECT COUNT(*) FROM derived_match_tokens WHERE token = 'ActivationBridgeV2'`).Scan(&leaked); err != nil {
		t.Fatal(err)
	}
	if leaked != 0 {
		t.Errorf("undeclared framework concept must not be forbidden; got %d rows", leaked)
	}
}
