# Bootstrap Contract YAML Migration

**Status**: `draft`
**世代**：Gen 3 子系統演進
**建立日期**：2026-05-25
**最後更新**：2026-05-25

> 本 plan 回應 2026-05-25 觀察：CORE_BOOTSTRAP.md 是 prose-only 結構，與系統其他 critical contracts（knowledge-update-flow / cognitive-modes / system-upgrade-governance / 各 enforcement）全部已 YAML 化的方向相反。Bootstrap 是最 critical 的契約，反而是最不結構化的。把它升級為 executable YAML + companion markdown，並讓 4 個 AI tool 的 entry points（CLAUDE.md / .cursor / .roomodes / **AGENTS.md**）一致指向 machine-readable contract。

---

## Decision Rationale

### Problem & Why Now

**Inverted priority**：系統把次要的 contracts YAML 化，最 critical 的 CORE_BOOTSTRAP 卻是 prose。具體後果：

1. **Per-turn obligation 無法可靠 fire**
   - 2026-05-26 測試證實：CLAUDE.md / CORE_BOOTSTRAP.md 的 `### Cognitive Mode 報告` per-turn obligation **無法強制**，agent 完成 Phase 1 work 後 final report 沒輸出 block
   - 即使加 IMPORTANT 框線、Fix A/B 都做了，prose obligation 仍未形成 forcing function
   - 對比：Bootstrap Receipt 之所以 work 是因為它要 query runtime.db 才寫得出來（machine-checkable）

2. **Obligations 不可機讀**
   - CORE_BOOTSTRAP.md 列了 Bootstrap Receipt + Cognitive Mode block 兩個 obligations，但散在 prose 中
   - Validator 沒辦法 iterate obligations 並 individual check
   - 新 obligation 要加：只能再寫一段 prose，又一條無法強制

3. **Entry files 漏列 AGENTS.md**
   - `bootstrap-entry-points.yaml` 列出 CLAUDE.md / .cursor / .roomodes 三條
   - **遺漏 AGENTS.md**（Codex 入口）→ thinness validator 對 AGENTS.md 完全沒檢查
   - Codex 用戶 fork 後加 obligation 到 AGENTS.md 不會被擋

4. **既有 YAML pattern 已驗證**
   - `knowledge-update-flow.yaml` + `.md` companion 模式 work
   - `cognitive-modes-*.yaml` (7 個) 全部 machine-iterable
   - CORE_BOOTSTRAP 是唯一例外

**Why now**：v2 cognitive-contract plan 才開到 Phase 1。如果先做 bootstrap YAML 化，v2 後續 phases 可直接用結構化 obligations 設計；反之 v2 完成後再做 YAML 化會多一次 migration。

### Decision

把 CORE_BOOTSTRAP.md 升級為 **executable YAML + companion markdown**，並讓所有 entry points 一致：

#### 1. CORE_BOOTSTRAP.yaml schema

```yaml
id: runtime.core-bootstrap
runtime_projection:
  enabled: true
  target_key: runtime.core_bootstrap.contract
  surface: generated_surfaces

required_reads:
  - id: rule-weight
    path: enforcement/rule-weight.md
    role: P0/P1/P2/P3 weight system
    estimated_tokens: 300
  - id: dependency-reading
    path: enforcement/dependency-reading.md
    role: dependency 鐵則 + writeback transaction
    estimated_tokens: 400
  - id: conversation-goal-ledger
    path: enforcement/conversation-goal-ledger.md
    role: .agent-goals/ 使用
    estimated_tokens: 100

per_session_obligations:
  - id: obligation.bootstrap.receipt
    fires: first_turn_before_any_non_read_tool
    format_template: "Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>"
    enforcement_layer: behavioral_forcing
    forcing_function: numeric_values_require_runtime_db_query
    severity: high
    blocking_gate_id: gate.bootstrap.receipt_present
    obligation_ledger_id: obligation.bootstrap.receipt_acknowledged

per_turn_obligations:
  - id: obligation.cognitive.mode_report
    fires: every_user_facing_response
    format_options:
      compact: "Cognitive: <e>·<c>·<g>·<m> / V:<v> / Cost:<cost> / Sig:<signal>"
      full: "### Cognitive Mode 報告 + 4-6 row table"
    trigger_full_when: any_dim_non_default OR high_risk
    enforcement_layer: commit_msg_hook + behavioral_per_turn
    severity: high

per_commit_obligations:
  - id: obligation.commit.cognitive_mode_block
    fires: every_commit
    validator: validateCognitiveContractFormat
    severity: block

phase_state_init:
  source: runtime/runtime.db
  required_tables:
    - phase_machine
    - obligations
    - gates
    - language_policy
    - output_rules
    - governance_gates

companion_markdown: CORE_BOOTSTRAP.md
companion_role: human-readable rationale + examples; YAML is canonical
```

