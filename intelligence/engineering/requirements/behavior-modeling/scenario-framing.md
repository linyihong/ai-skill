# Scenario Framing

**Status**: `candidate-intelligence`

## 判斷原則

Scenario framing 將需求轉成可驗證的行為片段，而不是把使用者一句話擴張成完整 feature set。

## 健康 scenario

- 有 actor 或 system role。
- 有 precondition。
- 有 observable action。
- 有 expected outcome。
- 有明確 out-of-scope 或 non-goal。

## 風險訊號

- Scenario 中出現使用者沒有要求的功能。
- Given / When / Then 缺少可驗證 target。
- 一個 scenario 同時宣稱整個 workflow 成功。
- Then 只是「系統正常」而不是具體結果。

## 規則

單一 scenario pass 不等於 full feature correctness。這是 behavior scope governance，不是 test count 問題。
