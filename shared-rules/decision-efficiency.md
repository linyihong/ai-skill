# Decision Efficiency

Use this rule to choose the next useful action without overloading context, tools, or the user. The goal is to preserve decision quality while reducing unnecessary reading, broad exploration, duplicate work, and token-heavy context.

This rule generalizes the decision routing pattern used in technical skills such as APK analysis: start from the current unknown, compare evidence paths, choose the highest-yield route, and load only the documents needed for that route.

## Core Rule

Before doing substantial work, state the current decision point in one sentence:

```text
Current unknown: <what must be learned or decided next>
```

Then choose the next action by comparing:

| Criterion | Question |
| --- | --- |
| Time to evidence | Which path can answer the unknown fastest? |
| Semantic distance | Which source is closest to the real decision, not just noisy symptoms? |
| Safety / reversibility | Which action is least destructive and easiest to roll back? |
| Validation signal | Which path gives a clear pass/fail or confidence update? |
| Context cost | Which files/tools are actually needed, and which can wait? |
| User value | Which result best advances the user's goal or removes a blocker? |

Prefer the route with the best evidence-to-cost ratio, not the route that happens to be first in a checklist.

## Context Loading

Load context in layers:

1. **Bootstrap:** read the shared-rule bootstrap set.
2. **Task frame:** read the user request, active `.agent-goals/` entry, and directly relevant open files.
3. **Skill entry:** read the matching `SKILL.md` and its routing guidance.
4. **Route-specific docs:** read only the workflow/tools/docs category needed for the current route.
5. **Deep references:** read examples, techniques, feedback lessons, or source files only after evidence says they matter.

Do not read every category or every technique by default. If a broad read seems necessary, explain why broad context is required and what decision it supports.

## Decision Routing

Use workflows as routing aids, not rigid scripts. When a workflow has many branches:

- Start with the highest-level triage.
- Stop once evidence clearly points to a branch.
- Read only that branch's detailed docs.
- Keep other branches as fallbacks, not active context.
- Re-route when evidence contradicts the current branch.

If an action has already answered the decision point, do not keep running broader or lower-level checks just because they are available.

## Token And Noise Reduction

Reduce context and output load by:

- Summarizing large evidence before expanding it.
- Reading indexes before child files.
- Using exact search for known names and semantic search only for broader discovery.
- Keeping raw logs, large payloads, screenshots, and generated dumps in project artifacts, then citing paths or sanitized excerpts.
- Recording open questions instead of reading unrelated files to fill speculative gaps.
- Moving reusable but lengthy material into focused child files per [`document-sizing.md`](document-sizing.md).

Do not use token reduction as an excuse to skip required dependencies. If a dependency is required by [`dependency-reading.md`](dependency-reading.md), read it or mark it blocked / not applicable.

## Stop Conditions

Stop the current route and reassess when:

- The current route produced enough evidence to answer the unknown.
- The route is generating noise without improving confidence.
- The route becomes destructive, unstable, or too slow compared with an available alternative.
- A higher-semantic source appears.
- A user priority or blocker changes the goal.
- Another active goal/owner/lock makes the work unsafe to continue in parallel.

When stopping, record the reason in the answer, active goal, or document TODO if the decision affects future work.

## Output Shape

For important decisions, report:

```text
Current unknown:
Options considered:
Chosen next action:
Why this is the highest-yield path:
What was deferred:
Validation signal:
```

Keep this short. The purpose is to make the route choice auditable, not to produce a long reasoning dump.

## Relationship To Other Rules

- Use [`dependency-reading.md`](dependency-reading.md) for required dependency scope.
- Use [`document-sizing.md`](document-sizing.md) when decision guidance grows too large or route-specific.
- Use [`goal-action-validation.md`](goal-action-validation.md) for goal/action/validation closure.
- Use [`conversation-goal-ledger.md`](conversation-goal-ledger.md) when route changes create new goals, blockers, or handoffs.

← [Back to shared rules index](README.md)
