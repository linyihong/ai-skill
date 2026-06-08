---
id: 2026-06-08-1047-feedback-learning-report-obligation
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-08
---

# Feedback / Learning Report Obligation

**Status**: `draft`
**Maturity**: Phase 0 checked; ready for Phase 1
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-08
**Priority**：P1

## Decision Rationale

### Problem & Why Now

Ai-skill 已有 final close-out 的 Cognitive Mode reporting obligation，用來讓 agent 每輪結束時顯性回報自己的 cognition state。但目前缺一個同級 close-out 檢查，要求 agent 顯性判斷「本輪是否產生可回饋 / 沉澱為 reusable knowledge 的內容」。

這補齊 runtime close-out 的第三個維度：

| 維度 | 問題 |
| --- | --- |
| Execution | 有沒有完成 |
| Cognition | 怎麼思考 |
| Learning | 有沒有值得留下 |

這在 repo context 不完整、無法確認 clean/pushed/readback、或無法直接寫入 project docs / feedback history 時特別重要。這些場景常見風險是：

- Agent 發現 workflow / validation / evidence gap，但 final response 只回報完成，沒有說是否需要沉澱。
- 使用者指出 agent 漏讀、漏驗證或誤判，但 agent 只修當下回覆，沒有觸發 failure learning loop。
- Agent 不是在本地 repo 啟動，無法完整執行 repo close-loop，卻沒有留下 feedback / learning handoff。
- 需要沉澱的知識被留在 chat 中，沒有進入 `feedback/history/`、`intelligence/`、`workflow/`、`enforcement/` 或 project docs。

### Decision

新增一個 per-turn close-out obligation：`obligation.feedback.learning_report`。

它與 `obligation.cognitive.mode_report`、`obligation.finality.close_loop_check` 同級，但目的不同：

- Execution close-loop：回報事情是否完成、repo 是否關閉。
- Cognitive Mode Report：回報 agent 本輪使用的 cognitive state。
- Feedback / Learning Report：回報本輪是否有值得留下的 learning disposition，以及若需要，目標層與寫回狀態是什麼。

核心原則：

```text
Every final response must report the learning decision.
Not every turn must write a feedback lesson.
```

也就是「必填 report，不強制每次 write」。

### Alternatives Considered

- **A. 維持現狀，靠 agent 自律判斷是否回饋** — reject。Cognitive Mode 之所以成為 obligation，就是因為 final close-out 類判斷靠自律容易漏。
- **B. 每次都強制新增 feedback lesson** — reject。會造成 `feedback/history/` 膨脹，把純查詢、已覆蓋規則與低價值紀錄都寫成 durable knowledge。
- **C. 只在 non-local repo 時要求 report** — reject。Repo context 與 learning decision 是正交資訊；local repo 中使用者指出 failure / workflow gap 時也應回報 learning decision。
- **D. 用 `NONE` / `CHECK` / `NEEDED` 單一欄位表達所有狀態** — reject。它混合「是否需要學習」「是否能判斷」「是否能寫回」「repo 是否可見」四個不同維度，會製造假訊號。
- **E. 拆成 `feedback_decision`、`repo_context`、`writeback_status` 三維 contract** — accept。讓 hook 可做 schema validation，semantic correctness 留給 scenarios；knowledge classification 留給 knowledge acquisition / economics / memory plans。

### Why Not an ADR Yet

此變更會影響 runtime/core-bootstrap obligation、tool stop hook、final response shape、validation scenarios 與 stop-hook validators。現在仍在設計 report 格式、repo context 判斷、hook fail-open / fail-closed 邊界，不宜先寫 accepted ADR。

### ADR Promotion Criteria

- [ ] `runtime/core-bootstrap.yaml` 新 obligation 已落地並被 tool stop hook 消費。
- [ ] Cursor / Claude / generic adapter 的 close-out behavior 已同步。
- [ ] 至少 3 個 validation scenarios 覆蓋：缺 report、enum/schema invalid、需要 feedback 但未回報。
- [ ] 實際使用一段時間後，沒有造成 feedback lesson spam。
- [ ] 確認這是 cross-session / cross-agent foundational obligation，而不是單一工具行為。

### Consequences

#### 正面

- Agent final response 會顯性回答 learning disposition。
- Repo context、learning decision、writeback capability 不再混成單一狀態。
- 使用者指出失效時，agent 比較不會只修當下而漏掉 reusable prevention。

#### 負面

- 每個 final response 會多一行或一個小區塊。
- Stop hook / commit hook / tool adapter 需要同步更新，否則會出現各工具 enforcement 不一致。

#### 風險

- 若格式太複雜，agent 可能 ritualize 填報。
- 若本 plan 開始分類 observation / lesson quality，會侵入 Knowledge Acquisition Layer。
- 若 hook 試圖判斷語義正確性，會把機械 close-out gate 變成 fragile semantic reviewer。

