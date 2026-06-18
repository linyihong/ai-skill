# Performance Governance（Delivery Model — active）

`workflow/software-delivery/perf-governance.md` 是 **execution 产品化** slice：把「何时测、怎么跑、怎么存证、如何决策」从 theory（[`validation.md`](validation.md)）与 PR gate（[`perf-risk-gate.md`](perf-risk-gate.md)）拆出来。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-perf-governance` |
| `purpose` | L0–L2 performance delivery：intake → smoke → evidence；`result` 可决策、`stability` 仅观测 |
| `type` | `execution` |
| `tags` | performance, governance, smoke, evidence |
| `load_when` | 效能敏感变更需要 L0 intake + L1 smoke + L2 summary |
| `do_not_load_when` | 纯文档、纯样式、尚未有 project-local runner 的 repo |
| `owner_layer` | workflow |
| `canonical_source` | 本檔（active）；来源 pilot 见下方 `evidence_scope` |
| `dependencies` | `sd-validation`（指标理论）、`sd-test-strategy`（测试类型选型）、`perf-risk-gate`（PR 热路径 gate） |

## Status

```yaml
status: active
evidence_scope:
  incidents: 2          # pagination (747fade) + player aggregation
  incident_classes: [pagination, aggregation]
  environments: 1       # http://16.163.215.93/h5
  cross_time_probe: closed   # P2.7 T0–T3 (2026-06-18)
  confidence: validated
promotion:
  p4b: promoted         # 2026-06-18 team consent
  cross_time: satisfied
scenario: validation/scenarios/software-delivery/perf-smoke-gate.yaml
```

**P4b canonical promote（2026-06-18）**：cross-incident + cross-time 门槛已满足；团队显式同意 `candidate` → `active`。

**Pilot pointer**：[`perf-governance-pilot.pointer.yaml`](../../governance/evidence-candidates/evidence-rules/perf-governance-pilot.pointer.yaml) → Vidoe-Test plan + evidence。

**P2.7 结案**：**Stability labels are execution-context sensitive.** T0 6/17 下午为 transient shared-environment noise；早+晚皆 STABLE。见 external `docs/evidence/perf/reports/2026-06-18-p27-closeout.md`。

---

## 与 validation / perf-risk-gate 的分工

| 层 | 文档 | 拥有 |
|---|---|---|
| **Theory** | `validation.md` | 测什么 metric、load/stress/spike/soak 选型引用 |
| **PR gate** | `perf-risk-gate.md` | 何时进 gate、anti-pattern、reviewer checklist |
| **Execution（本档）** | `perf-governance.md` | L0–L3 怎么跑、证据放哪、result vs stability 语义 |

```text
validation.md (theory sufficient)
        ↓
perf-governance.md (execution productization)   ← active
        ↓
project runner + docs/evidence/perf/            ← Vidoe-Test 实证
```

---

## 四阶段 Performance Delivery Model

| Level | 名称 | 目的 |
|---|---|---|
| **L0** | Perf Intake | ~30s 判定是否触发 perf 流程；存档 YAML |
| **L1** | Perf Smoke | warmup → collect → compare → report；fake 先于 real |
| **L2** | Evidence Store | commit **summary + baseline**；raw 放 gitignore / CI artifact |
| **L3** | Budget Gate | relative regression first → 累积后再 absolute + CI block |

**成熟度顺序**：L0 → L1 fake → G5 adoption → L1 real → L2 summary → review →（延后）L3 / soak / k6。

---

## L0 — Perf Intake

触发：hot-path API、分页/聚合查询、SSR payload、auth 列表、外部扇出、cache 结构变更。

最小 YAML 形状（project-local，路径自定）：

```yaml
perf_intake:
  incident: "<slug>"
  hot_paths:
    - "<api or page>"
  decision: smoke_required | defer
  rationale: "<one line>"
```

---

## L1 — Perf Smoke

**渐进路径**：

```text
P1a fake runner  →  G5 adoption  →  P1b real HTTP
```

**执行语义**：

1. **warmup** — 不计入统计
2. **collect** — measured runs → raw jsonl（gitignored）
3. **aggregate** — p50 / p95 / p99、error_rate
4. **compare** — vs baseline；超 regression → `result: fail`
5. **stability** — bucket p95 variance；**advisory only**（见下）

**固定测试账号**优于 ephemeral register（auth target 可复现）。

参考实现：`Vidoe-Test/tests/perf/perf-smoke.mjs` · `npm run perf:smoke`。

---

## L2 — Evidence 分层

```text
.tmp/evidence/perf/     # raw runs — gitignored
docs/evidence/perf/     # summary + baseline — commit
  intake/
  summaries/
  baseline.yaml
  reports/
