# Integration Pressure Analysis

**Status**: `candidate-intelligence`

## 判斷原則

Integration pressure 來自外部模型、第三方 API、legacy schema、跨 bounded context 協調與事件一致性需求。它決定是否需要 ACL、domain events 或更強的 context boundary。

## 訊號

| 訊號 | 可能策略 |
| --- | --- |
| 外部模型語言污染內部語言 | anti-corruption layer |
| 多 provider 需要統一內部語意 | adapter + ACL |
| 下游需要業務事實通知 | domain events |
| 同步呼叫造成 context coupling | asynchronous coordination |
| 外部狀態不可靠 | explicit state mapping + recovery policy |

## 避免

不要因為有外部 API 就建立 event-driven architecture。先判斷語意污染、協調需求與一致性成本。
