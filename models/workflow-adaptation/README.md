# Model Workflow Adaptation

`models/workflow-adaptation/` 定義不同 execution profile 下的 workflow shape。它調整讀取深度、patch scope、validation density 與 handoff 方式；不宣稱 provider model 已切換。

## 入口

- [`small-model-workflows.md`](small-model-workflows.md)：checklist-first 與 bounded context。
- [`large-model-workflows.md`](large-model-workflows.md)：source-backed、exploratory 與 cross-layer reasoning。
- [`coding-workflows.md`](coding-workflows.md)：diff precision、tests、lints、patch scope。
- [`architecture-workflows.md`](architecture-workflows.md)：planning、tradeoff、contradiction analysis。
- [`validation-workflows.md`](validation-workflows.md)：evidence-first、claim scope、no premature success。

## Routing Rule

先由 [`../routing/task-routing.md`](../routing/task-routing.md) 選 strategy，再用 [`../capabilities/README.md`](../capabilities/README.md) 檢查 capability dimensions，最後套用本目錄的 workflow shape。

## 禁止行為

- 不得用 workflow adaptation 取代 workflow primary source。
- 不得因模型能力推測而跳過 source-of-truth gate。
- 不得把 checklist-first 當成低品質 shortcut。
- 不得把 large-profile workflow 當成無限制廣泛探索。
