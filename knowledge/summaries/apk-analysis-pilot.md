# architecture.apk-analysis-pilot

| 欄位 | 值 |
| --- | --- |
| Atom ID | `architecture.apk-analysis-pilot` |
| Source path | [`../../plans/archived/2026-05-11-1129-apk-analysis-pilot-migration.md`](../../plans/archived/2026-05-11-1129-apk-analysis-pilot-migration.md) |
| Lifecycle | `new-layer-promoted` |
| Summary | `apk-analysis` 作為第一個 Workflow / Analysis / Intelligence 分離 pilot 的 migration map。新分層已 promoted：`workflow/apk-analysis/` 是端到端執行入口，`analysis/apk/` 保存可重用觀察、拆解與證據取得方法，`intelligence/engineering/analytical-reasoning/` 保存 reusable decision intelligence。 |
| When to read | 規劃 `apk-analysis` 內容抽取、維護舊 skill 與新分層的相容關係，或判斷哪些內容應進 `analysis/apk/`、`workflow/apk-analysis/`、`intelligence/engineering/analytical-reasoning/` 時。 |
| Do not use for | 不可把 `analysis/` 當成 raw case archive；raw logs、pcap、Frida output、class dump、host、endpoint 與一次性 reverse 過程留在業務專案 evidence，去敏抽象後才進 `feedback/history/` 或 promoted atoms。 |
| Validation signal | New reference-first paths 可找到；`knowledge/indexes/README.md`、routing registry 與 graph 均已更新為新分層路徑。 |
| Last checked | 2026-05-19 |

## Heuristics

| Heuristic | 說明 | 來源 lessons |
|-----------|------|-------------|
| [`ui-operation-stability.md`](../../intelligence/engineering/analytical-reasoning/heuristics/ui-operation-stability.md) | UI 操作穩定性啟發式 — 決定何時該用 bounded scroll、operation script、API-first replay | 4 lessons（ui-architecture-map、ui-automation-operation-scripts、scrollable-clickable-screen-mapping、ui-fast-path-bounded-scroll） |
| [`ui-to-api-attribution.md`](../../intelligence/engineering/analytical-reasoning/heuristics/ui-to-api-attribution.md) | UI-to-API 歸因啟發式 — 確保 UI 操作能正確對應到 API 請求，避免 attribution 錯誤 | 4 lessons（screen-reachability-operation-recipes、ui-route-backfill、foreground-package-validation、feature-context-validation） |

## Checklist

- 先讀 `workflow/apk-analysis/execution-flow.md`。
- 需要 migration context 時讀 pilot map。
- `analysis/apk/` 回答如何取得與拆解證據；`workflow/apk-analysis/` 回答如何執行任務順序；`intelligence/engineering/analytical-reasoning/` 回答如何判斷與避錯。
- 任何 promotion 都要補 metadata、knowledge index、validation 與 old entrypoint compatibility。
