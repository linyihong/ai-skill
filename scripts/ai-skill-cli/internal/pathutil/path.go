package pathutil

import (
	"errors"
	"path"
	"path/filepath"
	"strings"
)

type PathListSummary struct {
	Entries      int
	EmptyEntries int
}

func NormalizeForReport(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("path is empty")
	}

	normalized := strings.ReplaceAll(input, "\\", "/")
	isUNC := strings.HasPrefix(normalized, "//")
	cleaned := path.Clean(normalized)
	if isUNC && strings.HasPrefix(cleaned, "/") && !strings.HasPrefix(cleaned, "//") {
		cleaned = "/" + cleaned
	}
	return cleaned, nil
}

func SummarizePathList(value string) PathListSummary {
	if value == "" {
		return PathListSummary{}
	}

	entries := filepath.SplitList(value)
	summary := PathListSummary{Entries: len(entries)}
	for _, entry := range entries {
		if strings.TrimSpace(entry) == "" {
			summary.EmptyEntries++
		}
	}
	return summary
}
