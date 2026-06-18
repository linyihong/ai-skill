# Software Delivery Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/requirements/README.md`](../../intelligence/engineering/requirements/README.md)
- [`intelligence/engineering/development/docs-first-bdd-closure.md`](../../intelligence/engineering/development/docs-first-bdd-closure.md)
- [`analysis/development-guidance/risk-translation.md`](../../analysis/development-guidance/risk-translation.md)
- [`intelligence/engineering/requirements/product-alignment/README.md`](../../intelligence/engineering/requirements/product-alignment/README.md)
- [`workflow/software-delivery/requirements/README.md`](../../workflow/software-delivery/requirements/README.md)
- [`workflow/software-delivery/development-process.md`](../../workflow/software-delivery/development-process.md)
- [`workflow/software-delivery/ui-contracts.md`](../../workflow/software-delivery/ui-contracts.md)
- [`workflow/software-delivery/ui-governance.md`](../../workflow/software-delivery/ui-governance.md)
- [`workflow/software-delivery/incident-observation.md`](../../workflow/software-delivery/incident-observation.md)
- [`workflow/software-delivery/ui-incident-governance-workflow.md`](../../workflow/software-delivery/ui-incident-governance-workflow.md)
- [`workflow/software-delivery/change-retrospective.md`](../../workflow/software-delivery/change-retrospective.md)

本文件把 pre-build interrogation、product alignment、requirements cognition、docs-first BDD closure、contract-first development、UI incident layer selection、change retrospective 與 development guidance 的風險翻譯方法轉譯成 AI runtime software-delivery governance。原始 intelligence 回答「如何確認產品方向、穩定 observable behavior、acceptance、traceability 與 validation target」；本文件定義 change intake、pre-build interrogation、product alignment、requirements cognition、contract precedence、BDD closure、artifact completeness、performance evidence 與 same-session documentation closure 的治理 gate。

## 觸發時機

在下列情況套用本治理：

- 使用者要求 implement、develop、SDK、API、embedded、plan、code review 或 design review。
- 變更可能影響 observable behavior、public contract、domain invariant、API/schema、UI / consumer surface、error handling、storage、安全性、ownership、tests 或 performance。
- Product brief、BDD、contract、implementation 或 tests 之間出現 mismatch。
- 回填已實作專案的 BDD、contract、test plan、hardware/embedded evidence 或 traceability。
- UI / consumer incident（Navigation / Continuation / Recovery 未決）且 agent 可能直接跳至 implementation。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Change intake | 已分類為新需求、bug、refactor / replacement、安全/強化、performance 或 planning-only，並確認 code 前需要的 artifacts。 |
| Pre-build interrogation | Plan 或 implementation 前已釐清 goal、scope、non-goals、acceptance、validation target、framework source-of-truth、duplication risk、open questions 與 assumptions。 |
| Brief validation | Product brief 的主要 claim 已標記 `validated`、`assumption`、`open question`、`scoped out` 或 `invalidated`。 |
| Product alignment | Impact Map 的 Why / Who / How / What 已與 Customer Journey 的 actor、timing、pain point、blocker 交叉驗證。 |
| Requirements cognition | Observable behavior 已有 actor intent、behavior boundary、acceptance criteria、validation target 與 ambiguity disposition。 |
| Contract precedence | 衝突時先判斷 governing contract、product plan、BDD、domain/API/error/hardware contract、implementation、tests 的優先序。 |
| Docs-first BDD closure | Observable behavior 變更前，owning contract、BDD scenario、executable validation 與 implementation slice 已同步或明確 scope out。 |
| Refactor parity | 重構、遷移、改寫或 replacement 會替代既有功能、入口、腳本、API、資料流程、runtime surface 或操作流程時，已建立新舊能力 parity inventory，並為每個舊入口標明新入口、parity 狀態、副作用、外部依賴與測試 / fixture 證據。 |
| Artifact completeness | Change brief、contract、BDD scenario、implementation plan、review report 或 project-local equivalent 已產出或標記 not applicable。 |
| UI governance advisory projection | UI compliance claim 已載入 `sd-ui-governance`，並分類 governance domain、validation mechanism、evidence class、severity、project-local design-system policy 與 visual / AI review scope；此 gate 在本階段是 workflow / review advisory，不是 runtime hard block。 |
| Test strategy | 新行為與舊行為 regression 分開驗證；contract、fixture、integration、property、targeted mutation check 或 hardware-in-loop evidence 依風險選擇。 |
| Performance evidence | 影響 latency、throughput、資源、concurrency、external-call fan-out 時，不以功能測試取代 performance budget / smoke / load / stress / spike / soak evidence。 |
| Same-session closure | Code、docs、contracts、BDD、tests、generated clients、fixtures 與 linked updates 在同一批次閉環，或留下明確 owner 與 scoped debt。 |
| Incident observation | UI / consumer incident 已完成 incident card（symptom、timeline、observable per step）；未讀 hook / storage / PR 來跳過 observable。Workflow: [`incident-observation.md`](../../workflow/software-delivery/incident-observation.md)。 |
| Incident classification | 恰好一個 domain：Navigation \| Continuation \| Recovery \| Out-of-scope；禁止雙 domain 或 implementation-first classify。Workflow: [`ui-incident-governance-workflow.md`](../../workflow/software-delivery/ui-incident-governance-workflow.md) §Stage 1。 |
| Incident layer selection | 恰好一個 primary modification layer：Contract \| Overlay \| Verification \| Integration；對照 [`layer-ownership-matrix.md`](../../workflow/software-delivery/layer-ownership-matrix.md)。**Single-layer convergence**: YES → 允許進入 Execute（Contract / Implementation）；NO → **允許僅擴 verification**（integration、evidence sheet）；⚠️ overlay 需 review；**禁止**升 contract、新 abstraction / hub；禁止以「不能改」開新 plan 代替驗證。 |
| Change retrospective | Ship 後已填 retrospective：哪層被改、哪層未改、vocabulary/consumer、promotion 建議 ∈ {keep local, promote project, candidate canonical}；**禁止** direct canonical promote。Workflow: [`change-retrospective.md`](../../workflow/software-delivery/change-retrospective.md)。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 為什麼 feature investment 對準 business impact 與 user journey | `intelligence/engineering/requirements/product-alignment/` |
| 為什麼 BDD 是 requirements cognition、如何處理 ambiguity / acceptance / traceability | `intelligence/engineering/requirements/` |
| 為什麼 contract 要先於 code、為什麼 BDD 是 behavior bridge | `intelligence/engineering/development/` |
| Software delivery 的 AI runtime gate 與 completion criteria | `governance/ai-runtime-governance/` |
| 實際執行順序、輸出模板選擇、工作流程步驟 | `workflow/software-delivery/` |
| 風險翻譯、owner layer selection、guidance classification | `analysis/development-guidance/` |
| Controls / implementation / platform / language catalog | `metadata/development-guidance/` |

