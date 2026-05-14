# Entrypoint Positioning Drift（入口定位漂移）

Status: validated
Class: validation-gap

## Trigger

當 repository、skill、shared rule、architecture document 或 tool adapter 被 rename、re-scope，或提升成新的 top-level concept 時，使用此 pattern。

## Failure Mode

Agent 更新了 secondary references、mid-document sections、links 或 filenames，卻漏掉 primary entrypoint positioning：root title、opening paragraph、index summary，或 future reader 第一眼會看到的 framing。

## Risk

即使 deeper links 已指向新概念，使用者與 future agents 仍會先看到 stale framing。這會混淆 source of truth、削弱 architecture adoption，並讓已完成的更新看起來不完整。

## Required Agent Action

針對 naming、branding、architecture 或 top-level scope changes：

1. 編輯前先辨識 primary entrypoint files。
2. 更新 title 與 opening positioning，不只更新 links 或 mid-document sections。
3. 編輯後搜尋 old names、old slugs 與 old framing。
4. 從第一行開始重讀 entrypoint，模擬使用者看到的順序。
5. 若使用者指出 drift，立刻跑 failure learning loop，而不是只補漏掉的那一行。

## Prevention Gate

宣稱完成前，回答：

- 第一個 heading 是否使用新的 canonical name？
- 第一段是否描述新的 role/scope？
- Root indexes 與 architecture/tool indexes 是否指向新的 canonical file？
- 搜尋 old names 與 old slugs 時，是否只剩 intentional historical references？

## 驗證

從第 1 行讀回 root entrypoint 與已變更的 architecture/tool/shared indexes。對 old name 與 old slug 做 exact search。Commit/push/readback 後確認 `git status --short --branch` 乾淨。

## Linked Rules

- [`../failure-learning-system.md`](../failure-learning-system.md)
- [`../goal-action-validation.md`](../goal-action-validation.md)
- [`../linked-updates.md`](../linked-updates.md)
- [`../content-layering.md`](../content-layering.md)

## Linked Validation Scenarios

- `validate_intelligence_classification_boundary` — 檢查 `intelligence/README.md` 的結構圖與實際目錄一致，防止入口定位漂移導致新目錄未在結構圖中註冊
