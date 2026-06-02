---
id: 2026-01-01-0002-example-spike
plan_kind: spike
status: draft
owner: example-owner
created: 2026-01-01
parent: 2026-01-01-0000-example-main-plan
required_for_completion: false
sub_plan_reason: >
  範例 spike：PoC 試做 X 演算法，時限 1 工作天，結果無論成敗都回寫主計畫
  對應 phase decision。不阻擋 parent archive（required_for_completion: false）。
---

# Example Spike（fixture）

> **Fixture 用途**：示範 spike 模板的最小合法 frontmatter。
> Spike 模板極簡：只需 Goal / Acceptance / 結果回寫，免 Phase 0 公版。

## Goal

時限內試做 `<X>`，回答「`<Y>` 路線是否可行」。

## Time-box

≤ 1 工作天（超時即視為 inconclusive，回寫主計畫並關閉本 spike）。

## Acceptance（結果回寫條件）

- [ ] 試做完成或時限到
- [ ] 結果（pass / fail / inconclusive）已回寫主計畫對應 phase
- [ ] 本 spike `status: completed`，可安全 archive（不阻擋 parent）

## 結果回寫位置

主計畫 `2026-01-01-0000-example-main-plan` §Phase `<N>` 的決策段落。
