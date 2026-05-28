# Context Language Glossary System

**Status**: `completed` (archived 2026-05-28)
**世代**：Gen 3 子系統擴充
**建立日期**：2026-05-25
**最後更新**：2026-05-28（Phase 0–7 全部完成；validator + projection + auto-detect signal + commit-msg gate 全部 live）

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
5. **詞條形狀未定義**：若沒有 Glossary Entry Schema，agent 可能把詞條寫成一句話、ADR、runtime spec、incident record 或 temporary implementation note。
6. **Semantic split-brain**：`context_mode`、`compression`、`memory_mode`、`generated_surface`、`projection` 可能被 runtime、workflow、knowledge、contract 文件各自重新定義。
7. **Resolution order 不明**：agent 面對 project glossary、knowledge glossary、ADR、workflow docs、intelligence heuristics、memory replay 時，可能讓 memory 或舊 ADR 覆蓋 active glossary。
8. **Document lookup 不足**：agent 需要回答「這個詞誰擁有、哪些 workflow / validation / route / plan 使用它、哪些 alias 或 deprecated term 指向它」，不是只找 Markdown 檔案。
9. **語義關聯未索引**：`context_mode`、`memory_mode`、`compression`、`projection`、`generated_surface`、`runtime_refresh`、`owner-layer` 這類短詞高度專案特化且 relationship-heavy，單靠全文搜尋或 embedding 不穩。

### Decision

建立 **Context Language Glossary System**：

| 類型 | Canonical location | 說明 |
| --- | --- | --- |
| Ai-skill 全庫 / 框架 / 可重用語彙 | `knowledge/glossary/` | 當前 canonical glossary，供 workflow、routing、architecture、requirements 引用。 |
| 單一專案、跨 session 但非 canonical 的語言脈絡 | `<PROJECT_ROOT>/memory/project/context-language.md` 或專案等價文件；Ai-skill 提供 `templates/project/context-language.template.md` | 只作 selective replay，不取代專案正式 docs 或 Ai-skill glossary；template 是建議格式，不是 hard contract。 |
| 判斷何時需要 shared language | `intelligence/engineering/requirements/behavior-modeling/` 與 `intelligence/engineering/architecture/domain-modeling/` | 保留為 reasoning source。 |
| 執行何時讀 glossary | `workflow/software-delivery/requirements/behavior-driven-discovery/` 與 `pre-build-interrogation` | Workflow gate 觸發 glossary 使用。 |
| Runtime semantic projection | Phase 5 依現有 runtime/index 架構決定位置（候選：`knowledge/runtime/sqlite/glossary.sqlite`、整合進既有 runtime index、或 runtime-owned projection table） | 從 Markdown canonical glossary 生成 normalized semantic index，用於 owner lookup、alias resolution、dependency traversal、drift detection 與 semantic routing；路徑可調，語義結構不可退化成 Markdown cache。 |
| Semantic governance | `governance/semantic/`（Phase 2 視 validation 決定是否建立） | 定義 naming、ownership、deprecation、aliasing、resolution order、drift detection 與 semantic conflict resolution。 |

`knowledge/glossary/` 不做百科全書；初期只收 **Ai-skill framework glossary**，也就是 framework terms、runtime semantics、cognitive vocabulary 與 architecture contracts。每個詞條必須有 schema 與 semantic owner；其他文件只能引用 owner definition，不能重新定義同一詞。`software-delivery.md`、`runtime.md`、`validation.md` 等 domain glossary 必須等 framework semantics 穩定後再拆，不在第一批建立 mega glossary。

### Glossary Entry Schema（計畫要求）

Phase 2 的 `knowledge/glossary/README.md` 必須先定義詞條 schema，再建立任何 glossary entries。

Required:

- `term`
- `status`
- `meaning`
- `affects`
- `owner-layer`

Optional:

- `aliases`
- `anti-meaning`
- `related-terms`
- `introduced-by`
- `deprecated-by`

Allowed `status` values:

- `canonical`
- `candidate`
- `deprecated`
- `superseded`
- `alias-only`
- `experimental`
- `project-local`

Forbidden:

- project-specific hosts、paths、class/test names、sample IDs 或 incident evidence
- temporary implementation detail
- current runtime state
- raw historical transcript / memory replay
- executable contract semantics that belong in workflow / enforcement / runtime source

### Semantic Ownership（計畫要求）

每個 glossary term 必須指定 semantic owner。Owner 不只沿用資料夾名，因為 `knowledge/`、`runtime/`、`workflow/` 是 storage / execution topology，不足以表達 semantic ownership。

Initial semantic owner domains:

- `runtime-cognition`
- `semantic-routing`
- `workflow-orchestration`
- `validation-governance`
- `memory-replay`
- `runtime-projection`
- `architecture-contracts`

例如：

```text
Term: context_mode
Owner-layer: runtime-cognition
```

Owner layer 擁有該詞的 canonical meaning；其他 layer 只能引用、alias 或標記 local usage，不得重新定義。若同一詞在不同 context 有不同 meaning，必須拆 bounded context 或改名，不能讓同一 term 承載兩個 canonical meanings。

### Vocabulary Resolution Priority（計畫要求）

當多個 source 對同一詞有不同說法時，agent 必須依下列順序解析：

1. Active project canonical docs（僅限該專案語境）
2. `knowledge/glossary/`
3. Accepted ADRs / `constitution/`
4. Workflow docs
5. Intelligence heuristics
6. Memory replay

限制：

- Project docs 只能覆蓋 project-local usage，不得改寫 Ai-skill framework term。
- Memory replay 永遠不能覆蓋 glossary / ADR / workflow current source。
- 舊 ADR 若與 active runtime docs 衝突，必須檢查 Framework Generation / Vocabulary Evolution section，不可直接採舊詞義。

### Semantic SQLite Projection（計畫要求）

Glossary 的 canonical source 仍是 Markdown，因為它需要 human readable、reviewable、git diffable 與 PR friendly。但 runtime 不應只做全文 cache；Phase 5 必須建立 **semantic-normalized SQLite projection**，把 glossary 轉成可查詢的語義索引層。

Projection 原則：

- `knowledge/glossary/*.md` 是 canonical source-of-truth。
- Projection location 由 Phase 5 architecture preflight 決定；`knowledge/runtime/sqlite/glossary.sqlite` 只是候選，不是硬性指定。
- 目前首選是整合進現有 `runtime-index.sqlite` generator，新增 glossary-specific tables，因為 glossary projection 本質仍是 knowledge projection，不是 runtime authority。
- 若未來 semantic routing 變成 runtime-heavy，再評估 dedicated `glossary.sqlite`；若成為 runtime-owned canonical config，必須走 `runtime/runtime.db` canonical boundary。
- 不論位置如何，projection 都不得保存第二份 rule body。
- SQLite 只保存 normalized semantic fields、relations、usage index 與 validation metadata。
- 不採用 `term/content` 單表全文 cache 作為主要設計；全文搜尋可以是輔助，不是核心。
- 不過早導入 vector DB / embeddings；先用 ownership、relations、lifecycle、routing、dependency 這些 strongly typed semantics。

