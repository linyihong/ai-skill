# App Development Guidance Quick Start（開發指引快速入門）

本文件提取 [`skills/app-development-guidance/SKILL.md`](../../skills/app-development-guidance/SKILL.md) 中 Quick Start 的操作步驟，作為 `runtime/onboarding/` 層的執行指引。

> **相容性規則**：`skills/app-development-guidance/SKILL.md` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## 快速入門步驟

### 步驟 1：識別來源

判斷輸入來源類型：
- Product brief
- Observed weakness
- Reverse-engineering lesson
- App/API design question

### 步驟 2：驗證規劃文件

在程式碼變更前，檢查並驗證專案的規劃文件：
- 企劃書、product brief、planning docs、issue、ticket、PRD、design note
- BDD、API contract 或同等 artifacts

Product Brief 的主要宣稱必須標記為以下之一：
- `validated`
- `assumption`
- `open question`
- `scoped out`
- `invalidated`

在這些標記完成前，不可開始實作。

### 步驟 3：分類請求

將請求分類為：
- New requirement
- Bug fix
- Refactor
- Hardening
- Documentation-only

### 步驟 4：新需求或行為變更

若為新需求或行為變更，先更新或建立規劃文件：
- Change brief
- BDD scenarios
- Impacted Domain Model Contract
- Architecture Contract
- API / Interface Contract
- Error Handling Contract
- Implementation slices
- Tests

For embedded/hardware work，額外更新：
- Datasheet/protocol references
- Hardware context
- Driver/service/application ownership
- Fixture or hardware-in-loop validation
- Bring-up notes

在 blocker questions 解決前，不可開始寫 code。

### 步驟 5：Bug Fix 流程

若為 bug fix：

1. **Before code**：
   - 確認 expected vs actual behavior
   - 取得 reproduction/evidence
   - 確認 affected or missing BDD scenario
   - 確認 impacted contract/error handling
   - 規劃 regression test plan

2. **After code**：
   - 若 fix 改變了 **observable** behavior（含 integration-visible semantics），在**同一工作 session** 更新 owning contracts、BDD 與 project Linked Updates
   - 僅 green tests 不足以作為 Definition of Done——當 durable docs 仍描述舊行為時

### 步驟 6：定義測試策略

在 production code 之前定義測試策略：

- 區分 existing-regression coverage 與 changed/new-code validation
- 優先 BDD first，然後在新行為實作前先有 failing unit/contract/property/integration tests
- 若 performance 可能受影響，定義 budget 並選擇適當的 performance validation 類型：
  - Load test
  - Stress test
  - Spike test
  - Soak test
  - Smoke-size performance validation

### 步驟 7：從 Product Brief 開始

若從 product brief 開始，使用 [`process/`](process/) 來草擬或討論初始開發文件：
- Product Brief validation
- Bounded Contexts
- BDD behavior
- Domain Model Contract
- Architecture Contract
- API / Interface Contract
- Error Handling Contract
- Implementation slices
- Tests

### 步驟 8：既有專案回填

若在既有已實作專案上開啟此 skill：

- Audit missing documents 並回填
- Missing Product Brief fields 可標記為 `unknown` / `open question`
- BDD behavior 必須從 UI、API、code、tests、logs、fixtures 或 observed behavior 完成
- 恢復以下項目：
  - Document precedence
  - Traceability
  - BDD validation status
  - Generated-client flow
  - Vendor excerpts
  - Canceled/out-of-scope decisions

### 步驟 9：阻斷問題處理

若 behavior、domain invariants、API/interface shape、error handling、security、storage、tests、ownership、document precedence、generated clients、vendor integration 或 tool diagnostics 有任何缺失，先詢問使用者或要求 evidence 再繼續。在 blocker questions 解決前，不可繼續開發。

### 步驟 10：轉換為風險陳述

將分析發現轉換為 developer-facing risk statement 或 implementation opportunity。

### 步驟 11：選擇控制層

選擇適當的控制層：

| 控制層 | 說明 |
|--------|------|
| API/server contract | API 與伺服器端的契約 |
| App runtime behavior | App 執行時的行為控制 |
| Full-stack schema/codegen | 全端 schema 與 provider/consumer contract |
| Tooling | IDE extension、CLI、linter、static-analysis kernel/adapter |
| Vendor/third-party API | 外部 API 整合行為 |
| Embedded firmware | 硬體 context、sensor/protocol driver、board bring-up |
| Build/release configuration | 建置與發布配置 |
| Monitoring/fraud signal | 監控或 fraud 信號 |

### 步驟 12：加入驗證方法

為每個 guidance 加入具體的驗證方法：
- Brief evidence check
- Unit test
- BDD scenario
- API contract test
- Integration test
- Release checklist item
- Fixture
- Mutation/property check
- Manual review step

### 步驟 13：分類 Guidance

將 guidance 分類到對應目錄：

| 類型 | 目錄 |
|------|------|
| Core security control | `controls/` |
| Platform/app type detail | `platforms/` |
| Language/runtime-specific trap | `languages/` |
| Concrete implementation pattern | `implementation/` |
| Product-to-contract development flow | `process/` |
| Repeatable review step | `checklists/` |

### 步驟 14：套用 Required Linked Updates

依據 [`shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md)，當 process、control、platform、language、checklist、implementation pattern 或 template 受影響時，在同一變更中更新或明確驗證相關檔案。

### 步驟 15：未成熟 Lesson 處理

若 lesson 可重用但尚未成熟：
1. 先加入對應的 `feedback_history/<category>/` 資料夾
2. 跨領域 lesson 使用 `feedback_history/common/`
3. 驗證後再 promotion 到結構化目錄

## 與其他層的關係

- `workflow/app-development-guidance/execution-flow.md` 提供執行流程，本文件提供快速入門的操作步驟。
- `workflow/app-development-guidance/artifact-gates.md` 提供產出規範與品質門檻。
- `workflow/app-development-guidance/development-process.md` 提供完整的開發流程。
- `skills/app-development-guidance/SKILL.md` 是原始來源，仍為 active entrypoint。
