# Repository Analysis Workflow

`workflow/repo-analysis/` 負責「Repository 分析的執行流程」。本目錄保存 agent 在分析陌生 repository 時可照著執行的 planning flow、discovery flow 與 handoff flow，讓分析過程系統化且可重複。

## Scope

本 workflow 涵蓋以下分析類型：

- **New repository onboarding**：第一次接觸 repository 時的快速理解流程。
- **Deep codebase analysis**：需要深入理解特定模組或架構時的系統化分析。
- **Migration impact analysis**：重構或遷移前的影響範圍評估。
- **Tech debt assessment**：技術債的系統化評估與優先級排序。
- **Security audit**：安全性審查的程式碼分析流程。

## 核心原則

1. **先建立心智模型，再深入細節**。先理解 repository 的整體結構、核心抽象與資料流，再深入特定模組。
2. **分析要有明確目標**。不是「讀完所有程式碼」，而是回答特定問題（如「這個模組的職責是什麼？」、「這個 bug 的根因在哪？」）。
3. **分析結果要可驗證**。每個結論應有對應的 source reference，讓其他人可以追溯。
4. **分析過程要可重複**。如果另一個 agent 用同樣的流程分析同一個 repository，應得到相似的結論。

## 與既有層的關係

- `workflow/repo-analysis/` 是 Repository 分析執行流程的主要入口。所有 agent 應優先參考本目錄的內容。
- `analysis/repo/README.md` 提供 repository 分析的具體方法（靜態結構分析、依賴分析、entrypoint 追蹤、技術債評估）。本 workflow 引用這些方法，但不複製方法細節。
- `intelligence/engineering/architecture/` 可承接從 repo 分析中萃取的架構判斷。
- `intelligence/engineering/domain/` 可承接從 repo 分析中萃取的領域模型理解。
- `skills/` 是原始 skill 目錄，內容已逐步遷移至本層。新內容應直接寫入本層。

## 第一批候選遷移來源

- `skills/app-development-guidance/process/README.md` — 已提取 Existing Project Documentation Backfill 到 `analysis/repo/documentation-backfill.md`。
- `architecture/next-stage-upgrade-plan.md` 中 `workflow/` 的分層說明。

## 已提取內容

| 來源 | 目標 | 內容 |
| --- | --- | --- |
| `skills/app-development-guidance/process/README.md` §Existing Project Documentation Backfill | `analysis/repo/documentation-backfill.md` | 已實作專案的文件恢復方法（8 種文件、6 種 pipeline artifact、7 步順序） |
| `skills/app-development-guidance/process/README.md` §Traceability Gate | `analysis/repo/traceability-gate.md` | 文件追溯性建立方法（5 種連結、stable ID、未實作標記） |
| `skills/app-development-guidance/process/README.md` §Contract Governance Gate | `analysis/repo/contract-governance.md` | 文件優先順序與衝突處理規則 |

## 建議 Workflow 流程

### New Repository Onboarding Flow

```
1. 讀取 repository 根目錄 README.md。
2. 掃描目錄結構，識別主要模組與職責。
3. 讀取建置配置（Cargo.toml、package.json、build.gradle 等），理解技術棧。
4. 識別 entrypoint 與核心流程。
5. 識別測試策略（測試目錄、測試框架、測試覆蓋率）。
6. 識別 CI/CD 流程。
7. 產出 repository 心智模型摘要。
```

### Deep Codebase Analysis Flow

```
1. 定義分析目標（理解模組 X / 追蹤功能 Y 的資料流 / 找出效能瓶頸）。
2. 使用 analysis/repo/ 的方法進行：
   ├─ 靜態結構分析：找出相關檔案與模組。
   ├─ 依賴分析：理解模組間的依賴關係。
   ├─ Entrypoint 追蹤：從 entrypoint 追蹤到目標功能。
   └─ 資料流分析：追蹤資料從輸入到儲存的完整路徑。
3. 記錄關鍵發現與 source reference。
4. 產出分析報告。
```

### Migration Impact Analysis Flow

```
1. 定義遷移範圍（模組 X / API Y / 資料庫 Z）。
2. 使用依賴分析找出所有受影響的模組與檔案。
3. 對每個受影響的項目評估：
   ├─ 修改難度（簡單 / 中等 / 困難）。
   ├─ 修改風險（低 / 中 / 高）。
   └─ 是否需要 API 版本遷移。
4. 產出影響範圍矩陣。
5. 建議遷移順序與 rollback 計畫。
```

### Documentation Backfill Flow

當分析一個已實作完成的 repository 時，使用以下流程恢復缺失的開發文件：

```
1. Inventory 現有 docs、source folders、tests、schemas、API specs、fixtures、release notes。
2. 建立 documentation gap table（exists / partial / missing / unknown）。
3. 先恢復 BDD Behavior（從 UI、API、code、tests、logs 恢復）。
4. 恢復 Domain Model、Architecture、API/Interface、Error Handling Contracts。
5. 分開標記 unknown product intent 與 observed behavior。
6. 若 BDD 無法完成，停止並要求 missing evidence。
7. 對缺乏 coverage 的 critical BDD scenario 新增 tests 或 test TODOs。
```

詳細分析方法參見 [`analysis/repo/documentation-backfill.md`](../analysis/repo/documentation-backfill.md)。

## 產出格式

每次 repo 分析應產出：

- **分析目標**（≤100 tokens）：這次分析要回答的問題。
- **心智模型摘要**（≤300 tokens）：repository 的整體結構、核心抽象、資料流。
- **關鍵發現**（每個 ≤200 tokens）：發現的問題、模式或風險，附 source reference。
- **下一步建議**（≤200 tokens）：基於分析結果的建議行動。