Minimum tables（Phase 5 詳化）：

```sql
CREATE TABLE glossary_terms (
  term TEXT PRIMARY KEY,
  meaning TEXT NOT NULL,
  status TEXT NOT NULL,
  owner_layer TEXT NOT NULL,
  canonical_source TEXT NOT NULL,
  introduced_at TEXT,
  deprecated_at TEXT
);

CREATE TABLE glossary_relations (
  source_term TEXT NOT NULL,
  relation_type TEXT NOT NULL,
  target_term TEXT NOT NULL,
  source_file TEXT NOT NULL
);

CREATE TABLE glossary_usage (
  term TEXT NOT NULL,
  source_file TEXT NOT NULL,
  source_type TEXT NOT NULL,
  usage_context TEXT
);
```

Allowed `glossary_relations.relation_type` values start with:

- `alias_of`
- `related_to`
- `conflicts_with`
- `owned_by`
- `used_by`
- `deprecated_by`
- `replaced_by`

Allowed `glossary_usage.source_type` values start with:

- `workflow`
- `validation`
- `runtime`
- `knowledge`
- `adr`
- `plan`
- `memory`

Runtime queries this projection should support:

- 這個詞誰擁有？
- 哪些 workflow / validation scenario / runtime route / active plan 使用它？
- 哪些 alias 指向它？
- 哪些 term 已 deprecated 或被 replaced？
- 哪些詞 meaning 相近但 owner 不同，可能構成 semantic split-brain？
- 使用者問 `compression runtime issue` 時，應載入哪些 glossary owner、workflow、validation scenarios、active plans 與 runtime contracts？

Usage discovery policy:

- Phase 1 只認 declared references，例如 `related-terms`、`references`、`uses`、`depends-on`。
- Repo-wide scan 不得直接建立 canonical usage relation，因為 `mode`、`context`、`runtime`、`memory` 等短詞會造成大量 false semantic relationships。
- Phase 2 之後可加入 heuristic semantic discovery，例如 similarity、co-occurrence、alias candidates，但只能產生 candidate relations，不能直接 canonicalize。

### Evolution Roadmap（計畫要求）

| Generation | Capability | Completion Signal |
| --- | --- | --- |
| Gen 3.1 | Markdown canonical glossary + SQLite semantic projection | `knowledge/glossary/*.md` 可產生 framework-approved semantic projection，並可查 owner / status / canonical source。 |
| Gen 3.2 | Semantic relation graph | `glossary_relations` 支援 alias、related、conflict、deprecation、replacement 查詢。 |
| Gen 3.3 | Runtime semantic routing | 使用者訊號或 file diff 可透過 glossary 找 owner-layer、related workflows、validation scenarios 與 active plans。 |
| Gen 3.4 | Glossary drift auto-detection | Validation 可自動偵測 owner missing、alias loop、deprecated term resurrection、conflicting meanings。 |
| Gen 4 | Context-aware cognitive loading | Agent 先查 touched semantic domains，再決定載入哪些 knowledge、workflow、runtime contracts 與 validation rules。 |

### Open Question Decisions（2026-05-25）

| # | Decision | Rationale | Plan impact |
| --- | --- | --- | --- |
| 1 | 第一批 glossary 只放 Ai-skill framework vocabulary。 | 目前最危險的 drift 是 `context_mode`、`compression`、`projection`、`runtime_refresh`、`generated_surface`、`owner-layer`、`memory_replay`，不是業務詞典。 | Phase 3 只建立 `ai-skill.md`；`software-delivery.md`、`runtime.md`、`validation.md` deferred。 |
| 2 | 提供 project-local context language template，但標成 non-canonical。 | Downstream project 需要穩定格式，但 project 差異大，不應硬性契約化。 | Phase 4 新增 `templates/project/context-language.template.md`，用途是 replay / context aid only。 |
| 3 | Glossary 先不要 YAML contract。 | Glossary 是 semantic knowledge，不是 executable contract；過早 YAML 會增加 schema rigidity、projection sync 與 migration burden。 | Gen 3 維持 Markdown canonical + SQLite semantic projection；只有 runtime gating / executable enforcement / blocking validation 真的需要時才 promotion。 |
| 4 | `runtime-cognitive-modes-system` 必須依賴 glossary canonical definitions。 | `context_mode`、`compression`、`memory_mode`、`reasoning_mode` 最容易形成 subsystem-local semantics。 | 本 plan 與 cognitive modes plan 都新增 semantic dependency boundary。 |
| 5 | `status` 加入 `alias-only`、`superseded`、`experimental`。 | `alias-only` 防止 alias 擁有自己的 meaning；`superseded` 保留歷史引用；`experimental` 支援 Gen 3 快速演化。 | Phase 2 schema 固定最小 status set。 |
| 6 | `owner-layer` 改採 semantic owner domain，不只沿用資料夾。 | Folder layer 是 storage / execution topology，不足以表達 ownership。 | Phase 2 定義 `runtime-cognition`、`semantic-routing`、`workflow-orchestration`、`validation-governance`、`memory-replay` 等初始 owner domains。 |
| 7 | Phase 2 不建立獨立 `governance/semantic/`。 | 目前還沒有足夠 semantic governance pressure；過早拆層會過度架構化。 | Semantic governance 先放 `knowledge/glossary/README.md`，等 drift incidents、ownership conflicts、alias explosion、semantic migration 出現再拆。 |
| 8 | Semantic projection 先整合現有 runtime index generator。 | Glossary projection 本質是 knowledge projection，不是 runtime authority。 | Phase 5 優先在現有 `runtime-index.sqlite` generator 增加 glossary tables；未來 runtime-heavy 時再拆 dedicated DB。 |
| 9 | `glossary_usage` Phase 1 只掃 declared references。 | Repo-wide scan 對短詞太 noisy，會產生 false semantic relationships。 | Phase 1 只認 `related-terms`、`references`、`uses`、`depends-on`；heuristic discovery 只能產 candidate relations。 |

### Alternatives Considered

