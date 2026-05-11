# Runtime Routing

`runtime/routing/` defines how an agent should choose which Ai-skill knowledge to load for a task. It is a routing design layer, not executable policy.

## Routing Pipeline

```text
task intent
  -> knowledge/indexes/README.md
  -> metadata/schema.md fields
  -> metadata/ranking + confidence + compatibility
  -> source-of-truth gate
  -> selected primary source
  -> validation and readback
```

## Step 1: Identify Task Intent

Classify the task into a small intent before loading deep context:

- Bootstrap / session takeover.
- Skill execution.
- Skill update or promotion.
- Knowledge index / metadata / runtime work.
- Governance, validation, or close-loop work.
- Tool adapter or compatibility work.

If the task maps to an active `.agent-goals/` entry, use that goal as the current user-visible objective. Do not let stale goals override the latest user request.

## Step 2: Use The Knowledge Index

Read `../../knowledge/indexes/README.md` and find the closest `Task intent` row.

- Load `Primary source` first.
- Load `Related sources` only when the task needs them.
- If no row matches, fall back to root `README.md`, `shared-rules/README.md`, and the relevant layer README.
- If the row points to a candidate path, also load the old source-of-truth entrypoint.

## Step 3: Apply Metadata Ranking

Use `../../metadata/ranking/README.md` to decide read order:

1. Required shared rules and safety/source/validation gates.
2. Current source-of-truth entrypoint.
3. Validated or stable atoms.
4. Candidate maps and summaries.
5. Optional background references.

Use `../../metadata/confidence/README.md` to avoid treating low-confidence candidates as stable behavior.

## Step 4: Apply Compatibility Gate

Use `../../metadata/compatibility/README.md` and `../../governance/lifecycle/README.md`:

- If old `skills/` or `shared-rules/` source remains active, it wins over candidate new-layer content.
- If a new layer path is only a candidate map, it can guide discovery but cannot override behavior.
- If a promoted atom exists, confirm old links still resolve before relying on it as a replacement.

## Step 5: Validate The Route

Before acting on a routed source, confirm:

- The source is canonical, not a tool mirror.
- Required dependencies have been read or marked not applicable.
- The old entrypoint is preserved when migration is in progress.
- The selected source gives a clear validation signal.
- The final response or commit can explain what was loaded and what was deferred.

## Runtime Output Shape

For important routing decisions, report briefly:

```text
Task intent:
Primary source:
Related sources loaded:
Candidate sources deferred:
Source-of-truth gate:
Validation signal:
```

## Stop Conditions

Stop loading more sources when:

- The primary source answers the current decision.
- More context would duplicate the same source-of-truth.
- A candidate path would conflict with an active skill entrypoint.
- Required validation fails.
- The user changes the task priority.
