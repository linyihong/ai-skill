# Recovery Metadata（恢復政策 metadata）

`metadata/recovery/` 定義 mismatch escalation 後的 domain-specific recovery policy。通用 recovery state machine 由 `runtime/compiler/compiler-rules.yaml` 擁有；本 metadata 層負責告訴 agent：在重建 execution graph 前，不同 domain 必須重新載入哪些 source-of-truth。

Recovery retry 的治理 gate 見 [`governance/ai-runtime-governance/recovery-retry-governance.md`](../../governance/ai-runtime-governance/recovery-retry-governance.md)；cognitive state / evidence qualification 的治理 gate 見 [`governance/ai-runtime-governance/cognitive-state-governance.md`](../../governance/ai-runtime-governance/cognitive-state-governance.md)。本目錄只保存 domain-specific reload set、forbidden behaviors 與 validation metadata；domain evidence authority / freshness / scope policy 放在 [`metadata/evidence/`](../evidence/README.md)。

## 檔案

| 檔案 | 用途 |
| --- | --- |
| [`escalation-levels.yaml`](escalation-levels.yaml) | L1-L5 escalation level metadata、預設 action 與最低 reload 要求。 |
| [`domain-policies.yaml`](domain-policies.yaml) | 各 domain 的 trigger classes、required reload set、forbidden behaviors 與 validation gates。 |

## Policy Schema

每個 domain policy 應包含：

- `domain`：穩定 domain id，例如 `apk-analysis` 或 `software-delivery`。
- `applies_when`：選中此 policy 的 trigger 條件。
- `trigger_classes`：此 policy 覆蓋的 escalation trigger classes。
- `required_reload_set`：必須讀取的 source-of-truth 檔案或來源類別；若不適用，必須明確標記 `not_applicable` / `source_missing`。
- `rebuild_graph`：恢復 execution 前必須具備的 execution graph 欄位。
- `forbidden_behaviors`：進入 recovery 後必須停止的 domain-specific 行為。
- `validation`：證明 recovery 已閉環、可以繼續 execution 的檢查。

## Runtime Boundary

這些檔案在 Phase 4 是 metadata-only，尚未編入 `runtime.db`。使用方式是透過 routing / validation 讀取；可執行的 recovery procedure 仍保留在 `runtime/compiler/compiler-rules.yaml`。

若未來 phase 需要 runtime enforcement，應明確新增 compiler target，不要假設所有 metadata YAML 都會自動編譯。
