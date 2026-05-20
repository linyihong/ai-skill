# Index-First Documentation（索引優先文件）

**Status**: `candidate-intelligence`
**Source**: documentation workflow extraction

## 原則

**Documentation that agents must repeatedly read should expose routeable indexes before detailed prose.**

會被 agent 反覆載入的文件，應先提供可路由的索引，再把細節放到 leaf 文件。README 的主要責任是幫 agent 和人類判斷「何時讀哪個子檔」，不是承載所有背景、規範與範例。

## 為什麼

1. **Agent 需要停止條件** — 若入口文件沒有「何時不必繼續讀」與子檔讀取條件，agent 容易展開過多 context。
2. **長文會稀釋注意力** — 把 policy、how-to、reference、example 全放在同一檔，會讓重點與任務無關細節競爭注意力。
3. **索引比全文更穩定** — 父層 README 可以長期保存 route map，leaf 文件可隨主題演進獨立更新。
4. **單一真相降低維護成本** — 同一規則只在 source 文件維護，其他位置保留連結與差異說明。

## 判斷信號

| 信號 | 意義 |
| --- | --- |
| README 超過索引責任 | 應拆 leaf，README 保留用途、讀取條件與摘要。 |
| 一檔混合 policy / how-to / reference / examples | 需要重新標記 `kind` 並拆分。 |
| Agent 每次都讀完整文件才知道下一步 | 需要補導航區或父層索引。 |
| 同一規則在多處有全文 | 需要選 source-of-truth，其它位置改連結。 |

## 路由邊界

Index-first 不是把文件切碎，而是讓入口文件只承擔「路由」責任。當 README 同時保存背景、治理原理、操作步驟與範例時，agent 會把整份文件視為任務必要 context，造成 attention dilution。

判斷方式：

| 內容 | 優先位置 |
| --- | --- |
| 為什麼 agent 需要停止條件、索引、單一真相 | `intelligence/engineering/agent-architecture/` |
| AI runtime 文件治理 gate、分類欄位、linked update gate | `governance/ai-runtime-governance/` |
| 實際建立或整理文件的步驟 | `workflow/documentation/` |
| 行數閾值與拆檔形狀 | `governance/document-sizing.md` |
| 可機讀欄位或 routing metadata | `metadata/` 或 `knowledge/runtime/` |

若某段文字會被多個 workflow 引用，它通常不是 workflow 本體；先抽成 source intelligence 或 governance gate，再由 workflow 連回去。

## 適用範圍

- 專案 `docs/`、wiki、ADR、runbook、onboarding。
- Ai-skill 的 workflow、governance、runtime、metadata 索引文件。
- 任何會被 routing registry、bootstrap、skill 或 agent 重複讀取的 durable Markdown。

## 不適用範圍

- 一次性 scratch notes。
- 已明確 scoped 的小型 leaf 文件。
- 需要完整連續敘事的人類報告；但仍應提供短摘要與目錄。

## Governance Translation

- [`governance/ai-runtime-governance/documentation-context-governance.md`](../../../governance/ai-runtime-governance/documentation-context-governance.md)

## 相關 atoms

- [`attention-budgeting.md`](attention-budgeting.md) — 注意力預算。
- [`task-routing.md`](task-routing.md) — 任務路由。
- [`decomposition-strategy-selection.md`](decomposition-strategy-selection.md) — 拆分策略選擇。

---

← [回到 agent-architecture/](README.md)
