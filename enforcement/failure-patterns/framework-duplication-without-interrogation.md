# Framework Duplication Without Interrogation

Status: validated
Class: `source-of-truth-duplication` / `requirements-cognition-gap`

## Trigger

當 agent 要修改 framework、governance、runtime、workflow、metadata、validation、compiler、generated artifact 或 tool adapter，而且同一語意可能同時存在於 owner source、runtime projection、mirror、generated report 或 compatibility surface 時，懷疑此失效。

## Failure Mode

Agent 沒有先做 pre-build interrogation，就直接設計或實作 framework migration。結果可能留下兩份 rule body、兩條 activation path、兩個 source-of-truth、stale generated output，或把 projection 誤當 canonical source。

## Risk

- 後續 agent 可能讀到舊 surface，和新 contract 做出相反決策。
- Runtime validation 只證明新 surface 存在，卻沒有證明舊 duplicate 已刪除或降級。
- 使用者以為規則已收斂，但 repo 實際上仍維護雙寫語意。

## Required Agent Action

1. 先讀 [`workflow/software-delivery/requirements/pre-build-interrogation.md`](../../workflow/software-delivery/requirements/pre-build-interrogation.md)。
2. 在 plan 或回覆中列出 canonical source、owner layer、projection、mirror / cache / generated output、compiler impact 與 linked updates。
3. 對每個 duplicate surface 做決策：remove、deprecate、explicit precedence、compatibility retained。
4. 在驗證中證明舊 duplicate 已移除、降級或不再具有 authority。

## Prevention Gate

Framework 改動不得直接進 implementation。必須先通過 pre-build interrogation 的 framework discovery 與 duplication risk gates；若有 blocker question，先問使用者。

## 驗證

- `validation/scenarios/failure-derived/framework-duplication-without-interrogation-v1.yaml`
- `validation/scenarios/failure-derived/plan-without-prebuild-interrogation-v1.yaml`
- Runtime / generated surface 變更時，執行 compile / refresh / validate，並查證新 surface synced、舊 duplicate 已移除或非 authoritative。

## Linked Rules

- [`failure-learning-system.md`](../failure-learning-system.md)
- [`dependency-reading.md`](../dependency-reading.md)
- [`linked-updates.md`](../linked-updates.md)
- [`governance/lifecycle/executable-contract-boundary.md`](../../governance/lifecycle/executable-contract-boundary.md)
- [`plans/README.md`](../../plans/README.md)
