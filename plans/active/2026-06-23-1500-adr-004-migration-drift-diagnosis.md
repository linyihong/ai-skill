---
id: 2026-06-23-1500-adr-004-migration-drift-diagnosis
plan_kind: main
status: frozen
execution_status: deferred
owner: linyihong
created: 2026-06-23
parent: null
---

# ADR-004 Migration Drift — Frozen Diagnosis (Evidence Baseline)

## Status

**Frozen diagnosis.** 本文件是證據基線（evidence baseline），不是實作計畫。
Implementation phases **DEFERRED** — 後續的 *ADR-004 Migration Completion Plan* 將引用本文件，
不在此處 draft YAML / validator / code。凍結原因：診斷已完整收斂，需在 transcript compaction
前固定成可引用的正式紀錄，避免「兩天後忘記這次挖出什麼」。

## 一句話

這不是 `feedback_history` 路徑筆誤，而是一個**自我認證的 drift propagation loop**：ADR-004 在
2026-05-13 把 canonical sink 改成 `feedback/history/<domain>/`，但 Runtime / Fixture / Health-Check
三層各自基於已廢棄的世界觀互相證明對方正確，使錯誤模型得以持續自我複製並通過 CI。

## Layer 狀態表

| 層級 | 狀態 | 證據 |
| --- | --- | --- |
| Constitution（ADR-004） | ✅ 正確 | `constitution/ADR-004-feedback-promotion-pipeline.md` Vocabulary Evolution：`skills/*/feedback_history/` → `feedback/history/<domain>/<category>/`，2026-05-13 完成搬遷、舊路徑已刪除 |
| Knowledge Layer | ✅ 正確 | `feedback/README.md`、`enforcement/content-layering.md` 均宣告單一 sink |
| Feedback Layer | ✅ 正確 | lesson 實體已在 `feedback/history/<domain>/` |
| Routing Layer | ✅ 正確 | `knowledge/runtime/routing-registry.yaml` `route.feedback.history` `primary_source: feedback/history/README.md` |
| Runtime Index Layer | ❌ 錯誤 | `scripts/ai-skill-cli/internal/app/runtime.go:1213` glob `skills/*/feedback_history`（已刪除目錄）→ 0 lessons indexed |
| Validation Fixture Layer | ❌ 錯誤 | `runtime_test.go:244,625,640` 製造 `skills/demo/feedback_history/` 驗證已不存在的世界 |
| Runtime Health-Check Layer | ❌ 錯誤 | `runtime.go:1876,1879` `COUNT(fts MATCH 'feedback') > 0` — token match，非 provenance |
| Agent Direct-Read Path（Path 1） | ✅ 正確 | 經 `route.feedback.history` 直接讀 `feedback/history/` |
| Runtime Query Path（Path 2） | ❌ 失效 | `ai-skill runtime query`（`runtime.go:156,160`；FTS at `1377-1379`）對 lesson 回傳空集合 |

## Consumer Trace 結論：誰真的依賴 feedback lessons？

- **Path 1 — Reference-first routing（agent 直接讀檔）：✅ 正常。** `route.feedback.history` 指向
  canonical `feedback/history/README.md`，agent 以 reference-first 直接讀目錄，**繞過 index**。
  這是目前**真正在運作**的 discovery path，也是「臨時分析時仍找得到 lesson」的原因。
- **Path 2 — `ai-skill runtime query`（FTS index）：❌ 失效但 agent-facing。** `runtime query` 子命令
  對 atoms/fts 做 FTS 搜尋（`runtime.go:1377-1379`），與 feedback records 同表（`runtime.go:1146`）。
  因 glob 掃空目錄，lesson 從未進 index，故 `runtime query <topic>` 對 lesson **靜默回傳空**。
- **Health-Check（`nativeRuntimeIndexFTSCheck`）：自我驗證，非真實 consumer。** 它是遮蔽者，不是依賴者。
- **無任何 workflow/doc 指示 agent 用 index 查 lesson** → Path 2 是 latent 次要 path，非 load-bearing。

## Discovery Path 分析 → Capability Partition Failure

風險評級因 Consumer Trace 改變：這**不是** System Outage（discovery 全壞），而是
**Capability Partition Failure**——

```
Feedback Discovery
├─ Direct Read (Path 1)    ✅ 使用者視角：agent 仍找得到 lesson
└─ Runtime Query (Path 2)  ❌ runtime 視角：低 token discovery 已失效
```

