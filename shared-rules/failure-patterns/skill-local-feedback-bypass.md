# Skill-Local Feedback Bypass（只做 skill-local 回饋）

Status: validated
Class: `dependency-miss` / `validation-gap`

## Trigger

當使用者指出「為什麼沒有更新規則」、「有沒有讀錯誤回饋系統」、「把剛剛案例加上去」、「這不是只補 skill lesson」等語意，且 agent 剛剛只更新了單一 skill 的 `feedback_history/`、`WORKFLOW.md`、`SKILL.md` 或工具 mirror 時，使用此 pattern。

## Failure Mode

Agent 把 feedback 視為單一 skill 的技巧沉澱，只補了局部 lesson 或 workflow，卻沒有讀取並執行全庫 failure-learning loop。結果是：技術技巧被記錄了，但「agent 為什麼漏掉回饋、下次如何防止只修局部」沒有變成 cross-skill prevention gate。

## Risk

- 同類 agent 失效會在其他 skill、project 或工具中重演。
- 使用者以為「錯誤回饋系統」已經吸收案例，但實際只有 skill-local 方法被更新。
- `feedback_history/` 與 `shared-rules/failure-patterns/` 之間失去分工：前者記技巧，後者記 agent failure prevention。
- 後續 agent 只會看到新增技巧，不會看到漏觸發 feedback 的 root cause 與 validation gate。

## Required Agent Action

1. 停止繼續專案工作或新增更多 skill-local lessons。
2. 讀取 `shared-rules/failure-learning-system.md`、`shared-rules/feedback-lessons.md`、`shared-rules/reusable-guidance-boundary.md`、`shared-rules/linked-updates.md` 與 `shared-rules/failure-patterns/README.md`。
3. 分類本次失效：通常是 `dependency-miss`（未讀 failure-learning dependency）加上 `validation-gap`（未驗證 durable prevention）。
4. 若問題可能跨 skill / project / agent 重演，新增或更新 `shared-rules/failure-patterns/`。
5. 再回到 skill-local feedback：保留技術 lesson，但明確區分「domain / technique lesson」與「agent failure pattern」。
6. 依 linked updates 更新索引、相關 shared rule 或 affected skill entrypoint。
7. 完成 canonical repo diff review、去敏搜尋、commit、push、必要 tool sync、readback 與 final status。

## Prevention Gate

當使用者指出回饋缺口時，agent 在宣稱已修正前必須能回答：

| Check | Required answer |
| --- | --- |
| Feedback layer | 這是 skill-specific technique lesson、cross-skill failure pattern，還是兩者都需要？ |
| Required rules read | 是否已讀 `failure-learning-system.md` 與 `feedback-lessons.md`，而不是只讀某個 skill 的 `SKILL.md`？ |
| Root cause | 這次為什麼漏掉：把 project-specific implementation 誤判為非 reusable、只做 local feedback、沒有開 writeback transaction，還是未讀 failure system？ |
| Durable prevention | 防止重犯的 gate 寫在哪個 shared rule、failure pattern、skill workflow 或 tool doc？ |
| Validation | 是否已更新索引、檢查去敏、commit/push/readback，並確認 runtime mirror 是否需要同步？ |

## 驗證

此 pattern 已套用時，應可反查：

- `shared-rules/failure-patterns/` 有對應 pattern 或既有 pattern 被更新。
- Skill-local `feedback_history/` 只保存 technique / domain lesson，不承擔 cross-skill agent failure prevention。
- `shared-rules/feedback-lessons.md` 或 affected skill rule 明確提醒：agent mistake 或 close-loop gap 需走 failure-learning system。
- Canonical `<AI_SKILL_REPO>` 完成 commit、push、readback；若有 runtime mirror，已同步或明確標記不適用。

## Linked Rules

- [`../failure-learning-system.md`](../failure-learning-system.md)
- [`../feedback-lessons.md`](../feedback-lessons.md)
- [`../reusable-guidance-boundary.md`](../reusable-guidance-boundary.md)
- [`../linked-updates.md`](../linked-updates.md)
- [`../dependency-reading.md`](../dependency-reading.md)

## Linked Validation Scenarios

- `validate_directory_structure` — 檢查 `feedback/history/<domain>/` 目錄結構是否正確，防止 skill-local feedback bypass

← [Back to failure patterns](README.md)
