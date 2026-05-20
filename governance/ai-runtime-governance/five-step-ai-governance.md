# Five-Step AI Governance

## Source Philosophy

source_philosophy:

- [`intelligence/engineering/philosophy/musk-five-step-algorithm.md`](../../intelligence/engineering/philosophy/musk-five-step-algorithm.md)

本文件是 Musk Five-Step Algorithm 的 AI runtime governance 轉譯版。原始思想屬於 `intelligence/`；本文件定義如何把思想套用到 skill、memory、workflow、rule、context、automation 與 runtime surface 的新增決策。

核心原則：

> 不要自動化混亂。先質疑、刪除、簡化與加速，再自動化。

## 觸發時機

在新增或擴張下列項目前，先套用本治理 gate：

- skill
- workflow
- rule
- memory
- context 常駐內容
- prompt artifact
- automation
- runtime guard / compiler surface
- validation scenario
- persistent metadata

## Step 1 — Question Every Requirement

新增前先問：這是否真的必要？

AI runtime gate：

| 問題 | 通過條件 |
| --- | --- |
| Model 是否已經知道這件事？ | 若是通用能力，不新增 durable knowledge。 |
| 這是 operational knowledge 還是 generic knowledge？ | 只有 operational / organization-specific / repo-specific reusable knowledge 才持久化。 |
| 是否提高 reliability、repeatability、cost、orchestration 或 verification？ | 至少命中一項，且有明確 validation signal。 |
| 是否可由更好的 routing 取代？ | 若可，只更新 routing / summary，不新增大文件。 |

若無法回答，先停在 plan / candidate，不直接新增常駐規則。

## Step 2 — Delete Unnecessary Parts

保留前先問：能不能刪？

AI runtime gate：

| 類型 | 刪除 / 降級信號 |
| --- | --- |
| Context | 長期未被 routing 命中、與其他文件重複、只增加背景雜訊。 |
| Memory | 沒有 TTL、沒有 repeated usefulness、無 validation evidence。 |
| Skill / workflow | 被新分層取代、只保存 generic model knowledge、無 active entrypoint。 |
| Summary / index | 指向 stale source、與 generated report 不一致、無 source-of-truth gate。 |

刪除不是第一反應；可先降級為 `candidate`、`deprecated`、cold lookup 或 summary-only。

## Step 3 — Simplify And Optimize

擴張前先簡化。

AI runtime gate：

| 待簡化項 | 轉換方向 |
| --- | --- |
| giant prompt | core bootstrap → routing → summary → selective expansion |
| monolithic workflow | phase / checkpoint / artifact gate |
| recursive agent loop | planner → executor → verifier → finalizer |
| hidden state | explicit state machine / trace / ledger |
| duplicated rules | source rule + metadata mirror + generated report |

不得最佳化尚未穩定、不可觀察、不可重放的流程。

## Step 4 — Accelerate Cycle Time

加速不是跳過 validation，而是降低 feedback loop 成本。

AI runtime gate：

| 成本來源 | 降低方式 |
| --- | --- |
| Token usage | summary-first、TTL、prompt cache alignment、lazy-load rules。 |
| Model cost | small model 做 routing / classification，frontier model 做困難 reasoning。 |
| Tool latency | 先查 index / registry，再展開 canonical source。 |
| Retry loops | 建立 checkpoint、trace、failure class 與 recovery path。 |
| Debug friction | 讓每一步可 replay、可 isolated validation。 |

若加速會跳過 source-of-truth、sanitization、validation 或 user goal，則不通過。

## Step 5 — Automate Last

只有穩定流程可以自動化。

Automation 前必須具備：

- routing 已穩定。
- verification 已存在。
- trace / replay 已存在。
- rollback 或 recovery path 已定義。
- failure class 可分類。
- observability 足以判斷成功/失敗。

禁止自動化：

- hidden-state workflow。
- 無 verifier 的 execution。
- hallucination-prone reasoning。
- 不可重放的手動流程。
- 尚未有 source-of-truth gate 的 context / memory promotion。

## Source Philosophy Mapping

| Philosophy | AI runtime governance |
| --- | --- |
| Question every requirement | Necessity gate before adding durable knowledge or automation |
| Delete unnecessary parts | Context TTL、cold lookup、deprecation、summary-only loading |
| Simplify and optimize | Decompose prompts, workflows, agents, and state |
| Accelerate cycle time | Reduce token/model/tool/retry cost without skipping validation |
| Automate last | Require replay, verifier, observability, rollback/recovery before automation |

## Validation Checklist

- 是否沒有把 generic model knowledge 持久化成 skill？
- 是否先考慮刪除、降級或 summary-only？
- 是否拆分 monolithic workflow / prompt？
- 是否降低 cycle time 而不跳過 validation？
- 是否在 automation 前具備 replay、verifier、observability 與 recovery？

## 與其他層的關係

- `intelligence/engineering/philosophy/` 保存 source philosophy。
- `governance/` 保存 AI 化後的治理轉譯。
- `workflow/` 在任務流程中引用治理 gate。
- `runtime/` 只承接已穩定、可 machine-enforce 的治理結果。
- `validation/` 可把違反本治理的模式變成 stateless scenarios。