→ 因此 P0 修復的 justification **不是** "restore discovery"，而是
**"restore consistency between two discovery paths"**。兩條 path 已產生認知分叉。

## Drift Propagation Loop（masking chain）

```
        constitution/ADR-004  (Accepted · 2026-05-13)
                  │  declares canonical sink
                  ▼
        feedback/history/<domain>/   ◄────── REALITY: lessons live here
                  │
                  ├──────────────────────────────────────────────┐
                  │  PATH 1  route.feedback.history               │
                  │  primary_source → feedback/history/README.md  │ ✅ agent finds lessons
                  └──────────────────────────────────────────────┘   (bypasses everything below)

   ════════  the self-certifying broken loop — every node points at the OLD world  ════════

   runtime_test.go ──seeds──►  skills/demo/feedback_history/   (manufactures a deleted world)
        │ green unit test                        ▲
        ▼                                        │ "confirms" indexer is correct
   runtimeIndexFeedbackRecords ──glob──► skills/*/feedback_history/   ✗ DELETED 2026-05-13
        │                                        └─► 0 lessons indexed
        ▼
   atoms / fts   ◄── only route + summary atoms that happen to contain the token "feedback"
        │
        ▼
   nativeRuntimeIndexFTSCheck:  COUNT(MATCH 'feedback') > 0   ✅ PASSES on route/summary atoms
        │      └── Semantic Validation Drift: asserts TOKEN presence, not lesson PROVENANCE
        ▼
   GREEN CI ───────────────────────────────────────────────────┐
        │                                                       │
        └──► next author (2026-05-21) trusts green, mirrors ────┘
             the fixture's mental model → writes MORE old-world code

   ── real-but-starved consumer ──
   ai-skill runtime query <topic> ──► fts ──► (would surface lessons if indexed) ──► returns none
```

關鍵：`fixture → indexer → health-check → green → next author → fixture` 是一個**封閉迴圈**，
每個節點都對著已退役的世界互相驗證；loop 內**沒有任何節點查 Constitution**。

## 命名的兩個失效模式

### Semantic Validation Drift
Validator 並沒有驗證錯 implementation，而是驗證了**錯誤的 success condition**。
想驗證的是「feedback lessons 可被發現」；實際驗證的是「存在 token = feedback」。
`route.feedback.history`（route atom）與 `feedback/history/apk/lesson-001.md`（lesson）在語義上
不是同種證據，但 `MATCH 'feedback'` 把兩者視為同一證據。

### Canonical Knowledge Propagation Failure（上位模式）
> 命名修正：原為 *Knowledge Source Drift*，但 source 並沒有錯——`constitution/ADR-004` 與
> `feedback/history/` 都正確。錯的是 **propagation**：正確的 source-of-truth 更新無法傳播到
> runtime / fixture / validator。失效點是「傳播」，不是「來源」。

ADR / Runtime / Fixture / Health-Check 在不同地方共享同一個**已退役**的認知 → Source of Truth 沒有真正被
機械化、單向傳播。三處 mental model 互證，形成封閉 loop，使 ADR 的世界更新無法傳播到 runtime/validation。
（本次 drift loop 本身就是 propagation 的失效證據——見上方 masking chain。）

### 可重用模式（未來一定會再出現）
```
Constitution Drift → Fixture Drift → Validation Drift → Green CI → Knowledge Reproduction
```
本次實例是 `feedback_history`；下次可能是 promotion status enum / routing taxonomy /
enforcement registry / runtime schema。

### Invariant（可重用治理資產）

把本次最關鍵的那句（「loop 內沒有任何節點查 Constitution」）從敘述升成 invariant，
讓它不只服務 ADR-004，而能直接套用於上述任一未來實例：

> **Canonical-Path Derivation Invariant**
> No runtime / validator / fixture may derive canonical paths without passing
> through the contract or the registry.
> （任何 runtime / validator / fixture 都不得自行推導 canonical path；必須經由
> contract 或 registry 取得。）

這條 invariant 是上方 *Contract Ownership* 鏈的可驗證表述：違反它＝出現一個繞過 registry 的
path 宣告＝drift loop 的入口。未來只要某層直接 hard-code 路徑（如本次 `skills/*/feedback_history`
死 glob），即違反此 invariant，應在 review / 機械化檢查階段被攔截。

## P0 順序調整（observability before repair）

