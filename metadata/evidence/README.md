# Evidence Metadata

`metadata/evidence/` defines metadata-only evidence qualification policies. It does not execute validation and is not compiled into `runtime.db` unless a future compiler target explicitly adds it.

## Files

| File | Purpose |
| --- | --- |
| [`domain-policies.yaml`](domain-policies.yaml) | Domain evidence authority, freshness, validity, scope, and observability policy. |

## Boundary

| Concern | Owner |
| --- | --- |
| Generic evidence hierarchy and forbidden behavior | `enforcement/evidence-hierarchy.md` |
| Cognitive governance model | `governance/ai-runtime-governance/cognitive-state-governance.md` |
| Domain recovery reload sets | `metadata/recovery/domain-policies.yaml` |
| Domain evidence qualification | `metadata/evidence/domain-policies.yaml` |
| Runtime enforcement | `runtime/guards/` only after explicit promotion |

## Policy Schema

Each policy should include:

- `domain`: stable domain id.
- `applies_when`: signals that select this policy.
- `authority_order`: strongest to weakest evidence owners for this domain.
- `freshness_requirements`: when evidence expires or needs refresh.
- `scope_boundaries`: what evidence can and cannot prove.
- `invalid_claims`: claim compressions that must be blocked.
- `validation`: how to verify the policy was applied.

## Runtime Boundary

This directory is metadata-only in the current plan. Agents may read it during recovery, validation, or planning, but no runtime guard should claim it is machine-enforced until:

1. A compiler target is added.
2. `runtime.db` includes the compiled surface.
3. Validation scenarios prove the behavior.

