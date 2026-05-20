# Architecture Fit Signals for DDD

**Status**: `candidate-intelligence`

## 判斷原則

DDD adoption 必須由 architecture fit signal 驅動，而不是由 agent 偏好或框架慣性驅動。

## 強 adoption signal

- 高 domain complexity。
- 高 invariant density。
- 高 business language instability。
- 多 bounded context 或 subdomain。
- 高 integration pressure。
- 長 lifecycle、多人協作或 team boundary 明確。

## 弱 adoption signal

- 只是 entity 多。
- 只是程式碼資料夾大。
- 只是想測試 pattern。
- 只是使用某個 framework。

## 輸出要求

Architecture recommendation 必須列出：採用策略、拒絕的更複雜策略、支持證據、以及後續升級條件。
