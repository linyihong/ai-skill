# ADR-002: Intelligence vs Knowledge Separation

## Status

**Accepted**

## Context

在設計新分層時，需要決定「工程智慧」與「百科知識」的邊界。最初的想法是將所有內容放在 `knowledge/` 下，但隨著 `intelligence/` 概念的出現，需要明確區分兩者的責任。

核心問題：

1. **知識（knowledge）** 是「事實性資訊」：Redis 支援 pub/sub、CQRS 是什麼、HTTP 狀態碼意義。
2. **技能（skills）** 是「可執行流程」：如何分析 APK、如何做 code review、如何規劃旅行。
3. **智慧（intelligence）** 是「經驗判斷」：什麼時候該用 CQRS、什麼跡象表示 Redis 連線池有問題、如何取捨效能與可維護性。

如果將 intelligence 放在 knowledge 下，會導致：
- Agent 無法區分「事實查詢」與「經驗判斷」。
- 路由決策模糊：同一條路徑同時承載百科與 heuristic。
- 難以設計不同的 lifecycle：knowledge 可穩定不變，intelligence 需要持續演化。

## Decision

將 `intelligence/` 獨立於 `knowledge/`，兩者為平行層。

### 責任邊界

| 維度 | `knowledge/` | `intelligence/` |
| --- | --- | --- |
| 內容類型 | 事實、定義、規範 | 經驗法則、取捨、判斷框架 |
| 變動頻率 | 低（除非技術本身改變） | 中高（隨經驗累積演化） |
| 使用方式 | 查詢、驗證、參考 | 判斷、決策、取捨 |
| 生命週期 | stable → deprecated | candidate → validated → promoted → updated |
| Token 成本 | 可壓縮為 summary | 需保留完整決策邏輯 |
| 範例 | Redis pub/sub 文件 | 何時不該用 Redis pub/sub |

### 與 `skills/` 的關係

- `skills/` 是執行流程，引用 `knowledge/` 的事實與 `intelligence/` 的判斷。
- 一個 skill 可以同時依賴 knowledge atom 與 intelligence atom。

## Consequences

### 正面

- **路由清晰**：agent 可以根據任務類型（查事實 vs 做判斷）選擇正確的層。
- **Lifecycle 差異化**：knowledge 可長期穩定，intelligence 可持續演化與版本管理。
- **Token 優化**：intelligence atom 保留完整決策邏輯（不可過度壓縮），knowledge 可壓縮為 summary。
- **與人類認知一致**：人類也是分開儲存「事實知識」與「經驗智慧」。

### 負面

- **多一層管理成本**：需要維護兩個平行層的索引與 navigation。
- **邊界模糊風險**：某些內容（如 failure pattern）同時有 knowledge 面向（什麼是 connection leak）與 intelligence 面向（如何判斷 connection leak）。需要 governance 規則處理 hybrid cases。
- **初期內容稀疏**：`intelligence/` 在初期 atoms 較少，agent 可能找不到需要的 intelligence。

## Alternatives Considered

- **Intelligence 放在 knowledge 下**：簡化目錄結構，但路由模糊。不採用。
- **Intelligence 放在 skills 下**：與 skill 綁定，無法跨 skill 重用。不採用。
- **Intelligence 作為 metadata 標籤**：不建立獨立目錄，只在 knowledge atom 上加 `type: intelligence` 標籤。但 lifecycle 與 routing 無法差異化。不採用。

## Related

- [`intelligence/README.md`](../intelligence/README.md) — intelligence 層定義
- [`knowledge/README.md`](../knowledge/README.md) — knowledge 層定義
- [`plans/archived/2026-05-11-next-stage-upgrade-plan.md`](../plans/archived/2026-05-11-next-stage-upgrade-plan.md) — 分層架構規劃
- [`governance/lifecycle/README.md`](../governance/lifecycle/README.md) — lifecycle 差異化管理
