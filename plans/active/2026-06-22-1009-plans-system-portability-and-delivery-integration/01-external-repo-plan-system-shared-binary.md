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

### Phase 2.3b — Shadow Confidence Window ✅（closed 2026-06-24 via Stage-2 replay）
> **Reframe（關鍵）**：2.3b 驗的是 **shadow accounting path（transport / context / parity-on-violation）被跑過**，**不是**收集世界事件。completion wording = **「transport/context path exercised」**，非「observed violation」。

**收窗規則（兩階段）**：
- **Stage 1（短自然窗）**：窗口開啟但設上限 = **14 天（至 2026-07-07）或下一個 Phase 開始，先到者為準**；來源限真 commit（violation ≥1 / opt-out ≥1）。**不無限等**。
- **Stage 2（超時 → fixture replay 收尾）**：只補 transport/context（+ parity-on-violation），**不重驗 correctness**。
  - **禁止**：人工製造 extra/missing divergence（那是 engine-challenge test，非 shadow confidence）。

**本輪以 Stage-2 fixture replay 收尾（Stage 1 僅得 normal parity，依「不變長期觀測」直接收）**：
- [x] **Replay A — violation parity**：sub `parent: ghost` → legacy 與 engine **都 fire** → `same=plan_tree.parent_reference`、missing/extra 空、exit 不變。`TestPlanValidateShadow_ViolationParity`（直接測 `planValidateShadowCheck` 整合，先前未測）。
- [x] **Replay B — opt-out transport**：同違規 + `[skip-plan-tree-parent-reference]` trailer → legacy 抑制、engine 仍 emit → `transport=plan_tree.parent_reference`、非 extra、converged。`TestPlanValidateShadow_OptOutTransport`。
- [x] **context path**：由 `compare_test.go TestCompare_ContextBucket` 覆蓋（ExecutionMode 尚無 rule 觸發，故以 Compare 層 exercise）。
- [x] **收斂門檻達成**：`missing=0 ∧ extra=0`（5 筆 normal parity-ok commit + 兩 replay 皆 converged）。transport/context entry 皆 hint-explained（無自動忽略）。
- 自然窗 datapoints（normal，全 parity）：`403fa73` / `4ae796b` / `36aa035` / `be7d8b8` / 後續 plan commits，皆 `missing=- extra=-`。

→ **2.3b 收窗完成**，accounting path 已 exercised。**解除 2.4 gate**。

### Phase 2.4 — CLI consumer ✅（2026-06-24，scope A）
> **Gate ✅ 已解除**：2.3b 收窗完成。**Scope A 鎖定**：CLI = transport surface，非 external interoperability。

- [x] CLI `plans validate --root <path> [--format text|json]`（`plans.go runPlansValidate`）：薄 consumer，呼叫 engine entrypoint `planvalidate.Validate`，**零驗證邏輯**。
- [x] **Discovery 固定 `<root>/plans/active|archived`**，reuse `scanAllPlanFrontmatter`（抽共用 `normalizedPlansFromRoot`，shadow + CLI 共用，**無新 traversal abstraction**）。
- [x] **Explicit non-goals**（未做、明文排除）：custom plans dir / schema dialect / loader plugins / external path convention / filtering / policy → 屬 Q8 / Phase 3。
- [x] Output `text | json`；blocking finding → exit 30（manual consumer transport），否則 0。
- [x] **Acceptance（CLI 不碰 validator internals，只呼叫 engine entrypoint）**：`TestPlansValidateCLI_MatchesEngineEntrypoint`（CLI json findings == `planvalidate.Validate` projection）、`_ValidTreeExitZero`、`_ViolationThreeWayEquivalence`。
- [x] **完整閉環達成（same repo / same tree / same findings）**：同一 canonical tree 上 **hook(legacy) ≡ engine(shadow) ≡ CLI(manual)** findings 等價——`ViolationThreeWayEquivalence` 驗 legacy validator 與 CLI 同 fire `parent_reference`；shadow 已驗 legacy≡engine。**hook / ci-able / manual 三 consumer → 同一 engine**。
- [x] command-contract.md 同步（`plans tree` + `plans validate`）；無新 `route.*` / runtime surface（engine/CLI-only，免 Runtime Execution Path 表）。
- smoke：`plans validate --root <Ai-skill>` → plans=30 findings=0。

### Phase 2.5 — Soak checkpoint（非 phase、不開 plan、不寫 shim）
> Phase 2 已回答「能否抽出可重用 engine」。**Phase 3 是換問題**：「another repo can adopt this without importing governance」——非自然延伸。故先讓已完成的抽象**活幾天看有沒有反噬**，不缺證據，缺 soak。

- **時間**：3–7 天 或 直到下一次真實使用需求（先到為準）。
- **觀察**：
  | Signal | 想看什麼 |
  |---|---|
  | CLI 使用 | `plans validate` 有沒有真的取代人工檢查 |
  | Shadow | 還有沒有 unexpected transport/context（missing/extra 必須維持 0） |
  | Engine API | 有沒有 consumer 想要求 engine 加 flag |
- **退出條件（硬）**：soak 期間 **沒有新增 engine surface**。若出現「CLI 想加 policy」或「hook 想特判」→ **先修 consumer，不碰 engine**。
- soak 平穩 → 才開 Phase 3（需先過 Phase 3.0 preflight，見下）。

