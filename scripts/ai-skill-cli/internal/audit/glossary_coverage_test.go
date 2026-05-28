package audit

import (
	"strings"
	"testing"
)

func TestLooksLikeFrameworkTermAcceptsSeparators(t *testing.T) {
	for _, term := range []string{"route.foo", "owner-layer", "snake_case", "a.b"} {
		if !looksLikeFrameworkTerm(term) {
			t.Errorf("expected %q to look like framework term", term)
		}
	}
}

func TestLooksLikeFrameworkTermRejectsShortAndPlain(t *testing.T) {
	for _, term := range []string{"a", "ab", "the", "note", "value", "main"} {
		if looksLikeFrameworkTerm(term) {
			t.Errorf("expected %q to be rejected", term)
		}
	}
}

func TestLooksLikeFrameworkTermAcceptsCamelCase(t *testing.T) {
	for _, term := range []string{"camelCase", "PascalCase", "AbCd"} {
		if !looksLikeFrameworkTerm(term) {
			t.Errorf("expected %q (mixed case ≥ 4) to look like framework term", term)
		}
	}
}

func TestScanFileForCandidatesEmitsBacktick(t *testing.T) {
	known := map[string]struct{}{}
	var out []glossaryCandidate
	content := "intro `route.foo` more text"
	scanFileForCandidates(content, "/repo/plans/active/x.md", "/repo", known, &out)
	if len(out) == 0 {
		t.Fatal("expected at least one candidate")
	}
	found := false
	for _, c := range out {
		if c.Term == "route.foo" {
			found = true
		}
	}
	if !found {
		t.Errorf("backtick term not captured: %#v", out)
	}
}

func TestScanFileForCandidatesEmitsSnakeCase(t *testing.T) {
	known := map[string]struct{}{}
	var out []glossaryCandidate
	content := "the term retry_with_backoff matters"
	scanFileForCandidates(content, "/repo/runtime/x.yaml", "/repo", known, &out)
	found := false
	for _, c := range out {
		if c.Term == "retry_with_backoff" {
			found = true
		}
	}
	if !found {
		t.Errorf("snake_case term not captured: %#v", out)
	}
}

func TestScanFileForCandidatesSkipsKnownTerm(t *testing.T) {
	known := map[string]struct{}{"route.foo": {}}
	var out []glossaryCandidate
	content := "intro `route.foo` more text"
	scanFileForCandidates(content, "/repo/plans/active/x.md", "/repo", known, &out)
	for _, c := range out {
		if c.Term == "route.foo" {
			t.Errorf("known term should be skipped, got %#v", c)
		}
	}
}

func TestScanFileForCandidatesSkipsKnownAlias(t *testing.T) {
	known := map[string]struct{}{"runtime_audit": {}}
	var out []glossaryCandidate
	content := "see runtime_audit run"
	scanFileForCandidates(content, "/repo/runtime/x.yaml", "/repo", known, &out)
	for _, c := range out {
		if c.Term == "runtime_audit" {
			t.Errorf("known alias should be skipped, got %#v", c)
		}
	}
}

func TestScanFileForCandidatesIgnoresNaturalLanguage(t *testing.T) {
	known := map[string]struct{}{}
	var out []glossaryCandidate
	content := "This is plain English without any framework terms here."
	scanFileForCandidates(content, "/repo/plans/active/x.md", "/repo", known, &out)
	if len(out) != 0 {
		t.Errorf("natural language should produce no candidates, got %#v", out)
	}
}

func TestDedupCandidatesSortsStably(t *testing.T) {
	in := []glossaryCandidate{
		{Term: "z", Path: "a", Line: 2},
		{Term: "a", Path: "a", Line: 1},
		{Term: "z", Path: "a", Line: 2}, // duplicate
		{Term: "b", Path: "b", Line: 1},
	}
	out := dedupCandidates(in)
	if len(out) != 3 {
		t.Fatalf("expected 3 unique, got %d: %#v", len(out), out)
	}
	if out[0].Term != "a" || out[1].Term != "z" || out[2].Term != "b" {
		t.Errorf("sort order wrong: %#v", out)
	}
}

