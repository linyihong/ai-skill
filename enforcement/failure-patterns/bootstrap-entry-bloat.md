# Bootstrap Entry Bloat（tool entry 把 canonical 內容複製成第二份）

Status: validated
Class: `source-of-truth-duplication` / `governance-drift`

## Trigger

當 agent 修改 AI tool 的 session-start entry file（repo-root `CLAUDE.md`、`.cursor/rules/ai-skill-bootstrap.mdc`、`.roomodes` 等），把 canonical source 的 obligation / format / enum 內容**複製進去**而非用**參照**，使用此 pattern。

具體觸發訊號：

- 修改 entry file 後檔案行數 > 30
- Entry file 含 mode enum 字串（如 `FAST/NORMAL/DEEP/FORENSIC/RECOVERY`）
- Entry file 含完整 markdown 表格（如 `| 維度 | 值 | 理由 |`）
- Entry file 含 Bootstrap Receipt 完整 example（如 `Bootstrap: rules=✓ phase=phase.bootstrap obligations=`）
- 同樣 obligation format 同時出現在 entry file 與 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)

## Failure Mode

把 entry file 當「便利的 cheat sheet」貼上 obligation format，導致：

1. **Duplication drift**：CORE_BOOTSTRAP.md 更新 obligation format 時，CLAUDE.md / .cursor / .roomodes 不會自動同步 → 不同 tool 看到不同版本
2. **Multi-tool blast radius**：[`init_project.go`](../../scripts/ai-skill-cli/internal/app/init_project.go) 把 bloated CLAUDE.md template 散播到下游所有專案的 Roo / Cursor / Claude entries → 每改一條 obligation 要改 3+ 處
3. **Source-of-truth ambiguity**：agent 不知該以 entry file 或 CORE_BOOTSTRAP.md 為準
4. **Tool-add 障礙**：新增工具（例如 Aider）時，又要複製一份 entry — 重複工作 + 重新出包

## Risk

- **每個 tool 都「以為自己有 canonical 規則」**：Roo Code 用戶 fork repo 後改 .roomodes 想加 obligation，永遠不會去動 CORE_BOOTSTRAP.md
- **Obligation 漂移**：兩處內容不一致時，agent 隨選一個版本執行
- **Fix A regression 證據**：2026-05-25 commit `8970d4d` 把 Cognitive Mode block 加到 CLAUDE.md 同時 CORE_BOOTSTRAP.md — duplicate；下次改其中一處就出包

## Required Agent Action

修改任何 entry file 時：

1. 判斷你要改的內容性質：
   - **Cross-tool obligation**（任何 tool 都該遵守）→ 改 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)
   - **Tool-specific note**（只該 tool 在意）→ 改 [`ai-tools/agent/<tool>.md`](../../ai-tools/agent/)
   - **Mode value / contract format** → 改 [`runtime/cognitive-modes-*.yaml`](../../runtime/) 或 [`models/cognitive-modes/`](../../models/cognitive-modes/)
2. **不要**把 1 的內容複製進 entry file
3. Entry file 只保留：reference link + IMPORTANT framing + 「請依 canonical source 執行」
4. Entry file 應 ≤ 30 行
5. 若 init_project.go 要產出下游 entry，同樣保持 thin

## Prevention Gate

修改 entry file 前 agent 必須能回答：

| Check | Required answer |
|---|---|
| 要改的內容是 cross-tool obligation 嗎 | 是 → 改 CORE_BOOTSTRAP.md 而非 entry file |
| 要改的內容是 tool-specific note 嗎 | 是 → 改 ai-tools/agent/<tool>.md |
| Entry file 改完 ≤ 30 行嗎 | 是 |
| Entry file 含 mode enum / format example 嗎 | 不應有 |
| Reference link 還在嗎 | 必須在 |

對應 runtime enforcement：

- YAML contract [`runtime/bootstrap-entry-points.yaml`](../../runtime/bootstrap-entry-points.yaml)（projected to `generated_surfaces`）
- `validateBootstrapEntryThinness` in [`scripts/ai-skill-cli/internal/app/hooks.go`](../../scripts/ai-skill-cli/internal/app/hooks.go) — commit-msg 階段偵測 staged entry file + thinness 違反

## Validation

符合下列條件時，此 pattern 已被防止：

- 修改 entry file 的 commit 通過 `validateBootstrapEntryThinness`
- Entry files 行數 ≤ 30
- Entry files 不含 mode enum / format example
- Cross-tool obligation 只存在於 CORE_BOOTSTRAP.md 一處

## Source

- 2026-05-25 session：Fix A 把 Cognitive Mode block obligation 加到 CLAUDE.md（重複 CORE_BOOTSTRAP.md 既有內容）。User 立即指出：「不要都把各個東西特殊化的設定(claude.md)都設定在那裡 不然其他工具用了 還是會有這個問題 ... 只希望要薄薄的設定入口就好 ... 改成強制 或是使用yaml」。
- 同 session 補 Fix B（thinness rule + YAML contract + Go validator）。

## Related

- [`runtime/bootstrap-entry-points.yaml`](../../runtime/bootstrap-entry-points.yaml) — 規則 canonical source
- [`source-mirror-write-drift.md`](source-mirror-write-drift.md) — 同類 source-of-truth-duplication
- [`framework-duplication-without-interrogation.md`](framework-duplication-without-interrogation.md) — 同類 duplication failure
- [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) — cross-tool obligation canonical
- [`ai-tools/agent/`](../../ai-tools/agent/) — tool-specific notes canonical

## Linked Validation Scenarios

- `bootstrap-entry-thinness-v1` — YAML + failure pattern + validator + thin CLAUDE.md 存在

← [Back to failure patterns](README.md)