Glossary Impact: yes — proposed terms `feedback_learning_report`, `learning_decision`, `repo_context`, `writeback_status`; glossary entry 是否需要新增待 Phase 0 決定。

## Runtime Execution Path

| Surface | Planned owner | Runtime / consumer |
| --- | --- | --- |
| `runtime/core-bootstrap.yaml` `per_turn_obligations` | runtime | stop hook / final close-out validation |
| `CORE_BOOTSTRAP.md` companion description | root companion | human-readable bootstrap docs |
| `scripts/ai-skill-cli/internal/app/hooks.go` stop hook | CLI hook runner | Cursor / Claude / future tools |
| `ai-tools/agent/*.md` close-out notes | tool adapters | tool-specific behavior docs |
| `validation/scenarios/runtime/*feedback*` | validation | runtime validate / scenario audit |

Trigger flow:

```text
final response or session stop
→ stop hook reads assistant final response
→ validate presence of Cognitive Mode block
→ validate presence of Feedback / Learning Report
→ if missing, emit repair follow-up instruction
→ if present, allow close-out
```

Repo context / writeback flow:

```text
repo context unavailable OR not local repo OR clean/push/readback cannot be verified
→ repo_context = NON_LOCAL or UNKNOWN
→ feedback_decision still independently resolves to NONE / NEEDED / UNKNOWN
→ if feedback_decision = NEEDED and writeback cannot happen, writeback_status = UNAVAILABLE or DEFERRED
```

## Proposed Report Contract

### Compact Form

```text
FeedbackDecision: NONE
RepoContext: LOCAL
Writeback: N/A

FeedbackDecision: NEEDED
RepoContext: LOCAL
Writeback: COMPLETED
Target: workflow

FeedbackDecision: NEEDED
RepoContext: NON_LOCAL
Writeback: UNAVAILABLE
Target: feedback-history

FeedbackDecision: UNKNOWN
RepoContext: UNKNOWN
Writeback: N/A
```

Compact form uses fixed key order for simple schema validation:

1. `FeedbackDecision`
2. `RepoContext`
3. `Writeback`
4. `Target` only when `FeedbackDecision: NEEDED`

### Full Form

Use full form when `FeedbackDecision: NEEDED`, `FeedbackDecision: UNKNOWN`, high-risk failure learning, or deferred/unavailable writeback:

```markdown
### Feedback / Learning Report

| Field | Value |
| --- | --- |
| feedback_decision | NONE / NEEDED / UNKNOWN |
| repo_context | LOCAL / NON_LOCAL / UNKNOWN |
| writeback_status | COMPLETED / DEFERRED / UNAVAILABLE / N/A |
| trigger | user correction / validation gap / workflow gap / repeated failure / none |
| target | none / feedback-history / intelligence / workflow / enforcement / project-docs |
| reason | <why> |
```

## Decision Rules

### Contract Dimensions

| Dimension | Values | Meaning |
| --- | --- | --- |
| `feedback_decision` | `NONE` / `NEEDED` / `UNKNOWN` | 是否有值得沉澱的 learning，以及是否有足夠 evidence 判斷。 |
| `repo_context` | `LOCAL` / `NON_LOCAL` / `UNKNOWN` | Agent 是否能確認目前 repo context。 |
| `writeback_status` | `COMPLETED` / `DEFERRED` / `UNAVAILABLE` / `N/A` | 若需要沉澱，是否已能寫回或只能延後。 |

### Decision Semantics

| Value | Meaning | Notes |
| --- | --- | --- |
| `feedback_decision: NONE` | 沒有 reusable learning 需要沉澱 | 可與 `repo_context: NON_LOCAL` 同時存在，例如純概念問答。 |
| `feedback_decision: NEEDED` | 有值得沉澱的 lesson / project feedback / rule gap | 必須提供 `target`，並說明 `writeback_status`。 |
| `feedback_decision: UNKNOWN` | evidence insufficient / cannot verify / repo inaccessible | 不等於需要沉澱；只是不能可靠判斷。 |

### Triggers for `NEEDED`

- 使用者指出 agent 漏讀、漏驗證、誤判、重複犯錯。
- 發現 workflow、validation、evidence chain、routing、close-loop 或 feedback gap。
- 本輪產生新的 reusable reasoning pattern。
- 現有規則不足以防止同類錯誤。
- Non-local context 下發現需要 project docs / feedback history / handoff，但無法直接寫回；此時 `writeback_status: UNAVAILABLE`。

### Triggers for `UNKNOWN`

`UNKNOWN` should be rare. It is an evidence state, not a safety fallback for cases where the agent is unsure what to report.

