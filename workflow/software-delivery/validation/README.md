# Software Delivery Validation Workflow

Validation stage 負責 proof acquisition：behavior correctness、business invariant correctness、execution correctness 分別需要對應 evidence。

## Test Ordering（test-first vs test-after）

| 變更類型 | 順序強制等級 | 出處 |
|---------|------------|------|
| Framework / runtime / governance / workflow / validation / scenario / metadata / compiler 改動 | **強制 test-first**（scenarios commit 早於實作 commit）| [`governance/lifecycle/system-upgrade-governance.md`](../../../governance/lifecycle/system-upgrade-governance.md) §3 規則 9 |
| 一般 production code（非 framework / runtime / governance）| 建議 test-first（test-driven design feedback）| [`intelligence/engineering/heuristics/test-driven-heuristic.md`](../../../intelligence/engineering/heuristics/test-driven-heuristic.md) |
| Doc-only trial / bug fix / typo / spike | 豁免（須明寫理由）| 規則 9 豁免清單 |

詳細原則：

- 通用 TDD（測試難寫 → 設計回饋）→ [`intelligence/engineering/heuristics/test-driven-heuristic.md`](../../../intelligence/engineering/heuristics/test-driven-heuristic.md)
- Framework 升級的順序強制（scenarios 必須先於實作）→ [`intelligence/engineering/development/test-first-framework-upgrade.md`](../../../intelligence/engineering/development/test-first-framework-upgrade.md)

兩者互補：通用 TDD 處理「測試難不難寫」（設計問題），test-first-framework-upgrade 處理「測試何時寫」（順序問題）。

## 進入本 stage 前

[`execution-flow.md`](../execution-flow.md) §4 測試策略定義 已確認測試策略；§7 驗證 執行 proof acquisition；§Test-First Ordering callout 確認 framework 升級已遵循順序。

## Evidence Types（gate vocabulary）

Gate `requires:` 只允許 **evidence_type** token（例如 `evidence:user_visible`）。`collection_method`（如 `browser_observation`）與 `artifact_shape`（如 `screenshot`）屬於 integration envelope 欄位，**不得**進 `requires:`。

- Catalog: [`validation/evidence-types/README.md`](../../../validation/evidence-types/README.md)
- Trace chain: **gate → claim → artifact**（禁止只有 artifact 無 claim）
- `browser_review` = activity summary only；不是 pass/fail token
- OQ-5: **reject token inheritance** — 見 catalog README

- Gate vocabulary: [`evidence-gate-vocabulary.md`](evidence-gate-vocabulary.md)
- Authority routing: [`authority-decision-table.md`](authority-decision-table.md)
- Failure catalog: [`failure-evolution-catalog.md`](failure-evolution-catalog.md)
- User-visible KPI depth gate: [`../validation.md`](../validation.md) § User-Visible Counter Depth Gate
- Experience runtime template: [`../../cross-cutting/experience-runtime/player.yaml`](../../cross-cutting/experience-runtime/player.yaml)

## Live Evidence Chain Guide

當變更命中 [`state-visibility-gap.md`](../../../intelligence/engineering/execution/validation-reasoning/state-visibility-gap.md) 時，validation stage 必須取得足以支持 claim 的 evidence chain，而不只是一個局部測試結果。

建議步驟：

1. 列出 state source：identity、entitlement、tenant、ownership、feature flag、payment state 或其他來源。
2. 列出 propagation chain：API、domain logic、DB、queue、external system、SSR/API readback、UI 或 business observable result。
3. 依 [`evidence-model.md`](../../../intelligence/engineering/execution/validation-reasoning/evidence-model.md) 標記每個證據的 confidence 與 scope。
4. 依 [`evidence-depth.md`](../../../intelligence/engineering/execution/validation-reasoning/evidence-depth.md) 確認最低 depth；高風險用 live system，critical path 加 independent observation。
5. 對 identity-coupled side effect flow，使用產品正式 API 或 UI/H5 path 取得真實身份材料，帶該身份呼叫 SSR/API readback，並驗證 DB/state、side effect 與 user-observable state 一致。
6. 對 email、payment、queue、storage、external API 等 proxy-prone path，不接受 adapter success 作為 final proof；需有 provider/consumer/readback/independent confirmation。

本機 email、SMTP、帳密、收件人、cookies、tokens 或測試環境細節必須留在 gitignored local env 或專案安全文件，不得寫入 reusable Ai-skill 文件。
