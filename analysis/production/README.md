# Production Issue Analysis

`analysis/production/` 負責「Production 問題分析與根因追蹤方法」。本目錄保存如何觀察、拆解與診斷 production 問題的分析框架，讓 agent 能系統化地從症狀追蹤到根因，而不只是隨機嘗試修復。

## 核心責任

- Production incident 的結構化分析流程。
- 根因追蹤方法（RCA framework）。
- 觀測性資料判讀（metrics、logs、traces）。
- 效能問題診斷（CPU、memory、IO、network）。
- 資料不一致與錯誤率異常分析。

## 與其他層的關係

- `workflow/` 可引用本層的分析步驟，但不複製分析方法細節。
- `intelligence/engineering/failure/` 承接從 production 問題中萃取的抽象化失敗模式。
- `intelligence/engineering/tradeoffs/` 承接從 production 問題中萃取的技術取捨教訓。
- `skills/` 目前仍是相容入口；本層只承接逐步抽出的分析方法。

## 第一批候選遷移來源

- `shared-rules/failure-learning-system.md` 中偏 production 分析的方法。
- `plans/active/next-stage-upgrade-plan.md` 中 `analysis/` 的分層說明。

## 建議分析方法

### 1. Incident 分類與優先級

```
1. 確認影響範圍（單一用戶 / 部分功能 / 全部服務 / 資料遺失）。
2. 確認錯誤類型（錯誤率上升 / 延遲增加 / 功能不可用 / 資料不一致）。
3. 確認時間範圍（何時開始、持續多久、是否間歇性）。
4. 確認是否有近期變更（deploy、config change、dependency update）。
5. 根據 severity 決定分析深度（hotfix 優先 vs 完整 RCA）。
```

### 2. 觀測性資料收集

```
1. Metrics：CPU / memory / disk / network / request rate / error rate / latency p50/p95/p99。
2. Logs：error log、warning log、slow query log、access log。
3. Traces：request trace、dependency call tree、database query trace。
4. 基礎設施狀態：pod / container / VM 健康狀態、資源使用趨勢。
5. 外部依賴狀態：第三方 API、CDN、DNS、資料庫複製狀態。
```

### 3. 根因追蹤流程

```
1. 建立時間軸（incident 開始 → 發現 → 緩解 → 解決）。
2. 列出所有可能的根因假設。
3. 對每個假設收集支持或反對的證據。
4. 排除無法解釋所有症狀的假設。
5. 對剩餘假設進行最小干預驗證（canary、rollback、feature flag toggle）。
6. 確認根因後，記錄完整的「症狀 → 假設 → 證據 → 結論」鏈。
```

### 4. 效能問題診斷

```
1. CPU：identify 高 CPU 進程 / thread、檢查是否有 infinite loop、GC 頻率。
2. Memory：heap / off-heap 使用趨勢、GC 暫停時間、memory leak 模式。
3. IO：disk read/write 延遲、connection pool 耗盡、file descriptor leak。
4. Network：latency 分佈、packet loss、connection timeout、TLS handshake 失敗。
5. Database：slow query、lock contention、connection pool 耗盡、replication lag。
```

### 5. 資料不一致分析

```
1. 確認不一致的範圍（單一 record / 部分資料 / 全部資料）。
2. 追蹤資料寫入路徑（哪個 service、哪個 transaction、哪個 API）。
3. 檢查是否有 concurrent write 或 partial write。
4. 檢查是否有 migration / backfill 錯誤。
5. 檢查是否有資料格式變更但未處理舊資料。
```

## 產出格式

每次 production 分析應產出：

- **Incident 摘要**（≤200 tokens）：影響範圍、時間範圍、severity。
- **時間軸**（≤300 tokens）：關鍵事件與時間點。
- **根因分析**（≤300 tokens）：根因、證據鏈、排除的假設。
- **緩解與修復**（≤200 tokens）：hotfix、rollback、長期修復計畫。
- **預防措施**（≤200 tokens）：monitoring 改進、alert 調整、code review 重點。