**Soak observations（持續記）**：
- **2026-06-25 — 第 2 個外部 repo 真實採用（positive signal）**：`unwrapping/apk-analysis-sdk/docs/plans/integration/2026-06-24-1748-sdk-platform-sqlite-identity-pool`（1 main + **8 sub**，全 canonical schema）engine 量測 `parsed=9 findings=0` → **CLEAN**。**未觸發任何 engine surface 變更需求** → soak 退出條件維持。第二個獨立外部 repo 採 canonical schema 即乾淨通過，adoption-pass 由 N=1 → **N=2**。（throwaway 量測，非 shadow；plan_profile FROZEN、Q8 deferred 不變。）

**Soak verdict（2026-06-25）✅ 收尾**：engine stability = **pass**（退出條件「no new engine surface」達成：2 外部 repo、更大 tree、0 findings、0 engine requests——穩定性訊號，非 correctness 證明）。
- **關鍵拆分（這輪最大收穫）**：現在證明的是 **engine portability**，**不是** delivery portability。**不可把 engine success 當 externalization success**。
  - **已證**：engine + schema-compat + 三 consumer 閉環 + 2 外部 repo canonical 乾淨。
  - **未證（屬 delivery）**：shim / install path / invocation contract / upgrade(version) behavior / rollback。
- adoption-pass ×2 + dialect-pressure ×1 **仍不能推出 `choose adoption`** → **Q8 維持 deferred**（理由反而更強）。
- → **允許開 Phase 3.0 preflight（定義 only）**；**不**直接做 Phase 3 impl；**不**關 Q1/Q3。

### Phase 3.0 — Preflight ✅（PASSED 2026-06-25，**定義 only，禁 impl**）
> 範圍：**只定義契約**。**禁止**：git hook shim implementation / CI implementation / binary packaging / daemon / service / background sync。先回答「**What exactly is being externalized?**」。
> **Review verdict（2026-06-25）**：contract PASS、externalized object 定義正確（validation capability + invocation contract，非 whole governance）。Consumer contract / Integration shape / Rollback 三項 PASS（含 review 收緊：no persistent state、remove binary reference、三面 clean）。可開始寫 **Phase 3 impl plan（不是 code）**；**未授權** shim/CI 落地；Q1/Q3 不提前關；Q8 不碰。先補 §Phase 3.1 再開工。

**0. What is externalized（邊界）**
- **Externalized**：(a) validation **engine**（`planvalidate` package：`ValidationContext` + `Normalize` + `Validate`）、(b) **invocation contract**（如何呼叫），**經共用 binary**（`ai-skill plans validate`）——**非 vendored code、非 governance**。
- **NOT externalized**：governance overlay、runtime.db、routing-registry、cognitive modes、commit-msg 治理 validators。
- 即：外部 repo 取得的是「**一個會驗 plan 結構的可呼叫 engine**」，不是「一套治理系統」。

**1. Consumer contract（每個 consumer 必須滿足）**
- **input**：`ValidationContext{ Root, ChangedSet, ExecutionMode(commit|ci|manual), Metadata }` + plan set（scope A：`<root>/plans/active|archived`）。
- **output**：`[]Finding{ RuleID, Message, Blocking }`（engine 唯一輸出；無 severity enum，Q7 deferred）。
- **opt-out transport**：consumer 解析來源（commit msg / config / flag）→ 傳 **effective policy**；engine policy-free（「engine receives effective policy / consumer resolves source」）。
- **exit semantics**：blocking finding → consumer 自決 transport（hook→block / ci→fail / manual→exit 30）。engine 不決定 exit。
- **integration shape**：consumer = **transport only**，呼叫 engine entrypoint，**零驗證邏輯**（Consumer Thinness）。

**2. Integration shape（允許形態）**
- 允許：**git hook shim**（外部 repo commit-msg 呼叫共用 binary）、**CI wrapper**（CI step 呼叫共用 binary）、manual CLI（已存在）。
- 禁止：daemon / service / background sync / 任何常駐。
- **no persistent installation state（review 補）**：不得有 local db / cache daemon / managed worker。外部 repo 生命週期必須是：`clone → add shim → validate → remove shim → disappear`（裝/卸都不留狀態）。

**3. Rollback criteria（adoption 前必先證）**
- `remove shim + remove config + **remove binary reference**（review 補）→ repo returns clean`。驗三面乾淨：**git status clean + runtime clean + hook clean**（rollback 最易漏的是 **toolchain residue**）。
- **no schema residue**：integration 不得改 plan frontmatter / 檔案格式（移除後不需 migration）。
- uninstall = 單一可逆步驟。

- [x] 0/1/2/3 契約定義完成 + review 收緊（本輪，定義 only）。
- [ ] **（Phase 3 impl，禁於 3.0）**：shim / CI wrapper / cross-version / 真實外部 repo acceptance + rollback proof 執行 → 關 **Q1（跨 repo 強制機制）** + **Q3（跨版本）**。

### Phase 3 — Success Contract（**完成唯一標準**；非 slice、非 implementation）
> **成功的定義**（review 新增，唯一一句）：success **≠** `external repo passed`；success **=** `external repo adopted AND removed without residue **with preserved validation semantics**`。賣的是 **reversible adoption**，不是 works-on-my-repo。此標準於 slice **3.4** 一次完整 realize。

