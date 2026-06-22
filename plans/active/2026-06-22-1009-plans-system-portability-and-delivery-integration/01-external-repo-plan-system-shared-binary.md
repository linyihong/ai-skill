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
  overlay 的邊界（哪些 validator 屬 plan_system_profile），02 / 03 才知道
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
- **In**：portable core vs governance overlay 邊界定義（`plan_system_profile`）；外部 repo 跑 plan-tree / archival validators 的跨 repo 機制；共用 binary 版本相容策略。
- **Out**：把整套 Ai-skill 治理（runtime.db / glossary / cognitive modes / ADR pipeline）搬進外部 repo；init-project 抽取安裝（保留為 future option）。
- **Affected**：`scripts/ai-skill-cli/internal/app/plans.go`、`plan_tree.go`、`hooks.go`、`plans/README.md`、新 `ai-tools/` 或 `scripts/ai-skill-cli/docs/` 外部 repo 使用說明。

## Decision Rationale（sub 層）
讀取面（`plans tree --root`）已跨 repo；唯一缺口是 **commit-time 強制**。現有 plan-tree validators 在本 repo 的 `commit-msg` hook 內呼叫 `runCommitMsgHook`，綁本 repo。外部 repo 兩種薄路徑：
- **(a) `ai-skill plans validate --root <path>` 子命令**：把 plan-tree / archival validators 從 commit-msg 內部邏輯抽成可獨立呼叫的 read-only validate；外部 repo 在自己的 git hook 或 CI 呼叫共用 binary。tool-neutral、不綁特定 hook。
- **(b) git hook shim**：外部 repo 安裝薄 `commit-msg` shim 呼叫共用 binary。
建議 (a) 為核心（tool-neutral、CI 也能用），(b) 為 convenience adapter（放 `ai-tools/`）。Phase 0 確認 validators 是否已可在不依賴 Ai-skill repo-local 狀態下執行。

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
| Q1 跨 repo 強制機制 | still-open | Phase 1 盤點 hooks.go 可抽性 |
| Q2 portable core 成員 | still-open | Phase 1 逐一分類 13 plan validators |
| Q3 版本相容 | still-open | Phase 2 決定 schema version 宣告方式 |

### Phase 0.1 — 架構盤點
- [ ] 讀 `hooks.go` `runCommitMsgHook`：plan-tree / archival validators 是否依賴 Ai-skill repo-local 檔（routing-registry / runtime.db）？
- [ ] 逐一分類所有 commit-msg validators → portable core（plan-tree 5 + archival 2）vs Ai-skill overlay（runtime trigger wiring / glossary / cli-doc-sync / runtime yaml projection 等）。

## Phase 1 — Portable core 邊界定義
- [ ] 產出 `plan_system_profile` 定義：列出 portable validators 清單 + 排除清單 + 各自 contract source。
- [ ] 文件化於 `plans/README.md`（或新 `governance/lifecycle/plan-system-profile.md`）。
- [ ] 完成條件：邊界清單 review 通過，Q2 標 resolved。

## Phase 2 — `ai-skill plans validate --root` 子命令
- [ ] 抽出 portable validators 為可獨立呼叫（read-only，不依賴 commit context）。
- [ ] 新增 `plans validate --root <path> [--format text|json]`，回傳 findings。
- [ ] 加 schema version 宣告（frontmatter 或 profile 檔）以支援版本相容（Q3）。
- [ ] 測試：tmp dir fixture（合法 tree + 各 violation）≥ 5 case。
- [ ] **若新增 `route.*` 或 runtime surface，補 Runtime Execution Path + Per-surface consumer 表**（否則明寫 doc/CLI-only，無新 route）。

## Phase 3 — 外部 repo 使用路徑 + git hook shim（convenience）
- [ ] `ai-tools/` 或 `scripts/ai-skill-cli/docs/` 寫外部 repo 使用說明（共用 binary 路徑、`plans validate` 接 CI / git hook）。
- [ ] 提供薄 `commit-msg` shim 範例（呼叫共用 binary，tool-neutral）。
- [ ] **Acceptance evidence**：在一個非 Ai-skill 的測試 repo（或 tmp fixture repo）真實跑過 `ai-skill plans validate --root`，截取 pass + fail 輸出。

## 完成條件
- [ ] portable core 邊界清單落地（Q2 resolved）
- [ ] `plans validate --root` 子命令 + 測試通過
- [ ] 外部 repo 使用說明 + shim 範例落地
- [ ] 至少一次外部 / tmp repo validate evidence
- [ ] Q1 / Q3 resolved 或 deferred 並回寫

## Glossary Impact
Glossary Impact: yes — 新增 `plan_system_profile`（portable core 邊界術語）；Phase 1 落地時註冊到 `knowledge/glossary/ai-skill.md`。

## 與其他 plans 的關係
- 依賴 [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) 的 5 validators。
- 依賴 [`archived/2026-05-28-1830-plan-archival-audit-validator.md`](../../archived/2026-05-28-1830-plan-archival-audit-validator.md) 的 archival audit。
