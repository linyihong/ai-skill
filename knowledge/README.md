# Knowledge

`knowledge/` 負責「知識導航與知識圖譜」。本層保存 Knowledge Atom、indexes、summaries、graphs 與 runtime navigation 的結構，讓 agent 能找到 task-relevant knowledge。

## 目前入口

- [`indexes/`](indexes/README.md)：第一版 task intent routing table 與 navigation index format。

## 放什麼

- Knowledge Atom 的放置與索引策略。
- Navigation indexes、summaries、graphs 與 runtime lookup 設計。
- 支援 Dynamic Context Composition 的知識路由資料。
- 知識之間的 related、depends、conflicts 與 discovery path。

## 不放什麼

- Atom metadata 欄位規格；放到 `metadata/`。
- 工程智慧正文；放到 `intelligence/`。
- Agent 執行流程；放到 `workflow/`。
- 可執行 shared rules；放到 `shared-rules/`。

## 與既有層的關係

- `skills/` 與 `shared-rules/` 仍是目前可直接讀取的主要內容來源。
- `metadata/` 定義 knowledge atom 的控制欄位。
- `runtime/` 使用本層 index、summary 與 graph 做 context routing。
- `governance/` 定義知識 lifecycle、清理與 validation。

## 第一批候選遷移來源

- `architecture/next-stage-upgrade-plan.md` 的 Knowledge Navigation System
- `skills/README.md` 與各 skill README 中可抽成全庫索引的入口資訊
- `knowledge/indexes/README.md` 的 navigation index 初版
