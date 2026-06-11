package app

import "strings"

// Link is an inline markdown link extracted by extractMarkdownLinks.
// Target is the raw path string (after escape resolution); Line and
// Column are 1-based and point at the opening '['.
type Link struct {
	Target string
	Line   int
	Column int
}

// extractMarkdownLinks parses content and returns inline markdown links
// `[text](path)` and `[text](path "title")` whose target is a relative
// repo path. Absolute URLs, mailto:/tel: links, and pure anchors are
// recognised but filtered out of the result.
//
// This is a bounded parser by design (see Parser Strategy in
// plans/active/2026-06-11-1100-plan-archival-link-integrity.md):
//
//	Supported:
//	  - inline link [text](path)
//	  - titled inline link [text](path "title")
//	  - escaped parens in path: [text](../a\(b\).md)
//	  - code-fence exclusion (``` and ~~~ at line start, leading spaces ok)
//
//	Not supported (ignored, not partially interpreted):
//	  - reference-style links [text][ref] + [ref]: path
//	  - HTML <a href="...">
//	  - autolinks <https://...>
func extractMarkdownLinks(content []byte) []Link {
	var links []Link
	lines := strings.Split(string(content), "\n")

	inFence := false
	fenceMarker := ""

	for lineIdx, lineText := range lines {
		trimmed := strings.TrimLeft(lineText, " \t")
		if strings.HasPrefix(trimmed, "```") {
			if !inFence {
				inFence = true
				fenceMarker = "```"
			} else if fenceMarker == "```" {
				inFence = false
				fenceMarker = ""
			}
			continue
		}
		if strings.HasPrefix(trimmed, "~~~") {
			if !inFence {
				inFence = true
				fenceMarker = "~~~"
			} else if fenceMarker == "~~~" {
				inFence = false
				fenceMarker = ""
			}
			continue
		}
		if inFence {
			continue
		}
		links = append(links, scanLineForInlineLinks(lineText, lineIdx+1)...)
	}
	return links
}

// scanLineForInlineLinks scans a single line for `[text](path)` patterns.
// Link text may not contain unescaped ']' or newlines (newlines are
// impossible since we operate per-line). Reference-style and shortcut
// links are recognised by the missing '(' after ']' and skipped.
func scanLineForInlineLinks(line string, lineNum int) []Link {
	var out []Link
	n := len(line)
	i := 0
	for i < n {
		if line[i] != '[' {
			i++
			continue
		}
		startCol := i + 1
		j := i + 1
		depth := 1
		for j < n {
			if line[j] == '\\' && j+1 < n {
				j += 2
				continue
			}
			if line[j] == '[' {
				depth++
			} else if line[j] == ']' {
				depth--
				if depth == 0 {
					break
				}
			}
			j++
		}
		if depth != 0 || j >= n {
			i++
			continue
		}
		if j+1 >= n || line[j+1] != '(' {
			i = j + 1
			continue
		}

		k := j + 2
		var pathBuf strings.Builder
		ok := false
		for k < n {
			ch := line[k]
			if ch == '\\' && k+1 < n {
				nx := line[k+1]
				if nx == '(' || nx == ')' || nx == '\\' {
					pathBuf.WriteByte(nx)
					k += 2
					continue
				}
				pathBuf.WriteByte(ch)
				k++
				continue
			}
			if ch == ' ' || ch == '\t' {
				kk := k
				for kk < n && line[kk] != ')' {
					kk++
				}
				if kk < n {
					ok = true
					k = kk + 1
				}
				break
			}
			if ch == ')' {
				ok = true
				k++
				break
			}
			pathBuf.WriteByte(ch)
			k++
		}
		if !ok {
			i++
			continue
		}

		target := strings.TrimSpace(pathBuf.String())
		if shouldKeepLinkTarget(target) {
			out = append(out, Link{
				Target: target,
				Line:   lineNum,
				Column: startCol,
			})
		}
		i = k
	}
	return out
}

// shouldKeepLinkTarget filters out non-relative-path targets. This is
// the boundary between "links we validate" and "links we ignore".
func shouldKeepLinkTarget(t string) bool {
	if t == "" {
		return false
	}
	if strings.HasPrefix(t, "#") {
		return false
	}
	if strings.Contains(t, "://") {
		return false
	}
	if strings.HasPrefix(t, "mailto:") || strings.HasPrefix(t, "tel:") {
		return false
	}
	return true
}
