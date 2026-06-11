package app

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// repository_topology.go — CANONICAL loader for runtime/repository-topology.yaml.
//
// Schema spec + migration trajectory: runtime/repository-topology-migration.md.
// Plan reference: plans/active/2026-06-06-1800-sanitization-mechanical-
// enforcement.md §Phase 1B.
//
// SCOPE DISCIPLINE — Phase 1B (this file) defines the canonical loader
// and the v2 contract. It is NOT wired into the live projection rule
// (compiler line 339 still hard-codes v1 field names against the live
// YAML which is still in v1 shape). Phase 1C migrates the live YAML to
// v2 AND rewrites the projection rule to call this loader. Phase 1D
// adapts sanitization_scan.go's consumer path. See migration notes.
//
// Reviewer constraint: this file MUST NOT import from or reference
// sanitization_scan.go symbols (the legacy consumer of the projected
// table). It also MUST NOT modify the existing projection rule in
// runtime_compiler.go. Even a one-line cross-file dependency would
// blur the 1A / 1B / 1C / 1D phase boundaries established 2026-06-08.

// RepositoryTopologyFile is the parsed in-memory representation of
// runtime/repository-topology.yaml. Both v1 and v2 YAML shapes normalize
// to this type; the loader records which schema version was on disk via
// SchemaVersion so writers can preserve fidelity.
//
// v1 fields (owner / purpose / consumer_tracking) are empty after a v1
// read — they have no equivalent on disk. The writer rejects empty
// owner/purpose for any v2 write, ensuring a round-trip through
// WriteRepositoryTopology always produces a complete v2 file.
// This type is the in-memory normalized form. It is never directly
// marshalled or unmarshalled by yaml.v3 — read goes through
// rawTopologyFile, write goes through marshalV2's anonymous outFile.
// Therefore no yaml tags appear on its fields.
type RepositoryTopologyFile struct {
	SchemaVersion     int                      // what was on disk: 1 or 2
	Status            string                   // top-level `status:` (e.g. "active")
	OwnerLayer        string                   // top-level `owner_layer:` (e.g. "runtime")
	RuntimeProjection *RuntimeProjectionConfig // optional `runtime_projection:` block
	ConsumerTracking  *ConsumerTracking        // v2 only; nil after v1 read
	Subtrees          []Subtree                // normalized to v2 shape regardless of input
	Invariants        []string                 // optional, both schemas
}

// Subtree is a single subtree classification entry. In v1 YAML the keys
// are `subtree:` + `shared:`; in v2 they are `path:` + `shared_layer:` +
// `owner:` + `purpose:`. Owner and Purpose are required for v2 writes
// (the writer rejects missing values); after a v1 read they are empty
// strings.
type Subtree struct {
	Path        string
	SharedLayer bool
	Owner       string // v2 only; empty after v1 read
	Purpose     string // v2 only; empty after v1 read
}

// ConsumerTracking is the v2 frozen governance block that replaced v1's
// manual `expected_consumers:` list. The Strategy must be the literal
// "code_reference" — adding new strategies (or re-adding expected_consumers)
// requires a v3 schema bump and an explicit governance decision.
type ConsumerTracking struct {
	Strategy  string `yaml:"strategy"`
	Rationale string `yaml:"rationale,omitempty"`
}

// RuntimeProjectionConfig mirrors the top-level `runtime_projection:` block
// shared by most runtime/*.yaml surfaces. Carried through reads/writes
// verbatim so the loader does not strip it during round-trips.
type RuntimeProjectionConfig struct {
	Enabled   bool   `yaml:"enabled"`
	TargetKey string `yaml:"target_key,omitempty"`
}

