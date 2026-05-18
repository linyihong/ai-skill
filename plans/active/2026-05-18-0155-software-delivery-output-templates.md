# Plan: Software Delivery Output Templates

## 緣起

分析 [github/spec-kit](https://github.com/github/spec-kit) 後發現，spec-kit 的核心價值在於**標準化輸出模板**（spec-template, plan-template, tasks-template），讓每次開發的產出格式一致、可被工具消費。

我們的 [`workflow/software-delivery/`](workflow/software-delivery/README.md) 在流程深度（BDD 閉環、Change Intake Gate、Performance Test Gate、Contract Governance）上已完全涵蓋 spec-kit，唯一缺少的就是**標準化輸出模板**。

本 plan 目標：為 software-delivery workflow 建立 5 個輸出模板，讓每個階段的產出格式一致、可追溯、可被後續階段自動消費。

## 狀態

- **Status**: Draft
- **Author**: Roo (architect → code)
- **Created**: 2026-05-18 01:55 UTC
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

## 完成條件

- [ ] 5 個模板檔案已建立於 `workflow/software-delivery/templates/`
- [ ] `execution-flow.md` 已更新，每個階段引用對應模板
- [ ] `development-process.md` 已更新，每個關卡引用對應模板
- [ ] `README.md` 已更新，產出格式章節包含模板索引
- [ ] `routing-registry.yaml` 已註冊模板路徑
- [ ] `development-guidance.md` 已更新
- [ ] linked-updates 驗證通過
