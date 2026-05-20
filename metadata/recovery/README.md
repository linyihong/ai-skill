# Recovery Metadata

`metadata/recovery/` defines domain-specific recovery policy for mismatch escalation. Runtime owns the generic recovery state machine in `runtime/compiler/embedded_data.rb`; this metadata layer tells agents which source-of-truth set to reload for a given domain before rebuilding the execution graph.

## Files

| File | Purpose |
| --- | --- |
| [`escalation-levels.yaml`](escalation-levels.yaml) | L1-L5 escalation level metadata, default actions, and minimum reload requirements. |
| [`domain-policies.yaml`](domain-policies.yaml) | Domain-specific trigger classes, required reload sets, forbidden behaviors, and validation gates. |

## Policy Schema

Each domain policy should include:

- `domain`: stable domain id such as `apk-analysis` or `software-delivery`.
- `applies_when`: trigger conditions that select the policy.
- `trigger_classes`: escalation trigger classes covered by the policy.
- `required_reload_set`: source-of-truth files or source categories that must be read or explicitly marked `not_applicable` / `source_missing`.
- `rebuild_graph`: fields that must be present before execution resumes.
- `forbidden_behaviors`: domain-specific actions that must stop during recovery.
- `validation`: checks proving recovery has closed before execution continues.

## Runtime Boundary

These files are metadata-only in Phase 4. They are not compiled into `runtime.db` yet. Use them through routing / validation and keep the executable recovery procedure in `runtime/compiler/embedded_data.rb`.

If a future phase needs runtime enforcement, add a compiler target deliberately instead of assuming all metadata YAML is compiled.
