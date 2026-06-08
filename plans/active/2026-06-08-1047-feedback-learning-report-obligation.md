---
id: 2026-06-08-1047-feedback-learning-report-obligation
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-08
---

# Feedback / Learning Report Obligation

**Status**: `draft`
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-08
**Priority**：P1

## Decision Rationale

### Problem & Why Now

Ai-skill 已有 final close-out 的 Cognitive Mode reporting obligation，用來讓 agent 每輪結束時顯性回報自己的 cognition state。但目前缺一個同級 close-out 檢查，要求 agent 顯性判斷「本輪是否產生可回饋 / 沉澱為 reusable knowledge 的內容」。

這在 non-local repo、repo context 不完整、無法確認 clean/pushed/readback、或無法直接寫入 project docs / feedback history 時特別重要。這些場景常見風險是：

- Agent 發現 workflow / validation / evidence gap，但 final response 只回報完成，沒有說是否需要沉澱。
- 使用者指出 agent 漏讀、漏驗證或誤判，但 agent 只修當下回覆，沒有觸發 failure learning loop。
- Agent 不是在本地 repo 啟動，無法完整執行 repo close-loop，卻沒有留下 feedback / learning handoff。
- 需要沉澱的知識被留在 chat 中，沒有進入 `feedback/history/`、`intelligence/`、`workflow/`、`enforcement/` 或 project docs。

### Decision

新增一個 per-turn close-out obligation：`obligation.feedback.learning_report`。

它與 `obligation.cognitive.mode_report` 同級，但目的不同：

- Cognitive Mode Report：回報 agent 本輪使用的 cognitive state。
- Feedback / Learning Report：回報本輪是否需要沉澱 reusable learning，以及若需要，目標層與狀態是什麼。

核心原則：

```text
Every final response must report the learning decision.
Not every turn must write a feedback lesson.
```

也就是「必填 report，不強制每次 write」。

### Alternatives Considered

- **A. 維持現狀，靠 agent 自律判斷是否回饋** — reject。Cognitive Mode 之所以成為 obligation，就是因為 final close-out 類判斷靠自律容易漏。
- **B. 每次都強制新增 feedback lesson** — reject。會造成 `feedback/history/` 膨脹，把純查詢、已覆蓋規則與低價值紀錄都寫成 durable knowledge。
- **C. 只在 non-local repo 時要求 report** — reject。Non-local 是高風險 trigger，但 local repo 中使用者指出 failure / workflow gap 時也應回報 learning decision。
- **D. 新增必填 report，但允許 `NONE` / `CHECK` / `NEEDED` 三種狀態** — accept。兼顧低噪音與高風險透明度。

### Why Not an ADR Yet

此變更會影響 runtime/core-bootstrap obligation、tool stop hook、final response shape、validation scenarios 與 possibly commit-msg / stop-hook validators。現在仍在設計 report 格式、non-local repo 判斷、hook fail-open / fail-closed 邊界，不宜先寫 accepted ADR。

### ADR Promotion Criteria

- [ ] `runtime/core-bootstrap.yaml` 新 obligation 已落地並被 tool stop hook 消費。
- [ ] Cursor / Claude / generic adapter 的 close-out behavior 已同步。
- [ ] 至少 3 個 validation scenarios 覆蓋：缺 report、non-local repo CHECK、需要 feedback 但未回報。
- [ ] 實際使用一段時間後，沒有造成 feedback lesson spam。
- [ ] 確認這是 cross-session / cross-agent foundational obligation，而不是單一工具行為。

### Consequences

#### 正面

- Agent final response 會顯性回答「有沒有需要沉澱」。
- Non-local / limited close-loop context 下，不會默默略過 feedback learning decision。
- 使用者指出失效時，agent 比較不會只修當下而漏掉 reusable prevention。

#### 負面

- 每個 final response 會多一行或一個小區塊。
- Stop hook / commit hook / tool adapter 需要同步更新，否則會出現各工具 enforcement 不一致。

#### 風險

- 若格式太複雜，agent 可能 ritualize 填報。
- 若 `NEEDED` 的 action 定義不清，可能變成每次都新增低價值 feedback lesson。
- 若 non-local repo 判斷過度嚴格，可能干擾純問答場景。

Glossary Impact: yes — proposed terms `feedback_learning_report`, `learning_decision`, `repo_context`; glossary entry 是否需要新增待 Phase 0 決定。

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

Non-local / limited repo context flow:

```text
repo context unavailable OR not local repo OR clean/push/readback cannot be verified
→ Feedback report must use CHECK or NEEDED
→ final response must explain whether learning/writeback is deferred
```

## Proposed Report Contract

### Compact Form

```text
Feedback: NONE / Repo:local / Reason:no reusable lesson
Feedback: CHECK / Repo:non-local / Reason:close-loop evidence limited
Feedback: NEEDED / Target:intelligence|workflow|enforcement|feedback-history|project-docs / Reason:<short>
```

### Full Form

Use full form when `Feedback: NEEDED`, high-risk failure learning, or non-local repo with deferred action:

```markdown
### Feedback / Learning Report

| Field | Value |
| --- | --- |
| repo_context | local / non-local / unknown |
| feedback_needed | yes / no / check / deferred |
| trigger | user correction / validation gap / workflow gap / repeated failure / none |
| target | none / feedback-history / intelligence / workflow / enforcement / project-docs |
| reason | <why> |
| action | written / not-needed / deferred-with-reason |
```

## Decision Rules

### Report Status

| Status | Meaning | Allowed when |
| --- | --- | --- |
| `NONE` | No reusable learning needed | Local repo context is clear, no new pattern/failure/gap discovered |
| `CHECK` | Agent must call out possible learning / handoff risk | Non-local repo, unknown repo state, limited close-loop evidence, or uncertainty |
| `NEEDED` | Durable learning or project feedback should be written | User correction, repeated failure, new validation/workflow gap, reusable pattern found |

