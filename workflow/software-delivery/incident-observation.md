# Incident Observation Slice（Stage 0 — Observe）

> **Cognitive Slice**：`sd-incident-observation`（UI / consumer continuity or navigation incident 的 observable-first 收集階段，對應 Phase B 2026-06-18 transfer）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-incident-observation` |
| `purpose` | 在分類或改 code 前，只收集 observable behavior：symptom、timeline、metrics；禁止 implementation-first |
| `type` | `execution` |
| `tags` | incident, observable, ui, navigation, continuation |
| `load_when` | 未知 UI / consumer incident；bug 修復但 domain 或 modification layer 未決；G4 / drill / integration 需 incident card |
| `do_not_load_when` | 已有 playbook 與 owning layer；純 contract 編輯；planning-only；scripted validation 已有固定 checklist |
| `owner_layer` | workflow |
| `layer_justification` | 規定 incident 工作「先觀察、後分類」的 ordering；不承載 domain taxonomy（project overlay）或 runtime gate（governance） |
| `canonical_source` | 本檔 |
| `dependencies` | `sd-intake`（Change Intake 已標記 bug）、[`ui-incident-governance-workflow.md`](ui-incident-governance-workflow.md)、[`layer-ownership-matrix.md`](layer-ownership-matrix.md) |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Incident card 存在且 Steps 1–3 有 URL / context proxy / viewport metrics 才允許進入 Classify |

## Boundary

```text
sd-intake
  = Change classified as bug; expected vs actual at high level

sd-incident-observation
  = Collect observable only — no hooks, storage, PR, root-cause guess

sd-ui-incident-governance
  = Classify domain + select modification layer

software-delivery-governance
  = Layer selection gate (single-layer convergence)
```

Project-specific N/C/R scope rules and verification hooks stay in `<PROJECT_ROOT>/.ai-skill/project/` until second independent incident.

---

## Input

| Field | Required |
| --- | --- |
| **Symptom** | yes — one sentence, user- or test-visible |
| **Timeline** | yes — repro steps in observable order |
| **Observable behavior** | yes — per-step URL, visible state, metrics |

---

## Forbidden before incident card complete

- Read hook / component internals to **decide** domain
- Read sessionStorage / localStorage to **skip** observable steps
- Open PR or blame to **guess** root cause
- Propose new abstraction or hub expansion

Allowed: CDP, integration read-only probes, screenshot, DOM metrics, G4 evidence scripts.

---

## Collection checklist

Per repro step record:

1. **Route** — pathname + query (`?tab=`, `?keyword=`, etc.)
2. **Context proxy** — chip, category, keyword bar, in-tab semantic state
3. **Viewport** — scrollTop, maxScrollTop, visible section / region
4. **Overlay** — player, modal, sheet visibility if relevant

---

## Output — incident card

Minimal template:

```markdown
## Incident card

- Symptom:
- Timeline:
- Step observables:
  | Step | URL | Context | Viewport | Overlay |
  | --- | --- | --- | --- | --- |
  | 1 | | | | |
- First broken authority (if obvious from observables only): Route | Context | Viewport | unknown
- Not yet classified — do not implement
```

Store in project incident note, drill memo, or G4 evidence sheet — not in Ai-skill canonical body.

---

## Next stage

When incident card is complete → [`ui-incident-governance-workflow.md`](ui-incident-governance-workflow.md) **Classify** (Stage 1).

Governance gate for layer choice → [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md) **Incident layer selection**.
