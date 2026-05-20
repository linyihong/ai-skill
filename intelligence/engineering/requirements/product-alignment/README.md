# Product Alignment Intelligence

`product-alignment/` 保存 product / requirements alignment 的判斷智慧。它在 BDD-lite 與 architecture cognition 之前，先確認需求是否對準 business impact 與真實 user journey。

## 核心

Impact Map 從目標往外推；Customer Journey Map 從使用者真實旅程往內拉。兩者交叉驗證，避免 AI 很快產出功能但偏離真正問題。

## 目前條目

| 文件 | 用途 |
| --- | --- |
| [`impact-map-alignment.md`](impact-map-alignment.md) | 用 Why / Who / How / What 檢查 product brief 是否有清楚 impact chain。 |
| [`customer-journey-validation.md`](customer-journey-validation.md) | 用 journey stage、pain point、情緒低谷與 blocker 檢查需求是否貼近使用者真實路徑。 |
| [`impact-journey-cross-check.md`](impact-journey-cross-check.md) | 將 Impact Map 與 Journey Map 互相驗證，找出 target / timing / pain / feature mismatch。 |

## 邊界

- 不取代 product discovery、user research 或 business decision。
- 不要求每個 bug fix 都跑 Impact Map / Journey Map。
- 不把 product alignment promotion 成 runtime primitive；只可能產生 `product_goal_mismatch`、`journey_pain_mismatch` 或 `feature_without_impact` 類壓縮信號。
