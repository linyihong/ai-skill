# Architecture Metadata

`metadata/architecture/` 保存 architecture fit、DDD adoption、overengineering、bounded context 與 domain architecture cognition 的可機讀 heuristics。這些資料預設是 metadata-only，不代表 runtime enforced。

## 目前檔案

| 檔案 | 用途 |
| --- | --- |
| [`architecture-fit-matrix.yaml`](architecture-fit-matrix.yaml) | 將 complexity signal 對應到 architecture strategy。 |
| [`ddd-adoption-signals.yaml`](ddd-adoption-signals.yaml) | 判斷 DDD Lite / Full DDD 的 adoption signal。 |
| [`overengineering-signals.yaml`](overengineering-signals.yaml) | 偵測 architecture inflation。 |
| [`bounded-context-heuristics.yaml`](bounded-context-heuristics.yaml) | 判斷 bounded context 是否成立。 |

## 邊界

- 本目錄不建立 runtime invariant。
- 本目錄不取代 `intelligence/engineering/architecture/` 的正文判斷。
- 若未來要 promotion 成 runtime-lite signal，必須另開 plan，確認 compiler / generated surface。
