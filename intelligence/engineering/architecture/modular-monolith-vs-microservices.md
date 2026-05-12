# Modular Monolith vs Microservices（模組化單體 vs 微服務）

**Status**: `candidate-intelligence`
**Source**: 通用軟體架構經驗

## 原則

**Microservices increase operational complexity rapidly. Prefer modular monolith unless independent scaling, separate deployment cadence, or team autonomy is required.**

微服務的營運複雜度增長極快。除非需要獨立擴展、獨立部署節奏或團隊自治，否則優先選擇模組化單體。

## 為什麼

1. **營運成本被低估**：微服務需要 CI/CD pipeline、service mesh、observability infrastructure、分散式 tracing、container orchestration。這些在 3-5 個 service 時還可控，超過 10 個後急劇增長。
2. **網路延遲是真實成本**：原本的 in-process method call 變成跨網路 RPC，latency 從微秒級變成毫秒級，且引入 partial failure、retry、timeout 等複雜度。
3. **資料一致性複雜度**：單體可以用 ACID transaction 保證一致性；微服務需要 saga、eventual consistency、compensating transaction。
4. **除錯難度**：跨 service 的 bug 需要追蹤多個 service 的 log 與 trace，開發者體驗顯著下降。

## 何時適用 Microservices

- **獨立擴展需求**：不同模組的負載模式差異極大（例如：寫入服務需要 10 個 instance，讀取服務需要 100 個）。
- **獨立部署節奏**：不同模組的發布頻率不同（例如：API 層每週發布，資料處理層每月發布）。
- **團隊自治**：多個團隊各自負責不同模組，需要獨立的開發與部署流程。
- **技術棧異質性**：不同模組需要使用不同的語言或資料庫。

## 何時不適用 Microservices

- **小型專案（< 5 個開發者）**：營運成本超過收益。
- **CRUD-heavy 系統**：大部分邏輯是簡單的資料增刪改查，微服務只增加複雜度。
- **啟動階段**：產品尚未找到 product-market fit，快速迭代比架構彈性更重要。
- **缺乏 DevOps 成熟度**：沒有 CI/CD、container orchestration、observability 基礎設施。

## 決策流程

```text
需要拆分架構？
  ├── 團隊人數 > 5？
  │     ├── 否 → 模組化單體
  │     └── 是 → 繼續評估
  ├── 有獨立擴展需求？
  │     ├── 否 → 模組化單體
  │     └── 是 → 繼續評估
  ├── 有獨立部署節奏需求？
  │     ├── 否 → 模組化單體
  │     └── 是 → 繼續評估
  ├── 有 DevOps 基礎設施？
  │     ├── 否 → 先建立基礎設施，再考慮微服務
  │     └── 是 → 考慮微服務，但從最少 service 開始（2-3 個）
  └── 以上皆符合？
        ├── 是 → 微服務（但從 bounded context 開始，不要過度拆分）
        └── 部分符合 → 模組化單體 + 預留拆分介面
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「用微服務以後 scale 方便」 | Scale 問題通常可以先用垂直擴展或模組化單體解決。微服務是 scale 到極限後的手段，不是預防性措施 |
| 「微服務可以獨立技術棧」 | 每個 service 用不同語言增加維護成本。除非有明確理由，否則同一組織內應統一技術棧 |
| 「先拆 microservices，以後再補基礎設施」 | 沒有基礎設施的微服務是災難。先建立 CI/CD、observability、container orchestration |

## Token Impact

避免在專案早期投入大量營運成本在微服務基礎設施上。一個不必要的微服務拆分可能花費 2-4 週建立基礎設施，並持續增加每次開發的 overhead。

---

← [回到 engineering/architecture/](README.md)
