# Project Memory

`memory/project/` 保存**跨 session 的專案脈絡**。不同於 `memory/episodic/`（情境經驗），project memory 記錄的是「這個專案是什麼、目標是什麼、有哪些已知約束與決策」，讓 agent 在處理同一專案的不同 session 之間保持上下文連續性。

Project memory 只在 same repo / same project / compatible architecture boundary 下 replay。Repo refactor、migration、dependency change 或 branch context 改變時，必須重新驗證 freshness 與 compatibility scope。

## 用途

- 記錄專案的目標、範圍與已知約束
- 保存專案特有的架構決策與取捨理由
- 追蹤專案中已知的技術債、風險與待辦事項
- 支援跨 session 的專案上下文銜接（例如：今天做到一半，明天繼續）
- 提供 `workflow/` 執行時的專案特定 context

## 不放什麼

- 專案私有 raw evidence（token、host、incident logs）→ 留在業務專案
- Active goal、owner、lock、next action → `.agent-goals/`
- 專案獨立的 architecture decision → `memory/decision/` 或 `decisions/`
- 可跨專案重用的 intelligence → `intelligence/`
- Session 進行中狀態 → `memory/working/`

## 格式

```markdown
# Project: {專案名稱}

## Status
{active | paused | completed | archived}

## Goals
- {goal 1}
- {goal 2}

## Known Constraints
- {constraint 1}（來源：{source}）
- {constraint 2}（來源：{source}）

## Architecture Decisions
- {decision 1}（→ memory/decision/ 或 decisions/）
- {decision 2}

## Tech Debt / Risks
- {item 1}（severity: {high|medium|low}）
- {item 2}

## Key Contacts / References
- {reference 1}
- {reference 2}

## Session History
- {date}：{brief summary}（→ memory/summary/{session-file}）
- {date}：{brief summary}
```

## 規則

1. **專案為單位**：每個專案一個檔案（或一個目錄），不混合多個專案。
2. **抽象化邊界**：只記錄「可跨 session 重用的專案脈絡」，不記錄 session-local 細節。
3. **連結到 summary**：每個 session 結束時，如果影響專案狀態，應更新 project memory 並連結到對應的 session summary。
4. **不重複 ADR**：重要的架構決策應放在 `memory/decision/` 或 `decisions/`，project memory 只保留摘要與連結。
5. **Lifecycle**：專案完成或歸檔後，project memory 應標記為 `archived` 並移至冷儲存。
6. **Token-aware**：每個 project record 不超過 500 tokens。
7. **Compatibility scope**：每個 project memory 應能說明適用 repo、architecture boundary、workflow family 與 expires condition。
8. **Not active state**：不得保存 active goal、owner、lock、next action 或 current blocker。
9. **Source revalidation**：Project memory 可提供 context，但不能取代 current files、tests、runtime reports 或 user goal。

## 與既有層的關係

- `memory/summary/`：session summary 連結到 project memory
- `memory/decision/`：專案相關的架構決策記錄
- `memory/episodic/`：專案中發生的特定情境經驗
- `workflow/`：workflow 執行時需要 project memory 作為 context
- `governance/lifecycle/`：project memory 的 lifecycle 管理（active → archived）
- `knowledge/`：project memory 可被 knowledge navigation 索引
- `memory/retrieval-governance/`：定義 project memory 的 freshness decay、compatibility scope 與 activation threshold
