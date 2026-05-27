# Intelligence Layer Bypass via Tool Adapter（工具 Adapter 直寫繞過 Intelligence 層）

Status: validated
Class: `process-gap` / `knowledge-routing-miss`

## Trigger

Agent 取得跨工具可重用的 agent 行為洞見（設計原理、prompt 反應特性、架構模式）後，因為主題「關於某工具」而直接寫進 `ai-tools/<tool>.md`（P3 tool adapter），**knowledge-update-flow 的 Step 1 觸發條件從未被自問**，整個 master flow 沒被進入。

具體觸發信號：

- Commit type 為 `docs(<tool>):` 但內容含跨工具設計原理
- `ai-tools/<tool>.md` 出現「為什麼這樣設計」「實測發現行為規律」「prompt 模式強度比較」等段落
- 同一 PR 內無 `intelligence/` 新增，無 feedback lesson，無 failure pattern
- 使用者在後續 session 指出「這種知識應該放到 intelligence 然後引用」

## Failure Mode

P3 tool adapter 是工具專屬實作細節的容器，不是可重用知識的 canonical home。直接寫進去時：

1. **Intelligence atom 從未建立** — 其他工具的設計者無法從 `intelligence/` 找到跨工具設計原理
2. **Feedback lesson 從未寫入** — 可重用洞見不在 feedback pipeline 中
3. **Knowledge-update-flow 完全跳過** — Step 1 觸發條件被任務 framing 繞過，而非「sub-pipeline 替代 master」
4. **Rule-weight 倒置** — P2（intelligence 層）被 P3（tool adapter）搶佔，違反分層設計

## Risk

- 設計原理被鎖在工具專屬文件，跨工具複用困難
- 未來 agent 無法從 `intelligence/` routing 找到相關洞見（routing registry miss）
- 同類錯誤在其他工具 adapter 上無法被偵測（無 failure pattern coverage）

## Required Agent Action

寫任何 `ai-tools/` 段落前，先自問：

```
這個段落是否包含「為什麼這樣設計」「行為規律」「設計原理」？
  是 → 先建 intelligence atom，tool adapter 只引用
  否 → 工具專屬細節，可直接寫進 ai-tools/
```

1. **先問** 「是否跨工具可重用」，而非「主題是哪個工具」
2. **是 → 執行 knowledge-update-flow Step 1**，觸發 master flow，建 intelligence atom
3. **Tool adapter 只寫** `> 設計原理：見 [intelligence/...](...)`，不重複設計說明
4. 不得因為 commit type 是 `docs(...)` 就認為「沒有新知識」

## Tool Adapter vs Intelligence 層邊界

| 內容類型 | 正確位置 |
|---------|---------|
| 設計原理、行為模式、prompt 強度位階、架構 trade-off | `intelligence/engineering/<domain>/` |
| 特定工具的設定格式、CLI flag、版本限制、已知 bug | `ai-tools/<tool>.md` |
| 工具如何對應到通用設計模式 | `ai-tools/<tool>.md` 引用 intelligence atom |

## Prevention Gate

- **Tool adapter 段落審查**：含「為什麼」「實測發現」「行為規律」「設計原理」→ 移到 intelligence atom
- **Commit type 不是豁免**：`docs(<tool>):` commit 同樣需要 knowledge-update-flow Step 1 自問
- **PR review signal**：同一 PR 有 `ai-tools/` 修改但無 `intelligence/` 新增 → 懷疑 bypass

## Validation

符合下列條件時，此 pattern 已被防止：

- `ai-tools/<tool>.md` 中無獨立「為什麼這樣設計」段落（只有工具實作細節 + intelligence 引用連結）
- 對應 `intelligence/` atom 已建立，含工具中立設計原理
- Feedback lesson 已寫入 `feedback/history/<domain>/`
- Master flow 11 step 對照已在 commit message 中

## Related

- [`knowledge-update-flow-bypassed-by-sub-pipeline.md`](knowledge-update-flow-bypassed-by-sub-pipeline.md) — 相關但不同：sub-pipeline 替代 master（已進入 master 但走錯入口）；本 pattern 是「從未進入 master」
- [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md) — Master flow Step 1 觸發條件
- [`enforcement/rule-weight.md`](../rule-weight.md) — P2 intelligence 高於 P3 tool adapter

## Source

- 2026-05-27 session：commit `db9b2d0` 把 Bootstrap 三層架構設計原理寫進 `ai-tools/agent/claude.md`，commit type `docs(claude):`，無 intelligence atom / lesson / failure pattern。使用者在下一 session 明確指出後補救（commit `917f671`）。
- Corresponding feedback lesson: `feedback/history/development-guidance/common/2026-05-27_092801-intelligence-layer-bypass-via-tool-adapter.md`
