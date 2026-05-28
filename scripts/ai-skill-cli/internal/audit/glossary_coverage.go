package audit

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// glossaryCoveragePaths lists repo-relative directories the heuristic scans.
// Plan §Resolved Decisions Q7: scan exactly these seven paths.
var glossaryCoveragePaths = []string{
	"plans/active",
	"architecture",
	"workflow",
	"analysis",
	"intelligence",
	"runtime",
	"ecosystem",
}

// backtickRe matches backtick-wrapped identifiers; the captured group is the
// candidate term (without the surrounding backticks). Path separators are
// excluded — backtick path references like `plans/active/foo.md` are not
// candidates for glossary canonicalisation.
var backtickRe = regexp.MustCompile("`([A-Za-z][A-Za-z0-9_.\\-]*)`")

// snakeCaseRe matches identifiers with at least two snake_case segments and
// no spaces. Anchored to non-word boundaries via lookbehind-free matching.
var snakeCaseRe = regexp.MustCompile(`\b([a-z][a-z0-9]+(?:_[a-z0-9]+){1,}[a-z0-9]*)\b`)

// glossaryCandidate is an unmatched framework-looking term in a scanned file.
type glossaryCandidate struct {
	Term string
	Path string
	Line int
}

// scanGlossaryCoverage walks the seven configured paths and emits warning
// strings for backtick / snake_case ≥ 2 segments candidates that are not in
// glossary_terms.term or glossary_terms.aliases. Warnings are returned in
// stable sorted order. The scan never blocks; missing SQLite or missing
// directories degrade to empty results.
func scanGlossaryCoverage(repo string) []string {
	knownTerms, err := readKnownGlossaryTerms(repo)
	if err != nil {
		return nil
	}
	var candidates []glossaryCandidate
	for _, rel := range glossaryCoveragePaths {
		root := filepath.Join(repo, rel)
		if _, statErr := os.Stat(root); statErr != nil {
			continue
		}
		_ = filepath.Walk(root, func(p string, info os.FileInfo, walkErr error) error {
			if walkErr != nil || info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(p, ".md") && !strings.HasSuffix(p, ".yaml") && !strings.HasSuffix(p, ".yml") {
				return nil
			}
			b, readErr := os.ReadFile(p)
			if readErr != nil {
				return nil
			}
			scanFileForCandidates(string(b), p, repo, knownTerms, &candidates)
			return nil
		})
	}
	if len(candidates) == 0 {
		return nil
	}
	return summariseCandidatesByFrequency(candidates, 50)
}

// summariseCandidatesByFrequency collapses multiple occurrences of the same
// term into one warning, sorted by frequency descending. Each warning records
// the first occurrence location and total count. The limit caps the noise
// from large repos; remaining terms are mentioned in a tail summary so the
// reviewer knows more exist.
func summariseCandidatesByFrequency(in []glossaryCandidate, limit int) []string {
	type stat struct {
		Term     string
		Path     string
		Line     int
		Count    int
		FirstSeq int
	}
	bucket := map[string]*stat{}
	for i, c := range in {
		s, ok := bucket[c.Term]
		if !ok {
			s = &stat{Term: c.Term, Path: c.Path, Line: c.Line, FirstSeq: i}
			bucket[c.Term] = s
		}
		s.Count++
	}
	stats := make([]*stat, 0, len(bucket))
	for _, s := range bucket {
		stats = append(stats, s)
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count != stats[j].Count {
			return stats[i].Count > stats[j].Count
		}
		return stats[i].FirstSeq < stats[j].FirstSeq
	})
	if len(stats) == 0 {
		return nil
	}
	out := make([]string, 0, len(stats)+1)
	keep := limit
	if keep > len(stats) {
		keep = len(stats)
	}
	for i := 0; i < keep; i++ {
		s := stats[i]
		out = append(out, fmt.Sprintf("glossary candidate `%s` (×%d, first at %s:%d) not in glossary_terms or aliases", s.Term, s.Count, s.Path, s.Line))
	}
	if len(stats) > limit {
		out = append(out, fmt.Sprintf("... and %d more unique candidate terms (output truncated to top %d by frequency)", len(stats)-limit, limit))
	}
	return out
}

func scanFileForCandidates(content, path, repo string, known map[string]struct{}, out *[]glossaryCandidate) {
	rel, _ := filepath.Rel(repo, path)
	for lineNo, line := range strings.Split(content, "\n") {
		for _, m := range backtickRe.FindAllStringSubmatch(line, -1) {
			term := m[1]
			if !looksLikeFrameworkTerm(term) {
				continue
			}
			if _, ok := known[strings.ToLower(term)]; ok {
				continue
			}
			*out = append(*out, glossaryCandidate{Term: term, Path: filepath.ToSlash(rel), Line: lineNo + 1})
		}
		for _, m := range snakeCaseRe.FindAllStringSubmatch(line, -1) {
			term := m[1]
			if _, ok := known[strings.ToLower(term)]; ok {
				continue
			}
			*out = append(*out, glossaryCandidate{Term: term, Path: filepath.ToSlash(rel), Line: lineNo + 1})
		}
	}
}

// looksLikeFrameworkTerm filters backtick matches to framework-looking
// identifiers. Requires either at least one separator (_./-) OR mixed case
// with length ≥ 4, so single English words like `note` or `the` do not fire.
func looksLikeFrameworkTerm(term string) bool {
	if len(term) < 3 {
		return false
	}
	if strings.ContainsAny(term, "_./-") {
		return true
	}
	hasUpper := false
	hasLower := false
	for _, r := range term {
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
		}
		if r >= 'a' && r <= 'z' {
			hasLower = true
		}
	}
	return hasUpper && hasLower && len(term) >= 4
}

func dedupCandidates(in []glossaryCandidate) []glossaryCandidate {
	seen := map[string]bool{}
	out := make([]glossaryCandidate, 0, len(in))
	for _, c := range in {
		key := c.Term + "|" + c.Path + "|" + fmt.Sprint(c.Line)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		if out[i].Line != out[j].Line {
			return out[i].Line < out[j].Line
		}
		return out[i].Term < out[j].Term
	})
	return out
}

// readKnownGlossaryTerms returns the set of known terms (lowercased) from the
// glossary_terms.term column plus expanded aliases. Aliases may be stored as
// JSON array, comma-separated, or null; we accept any of those shapes.
func readKnownGlossaryTerms(repo string) (map[string]struct{}, error) {
	dbPath := filepath.Join(repo, "knowledge", "runtime", "sqlite", "runtime-index.sqlite")
	if _, err := os.Stat(dbPath); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT term, COALESCE(aliases, '') FROM glossary_terms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]struct{}{}
	for rows.Next() {
		var term, aliases string
		if err := rows.Scan(&term, &aliases); err != nil {
			return nil, err
		}
		if term != "" {
			out[strings.ToLower(term)] = struct{}{}
		}
		for _, a := range splitAliases(aliases) {
			if a != "" {
				out[strings.ToLower(a)] = struct{}{}
			}
		}
	}
	return out, nil
}

// splitAliases accepts JSON-array, comma-separated, or whitespace-separated
// aliases and returns the cleaned list. Bracket / quote characters are
// stripped; empties are skipped.
func splitAliases(raw string) []string {
	if raw == "" {
		return nil
	}
	cleaned := strings.NewReplacer("[", "", "]", "", "\"", "", "'", "").Replace(raw)
	parts := strings.FieldsFunc(cleaned, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
