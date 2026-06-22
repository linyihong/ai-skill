> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Standalone HTTP JSON field names may differ from Gson bean paths

Status: candidate

#### One-line Summary

Frida Gson hook 見 `PlayletEntity.catalogList[]`，但 **standalone signed HTTP** `getBookDetail` 回應可能是 **`data.chapterList.chapter_lists[]`**（不同巢狀鍵名）。Downloader/SDK 須以 **實際 wire JSON** 為準，不可只複製 Gson 類別欄位名。

#### Human Explanation

`141500` 處理「list endpoint 0 hit → detail 內嵌 catalog」。實作 downloader 時還會遇到：同一業務資料在 wire JSON 與 Gson 反序列化類名不一致（蛇形/巢狀 wrapper）。若 SDK 只解析 `catalogList` 會得到空陣列，但 API `code=0`。規則：hybrid client 成功後，用實際 response keys 更新 parser；Gson log 僅作結構參考。

#### Trigger

- Gson hook 有 `catalogList[]`
- Signed POST `getBookDetail*` 回 `chapterList` / `chapter_lists` wrapper
- Downloader `list_episodes()` 回空但 API success

#### Evidence

- Tool: hybrid signed POST + Gson capture diff
- Sanitized excerpt: wire `data.chapterList.chapter_lists[]` with `is_lock` / `chapter_id`
- Evidence path: `<PROJECT_ROOT>/scripts/sign/hybrid_client.py`

#### Generalized Lesson

```text
Episode list for SDK:
  1. Prefer wire JSON from signed client you control
  2. Gson class/field names = hint only
  3. Log top-level data keys on first success
  4. Cross-check 141500 embedded-catalog rule still applies (detail not list endpoint)
```

#### Agent Action

1. Project parser 註明 wire path vs Gson bean。
2. Ai-skill 不寫 package/endpoint 真值。

#### Goal / Action / Validation

- Goal: downloader 不因欄位名不一致空跑。
- Validation: `list_episodes` 非空且 `is_lock` 可 filter。

#### Applies When

- Building standalone client from dynamic Gson logs
- Detail API returns nested chapter array

#### Does Not Apply When

- Wire JSON already matches Gson field names
- Using in-app relay only (no parser)

#### Validation

- Documented wire key path in project client code

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §SDK response parsing

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 交叉引用 `141500`
- 已依 sanitization / reusable-guidance-boundary 自查
