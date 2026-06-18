# UI Incident Governance Workflow（Stage 1 — Classify · Stage 2 prep — Select Layer）

> **Cognitive Slice**：`sd-ui-incident-governance`（Classify domain + route to layer selection；operational companion to `sd-incident-observation`）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-ui-incident-governance` |
| `purpose` | 在 implementation 前完成 blind classify 與 authority-first mapping；產出 classification 與 primary modification layer 候選 |
| `type` | `execution` |
| `tags` | incident, ui, navigation, continuation, classification, layer-selection |
| `load_when` | Incident card 已完成；Navigation / Continuation / Recovery 未決；layer 選擇有爭議 |
| `do_not_load_when` | Playbook 已知且 layer 已寫在 project overlay；非 UI / consumer surface |
| `owner_layer` | workflow |
| `layer_justification` | 規定 classify 順序與產物；single-layer gate 正文在 governance |
| `canonical_source` | 本檔 |
| `dependencies` | `sd-incident-observation`、[`layer-ownership-matrix.md`](layer-ownership-matrix.md)、[`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md) |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | `classification.md` 等價產物：恰好一個 domain + 恰好一個 primary layer |

---

## Placement in software-delivery lifecycle

**Architecture shift**: incident path is **evidence-driven change**, not linear requirement-driven delivery.

```text
Discover          ← sd-intake (Change Intake: bug)
  ↓
Observe           ← sd-incident-observation (incident card)
  ↓
Classify          ← this file §Stage 1
  ↓
Select Layer      ← this file §Stage 2 + layer-ownership-matrix + governance gate
  ↓
Execute           ← Contract · Implementation · Verification (sd-contracts / sd-implementation / sd-validation)
  ↓
Ship              ← sd-closure DoD
  ↓
Retrospective     ← change-retrospective.md (keep local | promote project | candidate canonical)
```

**Anti-pattern (blocked)**:

```text
incident → implementation
```

---

## Stage 1 — Classify

### Purpose

Choose **one** governing domain before naming components or opening storage.

### Output (classification)

Allowed values — **exactly one**:

| Domain | When |
| --- | --- |
| **Navigation** | Route / back-stack / entry URL / search exit / secondary page return |
| **Continuation** | Tab suspend-resume, in-tab state restore, overlay return with keep-alive scope |
| **Recovery** | Cold load, bookmark, refresh, empty history — bounded recovery not keep-alive |
| **Out-of-scope** | Not a UI continuity/navigation incident — exit to intake or other workflow |

### Rules

- **Forbidden**: select two domains simultaneously
- **Forbidden**: classify by opening hooks or sessionStorage first
- **Forbidden**: open new abstraction to «resolve» ambiguity — document gray zone and use authority mapping

### Authority-first (before domain final)

From incident card, name **first broken authority** (see [`layer-ownership-matrix.md`](layer-ownership-matrix.md)):

1. Route
2. Context
3. Viewport
4. Persistence — only after 1–3 ruled out

Gray zone: same observable may map to Navigation **or** Continuation (e.g. `?tab=` on back) — authority sequence decides per incident; no global rule «context lost → always Continuation».

### Project scope overlays

Repo-specific classify order (e.g. search-is-not-keep-alive) lives in `<PROJECT_ROOT>/.ai-skill/project/rules/ui-incident-governance-workflow.md` — load when present; do not duplicate into Ai-skill body.

### Artifact — classification.md

```markdown
## Classification

- Incident ref:
- Domain: Navigation | Continuation | Recovery | Out-of-scope
- First broken authority: Route | Context | Viewport | Persistence
- Gray zones noted:
- Rejected domains:
```

---

## Stage 2 — Select Layer (workflow side)

### Purpose

Choose **one primary** modification layer before contract or code work.

### Allowed primary layers

| Layer | Meaning |
| --- | --- |
| **Contract** | Owning frontend-contract / BDD / invariant |
| **Overlay** | Project `.ai-skill/project/rules/` implementation pattern |
| **Verification** | Test / G4 / integration evidence only |
| **Integration** | App wiring, navigation helper, storage timing fix |

Use [`layer-ownership-matrix.md`](layer-ownership-matrix.md) to constrain allowed layers for the broken authority.

### Rules

- One primary layer per incident
- Secondary layers documented but must not dilute primary fix path
- **Gate**: pass [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md) **Incident layer selection** before implementation

### Artifact — layer selection record

```markdown
## Layer selection

- Primary layer: Contract | Overlay | Verification | Integration
- Single-layer convergence: YES | NO
- If NO: see verification expansion exception (below) — not «open new plan» by default
- Rejected: new abstraction / invariant promote / hub expand
- Opens: (contract path | overlay path | test path | code path)
```

### Single-layer convergence — NO exception

When convergence is **NO** (evidence insufficient to pick one primary fix layer):

| Action | Allowed |
| --- | --- |
| Add integration test / G4 evidence sheet | ✅ |
| Expand verification only | ✅ |
| Add project overlay | ⚠️ review — not default |
| Promote / change contract | ❌ |
| New abstraction / runtime hub | ❌ |

**Rationale**: insufficient evidence ≠ cannot expand verification. Forbidden escape: «cannot change → new plan / new abstraction».

Governance gate: [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md) §Incident layer selection.

## Promotion boundary

| Content | Ai-skill workflow | Project overlay |
| --- | --- | --- |
| Observe → Classify → Select Layer process | ✅ this file + incident-observation | Scope tables, NAV ids, drills |
| Layer ownership matrix (generic) | ✅ layer-ownership-matrix.md | Project overrides |
| Experience Runtime Governance | ❌ defer | ❌ |

Promotion to enforcement / runtime requires second independent incident + plan — not automatic after one pilot.
