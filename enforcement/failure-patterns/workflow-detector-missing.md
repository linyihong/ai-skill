# Workflow Detector Missing（registry 有 activation_triggers，但沒有 runtime detector 讀它）

Status: validated（2026-05-31 session 暴露；fix landed via plan 1900 Phase 3-5）
Class: `meta-governance-gap` / `routing-miss`

> 這是 [`rule-without-executor`](rule-without-executor.md) meta-pattern 的
> **instance #2 的領域特化版**：`knowledge/runtime/routing-registry.yaml`
> 宣告了 `activation_triggers`，但 runtime 沒有 executor 去讀它、把對應
> workflow 設為 active。

## Trigger

具體觸發訊號：

- 使用者描述一個明顯屬於某既有 workflow 的任務（例如「幫我規劃旅遊行程」對應
  `route.workflow.travel-planning`），但 agent 沒有先讀該 workflow 的
  `primary_source` 就直接開工。
- `routing-registry.yaml` 該 route 已有 `activation_triggers.user_signals`，
  訊號詞也明明出現在使用者訊息裡，卻沒有任何 forcing function 讓它生效。
- 使用者反覆追問「為什麼 X workflow 沒有自動啟動 / 你怎麼沒走 workflow？」
- registry 加了新 route + triggers，但 commit 內沒有任何 Go executor
  （`detector.go`）變更 —— triggers 寫了等於沒寫。

## Failure Mode

「Registry 宣告 activation_triggers，但沒有 detector 讀它」→ activation 只是
一張**沒有讀取器的對照表**。具體後果：

1. **Silent routing miss**：使用者任務命中 trigger 詞，但沒有任何機制把
   workflow 設為 active；agent 憑直覺執行，跳過 workflow 的 artifact gates /
   ordering / acceptance。
2. **2026-05-31 實例**：travel-planning route 的 triggers 存在，但 session
   全程沒觸發，agent 直接以一般對話方式回應，未載入 workflow primary_source。
3. **每加一個 route 就累積一份空殼**：新增 triggers 的人以為「加了就會生效」，
   但 runtime 根本沒讀。
4. **與 Discovery 混淆**：誤以為「沒走 workflow」是 Discovery 沒做，實際是
   detector 這層從不存在 —— known route 的 known trigger 本該由 cheap、
   deterministic 的 detector 處理，不該丟給 expensive 的 Discovery。

## Risk

- **Workflow 治理失效**：artifact gates / docs-first / acceptance 等 workflow
  紀律全靠 agent 自律，registry 的 routing 投資白費。
- **信任崩塌**：使用者看到 triggers 寫在 registry 卻不生效，會質疑整個
  routing 層。
- **錯誤歸因**：把 detector gap 誤判成 Discovery gap 或 prompt 問題，
  patch 錯地方。

## Required Agent Action

### 新增 / 修改 route 的 activation_triggers 時

1. 確認 runtime 有 executor 真正讀它：`scripts/ai-skill-cli/internal/app/detector.go`
   `DetectWorkflows`（deterministic any-hit，two-phase，無加權）。
2. 若是 mixed-layer route（analysis / intelligence），確認 `activation_mode`
   已顯式宣告（`must-declare`），否則 detector 用 route_type 預設可能誤分類。
3. 在 `enforcement/enforcement-registry.yaml` 的 `workflow_activation`
   rule_class 確認 binding 與覆蓋狀態一致（detector 已 land → 應反映
   `symbol_exists: ✓`；coverage 在 Phase 7 scenarios + 30 天 runtime 觀察後
   promote 到 `mechanical`）。

### 設計「為什麼沒啟動 workflow」的修法時

- known route + known trigger → 補 / 修 **detector**（cheap, deterministic,
  per-task）。
- unknown capability（registry 根本沒這個 route）→ 走 **Discovery** fallback，
  不是改 detector。兩者分工見
  [`capability-discovery-philosophy.md`](../../governance/lifecycle/capability-discovery-philosophy.md)
  §Discovery → Detector Feedback Loop。

## Fix Landed（plan 1900）

| Phase | 修了什麼 |
|---|---|
| 3 | `detector.go` `DetectWorkflows` — deterministic、two-phase、legacy 正規化、單元測試 |
| 4.0 | `runtime_context.go` `BuildRuntimeContext` — active_route lifecycle（substantive / pivot / manual-lock / NO implicit drift） |
| 5 | `obligation.workflow.activation_evidence` + PreToolUse `workflowPrimarySourceGate` — detector 鎖定 single active_route 後，要求先 Read primary_source（**不誤殺**：miss / conflict / 無 registry 一律 fail-open） |
| 6.1 | Discovery → Detector feedback loop（detector miss → 提案新 route candidate） |

## Cross-References

- [`rule-without-executor`](rule-without-executor.md) — 上層 meta-pattern（本 pattern 是其 instance #2 特化）
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) §`activation_triggers_spec`
- [`scripts/ai-skill-cli/internal/app/detector.go`](../../scripts/ai-skill-cli/internal/app/detector.go) — the executor that closes this gap
- [`governance/workflow-activation-engine.md`](../../governance/workflow-activation-engine.md) — operational spec
- Source plan：[`plans/active/2026-05-31-1900-workflow-activation-engine.md`](../../plans/active/2026-05-31-1900-workflow-activation-engine.md)

← [Back to failure-patterns index](README.md)
