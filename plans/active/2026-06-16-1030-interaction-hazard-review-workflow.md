---
id: 2026-06-16-1030-interaction-hazard-review-workflow
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-16
priority: P1
required_for_completion: false
---

# State Trust Transition — Promotion Discipline (Decision Framework)

**Status**: `draft` — **promotion discipline / decision framework**（非 workflow 设计稿）
Owner: framework maintainer (linyihong)
**建立日期**：2026-06-16
**最後修訂**：2026-06-16（review #4：downstream pilot gate 勾选；Vidoe-Test 证据回写 A0–B.5 + C partial）
**Priority**：**P1**
**Downstream pilot**（canonical evidence consumer — 细节不复制进本 plan）：Vidoe-Test `docs/plans/2026-06-16-state-trust-transition-pilot.md`；commits `bcce737`（A0 overlay）、`6665b77`（A1/B/B.5 + C BDD partial）

## Executive summary

本 plan 已从「设计 workflow」转为 **升级判准（promotion discipline）**：

```text
primitive → evidence → consumer → 是否值得升 workflow
```

**不是**：问题 → 建 workflow。

**Primitive / consumer 分离**（本版最成熟处）：

- **Primitive**：State Trust Transition table（Ownership Map + Invalidation + Recovery）
- **Evidence**：`temporal_behavior` 子形状（event_trace / dom_presence / ownership）
- **Consumer**：software-delivery、UI overlay、integration claims — **仅为消费者**

**真正在验的题**（不是 O2/O3 标签本身）：

```text
同一 trust transition（invalidate → recover）能否形成最小闭环？
```

若 Dialog、rollback、websocket 都能用 **同一四栏、不改字** 表达 trust lost → trust restored ⇒ lean **O3**（generic model）。否则 lean **O2**（UI conditional gate）。

**Blocking decision**（O1 = future promotion only）：

| Path | 选 when | Pilot status (2026-06-16) |
|---|---|---|
| **O2 — Conditional gate** | B + B.5 失败：栏位不能跨域复用 | not selected |
| **O3 — Generic trust transition model** | B + B.5 通过：四栏不改字可套 rollback + websocket | **tentative lean** — see §Downstream pilot gate |

---

## Downstream pilot gate

