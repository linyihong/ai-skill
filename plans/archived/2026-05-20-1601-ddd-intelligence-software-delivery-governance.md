# DDD Integration Plan for Intelligence & Software Delivery Governance

> **狀態**: completed / archived
> **建立日期**: 2026-05-20
> **目的**: 將 Domain-Driven Design（DDD）納入 intelligence / software-delivery 知識體系，但保持其作為「architecture strategy」而非 runtime invariant，並與現有 Cognitive State & Evidence Governance、Runtime Recovery、Workflow Governance 系統建立清楚邊界與協作關係。

---

## 0. Architecture Compatibility Preflight

| 欄位 | 內容 |
| --- | --- |
| Trigger | 使用者要求建立 DDD integration plan，並檢查是否與現有框架衝突。 |
| Checked sources | `plans/README.md`、`plans/active/2026-05-20-1501-cognitive-state-evidence-governance.md`、`workflow/software-delivery/README.md`、`governance/ai-runtime-governance/software-delivery-governance.md`、`intelligence/engineering/architecture/README.md`、`intelligence/engineering/domain/README.md`、`runtime/runtime.db`、`metadata/README.md`、`knowledge/graphs/intelligence-domain.yaml`。 |
| Conflicts | 原 proposal 的 `software-delivery/architecture/` root path 與現行架構不一致；應改為 `workflow/software-delivery/architecture/` 或 `governance/ai-runtime-governance/software-delivery-architecture-governance.md`。DDD tactical modeling 已有現行入口 `intelligence/engineering/domain/`，不宜全部放到 `intelligence/engineering/architecture/`。DDD 不應直接寫入 `runtime/compiler/embedded_data.rb` 或 `runtime.db`，除非後續明確 promotion 成 runtime-lite signal。 |
| Decision | Proceed with revised plan. DDD 作為 selectable architecture strategy 納入 intelligence / workflow / governance / metadata / validation；不進 runtime primitive。 |
| Validation | 本 plan 完成後需做 diff review、Markdown link check、ReadLints；若後續實作 metadata / routing / validation source，再跑 `ruby scripts/refresh-knowledge-runtime.rb`。 |

### 0.1 Boundary Corrections

| Proposal | Current framework fit | Adjustment |
| --- | --- | --- |
| `intelligence/engineering/architecture/domain-driven-design/` | Partially compatible. `architecture/` 放架構思考模式，但現有 `intelligence/engineering/domain/` 已明確負責 DDD / 業務模型智慧。 | DDD tactical modeling 放 `intelligence/engineering/domain/domain-driven-design/`；architecture selection / architecture fit 放 `intelligence/engineering/architecture/architecture-selection/`。 |
| `software-delivery/architecture/` | Conflict. Repo 沒有 root `software-delivery/`。 | 使用 `workflow/software-delivery/architecture/` 保存執行流程；治理條文放 `governance/ai-runtime-governance/software-delivery-architecture-governance.md`。 |
| `governance/software-delivery-governance/` | Conflict. 現行位置是 `governance/ai-runtime-governance/software-delivery-governance.md`。 | 若新增 architecture governance，放同層 `governance/ai-runtime-governance/software-delivery-architecture-governance.md`，並連回既有 software delivery governance。 |
| `metadata/architecture/` | Compatible as new metadata namespace. | 明確標為 metadata-only；不可宣稱 runtime enforced。 |
| `validation/scenarios/architecture/` | Compatible as new validation namespace. | 建立後需更新 validation README / knowledge runtime refresh。 |
| DDD runtime-lite signals | Compatible only if compressed. | 只允許 `architecture_reliability_degradation`、`boundary_mismatch_recovery_loop` 這類 compressed signal；禁止 full DDD runtime enforcement。 |

---

## 1. Problem Statement

目前系統已具備：

- runtime governance
- recovery orchestration
- evidence qualification
- autonomy downgrade
- execution stabilization
- workflow governance
- escalation handling

