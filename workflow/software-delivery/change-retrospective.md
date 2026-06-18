# Change Retrospective Slice（Ship → Retrospective）

> **Cognitive Slice**：`sd-change-retrospective`（incident-driven change 的正式學習出口；防止 overlay 堆積或直接 canonical promote）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-change-retrospective` |
| `purpose` | Ship 後記錄哪層被改、哪層刻意沒改、vocabulary/consumer 變化與 promotion 建議 — 僅三種結論 |
| `type` | `execution` |
| `tags` | closure, retrospective, promotion, incident |
| `load_when` | UI / consumer incident 路徑已 Ship；Phase / drill 關閉；考慮寫入 project overlay 或 feedback |
| `do_not_load_when` | 變更尚未 Ship；純 greenfield feature 無 incident 決策鏈 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「學到什麼」的正式出口與禁止 direct promote；不承載 runtime state model |
| `canonical_source` | 本檔 |
| `dependencies` | `sd-closure`（DoD passed）、`sd-ui-incident-governance`、[`layer-ownership-matrix.md`](layer-ownership-matrix.md) |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Retrospective 存在；promotion ∈ {keep local, promote project, candidate canonical}；**allowed_outputs ≤ 1** |

---

## Placement — decision-driven delivery loop

Software-delivery 對 incident 路徑不是線性「需求 → 合約 → 實作」，而是 **證據驅動變更** 決策鏈：

```text
Discover
  ↓
Observe → Classify → Select Layer    ← 決策（evidence-driven）
  ↓
Execute                              ← Contract · Implementation · Verification
  ↓
Ship                                 ← sd-closure DoD
  ↓
Retrospective                        ← this file（決定知識去哪 — 非第二個 governance hub）
```

**Closure 語意變更**：Ship 不是終點 — closure = 交付完成 **+** 決定知識去向（project · feedback · canonical 三出口正式拆分）。

**Learning system boundary**: Retrospective 是 **learning outlet**，不是 accumulation hub。一次变更一条主去向。

**Forbidden**:

```text
Ship → 直接新增 overlay / 直接 canonical promote / 新 abstraction hub
Ship → 同時 promote project + candidate canonical + feedback（知識分叉）
```

---

## allowed_outputs ≤ 1（硬规则）

一次变更的 Retrospective **只能有一个主去向**：

| 主去向 | 允许写入 |
| --- | --- |
| **keep local** | plan § / incident note only |
| **promote project** | `.ai-skill/project/rules/` **或** `feedback/` — 择一 primary |
| **candidate canonical** | `feedback/history/` lesson + plan pointer |

**Forbidden in same retrospective**:

- promote project **and** candidate canonical
- candidate canonical **and** new project overlay body
- 多个 parallel promotion 动作

Supporting artifacts（G4 sheet、commit、drill memo）不算 second output — 它们是 evidence，不是 knowledge destination。

若 lesson 值得记录但主去向是 keep local，可在 incident note 内嵌 «rejected promote» 理由 — **不**另开 feedback 文件。

---

## When to run

After **Ship** (merge, deploy, or bounded task closure) when the work path included:

- Incident card ([`incident-observation.md`](incident-observation.md))
- Classification + layer selection ([`ui-incident-governance-workflow.md`](ui-incident-governance-workflow.md))
- Any modification through Contract / Overlay / Verification / Integration

Greenfield features without incident decision chain may use [`closure.md`](closure.md) §Feed Back only.

---

## Retrospective record

| 問題 | 回答 |
| --- | --- |
| **哪層被修改** | Primary layer actually changed: Contract \| Overlay \| Verification \| Integration |
| **哪層刻意沒改** | Layers considered and rejected (with reason) |
| **有沒有新增 vocabulary** | yes/no — terms added; if yes, list and check synonymy |
| **有沒有新增 consumer** | yes/no — consumer #N; if yes, does not alone justify abstraction |
| **promotion 建議** | **exactly one** primary output — `allowed_outputs ≤ 1` |

---

## Promotion 建議（僅三種 — 禁止 direct promote）

| 建議 | 含义 | 下一步 |
| --- | --- | --- |
| **keep local** | 單次 incident；無 reusable 決策 |  incident note / plan § only；不新增 overlay |
| **promote project** | Repo-specific pattern 穩定 | `<PROJECT_ROOT>/.ai-skill/project/rules/` or `feedback/` per project `overlay-lifecycle.md` |
| **candidate canonical** | 可能跨 project；**未證實** | `feedback/history/` lesson + plan；**禁止**直接写 `governance/` or Experience Runtime |

**Blocked outcomes**:

- Direct canonical promote
- Experience Runtime Governance from one pilot
- New hub/abstraction without second independent incident

Canonical promotion path: [`decision-promotion-pipeline.md`](../../governance/lifecycle/decision-promotion-pipeline.md).

---

## Template

```markdown
## Change retrospective

- Incident / task ref:
- Ship evidence: (commit, deploy, drill sign-off)
- Primary output (exactly one): keep local | promote project | candidate canonical

| Question | Answer |
| --- | --- |
| Layer modified | |
| Layers deliberately not changed | (from Select Layer rejection table) |
| New vocabulary | |
| New consumer | |
| Promotion suggestion | (must match Primary output above) |

- Rejected promotion paths: (e.g. «did not promote project — single consumer»)
- Notes for next incident:
```

Store in project plan §、incident closure memo、或 `feedback/` — not in Ai-skill workflow body.

---

## Refs

- Layer ownership (core capability): [`layer-ownership-matrix.md`](layer-ownership-matrix.md)
- Project overlay lifecycle: `<PROJECT_ROOT>/.ai-skill/project/overlay-lifecycle.md`
- Governance gate: [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md) §Change retrospective
