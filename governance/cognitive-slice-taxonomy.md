# Cognitive Slice / Surface Taxonomy（Phase 1 complete）

**Status**: `phase-3-loading-linked`（taxonomy 已定義並套用 software-delivery pilot；Phase 2 已抽出 6 個 focused surfaces + 1 pre-existing examples surface，`sd-implementation` 依 stakeholder 決定暫留 execution-flow core；Phase 3 已將 focused surfaces 掛到 existing hierarchical route / executable contract / graph / summary，fixtures 為 test-first 草稿待 Phase 4 執行）
**Owner layer**: governance
**來源 plan**: [`plans/archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](../plans/archived/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md) §Phase 1
**命名決定**：見 §6。framework vocabulary 正式註冊**延後至 Phase 4 validation**；本檔過渡期一律用 `execution / evidence surface` 措辭。

> 本檔源自 Phase 1 taxonomy 產出：定義 slice schema + 5 條治理規則，並把它們**套用到 software-delivery pilot**（§7 slice 盤點）。Phase 2 已完成 focused surface extraction 的主要部分；Phase 3 已完成 loading/routing link 同步；scenario 執行在 Phase 4。

---

## 1. Slice schema 欄位（spec）

| 欄位 | 必填 | 說明 |
|---|---|---|
| `id` | 是 | slice 穩定識別碼（pilot 用 `sd-<phase>`） |
| `purpose` | 是 | 此 slice 要讓 agent 完成的認知目的 |
| `type` | 是 | primary，**只允許 4 值**：`execution` / `evidence` / `examples` / `failure` |
| `tags` | 否 | secondary 自由標註（artifact-gate / closure / handoff / templates / observation-triage / tool-procedure / domain-specific / extraction-to-intelligence …） |
| `load_when` | 是 | 何種 task intent 應載入 |
| `do_not_load_when` | 是 | 何種任務不應載入（suppression） |
| `owner_layer` | 是 | workflow / analysis / intelligence（依 §4 三層邊界規則判定） |
| `layer_justification` | 是 | 歸層的 falsifiable 理由，須通過該層 membership predicate（§4） |
| `evidence_refs` | intelligence 必填 | ≥2 個獨立、已驗證、可解析的 analysis 觀察 / failure case 指標 |
| `canonical_source` | 是 | 正文 canonical 來源（slice 只導航，不重定義） |
| `dependencies` | 否 | 依賴的其他 slice / source |
| `dependency_budget` | 是 | heuristic default `max_depth:2`/`max_runtime_dependencies:4` + `override_when: task_complexity=high`（非 rigid） |
| `summary_path` | 否 | 對應 summary-first 入口 |
| `validation_signal` | 是 | Phase 4 用哪個 scenario 驗證 |

---

## 2. type+tags 收斂規則

primary `type` 固定 4 種（`execution` / `evidence` / `examples` / `failure`），不得擴張為 first-class taxonomy；其餘責任一律降為 `tags`。新需求預設加 tag，不加 type。新增第 5 個 primary type 須回 plan 重評。

**套用 pilot**：software-delivery 是 execution path，slice 只落在 `execution`、`failure`（surgical caveats）、`examples` 三型；`evidence` 型不出現在此 domain（屬 `analysis/apk`、`analysis/travel`）。這正好印證 4 個 type 能橫跨 workflow / analysis 兩層而**不需擴張**——artifact-gate、contract、traceability、bdd、embedded、backfill 等全部降為 `tags`，沒有任何一個被升成新 primary type。

## 3. Granularity 原則

slice 最小單位 = **能獨立完成一個 cognitive phase**（非 step、非 concept）。判準：載入後 agent 能完成一個自足認知階段而不需瘋狂 cross-reference。

**套用 pilot**：software-delivery 的 9 步 `execution-flow` + 12 gate `development-process` **不**逐步 / 逐 gate 切（那會 over-fragmentation），而是收斂成 **6 個生命週期 cognitive phase**（intake → contracts → test-strategy → implementation → validation → closure）+ 1 個跨階段紀律 caveat（surgical）+ 1 個 examples。條件性子流程（embedded / hardware、backfill for existing project）**不另開 slice**，而以 `tags: domain-specific` 掛在對應 phase slice，避免 route fragmentation。

## 4. 三層邊界規則 + placement 可驗證 predicate

- `workflow` = 「要做什麼順序」；`analysis` = 「如何取得與驗證證據」；`intelligence` = 「為何這種模式長期有效 / 失敗」。
- **Extraction direction（單向）**：analysis → intelligence；intelligence 只接受 validated repeated patterns。
- **Falsifiable membership predicate**（歸層不是 honor-system 標籤）：
  - **workflow membership test**：內容規定「做什麼、什麼順序、過哪些 gate」，是 procedure / ordering / gate；不承載證據取得方法，也不論證長期模式。
  - **analysis membership test**：回答「如何取得 / 驗證證據」，task-instance 級 observation/signal/evidence，**不得**斷言跨實例通則。
  - **intelligence membership test**：是一個 generalization，**且** `evidence_refs` 含 ≥2 個獨立、已驗證、可解析來源；不足 → premature promotion → 強制退回 analysis。
  - 限制：無完全機械 oracle；目標是「misplacement 可偵測、可逆、便宜修正」，非「證明每次放對」。
- **套用 pilot**：8 個 pilot slice 全數通過 **workflow** membership（都是 order / gate / 紀律），無一是 evidence 取得方法或長期模式論證 → 全 `owner_layer: workflow`。唯一灰區是 `sd-closure` 的「Feed Back Reusable Lessons」(execution-flow §8)：它**產生** intelligence 候選，但本身是 workflow 的閉環步驟，故留在 workflow 並標 `tags: extraction-to-intelligence`；真正的 intelligence 內容（為何某模式長期有效）不在此 slice，須另經 evidence_refs gate 升層。

## 5. Examples suppression bias 規則

`type: examples` 的 slice 預設 `default_load: false`，只在 `user_requested_examples` 或 `ambiguity_detected` 時載入（防 example-driven loading contamination / override doctrine；對應 Watch-Out Wall 5）。

**套用 pilot**：`sd-examples`（`examples/EXAMPLES.md`，528 行、token 密度高）標 `default_load: false`。execution-only / mixed 任務（Scenario A/C）的 `forbidden_load` 必含 `sd-examples`，除非 user 明確要範例或偵測到 ambiguity。

---

## 6. 命名 / glossary 決定（Phase 1 resolved）

- **決定**：過渡期 operational wording 採 `execution surface` / `evidence surface`（pilot 文件用語）。**不在 Phase 1 註冊 `Cognitive Slice` 到 glossary**，亦不鎖定 `capability surface` / `cognitive surface` 任一候選為 framework vocabulary。
- **理由**：`slice` 易被聯想成 arbitrary chunk / static partition，但其本質是 routable cognition surface；命名屬難逆轉決定，須等 Phase 4 validation 證明 taxonomy 穩定後再定。本檔內部仍用「slice」作 working term，但對外文件用 `surface`。
- **正式 glossary 註冊**：延後至 Phase 4 validation 之後（對應 plan §Open Questions glossary 條目）。

---

## 7. Software-delivery pilot slice 盤點（Phase 1 taxonomy，尚未實體拆檔）

> 所有 slice 留在既有 `workflow/software-delivery/` owner layer（不新增 `slices/` 子目錄，對應 plan §Open Questions resolved）。`dependency_budget` 全採 default 2/4，未宣告 high override。`canonical_source` 為 Phase 2 拆檔前的 heading 範圍對映。

| id | type | tags | load_when | do_not_load_when | canonical_source（現況 heading） |
|---|---|---|---|---|---|
| `sd-intake` | execution | requirements, parity, intake, domain-specific（backfill） | 接收新需求 / 變更 / 重構意圖、需求認知盤點、product brief 驗證、既有專案回填 | 已有明確 contract、純執行既定改動 | **`intake.md`（Phase 2 已實體拆檔，跨檔同批：原 execution-flow §1 + §6 Backfill + development-process §Initial Doc Pack / §Product Brief Validation Gate / §Change Intake Gate / §Missing Information Gate / §Existing Project Documentation Backfill）** |
| `sd-contracts` | execution | artifact-gate, contract, traceability | 需建立 / 治理 contract 與可追溯性 | 無 contract 異動的小改 | **`contracts.md`（Phase 2 已實體拆檔，原 development-process §Required Contracts / Contract Governance / Traceability / Contract-First Rules）** |
| `sd-test-strategy` | execution | artifact-gate, test, bdd | 定義測試策略 / BDD 閉環 / test-first ordering | 不涉測試設計的純文件改動 | **`test-strategy.md`（Phase 2 已實體拆檔，跨檔同批：原 execution-flow §2 + §4 子節「測試策略定義」+「Test-First Ordering」+ development-process §BDD Execution Closure + §Test Strategy Gate 含 Mutation Testing）** |
| `sd-implementation` | execution | execution-order, domain-specific（embedded） | 實際進行程式碼變更（核心執行順序） | evidence-only / 純分析任務 | **暫留 `execution-flow.md` §3/§4 + `development-process.md` embedded / producer-consumer fallback**。依 stakeholder 2026-05-30 決定，Phase 3/4 用 routing / validation evidence 判斷是否需要獨立 `implementation.md`；目前不拆以避免 over-fragmentation。 |
| `sd-surgical-caveats` | failure | caveat, surgical, diff-purity | 進行外科手術式小改、需控制 diff 純度 / orphan | 大型新功能初始實作 | **`surgical-changes.md`（Phase 2 已實體拆檔，原 execution-flow §9.1–9.5）** |
| `sd-validation` | execution | artifact-gate, validation, performance | 驗證變更 / 效能關卡 | 尚未實作完成前 | **`validation.md`（Phase 2 已實體拆檔，原 execution-flow §5 Perf Gate + §7 Validate）** |
| `sd-closure` | execution | closure, handoff, extraction-to-intelligence | 收尾、DoR/DoD 檢核、回饋可重用課程 | 任務中段 | **`closure.md`（Phase 2 已實體拆檔，原 execution-flow §8 + development-process §DoR / §DoD）** |
| `sd-examples` | examples | (default_load:false) | user 明確要求範例 / 偵測到 ambiguity | 預設一律 suppress（execution-only / mixed） | examples/EXAMPLES.md |

**layer_justification（全 slice 共通）**：每條都規定「做什麼 / 什麼順序 / 過哪些 gate」，通過 workflow membership test；無一承載 evidence 取得方法（非 analysis）或長期模式論證（非 intelligence）。`sd-closure` 的 extraction-to-intelligence 僅為候選標記，升 intelligence 須補 `evidence_refs`≥2。

**條件性子流程（不另開 slice）**：Embedded / Hardware Product Flow → 掛 `sd-intake` + `sd-implementation` 的 `tags: domain-specific,embedded`；Backfill for existing project（development-process §Backfill、execution-flow §6）→ 掛 `sd-intake` 的 `tags: domain-specific,backfill`。

---

## 8. Phase 4 test-first fixtures（草稿，待 Phase 4 執行）

> 形狀對齊 plan §Phase 4。斷言：`expected_load` ⊆ loaded、`forbidden_load` ∩ loaded = ∅、載入深度/廣度未超 `dependency_budget`。驗證須檢查**實際載入的 surface**，非僅 route 存在。

```yaml
# Scenario A — execution-only：小型 API validation 變更
scenario: A-execution-only
task_intent: "為既有 API 加一個輸入驗證，無 contract / 測試策略變動"
expected_load: [sd-implementation, sd-validation]
forbidden_load: [sd-examples, sd-intake, sd-contracts, analysis/**, intelligence/**]
dependency_budget: { default: { max_depth: 2, max_runtime_dependencies: 4 } }

# Scenario B — evidence-only：分析 APK 網路行為（analysis 層）
scenario: B-evidence-only
task_intent: "分析某 APK 的網路流量行為"
expected_load: ["analysis/apk/<evidence-acquisition surface>"]
forbidden_load: [sd-intake, sd-contracts, sd-test-strategy, sd-implementation, sd-validation, sd-closure, sd-examples]
dependency_budget: { default: { max_depth: 2, max_runtime_dependencies: 4 } }

# Scenario C — mixed：debug 失敗的 deployment pipeline
scenario: C-mixed
task_intent: "deployment pipeline 失敗，需同時看執行步驟與失敗證據"
expected_load: [sd-validation, sd-surgical-caveats, "analysis/<failure-caveat surface>"]
forbidden_load: [sd-examples, intelligence/**, "其他 domain slice"]
dependency_budget: { override_when: { task_complexity: high } }  # 高複雜任務，允許放寬至 depth3/deps6

# Scenario D — placement / misplacement 負向驗證
scenario: D-misplacement
task_intent: "嘗試把一條無 evidence 的 heuristic 標成 intelligence"
assert:
  - "placement predicate 擋下：evidence_refs < 2 → 強制退回 analysis"
  - "正確的 analysis 證據 slice 通過 analysis membership test"
  - "若誤標 slice 仍存在，會在 Scenario B/C 的 forbidden_load 洩漏（contamination 探針）"
```
