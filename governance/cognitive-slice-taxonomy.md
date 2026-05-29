# Cognitive Slice / Surface Taxonomy（Phase 1 skeleton）

**Status**: `skeleton`（Phase 1 down-payment；尚未套用 pilot、尚未草擬 validation fixture）
**Owner layer**: governance
**來源 plan**: [`plans/active/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md`](../plans/active/2026-05-29-0916-gen3-workflow-analysis-cognitive-slice-decomposition.md) §Phase 1
**命名注意**：framework vocabulary（`Cognitive Slice` vs `surface`）正式註冊**刻意延後至 Phase 4 validation**；本檔過渡期一律用 `loading / execution / evidence surface` 措辭，不把 `Cognitive Slice` 當已確立詞彙散播。

> 本檔目前是**結構骨架**：定義 slice schema 欄位與 4 條治理規則的位置，**內容值（套用到 software-delivery pilot 的實際分類、fixture 草稿）留待 Phase 1 完整 session 補**。骨架本身不分類任何既有檔案。

---

## 1. Slice schema 欄位（spec）

| 欄位 | 必填 | 說明 | 值 |
|---|---|---|---|
| `purpose` | 是 | 此 slice 要讓 agent 完成的認知目的 | TODO Phase 1 |
| `type` | 是 | primary，**只允許 4 值**：`execution` / `evidence` / `examples` / `failure` | TODO Phase 1 |
| `tags` | 否 | secondary 自由標註（artifact-gate / closure / handoff / templates / observation-triage / tool-procedure / domain-specific / extraction-to-intelligence …） | TODO Phase 1 |
| `load_when` | 是 | 何種 task intent 應載入 | TODO Phase 1 |
| `do_not_load_when` | 是 | 何種任務不應載入（suppression） | TODO Phase 1 |
| `owner_layer` | 是 | workflow / analysis / intelligence（依三層邊界規則判定） | TODO Phase 1 |
| `layer_justification` | 是 | 歸層的 falsifiable 理由，須通過該層 membership predicate（規則 4） | TODO Phase 1 |
| `evidence_refs` | intelligence 必填 | ≥2 個獨立、已驗證、可解析的 analysis 觀察 / failure case 指標 | TODO Phase 1 |
| `canonical_source` | 是 | 正文 canonical 來源（slice 只導航，不重定義） | TODO Phase 1 |
| `dependencies` | 否 | 依賴的其他 slice / source | TODO Phase 1 |
| `dependency_budget` | 是 | heuristic default `max_depth:2`/`max_runtime_dependencies:4` + `override_when: task_complexity=high`（非 rigid） | TODO Phase 1 |
| `summary_path` | 否 | 對應 summary-first 入口 | TODO Phase 1 |
| `validation_signal` | 是 | Phase 4 用哪個 scenario 驗證 | TODO Phase 1 |

---

## 2. type+tags 收斂規則

primary `type` 固定 4 種（`execution` / `evidence` / `examples` / `failure`），不得擴張為 first-class taxonomy；其餘責任一律降為 `tags`。新需求預設加 tag，不加 type。新增第 5 個 primary type 須回 plan 重評。

- 套用到 pilot：TODO Phase 1

## 3. Granularity 原則

slice 最小單位 = **能獨立完成一個 cognitive phase**（非 step、非 concept）。判準：載入後 agent 能完成一個自足認知階段而不需瘋狂 cross-reference。

- 套用到 pilot：TODO Phase 1

## 4. 三層邊界規則 + placement 可驗證 predicate

- `workflow` = 「要做什麼順序」；`analysis` = 「如何取得與驗證證據」；`intelligence` = 「為何這種模式長期有效 / 失敗」。
- **Extraction direction（單向）**：analysis → intelligence；intelligence 只接受 validated repeated patterns。
- **Falsifiable membership predicate**（歸層不是 honor-system 標籤）：
  - **analysis membership test**：回答「如何取得 / 驗證證據」，task-instance 級 observation/signal/evidence，**不得**斷言跨實例通則。
  - **intelligence membership test**：是一個 generalization，**且** `evidence_refs` 含 ≥2 個獨立、已驗證、可解析來源；不足 → premature promotion → 強制退回 analysis。
  - 限制：無完全機械 oracle；目標是「misplacement 可偵測、可逆、便宜修正」，非「證明每次放對」。
- 套用到 pilot：TODO Phase 1

## 5. Examples suppression bias 規則

`type: examples` 的 slice 預設 `default_load: false`，只在 `user_requested_examples` 或 `ambiguity_detected` 時載入（防 example-driven loading contamination / override doctrine；對應 Watch-Out Wall 5）。

- 套用到 pilot：TODO Phase 1

---

## 6. 命名 / glossary 決定

- 過渡期 operational wording：`loading / execution / evidence surface`。
- 候選評估（`capability surface` / `cognitive surface` / `execution surface` vs `slice`）：TODO Phase 1（記錄理由，不在此鎖定 framework vocabulary）。
- 正式 glossary 註冊：延後至 Phase 4 validation 之後。

## 7. 延後項（不在本骨架）

- software-delivery pilot 的實際 slice 清單與每個 slice 的 schema 值 → Phase 1 完整 session。
- Phase 4 test-first fixture（`expected_load` / `forbidden_load` / Scenario A/B/C/D）草稿 → Phase 1/2。
