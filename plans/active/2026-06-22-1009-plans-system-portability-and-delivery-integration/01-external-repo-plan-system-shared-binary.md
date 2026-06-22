---
id: 2026-06-22-1009-external-repo-plan-system-shared-binary
plan_kind: sub
status: draft
owner: linyihong
created: 2026-06-22
parent: 2026-06-22-1009-plans-system-portability-and-delivery-integration
required_for_completion: true
sub_plan_reason: >
  外部化是三件事的地基：必須先切出 portable core vs Ai-skill governance
  overlay 的邊界（哪些 validator 屬 plan_profile），02 / 03 才知道
  自己改的東西算 portable 還是 Ai-skill-only。獨立成 sub-plan 以便先 graduate
  邊界共識（Q2）與跨 repo 強制機制（Q1/Q3），且其 acceptance（外部 repo
  真實跑過 validate）可獨立 sign-off。
---

# External-repo Plan System via Shared Binary（sub-plan）

**Status**: `draft`
**Owner**: linyihong
**Parent**: [`_plan.md`](_plan.md)

## Source Request
讓外部 repo 透過共用 `ai-skill` binary 使用 plans 系統（使用者選「共用 binary 指向外部 repo」，非 init-project 抽取）。

## Scope
- **In**：portable 邊界**推導模型**（`plan_profile` capability + `plan_schema` 相容契約，非預設清單）；可重用 **validator engine package**（與 commit-msg hook 解耦）；外部 repo 跑 plan validators 的跨 repo consumer 路徑；schema / 版本相容策略。
- **Out**：把整套 Ai-skill 治理（runtime.db / glossary / cognitive modes / ADR pipeline）搬進外部 repo；init-project 抽取安裝（保留為 future option）。
- **Affected**：`scripts/ai-skill-cli/internal/app/plans.go`、`plan_tree.go`、`hooks.go`、（新）validator engine package、`plans/README.md` 或 `governance/lifecycle/`、新 `ai-tools/` 或 `scripts/ai-skill-cli/docs/` 外部 repo 使用說明。

## Decision Rationale（sub 層）

讀取面（`plans tree --root`）已跨 repo；唯一缺口是 **commit-time 強制**。現有 plan-tree validators 在本 repo 的 `commit-msg` hook 內由 `runCommitMsgHook` 呼叫，綁本 repo。

**抽象層決定（回應 review #3）：要抽的不是 CLI command，而是 validator engine。**

```
validator engine package          ← 核心，read-only，吃 (root, staged-set)
        ↓ consumers (薄)
  ├── commit-msg hook (本 repo 既有)
  ├── git hook shim   (外部 repo, 放 ai-tools/)
  ├── CI              (任何 repo, 直接呼叫)
  └── plans validate  (CLI surface, 也只是 consumer)
  └── future API
```

若把 `ai-skill plans validate` 當核心，半年後它會長成另一個 orchestration layer（CLI 累積 flag / 狀態 / 隱性 contract）。因此 CLI 只是其中一個薄 consumer。

**portable 邊界（回應 review #2）：不預設「plan-tree 5 + archival 2 = portable」**。portable 不是看 validator 類型，而是看 contract → dependency → execution context。Phase 1 必須**先建分類模型再分類**，否則會變「先決定 portable 再找理由」。

## Open Questions（本 sub）
- Q1（跨 repo 強制機制）/ Q2（portable core 成員）/ Q3（版本相容）— 見 main plan §Open Questions。

## Phase 0 — Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）
- [ ] 已讀 main + 本 sub §Open Questions 全部條目
- [ ] 對每條標記 `resolved` / `still-open` / `deferred`
- [ ] `resolved` 條目回寫
- [ ] 新問題已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q1 跨 repo 強制機制 | still-open | Phase 0.1 盤點 hooks.go engine 可抽性；engine→consumers 抽象 |
| Q2 portable 邊界 | still-open | Phase 1 先建分類表（contract/dependency/context）再推導，非預設清單 |
| Q3 schema/版本相容 | still-open | Phase 2 `plan_schema` version 宣告 + Phase 3 跨版本 evidence |