- **Acceptance 必含四段（缺一不算過）**：
  1. **install** — 加 invocation adapter（無持久狀態）。
  2. **validate** — 真實 commit 觸發、findings 行為符合 consumer contract。
  3. **upgrade once** — 升一次共用 binary（或 `plan_schema` version），相容行為符合 Q3 策略。
  4. **rollback** — 移除 adapter+config+binary reference → git/runtime/hook 三面 clean、no schema residue。
- **+ preserved validation semantics（review 補）**：upgrade 前後 **findings 的 meaning 不得變**。驗收明確拆**三層**（防「RuleID 沒變就說 semantics preserved」但 Blocking 已變）：
  - **Structural**：same RuleID set。
  - **Behavioral**：same Blocking result。
  - **Policy**：same opt-out interpretation。
- **+ removal must be monotonic（review ④）**：`remove adapter → no new validation errors → only capability disappears`。rollback 不得反而留下新錯（preserved semantics 的逆向版）。
- 四段 + 三層語意保真 + monotonic removal，缺任一 → Phase 3 **不算完成**，Q1/Q3 **不得 close**。

**Q-close 映射（evidence-slice 版）**：Q2 ✅（Phase 2.2）；**Q3 → slice 3.2（upgrade once）**；**Q1 → slice 3.3（consumer equivalence）**；Phase 3 complete → slice 3.4（四段 acceptance）；Q8 → deferred。

**Q-close 映射**：Q2 → Phase 2.2 後可 close；Q1 → Phase 3；Q3 → 跨版本 evidence。

## Phase 3 — 外部 repo consumer 路徑（git hook shim / CI）
> **Gated**：須先 (1) Phase 2.5 soak 平穩（無新增 engine surface）+ (2) Phase 3.0 preflight ✅ 通過 + Success Contract 補上才開工。Phase 3 換的是新問題「adopt without importing governance」，非 Phase 2 自然延伸。**按 evidence slice 切（見下 impl plan）**：Q3 → slice 3.2、Q1 → slice 3.3；Q8（external schema policy）仍 deferred。**impl 未授權落地**。

### External Evidence（3 個獨立 bucket，**不可互相污染**）
> 證據分三桶，**各自獨立、不可推導彼此**（回應 review：避免「adoption-pass → adoption selected」偷跑）：
>
> | Bucket | Meaning | 來源 | 掛點 |
> |---|---|---|---|
> | **adoption-pass** | 外部採 canonical **schema** → engine 乾淨驗（**caveat**：經 path-override EXT_PLANS_DIR=docs/plans 量測；scope-A CLI 需 `plans/active\|archived` layout，外部 repo 用 docs/plans → CLI `plans=0`，見 §3.4a location-convention gap） | **2 external repos**：Vidoe-Test h5 tree（6-node）+ apk-analysis-sdk identity-pool tree（9-node），皆 findings=0 | **Phase 3 acceptance evidence anchor**（**非** Q8） |
> | **dialect-pressure** | 非 canonical metadata → 觀察到 mismatch | Vidoe-Test flat plans（semantic mismatch：parent path vs id） | Q8 證據（揭露邊界，**不**決策） |
> | **compatibility-policy** | adoption / normalization / explicit-unsupported 三選一 | — | **Q8 = deferred**（尚未決） |
>
> 三桶現況：adoption evidence = yes／dialect evidence = yes／**policy decision = no**。`adoption-pass` 只證明「一個 branch 可行」，**不等於**「該選 adoption」——policy 仍 deferred。
>
> **排除**：runtime-index refresh（commit `612909c`）屬 **environment maintenance**，**不計入** 01 externalization evidence（避免污染進度）。
>
> **全 corpus 評斷（遞迴整個 `docs/plans`，measured 2026-06-24，throwaway 未 commit）**：files=20 → parsed=19（README 跳過）、normalizeErr=0。**canonical=6**（h5 tree：1 main + 5 sub）／**dialect=13**（flat plans）。全 corpus **findings=2，全部來自 dialect 的 path-parent**；canonical 6 個全 clean。
>
> 下表 coverage（記 coverage，不寫 pass/fail）：

| Bucket | Plans | Loader | Engine 結果 |
|--------|-------|--------|------------|
| baseline | Ai-skill native | native | 4/4 rules |
| **adoption-pass**（canonical） | Vidoe-Test h5 tree：6（1 main+5 sub） | read-only harness | **findings=0**；4 rules exercised（sub frontmatter pass、parent→main id resolve、unique_id、archive_order inert） |
| **dialect-pressure** | Vidoe-Test flat：13 | read-only harness | **findings=2**（2 plan 的 `parent` 持 path → semantic mismatch false-positive）；其餘 inert（plan_kind/required/sub_reason 欄缺）。inert 來源 = schema/semantics，非 loader |

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

#### [bucket: adoption-pass] Canonical external **plan tree** — Phase 3 acceptance anchor（measured 2026-06-24）
> **掛點 = Phase 3 acceptance evidence anchor，非 Q8**。意義限定為「canonical schema 這條 branch 可行（engine 直接乾淨驗過）」，**不**蘊含「該選 adoption」——compatibility-policy 仍 Q8-deferred。不因此重開 `plan_profile`、不提前決 a/b。

外部團隊把 h5 plan **升級成 plan tree**：`_plan.md`（main, `parent: null`）+ **5 sub**（`plan_kind: sub` / `parent`=main id / `required_for_completion: true` / `sub_plan_reason` / `schema_version "1"`；含後加的 `05-load-stress-redis-cache`）。補上先前缺的「外部 plan tree」coverage。