- **A. 建 root `CONTEXT.md`**：拒絕。Root context 容易變成 always-loaded 大檔，且會與 `knowledge/`、`workflow/`、`memory/` 形成平行 source。
- **B. 放 `memory/project/context-language.md` 作 canonical**：拒絕。Memory 是 selective replay，不是 current truth；可作 project-local replay，但不可作 Ai-skill canonical。
- **C. 只保留 intelligence，不新增 glossary**：拒絕。Intelligence 只回答何時需要 ubiquitous language，不提供 agent 可查的穩定詞彙 source。
- **D. 建 `knowledge/glossary/` 作 canonical，memory 只作 project replay**：接受。
- **E. 只建 glossary，不建 semantic governance**：接受但收斂範圍。Phase 2 不建立 `governance/semantic/`，先在 `knowledge/glossary/README.md` 內承載 semantic ownership / lifecycle / conflict rules；等 drift incidents、ownership conflicts、deprecation lifecycle、alias explosion 或 semantic migration 出現，再拆出 `governance/semantic/`。

### Why Not an ADR Yet

依 `plans/README.md` 與 ADR-007：

- 目前是可逆的 layer / workflow 擴充，尚未證明需要 constitution 級不可逆決策。
- 是否需要 runtime projection、validation scenarios、project-local template 的具體實作形狀，仍需透過本 plan 驗證。
- 若 completed 後證明 glossary 成為 foundational + cross-session + cross-project + expensive-to-reverse 的 architecture boundary，再評估 ADR promotion。

### ADR Promotion Criteria（completed 時驗證）

- [ ] `knowledge/glossary/` 已被 workflow 或 routing 實際使用。
- [ ] `knowledge/glossary/README.md` 已定義 Entry Schema、Semantic Ownership 與 Vocabulary Resolution Priority。
- [ ] Semantic drift validation 能偵測 duplicate meaning、conflicting ownership、alias loops、deprecated term resurrection。
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
- SQLite semantic projection 讓 agent 查語義關聯，而不是只查檔案路徑。
- 未來可支援 runtime semantic routing：由 term 找 owner-layer、workflow、validation、active plan 與 runtime contract。

#### 負面

- 新增 `knowledge/glossary/` 需要 routing、graph、summary 與 README 連動。
- 新增 semantic projection 需要 compiler / refresh 邏輯或 dedicated generator，避免 projection drift；實際位置需符合 runtime/index ownership boundary。
- 若 glossary 過度膨脹，會變成低價值百科，需要嚴格 scope。
- Project-local glossary template 需要說清楚「專案 source 優先」，避免污染 reusable docs。

#### 風險

| 風險 | 緩解 |
| --- | --- |
| Glossary 變成百科 | 只收影響 behavior / contract / validation / runtime decision 的詞；其他詞留在 project docs。 |
| Memory 被當 canonical | 在 memory README、glossary README 與 validation scenario 明確禁止。 |
| 與 intelligence/domain-modeling 重複 | Intelligence 保留判斷原則；glossary 保存穩定詞條。 |
| 與 active runtime cognitive modes plan 詞彙衝突 | Phase 0 先盤點 `context_mode` / `compression` / `memory_mode`，避免先寫兩套。 |
| Semantic split-brain | 每個 term 有 `owner-layer`；其他 layer 只能引用，不得重新定義。 |
| Alias / deprecated term 復活 | Phase 1 先寫 `semantic-term-overlap-v1.yaml`，Phase 2 定義 alias / deprecation lifecycle。 |
| SQLite 變成 Markdown cache | Phase 5 明確禁止 `term/content` 單表作核心設計；必須 normalized terms + relations + usage index。 |
| 過早 vectorization | Gen 3.x 先做 structured semantics；只有在 relation / usage index 無法滿足查詢時才評估 embeddings。 |
| Software-delivery mega glossary 爆炸 | 第一批只收 Ai-skill framework terms；domain glossary 必須等 framework semantics 穩定後拆分。 |
| Repo-wide usage scan noise | Phase 1 只收 declared references；heuristic scan 只能產生 candidate relations。 |

---

## Runtime Execution Path

| 欄位 | 內容 |
| --- | --- |
| Runtime owner | Phase 1-4：Markdown canonical source + manual workflow reads；Phase 5：existing runtime index generator + framework-approved semantic projection + runtime reports；Phase 5+ 若需要 blocking gate，再新增 owner-layer executable contract。 |
| Trigger flow | Phase 1-4：user asks for framework / workflow / term / source-of-truth clarification → `pre-build-interrogation` or software-delivery route loads candidate sources → glossary source is manually read as dependency evidence → plan final report lists glossary decision / no duplication evidence.<br>Phase 5 semantic projection：file_diff or user_signal contains glossary / term conflict / ubiquitous language / source-of-truth wording → `route.knowledge.glossary` or semantic index query detects glossary candidate → query `glossary_terms` / `glossary_relations` / `glossary_usage` → load canonical `knowledge/glossary/*.md` owner entry + related workflow / validation / plan sources → run glossary drift / semantic overlap scenarios → block completion if owner, alias, deprecated term, or memory/project boundary conflicts. |
| Trigger location | Phase 1-4：`pre-build-interrogation`、`route.workflow.software-delivery`、agent dependency-read ledger；Phase 5：`file_diff`、`user_signal`、runtime semantic index lookup、future `route.knowledge.glossary`。 |
| Activation contract | Phase 1-5：無新 executable contract；使用既有 `workflow.software_delivery.pre_build_interrogation.contract` 作 plan 前 gate。Phase 6：新增 (a) `cognitive-modes-discovery.yaml` 中 `file_diff_glossary_touched` + `user_keyword_term_conflict` 兩個 discovery signals 與 `route.knowledge.glossary` 關聯；(b) `commit-msg` validator `validateGlossaryRetroOwn`，當 staged diff 動到 `runtime/cognitive-modes*.yaml` / `runtime/economics/**` / `ecosystem/**` 但沒同時 stage `knowledge/glossary/ai-skill.md` 時 block，opt-out marker `[skip-glossary-retro-own]`。 |
| Generated surface | 初期不投影 executable contract；Phase 5 建立 framework-approved semantic projection（候選為 dedicated SQLite、既有 runtime index extension 或 runtime-owned projection table），並在 runtime reports 記錄 glossary source / route / validation coverage。若新增 executable contract，target_key 需另定。 |
| Validation scenarios | Phase 1 先新增：`validation/scenarios/failure-derived/glossary-source-duplication-v1.yaml`、`validation/scenarios/failure-derived/memory-context-language-as-canonical-v1.yaml`、`validation/scenarios/failure-derived/semantic-term-overlap-v1.yaml`。 |
| Test passing evidence | `ai-skill runtime refresh --native-index --native-reports`、`ai-skill runtime validate --json`、scenario / lints / diff review。 |

