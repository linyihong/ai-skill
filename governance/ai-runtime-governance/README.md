# AI Runtime Governance

`governance/ai-runtime-governance/` 保存「工程哲學 → AI infrastructure governance」的轉譯層。它把 `intelligence/engineering/philosophy/` 中的原始思想，轉成 context、token、activation、replay、automation 等 AI runtime 決策 gate。

本層是治理設計，不直接取代 `enforcement/`。當某條治理要求需要成為可執行 blocking rule，再 promotion 到 `enforcement/`、`runtime/` 或 `validation/`。

## 目前條目

| 文件 | Source philosophy | 用途 |
| --- | --- | --- |
| [`five-step-ai-governance.md`](five-step-ai-governance.md) | [`musk-five-step-algorithm.md`](../../intelligence/engineering/philosophy/musk-five-step-algorithm.md) | 新增 skill / memory / workflow / rule / automation 前的 necessity、deletion、simplification、acceleration、automation-last gate。 |

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
