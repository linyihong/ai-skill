# External Plan Validation — Adoption Guide (Phase 3.1)

External repositories validate their plan trees by **invoking the shared
`ai-skill` binary**, not by vendoring code or importing governance. What you get
is a callable **validation capability** + an **invocation contract** — nothing
else. See `plans/active/2026-06-22-1009-plans-system-portability-and-delivery-integration/01-...md`.

## What is (and is not) externalized

- **Externalized**: the plan-validation engine (`plans validate`) + how to call it.
- **NOT externalized**: governance, runtime.db, routing, cognitive modes, the
  commit-msg governance validators. None of these are installed in your repo.

## Non-goals (hard)

Adoption must NOT add to your repo: an external registry, external runtime state,
or external plan metadata. The integration is a thin call-out and nothing more.

## Invocation contract

```
<AI_SKILL_BINARY> plans validate --root <REPO_ROOT> [--format text|json]
```

- **Input**: your repo root (the engine reads `<REPO_ROOT>/plans/active` and
  `plans/archived`).
- **Output**: findings (`text` or `json`). The engine is policy-free.
- **Exit**: non-zero iff a blocking finding exists. *You* (the consumer) decide
  what that means — block the commit, fail CI, or just warn.

## Invocation adapter (replaceable)

The adapter is a **thin transport** you own; it must stay replaceable and must
never push policy back into the engine. Two equivalent forms:

### git hook shim — `<REPO_ROOT>/.git/hooks/commit-msg`

```sh
#!/bin/sh
# Thin invocation adapter — calls the shared binary, owns no validation logic.
exec "<AI_SKILL_BINARY>" plans validate --root "$(git rev-parse --show-toplevel)"
```

### CI wrapper (any CI)

```sh
"<AI_SKILL_BINARY>" plans validate --root "$PWD" --format json
```

> Replace `<AI_SKILL_BINARY>` with the path to the platform binary under
> `scripts/ai-skill-cli/bin/`. No daemon, service, or background process is
> involved; nothing persistent is installed.

## Install → validate → remove (reversible adoption)

```sh
# install
cp commit-msg <REPO_ROOT>/.git/hooks/commit-msg && chmod +x <REPO_ROOT>/.git/hooks/commit-msg
# validate (a commit now runs the engine over your plan tree)
git commit ...
# remove (monotonic: removing the adapter only removes the capability)
rm <REPO_ROOT>/.git/hooks/commit-msg
```

**Rollback guarantee**: removing the adapter returns the repo to clean —
`git status` clean, no runtime residue, no hook residue, and **no schema
residue** (your plan frontmatter is never migrated or rewritten). Removal is
**monotonic**: it produces no new validation errors; only the capability
disappears.

## Schema

Plans use the canonical frontmatter (`id` / `plan_kind` / `parent` /
`required_for_completion` / `sub_plan_reason` / `schema_version`). Whether
non-canonical dialects are supported is an open policy question (deferred); this
guide assumes the canonical schema.
