---
id: 2026-01-01-0001-example-sub-plan
plan_kind: sub
status: draft
owner: example-owner
created: 2026-01-01
parent: 2026-01-01-0000-example-main-plan
required_for_completion: true
sub_plan_reason: >
  範例 sub-plan：示範必填欄位 parent / required_for_completion / sub_plan_reason
  三件套。實際使用時請寫具體拆分理由（例如「DSL schema 設計獨立於 executor wiring，
  需要獨立 stakeholder sign-off」）。
---

# Example Sub-Plan（fixture）

> **Fixture 用途**：示範 sub-plan 的最小合法 frontmatter。Sub-plan 不需要
> §Decision Rationale，繼承 parent；只需要 Purpose / Acceptance / Runtime Impact。

## Purpose

子計畫負責回答「How」 — 具體實作步驟、acceptance criteria、技術細節。

## Acceptance Criteria

- [ ] 步驟 1 完成 + validation
- [ ] 步驟 2 完成 + validation

## Runtime Impact

無（純文件交付），或：
- 新增 `route.<domain>.<name>` 到 routing-registry
- 新增 generated_surface `<key>`，consumer = `<validator/hook/CLI>`