// rawTopologyFile is the union YAML shape accepted by the loader. v1
// fields and v2 fields are both declared so a single unmarshal handles
// either input; the loader then normalizes into RepositoryTopologyFile.
type rawTopologyFile struct {
	SchemaVersion     int                       `yaml:"schema_version"`
	Status            string                    `yaml:"status"`
	OwnerLayer        string                    `yaml:"owner_layer"`
	RuntimeProjection *RuntimeProjectionConfig  `yaml:"runtime_projection,omitempty"`

	// v1 fields
	SharedLayerClassification []rawSubtreeV1 `yaml:"shared_layer_classification,omitempty"`
	ExpectedConsumers         []string       `yaml:"expected_consumers,omitempty"`

	// v2 fields
	ConsumerTracking *ConsumerTracking `yaml:"consumer_tracking,omitempty"`
	Subtrees         []rawSubtreeV2    `yaml:"subtrees,omitempty"`

	// Shared
	Invariants []string `yaml:"invariants,omitempty"`

	// Forward-compat catch-all for fields neither v1 nor v2 declared
	// (e.g. v3 might add fields). Captured but not consumed.
	Unknown map[string]any `yaml:",inline"`
}

type rawSubtreeV1 struct {
	Subtree string `yaml:"subtree"`
	Shared  bool   `yaml:"shared"`

	Unknown map[string]any `yaml:",inline"`
}

type rawSubtreeV2 struct {
	Path        string `yaml:"path"`
	SharedLayer bool   `yaml:"shared_layer"`
	Owner       string `yaml:"owner"`
	Purpose     string `yaml:"purpose"`

	Unknown map[string]any `yaml:",inline"`
}

// LoadRepositoryTopology reads runtime/repository-topology.yaml (or any
// file conforming to the v1 or v2 schema) and returns the normalized
// in-memory representation. Both schema versions are accepted on read.
//
// I/O errors (file not found, YAML parse failure) surface as plain errors.
// Schema-violation errors (e.g. v1 file with no `shared_layer_classification`
// key; v2 file with subtrees missing both v1 and v2 paths) are returned
// as *RepositoryTopologyValidationError so callers can branch.
func LoadRepositoryTopology(path string) (*RepositoryTopologyFile, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read repository topology file %s: %w", path, err)
	}
	return ParseRepositoryTopology(body)
}

// ParseRepositoryTopology is LoadRepositoryTopology's byte-slice form.
func ParseRepositoryTopology(body []byte) (*RepositoryTopologyFile, error) {
	var raw rawTopologyFile
	if err := yaml.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse repository topology yaml: %w", err)
	}
	return normalizeTopology(&raw)
}

func normalizeTopology(raw *rawTopologyFile) (*RepositoryTopologyFile, error) {
	// Detect schema version. The on-disk schema_version field is
	// authoritative; if absent, fall back to shape inference (v2 has
	// `subtrees:`; v1 has `shared_layer_classification:`).
	version := raw.SchemaVersion
	if version == 0 {
		switch {
		case len(raw.Subtrees) > 0 || raw.ConsumerTracking != nil:
			version = 2
		case len(raw.SharedLayerClassification) > 0 || len(raw.ExpectedConsumers) > 0:
			version = 1
		default:
			// Empty file — pick v2 so writers default to the canonical
			// shape; loader callers can override via SchemaVersion if
			// they need v1 fidelity.
			version = 2
		}
	}

	file := &RepositoryTopologyFile{
		SchemaVersion:     version,
		Status:            raw.Status,
		OwnerLayer:        raw.OwnerLayer,
		RuntimeProjection: raw.RuntimeProjection,
		Invariants:        raw.Invariants,
	}

	var failures []RepositoryTopologyValidationFailure

	switch version {
	case 1:
		// Both shapes can theoretically co-exist in one file; the v1
		// path takes priority on read so legacy YAMLs work as-is. If a
		// caller mixes shapes, the v2 fields are silently ignored on
		// the v1 read path (a writer would re-emit only v2).
		for i, s := range raw.SharedLayerClassification {
			if strings.TrimSpace(s.Subtree) == "" {
				failures = append(failures, RepositoryTopologyValidationFailure{
					RuleID:  "v1.shared_layer_classification.subtree.required",
					Path:    fmt.Sprintf("shared_layer_classification[%d].subtree", i),
					Message: "v1 entry missing subtree path",
				})
				continue
			}
			file.Subtrees = append(file.Subtrees, Subtree{
				Path:        s.Subtree,
				SharedLayer: s.Shared,
				// Owner / Purpose left empty: no v1 equivalent on disk.
			})
		}
		// ExpectedConsumers from v1 is intentionally NOT carried into
		// the normalized form. v2's frozen governance decision is that
		// the consumer list is derived from code references, not
		// maintained inline. Phase 1C migration drops this field; the
		// loader simulates that drop here by ignoring it on v1 reads.

	case 2:
		// v2: subtrees + consumer_tracking. Both can be empty (a brand
		// new project may have no subtrees declared yet); but if a
		// subtree IS declared it must have a path.
		file.ConsumerTracking = raw.ConsumerTracking
		for i, s := range raw.Subtrees {
			if strings.TrimSpace(s.Path) == "" {
				failures = append(failures, RepositoryTopologyValidationFailure{
					RuleID:  "v2.subtrees.path.required",
					Path:    fmt.Sprintf("subtrees[%d].path", i),
					Message: "v2 entry missing path",
				})
				continue
			}
			file.Subtrees = append(file.Subtrees, Subtree{
				Path:        s.Path,
				SharedLayer: s.SharedLayer,
				Owner:       s.Owner,
				Purpose:     s.Purpose,
			})
		}

	default:
		failures = append(failures, RepositoryTopologyValidationFailure{
			RuleID:  "schema_version.unsupported",
			Path:    "schema_version",
			Message: fmt.Sprintf("unsupported schema_version %d (loader supports 1 and 2)", version),
		})
	}

	if len(failures) > 0 {
		return file, &RepositoryTopologyValidationError{Failures: failures}
	}
	return file, nil
}

