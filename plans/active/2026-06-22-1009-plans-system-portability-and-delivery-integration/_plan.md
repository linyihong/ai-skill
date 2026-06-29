---
id: 2026-06-22-1009-plans-system-portability-and-delivery-integration
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-22
parent: null
---

# Plans System Portability & Delivery Integration（plans 系統外部化與交付接入）

**Status**: `draft`
**Owner**: linyihong
**建立日期**: 2026-06-22
**Source**: 2026-06-22 對話 — 使用者觀察 plans 系統（plan-tree + 驗證機制 + 子計畫）已足夠成熟，希望 (1) 讓外部 repo 也能使用、(2) 把 software-delivery 接入「開發前先寫 plan」、(3) 用子計畫系統讓一個 sub-plan 可交給其他 agent 執行。
**Glossary Impact**: yes — 預期新引入 framework vocabulary，**刻意拆成單一責任術語避免概念漂移**：`plan_profile`（capability / portable core 邊界：哪些 validator 對外部 repo 適用）、`plan_schema`（frontmatter schema + version 相容契約）、`delegation`（子計畫委派 schema：modes + brief + constraints）。落地前須註冊到 `knowledge/glossary/ai-skill.md`；本 main plan 階段先宣告，sub-plan graduate 時才註冊。

> **Watch-Out List citation**：本 plan 對應 [`architecture/ai-native-cognitive-ecosystem-system.md`](../../../architecture/ai-native-cognitive-ecosystem-system.md) §Watch-Out List 的「process bloat」「premature abstraction」「over-engineering」防呆 — plan-first 採 advisory workflow ordering 而非機械 block（避免誤擋小修補）；委派只先做 schema + 雙路徑（人工 / agent），不建自動 orchestrator；外部化採共用 binary 而非整套治理搬遷，避免外部 repo 被 Ai-skill governance 綁架。

---

## Decision Rationale

### Problem & Why Now

plans 系統目前的價值（plan-tree 階層、frontmatter 單一 source-of-truth、5 個 commit-msg validators、archival audit / link integrity、pre-build interrogation → preflight → completion closure）都封在 Ai-skill repo 內部，且與 Ai-skill 專屬治理（`runtime.db`、`generated_surfaces`、routing-registry、Glossary、Cognitive Modes、ADR pipeline、intelligence atoms）混在同一份 plan 模板。三件事同時成熟到值得動手：

1. **外部 repo 想用**：`ai-skill plans tree --root <path>` 已支援任意 repo root（讀取面已跨 repo），但 commit-time 強制（validators）仍是本 repo 的 git hook，外部 repo 無法享受。
2. **software-delivery 缺 plan-first ordering**：現有 intake 有 pre-build-interrogation / Test-First Ordering，但沒有明文要求「實作前先有 plan」的順序，導致 plan 與 delivery 各走各的。
3. **子計畫已有獨立 acceptance / archive，但沒有委派出口**：sub-plan 可獨立 sign-off，卻沒有「自足 brief + 可交給其他 agent（人工或 Agent 工具）執行」的 schema。

### Decision

開一個 **plan tree（main + 3 sub-plan）**，dogfood plan-tree 本身，分階段 graduate：

- **Sub-plan 01 — External-repo plan system via shared binary**：外部 repo 透過**共用 ai-skill binary**（非 init-project 抽取安裝）使用 plans 系統。核心缺口是 commit-time 強制跨 repo。**抽象層關鍵：要抽的是 validator engine package（被 git hook shim / CI / CLI / 未來 API 共用），CLI `plans validate` 只是其中一個 consumer，不是核心**，否則半年後 `plans validate` 會長成另一個 orchestration layer。同時以**分類模型推導** portable 邊界（`plan_profile` capability + `plan_schema` 相容契約），不預設「哪些 validator 屬 portable」。
- **Sub-plan 02 — software-delivery plan-first ordering**：把 plan-first 寫進 `workflow/software-delivery/` intake 段，接在 pre-build-interrogation / Test-First Ordering 之後，**advisory + review 檢查，不做機械 block**。
- **Sub-plan 03 — Sub-plan agent delegation**：在 sub-plan frontmatter 加 **nested `delegation` 物件**（`enabled` + `modes:[human, agent]` + `brief` + `constraints`），避免把「可委派」與「brief 存在」綁死，並支援 manual-only / agent-only / hybrid / forbidden 不破 schema。**同時支援人工派發與 Agent/Task 工具派發**（依專案需求選用或並用），執行端兩條路都通。

### Alternatives Considered

