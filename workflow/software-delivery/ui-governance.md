# UI Governance Slice

> **Cognitive Slice**：`sd-ui-governance`（從 UI contract surface 延伸出的 workflow governance slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §4 的 workflow membership test）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-ui-governance` |
| `purpose` | 在 UI / consumer contract 已存在或正在建立時，定義 UI compliance 如何被分類：governance domain、render context、collection method、validation mechanism、evidence class、severity policy 與 runtime projection boundary |
| `type` | `execution` |
| `tags` | artifact-gate, ui, design-system, accessibility, behavior-governance, responsive, visual-validation |
| `load_when` | UI / consumer surface 需要 design-system enforcement、accessibility evidence、responsive evidence、behavior pattern checks、visual baseline review、AI visual review scoping，或 completion claim 依賴 UI compliance |
| `do_not_load_when` | 無 user-visible consumer surface、純 provider 內部變更、只修不影響 UI behavior / contract / design-system compliance 的小錯 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「UI compliance 評估前要分類哪些 domain、collection method、mechanism、evidence、severity，並保持 runtime projection advisory-only」的 ordering / gate；通過 workflow membership test。Phase 1 只擁有 UI-local taxonomy usage；長期 shared taxonomy owner 可能是 validation-reasoning |
| `canonical_source` | 本檔 |
| `dependencies` | `sd-ui-contracts`（UI contract expectations）、`sd-test-strategy`（proof target selection）、`sd-validation`（evidence acquisition / evaluation execution） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | UI governance scenarios can prove domain/mechanism separation, evidence classification, and AI visual warning boundary before runtime/enforcement promotion |

## Boundary

`sd-ui-contracts` defines what the UI should be. `sd-ui-governance` defines how UI compliance is classified before validation executes acquisition and evaluation.

```text
sd-ui-contracts
  = What UI should be

sd-ui-governance
  = How UI compliance is classified

validation-reasoning (future shared owner)
  = Cross-governance evidence taxonomy, if Phase 1 proves reusable

sd-validation
  = How evidence is acquired and evaluated

runtime-governance
  = How violations affect execution
```

Current scope: workflow governance only. Runtime signals are advisory projections until a separate promotion decision names an executor, consumer, evidence threshold, and registry transition path.

Phase 1 owner boundary:

- `sd-ui-governance` owns local UI compliance classification: domain, render context, collection method, validation mechanism, evidence class, severity, and policy scope.
- `sd-validation` owns execution: acquiring evidence, running evaluation, and recording results.
- Future shared `validation-reasoning` may own the cross-governance taxonomy if the same model proves useful for Architecture, Runtime, Documentation, or Security governance.

## Governance Domains

Governance domains answer what is being governed. Do not replace these with tool names or validation mechanisms.

| Domain | Governs | Typical blocking candidates |
|---|---|---|
| Contract | Required UI states, screen actions, labels, traceability, consumer needs, view model projection | Missing loading/error/success state, missing label, untraceable screen action, unstable ViewModel derivation |
| Design System | Tokenized style usage, approved primitives, component variants, layout constraints | Raw color/spacing/font when project policy requires tokens, custom component bypassing required primitive |
| Accessibility | Keyboard, focus, semantics, labels, assistive feedback, contrast/motion expectations | Missing label, broken focus flow, missing role, objective contrast violation |
| Responsive | Whether the same UI contract remains valid across declared render contexts | Mobile-only overflow, modal outside viewport, CTA hidden by keyboard, sticky header covering content, sidebar cannot collapse, touch target too small, safe-area or orientation failure |
| Behavior | Interaction state transitions, validation feedback, destructive confirmation, retry/end-of-list/offline handling | Submit without loading state, delete without confirmation, retry path missing, permission denied state missing |
| Closure | Whether UI governance evidence is linked into review, DoD, tests, or deferred scope | Completion claim without evidence class, unresolved UI governance blocker hidden as done |

Observation: `Consumer / ViewModel` may later collapse under Contract as subdomains rather than remain a parallel governance domain:

```text
Contract
  ├─ Screen Contract
  ├─ Consumer Contract
  └─ ViewModel Contract
