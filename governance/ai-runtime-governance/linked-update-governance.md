# Linked Update Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/agent-architecture/linked-updates-completeness.md`](../../intelligence/engineering/agent-architecture/linked-updates-completeness.md)
- [`intelligence/engineering/agent-architecture/context-collapse.md`](../../intelligence/engineering/agent-architecture/context-collapse.md)
- [`intelligence/engineering/agent-architecture/cognitive-boundaries.md`](../../intelligence/engineering/agent-architecture/cognitive-boundaries.md)

本文件把 linked-updates completeness 的 agent architecture intelligence 轉譯成 AI runtime governance。原始思想回答「為什麼 agent 會漏掉間接引用、索引與 generated surface」；本文件定義改動前後如何判斷連動更新完整，並界定何時必須進入 generated refresh、route update 或 explicit not-applicable。

## 觸發時機

在下列情況套用本治理：

- 新增、刪除、搬移或重命名可重用文件、目錄、route、summary、graph、template 或 workflow entrypoint。
- 修改 `enforcement/`、`governance/`、`workflow/`、`intelligence/`、`knowledge/`、`runtime/` 的 source-of-truth。
- 變更會影響 README table、routing registry、knowledge index、graph、summary、runtime report 或 generated SQLite index。
- 使用者或 review 指出 dead link、missing index、stale route、stale summary 或 incomplete close loop。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Reference discovery | 已搜尋新舊路徑、slug、route id、標題與主要 anchor 的 inbound references。 |
| Owner surface update | 所屬目錄 README、root/index README、routing registry 或 knowledge index 已同步，或標 `not_applicable`。 |
| Generated surface decision | Summaries、graphs、runtime reports、SQLite index、runtime.db 是否需 refresh 已判斷並執行或記錄理由。 |
| Source-of-truth preservation | 同一規則不在多處維護全文；次要位置只保留摘要與連結。 |
| Explicit non-update | 相關 surface 不改時，必須能說明為何不受影響，不可沉默略過。 |
| Commit boundary | 主變更與必要 linked updates 同 commit；不把必要 linked update 留到「之後」。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 為什麼 linked update 容易遺漏、遺漏會造成什麼系統退化 | `intelligence/engineering/agent-architecture/` |
| Completeness gate、generated surface decision、explicit non-update criteria | `governance/ai-runtime-governance/` |
| 全庫必須同步更新的具體表格與 commit/push/readback 條文 | `enforcement/linked-updates.md` |
| Knowledge update phase / transaction / generated surface runtime view | `governance/lifecycle/`、`runtime/`、`knowledge/runtime/` |
| Architecture-layer validation checklist | `governance/validation/` |

## Runtime Mapping

- [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) — 可執行 mandatory linked-update rule。
- [`enforcement/dependency-reading.md`](../../enforcement/dependency-reading.md) — dependency read ledger 與 writeback transaction close loop。
- [`governance/lifecycle/knowledge-update-flow.md`](../lifecycle/knowledge-update-flow.md) — knowledge update step orchestration。
- [`governance/validation/README.md`](../validation/README.md) — validation checklist 與 generated refresh gates。

## Validation Candidate

後續若要 promotion 到 `validation/`，可建立 scenario 檢查：

- 新增 governance doc 但未更新 `governance/ai-runtime-governance/README.md`。
- 新增 route 但未 refresh generated runtime reports。
- 搬移 source path 後舊 inbound links 未更新。
- Agent 宣稱 linked updates complete，但沒有搜尋 reference 或說明 non-update。
