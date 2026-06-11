---
id: plan-archival-link-drift
status: active
owner_layer: enforcement
empirical_origin: 2026-06-11
related_rule_class: plan_archival_link_integrity
---

# Plan Archival Link Drift

> Family: **Reference Integrity** (sibling: `runtime_index_freshness`, topology references, plan-tree parent/child references).

## Failure pattern

Archiving a plan (move `plans/active/<id>.md` → `plans/archived/<id>.md`) silently breaks two surfaces:

1. **Outbound** — the moved file's own relative links (parent / sibling plan references, runtime / metadata cross-links) were written assuming the file lived in `plans/active/`. After the move, `../archived/`, same-dir, and `../active/` paths all need to be recomputed against the new location.
2. **Inbound** — every other repo file (README, runtime yamls with `source_plan`, failure-patterns, metadata, topology references) that pointed at the **old** path now resolves to a non-existent target.

Both surfaces drift without any code-level signal. Manual grep + per-link edit is the only mitigation absent mechanical enforcement. Half a year later, someone clicks a 404 or a downstream tool's path resolver fails.

## Empirical evidence

**2026-06-11** — archiving `plans/active/2026-06-06-1800-sanitization-mechanical-enforcement.md` exposed exactly this pattern:

- **8 inbound references** required hand fix (`enforcement/README.md`, 2 runtime yaml `source_plan` fields, 2 failure-pattern entries, 2 metadata files, topology-migration plan).
- **3 outbound links** in the moved file itself required hand fix (parent plan reference + sibling plan reference + an `../active/` link to a peer that had not been archived).

Resolution commit: `3f7c4b4 chore(plans): archive completed sanitization plan + retarget references`. The fix was correct, but entirely manual — there was no executor to surface the drift at commit time.

## Why pre-existing validators did not catch it

`validatePlanArchivalAudit` (under `plan_governance` rule_class) checks **archive workflow completeness** — that the archived plan has no unchecked `- [ ]` items without justification. It does **not** parse markdown links, and does not look at the rest of the repo for inbound references. That validator is a *workflow* check; link integrity is a *Reference Integrity* check. Conflating the two in one validator would have made `plan_archival_audit` a "kitchen sink" with mixed concerns.

`validateRuntimeIndexFreshness` covers staged-blob ↔ index checksum drift for files in the runtime index `sources` table. Plans are not generally tracked there, so it does not surface this case.

## Mechanical enforcement (landed)

Validator: `validatePlanArchivalLinkIntegrity` (commit-msg).

- **Rename detection**: `git diff --cached --find-renames=90% --name-status` over `plans/active/` ↔ `plans/archived/`. Batch rename map built before any resolution — multi-archive in the same commit (A and B both moving, A linking B) is handled correctly.
- **Outbound check** (block): for each moved plan, parse markdown links via a bounded parser; resolve each `[text](path)` from the **new** location; flag any target that does not exist in the repo or that hits the rename map.
- **Inbound check** (block): scan every repo `.md` for markdown links whose target resolves to a renamed plan's **old** path; emit finding with `suggested_replacement` derived from the rename map.
- **Bare textual scan** (warning): plain-text occurrences of the old path that are **not** inside a markdown link target are reported as `stale_textual_reference`. An opt-in `<!-- archival-provenance -->` marker on the same or preceding line downgrades to `historical_provenance_reference` (info, suppressed from output) for intentionally retained historical mentions.
- **Staged-blob first** (TD-1 Resolution Gate outcome): `readFileForScan` reads `git show :<path>` first and only falls back to the worktree on read failure, so partial-stage (`git add -p`) and post-stage worktree edits cannot produce false-pass or false-block.
- **Payload**: every finding carries `{Severity, Category, File, Line, Column, Target, SuggestedReplacement}` so downstream IDE / auto-fix tools can rewrite without re-resolving.

Opt-out: standalone `[skip-plan-archival-link-integrity]` trailer for emergency archives.

## How to recognise this pattern in future failures

- Symptom: post-archive bug report or 404 in docs that links to `plans/active/<old-id>.md`.
- Symptom: runtime tool fails to resolve a `source_plan` field after a plan archive.
- Symptom: hand-grep + per-file edit during archive review.

If any of the above recurs, treat it as a regression in `plan_archival_link_integrity` enforcement, not a generic "broken link" issue.

## Related

- Rule class: [`plan_archival_link_integrity`](../enforcement-registry.yaml) (mechanical).
- Plan: [`plans/active/2026-06-11-1100-plan-archival-link-integrity.md`](../../plans/active/2026-06-11-1100-plan-archival-link-integrity.md) (Phase 1–3, TD-1 Resolution Gate).
- Sibling family member: `runtime_index_freshness` (source ↔ index checksum drift).
- Workflow-completeness sibling (different concern): `validatePlanArchivalAudit` under `plan_governance`.
