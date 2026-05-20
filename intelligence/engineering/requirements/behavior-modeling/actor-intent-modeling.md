# Actor Intent Modeling

**Status**: `candidate-intelligence`

## 判斷原則

Actor intent modeling 將「誰要什麼結果」與「系統如何實作」分離。Agent 不應把 inferred implementation 當成 user intent。

## 檢查問題

- Actor 是 user、admin、system、external provider 還是 background worker？
- Actor 是否有權執行該行為？
- Outcome 是狀態改變、訊息、資料回傳、錯誤處理還是 side effect？
- 哪些行為是明確要求，哪些只是推論？

## 防止

```text
user intent -> inferred behavior -> feature inflation
```

若 inferred behavior 會影響 scope、contract、security、cost 或 UX，必須標成 assumption 或 open question。
