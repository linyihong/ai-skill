> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-01 - API field documentation after analysis

Status: promoted

#### One-line Summary

分析完 API 後要文件化 HTTP headers、request、response，並逐字段說明 type、meaning、required/source；截圖只作為 UI 觸發輔助證據。

#### Human Explanation

APK 分析常會先追到 endpoint、解密點或 response JSON，但如果沒有把 API 寫成可讀文件，後續實作 SDK、重放測試或交接時仍然需要重新看 hook log。可靠的 API 文件應該能回答：這支 API 從哪個 UI 操作來、需要哪些 header、request 每個字段怎麼來、response wrapper 與 inner payload 每個字段代表什麼。

截圖可以幫助理解 tab、screen 與操作來源，但它不能取代 HTTP/API 文件。若 UI binding 尚未完成，先把 UI path 標成 unknown 或低信心，等核心 API 穩定後再補截圖與操作時間窗。

#### Trigger

- 已取得 HTTP request object、MITM 明文、hook log、decrypted response 或 replay fixture。
- Endpoint 已知，但文件只有 path 或 raw log，缺少 headers/request/response 字段說明。
- 需要把 APK 分析結果轉成 SDK mapping、contract test、API reference 或交接文件。

#### Evidence

- Tool: MITM export, Frida request/response hook, decrypted fixture, replay log.
- Sanitized excerpt: `POST /<path>` with header names, request field shapes, response top-level keys/types, and UI path if known.
- Evidence path: 專案 API/reference docs；reusable skill 只保存模板與去敏方法。

#### Generalized Lesson

After an API has been observed or decoded, the analysis is not complete until the project documentation records HTTP headers, request query/body fields, response headers, response wrapper, inner payload fields, and validation evidence. Each field should include type/shape, meaning, required/optional, source or derivation, sensitivity, and notes. Screenshots may support UI attribution but do not replace API field analysis.

#### Agent Action

下次 agent 完成 API 分析時，應主動檢查：

- API 文件是否已回填，不只留在 hook/MITM log。
- HTTP headers 是否記錄名稱、用途、來源與是否敏感。
- Request query/body 是否逐字段分析 type、meaning、required、source、是否參與簽章/加密。
- Response wrapper 與 decrypted/inner payload 是否逐字段分析 type、meaning、optional/nullability。
- UI path / screenshot 是否只作為輔助 attribution，未確認時標示 confidence。

#### Applies When

- HTTP/HTTPS API、local proxy HTTP route、Dio/OkHttp/WebView/Cronet/native client 可觀測。
- 已取得 schema-only response、decrypted JSON、request metadata 或 replay fixture。
- 分析成果需要被人或測試重用。

#### Does Not Apply When

- 仍只知道 pcap SNI/IP/timing，沒有 path/header/request/response schema。
- 授權範圍不允許保存或文件化特定 API 結論；此時只記方法，不寫 target facts。

#### Validation

- 任一 endpoint 文件可讓讀者知道 method/path、auth/header、request 字段、response 字段、證據來源與 UI attribution 信心。
- 去敏檢查通過：沒有 token、cookie、device id、個資、私有 host 或可重放 URL value。
- 至少有 hook/MITM/replay/fixture 之一支撐 request/response 對齊。

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`
- `SKILL.md`

#### Required Linked Updates

- 已更新 `SKILL.md` durable assets，加入 HTTP/API docs with headers/request/response fields。
- 已更新 `WORKFLOW.md`，新增 API 文件化步驟與分析結束條件。
- 已更新 `DOCUMENTATION.md`，擴充 API/Schema 模板與文件要求。
- 已更新 `feedback_history/README.md` 索引。
