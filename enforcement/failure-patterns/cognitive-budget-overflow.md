# Cognitive Budget Overflow（任務超出 cognitive mode 預算）

Status: candidate
Class: `resource-bound` / `process-gap`

## Trigger

當 agent 執行任務時，**實際 token 用量超過** 宣告的 cognitive mode 組合對應的 `max_tokens` 預算（per [`runtime/cognitive-modes-token-budget.yaml`](../../runtime/cognitive-modes-token-budget.yaml)），使用此 pattern。

具體觸發訊號：

- Commit body 宣告 `Token Estimate: <n>` 但 n 超過 mode tuple 的 `max_tokens`
- Session 內累積 tool call tokens 超過 mode tuple 預算（Phase 5 adaptive runtime 才能偵測）
- Agent 用 FAST mode 但實際展開了 source-backed reads（mode tuple 不一致 → 隱性超 budget）
- 任務沒走 downgrade_path（GRAPH_ASSISTED → SOURCE_BACKED → CHECKLIST_FIRST → SUMMARY_FIRST → INDEX_ONLY）就硬上

## Failure Mode

把高 budget mode（DEEP/FORENSIC）當預設，或忽略 budget 提示直接展開 context，導致：

1. **Context bloat**：session 累積大量無關 context，attention dilution
2. **Cost runaway**：未受控的 token 消耗
3. **Wrong mode**：實際用了 DEEP 行為卻宣告 FAST，governance gate 該開沒開
4. **Audit gap**：commit 沒留 token 估算，事後無法判斷是否合理

## Risk

- **長 session 漂移**：budget 不受控 → context window 滿 → summary triggered → 回到 bootstrap-bypass-on-resume 循環
- **隱性 mode 升級**：宣告 NORMAL 但行為等同 DEEP，governance 落差累積
- **退化路徑遺忘**：應該先 downgrade context_mode 而不是直接放棄任務

## Required Agent Action

任務開始前：

1. 評估 cognitive mode tuple
2. 查 [`runtime/cognitive-modes-token-budget.yaml`](../../runtime/cognitive-modes-token-budget.yaml) 對應 `max_tokens`
3. 在 final report / commit body 宣告 `Token Estimate: <n>` (大致估算)
4. 接近 budget 上限時：
   - **先 downgrade context_mode**（依 `downgrade_path`）
   - 不夠再 split 工作成多 commit
   - 最後才用 `[skip-token-budget]` opt-out（需明確理由）
5. 超 budget 但 commit body 沒宣告 estimate → 視為 declaration miss，補上

## Prevention Gate

任務開始 / commit 前，agent 必須能回答：

| Check | Required answer |
|-------|-----------------|
| 宣告的 mode tuple 是什麼 | 4 維值明確 |
| 對應的 `max_tokens` 是多少 | 查 YAML contract 取得 |
| 預估 token 用量 | 與實際工作量一致 |
| 超 budget 怎麼辦 | 先 downgrade，再 split，最後 opt-out |

對應 runtime enforcement：

- YAML contract `runtime/cognitive-modes-token-budget.yaml`（projected to `generated_surfaces`）
- `validateTokenBudget` in `scripts/ai-skill-cli/internal/app/hooks.go` — commit-msg 偵測宣告超 budget
- Phase 5 adaptive runtime（deferred）才能偵測**實際**用量超 budget

## Validation

符合下列條件時，此 pattern 已被防止：

- 重要任務 commit body 含 `Token Estimate: <n>` 宣告
- 宣告 vs 實際用量差距小（事後 audit）
- 超 budget 時看到 downgrade 紀錄而非直接 skip

## Source

- 2026-05-25 session：完成 Phase 4 時建立此 pattern 作為 cognitive-modes-token-budget.yaml 的 prevention companion。Behavioral 層仍以 commit-msg 為粒度；per-tool-call enforcement 等 Phase 5 adaptive runtime。

## Related

- [`runtime/cognitive-modes-token-budget.yaml`](../../runtime/cognitive-modes-token-budget.yaml) — budget table + 觸發規則
- [`cognitive-mode-resolution-bypass.md`](cognitive-mode-resolution-bypass.md) — mode 沒解析就跑（常與 budget overflow 連鎖）
- [`runtime/cognitive-modes-phase-integration.yaml`](../../runtime/cognitive-modes-phase-integration.yaml) — execution_mode 對 phase_machine 的影響

## Linked Validation Scenarios

- `phase4-token-budget-v1` — YAML contract + projection + validator + failure pattern 存在性

← [Back to failure patterns](README.md)
