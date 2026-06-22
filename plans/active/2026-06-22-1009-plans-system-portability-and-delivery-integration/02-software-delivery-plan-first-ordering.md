---
id: 2026-06-22-1009-software-delivery-plan-first-ordering
plan_kind: sub
status: draft
owner: linyihong
created: 2026-06-22
parent: 2026-06-22-1009-plans-system-portability-and-delivery-integration
required_for_completion: true
sub_plan_reason: >
  plan-first ordering 只動 software-delivery workflow 文件（intake 段），與 01
  的 CLI / validator 工作完全不同 owner / 不同 review 焦點，可獨立 sign-off 並
  與 01 並行。獨立成 sub-plan 避免把 workflow doc 改動混進 01 的 Go / CLI commit。
---

# Software-delivery Plan-First Ordering（sub-plan）

**Status**: `draft`
**Owner**: linyihong
**Parent**: [`_plan.md`](_plan.md)

## Source Request
讓 software-delivery 接入 plans 系統：「以後開發需要先寫 plans 系統把所有規劃好」。使用者選 **workflow 層 ordering（advisory）**，不做機械 block。

## Scope
- **In**：在 `workflow/software-delivery/` intake 段加入 plan-first ordering，明文接在 pre-build-interrogation / Test-First Ordering 之後；advisory + review checklist；可選 validation scenario。
- **Out**：commit-msg 機械 block（無 active plan 不准 commit code）— 保留為後續 maturity ladder 升級候選。
- **Affected**：`workflow/software-delivery/intake.md`、`workflow/software-delivery/execution-flow.md`（導航）、可能 `workflow/software-delivery/test-strategy.md`（Test-First Ordering 接點）、可能新增 validation scenario。

## Decision Rationale（sub 層）
現有 intake 已有 pre-build-interrogation（goal/scope/non-goals/acceptance/source-of-truth/duplication risk）與 Test-First Ordering（framework/runtime/governance 升級強制順序）。plan-first **不是新 gate，而是把「實作前先有 plan artifact」明文化為 intake 的一環**，並連到 plans 系統（`plans/active/` + plan-tree）。

**關鍵修正（回應 review #4）：plan 是 artifact，不是一個一次性 stage。** 線性「interrogation → plan → preflight → implement」是錯的，因為 preflight 會**回改 plan**（架構相容性檢查發現衝突就更新 plan）。正確模型是 plan 在 preflight 間反覆收斂：

```
Discover → Interrogate → Draft Plan ⟲ Preflight → Execute
                              └──────────┘
                         preflight 可回改 plan
                         （plan 非一次生成）
```

不重複 pre-build-interrogation（Q4）的分工：pre-build-interrogation = 需求拷問（產出 plan 的輸入）；plan-first = 拷問結果落成可收斂的 plan artifact；Architecture Compatibility Preflight = 對 draft plan 做相容檢查並回饋修正。三者不是序列三段，而是 interrogation 餵入、plan 為中心 artifact、preflight 反覆驗證。

### Alternatives
- A. 硬機械 gate：reject（本輪）— 易誤擋小修補，使用者已選 advisory。
- B. 完全不接、只靠既有 preflight：reject — 缺「先有 plan artifact」的明文順序，plan 與 delivery 仍脫節。
- C. advisory ordering + review checklist（accept）。

## Open Questions（本 sub）
- Q4（與 pre-build-interrogation 不重複）— 見 main §Open Questions。
- 是否所有 software-delivery 任務都要 plan，還是依規模分級（小修補豁免）？

## Phase 0 — Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）
- [ ] 已讀 main + 本 sub §Open Questions 全部條目
- [ ] 對每條標記 `resolved` / `still-open` / `deferred`
- [ ] `resolved` 條目回寫
- [ ] 新問題已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q4 不重複 pre-build-interrogation | still-open | Phase 0.1 讀 intake.md 後界定分工 |
| 規模分級豁免 | still-open | Phase 1 決定 advisory 文字 |

### Phase 0.1 — 架構盤點
- [ ] 讀 `workflow/software-delivery/intake.md` 現行 intake 順序與 pre-build-interrogation 內容。
- [ ] 讀 `workflow/software-delivery/test-strategy.md` Test-First Ordering，確認接點。
- [ ] 讀 `plans/README.md` Architecture Compatibility Preflight（已要求實作前 preflight），界定 plan-first 與 preflight 的關係（避免三重 gate）。

## Phase 1 — Plan-first ordering 文件化
- [ ] 在 intake 段新增「Plan-First Ordering」小節：明文「會導向 code/workflow/governance/runtime 改動的任務，實作前須有對應 `plans/active/` plan（可 inline 小 plan 或 plan-tree）」。
- [ ] **用 loop 模型描述（非線性三段）**：`Discover → Interrogate → Draft Plan ⟲ Preflight → Execute`；明寫 plan 是 artifact、preflight 可回改 plan，避免讀者誤以為 plan 一次生成後不再變。
- [ ] 分工說明：interrogation 餵入 / plan 為中心 artifact / preflight 反覆驗證並回饋（Q4）。
- [ ] 規模分級：< 1 session / 純文件補強 / surgical 小修補的豁免條件（接 plan-tree「何時不開 sub-plan」既有規則）。
- [ ] advisory 語氣（「應 / 建議」而非 commit-time block）。

## Phase 2 — Review checklist + 可選 scenario
- [ ] 在對應 review / DoR checklist 加一條 plan-first 檢查項。
- [ ] （可選）新增 validation scenario 描述「實作前有 plan / 豁免條件成立」的可觀察判定；**若不接 runtime gate，明寫 doc-only + 未來升級條件**。
- [ ] 更新 `execution-flow.md` 導航指向新小節。

## 完成條件
- [ ] intake plan-first ordering 小節落地，用 loop 模型、與 pre-build-interrogation / preflight 分工清楚（Q4 resolved）
- [ ] 規模分級豁免條件落地
- [ ] review checklist 更新
- [ ] doc-only 宣告明確（若未接 runtime gate）
- [ ] **Acceptance evidence（回應 review #6）**：至少一次**真實的 software-delivery intake** 走過 plan-first loop（含一次 preflight 回改 plan 的實例），而非僅文件範例
- [ ] linked-updates 檢查（README / 導航 / checklist 同步）

## Glossary Impact
Glossary Impact: no — plan-first ordering 復用既有 plans 系統與 software-delivery 詞彙，no new framework vocabulary introduced。

## 與其他 plans 的關係
- 接入 [`workflow/software-delivery/intake.md`](../../../workflow/software-delivery/intake.md) 與 [`workflow/software-delivery/test-strategy.md`](../../../workflow/software-delivery/test-strategy.md)。
- 復用 [`plans/README.md`](../../README.md) Architecture Compatibility Preflight 與 plan-tree「何時不開 sub-plan」規則。