但在 software delivery / architecture decision 層，尚未建立完整的：

- architecture selection governance
- business complexity evaluation
- bounded context heuristics
- overengineering detection
- domain decomposition policy
- long-lived system modeling strategy

導致 agent 可能：

- 對所有專案強制套用 DDD。
- 在低 complexity 專案過度抽象化。
- 將 CQRS / aggregate / event sourcing 當作 default architecture。
- 缺乏 architecture fit analysis。
- 無法區分 runtime governance problem、domain modeling problem、delivery complexity problem。

本計畫目標是建立 DDD 作為「可選 architecture strategy」的治理與知識層，而不是 mandatory workflow。

---

## 2. Core Principle

DDD 的定位：

```text
DDD = architecture strategy
不是 runtime governance primitive
不是 universal workflow invariant
```

DDD 的責任：

- business domain modeling
- bounded context management
- ubiquitous language
- invariant protection
- long-lived complexity management

DDD 不負責：

- runtime recovery
- evidence governance
- autonomy control
- contradiction propagation
- execution stabilization

---

## 3. Layer Responsibility

| Concern | 建議位置 | 不該做的事 |
| --- | --- | --- |
| DDD philosophy / modeling theory | `intelligence/engineering/domain/domain-driven-design/` | 不寫 runtime MUST rule。 |
| Architecture selection heuristics | `intelligence/engineering/architecture/architecture-selection/` | 不把 DDD 當 default。 |
| Architecture selection governance | `governance/ai-runtime-governance/software-delivery-architecture-governance.md` | 不直接綁死 DDD。 |
| Software delivery execution steps | `workflow/software-delivery/architecture/` | 不取代 `workflow/software-delivery/execution-flow.md` 的主要交付流程。 |
| Complexity evaluation | `workflow/software-delivery/architecture/architecture-fit-analysis.md` 或 `metadata/architecture/` | 不直接變 runtime gate。 |
| Overengineering detection | `governance/ai-runtime-governance/software-delivery-architecture-governance.md` | 不阻塞低風險 delivery；只要求 simplification review。 |
| Runtime recovery | `runtime/` / `governance/ai-runtime-governance/recovery-retry-governance.md` | 不理解 bounded context，不處理 aggregate modeling。 |
| Cognitive governance | `governance/ai-runtime-governance/cognitive-state-governance.md`（planned） | 不處理 aggregate modeling，只處理 belief/evidence/claim scope。 |
| Domain heuristics | `metadata/architecture/` | metadata-only；不得假裝 runtime enforced。 |
| Validation scenarios | `validation/scenarios/architecture/` | 測 routing / decision behavior，不測私有 project architecture。 |

---

## 4. Proposed Directory Structure

### 4.1 Intelligence Layer

```text
intelligence/
└── engineering/
    ├── domain/
    │   └── domain-driven-design/
    │       ├── README.md
    │       ├── bounded-context.md
    │       ├── aggregate-design.md
    │       ├── ubiquitous-language.md
    │       ├── domain-events.md
    │       ├── anti-corruption-layer.md
    │       ├── domain-services.md
    │       ├── repository-pattern.md
    │       ├── tactical-vs-strategic-design.md
    │       ├── ddd-lite-vs-full-ddd.md
    │       ├── event-storming.md
    │       ├── overengineering-risks.md
    │       └── architecture-fit-signals.md
    │
    ├── architecture/
    │   └── architecture-selection/
    │       ├── README.md
    │       ├── architecture-fit-analysis.md
    │       ├── complexity-evaluation.md
    │       ├── lifecycle-evaluation.md
    │       ├── team-topology-analysis.md
    │       ├── integration-pressure-analysis.md
    │       └── architecture-tradeoff-matrix.md
    │
    └── anti-patterns/
        ├── cargo-cult-ddd.md
        ├── premature-cqrs.md
        ├── aggregate-explosion.md
        ├── repository-abstraction-overuse.md
        └── architecture-absolutism.md
```