```

Keep this open until scenarios show whether separate handling improves routing or only adds taxonomy cost.

## Render Contexts

`render_context` answers where the same UI contract must continue to hold. It is a presentation / experience context, not a device inventory.

```yaml
supported_render_contexts:
  - desktop
  - tablet
  - mobile
  - narrow_mobile
```

Use context names before tool or hardware names. A validation suite may implement `mobile` with a browser viewport, a device emulation profile, a simulator, or a real device, but the governance finding should remain `domain: Responsive` + `render_context: mobile` rather than `device_testing: iPhone`.

Longer term, `render_context` may become one family inside a typed context taxonomy, alongside interaction, accessibility, environment, appearance, or locale contexts. Do not promote that broader taxonomy until scenarios show it is needed outside responsive UI validation.

Responsive domain downgrade watch:

- Keep `Responsive` as a workflow-local governance domain while it helps route viewport-specific failures clearly.
- Consider downgrading `Responsive` into a cross-cutting `context.render` dimension when most responsive failures are better explained as another governance domain evaluated under a render context, such as `domain: Accessibility` + `render_context: mobile`.
- Do not introduce a broad context taxonomy until multiple context families and workflow domains consume it.

## Collection Methods

Collection methods answer how evidence is acquired before it is evaluated. A collection method is not a governance domain, not a validation mechanism, and not proof of compliance by itself.

```yaml
supported_collection_methods:
  - contract_readback
  - static_analysis
  - runtime_trace
  - browser_review
  - human_observation
```

| Collection method | Acquires | Typical consumers |
|---|---|---|
| `contract_readback` | Screen / Consumer / ViewModel / Accessibility contract rows, generated surface readback, declared policy | deterministic mechanism, Closure domain |
| `static_analysis` | Code/config/lint/build assertions, token usage, component primitive usage | deterministic mechanism |
| `runtime_trace` | Component state, route behavior, UI event path, interaction trace | deterministic or manual review mechanisms |
| `browser_review` | Screenshot, DOM snapshot, accessibility tree, interaction trace, responsive capture, observed interaction state | screenshot_diff, ai_review, deterministic accessibility, responsive, or behavior validation |
| `human_observation` | Designer / reviewer observation, UX heuristic note, manual behavior evidence | manual_review mechanism, human_review evidence |

Browser Review invariants:

- Browser Review is not a governance domain.
- Browser Review is not a validation mechanism.
- Browser Review is an evidence collection method.
- A collection method may support multiple mechanisms.
- A mechanism may consume evidence from multiple collection methods.
- Evidence acquisition must not imply validation success by itself.

Visual reference review:

- When a UI change cites a screenshot, mock, product reference, or design artifact, evaluate layout by relative proportions and visual anchors before applying fixed pixel offsets. Compare edge insets, element-to-surface ratios, alignment to safe areas, and relationship to nearby content rather than asking only whether the element is visible.
- For mobile or constrained surfaces, verify proportional decisions with real viewport metrics such as DOM bounding boxes, computed styles, or device emulation. Headless screenshots can be misleading when the browser viewport, screenshot crop, or device scale does not match the target surface.
- Do not encode project-specific pixel values as reusable guidance. Preserve the transferable rule: derive spacing from the reference's proportions, then validate the chosen implementation against the target viewport and interaction constraints.

Future refinement: `collection_method` may later split from `collection_channel` if browser / human / telemetry channels need separate treatment. Do not add that fifth taxonomy layer in Phase 1.

## Validation Mechanisms

Validation mechanisms answer how a governance domain is evaluated. A mechanism can support multiple domains and must not own a domain.

```yaml
supported_mechanisms:
  - deterministic
  - screenshot_diff
  - ai_review
  - manual_review
```

| Mechanism | Use when | Boundary |
|---|---|---|
| deterministic | Rule can be checked from contract, code, config, lint, fixture, runtime trace, accessibility tree, or scan output | Can become block candidate when project policy and evidence are objective |
| screenshot_diff | UI appearance has deterministic capture and a meaningful baseline | Supports Design System or Behavior; not a governance domain or collection method |
| ai_review | Visual hierarchy, spacing consistency, CTA prominence, or heuristic feedback needs model review | Advisory by default; requires scoped rubric before it can affect release decision |
| manual_review | Human design / UX judgment is required | Useful evaluation, but not mechanical unless mapped to objective checklist items |

Two invariants must hold:

- No governance domain is represented solely by a validation mechanism.
- No validation mechanism is owned by a single governance domain.
- No collection method is treated as a validation mechanism.
- No collection method is treated as proof quality by itself.

## Severity Policy

```yaml
subjective_judgement:
  default_severity: warn

