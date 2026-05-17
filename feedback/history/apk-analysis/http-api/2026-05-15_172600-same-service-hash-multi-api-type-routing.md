> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-15 - 同一 Service Hash 對應多個 API，透過 type 參數區分

Status: candidate

#### One-line Summary

同一個 service hash 可能對應多個不同的 API endpoint，透過 query string 中的 `type` 參數（如 `type=list` vs `type=detail`）來路由到不同的 handler，query string 的長度和 keys 也會完全不同。

#### Human Explanation

在分析 API 路由時，常假設一個 service hash 唯一對應一個 API endpoint。但實務上，同一個 service hash 可以對應多個 API，透過 query string 中的 `type` 或其他路由參數來區分。

例如，`serviceHash=468abf8fac324d8c` 同時用於：
- **LIST API**：query string 344 chars，15 個 keys（`_sign,cate,device,...,type,uid,version`），`type=list`
- **DETAIL API**：query string 69 chars，4 個 keys（`l,page,service,type`），`type=detail`

兩個 API 使用完全不同的參數集，但共用同一個 service hash。這表示 service hash 可能對應到一個 controller class 或 handler group，而 `type` 參數決定具體呼叫哪個 method。

#### Trigger

Frida capture 顯示兩個不同的 API call 有相同的 service hash，但 query string 長度和 keys 完全不同。

#### Evidence

- Tool: Frida hook on `RequestInterceptor._generateEhHeader`
- Sanitized excerpt:
  - evt=5: `chars=344 queryKeys=_sign,cate,device,...,type,uid,version serviceHash=468abf8fac324d8c`
  - evt=8: `chars=69 queryKeys=l,page,service,type serviceHash=468abf8fac324d8c`
- Evidence path: `<PROJECT_ROOT>/capture/short_drama_20260515_1255.log` evt=5, evt=8

#### Generalized Lesson

當分析 Frida capture 中的 API 呼叫時：

1. **不要假設 service hash 唯一對應一個 API** — 同一個 hash 可能對應多個 endpoint。
2. **比對 query string 的 keys 和長度** — 如果兩個呼叫有相同 hash 但不同 keys/長度，它們是不同的 API。
3. **尋找路由參數** — 通常 `type` 參數（或其他枚舉值）決定具體路由。
4. **分別分析每個 API 的 response schema** — 即使共用 hash，response 結構也可能完全不同。

#### Agent Action

1. 收集所有具有相同 service hash 的 Frida capture entries。
2. 比對它們的 query string 長度、keys、以及 response schema。
3. 如果 keys 和長度不同，視為不同的 API endpoint。
4. 分別測試每個 API 的 service name candidates。

#### Goal / Action / Validation

- Goal: 正確識別同一 service hash 下的多個 API endpoint
- Action: 比對 query string 的 keys 和長度，尋找路由參數
- Validation or reference source: 用 live test 分別驗證每個 API 的 response schema

#### Applies When

- Frida capture 顯示多個 API call 有相同 service hash
- query string 的 keys 或長度在不同呼叫之間有顯著差異

#### Does Not Apply When

- 所有具有相同 service hash 的呼叫都有相同的 query string keys 和長度
- service hash 已經過完整逆向工程，確認唯一對應一個 handler

#### Validation

對每個疑似不同的 API endpoint 分別發送請求，確認 response schema 不同。

#### Promotion Target

- `analysis/apk-analysis/frida-capture-analysis.md`

#### Required Linked Updates

- 無需連動更新。