### Phase 0.1 — 架構盤點
- [ ] 讀 `hooks.go` `runCommitMsgHook`：哪些 validators 依賴 Ai-skill repo-local 狀態（routing-registry / runtime.db / commit context），哪些只吃 (root, staged-set)？
- [ ] 確認 validators 是否已可在不依賴 commit context 下執行（決定 engine 抽取成本）。

## Phase 1 — Portable 分類模型 + 邊界推導（先模型後分類）
- [ ] 先產**分類表**，每個 commit-msg validator 一列，**從 contract 推導 portable**，不從類型直覺：

  ```
  | validator | contract_source | runtime_dependency | execution_context | portable | reason |
  ```

- [ ] 依分類表決定 `plan_profile`（哪些 validator 對外部 repo 適用）與排除清單；`plan_schema` 記錄 frontmatter schema + version。
- [ ] 文件化於 `plans/README.md` 或新 `governance/lifecycle/plan-profile.md`。
- [ ] 完成條件：分類表 + 邊界 review 通過，Q2 標 resolved（且邊界由分類表推導，非預設）。

## Phase 2 — Validator engine package（核心）+ thin consumers
- [ ] **抽出 validator engine package**：read-only，input = (root, staged-set)，output = findings；不依賴 commit context、不綁特定 hook。
- [ ] 既有 commit-msg hook 改為呼叫 engine（重構，行為不變；保護既有測試）。
- [ ] CLI `plans validate --root <path> [--format text|json]` 作為**薄 consumer**呼叫 engine（CLI 不持有驗證邏輯）。
- [ ] schema version 宣告（`plan_schema`）以支援版本相容（Q3）。
- [ ] 測試：engine 單元測試（合法 tree + 各 violation）≥ 5 case；CLI consumer smoke test。
- [ ] **若新增 `route.*` 或 runtime surface，補 Runtime Execution Path + Per-surface consumer 表**（否則明寫 engine/CLI-only，無新 route）。

## Phase 3 — 外部 repo consumer 路徑（git hook shim / CI）
- [ ] `ai-tools/` 或 `scripts/ai-skill-cli/docs/` 寫外部 repo 使用說明（共用 binary 路徑、engine 接 CI / git hook）。
- [ ] 提供薄 `commit-msg` shim 範例（呼叫共用 binary，tool-neutral）。
- [ ] **Acceptance evidence（回應 review #6，收緊）**：
  - [ ] tmp fixture repo：engine + CLI pass / fail 輸出
  - [ ] 一個**真實的非 Ai-skill repo**：實裝 shim，真實 commit 觸發一次 pass + 一次 block
  - [ ] **跨 binary 版本驗一次**：升一次共用 binary（或改一次 `plan_schema` version），確認外部 repo 相容行為符合 Q3 策略
  - [ ] **rollback evidence（回應 review #6）**：外部 repo 可移除 integration（remove hook shim + config）並恢復 clean，證明接入非侵入、可逆

## 完成條件
- [ ] portable 分類表 + `plan_profile` / `plan_schema` 邊界落地（Q2 resolved，且由分類表推導）
- [ ] validator engine package + thin consumers（hook / CLI）+ 測試通過
- [ ] 既有 commit-msg hook 行為不變（重構回歸驗證）
- [ ] 外部 repo 使用說明 + shim 範例落地
- [ ] Acceptance evidence 四項（tmp / 真實 repo / 跨版本 / rollback 可逆）齊備
- [ ] Q1 / Q3 resolved 或 deferred 並回寫

## Glossary Impact
Glossary Impact: yes — 新增 `plan_profile`（capability / portable 邊界）與 `plan_schema`（frontmatter schema + version 相容契約），刻意拆成兩個單一責任術語；Phase 1 落地時註冊到 `knowledge/glossary/ai-skill.md`。

## 與其他 plans 的關係
- 依賴 [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) 的 5 validators。
- 依賴 [`archived/2026-05-28-1830-plan-archival-audit-validator.md`](../../archived/2026-05-28-1830-plan-archival-audit-validator.md) 的 archival audit。
