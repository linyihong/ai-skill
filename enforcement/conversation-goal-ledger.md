# 對話目標帳本

本規則定義工具中立的 temporary ledger，用來追蹤 active conversation goals。它讓後續 agent 能在中斷、context compaction、model switch 或 multi-agent handoff 後恢復未完成工作。

Goal ledger 是 project-local state，不是 reusable knowledge。它應放在目前 project 底下、排除於 git 之外，並在目標完成且驗證後刪除。

Goal ledger 讓人和 agent 在長對話後快速知道：哪些工作未完成、哪些決策待確認、下一步優先做什麼、哪些內容仍需補強才能算完成。

Goal ledger 不是長期 roadmap、知識生命週期紀錄或完成工作 archive。若一個 conversation goal 完成後仍代表長期方向、後續 phase、未完成能力、promotion / deprecation 狀態或系統治理決策，agent 必須在刪除 `.agent-goals/` 前確認該長期狀態已寫入 durable planning 文件，例如 `architecture/` roadmap、相關 layer README、`governance/`、`knowledge/`、`metadata/`、issue tracker 或專案正式 planning docs。

## 目的

當 conversation work 跨越一個以上 action、包含多個 goals、可能被中斷，或需要明確 completion criteria 時，使用 goal ledger。

Ledger 必須回答：

| 欄位 | 必填內容 |
| --- | --- |
| Goal | 正在追求的 user-visible outcome。 |
| Priority | `P0`、`P1`、`P2` 或 `P3`。 |
| Status | `active`、`paused`、`blocked`、`needs-validation`、`superseded` 或 `complete-pending-delete`。 |
| Parallelization mode | `parallelizable`、`single-owner` 或 `non-parallelizable`；若重疊工作不安全，附簡短原因。 |
| Owner | 目前 agent/tool owner 與 timestamp。 |
| Owner / lock decision | 目前 agent 是否可編輯、需 acquire/refresh lock，或因其他 owner/lock 重疊而必須停止。 |
| Source | 建立 goal 的使用者要求或指令。 |
| Scope | In scope、out of scope、affected project/repo。 |
| Subgoals | 目標拆解後的 child goals 或 checklist items。 |
| Planning / todo links | 相關 planning document path、plan section、TodoWrite IDs、checklist items 或 external issue links。 |
| Open work / decisions | 尚未完成內容、需要的 decision、或需 strengthening 的項目。 |
| Dependencies | 所需 user answer、external command、file、agent 或 upstream goal。 |
| Next action | 新 agent 應採取的下一個具體步驟。 |
| Completion criteria | 刪除 goal file 前必須成立的條件。 |
| Validation | 完成狀態如何驗證，或之後要如何驗證。 |

## Goal 層級邊界

| 層級 | 放置位置 | 用途 | 完成後處理 |
| --- | --- | --- | --- |
| Active conversation goal | `<PROJECT_ROOT>/.agent-goals/` | 本輪對話中尚未完成、可中斷、需接手的 user-visible 工作。 | 完成條件與驗證成立後刪除。 |
| Durable roadmap goal | `architecture/`、layer README、`governance/`、`knowledge/`、`metadata/`、issue / planning docs | 長期方向、phase、未完成能力、migration 狀態、promotion / deprecation 決策。 | 保留在正式文件；不要放在 `.agent-goals/` 當 archive。 |
| Implementation task | `.agent-goals/` + todo / plan link | 從 durable roadmap 拉進本輪執行的具體工作。 | 結果回寫 durable 文件後，再刪除 active goal。 |

判斷規則：

- 若目標完成後只剩「已做完」的歷史紀錄，使用 git commit、PR、issue 或正式文件 changelog；不要保留 completed goal row。
- 若目標完成後仍有下一階段、開放問題、待決策、未完成能力或長期治理狀態，必須寫入 durable planning 文件，再刪除 `.agent-goals/`。
- `.agent-goals/` 可以連到 durable roadmap goal，但不能取代 durable roadmap goal。
- Durable roadmap goal 被拉進本輪工作時，建立或更新 active conversation goal，並在 `Planning / Todo Links` 連回正式文件。

## 位置

在正在工作的 project 內保存 ledger：

```text
<PROJECT_ROOT>/.agent-goals/
  README.md              # main goal table / quick locator
  goals/
    P1-<slug>.md
  locks/
    <goal-id>.lock/
```

