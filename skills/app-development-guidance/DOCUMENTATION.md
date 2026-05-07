# Documenting App Development Guidance Notes

Use this guide when turning an APK analysis lesson, implementation observation, embedded/firmware observation, hardware-product lesson, or design review into reusable development guidance.

## Reusable Note Structure

```markdown
### Short title

Status: candidate | validated | promoted | deprecated | experimental

#### Observed Risk

What pattern was observed, without target-specific secrets.

#### Development Consequence

Why this matters when building our own app, API, SDK, firmware, or hardware-backed product.

#### Recommended Implementation / Control

What to implement, and which layer owns it.

#### Validation

How to prove the control exists or fails safely.

#### Limits

What this does not solve.
```

## Keep Separate

| Content | Put it in |
| --- | --- |
| Cross-platform security principle | `controls/`. |
| Platform, app-type, embedded, firmware, or hardware-product implementation detail | `platforms/`. |
| Language/runtime-specific pitfall | `languages/`. |
| Concrete buildable implementation pattern | `implementation/`. |
| Repeatable review process | `checklists/`. |
| Reusable development principle not yet promoted | Matching `feedback_history/<category>/`, or `feedback_history/common/` for cross-cutting lessons. |
| APK analysis method or hook technique | [`apk-analysis`](../apk-analysis/). |
| Project-specific board wiring, calibration logs, firmware dumps, bench measurements, device identifiers, or target hardware conclusions | Project repository docs. |
| Product-specific API host, endpoint, schema, or token detail | Project repository docs. |
| Raw vendor documents, account-specific terms, credentials, sandbox/live hosts, private webhook payloads, or real customer data | Project repository docs with sanitization and access control. |
| Generated clients, SDKs, fixtures, and provider/consumer contract checks | `implementation/` and the project repository. |
| Product Brief validation, document precedence, traceability, and BDD closure process | `process/`, templates, and checklists. |
| Performance budgets, load/stress/spike/soak strategy, CI smoke checks, and release evidence | `process/`, `CHECKLIST.md`, templates, and the project repository's test or release notes. |
| Shared sanitization or feedback rules | [`shared-rules`](../../shared-rules/README.md). |
| Local-only scratch notes, credentials, or ephemeral process artifacts | Project repository only: **gitignored** paths with **neutral** directory naming; configuration via env vars and optional untracked files; keep tracked READMEs free of filesystem tours and internal codenames. |

## Reusable Guidance Boundary

This section applies the global rule in [`reusable-guidance-boundary.md`](../../shared-rules/reusable-guidance-boundary.md).

Skill documentation should describe the reusable reason, failure mode, decision rule, and validation method. Do not promote a project incident by copying its app name, module name, endpoint path, sample payload, class name, local path, host, or live-data quirk into the skill.

If an incident teaches a useful lesson, split it:

- **Skill:** generalized rule, such as "SDK bugs reported against a live service should be reproduced through the SDK public surface and then pinned with behavior specs and regression tests."
- **Project docs:** concrete reproduction target, affected feature, sample IDs, live environment notes, BDD file names, test class names, and execution results.

Start with [`templates/README.md`](templates/README.md) when choosing a copyable documentation shape. Use [`templates/initial-development-docs.md`](templates/initial-development-docs.md) when starting from or validating a product brief, [`templates/hardening-note.md`](templates/hardening-note.md) for reusable guidance, and [`templates/threat-model-lite.md`](templates/threat-model-lite.md) for quick feature reviews.

## Required Linked Update Statement

Every reusable note that affects multiple folders must follow [`../../shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md) and include a short linked-update statement:

```markdown
#### Required Linked Updates

- `controls/...`: updated or checked because ...
- `implementation/...`: updated or checked because ...
- `checklists/...`: updated or checked because ...
```

If linked updates are not needed, state why. This makes it clear that related docs are required to stay in sync, not optional follow-ups.

## Good Guidance

Good development guidance is:

- Actionable by an engineer.
- Explicit about the owner layer.
- Clear enough to turn into code, configuration, tests, or review checklist items.
- Testable.
- Honest about residual risk.
- Sanitized and free of target-specific details.
- Explicit when performance evidence is required, including the metric, budget, environment, runner, and release gate.

## Avoid

- Tracked markdown that maps unpublished workflow to evocative directory names, lists developer-machine paths, or repeats internal investigation stories—those belong in local-only notes under gitignore, not in the default-branch narrative.
- "Use obfuscation" without naming what it protects and what it does not.
- "Performance is fine" based only on functional tests or average latency, without P95/P99, throughput, error-rate, resource, baseline, or environment context.
- "Add pinning" without a rotation plan or threat model.
- "Detect root" as a hard authorization decision.
- "Hide the secret in the app" as a durable security control.
- Copying raw findings from a third-party APK into reusable docs.
