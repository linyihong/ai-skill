# UI Governance Slice

> **Cognitive Slice**：`sd-ui-governance`（從 UI contract surface 延伸出的 workflow governance slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §4 的 workflow membership test）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-ui-governance` |
| `purpose` | 在 UI / consumer contract 已存在或正在建立時，定義 UI compliance 如何被評估：governance domain、validation mechanism、evidence class、severity policy 與 runtime projection boundary |
| `type` | `execution` |
| `tags` | artifact-gate, ui, design-system, accessibility, behavior-governance, visual-validation |
| `load_when` | UI / consumer surface 需要 design-system enforcement、accessibility evidence、behavior pattern checks、visual baseline review、AI visual review scoping，或 completion claim 依賴 UI compliance |
| `do_not_load_when` | 無 user-visible consumer surface、純 provider 內部變更、只修不影響 UI behavior / contract / design-system compliance 的小錯 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「UI compliance 評估前要分類哪些 domain、mechanism、evidence、severity，並保持 runtime projection advisory-only」的 ordering / gate；通過 workflow membership test，不承載 evidence 取得方法（非 analysis），不論證長期 pattern（非 intelligence） |
| `canonical_source` | 本檔 |
| `dependencies` | `sd-ui-contracts`（UI contract expectations）、`sd-test-strategy`（proof target selection）、`sd-validation`（evidence execution） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | UI governance scenarios can prove domain/mechanism separation, evidence classification, and AI visual warning boundary before runtime/enforcement promotion |

## Boundary

`sd-ui-contracts` defines what the UI should be. `sd-ui-governance` defines how UI compliance is evaluated.

```text
sd-ui-contracts
  = What UI should be

sd-ui-governance
  = How UI compliance is evaluated

sd-validation
  = How evidence is executed

runtime-governance
  = How violations affect execution
```

Current scope: workflow governance only. Runtime signals are advisory projections until a separate promotion decision names an executor, consumer, evidence threshold, and registry transition path.

## Governance Domains

Governance domains answer what is being governed. Do not replace these with tool names or validation mechanisms.

| Domain | Governs | Typical blocking candidates |
|---|---|---|
| Contract | Required UI states, screen actions, labels, traceability, consumer needs, view model projection | Missing loading/error/success state, missing label, untraceable screen action, unstable ViewModel derivation |
| Design System | Tokenized style usage, approved primitives, component variants, layout constraints | Raw color/spacing/font when project policy requires tokens, custom component bypassing required primitive |
| Accessibility | Keyboard, focus, semantics, labels, assistive feedback, contrast/motion expectations | Missing label, broken focus flow, missing role, objective contrast violation |
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
| deterministic | Rule can be checked from contract, code, config, lint, fixture, or scan output | Can become block candidate when project policy and evidence are objective |
| screenshot_diff | UI appearance has deterministic capture and a meaningful baseline | Supports Design System or Behavior; not a governance domain |
| ai_review | Visual hierarchy, spacing consistency, CTA prominence, or heuristic feedback needs model review | Advisory by default; requires scoped rubric before it can affect release decision |
| manual_review | Human design / UX judgment is required | Useful evidence, but not mechanical unless mapped to objective checklist items |

Two invariants must hold:

- No governance domain is represented solely by a validation mechanism.
- No validation mechanism is owned by a single governance domain.

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

| Evidence class | Source examples | Use |
|---|---|---|
| `contract` | Screen Contract, Consumer Contract, ViewModel Contract | Defines expected behavior; not execution proof by itself |
| `runtime` | component state, route behavior, UI event path | Proves behavior exists in implementation |
| `accessibility_scan` | axe, Lighthouse, Pa11y | Strong objective evidence for accessibility gates |
| `visual_diff` | Percy, Applitools, Chromatic, screenshot golden diff | Strong when baseline and capture are deterministic |
| `screenshot` | Playwright screenshot, manual capture | Useful input; weaker unless paired with baseline or rubric |
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
| Domain | Contract, Design System, Accessibility, Behavior, Closure |
| Mechanism | deterministic, screenshot_diff, ai_review, manual_review |
| Evidence class | contract, runtime, accessibility_scan, visual_diff, screenshot, ai_review, human_review |
| Severity | block_candidate, warn, research, not_applicable |
| Project policy | token policy, component primitive policy, accessibility target, screenshot baseline, or not configured |
| Deferred scope | explicit owner and reason when evidence is missing |

## Non-goals

- Do not define a global design-token scale.
- Do not add a runtime/enforcement rule_class in the first landing.
- Do not treat visual regression or AI review as governance domains.
- Do not make AI visual review blocking without objective rubric, deterministic capture, and project opt-in.
- Do not duplicate `sd-ui-contracts`; contracts remain the source for expected UI behavior.

