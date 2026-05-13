# Language Preference Drift（語言偏好漂移）

Status: validated
Class: `configuration-gap` / `instruction-conflict`

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

- 用中文提問，確認 agent 用中文回應
- 切換到日文提問，確認 agent 用日文回應
- 切回英文提問，確認 agent 用英文回應

## Linked Failure Patterns

- (none)

## Linked Feedback Lessons

- `feedback/history/roo-code/2026-05-13_124800-language-preference-drift.md`

## Linked Validation Scenarios

- (none yet — 需要建立 cross-domain scenario 測試語言跟隨行為)

---

← [Back to failure patterns index](README.md)
