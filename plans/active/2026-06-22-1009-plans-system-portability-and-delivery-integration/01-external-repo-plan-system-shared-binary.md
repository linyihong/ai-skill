---
id: 2026-06-22-1009-external-repo-plan-system-shared-binary
plan_kind: sub
status: draft
owner: linyihong
created: 2026-06-22
parent: 2026-06-22-1009-plans-system-portability-and-delivery-integration
required_for_completion: true
sub_plan_reason: >
  本 sub-plan 切出 portable core vs Ai-skill governance overlay 的邊界
  （哪些 validator 屬 plan_profile）。此邊界**有助於** 03 後續判斷新欄位
  是否 portable（非硬依賴 — 03 可平行設計 schema，僅與本 plan 對齊 frontmatter；
  02 完全獨立）。獨立成 sub-plan 以便 graduate 邊界共識（Q2）與跨 repo
  強制機制（Q1/Q3），且其 acceptance（外部 repo 真實跑過 validate）可獨立 sign-off。
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
plan artifact (frontmatter + checkbox)
        ↓
schema compatibility layer    ← normalize / version-resolve（plan_schema 住這層）
        ↓ normalized model
validator engine              ← 核心，read-only，吃 ValidationContext，不吃 raw tuple
        ↓ consumers (薄)
  ├── commit-msg hook (本 repo 既有)
  ├── git hook shim   (外部 repo, 放 ai-tools/)
  ├── CI              (任何 repo, 直接呼叫)
  ├── plans validate  (CLI surface, 也只是 consumer)
  └── future API