### 4.2 Software Delivery Governance / Workflow

```text
governance/
└── ai-runtime-governance/
    └── software-delivery-architecture-governance.md

workflow/
└── software-delivery/
    └── architecture/
        ├── README.md
        ├── architecture-selection-governance.md
        ├── architecture-decision-framework.md
        ├── bounded-context-evaluation.md
        ├── architecture-escalation-policy.md
        ├── architecture-minimality-principle.md
        └── architecture-overengineering-detection.md
```

### 4.3 Metadata and Validation

```text
metadata/
└── architecture/
    ├── README.md
    ├── architecture-fit-matrix.yaml
    ├── ddd-adoption-signals.yaml
    ├── overengineering-signals.yaml
    └── bounded-context-heuristics.yaml

validation/
└── scenarios/
    └── architecture/
        ├── cargo-cult-ddd.yaml
        ├── architecture-fit-mismatch.yaml
        ├── overengineering-detection.yaml
        ├── bounded-context-collapse.yaml
        └── aggregate-explosion.yaml
```

---

## 5. Architecture Selection Governance

Agent 不應預設：

```text
project => full DDD
```

而應先做 architecture fit analysis。

### 5.1 Evaluation Axes

| Axis | Meaning |
| --- | --- |
| `domain_complexity` | business rules complexity |
| `invariant_density` | critical business invariants |
| `business_language_instability` | vocabulary churn |
| `workflow_complexity` | orchestration complexity |
| `integration_pressure` | external system interaction |
| `lifecycle_length` | expected system lifespan |
| `team_scale` | number of contributors |
| `bounded_context_count` | number of subdomains |
| `event_coordination_need` | async/event-driven pressure |
| `delivery_speed_priority` | MVP vs long-term maintainability |

### 5.2 Suggested Routing

| Complexity | Suitable strategies | Avoid by default |
| --- | --- | --- |
| Low | CRUD, Vertical Slice, Feature Modules, Simple Service Layer | full aggregate design, CQRS everywhere, event sourcing |
| Medium | DDD Lite, explicit domain services, selective aggregate, modular monolith | microservices by default, event sourcing without audit/business need |
| High | Full DDD, bounded contexts, anti-corruption layers, event-driven coordination, aggregate consistency boundaries | single global model, CRUD-only model over high-invariant domain |

---

## 6. DDD Lite Governance

DDD Lite exists to avoid:

```text
small project → full enterprise architecture
```

DDD Lite does not require:

- CQRS
- event sourcing
- repository abstraction
- aggregate everywhere
- microservice split

DDD Lite keeps:

- ubiquitous language
- domain boundary awareness
- critical invariant protection
- explicit domain modeling
- minimal architecture decision record

---

## 7. Overengineering Governance

### 7.1 Purpose

Prevent AI architecture inflation.

### 7.2 Overengineering Signals

- aggregate count rapidly grows.
- CQRS added without scale / read-write divergence requirement.
- event sourcing without audit, temporal query, replay, or business event need.
- repository abstraction with single implementation and no test seam value.
- microservice split without deployment, scaling, ownership, or compliance boundary.
- excessive domain event chaining.
- abstraction count exceeds domain complexity.

### 7.3 Required Action

- architecture simplification review.
- bounded-context merge evaluation.
- remove speculative abstractions.
- return to minimal viable architecture.
- record why higher architecture complexity is justified if retained.

---

## 8. Relationship with Cognitive Governance

| Cognitive Governance | DDD |
| --- | --- |
| evidence qualification | domain modeling |
| recovery stabilization | business decomposition |
| autonomy downgrade | aggregate boundary selection |
| contradiction propagation | domain invariant review |
| claim scope governance | bounded context scope |
| intent stability | workflow / business intent |
| runtime safety | architecture maintainability |

### Shared Concepts

