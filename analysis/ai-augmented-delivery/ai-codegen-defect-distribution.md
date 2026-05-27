# AI Codegen Defect Distribution（AI 生成程式碼的缺陷分布）

## 觀察的核心問題

當 AI codegen 工具（Copilot / Claude / Cursor / 內部 agent）大幅加速程式碼產出時，**缺陷在哪個階段被抓到的分布會改變**。傳統 pipeline 假設大多數缺陷在 unit / integration test 階段就被攔截；AI 輔助開發下，缺陷會更多流向 staging 與 production。

## 量化參考點（2025 業界觀察）

> **資料來源說明**：以下量化來自 2025 公開發表的業界研究／infographic 摘要（外部素材），非 Ai-skill repo 內 first-party 觀察。引用時請保留來源並標註觀察日期。

| 指標 | 數值 | 含義 |
|---|---|---|
| 程式碼產出倍率 | **3.2x**（4,450 → 14,148 行 / 開發者 / 月） | AI 工具讓 raw 產出量大幅上升 |
| Production 仍需手動 debug 的 AI 生成程式碼 | **43%** | 通過 QA / staging 不代表 production 安全 |
| 確認 AI 修復有效所需的 redeploy 次數 | **2–3 次**（88% 公司） | 「修好了」的訊號比以前更不可信，需要多次驗證 |
| 開發者每週花在 debug 與驗證的時間 | **約 2 天 / 5 天**（38% 開發者） | 驗證階段已是主要時間 sink |

## 拆解這四個數字

### 1. 3.2x 產出 ≠ 3.2x 交付

產出（lines written）是上游指標；交付（lines safely in production）是下游指標。若驗證能力不變，3.2x 產出只是把更多東西塞進更窄的下游瓶頸。

### 2. 43% production debug 揭露什麼

通過 QA / staging 的 AI 生成程式碼仍有 43% 要在 production 手動 debug，意味：

- **CI 涵蓋度不足**：unit / integration test 沒覆蓋到 AI 容易犯的缺陷類型
- **AI 的盲點與人類不同**：AI 寫的 code 結構上看起來合理，但對 load、邊界 case、長尾資料分布缺乏直覺
- **「看起來像對的」程式碼比例上升**：reviewer 與測試的判別成本上升

### 3. 88% 公司需 2–3 redeploy 揭露什麼

修復一次性成功率下降，意味：

- **AI 修 bug 容易引入新 bug**：修 A 影響 B 的隱性副作用變多
- **回歸測試訊號不夠**：CI 通過不代表 fix 真的有效
- **信任成本上升**：每次「應該修好了」都要花更多次部署來確認

### 4. 38% 開發者 2 天 / 週在 debug

開發者實際時間分布從「寫程式」往「驗證程式」傾斜，意味：

- **產出加速但人不能放鬆**：表面看時間省了，實際被驗證吃掉
- **debug 是高認知成本工作**：比 write code 更耗心智頻寬
- **Net 生產力不一定上升**：要看「寫 3.2x 但花 1.6x 時間 debug」的淨值

## 對應的工程行動

| 觀察 | 工程行動 |
|---|---|
| 43% production debug | 投資 APM / distributed tracing / SLO alerting，把缺陷在 production 階段也能被快速定位 |
| 2–3 次 redeploy | Canary / progressive rollout / feature flag，讓 fix 上線不需要 full deploy |
| 38% debug 時間 | PR 階段 reviewer checklist + hot-path micro-benchmark，把缺陷往上游推 |
| 3.2x 產出 | Reviewer 心智頻寬無法 3.2x，需要 detection 自動化 + 風險分級 |

## 引用方式

引用本資料時，請保留：
- 觀察年份（2025）
- 來源類型（業界 infographic / 公開研究，非 first-party）
- 數值區間（單一數字只是中位估計，實際差異會隨團隊、語言、工具差很多）

## Related

- [`generation-validation-rate-parity.md`](../../intelligence/engineering/ai-augmented-delivery/generation-validation-rate-parity.md) — 抽象原則
- [`perf-test-bottleneck-anatomy.md`](perf-test-bottleneck-anatomy.md) — 為什麼 unit test 抓不到 perf 缺陷
- [`analysis/production/`](../production/README.md) — Production 問題分析的通用方法

## Source

- 2026-05-27 session：使用者提供外部 infographic 摘要四項量化資料；本檔將其轉為工具中立觀察並標註 candidate 狀態。