- **A. 維持現狀（plans 只在 Ai-skill 內用）**：reject — 使用者明確要外部化，且讀取面已跨 repo，差最後一哩 commit-time 強制。
- **B. init-project 把整套 plan 模板抽取安裝進外部 repo**：reject（本輪）— 使用者選「共用 binary 指向外部 repo」；整套抽取會把 Ai-skill governance overlay 一起帶過去，違反 portability 初衷且工作量最大。保留為 01 的 future option。
- **C. plan-first 做成硬機械 gate（無 active plan 不准 commit code）**：reject（本輪）— 使用者選 workflow 層 ordering；硬 gate 易誤擋小修補、需複雜 opt-out。保留為 02 的後續升級候選（maturity ladder：ordering 觀察 → 視 evidence 再升級）。
- **D. 一次做完三件**：reject — 三件可分開 graduate，分階段降低 blast radius，且 plan-tree 本身就是用來表達這種 main↔sub 拆分。

### Why Not an ADR Yet

scope 仍會隨 Phase 0 盤點調整（特別是 01 的跨 repo 強制機制、03 的 delegation schema 形狀）；多個 Open Questions 未解；可能有更輕的 promotion target（例如 02 可能只需 workflow doc 更新，不需任何 runtime contract）。三條 sub-plan 各自 graduate 後，若浮現需跨 session / 跨 project 固化的決策，再評估 ADR。

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] 三條 sub-plan 結果證實 portable core 邊界、plan-first ordering、delegation schema 三者可行且被真實使用
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用（per ADR-007）
- [ ] 至少一個外部 repo 真實透過共用 binary 跑過 plan validate（具體 evidence）

### Consequences（預期）

#### 正面
- 外部 repo 不需 fork Ai-skill governance 就能用成熟的 plan-tree + 驗證閉環。
- software-delivery 與 plans 系統對齊，「先規劃後實作」成為可見的 intake 順序。
- 子計畫成為可委派的執行單元，支援多 agent / 人工混合交付。

#### 負面
- 共用 binary 路徑讓外部 repo 對 Ai-skill binary 版本產生依賴（需版本相容策略）。
- delegation schema 增加 sub-plan frontmatter 表面積。

#### 風險
- portable 邊界若用「validator 類型」直覺切而非用 contract/dependency/execution-context 分類模型推導，會變成「先決定 portable 再找理由」→ 01 Phase 1 必須先產分類表（validator → contract_source → runtime_dependency → portable → reason）再分類。
- plan-first ordering 若沒接好既有 pre-build-interrogation，會變成重複 gate（process bloat）。

---

## Runtime Execution Path

本 main plan 為 **規劃容器**，自身不直接接入 runtime；具體 runtime / workflow / validation 改動由各 sub-plan 宣告：

- **01** 核心是 **validator engine**（consumed by hook / CI / CLI / future API — 不以任一 consumer 為標準入口；CLI `plans validate` 只是其中一個 surface）。若新增 commit-msg validator 行為或 `route.*`，須在 01 自己的 Runtime Execution Path + Per-surface consumer 表宣告，並由 `validateRuntimeTriggerWiring` 機械驗證。
- **02** 預期僅改 `workflow/software-delivery/` 文件（intake ordering）+ 可能新增 validation scenario；若不接 runtime，須在 02 明寫 doc-only + 未來接入條件。
- **03** 預期改 plan-tree frontmatter schema + 可能擴充 `validatePlanTreeFrontmatter`（新增 optional delegation 欄位驗證）；若擴充 validator 須宣告 trigger flow。

**doc-only trial 宣告**：本 main plan `status: draft` 階段不接入 runtime；接入由 sub-plan graduate 時各自落地。

---

## Open Questions

> 本表是 canonical Open Questions registry；sub-plan 的「已讀 §Open Questions」核對與其 Phase 0 公版 checklist 以此表為錨。
>
> **關閉規則（避免 `resolved = 作者覺得差不多`）**：
> - **誰可關閉**：僅 `Resolved By` 指定的 sub-plan owner 可把 Status 改 `resolved`。
> - **何時可關閉**：`Closed Criteria` 全部成立時，不早於該 sub-plan 對應 Phase 完成。
> - **關閉證據去哪**：`Resolution Evidence` 必須指向具體 commit / 檔案 / 測試（非「已討論」）；其他 sub-plan 對結論有異議時，在此表加 `disputed` 註記而非各自繞過。

