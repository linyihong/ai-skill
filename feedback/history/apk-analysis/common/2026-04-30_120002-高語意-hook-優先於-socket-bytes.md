# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md) (Section 2: Quick Start, step 6) and [`intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md`](../../../../intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md)

### 2026-04-30 - 高語意 hook 優先於 socket bytes

Status: promoted

#### One-line Summary

能 hook request/response 物件，就不要先從 socket bytes 開始拼。

#### Human Explanation

socket、TLS read/write、`send`/`recv` 事件很多，容易卡 App，也需要自己重組 HTTP、解壓縮、切分 body。高語意 hook 例如 request options、response interceptor、decrypt function，通常事件少、內容接近業務語意，更適合建立 API 文件與測試 fixture。

#### Trigger

低層 socket / TLS hook 事件量大，容易造成 App 卡頓，也需要自行重組 HTTP。

#### Generalized Lesson

優先找 request options、response interceptor、decode/decrypt function。只有在高語意點找不到或需要補證據時，才降到 socket / TLS 層。

#### Agent Action

看到 socket hook 卡頓、ANR、輸出爆量時，停止擴大低層 hook，改回靜態搜尋 request builder / interceptor / decoder，或縮小 hook 條件。

#### Promotion Target

已整理到 `WORKFLOW.md`。
