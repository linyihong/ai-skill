# BDD + DDD Cognition-Aligned Reframe

> **狀態**: completed / archived
> **建立日期**: 2026-05-20
> **目的**: 將 BDD / DDD integration 從 methodology-first 重新整理為 cognition-first software delivery operating model：BDD 歸入 requirements cognition，DDD 歸入 domain architecture cognition，workflow 只保留 delivery stages，runtime 只接收壓縮後的 minimum viable delivery safety signals。

---

## 0. Architecture Compatibility Preflight

| 欄位 | 內容 |
| --- | --- |
| Trigger | 使用者要求執行 attached `bdd-ddd-reframe` plan，且明確要求不要修改 `.cursor/plans/` plan file。 |
| Checked sources | `plans/README.md`、`enforcement/dependency-reading.md`、`enforcement/linked-updates.md`、`workflow/software-delivery/README.md`、`workflow/software-delivery/execution-flow.md`、`workflow/software-delivery/development-process.md`、`governance/ai-runtime-governance/software-delivery-governance.md`、`governance/ai-runtime-governance/software-delivery-architecture-governance.md`、`intelligence/engineering/development/README.md`、`intelligence/engineering/architecture/README.md`、`intelligence/engineering/architecture/domain-modeling/README.md`、`metadata/architecture/README.md`、`validation/README.md`、`knowledge/runtime/routing-registry.yaml`。 |
| Conflicts | attached plan 的方向與現行框架相容；需修正既有 DDD integration 的 methodology-first path。不得新增 root `software-delivery/`；canonical workflow 仍是 `workflow/software-delivery/`。不得把 BDD / DDD syntax 或 tactical pattern 寫入 runtime primitive。 |
| Decision | Proceed. 建立 repo-local active plan，執行 cognition-first reframe；不修改 `.cursor/plans/bdd-ddd-reframe_6ea009a6.plan.md`。 |
| Validation | `ruby scripts/refresh-knowledge-runtime.rb`、touched Markdown link check、ReadLints、diff review、commit/push/readback、clean status。 |

---

## 1. Design Decision

採用 cognition-first 分層：

- BDD = requirements cognition：治理 ambiguity、intent mismatch、acceptance drift、behavior scope、requirement traceability。
- DDD = domain architecture cognition：治理 domain complexity、ubiquitous language、bounded context、invariant、consistency boundary、architecture fit。
- Workflow = delivery stages：requirements → architecture → implementation → validation → release → recovery。
- Runtime = minimum viable safety signal：requirement contradiction、missing validation target、invariant violation、stale acceptance criteria。

---

## 2. Existing DDD Integration Mapping

| Current source | New location / decision |
| --- | --- |
| `intelligence/engineering/architecture/domain-modeling/` | Move into `intelligence/engineering/architecture/domain-modeling/` and `architectural-fit/` where appropriate. |
| `intelligence/engineering/architecture/architectural-fit/` | Rename / consolidate into `intelligence/engineering/architecture/architectural-fit/`. |
| `workflow/software-delivery/architecture/` | Keep as workflow stage, but reorganize around fit/discovery/design/escalation instead of DDD methodology bucket. |
| `metadata/architecture/` | Keep, extend with behavior governance metadata. |
| `validation/scenarios/architecture/` | Keep architecture scenarios, add software-delivery requirement/behavior scenarios. |

---

## 3. Target Structure

```text
intelligence/engineering/
  requirements/
    behavior-modeling/
    specification-quality/
    validation-thinking/
  architecture/
    domain-modeling/
    architectural-fit/
    delivery-alignment/
    system-boundaries/
    modularity/
    coupling-tradeoffs/
    consistency-models/
  execution/
    software-delivery-cognition/
    task-decomposition/
    validation-reasoning/
    recovery-thinking/
```

```text
workflow/software-delivery/
  requirements/
  architecture/
  implementation/
  validation/
  release/
  recovery/
```

---

## 4. Runtime Boundary

不建立 runtime BDD / DDD primitive。只允許後續另案評估下列 runtime-lite signal：

- `requirement_contradiction`
- `missing_validation_target`
- `invariant_violation`
- `stale_acceptance_criteria`

---

## 5. Completion Checklist

- [x] 新 active plan 建立並含 preflight。
- [x] DDD integration files moved / renamed / references updated。
- [x] Requirements cognition intelligence / workflow / governance 補齊。
- [x] Architecture cognition paths 補齊並更新 README。
- [x] Metadata、validation scenarios、routing registry、graphs、summaries 更新。
- [x] Knowledge runtime refresh 通過。
- [x] Plan completion closure：標 completed、archive、更新 `plans/README.md`。

---

## 6. Closure Reconciliation

本 plan 已執行完成：

- BDD 被落在 `intelligence/engineering/requirements/` 與 `workflow/software-delivery/requirements/`，作為 requirements cognition / BDD-lite，而不是 Gherkin everywhere。
- DDD tactical modeling 已從 `intelligence/engineering/domain/domain-driven-design/` 重新歸位到 `intelligence/engineering/architecture/domain-modeling/`，architecture selection 已收斂為 `intelligence/engineering/architecture/architectural-fit/`。
- Workflow 保留 delivery stages：requirements、architecture、implementation、validation、release、recovery。
- Runtime 邊界維持 metadata-only / validation-scenario surface；未新增 BDD / DDD runtime primitive。
- 已更新 governance、metadata、validation scenarios、knowledge graphs、summaries、routing registry 與 generated runtime reports。
- 驗證通過：`ruby scripts/refresh-knowledge-runtime.rb`、touched Markdown link check、ReadLints。