不要把 canonical ledger 放在工具專屬設定目錄；此 workflow 必須可跨不同 agent tools 使用。工具可以讀取或提醒此目錄，但它不是 source of truth。

`.agent-goals/` 是 temporary project state，不應 commit。優先用 `.git/info/exclude` 排除，避免業務 repo 產生 policy churn。只有團隊希望追蹤此慣例時，才將 `.agent-goals/` 加入 `.gitignore`。

## 何時建立或更新

實質工作前，先執行或等效完成：

```text
<AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> status
```

接著讀 `<PROJECT_ROOT>/.agent-goals/README.md` 與相關 active goal file，確認：

- Active / blocked / needs-validation goals。
- Priority，以及另一個 active `P1` 是否與最新使用者要求衝突。
- Owner 與 lock state。
- Parallelization mode，以及是否允許重疊工作。
- 既有 planning / todo links。
- Open missing work、decisions 與 needs-strengthening items。

Ledger 是 recovery aid，不是自動切換任務的來源。最新使用者要求、目前 accepted plan、以及進行中的 tool todo list 才定義本輪 active task。若使用者說「continue」但有多個合理 goal，或 ledger active row 指向與最新 task 不同的 project，先詢問要繼續哪個 goal，不要讀取或操作無關 project。

若 ledger 不存在，且符合下列任一 trigger，先初始化：

```text
<AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> init
```

建立或更新 goal file 的情境：

- 使用者要求 implementation、analysis、planning、review、debugging 或 repository updates，且可能跨多個 tool call。
- 任務有多個 goals 或 priorities。
- Agent 看到 active project 有 modified、staged、untracked 或其他 dirty files，且打算繼續工作。
- Agent 建立 tool-level todo list，或恢復未完成 todo list。
- 使用者要求繼續 prior multi-step task，尤其在 context compaction、中斷或 side quest 後。
- 任務 paused、blocked、superseded 或等待使用者輸入。
- Goal 被拆成小目標。
- Agent 即將停止、compact context、switch mode、launch subagents 或 hand off。
- 使用者改變 priority、新增 target 或重導對話。

非常小的一次性回答可不建立 ledger。只要回覆後仍有工作、已有檔案變更，或 active task 的 working tree 是 dirty，ledger 就不再是 optional。不要把 tool todo list 當成 ledger 替代品；todo 追蹤執行步驟，`.agent-goals/` 追蹤 user-visible goals 與 handoff state。

修改 ledger-tracked task 的檔案前，目前 goal 必須有足夠 handoff 結構：

- frontmatter 中有 `parallelization`，或正文有明確 parallelization line。
- Owner / lock decision 已記錄，或由 unlocked current-owner goal 可推導。
- 若有 plan、checklist、TodoWrite item 或 issue，已連到 Planning / Todo Links。
- `Missing work`、`Decision needed`、`Needs strengthening` 有實際內容或 `none`。
- Next action 與 completion criteria 具體到 future agent 可驗證。

## Goal file 範本

使用 Markdown，讓任何工具都能讀：

```markdown
---
id: P1-short-slug
priority: P1
status: active
parallelization: single-owner
owner: <agent/tool/session>
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
project: <PROJECT_ROOT or project label>
---

# <Goal title>

## Source Request
<User request or concise quote.>

## Scope
- In:
- Out:
- Affected paths/repos:

## Subgoals
- [ ] <subgoal>

## Planning / Todo Links
| Type | Reference | Status / Note |
| --- | --- | --- |
| plan | <path#section or none> | <why it matters> |
| todo | <todo id / checklist item / issue> | <pending / in_progress / completed / blocked> |

## Owner / Lock Decision
- Owner: <current agent/tool/session>
- Lock: <unlocked / lock acquired / blocked by owner>
- Parallelization: <parallelizable / single-owner / non-parallelizable> because <reason>

## Open Work / Decisions
- Missing work: <none / concrete unfinished work>
- Decision needed: <none / decision required>
- Needs strengthening: <none / weak rule / validation / docs>

## Dependencies
- <none / user answer / external state / parent goal>

## Progress
- <timestamp>: <what changed>

## Next Action
<The next concrete action for a future agent.>

## Completion Criteria
- [ ] <observable completion condition>

## Validation
- <diff review / test / lint / source checked / user confirmation / not yet validated>

## Handoff Notes
<Risks, blockers, assumptions, and recovery hints.>
```

不要把 secrets、tokens、raw private data、reservation codes、personal addresses 或 private host details 寫入 ledger。使用 redacted labels 或 project-local references。