- **全 corpus 重跑（遞迴）**：6-node tree `findings=0` → **CLEAN**（5 sub frontmatter pass、5 parent→main id resolve、unique_id pass、archive_order inert）。全 `docs/plans` 19 parsed 中，**唯二 findings 來自 dialect path-parent，canonical 6 個全綠**。這是目前最強的 **Phase 3 acceptance-style 證據**：真實外部 repo、canonical schema、完整 tree、engine 乾淨驗過。
- **throwaway 量測的 false-positive（歸因更正，不可誇大）**：首跑時 main 的 `parent: null` 被判不 resolve——但這是**我 throwaway naive parser 的產物，非真 pipeline bug**。app 的真實 parser（`plan_tree.go:147`）**已**把 `parent: null/~/""`→`Parent=""`，所以 **shadow pipeline 從不受影響**。（草稿一度誤寫成「engine 漏了 legacy 有的東西」，已更正。）
- **fix 保留為 defense-in-depth（非 bug 修補）**：`Normalize` 加 `normalizeNullScalar`，讓 engine **不依賴任一 loader 的 null-awareness**（未來真實外部 loader 若 naive 也安全），與 app 一致、Gate B 不讓 engine 看 YAML idiom；測試 `TestNormalize_ParentNullIsEmpty`。屬 compat-layer normalization 家族（與 `"1"` 引號同類）。
- **護欄**：compat-layer 變更，非 `plan_profile` membership 變更 → FROZEN 不受影響；Q8 仍 deferred。

### Phase 3 impl plan — **evidence slices**（plan-only，**未授權落地** shim/CI/真實 repo）
> **🔒 骨架 FROZEN（2026-06-25，review 通過）**：不再拆 slice；evidence flow 不得切回技術分層。
> **切法（回應 review）：按 evidence stage 切，不按技術元件（shim/CI/wrapper）切**——避免 transport 反客為主、升格成 architecture。shim/CI 只是**證據載體（invocation adapter，replaceable，contract 不入 engine）**，掛在 slice 底下。
>
> **Phase 3 Non-goal（鎖死，外部化最易偷長）**：**不得新增 external registry、不得新增 external runtime state、不得新增 external plan metadata**。外部 repo 取得的永遠是「可呼叫的 engine」，不是一套新狀態。

| Slice | 目標 | 產物（載體） | 關閉 |
|---|---|---|---|
| **3.1 Adoption Slice**（install → validate） | 證外部 repo **可採用且不留狀態** | **invocation adapter**（git hook shim / CI wrapper / equivalent，**必須 replaceable；adapter contract 不得寫入 engine**）、install/remove 指南、validate evidence | reversible adoption 成立、rollback clean 成立；**不關 Q3** |
| **3.2 Compatibility Slice**（upgrade once）✅ DONE | 證 invocation contract **穩定非碰巧**（consumer compatibility，非 fixture） | **1 axis × 1 subject** upgrade（禁混升）、supported→preserved（三層）、unsupported→**deterministic + diagnosable** reject（same stage+reason class）、**no-change baseline**（negative proof）、3-axis+subject notes | **Q3 ✅ CLOSED**（見 §Phase 3.2 Preflight） |
| **3.3 Consumer Equivalence**（見 §Phase 3.3 Preflight matrix）| 同 repo 同 tree：**manual ≡ hook ≡ CI**（比 observation boundary 非 execution trace）；**equivalence ≠ coupling ≠ execution identity** | 依 Equivalence Matrix 驗（RuleID/Blocking/opt-out/discovery-scope MUST equal；exit/mode/input-snapshot MAY differ）、**removal-independence proof**、**≥1 consumer replaced without engine change** | **Q1 close（收緊，見 Q1 row）** |
| **3.4 Real Repo Acceptance**（orchestration；**3.4a selection 🟢 authorized / 3.4b execution ⛔ NOT YET**） | 完整四段 **install → validate → upgrade → rollback (→ remove)** on a **persistent** real repo（**MUST persist，禁 throwaway**；需時間維度）；先過 **3.4 Entry Gate** | acceptance evidence、rollback proof、residue check、acceptance metadata | **Phase 3 complete** |

**原 checklist 收編對照**：shim design → 3.1；CI wrapper → 3.1 / 3.3；cross-version → 3.2；rollback proof → 3.4；real repo acceptance → 3.4。真實外部 repo 維護中繼資料（`repo_owner` / `repo_type` / `removal_policy`，不寫 repo 名）掛 3.4。

**好處**：(a) transport 不升格 architecture；(b) Q-close 更乾淨——**Q3 在 upgrade 完即關（3.2），不等真實 repo**；**Q1 在 consumer equivalence 關（3.3）**；(c) `install→validate→upgrade→rollback` 四段 acceptance **只出現一次（3.4）**，前面都是鋪路。

- [x] **3.1 Adoption Slice — STARTED & demonstrated（authorized 2026-06-25）**（見下）。
- [ ] **（3.2–3.4 仍 plan-only，未授權落地）**：upgrade / consumer-equivalence / 真實 repo 四段**待授權**才寫 code。