### Doc-only Trial 聲明（Phase 1-5）+ Runtime Graduation（Phase 6）

**Phase 1-5 為 doc-only trial**，初期不接入新的 runtime execute layer。理由：

- Shared language 先是 knowledge source，不是 blocking runtime gate。
- 現有 `pre-build-interrogation` 已能在 plan 前阻擋 source-of-truth duplication。
- 若未先驗證 glossary scope，直接做 runtime blocking contract 會增加 noise。

**Graduation 條件已達成（2026-05-28）**：

- Schema spec、validator、19 entries、SQLite projection 與 5 條 query 已全綠
- `route.knowledge.glossary` 已存在於 `routing-registry.yaml`
- 沒有 entry schema 之外的機械強制 → 仍構成 [`governance/lifecycle/system-upgrade-governance.yaml`](../../governance/lifecycle/system-upgrade-governance.yaml) §`define_runtime_trigger_flow` 的 `routing_only_without_detector_or_action` 與 `sqlite_projection_without_routable_or_validator_consumer` 違規

→ 因此本 plan **不能在 Phase 5 後直接 close-loop**，必須先完成 **Phase 6 Runtime Auto-Detect Integration**（見下方）才能 graduate；否則違反 doc-only trial 的 `graduation_deadline_or_signal` 與 `explicit_acknowledgement_doc_only_trial_does_not_count_as_runtime_integration` 規則。

未來接入時機：

- Phase 5 先把 semantic projection 整合進現有 runtime index generator；若 projection + validation 仍不足以穩定阻擋 drift，再新增 glossary owner-layer executable contract。

---

## Phase 0 Pre-Build Interrogation

| 欄位 | 內容 |
| --- | --- |
| Trigger | 使用者要求把 `CONTEXT.md / ubiquitous language system` 寫成計畫，並提醒架構已調整需重讀規則。 |
| Checked sources | `CORE_BOOTSTRAP.md`、`enforcement/dependency-reading.md`、`enforcement/linked-updates.md`、`workflow/software-delivery/requirements/pre-build-interrogation.md`、`plans/README.md`、`governance/lifecycle/system-upgrade-governance.md`、`architecture/README.md`、ADR-007、knowledge / memory README、ubiquitous language intelligence。 |
| Goal | 建立 shared language canonical placement plan，吸收 `CONTEXT.md` 技巧但不引入第二份 source-of-truth。 |
| Scope | Plan only；規劃 framework-only `knowledge/glossary/`、non-canonical project template、project memory boundary、semantic SQLite projection、workflow/routing/validation linked updates。 |
| Non-goals | 本 plan 不立即建立完整 glossary、不中斷 active runtime cognitive modes plan、不建立 proposed ADR、不新增 root `CONTEXT.md`。 |
| Acceptance | Plan 符合新 `plans/README.md` 必填章節；明確回答 `knowledge/glossary/` vs `memory/project/context-language.md`；列出 test-first validation。 |
| Framework discovery | `knowledge/` 是 navigation / atom / glossary source 候選；`knowledge/runtime/sqlite/` 可作 generated semantic projection；`memory/` 不保存 current truth；`intelligence/` 保存判斷原則；`workflow/` 保存執行 gate；`governance/semantic/` 是 semantic lifecycle / conflict resolution 候選。 |
| Duplication risk | Root `CONTEXT.md`、memory canonical、knowledge glossary 三者不可並存為同一語彙 source。Plan 採 `knowledge/glossary/` canonical + memory replay 非 canonical；每個 term 另需 `owner-layer` 防止 semantic split-brain。 |
| Open questions | 見下一節。 |
| Decision | proceed with draft plan；implementation 需另行啟動。 |

---

## Open Questions

**全部於 2026-05-25 resolved**。決議詳見 §Open Question Decisions；新增問題只能在 Phase 0 / Phase 1 test-first validation 中提出，且不得重新打開 root `CONTEXT.md`、software-delivery mega glossary、repo-wide semantic crawling 或 vector-first architecture，除非先新增 validation evidence。

---

## Phase 1 Test-First Validation

先建立 validation scenarios，再實作 glossary 結構。

### 期望可觀察行為

- Agent 不會建 root `CONTEXT.md` 作 Ai-skill canonical glossary。
- Agent 不會把 `memory/project/context-language.md` 當 current truth。
- Agent 在 framework / requirements plan 中能找到 `knowledge/glossary/` 作 shared language source。
- Agent 能偵測 near-duplicate terms、conflicting owner-layer、alias loops 與 deprecated term resurrection。
- Agent 能偵測 semantic projection 設計是否退化為 `term/content` Markdown cache。

### Tasks

- [ ] 新增 `validation/scenarios/failure-derived/glossary-source-duplication-v1.yaml`
- [ ] 新增 `validation/scenarios/failure-derived/memory-context-language-as-canonical-v1.yaml`
- [ ] 新增 `validation/scenarios/failure-derived/semantic-term-overlap-v1.yaml`
- [ ] 新增 `validation/scenarios/failure-derived/glossary-semantic-projection-shape-v1.yaml`
- [ ] 新增 `validation/scenarios/failure-derived/cognitive-core-vs-ecosystem-boundary-v1.yaml`（cross-plan pre-architecture：同時服務 `plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md` Phase 2 — 防 `runtime/cognitive-modes*.yaml` 被塞入 adaptation/economics 邏輯，且防 `ecosystem/` 或 `runtime/economics/` 被誤用成第二個 runtime control plane。在 glossary plan 屬 semantic ownership boundary 驗證；在 economics plan 屬 Phase 2 前置 scenario。）
- [ ] 更新 relevant graph / summary / routing candidate（若 scenario 需要）

### Phase 1 完成條件

- [ ] Scenarios 符合 `validation/scenario.schema.json`
- [ ] `ai-skill runtime refresh --native-index --native-reports` 通過
- [ ] `ai-skill runtime validate --json` 通過

---

## Phase 2 Glossary Schema And Semantic Governance Boundary

### Tasks

