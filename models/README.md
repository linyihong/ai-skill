# Models

`models/` 負責「不同模型如何協作」。本層保存 model capability profile、routing strategy、compression strategy 與 prompt adaptation 的工具中立設計。

## 目前入口

- [`profiles/`](profiles/README.md)：定義 `small`、`large`、`specialized` model profiles 與 context loading 深度。
- [`routing/`](routing/README.md)：定義 task class、cognitive state、tool capability 與 fallback 對應的 execution strategy。
- [`capabilities/`](capabilities/README.md)：把粗略 profile 拆成 reasoning depth、context stability、tool reliability、hallucination risk 與 compression resilience 等能力維度。
- [`workflow-adaptation/`](workflow-adaptation/README.md)：定義 checklist-first、source-backed、coding、architecture 與 validation workflow shapes。
- [`governance/`](governance/README.md)：定義 model selection、hallucination、context budget 與 confidence governance。
- [`runtime/`](runtime/README.md)：定義可供 lookup 使用的 minimal routing primitives，不保存 provider model state。
- [`compression/`](compression/README.md)：定義 index-only、summary-first、checklist-first、source-backed、graph-assisted 等壓縮層級。
- [`../knowledge/runtime/model-context-report.md`](../knowledge/runtime/model-context-report.md)：由 routing registry 產生的 model-aware context loading view。
- [`../knowledge/runtime/model-checklists.md`](../knowledge/runtime/model-checklists.md)：由 routing registry 產生的 per-model context-loading checklist。

## 放什麼

- Model capability profile 與適用任務類型。
- Task class、cognitive state、autonomy mode 與 model capability 對 execution strategy 的影響。
- Large-model、small-model 與 specialized-model 的 context loading defaults。
- Context compression、summary loading 與 checklist-first strategy。
- Prompt adaptation 與 model-aware workflow design。

## 不放什麼

- 特定工具的 model selector UI 或設定路徑；放到 `ai-tools/`。
- 未經工具證實的「實際模型已切換」主張；只能記錄 behavior-only adaptation。
- Skill workflow 正文；放到 `workflow/` 或仍保留在 `skills/`。
- Metadata schema 欄位定義；放到 `metadata/`。
- 對模型能力的未驗證主張；需先標示 confidence 或留在 TODO。

## 誰會參考這裡（Inbound References）

- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:312) — primary_source 為 `models/profiles/README.md`，required_dependencies 引用 `models/compression/README.md`
- [`knowledge/runtime/model-context-report.md`](../knowledge/runtime/model-context-report.md) — generated view，由 routing registry 產生的 model-aware context loading view
- [`knowledge/runtime/model-checklists.md`](../knowledge/runtime/model-checklists.md) — generated view，由 routing registry 產生的 per-model context-loading checklist

## 與既有層的關係

- `metadata/` 可記錄知識適合哪些 model profile。
- `runtime/` 只能在 validation 後接收 minimal routing primitives；初期 model-aware routing 保持 design-layer-only。
- `workflow/` 可引用本層策略，定義大模型與小模型的不同讀取深度。
- `ai-tools/` 保存工具如何實際選用或設定模型。
- `knowledge/runtime/model-context-report.md` 可用來快速檢視目前 route 分別採用的 profile 與 compression level；它是 generated view，不取代 profiles / compression source。
- `knowledge/runtime/model-checklists.md` 可作為 agent 執行前的壓縮 checklist；需要修改或高信心判斷時仍讀回 profiles / compression source。

## Phase 0-2 Contract

Model-aware routing 目前是 execution strategy contract：

1. 先用 [`routing/task-routing.md`](routing/task-routing.md) 判斷 task class。
2. 再用 [`routing/autonomy-routing.md`](routing/autonomy-routing.md) 套用 cognitive state / autonomy mode。
3. 用 [`capabilities/README.md`](capabilities/README.md) 把粗略 profile 轉成可驗證能力維度。
4. 若工具不能實際選 model，依 [`routing/fallback-routing.md`](routing/fallback-routing.md) 使用 behavior-only adaptation。
5. 只有 tool adapter 證實可指定模型時，才依 [`routing/multi-model-handoff.md`](routing/multi-model-handoff.md) 進行 explicit model / subagent handoff。

## Phase 3-6 Surfaces

後續 workflow / governance / runtime / validation 接入遵守：

- Workflow adaptation 只調整 execution shape，不取代 workflow primary source。
- Governance 禁止 unsupported model switch claims 與 silent substitution。
- Runtime surface 只保存 minimal strategy primitives，不保存 provider state。
- Validation scenarios 必須覆蓋 unavailable model、explicit model request、uncertain / contaminated state、small-model workflow、architecture workflow 與 source-of-truth override prevention。

## 第一批候選遷移來源

- `plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md` 的 Multi-model Runtime Architecture
- `enforcement/decision-efficiency.md` 中與 context cost、compression 相關的 routing 概念
- 未來各 tool adapter 中可抽象成工具中立 model profile 的內容