Phase 勾选 **不以 Ai-skill 自证**；以 downstream validation pilot 产物 + commit 为准（[`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md)）。

| Phase | Gate pass when | Downstream evidence | Status |
|---|---|---|---|
| A0 | 四栏 + Recovery Boundary 定稿并 sync 到 consumer workflow | Vidoe-Test `bcce737` — `framework-development-workflow.yaml`, `interaction-hazard-review.md` | **pass** |
| A1 | Coupon 四栏 trust table | Vidoe-Test `6665b77` — `screen-mapping/episode-coupon-redeem-journey.md` | **pass** |
| B | 非 player invalidate↔recover 闭环 | Vidoe-Test `6665b77` — `screen-mapping/membership-payment-sync-trust-journey.md` | **pass** |
| B.5 | 四栏不改字压力测试 | Vidoe-Test `6665b77` — `screen-mapping/websocket-subscription-trust-journey.md` | **pass** |
| C | field survival scenario + predictive prevention | Vidoe-Test `6665b77` — `tests/bdd/state-trust-transition-pilot.test.mjs` (5/5); predictive **pending** | **partial** |
| D | 全部 ADR criteria | — | **blocked** (criteria 4, 6; integration journey pending) |

**Promotion gate summary**（ADR criteria 1–6）：

| # | Criterion | Status |
|---|---|---|
| 1 | ≥2 cases, four-column table | **pass** (coupon + payment sync) |
| 2 | ≥1 validation scenario consumes trust evidence | **partial** (BDD field survival; Ai-skill validation YAML ids pending) |
| 3 | O2 or O3 resolved | **tentative O3** (B + B.5 pass) |
| 4 | ≥1 previously unknown prevention | **pending** |
| 5 | template field survives renaming | **pass** (B.5 websocket, same headers) |
| 6 | no rubber-stamp | **ongoing** |

---

## Core primitive — State Trust Transition

### Ownership Map（四栏 — A0 定稿）

| State | Owner | Invalidation Event | Recovery Boundary |
|---|---|---|---|

**Invalidation Event** — any event after which the state must no longer be trusted.

**Recovery Boundary** — what evidence makes this state trustworthy again.

```text
trust lost  →  trust restored
```

| State | Owner | Invalidation Event | Recovery Boundary |
|---|---|---|---|
| `playbackAllowed` | Entitlement | refresh | grant readback |
| `optimisticBalance` | QueryCache | rollback | server sync |
| `dialogOpen` | CouponPanel | unmount | explicit reopen |
| `websocketReady` | Connection | reconnect | handshake complete |

**Why Recovery Boundary matters**：invalidated ≠ consumer must immediately stop — e.g. refresh 后是否可暂时信、等 readback？无 Recovery 栏，Phase B 只能描述问题，不能描述 **结束**。

Hazard class：`owner-invalidation-before-complete` when recovery boundary crossed while async work still in flight.

Ownership 是 trust transition 的 **一个来源**，不是全部 — 模型名倾向 **State Trust Transition**，Ownership Map 为表格名。

---

## Promotion discipline（ADR criteria）

Phase D graduation **全部**满足：

1. ≥ **2** independent cases，四栏 finalized template 填完
2. ≥ **1** validation scenario 机械消费 interception / trust evidence
3. **O2 or O3** resolved（B + B.5 证据）
4. ≥ **1 previously unknown prevention** — not post-hoc explanation of shipped bugs
5. ≥ **1 template field survives renaming pressure** — 若 Dialog → rollback → websocket 后 **Invalidation Event** 与 **Recovery Boundary** 列名仍存在、不需改栏位 ⇒ primitive 有生命力；若每案改栏 ⇒ abstraction noise，**勿 O3**
6. No sustained rubber-stamp on empty sections

```text
Post-hoc explanatory power ≠ predictive interception
Abstraction that renames every case ≠ primitive
```

---

## Evidence（consumer 层 — 不膨胀 taxonomy）

```yaml
temporal_behavior:
  event_trace:
  dom_presence:
  ownership:    # trust boundary / no-invalidate-before-complete
```

**新增 consumer ≠ 新增 taxonomy.**

Governance invariant（consumer）：`observable_outcome_must_survive_owner_refresh` — outcome trusted across refresh window until recovery boundary evidence arrives.

---

## Evidence Rule

> Machine-readable evidence-rule（schema `evidence-rule-v1`），索引於
> [`governance/evidence-candidates/evidence-rules/interaction-hazard.pointer.yaml`](../../governance/evidence-candidates/evidence-rules/interaction-hazard.pointer.yaml)。
> **Phase 1A Step 2（consumer attach）**：本 section 成立 = consumer hook 建立；criterion 內容是
> **Step 3（criteria authoring）**，下方刻意留 placeholder。rule 定義 owner = 本 plan。acceptance-gate
> 形狀候選 `pilot_complete + criteria_pass >= 6`；證據可跨 repo（下游 Vidoe-Test commit）。notify
> 屬 acceptance-gate（gate projection），不在 evidence_rule。設計來源見
> [`evidence-candidate-system`](2026-06-16-1131-evidence-candidate-system.md)。

```yaml
evidence_rule:
  collect: true
  match:
    artifact_types: []   # Step 3
    criteria: []         # Step 3（候選方向：四欄 trust table case / downstream pilot gate / field survival）
  exclusions: []         # Step 3
```

## Roadmap

```text
A0 — Template + Recovery Boundary（四栏定稿 + downstream sync）
A1 — Coupon trust transition table
B  — Optimistic rollback（invalidate ↔ recover 闭环）
B.5 — Rename pressure test（websocket 或第三域；不必新 incident）
C  — Scenario spike（predictive + field survival）
D  — IF ADR criteria → O2/O3 graduation（optional slice / glossary）
E  — Project overlay advisory
F  — Mechanical promotion（deferred）
```

### A0 — Template + Recovery Boundary

- [x] 四栏定义写入本 plan
- [x] Downstream project overlay / screen-mapping sync Recovery Boundary — Vidoe-Test `bcce737`
- [x] Side Effect Chain：`invalidation_events` + **`recovery_evidence`** per step — downstream workflow yaml + screen mappings

### A1 — Coupon

- [x] 四栏填 coupon — Vidoe-Test `episode-coupon-redeem-journey.md` (`6665b77`); aligns Appendix A
- [x] Counterfactual documented — coupon unmount hazard + recovery (post-ship fix recorded; predictive record pending criterion 4)

### B — Optimistic rollback / payment sync trust

- [x] Primary：invalidate + server sync recovery 落入四栏 — membership payment sync mapping (`6665b77`)
- [x] 验 **trust transition 闭环** — pending UI vs sync + `router.refresh` recovery boundary

### B.5 — Rename pressure test

- [x] **不用同名案例** — websocket subscription sketch (`6665b77`)
- [x] 四栏 **不改字** 套用 — BDD `template field survival` test passes
- [x] **B.5 pass ⇒ tentative O3**（非最终 Phase D graduation）

### C — Scenario spike

- [ ] ≥1 previously unknown prevention — **gate blocked**
- [x] ≥1 scenario asserts template field survival — Vidoe-Test `state-trust-transition-pilot.test.mjs` (`6665b77`)
- [ ] Draft ids promoted to Ai-skill `validation/scenarios/` — deferred; downstream BDD equivalent green

### D — Graduation

- [ ] O2 / O3 **final** written decision — tentative O3 only; ADR 4 incomplete
- [ ] **Do not** register runtime surface until criteria 4 + integration journey evidence met

---

## Watch-out: execution primitive drift

Many primitives are **not designed** — they are **pulled out by three or four consumers**.

**Do not optimize for becoming a primitive.**

- Phase D 前：不称 slice / lifecycle phase
- 观察最小闭环：**Ownership + Invalidation + Recovery** 能否跨 consumer 成立
- 若能：它不是 UI workflow — 单独 ADR 再议 naming / owner layer
- 若不能：stay O2 conditional gate

---

## Stakeholder review log

| Review | Key outcome |
|---|---|
| #1 | scenario→slice；evidence 不膨胀 |
| #2 | O2 vs O3；Invalidation Event；predictive ADR |
| #3 | Recovery Boundary；trust transition 验题；B.5 rename pressure；field survival ADR；anti optimize-for-primitive |
| #4 | Downstream pilot gate 回写；A0–B.5 pass；C partial；tentative O3 |

---

## 完成条件

- [x] A0 四栏 synced（downstream `bcce737`）
- [x] A1 + B + B.5 填表（downstream `6665b77`）
- [ ] C scenario complete — partial: field survival **pass**; predictive prevention **pending**
- [ ] O2/O3 **final** 书面决策 — tentative O3 documented; await criterion 4
- [ ] Phase D — **blocked** until full ADR criteria

---

## Appendix A — Coupon

| State | Owner | Invalidation Event | Recovery Boundary |
|---|---|---|---|
| `previewLimitReached` | ImmersivePlayerFrame | early redeem success clears mask | refresh + entitlement readback |
| `dialogOpen` | Hoisted Frame (was CouponPanel) | panel unmount | explicit reopen after terminal redeem state |
| `pendingRedeem` | Redeem mutation | cancel / onOpenChange(false) | idle / new user action |
| `playbackAllowed` | Entitlement gate | refresh before grant readback | grant readback + poll window clear |

## Appendix B — Optimistic rollback

| State | Owner | Invalidation Event | Recovery Boundary |
|---|---|---|---|
| `optimisticBalance` | QueryCache | rollback | server sync / query refetch |
| `optimisticUISuccess` | View layer | invalidate (server truth differs) | reconciled view model |
| `pendingMutation` | Mutation hook | supersede | new mutation settled |

## Appendix C — B.5 websocket sketch（rename pressure）

| State | Owner | Invalidation Event | Recovery Boundary |
|---|---|---|---|
| `websocketReady` | Connection manager | reconnect / drop | handshake complete |
| `subscriptionActive` | Client session | connection invalidate | resubscribe ack |

**B.5 pass** = 上表无需改列名即可填写。
