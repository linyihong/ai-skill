# Decision Promotion Pipeline

This document defines how a session/runtime decision is promoted. The endpoint is not always an ADR.

## Core Rule

Promote decisions by content type, not by a fixed ladder.

| Decision content | Promotion target |
| --- | --- |
| Executable policy or cross-agent rule | `enforcement/` |
| Reasoning heuristic, tradeoff, signal, anti-pattern, or failure judgment | `intelligence/` |
| Operational procedure or repeatable work sequence | `workflow/` |
| Runtime gate, activation, phase, obligation, policy surface, or executable contract projection | `runtime/runtime.db` |
| Architecture-level irreversible or foundational decision | `constitution/ADR-*` |
| Session-scoped decision with future replay value | `memory/decision/` |
| Project-specific decision | `<PROJECT_ROOT>/docs/decisions/` |

## Promotion Flow

```text
runtime/session decision
  -> classify content type
  -> validate repeated or durable value
  -> choose target by content
  -> update linked surfaces
  -> write feedback lesson when reusable
  -> refresh runtime projection when execution-affecting
```

## ADR Boundary

Use `constitution/ADR-*` only when the decision is:

- foundational to the Ai-skill architecture,
- cross-session and cross-project,
- expected to remain stable,
- expensive to reverse,
- needed to explain why the system is shaped this way.

Do not promote every repeated decision to ADR. Many decisions are better represented as enforcement rules, intelligence atoms, workflow contracts, or runtime gates.

## No-Proposed-ADR Rule（2026-05-22 確立）

`constitution/` 只放 **accepted** ADRs。架構決策的提案階段 **不**寫 proposed ADR，改在對應 `plans/active/<plan>.md` 的 `Decision Rationale` section 完成提案、討論、alternatives 評估與 open questions。

### 規則

| 階段 | 位置 |
|------|------|
| 提案（problem + decision + alternatives + open questions） | `plans/active/<plan>.md` §Decision Rationale |
| 實作（phase 0-N） | 同一 plan |
| 完成 | plan 通過 §ADR Promotion Criteria |
| 升級為 ADR（直接 accepted） | `constitution/ADR-<n>-<slug>.md`，引用 completed plan 為 evidence |

### 為什麼

- 「Proposed ADR」是「未驗證的憲法」，與 constitution 的定位矛盾
- 失敗的 proposed ADR 留下「廢棄憲法」累積成噪音
- ADR-007 §Decision 已明文「ADR is NOT the default endpoint」— proposed ADR 違反此原則
- Plan 已是提案的天然容器，不需要平行 ADR 文件
- `Decision Rationale` section 強制 alternatives + open questions，保留 ADR 的 forcing function

### 失誤模式

未遵守此規則的失誤已沉澱為 [`enforcement/failure-patterns/premature-adr-promotion.md`](../../enforcement/failure-patterns/premature-adr-promotion.md)。

### 過渡

- 既有 accepted ADRs（ADR-001 ~ ADR-007）不受影響，仍為 canonical
- 既有 proposed ADRs（若有）需在 2026-05-22 後撤回並轉入 plan
- ADR 編號保留為流水號；升級為 accepted 時取下一個可用編號（不預留給特定 plan）

## Runtime Rule

If the target affects agent execution, it must either:

- have an executable YAML contract with `runtime_projection.enabled: true`, or
- update a canonical runtime document inside `runtime/runtime.db`.

## Related

- [`decision-promotion-pipeline.yaml`](decision-promotion-pipeline.yaml)
- [`executable-contract-boundary.md`](executable-contract-boundary.md)
- [`../../constitution/README.md`](../../constitution/README.md)
- [`../../memory/decision/README.md`](../../memory/decision/README.md)
