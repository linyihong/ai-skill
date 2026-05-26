# Bootstrap Check (session trigger)

This file is auto-loaded into every Claude Code session as a rehydration trigger.

## Required check (before any response)

If you have NOT yet emitted a Bootstrap Receipt in this session, you MUST execute
the runtime_loader defined in `CLAUDE.md` BEFORE responding to the user — even if
the user's question seems trivial.

The runtime_loader points to `runtime/core-bootstrap.yaml`, which declares all
`required_reads`, `per_session_obligations`, and `per_turn_obligations`.

## Why this file exists

- `CLAUDE.md` is "entry pointer" (architecture).
- This file is "session trigger" (behavior reminder, loaded every session).
- `PreToolUse` hook is mechanical enforcement (blocks tools without Bootstrap Receipt).

All three layers cooperate. This file's only job: remind Claude on every session
to check whether bootstrap has run, and if not, run it.