#### 2. CORE_BOOTSTRAP.md 變 companion

像 `knowledge-update-flow.md` 對 `.yaml` — 解釋背景、提供範例，但 canonical contract 在 YAML。所有舊 link 不動（向後相容）。

#### 3. AGENTS.md 加進 bootstrap-entry-points.yaml

把 AGENTS.md 補進 entry_files，bootstrap-entry-thinness validator 自動覆蓋。

#### 4. Bootstrap Receipt 升級

Phase 5 task：Receipt 額外列 active per-turn obligations：

```
Bootstrap: rules=✓ phase=phase.bootstrap obligations=2 gates=3
Active per-turn obligations: cognitive.mode_report
```

Agent first turn 就看到自己這 session 要遵守哪些 per-turn obligations。

#### 5. 全部 4 個 entry files 改 query runtime.db / generated_surfaces

每個 entry 改為：
> 啟動時讀 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md)（companion） + query `runtime/runtime.db generated_surfaces[target_key=runtime.core_bootstrap]` 取得 machine-readable obligations list。

### Alternatives Considered

- **A. 維持 prose CORE_BOOTSTRAP.md（do nothing）**：reject — 已證實 per-turn obligation 無法強制；架構不一致持續
- **B. 純 prose 改強：把 obligation 移到 .md 頭部 + 加更狠 IMPORTANT**：reject — prose 仍是 prose；不解決 machine-iterability
- **C. 把 obligations 散到既有 YAMLs**（cognitive-modes.yaml + bootstrap-entry-points.yaml etc）：reject — 沒有單一 canonical entry，agent 不知該 query 哪個 surface
- **D. CORE_BOOTSTRAP.yaml + .md companion（current decision）**：accept — 對齊既有 pattern（knowledge-update-flow / system-upgrade-governance）
- **E. 直接寫 ADR-009「YAML-first runtime contracts」**：defer — 本 plan 完成後評估，若 promotion criteria 達標再 promote

### Why Not an ADR Yet

- 設計範圍涵蓋 7 phases + 4 tools，可能在 Phase 2-5 發現需要調整 schema
- 既有 7 個 YAML contracts 沒有跟 CORE_BOOTSTRAP 同等 critical，schema 是否完全照搬未知
- Bootstrap Receipt enhancement 是否真的提升 obligation 內化率需要多 session evidence

### ADR Promotion Criteria（completed 時驗證）

- [ ] foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中
- [ ] Plan 結果證實 YAML migration + 4 entry points 一致 + Receipt enhancement 可行
- [ ] Open Questions 全解
- [ ] 沒有更輕的 promotion target 適用
- [ ] 系統真實使用：至少 3 個新 session 透過 runtime.db 查 obligation list（非 docs 路徑）
- [ ] Per-turn obligation 內化率上升（量化：新 session 的 final response 含 Cognitive Mode block 比例 > v1 baseline）

### Consequences（預期）

#### 正面

- **Machine-iterable obligations**：validator 可 individual check 每條 obligation
- **Cross-tool 一致**：4 個 entry points 都指向同一 generated_surface
- **AGENTS.md 進防護網**：thinness 規則覆蓋 Codex 入口
- **Receipt enhancement**：agent first turn 看到 active obligations 名單
- **Backward compat**：舊 link 全部保留（CORE_BOOTSTRAP.md 變 companion）

#### 負面

