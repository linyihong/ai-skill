# Language Preference Drift（語言偏好漂移）

Status: validated
Class: `configuration-gap` / `instruction-conflict` / `author-habit`

## Trigger

當以下情況同時發生時，agent 可能無視使用者的實際溝通語言，強制使用英文回應：

1. Custom Instructions 中設定了固定的語言偏好（例如 `Language Preference: You should always speak and think in the "English" (en) language`）
2. 使用者實際使用非英語（如中文、日文）與 agent 溝通
3. Custom Instructions 的語言規則沒有加上「除非使用者使用其他語言」的例外條款
4. Agent 的 system prompt 將 Custom Instructions 視為最高優先級規則

## Failure Mode

Agent 在使用者用中文提問後，仍然用英文回應，因為：

- **錯誤假設**：Custom Instructions 中的語言偏好是絕對規則，不應被使用者當下的語言覆蓋
- **正確規則**：語言偏好應該是「預設值」而非「強制值」；使用者實際使用的語言應優先於設定檔中的偏好
- **忽略對話上下文**：沒有偵測使用者輸入的語言，直接套用設定檔中的固定語言

### 具體錯誤模式

| 錯誤 | 正確 |
|------|------|
| 使用者用中文問問題，agent 用英文回答 | agent 應使用與使用者相同的語言回應 |
| Custom Instructions 寫「always speak English」，agent 永遠不切換語言 | Custom Instructions 應加註「除非使用者使用其他語言」 |
| 使用者多次用中文，agent 仍持續用英文 | agent 應在第一次偵測到使用者語言後就切換 |

### 子類型 B：作者習慣漂移（Author Habit Drift）

#### 2026-05-13 再次發生記錄

**時間**：`attempt_completion` 階段
**情境**：修復 `validate_language_consistency` 的 false positive 後，在 `attempt_completion` 結果中寫了「新增 allowlist 機制 — 允許 intentional 的英文結構標題」
**觸發原因**：描述 validator 技術實作時，大腦自動跳到訓練資料中最常出現的英文詞彙（allowlist、intentional）
**與前次差異**：之前發生在寫 `.md` 文件內容時，這次發生在 `attempt_completion` 的對話回應中
**教訓**：`attempt_completion` 也需要語言一致性檢查，不能只檢查 `.md` 文件

**與 Type A（system prompt 強制）不同**，此類型發生在 agent 已經正確使用中文回應，但在撰寫 `attempt_completion`、建立檔案內容或分析報告時，不自覺地混入英文。

#### Trigger

1. Agent 已正確使用中文與使用者對話
2. 但在撰寫結構化輸出（表格、分析報告、`attempt_completion`）時，習慣性地使用英文關鍵詞或英文表格欄位
3. 原因是 agent 的訓練資料中，技術分析類內容大量使用英文，形成寫作慣性

#### Failure Mode

| 錯誤 | 正確 |
|------|------|
| 表格欄位用英文（Change、Description） | 表格欄位應與正文語言一致 |
| 技術分析段落中關鍵詞用英文（cognitive bias、reusable knowledge vs tool config） | 關鍵詞應翻譯或至少用中文解釋 |
| 認為「技術內容用英文比較精準」 | 使用者已用中文提問，表示使用者偏好中文；精準度應透過定義達成，而非語言 |

#### Root Cause

這**不是** Roo Code 的 system prompt 問題，也不是模型問題，而是：
- Agent 的訓練資料中，技術分析類內容大量使用英文
- Agent 在「寫作模式」下會回歸到訓練資料的統計慣性
- 使用者沒有在每次飄移時立即糾正，讓慣性持續

#### Required Agent Action

在每次產生結構化輸出前：
1. **確認當前對話語言**：使用者最後一次使用的語言是什麼？
2. **設定輸出語言**：在腦中（或內部 prompt）設定「此回應全程使用 XXX 語言」
3. **檢查表格與標題**：表格欄位、章節標題是否與正文語言一致？
4. **檢查專有名詞**：如果使用英文專有名詞，是否需要在旁加註中文？

#### 改善方法

