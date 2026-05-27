> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-27 — 可重用設計洞見不得直接寫入工具 adapter；必須先提到 intelligence 層再引用

Status: validated

#### One-line Summary

當 agent 取得跨工具可重用的 agent 行為洞見（如「多層強制執行設計模式」）時，不得因為主題「關於某工具」就直接寫進 `ai-tools/<tool>.md`（P3 tool adapter）；必須先依 knowledge-update-flow 提到 `intelligence/` 層，再由 tool adapter 引用 intelligence atom。

#### Human Explanation

本次問題的根本是**知識分類錯誤**，而非「sub-pipeline 替代 master flow」。Agent 看到「關於 Claude 的 empirical 觀察」就直覺跳到 `ai-tools/agent/claude.md`，把 commit type 設為 `docs(claude):`，結果：

1. knowledge-update-flow 的 Step 1 觸發條件（「是否有新的可重用知識？」）**從未被自問**
2. 整個 11-step master flow 沒被觸發，不是「sub-pipeline 替代 master」，而是「完全沒進入 master」
3. 洞見被鎖在 P3 工具 adapter 中，其他工具的設計者無法從 intelligence 層找到它

跨工具可重用的知識 vs. 工具專屬細節的邊界：

| 類型 | 正確位置 | 例子（本次） |
|------|---------|------------|
| 設計原理、行為模式、prompt 強度位階、primacy effect | `intelligence/engineering/agent-architecture/` | 多層強制執行設計（L1/L2/L3）、prompt 強度 4 級、primacy effect |
| 特定工具的檔案路徑、格式要求、版本限制 | `ai-tools/<tool>.md` | `.claude/settings.json` nested 格式、SessionStart matcher 字串 |

#### Trigger

2026-05-27：agent 在 empirical 測試後把「Bootstrap 三層架構設計原理」直接寫進 `ai-tools/agent/claude.md`，commit type 為 `docs(claude):`。使用者在下一個 session 明確指出「這種知識應該放到 intelligence 然後引用」且「為什麼之前沒照 knowledge-update-flow 跑」，才觸發補救。

#### Evidence

- Commit `db9b2d0`: `docs(claude): record three-layer bootstrap architecture` — 整個洞見寫進 `ai-tools/agent/claude.md`，無 intelligence atom，無 feedback lesson，無 failure pattern
- 補救 commit `917f671`: `refactor(intelligence): extract multi-layer-enforcement atom from claude.md` — 抽取至 `intelligence/engineering/agent-architecture/multi-layer-enforcement.md`，tool adapter 改為引用
- Evidence path: `<AI_SKILL_REPO>/ai-tools/agent/claude.md`、`<AI_SKILL_REPO>/intelligence/engineering/agent-architecture/multi-layer-enforcement.md`

#### Generalized Lesson

**在工具 adapter 寫任何段落之前，必須先問：「這個洞見的適用範圍是這個工具，還是跨工具可重用？」**

判斷句：

| 問題 | 判斷 | 行動 |
|------|------|------|
| 這個設計原理、行為特性、trade-off 只適用於這個工具嗎？ | 否（跨工具） | 先建 intelligence atom，tool adapter 只留引用 |
| 這個設定格式、版本限制、CLI flag 只在這個工具上發生？ | 是（工具專屬） | 直接寫進 tool adapter，無需 intelligence atom |
| 任務 commit type 是 `docs(<tool>):`？ | → 先自問上兩個問題再決定 | 不得因為 commit type 是 docs 就跳過 knowledge-update-flow |

#### Agent Action

1. 寫進 `ai-tools/` 之前先問：「這個段落是否包含跨工具可重用的設計洞見？」
2. 若是，**先執行 knowledge-update-flow Step 1**（觸發 master flow），再建 intelligence atom，最後 tool adapter 只寫 `> 設計原理：見 [intelligence/...](...)`
3. 不得因為 commit type 是 `docs(...)` 就認為「沒有新知識」，docs commit 同樣可能包含可重用洞見

#### Goal / Action / Validation

- Goal: 確保可重用的 agent 行為洞見永遠先到 intelligence 層，tool adapter 只含工具專屬實作細節
- Action: 寫 tool adapter 前先問「是否跨工具可重用」，是→先建 intelligence atom
- Validation: tool adapter 中不應有「為什麼這樣設計」的段落，只有「這個工具如何實作」；設計原理應在 intelligence atom，有引用連結

#### Applies When

- Agent 即將在 `ai-tools/<tool>.md` 寫入「為什麼這樣設計」「實測發現行為規律」「prompt 模式強度比較」等內容
- Agent 做 empirical 測試後要記錄觀察結果
- Commit type 被設為 `docs(<tool>):` 但內容包含跨工具設計原理

#### Does Not Apply When

- 純粹的工具版本說明、CLI 格式要求、環境設定
- bug 修正記錄（不含設計原理）

#### Validation

下次同類任務（empirical 觀察寫進工具 adapter）執行時，tool adapter 應只有引用連結指向 intelligence atom，無「為什麼」段落；如有「為什麼」段落，視為 lesson 未生效。

#### Promotion Target

- `enforcement/failure-patterns/intelligence-layer-bypass-via-tool-adapter.md`（新增 failure pattern，跨 skill/tool 重演風險高）

#### Required Linked Updates

- Step 6（Intelligence Extraction）：**是** — 已建 `intelligence/engineering/agent-architecture/multi-layer-enforcement.md`（commit `917f671`）
- Step 7（Failure Learning）：**是** — 本 lesson → promotion target `enforcement/failure-patterns/intelligence-layer-bypass-via-tool-adapter.md`（本輪新建）
- `feedback/history/development-guidance/README.md`：common count 32 → 33
- `enforcement/failure-patterns/README.md`：加入新 pattern 索引列
