# Refactor Parity Feedback Miss（重構對照回饋漏記）

Status: candidate
Class: `validation-gap` / `dependency-miss`

## Trigger

當使用者指出「這個重構是否應該回饋到 software-delivery」、「為什麼沒有把新舊功能對照表沉澱成流程」、「replacement 會不會漏掉舊能力」等語意，且 agent 剛剛只在單一計畫或文件補了 parity inventory，卻沒有把通用 gate 回饋到 workflow / governance 時，使用此 pattern。

## Failure Mode

Agent 把 parity inventory 當成單一專案的補文件工作，只修當前計畫，沒有識別出這是可重用的 software-delivery 規則：任何重構、遷移、改寫或 replacement 若要替代既有功能，都需要先盤點舊入口、現有能力、副作用、外部依賴、新入口、parity 狀態與測試證據。

## Risk

- 新功能看似規劃完整，但沒有可驗證方式證明舊能力已被覆蓋。
- 後續 replacement 類任務重複漏掉舊功能、flags、side effects、hooks、generated surfaces 或環境依賴。
- 使用者只能用人工記憶檢查遺漏，而 workflow 沒有強制 agent 在 implementation 前建立對照表。
- 可重用 lesson 停留在單一計畫文件，下一個 software-delivery 任務不會讀到。

## Required Agent Action

1. 將目前工作重新分類：若新入口取代舊入口，這不是普通純 refactor，而是 replacement / migration parity 問題。
2. 讀取 `workflow/software-delivery/execution-flow.md`、`artifact-gates.md`、`governance/ai-runtime-governance/software-delivery-governance.md`、`enforcement/failure-learning-system.md` 與 `enforcement/linked-updates.md`。
3. 在 software-delivery workflow 補上 replacement parity gate，而不是只補當前計畫。
4. 若 failure 來自 agent 未主動 feedback reusable lesson，新增或更新 `enforcement/failure-patterns/`。
5. 回到原計畫確認它已連到新的 workflow gate，並保留計畫自己的 project-specific inventory。
6. 完成 diff review、語言檢查、runtime compiler / validator、commit、push、readback 與 clean status。

## Prevention Gate

當任務包含「重構、遷移、替換、重寫、跨平台改寫、新 CLI 取代舊 script、新 API 取代舊 API、新資料流取代舊資料流」任一訊號時，agent 在實作前必須能回答：

| Check | Required answer |
| --- | --- |
| Replacement scope | 哪些舊入口或舊能力會被新入口取代？ |
| Parity inventory | 是否已有新舊能力對照表，包含輸入、輸出、副作用、外部依賴與新入口？ |
| Deferred handling | 未覆蓋、延後或 tool-specific 的能力為何不阻擋本 phase？ |
| Test evidence | 每個高風險舊能力是否有 BDD、fixture、contract test、golden output 或人工審查 gate？ |
| Reusable feedback | 若這是 workflow-level 缺口，是否已回饋到 `workflow/software-delivery/` 或治理 gate？ |

## Validation

此 pattern 已套用時，應可反查：

- `workflow/software-delivery/execution-flow.md` 有 replacement / refactor parity gate。
- `governance/ai-runtime-governance/software-delivery-governance.md` 有對應 runtime gate 或 validation candidate。
- 當前 project plan 或 implementation docs 仍保留具體 parity inventory。
- `enforcement/failure-patterns/README.md` 已索引本 pattern。
- Canonical `<AI_SKILL_REPO>` 完成 commit、push、readback；reference-first 時 tool sync 標為不適用。

## Linked Rules

- [`../failure-learning-system.md`](../failure-learning-system.md)
- [`../linked-updates.md`](../linked-updates.md)
- [`../dependency-reading.md`](../dependency-reading.md)
- [`../../workflow/software-delivery/execution-flow.md`](../../workflow/software-delivery/execution-flow.md)
- [`../../workflow/software-delivery/artifact-gates.md`](../../workflow/software-delivery/artifact-gates.md)
- [`../../governance/ai-runtime-governance/software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md)

## Linked Validation Scenarios

- `refactor-parity-feedback-miss-v1` — 檢查 replacement / refactor parity lesson 是否從單一計畫回饋到 software-delivery workflow 與 governance gate
- `validate_failure_pattern_validator_coverage` — 檢查每個 failure pattern 的 Linked Validation Scenarios 是否為空

← [Back to failure patterns](README.md)