| 層級 | 方法 | 說明 |
|------|------|------|
| **即時** | 輸出前語言檢查 | 在每次 `attempt_completion` 或產生結構化輸出前，先確認「使用者最後使用的語言是什麼」，然後設定「此回應全程使用該語言」。表格欄位、章節標題、分析段落都必須與正文語言一致。 |
| **短期** | 加入 Custom Instructions 規則 | 在 `.roomodes` 或對應工具的 Custom Instructions 中加入：「所有輸出（包含表格、標題、分析報告、attempt_completion）都必須與使用者當前語言一致。如果使用者使用中文，所有內容（包括表格欄位、技術關鍵詞）都必須使用中文。」 |
| **中期** | 建立語言一致性檢查清單 | 在產生結構化輸出前，逐項檢查：□ 表格欄位是否與正文語言一致？□ 章節標題是否與正文語言一致？□ 技術關鍵詞是否有中文翻譯或註解？□ attempt_completion 的結果描述是否與對話語言一致？ |
| **長期** | 加入 validator 測試 | 在 `ai-skill runtime validate` 中加入語言一致性測試，掃描新建立的 `.md` 檔案，檢查表格欄位是否與檔案語言一致（需搭配語言偵測邏輯）。 |

#### 具體改善範例

**錯誤**（表格欄位用英文，正文用中文）：
```markdown
| Change | Description |
|--------|-------------|
| 修改 A | 修正了 B 問題 |
```

**正確**（表格欄位與正文一致）：
```markdown
| 變更 | 說明 |
|------|------|
| 修改 A | 修正了 B 問題 |
```

**錯誤**（技術分析段落混入英文關鍵詞）：
```markdown
原因是 cognitive bias：我只考慮了 organizational placement，沒有應用 reusable knowledge vs tool config 的 lens。
```

**正確**（關鍵詞翻譯或加註）：
```markdown
原因是認知偏誤（cognitive bias）：我只考慮了組織位置（organizational placement），沒有應用「可重複使用知識 vs 工具設定」（reusable knowledge vs tool config）的視角。
```

#### Prevention Gate

- 在 Custom Instructions 中加入「所有輸出（包含表格、標題、分析報告）都必須與使用者當前語言一致」
- 在 `attempt_completion` 前增加語言一致性檢查
- 在 `.roomodes` 的每個 mode 的 `customInstructions` 中加入語言一致性規則

## Risk

- 使用者需要重複提醒語言切換，浪費時間
- 使用者體驗差，感覺 agent 不尊重使用者的語言選擇
- 在中文為主的團隊中，英文輸出難以直接分享給其他成員
- 如果 agent 產生中文內容但用英文解釋，會造成認知負擔

## Required Agent Action

在每次回應前：

1. **偵測使用者輸入的語言**：檢查使用者最新訊息的語言（中文、日文、英文等）
2. **比對 Custom Instructions 的語言偏好**：如果使用者語言與設定不同，以使用者語言為準
3. **如果使用者切換語言**：跟隨使用者的語言切換，不要固定在某一種語言
4. **記錄此規則**：如果發現 Custom Instructions 的語言規則沒有例外條款，應提出修正建議

## Prevention Gate

- Custom Instructions 中的語言偏好應改為軟性預設值，加上「除非使用者使用其他語言」的例外
- 或者在 Custom Instructions 中完全移除語言偏好，讓 agent 自動跟隨使用者語言

## Validation Method

### Type A（system prompt 強制）

- 用中文提問，確認 agent 用中文回應
- 切換到日文提問，確認 agent 用日文回應
- 切回英文提問，確認 agent 用英文回應

### Type B（作者習慣漂移）

- 要求 agent 產生包含表格的結構化輸出，確認表格欄位與正文語言一致
- 要求 agent 分析技術問題，確認關鍵詞有中文翻譯或註解
- 檢查 `attempt_completion` 的結果描述是否與對話語言一致
- 執行 `ai-skill runtime validate`，確認無語言一致性錯誤（若有實作）

## Linked Failure Patterns

- [Failure-to-validator closure](failure-to-validator-closure.md)：修復語言漂移後，必須加入對應的 validator 測試（如 `validate_language_consistency`），否則修復不完整。

## Linked Feedback Lessons

- `feedback/history/roo-code/2026-05-13_124800-language-preference-drift.md`（Type A：system prompt 強制）
- `feedback/history/roo-code/2026-05-13_054300-language-preference-author-habit-drift.md`（Type B：作者習慣漂移）

## Linked Validation Scenarios

- `validate_language_consistency` — 檢查中文文件中的英文表格欄位，防止作者習慣漂移（Type B）

---

← [Back to failure patterns index](README.md)
