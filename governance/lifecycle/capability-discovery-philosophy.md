# Capability Discovery Philosophy

## Problem

現有系統的 lazy-loading 機制（[`activation-engine.rb`](../../runtime/router/activation-engine.rb)）非常有效——它確保 agent 只載入需要的知識。但這也帶來一個根本問題：

> **Agent 不知道自己不知道什麼。**

當一個 workflow、intelligence atom 或 validation rule 沒有被任何 activation condition 觸發時，agent 就不會知道它的存在。這就是 **Capability Discovery Problem**。

## Solution: Discovery Checkpoints

Discovery Runtime 在每個 phase 的關鍵時機插入**主動探索檢查點**，讓 agent 在執行前搜尋可能相關但尚未載入的能力。

### 與 Activation Engine 的互補關係

| 面向 | Activation Engine（被動） | Discovery Runtime（主動） |
|------|--------------------------|--------------------------|
| 觸發方式 | 根據已知條件（file_change、task_intent、user_signal） | 根據 phase 時機（before_execution、before_validation） |
| 搜尋範圍 | 已知的 activation rules（14 條） | 未知的能力（workflow、intelligence、validation） |
| 載入策略 | 精準載入（精確匹配條件） | 廣泛搜尋（探索可能相關） |
| 失敗處理 | block 或 warn | warn + continue（不阻擋流程） |
| 知識來源 | activation-rules.yaml | knowledge indexes、graphs、routing registry |

### 與 Knowledge Graphs 的整合

Discovery 使用 [`knowledge/graphs/`](../../knowledge/graphs/README.md) 的 edge relationships 作為搜尋索引：

- `depends_on`：發現與當前任務相依的知識
- `related_to`：發現與當前任務相關的知識
- `routes_to`：發現可路由的知識路徑

當 graph 中缺少 edge 時，discovery 會 fallback 到 directory scan 或 grep search。

## Runtime Surface

| 元件 | 整合方式 |
|------|---------|
| [`phase-machine.yaml`](../../runtime/phases/phase-machine.yaml) | checkpoint phase 新增 `run_capability_discovery` allowed_action |
| [`obligation-ledger.yaml`](../../runtime/obligations/obligation-ledger.yaml) | 新增 `obligation.checkpoint.run_capability_discovery` |
| [`blocking-gates.yaml`](../../runtime/gates/blocking-gates.yaml) | 新增 `gate.checkpoint.capability_discovery_completed`（medium severity） |
| [`capability-checkpoints.yaml`](../../runtime/discovery/capability-checkpoints.yaml) | Phase-aware discovery checkpoint definitions |

## 與既有文件的關係

- [`runtime/discovery/`](../../runtime/discovery/) — Runtime navigation entry point (data file: `capability-checkpoints.yaml`)
- [`runtime/discovery/capability-checkpoints.yaml`](../../runtime/discovery/capability-checkpoints.yaml) — Structured checkpoint definitions
- [`runtime/router/activation-engine.rb`](../../runtime/router/activation-engine.rb) — 被動 lazy-loading，與 discovery 互補
- [`knowledge/graphs/README.md`](../../knowledge/graphs/README.md) — Graph edges 作為 discovery 的搜尋索引
- [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md) — Task intent routing table 作為 discovery 的快速路徑
