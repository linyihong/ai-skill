# AI Codegen Validation Rate Parity（外部素材沉澱：AI 寫程式爆量但 Perf Test 跟不上）

Lesson date: 2026-05-27
Trigger: 使用者提供外部 2025 業界 infographic「AI 寫程式爆量的時代，Performance Test 已經跟不上了」
Status: lesson-captured

## 觀察到的訊號

外部 infographic 摘要：
- AI codegen 把每位開發者每月程式碼產出推到 ~3.2x（4.5k → 14k 行）
- 43% AI 生成程式碼通過 QA/staging 但 production 仍需手動 debug
- 88% 公司需要 2–3 次 redeploy 才能確認 AI 修復有效
- 38% 開發者每週約 2 天在 debug 與驗證

四個 perf anti-pattern：迴圈藏 DB query、collection 無界、外部呼叫無 timeout、SQL 字串拼接。

## Lesson

當外部素材描述「某類加速工具導致驗證跟不上」時，應該識別其為**meta-tool-risk** 類型的觀察，跨工具可重用。處理路徑不是單一層級的更新，而是 6 層同步：

| 層 | 產物 |
|---|---|
| `intelligence/engineering/<domain>/` | 抽象原則 atom（為什麼這樣判斷） |
| `analysis/<domain>/` | 量化資料 + 解剖（如何觀察） |
| `enforcement/failure-patterns/` | Trigger + detection rule（如何防止） |
| `validation/scenarios/<workflow>/` | 機械化 scenario（如何驗證） |
| `workflow/<workflow>/` | 執行流程（何時執行 + 步驟） |
| `feedback/history/<domain>/` | 本 lesson（觀察記錄） |

`governance/` 為可選第 7 層，當 production-gate 需要 cross-tool 規範時加。

## 容易踩到的反模式（meta）

1. **只放 intelligence**：抽象原則沒對應觀察方法 → 沒人找得到、用不出來
2. **只放 enforcement**：detection rule 沒對應原理 → 規則更新時失去思考脈絡
3. **沒有 analysis layer**：量化資料散在 intelligence 或 enforcement 文件中 → 引用時混淆「是觀察還是判斷」
4. **沒有 workflow integration**：知識存在但 agent / reviewer 在實際流程中不知道何時觸發
5. **沒有 validation scenario**：流程只能靠人腦執行，無機械化路徑

## 對應到 Ai-skill 自身的關聯

本 repo 的 cognitive contract / hooks / runtime validation stack 就是「同步加速驗證」的具體實作。每次新加 codegen 自動化能力（例如本 session 的 Go-native hooks 取代 .sh），都應檢查 validators 是否需要同步擴張。

## 後續行動

- 候選狀態 `candidate-intelligence`：待 repo 內 first-party 觀察出現（例如：本 repo 自身或夥伴專案中遇到 AI 生成的程式碼通過 CI 但 production 出 perf bug），promote 為 `validated`
- 若高頻使用 perf-risk-gate workflow，後續加 `validatePerfRisks` pre-commit validator（Phase 7 候選工作）

## Related

- [`intelligence/engineering/ai-augmented-delivery/generation-validation-rate-parity.md`](../../../../intelligence/engineering/ai-augmented-delivery/generation-validation-rate-parity.md)
- [`analysis/ai-augmented-delivery/`](../../../../analysis/ai-augmented-delivery/README.md)
- [`enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md`](../../../../enforcement/failure-patterns/ai-codegen-passes-ci-fails-production.md)
- [`workflow/software-delivery/perf-risk-gate.md`](../../../../workflow/software-delivery/perf-risk-gate.md)
- [`validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml`](../../../../validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml)
