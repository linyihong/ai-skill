# Requirement Hallucination

**Status**: `candidate-intelligence`

## 反模式

Agent 把使用者未確認的推論寫成需求、feature 或 acceptance criteria。

## 訊號

- 新增行為沒有 upstream requirement。
- 使用「順便」「通常會」「應該也要」擴張 scope。
- 實作先於 behavior contract。
- 測試證明的是 agent 自己想像的行為。

## 修正

將推論標記為 assumption 或 open question，回到 behavior contract / acceptance boundary。