- [ ] 建立 `knowledge/glossary/README.md`，先定義 Entry Shape，不先建立大量詞條。
- [ ] 定義 required / optional / forbidden fields。
- [ ] 定義 `owner-layer` semantics：owner 擁有 canonical meaning，其他 layer 只能引用或 alias。
- [ ] 定義初始 semantic owner domains：`runtime-cognition`、`semantic-routing`、`workflow-orchestration`、`validation-governance`、`memory-replay`、`runtime-projection`、`architecture-contracts`。
- [ ] **Cross-plan pre-architecture**：同步定義 candidate semantic owner domains，status=`candidate`：`ecosystem-adaptation`（承載 economics / pressure / adaptation 詞彙）、`runtime-economics`（若 economics plan Phase 0 確認 `runtime/economics/` 路徑則 activate）。此舉讓 `plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md` Phase 1 不必回頭改 glossary owner domain；若 economics plan 最終選用 `ecosystem/` top-level layer，`runtime-economics` 可降為 deprecated。
- [ ] 定義 Vocabulary Resolution Priority。
- [ ] 定義最小 status set：`canonical`、`candidate`、`deprecated`、`superseded`、`alias-only`、`experimental`、`project-local`。
- [ ] 定義 relation lifecycle：`alias_of`、`related_to`、`conflicts_with`、`owned_by`、`used_by`、`deprecated_by`、`replaced_by`、`derived_from`、`aggregates`。`derived_from` / `aggregates` 為 economics plan 預留（cross-plan pre-architecture，2026-05-27 confirmed）— `cognitive_cost` aggregates `thinking_cost` + `context_cost` + `execution_cost` + `knowledge_cost`，這層關係必須 first-class 才不會在 Phase 3 entries 落入 free-form `related_to`。
- [ ] 定義 usage index source types：`workflow`、`validation`、`runtime`、`knowledge`、`adr`、`plan`、`memory`。
- [ ] 定義 drift detection categories：duplicate meaning、conflicting ownership、alias loop、deprecated term resurrection、near-duplicate concept fork。
- [ ] 將 semantic governance 暫放在 `knowledge/glossary/README.md`；只記錄未來拆出 `governance/semantic/` 的 promotion triggers。

### Phase 2 Schema Commitments（2026-05-27 confirmed）

下列決策已於 plan execution 前確認，必須寫進 `knowledge/glossary/README.md`，並由 Phase 2 validator 機械強制（見下方 Validator commitment）：

| # | 決策 | 內容 |
| --- | --- | --- |
| 1 | **Entry physical shape** | Markdown H2 heading per term + 緊跟一個 YAML code block 承載 schema fields。H2 提供 GitHub anchor link 與人類閱讀層；YAML block 供 Phase 5 SQLite projection parser 機械解析。禁止：純 H2 + 散落欄位列表、單一大 frontmatter / 全文件 YAML array、一 term 一檔。 |
| 2 | **Schema validation 機制** | Phase 2 必須在 `scripts/ai-skill-cli/` 新增 `glossary validate` subcommand（Go validator），讀 `knowledge/glossary/*.md` 內的 H2 + YAML block，檢查 required / optional / forbidden fields、status enum、owner-layer enum、alias 規則、relation enum、symmetric relation 對稱性。並接進 `ai-skill runtime validate` 流程，使 Phase 2 之後新增 entry 必須通過。純人類 review 不被接受。 |
| 3 | **Term naming convention** | `term:` 欄位 canonical 寫法為 **snake_case**（例：`context_mode`、`cognitive_cost`）。kebab-case（例：`context-mode`）由 validator 自動列為 alias，無需手動寫入 `aliases:`。Markdown H2 heading 使用 snake_case 與 `term:` 完全一致，確保 GitHub anchor stable。 |
| 4 | **Alias 表達方式** | Alias 不建獨立 entry。每個 alias 只在 canonical term 的 `aliases:` 欄位列出（flat string array）。Validator 強制：(a) `aliases:` 中的字串不得出現為其他 term 的 `term:` 欄位；(b) alias chain 不得形成 cycle；(c) `status: alias-only` 在新 schema 下視為禁用（已由 `aliases:` 取代）。狀態 set 仍保留 `alias-only` 但僅供 deprecated entry 過渡使用，新 entry 禁用。 |
| 6 | **`introduced-by` / `deprecated-by` 引用格式** | 只允許兩種值：`plans/<path>`（active 或 archived plan 路徑）或 `constitution/ADR-XXX.md`（accepted ADR）。禁止 commit SHA、issue 編號、PR URL（會 rot 或脫離 repo）。Validator 強制此 enum 形狀。 |
| 8 | **Relation enum 預留** | 初始 allowed `relation_type` 包含 `derived_from` 與 `aggregates`（為 economics plan §Split cost model 預留），即使 Phase 3 第一批 entries 暫未使用。理由：避免 Phase 3 寫 `cognitive_cost` entry 時退回 free-form `related_to`、損失 split cost 語義。 |
| 5 | **Relation 對稱性：單向寫入 + projection 推導** | Relation 在 Markdown / YAML 只寫單向。Phase 5 SQLite projection 自動 derive 反向關係（query 時換方向，不寫第二筆 row）。`knowledge/glossary/README.md` 必須包含 **symmetry 分類表**（哪些 relation symmetric、哪些 asymmetric + 反向語義名稱），與 **人類查找反向關係範例**（同時提供 (a) `ai-skill glossary inspect <term>` CLI 預期輸出，顯示 incoming relations；(b) 直接 SQL 範例如 `SELECT source_term FROM glossary_relations WHERE target_term = 'compression' AND relation_type = 'alias_of'`；(c) Markdown 內肉眼查找路徑：去 target term 的 entry 看 `related-terms` 與 alias，或 grep `target_term: <name>` 找上游引用）。Validator 強制 symmetric relation 雙向出現、asymmetric relation 僅單向；symmetry 分類表為 YAML 內嵌（machine-readable），不能只是散文。 |
| 7 | **`anti-meaning` 拆為兩欄（B 案）** | Schema 拆成兩個 optional 欄位：(a) `anti-meaning` 純為 disambiguator — 「這個詞不是另一個聽起來像它的詞」（純人類閱讀，validator 不強制引用）；(b) `excludes` 為 scope negative — 「這個詞不涵蓋語義鄰居的職責」，必須引用其他 glossary term（validator 強制 `excludes:` 內字串為合法 term name），可被 `cognitive-core-vs-ecosystem-boundary-v1` 等 owner boundary scenario 機械引用。Phase 3 entries 寫入時，若 owner 邊界有風險（例：`context_mode` vs `discovery_mode`），`excludes:` 必填；純命名混淆才用 `anti-meaning`。 |

### Validator Commitment（Phase 2 新增）

Phase 2 必須建立 `scripts/ai-skill-cli/internal/glossary/validator.go`（或等效路徑），由 `ai-skill glossary validate` 與 `ai-skill runtime validate` 兩處入口執行。最小檢查集：

