> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-08 - Agent goals before transcripts when resuming

Status: candidate

#### One-line Summary

接續中斷工作時，先讀 project-local `.agent-goals`，再用 transcripts、terminal output 與 git 狀態交叉確認。

#### Human Explanation

Transcript 和 git status 只能回答「最近說過什麼、磁碟目前變成什麼」，不一定能回答「使用者認定哪個 goal 還沒完成、下一步決策是什麼、是否有 lock 或 single-owner 限制」。若 agent 中斷後先看 transcript，容易漏掉 `.agent-goals` 裡已整理好的 active goal、open decisions、next action 與 validation gap。

#### Trigger

使用者說 agent 中斷、突然關閉、要從哪裡重做、剩下什麼、下一步是什麼，或質疑 agent 沒有讀 goal ledger。

#### Evidence

- Tool: Cursor agent conversation recovery.
- Sanitized excerpt: A resume answer was first based on transcripts and git state; the user pointed out that the existing `.agent-goals` ledger should have been consulted first.
- Evidence path: project-local `.agent-goals/README.md` and referenced goal files; do not copy project-specific goal content into reusable skill files.

#### Generalized Lesson

For resumable analysis or SDK work, `.agent-goals` is the primary handoff source. Transcripts, terminal output, and git status are confirmation sources. This ordering prevents the agent from missing owner/lock state, active priority, open decisions, next action, and completion criteria.

#### Agent Action

When resuming interrupted work:

1. Read `<PROJECT_ROOT>/.agent-goals/README.md`.
2. Read the referenced active goal under `<PROJECT_ROOT>/.agent-goals/goals/`.
3. Check owner, lock, priority, open work, decisions, next action, completion criteria, and validation.
4. Then read transcripts / terminal output and run git status to confirm recent actions and disk state.
5. If an overlapping lock or different active owner is present, stop and ask before editing.

#### Goal / Action / Validation

- Goal: Make interrupted-session recovery start from the durable project goal ledger instead of chat history alone.
- Action: Promote the ordering into Cursor tool guidance and keep this feedback lesson as the reusable reminder.
- Validation or reference source: `shared-rules/conversation-goal-ledger.md` defines `.agent-goals` as the project-local handoff source; `ai-tools/cursor.md` now states the Cursor-specific resume ordering.

#### Applies When

- Work spans multiple tool calls, sessions, agents, or model/context transitions.
- The user asks where to resume or what remains unfinished.
- The project has a `.agent-goals/` directory or the task meets the trigger for initializing one.

#### Does Not Apply When

- The user asks a single factual question with no pending project work.
- The project has no ledger and the task is a one-message answer with no files changed.

#### Validation

Confirm by checking that Cursor tool guidance says `.agent-goals` comes before transcripts/git for interruption recovery, and that future resume answers cite active goal state before recent transcript summary.

#### Promotion Target

- `ai-tools/cursor.md`
- `shared-rules/conversation-goal-ledger.md`

#### Required Linked Updates

- Updated `ai-tools/cursor.md` because this is a Cursor-specific execution failure.
- Checked `shared-rules/conversation-goal-ledger.md`; it already defines `.agent-goals` as the project-local source of truth and requires reading it before editing, so no generic rule change is needed.
- Updated `skills/apk-analysis/feedback_history/common/README.md` index.
