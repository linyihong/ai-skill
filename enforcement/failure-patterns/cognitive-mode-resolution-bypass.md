# Cognitive Mode Resolution Bypass（跳過 cognitive mode 解析直接執行）

Status: candidate
Class: `process-gap` / `governance-drift`

## Trigger

當 agent 開始執行任務，卻沒有依 discovery signals 評估 cognitive mode（execution/context/governance/memory 四維），或 `cognitive_modes` 表中沒有本 session 的 mode resolution row，使用此 pattern。

具體觸發訊號：

- 收到任務後直接進入 execution，沒有讀 discovery signals
- 執行前未確認 `execution_mode`（FAST/NORMAL/DEEP/FORENSIC/RECOVERY）
- `context_mode` 未設定而使用預設壓縮行為（不等於 DEFAULT fallback 已套用）
- `governance_mode` 未設定而跳過 gate check
- `memory_mode = NONE` 被忽略，仍讀了 `memory/` subdir
- Chat / session final response 沒有 compact `Cognitive:` 或 full `### Cognitive Mode 報告`，或 final-response hook 使用錯誤 blocking protocol 導致缺報仍被顯示，代表 per-turn resolution 沒有 user-facing evidence

## Failure Mode

在沒有明確 mode 解析的情況下執行，導致：

1. **壓縮策略錯誤**：複雜 architecture task 用 FAST mode 跳過 source-backed reading
2. **governance gate 漏掉**：STRICT/LOCKDOWN 模式該啟用的 gates 未激活
3. **memory 誤讀**：`memory_mode = NONE` 應禁止讀 memory，卻仍進行 replay
4. **phase_machine floor 被忽略**：DEEP mode 要求 `source_backed_before_decision`，FAST 的 `allowed_actions` floor 不一樣
5. **後驗難以追蹤**：`cognitive_modes` 表沒有 row，無法事後 audit 哪個 task 用了哪個 mode
6. **Chat close-out 漏報**：commit hook 之外的 session final response 沒有 Cognitive evidence，使用者無法知道本輪實際套用的 mode / validation depth

## Risk

- **Silent default contamination**：tool 用錯誤 compression level 載入 context，agent 不知道
- **Gate bypass accumulation**：多次無 mode 執行後，blocking_gates 的 STRICT/LOCKDOWN category 從未激活
- **Memory isolation failure**：memory_mode NONE 的場景（contaminated state）被跳過，舊 memory 污染當前判斷
- **Commit-only blind spot**：只在 commit-msg 驗證 Cognitive，會漏掉 read-only、diagnostic、aborted、或 session final 的 user-facing response

## Required Agent Action

每個任務開始前：

1. **讀 discovery signals**（`runtime/cognitive-modes-discovery.yaml` 或 `discovery_signals` 表）
2. **評估 task context**，選出 execution_mode / context_mode / governance_mode / memory_mode
3. **記錄 resolution**（或確認 `cognitive_modes` 表已有本 task row）
4. 若信號不足，套用 **fallback**：NORMAL / SUMMARY_FIRST / STANDARD / NONE
5. 依 governance_mode 啟用對應 gate_set（LIGHT/STANDARD/STRICT/LOCKDOWN）
6. 若 memory_mode = NONE，跳過所有 memory subdir
7. 最終 user-facing response 必須附 compact 或 full Cognitive Mode 報告；工具若支援 stop/final-response hook，必須用 hook 檢查這個 close-out

## Prevention Gate

開始執行前，agent 必須能回答：

| Check | Required answer |
|-------|-----------------|
| 是否已評估 discovery signals | 已評估，或明確套用 fallback |
| execution_mode 是什麼 | 明確值（FAST/NORMAL/DEEP/FORENSIC/RECOVERY）|
| context_mode 是什麼 | 明確值（FULL/SUMMARY_FIRST/INDEX_ONLY/CHECKLIST_ONLY/GRAPH_ONLY）|
| governance_mode 對應哪個 gate_set | 已確認 gate_set_activation |
| memory_mode 是否為 NONE | 若是，已確認不讀任何 memory subdir |
| final response 是否有 Cognitive evidence | 已附 compact/full block，或 stop hook 會 block/loop back |

若回答不確定，先讀 [`runtime/cognitive-modes.yaml`](../../runtime/cognitive-modes.yaml) + [`runtime/cognitive-modes-discovery.yaml`](../../runtime/cognitive-modes-discovery.yaml)。

## Validation

符合下列條件時，此 pattern 已被防止：

- 每個任務執行前有可追蹤的 mode resolution（`cognitive_modes` 表有 row 或 fallback 明確記錄）
- governance_mode 對應的 gate_set 已激活
- memory_mode = NONE 時，沒有讀 `memory/` 任何 subdir
- Cursor / Claude 等支援 stop/final-response hook 的 adapter 有 regression 覆蓋「final response 缺 Cognitive 時會 block 或 loop back」；Cursor stop 必須驗證 `followup_message` protocol

## Source

- 2026-05-22~25 session：Phase 3 設計 governance-structure enforcement 時，確認「mode resolution bypass」是最高風險場景，因為現有 gate 只檢查 mode 是否 resolved，若 agent 完全跳過解析則 gate 無法作用。Phase 3-B 已加入 `obligation.execution.resolve_cognitive_mode` + `gate.execution.cognitive_mode_resolved`，但 behavioral enforcement（Go 代碼層）已延後到 Phase 4。
- Deferred: Phase 4 behavioral enforcement — pre-commit hook Go 代碼檢查 `cognitive_modes` 表是否有 current-session row。

## Related

- [`runtime/cognitive-modes.yaml`](../../runtime/cognitive-modes.yaml) — 4 維 mode primitives + defaults + gate_activation
- [`runtime/cognitive-modes-discovery.yaml`](../../runtime/cognitive-modes-discovery.yaml) — 14 discovery signals，決定 mode 不靠文件查詢
- [`runtime/cognitive-modes-cost-class.yaml`](../../runtime/cognitive-modes-cost-class.yaml) — v2 cognitive_cost derivation，避免 agent self-claim cost
- [`runtime/cognitive-modes-governance-integration.yaml`](../../runtime/cognitive-modes-governance-integration.yaml) — governance_mode → gate_set_activation
- [`runtime/cognitive-modes-memory-integration.yaml`](../../runtime/cognitive-modes-memory-integration.yaml) — memory_mode → subdir + AND compose
- [`memory/retrieval-governance/activation-thresholds.md`](../../memory/retrieval-governance/activation-thresholds.md) — memory activation levels + memory_mode compose
- [`inflated-cognitive-mode-reporting.md`](inflated-cognitive-mode-reporting.md) — v2 防止 mode / cost / activation_reason 被 agent 自我膨脹

## Linked Validation Scenarios

- `cognitive-modes-enforcement-gate-exists-v1` — 驗證 `gate.execution.cognitive_mode_resolved` 在 `gates` 表存在
- `phase6-cognitive-contract-v2-inflated-rejection-v1` — 驗證 inflated reporting 被 cost / signal validators 擋下
- `cursor-stop-hook-final-cognitive-required-v1` — 驗證 Cursor-style stop payload 缺 Cognitive 時輸出 `followup_message` loop back

← [Back to failure patterns](README.md)