- H2 heading 與 YAML block 配對；H2 文字必須等於 YAML block 的 `term:`
- Required fields 完整、forbidden patterns 未出現
- `status` / `owner-layer` / `relation_type` / `aliases` source type 等 enum 合法
- Term naming convention：`term:` 為 snake_case
- `introduced-by` / `deprecated-by` 形狀為 `plans/<path>` 或 `constitution/ADR-XXX.md`
- Alias 規則：alias 不得為另一 canonical term；alias chain 無 cycle
- Symmetric relation 對稱規則：依 README.md 內 symmetry 分類表，symmetric relation 必須雙向出現於兩端 entry；asymmetric relation 僅單向寫入，反向由 projection 推導
- `excludes:` 欄位內字串必須為現有合法 term name（防止打字錯誤或孤立引用）

### 反向關係查找（人類路徑，必寫進 README.md）

`knowledge/glossary/README.md` 必須包含下列三條人類查找範例，作為 Phase 5 projection 上線前後皆可用的 fallback：

1. **CLI**：`ai-skill glossary inspect <term>` 預期輸出格式 — 顯示該 term 的 outgoing + incoming relations（incoming 由 SQLite query 即時 derive）。Phase 5 實作完成前先給 stub usage example。
2. **SQL**：直接查 SQLite，例如 `SELECT source_term FROM glossary_relations WHERE target_term = 'compression' AND relation_type = 'alias_of';` 顯示誰指向 `compression` 作為 alias。
3. **Markdown 肉眼**：去 target term 的 entry 看 `related-terms` / `aliases` 欄位；或 `grep -r "target_term: compression" knowledge/glossary/` 找上游引用。在 projection 尚未生成時，這是唯一可用路徑。

Linked update：若新增 validator，必須同步更新 `scripts/ai-skill-cli/internal/app/hooks.go` registry、`scripts/ai-skill-cli/README.md`、`runtime/runtime.db` 對應 per-commit validator 條目（若 promotion 為 mandatory）。

### Phase 2 完成條件

- [ ] 任何 glossary entry 建立前，schema 已存在。
- [ ] Entry physical shape（H2 + YAML block）已寫明，含 1 個範例。
- [ ] Schema validation 機制：`glossary validate` subcommand 已存在並通過 sample test。
- [ ] Term naming convention（snake_case canonical、kebab-case auto alias）已寫明，validator 強制。
- [ ] Alias 表達方式（flat array、無 alias entry）已寫明，validator 強制。
- [ ] `introduced-by` / `deprecated-by` 形狀已寫明，validator 強制。
- [ ] Relation enum 含 `derived_from` / `aggregates`。
- [ ] Symmetry 分類表已內嵌於 README.md（YAML 形式，machine-readable），validator 引用。
- [ ] 反向關係人類查找三條路徑（CLI / SQL / Markdown 肉眼）已寫進 README.md。
- [ ] `anti-meaning` / `excludes` 拆兩欄，`excludes` validator 強制引用合法 term。
- [ ] Semantic ownership 與 resolution priority 已寫明。
- [ ] Status set 與 semantic owner domains 已寫明。
- [ ] Relation types 與 usage source types 已寫明。
- [ ] `governance/semantic/` 已明確 deferred，且列出 drift incidents、ownership conflicts、deprecation lifecycle、alias explosion、semantic migration 作為 promotion trigger。

---

## Phase 3 Knowledge Glossary Structure

### Proposed files

```text
knowledge/glossary/
  README.md
  ai-skill.md
```

Deferred candidates after framework semantics stabilize:

```text
knowledge/glossary/
  software-delivery.md
  runtime.md
  validation.md
```

### Tasks

- [ ] 建立 `knowledge/glossary/ai-skill.md`，收 Ai-skill framework 詞彙。
- [ ] **Retro-own archived cognitive vocabulary**：`ai-skill.md` 必須涵蓋 `plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md` 引入的 cognitive vocabulary（至少 `context_mode`、`compression`、`memory_mode`、`reasoning_mode`），每詞補上 `owner-layer`（候選：`runtime-cognition`）；否則這些詞會在 archived plan 中成為 silent owner，違反 §Semantic Ownership。來源依據：Phase 0 Pre-Build Interrogation §7 risk 重新檢視。
- [ ] **Cross-plan pre-architecture — Economics-adjacent candidate terms**：為避免 `plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md` 執行時詞彙真空，於 `ai-skill.md` 預收下列詞條，全部 status=`candidate`（等 economics plan Phase 1 確認 owner path 再 promote 到 `canonical`）：
  - `cognitive_cost`（owner-layer: `runtime-cognition`）— 現行 6-dim report 已有，補 canonical 定義
  - `thinking_cost`、`context_cost`、`execution_cost`、`knowledge_cost`（owner-layer: `ecosystem-adaptation`）— split cost model 子項
  - `knowledge_mode`、`discovery_mode`、`intelligence_mode`（owner-layer: `ecosystem-adaptation`）— economics plan Phase 9 引入
  - `ecosystem`（owner-layer: `ecosystem-adaptation`）— cross-layer interaction layer 概念
  - `pressure_model`（owner-layer: `ecosystem-adaptation`）— 例：`context_explosion`、`memory_amplification`
  Anti-scope：不為 economics plan 預建 `runtime/economics/` 或 `ecosystem/` 結構；不把上述詞條設為 `canonical`（economics plan 尚未驗證 owner path）。
- [ ] 明確禁止第一批建立 `software-delivery.md` / `runtime.md` / `validation.md`，除非 Phase 3 完成後另有 validation evidence。
- [ ] 更新 `knowledge/README.md`，將 glossary 加入目前入口。
- [ ] 更新 `knowledge/graphs/`，把 glossary 連到 requirements / domain modeling / software delivery。

### Phase 3 完成條件

- [ ] Glossary 只收影響 behavior / contract / validation / runtime decision 的詞。
- [ ] 第一批 glossary 只收 Ai-skill framework terms、runtime semantics、cognitive vocabulary、architecture contracts。
- [ ] Archived `runtime-cognitive-modes-system.md` 引入的 cognitive vocabulary 已 retro-own，無 silent owner。
- [ ] 每個 term 都符合 Entry Schema，且有 `owner-layer`。
- [ ] 沒有把 project-specific 詞、host、class、incident evidence 寫入 reusable glossary。
- [ ] 沒有建立 software-delivery mega glossary。
- [ ] `knowledge/README.md` 可導向 glossary。

---

## Phase 4 Workflow And Memory Boundary Integration

### Tasks

