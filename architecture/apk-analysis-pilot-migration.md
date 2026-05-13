# APK Analysis Pilot Migration Map

本文件定義 `apk-analysis` 作為下一階段 Workflow / Analysis / Intelligence 分離的 pilot。目標是先建立 reference-first 目的地與 mapping，不搬移大量內容、不破壞既有 `skills/apk-analysis/` 入口。

## 遷移狀態

| 欄位 | 值 |
| --- | --- |
| Pilot skill | `skills/apk-analysis/` |
| 狀態 | `techniques-deleted` |
| 舊入口 | `skills/apk-analysis/SKILL.md` 仍為 active skill entrypoint |
| 新參考路徑 | `analysis/apk/`, `workflow/apk-analysis/`, `intelligence/engineering/apk-analysis/` |
| 大量搬移 | Phase C 已完成 — `skills/apk-analysis/techniques/` 和 `analysis/apk/techniques/` 已刪除。新檔案為 canonical source。 |
| 已提取檔案 | `workflow/apk-analysis/execution-flow.md`, `workflow/apk-analysis/artifact-gates.md`, `analysis/apk/traffic-triage.md`, `analysis/apk/tools-and-failures.md`, `analysis/apk/workflows/frida-hook-flow.md`, `analysis/apk/workflows/http-api-documentation-flow.md`, `analysis/apk/workflows/local-proxy-hook-flow.md`, `analysis/apk/workflows/media-hls-analysis-flow.md`, `intelligence/engineering/apk-analysis/evidence-first-routing.md`, `intelligence/engineering/apk-analysis/live-readiness-gates.md`, `intelligence/engineering/apk-analysis/heuristics/hook-selection.md`, `intelligence/engineering/apk-analysis/heuristics/api-documentation-completeness.md`, `intelligence/engineering/apk-analysis/heuristics/local-proxy-routing-diagnosis.md`, `intelligence/engineering/apk-analysis/anti-patterns/early-hook-instability.md`, `intelligence/engineering/apk-analysis/failure/frida-spawn-race.md`, `intelligence/engineering/apk-analysis/signals/flutter-dart-aot-detection.md`, `intelligence/engineering/apk-analysis/signals/local-proxy-detection.md`, `intelligence/engineering/apk-analysis/signals/media-type-detection.md` |
| 索引已更新 | `knowledge/indexes/README.md` — 新增 4 條 routing entries |

## 遷移狀態

- `skills/apk-analysis/SKILL.md` 已不再作為 active entrypoint。新內容請直接寫入新分層路徑。
- 新 top-level layer 路徑已承接對應內容，並通過 validation。
- 既有 `skills/apk-analysis/` 檔案保留作為歷史參考，不應作為新內容的寫入目標。
- 特定 APK hosts、endpoints、raw responses、tokens、device identifiers 與 private run evidence 留在 project docs。

## 來源到目標對應表

| 既有來源 | 目前角色 | 候選目標 | 遷移動作 | 狀態 |
| --- | --- | --- | --- | --- |
| `skills/apk-analysis/SKILL.md` | Tool skill trigger、authorization boundary、output style、cross-skill handoff | `workflow/apk-analysis/` | 已遷移至新分層，skills/ 不再作為 active entrypoint | ✅ migrated |
| `skills/apk-analysis/WORKFLOW.md` | Evidence-first traffic/runtime decision tree | `analysis/apk/` + `workflow/apk-analysis/` | 將 observation / triage methods 拆分到 analysis；execution sequencing 保留在 workflow | candidate |
| `skills/apk-analysis/TOOLS.md` | 工具選擇、命令模板、失敗判讀 | `analysis/apk/` + `workflow/apk-analysis/` | 將 tool-selection reasoning 移到 analysis；setup steps 保留為 workflow references 或 tool docs | candidate |
| `skills/apk-analysis/DOCUMENTATION.md` | 專案 artifact templates 與 documentation gates | `workflow/apk-analysis/` + `intelligence/engineering/apk-analysis/` | Artifact production flow 保留在 workflow；stable engineering lessons 提取到 intelligence | candidate |
| `skills/apk-analysis/techniques/` | Route-specific analysis methods | `analysis/apk/techniques/` | 已提取 4 個 categories（flutter-dart-aot、http-api、local-proxy、media-hls）到 `analysis/apk/techniques/`，含 routing rules 與 migration notes | ✅ extracted |
| `skills/apk-analysis/feedback_history/` | Lesson history 與 validated/candidate learning | `intelligence/engineering/apk-analysis/` + `feedback/` + `memory/` | Lesson files 保留在 skill history；stable conclusions 以 reference 方式 promotion，非 bulk copy | candidate |
| Cross-skill handoff 到 `app-development-guidance` | Evidence recovery 與 development guidance 之間的邊界 | `workflow/apk-analysis/` | 保留 handoff artifact 與 ownership boundary | candidate |

## 首批 Reference-First 路徑

| 任務意圖 | 新路徑 | 仍需讀取 |
| --- | --- | --- |
| 決定如何觀察 APK traffic/runtime 路徑 | `analysis/apk/README.md` | `skills/apk-analysis/WORKFLOW.md`、相關 `skills/apk-analysis/techniques/` category |
| 執行 APK analysis session 或 handoff | `workflow/apk-analysis/README.md` | `skills/apk-analysis/SKILL.md`、`WORKFLOW.md`、`DOCUMENTATION.md` |
| 重複使用 APK analysis 的工程經驗 | `intelligence/engineering/apk-analysis/README.md` | `skills/apk-analysis/feedback_history/README.md`、validated lesson files |

## 未來遷移的完成條件

- 每個被搬移的 atom 使用 `metadata/schema.md` 定義 metadata。
- `knowledge/indexes/README.md` 能將 task intents 路由到新路徑與舊 skill entrypoint。
- 舊 `skills/apk-analysis/` 連結仍可解析。
- Shared rules 與 skill dependencies 仍可從舊 entrypoint 讀取。
- 驗證包含 Markdown link check、diff review、close-loop dry run、commit、push、readback 與 clean status。

## 待辦事項

- 為最高價值的 `apk-analysis` workflow 與 technique entries 建立 Knowledge Atom candidates。
- 定義 feedback lesson 何時從 `skills/apk-analysis/feedback_history/` 畢業到 `intelligence/engineering/apk-analysis/`。
- 評估 `skills/repo-analysis/` 是否為下一批遷移目標。
