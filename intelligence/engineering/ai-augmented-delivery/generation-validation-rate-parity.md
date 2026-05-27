# Generation–Validation Rate Parity（產出—驗證速度對稱原則）

Status: candidate-intelligence
Layer: `intelligence/engineering/ai-augmented-delivery/`

## 主張（Claim）

> 任何加速「產出」的工具或流程，必須同步加速「驗證」；否則淨效益會被驗證瓶頸吃掉，且風險會集中流向 production。

這是個跨工具、跨領域的工程原則。AI codegen 只是最近的具體案例，同樣的原則適用於低程式碼平台、自動化 ETL、scaffolding generator、CI 加速、PR 自動分派、agent 框架等任何「產出加速」場景。

## 為什麼成立

開發流程是一個 throughput pipeline：

```
產出（write）→ 驗證（test / review / observe）→ 部署 → production
```

若只加速第一段，下游每一段都會塞車：

1. **驗證階段塞車**：reviewer 看不完、CI 等不完、QA 抓不完
2. **資訊延遲**：問題被往下游推，越下游成本越高（修一個 production bug 比修一個 PR bug 貴 10–100 倍）
3. **風險集中**：原本散佈在 PR / staging / canary 的問題集中爆發在 production
4. **回饋循環變慢**：開發者收到「這段程式碼錯了」的訊號時間從分鐘 → 小時 → 天

外部觀察支持此原則：AI codegen 把每位開發者每月程式碼產出從 ~4.5k 行推到 ~14k 行（約 3.2x），但 43% 的 AI 生成程式碼即使通過 QA/staging 仍要在 production 手動 debug，88% 公司需要 2–3 次重新部署才能確認修復有效。詳細量化見 [`analysis/ai-augmented-delivery/ai-codegen-defect-distribution.md`](../../../analysis/ai-augmented-delivery/ai-codegen-defect-distribution.md)。

## 推論（Implications）

### 1. 加速產出 ≠ 提升 throughput

如果驗證沒同步加速，pipeline 整體 throughput 由最慢的環節決定（Little's Law / TOC）。淨產出可能比 baseline 更差，因為：
- production 事故成本爆掉
- 開發者更多時間在 debug 而非寫程式（外部觀察：38% 開發者每週花約 2 天在 debug 與驗證）
- 修復信任成本上升（2–3 次重新部署才確認）

### 2. 驗證投資要先於產出加速，至少要同步

若要採用會加速產出的工具，應先盤點驗證能力是否能跟上：
- Unit / integration test 是否能抓到該工具會犯的錯？
- Reviewer 心智頻寬是否能消化新流量？
- Observability（metrics / tracing / SLO）是否能在 production 抓到漏網問題？

### 3. 驗證投資的優先序

從事故發生地點往回推，越靠近 production 的環節投資 ROI 越高（事故已過濾過上游）：

```
production observability (APM / tracing / SLO alerting) ← 最高 ROI
        ↑
canary / progressive rollout
        ↑
integration / load / perf test
        ↑
PR review checklist + hot-path micro-benchmark
        ↑
unit test ← 最容易做但漏網率最高
```

### 4. 不同類型缺陷的捕捉位置

- **功能正確性**：unit / integration test 通常能抓
- **效能特性**（迴圈藏 DB query、collection 無界、外部呼叫無 timeout、SQL 字串拼接）：unit test 抓不到，integration test 通常也抓不到，要等 load test 或 production
- **資料一致性與 race**：要 chaos / 真實流量才浮現

AI codegen 的失敗分布偏向「效能特性」與「邊界 case」，所以 unit test 高通過率不代表生產安全。

## 不適用情境

- **驗證已過剩**：若 pipeline 瓶頸不在驗證階段（例如 deploy / release coordination），加速產出可能反而是淨增益
- **產出量極低**：每月只有少量 PR 的專案，driver 是別處
- **拋棄式 prototype / spike**：驗證需求被刻意降級，原則不適用（但需明確標註）

## 對應到 Ai-skill 自身

本 repo 的 cognitive contract / hooks / runtime 驗證 stack（pre-commit / commit-msg / pre-push validators + go test）就是「同步加速驗證」的具體實作。每次新加 codegen 自動化能力，都應檢查 validators 是否需要同步擴張。本 session 的 `validateNoNewShellScripts` 與 `validateNoNewShellScripts` → Go hooks migration 是典型例子：產出加速（script-as-Go 跨平台）必須伴隨對應 governance gate，否則治理會被新流量繞過。

## Related

- [`analysis/ai-augmented-delivery/`](../../../analysis/ai-augmented-delivery/README.md) — 量化觀察與解剖
- [`enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md`](../../../enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md) — 具體 trigger 與 detection
- [`workflow/software-delivery/perf-risk-gate.md`](../../../workflow/software-delivery/perf-risk-gate.md) — 執行流程整合
- [`validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml`](../../../validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml) — Reviewer checklist 機械化

## Source

- 2026-05-27 session：使用者提供外部 infographic「AI 寫程式爆量的時代，Performance Test 已經跟不上了」，含三組量化資料（43% / 88% / 38%）與四個 perf anti-pattern。原始素材為外部研究，本 atom 將其抽象為工具中立原則並標註 candidate-intelligence，待 repo 內 first-party 證據出現後再 promote 為 validated。
