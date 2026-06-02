---
id: 2026-01-01-0000-example-main-plan
plan_kind: main
status: draft
owner: example-owner
created: 2026-01-01
parent: null
---

# Example Main Plan（fixture）

> **Fixture 用途**：示範 main plan 的最小合法 frontmatter。Phase 2 validator
> 將以此作為 happy-path testdata。請勿在實際 plan 中直接複用此 `id`。

## Purpose

主計畫負責回答「Why」 — 拆分理由、業務目標、acceptance criteria 摘要、
sub-plan 驗證要點。詳細 How 留給 sub-plan。

## Sub-plan 驗證要點表（範例）

| Sub-plan | 完成條件摘要 | required_for_completion | 驗證方式 |
|---|---|---|---|
| `01-foo` | 範例 sub 完成條件 | true | unit test pass |

## Acceptance

- [ ] 全部 `required_for_completion: true` 的 sub-plan `status: completed`
- [ ] 主計畫 Why / Decision Rationale 已關閉