| Shared concept | DDD meaning | Cognitive / runtime meaning | Boundary |
| --- | --- | --- | --- |
| Bounded Context / Claim Scope | local domain language and model boundary | local claim must not become global conclusion | Shared analogy only; not same enforcement mechanism. |
| Aggregate Invariant / Runtime Safety Gate | business consistency boundary | execution invariant / phase gate | Similar consistency idea; implementation belongs to different layers. |
| Anti-Corruption Layer / Evidence Qualification | protect domain model from external model contamination | protect execution belief from low-quality evidence contamination | Shared contamination metaphor; DDD does not govern evidence authority. |

---

## 9. Runtime Boundary

DDD must not be promoted directly into:

- runtime primitive
- every-task enforcement
- `runtime.db` cognition state
- phase transition invariant
- blocking gate for low-risk delivery

Reason:

```text
DDD = delivery architecture strategy
not execution safety primitive
```

### 9.1 Concepts That Must Not Become Runtime Invariants

| Concept | Reason |
| --- | --- |
| bounded context purity | architecture concern |
| aggregate modeling style | implementation strategy |
| ubiquitous language completeness | delivery optimization |
| event storming workflow | planning methodology |
| tactical pattern preference | non-runtime concern |

### 9.2 Runtime-Relevant DDD Signals

Only a few compressed signals may become runtime-lite candidates:

| Signal | Reason |
| --- | --- |
| cross-context invariant violation | may affect execution correctness |
| architecture drift causing delivery instability | impacts execution reliability |
| boundary mismatch causing recovery loops | affects runtime governance |

Even then:

```text
runtime-lite only
not full DDD runtime enforcement
```

This aligns with the Minimal Runtime Principle in the cognitive governance plan.

---

## 10. Metadata

Suggested metadata files:

```text
metadata/architecture/
├── README.md
├── architecture-fit-matrix.yaml
├── ddd-adoption-signals.yaml
├── overengineering-signals.yaml
└── bounded-context-heuristics.yaml
```

Boundary:

- metadata-only unless explicitly promoted.
- metadata describes fit signals, not mandatory architecture.
- compiler / runtime integration must be a separate plan if any signal is promoted.

---

## 11. Validation Scenarios

Candidate scenarios:

| Scenario | Tests |
| --- | --- |
| `cargo-cult-ddd.yaml` | 小專案被錯誤升級 full DDD。 |
| `architecture-fit-mismatch.yaml` | 高 complexity 專案錯誤使用 CRUD-only model。 |
| `overengineering-detection.yaml` | abstraction growth 超過 domain complexity。 |
| `bounded-context-collapse.yaml` | multiple domains 被錯誤混成單一 model。 |
| `aggregate-explosion.yaml` | aggregate boundaries 過度切分。 |

Validation should check agent routing and decision behavior, not project-specific architecture.

---

## 12. Open Questions

| Question | Current recommendation |
| --- | --- |
| DDD 是否預設啟用？ | No. 先做 architecture fit analysis。 |
| 是否允許 DDD Lite？ | Yes. 大部分中型系統更適合。 |
| 是否把 DDD promotion 成 runtime primitive？ | Mostly no. 僅極少數 architecture instability signal 可 runtime-lite。 |
| 是否建立 architecture minimality principle？ | Yes. 避免 AI architecture inflation。 |
| 是否建立 architecture governance tier？ | Yes. 避免 architecture optimization 阻塞 delivery。 |
| DDD tactical docs 放 architecture 還是 domain？ | Tactical modeling 放 `intelligence/engineering/domain/`；architecture selection 放 `intelligence/engineering/architecture/`。 |

---

## 13. Suggested Implementation Phases

### Phase 0 — Boundary Definition

- Define DDD vs runtime governance.
- Define architecture vs execution.
- Define modeling vs recovery.
- Define architecture minimality principle.
- Update route / layer docs only if needed for discoverability.

### Phase 1 — Intelligence Layer

