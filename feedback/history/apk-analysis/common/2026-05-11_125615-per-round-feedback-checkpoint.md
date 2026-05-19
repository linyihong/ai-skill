> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`enforcement/feedback-lessons.md`](../../../../enforcement/feedback-lessons.md)

### 2026-05-11 - Per-Round Feedback Checkpoint

Status: promoted

#### One-line Summary

每個有實質進展的 APK 分析回合結束前，agent 都要檢查是否產生新技巧或錯誤回饋，而不是等專案收尾。

#### Human Explanation

APK 分析常以多輪「繼續」推進：一輪新增 runner 參數、下一輪補 hook guard、再下一輪修文件或回歸測試。若 feedback 檢查只放在最終任務結束，容易在中途漏掉可重用技巧、validation rule 或 agent close-loop gap。把檢查放在每輪結尾，可以在新技巧仍清楚時立刻泛化，並避免 project-only commit 先完成後才回頭追補。

#### Trigger

- 使用者說「繼續」展開下一輪 APK 分析。
- 本輪新增 replay runner option、hook flag、short-window override、response-shape classifier、context guard、validation matrix row、evidence-attribution rule 或錯誤分類。
- 準備提交 project-only evidence、切回長時間動態測試、或回覆本輪進度前。
- 使用者詢問「有沒有加到 skill / 錯誤回饋 / 對話提示」。

#### Evidence

- Tool: APK 分析對話中的回合式工作流程。
- Sanitized excerpt: 一輪動態 capture 可能同時產生 project evidence 與可重用 runner/context guard；若只提交 project docs，下一輪 agent 不一定知道這個 guard 應沉澱。
- Evidence path: 具體 project incident 留在 `<PROJECT_ROOT>` docs；本 lesson 只保留泛化規則。

#### Generalized Lesson

將 feedback 檢查設為 per-round checkpoint，而不是 final-only checklist。每個有實質進展的 work round 結束前，agent 必須自問：

```text
本輪是否新增可重用技巧、validation rule、replay/runner guard、hook option、error pattern 或 close-loop gap？
```

若答案是 yes，先開啟 canonical `<AI_SKILL_REPO>` writeback transaction，再繼續長時間專案工作。若答案是 no，回覆時可簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。

#### Agent Action

1. 在每輪結束、使用者說「繼續」前後、project-only commit 前做 feedback checkpoint。
2. 將結果分成 `new reusable lesson`、`failure-learning needed`、`project-specific only`、`unvalidated candidate`。
3. 對 `new reusable lesson` 新增 `feedback_history/<category>/` lesson；成熟時同步 promotion target。
4. 對 `failure-learning needed` 讀 `failure-learning-system.md` 並判斷是否需要 `enforcement/failure-patterns/`。
5. 不把具體 App 名稱、sample id、host、payload 或本機路徑寫入 reusable lesson。

#### Goal / Action / Validation

- Goal: 降低長對話 APK 分析中漏寫 skill feedback 或錯誤回饋的機率。
- Action: 在 `SKILL.md` 與 `enforcement/feedback-lessons.md` 加入 per-round feedback checkpoint，並以本 lesson 記錄泛化規則。
- Validation or reference source: 回合結束時能反查 `SKILL.md` Quick Start §8、Feedback Loop checklist，以及 shared feedback rule 的每輪檢查要求。

#### Applies When

- APK 分析或動態 capture 以多輪方式持續推進。
- 本輪有新增工具選項、hook 行為、驗證 gate、文件化規則或 agent close-loop 修正。
- 使用者要求確認 skill feedback、錯誤回饋或規則提示是否已加入。

#### Does Not Apply When

- 本輪只是純狀態回報，沒有新證據、新方法、新錯誤或新驗證規則。
- 內容只對單一 project 的私有資料成立，無法抽象成通用方法。

#### Validation

- `SKILL.md` 的 Automatic skill feedback 明確要求 per-round checkpoint。
- `enforcement/feedback-lessons.md` 明確要求每輪回饋檢查。
- 若本輪沒有新增 lesson，最終回覆能說明為何只是 project-specific 或尚未驗證。

#### Revision - 2026-05-19

如果使用者問「現在有什麼可以回饋的技巧嗎」，而 agent 已能列出具體可重用技巧、validation rule 或 hook guard，不能只口頭回答清單。這代表 checkpoint 已判定 `new reusable lesson`，agent 應立即寫入 `feedback/history/<domain>/<category>/` 並同步索引；若暫不寫，必須明確說出 blocker（例如證據不足、需要使用者授權、或 canonical repo 不可寫）。

#### Promotion Target

- `SKILL.md`
- `enforcement/feedback-lessons.md`

#### Required Linked Updates

- 已依 [`linked-updates.md`](../../../../enforcement/linked-updates.md) 檢查：本 lesson 已 promoted 到 `SKILL.md` 與 shared feedback rule；`feedback_history/README.md` 需要新增索引列。
- 已依 [`reusable-guidance-boundary.md`](../../../../enforcement/reusable-guidance-boundary.md) 檢查：正文只保留 generalized checkpoint，具體 project evidence 留在 project docs。
