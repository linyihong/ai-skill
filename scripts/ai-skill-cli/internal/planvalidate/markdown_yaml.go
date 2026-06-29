package planvalidate

import "strings"

// ExtractFencedYAMLBlocks returns the inner text of ```yaml ... ``` fences in order.
func ExtractFencedYAMLBlocks(markdown string) []string {
	lines := strings.Split(markdown, "\n")
	var blocks []string
	inFence := false
	var current []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "```yaml" {
			inFence = true
			current = nil
			continue
		}
		if inFence && trimmed == "```" {
			inFence = false
			blocks = append(blocks, strings.Join(current, "\n"))
			continue
		}
		if inFence {
			current = append(current, line)
		}
	}
	return blocks
}
