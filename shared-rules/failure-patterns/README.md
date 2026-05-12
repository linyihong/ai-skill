# 失效模式 Patterns

本目錄存放跨 skill 可重用的 agent failure patterns。每個 pattern 記錄泛化後的 failure mode、trigger、required action、prevention gate 與 validation method。

當 [`failure-learning-system.md`](../failure-learning-system.md) 要求 promote 或查找 reusable failure pattern 時，先讀本索引。

| Pattern | Class | Status | 摘要 |
| --- | --- | --- | --- |
| [Correction loop bypass](correction-loop-bypass.md) | `validation-gap` | validated | 防止 agent 在使用者指出修正不完整時，只修當下文字，卻漏掉 `.agent-goals`、failure learning、linked updates、validation、commit/push/readback。 |
| [Entrypoint positioning drift](entrypoint-positioning-drift.md) | `validation-gap` | validated | 防止 agent 在命名或架構變更後，只更新次要連結或段落，卻留下 root title、opening paragraph 或主要入口 framing 過期。 |
| [Shared-rules architecture drift](shared-rules-architecture-drift.md) | `dependency-miss` / `validation-gap` | validated | 防止 agent 在架構重構後，只更新主要檔案（workflow、intelligence、analysis）卻漏掉 shared-rules/ 中的路徑參考同步。 |
| [Skill-local feedback bypass](skill-local-feedback-bypass.md) | `dependency-miss` / `validation-gap` | validated | 防止 agent 只補單一 skill 的 feedback lesson，卻沒有讀取全庫 failure-learning system 並沉澱 cross-skill prevention gate。 |
| [Source / mirror write drift](source-mirror-write-drift.md) | `source-mirror-drift` | validated | 防止 agent 更新 project-local tool mirrors 或 runtime copies，而不是 canonical source repo。 |
| [Tool config design without rule check](tool-config-design-without-rule-check.md) | `tool-strategy-gap` | candidate | 防止 agent 設計新工具配置時漏讀 `ai-tools/<tool>.md` 的現有規則，導致重複或邊界混淆。 |

## 維護

- 不要把 project-specific evidence 放進本目錄。
- 當 failure mode 可能跨 projects、tools、skills 或 agents 重演時，新增 pattern。
- 若 pattern 變成 skill-specific，把 lesson 移到該 skill 的 `feedback_history/`；只有 cross-skill trigger 仍有價值時，才從這裡連回。
- 若 pattern 變長，拆出獨立 examples，不要膨脹索引。

← [Back to shared rules index](../README.md)
