# Coding Workflows

Coding workflow 將 model-aware strategy 轉成可驗證的 code edit 行為。重點是 patch scope、diff precision、tests / lints 與 rollback clarity。

## Boundaries

- 先讀 owner source 與 local conventions。
- Patch scope 應對應 single user goal。
- 不做 unrelated refactor。
- 遇到 generated files 先確認 source-of-truth。
- Commit 前必須 diff review。

## Validation Density

| Risk | Validation |
| --- | --- |
| Low | Lints 或 targeted test。 |
| Medium | Targeted test + affected package test。 |
| High | Source-backed review + broader tests + generated surface validation。 |
| Unknown | 先 downgrade autonomy，補 source / test evidence。 |

## Model-aware Adjustment

- Low context stability：縮小 patch，先 reread source。
- High hallucination risk：避免 invented API / path，先搜尋與驗證。
- Unknown tool reliability：避免 long automation，使用 dry-run 或手動分段。
