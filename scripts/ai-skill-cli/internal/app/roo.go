package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const rooStorageKey = "RooVeterinaryInc.roo-cline"

const defaultRooCustomInstructions = `你是一個 AI 助手，遵循以下啟動流程：

## 啟動流程

1. 首先讀取 CORE_BOOTSTRAP.md 了解核心啟動流程
2. 依照 CORE_BOOTSTRAP.md 的指示載入依賴文件（dependency-reading.md）
3. 依照 dependency-reading.md 的規則載入 shared-rules 目錄下的必要規則
4. 根據任務類型載入對應的 workflow、intelligence、analysis 文件

## 語言偏好

- 預設使用英文回應
- 如果使用者使用中文提問，則用中文回應
- 如果使用者切換語言，跟隨其切換
- 所有輸出（包含技術分析、表格欄位、章節標題）都必須與使用者當前語言一致

## 對話目標閉環

- 每個對話必須有明確的目標
- 完成目標後使用 attempt_completion 呈現結果
- 不要在結果結尾提出問題或要求進一步協助

## 專案設定

此 Ai-skill 專案位於使用者的 Documents/Ai-skill 目錄。
當使用者在此專案中工作時，遵循專案內的 .roomodes 和 shared-rules 設定。
當使用者在其他專案中工作時，此全域設定確保啟動流程仍然有效。`

type rooOptions struct {
	command            string
	dbPath             string
	instructionsFile   string
	allowRunningVSCode bool
	jsonOutput         bool
	plainOutput        bool
}

func runRoo(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: ai-skill roo <set-global-custom-instructions> [flags]")
		return ExitInvalidUsage
	}
	opts := rooOptions{command: args[0]}
	if opts.command != "set-global-custom-instructions" {
		_, _ = fmt.Fprintf(stderr, "unsupported roo command: %s\n", opts.command)
		return ExitInvalidUsage
	}
	fs := newFlagSet("roo "+opts.command, stderr)
	fs.StringVar(&opts.dbPath, "db", defaultRooDBPath(), "VS Code globalStorage state.vscdb path")
	fs.StringVar(&opts.instructionsFile, "instructions-file", "", "optional file containing customInstructions content")
	fs.BoolVar(&opts.allowRunningVSCode, "allow-running-vscode", false, "allow writing while VS Code appears to be running")
	fs.BoolVar(&opts.jsonOutput, "json", false, "write machine-readable JSON output")
	fs.BoolVar(&opts.plainOutput, "plain", false, "write human-readable output")
	if err := fs.Parse(args[1:]); err != nil {
		return ExitInvalidUsage
	}
	if opts.jsonOutput && opts.plainOutput {
		_, _ = fmt.Fprintln(stderr, "--json and --plain are mutually exclusive")
		return ExitInvalidUsage
	}
	result := buildRooGlobalInstructionsResult(opts)
	if opts.jsonOutput {
		if err := writeJSON(stdout, result); err != nil {
			_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
			return ExitGeneralFailure
		}
		return result.ExitCode
	}
	if err := writePlain(stdout, result); err != nil {
		_, _ = fmt.Fprintf(stderr, "write output: %v\n", err)
		return ExitGeneralFailure
	}
	return result.ExitCode
}

