# Knowledge

`knowledge/` 負責「知識導航與知識圖譜」。本層保存 Knowledge Atom、indexes、summaries、graphs 與 runtime navigation 的結構，讓 agent 能找到 task-relevant knowledge。

## 目前入口

- [`indexes/`](indexes/README.md)：第一版 task intent routing table 與 navigation index format。
- [`summaries/`](summaries/README.md)：未來 Knowledge Atom 與 source-of-truth 文件摘要格式；目前只定義格式。
- [`graphs/`](graphs/README.md)：未來 atom、source、dependency、conflict、promotion 關係圖格式；目前只定義格式。
- [`glossary/`](glossary/README.md)：Ai-skill framework / runtime / cognitive / architecture 共享語彙的 canonical 定義；schema spec + entries（`ai-skill.md`）；由 `ai-skill glossary validate` 機械強制。
- [`runtime/`](runtime/README.md)：runtime-facing knowledge view 格式、邊界與第一版 machine-readable routing registry。
  - [`runtime/sqlite/`](runtime/sqlite/README.md)：SQLite / FTS generated lookup cache 規劃；用於低 token 搜尋候選 source，不作 source-of-truth。

## 放什麼

- Knowledge Atom 的放置與索引策略。
- Navigation indexes、summaries、graphs 與 runtime lookup 設計。
- 支援 Dynamic Context Composition 的知識路由資料。
- 知識之間的 related、depends、conflicts 與 discovery path。

## 不放什麼

- Atom metadata 欄位規格；放到 `metadata/`。
- 工程智慧正文；放到 `intelligence/`。
- Agent 執行流程；放到 `workflow/`。
- 可執行 enforcement rules；放到 `enforcement/`。

## 誰會參考這裡（Inbound References）

- [`route.bootstrap.ai-skill`](../knowledge/runtime/routing-registry.yaml:21) — candidate_sources 引用 `knowledge/indexes/README.md`
- [`route.skill.discovery`](../knowledge/runtime/routing-registry.yaml:49) — candidate_sources 引用 `knowledge/indexes/README.md`、`knowledge/summaries/README.md`
- [`route.runtime.context-ttl`](../knowledge/runtime/routing-registry.yaml:102) — candidate_sources 引用 `knowledge/runtime/README.md`、`knowledge/runtime/routing-registry.yaml`
- [`route.runtime.context-loading`](../knowledge/runtime/routing-registry.yaml:161) — candidate_sources 引用 `knowledge/runtime/README.md`、`knowledge/summaries/README.md`、`knowledge/graphs/README.md`
- [`route.metadata.knowledge-atom-schema`](../knowledge/runtime/routing-registry.yaml:191) — candidate_sources 引用 `knowledge/runtime/README.md`、`knowledge/summaries/README.md`、`knowledge/graphs/README.md`
- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:319) — candidate_sources 引用 `knowledge/runtime/README.md`
- [`route.runtime.router-flow`](../knowledge/runtime/routing-registry.yaml:348) — candidate_sources 引用 `knowledge/runtime/README.md`
- [`route.runtime.context-ttl-doc`](../knowledge/runtime/routing-registry.yaml:407) — candidate_sources 引用 `knowledge/runtime/README.md`
- 共 8 條 routing records 直接引用 `knowledge/` 子目錄作為 candidate_sources

## 與既有層的關係

- `workflow/`、`analysis/`、`intelligence/` 與 `enforcement/` 是目前可直接讀取的主要內容來源。
- `metadata/` 定義 knowledge atom 的控制欄位。
- `runtime/` 使用本層 index、summary 與 graph 做 context routing。
- SQLite / FTS index 屬於 runtime lookup cache，應由 canonical Markdown / YAML 產生並可重建。
- `governance/` 定義知識 lifecycle、清理與 validation。

## 第一批候選遷移來源

- `plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md` 的 Knowledge Navigation System
- 舊 skill README 已遷移出的 workflow / analysis / intelligence 入口資訊
- `knowledge/indexes/README.md` 的 navigation index 初版
