# Shell Script Added Without Go Migration（新增 .sh 腳本而非 Go 實作）

Status: validated
Class: `process-gap` / `platform-governance-miss`

## Trigger

Agent 在新增 hook、工具腳本、CLI 輔助功能時，直接建立 `.sh` 檔案，未先查閱平台治理政策（`plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md`），導致跨平台問題被引入，且 Go CLI 遷移路徑未被規劃。

具體觸發信號：

- `git diff --cached --diff-filter=A --name-only` 顯示新增 `.sh` 檔案
- Commit message 無說明「為何不用 Go 實作」或 Go migration plan
- 新增 `.claude/hooks/*.sh`、`scripts/*.sh` 等平台相依腳本

## Failure Mode

Bash scripts 只能在 macOS / Linux 執行，Windows 需要 Git Bash / WSL。本 repo 的治理政策是：**所有新增腳本邏輯必須在 Go CLI 中實作**，並透過 `ai-skill hooks install` 將平台正確的 binary 路徑寫入 `.claude/settings.json`。

1. **跨平台 breakage**：.sh 在 Windows 無法直接執行
2. **遷移債務累積**：每新增一個 .sh，Go migration 工作量就增加
3. **治理政策被繞過**：platform governance 形同虛設

## Required Agent Action

1. **修改 `.claude/hooks/*.sh` 或新增 `scripts/*.sh` 前**，先讀 [`plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md`](../../plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md) §platform-governance
2. 新邏輯必須實作在 `scripts/ai-skill-cli/internal/app/hooks.go` 或新 subcommand
3. `ai-skill hooks install` 負責在 `.claude/settings.json` 寫入平台正確的 binary 路徑
4. 若確有短暫過渡需要（例如等待 Go 實作完成前的 wrapper），在 commit message 加 `[skip-go-migration]` 並說明 Go migration plan

## Prevention Gate

Pre-commit hook `validateNoNewShellScripts`（`scripts/ai-skill-cli/internal/app/hooks.go`）：
- 偵測 Added `.sh` 檔案 → block with `new_shell_script_forbidden`
- Opt-out：`[skip-go-migration]` 獨立行於 commit message body
- 修改既有 `.sh`（pending migration）不觸發

## Validation

符合下列條件時，此 pattern 已被防止：

- `git diff --cached --diff-filter=A --name-only | grep '\.sh$'` 為空，或
- Commit message 含 `[skip-go-migration]` 且有 Go migration plan 說明
- 新邏輯在 `hooks.go` 或 subcommand 中實作，`.claude/settings.json` 指向 binary

## Source

- 2026-05-27 session：修正 double-bootstrap 與 Cognitive Mode double-output 兩個 bug 時，直接修改 `.claude/hooks/*.sh`，未考慮 Go 遷移政策。使用者在下一個 turn 指出應查平台治理文件。
- Pre-commit gate 實作：`validateNoNewShellScripts` in `scripts/ai-skill-cli/internal/app/hooks.go`

## Related

- [`plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md`](../../plans/archived/2026-05-21-0834-cross-platform-go-script-runtime.md) — 平台治理 master plan
- [`runtime/cli-modification-policy.yaml`](../../runtime/cli-modification-policy.yaml) — CLI 修改政策
- [`enforcement/failure-patterns/intelligence-layer-bypass-via-tool-adapter.md`](intelligence-layer-bypass-via-tool-adapter.md) — 同類問題（tool adapter 寫入而非 canonical layer）
