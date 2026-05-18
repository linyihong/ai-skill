# Greenfield Execution Flow（新專案標準化執行流程）

## 流程總覽

```
使用者需求 → [Specify] → Feature Specification → [Plan] → Implementation Plan → [Tasks] → Task Breakdown → [Implement] → 實作程式碼 + 測試
```

---

## Phase 1: Specify（規格定義）

### Entry Condition
- 使用者已提供需求描述（文字、產品簡報、issue）

### Process
1. 解析使用者需求，識別核心功能與邊界
2. 拆解為多個 **Independent User Story**，每個 story 必須：
   - 可獨立測試
   - 可獨立交付
   - 有明確的 Priority（P1/P2/P3）
3. 為每個 User Story 撰寫 **Acceptance Scenarios**（Given/When/Then）
4. 列出 **Edge Cases**
5. 定義 **Functional Requirements**（FR-001, FR-002...）
6. 定義 **Success Criteria**（SC-001, SC-002...）
7. 記錄 **Assumptions**

### 產出
- Feature Specification（使用 [`spec-template.md`](templates/spec-template.md)）

### Gate
- [ ] 所有 User Story 有明確 Priority
- [ ] 每個 P1 Story 有至少 1 個 Acceptance Scenario
- [ ] Functional Requirements 已列出
- [ ] Success Criteria 可衡量

### Template Reference
使用 [`spec-template.md`](templates/spec-template.md) 作為輸出格式。

---

## Phase 2: Plan（技術計畫）

### Entry Condition
- Feature Specification 已通過 Gate

### Process
1. 從 Feature Specification 提取核心需求
2. 決定 **Technical Context**：
   - 語言/版本
   - 主要依賴
   - 儲存方案
   - 測試框架
   - 目標平台
   - 專案類型
   - 效能目標
   - 限制條件
3. 執行 **Constitution Check**：
   - 架構一致性
   - 依賴授權相容性
   - 安全與合規需求
4. 設計 **Project Structure**（文件樹 + 原始碼樹）
5. 記錄 **Complexity Tracking**（如有架構違規需說明理由）

### 產出
- Implementation Plan（使用 [`plan-template.md`](templates/plan-template.md)）

### Gate
- [ ] Technical Context 完整填寫
- [ ] Constitution Check 全部通過（或違規已記錄理由）
- [ ] Project Structure 已定義

### Template Reference
使用 [`plan-template.md`](templates/plan-template.md) 作為輸出格式。

---

## Phase 3: Tasks（任務拆解）

### Entry Condition
- Implementation Plan 已通過 Gate

### Process
1. 根據 User Story 優先級拆解任務
2. 任務分組：
   - **Phase 1: Setup** — 專案初始化、依賴安裝、linting 設定
   - **Phase 2: Foundational** — 阻斷性基礎設施（DB schema、auth、routing）
   - **Phase 3+: User Stories** — 依 Priority 順序（P1 → P2 → P3）
   - **Final Phase: Polish** — 跨 story 的優化、文件、效能
3. 標記可並行任務（`[P]`）
4. 標記任務所屬 User Story（`[US1]`, `[US2]`...）
5. 定義 **Dependencies & Execution Order**
6. 選擇 **Implementation Strategy**（MVP First / Incremental / Parallel Team）

### 產出
- Task Breakdown（使用 [`tasks-template.md`](templates/tasks-template.md)）

### Gate
- [ ] 每個 User Story 有對應的任務群組
- [ ] 任務有明確的檔案路徑
- [ ] Dependencies 已定義
- [ ] Implementation Strategy 已選擇

### Template Reference
使用 [`tasks-template.md`](templates/tasks-template.md) 作為輸出格式。

---

## Phase 4: Implement（實作）

### Entry Condition
- Task Breakdown 已通過 Gate

### Process
1. 依 Task Breakdown 的順序實作
2. 每個 User Story 完成後執行獨立驗證
3. 使用 [`software-delivery`](../software-delivery/execution-flow.md) 的 BDD Closure Loop 確保測試覆蓋
4. 使用 [`software-delivery`](../software-delivery/execution-flow.md) 的 Review 流程進行 Code Review

### 產出
- 實作程式碼
- 測試（unit / integration / contract）
- 文件（README、API docs）

### Gate
- [ ] 所有 P1 User Story 已完成並通過測試
- [ ] BDD Closure 已執行
- [ ] Code Review 已完成
- [ ] 文件已更新

### Template Reference
- BDD Scenarios: [`bdd-scenario-template.md`](../software-delivery/templates/bdd-scenario-template.md)
- Review Report: [`review-report-template.md`](../software-delivery/templates/review-report-template.md)

---

## 快速參考

| 階段 | 主要問題 | 產出 | 關鍵 Gate |
|------|---------|------|-----------|
| Specify | 做什麼？ | Feature Specification | User Story 完整性 |
| Plan | 怎麼做？ | Implementation Plan | Constitution Check |
| Tasks | 誰做？何時做？ | Task Breakdown | Dependency 定義 |
| Implement | 做出來！ | 程式碼 + 測試 | BDD Closure + Review |