## 主要 goal 表格

`<PROJECT_ROOT>/.agent-goals/README.md` 是 active goals 的主要 locator，應包含連到各 goal file 的 compact table：

```markdown
| Priority | Status | Mode | Owner | Lock | Goal | Open Work / Decisions | Planning / Todo Links | Next Action | Updated |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| P1 | active | single-owner | agent/session | unlocked | [Short title](goals/P1-short-slug.md) | decision: choose live gate | plan: docs/plan.md#section; todo: implement-api | Run validation | 2026-05-08T00:00:00Z |
```

Main table 用於快速恢復，不取代各 goal file 的 detail。

當 goal 被建立、暫停、拆分、連到 todo、owner/lock state 改變、parallelization mode 改變或完成時，更新 main table。Goal file 在驗證後刪除時，也從表格移除。

## Planning 與 Todo 連結

若存在 planning document、checklist 或 tool-level todo list，將它們連到 goal ledger：

1. 可行時，在相關 plan section、checklist item 或 todo 旁標出 goal ID。
2. 在 goal file 的 `Planning / Todo Links` 記錄 plan path、section anchor、TodoWrite ID、checklist item 或 issue ID。
3. 若 todo 變成獨立可恢復工作，將它加入 subgoal，或拆成 child goal。
4. Todo 完成時，先更新 goal progress 與 validation notes，再刪除 goal。
5. Todo 因使用者改變方向而取消時，將 linked goal 標成 `paused` 或 `superseded`，並記錄原因。

Goal ledger 追蹤 user-facing intent；todo tools 追蹤 execution steps。兩者要互相連結，讓 future agent 能從高層 goal 跳到具體 plan/todo item，再回到 goal。

Document-level TODO lists 屬於單一文件，應出現在該文件前段。見 [`document-todo-list.md`](document-todo-list.md)。當 document TODO 屬於更大 user goal 時，從 `Planning / Todo Links` 或 `Open Work / Decisions` 連回。

## Priority 規則

使用下列 priorities：

| Priority | 意義 |
| --- | --- |
| `P0` | User-blocking、safety/secret risk、data-loss risk，或明確 urgent request。 |
| `P1` | 目前主要 user goal。 |
| `P2` | Primary path 後重要 follow-up 或 validation。 |
| `P3` | Nice-to-have cleanup、optional refactor 或 low-risk follow-up。 |

通常一個 conversation 只應有一個 active `P1`。若新的 `P1` 到來，將前一個 `P1` pause 或 supersede，並記錄原因與 next action。

Ledger 是 recovery context，不可覆蓋最新 user message。遵循 `active` row 前，先比對目前 user request、目前 viewed/referenced files 與最近 conversation thread。若最新要求清楚延續另一個 project、repo 或 goal，跟隨最新要求並更新 ledger，而不是切到 stale active row。

如果照 ledger 會移動到無關 repo/project，或離開使用者剛詢問的 project，先停下來確認。不要把 promoted `active` goal、stale todo list 或舊 handoff summary 視為 silent project switch 的許可。

若 agent 因 stale 或 over-promoted ledger state 切錯 project，必須停止錯誤工作、回報它建立的 dirty files、更新 ledger/rule 記錄 root cause，並回到 user-requested project。

## Decomposition

Goal 太大時拆分：

1. 保留 parent goal，代表 user-facing outcome。
2. 為可獨立恢復的工作加入 child goals 或 checklist items。
3. 記錄 child goals 之間的 dependencies。
4. 只有當 child goal 是目前工作焦點時，才提升為 `P1`。

不要只在 chat 裡隱藏 discovered subgoal。若它影響完成條件，就記錄在 ledger。

## Owner、Lock 與 Parallelization

修改 active goal 相關檔案前，決定並記錄工作如何重疊：

| Mode | 意義 | 必要行為 |
| --- | --- | --- |
| `parallelizable` | 獨立檔案或清楚分離的 subgoals 可並行。 | 拆 child goals 或 todos，記錄 ownership，避免沒有 lock 就編輯同一檔。 |
| `single-owner` | 可跨 session 延續，但同一時間只應有一個 active owner 編輯。 | 編輯前檢查 lock state；handoff 時更新 owner 與 next action。 |
| `non-parallelizable` | Shared state、secrets、live captures、migrations 或 fragile workflows 讓重疊工作不安全。 | Acquire/refresh lock；若有其他 owner/lock，停止並詢問。 |