- [ ] 更新 `workflow/software-delivery/requirements/behavior-driven-discovery/README.md`：shared language 對齊時先查 `knowledge/glossary/`。
- [ ] 更新 `workflow/software-delivery/requirements/pre-build-interrogation.md`：framework discovery 若有詞彙分叉，先查 glossary。
- [ ] 更新 `memory/README.md` 或 `memory/project/README.md`：project context-language 只作 replay，不作 canonical。
- [ ] 新增 `templates/project/context-language.template.md`，標題必須是 `Non-Canonical Project Context Language`，並列明不是 canonical source、不是 runtime truth、不是 architecture contract。
- [ ] 更新 `intelligence/engineering/requirements/behavior-modeling/ubiquitous-language-alignment.md`：判斷原則指向 glossary 作 source。
- [ ] 更新 `intelligence/engineering/architecture/domain-modeling/ubiquitous-language.md`：domain glossary 與 Ai-skill glossary 的邊界。
- [ ] 若建立 `governance/semantic/`，更新 lifecycle / routing / linked updates 的 semantic governance boundary。

### Phase 4 完成條件

- [ ] Workflow、memory、intelligence 三層沒有互相取代。
- [ ] `memory/project/context-language.md` 的用途被限制為 project-local replay / context aid。
- [ ] Project-local context language template 是 recommended format，不是 hard contract。
- [ ] Agent 能從 behavior-driven discovery 找到 canonical glossary source。
- [ ] Resolution priority 已在 workflow / memory / intelligence 引用。

---

## Phase 5 Semantic Projection, Routing, And Runtime Reports

### Tasks

- [ ] 更新 `knowledge/runtime/routing-registry.yaml`，新增或擴充 glossary route。
- [ ] 更新 `knowledge/summaries/requirements-cognition.md` / `development-guidance.md` / related summaries。
- [ ] 優先在現有 `runtime-index.sqlite` generator 內建立 glossary projection tables；若不採用，必須記錄 architecture preflight 理由。
- [ ] 建立或規劃 framework-approved semantic projection schema，至少支援 `glossary_terms`、`glossary_relations`、`glossary_usage`。
- [ ] 確認 projection 從 `knowledge/glossary/*.md` 生成，不成為第二份 canonical source。
- [ ] `glossary_usage` Phase 1 只收 declared references；repo-wide scan 只能生成 candidate relations。
- [ ] 執行 `ai-skill runtime refresh --native-index --native-reports`，讓 glossary source、semantic projection 與 route 進 generated lookup。
- [ ] 補查詢驗證：owner lookup、alias resolution、usage traversal、deprecated / replaced term lookup、owner missing drift query。
- [ ] 若決定需要 executable contract，再補 `knowledge/glossary/*.yaml` 或 owner-layer contract；否則明確記錄 Markdown-only / knowledge source。

### Phase 5 完成條件

- [ ] Runtime reports / SQLite index 可找到 glossary source。
- [ ] Semantic projection 可查 term owner、status、canonical source、relations 與 usage。
- [ ] Projection schema 不是 `term/content` Markdown cache。
- [ ] Usage index 沒有把 repo-wide textual matches 直接 canonicalize。
- [ ] 若無 runtime executable contract，已記錄原因：glossary 是 knowledge source，不是 workflow gate。
- [ ] 若有 contract，`generated_surfaces` target_key 已 synced。
- [ ] Semantic drift scenarios 進入 validation inventory 或 runtime reports 可查來源。

---

## Phase 6 Runtime Auto-Detect Integration（graduation gate）

> 為什麼新增此 phase：Phase 5 完成 SQLite projection 與 `route.knowledge.glossary` 後，[`system-upgrade-governance.yaml`](../../governance/lifecycle/system-upgrade-governance.yaml) §`define_runtime_trigger_flow` 的 `routing_only_without_detector_or_action` 與 `sqlite_projection_without_routable_or_validator_consumer` 仍未被滿足 — route 存在但無 detector / projection 存在但無 consumer。本 phase 把 trigger chain（event → detector → loaded source → runtime action → evidence）的缺失環節補齊。

### Trigger Chain（為何此 phase 是必要的）

| Chain 元素 | Phase 5 完成度 | Phase 6 補上 |
| --- | --- | --- |
| event_or_signal | ❌ 沒 signal | ✅ `file_diff_glossary_touched`、`user_keyword_term_conflict` |
| detector | ⚠️ route 有，但沒 signal 拉它 | ✅ discovery signal wire route + `validateGlossaryRetroOwn` commit-msg validator |
| loaded_source_or_executable_contract | ✅ `route.knowledge.glossary` 已在 routing-registry | （已完成）|
| runtime_action_or_blocking_gate | ❌ 無 | ✅ Validator `validateGlossaryRetroOwn` block commit；signal escalate `context_mode` 至 `SOURCE_BACKED` |
| observable_evidence | ⚠️ 只能用 SQL 手動驗 | ✅ 新 validation scenario `glossary-retro-own-missing-v1.yaml` |

### Tasks

- [ ] 新增 discovery signal：在 [`runtime/cognitive-modes-discovery.yaml`](../../runtime/cognitive-modes-discovery.yaml) 加 `file_diff_glossary_touched`（trigger: staged diff includes `knowledge/glossary/**` 或 `runtime/cognitive-modes*.yaml`）與 `user_keyword_term_conflict`（trigger: user prompt 含 "glossary" / "term conflict" / "ubiquitous language" / "shared language" / "owner-layer"）；兩個 signal 都 escalate `context_mode → SOURCE_BACKED` 並 load `route.knowledge.glossary`。
- [ ] 新增 commit-msg validator `validateGlossaryRetroOwn`：staged diff 動到 `runtime/cognitive-modes*.yaml` / `runtime/economics/**` / 未來 `ecosystem/**` 但沒同時 stage `knowledge/glossary/ai-skill.md` → block with exit 30；opt-out trailer `[skip-glossary-retro-own]`，需在 commit message 說明為何此次 framework 改動不需 retro-own term。Validator dispatch 走 `scripts/ai-skill-cli/internal/app/hooks.go` registry。
- [ ] 新增 validation scenario `validation/scenarios/failure-derived/glossary-retro-own-missing-v1.yaml`：framework term 引入 runtime YAML 但無對應 glossary entry → 預期被 validator block。
- [ ] 升級 [`runtime/cli-modification-policy.yaml`](../../runtime/cli-modification-policy.yaml)：把 `gate.glossary.retro_own_required` 列入 commit-msg gates，opt-out marker `[skip-glossary-retro-own]`。
- [ ] 更新 [`scripts/ai-skill-cli/docs/command-contract.md`](../../scripts/ai-skill-cli/docs/command-contract.md) 加新 validator description；更新 [`bdd-scenarios.md`](../../scripts/ai-skill-cli/docs/bdd-scenarios.md) 加 retro-own 場景；更新 [`test-fixture-plan.md`](../../scripts/ai-skill-cli/docs/test-fixture-plan.md) 加 fixture。
- [ ] Go tests for `validateGlossaryRetroOwn` (happy path / block path / opt-out path)；bin rebuild 兩段 commit。