- Build DDD theory and tactical heuristics under `intelligence/engineering/domain/domain-driven-design/`.
- Build architecture selection heuristics under `intelligence/engineering/architecture/architecture-selection/`.
- Add anti-patterns: cargo-cult DDD, premature CQRS, aggregate explosion, repository abstraction overuse, architecture absolutism.
- Link existing `intelligence/engineering/domain/aggregate-boundary-heuristics.md`.

### Phase 2 — Software Delivery Governance

- Add architecture fit analysis to software delivery workflow.
- Add architecture routing and overengineering governance.
- Add architecture escalation policy for fit mismatch, not for every project.
- Keep software delivery flow minimal for low-risk changes.

### Phase 3 — Metadata & Heuristics

- Add architecture-fit metadata.
- Add DDD adoption signals.
- Add complexity evaluation heuristics.
- Keep metadata as metadata-only.

### Phase 4 — Validation Scenarios

- Add cargo cult DDD detection.
- Add architecture mismatch scenario.
- Add bounded-context contamination scenario.
- Add abstraction explosion scenario.
- Run knowledge runtime refresh after validation sources are added.

---

## 14. Completion Definition

完成後系統應做到：

- 不會對所有專案強制 DDD。
- 能做 architecture fit analysis。
- 能區分 CRUD / DDD Lite / Full DDD。
- 能檢測 architecture overengineering。
- 能把 bounded context 當成 architecture / domain modeling concern。
- 不會把 DDD promotion 成 runtime invariant。
- 能避免 AI architecture absolutism。
- 能讓 architecture complexity 與 business complexity 對齊。
- 能維持 minimal viable architecture。
- 能讓 runtime governance 與 architecture governance 解耦。
- 能讓 DDD 成為 selectable strategy，而非 universal doctrine。

---

## 15. Current Compatibility Summary

No blocking conflict found after boundary correction.

Non-blocking issues to respect during implementation:

- `runtime/runtime.db` already has `route.intelligence.domain`; DDD expansion should extend or route through this rather than create a competing route without registry review.
- `intelligence/engineering/domain/aggregate-boundary-heuristics.md` already exists; new DDD docs should reference it instead of duplicating aggregate boundary guidance.
- `workflow/software-delivery/README.md` already has Simplicity First; DDD governance must reinforce, not weaken, this principle.
- `governance/ai-runtime-governance/software-delivery-governance.md` already defines delivery runtime gates; architecture governance must not add mandatory DDD gates for low-risk delivery.
- `plans/active/2026-05-20-1501-cognitive-state-evidence-governance.md` already owns cognitive evidence / confidence / claim scope; DDD should only share analogies, not take over cognitive governance.

---

## 16. Closure Reconciliation

執行日期：2026-05-20

完成狀態：

- Phase 0 Boundary Definition：完成。DDD 被明確定位為 selectable architecture strategy，不是 runtime primitive；原 proposal 的 `software-delivery/architecture/` root path 已修正為 `workflow/software-delivery/architecture/` 與 `governance/ai-runtime-governance/software-delivery-architecture-governance.md`。
- Phase 1 Intelligence Layer：完成。新增 `intelligence/engineering/domain/domain-driven-design/`、`intelligence/engineering/architecture/architecture-selection/` 與 DDD / architecture anti-patterns。
- Phase 2 Software Delivery Governance：完成。新增 software-delivery architecture workflow 與 governance gate，並接入 `workflow/software-delivery/README.md` / `execution-flow.md`。
- Phase 3 Metadata & Heuristics：完成。新增 `metadata/architecture/`，並標記為 metadata-only。
- Phase 4 Validation Scenarios：完成。新增 5 個 architecture scenarios，並更新 routing registry / graphs / summaries / generated runtime surfaces。

Validation：

- `ruby scripts/refresh-knowledge-runtime.rb` 通過。
- Knowledge runtime report、model context report、model checklists 與 SQLite runtime index 已重新生成。
- DDD 未 promotion 成 `runtime.db` runtime invariant；只有 metadata / validation / routing surface 被更新。
