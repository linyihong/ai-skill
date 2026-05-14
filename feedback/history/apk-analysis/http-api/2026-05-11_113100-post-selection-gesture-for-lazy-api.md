> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-11 - Post-Selection Gesture for Lazy API

Status: candidate

#### One-line Summary

選中分類、filter、grid label 或 tab 後若沒有立刻出現 feature API，先補一個低風險後續 gesture，再判定是否真的無 API。

#### Human Explanation

有些 App 在 selection 當下只更新 UI state、cache key 或 selected route，真正的列表 API 會延後到下一個列表 scroll、refresh、曝光或 layout settle 才發出。只看 selection-only 窗口，容易把尚未 lazy-load 誤判成 `no-network-observed`。

#### Trigger

- UI 已選中分類、filter、grid label、chip 或 tab。
- foreground package 與 feature context 都正確。
- 高語意 request/decrypt hook 或 pcap 在 selection 後短窗口沒有新增 feature API。

#### Evidence

- Tool: UI replay script + high-semantic request/decrypt hook.
- Sanitized excerpt: selection-only window did not reliably produce the expected list request; adding one bounded post-selection list gesture produced the route key class and list-array wrapper shape.
- Evidence path: concrete project evidence stays under `<PROJECT_ROOT>` project docs or capture summaries; reusable lesson keeps only the generalized trigger.

#### Generalized Lesson

For UI-to-API attribution, separate `selection-only` evidence from `post-selection-triggered` evidence. A `no additional network observed` result after a UI selection does not rule out the scoped API until at least one bounded follow-up gesture has been tested.

#### Agent Action

Add optional replay-runner knobs such as `--post-select-scrolls`, `--post-select-wait`, or `--post-select-refresh`; keep defaults unchanged; capture before selection, after selection, and after the follow-up gesture with package/context guards.

#### Goal / Action / Validation

- Goal: avoid false negatives when category/filter APIs lazy-load after a later gesture.
- Action: run a bounded post-selection gesture and observe request/response shapes.
- Validation or reference source: compare request/decrypt sequence before selection, immediately after selection, and after the follow-up gesture.

#### Applies When

- The UI element changes list scope, category, filter, result tab, or feed segment.
- The operation is read-only and a bounded scroll/refresh is in authorization scope.

#### Does Not Apply When

- The selected action is a write action, payment, login, destructive operation, or outside authorization scope.
- The app leaves the target package or feature context.

#### Validation

Confirm that the follow-up gesture produces a new feature request, a stable no-network-after-follow-up result, or a documented wrong-screen/external transition. Do not count startup, preload, or background requests without timing and context alignment.

#### Promotion Target

- `WORKFLOW.md`

#### Required Linked Updates

- Updated `feedback_history/README.md` and `feedback_history/http-api/README.md`.
- Promoted the rule into `WORKFLOW.md`.
- Checked reusable guidance boundary: project-specific evidence remains in project docs.
