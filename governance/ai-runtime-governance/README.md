# AI Runtime Governance

`governance/ai-runtime-governance/` 保存「工程哲學 → AI infrastructure governance」的轉譯層。它把 `intelligence/engineering/philosophy/` 中的原始思想，轉成 context、token、activation、replay、automation 等 AI runtime 決策 gate。

本層是治理設計，不直接取代 `enforcement/`。當某條治理要求需要成為可執行 blocking rule，再 promotion 到 `enforcement/`、`runtime/` 或 `validation/`。

## 治理編譯模式

當某個想法一開始只是「人類工程直覺」或「使用者想要的系統思考方式」時，不要直接把它塞進 workflow 或 runtime。先判斷它屬於哪一層：

```text
intelligence source
  -> governance translation
  -> workflow application
  -> runtime / validation enforcement
```

分層規則：

| 層 | 放什麼 | 例子 |
| --- | --- | --- |
| `intelligence/` | 原始思想、判斷智慧、why | Musk Five-Step、index-first documentation thinking、systems thinking |
| `governance/` | AI 化後的可執行治理語意 | context governance、documentation context governance、automation readiness gate |
| `workflow/` | 具體任務中如何套用 | 寫文件時先選 `kind/audience/stability/routing`，再寫 README / leaf docs |
| `runtime/` | 已穩定且可 machine-enforce 的狀態或 guard | TTL、activation rules、recovery state machine |
| `validation/` | 可測的 failure mode | generic skill bloat、automation-before-verification、documentation route drift |

這個模式用來避免 agent 把「哲學」、「治理」和「操作步驟」混成一份文件。若某個 workflow 內開始出現可跨 workflow 重用的治理原則，應優先抽到 `intelligence/` 或 `governance/`，再讓 workflow 引用它。

## 目前條目

| 文件 | Source philosophy | 用途 |
| --- | --- | --- |
| [`five-step-ai-governance.md`](five-step-ai-governance.md) | [`musk-five-step-algorithm.md`](../../intelligence/engineering/philosophy/musk-five-step-algorithm.md) | 新增 skill / memory / workflow / rule / automation 前的 necessity、deletion、simplification、acceleration、automation-last gate。 |
| [`documentation-context-governance.md`](documentation-context-governance.md) | [`index-first-documentation.md`](../../intelligence/engineering/agent-architecture/index-first-documentation.md) | 將 README-as-router、分類欄位、停止條件、單一真相與 leaf expansion 轉成文件 context governance gate。 |
| [`routing-signal-governance.md`](routing-signal-governance.md) | [`task-routing.md`](../../intelligence/engineering/agent-architecture/task-routing.md) | 將 signal strength、primary source、negative signal 與 multi-route disambiguation 轉成 routing decision gate。 |
| [`context-attention-governance.md`](context-attention-governance.md) | [`context-collapse.md`](../../intelligence/engineering/agent-architecture/context-collapse.md), [`attention-budgeting.md`](../../intelligence/engineering/agent-architecture/attention-budgeting.md) | 將 summary-first、attention budget、recap checkpoint 與 task-boundary prune 轉成 context loading gate。 |
| [`validation-scenario-governance.md`](validation-scenario-governance.md) | [`stateless-validation-necessity.md`](../../intelligence/engineering/agent-architecture/stateless-validation-necessity.md), [`failure-to-scenario-closure.md`](../../intelligence/engineering/agent-architecture/failure-to-scenario-closure.md) | 將 stateless reproduction、answer leakage、failure class 與 traceability 轉成 scenario promotion gate。 |

## 放什麼

- AI runtime governance principles。
- Context / token / activation / replay / automation governance。
- Source philosophy 到 workflow/runtime/validation 的 mapping。
- 可被後續 validation scenario 或 runtime gate 使用的治理 criteria。

## 不放什麼

- 原始工程哲學全文；放到 `intelligence/engineering/philosophy/`。
- 當前可執行 enforcement 條文；放到 `enforcement/`。
- 單一 workflow 的操作步驟；放到 `workflow/`。
- runtime machine-readable state；放到 `runtime/`。
