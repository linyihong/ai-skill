# Core Bootstrap

本文件是 AI-native Knowledge Operating System 的**最小必讀啟動集合**。每個 session 啟動時，Agent 只需載入這 3 條核心規則，其餘規則依任務 lazy-load。

## 必讀規則

| 順序 | 規則 | 用途 | 預估 tokens |
| --- | --- | --- | --- |
| 1 | [`shared-rules/rule-weight.md`](shared-rules/rule-weight.md) | 規則衝突時如何判斷優先序（P0/P1/P2/P3） | ~300 |
| 2 | [`shared-rules/dependency-reading.md`](shared-rules/dependency-reading.md) | 依賴文件讀取鐵則與 Ai-skill writeback transaction | ~400 |
| 3 | [`shared-rules/conversation-goal-ledger.md`](shared-rules/conversation-goal-ledger.md) | 對話目標閉環與 `.agent-goals/` 使用方式 | ~100 |

**總計：~800 tokens**

## 啟動流程

```
1. 讀取 CORE_BOOTSTRAP.md（本檔）
2. 讀取 README.md（超短入口，了解 OS layout）
3. 依任務 intent 查詢 skills-index.yaml 找到對應 skill
4. 檢查該 skill 是否有 primary_entrypoint 欄位：
   - 有 → 優先讀 primary_entrypoint 指向的新分層路徑
   - 無 → 讀 entrypoint 指向的舊路徑（向後相容）
5. 依 activation rules 決定哪些 lazy-load rules 需要載入
6. 先讀 knowledge/summaries/ 對應 summary（300-500 tokens）
7. 需要時才展開完整 source
```

## 與舊 Default Bootstrap 的關係

舊 Default Bootstrap（12 條規則）已拆分為：

- **Core Bootstrap**（本檔）：3 條必讀規則，每個 session 載入
- **Lazy-load rules**：9 條規則，只在特定情境 activate（詳見 [`shared-rules/README.md`](shared-rules/README.md) 的 activation model）

## 不變的原則

- `reference-first` 仍是預設：agent 直接讀本 repository
- `rule-weight` 的 P0/P1/P2/P3 權重體系不變
- `dependency-reading` 的 dependency read ledger 與 writeback transaction 不變
- `conversation-goal-ledger` 的 `.agent-goals/` 使用方式不變

## 驗證

- Core Bootstrap 三條規則已讀
- 任務 intent 已對應到 skills-index.yaml
- Lazy-load rules 的 activation conditions 已檢查
