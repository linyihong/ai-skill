# Impact Map Alignment

**Status**: `candidate-intelligence`

## 判斷原則

Impact Map 用 `Why / Who / How / What` 檢查需求是否能從 business goal 連到具體 feature investment。AI 產生功能前，必須知道該功能試圖改變誰的什麼行為，以及為什麼這會影響目標。

## 檢查項

| Impact Map 欄位 | 問題 |
| --- | --- |
| Why | 目標或 metric 是什麼？是否可觀察、可驗證、有時間範圍？ |
| Who | 哪個 actor、segment、role 或 customer type 會造成 impact？ |
| How | 期待這個 actor 改變什麼行為？ |
| What | 哪些 feature、experiment、message、workflow 或 intervention 支援該行為改變？ |

## 風險訊號

- `What` 很多，但 `Why` 不清楚。
- `Who` 只是「所有使用者」。
- `How` 是 feature delivery，而不是 user behavior change。
- metric 無法追到 behavior 或 journey stage。

## 行動

若 Impact Map 缺口影響 scope、priority、acceptance criteria 或 architecture investment，先標記 `open question` 或回到 product alignment，不直接開始 BDD / implementation。
