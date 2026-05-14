# Agent-Friendly Documentation Workflow

`workflow/documentation/` 負責「在 **業務專案或其他 repository** 撰寫、整理、分類文件時，讓內容對 **人類與 AI agent** 都好讀、好路由、且 **降低無效 token**」。本層是**可執行的步驟與分類表**；政策與閾值仍以 `enforcement/` 與 `governance/` 為準。

## 何時讀這裡

- 要在 `<PROJECT_ROOT>` 建立或擴充 `docs/`、`wiki/`、ADR、onboarding、runbook、知識庫型筆記。
- 文件會被 **Ai-skill 或 agent 反覆載入**，需要 **index-first**、可選讀子檔，而不是單一巨大 Markdown。
- 需要 **一致的分類維度**，讓之後搜尋、registry、或人類瀏覽都穩定。

## 不取代什麼

- **語言與中性用語**：見 [`../../enforcement/neutral-language.md`](../../enforcement/neutral-language.md)。
- **工具專屬路徑與 sync**：見 [`../../enforcement/tool-neutral-documentation.md`](../../enforcement/tool-neutral-documentation.md) 與 [`../../ai-tools/README.md`](../../ai-tools/README.md)。
- **何時拆檔、行數閾值**：見 [`../../governance/document-sizing.md`](../../governance/document-sizing.md)。
- **內容應落在哪一層（durable vs `.agent-goals/`）**：見 [`../../enforcement/content-layering.md`](../../enforcement/content-layering.md)。
- **先界定未知、再載入 context**：見 [`../../enforcement/decision-efficiency.md`](../../enforcement/decision-efficiency.md)。

## 核心原則（精簡版）

1. **README 當路由器**：目錄下的 `README.md` 只放「何時讀哪個子檔」與短摘要，不放長篇教學。
2. **一檔一事**：同一檔不要混「規範 + 長範例 + checklist + 工具操作」；混了就拆（見 document-sizing）。
3. **先索引再全文**：Agent 先讀索引表或 summary 區塊，證據需要時才打開子檔或專案內 deep link。
4. **可攜、可去敏**：可重用段落避免綁定單一工具、內網主機、或客戶證據；證據留在專案區並抽象成規則。
5. **分類可機讀**：檔名路徑或 front-matter（若專案慣例有）能表達「類型 + 生命週期」，方便之後自動化或手動路由。

## 已提取內容

| 檔案 | 用途 |
| --- | --- |
| [`execution-flow.md`](execution-flow.md) | 從分類、選位、檔案形狀到驗證與連動更新的執行步驟與分類維度表。 |

## 與既有層的關係

- 本 workflow **引用** `enforcement/`、`governance/document-sizing.md`，不重複貼全文。
- `workflow/repo-analysis/` 與 [`../../analysis/repo/documentation-backfill.md`](../../analysis/repo/documentation-backfill.md) 偏「已存在專案的補文件」；本目錄偏「從一開始就寫成 agent 友善形狀」。
- `workflow/software-delivery/` 的 [`artifact-gates.md`](../software-delivery/artifact-gates.md) 有 reusable note 結構；可與本 flow 並用。

## 誰會參考這裡（Inbound References）

- [`route.workflow.documentation-ai-native`](../../knowledge/runtime/routing-registry.yaml) — `primary_source` 指向本 `README.md`。
- [`skills-index.yaml`](../../skills-index.yaml) 的 `documentation-ai-native` skill 條目。

← [回到 workflow 索引](../README.md)