#### Phase 3.1 — Adoption Slice evidence（2026-06-25）
- **Preflight finding（go-first 相容）**：`validateNoNewShellScripts` 擋任何新 `.sh` → **adapter 以 documented template 交付（guide 內 shim snippet），不在 Ai-skill commit `.sh`**；外部 repo 自行 materialize。同時滿足 no-new-shell / no persistent state（Ai-skill 端）/ replaceable / contract 不入 engine。
- **產物**：[`scripts/ai-skill-cli/docs/external-plan-validation.md`](../../../scripts/ai-skill-cli/docs/external-plan-validation.md)——invocation contract + adapter（git hook shim / CI wrapper）+ install/remove 指南 + non-goals。
- **Fixture evidence（真實 external git repo，throwaway）**：install 薄 commit-msg adapter（呼叫共用 binary `plans validate --root`）→
  - **validate PASS**：valid canonical tree `plans=2 findings=0 blocking=0`，commit exit 0。
  - **validate BLOCK**：注入 `parent: ghost` 違規 → `[BLOCK] plan_tree.parent_reference ... "ghost" does not resolve`、`blocking=1` → commit **aborted**。
  - **monotonic removal**：移除 adapter 後同一違規 commit exit 0 → **只有 capability 消失、無新 validation error**（符合 Success Contract removal-monotonic）。
  - **no residue**：hook 消失、plan frontmatter 未被改寫（no schema residue）、`git status` clean。
- **3.1 = CLOSED**：reversible adoption 成立、rollback clean 成立、monotonic removal + no residue 成立。**Q3 OPEN**（upgrade-once 屬 3.2）。
- **邊界（review 補）**：此證據是 **fixture compatibility**，**不是 consumer compatibility**——同一機器、同一 binary、throwaway repo。真實長期 consumer + 四段完整在 3.4。此差別直接影響 3.2 的歸因（見 3.2 Preflight）。

#### Phase 3.2 — Compatibility ✅ DONE（impl landed 2026-06-25，commit `2c26f6e`；Q3 CLOSED）
> **impl evidence**：compat layer 改 `supportedSchemaVersions` set（可擴充，{1,2} shape-identical 證 extensible 不破 findings）+ `CompatError{Stage,ReasonClass}` typed reject；**端到端 wire**：`PlanFrontmatter.SchemaVersion`（parser 已 strip quote → Q3 引號需求滿足）→ `normalizedPlansFromRoot` 帶進 `RawPlan.SchemaVersion`，unsupported 轉成 **blocking `compat.unsupported_schema_version` finding**（非默默降級）。測試：supported 1→2 preserved、unsupported deterministic+diagnosable reject、no-change baseline、CLI end-to-end（v2 clean / v99 exit 30 + blocking）。subject=artifact、axis=schema、單軸單 subject。

#### Phase 3.3 — Preflight 🟡 READY（definition only，**impl NOT YET**）
> **核心澄清（review）**：Consumer Equivalence **≠ Consumer Execution Identity**。equivalence 比的是 **observation boundary（findings）**，**不是 execution trace**。否則會誤收 hook(staged set+commit metadata) / CLI(worktree snapshot) / CI(checkout snapshot) 因輸入不同而出 finding 差異，最後有人把 context 搬回 engine。

**Consumer Equivalence Matrix（definition only）**：
| Dimension | Required? |
|---|---|
| RuleID set | **MUST equal** |
| Blocking | **MUST equal** |
| Opt-out effect | **MUST equal** |
| Discovery scope | **MUST equal**（都只看 `plans/active\|archived`） |
| Exit semantics | MAY differ |
| ExecutionMode | MAY differ |
| Input snapshot | MAY differ（hook=staged / CLI=worktree / CI=checkout） |
| Message text | IGNORE |
| Timing | IGNORE |

最重要兩條：**Discovery scope MUST equal**（同一發現面）+ **Input snapshot MAY differ**（輸入快照可不同，因為比的是 observation 不是 trace）。

**Reject 條件**：`若 equivalence 需要把 consumer-specific context 搬進 engine → FAIL 3.3`（守 engine=capability / consumer=transport，3.3 最易在此倒退）。

**Preflight ✅ APPROVED（2026-06-25）→ impl 可開始（3.3a 起）。** 再補兩條 acceptance（impl 遵守）：

- **E.1 — Canonical Observation Record (COR)**：三 consumer 的比較資料先收斂成同一 record 再比，避免 transport 資訊混進 equivalence：
  ```
  ObservationRecord { Findings[](RuleID+Blocking), OptOutEffect, DiscoveryScope }
  // 明確排除：ExitCode / ExecutionMode / SnapshotOrigin / Timing / Message
  ```
- **E.2 — Replacement must be directional**（收緊 replaceability）：`replace consumer → engine unchanged → remaining consumers unchanged → observation preserved`。**不接受** `replace consumer → patch engine → patch other consumer → pass`（那是整體重編，非 adapter replaceable）。

**3.3 子順序（不一次三向）**：
- **3.3a** manual ↔ hook（已有 2.3/2.4 大量前期證據）。
- **3.3b** hook ↔ CI（CI 是新 transport）。
- **3.3c** replacement proof（最後驗，避免把 CI 接入問題誤歸因 engine）。
- **Q1 維持 OPEN 到 3.3c 結束才關。**