若 `.agent-goals/README.md` 顯示重疊工作有其他 active lock，編輯前停止。回報 lock owner、age、affected goal 與預計 next step。不要只因為目前 chat 有上下文就假設 lock stale；使用設定好的 cleanup rule 或詢問使用者。

當新要求與不同 id 的 active goal 重疊、`single-owner` goal 的 recorded owner 不同、workflow 是 `non-parallelizable`，或 agent 無法判斷兩個 goals 是否重疊但涉及同一 files、git branch、database、release 或 shared state 時，也要停止並詢問使用者。詢問應包含 existing goal id/title、owner 與 lock age（若有）、affected files/resources、平行工作風險，以及等待、接手、拆 child goal 或建立非重疊 goal 等具體選項。

下列 workflow 應標為 `non-parallelizable`：git history operations、merge conflict resolution、release tagging、deploys、migration sequencing、shared rule / skill writeback transactions、data migrations、destructive operations、credential rotation、production configuration，或任何兩個 agents 獨立編輯可能讓 validation 失效、重複 commit 或產生矛盾 user-facing decisions 的任務。

當使用者重導工作、改 priority，或在 side task 後要求繼續時，實質編輯前同步更新 owner/lock/parallelization、`Missing work`、`Decision needed`、`Needs strengthening`、`Planning / Todo Links` 與 `Next Action`。

## Goal 轉移

使用者重導任務時：

1. 將舊 goal 更新為 `paused` 或 `superseded`。
2. 記錄暫停原因與恢復條件。
3. 建立或提升新 goal，給正確 priority。
4. Final response 明確說明目前 active goal。

若新 goal 與高風險未完成 goal 衝突，切換前先標出衝突。

## Multi-agent Safety

Agents 必須透過 lock directories 協調：

```text
<PROJECT_ROOT>/.agent-goals/locks/<goal-id>.lock/
  owner
  pid
  startedAt
```

使用 atomic directory creation 建立 locks。若另一個 active lock 存在，不要修改該 goal。回報 owner、age 與預計 next step。只有在檢查 recorded PID/session 已不再 active，或取得使用者同意後，才移除 stale lock。

建議 default TTL：30 minutes。工具可在合法 long-running task active 時覆寫 TTL。

## 完成與刪除

只有全部成立時才刪除 goal file：

1. Completion criteria 已滿足。
2. Validation 已執行，或使用者明確接受結果。
3. 沒有 child goal 仍為 active、blocked 或 needs validation。
4. 若 goal 完成後仍有 durable roadmap / lifecycle / governance / follow-up 狀態，已回寫到正式 planning 文件或 issue。
5. Final answer 或 handoff 已說明 outcome。

當所有刪除條件成立時，必須在同一 close-out turn 刪除 goal file，並刷新 `<PROJECT_ROOT>/.agent-goals/README.md`，避免 completed work 長期留在 active recovery table。不要把 `completed` row 當成長期記憶或 archive；durable lessons 應放在 project docs、commits、issues 或 reusable skill feedback，不放 `.agent-goals/`。Durable roadmap、phase 與 migration 狀態應放在正式 planning 文件，不以 `.agent-goals/` 保存。

如果工作完成但 validation 缺失，設為 `needs-validation`，不要刪除。

若 goal 被 superseded，保留到使用者接受新方向，或原因足夠清楚讓 future agent 接手。之後可依 project preference 刪除或歸檔。

## 與 Canonical Repository Writeback Transactions 的關係

本 ledger 與 [`dependency-reading.md`](dependency-reading.md) 中的 canonical repository writeback transaction 是不同機制。

- Conversation goal ledger：project-local、temporary、不提交，追蹤 user goals 與 handoff state。
- Canonical repository writeback transaction：repository-specific、需 commit/push，追蹤本 knowledge OS 的變更與 configured tool sync/mirrors。

更新本 repository 時，兩者可能同時適用：project goal ledger 追蹤 user-facing task，而 repository transaction 仍必須完成 diff review、linked updates、sync、commit、push、reread 與 clean status。

## Tool Integration

工具可自動化 ledger checks，但除非能 deterministic validate goal，否則 automation 只能作為提醒。Hook 或 script 可以提醒、建立或檢查 goals；沒有 completion criteria 與 validation evidence 時，不應 silently mark goals complete。

工具專屬處理見 `ai-tools/` 中對應文件。

← [Back to shared rules index](README.md)
