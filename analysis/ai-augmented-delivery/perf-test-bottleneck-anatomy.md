# Performance Test Bottleneck Anatomy（為什麼 AI 生成程式碼的 perf 缺陷漏網）

## 核心觀察

AI codegen 的功能正確性通常會通過 unit / integration test，但**效能特性常常是錯的**。這不是 AI 特有，而是 perf 缺陷本身的偵測成本特性：許多 perf 缺陷在小資料量、單執行緒、無延遲注入的測試環境裡看不出來，要等真實 load 或 production 才浮現。

## 四個常見 anti-pattern 的解剖

以下四項是 AI 生成程式碼最常出現的 perf 隱患。共通特徵：**unit test 通過率不變，integration test 經常也過，要到 load test 或 production 才爆**。

### 1. 迴圈裡藏 DB query（N+1）

```python
# AI 寫的代碼看起來合理
for user in users:
    orders = db.query("SELECT * FROM orders WHERE user_id = ?", user.id)
    process(user, orders)
```

- **單元測試**：users 只有 3 筆，跑 3 次 query 也是 30ms，通過
- **整合測試**：跑 fixture 資料，量級類似，通過
- **Load test 才浮現**：production users 10 萬筆 → 10 萬次 DB round-trip → DB 飽和
- **Detection 成本**：靜態 AST scan 可抓（迴圈內含 ORM / db.Query 呼叫），cost 低

### 2. Collection 沒綁大小（unbounded growth）

```go
// AI 寫的快取/累積邏輯
var cache = make(map[string]Result)
func Process(key string) Result {
    if v, ok := cache[key]; ok { return v }
    r := compute(key)
    cache[key] = r // 永遠不刪
    return r
}
```

- **單元測試**：呼叫 10 次，map 大小 10，通過
- **整合測試**：fixture 不會撐大記憶體，通過
- **Production 才浮現**：跑 30 天後 OOM 或 GC 暴衝
- **Detection 成本**：靜態 scan 可抓「無 eviction 的 map / slice append in long-lived scope」，cost 中

### 3. 外部呼叫沒設 timeout

```javascript
// AI 寫的 HTTP 呼叫
const response = await fetch(url);  // 預設 no timeout
const data = await response.json();
```

- **單元測試**：mock 立即回應，通過
- **整合測試**：testing endpoint 通常快，通過
- **Production 才浮現**：下游慢回應 → connection pool 耗盡 → 整個服務 hang
- **Detection 成本**：grep / linter rule 可抓（http.Client 無 Timeout、fetch 無 AbortSignal），cost 低

### 4. SQL 用字串拼接

```python
# AI 寫的查詢
sql = "SELECT * FROM users WHERE name = '" + name + "'"
db.execute(sql)
```

- **單元測試**：mock data，通過
- **整合測試**：純功能驗證，通過
- **Production 才浮現**：(a) SQL injection 安全事故；(b) 每個 query 字串不同 → DB query plan cache miss → DB CPU 暴衝
- **Detection 成本**：grep / linter / SAST 可抓（字串拼接後接 SQL keyword），cost 低

## 為什麼 unit / integration test 抓不到

| 缺陷類型 | unit test 假設 | 為什麼漏網 |
|---|---|---|
| N+1 query | 資料量小 | 沒有真實量級的 fixture |
| Unbounded collection | 跑得快結束 | 沒有長時間運行 |
| No timeout | mock 立即回應 | 沒有模擬慢回應 / 失敗 |
| String SQL | 功能正確就好 | 沒測 query plan cache、沒做 injection check |

共通根因：**單元/整合測試的目標是「功能對不對」，不是「效能特性對不對」**。要抓 perf 缺陷，需要的是：

1. **PR 階段靜態 detection**（grep / AST scan / linter） — 對上述 4 類最有效
2. **熱路徑 micro-benchmark** — 對 N+1、algorithmic complexity 有效
3. **Load test / k6 / locust** — 對 unbounded、timeout cascade 有效
4. **Production observability**（APM / tracing） — 對混合或環境相依問題有效

## 對應的 detection 投資

| Anti-pattern | 最便宜的 detection | 升級路徑 |
|---|---|---|
| N+1 | Linter rule 或 AST scan 偵測「迴圈內 DB 呼叫」 | + ORM 內建 N+1 detector + tracing query count metric |
| Unbounded collection | Linter rule 偵測「long-lived map / slice 無 eviction」 | + JVM / runtime memory profiler in canary |
| No timeout | Grep / linter rule 偵測 HTTP client 無 timeout 設定 | + chaos injection 模擬慢回應 |
| String SQL | SAST 偵測字串拼接後接 SQL keyword | + ORM / query builder 強制 + DB query plan cache hit rate metric |

四項中三項可在 PR 階段靜態抓到，這是 ROI 最高的投資。對應的 reviewer checklist 機械化見 [`validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml`](../../validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml)。

## Related

- [`ai-codegen-defect-distribution.md`](ai-codegen-defect-distribution.md) — 量化資料
- [`generation-validation-rate-parity.md`](../../intelligence/engineering/ai-augmented-delivery/generation-validation-rate-parity.md) — 抽象原則
- [`enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md`](../../enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md) — Trigger 與 prevention gate
- [`analysis/production/`](../production/README.md) — Production 問題的通用分析方法

## Source

- 2026-05-27 session：使用者提供外部 infographic「為什麼 Performance Test 特別慘」段落，列出 4 個 anti-pattern；本檔將其轉為解剖 + detection cost 對照。
