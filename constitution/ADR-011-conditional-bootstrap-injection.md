# ADR-011: UserPromptSubmit Hook Conditional Bootstrap Injection

## Status

**Accepted**

## Framework Generation

- **世代分類**：Gen 3 Layer 2.5 (Meta Governance) — `bootstrap_integrity` rule_class
  enforcement-strength 調整（非世代升級）
- **當前世代文件**：[`enforcement/enforcement-registry.yaml`](../enforcement/enforcement-registry.yaml)
  `rule_classes[id=bootstrap_integrity]`；[`scripts/ai-skill-cli/internal/app/hooks.go`](../scripts/ai-skill-cli/internal/app/hooks.go)
  `runUserPromptSubmitHook`
- **適用狀態**：本 ADR 將 `bootstrap_integrity` 從 `coverage: mechanical` demote
  為 `coverage: behavioral_only`，並引用本 ADR 作為 R2 demotion rationale。
  PreToolUse mechanical gate（`gate.bootstrap.receipt_present`）不變。

## Date

2026-06-03

## Source Plan

[`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](../plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md)
§Phase 5.x — Hook Injection Economics（inaugural self-governance test case）

## Context

2026-05-31 session 揭露 `runUserPromptSubmitHook` 每輪都把整份 `CORE_BOOTSTRAP.md`
注入為 `additionalContext`，並把 cognitive-mode close-out reminder 與
「If bootstrap has not yet been acknowledged」conditional 黏在同一行字串。實測造成兩個問題：

1. **Token 浪費**：CORE_BOOTSTRAP.md 約 2–3K tokens，乘以 N turns 後 cache miss
   + window 壓縮成本明顯（同 session 內，bootstrap.md 不會改變）。
2. **Conditional 由 agent 解讀，不可機械驗證**：hook 字串說「if not yet acknowledged」，
   但 hook 自身沒 scan transcript，只是把 conditional 推給 agent。agent 在 resume /
   compaction 後容易誤判而 over-emit Bootstrap Receipt（本 session 自己就是這個失效模式）。

D1-D4 詳細 defect 表見 plan §Phase 5.x。

## Decision

採用 **transcript-aware conditional injection**：

1. `runUserPromptSubmitHook` 從 stdin 讀 Claude Code payload 取 `transcript_path`。
2. 新增 helper `transcriptHasBootstrapAcknowledgment(transcriptPath, lastN)`：
   scan JSONL transcript 中最近 N 個 assistant text turn，比對
   `Bootstrap: rules=✓ phase=… obligations=… gates=…` 的 Receipt signature。
   結構複用 `transcriptHasRequiredBootstrapReads`（差別在比對 text block 而非 tool_use Read）。
3. Hook 注入內容拆兩個獨立 block：
   - **Block 1（永遠注入）**：Cognitive Mode 報告 close-out MUST。
   - **Block 2（條件注入）**：transcript 已觀察到 Receipt → 省略；否則注入完整
     CORE_BOOTSTRAP.md + 「Receipt not yet observed; emit at start of next response」prompt。
4. Git report block 邏輯不變。

`bootstrap_integrity` rule_class coverage 從 `mechanical` demote 為 `behavioral_only`：
hook 注入不再是「always full-context forcing function」；真正機械強制留給 PreToolUse
`gate.bootstrap.receipt_present`（仍然 block agent 工具呼叫直到 Receipt + canonical Read events
都在 transcript 出現）。Registry 標 `behavioral_only` 反映「主要 forcing function 已變條件式」，
而非「完全失去機械保證」。Sunset 條件：當 hook injection 經 30 天 runtime_metrics 觀察
證實 Receipt over-emit 率明顯下降，且 PreToolUse gate 持續攔截偽 Receipt，未來可考慮
re-promote 為 mechanical。

## Consequences

### 正向

- Token 節省：subsequent turn 不再注入 ~2–3K tokens 的 bootstrap.md。
- Conditional 機械化：hook 自己 scan transcript 判定 acknowledgment，agent 不再做
  conditional 解讀，減少 over-emit Receipt 的失效模式。
- MUST / conditional 分塊：cognitive-mode close-out 永遠在第一段，不會被 bootstrap
  block 蓋過。
- 是 **registry inaugural self-governance test case**：本 ADR + registry demotion +
  hooks.go fix 在同一 commit landing，會觸發 Phase 4.5 R1 + R2 validator
  (`validateEnforcementRegistryTransition`) ；commit message 帶
  `[registry-status-change]` trailer + `rationale:` line + ADR 解析成功 → R2 PASS →
  自我驗證 Phase 4.5 validator 在真實 commit path 正常 fire。

### 負向 / 風險

- PreToolUse gate 仍是真正 enforcement；若該 gate 失效，這次 demotion 等於 weakening
  整個 bootstrap chain。緩解：PreToolUse `gate.bootstrap.receipt_present` 有獨立
  regression test (`TestPreToolUseHookBlocksReceiptWithoutReads`)，change set 不動 PreToolUse 路徑。
- `transcriptHasBootstrapAcknowledgment` 使用正則匹配 Receipt 文字格式；若 Receipt format 改
  (e.g. 變 emoji 或欄位順序)，helper 會 false-negative 而退回 always-inject 模式。緩解：
  format 由 `runtime/core-bootstrap.yaml` per_session_obligations canonical 定義，這次
  helper 用最寬鬆 regex (`Bootstrap:\s*rules=✓\s*phase=...`) 保留彈性。
- Stdin 讀取在無 payload 時不能 hang。實作上 tolerant parse：empty / malformed →
  視為 transcript_path 未知 → 退回 always-inject（安全側）。

## Alternatives Considered

1. **完全不注入 bootstrap.md，只靠 SessionStart hook**：放棄 mid-session repair 能力。
   resume / compaction 後若 agent 沒帶 bootstrap context，PreToolUse 會 block 但沒有
   reminder 告訴 agent 該讀什麼。捨棄。
2. **依 turn count 而非 transcript scan**：依 turn N≥2 跳過。同樣會在 resume / 多 agent
   handoff 情境失準。捨棄。
3. **Coverage 維持 `mechanical`，只改 hooks.go**：不誠實反映「主要 forcing function 變條件式」
   的事實，且不會觸發 R2 → 失去 inaugural self-governance test 機會。捨棄。

## Related

- [`enforcement/enforcement-registry.yaml`](../enforcement/enforcement-registry.yaml) §self_governance R1/R2
- [`enforcement/enforcement-registry.md`](../enforcement/enforcement-registry.md) §Status Transition Matrix
- [`enforcement/failure-patterns/bootstrap-bypass-on-resume.md`](../enforcement/failure-patterns/bootstrap-bypass-on-resume.md)
- ADR-010 `upstream_classes` scope freeze（同 registry 治理脈絡）