// WriteRepositoryTopology serializes the in-memory representation to YAML
// at the target path. Output is always v2 — Phase 1B is one-way migration
// at the writer boundary; v1 files round-tripping through Load → Write
// become v2 files.
//
// Validation: every subtree MUST declare a non-empty owner and purpose.
// A WriteRepositoryTopology call on a file freshly loaded from v1 YAML
// will fail this validation (because v1 has no owner/purpose on disk);
// callers performing a v1 → v2 migration must populate owner+purpose
// before write. This is intentional: Phase 1B explicitly does NOT
// auto-fabricate owner/purpose values during migration. Phase 1C's
// migration step provides them via the curated table in
// runtime/repository-topology-migration.md §v2 schema.
func WriteRepositoryTopology(path string, file *RepositoryTopologyFile) error {
	if err := validateForV2Write(file); err != nil {
		return err
	}
	body, err := marshalV2(file)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return fmt.Errorf("write repository topology file %s: %w", path, err)
	}
	return nil
}

func validateForV2Write(file *RepositoryTopologyFile) error {
	var failures []RepositoryTopologyValidationFailure
	for i, s := range file.Subtrees {
		if strings.TrimSpace(s.Path) == "" {
			failures = append(failures, RepositoryTopologyValidationFailure{
				RuleID:  "v2.subtrees.path.required",
				Path:    fmt.Sprintf("subtrees[%d].path", i),
				Message: "subtree path must be present and non-empty",
			})
		}
		if strings.TrimSpace(s.Owner) == "" {
			failures = append(failures, RepositoryTopologyValidationFailure{
				RuleID:  "v2.subtrees.owner.required",
				Path:    fmt.Sprintf("subtrees[%d].owner", i),
				Message: "v2 write requires non-empty owner for every subtree",
			})
		}
		if strings.TrimSpace(s.Purpose) == "" {
			failures = append(failures, RepositoryTopologyValidationFailure{
				RuleID:  "v2.subtrees.purpose.required",
				Path:    fmt.Sprintf("subtrees[%d].purpose", i),
				Message: "v2 write requires non-empty purpose for every subtree",
			})
		}
	}
	if len(failures) == 0 {
		return nil
	}
	return &RepositoryTopologyValidationError{Failures: failures}
}