### Phase 6 完成條件

- [ ] Discovery signals 已 wire 並可被 `cognitive-modes-discovery.yaml` parser 讀到；route auto-load 可由 file_diff 觸發
- [ ] `validateGlossaryRetroOwn` commit-msg hook 上線；本地 sample test 證明 block + opt-out
- [ ] `glossary-retro-own-missing-v1.yaml` scenario 通過 `runtime validate`
- [ ] §Runtime Execution Path Activation contract 描述與實際 wired 行為一致
- [ ] [`governance/lifecycle/system-upgrade-governance.yaml`](../../governance/lifecycle/system-upgrade-governance.yaml) §`define_runtime_trigger_flow` 的 5-element chain 全部滿足，無 forbidden 違規

---

## Phase 7 Close Loop

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
- [ ] Glossary Entry Schema、Semantic Ownership、Vocabulary Resolution Priority 已建立或明確 deferred with blocker。
- [ ] 第一批 glossary 限定 Ai-skill framework terms，且 domain glossary 已 deferred。
- [ ] Non-canonical project context language template 已建立或明確 deferred with blocker。
- [ ] Semantic SQLite projection 已建立或明確 deferred with blocker。
- [ ] Semantic drift validation scenario 已建立。
- [ ] `memory/project/context-language.md` 非 canonical 邊界已寫入。
- [ ] Validation scenarios 可防止 root `CONTEXT.md` / memory canonical duplication。
- [ ] Runtime refresh / validate 通過。
- [ ] Linked updates 全部完成或明確 not applicable。

---

## Stakeholder 同意項目

- [ ] 接受 `knowledge/glossary/` 作 Ai-skill canonical glossary。
- [ ] 接受第一批只放 Ai-skill framework glossary，不建立 software-delivery mega glossary。
- [ ] 接受 `memory/project/context-language.md` 只作 project-local replay / context aid。
- [ ] 接受提供 non-canonical project context language template。
- [ ] 接受不建立 root `CONTEXT.md`。
- [ ] 接受 glossary 暫不建立 YAML contract，先採 Markdown canonical + SQLite semantic projection。
- [ ] 接受 semantic owner domains 不只沿用資料夾 layer。
- [ ] 接受 Phase 2 暫不建立 `governance/semantic/`。
- [ ] 接受 semantic projection 優先整合現有 runtime index generator。
- [ ] 接受 Phase 1 usage index 只掃 declared references。
- [ ] 接受先 test-first scenarios，再建立 glossary source。
- [ ] 接受先做 structured semantic projection，不過早導入 vector DB / embeddings。

---

## 與其他 plans 的關係

| Plan | 關係 |
| --- | --- |
| `plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md` | 必須引用 glossary canonical definitions for runtime semantic vocabulary；`context_mode`、`compression`、`memory_mode`、`reasoning_mode` 不得在 cognitive modes plan 內形成 subsystem-local semantics。 |
| `plans/archived/2026-05-22-0855-executable-yaml-contract-migration.md` | 若 glossary 未來成為 executable contract，需遵守 owner-layer YAML projection 規範。 |
| `plans/archived/2026-05-20-1635-bdd-ddd-cognition-aligned-reframe.md` | Requirements cognition 與 DDD cognition 的 shared language 邊界來源。 |
| `plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md` | **Downstream consumer**。本 plan 預先在 Phase 1 / 2 / 3 為 economics plan 架設 candidate scenario（`cognitive-core-vs-ecosystem-boundary-v1`）、candidate owner domains（`ecosystem-adaptation`、`runtime-economics`）與 candidate terms（split cost model、`knowledge_mode` / `discovery_mode` / `intelligence_mode`、`ecosystem`、`pressure_model`）。Economics plan 執行時可直接引用，不需要回頭擴 glossary schema。 |

---

## Session Handoff（下次接手必讀）

**目前進度斷點**：Phase 0 已執行完成，**尚未進入 Phase 1**。沒有檔案系統實際變更（除本 plan 自身）。

### 已完成
- ✅ Phase 0 Pre-Build Interrogation 已驗證：所有 checked sources 存在；`knowledge/glossary/`、`governance/semantic/`、`templates/`、`memory/project/context-language.md`、root `CONTEXT.md` **全部尚未存在**（greenfield，符合預期）。
- ✅ Open Questions 9 題全部 resolved（2026-05-25）。
- ✅ 與其他 plans cross-reference 已核實：`runtime-cognitive-modes-system.md`、`executable-yaml-contract-migration.md`、`bdd-ddd-cognition-aligned-reframe.md` 全在 `plans/archived/`，無 active 衝突。
- ✅ Phase 3 補強：retro-own archived cognitive vocabulary（`context_mode` / `compression` / `memory_mode` / `reasoning_mode`）— 防 silent owner。

### 下一個動作（從這裡開始）
**Phase 1 Test-First Validation** — 建立 4 個 validation scenarios（**先**寫 scenarios，**再**建 glossary 結構）：
1. `validation/scenarios/failure-derived/glossary-source-duplication-v1.yaml`
2. `validation/scenarios/failure-derived/memory-context-language-as-canonical-v1.yaml`
3. `validation/scenarios/failure-derived/semantic-term-overlap-v1.yaml`
4. `validation/scenarios/failure-derived/glossary-semantic-projection-shape-v1.yaml`

Schema 在 `validation/scenario.schema.json`；參考既有 `validation/scenarios/failure-derived/*.yaml` 取得格式範例。完成後跑：
```
ai-skill runtime refresh --native-index --native-reports
ai-skill runtime validate --json
```

### 不可重開的封閉決議（碰到請拒絕）
- ❌ 不建 root `CONTEXT.md`
- ❌ 第一批不建 software-delivery / runtime / validation mega glossary
- ❌ 不做 repo-wide semantic crawling 作 canonical usage
- ❌ 不導入 vector DB / embeddings（Gen 3 結束前）
- ❌ Memory 永遠不能 override glossary / ADR / workflow current source

### Runtime State 提醒
本 plan 觸發的 runtime phase 流程：`phase.bootstrap → phase.checkpoint → phase.execution`（Phase 1 起算）。每個 session 啟動仍須 Bootstrap Receipt（見 [CORE_BOOTSTRAP.md](../../CORE_BOOTSTRAP.md)）。

### 關鍵檔案速查
- 本 plan：`plans/active/2026-05-25-1000-context-language-glossary-system.md`
- Schema 參照：`validation/scenario.schema.json`
- Pre-build gate（已通過）：`workflow/software-delivery/requirements/pre-build-interrogation.md`
- 待 retro-own 的 archived plan：`plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md`
