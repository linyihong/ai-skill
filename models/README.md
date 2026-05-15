# Models

`models/` 負責「不同模型如何協作」。本層保存 model capability profile、routing strategy、compression strategy 與 prompt adaptation 的工具中立設計。

## 目前入口

- [`profiles/`](profiles/README.md)：定義 `small`、`large`、`specialized` model profiles 與 context loading 深度。
- [`compression/`](compression/README.md)：定義 index-only、summary-first、checklist-first、source-backed、graph-assisted 等壓縮層級。
- [`../knowledge/runtime/model-context-report.md`](../knowledge/runtime/model-context-report.md)：由 routing registry 產生的 model-aware context loading view。
- [`../knowledge/runtime/model-checklists.md`](../knowledge/runtime/model-checklists.md)：由 routing registry 產生的 per-model context-loading checklist。

## 放什麼

- Model capability profile 與適用任務類型。
- Large-model、small-model 與 specialized-model 的 routing strategy。
- Context compression、summary loading 與 checklist-first strategy。
- Prompt adaptation 與 model-aware workflow design。

## 不放什麼

- 特定工具的 model selector UI 或設定路徑；放到 `ai-tools/`。
- Skill workflow 正文；放到 `workflow/` 或仍保留在 `skills/`。
- Metadata schema 欄位定義；放到 `metadata/`。
- 對模型能力的未驗證主張；需先標示 confidence 或留在 TODO。

## 誰會參考這裡（Inbound References）

- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:312) — primary_source 為 `models/profiles/README.md`，required_dependencies 引用 `models/compression/README.md`
- [`knowledge/runtime/model-context-report.md`](../knowledge/runtime/model-context-report.md) — generated view，由 routing registry 產生的 model-aware context loading view
- [`knowledge/runtime/model-checklists.md`](../knowledge/runtime/model-checklists.md) — generated view，由 routing registry 產生的 per-model context-loading checklist

## 與既有層的關係

- `metadata/` 可記錄知識適合哪些 model profile。
- `runtime/` 會使用本層 profile 做 task routing 與 context loading。
- `workflow/` 可引用本層策略，定義大模型與小模型的不同讀取深度。
- `ai-tools/` 保存工具如何實際選用或設定模型。
- `knowledge/runtime/model-context-report.md` 可用來快速檢視目前 route 分別採用的 profile 與 compression level；它是 generated view，不取代 profiles / compression source。
- `knowledge/runtime/model-checklists.md` 可作為 agent 執行前的壓縮 checklist；需要修改或高信心判斷時仍讀回 profiles / compression source。

## 第一批候選遷移來源

- `plans/archived/2026-05-11-next-stage-upgrade-plan.md` 的 Multi-model Runtime Architecture
- `enforcement/decision-efficiency.md` 中與 context cost、compression 相關的 routing 概念
- 未來各 tool adapter 中可抽象成工具中立 model profile 的內容
