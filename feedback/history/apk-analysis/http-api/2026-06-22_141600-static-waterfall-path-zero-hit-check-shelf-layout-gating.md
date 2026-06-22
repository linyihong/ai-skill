> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Static waterfall pagination path with zero hits — check hall shelf layout gating

Status: candidate

#### One-line Summary

靜態 DEX 同時存在 **hall info** 與 **waterfall pagination** path，但 Home deep scroll capture 對 pagination path **0 hit** 時，不要判定 API 不存在；應對照 **analytics scene / shelf type**（grid vs waterfall）與 info response 的 layout 配置，waterfall pagination 可能是 **server layout gated** 的可選分支。

#### Human Explanation

Feed 首頁常有多種 shelf renderer（grid、waterfall、rank tab）。靜態 path catalog 會列出所有歷史/AB 分支 endpoint；實際帳號可能只走 grid + info 一次載入。Deep scroll 只驗證「當前 layout 的分頁策略」，不能證明 waterfall path 全局不存在。常見誤判：scroll 沒打出 waterfall → hook 或 UI 自動化錯了。

#### Trigger

- 靜態有 waterfall pagination + hall info path
- Info/hall config API 已命中；waterfall pagination 0 hit
- Telemetry 或 log 出現 gridShelf / recommendedGrid 類 scene 名
- DEX 有並存的 grid vs waterfall UI style 字串

#### Evidence

- Tool: Frida REQ + analytics JSON（需過濾）
- Sanitized excerpt: scene 名模式 `*Grid*` vs `*Waterfall*`；info hit、pagination 0 hit 同窗口
- Evidence path: `<PROJECT_ROOT>/docs/domain-baseline.md`、`<PROJECT_ROOT>/api/dynamic-*.md`

#### Generalized Lesson

```text
waterfall pagination 0 hit？
  1. 確認 hall/info config API 已命中
  2. 讀 telemetry shelf/scene type
  3. grid layout → 分頁可能在 info 一次帶回或不同 path
  4. waterfall pagination → 需 server 配置 waterfall shelf 或 AB 帳號
  5. 文件標「layout-dependent optional」，非 W2 硬性 gate
```

#### Agent Action

1. Project baseline 記錄「已觀察 layout + 已命中 path 集合」。
2. 第二帳號/forced layout 驗證前，不 promote 成「App 無 waterfall API」。
3. Ai-skill 不寫 hall_id、host、具體 path 字串。

#### Goal / Action / Validation

- Goal: 避免 layout 分支誤判為 capture 失敗。
- Action: execution-flow 區分「static catalog path」與「layout-observed path」。
- Validation: 文件同時記錄 info hit + layout type + pagination hit count。

#### Applies When

- 內容/電商 App 多 shelf 首頁
- 靜態 path 多於動態 hit

#### Does Not Apply When

- Waterfall pagination 已穩定命中
- 首頁單一 infinite scroll endpoint（無 shelf 分支）

#### Validation

- 同 session path hit 表 + telemetry scene 欄位寫入 project baseline

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §hall/feed capture
- `analysis/apk/traffic-triage.md` §static vs dynamic path

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