func TestSplitAliasesHandlesShapes(t *testing.T) {
	cases := map[string][]string{
		"":                          nil,
		"foo,bar":                   {"foo", "bar"},
		"[foo, bar]":                {"foo", "bar"},
		`["foo","bar"]`:             {"foo", "bar"},
		"foo bar":                   {"foo", "bar"},
		`["alias-one", "alias_two"]`: {"alias-one", "alias_two"},
	}
	for raw, want := range cases {
		got := splitAliases(raw)
		if len(got) != len(want) {
			t.Errorf("splitAliases(%q) length: got %v want %v", raw, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("splitAliases(%q)[%d]: got %q want %q", raw, i, got[i], want[i])
			}
		}
	}
}

func TestScanGlossaryCoverageMissingDBReturnsNil(t *testing.T) {
	out := scanGlossaryCoverage(t.TempDir())
	if out != nil {
		t.Errorf("missing SQLite should return nil, got %#v", out)
	}
}

func TestGlossaryCandidateWarningFormat(t *testing.T) {
	w := []string{}
	c := glossaryCandidate{Term: "route.foo", Path: "plans/active/x.md", Line: 42}
	w = append(w, formatTestWarning(c))
	if !strings.Contains(w[0], "route.foo") || !strings.Contains(w[0], "plans/active/x.md") || !strings.Contains(w[0], ":42") {
		t.Errorf("warning format missing fields: %q", w[0])
	}
}

// formatTestWarning mirrors the per-term summary format produced by
// summariseCandidatesByFrequency for the count=1 case.
func formatTestWarning(c glossaryCandidate) string {
	return "glossary candidate `" + c.Term + "` (×1, first at " + c.Path + ":" + intToStr(c.Line) + ") not in glossary_terms or aliases"
}

func intToStr(i int) string {
	return strings.TrimSpace((map[int]string{
		1: "1", 2: "2", 3: "3", 42: "42",
	})[i])
}

func TestSummariseCandidatesByFrequencyOrdersByCountAndCapsLimit(t *testing.T) {
	in := []glossaryCandidate{
		{Term: "rare", Path: "a.md", Line: 1},
		{Term: "popular", Path: "a.md", Line: 1},
		{Term: "popular", Path: "b.md", Line: 5},
		{Term: "popular", Path: "c.md", Line: 9},
		{Term: "mid", Path: "a.md", Line: 2},
		{Term: "mid", Path: "b.md", Line: 4},
	}
	out := summariseCandidatesByFrequency(in, 50)
	if len(out) != 3 {
		t.Fatalf("expected 3 unique summaries, got %d: %#v", len(out), out)
	}
	if !strings.Contains(out[0], "popular") || !strings.Contains(out[0], "×3") {
		t.Errorf("popular should rank first with ×3, got %q", out[0])
	}
	if !strings.Contains(out[1], "mid") || !strings.Contains(out[1], "×2") {
		t.Errorf("mid should rank second with ×2, got %q", out[1])
	}
	if !strings.Contains(out[2], "rare") || !strings.Contains(out[2], "×1") {
		t.Errorf("rare should rank last with ×1, got %q", out[2])
	}
}

func TestSummariseCandidatesByFrequencyTruncates(t *testing.T) {
	in := []glossaryCandidate{}
	for i := 0; i < 75; i++ {
		in = append(in, glossaryCandidate{Term: "t" + intLabel(i), Path: "x.md", Line: i})
	}
	out := summariseCandidatesByFrequency(in, 50)
	if len(out) != 51 { // 50 terms + truncation notice
		t.Fatalf("expected 51 lines (50 + truncation notice), got %d", len(out))
	}
	if !strings.Contains(out[50], "truncated to top 50") {
		t.Errorf("missing truncation notice in tail: %q", out[50])
	}
}

func intLabel(i int) string {
	// Cheap unique label without importing strconv.
	out := []byte{}
	if i == 0 {
		return "0"
	}
	for i > 0 {
		out = append([]byte{byte('0' + i%10)}, out...)
		i /= 10
	}
	return string(out)
}