objective_judgement:
  default_severity: block_candidate
```

Do not make "AI thinks the UI looks bad" equivalent to "build failed". Missing loading state, missing label, missing destructive confirmation, or raw design-token violation can become blocking candidates. Vague visual taste remains warning / research until translated into objective criteria.

## Evidence Classes

Evidence classes distinguish source trust. They prevent an AI review, a Percy diff, an axe violation, and a contract row from collapsing into one undifferentiated UI finding.

Phase 1 warning: this list is provisional. It currently mixes evidence bodies (`contract`, `runtime`), acquisition artifacts (`screenshot`), evaluation outputs (`accessibility_scan`, `visual_diff`), and review outputs (`ai_review`, `human_review`). Do not use `evidence_class` as a durable finding taxonomy until a future validation-reasoning model splits Artifact / Evaluation / Finding.

| Evidence class | Source examples | Use |
|---|---|---|
| `contract` | Screen Contract, Consumer Contract, ViewModel Contract | Defines expected behavior; not execution proof by itself |
| `runtime` | component state, route behavior, UI event path | Proves behavior exists in implementation |
| `accessibility_scan` | axe, Lighthouse, Pa11y | Strong objective evidence for accessibility gates |
| `visual_diff` | Percy, Applitools, Chromatic, screenshot golden diff | Strong when baseline and capture are deterministic |
| `screenshot` | Playwright screenshot, manual capture | Useful input; weaker unless paired with baseline or rubric |
| `responsive_capture` | Per-context screenshot, DOM snapshot, viewport metrics, layout bounds, safe-area/focus/keyboard evidence | Strong when at least two required render contexts are captured and evaluated |
| `ai_review` | vision model review output | Advisory by default; requires scoped criteria |
| `human_review` | design review, UX heuristic review | Useful for subjective or strategic judgment; not mechanical by default |

This taxonomy is currently UI-scoped. Promote it to shared validation reasoning only after reuse across multiple governance domains, such as API governance, architecture governance, or runtime governance.

## Project-local Design System Policy

Ai-skill governs that design-system compliance is explicit, not the global shape of a design system.

Allowed global governance:

- Require tokenized design when a project declares token policy.
- Forbid raw style escape hatches when project policy marks them invalid.
- Require component primitive usage when project policy declares approved primitives.
- Require deferred scope when no token policy exists.

Not global governance:

- Hard-code spacing scales such as `[4, 8, 12, 16]`.
- Require a specific color token taxonomy.
- Assume Material, Carbon, Fluent, Tailwind, Ant Design, or any framework-specific primitive set.

## Minimum UI Governance Review

Before claiming UI compliance, classify:

| Field | Required answer |
|---|---|
| Domain | Contract, Design System, Accessibility, Responsive, Behavior, Closure |
| Render context | desktop, tablet, mobile, narrow_mobile, not_applicable, or project-declared equivalent |
| Collection method | contract_readback, static_analysis, runtime_trace, browser_review, human_observation |
| Mechanism | deterministic, screenshot_diff, ai_review, manual_review |
| Evidence class | contract, runtime, accessibility_scan, visual_diff, screenshot, responsive_capture, ai_review, human_review |
| Severity | block_candidate, warn, research, not_applicable |
| Project policy | token policy, component primitive policy, accessibility target, screenshot baseline, or not configured |
| Deferred scope | explicit owner and reason when evidence is missing |

## Non-goals

- Do not define a global design-token scale.
- Do not add a runtime/enforcement rule_class in the first landing.
- Do not treat responsive validation as a device SKU checklist. Govern the render context first; map it to tooling only during evidence acquisition.
- Do not treat visual regression or AI review as governance domains.
- Do not treat Browser Review as a governance domain, validation mechanism, or standalone validator.
- Do not make AI visual review blocking without objective rubric, deterministic capture, and project opt-in.
- Do not treat acquired evidence as validated finding without an evaluation mechanism and scoped evidence class.
- Do not duplicate `sd-ui-contracts`; contracts remain the source for expected UI behavior.