- **Migration cost**：CORE_BOOTSTRAP 拆 YAML + companion、4 個 entry files 各更新、bootstrap-entry-points.yaml 加 AGENTS.md
- **Schema 演進負擔**：obligations 新增需同時改 YAML + projection + companion

#### 風險

- **Receipt enhancement 不一定有效**：列出 obligation 名單 ≠ agent 自願遵守 → mitigation：accept partial gain，commit-msg hook 仍是強制基線
- **與 v2 plan 重疊**：v2 也在改 cognitive-mode 報告 schema → mitigation：本 plan 不動 cognitive-mode schema 內部（v2 主管），只負責讓 cognitive-mode obligation 變 machine-iterable entry

---

## Runtime Execution Path

| 欄位 | 內容 |
|---|---|
| Runtime owner | `CORE_BOOTSTRAP.yaml`（新）+ `runtime/bootstrap-entry-points.yaml`（加 AGENTS.md）+ `scripts/ai-skill-cli/internal/app/hooks.go`（Receipt validator 升級） |
| Trigger flow | session start → entry file（CLAUDE.md/.cursor/.roomodes/AGENTS.md）→ 點 CORE_BOOTSTRAP.md companion + query runtime.db generated_surfaces[runtime.core_bootstrap] → 列舉 obligations → Bootstrap Receipt 附 active per-turn obligations → commit-msg hook 依 obligation list 個別 check |
| Trigger location | session bootstrap (first turn) + commit-msg hook |
| Activation contract | `runtime.core_bootstrap.contract` projected to `generated_surfaces` |
| Generated surface | `runtime.core_bootstrap.contract`（含 obligations list） + 既有 `runtime.bootstrap_entry_points`（含 AGENTS.md 後） |
| Validation scenarios | `bootstrap-yaml-contract-exists-v1`、`bootstrap-yaml-projected-v1`、`bootstrap-obligations-enumerable-v1`、`agents-md-in-entry-list-v1`、`receipt-includes-active-obligations-v1`、`all-tool-entries-thin-v1` |
| Test passing evidence | 全部 scenarios PASS + ≥3 新 session evidence of querying runtime.db obligations + per-turn obligation 內化率量化提升 |

---

## Open Questions

| # | Question | Status |
|---|---|---|
| 1 | YAML schema：obligations 是否需 phase scoping（some obligation only fires in certain phase）？ | TBD Phase 2 |
| 2 | Receipt enhancement format：完整列名 vs 只列數量？token cost vs visibility tradeoff | TBD Phase 5 |
| 3 | Agent 真的會 query generated_surfaces 嗎？還是直接讀 .md companion？需要 Receipt 提示？ | TBD Phase 4 |
| 4 | Lazy-load rules（9 條）也要結構化還是只 Phase 1 處理必讀 3 條 + 2 個 obligations？ | TBD Phase 2 |
| 5 | 既有 CORE_BOOTSTRAP.md prose obligations 與 v2 cognitive-contract plan 重疊區的所有權邊界？ | TBD Phase 2 |

---

## 完成條件

### 計畫書本身

- [ ] 計畫書狀態：`draft` → `in-progress`（Phase 0 通過後）→ `completed`
- [ ] 5 Open Questions 全部 resolved
- [ ] Phase 7 close-loop 完成 ADR 決議

### Behavioral evidence

- [ ] ≥3 新 session 觀察到 Bootstrap Receipt 含 active obligations 列表
- [ ] Per-turn obligation 內化率：≥3 個新 session 的 final response 含 Cognitive Mode block（base rate 對比目前 0/2 = 0%）

### Validation

- [ ] 全部新 scenarios PASS
- [ ] 既有 ≥21 scenarios 仍 PASS（不 regress）

---

## Phase 0 Pre-Build Interrogation

### 目的

驗證 schema 設計與既有 7 個 YAML contracts 相容；確認 4 個 tool entry points 與 init_project.go 產出可一致改造。

### Tasks