## Workflow Mapping

- [`workflow/software-delivery/requirements/README.md`](../../workflow/software-delivery/requirements/README.md) — product alignment and requirements cognition stage。
- [`workflow/software-delivery/requirements/pre-build-interrogation.md`](../../workflow/software-delivery/requirements/pre-build-interrogation.md) — plan / implementation 前的需求拷問、framework discovery 與 source-of-truth duplication gate。
- [`workflow/software-delivery/execution-flow.md`](../../workflow/software-delivery/execution-flow.md) — software-delivery workflow entry and execution order。
- [`workflow/software-delivery/development-process.md`](../../workflow/software-delivery/development-process.md) — contract-first development process and detailed gates。
- [`workflow/software-delivery/ui-contracts.md`](../../workflow/software-delivery/ui-contracts.md) — expected UI / consumer behavior contracts。
- [`workflow/software-delivery/ui-governance.md`](../../workflow/software-delivery/ui-governance.md) — UI compliance classification surface for domains, mechanisms, evidence, severity, and advisory runtime projection boundary。
- [`workflow/software-delivery/incident-observation.md`](../../workflow/software-delivery/incident-observation.md) — Stage 0 Observe: incident card before classify。
- [`workflow/software-delivery/ui-incident-governance-workflow.md`](../../workflow/software-delivery/ui-incident-governance-workflow.md) — Stage 1 Classify + Stage 2 Select Layer workflow。
- [`workflow/software-delivery/layer-ownership-matrix.md`](../../workflow/software-delivery/layer-ownership-matrix.md) — authority → domain owner → allowed modifications。
- [`workflow/software-delivery/change-retrospective.md`](../../workflow/software-delivery/change-retrospective.md) — Ship → Retrospective；promotion 三選一。
- [`workflow/software-delivery/artifact-gates.md`](../../workflow/software-delivery/artifact-gates.md) — reusable note structure and artifact quality gates。
- [`analysis/development-guidance/README.md`](../../analysis/development-guidance/README.md) — development guidance analysis methods。

## Runtime-Lite Boundary

本治理不讓 runtime 理解 BDD syntax、Gherkin、scenario grammar 或 universal requirement schema。只有下列壓縮訊號可作為未來 runtime-lite 候選：

