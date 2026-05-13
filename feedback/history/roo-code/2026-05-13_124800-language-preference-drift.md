# Language Preference Drift（語言偏好漂移）

## 觀察

在使用 Roo Code 的過程中，發現即使使用者用中文提問，agent 仍然會用英文回應。這是因為 Roo Code 的 Custom Instructions 中設定了固定的語言偏好。

## 原因

Roo Code Extension 的 **Custom Instructions** 設定中包含以下規則：

```
Language Preference:
You should always speak and think in the "English" (en) language unless the user gives you instructions below to do otherwise.
```

這條規則的問題：
1. 它是**絕對規則**而非**軟性預設值**
2. 它沒有考慮使用者實際使用的語言
3. 它與「跟隨使用者語言」的直覺行為衝突

## 影響

- 使用者需要重複提醒「用中文回答」
- 中英文夾雜的對話造成認知負擔
- 團隊協作時，英文輸出難以直接分享

## 修正方式

將 Custom Instructions 中的語言偏好改為：

```
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.
```

## 預防

1. Custom Instructions 中的語言規則應加上「除非使用者使用其他語言」的例外
2. 或者在 Custom Instructions 中完全移除語言偏好，讓 agent 自動跟隨使用者語言
3. 已建立 failure pattern: `shared-rules/failure-patterns/language-preference-drift.md`

## 相關檔案

- Failure pattern: [`shared-rules/failure-patterns/language-preference-drift.md`](../../shared-rules/failure-patterns/language-preference-drift.md)
- Roo Code 設定文件: [`ai-tools/agent/roo.md`](../../ai-tools/agent/roo.md)
