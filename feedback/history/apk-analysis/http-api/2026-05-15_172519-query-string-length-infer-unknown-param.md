> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-15 - Query String 長度分析推算未知參數值

Status: candidate

#### One-line Summary

當 Frida capture 只顯示 query string 總長度和 keys（不顯示 values）時，可透過總長度減去已知 key-value 長度來推算未知參數（如 service name）的長度，縮小暴力搜尋範圍。

#### Human Explanation

在分析 Dart AOT 應用的 Frida capture 時，`_generateEhHeader` hook 的 `x2` 參數是 query string 的 Dart String 物件，但 Frida 只輸出 `chars=N bytes=N`（總長度）和 `queryKeys=...`（所有 key 名稱），**不輸出實際的 value**。這在分析未知 API 時是常見的障礙。

解法是：如果知道大部分參數的值（如 `l=zh-cn`、`page=1`、`type=detail`），可以計算這些已知 key-value 的長度總和，再從 query string 總長度減去，得到未知參數（通常是 `service=XXX`）的 value 長度。這可以大幅縮小 service name 的暴力搜尋範圍。

例如：69-char query string 有 keys `l,page,service,type`。若 `l=zh-cn`（6 chars）、`page=1`（6 chars）、`type=detail`（11 chars），加上 3 個 `&` 分隔符，則 `service=XXX` 佔 69 - 6 - 6 - 11 - 3 = 43 chars，扣除 `service=`（8 chars），service name 長度為 35 chars。

#### Trigger

Frida `_generateEhHeader` hook 輸出 query string 的 `chars=N` 和 `queryKeys=k1,k2,k3`，但 values 被截斷或無法直接讀取。

#### Evidence

- Tool: Frida hook on `RequestInterceptor._generateEhHeader`
- Sanitized excerpt: `dartString chars=69 bytes=69 queryKeys=l,page,service,type serviceHash=468abf8fac324d8c`
- Evidence path: `<PROJECT_ROOT>/capture/short_drama_20260515_1255.log` evt=8

#### Generalized Lesson

當 Frida capture 只顯示 query string 長度和 keys 時，可透過以下公式推算未知參數值長度：

```
unknown_value_length = total_length
  - sum(known_key.length + 1 + known_value.length)  // +1 for '='
  - (num_keys - 1)                                   // '&' separators
  - (unknown_key.length + 1)                         // '=' for unknown key
```

若有多個未知參數，需聯立方程組求解。

#### Agent Action

1. 從 Frida capture 取得 query string 總長度（`chars=N`）和所有 keys（`queryKeys=...`）。
2. 從已知的 API 模式推測已知參數的值（如 `l=zh-cn`、`page=1`、`type=list`）。
3. 用上述公式計算未知參數（通常是 `service`）的 value 長度。
4. 用該長度過濾 service name candidates，只測試符合長度的候選。

#### Goal / Action / Validation

- Goal: 從有限的 Frida capture 資訊推測 API 的 service name 長度
- Action: 用 query string 總長度減去已知 key-value 長度
- Validation or reference source: 比對已知 API（如 LIST API 的 344-char query string）驗證公式正確性

#### Applies When

- Frida `_generateEhHeader` hook 輸出 `chars=N` 和 `queryKeys=...` 但 values 不可讀
- query string 中大部分參數的值可從其他來源推測（如固定值 `l=zh-cn`、`plat=android`）

#### Does Not Apply When

- query string 中所有參數的值都未知（無法建立方程式）
- Frida capture 直接輸出了完整的 query string values

#### Validation

用已知的 API（如 LIST API，344 chars，已知所有 values）驗證公式：計算結果應與實際 service name 長度一致。

#### Promotion Target

- `analysis/apk-analysis/frida-capture-analysis.md`

#### Required Linked Updates

- 無需連動更新。本 lesson 是獨立技巧。
