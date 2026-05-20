# Impact Journey Cross-Check

**Status**: `candidate-intelligence`

## 判斷原則

Impact Map 與 Customer Journey Map 必須交叉驗證。Impact Map 防止功能沒有 business target；Journey Map 防止 business target 忽略使用者真實痛點。

## 四個關鍵問題

| 問題 | 檢查 |
| --- | --- |
| Who 對得上嗎？ | Impact Map 的 actor / segment 是否真的出現在 Journey Map？ |
| 行為時機對了嗎？ | 期待改變的行為是否發生在正確 journey stage？ |
| 痛點影響目標嗎？ | Journey pain / blocker 是否真的影響 Impact Map 的 goal 或 metric？ |
| 資源放對地方了嗎？ | Feature investment 是否對準最高摩擦點，而不是做容易做的功能？ |

## 典型 mismatch

- `Who mismatch`：business target 指向 buyer，但 feature 設計給 operator。
- `Timing mismatch`：目標是 activation，但功能改善的是 post-retention workflow。
- `Pain mismatch`：使用者卡在 trust / permission，feature 卻增加 dashboard。
- `Investment mismatch`：高成本功能只解低頻 journey step。

## 輸出

交叉檢查後給出：

- `proceed`：Impact 和 journey 對齊，可進入 BDD-lite。
- `revise`：需要修改 Who / How / What 或 journey assumption。
- `reject`：feature 不支持目標或痛點。
- `ask_user`：缺少阻擋性 product / journey evidence。

## 邊界

這不是 runtime gate，也不是大型 product research 流程。它是 AI coding 前的 requirements alignment check，用於避免「做得很快但做錯問題」。
