# Context Language Glossary System

**Status**: `draft`
**世代**：Gen 3 子系統擴充
**建立日期**：2026-05-25
**最後更新**：2026-05-25（建立 draft plan；尚未開始 implementation）

> 本 plan 回應外部 `CONTEXT.md / ubiquitous language` 建議，參考 [mattpocock/skills](https://github.com/mattpocock/skills/tree/main) 將 shared language 作為 agent alignment 技巧；但在 Ai-skill 中不建立 root `CONTEXT.md` 作為第二份 source-of-truth，而是依 Gen 3 分層落到 `knowledge/glossary/` 與 project-local memory 邊界。

---

## Decision Rationale

### Problem & Why Now

Matt Pocock 的 skills 將 `CONTEXT.md` / shared language 視為降低 agent 誤解與冗詞的重要能力：agent 先學會專案語彙，後續 plan、code、review 都能用同一組概念溝通。

Ai-skill 目前已有相關 intelligence：

- `intelligence/engineering/requirements/behavior-modeling/ubiquitous-language-alignment.md`
- `intelligence/engineering/architecture/domain-modeling/ubiquitous-language.md`
- `workflow/software-delivery/requirements/behavior-driven-discovery/README.md`

但目前缺少一個明確的 canonical glossary layer。若直接新增 `CONTEXT.md` 或 `memory/project/context-language.md`，會產生下列風險：

1. **Source-of-truth duplication**：同一詞彙可能同時存在於 root context、knowledge atom、memory replay、requirements docs。
2. **Memory 被誤當 current truth**：`memory/README.md` 明確規定 memory 不保存 canonical source / current truth。
3. **Workflow 找不到可查 source**：behavior-driven discovery 目前只說「對齊 shared language」，但沒有穩定路徑可讀。
4. **框架詞彙漂移**：近期 active plan 已出現 `context_mode` / `compression` 詞彙重疊風險；需要更明確的 glossary owner。

### Decision

建立 **Context Language Glossary System**：

| 類型 | Canonical location | 說明 |
| --- | --- | --- |
| Ai-skill 全庫 / 框架 / 可重用語彙 | `knowledge/glossary/` | 當前 canonical glossary，供 workflow、routing、architecture、requirements 引用。 |
| 單一專案、跨 session 但非 canonical 的語言脈絡 | `<PROJECT_ROOT>/memory/project/context-language.md` 或專案等價文件 | 只作 selective replay，不取代專案正式 docs 或 Ai-skill glossary。 |
| 判斷何時需要 shared language | `intelligence/engineering/requirements/behavior-modeling/` 與 `intelligence/engineering/architecture/domain-modeling/` | 保留為 reasoning source。 |
| 執行何時讀 glossary | `workflow/software-delivery/requirements/behavior-driven-discovery/` 與 `pre-build-interrogation` | Workflow gate 觸發 glossary 使用。 |

`knowledge/glossary/` 不做百科全書；只收會影響 behavior、contract、runtime surface、validation、routing 或 framework decision 的詞。

### Alternatives Considered

- **A. 建 root `CONTEXT.md`**：拒絕。Root context 容易變成 always-loaded 大檔，且會與 `knowledge/`、`workflow/`、`memory/` 形成平行 source。
- **B. 放 `memory/project/context-language.md` 作 canonical**：拒絕。Memory 是 selective replay，不是 current truth；可作 project-local replay，但不可作 Ai-skill canonical。
- **C. 只保留 intelligence，不新增 glossary**：拒絕。Intelligence 只回答何時需要 ubiquitous language，不提供 agent 可查的穩定詞彙 source。
- **D. 建 `knowledge/glossary/` 作 canonical，memory 只作 project replay**：接受。

### Why Not an ADR Yet

依 `plans/README.md` 與 ADR-007：

- 目前是可逆的 layer / workflow 擴充，尚未證明需要 constitution 級不可逆決策。
- 是否需要 runtime projection、validation scenarios、project-local template，仍需透過本 plan 驗證。
- 若 completed 後證明 glossary 成為 foundational + cross-session + cross-project + expensive-to-reverse 的 architecture boundary，再評估 ADR promotion。

### ADR Promotion Criteria（completed 時驗證）

- [ ] `knowledge/glossary/` 已被 workflow 或 routing 實際使用。
- [ ] `memory/project/context-language.md` 邊界未被誤用為 canonical source。
- [ ] 至少一個 validation scenario 能防止 glossary source duplication 或 memory-as-truth。
- [ ] Open Questions 全解。
- [ ] 沒有更輕的 promotion target 適用（knowledge README / workflow docs / validation scenario 足夠時不升 ADR）。

### Consequences（預期）

#### 正面

- Agent 有穩定 shared language source，不必從多處 inference。
- Framework 詞彙如 `context_mode`、`generated_surfaces`、`owner-layer contract` 可集中定義。
- Requirements / BDD / domain modeling 可引用同一 glossary，降低語義漂移。
- Project memory 與 canonical knowledge 邊界清楚。

#### 負面

- 新增 `knowledge/glossary/` 需要 routing、graph、summary 與 README 連動。
- 若 glossary 過度膨脹，會變成低價值百科，需要嚴格 scope。
- Project-local glossary template 需要說清楚「專案 source 優先」，避免污染 reusable docs。

#### 風險

| 風險 | 緩解 |
| --- | --- |
| Glossary 變成百科 | 只收影響 behavior / contract / validation / runtime decision 的詞；其他詞留在 project docs。 |
| Memory 被當 canonical | 在 memory README、glossary README 與 validation scenario 明確禁止。 |
| 與 intelligence/domain-modeling 重複 | Intelligence 保留判斷原則；glossary 保存穩定詞條。 |
| 與 active runtime cognitive modes plan 詞彙衝突 | Phase 0 先盤點 `context_mode` / `compression` / `memory_mode`，避免先寫兩套。 |

---

## Runtime Execution Path

| 欄位 | 內容 |
| --- | --- |
| Runtime owner | Phase 1-2：無，doc / knowledge layer plan；Phase 3 後視結果決定是否新增 `generated_surfaces` glossary index 或只走 knowledge runtime refresh。 |
| Trigger location | `workflow/software-delivery/requirements/behavior-driven-discovery/`、`pre-build-interrogation`、`route.workflow.software-delivery`、未來 glossary route。 |
| Activation contract | 初期無新 executable contract；使用既有 `workflow.software_delivery.pre_build_interrogation.contract` 作 plan 前 gate。若 Phase 3 決定 glossary 使用需要 blocking gate，再補 owner-layer YAML。 |
| Generated surface | 初期不投影 executable contract；`knowledge/runtime/sqlite/runtime-index.sqlite` 與 runtime reports 會收錄 glossary source。若新增 executable contract，target_key 需另定。 |
| Validation scenarios | Phase 1 先新增：`validation/scenarios/failure-derived/glossary-source-duplication-v1.yaml`、`validation/scenarios/failure-derived/memory-context-language-as-canonical-v1.yaml`。 |
| Test passing evidence | `ai-skill runtime refresh --native-index --native-reports`、`ai-skill runtime validate --json`、scenario / lints / diff review。 |

### Doc-only Trial 聲明

本 plan 初期不接入新的 runtime execute layer。理由：

- Shared language 先是 knowledge source，不是 blocking runtime gate。
- 現有 `pre-build-interrogation` 已能在 plan 前阻擋 source-of-truth duplication。
- 若未先驗證 glossary scope，直接做 runtime contract 會增加 noise。

未來接入時機：

- Phase 3 證明 agent 需要 glossary route 或 blocking gate 才能穩定使用時，新增 routing / generated surface / executable contract。

---

## Phase 0 Pre-Build Interrogation

| 欄位 | 內容 |
| --- | --- |
| Trigger | 使用者要求把 `CONTEXT.md / ubiquitous language system` 寫成計畫，並提醒架構已調整需重讀規則。 |
| Checked sources | `CORE_BOOTSTRAP.md`、`enforcement/dependency-reading.md`、`enforcement/linked-updates.md`、`workflow/software-delivery/requirements/pre-build-interrogation.md`、`plans/README.md`、`governance/lifecycle/system-upgrade-governance.md`、`architecture/README.md`、ADR-007、knowledge / memory README、ubiquitous language intelligence。 |
| Goal | 建立 shared language canonical placement plan，吸收 `CONTEXT.md` 技巧但不引入第二份 source-of-truth。 |
| Scope | Plan only；規劃 `knowledge/glossary/`、project memory boundary、workflow/routing/validation linked updates。 |
| Non-goals | 本 plan 不立即建立完整 glossary、不中斷 active runtime cognitive modes plan、不建立 proposed ADR、不新增 root `CONTEXT.md`。 |
| Acceptance | Plan 符合新 `plans/README.md` 必填章節；明確回答 `knowledge/glossary/` vs `memory/project/context-language.md`；列出 test-first validation。 |
| Framework discovery | `knowledge/` 是 navigation / atom / glossary source 候選；`memory/` 不保存 current truth；`intelligence/` 保存判斷原則；`workflow/` 保存執行 gate。 |
| Duplication risk | Root `CONTEXT.md`、memory canonical、knowledge glossary 三者不可並存為同一語彙 source。Plan 採 `knowledge/glossary/` canonical + memory replay 非 canonical。 |
| Open questions | 見下一節。 |
| Decision | proceed with draft plan；implementation 需另行啟動。 |

---

## Open Questions

1. `knowledge/glossary/` 是否先只放 Ai-skill framework glossary，還是同時放 software-delivery glossary？
2. Project-local `context-language.md` 是否由 Ai-skill 提供 template，或只寫規範讓 downstream project 自行建立？
3. Glossary 是否需要 companion YAML contract，或維持 Markdown + generated knowledge index 即可？
4. 是否需要在 active `runtime-cognitive-modes-system` plan 中加入詞彙對齊 dependency，避免 `context_mode` 與 `compression` 分叉？
5. Glossary entry 是否需要 status 欄位：`canonical` / `candidate` / `deprecated` / `project-local`？

---

## Phase 1 Test-First Validation

先建立 validation scenarios，再實作 glossary 結構。

### 期望可觀察行為

- Agent 不會建 root `CONTEXT.md` 作 Ai-skill canonical glossary。
- Agent 不會把 `memory/project/context-language.md` 當 current truth。
- Agent 在 framework / requirements plan 中能找到 `knowledge/glossary/` 作 shared language source。

### Tasks

- [ ] 新增 `validation/scenarios/failure-derived/glossary-source-duplication-v1.yaml`
- [ ] 新增 `validation/scenarios/failure-derived/memory-context-language-as-canonical-v1.yaml`
- [ ] 更新 relevant graph / summary / routing candidate（若 scenario 需要）

### Phase 1 完成條件

- [ ] Scenarios 符合 `validation/scenario.schema.json`
- [ ] `ai-skill runtime refresh --native-index --native-reports` 通過
- [ ] `ai-skill runtime validate --json` 通過

---

## Phase 2 Knowledge Glossary Structure

### Proposed files

```text
knowledge/glossary/
  README.md
  ai-skill.md
  software-delivery.md
  runtime.md
```

### Tasks

- [ ] 建立 `knowledge/glossary/README.md`，定義 glossary scope、entry shape、status、non-goals。
- [ ] 建立 `knowledge/glossary/ai-skill.md`，收 Ai-skill framework 詞彙。
- [ ] 視 Open Question 1 決定是否同批建立 `software-delivery.md` / `runtime.md`。
- [ ] 更新 `knowledge/README.md`，將 glossary 加入目前入口。
- [ ] 更新 `knowledge/graphs/`，把 glossary 連到 requirements / domain modeling / software delivery。

### Phase 2 完成條件

- [ ] Glossary 只收影響 behavior / contract / validation / runtime decision 的詞。
- [ ] 沒有把 project-specific 詞、host、class、incident evidence 寫入 reusable glossary。
- [ ] `knowledge/README.md` 可導向 glossary。

---

## Phase 3 Workflow And Memory Boundary Integration

### Tasks

- [ ] 更新 `workflow/software-delivery/requirements/behavior-driven-discovery/README.md`：shared language 對齊時先查 `knowledge/glossary/`。
- [ ] 更新 `workflow/software-delivery/requirements/pre-build-interrogation.md`：framework discovery 若有詞彙分叉，先查 glossary。
- [ ] 更新 `memory/README.md` 或 `memory/project/README.md`：project context-language 只作 replay，不作 canonical。
- [ ] 更新 `intelligence/engineering/requirements/behavior-modeling/ubiquitous-language-alignment.md`：判斷原則指向 glossary 作 source。
- [ ] 更新 `intelligence/engineering/architecture/domain-modeling/ubiquitous-language.md`：domain glossary 與 Ai-skill glossary 的邊界。

### Phase 3 完成條件

- [ ] Workflow、memory、intelligence 三層沒有互相取代。
- [ ] `memory/project/context-language.md` 的用途被限制為 project-local replay / context aid。
- [ ] Agent 能從 behavior-driven discovery 找到 canonical glossary source。

---

## Phase 4 Routing, Runtime Reports, And Generated Lookup

### Tasks

- [ ] 更新 `knowledge/runtime/routing-registry.yaml`，新增或擴充 glossary route。
- [ ] 更新 `knowledge/summaries/requirements-cognition.md` / `development-guidance.md` / related summaries。
- [ ] 執行 `ai-skill runtime refresh --native-index --native-reports`，讓 glossary 進 generated lookup。
- [ ] 若決定需要 executable contract，再補 `knowledge/glossary/*.yaml` 或 owner-layer contract；否則明確記錄 Markdown-only / knowledge source。

### Phase 4 完成條件

- [ ] Runtime reports / SQLite index 可找到 glossary source。
- [ ] 若無 runtime executable contract，已記錄原因：glossary 是 knowledge source，不是 workflow gate。
- [ ] 若有 contract，`generated_surfaces` target_key 已 synced。

---

## Phase 5 Close Loop

### Tasks

- [ ] Diff review。
- [ ] ReadLints 檢查。
- [ ] `ai-skill runtime validate --json`。
- [ ] 必要時 `go test ./...`（若 CLI/compiler/validator 變更）。
- [ ] 更新 `plans/README.md` 狀態。
- [ ] Plan completed 後搬到 `plans/archived/`，或若持續生效則標明例外。
- [ ] Commit / push / readback / clean status。

### 完成條件

- [ ] `knowledge/glossary/` canonical placement 已建立或明確 deferred。
- [ ] `memory/project/context-language.md` 非 canonical 邊界已寫入。
- [ ] Validation scenarios 可防止 root `CONTEXT.md` / memory canonical duplication。
- [ ] Runtime refresh / validate 通過。
- [ ] Linked updates 全部完成或明確 not applicable。

---

## Stakeholder 同意項目

- [ ] 接受 `knowledge/glossary/` 作 Ai-skill canonical glossary。
- [ ] 接受 `memory/project/context-language.md` 只作 project-local replay / context aid。
- [ ] 接受不建立 root `CONTEXT.md`。
- [ ] 接受先 test-first scenarios，再建立 glossary source。

---

## 與其他 plans 的關係

| Plan | 關係 |
| --- | --- |
| `plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md` | 需對齊 `context_mode` / `compression` / `memory_mode` 等詞彙，避免 active plan 內詞彙分叉。 |
| `plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md` | 若 glossary 未來成為 executable contract，需遵守 owner-layer YAML projection 規範。 |
| `plans/archived/2026-05-20-1635-bdd-ddd-cognition-aligned-reframe.md` | Requirements cognition 與 DDD cognition 的 shared language 邊界來源。 |