```

**engine contract 不過早固定（回應 review #1）**：Q1（context 最小集合）尚未解，**不鎖 `(root, staged-set)` tuple**，否則 Q1 一解就破 API。改用可演化 struct：

```go
type ValidationContext struct {
    Root          string            // repo root
    ChangedSet    []string          // staged / changed files
    ExecutionMode string            // commit | ci | manual
    SchemaVersion string            // 由 compat layer 解析後填入
    Metadata      map[string]string // HEAD / working-tree state 等，演化用
}
```

**schema version 放對層（回應 review #3）**：`plan_schema` version **不屬 engine**，屬其上的 **schema compatibility layer**；engine 只吃 normalized model。否則 engine v4 ↔ schema v2 會綁死。

若把 `ai-skill plans validate` 當核心，半年後它會長成另一個 orchestration layer（CLI 累積 flag / 狀態 / 隱性 contract）。因此 CLI 只是其中一個薄 consumer。**（候選治理原則：Consumer Thinness Rule — surface = transport only / logic = engine；現不 promote，記為 deferred candidate，待 N≥5 consumer 再評估抽共用規則。）**

**portable 邊界（回應 review #2）：不預設「plan-tree 5 + archival 2 = portable」**。portable 不是看 validator 類型，而是看 contract → dependency → execution context。Phase 1 必須**先建分類模型再分類**，否則會變「先決定 portable 再找理由」。

## Open Questions（本 sub）
- Q1（跨 repo 強制機制）/ Q2（portable 邊界）/ Q3（版本相容）/ **Q7（validator failure semantics — severity 映射歸 engine contract）** — 見 main plan §Open Questions。Q7 於 Phase 1 facts 碰到 severity + opt-out 後正式立案。

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

### Phase 0.1 — 架構盤點 ✅（2026-06-22 完成）
- [x] 讀 `hooks.go` `runCommitMsgHook`：哪些 validators 依賴 Ai-skill repo-local 狀態（routing-registry / runtime.db / commit context），哪些只吃 (root, staged-set)？
- [x] 確認 validators 是否已可在不依賴 commit context 下執行（決定 engine 抽取成本）。

#### 盤點發現（survey findings）

| validator | input | 內容來源 | commit-context 耦合 | runtime.db / routing | 抽取分級 |
|---|---|---|---|---|---|
| `validatePlanTreeFrontmatter` | `(text, staged, root)` | `scanAllPlanFrontmatter(root)` walk **working-tree**（`os.ReadFile`） | 僅 opt-out marker 來自 `text` | 無 | **context-free core** |
| `validatePlanTreeArchiveOrder` | 同上 | 同上（filesystem walk） | 僅 opt-out marker | 無 | context-free core |
| `validatePlanTreeParentReference` | 同上 | 同上 | 僅 opt-out marker | 無 | context-free core |
| `validatePlanTreeUniqueID` | 同上 | 同上 | 僅 opt-out marker | 無 | context-free core |
| `validatePlanTreeFolderConvention` | 同上 | 同上 | 僅 opt-out marker | 無 | context-free core |
| `validatePlanArchivalAudit` | `(text, staged, root)` | `findArchivedPlans(staged)` + filesystem | opt-out **與 body-justification** 皆讀 `text` | 無 | core，但 text-coupled |
| `validatePlanArchivalLinkIntegrity` | `(text, staged, root)` | `readStagedFileContent` = `git show :<path>` **有 working-tree fallback** | **staged-blob 讀取**（可降級） | 無 | core，需 ExecutionMode 決定 staged vs worktree |

**三點結論：**

1. **engine 抽取低風險（解 Q1 大半）**：所有 plan validators 的核心邏輯已近乎 context-free（filesystem walk @ root + changedSet）。**無任何 plan validator 觸碰 routing-registry / runtime.db** → portable core 與 Ai-skill governance overlay（`route.*` / `validateRuntimeTriggerWiring` 等）在**依賴層就乾淨可分**，Q2 分類可行。
2. **僅兩處 commit-context 耦合，且正好對應 round-4 的 `ValidationContext` 欄位**：(a) link-integrity 的 staged-blob 讀取（→ `ExecutionMode` 決定 `git show` vs working-tree）；(b) opt-out marker / archival body-justification 讀 commit message（→ `Metadata` / `ExecutionMode`）。這兩處正是「不可鎖 `(root, staged-set)` tuple」的實證理由。
3. **opt-out transport 是最大跨 consumer 設計含意**：`[skip-plan-tree-*]` 與 archival body-justification 在 commit mode 來自 commit message；CI / manual mode **沒有 commit message**，opt-out 須改走 config / flag transport。Phase 1 分類表的 `consumer_surface` 欄需明確標示此差異，Phase 2 engine 須把 opt-out 解析下放到 consumer adapter，不寫死在 engine。

> Q1 在 §Open Questions 仍標 `open`（依關閉規則，需待 engine→consumer 抽象實際落地 + shim/CI 路徑文件化才可 close）；本盤點為其 Resolution Evidence 的前置證據。

## Phase 1 — Portable 分類（兩層：先 facts 後 decisions）

> **拆兩層（回應 review）**：不直接從一張表推 `portable`，否則 `consumer_surface=CI → portable=no` 這種錯誤會偷偷發生。Layer A 只寫**觀察事實**、不下結論；Layer B 才**推導**。
>
> **opt-out 治理句（這輪升格）**：`engine receives effective policy / consumer resolves policy source` — opt-out（`[skip-*]`）與 archival body-justification 是 **transport concern 不是 engine concern**：commit→message、ci→config、manual→flag；engine 只接收已解析的 effective policy，由 consumer adapter 解析來源。commit / ci / api 因此不互相污染。

### Layer A — Facts inventory（觀察，不下結論；Phase 0.1 已驗證 9 列）

| validator | contract_source | runtime_dependency | execution_context | consumer_surface |
|---|---|---|---|---|
| `validatePlanTreeFrontmatter` | plan-tree-hierarchy plan | none | working-tree walk + changedSet；opt-out 來自 text | hook, cli, ci |
| `validatePlanTreeArchiveOrder` | 同上 | none | 同上 | hook, cli, ci |
| `validatePlanTreeParentReference` | 同上 | none | 同上 | hook, cli, ci |
| `validatePlanTreeUniqueID` | 同上 | none | 同上 | hook, cli, ci |
| `validatePlanTreeFolderConvention`(warn) | 同上 | none | 同上 | hook, cli, ci |
| `validatePlanArchivalAudit` | archival-audit plan | none | filesystem；opt-out **+ body-justification** 來自 text | hook, ci(降級) |
| `validatePlanArchivalLinkIntegrity` | archival-link-integrity plan | none | **staged-blob (git show)** + worktree fallback | hook, ci, cli（ExecutionMode-gated） |
| `validatePlanCheckboxSync` | gen3 plan | none | **commit-message 解析 plan refs** + staged diff | hook only |
| `validatePlanStatusSync` | plan-status-sync-enforcement.yaml | none | **commit-message 解析 phase/refs** + staged | hook only |
| overlay — cognitive 家族（ExecutionModeFloors / GovernanceModeConsistency / MemoryModeSubdir / CognitiveCost / ActivationSignals / CapabilitySnippet / TokenBudget / AdaptiveTriggers / CognitiveContractFormat） | cognitive-modes*.yaml | **runtime.db**（讀 `generated_surfaces[runtime.cognitive_modes.discovery]` @ hooks.go:2934；fresh-clone fallback） | 解析 commit-message cognitive block + staged | internal only |
| overlay — runtime/routing wiring（RuntimeTriggerWiring / RuntimeYamlProjects / RuntimeIndexFreshness） | system-upgrade-governance / runtime | **routing-registry**（hooks.go:3624/3761）/ **runtime.db / runtime index** | staged yaml/registry + 投影狀態 | internal only |
| overlay — Ai-skill repo-structure（CLIDocSync / GlossaryRetroOwn / BootstrapEntryThinness / MarkdownYamlSync / EvidenceHierarchy / EnforcementRegistry*） | cli-modification-policy / glossary / enforcement-registry | repo-local 檔（CLI docs / glossary / enforcement registry / sibling yaml）；無 runtime.db | Ai-skill 專屬路徑結構 | internal only |
| overlay — safety（SanitizationStagedContent / NoNewShellScripts） | sanitization / go-first policy | metadata-derived forbidden tokens（Ai-skill 政策清單） | staged content scan | internal only |

- [x] **Layer A 殘留已驗證（2026-06-22）**：逐一 empirically 確認 overlay validators 的 runtime_dependency（不只信 contract_source）。結論：**無任何 overlay validator 驗證 plan-tree 結構**；cognitive 家族經 runtime.db、wiring 群組經 routing-registry/runtime.db、repo-structure 群組經 Ai-skill 專屬路徑——全部 `internal only`。**修正先前 Layer A 把 cognitive 家族標 `none` 的錯誤**（實際讀 runtime.db discovery surface）。無漏網的 plan 相關 validator。

### Layer B — Decisions（推導，引用 Layer A）

| validator | portable_candidate | profile | excluded_reason |
|---|---|---|---|
| 5 個 plan-tree validators | yes | `plan_profile.core` | — |
| `validatePlanArchivalAudit` / `LinkIntegrity` | yes（ExecutionMode-gated） | `plan_profile.archival` | body-justification / staged-blob 須由 ExecutionMode 決定 transport |
| `validatePlanCheckboxSync` / `validatePlanStatusSync` | **no** | — | **commit-message discipline**：解析 commit text，CI/manual 無 commit message 無法轉移（≠ plan 結構驗證） |
| governance overlay group | no | — | Ai-skill 治理 contract / 部分 runtime.db / routing 依賴 |

- [x] **`plan_profile` membership frozen（2026-06-22，pending consumer validation）**：
  - `plan_profile.core` = 5 plan-tree validators（frontmatter / archive-order / parent-ref / unique-id / folder-convention）
  - `plan_profile.archival` = archival-audit / link-integrity（ExecutionMode-gated：staged-blob vs worktree、body-justification transport）
  - **excluded**（evidence in Layer A）：checkbox-sync / status-sync（commit-message discipline，hook-only）、cognitive 家族（runtime.db）、wiring 群組（routing-registry/runtime.db）、repo-structure（Ai-skill 路徑）、safety（policy tokens）
  - `consumer_surface` 為獨立 execution 維度，**不併入** `plan_profile`；`plan_schema`（frontmatter schema + version）住 compat layer，非 engine。
  - **`plan_profile` FROZEN（回應 review）**：membership 凍結。**只有三種事件可重開**：(1) shadow 出現 `missing`、(2) shadow 出現 `extra`、(3) Phase 3 acceptance 被卡住。其餘一律當 observation。**禁止**因 Vidoe-Test 發現（folder-convention / loader / dialect）回頭擴 `plan_profile.core`。
- [ ] **canonical `governance/lifecycle/plan-profile.md` 暫不建（避免 premature canonical surface）**：frozen membership 先留在本 plan；待 Phase 2 engine 抽出 + **一個 consumer 成功跑**（Q2 close 條件）再 promote 成 canonical doc。符合 maturity ladder（observation → runtime）。
- [ ] **完成條件**：Layer A facts 完整（✅ 含殘留驗證）+ Layer B decisions review 通過（✅ membership frozen）+ canonical promote 待 consumer。**Q2 不在此 close**（見 §Open Questions 收緊後條件：尚需一個 consumer 成功執行）。

## Phase 2 — Engine 抽取（gated sub-phases 2.0→2.4）

> **Architecture Compatibility Gates（2026-06-22 review，動 Go 前必過）**：
> - **Gate A**：先立 `ValidationContext` contract，**不寫 engine / 不實作 method / 不解析 commit / 不碰 schema**；證明 hook / cli / ci 三邊都能餵。卡住 = engine contract 未成熟。
> - **Gate B**：schema compat layer 與 engine **分開（不同 commit）**；順序 schema loader → normalized plan model → engine。驗收：**engine 不知道 frontmatter version**。
> - **Gate C**：commit hook **不直接替換**，先 shadow path（legacy ‖ engine-shadow），比對 findings equality 至少 pass / fail / opt-out 三組（opt-out 最易被重構掉行為）。
> - **Gate D**：第一個 consumer **不是 CLI**，是 engine integration test（fixture → context → findings），再接 CLI（CLI 易把 transport 細節偷帶進 engine）。
>
> **Gate B 修正 round-4**：`SchemaVersion` **不放** `ValidationContext`（engine 不該知道 version）；改由 compat layer 解析成 normalized model。
>
> **package 位置**：新 `scripts/ai-skill-cli/internal/planvalidate/`（不在 `internal/app/` 下 → 不觸發 `validateCLIDocSync`）。
> **暫不碰外部 repo**（Phase 3 太早）。

### Phase 2.0 — `ValidationContext` contract（= Gate A）✅（2026-06-22）
- [x] 新 package `internal/planvalidate/`（`context.go`）：定義 `ValidationContext{ Root, ChangedSet, ExecutionMode, Metadata }` + 具名型別（`ChangedSet` / `ExecutionMode` enum: commit|ci|manual / `ValidationMetadata`）。**types only，無 method、無 parsing、無 schema、無 SchemaVersion 欄**（Gate B 修正已落實）。
- [x] 測試證明 **hook / cli / ci 三邊都能 construct**（`context_test.go` 三個 construction site，斷言欄位 + 各 mode 的 opt-out transport metadata）。
- [x] `go build ./...` 綠 + `go test ./...` 全 suite 綠（planvalidate 3 PASS，app/audit/glossary/pathutil 未受影響）。Gate A 通過：contract 可被三 consumer 餵。

### Phase 2.1 — schema compatibility layer（= Gate B）✅（2026-06-23）
- [x] schema loader（`schema.go`）：`RawPlan`（唯一帶 `SchemaVersion`）→ `Normalize()` → `NormalizedPlanModel`（version-free）。`plan_schema` version 解析**只**住此層。
- [x] **B.1**：engine-facing `NormalizedPlanModel` **無 version 欄**——reflection test `TestNormalizedPlanModel_HasNoVersionField` 機械鎖；grep `version` 僅命中 `Normalize` 區域變數。
- [x] **B.2**：所有版本判斷集中於 `Normalize`（unsupported version 在此 reject），validator 永不 `if plan.Version`。
- [x] **B.3**：雙向 fixture `TestNormalize_AbsentAndExplicitVersionProduceSameModel`——「無 `schema_version`（既有 plan）」與「顯式 `schema_version: "1"`」normalize 成同一 model。誠實起點：今天唯一真實相容邊界是 absent==baseline，非捏造 v2。
- [ ] **延後（記，不做）— CompatibilityResult{ Model, Warnings }**：未來出現 deprecated-but-tolerated 欄位時，須區分 **missing field vs deprecated field**，並回 warnings 讓 hook/CI/CLI 各自決定（不是 error-only）。目前 `Normalize` 回 `(model, error)`，warning channel 保留待第二個 schema version 落地再加（避免 over-build）。

### Phase 2.2 — engine + integration test（= Gate D；Q2 close 點）✅（2026-06-23）
- [x] engine（`engine.go`）：`Validate(ValidationContext, []NormalizedPlanModel) []Finding`。
  - **Gate D.1**：`Finding{ RuleID, Message, Blocking bool }` — **minimal，不開 severity enum**（修正先前「severity block/warn/info」草案；Q7 未成熟，transport policy 不寫進 engine）。
  - **Gate D.2**：collect-all，無 fail-fast（`TestEngine_Mixed_CollectsAllFindings` 一次回多個 rule class）。
  - 實作 4 條純 core 規則：frontmatter / unique-id / parent-reference / archive-order。
- [x] **第一個 consumer = engine integration test（no CLI，Gate D）**：`engine_test.go` **按行為分組（Gate D.3）**：valid / missing / broken_link / archive_required / opt_out / mixed。opt-out 證明為 **consumer-side**（engine 仍出 finding，consumer filter 套 effective policy）。
- [x] **Gate D.4 negative evidence**：`TestEngine_CannotExpressExcludedValidators` — `NormalizedPlanModel` 無 route/registry/commit/message/runtime/discovery/diff 欄，excluded validators（runtime-trigger-wiring、checkbox/status-sync）**結構上無法表達** → portable 邊界 by construction 成立。
- [x] **doc 補**：engine.go 明寫「Finding transport policy DEFERRED」，防後人見 `Blocking` 就加 Warning/Info。
- [x] **Q2 close**：分類 + `plan_profile` frozen + 首個 consumer（integration test）綠 + negative evidence → §Open Questions Q2 標 resolved。
- [ ] **supplemental input rule（非「殘留」，回應 review）**：folder-convention（warning）**不是 NormalizedPlanModel 不夠完整**，而是 **rule needs additional context**（dir listing）。**禁止把 filesystem 搬進 model**；未來若支援，走 `ValidationContext{ Root, ChangedSet, SupplementalInputs }`，不污染 model。archival 規則的 staged-blob vs worktree 同類，由 `ctx.ExecutionMode` + supplemental input 提供，待接 consumer 時補。

### Phase 2.3 — shadow hook（= Gate C）✅（2026-06-23）
- [x] `planvalidate.Compare(legacy, engine, hints) Divergence` + `planValidateShadowCheck` 接入 `runCommitMsgHook`（不替換 legacy；只 append 一個 informational Check）。
- [x] **Gate C.1**：shadow **不影響 exit code** — helper 只回傳 `Check`，從不 set `result.ExitCode/Status`；commit 結果維持 legacy-only。既有 app 測試全綠（行為未變）。
- [x] **Gate C.2**：equality 比 `RuleID + Blocking`，**忽略 Message**（`TestCompare_IgnoresMessage`）。
- [x] **Gate C.3**：divergence 分 5 桶 `same / missing / extra / transport / context`——opt-out 差異歸 **Transport**（engine policy-free 仍出 finding，legacy 經 `[skip-*]` 抑制）、staged/worktree/execmode 歸 **Context**；只有 missing/extra 算真 gap（`Converged()`）。
- [x] 行為測試覆蓋 pass(valid) / fail(genuine gap) / **opt-out→transport** / context 四情境。
- [x] **live evidence**：本 plan 的 commit 會觸發 shadow（plans staged + full-table），commit hook 輸出 `planvalidate_shadow: ok engine parity ...` 即首個真實 hook-context parity 證據。
- [ ] **收斂門檻移至 Phase 2.3b 觀察**（見下）。

> **Observation（記，不升 governance）**：2.3 真正建立的不是 shadow 功能本身，而是 **observability before replacement** ——在切換實作前先有 production-grade 比對通道。先記為 observation；N≥? 再考慮是否抽成 reusable pattern，本輪不 promote。

### Phase 2.3b — Shadow Confidence Window（觀察窗，非新功能）
> 不接 CLT、不切換 hook；只收集真實 commit 的 shadow 輸出，避免太早關掉觀測窗口。

- [ ] **證據量（最少）**：normal commit ≥3、真違規 commit ≥1、opt-out commit ≥1。
- [ ] **收斂門檻（唯一硬條件）**：`missing=0 ∧ extra=0`。**不要求 same=100%**：transport（opt-out 預期）、context（ExecutionMode 差異）允許非空。
- [ ] **收斂規則（強化，回應 review）**：transport / context 的每個 entry **必須可解釋**，不得「非 genuine 就自動忽略」。機械保證：`Compare` 只在 hint（`OptedOut` / `ContextSensitive`）命中時才歸 transport/context，無 hint 一律 genuine（`compare_test.go` `GenuineGaps` / `OptOutBecomesTransport` / `ContextBucket` 已鎖）。
- [ ] **不人造 commit 湊數**：靠自然 commit 流量。**14 天 fallback（回應 review）**：若自 2.3b 開啟（2026-06-23）起 **14 天內（至 2026-07-07）** 自然流量仍未出現 violation / opt-out，**允許 fixture replay**（不進主線）補齊樣本——目的是讓窗口能關，不是湊數。fixture replay 仍屬 observability（shadow），非行為驗證。
- [ ] **觀察記錄**：累積 shadow Check 輸出（commit hash + 5 桶）於本節，達門檻後才允許評估進 2.4。
- 進度：normal #1 = commit `403fa73`（`same=- missing=- extra=- transport=- context=-`，valid plan parity）。

### Phase 2.4 — CLI consumer
> **Gate：須先過 Phase 2.3b 收斂門檻**（`missing=0 ∧ extra=0` + 證據量達標）才開工，避免太早關掉 shadow 觀測窗口。
- [ ] CLI `plans validate --root <path> [--format text|json]` 作為薄 consumer（transport only）。
- [ ] **若新增 `route.*` 或 runtime surface，補 Runtime Execution Path + Per-surface consumer 表**（否則明寫 engine/CLI-only，無新 route）。

**Q-close 映射**：Q2 → Phase 2.2 後可 close；Q1 → Phase 3；Q3 → 跨版本 evidence。

## Phase 3 — 外部 repo consumer 路徑（git hook shim / CI）

### External Pressure Samples（壓力樣本，**非** acceptance samples）
> Vidoe-Test 定位為 **pressure sample**（揭露邊界），不是 acceptance sample。本表記 **coverage**，**不寫 pass/fail**——這不是 correctness 問題。

| Repo | Loader | Schema Match | Engine Coverage | Notes |
|------|--------|-------------|-----------------|-------|
| Ai-skill | native | full | 4/4 | baseline |
| Vidoe-Test | read-only throwaway harness（**measured 2026-06-24**，未 commit） | partial | 13/14 parsed（README 跳過），normalizeErr=0；**2/4 rules exercised** | **實測修正先前「1/4」預測**。`unique_id` 跑（0 findings）；`parent_reference` 跑出 **2 findings**。schema 欄位實測：plan_kind=0 / parent=**2** / required_for_completion=0 / sub_plan_reason=0 → frontmatter & archive_order inert（欄位缺）。**inert 來源已歸因 = schema/semantics，非 loader**（loader 0 error 解析 13 個）。 |

**Boundary finding**：**portable ≠ schema-agnostic**。01 目前完成的是 `shared engine + portable profile + same plan contract`，**不是** universal plan language / interoperability framework。

**新證據（measured，回應「能否按我們流程跑」）**：Vidoe-Test 架構**機械上跑得動**（loader 解析 13/13、engine 執行），但輸出含 **2 個 false-positive**：其 `parent:` 欄持有 **path**（`docs/plans/xxx.md`），我們 engine 契約期待 **id** → 判不解析。這是 **semantic mismatch on a shared field name**（同名 `parent`、語意不同），**不再只是 absence**。意涵：
- 對 Q8——上輪「只有 absence、無 semantic mismatch，證據不足以支撐 normalize」的前提**已被此證據更新**；mapping 不再因「缺證據」被排除，但 **a/b/explicit-unsupported 決定仍 deferred-to-phase-3**，本輪不決、不做 mapping。
- 此為 **external throwaway 量測**，**非 shadow** → **不觸發 `plan_profile` reopen**（reopen 只認 shadow missing/extra 或 Phase 3 blocked）。`plan_profile` 維持 FROZEN，純當 Q8 observation。

#### Canonical-authored external sample（acceptance-leaning，measured 2026-06-24）
外部團隊改用**我們 schema 結構**新寫 `Vidoe-Test/docs/plans/2026-06-22-1600-h5-redis-read-cache.md`（採 `schema_version: "1"`、`status: draft`、完整必填 sections：Decision Rationale 全子節 / Open Questions / Phase 0.0 公版 / Stakeholder / 完成条件 / Runtime Execution Path）。

| 性質 | 實測 |
|---|---|
| 解析 | 成功（dir 15 files → 14 parsed，README 跳過），normalizeErr=0 |
| 此 plan 的 findings | **0（乾淨、無 false-positive）**——無 top-level `parent:`（用 nested `upstream` id+path），故不像 2 個舊 path-parent plan 觸發假陽性 |
| 規則覆蓋 | 單一 main plan → frontmatter/parent/archive **正確 inert**（applicability，非 mismatch）；unique_id pass |
| 定位 | **acceptance-leaning**（不再只是 pressure sample）：首個外部 repo 以 canonical 結構撰寫、engine 乾淨驗過 → 支撐 Q8 的 **adoption** path 可行 |

**誠實 caveat（不可誇大）**：
1. 本次 throwaway **未把 `schema_version` 接進 `Normalize` 的 version 路徑**（只當 absent→baseline 正規化）→「compat layer 接受顯式 `schema_version: "1"`」**尚未機械證明**，留待真實 loader（2.4 / Phase 3）。
2. 觀測到 raw 值為 `"1"`（**含 YAML 引號**）。真實 loader 必須 strip quote（`"1"` ↔ `1`），否則對 `currentSchemaVersion` 比對會誤拒——**記為 Q3 / loader 需求**。
3. 仍是**單一 main plan**，未涵蓋 sub-plan tree 規則；完整外部 coverage 需一份外部 **plan tree（main+sub）** 以 canonical schema 撰寫。
4. 仍 throwaway、非 shadow → 不動 `plan_profile` FROZEN、不提前決 Q8（此為 adoption-path 正向證據，非強制決策）。

#### Canonical external **plan tree** + caught fidelity bug（DECISIVE，measured 2026-06-24）
外部團隊把 h5 plan **升級成 plan tree**：`_plan.md`（main, `parent: null`）+ 4 sub（`plan_kind: sub` / `parent`=main id / `required_for_completion: true` / `sub_plan_reason` / `schema_version "1"`）。補上先前缺的「外部 plan tree」coverage。

- **首跑抓到 1 個 false-positive**：main 的 `parent: null` 被當字串 "null" → `parent_reference` 判不 resolve。**真 engine fidelity bug**（legacy 用 `HasParentField`/null 處理，engine port 漏了）。
- **依 Gate B 修在 compat layer**：`Normalize` 加 `normalizeNullScalar`（null/~/Null/NULL → 空），engine 維持不看 YAML idiom；加測試 `TestNormalize_ParentNullIsEmpty`（含 tree-level Validate）。
- **修後重跑**：external tree `parsed=5 findings=0` → **CLEAN**。4 條規則全被真實外部 tree 觸發（4 sub frontmatter pass、4 parent→main id resolve、unique_id pass、archive_order inert）。
- **意義**：目前最強的 **Phase 3 acceptance-style 證據**——真實外部 repo、canonical schema、完整 tree、engine 乾淨驗過；且外部真實資料**當場抓到 replacement fidelity bug**（正是用外部素材的價值）。bug 屬 compat-layer normalization 家族（與 `"1"` 引號同類）。
- **護欄**：此為 compat-layer **correctness fix**，非 `plan_profile` membership 變更 → FROZEN 不受影響；Q8 仍 deferred。

- [ ] `ai-tools/` 或 `scripts/ai-skill-cli/docs/` 寫外部 repo 使用說明（共用 binary 路徑、engine 接 CI / git hook）。
- [ ] 提供薄 `commit-msg` shim 範例（呼叫共用 binary，tool-neutral）。
- [ ] **Acceptance evidence（回應 review #6，收緊）**：
  - [ ] tmp fixture repo：engine + CLI pass / fail 輸出
  - [ ] 一個**真實的非 Ai-skill repo**：實裝 shim，真實 commit 觸發一次 pass + 一次 block。**記錄維護中繼資料（回應 review #5，不寫 repo 名）**：`repo_owner` / `repo_type: internal|public|fixture` / `removal_policy`，避免一年後 repo 不存在無法追溯。
  - [ ] **跨 binary 版本驗一次**：升一次共用 binary（或改一次 `plan_schema` version），確認外部 repo 相容行為符合 Q3 策略
  - [ ] **rollback evidence（回應 review #6）**：外部 repo 可移除 integration（remove hook shim + config）並恢復 clean，且 **no schema residue** — 移除後 plan frontmatter 不需 migration、檔案格式不被永久改變，證明接入非侵入、完全可逆

## 完成條件
- [ ] portable 分類表（含 `consumer_surface` 欄）+ `plan_profile`（capability）/ `plan_schema`（compat）邊界落地（Q2 resolved，由分類表推導，capability 與 execution 維度分離）
- [ ] validator engine（吃 `ValidationContext`）+ schema compat layer + thin consumers（hook / CLI）+ 測試通過（含 no-CLI engine integration test）
- [ ] 既有 commit-msg hook 行為不變（重構回歸驗證）
- [ ] 外部 repo 使用說明 + shim 範例落地
- [ ] Acceptance evidence 四項（tmp / 真實 repo＋維護中繼資料 / 跨版本 / rollback 可逆＋no schema residue）齊備
- [ ] Q1 / Q3 resolved 或 deferred 並回寫

## Glossary Impact
Glossary Impact: yes — 新增 `plan_profile`（capability / portable 邊界）與 `plan_schema`（frontmatter schema + version 相容契約），刻意拆成兩個單一責任術語；Phase 1 落地時註冊到 `knowledge/glossary/ai-skill.md`。

## 與其他 plans 的關係
- 依賴 [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) 的 5 validators。
- 依賴 [`archived/2026-05-28-1830-plan-archival-audit-validator.md`](../../archived/2026-05-28-1830-plan-archival-audit-validator.md) 的 archival audit。
