package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRooSetGlobalCustomInstructionsWritesFakeStateDB(t *testing.T) {
	dbPath := createRooStateDBFixture(t, `{"language":"zh","customInstructions":"old"}`)
	instructionsPath := filepath.Join(t.TempDir(), "instructions.txt")
	if err := os.WriteFile(instructionsPath, []byte("new custom instructions"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"roo", "set-global-custom-instructions", "--db", dbPath, "--instructions-file", instructionsPath, "--allow-running-vscode", "--json"}, &stdout, &stderr)
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}

	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !hasCheckStatus(result.Checks, "roo_custom_instructions", "ok") {
		t.Fatalf("expected roo custom instructions check, got %#v", result.Checks)
	}
	if len(result.Mutations) != 1 || result.Mutations[0] != dbPath {
		t.Fatalf("expected db mutation, got %#v", result.Mutations)
	}
	if got := rooStateJSONValue(t, dbPath, "customInstructions"); got != "new custom instructions" {
		t.Fatalf("customInstructions = %q", got)
	}
	if got := rooStateJSONValue(t, dbPath, "language"); got != "zh" {
		t.Fatalf("language = %q", got)
	}
}

func TestRooSetGlobalCustomInstructionsBlocksMissingDB(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"roo", "set-global-custom-instructions", "--db", filepath.Join(t.TempDir(), "missing.vscdb"), "--allow-running-vscode", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}
	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "missing_roo_state_db" {
		t.Fatalf("expected missing_roo_state_db, got %#v", result.Error)
	}
}

func TestRooSetGlobalCustomInstructionsBlocksMissingRow(t *testing.T) {
	dbPath := createRooStateDBFixture(t, "")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"roo", "set-global-custom-instructions", "--db", dbPath, "--allow-running-vscode", "--json"}, &stdout, &stderr)
	if code != ExitValidationFailed {
		t.Fatalf("expected validation failure, got %d; stderr=%s; stdout=%s", code, stderr.String(), stdout.String())
	}
	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if result.Error == nil || result.Error.Code != "roo_write_failed" || !strings.Contains(result.Error.Message, "settings row not found") {
		t.Fatalf("expected missing row failure, got %#v", result.Error)
	}
}

func createRooStateDBFixture(t *testing.T, value string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "state.vscdb")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE ItemTable (key TEXT PRIMARY KEY, value TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	if value != "" {
		if _, err := db.Exec(`INSERT INTO ItemTable (key, value) VALUES (?, ?)`, rooStorageKey, value); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

func rooStateJSONValue(t *testing.T, path string, field string) string {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var raw string
	if err := db.QueryRow(`SELECT value FROM ItemTable WHERE key = ?`, rooStorageKey).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	data := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		t.Fatal(err)
	}
	return runtimeString(data[field])
}