##### 3.3a — manual ↔ hook ✅（2026-06-25）
- **COR 落地**：`ObservationRecord{ Findings(RuleID+Blocking), OptOutEffect, DiscoveryScope }`（`observation.go`），**結構上無 transport 欄**（reflection test `TestObservationRecord_ExcludesTransport` 鎖 exit/mode/snapshot/timing/message）。
- **manual(engine)↔hook(legacy) 等價**：三情境 COR `reflect.DeepEqual` 相等——valid（皆空）、violation（皆 `{parent_reference:true}`）、opt-out（皆 Findings 空 + `OptOutEffect{parent_reference}`）。engine policy-free：consumer 自行從 transport 解析 opt-out 再套用。
- DiscoveryScope 兩者同為 `plans/active|archived`（by construction equal）。
- **未做**：3.3b（hook↔CI）、3.3c（directional replacement proof）。Q1 仍 OPEN。

##### 3.3b — hook ↔ CI ✅（2026-06-25）
- **F.1 CI=adapter 非新 authority**：CI consumer = engine-backed（build ValidationContext → 同 engine entrypoint → consumer policy → render），不持 validation/schema、不 mutate findings。判定用 **CI COR == hook COR**（非 exit）。
- **F.2 snapshot normalization 在 adapter 不進 engine**：CI 讀 checkout 全樹由 loader normalize；guard test `TestValidationContext_NoSnapshotOriginField`（engine input 無 snapshot/origin/staged 欄）。
- **F.3 asymmetric proof**：violation tree 上 **raw CI（naive staged-empty）≠ hook**（漏看 violation）、**normalized CI == hook**（`{parent_reference}`）→ 證 equivalence 是 contract 非雙空巧合。三情境（valid/violation/opt-out）COR DeepEqual。
- **未做**：3.3c。Q1 仍 OPEN。

##### 3.3c — directional replacement ✅（2026-06-25；Q1 close package assembled）
- **真 replacement（非 mock）**：manual consumer 的兩個**真實 adapter** 互換——CLI `plans validate`（`runPlansValidate` via `Run`）↔ direct engine（`engineObservation`），皆過同一 engine entrypoint。`TestDirectionalReplacement_ManualAdapterSwap`：observation（COR + applicability）preserved。
- **R.1 removal-independence**：移除/不呼叫 manual → hook & CI COR 不變、無 fallback/shared-state（`TestDirectionalReplacement_RemovalIndependence`）。
- **R.2 contract stability**：consumer-facing contract fingerprint（`ValidationContext` + `Finding` 欄位集）== golden → engine 內部演化（同 contract）不需改 adapter（`TestDirectionalReplacement_ContractStable`）。
- **R.3 applicability（防雙空作弊）**：observation preservation ≠ set-equality；新增 `ruleApplicability`（3.3c-only，不入 COR）。`ApplicabilityNotSetEquality`：parent_reference「applicable+pass」vs「silently inapplicable」兩者 findings 皆空但 applicability 不同 → 非 preserved。

**Q1 close package（assembled，待 review 一起看；本輪不自行 close）**：
| 證據 | slice | 狀態 |
|---|---|---|
| equivalence proven（manual≡hook） | 3.3a | ✅ COR DeepEqual |
| equivalence proven（hook≡CI）+ asymmetric | 3.3b | ✅ + F.3 |
| removal-independence | 3.3c R.1 | ✅ |
| replaceability（≥1 真 replace） | 3.3c swap | ✅ CLI↔direct |
| contract stability | 3.3c R.2 | ✅ fingerprint |
| observation preserved（含 applicability） | 3.3c R.3 | ✅ |

→ **Q1 ✅ CLOSED（2026-06-25，review 通過）**。四層獨立鏈成立。Close Note A：fingerprint=guard 非 compatibility authority（meaning 保真靠 COR/applicability，勿當 version system）。Close Note B：未證 plurality（多外部 consumer 並存）→ 留 3.4，不升新 Q。

#### Phase 3.4 — Real Repo Acceptance（orchestration；authorized 2026-06-25）
> **重心轉移**：3.3 後已非 architecture work，3.4 是 **acceptance orchestration**。**不回頭重抽 engine、不再加 consumer、不提前處理 Q8。**
> **persistence 限制（review）**：acceptance repo **MUST persist（禁 throwaway tmp）**——3.4 驗 `install → validate → upgrade → rollback → remove`，需**時間維度**（跨 commit / 跨 session），非單次狀態。

**正式 split + 結論（review 2026-06-25，收在 Batch A′）**：
- **3.4a Real Repo Lifecycle Evidence → ✅ DONE**（adopt / validate / archive gate / rollback / monotonic removal，real repo、net-zero）。**未覆蓋**：invocation adapter、time-window install、upgrade-once、persistence across consumer transport。
- **3.4b Operational Acceptance → ⏸ DEFERRED**。**條件**：有明確窗口 + repo churn 降低 + 可控制 install/remove 時機。**覆蓋**：adapter install、real commit interception、upgrade-once、preserved semantics across operation。
- **Batch B（dialect canonicalize）→ ⛔ BLOCKED BY Q8**。

> **判定**：**Phase 3 不宣告 complete**；但 **architecture externalization 已完成**，剩下的是 **operational acceptance debt**（3.4b）+ schema policy（Q8）。
> **原則**：不為了完成階段而改變正在觀察的系統（同 soak 邏輯）——高 churn 期不裝 adapter；等低 churn 窗口再做 3.4b。

**3.4 Entry Gate（acceptance guard，非新 phase；acceptance repo 必須先滿足）**：
1. **非生產交付 repo**：不承擔正式 release；可接受 hook/workflow 暫時存在。
2. **獨立 rollback checkpoint**：install 前記錄 `git status clean` + `hooks absent` + `validation baseline`；rollback 後三者須回到同一狀態。
3. **acceptance evidence 不進產品語意**：不新增 plan metadata、不新增 schema 欄位、不為測試改 validation 行為。
4. **time-box**：install→validate→upgrade→rollback 在固定窗口完成（避免長期漂移）。