- 無法確認必要 evidence。
- Repo inaccessible 或 repo context unknown，且 task outcome 可能包含 reusable learning。
- 只取得 chat evidence，沒有足夠 source / project docs 判斷 learning disposition。

### Repo Context Rule

`repo_context` 不決定 `feedback_decision`。例如：

```text
FeedbackDecision: NONE
RepoContext: NON_LOCAL
Writeback: N/A
```

這可用於純概念問答、沒有新 pattern、沒有 failure、沒有 workflow gap 的情境。

### Writeback Status Rule

`writeback_status` 只描述需要沉澱時的寫回能力，不替代 `feedback_decision`。

```text
FeedbackDecision: NEEDED
RepoContext: NON_LOCAL
Writeback: UNAVAILABLE
Target: project-docs
```

此情境代表「確實需要沉澱，但本輪無法寫回」，不是 `UNKNOWN`。

## Enforcement Boundary

Stop hook should validate only:

- Presence.
- Schema shape.
- Enum values.
- Required field combinations, e.g. `NEEDED` requires `Target`.

Stop hook should not validate semantic correctness:

- It must not decide whether the agent should have used `NONE` vs `NEEDED`.
- It must not infer from user correction that `NONE` is wrong.
- Semantic failures belong in validation scenarios and reviews.

Commit messages are explicitly out of scope. This report is a chat/session runtime close-out artifact, not a code artifact.

## Out of Scope Boundary

This plan reports learning disposition and writeback capability only. It does not track:

- knowledge classification quality
- promotion lifecycle
- knowledge acquisition lifecycle
- memory integration
- economics / knowledge cost
- activation fitness
- telemetry
- linked update completion