// marshalV2 emits the file in v2 YAML shape. Subtrees are sorted by Path
// for deterministic output (round-trip stability + clean diff in version
// control); preserving the original order is not a v2 contract.
func marshalV2(file *RepositoryTopologyFile) ([]byte, error) {
	type outSubtree struct {
		Path        string `yaml:"path"`
		SharedLayer bool   `yaml:"shared_layer"`
		Owner       string `yaml:"owner"`
		Purpose     string `yaml:"purpose"`
	}
	type outFile struct {
		SchemaVersion     int                      `yaml:"schema_version"`
		Status            string                   `yaml:"status,omitempty"`
		OwnerLayer        string                   `yaml:"owner_layer,omitempty"`
		RuntimeProjection *RuntimeProjectionConfig `yaml:"runtime_projection,omitempty"`
		ConsumerTracking  *ConsumerTracking        `yaml:"consumer_tracking,omitempty"`
		Subtrees          []outSubtree             `yaml:"subtrees"`
		Invariants        []string                 `yaml:"invariants,omitempty"`
	}

	subtrees := make([]outSubtree, 0, len(file.Subtrees))
	for _, s := range file.Subtrees {
		subtrees = append(subtrees, outSubtree{
			Path:        s.Path,
			SharedLayer: s.SharedLayer,
			Owner:       s.Owner,
			Purpose:     s.Purpose,
		})
	}
	sort.SliceStable(subtrees, func(i, j int) bool {
		return subtrees[i].Path < subtrees[j].Path
	})

	out := outFile{
		SchemaVersion:     2,
		Status:            file.Status,
		OwnerLayer:        file.OwnerLayer,
		RuntimeProjection: file.RuntimeProjection,
		ConsumerTracking:  file.ConsumerTracking,
		Subtrees:          subtrees,
		Invariants:        file.Invariants,
	}

	return yaml.Marshal(&out)
}

// RepositoryTopologyValidationError aggregates schema-violation findings.
type RepositoryTopologyValidationError struct {
	Failures []RepositoryTopologyValidationFailure
}

// RepositoryTopologyValidationFailure is one rule failure occurrence.
type RepositoryTopologyValidationFailure struct {
	RuleID  string
	Path    string
	Message string
}

func (e *RepositoryTopologyValidationError) Error() string {
	if len(e.Failures) == 0 {
		return "repository topology validation: no failures (this is a bug)"
	}
	lines := []string{fmt.Sprintf("repository topology validation: %d failure(s)", len(e.Failures))}
	for _, f := range e.Failures {
		lines = append(lines, fmt.Sprintf("  - [%s] %s: %s", f.RuleID, f.Path, f.Message))
	}
	return strings.Join(lines, "\n")
}

// SubtreeForPath returns the declared subtree that governs a repo-relative
// path, by longest-prefix match (the most specific declared subtree wins), and
// whether any subtree matched. Used by authority classification to resolve a
// file's shared_layer / owner from topology. Matching is slash-normalized and
// treats a subtree path as a directory prefix (with exact-path match allowed).
func (f *RepositoryTopologyFile) SubtreeForPath(rel string) (Subtree, bool) {
	rel = strings.ReplaceAll(strings.TrimSpace(rel), "\\", "/")
	var best Subtree
	bestLen := -1
	for _, s := range f.Subtrees {
		p := strings.ReplaceAll(strings.TrimSpace(s.Path), "\\", "/")
		if p == "" {
			continue
		}
		prefix := p
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		if rel == strings.TrimSuffix(p, "/") || strings.HasPrefix(rel, prefix) {
			if len(p) > bestLen {
				best = s
				bestLen = len(p)
			}
		}
	}
	return best, bestLen >= 0
}

// IsRepositoryTopologyValidationError reports whether err is (or wraps)
// a *RepositoryTopologyValidationError. Uses errors.As so callers that
// wrap the validation error via fmt.Errorf("ctx: %w", err) still get the
// right answer. Mirrors IsValidationError in the project_metadata.go
// sibling parser.
func IsRepositoryTopologyValidationError(err error) bool {
	var v *RepositoryTopologyValidationError
	return errors.As(err, &v)
}
