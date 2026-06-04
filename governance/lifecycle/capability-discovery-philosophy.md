# Capability Discovery Philosophy

## Problem

現有系統的 contract-backed activation（owner-layer executable YAML contracts + [`routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml)）非常有效——它確保 agent 只載入需要的知識。但這也帶來一個根本問題：

> **Agent 不知道自己不知道什麼。**

當一個 workflow、intelligence atom 或 validation rule 沒有被任何 activation condition 觸發時，agent 就不會知道它的存在。這就是 **Capability Discovery Problem**。

## Solution: Discovery Checkpoints

Discovery Runtime 在每個 phase 的關鍵時機插入**主動探索檢查點**，讓 agent 在執行前搜尋可能相關但尚未載入的能力。

### 與 Contract Activation 的互補關係

| 面向 | Contract Activation（被動） | Discovery Runtime（主動） |
|------|--------------------------|--------------------------|
| 觸發方式 | 根據已知條件（file_change、task_intent、user_signal） | 根據 phase 時機（before_execution、before_validation） |
| 搜尋範圍 | 已知的 executable contract activation | 未知的能力（workflow、intelligence、validation） |
| 載入策略 | 精準載入（精確匹配條件） | 廣泛搜尋（探索可能相關） |
| 失敗處理 | block 或 warn | warn + continue（不阻擋流程） |
| 知識來源 | owner-layer executable YAML contracts / generated_surfaces | knowledge indexes、graphs、routing registry |

### 與 Knowledge Graphs 的整合

Discovery 使用 [`knowledge/graphs/`](../../knowledge/graphs/README.md) 的 edge relationships 作為搜尋索引：

- `depends_on`：發現與當前任務相依的知識
- `related_to`：發現與當前任務相關的知識
- `routes_to`：發現可路由的知識路徑

當 graph 中缺少 edge 時，discovery 會 fallback 到 directory scan 或 grep search。

## Runtime Surface

| 元件 | 整合方式 |
|------|---------|
| [`runtime/runtime.db`](../../runtime/runtime.db) / [`runtime/runtime.db`](../../runtime/runtime.db) | checkpoint phase、`run_capability_discovery` allowed_action、obligation 與 gate 的 compiled source |
| `capability_checkpoints` table | Phase-aware discovery checkpoint definitions |

## Discovery → Detector Feedback Loop

Detector 與 Discovery 不是互斥替代，而是 **known vs unknown** 的分工，且在
detector miss 時形成回饋環（plan 1900 Phase 6.1）：

| | Detector | Capability Discovery |
|---|---|---|
| 處理 | known route 的 known trigger | unknown capability（registry 沒有的 route） |
| 成本 | cheap、deterministic、per-task | expensive、exploratory |
| 觸發 | 每個 substantive task | **只在 detector miss 時** fire |
| 輸出 | `RuntimeContext.ActiveRoute` | route-candidate proposal（Registry growth） |

```text
User Request → Detector
  ├─ hit  → RuntimeContext.ActiveRoute → Execution
  └─ miss → Capability Discovery（graph traversal 找相關 atom）
              → 提案新 route candidate（route-candidate-proposals.yaml）
              → occurrence 累積到門檻 → user / governance review → promote 成正式 route
```

**為什麼這樣分工**：detector miss 不代表「沒有 workflow 可用」，可能是 Registry
還沒收錄此 task type。Discovery 跑 graph traversal 找出可能相關的 capability
atom；若多次一致指向同一群 capability → 暗示應新增一個 route。這形成 **Registry
自我成長機制**：使用越多、coverage 越廣，而非靠人預先窮舉所有 route。

**反濫用**：proposals 採 occurrence-tracking schema（`occurrence_count` /
`first_seen` / `last_seen` / `status`），只有累積到門檻（`occurrence_count >= 5`
且近 30 天活躍）才升為 `ready_for_review`，避免一次性需求污染 Registry。狀態機
與 `ai-skill router proposals` CLI 見 plan 1900 Phase 6.1
（`runtime/router/route-candidate-proposals.yaml`）。

**Sunset 綁定**：本 philosophy 在 `enforcement/enforcement-registry.yaml`
§`capability_discovery` 的 `behavioral_only` 覆蓋，其 revisit 條件即
「workflow_activation Phase 6.1 lands」—— Discovery 的 mechanical 整合點正是
detector 的 miss path，而非把 Discovery 獨立做成 per-turn executor。

## 與既有文件的關係

- [`runtime/runtime.db`](../../runtime/runtime.db) — Compiled discovery checkpoint runtime surface
- [`runtime/runtime.db`](../../runtime/runtime.db) — Source for discovery checkpoints
- [`runtime/runtime.db`](../../runtime/runtime.db) — 被動 lazy-loading，與 discovery 互補
- [`knowledge/graphs/README.md`](../../knowledge/graphs/README.md) — Graph edges 作為 discovery 的搜尋索引
- [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md) — Task intent routing table 作為 discovery 的快速路徑
