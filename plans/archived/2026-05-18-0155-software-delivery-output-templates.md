# Plan: Software Delivery Output Templates

## 緣起

分析 [github/spec-kit](https://github.com/github/spec-kit) 後發現，spec-kit 的核心價值在於**標準化輸出模板**（spec-template, plan-template, tasks-template），讓每次開發的產出格式一致、可被工具消費。

我們的 [`workflow/software-delivery/`](workflow/software-delivery/README.md) 在流程深度（BDD 閉環、Change Intake Gate、Performance Test Gate、Contract Governance）上已完全涵蓋 spec-kit，唯一缺少的就是**標準化輸出模板**。

本 plan 目標：為 software-delivery workflow 建立 5 個輸出模板，讓每個階段的產出格式一致、可追溯、可被後續階段自動消費。

## 狀態

- **Status**: ✅ completed
- **Author**: Roo (architect → code)
- **Created**: 2026-05-18 01:55 UTC
- **Completed**: 2026-05-18 11:30 UTC+9
- **Priority**: Medium

## 相依性

- [`workflow/software-delivery/execution-flow.md`](workflow/software-delivery/execution-flow.md) — 需在對應階段引用模板
- [`workflow/software-delivery/development-process.md`](workflow/software-delivery/development-process.md) — 需在對應關卡引用模板
- [`workflow/software-delivery/README.md`](workflow/software-delivery/README.md) — 需在「產出格式」章節引用模板
- [`knowledge/runtime/routing-registry.yaml`](knowledge/runtime/routing-registry.yaml) — 需註冊新模板路徑

## 執行計畫

### Phase 1: 建立模板檔案（5 個）

在 `workflow/software-delivery/templates/` 下建立以下模板：

#### 1. [`change-brief-template.md`](workflow/software-delivery/templates/change-brief-template.md)

Change Intake 階段的輸出模板。對應 [`execution-flow.md#變更接收`](workflow/software-delivery/execution-flow.md:18) 與 [`development-process.md#Change Intake Gate`](workflow/software-delivery/development-process.md:89)。

```markdown
# Change Brief: <title>

## Metadata
- **Change Type**: <feature|bugfix|refactor|perf|docs|chore>
- **Priority**: <p0|p1|p2|p3>
- **Evidence Source**: <issue|incident|product-brief|customer-feedback>
- **Date**: <YYYY-MM-DD>

## Evidence Summary
<Link to source evidence, max 3 bullet points>

## Scope
<What is in scope / out of scope>

## Blocker Assessment
- [ ] No blocker — proceed to Contract phase
- [ ] Blocker identified: <description>
```

#### 2. [`contract-template.md`](workflow/software-delivery/templates/contract-template.md)

Contract Governance 階段的輸出模板。對應 [`development-process.md#Contract Governance Gate`](workflow/software-delivery/development-process.md:102)。

```markdown
# Contract: <feature/component>

## Domain Model
<Key entities and relationships>

## Architecture Decision
<ADR-style record: context → decision → consequence>

## API Contract
<Endpoint / method / request / response / error codes>

## Error Handling
<Expected error scenarios and handling strategy>

## Data Model Changes
<Schema changes, migrations, backward compatibility>

## Traceability
- **Change Brief**: <link>
- **BDD Scenarios**: <link>
```

#### 3. [`bdd-scenario-template.md`](workflow/software-delivery/templates/bdd-scenario-template.md)

BDD Closure Loop 階段的輸出模板。對應 [`execution-flow.md#文件優先 BDD 閉環`](workflow/software-delivery/execution-flow.md:37)。

```markdown
# BDD Scenario: <feature/behavior>

## Scenario: <title>
**Given** <precondition>
**When** <action>
**Then** <expected outcome>

## Scenario: <title>
**Given** <precondition>
**When** <action>
**Then** <expected outcome>

## Edge Cases
- <edge case 1>
- <edge case 2>

## Regression Scope
- [ ] Existing tests affected: <list>
- [ ] New tests required: <count>
```

#### 4. [`implementation-plan-template.md`](workflow/software-delivery/templates/implementation-plan-template.md)

Implementation 階段的任務拆解模板。對應 [`development-process.md#Default Flow`](workflow/software-delivery/development-process.md:7)。

```markdown
# Implementation Plan: <feature>

## Task Breakdown
### Task 1: <title>
- **File(s)**: <path>
- **Description**: <what to do>
- **Acceptance**: <how to verify>

### Task 2: <title>
- **File(s)**: <path>
- **Description**: <what to do>
- **Acceptance**: <how to verify>

## Dependencies
- <prerequisite task or external dependency>

## Risk Assessment
- <potential risk and mitigation>
```

#### 5. [`review-report-template.md`](workflow/software-delivery/templates/review-report-template.md)

Review 階段的輸出模板。對應 [`README.md#審查類型`](workflow/software-delivery/README.md:21) 的 6 種 review。

```markdown
# Review Report: <type> — <feature>

## Review Type
<Design Review | Code Review | Release Review | Security Review | Contract Governance Review | Embedded Firmware Review>

## Findings
### Critical
- <finding>

### Warning
- <finding>

### Suggestion
- <finding>

## Decision
- [ ] Approved
- [ ] Changes requested: <list>
- [ ] Blocked: <reason>

## Artifacts Reviewed
- <link to contracts, code, docs>
```

### Phase 2: 更新 Workflow 文件

#### 2a. [`execution-flow.md`](workflow/software-delivery/execution-flow.md)

在每個階段的結尾加入「輸出模板」小節，引用對應模板：

| 階段 | 插入位置 | 模板 |
|------|---------|------|
| 1. 從證據開始 → 變更接收 | 在 Change Intake 段落後 | `change-brief-template.md` |
| 2. 文件優先 BDD 閉環 | 在 BDD Closure Loop 段落後 | `bdd-scenario-template.md` |
| 4. 同工作階段閉環 | 在 Same-Session Closure 段落後 | `implementation-plan-template.md` |
| 7. 驗證 | 在 Validate 段落後 | `review-report-template.md` |

#### 2b. [`development-process.md`](workflow/software-delivery/development-process.md)

在對應關卡加入模板引用：

| 關卡 | 插入位置 | 模板 |
|------|---------|------|
| Change Intake Gate | 在 gate 說明後 | `change-brief-template.md` |
| Contract Governance Gate | 在 gate 說明後 | `contract-template.md` |
| BDD Execution Closure | 在 closure 說明後 | `bdd-scenario-template.md` |
| Default Flow | 在 flow 步驟後 | `implementation-plan-template.md` |

#### 2c. [`README.md`](workflow/software-delivery/README.md)

更新「產出格式」章節，加入模板索引表。

### Phase 3: 更新 Routing Registry

#### 3a. [`routing-registry.yaml`](knowledge/runtime/routing-registry.yaml)

在 `route.workflow.software-delivery` 下新增 `templates` 欄位：

```yaml
route.workflow.software-delivery:
  templates:
    change-brief: workflow/software-delivery/templates/change-brief-template.md
    contract: workflow/software-delivery/templates/contract-template.md
    bdd-scenario: workflow/software-delivery/templates/bdd-scenario-template.md
    implementation-plan: workflow/software-delivery/templates/implementation-plan-template.md
    review-report: workflow/software-delivery/templates/review-report-template.md
```

#### 3b. [`development-guidance.md`](knowledge/summaries/development-guidance.md)

在 summary 中新增一行描述模板的存在。

### Phase 4: 驗證

1. 確認每個模板的 markdown 格式正確
2. 確認所有 cross-reference link 可解析
3. 確認 routing registry 的 template path 與實際檔案一致
4. 執行 `linked-updates` 檢查是否有遺漏的更新

---

## 擴展方向（Phase 5–8）

以下 3 個方向來自 [spec-kit](https://github.com/github/spec-kit) 的啟發，補足我們在 **Greenfield 標準化流程**、**Slash Command 模式** 與 **spec→plan→tasks 模板系統** 的缺口。

---

### Phase 5: Greenfield 標準化流程

#### 緣起

spec-kit 的核心流程是 `specify → plan → tasks → implement`，這是一條專為 **從零開始的新專案（Greenfield）** 設計的標準化路徑。我們的 [`workflow/software-delivery/`](workflow/software-delivery/README.md) 雖然涵蓋既有專案的變更管理（Change Intake Gate、Contract Governance、BDD Closure），但缺少一條**從零開始的標準化 onboarding 流程**。

#### 目標

建立 `workflow/greenfield/` 目錄，包裝 spec-kit 的 specify → plan → tasks → implement 流程為一個可重複使用的 skill，讓 agent 在開新專案時可以走這條標準化路徑。

#### 執行項目

##### 5a. 建立 `workflow/greenfield/` 目錄結構

```
workflow/greenfield/
├── README.md              # Entry point：何時使用此 workflow
├── execution-flow.md      # 4 階段流程：specify → plan → tasks → implement
└── templates/
    ├── spec-template.md   # 規格模板（改編自 spec-kit）
    ├── plan-template.md   # 計畫模板（改編自 spec-kit）
    └── tasks-template.md  # 任務模板（改編自 spec-kit）
```

##### 5b. [`README.md`](workflow/greenfield/README.md)

定義何時進入此 workflow：

- **觸發條件**：使用者要求「開新專案」、「從頭建立」、「Greenfield」、「新功能從零開始」
- **不適用**：既有專案的變更、bugfix、refactor（應走 `software-delivery` workflow）
- **與 `software-delivery` 的關係**：Greenfield 流程的產出（spec、plan、tasks）可作為 `software-delivery` 的 Change Intake 輸入

##### 5c. [`execution-flow.md`](workflow/greenfield/execution-flow.md)

4 階段流程：

| 階段 | 名稱 | 輸入 | 輸出 | 對應模板 |
|------|------|------|------|---------|
| 1 | **Specify**（規格定義） | 使用者需求描述 | Feature Specification | `spec-template.md` |
| 2 | **Plan**（技術計畫） | Feature Specification | Implementation Plan | `plan-template.md` |
| 3 | **Tasks**（任務拆解） | Implementation Plan | Task Breakdown | `tasks-template.md` |
| 4 | **Implement**（實作） | Task Breakdown | 實作程式碼 + 測試 | —（引用 `software-delivery` 的 BDD Closure） |

每個階段包含：
- **Entry Condition**：前一階段的產出必須 complete
- **Process**：階段內要做的事
- **Output**：必須產出的文件/程式碼
- **Gate**：必須通過的檢查點才能進入下一階段
- **Template Reference**：引用對應模板

##### 5d. 模板檔案（3 個）

改編自 spec-kit 的 3 個模板，融入我們系統的既有元素（Contract Governance、BDD Closure、Change Intake Gate）：

###### 5d-1. [`spec-template.md`](workflow/greenfield/templates/spec-template.md)

```markdown
# Feature Specification: <feature name>

## Metadata
- **Status**: Draft | Review | Approved
- **Created**: <YYYY-MM-DD>
- **Input Source**: <user description | product brief | issue>

## User Scenarios & Testing

### User Story 1 — <title> (Priority: P1)
**Description**: <plain language description>
**Why P1**: <value justification>
**Independent Test**: <how to verify this story alone>

**Acceptance Scenarios**:
1. **Given** <precondition> **When** <action> **Then** <expected>
2. **Given** <precondition> **When** <action> **Then** <expected>

### User Story 2 — <title> (Priority: P2)
...

### User Story 3 — <title> (Priority: P3)
...

### Edge Cases
- <boundary condition>
- <error scenario>

## Requirements

### Functional Requirements
- **FR-001**: System MUST <capability>
- **FR-002**: System MUST <capability>

### Key Entities
- **<Entity>**: <description, attributes, relationships>

## Success Criteria
- **SC-001**: <measurable outcome>
- **SC-002**: <measurable outcome>

## Assumptions
- <assumption about scope, environment, or dependencies>
```

###### 5d-2. [`plan-template.md`](workflow/greenfield/templates/plan-template.md)

```markdown
# Implementation Plan: <feature>

## Summary
<primary requirement + technical approach>

## Technical Context
- **Language/Version**: <e.g., Python 3.11>
- **Primary Dependencies**: <e.g., FastAPI, SQLAlchemy>
- **Storage**: <e.g., PostgreSQL>
- **Testing**: <e.g., pytest>
- **Target Platform**: <e.g., Linux server>
- **Project Type**: <library | cli | web-service | mobile-app>
- **Performance Goals**: <e.g., 1000 req/s>
- **Constraints**: <e.g., <200ms p95>

## Constitution Check
*GATE: Must pass before proceeding to Tasks phase.*
- [ ] Architecture consistency with existing system
- [ ] Dependency license compatibility
- [ ] Security & compliance requirements

## Project Structure

### Documentation
```
specs/<feature>/
├── spec.md       # This file
├── plan.md       # This file
└── tasks.md      # Next phase output
```

### Source Code
```
<project structure based on project type>
```

## Complexity Tracking
| Violation | Why Needed | Simpler Alternative Rejected |
|-----------|------------|------------------------------|
| <e.g., 4th service> | <current need> | <why 3 insufficient> |
```

###### 5d-3. [`tasks-template.md`](workflow/greenfield/templates/tasks-template.md)

```markdown
# Tasks: <feature name>

## Format
`[ID] [P?] [Story] Description — <file path>`
- **[P]**: Can run in parallel
- **[Story]**: Maps to user story (US1, US2, US3)

## Phase 1: Setup (Shared Infrastructure)
- [ ] T001 Create project structure
- [ ] T002 [P] Initialize project with dependencies
- [ ] T003 [P] Configure linting and formatting

## Phase 2: Foundational (Blocking Prerequisites)
- [ ] T004 Setup database schema and migrations
- [ ] T005 [P] Implement auth framework
- [ ] T006 [P] Setup API routing and middleware

## Phase 3: User Story 1 — <title> (Priority: P1) 🎯 MVP
### Tests (if requested)
- [ ] T010 [P] [US1] Contract test — tests/contract/test_<name>.py
- [ ] T011 [P] [US1] Integration test — tests/integration/test_<name>.py
### Implementation
- [ ] T012 [P] [US1] Create <Entity> model — src/models/<entity>.py
- [ ] T013 [US1] Implement <Service> — src/services/<service>.py
- [ ] T014 [US1] Implement endpoint — src/api/<endpoint>.py

## Phase 4: User Story 2 — <title> (Priority: P2)
...

## Phase N: Polish & Cross-Cutting
- [ ] TXXX [P] Documentation — docs/
- [ ] TXXX Code cleanup and refactoring
- [ ] TXXX Performance optimization

## Dependencies & Execution Order
- Setup → Foundational (BLOCKS all stories)
- Foundational → User Stories (can proceed in parallel or sequentially)
- All stories → Polish

## Implementation Strategy
- **MVP First**: Complete US1 only → validate → deploy
- **Incremental**: Add US2 → US3 one by one
- **Parallel Team**: Multiple devs on different stories after Foundational
```

##### 5e. 更新 Routing Registry

在 [`routing-registry.yaml`](knowledge/runtime/routing-registry.yaml) 新增 `route.workflow.greenfield`：

```yaml
- id: route.workflow.greenfield
  task_intent: 從零開始的新專案標準化流程（specify → plan → tasks → implement）
  primary_source: workflow/greenfield/execution-flow.md
  required_dependencies:
    - workflow/greenfield/README.md
    - enforcement/README.md
  candidate_sources:
    - workflow/greenfield/templates/spec-template.md
    - workflow/greenfield/templates/plan-template.md
    - workflow/greenfield/templates/tasks-template.md
    - workflow/software-delivery/execution-flow.md
  templates:
    spec: workflow/greenfield/templates/spec-template.md
    plan: workflow/greenfield/templates/plan-template.md
    tasks: workflow/greenfield/templates/tasks-template.md
```

---

### Phase 6: Slash Command 模式

#### 緣起

spec-kit 提供 `/speckit.spec`、`/speckit.plan`、`/speckit.tasks` 等 Slash Command，讓使用者直接在對話中觸發特定 workflow 階段。我們的系統目前缺少這種**指令式觸發**機制。

#### 目標

定義 `/skill.*` 指令模式，讓使用者可以直接觸發特定 workflow，無需手動描述需求。

#### 執行項目

##### 6a. 定義指令格式

在 [`enforcement/`](enforcement/README.md) 或 [`ai-tools/`](ai-tools/) 下建立指令定義：

```
/skill.spec     → 觸發 Greenfield Specify 階段（產出 spec）
/skill.plan     → 觸發 Greenfield Plan 階段（產出 plan）
/skill.tasks    → 觸發 Greenfield Tasks 階段（產出 tasks）
/skill.brief    → 觸發 Software Delivery Change Intake（產出 change brief）
/skill.contract → 觸發 Software Delivery Contract Governance（產出 contract）
/skill.bdd      → 觸發 Software Delivery BDD Closure（產出 BDD scenarios）
```

##### 6b. 指令 → Workflow 對應表

| 指令 | 對應 Workflow | 對應階段 | 產出模板 |
|------|---------------|---------|---------|
| `/skill.spec` | `greenfield` | Specify | `spec-template.md` |
| `/skill.plan` | `greenfield` | Plan | `plan-template.md` |
| `/skill.tasks` | `greenfield` | Tasks | `tasks-template.md` |
| `/skill.brief` | `software-delivery` | Change Intake | `change-brief-template.md` |
| `/skill.contract` | `software-delivery` | Contract Governance | `contract-template.md` |
| `/skill.bdd` | `software-delivery` | BDD Closure | `bdd-scenario-template.md` |

##### 6c. 建立指令文件

建立 [`ai-tools/slash-commands.md`](ai-tools/slash-commands.md)：

```markdown
# Slash Commands

## 格式
`/skill.<command> [arguments]`

## 指令列表

### Greenfield Workflow
| 指令 | 功能 | 範例 |
|------|------|------|
| `/skill.spec` | 建立 Feature Specification | `/skill.spec 建立使用者登入功能` |
| `/skill.plan` | 建立 Implementation Plan | `/skill.plan`（需先有 spec） |
| `/skill.tasks` | 建立 Task Breakdown | `/skill.tasks`（需先有 plan） |

### Software Delivery Workflow
| 指令 | 功能 | 範例 |
|------|------|------|
| `/skill.brief` | 建立 Change Brief | `/skill.brief 修復登入頁面 500 錯誤` |
| `/skill.contract` | 建立 Contract | `/skill.contract`（需先有 change brief） |
| `/skill.bdd` | 建立 BDD Scenarios | `/skill.bdd`（需先有 contract） |

## 實作方式
Agent 在收到 `/skill.<command>` 時：
1. 解析 command 名稱，對應到 routing-registry.yaml 的 route
2. 載入對應 workflow 的 execution-flow
3. 載入對應模板
4. 根據 arguments 或 conversation context 填寫模板
5. 輸出結構化文件
```

##### 6d. 更新 Routing Registry

在 [`routing-registry.yaml`](knowledge/runtime/routing-registry.yaml) 為每個 route 加上 `slash_commands` 欄位：

```yaml
route.workflow.greenfield:
  slash_commands:
    - /skill.spec
    - /skill.plan
    - /skill.tasks

route.workflow.software-delivery:
  slash_commands:
    - /skill.brief
    - /skill.contract
    - /skill.bdd
```

---

### Phase 7: spec→plan→tasks 模板系統整合

#### 緣起

Phase 5 建立了 3 個模板檔案，但這些模板需要與我們系統的既有元素（Contract Governance、BDD Closure、Change Intake Gate）深度整合，形成一個**完整的模板生態系**。

#### 目標

建立模板之間的 traceability chain，讓 spec → plan → tasks 的產出可被後續階段自動消費，並與 `software-delivery` workflow 的 5 個模板形成完整閉環。

#### 執行項目

##### 7a. 建立模板 traceability chain

```
spec-template.md  ──→  plan-template.md  ──→  tasks-template.md
      │                      │                       │
      │ (user stories)       │ (architecture)        │ (task IDs)
      ▼                      ▼                       ▼
change-brief-template.md  contract-template.md   bdd-scenario-template.md
      │                      │                       │
      └──────→  implementation-plan-template.md  ←──┘
                              │
                              ▼
                      review-report-template.md
```

##### 7b. 更新 `software-delivery` 模板的 traceability

在每個現有模板的 Metadata 區塊加入 `Traceability` 欄位：

**`change-brief-template.md`** 新增：
```markdown
## Traceability
- **Spec Source**: <link to spec-template.md if from Greenfield>
- **Downstream**: → Contract → BDD Scenarios → Implementation Plan
```

**`contract-template.md`** 新增：
```markdown
## Traceability
- **Upstream**: <link to change-brief-template.md>
- **Downstream**: → BDD Scenarios → Implementation Plan
```

**`bdd-scenario-template.md`** 新增：
```markdown
## Traceability
- **Upstream**: <link to contract-template.md>
- **Downstream**: → Implementation Plan → Review Report
```

**`implementation-plan-template.md`** 新增：
```markdown
## Traceability
- **Upstream**: <link to change-brief | contract | bdd-scenario>
- **Downstream**: → Review Report
```

**`review-report-template.md`** 新增：
```markdown
## Traceability
- **Artifacts Reviewed**: <links to all upstream templates>
```

##### 7c. 更新 `development-guidance.md`

在 summary 中補充 Greenfield workflow 與 Slash Commands 的存在：

```markdown
| Summary | ...提供 5 個標準化輸出模板（change-brief / contract / bdd-scenario / implementation-plan / review-report），位於 `workflow/software-delivery/templates/`。另提供 Greenfield 標準化流程（`workflow/greenfield/`）與 Slash Command 模式（`ai-tools/slash-commands.md`）。 |
```

##### 7d. 更新 `plans/README.md`

在「目前狀態」章節記錄本 plan 的擴展方向。

---

### Phase 8: 驗證（擴展部分）

1. 確認 `workflow/greenfield/` 目錄結構完整
2. 確認 3 個新模板的 markdown 格式正確
3. 確認 `ai-tools/slash-commands.md` 的指令對應正確
4. 確認 `routing-registry.yaml` 的 `route.workflow.greenfield` 可解析
5. 確認所有模板之間的 traceability chain 完整
6. 確認 `development-guidance.md` 已更新
7. 執行 `linked-updates` 檢查是否有遺漏的更新

## 完成條件（擴展）

- [ ] Phase 5: `workflow/greenfield/` 目錄建立（README、execution-flow、3 個模板）
- [ ] Phase 5: `routing-registry.yaml` 新增 `route.workflow.greenfield`
- [ ] Phase 6: `ai-tools/slash-commands.md` 建立
- [ ] Phase 6: `routing-registry.yaml` 的 slash_commands 欄位更新
- [ ] Phase 7: 所有模板的 traceability 欄位更新
- [ ] Phase 7: `development-guidance.md` 更新
- [ ] Phase 8: 驗證通過
- [ ] linked-updates 驗證通過
- [ ] commit & push
