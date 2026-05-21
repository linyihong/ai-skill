# Hallucination Boundaries

Hallucination boundary 將 model output 視為 candidate evidence，而不是 source-of-truth。

## High-risk Outputs

- 新檔案路徑、API、command 或 schema。
- Provider model capability。
- Current repo architecture claim。
- Completion / success claim。
- Runtime / generated surface claim。

## Required Handling

1. 查 current source。
2. 限縮 claim scope。
3. 標出 uncovered areas。
4. 執行或命名 validation。
5. 若 evidence local，只做 local claim。

## Forbidden

- 用 model confidence 替代 tests / lints / source。
- 用 summary 或 memory replay 直接覆蓋 current source。
- 從 model 名稱推測實際能力。