- [ ] 讀全部 7 個既有 cognitive-modes-*.yaml schema，確認 CORE_BOOTSTRAP.yaml schema 樣式對齊
- [ ] 讀 `ai-tools/agent/{claude,cursor,roo,codex}.md` 確認 4 個工具的 entry 機制
- [ ] 確認 `scripts/ai-skill-cli/internal/app/init_project.go` 產出 4 個 tool 的 entry：CLAUDE.md / .cursor/rules/*.mdc / .roomodes / AGENTS.md
- [ ] 驗證 `runtime/bootstrap-entry-points.yaml` 加 AGENTS.md 不會 break 既有 validator
- [ ] 確認 Bootstrap Receipt format 是否要改（向後相容 vs 加 active obligations 行）

### Phase 0 完成條件

- [ ] 5 tasks all done
- [ ] No breaking-change conflict
- [ ] 若有 conflict 更新 §Decision 或加 Open Question

---

## Phase 1 Test-First Validation

### Tasks

- [ ] Scenario `bootstrap-yaml-contract-exists-v1` — CORE_BOOTSTRAP.yaml 存在 + schema_version + runtime_projection.enabled
- [ ] Scenario `bootstrap-yaml-projected-v1` — generated_surfaces 含 target_key=runtime.core_bootstrap.contract
- [ ] Scenario `bootstrap-obligations-enumerable-v1` — sqlite query 可取得 obligations 列表
- [ ] Scenario `agents-md-in-entry-list-v1` — `runtime/bootstrap-entry-points.yaml` entry_files 含 AGENTS.md
- [ ] Scenario `receipt-includes-active-obligations-v1` — commit-msg hook OR Receipt validator 可檢查 Receipt 含 active obligations 行（softness：可能是建議而非強制）
- [ ] Scenario `all-tool-entries-thin-v1` — 4 個 entry points（含 AGENTS.md）都通過 thinness check

### Phase 1 完成條件

- [ ] 6 scenarios 寫好且 initial state = FAIL
- [ ] Atomic test-first commit

---

## Phase 2 CORE_BOOTSTRAP.yaml Schema And Migration

### Tasks

- [ ] 寫 `CORE_BOOTSTRAP.yaml`（per §Decision schema）
- [ ] CORE_BOOTSTRAP.md prose 改為 companion role：保留說明性段落，標明 canonical contract 在 YAML
- [ ] 確認 compiler 自動 walk `runtime/` 也應加 walk repo-root（看 sourceRoots 是否要加新 path，或把 CORE_BOOTSTRAP.yaml 放 `runtime/core-bootstrap.yaml`）
- [ ] 決議 Open Question 4：lazy-load rules 範圍

### Phase 2 完成條件

- [ ] YAML 通過 `ai-skill runtime compile + validate`
- [ ] Scenarios `bootstrap-yaml-contract-exists-v1` + `bootstrap-yaml-projected-v1` PASS

---

## Phase 3 Obligations Enumeration

### Tasks

- [ ] 確認 YAML 內 per_session_obligations / per_turn_obligations / per_commit_obligations 在 generated_surfaces JSON 內可被 SQL query 解析（用 SQLite JSON1 函數）
- [ ] 寫 helper SQL/CLI：`ai-skill bootstrap obligations --phase <id>` 列出當前 active obligations
- [ ] Scenario `bootstrap-obligations-enumerable-v1` PASS

### Phase 3 完成條件

- [ ] Obligations 可從 runtime.db 個別取出
- [ ] CLI helper 通過

---

## Phase 4 Entry Files Audit & AGENTS.md Inclusion

### Tasks

- [ ] 更新 `runtime/bootstrap-entry-points.yaml` 加 AGENTS.md（含 description + tool=Codex）
- [ ] 更新 `validateBootstrapEntryThinness` 對 AGENTS.md 同 30-line + forbidden_substring check
- [ ] 4 個 entry files 內容檢查：CLAUDE.md / `.cursor/rules/*.mdc` / `.roomodes` / AGENTS.md 是否都符合 thinness + 都 point to CORE_BOOTSTRAP companion + 都提示 query runtime.db generated_surfaces
- [ ] 更新 `init_project.go` 對 4 個下游 entries 的 template 一致
- [ ] Scenarios `agents-md-in-entry-list-v1` + `all-tool-entries-thin-v1` PASS

### Phase 4 完成條件

- [ ] 4 個 entry points 通過 thinness
- [ ] init_project.go 產出 4 個 thin entries
- [ ] AGENTS.md 進 bootstrap-entry-points

---

## Phase 5 Bootstrap Receipt Enhancement

### Tasks

- [ ] 決議 Open Question 2：active obligations 列名 vs 列數
- [ ] 更新 CORE_BOOTSTRAP.yaml format_template 加 `Active per-turn obligations: ...` 行
- [ ] 更新 Bootstrap Receipt 範例
- [ ] commit-msg hook validator 可接受新舊兩種 Receipt format（backward compat）
- [ ] Scenario `receipt-includes-active-obligations-v1` PASS

### Phase 5 完成條件

- [ ] Receipt 升級 backward-compat
- [ ] Scenario PASS

---

## Phase 6 Validators And Per-obligation Enforcement

### Tasks

- [ ] hook 從 `generated_surfaces[runtime.core_bootstrap.contract]` 取 obligations list
- [ ] 對每條 per_commit_obligations 個別 dispatch validator（取代目前 hardcode list）
- [ ] 既有 6 個 validators 保留為 per-obligation handlers
- [ ] Failure pattern `bootstrap-yaml-bypass.md`（agent 跳過 query generated_surfaces 直接讀 .md prose）

### Phase 6 完成條件

- [ ] Hook 從 runtime.db 動態載入 obligation list
- [ ] Failure pattern 文件

---

## Phase 7 Close-Loop & ADR Decision

### Tasks

- [ ] 全部 scenarios PASS
- [ ] 累積 evidence：≥3 新 session 觀察 Receipt 含 active obligations
- [ ] 量化 per-turn obligation 內化率（base = 0/2，目標 ≥ 3/3）
- [ ] 評估 ADR Promotion Criteria
- [ ] 決議：(a) 寫 ADR-009「YAML-first runtime contracts」 / (b) amend ADR-008 add runtime contract surface / (c) plan-only
- [ ] Plan status → completed，移到 plans/archived/
- [ ] Plan Completion Closure

### Phase 7 完成條件

- [ ] Scenarios PASS
- [ ] Behavioral evidence 累積達標
- [ ] ADR decision recorded

---

## Stakeholder 同意項目

- [ ] User confirms: CORE_BOOTSTRAP.yaml schema 含 required_reads / per_session / per_turn / per_commit obligations
- [ ] User confirms: 把 CORE_BOOTSTRAP.md 降級為 companion 是可接受的（vs deprecate）
- [ ] User confirms: AGENTS.md 需補進 bootstrap-entry-points
- [ ] User confirms: Bootstrap Receipt 升級為含 active obligations（可能加 30-50 tokens / session）
- [ ] User confirms: 4 個 entry files 都改 query runtime.db 而非 dock to prose

---

## 與其他 plans 的關係

| Plan | 關係 |
|---|---|
| [`active/2026-05-25-2100-runtime-cognitive-contract-v2.md`](2026-05-25-2100-runtime-cognitive-contract-v2.md) | v2 設計 cognitive mode schema 內部；本 plan 提供 CORE_BOOTSTRAP YAML 容器讓 cognitive-mode obligation 成為 enumerable item。兩 plan 互補；v2 Phase 2 + 本 plan Phase 3 可能需要協調 schema |
| [`active/2026-05-25-1000-context-language-glossary-system.md`](2026-05-25-1000-context-language-glossary-system.md) | independent；glossary 是 ubiquitous language，bootstrap 是 obligations |
| [`archived/2026-05-22-1629-runtime-cognitive-modes-system.md`](../archived/2026-05-22-1629-runtime-cognitive-modes-system.md) | v1 cognitive modes plan；本 plan 把它的 obligations machine-iterable 化 |
| [`archived/2026-05-22-0855-executable-yaml-contract-migration.md`](../archived/2026-05-22-0855-executable-yaml-contract-migration.md) | parent pattern — 本 plan 是同 pattern 套到 CORE_BOOTSTRAP |
| [`constitution/ADR-008-runtime-cognitive-modes.md`](../../constitution/ADR-008-runtime-cognitive-modes.md) | obligations enumeration 可能影響 ADR-008 補強 |
