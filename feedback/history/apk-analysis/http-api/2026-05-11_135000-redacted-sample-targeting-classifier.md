> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-11 - Redacted Sample Targeting Classifier

Status: candidate

#### One-line Summary

UI 盲抽低收益時，可用 disabled-by-default classifier 只輸出 value class 與 item index，先找候選樣本再回到 UI replay。

#### Human Explanation

列表、分類、留言、媒體或庫存頁常需要「挑到某種樣本」才會觸發目標 API 或非空 schema。若直接盲點 UI，容易多輪都落在空樣本；若把完整 decrypted payload 印出，又會暴露內容、id、user、URL 或計數值。較好的折衷是只把可用性欄位轉成類別，例如 `zero/nonzero/missing/other`、`available/unavailable`、`hasMedia/noMedia`，並保留列表索引，讓後續 UI replay 可以點候選 card。

#### Trigger

遇到 UI 盲抽多次沒有命中目標樣本，但 response schema 中有欄位可判斷哪個 item 比較可能觸發目標路徑。

#### Evidence

- Tool: Frida response/decrypt hook or schema-only decoder with a disabled-by-default classifier.
- Sanitized excerpt: `itemClassified=true zeroIdx=0,1 nonzeroIdx=2 missingIdx= otherIdx=` 這類索引摘要；不包含 raw id、標題、本文、留言、使用者、URL、token、完整 count 或 payload value。
- Evidence path: project docs may reference `<PROJECT_ROOT>/capture/...` artifacts and summarize only value classes / indexes.

#### Generalized Lesson

在授權 APK 分析中，sample targeting 可以先靠去敏 value-class 索引完成，再由已驗證的 UI replay 選取候選項。classifier 應預設關閉、短窗啟用、只輸出分類與 index，並明確標成 targeting aid，而不是資料擷取、欄位真值證明或 standalone replay parity。

#### Agent Action

下次遇到低收益 UI 盲抽時，先找 response 中是否存在可去敏分類的樣本選擇欄位。若有，新增或啟用 disabled-by-default classifier，跑短窗取得 class/index，再用 target package + feature-context guard 的 replay 點擊候選項。若 classifier 需要輸出 raw value 才有用，就不要把它放進 public/tracked docs。

#### Goal / Action / Validation

- Goal: 提高稀有樣本命中率，同時維持 raw-value redaction。
- Action: 把可用性欄位轉成 value classes and indexes, then replay a candidate through the app-owned UI/request/decrypt path.
- Validation or reference source: 候選 replay 必須通過 package/feature/detail guards，並由後續 API schema 證明目標樣本已命中。

#### Applies When

- 目標 API/schema 需要特定 list item、分類、媒體、留言或狀態樣本才會出現。
- Decrypted/schema hook 已有合法授權且可在短窗中只輸出分類摘要。
- UI replay 能用 index/visible card/selector 對齊候選項。

#### Does Not Apply When

- 需要保存 raw ids、raw content、完整計數值或使用者資料才能選樣。
- classifier 會改動 App 狀態、改寫 request、或繞過 App 自己的 signing/decrypt path。
- UI replay 無法驗證 target package / feature context，導致 candidate attribution 不可靠。

#### Validation

1. 語法或 hook 健康檢查通過。
2. classifier log 只含 class/index，不含 raw sensitive/content values。
3. follow-up replay 命中候選樣本並捕獲目標 API/schema。
4. project docs 保留具體 evidence；skill 只保留 generalized method。

#### Promotion Target

- `WORKFLOW.md`

#### Required Linked Updates

- 已同步 `WORKFLOW.md` 的 Redacted sample-targeting classifier rule。
- 已依 reusable-guidance-boundary 檢查：具體 App/category/capture 結論留在 project docs，本 lesson 只記錄通用方法。
