# Intelligence Extraction

`feedback/extraction/` 定義「智慧抽取」的系統設計。本目錄保存如何從 feedback lesson、replay 結果、review comment 與實作經驗中，提煉出可重複使用的工程智慧（intelligence atom），讓一次性觀察能被結構化為長期可用的判斷力。

## 核心責任

- 從 raw feedback 到 intelligence atom 的抽取流程。
- Extraction 的品質門檻（什麼樣的 lesson 值得變成 intelligence）。
- Intelligence atom 的結構化模板與必備欄位。
- Extraction 結果的驗證與 promotion 準備。
- 與 `intelligence/` 各子目錄的對應關係。

## 核心原則

1. **Extraction 不是搬運**。不是把 feedback lesson 複製到 intelligence/，而是提煉出可泛化的工程判斷。
2. **Extraction 必須去敏**。包含專案名稱、客戶名稱、raw evidence 的內容不能進入 intelligence。
3. **Extraction 必須有邊界條件**。每個 intelligence atom 必須說明「何時適用」與「何時不適用」，否則只是 opinion 不是 intelligence。
4. **Extraction 是選擇性的**。不是每個 feedback lesson 都需要 extraction。只有跨專案、跨 skill 或高重複價值的 lesson 才值得。

## Extraction 門檻

符合以下任一條件時，應考慮 extraction：

| 條件 | 說明 |
| --- | --- |
| 同一 lesson 在 ≥2 個 skill 或專案中出現 | 表示是跨領域的通用智慧。 |
| Lesson 涉及 trade-off 或架構判斷 | 這類知識最適合放在 intelligence/。 |
| Lesson 可節省大量 token 或時間 | 如果 agent 每次都要重新摸索，值得 extraction。 |
| Lesson 是 failure prevention | 預防性知識比修復性知識更有價值。 |
| Lesson 來自資深工程師的 code review | Review comment 通常包含高密度的工程判斷。 |

不符合以下任一條件時，不應 extraction：

- Lesson 只適用於單一專案的特定實作細節。
- Lesson 無法泛化（「這個 API 的參數是 X」不是 intelligence）。
- Lesson 沒有明確的邊界條件（「永遠不要用微服務」不是 intelligence）。

## Extraction 流程

```
1. 識別 extraction 候選
   ├─ 來自 replay 結果的 promotion_candidate
   ├─ 來自 feedback_history 中標記為 high-value 的 lesson
   ├─ 來自 code review 或架構審查的 comment
   └─ 來自 production incident 的根因分析

2. 判斷 intelligence 類型
   ├─ Heuristic（經驗法則）→ intelligence/engineering/heuristics/
   ├─ Tradeoff（取捨）→ intelligence/engineering/tradeoffs/
   ├─ Failure pattern（失敗模式）→ intelligence/engineering/failure/
   ├─ Architecture pattern（架構判斷）→ intelligence/engineering/architecture/
   ├─ Domain pattern（領域知識）→ intelligence/engineering/domain/
   ├─ Anti-pattern（錯誤設計）→ intelligence/engineering/anti-patterns/
   ├─ Distributed systems（分散式系統）→ intelligence/engineering/distributed-systems/
   ├─ Business decision（商業判斷）→ intelligence/business/
   └─ Domain intelligence（領域經驗）→ intelligence/travel/

3. 撰寫 intelligence atom
   ├─ 使用 intelligence/ 對應子目錄的格式
   ├─ 包含：原則、為什麼、何時適用、何時不適用
   ├─ 包含：決策流程或 decision tree（如適用）
   ├─ 包含：Token Impact 評估
   └─ 不包含：專案特定 raw evidence、客戶名稱、未去敏數據

4. 設定 lifecycle state
   ├─ candidate-intelligence（預設）
   ├─ validated-intelligence（經過至少一次實戰驗證）
   └─ promoted-intelligence（已進入正式 routing）

5. 執行 linked updates
   ├─ 更新對應 intelligence/ 子目錄的 README（目前 atoms 表格）
   ├─ 更新 knowledge/indexes/README.md（路由）
   ├─ 更新 knowledge/runtime/routing-registry.yaml
   ├─ 更新 knowledge/summaries/（如需要）
   └─ 更新 knowledge/graphs/（如需要）
```

## Intelligence Atom 必備欄位

每個從 extraction 產出的 intelligence atom 應包含：

| 欄位 | 說明 | 範例 |
| --- | --- | --- |
| 原則 | 一句話描述核心判斷 | "If abstraction removes more clarity than duplication, do not abstract." |
| 為什麼 | 這個原則背後的理由 | Abstraction 增加 indirection，如果好處只有消除重複但降低可讀性，淨效益為負。 |
| 何時適用 | 這個原則適用的情境 | 當你正在考慮引入新抽象層時。 |
| 何時不適用 | 這個原則不適用的情境 | 當重複次數 > 3 且每次修改都需要同步變更時。 |
| 決策流程 | 結構化的判斷步驟 | Decision tree 或 checklist。 |
| Token Impact | 使用這個 atom 對 token 的影響 | 約 200-400 tokens。 |

## 與其他層的關係

- `feedback/replay/`：Replay 結果是 extraction 的主要輸入來源之一。
- `feedback/refinement/`：Extraction 可能發現 workflow 需要調整，進入 refinement 流程。
- `feedback/promotion/`：Extraction 完成的 intelligence atom 需通過 promotion pipeline 才能正式上線。
- `intelligence/`：Extraction 的最終目的地。各子目錄的 README 定義了 atom 格式與 scope。
- `governance/lifecycle/`：Intelligence atom 的 lifecycle state 由此管理。
- `governance/validation/`：Extraction 結果需通過 validation gates。
- `knowledge/indexes/README.md`：Extraction 完成後需更新 index 讓 agent 可路由到此 atom。
