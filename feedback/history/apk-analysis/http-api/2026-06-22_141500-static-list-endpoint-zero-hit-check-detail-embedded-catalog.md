> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Static chapter-list path with zero dynamic hits — check detail response for embedded catalog

Status: candidate

#### One-line Summary

靜態 DEX 列出 standalone **chapter/episode list** path，但動態 capture 在「開詳情 → 播放器 → catalog UI」路徑上 **0 hit** 時，不要先假設 hook 失敗；應檢查 **book/detail API 回應是否已內嵌 `catalogList` / `chapter_list` 陣列**，並以 Gson/converter hook 的實際 JSON 作 SDK schema 真相來源。

#### Human Explanation

短劇/連載類 App 常把集數 metadata 嵌在 detail response，standalone list API 僅用於 refresh、互動分支或 stale cache。常見誤判：靜態 path catalog 有 list endpoint → capture 沒看到 → 判定「list 未知」。實務上 detail JSON 可能已含 `chapter_id`、display name、order、duration、lock flag 等欄位；UI catalog drawer 可能只 re-fetch detail（第三參數 refresh flag），而非打獨立 list。

#### Trigger

- DEX/strings 有 chapter-list path，Frida REQ 0 hit
- Detail/book API 已在 log 中命中
- Gson/Retrofit converter hook 可用
- 使用者打開播放器 catalog / episode bar 前後仍無 list path

#### Evidence

- Tool: OkHttp REQ + Gson.fromJson hook（或等效 converter hook）
- Sanitized excerpt: detail response 含 array 欄位（catalog/chapter 語意）；standalone list path 同 session 0 次
- Evidence path: `<PROJECT_ROOT>/api/dynamic-*.md`、`<PROJECT_ROOT>/capture/`（gitignore 依專案）

#### Generalized Lesson

```text
list endpoint 0 hit？
  1. grep detail response JSON 是否含 catalogList / chapter_list / episodes[]
  2. 對照 UI：catalog drawer 是否只 refresh detail API
  3. standalone list 標為 optional / refresh-only，直到第二路徑驗證
  4. SDK schema 以動態 JSON 為準，靜態 path 目錄為索引而非唯一契約
```

#### Agent Action

1. Capture checklist 增加「detail embedded catalog 欄位表」寫入 project docs。
2. Ai-skill 只寫判斷樹；endpoint 名、bean 名、Retrofit tag 留 project。
3. 勿因 0 list hit 阻塞 playback/content capture 進度。

#### Goal / Action / Validation

- Goal: 避免錯誤的「episode list unknown」結論與重複 UI 自動化。
- Action: project `paths.md` 將 standalone list 標註「optional until observed」。
- Validation: ≥1 run 中 detail JSON 含可播放所需的 chapter id + lock flag。

#### Applies When

- 內容 App：detail + chapter content 分離的 REST 族
- 靜態多 path、動態 list path 未觀察

#### Does Not Apply When

- List path 已穩定命中（以 REQ 為準）
- 完全無 detail API（純 GraphQL / 單一 feed endpoint）

#### Validation

- 同 session：detail path hit + Gson catalog 欄位 + list path 計數對照表

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §API inventory vs dynamic truth
- `analysis/apk/traffic-triage.md` §response schema

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查；target incident 證據僅在 `<PROJECT_ROOT>/`