```

**禁止**：将全量 raw JSON commit 进 repo。

Summary 须含 `schema_version`、`env`、`comparison` per target、`variance_source`。

---

## 信息模型（governance principle — active）

> **Performance result 可决策；stability 仅观测。**

| 字段 | 决策权 | 说明 |
|---|---|---|
| `result` | **可作 merge 决策** | pass / fail — error_rate + baseline regression |
| `stability.status` | **advisory only** | 环境观测；短窗口可能不可重现 |

```yaml
stability:
  semantics:
    type: advisory
  escalation:
    require_confirmation_runs: 2
  merge_policy:
    hard_gate: false
  execution_context_required: true   # P2.7 — record started_at, time_window, weekday, env
```

**模型（P2.7 验证）**：`stability ≈ f(service variance, shared load, execution timing)` — **不是** `f(api quality)`。

**不能说**：「早上比较快 / 比较准」→ **只能说**：「这次 run 的 stability 标签受 execution context 影响」。

单次 `unstable` **不具有决策权**。不可因单次 UNSTABLE 优化 API 或未经校准调 threshold。

**execution_context**（summary / report 须记录；runner 实现待 project）：

```yaml
execution_context:
  started_at: ISO-8601
  time_window: morning | afternoon | evening
  weekday: <day>
  env: <base_url>
  # 预留：host_load, cache_state
```

**延后**：load/stress/soak — 共享环境噪音可让 smoke 翻盘；先收稳 execution context + baseline。

**本档明确不 promote 的数值/结论**（留在 project evidence）：

- variance threshold 具体百分比
- tail detector / `tail_ratio` 标定
- `stability_flip_rate` promotion gate
- 单次 incident 的诊断结论（如 localized tail-event pattern）

---

## G5 Adoption（先于 infra 扩建）

流程价值验证优先于工具膨胀。最低检查：

| ID | 检查 |
|---|---|
| G5-1 | 全 suite < 3 min |
| G5-2 | stdout 一屏 + summary markdown |
| G5-3 | 一条命令可跑（`npm run perf:smoke`） |
| G5-4 | stability 可复现性 — **open / advisory**；不作 hard pass 条件 |
| G5-5 | summary 含 env + auth 说明 |
| G5-6 | 参与者 < 10 min 可跑完并解释 |

G5 Pass → project-local governance **值得继续**；P4b promote 后 canonical 拥有流程与语义，数值标定仍留在 project evidence。

---

## L3 — Budget（延后）

Pilot **不** 将 CI block 写入 canonical。需 **2+ baseline** 跨 incident 后再定义 `budgets/` + CI job。

---

## 三层边界（project vs canonical）

```text
Evidence（project docs/evidence/perf/）     ← 审计链，保留数字与 run 细节
        ↓
Project governance（runner + plan）        ← 已验证
        ↓
Canonical theory（本档 active）         ← 流程与语义；P4b promoted 2026-06-18
```

**可迁移知识**：如何把 performance **验证产品化** — 不是某个 `threshold = 15%`。

---

## Related

- [`validation.md`](validation.md) — 指标与测试类型 theory
- [`perf-risk-gate.md`](perf-risk-gate.md) — PR 阶段 perf gate（`active`）
- [`perf-smoke-gate.yaml`](../../validation/scenarios/software-delivery/perf-smoke-gate.yaml) — L1 pass/fail scenario
- [`test-strategy.md`](test-strategy.md) — load/stress/spike/soak 选型
- Vidoe-Test pilot plan（external）：`docs/plans/2026-06-17-1400-performance-validation-architecture-pilot.md`
- P3 review（external）：`docs/evidence/perf/reports/2026-06-17-p3-review-final.md`
- P2.7 close-out（external）：`docs/evidence/perf/reports/2026-06-18-p27-closeout.md`
- Player generalization（external）：`docs/evidence/perf/reports/2026-06-17-player-episode-generalization.md`

← [Back to software-delivery workflow](README.md)
