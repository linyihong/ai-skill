# Greenfield Standardization Workflow

## 文件結構

- [`README.md`](README.md) — 本文件，entry point
- [`execution-flow.md`](execution-flow.md) — 4 階段執行流程（Specify → Plan → Tasks → Implement）
- [`execution-flow.yaml`](execution-flow.yaml) — 4 階段 executable contract 與 gate
- [`templates/spec-template.md`](templates/spec-template.md) — Feature Specification 模板
- [`templates/plan-template.md`](templates/plan-template.md) — Implementation Plan 模板
- [`templates/tasks-template.md`](templates/tasks-template.md) — Task Breakdown 模板

## 何時進入此 Workflow

當使用者要求**從零開始的新專案**時，使用此 workflow。典型觸發條件：

- 「開新專案」
- 「從頭建立一個 ...」
- 「Greenfield 專案」
- 「新功能從零開始」
- 「建立一個全新的 ...」

## 不適用

以下情況**不應**使用此 workflow，應走 [`software-delivery`](../software-delivery/README.md) workflow：

- 既有專案的變更（bugfix、refactor、feature addition）
- 既有程式碼庫的分析與文件回填
- 需要 Change Intake Gate 的變更管理

## 與 software-delivery 的關係

Greenfield 流程的產出可作為 `software-delivery` 的 Change Intake 輸入：

```
Greenfield:  Specify → Plan → Tasks → Implement
                │         │       │         │
                ▼         ▼       ▼         ▼
software-delivery:  Change Intake → Contract → BDD → Implementation → Review
```

## 流程概述

| 階段 | 名稱 | 輸入 | 輸出 | 模板 |
|------|------|------|------|------|
| 1 | **Specify**（規格定義） | 使用者需求描述 | Feature Specification | [`spec-template.md`](templates/spec-template.md) |
| 2 | **Plan**（技術計畫） | Feature Specification | Implementation Plan | [`plan-template.md`](templates/plan-template.md) |
| 3 | **Tasks**（任務拆解） | Implementation Plan | Task Breakdown | [`tasks-template.md`](templates/tasks-template.md) |
| 4 | **Implement**（實作） | Task Breakdown | 實作程式碼 + 測試 | 引用 `software-delivery` 的 BDD Closure |

## 核心原則

1. **Spec-first**：先定義規格（user stories + acceptance scenarios），再談技術
2. **Independent stories**：每個 user story 必須可獨立測試、獨立交付
3. **Gate-based progression**：每個階段有 entry condition 與 exit gate，未通過不得進入下一階段
4. **Template-driven output**：每個階段的產出使用標準化模板，確保格式一致、可被工具消費

## 遷移狀態

- **Status**: `draft`
- **Created**: 2026-05-18
- **Source**: Adapted from [github/spec-kit](https://github.com/github/spec-kit) specify → plan → tasks → implement pipeline