Those belong to Runtime Cognitive State / Knowledge Acquisition / Economics / Memory / Fitness contracts, especially [`2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md) and [`2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md`](2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md).

## Open Questions

- [x] `repo_context: LOCAL` 的判斷是否要求 `git status` clean，還是只要求有 local repo root？Resolved: local repo / project context is enough. Clean/pushed/readback belongs to execution close-loop or `writeback_status`, not `repo_context`.
- [x] `feedback_decision: NONE` 是否允許在 non-local repo 出現？Resolved: yes. Repo context and feedback decision are orthogonal.
- [x] Stop hook 應 fail-closed 還是 fail-open？Resolved: follow existing host-specific stop-hook behavior. Final assistant response present + missing required report loops/blocks; user-aborted, audit-only, non-final tool/status payloads fail-open; Cursor missing assistant text currently fail-opens, Claude missing assistant text blocks.
- [x] Compact block 是否需要固定行順序以便 schema validator 簡單檢查？Resolved: yes. Fixed order is `FeedbackDecision` → `RepoContext` → `Writeback` → optional `Target`.
- [x] Commit-msg 是否也要要求 Feedback report，或只要求 chat/session final response？Resolved: chat/session final response only. Commit message 不強制。

## Phase 0 — Preflight

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [x] 已讀本 plan §Open Questions 全部條目
- [x] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [x] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [x] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Repo local 判斷 | resolved | `hooks.go` separates local repo detection (`collectDirtyGitRepoReports`, `isAiSkillRepoRoot`) from dirty/ahead status; clean/pushed belongs to close-loop/writeback, not repo_context |
| Non-local 是否允許 NONE | resolved | Repo context 與 feedback decision 正交；non-local 可為 NONE |
| Stop hook fail policy | resolved | `hooks.go` and adapters show host-specific transport: Cursor uses followup loop, Claude uses `decision:block`; non-final / audit-only events fail-open |
| Compact line order | resolved | Fixed order supports presence/schema/enum validation without semantic judgment |
| Commit-msg 是否要求 | resolved | Chat/session final response only；commit message 不強制 |

### Phase 0.1 — Architecture Compatibility Preflight

- [x] 讀 `runtime/core-bootstrap.yaml` per-turn obligations。
- [x] 讀 stop hook implementation and tests：`scripts/ai-skill-cli/internal/app/hooks.go`、`hooks_test.go`。
- [x] 讀 tool adapters：`ai-tools/agent/cursor.md`、`ai-tools/agent/claude.md`、generic adapters as needed。
- [x] 讀 feedback learning rules：`enforcement/failure-learning-system.md`、`feedback/feedback-lessons.md`、`enforcement/reusable-guidance-boundary.md`。
- [x] 確認 report belongs to runtime obligation, not workflow-specific checklist。
- [x] 確認 linked updates：CORE_BOOTSTRAP companion、tool adapters、validation scenarios、CLI docs/tests。

Phase 0 result: proceed to Phase 1. No architecture conflict found. Existing active plans overlap only at future knowledge acquisition/economics surfaces; this plan remains the narrow close-out reporting contract and should land before broader Runtime Cognitive State work.

## Phase 1 — Runtime Contract

- [ ] Add `obligation.feedback.learning_report` to `runtime/core-bootstrap.yaml`.
- [ ] Update `CORE_BOOTSTRAP.md` companion with compact/full report examples and rationale.
- [ ] Decide final compact grammar and enum values.
- [x] Decide whether this obligation is chat/session-only or also commit-msg. Decision: chat/session final response only; commit message not required.

## Phase 2 — Hook Enforcement

- [ ] Extend stop hook final response validation to require Feedback / Learning Report.
- [ ] Add repair follow-up wording parallel to missing Cognitive Mode block.
- [ ] Add tests for missing report, compact report, full report, enum validation, and required field combinations.
- [ ] Preserve existing Bootstrap Receipt and Cognitive checks.

## Phase 3 — Feedback Learning Routing

- [ ] Define how `feedback_decision: NEEDED` maps to durable targets: `feedback-history`, `intelligence`, `workflow`, `enforcement`, `project-docs`.
- [ ] Define how `writeback_status` reports completed / deferred / unavailable writeback capability without tracking promotion lifecycle.
- [ ] Update feedback / failure learning docs only if Phase 0 finds the existing rules insufficient.
- [ ] Ensure report does not force writing low-value lessons.

## Phase 4 — Validation Scenarios

- [ ] Add `validation/scenarios/runtime/feedback-report-required-v1.yaml`.
- [ ] Add `validation/scenarios/runtime/feedback-report-schema-v1.yaml`.
- [ ] Add `validation/scenarios/runtime/non-local-repo-feedback-none-allowed-v1.yaml`.
- [ ] Add `validation/scenarios/failure-derived/feedback-needed-but-not-reported-v1.yaml`.
- [ ] Run `ai-skill runtime refresh` and `ai-skill runtime validate`.

## Phase 5 — Tool Adapter Docs

- [ ] Update `ai-tools/agent/cursor.md` stop close-out description.
- [ ] Update `ai-tools/agent/claude.md` if Claude adapter has stop/final behavior.
- [ ] Update `ai-tools/README.md` only if routing / adapter summary changes.

## 完成條件

- [ ] New per-turn obligation exists in runtime contract.
- [ ] Final close-out hook blocks or repairs missing Feedback / Learning Report.
- [ ] Compact and full formats are documented.
- [ ] Repo context, feedback decision, and writeback status are separate dimensions.
- [ ] Non-local repo can report `feedback_decision: NONE` when no learning exists.
- [ ] Hook validates report presence/schema/enums only, not semantic correctness.
- [ ] Commit message does not require Feedback / Learning Report.
- [ ] Validation scenarios exist and pass.
- [ ] Tool adapters document the new final close-out requirement.
- [ ] No secret, local path, project incident evidence, or private repo detail is written into reusable docs.
- [ ] Plan Completion Closure executed when implementation finishes.

## Stakeholder 同意項目

- [ ] Report is mandatory; writing a feedback lesson is conditional.
- [ ] Repo context must be explicit in final response.
- [ ] `Feedback: NONE` remains allowed for ordinary local low-risk turns.
- [ ] `feedback_decision: NONE` remains allowed for non-local low-risk turns.
- [ ] `feedback_decision: NEEDED` must name a durable target and writeback status.

## Per-surface consumer 表

| Generated surface key | Named consumer(s) | Consumer 類型 |
| --- | --- | --- |
| `runtime.core_bootstrap.contract` updated obligation | stop hook final response validator | Go hook consumer |
| `validation/scenarios/runtime/*feedback*` | runtime audit / validation scenario inventory | validation scenario |

## 與其他 plans 的關係

- Related to [`archived/2026-05-25-2100-runtime-cognitive-contract-v2.md`](../archived/2026-05-25-2100-runtime-cognitive-contract-v2.md): mirrors the idea of final response reporting obligation, but for learning decision rather than cognition state.
- Related to [`2026-05-27-1557-tool-runtime-signal-economics-integration.md`](2026-05-27-1557-tool-runtime-signal-economics-integration.md): that plan discusses knowledge acquisition inside future Runtime Cognitive State / economics surfaces. This plan should run first as a narrow close-out obligation and avoid absorbing the broader economics / telemetry scope.
- Related to [`2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md`](2026-05-28-1636-gen4-fitness-optimization-memory-interface-reservation.md): fitness / optimization memory remains future Gen4 interface work; this plan only reports learning disposition, not outcome fitness.
- Related to [`archived/2026-05-31-1900-workflow-activation-engine.md`](../archived/2026-05-31-1900-workflow-activation-engine.md): final close-out hooks are one enforcement surface for runtime agent behavior.
- Related to feedback/failure learning rules: this plan adds final reporting, not a replacement for failure learning loop.
