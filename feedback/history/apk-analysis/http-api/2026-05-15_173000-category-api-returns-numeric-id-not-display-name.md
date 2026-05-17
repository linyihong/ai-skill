> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-15 - Category API 返回 Numeric ID，LIST API 需要 Numeric ID 而非 Display Name

Status: candidate

#### One-line Summary

某些 API 的 category 參數需要 numeric ID（如 `cate=41`）而非 UI 上顯示的 display name（如 `cate=全部`）。必須先呼叫 categories API 取得 ID mapping，再將 numeric ID 用於 LIST API。

#### Human Explanation

在分析短劇列表 API 時，UI 上顯示的分類名稱是中文（如「全部」、「穿越」、「古裝」），直覺上會認為 API 的 `cate` 參數也接受這些名稱。但實際上 API 需要的是 numeric ID（如 `41`、`42`、`43`）。

這些 numeric ID 來自 categories API 的回應中的 `cate_info` array，格式為 `[{"id":41,"name":"全部"}, ...]`。LIST API 的 `cate` 參數必須使用 `id` 欄位的值，而非 `name`。

這是一個常見的 API 設計模式：UI 顯示 human-readable name，但 API 傳輸 machine-readable ID。如果直接用 display name 呼叫 LIST API，會得到空結果（`"info":[]`）而非錯誤。

#### Trigger

LIST API 回傳 HTTP 200 但 `data.info` 為空陣列，而 categories API 回傳的 `cate_info` 中有對應的分類名稱。

#### Evidence

- Tool: Direct HTTP call to production API
- Sanitized excerpt:
  - With `cate=全部`（Chinese name）: `{"ret":200,"data":{"code":0,"msg":"","info":[]}}`
  - With `cate=41`（numeric ID from cate_info）: 回傳正常結果
- Evidence path: `<PROJECT_ROOT>/capture/short_drama_20260515_1255.log` evt=6 (LIST response schema: `data.info{list: array[10]}`)

#### Generalized Lesson

當 LIST API 回傳空結果時：

1. **檢查 category 參數格式** — 確認 API 需要 numeric ID 還是 display name。
2. **先呼叫 categories API** — 取得 `cate_info` 或類似結構中的 ID mapping。
3. **從 mapping 中提取 numeric ID** — 使用 `id` 欄位而非 `name` 欄位。
4. **用 numeric ID 重新測試 LIST API**。

更一般化：任何有「分類」概念的 API，都應先確認分類參數的格式（numeric ID vs string name），並從 categories API 的 response 中提取正確的格式。

#### Agent Action

1. 如果 LIST API 回傳空結果，先檢查是否使用了 display name 作為 category 參數。
2. 找出對應的 categories API endpoint。
3. 解析 categories API 的 response，提取 `cate_info` 或類似 array 中的 `id` 欄位。
4. 使用 numeric `id` 重新呼叫 LIST API。

#### Goal / Action / Validation

- Goal: 正確使用 category 參數呼叫 LIST API
- Action: 先呼叫 categories API 取得 ID mapping，使用 numeric ID
- Validation or reference source: 用 numeric ID 呼叫 LIST API 應回傳非空結果

#### Applies When

- API 有 category/filter 參數
- categories API 回傳 `id` 和 `name` 兩個欄位
- 用 display name 呼叫 LIST API 回傳空結果

#### Does Not Apply When

- API 直接接受 display name 作為 category 參數
- 沒有對應的 categories API

#### Validation

用 numeric ID 和 display name 分別呼叫 LIST API，比對結果是否不同。

#### Promotion Target

- `analysis/apk-analysis/http-api/api-analysis-checklist.md`

#### Required Linked Updates

- 無需連動更新。
