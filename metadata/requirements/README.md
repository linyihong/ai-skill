# Requirements Metadata

`metadata/requirements/` 保存 requirements cognition、behavior scope、acceptance freshness 與 validation target 的可機讀 heuristics。這些資料預設是 metadata-only，不代表 runtime enforced。

## 目前檔案

| 檔案 | 用途 |
| --- | --- |
| [`behavior-governance-signals.yaml`](behavior-governance-signals.yaml) | 將 requirement contradiction、missing validation target、stale acceptance criteria 等信號分類。 |
| [`acceptance-quality-thresholds.yaml`](acceptance-quality-thresholds.yaml) | 定義 acceptance criteria quality threshold。 |

## 邊界

- 不要求所有專案使用 Gherkin runner。
- 不建立 universal requirement database。
- 不讓 runtime 理解 BDD syntax。
- 可 promotion 的只有壓縮後 runtime-lite signal。