**Repo 選擇（review 偏好）**：
- ✅ **專用 acceptance repo**（最符合 Success Contract）。
- ⚠️ `apk-analysis-sdk`（僅當非高變動交付 repo）。
- ❌ `Vidoe-Test`（已提供 dialect pressure，角色乾淨，不混用）。

**3.4a — Acceptance Environment Selection（plan / selection only，🟢 authorized；不裝 adapter）**
- [x] **候選評估（2026-06-25，read-only）**：
  | repo | 14d churn | releases | 既有 evidence 角色 | scope-A `plans validate` | 判定 |
  |---|---|---|---|---|---|
  | Vidoe-Test | **169**（極高）| 0 | dialect-pressure + adoption-pass | `plans=0` | ❌ churn 違反 time-box + 角色已滿 |
  | apk-analysis-sdk | 27（中）| 0 | adoption-pass #2 | `plans=0` | ⚠️ churn OK，但仍 plans=0 |
- [x] **BLOCKER finding（selection 攔截，未裝 adapter）**：兩 repo plans 皆在 **`docs/plans/`**，scope-A CLI 固定 discover **`plans/active|archived`** → 實跑 `plans=0`。**先前 adoption-pass 是 path-override（EXT_PLANS_DIR→docs/plans）量測**，真實 scope-A CLI/adapter 掃不到。→ **兩 repo 皆非 acceptance-ready as-is**。
- **這是 location-convention gap（Q8 的 sibling，非 dialect）**：外部 repo 採了 canonical **schema** 但用了自己的 **location**（docs/plans）。修法選項：
  - ✅ **dedicated acceptance repo，用 canonical `plans/active|archived` layout**（sidestep mismatch、角色乾淨、persist、可 time-box）——**recommended**。
  - ⚠️ 在外部 repo 把 plans 搬到 `plans/active|archived`（invasive，改其 layout，違 Entry Gate #3 精神）。
  - ⛔ 給 CLI 加 `--plans-dir docs/plans` discovery option（location-convention policy，**Q8-adjacent → deferred**，不現在做）。
- [x] **selection 結論（2026-06-25，使用者覆寫 ❌）**：**選 Vidoe-Test**（使用者授權修改 + 提供「plan 完成狀態不確定」當測試材料）。
  - **checkpoint baseline 已捕捉（read-only）**：`git status CLEAN` / **commit-msg hook absent** / **validation baseline `plans=0`**（scope-A 尚無 `plans/active|archived`）。rollback 後須回到此三者。
  - **caveat 1 — 高 churn（169/14d）= 使用者正在密集使用**：commit-msg adapter **不可裝著不動**（會攔截真實 commit）；必須 **tight-window**（install→一次測試 commit→remove）。validate 主力走 **CLI/CI adapter**（`plans validate --root`，不攔 commit）。
  - **caveat 2 — 角色 overload**：Vidoe-Test 已是 dialect-pressure + adoption-pass；3.4 acceptance evidence 須**明確標為獨立角色**（acceptance = install/upgrade/rollback 機制 + 可逆，**不**重算 adoption/dialect）。
  - **location-gap 修法**：在 Vidoe-Test **新增可逆 acceptance artifact** `plans/active/`（放 canonical 測試 plans），不動其 `docs/plans/`；rollback 連 `plans/active/` 一併移除 → 回 baseline。
- [ ] **不安裝 adapter**（3.4a 仍只 selection；adapter install 屬 3.4b）。

##### Batch A′ — canonical h5 tree reversible acceptance ✅（2026-06-25，Vidoe-Test，net-zero）
> 範圍鎖死：僅 h5-redis canonical tree；**不碰 frontmatter / 不 canonicalize dialect / 不加 schema_version / 不改 parent / 不修 unrelated refs**。全程未 commit Vidoe-Test、結尾 restore。dialect plans 原地不動。
- **完成 required subs 確認**：main + 5 subs 皆 `status: completed`、`required_for_completion: true`。inbound ref 僅 README（docs 索引，restore 後即還原，未動）。
- **可逆 round-trip evidence**（filesystem move，非 commit）：
  - baseline（未移）`plans=0` → move→`plans/active/` `plans=6 findings=0`（discover + clean）→ move→`plans/archived/` `plans=6 findings=0`（**archive_order 生效**：全 required subs completed → gate 通過、無誤擋）→ rollback restore `plans=0`。
  - **residue check**：Vidoe-Test `git status` clean、tree 回 `docs/plans/` 原位、`plans/` 目錄移除 → **零殘留**。
- **證明（Success Contract 子集）**：reversible adoption ✅ + archive_order 當 completion gate ✅（真 repo）+ monotonic removal/no-residue ✅。
- **本批未含**：adapter(commit-msg hook) install + **upgrade once**（屬完整 3.4b 四段；archive_order「抓未完成」由 unit test 覆蓋，本批因全 completed 故 gate pass 非 block）。
- **Batch B（dialect plans canonicalize）= ⛔ BLOCKED BY Q8**（adoption / normalization / explicit-unsupported 未決），不碰。

