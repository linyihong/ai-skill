# Documentation Context Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/agent-architecture/index-first-documentation.md`](../../intelligence/engineering/agent-architecture/index-first-documentation.md)
- [`intelligence/engineering/agent-architecture/attention-budgeting.md`](../../intelligence/engineering/agent-architecture/attention-budgeting.md)

本文件把 index-first documentation 的 agent architecture intelligence 轉譯成 AI runtime 文件治理。原始思想回答「為什麼 README 要當路由器」；本文件定義 agent 寫文件或整理文件時要通過的 context governance gate。

## 觸發時機

在新增或整理下列 durable 文件前套用本治理：

- project `docs/`、wiki、ADR、runbook、onboarding。
- Ai-skill workflow、governance、runtime、metadata、analysis、intelligence 文件。
- 會被 routing registry、bootstrap、skill 或 agent 重複載入的 Markdown。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Audience and lifecycle | 文件已標明主要讀者與生命週期；本輪暫存不升格為 durable policy。 |
| Kind classification | 文件有 `kind` 或等價分類；不要把 policy、how-to、reference、example 混成一檔。 |
| README-as-router | 父層 README 保留入口、摘要、讀取條件與子檔索引，不承載長篇正文。 |
| Stop condition | Agent-facing 文件在頂部說明何時不用繼續讀。 |
| Single source of truth | 同一規則只維護一份全文；其它位置只保留連結、摘要或差異。 |
| Leaf expansion | 細節、範例、長 log、API schema、操作步驟放到 leaf 文件或專案證據區。 |
| Linked update | 新增、搬移或拆分文件後，父層 README、routing、summary 或 registry 同步更新。 |

## Workflow Boundary Gate

當 workflow 文件開始承載可跨任務重用的文件治理原理時，先做分層檢查：

| 檢查 | 通過條件 |
| --- | --- |
| Why vs how | `workflow/` 只保留操作順序；why 放回 source intelligence。 |
| Gate vs checklist | 可重用治理 gate 放在本文件；單次任務 checklist 留在 workflow。 |
| Source-of-truth | 同一規則只在一層維護全文；其它層只保留摘要與連結。 |
| Runtime readiness | 只有穩定、可觀察、可驗證的 gate 才考慮升到 `runtime/` 或 `validation/`。 |

若不確定，預設保留 workflow 入口路徑不搬移，先抽薄正文並補強上游 `intelligence/` 或 `governance/`。這可避免 routing registry 的 primary source 失效。

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 文件為什麼要 index-first、為什麼 agent 會讀太多 | `intelligence/engineering/agent-architecture/` |
| 文件 context、分類、索引與停止條件的治理 gate | `governance/ai-runtime-governance/` |
| 實際在業務專案建立文件的步驟 | `workflow/documentation/` |
| 可機讀欄位、routing metadata、schema | `metadata/` 或 `knowledge/runtime/` |
| 可測 failure mode，例如 README 變成長文、leaf 未索引 | `validation/` |

## Failure Modes

| Failure mode | 風險 | 修復方向 |
| --- | --- | --- |
| README becomes article | Agent 每次被迫讀完整背景。 | 拆 leaf，README 保留路由表。 |
| Mixed document kind | 審查、搜尋與 agent route 都不穩定。 | 先標 `kind`，再依主用途拆分。 |
| Duplicate policy text | 規則 drift，更新一處漏另一處。 | 選 source-of-truth，其它位置改連結。 |
| No stop condition | Agent 不知道何時停止展開 context。 | 在導航區加入讀取條件與停止條件。 |
| Durable rule in scratch note | 重要規則無法被 routing 找到。 | 升格到正確 durable layer，scratch note 只留來源連結。 |

## Workflow Mapping

- [`workflow/documentation/README.md`](../../workflow/documentation/README.md) — 實際執行入口。
- [`workflow/documentation/execution-flow.md`](../../workflow/documentation/execution-flow.md) — 將本治理套用到業務專案文件的步驟。

## Validation Candidate

後續若要 promotion 到 `validation/`，可建立 scenario 檢查：

- 新增 leaf 文件但父層 README 未索引。
- README 同時包含長篇教學與多個主題細節。
- workflow 文件新增通用治理原則但未抽到 `intelligence/` 或 `governance/`。
