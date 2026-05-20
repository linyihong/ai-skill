# Musk Five-Step Algorithm

## 核心思想

Musk Five-Step Algorithm 是一個 complexity reduction framework。它的核心不是「做更多」，而是在工程系統擴張前反覆確認：

```text
Question -> Delete -> Simplify -> Accelerate -> Automate
```

最重要的警告是：

> Do not automate chaos.

自動化會放大效率，也會放大混亂。若 requirement 不清、流程不穩、驗證不存在，自動化只會更快地製造錯誤。

## Step 1 — Question Every Requirement

先質疑需求本身是否必要。許多系統膨脹不是因為缺少能力，而是把不必要的要求包裝成「必須」。

對 AI knowledge system 的啟發：

- 不要因為模型能讀更多 context，就把所有知識外部化。
- 不要新增 generic skill 取代 frontier model 已經知道的通用能力。
- 只有能提高 reliability、repeatability、cost control、orchestration 或 verification 的知識，才值得持久化。

## Step 2 — Delete Unnecessary Parts

刪除未使用、重複、過期、低信心或低價值的部件。Context、memory、rules、workflow 和 automation 都有維護成本。

對 AI knowledge system 的啟發：

- Context 是昂貴記憶，不是免費背景。
- Stale memory 會稀釋推理並增加 hallucination。
- 重複 summaries、dead skills、obsolete prompts 應該被清理或降級。

## Step 3 — Simplify And Optimize

先降低系統複雜度，再最佳化。不要最佳化一個不穩定、不可觀察、不可重放的流程。

對 AI knowledge system 的啟發：

- 巨型 prompt 應拆成 routing、summary、selective expansion。
- Monolithic agent 應拆成 planner、executor、verifier、finalizer。
- Hidden state 應轉成明確 state machine 或 trace。

## Step 4 — Accelerate Cycle Time

降低 feedback loop 的成本：token、latency、retry、debugging friction。快速可靠的迭代通常勝過一次超大型 reasoning chain。

對 AI knowledge system 的啟發：

- 使用 summary-first retrieval。
- 低成本模型可處理 routing、tagging、classification。
- Frontier model 應保留給高難度 reasoning。
- 避免每次 full-repo 或 full-context analysis。

## Step 5 — Automate Last

只有穩定、可觀察、可驗證、可回放、可恢復的流程才值得自動化。

對 AI knowledge system 的啟發：

- Automation 前要有 verifier。
- Automation 前要有 replay / trace。
- Automation 前要有 rollback 或 recovery。
- Persistent memory 應該在重複有用、驗證成功、reuse pattern 穩定後才 promotion。

## 適用邊界

適用於：

- 新增 skill、workflow、rule、memory、runtime surface 前。
- 大型 context / prompt 常駐化前。
- 將手動流程升級成 automation 前。
- 拆分或簡化 AI infrastructure 前。

不適用於：

- 需要立刻處理的 P0 safety / secrets / data loss 事件。
- 已有明確 compliance requirement 的必要流程。
- 使用者明確要求保留的 domain-specific operational requirement。

## Governance Translation

AI runtime 的可執行轉譯見：

- [`governance/ai-runtime-governance/five-step-ai-governance.md`](../../../governance/ai-runtime-governance/five-step-ai-governance.md)