| ID | Question | Owner / Resolved By | Status | Closed Criteria | Resolution Evidence |
|----|----------|---------------------|--------|-----------------|---------------------|
| Q1 | 外部 repo 跑 plan validators 的最薄強制機制（validator engine 被哪些 consumer 呼叫）？ | 01 | **✅ CLOSED（2026-06-25）** | 四層獨立證明鏈：**Equivalence**（3.3a/3.3b COR + asymmetric）+ **Replaceability**（≥1 真 replace CLI↔direct）+ **Removal-Independence**（R.1）+ **Semantic-Preservation**（R.3 COR+applicability，防雙空）→ 支持 `validation capability ≠ consumer transport` | 01 §Phase 3.3 close package。**Close Note A**：fingerprint 是 guard（保 contract-unchanged）非 compatibility authority（保 meaning 仍靠 COR/applicability/preserved-semantics）；勿當 version system。**Close Note B**：已證 replace-one，未證 plurality（多外部 consumer 並存）——不升新 Q，留 3.4 orchestration |
| Q2 | portable 邊界如何**推導**（非預設 plan-tree 5 + archival 2）？ | 01 | **resolved（2026-06-23）** | **收緊**：Layer A facts 完整 + Layer B decisions review + `plan_profile` committed **+ 至少一個 consumer 成功執行**（taxonomy ≠ 驗證 portable）✅ 全達成 | 01 §Phase 1 兩層分類 + §Phase 2.2 engine integration test（首個 consumer 綠）+ Gate D.4 negative evidence（excluded validators 結構上無法表達）。邊界由分類推導 + consumer 驗證 + 反證三重確立 |
| Q3 | schema / 版本相容策略（pin 哪個 binary、schema version 怎麼宣告與演進）？ | 01 | **resolved（2026-06-25，commit `2c26f6e`）** | `plan_schema` version 宣告 + 跨版本 acceptance pass ✅ | Phase 3.2：extensible `supportedSchemaVersions` set + `CompatError` deterministic+diagnosable reject + end-to-end loader wiring（`SchemaVersion`，parser strip quote 滿足 YAML 引號需求）+ no-change baseline。測試：supported 1→2 preserved / unsupported v99 → blocking reject / CLI exit 30。subject=artifact、單軸單 subject |
| Q4 | plan-first 與 pre-build-interrogation / Preflight 分工（plan 是 artifact、preflight 回改 plan，loop 非線性）？ | 02 | open | intake loop 段落落地 + 一次真實 intake 含 preflight 回改實例 | <commit + intake evidence> |
| Q5 | delegation 最小契約：nested `delegation: { enabled, modes, brief, constraints }`，不綁死、支援 manual/agent/hybrid/forbidden？ | 03 | open | schema + validator + 測試 committed | <commit + test path> |
| Q6 | tool-neutral 邊界：brief 只定義自足資訊，工具細節歸 `constraints` / `ai-tools/` 不進 schema？ | 03 | open | brief 契約文件化 tool-neutral + 雙路徑 SOP | <commit + ai-tools path> |
| Q7 | **validator failure semantics 是 contract 不是 consumer detail**：同一 engine 在 hook→block / CI→fail / manual→warning，severity 映射歸誰？（Phase 1 facts 已碰 severity + opt-out，依升格準則正式立案） | 01 | open | failure-semantics 映射在 engine contract 層定義（severity 為 engine 輸出，consumer 只決定 transport 行為） | Phase 1 Layer A 含 severity(block/warn) + opt-out transport → 觸發升格 |
| Q8 | **external schema compatibility boundary**（= **compatibility-policy** bucket）：For repositories using non-canonical plan metadata, should compatibility be enforced by **adoption**, **normalization**, or **explicit unsupported declaration**?（刻意保留第三條；不預設 mapping） | Phase 3（**deferred-to-phase-3**，不掛 01 current phases） | deferred-to-phase-3 | 待 Phase 3 跑真實外部 repo 後，依證據三選一 | **僅 dialect-pressure 證據掛此**（Vidoe-Test flat plans semantic mismatch：`parent` path vs id，measured 2026-06-24）。**adoption-pass 證據（canonical tree clean）NOT 掛 Q8**——它是 Phase 3 acceptance anchor，「一 branch 可行 ≠ 該選 adoption」。policy decision = **no**（deferred）；不觸發 plan_profile reopen |

---

## 完成條件

- [ ] 三條 sub-plan（01 / 02 / 03）皆 `status: completed`
- [ ] 各 sub-plan 的 acceptance 達成並通過其宣告的 validation
- [ ] Open Questions 全部標記 `resolved` / `deferred`（附原因）並回寫
- [ ] Glossary Impact 落實：新 vocabulary 已註冊或明確不註冊
- [ ] 執行 Plan Completion Closure（含 `ai-skill runtime refresh` 若涉 knowledge/validation 層）
- [ ] plan tree 通過 `ai-skill plans tree` 檢視（main + 3 sub 階層正確）

