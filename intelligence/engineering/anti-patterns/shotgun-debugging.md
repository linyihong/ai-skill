# Shotgun Debugging（散彈式除錯反模式）

Status: candidate-intelligence
Layer: `intelligence/engineering/anti-patterns/`

## 反模式定義

當症狀來源未被定位時，同時改多個變數試圖「修好它」，沒有 baseline、沒有單變數變更、沒有 after 對照。在 AI 輔助開發下，這個反模式被放大十倍——agent 可以一次改 N 個檔案，但無法回答「N 個改動中哪個生效」。

## 觸發信號

- 單一 commit / PR 同時動多個 unrelated config / 設定 / 程式碼，宣稱「修一個 symptom」
- 描述「應該修好了」「試試看」「reset / restart 再看」，但沒有 before / after 量化證據
- Agent 對「為什麼這樣改有用」無法明確回答，只能說「組合起來應該會解決」
- 修復成功後仍無法回答「下次再發生時要從哪個變數找起」

## 為什麼會發生

1. **快速嘗試的成本被誤估**：改 5 個變數的「嘗試成本」看起來低（自動化、批次），但「事後歸因成本」高到無法支付
2. **觀察—介入—驗證的順序倒置**：先介入再觀察，把 noise 帶進 baseline，後續無法判斷是 fix 還是 noise
3. **Agent 沒有強制 baseline 的機制**：人類 troubleshooter 會在白板寫 baseline 數字；agent 在 prompt 沒明確要求時容易跳過
4. **「修好就好」的短期 reward 太強**：production 火災時尤其；但這把學習機會（哪個變數是真因）丟掉了

## 與其他反模式的差別

- vs `architecture-absolutism`：那是設計階段反模式；shotgun debugging 是 troubleshooting / fix 階段反模式
- vs `migration-feature-bundling`：那是 migration 時把搬遷+新功能綁一起；shotgun debugging 是 fix 時把多個變更綁一起
- vs `framework-duplication-without-interrogation`：那是修改 framework 前沒拷問需求；shotgun debugging 是修改任何東西前沒拷問「我在量什麼」

## 反例與正例

### 反例（shotgun）

```
症狀：API latency p95 突然飆高
agent 動作：同時改 connection pool 大小 + DB index + 升級 ORM + 加 Redis cache + 調 Nginx timeout
結果：latency 回到正常
事後問題：到底是哪個改動修好的？下次再發生要從哪裡開始？
```

### 正例（measure → process → verify）

```
症狀：API latency p95 突然飆高
Step 0 baseline：記錄當前 p95 數字 / DB query count / connection pool 使用率 / cache hit rate
Step 1 觀察：APM trace 顯示 80% 時間在某個 SQL query
Step 2 假設 + 單變數變更：對該 query 加 covering index
Step 3 verify：再跑 baseline 測，p95 是否回落
Step 4 結論：找到根因（missing index）+ 留下未來 detection rule
```

## Required Agent Action

當 agent 被指派「修一個 symptom」類任務時：

1. **強制 baseline 階段**：在做任何改動前，回答「我現在量到的數字是什麼」「修好的判斷標準是什麼」
2. **強制單變數變更**：每次 commit / iteration 只改一個假設來源的變數；多個假設要排隊
3. **強制 verify 階段**：每次改動後用同樣的方法重測，比對 before / after
4. **拒絕「組合修法」**：若無法說明每個改動的獨立貢獻，就是違反此原則

## Prevention Gate

- **PR / commit 層**：troubleshooting commit body 應有 `Before:` / `After:` 區段並含可量化數字
- **Workflow 層**：troubleshooting / optimization 類任務套用 measure → process → verify 模板（見 [`workflow/software-delivery/perf-risk-gate.md`](../../../workflow/software-delivery/perf-risk-gate.md)）
- **Reviewer 層**：reviewer 看到 commit 動多個 unrelated 設定且只解一個 symptom，要求拆 commit 並補 before/after 數據
- **Runtime 層（已部分覆蓋）**：本 repo cognitive contract v2 的 `validation_mode != NONE` 是這個原則的 runtime projection

## 與本 repo 既有設計的關聯

本 repo 的 cognitive contract v2（`runtime/cognitive-modes-*.yaml`）已強制 `validation_mode` enum，並在 commit-msg hook 驗證 cognitive cost / activation signal。這是 shotgun debugging 的部分 runtime 防護。本 atom 把該防護背後的工程原理命名，便於跨工具引用（其他工具 / 其他專案沒有同等 runtime 時，可手動套用本原則）。

## Related

- [`wish-to-task-list-translation.md`](../agent-architecture/wish-to-task-list-translation.md) — Agent 層的姊妹 atom：user wish 進來時應先翻成 task list，避免 shotgun
- [`workflow/software-delivery/perf-risk-gate.md`](../../../workflow/software-delivery/perf-risk-gate.md) — measure → process → verify 的具體執行模板
- [`enforcement/failure-patterns/correction-loop-bypass.md`](../../../enforcement/failure-patterns/correction-loop-bypass.md) — 鄰近：使用者指出修正不完整時，agent 只改當下文字（也是 fix 階段反模式）
- [`migration-feature-bundling.md`](migration-feature-bundling.md) — Migration 階段的同類綁定反模式

## Source

- 2026-05-27 session：使用者提供外部 infographic「想讓網路變快，不要靠感覺」，列出 measure → process → verify 三段；本 atom 將其反面（shotgun debugging）抽象為跨工具反模式並命名。Status `candidate-intelligence` 至 repo 內首次觀察到具體 shotgun-debugging commit 後 promote 為 `validated`。
