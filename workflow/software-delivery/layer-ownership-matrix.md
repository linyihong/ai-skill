# Layer Ownership Matrix — UI Incident Authorities

Maps **which authority broke first** → **domain owner** → **allowed modification layers**. Used after [`incident-observation.md`](incident-observation.md) and during [`ui-incident-governance-workflow.md`](ui-incident-governance-workflow.md) Stage 2.

**Scope**: UI / consumer continuity and navigation incidents during software delivery. Not security control owner layers (see [`analysis/development-guidance/risk-translation.md`](../../analysis/development-guidance/risk-translation.md)).

Project repos may extend rows in `.ai-skill/project/rules/authority-mapping.md` — do not fork this matrix into project without divergence evidence.

---

## Matrix

| Authority | Observable proxy | Domain owner | Allowed primary modifications |
| --- | --- | --- | --- |
| **Route** | pathname + query; back lands wrong URL | **Navigation** | Contract · Verification · Integration |
| **Context** | chip, category, in-tab URL semantic state | **Continuation** (gray: Navigation) | Contract · Overlay · Integration |
| **Viewport** | scrollTop, visible region | **Continuation** | Overlay · Verification |
| **Persistence** | storage read/write timing (after 1–3 mapped) | **Runtime / Integration** | Overlay · Integration |

---

## Mis-routing traps (blocked)

| Symptom jump | Wrong layer | Correct first step |
| --- | --- | --- |
| scroll bad → change contract | Contract | Confirm viewport authority; often Overlay |
| tab bad → only Navigation overlay | Overlay | Check context authority in gray zone |
| storage key wrong → new runtime hub | Abstraction | Integration or Overlay primary |

---

## Gray zone handling

When observable fits two domain owners:

1. Record both candidates in classification.md
2. Use **first broken authority** row, not symptom label
3. Pick one primary layer; document rejected branch
4. Do not promote «always Continuation» or «always Navigation» to Ai-skill invariant from one incident

---

## Refs

- Classify workflow: [`ui-incident-governance-workflow.md`](ui-incident-governance-workflow.md)
- Governance gate: [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md) §Incident layer selection
- Responsive UI authority (different domain): [`responsive-ui/authority-mapping.md`](../../intelligence/engineering/governance/responsive-ui/authority-mapping.md)
