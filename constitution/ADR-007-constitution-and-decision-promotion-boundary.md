# ADR-007: Constitution and Decision Promotion Boundary

## Status

**Accepted**

## Context

原本外層的 `decisions/` 目錄存放正式 ADR，而 `memory/decision/` 存放 session 等級的決策記憶。`decisions/` 這個名稱太廣，容易與 runtime decisions、project decisions 或 memory decisions 混淆。

同時，決策的 promotion 也需要更清楚的 target 規則 — 重複出現或具有持久性的決策**不應**自動升級為 ADR。

## Decision

將正式 ADR 層由 `decisions/` 改名為 `constitution/`。

採用「依內容決定 promotion target」的規則：

| 決策內容 | Promotion target |
| --- | --- |
| 可執行規則或 cross-agent policy | `enforcement/` |
| 推理 heuristic、tradeoff、signal、anti-pattern 或 failure 判斷 | `intelligence/` |
| 操作流程或可重複的 workflow | `workflow/` |
| Runtime gate、activation、phase、obligation、或 executable contract projection | `runtime/runtime.db` |
| 架構級不可逆或基礎性決策 | `constitution/ADR-*` |
| Session 範疇的 replay 決策 | `memory/decision/` |
| 專案專屬的決策 | `<PROJECT_ROOT>/docs/decisions/` |

Runtime 決策記錄的 canonical config 改名為 `runtime/constitution/decision-recording.yaml`，與架構層的命名對齊。

## Consequences

- `constitution/` 是正式 ADR / 架構憲法層。
- `memory/decision/` 仍為 session 等級的 decision replay 層。
- `<PROJECT_ROOT>/docs/decisions/` 仍為專案本地的 decision tier。
- 建立 ADR 不再是每個 promoted decision 的預設終點。
- 影響 runtime 行為的決策必須更新 `runtime.db`，或更新會投影到 `runtime.db` 的 executable YAML contract。

## Alternatives Considered

- 保留 `decisions/`：拒絕，因為它仍會讓正式 ADR 與 memory / runtime 決策混淆。
- 把每個具持久性的決策都 promote 為 ADR：拒絕，因為可執行規則、推理 heuristic、workflow 與 runtime gate 都有更適合的 owner layer。
- 把專案專屬的 decision 資料夾搬到 `constitution/`：拒絕，因為專案本地的決策不屬於 Ai-skill 憲法層。

## Related

- [`README.md`](README.md)
- [`../governance/lifecycle/decision-promotion-pipeline.md`](../governance/lifecycle/decision-promotion-pipeline.md)
- [`../governance/lifecycle/decision-promotion-pipeline.yaml`](../governance/lifecycle/decision-promotion-pipeline.yaml)
- [`../memory/decision/README.md`](../memory/decision/README.md)
- [`../runtime/runtime.db`](../runtime/runtime.db)
