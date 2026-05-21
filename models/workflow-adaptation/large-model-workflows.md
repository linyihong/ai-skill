# Large / High Reasoning Workflows

Large 或 high-reasoning profile 適合跨層規劃、migration、recovery、architecture tradeoff 與多 source synthesis。它仍需遵守 context minimality，不代表可以無限制讀取。

## Workflow Shape

1. 讀 primary source 與 required dependencies。
2. 建立 task class、assumptions、validation target。
3. 比對 layer responsibility 與 linked updates。
4. 分批 edit，避免同時改多個 owner group。
5. 每批完成後執行 relevant validation。

## 適合任務

- Architecture planning / governance design。
- Runtime / knowledge / metadata integration。
- Legacy migration / deprecation。
- Evidence contradiction / recovery。
- Multi-document linked updates。

## Guardrails

- 高 reasoning 不取代 evidence。
- Cross-layer work 必須說明 source-of-truth。
- Long context 後要 reread changed sources。
- 若輸出只是建議，明確標記未驗證區域。
