# Agent-Friendly Documentation Workflow

`workflow/documentation/` 負責「在 **業務專案或其他 repository** 撰寫、整理、分類文件時，讓內容對 **人類與 AI agent** 都好讀、好路由、且 **降低無效 token**」。本層只保存**可執行的步驟與分類表**；背後的 index-first 判斷智慧見 [`index-first-documentation.md`](../../intelligence/engineering/agent-architecture/index-first-documentation.md)，AI runtime 治理 gate 見 [`documentation-context-governance.md`](../../governance/ai-runtime-governance/documentation-context-governance.md)，拆檔與篇幅閾值見 [`governance/document-sizing.md`](../../governance/document-sizing.md)。**語言與是否工具中立**不由本 workflow 預設，由業務專案自訂（見 [不取代什麼](#不取代什麼)）。

本 workflow 也提供 project documentation 的結構化契約判斷：一般說明、背景、決策脈絡保留 Markdown；若文件定義 agent 要反覆執行的流程、runbook、release checklist、required evidence、failure action 或 blocking gate，應建立 companion YAML 或專案等價 structured contract。欄位語意可對照 [`../../metadata/executable-contract-schema.md`](../../metadata/executable-contract-schema.md)，但不強制其他專案照抄 Ai-skill schema。

## 何時讀這裡

- 要在 `<PROJECT_ROOT>` 建立或擴充 `docs/`、`wiki/`、ADR、onboarding、runbook、知識庫型筆記。
- 文件會被 **Ai-skill 或 agent 反覆載入**，需要 **index-first**、可選讀子檔，而不是單一巨大 Markdown。
- 需要 **一致的分類維度**，讓之後搜尋、registry、或人類瀏覽都穩定。

若專案需要把 durable Markdown 寫作設為預設觸發，應在專案自己的 tool config 或 overlay 中指向本目錄與 `execution-flow.md`。`ai-skill init-project` 產生的通用 bootstrap 只保留 thin pointer，不複製此 workflow 規則。

## 不取代什麼

- **語言、語氣、是否綁定特定工具**：由**業務專案**自定（CONTRIBUTING、style guide、或專案 README）；本 workflow **不**預設要求工具中立或特定語言政策。若專案要對齊 Ai-skill 可重用規範，可另讀 [`../../enforcement/neutral-language.md`](../../enforcement/neutral-language.md)、[`../../enforcement/tool-neutral-documentation.md`](../../enforcement/tool-neutral-documentation.md)。
- **何時拆檔、行數閾值**：見 [`../../governance/document-sizing.md`](../../governance/document-sizing.md)。
- **內容應落在哪一層（durable vs `.agent-goals/`）**：見 [`../../enforcement/content-layering.md`](../../enforcement/content-layering.md)。
- **先界定未知、再載入 context**：見 [`../../enforcement/decision-efficiency.md`](../../enforcement/decision-efficiency.md)。

## 執行記憶

詳細治理 gate 以 [`documentation-context-governance.md`](../../governance/ai-runtime-governance/documentation-context-governance.md) 為準；判斷智慧以 [`index-first-documentation.md`](../../intelligence/engineering/agent-architecture/index-first-documentation.md) 為準。本 workflow 只保留執行時需要快速檢查的版本：

1. 先標讀者與生命週期。
2. 選 `kind` / `audience` / `stability` / `routing`。
3. 父層 `README.md` 保持 router；細節放 leaf。
4. 重複規則只留一份 source-of-truth，其它位置連結。
5. 變更後更新父層 README、routing 或索引。

## 已提取內容

| 檔案 | 用途 |
| --- | --- |
| [`execution-flow.md`](execution-flow.md) | 從分類、選位、檔案形狀到驗證與連動更新的執行步驟與分類維度表。 |
| [`execution-flow.yaml`](execution-flow.yaml) | 本 workflow 的 executable contract，結構化 reader/lifecycle、classification、YAMLization decision、validation 與 final report gates。 |
| [`../../metadata/executable-contract-schema.md`](../../metadata/executable-contract-schema.md) | 判斷其他專案文件何時需要 companion YAML 或等價 structured contract 的欄位語意。 |
| [`../../intelligence/engineering/agent-architecture/index-first-documentation.md`](../../intelligence/engineering/agent-architecture/index-first-documentation.md) | 為什麼 agent-facing 文件要 index-first 的判斷智慧。 |
| [`../../governance/ai-runtime-governance/documentation-context-governance.md`](../../governance/ai-runtime-governance/documentation-context-governance.md) | 文件 context、分類、停止條件與單一真相的 AI runtime governance。 |

## 與既有層的關係

- 本 workflow 可依需要 **引用** [`content-layering.md`](../../enforcement/content-layering.md)、[`decision-efficiency.md`](../../enforcement/decision-efficiency.md)；**不**預設要求讀 [`tool-neutral-documentation.md`](../../enforcement/tool-neutral-documentation.md)。
- `workflow/repo-analysis/` 與 [`../../analysis/repo/documentation-backfill.md`](../../analysis/repo/documentation-backfill.md) 偏「已存在專案的補文件」；本目錄偏「從一開始就寫成 agent 友善形狀」。
- `workflow/software-delivery/` 的 [`artifact-gates.md`](../software-delivery/artifact-gates.md) 有 reusable note 結構；可與本 flow 並用。

## 誰會參考這裡（Inbound References）

- [`route.workflow.documentation-ai-native`](../../knowledge/runtime/routing-registry.yaml) — `primary_source` 指向本 `README.md`。

← [回到 workflow 索引](../README.md)