func buildRooGlobalInstructionsResult(opts rooOptions) Result {
	result := Result{
		Command:        "roo set-global-custom-instructions",
		Mode:           "guarded_write",
		Status:         "success",
		ExitCode:       ExitSuccess,
		Checks:         []Check{},
		PlannedActions: []string{"update Roo Code global customInstructions in VS Code state database"},
		Mutations:      []string{},
	}
	if !opts.allowRunningVSCode {
		check := rooVSCodeNotRunningCheck()
		result.Checks = append(result.Checks, check)
		if check.Status != "ok" {
			result.Status = "blocked"
			result.ExitCode = ExitValidationFailed
			result.Error = &CommandError{Code: "vscode_running", Message: check.Message, Remediation: "Close VS Code completely, then rerun the command."}
			return result
		}
	}
	dbCheck := fileExistsCheck("roo_state_db", opts.dbPath)
	result.Checks = append(result.Checks, dbCheck)
	if dbCheck.Status != "ok" {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "missing_roo_state_db", Message: dbCheck.Message, Remediation: "Install and run Roo Code at least once, or pass --db with the correct state.vscdb path."}
		return result
	}
	instructions := defaultRooCustomInstructions
	if strings.TrimSpace(opts.instructionsFile) != "" {
		content, err := os.ReadFile(opts.instructionsFile)
		if err != nil {
			result.Status = "blocked"
			result.ExitCode = ExitInvalidUsage
			result.Error = &CommandError{Code: "missing_instructions_file", Message: err.Error(), Remediation: "Pass an existing --instructions-file path."}
			return result
		}
		instructions = string(content)
	}
	previousLength, newLength, err := writeRooCustomInstructions(opts.dbPath, instructions)
	if err != nil {
		result.Status = "blocked"
		result.ExitCode = ExitValidationFailed
		result.Error = &CommandError{Code: "roo_write_failed", Message: err.Error(), Remediation: "Check ItemTable, JSON content, and database permissions."}
		return result
	}
	result.Checks = append(result.Checks, Check{Name: "roo_custom_instructions", Status: "ok", Message: fmt.Sprintf("customInstructions length %d; previous length %d", newLength, previousLength)})
	result.Mutations = append(result.Mutations, opts.dbPath)
	return result
}

func writeRooCustomInstructions(dbPath string, instructions string) (int, int, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0, 0, err
	}
	defer db.Close()
	var raw string
	if err := db.QueryRow(`SELECT value FROM ItemTable WHERE key = ?`, rooStorageKey).Scan(&raw); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("Roo Code settings row not found: %s", rooStorageKey)
		}
		return 0, 0, err
	}
	data := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return 0, 0, fmt.Errorf("invalid Roo settings JSON: %w", err)
	}
	previousLength := len(runtimeString(data["customInstructions"]))
	data["customInstructions"] = instructions
	updated, err := json.Marshal(data)
	if err != nil {
		return 0, 0, err
	}
	if _, err := db.Exec(`UPDATE ItemTable SET value = ? WHERE key = ?`, string(updated), rooStorageKey); err != nil {
		return 0, 0, err
	}
	if _, err := db.Exec(`PRAGMA wal_checkpoint(TRUNCATE)`); err != nil {
		return 0, 0, err
	}
	var verifyRaw string
	if err := db.QueryRow(`SELECT value FROM ItemTable WHERE key = ?`, rooStorageKey).Scan(&verifyRaw); err != nil {
		return 0, 0, err
	}
	verify := map[string]any{}
	if err := json.Unmarshal([]byte(verifyRaw), &verify); err != nil {
		return 0, 0, err
	}
	if runtimeString(verify["customInstructions"]) != instructions {
		return 0, 0, fmt.Errorf("customInstructions verification failed")
	}
	return previousLength, len(instructions), nil
}

func defaultRooDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "state.vscdb")
}

func rooVSCodeNotRunningCheck() Check {
	path, err := exec.LookPath("pgrep")
	if err != nil {
		return Check{Name: "vscode_not_running", Status: "ok", Message: "pgrep unavailable; skipped VS Code process check"}
	}
	cmd := exec.Command(path, "-f", "Visual Studio Code")
	output, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		return Check{Name: "vscode_not_running", Status: "failed", Message: "VS Code appears to be running"}
	}
	return Check{Name: "vscode_not_running", Status: "ok", Message: "VS Code process not detected"}
}

func fileExistsCheck(name string, path string) Check {
	if strings.TrimSpace(path) == "" {
		return Check{Name: name, Status: "missing", Message: "path is empty"}
	}
	info, err := os.Stat(path)
	if err != nil {
		return Check{Name: name, Status: "missing", Message: path}
	}
	if info.IsDir() {
		return Check{Name: name, Status: "failed", Message: path + " is a directory"}
	}
	return Check{Name: name, Status: "ok", Message: path}
}
