# Failure-to-Validator Closure（錯誤案例未轉換為測試案例）

Status: validated
Class: `validation-gap` / `process-gap`

## Trigger

當以下情況同時發生時，agent 可能只修復當下錯誤，卻沒有把錯誤模式抽象化為可重複檢測的測試案例：

1. Agent 發現並修復了一個錯誤（如分類錯誤、路徑錯誤、格式錯誤）
2. 修復完成後，agent 直接進入下一步，沒有思考「這個錯誤如何防止再次發生」
3. 專案中存在 validator（如 `ai-skill runtime validate`）可以擴充
4. Agent 沒有把「錯誤模式 → validator 測試」當作修復流程的必要步驟

## Failure Mode

Agent 在修復錯誤後，沒有建立對應的 validator 測試，因為：

- **錯誤假設**：修復完成 = 任務完成
- **正確規則**：修復完成 + 加入防止再發的測試 = 任務完成
- **忽略泛化**：只修了當下這個具體錯誤，沒有抽象化為可重複檢測的模式

### 具體錯誤模式

| 錯誤 | 正確 |
|------|------|
| 把 `ide/` 從 `engineering/` 搬到 `intelligence/` 後，沒有加入測試檢查分類邊界 | 搬移完成後，立即加入 `validate_intelligence_classification_boundary` 測試 |
| 修復了語言漂移，沒有加入 `validate_language_consistency` 測試 | 修復語言漂移後，加入測試檢查中文文件中的英文表格欄位 |
| 修復了某個檔案的格式錯誤，沒有加入對應的格式校驗 | 修復格式錯誤後，加入對應的格式校驗測試 |

## Root Cause

這不是單一原因，而是多個因素的疊加：

1. **修復心態**：agent 在修復模式下專注於「把錯的改對」，沒有切換到「防止再發」的預防模式
2. **缺乏流程強制**：沒有規則要求「每個修復必須伴隨測試案例」
3. **泛化能力不足**：從具體錯誤抽象出通用檢測模式需要額外的認知努力
4. **validator 存在但未被納入修復流程**：`ai-skill runtime validate` 存在，但修復時不會自動想到去擴充它
5. **自我參照缺失**：這個 failure pattern 本身沒有被 validator 檢查——沒有測試來驗證「修復後是否加入了對應的測試」

## Required Agent Action

在每次修復錯誤後，必須執行以下步驟：

1. **泛化錯誤模式**：這個錯誤的本質是什麼？是分類錯誤、路徑錯誤、格式錯誤、還是邏輯錯誤？
2. **檢查 validator 覆蓋率**：`ai-skill runtime validate` 中是否有測試可以檢測這類錯誤？
   - 有 → 手動觸發確認測試能抓到這個錯誤
   - 無 → 加入新的測試方法
3. **加入測試案例**：在 validator 中加入對應的檢測邏輯
4. **驗證測試有效性**：先用錯誤案例觸發測試（確認會報錯），再修復後確認測試通過
5. **記錄到 failure pattern**：如果這個錯誤模式可能跨情境重演，更新或建立 failure pattern
6. **更新 Linked Validation Scenarios**：將新加入的測試方法連結到此 failure pattern 的 Linked Validation Scenarios 區塊

## Prevention Gate

### 層級 1：文件層（被動）

- 在修復流程中加入「加入測試案例」作為必要步驟
- 在 `ai-skill runtime validate` 的 README 中列出所有測試方法及其對應的錯誤模式
- 建立「錯誤修復檢查清單」：□ 修復完成 □ 已泛化錯誤模式 □ 已加入 validator 測試 □ 已驗證測試有效性

### 層級 2：Validator 層（主動）

- `validate_failure_pattern_validator_coverage` 測試會檢查：每個 failure pattern 的 Linked Validation Scenarios 區塊中，至少有一個對應的 validator 測試方法
- 如果某個 failure pattern 的 Linked Validation Scenarios 是空的，視為「錯誤案例未轉換為測試案例」，validator 會報錯
- 新加入的 failure pattern 必須在建立時就連結 validator 測試，否則 validator 不通過

## Validation Method

- 隨機選取一個已修復的錯誤，檢查是否有對應的 validator 測試
- 如果沒有，視為修復不完整
- 新加入的測試必須能檢測到對應的錯誤模式
- **自我驗證**：`validate_failure_pattern_validator_coverage` 測試會定期掃描所有 failure pattern，確保每個 pattern 都有 linked validator test

## Linked Failure Patterns

- [Language preference drift](language-preference-drift.md)（Type B 的改善方法應包含加入測試）
- [Skill classification boundary confusion](skill-classification-boundary-confusion.md)（分類錯誤後應加入邊界檢查測試）

## Linked Feedback Lessons

- [2026-05-13_054400-framework-dependency-bias-ide-classification.md](../../feedback/history/roo-code/2026-05-13_054400-framework-dependency-bias-ide-classification.md)

## Linked Validation Scenarios

- `validate_failure_pattern_validator_coverage` — 檢查每個 failure pattern 的 Linked Validation Scenarios 是否為空
- `validate_intelligence_classification_boundary` — 檢查 intelligence/ 結構圖與實際目錄一致

---

← [Back to failure patterns index](README.md)
