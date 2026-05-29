package app

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- ScanCheckboxes unit tests ---

func TestScanCheckboxesPureContent(t *testing.T) {
	content := "# Plan\n\n- [ ] task A\n- [x] task B\n- [ ] task C\n"
	result := ScanCheckboxes(content)
	if len(result.UncheckedLines) != 2 {
		t.Errorf("expected 2 unchecked, got %d", len(result.UncheckedLines))
	}
	if len(result.CheckedLines) != 1 {
		t.Errorf("expected 1 checked, got %d", len(result.CheckedLines))
	}
}

func TestScanCheckboxesAllChecked(t *testing.T) {
	content := "- [x] done\n- [X] also done\n"
	result := ScanCheckboxes(content)
	if result.HasUnchecked() {
		t.Error("expected no unchecked items")
	}
	if len(result.CheckedLines) != 2 {
		t.Errorf("expected 2 checked, got %d", len(result.CheckedLines))
	}
}

func TestScanCheckboxesEmpty(t *testing.T) {
	result := ScanCheckboxes("")
	if result.HasUnchecked() {
		t.Error("expected no unchecked on empty input")
	}
}

func TestScanCheckboxesIndented(t *testing.T) {
	content := "  - [ ] indented unchecked\n    - [x] deeper checked\n"
	result := ScanCheckboxes(content)
	if len(result.UncheckedLines) != 1 {
		t.Errorf("expected 1 unchecked (indented), got %d", len(result.UncheckedLines))
	}
	if len(result.CheckedLines) != 1 {
		t.Errorf("expected 1 checked (indented), got %d", len(result.CheckedLines))
	}
}

func TestScanCheckboxesInFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checklist.md")
	_ = os.WriteFile(path, []byte("- [ ] open\n- [x] done\n"), 0o644)
	result, err := ScanCheckboxesInFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.File != path {
		t.Errorf("expected File=%q, got %q", path, result.File)
	}
	if len(result.UncheckedLines) != 1 {
		t.Errorf("expected 1 unchecked, got %d", len(result.UncheckedLines))
	}
}

func TestScanCheckboxesInFileMissing(t *testing.T) {
	_, err := ScanCheckboxesInFile("/nonexistent/path/checklist.md")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// --- runScanCheckboxes CLI tests ---

func TestRunScanCheckboxesPlain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plan.md")
	_ = os.WriteFile(path, []byte("- [ ] task A\n- [x] task B\n"), 0o644)

	var out bytes.Buffer
	code := runScanCheckboxes([]string{path}, &out, io.Discard)
	if code != ExitSuccess {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), "unchecked: 1") {
		t.Errorf("expected unchecked count in output, got %q", out.String())
	}
}

func TestRunScanCheckboxesExitCodeWithUnchecked(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plan.md")
	_ = os.WriteFile(path, []byte("- [ ] unfinished\n"), 0o644)

	var out bytes.Buffer
	code := runScanCheckboxes([]string{"--exit-code", path}, &out, io.Discard)
	if code != ExitGeneralFailure {
		t.Errorf("expected exit 1 on unchecked items with --exit-code, got %d", code)
	}
}

func TestRunScanCheckboxesExitCodeAllChecked(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plan.md")
	_ = os.WriteFile(path, []byte("- [x] all done\n"), 0o644)

	var out bytes.Buffer
	code := runScanCheckboxes([]string{"--exit-code", path}, &out, io.Discard)
	if code != ExitSuccess {
		t.Errorf("expected exit 0 when all checked, got %d", code)
	}
}

func TestRunScanCheckboxesJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plan.md")
	_ = os.WriteFile(path, []byte("- [ ] task\n"), 0o644)

	var out bytes.Buffer
	code := runScanCheckboxes([]string{"--format", "json", path}, &out, io.Discard)
	if code != ExitSuccess {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), `"unchecked_lines"`) {
		t.Errorf("expected JSON output with unchecked_lines, got %q", out.String())
	}
}
