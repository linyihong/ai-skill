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
**最後修訂**：2026-06-16（review #3：Recovery Boundary；trust transition 验题；B.5 rename pressure；template field survival ADR；do not optimize for primitive）
**Priority**：**P1**

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

| Path | 选 when |
|---|---|
| **O2 — Conditional gate** | B + B.5 失败：栏位不能跨域复用 |
| **O3 — Generic trust transition model** | B + B.5 通过：四栏不改字可套 rollback + websocket |

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

### A0 — Template + Recovery Boundary（当前）

- [x] 四栏定义写入本 plan
- [ ] Downstream project overlay / screen-mapping sync Recovery Boundary
- [ ] Side Effect Chain：`invalidation_events` + **`recovery_evidence`** per step where applicable

### A1 — Coupon

- [ ] 四栏填 coupon（Appendix A）
- [ ] Counterfactual：A0 前能否 predict unmount + recovery（grant readback / refresh window）

### B — Optimistic rollback

- [ ] Primary：rollback invalidate + server sync recovery 是否自然落入四栏
- [ ] 验 **trust transition 闭环**，非仅「三栏够不够」

### B.5 — Rename pressure test

- [ ] **不用同名案例** — 例如 websocket reconnect / handshake complete
- [ ] 问：四栏能否 **不改字** 套用？
- [ ] **不能 ⇒ 先别 O3**

### C — Scenario spike

- [ ] ≥1 previously unknown prevention
- [ ] ≥1 scenario asserts template field survival across cases
- [ ] Draft ids：`state-trust-transition-invalidate-recover-v1`、`observable-outcome-survives-owner-refresh-v1`、`temporal-behavior-ownership-subshape-v1`

### D — Graduation

- O2 / O3 only；O1 optional future promotion
- **Do not** register runtime surface until criteria 4+5 met

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

---

## 完成条件

- [ ] A0 四栏 synced
- [ ] A1 + B + B.5 填表
- [ ] C scenario（predictive + field survival）
- [ ] O2/O3 书面决策
- [ ] Phase D only if full ADR criteria

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