### Triggers for `NEEDED`

- 使用者指出 agent 漏讀、漏驗證、誤判、重複犯錯。
- 發現 workflow、validation、evidence chain、routing、close-loop 或 feedback gap。
- 本輪產生新的 reusable reasoning pattern。
- 現有規則不足以防止同類錯誤。
- Non-local context 下發現需要 project docs / feedback history / handoff，但無法直接寫回。

### Triggers for `CHECK`

- Agent 不是在本地 repo 啟動。
- 無法確認 `git status` / clean / pushed / readback。
- 只取得 chat evidence，沒有 source repo / project docs 寫回能力。
- 使用者要求只提供建議，不授權寫入，但內容可能值得沉澱。

## Open Questions

- [ ] `Repo:local` 的判斷是否要求 `git status` clean，還是只要求有 local repo root？
- [ ] `Feedback: NONE` 是否允許在 non-local repo 出現？初步建議：不允許，至少 `CHECK`。
- [ ] Stop hook 應 fail-closed 還是 fail-open？初步建議：Ai-skill repo 內 fail-closed；外部專案若找不到 Ai-skill repo，依現有 stop hook policy。
- [ ] Compact line 是否需要固定順序以便 regex validator 簡單檢查？
- [ ] Commit-msg 是否也要要求 Feedback report，或只要求 chat/session final response？

## Phase 0 — Preflight

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [ ] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [ ] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Repo local 判斷 | pending | Phase 0 inspect current stop hook / repo detection |
| Non-local 是否允許 NONE | pending | Phase 0 define contract |
| Stop hook fail policy | pending | Phase 0 compare Cursor / Claude adapters |
| Compact line order | pending | Phase 0 validator design |
| Commit-msg 是否要求 | pending | Phase 0 scope decision |

### Phase 0.1 — Architecture Compatibility Preflight

- [ ] 讀 `runtime/core-bootstrap.yaml` per-turn obligations。
- [ ] 讀 stop hook implementation and tests：`scripts/ai-skill-cli/internal/app/hooks.go`、`hooks_test.go`。
- [ ] 讀 tool adapters：`ai-tools/agent/cursor.md`、`ai-tools/agent/claude.md`、generic adapters as needed。
- [ ] 讀 feedback learning rules：`enforcement/failure-learning-system.md`、`feedback/feedback-lessons.md`、`enforcement/reusable-guidance-boundary.md`。
- [ ] 確認 report belongs to runtime obligation, not workflow-specific checklist。
- [ ] 確認 linked updates：CORE_BOOTSTRAP companion、tool adapters、validation scenarios、CLI docs/tests。

## Phase 1 — Runtime Contract

- [ ] Add `obligation.feedback.learning_report` to `runtime/core-bootstrap.yaml`.
- [ ] Update `CORE_BOOTSTRAP.md` companion with compact/full report examples and rationale.
- [ ] Decide final compact grammar and enum values.
- [ ] Decide whether this obligation is chat/session-only or also commit-msg.

## Phase 2 — Hook Enforcement

- [ ] Extend stop hook final response validation to require Feedback / Learning Report.
- [ ] Add repair follow-up wording parallel to missing Cognitive Mode block.
- [ ] Add tests for missing report, compact report, full report, and non-local CHECK.
- [ ] Preserve existing Bootstrap Receipt and Cognitive checks.

## Phase 3 — Feedback Learning Routing

- [ ] Define how `NEEDED` maps to durable targets: `feedback-history`, `intelligence`, `workflow`, `enforcement`, `project-docs`.
- [ ] Update feedback / failure learning docs only if Phase 0 finds the existing rules insufficient.
- [ ] Ensure report does not force writing low-value lessons.

## Phase 4 — Validation Scenarios

- [ ] Add `validation/scenarios/runtime/feedback-report-required-v1.yaml`.
- [ ] Add `validation/scenarios/runtime/non-local-repo-feedback-check-v1.yaml`.
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
- [ ] Non-local / limited repo context requires `CHECK` or `NEEDED`, not silent omission.
- [ ] Validation scenarios exist and pass.
- [ ] Tool adapters document the new final close-out requirement.
- [ ] No secret, local path, project incident evidence, or private repo detail is written into reusable docs.
- [ ] Plan Completion Closure executed when implementation finishes.

## Stakeholder 同意項目

- [ ] Report is mandatory; writing a feedback lesson is conditional.
- [ ] Non-local repo / limited close-loop context must be explicit in final response.
- [ ] `Feedback: NONE` remains allowed for ordinary local low-risk turns.
- [ ] `NEEDED` must name a durable target or explain deferred action.

## Per-surface consumer 表

| Generated surface key | Named consumer(s) | Consumer 類型 |
| --- | --- | --- |
| `runtime.core_bootstrap.contract` updated obligation | stop hook final response validator | Go hook consumer |
| `validation/scenarios/runtime/*feedback*` | runtime audit / validation scenario inventory | validation scenario |

## 與其他 plans 的關係

- Related to [`archived/2026-05-25-2100-runtime-cognitive-contract-v2.md`](../archived/2026-05-25-2100-runtime-cognitive-contract-v2.md): mirrors the idea of final response reporting obligation, but for learning decision rather than cognition state.
- Related to [`archived/2026-05-31-1900-workflow-activation-engine.md`](../archived/2026-05-31-1900-workflow-activation-engine.md): final close-out hooks are one enforcement surface for runtime agent behavior.
- Related to feedback/failure learning rules: this plan adds final reporting, not a replacement for failure learning loop.