## Phase 0 — Pre-Build Interrogation

### Phase 0.0 — Open Questions 核對（公版，必填）

逐條核對本 plan §Open Questions，標記處置並回寫：

- [ ] 已讀本 plan §Open Questions 全部條目
- [ ] 對每條標記 `resolved`（附 Phase 0 證據）/ `still-open` / `deferred`（附原因）
- [ ] `resolved` 的條目已同步勾選 / 附註於 §Open Questions
- [ ] 若盤點新發現問題，已加入 §Open Questions

| Open Question | 處置 | 證據 / 原因 |
|---|---|---|
| Q1 / Q2 / Q3 | open | 由 01 Phase 0 盤點解決並回寫 §Open Questions Status |
| Q4 | open | 由 02 Phase 0 解決 |
| Q5 / Q6 | open | 由 03 Phase 0 解決 |

### Phase 0.1 — 架構相容性 preflight（main 層）

- [ ] 確認 `plans/README.md` Plan Tree Hierarchy 規則與本 tree frontmatter 相容
- [ ] 確認 `scripts/ai-skill-cli/internal/app/plan_tree.go` / `plans.go` 現行 schema（避免 03 改 frontmatter 撞既有 validator）
- [ ] 確認 `workflow/software-delivery/intake.md` 現行 intake 順序（02 接入點）
- [ ] 確認 `governance/lifecycle/plan-tree-hierarchy.md` 治理規則

## Phase 1-N

實作步驟拆進三條 sub-plan，各自有 Phase 0 / Phase 1-N / 完成條件：

> **Sequencing 是 recommended，不是 hard dependency（回應 review #7，避免不必要 serialization）**：三條無強制先後。01 先做只是因為它產出的 `plan_schema` frontmatter 共識對 03 有幫助；02 **完全獨立**（純 workflow doc，不依賴 portable 邊界）；03 只需與 01 對齊 frontmatter schema，**不依賴外部 repo 能力**。可並行開工。

| Sub-plan | 主題 | required_for_completion | Sequencing（recommended，非依賴） |
|---|---|---|---|
| [`01-external-repo-plan-system-shared-binary.md`](01-external-repo-plan-system-shared-binary.md) | 外部 repo 經共用 binary 使用 plans 系統 + portable 邊界 | true | recommended first（產 `plan_schema` 共識，利於 03） |
| [`02-software-delivery-plan-first-ordering.md`](02-software-delivery-plan-first-ordering.md) | software-delivery plan-first workflow ordering（advisory） | true | independent（可隨時開工） |
| [`03-subplan-agent-delegation.md`](03-subplan-agent-delegation.md) | sub-plan 委派 schema（人工 + agent 雙路徑） | true | 需與 01 對齊 frontmatter schema；不依賴外部 repo |

## Stakeholder 同意項目

> 描述**現行選定策略**（治理現況），非聊天紀錄。改方向時直接更新本表，不視為「推翻歷史決策」。

| 決策面 | Current selected strategy |
|--------|---------------------------|
| 外部化機制 | shared binary（非 init-project 抽取；抽取保留為 future option） |
| plan-first gate | workflow ordering advisory（非機械 block；保留 maturity-ladder 升級候選） |
| 委派路徑 | human + agent 雙路徑（nested `delegation.modes`，依專案選用或並用） |
| 交付節奏 | 分階段、dogfood plan tree |

## Glossary Impact

Glossary Impact: yes — 預期新增 `plan_profile`、`plan_schema`、`delegation`（各為單一責任術語，避免 `plan_system_profile` 一詞同時背 capability / validator set / 相容契約 / schema version 四個責任）；本 main plan 階段僅宣告，實際註冊由 sub-plan graduate 時落地（避免提前註冊未定稿術語）。

## 與其他 plans 的關係

- [`archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md`](../../archived/2026-06-02-1200-plan-tree-hierarchy-governance/_plan.md) — 本 tree 直接建立於其 frontmatter schema + 5 validators 之上；03 將擴充其 frontmatter。
- [`archived/2026-05-28-1830-plan-archival-audit-validator.md`](../../archived/2026-05-28-1830-plan-archival-audit-validator.md) — archival audit validator，01 portable core 候選成員。
- [`workflow/software-delivery/intake.md`](../../../workflow/software-delivery/intake.md) — 02 的接入點。
