# Acceptance Boundaries

**Status**: `candidate-intelligence`

## 判斷原則

Acceptance boundary 定義「這個需求何時算完成」與「哪些行為不在本次範圍」。它防止 acceptance criteria 漂移成無限 feature list。

## 必填內容

- In-scope behavior。
- Out-of-scope behavior。
- Validation target。
- Regression scope。
- Manual-only 或 not-automatable 的理由。

## 風險訊號

- Acceptance criteria 只寫「正常運作」。
- 沒有 negative path。
- 沒有 validation target。
- 沒有說明哪些相鄰行為不處理。
