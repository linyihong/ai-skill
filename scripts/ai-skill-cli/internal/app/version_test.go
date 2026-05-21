package app

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionJSONReportsBuildMetadata(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, Date
	Version, Commit, Date = "v1.2.3", "abc123", "2026-05-21T00:00:00Z"
	defer func() {
		Version, Commit, Date = oldVersion, oldCommit, oldDate
	}()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"version", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s", code, stderr.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Command != "version" || result.Mode != "info" {
		t.Fatalf("unexpected result identity: %#v", result)
	}
	if !hasCheckStatus(result.Checks, "version", "ok") {
		t.Fatalf("expected version check, got %#v", result.Checks)
	}
}
