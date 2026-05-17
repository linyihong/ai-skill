# Capability Discovery Runtime

`runtime/discovery/` 定義 AI agent 的**能力探索機制**，解決 **Capability Discovery Problem**：lazy loading 本身不會產生探索意識，agent 需要主動搜尋「可能不知道的能力」。

## 問題

現有系統的 lazy-loading 機制（[`activation-engine.rb`](../router/activation-engine.rb)）非常有效——它確保 agent 只載入需要的知識。但這也帶來一個根本問題：

> **Agent 不知道自己不知道什麼。**

當一個 workflow、intelligence atom 或 validation rule 沒有被任何 activation condition 觸發時，agent 就不會知道它的存在。這就是 **Capability Discovery Problem**。

## 解決方案：Discovery Checkpoints

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

## 檔案結構

```
runtime/discovery/
├── README.md                    # 本文件：能力探索機制的說明
└── capability-checkpoints.yaml  # Phase-aware discovery checkpoints 定義
```

## 整合點

Discovery Runtime 整合到以下系統元件：

| 元件 | 整合方式 |
|------|---------|
| [`phase-machine.yaml`](../phases/phase-machine.yaml) | checkpoint phase 新增 `run_capability_discovery` allowed_action |
| [`obligation-ledger.yaml`](../obligations/obligation-ledger.yaml) | 新增 `obligation.checkpoint.run_capability_discovery` |
| [`blocking-gates.yaml`](../gates/blocking-gates.yaml) | 新增 `gate.checkpoint.capability_discovery_completed`（medium severity） |

## 使用方式

Agent 在 checkpoint phase 中：

1. 讀取 [`capability-checkpoints.yaml`](capability-checkpoints.yaml) 了解本輪需要探索的能力類型
2. 根據任務意圖，從 `knowledge/indexes/README.md` 搜尋相關 workflow
3. 從 `knowledge/graphs/` 搜尋相關 intelligence atoms
4. 從 `enforcement/failure-patterns/` 搜尋相關 failure patterns
5. 從 `governance/lifecycle/` 搜尋相關治理規則
6. 標記 `capability_discovery_completed = true`

## 誰會參考這裡（Inbound References）

- [`gate.checkpoint.capability_discovery_completed`](../gates/blocking-gates.yaml) — checkpoint phase 的 blocking gate
- [`obligation.checkpoint.run_capability_discovery`](../obligations/obligation-ledger.yaml) — checkpoint phase 的 obligation
- [`phase.checkpoint`](../phases/phase-machine.yaml) — checkpoint phase 的 allowed_actions

## 與既有層的關係

- `runtime/router/activation-engine.rb` — 被動 lazy-loading，與 discovery 互補
- `knowledge/graphs/` — 提供 graph edges 作為 discovery 的搜尋索引
- `knowledge/indexes/` — 提供 task intent routing table 作為 discovery 的快速路徑
- `knowledge/runtime/routing-registry.yaml` — 提供 routing records 作為 discovery 的完整路徑
- `enforcement/failure-patterns/` — 提供已知失效模式供 discovery 搜尋
