# architecture.apk-analysis-pilot

| 欄位 | 值 |
| --- | --- |
| Atom ID | `architecture.apk-analysis-pilot` |
| Source path | [`../../plans/archived/apk-analysis-pilot-migration.md`](../../plans/archived/apk-analysis-pilot-migration.md) |
| Lifecycle | `candidate` |
| Summary | `apk-analysis` 作為第一個 Workflow / Analysis / Intelligence 分離 pilot 的 migration map。它建立新 reference-first 候選目的地，但保留 `skills/apk-analysis/SKILL.md` 作為 active skill entrypoint。 |
| When to read | 規劃 `apk-analysis` 內容抽取、維護舊 skill 與新分層的相容關係，或判斷哪些內容應進 `analysis/apk/`、`workflow/apk-analysis/`、`intelligence/engineering/analytical-reasoning/` 時。 |
| Do not use for | 不可把新候選路徑當成正式 replacement；pilot 期間不搬移大量內容，不刪除舊 skill entrypoint。 |
| Validation signal | 舊入口仍可讀；new reference-first paths 可找到；`knowledge/indexes/README.md` 與 routing registry 均保留 old skill source-of-truth gate。 |
| Last checked | 2026-05-14 |

## Heuristics

| Heuristic | 說明 | 來源 lessons |
|-----------|------|-------------|
| [`ui-operation-stability.md`](../../intelligence/engineering/analytical-reasoning/heuristics/ui-operation-stability.md) | UI 操作穩定性啟發式 — 決定何時該用 bounded scroll、operation script、API-first replay | 4 lessons（ui-architecture-map、ui-automation-operation-scripts、scrollable-clickable-screen-mapping、ui-fast-path-bounded-scroll） |
| [`ui-to-api-attribution.md`](../../intelligence/engineering/analytical-reasoning/heuristics/ui-to-api-attribution.md) | UI-to-API 歸因啟發式 — 確保 UI 操作能正確對應到 API 請求，避免 attribution 錯誤 | 4 lessons（screen-reachability-operation-recipes、ui-route-backfill、foreground-package-validation、feature-context-validation） |

## Checklist

- 先讀 `skills/apk-analysis/SKILL.md`。
- 需要 migration context 時讀 pilot map。
- 候選新路徑只作 mapping / promotion target，不覆蓋 skill 行為。
- 任何 promotion 都要補 metadata、knowledge index、validation 與 old entrypoint compatibility。