- `product_goal_mismatch`：business goal、target actor、journey pain 或 feature investment 不一致。
- `feature_without_impact`：feature list 缺少 Why / Who / How impact chain。
- `requirement_contradiction`：requirements、BDD、contract、implementation 或 tests 出現互斥 claim。
- `missing_pre_build_interrogation`：模糊需求、plan、framework 改動或 source-of-truth 風險未先經過需求拷問。
- `framework_source_duplication`：同一 executable semantics 同時存在兩份 rule body、activation path、runtime table、mirror 或 generated surface。
- `missing_validation_target`：acceptance criteria 沒有可執行或可審查的 proof target。
- `stale_acceptance_criteria`：product intent、domain invariant 或 implementation truth 改變後，acceptance baseline 未同步。
- `behavior_scope_overclaim`：local scenario pass 被宣稱為 global feature correctness。
- `implementation_first_incident_classify`：未產出 incident card 或未寫 primary layer 就開 hook / storage / code。
- `incident_layer_not_converged`：primary layer 無法單層收斂卻仍宣稱 ready to implement。
- `authority_layer_mismatch`：scroll / viewport 問題直接改 contract，或 route 問題直接改 continuation overlay，違反 layer-ownership-matrix。

**Incident signals — conservative boundary (Phase B)**:

- These are **advisory metadata only** — `signal ≠ lifecycle`.
- Validated this pilot: **authority trace** (observable → first broken authority → allowed layer).
- **Not** validated: runtime state model, Experience Runtime, or persistent runtime tables.
- Do **not** promote signals to `runtime/` generated surfaces or Experience Runtime Governance without a separate plan and second independent incident.

任何 promotion 都必須另開 plan，確認 compiler / generated surface；預設維持 metadata-only。

### UI Governance Advisory Candidates

UI governance signals are runtime-lite advisory projections. They name review pressure but do not create a runtime hard block, generated surface, detector, or enforcement rule_class in this landing.

Path:

```text
event: UI / consumer surface compliance claim
→ detector / route / query: route.workflow.software-delivery loads execution-flow.md
→ loaded source: workflow/software-delivery/ui-governance.md
→ advisory projection: candidate signal named below
→ evidence: UI governance domain + mechanism + evidence class + severity + project-local policy
```

Candidate signals:

- `missing_ui_contract_state`：UI contract / screen evidence lacks required loading, empty, error, success, permission, or retry state.
- `raw_design_token_detected`：project-local policy requires tokens / approved primitives but artifact or review evidence shows raw style escape hatch.
- `destructive_action_missing_confirmation`：delete, revoke, reset, irreversible submit, or permission-sensitive action lacks objective confirmation / recovery behavior.
- `visual_claim_without_baseline`：visual quality, regression, or parity is claimed without deterministic baseline, capture context, or scoped review evidence.
- `ai_visual_validator_without_scope`：AI visual review is used without rubric, prompt scope, warning/research severity, or human/project opt-in boundary.

Promotion boundary:

- Advisory signals may trigger `sd-ui-governance` loading, review checklist attention, or evidence-template completion.
- Blocking enforcement requires a separate promotion plan naming executor, detector input, evidence threshold, generated surface, registry transition, and rollback path.
- `enforcement/enforcement-registry.yaml` is intentionally unchanged in the first landing.

## Validation Candidate

後續若要 promotion 到 `validation/`，可建立 scenario 檢查：

- Requirement / BDD / tests 互相矛盾卻繼續 implementation。
- 模糊需求或 framework 改動沒有 pre-build interrogation 就產出 implementation plan。
- Framework 改動沒有識別 canonical source / projection / duplicate surface 就開始實作。
- Product goal、target actor、journey pain 或 feature investment 不一致卻直接產生 implementation plan。
- Acceptance criteria 缺 validation target 卻宣稱 ready。
- Observable behavior change 只改 code，未更新 owning contract / BDD / tests。
- Product brief claim 未驗證卻被當成 implementation input。
- API/schema contract 改變但 generated client、fixtures 或 consumer tests 未同步。
- Refactor / replacement 只描述新設計，沒有盤點舊入口、舊能力、副作用、外部依賴、new surface mapping 與 parity 測試證據。
- Performance-sensitive change 只用 unit/functional tests 宣稱可 release。
- Existing project backfill 憑空創造 product intent，而不是標記 `unknown` 或 `open question`。
- UI compliance claim 缺少 domain / mechanism / evidence class / severity classification。
- Visual quality claim 沒有 screenshot baseline 或 AI visual review scope。
