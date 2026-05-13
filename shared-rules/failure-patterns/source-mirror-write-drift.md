# Source / Mirror Write Drift（source / mirror 寫入漂移）

Status: validated
Class: `source-mirror-drift`

## Trigger

當使用者要求 agent 更新 rules、skills、tool setup、feedback lessons、templates 或 OS guidance，且可見檔案位於 project `.cursor`、`~/.cursor`、generated bundle、runtime copy 或其他 tool mirror 時，使用此 pattern。

## Failure Mode

Agent 編輯目前工具看得到的 copy，而不是 canonical `<AI_SKILL_REPO>` source repository。當前 session 看起來已修好，但 reusable knowledge 不會可靠 persist、sync、commit，或傳播到其他 projects 和 tools。

## Risk

- 下一個 agent 讀到 stale canonical rules，並重複同一錯誤。
- Project-local `.cursor` copy 與 shared knowledge base 分歧。
- Agent 誤以為 mirror 是 source，跳過 tool sync、commit、push 與 readback gates。
- Private project paths 或 tool-local details 可能洩漏進 reusable docs。

## Required Agent Action

1. 一旦懷疑 source/mirror confusion，停止廣泛編輯。
2. 定位 `<AI_SKILL_REPO>`，並確認它是 git root。
3. 在 `<AI_SKILL_REPO>` 檢查 `git status --short --branch`。
4. 在 canonical source file 套用 reusable change。
5. 將 project `.cursor`、`~/.cursor/skills*`、`~/.cursor/shared-rules`、`~/.cursor/bundles/*` 與 generated bundles 視為 deployment 或 runtime surfaces。
6. 只有 source changes 正確後，才執行 configured tool sync。
7. 宣稱完成前，commit、push、reread updated entries，並確認 clean status。

## Prevention Gate

第一次寫入 rules、skills、templates、feedback lessons 或 tool deployment paths 前，agent 必須能回答：

| Check | Required answer |
| --- | --- |
| Canonical repo | 哪個路徑是 `<AI_SKILL_REPO>`，且它是否為 git root？ |
| Current file role | 目前檔案是 source、project config、tool mirror、runtime copy，還是 generated output？ |
| Source edit | 會先編輯哪個 canonical source file？ |
| Sync strategy | 目前是 reference-only、symlink/bundle，還是 copy snapshot？ |
| Close loop | 會跑哪些 diff review、sync、commit、push、readback 與 clean-status checks？ |

若任何答案未知，不要編輯 mirror。先讀 [`dependency-reading.md`](../dependency-reading.md) 與相關 `ai-tools/` 文件。

## 驗證

符合下列條件時，此 pattern 已被驗證：

- Canonical source file 包含 reusable change。
- Mirror 或 runtime paths 未變更、symlinked to source，或只由 configured sync 更新。
- 已在 `<AI_SKILL_REPO>` 檢查 `git diff` 與 `git status --short --branch`。
- 變更影響 repository 時，已 commit、push、read back。
- Final response 將 project-local mirror changes 與 source update 分開命名。

## Linked Rules

- [`failure-learning-system.md`](../failure-learning-system.md)
- [`dependency-reading.md`](../dependency-reading.md)
- [`tool-neutral-documentation.md`](../tool-neutral-documentation.md)
- [`linked-updates.md`](../linked-updates.md)

← [Back to failure patterns](README.md)