| 新序 | 項目 | 理由 |
| --- | --- | --- |
| **P0-A** | Health-check tightening（`nativeRuntimeIndexFTSCheck`，`runtime.go:1879`）：由 `MATCH 'feedback' > 0` 改為「≥1 atom 且其 `source_path` 在 canonical lesson sink 之下」（provenance，非 token） | 目前最大風險不是「Path 2 壞」，而是「Path 2 壞但看起來正常」。沒有可靠量測，後續所有修復皆不可觀測。 |
| **P0-B** | Runtime indexer repoint（`runtime.go:1213`）：`skills/*/feedback_history` → `feedback/history/<domain>/**`，路徑改由 registry 提供 | 修復 Path 2（次要低 token path）；恢復兩條 discovery path 的一致性。 |
| **P0-C** | Test reality alignment（`runtime_test.go:244,625,640`）：seed `feedback/history/<domain>/`；新增 negative test 確保舊路徑**不**被索引 | 停止 fixture 對 deleted world 的教學；鎖門防回歸。 |

> **P0-A 不保證 repair。** Success of P0-A only establishes truthful observability;
> it does NOT imply feedback lessons are indexed. P0-A 修的是「量測正確性」
> (measurement correctness)，不是「runtime 正確性」(runtime correctness)——health-check
> 綠燈在 P0-B/C 落地前，只代表「我們現在能誠實看到 Path 2 仍是空的」，不代表 lesson 已可被 query。

（原序為 runtime.go → test → health-check；調整理由＝Observability 先於 Repair。）

## Contract Layer（後續方向，非本文件實作）

本次找到的是「Source of Truth 未機械化」，故修復應上升為 contract，而非單一 path rule。
建議位置 `knowledge/runtime/contracts/`，形成 **Constitution → Contract → Enforcement**：

```
constitution/ADR-004                         ← human-readable constitution（知識）
knowledge/runtime/contracts/<name>.yaml      ← machine-readable contract（可執行知識）
enforcement/ validator implementation        ← executable enforcement（執行器）
```

Validator / runtime / docs / tests 一律讀 contract（與其指向的 canonical-path registry），
不直接耦合 ADR；未來改 sink 只動一處。此 contract layer 補上 loop 缺的那條邊：`Constitution ⊨ Runtime`。

### Contract Ownership（避免 contract 自己再 drift）

四個 source（ADR / runtime contract / routing registry / validator yaml）若各自宣告 canonical path，
contract layer 本身就會變成下一個 drift 點。因此 **ownership 必須寫死、單向 derive**——
**Contract Owner = Runtime Registry**：

```
constitution/ADR-004              （human-readable 知識）
        │ derives
        ▼
knowledge/runtime/contracts/<name>.yaml   （machine-readable contract）
        │ materializes
        ▼
knowledge/runtime/routing-registry.yaml   （canonical-path registry · 唯一 path owner）
        │ consumed by
        ▼
runtime / validator / tests       （只讀，不得自行宣告 path）
```

關鍵：path 的唯一擁有者是 registry，contract 只 materialize 不 fork；runtime/validator/tests 一律
下游消費，禁止平行宣告。這把「四個 source 各自為政」收斂成一條單向 derive 鏈。

## 證據錨點（file:line）

- `constitution/ADR-004-feedback-promotion-pipeline.md` — Vocabulary Evolution 表（canonical sink）
- `scripts/ai-skill-cli/internal/app/runtime.go:1146` — feedback records 併入 index
- `scripts/ai-skill-cli/internal/app/runtime.go:1212-1213` — `runtimeIndexFeedbackRecords` / 死 glob
- `scripts/ai-skill-cli/internal/app/runtime.go:156,160,1377-1379` — `runtime query` + FTS（Path 2 consumer）
- `scripts/ai-skill-cli/internal/app/runtime.go:1874-1898` — `nativeRuntimeIndexFTSCheck`（masking gate）
- `scripts/ai-skill-cli/internal/app/runtime_test.go:244,625,640` — old-world fixtures
- `knowledge/runtime/routing-registry.yaml:1951` — `route.feedback.history`（Path 1，正確）
- glob 引入時間：2026-05-21（migration 後 8 天，authored dead）
- 已落地的部分 Class A 文件修正：commit `4f50822f`（僅 `workflow/` + `analysis/`，`enforcement/` 仍 drift）

## Scope / Next

- 本文件僅凍結診斷與 P0 順序；**不**含 implementation。
- 後續開 *ADR-004 Migration Completion Plan*（separate plan，引用本 id）時，再展開：Step 0 完整
  99-file Classification + Intent Mapping、P0-A/B/C 實作、Class A/D canonicalization、Contract Layer
  與 canonical-path registry、validator、rollback / success criteria。