**3.4b — Operational Acceptance ⏸ DEFERRED（target = Vidoe-Test；entry conditions：明確窗口 + repo churn 降低 + 可控制 install/remove 時機）**
> 覆蓋：adapter install / real commit interception / upgrade-once / preserved semantics across operation。**不為了完成階段而改變正在觀察的系統**——高 churn 期不裝 adapter，等低 churn 窗口再做。
> **非侵入式**（因高 churn）：validate 主力走 **CLI/CI adapter**（不攔 commit）；**commit-msg hook 只 tight-window 驗一次**（install→1 test commit→remove）。
- [ ] **setup（可逆）**：在 Vidoe-Test 新增 `plans/active/`（canonical 測試 plans，含「完成狀態不確定」的樣本）。
- [ ] **validate**：`plans validate --root <Vidoe-Test>` → 預期能 discover（plans>0）並就「**plan 完成狀態**」給 findings（archive_order：archived main 有未完成 required sub 則 block；frontmatter/parent/unique 同步）——這正是使用者要的「哪些 plan 沒完成」測試。
- [ ] **hook tight-window**：暫裝 commit-msg adapter → 一次測試 commit 驗 block/pass → 立即 remove（不留著攔真實 commit）。
- [ ] **upgrade once**：單軸單 subject（per 3.2，例：schema_version）→ findings 三層保真。
- [ ] **rollback**：remove adapter + `plans/active/` + config + binary ref → git/runtime/hook clean + no schema residue + monotonic → **回 3.4a checkpoint（status clean / no hook / plans=0）**。
- [ ] acceptance metadata：`repo_owner` / `repo_type: internal|public|fixture` / `removal_policy`（不寫 repo 名）。
- [ ] plurality（Close Note B，observation-only，不升 Q）。
- **完成 = Phase 3 complete**（Success Contract 四段 + 三層保真 + monotonic 全達成）。

> 第一次拿到「成功 adoption 證據」後最易把 compatibility 當單純升版測試。先拆軸，否則 3.2 測出綠燈卻不知哪層相容。**doc 內現存三個 version 必須分開**：binary version / `plan_schema` version / invocation-contract version（目前隱含）。

**三軸（升級前先回答）**：
- **Axis A — 什麼被升級？** binary / schema / both。
- **Axis B — 誰持有版本宣告？** adapter / CLI / `plan_schema`。
- **Axis C — 何時判不相容？** load / validate / invoke。

**Compatibility Subject（補口 1，明文化）**：每次升級必標 subject ∈ `{engine, consumer, artifact}`：
| Subject | 例子 |
|---|---|
| engine | binary 升版 |
| consumer | adapter / CLI 換版本 |
| artifact | `schema_version` 變動 |

**3.2 Success Contract（收緊；比 single-axis 再嚴一層）**：`upgrade once` **AND** `**1 compatibility axis × 1 compatibility subject**` **AND** `same validation semantics`（三層 Structural/Behavioral/Policy）。
- ✅ `binary(v1→v2) × engine`；❌ `binary+schema`；❌ `adapter+binary 同升`（invocation-contract 是第三個版本面，混升即失歸因）。

**3.2 Acceptance（≥3 例：supported / unsupported / baseline）**：
- **supported upgrade → findings preserved**（語意三層不變）。
- **unsupported combination → deterministic reject** **+ diagnosable（補口 2）**：同一 unsupported case 須 **same rejection stage + same rejection reason class**（分類即可，如 `NormalizeReject` / `InvokeReject` / `TransportReject`，不比字串），防診斷面漂移（Run A load-reject / Run B invoke-reject）。
- **no-version-change baseline → findings unchanged（補口 3，negative proof）**：跑一次 upgrade machinery 但**不改任何 version**，證明 findings 不變——確認驗的是相容性、不是「跑到新路徑」。

- [ ] **（3.2 impl 未授權）**：三補丁（compatibility subject / reject stage-stability / no-change baseline）已入 plan；待授權才升 binary 落地。

## 完成條件
- [ ] portable 分類表（含 `consumer_surface` 欄）+ `plan_profile`（capability）/ `plan_schema`（compat）邊界落地（Q2 resolved，由分類表推導，capability 與 execution 維度分離）
- [ ] validator engine（吃 `ValidationContext`）+ schema compat layer + thin consumers（hook / CLI）+ 測試通過（含 no-CLI engine integration test）
- [ ] 既有 commit-msg hook 行為不變（重構回歸驗證）
- [ ] Phase 3 evidence slices：3.1 Adoption（reversible adoption + rollback clean）/ 3.2 Compatibility（→ Q3 close）/ 3.3 Consumer Equivalence（manual≡hook≡CI → Q1 close）/ 3.4 Real Repo Acceptance（四段 install→validate→upgrade→rollback，唯一完成標準）
- [ ] Phase 3 Non-goal 守住：無 external registry / external runtime state / external plan metadata
- [ ] Q1 / Q3 依 slice 3.3 / 3.2 close 並回寫；Q8 deferred

## Glossary Impact
Glossary Impact: yes — 新增 `plan_profile`（capability / portable 邊界）與 `plan_schema`（frontmatter schema + version 相容契約），刻意拆成兩個單一責任術語；Phase 1 落地時註冊到 `knowledge/glossary/ai-skill.md`。

## 與其他 plans 的關係
- 依賴 [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) 的 5 validators。
- 依賴 [`archived/2026-05-28-1830-plan-archival-audit-validator.md`](../../archived/2026-05-28-1830-plan-archival-audit-validator.md) 的 archival audit。
