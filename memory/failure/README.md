# Failure Memory

`memory/failure/` 保存**抽象化的失效模式記憶**。不同於 `enforcement/failure-patterns/`（可執行的 failure detection 規則），failure memory 記錄的是「發生過什麼失效、為什麼發生、學到了什麼」，作為 failure intelligence 的原始素材與回顧依據。

Failure memory 可提示風險與 prevention direction，但不能取代 enforcement rule、current source、validator 或 user-facing completion evidence。

## 用途

- 記錄跨 session 的失效事件與其根本原因
- 保存失效的 detection signals 與 context
- 追蹤失效的緩解措施與 prevention 策略
- 提供 `intelligence/engineering/failure/` 的抽象化素材
- 支援 failure pattern 的長期演化追蹤

## 不放什麼

- 可執行的 failure detection 規則 → `enforcement/failure-patterns/`
- 抽象化的 failure intelligence atom → `intelligence/engineering/failure/`
- 專案私有的 incident raw logs → 留在業務專案
- Session-local 的錯誤記錄 → `memory/working/`
- Feedback lesson 的 promotion workflow → `feedback/`

## 格式

```markdown
# Failure: {失效模式名稱}

## Status
{active | monitored | resolved | archived}

## Trigger Context
{什麼情境下發生了這個失效}

## Symptoms
- {symptom 1}
- {symptom 2}

## Root Cause
{根本原因是什麼}

## Impact
{影響範圍與嚴重程度}

## Detection Signals
- {signal 1}：{如何觀察到}
- {signal 2}：{如何觀察到}

## Mitigation
{採取了什麼緩解措施}

## Prevention
{如何防止再次發生}

## Generalized Lesson
{可抽象化的教訓（→ intelligence/engineering/failure/ 的候選）}

## Occurrences
- {date}：{brief context}（→ memory/summary/{session-file}）
- {date}：{brief context}

## Linked Patterns
- {related failure pattern 1}（→ enforcement/failure-patterns/）
- {related failure pattern 2}
```

## 規則

1. **抽象化優先**：記錄失效時應著重「可泛化的教訓」，而非專案特定的細節。
2. **去敏**：Failure memory 不保存 token、host、private key 等敏感資訊。
3. **連結到 source**：每個 failure record 應連結到發生時的 session summary，以便追溯。
4. **Promotion 路徑**：成熟的 failure lesson 應考慮 promotion 到：
   - `enforcement/failure-patterns/`（可執行規則）
   - `intelligence/engineering/failure/`（抽象化 intelligence）
5. **演化追蹤**：同一 failure 模式多次發生時，應更新 occurrences 列表並評估是否需要升級 prevention 策略。
6. **Token-aware**：每個 failure record 不超過 400 tokens。
7. **Risk hint only**：Replay failure memory 時，預設只作 weak / scoped risk hint。
8. **Current validation required**：若 failure memory 影響本次 patch、commit、runtime 或規則更新，必須重新驗證 current source。
9. **No workaround promotion**：不得把舊 workaround 直接 promotion 成 permanent policy；需先抽象化並驗證 recurring failure。

## 與既有層的關係

- `enforcement/failure-learning-system.md`：定義 failure learning 的可執行流程
- `enforcement/failure-patterns/`：可執行的 failure detection 規則（promotion 目標）
- `intelligence/engineering/failure/`：抽象化的 failure intelligence（promotion 目標）
- `memory/episodic/`：failure 發生的情境細節
- `memory/summary/`：failure 發生時的 session 摘要
- `feedback/replay/`：failure 可觸發 replay 流程
- `governance/lifecycle/`：failure memory 的 lifecycle 管理
- `memory/retrieval-governance/`：定義 repeated failure trigger、replay budget 與 contamination response
