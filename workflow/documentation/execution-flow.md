# Agent-Friendly Documentation — Execution Flow

在 **業務專案** 或其他 repo 寫文件時，依序執行下列步驟。目標：**人類可維護**、**agent 可路由且少讀無關字**、**分類穩定**。

上層依據：

- 判斷智慧：[`index-first-documentation.md`](../../intelligence/engineering/agent-architecture/index-first-documentation.md)
- 治理 gate：[`documentation-context-governance.md`](../../governance/ai-runtime-governance/documentation-context-governance.md)

## 0. 釐清讀者與生命週期

用一句話寫在 PR 或 doc 頂部註記（可刪）：

```text
讀者：<人類 only | agent+人類>
生命週期：<本輪對話 | 專案期間 | 組織長期>
```

- **本輪對話**：優先放 `<PROJECT_ROOT>/.agent-goals/` 或暫存筆記，不要當成長期規範。
- **專案期間**：放 `docs/`、`design/`、`adr/` 等專案內 durable 路徑。
- **組織長期 / 多專案重用**：抽象成與 repo 無關的規則語意；具體 host、ticket、log 留在專案證據區。對照 [`../../enforcement/reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md)。

## 1. 分類維度（寫入前先選標籤）

為每份新文件選 **一至兩個主分類**（列在檔案最上方短表或目錄 README 的索引列），避免「全部丟進 misc」。

| 維度 | 範例值 | 用途 |
| --- | --- | --- |
| `kind` | `how-to` / `reference` / `policy` / `decision` / `runbook` / `glossary` | 決定段落結構與詳略；`policy` 要短而可驗證。 |
| `audience` | `human` / `agent` / `both` | `agent` 為主時必須有「何時讀、讀到哪停」的明確停止條件。 |
| `stability` | `draft` / `active` / `deprecated` | 讓搜尋與 agent 可過濾過期文。 |
| `routing` | `index` / `leaf` | `index` 只連結與摘要；`leaf` 承載細節。 |

專案若已有 front-matter（YAML）慣例，可把上表欄位映射過去；若沒有，**在父層 README 用表格維護同一資訊**即可。

## 2. 選目錄與檔名（對照 content layering）

依 [`../../enforcement/content-layering.md`](../../enforcement/content-layering.md) 的意圖，映射到專案約定路徑（範例）：

| 意圖 | 建議路徑形狀 |
| --- | --- |
| 入口與路由 | `docs/README.md` 或主題資料夾的 `README.md` |
| 操作步驟 | `docs/how-to/<topic>.md` |
| 規格與欄位字典 | `docs/reference/<topic>.md` |
| 決策紀錄 | `docs/adr/` 或 `constitution/` |
| 本輪實驗／待驗證 | `.agent-goals/` 或 `notes/scratch/`（不當作單一真相） |

**檔名**：優先 `kebab-case`、主題明確；避免 `notes2.md`、`misc.md`。

## 3. 檔案形狀（降低 token）

1. **頂部 8–15 行「導航區」**：目的、讀者、`kind`、主要依賴連結、**何時不必讀下去**。
2. **索引優先**：長主題拆成多 `leaf`；父檔只保留表格列「意圖 → 子檔路徑 → 讀取條件」。
3. **少貼全文**：程式碼、OpenAPI、長 log 用「路徑 + 節圍（行號或 anchor）」引用，不整段複製進規範文。
4. **單一真相**：同一規則不在兩處維護兩份全文；第二處只留一行連結與差異說明。
5. **拆分閾值**：見 [`../../governance/document-sizing.md`](../../governance/document-sizing.md)（行數與主題單一性）。

若上述規則開始被其它 workflow 重複引用，不在本檔擴寫治理原理；回到 [`documentation-context-governance.md`](../../governance/ai-runtime-governance/documentation-context-governance.md) 補充 gate，workflow 只保留操作順序。

## 3a. 判斷是否需要結構化流程契約

其他專案不需要把所有文件都 YAML 化。寫入或整理文件時，先做以下判斷：

| 文件內容 | 標準做法 |
| --- | --- |
| 說明、背景、決策脈絡、設計理由、操作備忘 | 保持 Markdown；可用專案既有 front-matter 或父層 README 索引。 |
| 文件只需要人類閱讀，agent 不會反覆依它執行步驟 | 保持 Markdown，並在導航區寫清楚讀者與停止條件。 |
| 文件定義 agent 要反覆執行的流程、runbook、release checklist、required evidence、failure action 或 blocking gate | 建立 companion YAML 或專案等價 structured contract。 |
| 專案已有 front-matter / workflow schema | 將欄位映射到 activation、required sources、steps、gates、success criteria、failure modes、final report。 |

Ai-skill 本庫的欄位語意見 [`../../metadata/executable-contract-schema.md`](../../metadata/executable-contract-schema.md)。其他專案可用自己的格式，但必須保留同等判斷語意：`contract_required`、`markdown_only` 或 `not_applicable`。

## 4. 驗證與連動

- **Link**：從本檔或 `README` 出去的相對連結可點、無斷鏈。
- **路由**：父層 README 的表格含新子檔；routing registry、workflow index 或專案內索引若有對應列，一併更新。
- **分層**：若新增的是 why / 治理 gate，而不是操作步驟，回到 [`index-first-documentation.md`](../../intelligence/engineering/agent-architecture/index-first-documentation.md) 或 [`documentation-context-governance.md`](../../governance/ai-runtime-governance/documentation-context-governance.md)，不要在 workflow 重複全文。
- **結構化流程契約**：若 Step 3a 判定為 `contract_required`，確認 companion YAML 或等價 structured contract 已建立；若暫不建立，記錄缺口與理由。
- **Ai-skill 本庫**：若同步規則或模板，依 [`../../enforcement/dependency-reading.md`](../../enforcement/dependency-reading.md) 與 [`../../enforcement/linked-updates.md`](../../enforcement/linked-updates.md) 走完整閉環。

## 5. 完成訊號（給 agent 自查）

- [ ] 新文件有 `kind` / `audience` / `stability`（或父 README 等價欄位）。
- [ ] 有 `index` 或導航區；長文已拆 `leaf`。
- [ ] 已判斷 `contract_required` / `markdown_only` / `not_applicable`，必要時有 companion YAML 或等價 structured contract。
- [ ] 無重複規範全文；斷鏈已修。
- [ ] 本輪暫存與長期 durable 分界清楚。
- [ ] Workflow 只保留操作步驟；可重用 why / gate 已連回 intelligence 或 governance。

← [回到本 workflow 索引](README.md) · [workflow 總索引](../README.md)
