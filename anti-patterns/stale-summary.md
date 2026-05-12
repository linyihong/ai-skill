# Stale Summary

## 症狀

- Summary 描述的內容與實際 source 不一致
- Agent 根據 summary 做決策但 source 已變更
- Summary 的 `last_checked` 日期過舊

## 根本原因

1. **Source 更新後未 refresh summary**。
2. **無自動 revalidation 機制**。
3. **Summary 被當作 source-of-truth**。

## 影響

- Agent 基於過期資訊做決策
- 需要手動比對 summary 與 source
- 降低對 summary layer 的信任

## 預防

1. Source 更新後執行 `scripts/refresh-knowledge-runtime.rb`。
2. Summary 加入 `last_checked` 與 `validation_signal` 欄位。
3. Summary 只做 routing，不做精確判斷。
4. 使用 `knowledge/runtime/refresh-policy.yaml` 的 revalidation 流程。

## 檢測

- `last_checked` 超過 7 天
- Source 的 git commit 比 summary 新
- Summary 的 links 無法解析

## 恢復

1. 重新讀取 canonical source。
2. 更新 summary 內容。
3. 執行 validation scripts 確認一致性。

## 相關 Guards

- `knowledge/runtime/refresh-policy.yaml`
- `scripts/refresh-knowledge-runtime.rb`
- `scripts/validate-knowledge-runtime.rb`
