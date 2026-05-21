# Memory Retrieval & Activation Governance

> **狀態**: draft
> **建立日期**: 2026-05-20
> **目的**: 將 `memory/` 從「歷史資料分類」升級為 selective cognitive replay system，補足 retrieval、activation、replay cost、freshness、contamination boundary、working-memory buffer 與 promotion pipeline，同時避免與既有 cognitive state / evidence governance 重複。

---

## 1. Problem Statement

目前 `memory/` 已有長期記憶分類：

- `working/`：session-local 工作記憶。
- `summary/`：壓縮 session 摘要。
- `decision/`：跨 session 決策記憶。
- `episodic/`：情境記憶。
- `project/`：專案脈絡記憶。
- `failure/`：抽象化失效記憶。

但目前仍偏向 storage taxonomy，而不是 runtime 可選擇 activation 的 cognitive replay system。主要缺口：

1. Retrieval trigger 不夠明確，agent 不知道何時該回放哪類 memory。
2. Replay qualification 不完整，memory 容易被當成 source-of-truth。
3. Replay token cost、freshness decay 與 replay depth 缺少治理。
4. `memory/working/`、`.agent-goals/` 與 runtime execution state 的邊界仍需收緊。
5. Memory promotion 目前有方向，但缺少從 working memory → summary → episodic/failure/project/decision → knowledge/intelligence/enforcement 的可治理 pipeline。
6. 與 `knowledge/` 的差異需要寫入 durable architecture boundary，避免 future agents 把 historical replay 與 structured navigation 混用。

本 plan 的目標不是讓 memory 自動載入更多 context，而是讓 agent 知道何時「不要回想」、何時只 weak replay、何時必須重新驗證。

---

## 2. Architecture Compatibility Preflight

| Field | Content |
| --- | --- |
| Trigger | 使用者提供 Memory Layer Upgrade Plan，要求先建立計畫、檢查是否與 active plan 重複，並具體拆分。 |
| Checked sources | `plans/active/2026-05-20-1501-cognitive-state-evidence-governance.md`、`plans/README.md`、`memory/README.md`、`memory/working/README.md`、`memory/summary/README.md`、`memory/episodic/README.md`、`memory/project/README.md`、`memory/failure/README.md`、`memory/decision/README.md`、`knowledge/README.md`、`runtime/README.md`、`governance/lifecycle/README.md`、`enforcement/conversation-goal-ledger.md`。 |
| Conflicts | 既有 active plan 已涵蓋 memory 作為 evidence source、stale execution memory、cognitive contamination、temporal confidence decay、belief GC、governance minimality 與 runtime reduction。本 plan 不重寫那些 cognitive-state governance；只定義 `memory/` 內部 retrieval / activation / replay / promotion lifecycle，以及其與 `knowledge/`、`working memory`、`.agent-goals/`、`runtime/` 的分層邊界。 |
| Decision | Proceed as separate active plan。Scope 應保持 memory-layer governance，不直接新增 runtime state machine，也不取代 cognitive-state plan 的 evidence/confidence governance。 |
| Validation | Plan readback、diff review、link check、ReadLints。若後續 phase 修改 `knowledge/runtime/routing-registry.yaml`、runtime compiler source 或 generated surfaces，必須執行 `ai-skill runtime refresh` 與相關 validator。 |

### 2.1 與現有 Active Plan 的切分

| Concern | Cognitive State & Evidence Governance | 本 Memory Plan |
| --- | --- | --- |
| Memory as evidence | 定義 memory evidence 的信心、freshness、scope 與 contamination 風險。 | 定義哪些 memory 可以被 retrieval、何時只能 weak replay、何時需要外部 revalidation。 |
| Cognitive contamination | 定義 stale frame 污染 execution 的治理與 autonomy downgrade。 | 定義 replay 前的 contamination boundary 與 forbidden replay cases。 |
| Confidence decay | 定義 belief / evidence / claim 的 confidence integrity。 | 定義 memory confidence defaults、freshness decay 與 replay qualification。 |
| Belief GC | 定義 execution-scoped belief 的 prune / deprecated lifecycle。 | 定義 memory pruning、compression、promotion 與 archived/cold lookup lifecycle。 |
| Runtime guard | 只 promotion 最小 viable cognitive safety primitives。 | 預設不新增 runtime state；只把 memory 作為 selective retrieval candidate。 |
| `.agent-goals/` | 檢查 execution intent stability，不取代 goal ledger。 | 明確禁止 active execution contract 長期化成 memory，除非經 abstraction / compression / promotion。 |

### 2.2 Recommended Execution Order

本 plan 依賴既有 active plan 的上游治理語意，不建議先完整執行 memory replay runtime 化。

建議順序：

1. 先執行 `plans/active/2026-05-20-1501-cognitive-state-evidence-governance.md`。
   - 原因：先定義 evidence qualification、confidence integrity、claim scope、intent stability、cognitive contamination、runtime reduction 與 minimal runtime principle。
   - Memory replay 的 freshness、confidence、contamination boundary 應引用這些通用治理語意，而不是在 `memory/` 重新發明一套。
2. 再執行本 plan 的 Batch A / Phase 0-2。
   - 只處理 `memory/`、`memory/working/`、`.agent-goals/`、`knowledge/`、`runtime/` 的分層邊界。
   - 暫不新增 runtime guard，不把 memory activation 變成 persistent execution state。
3. 最後執行本 plan 的 Batch B / C / Phase 3-5。
   - 在上游 cognitive-state governance 定穩後，再補 replay cost、activation threshold、freshness decay、routing、graphs、validation scenarios 與 lifecycle/enforcement linked updates。

Blocking note：

- 若 cognitive-state plan 尚未完成，memory plan 仍可執行 Batch A 的文件邊界整理。
- 但 memory replay 的 runtime guard、confidence enforcement、contamination enforcement 或 activation rule promotion 必須等 cognitive-state plan 的 runtime reduction / signal normalization 完成後再決定。

---

## 3. Target Layer Model

目標分層：

```text
memory/
  historical replay archive
  ↓ retrieval + qualification
memory/working/
  session cognition buffer
  ↓ selected activation
.agent-goals/
  active execution contract
  ↓ execution
runtime/
  executable / machine-readable state and lookup
```

### 3.1 Layer Responsibility

| Layer | Responsibility | Must Not Become |
| --- | --- | --- |
| `memory/` | 長期 historical cognition archive；保存可回放、可壓縮、可治理的歷史脈絡。 | 永久 context dump、active execution state、canonical source。 |
| `memory/working/` | Session cognition buffer；保存本 session 可丟棄的 activation frame、recent evidence、temporary workflow context。 | `.agent-goals/` 替代品、長期 memory、runtime-state。 |
| `.agent-goals/` | Active execution contract；保存 current goal、next action、success criteria、blockers、owner、execution status。 | Long-term archive、memory summary、decision record。 |
| `runtime/` | 可執行、可查詢、deterministic 的 runtime state / lookup surface。 | 歷史經驗、長篇 reasoning、memory archive。 |
| `knowledge/` | Knowledge navigation、summaries、graphs、routing registry、runtime-facing lookup。 | Historical replay、session summary、project recap。 |
| `intelligence/` | 抽象化 reasoning 與可重用判斷智慧。 | Raw memory、incident log、execution state。 |

### 3.2 Knowledge vs Memory Boundary

| Question | `knowledge/` | `memory/` |
| --- | --- | --- |
| 回答什麼 | 一般情況下有哪些可重用結構、摘要、graph、routing path？ | 以前發生過什麼、是否值得參考？ |
| 本質 | Structured reusable navigation / reference surface。 | Historical replayable experience。 |
| 預設使用 | 可被 routing 作為低 token navigation surface。 | 預設 dormant，需 trigger + qualification 才 replay。 |
| 信心 | 經 abstraction / normalization 後較穩定。 | 帶有 historical bias，需 freshness / scope / confidence qualification。 |
| Promotion | 接收被抽象化且可重用的 memory-derived structure。 | 保存尚保留歷史痕跡的經驗。 |

判斷規則：

> 如果某內容即使沒有特定 incident 也仍成立，放 `knowledge/`、`intelligence/`、`workflow/` 或 `governance/`；如果它依賴「曾經發生過什麼」，放 `memory/`，且不得直接當作 source-of-truth。

---

## 4. Proposed Components

### 4.1 Memory Retrieval Governance

新增 candidate directory：

```text
memory/retrieval-governance/
├── README.md
├── activation-thresholds.md
├── retrieval-routing.md
├── replay-cost-governance.md
├── replay-budget.md
├── freshness-and-decay.md
├── contamination-boundary.md
└── memory-promotion-policy.md
```

Responsibilities：

- 定義 retrieval trigger。
- 定義 activation threshold。
- 定義 replay qualification。
- 定義 replay cost class。
- 定義 replay budget 與最大 replay depth。
- 定義 memory confidence / freshness / scope defaults。
- 定義禁止 replay 的情境。
- 定義 memory promotion 與 pruning policy。

### 4.2 Retrieval Routing

Memory 不應 always-loaded，而應作為 candidate retrieval source。

初版 retrieval signals：

| Signal | Candidate Memory | Qualification |
| --- | --- | --- |
| repeated failure class | `memory/failure/`、`memory/episodic/` | 必須符合相同 failure class 或相似 execution graph。 |
| same repo / same project | `memory/project/` | 必須確認 repo / architecture boundary 仍相容。 |
| architecture decision recall | `memory/decision/` | 檢查 status 是否 accepted / superseded / deprecated。 |
| context compaction recovery | `memory/summary/`、`memory/working/` archive | 只回放最小摘要，不回放 full transcript。 |
| stale assumption suspicion | `memory/failure/`、`memory/episodic/` | 只作 weak hint，必須重新驗證 current source。 |
| workflow family match | `memory/episodic/`、`memory/project/` | 需確認 workflow scope 與 domain boundary。 |

### 4.3 Activation Pipeline

Memory replay 必須通過：

```text
trigger
→ retrieval
→ qualification
→ replay budget check
→ activation
→ memory/working/ buffer
→ execution usage
→ revalidation / discard / promotion
```

Rules：

- 未通過 qualification 的 memory 不得進入 active execution frame。
- Episodic memory 預設只能作 weak guidance。
- Project memory 只在同 repo / 同 architecture boundary / 同 workflow family 下使用。
- Decision memory 必須檢查 status 與 supersession。
- Summary memory 只能用於恢復脈絡，不得取代 current source reading。
- Full session replay 預設禁止，除非使用者明確要求或沒有其他足夠 source。

### 4.4 Replay Cost Governance

| Memory Type | Cost | Default |
| --- | --- | --- |
| failure pattern memory | low | frequent but scoped。 |
| decision memory | low | allowed after status check。 |
| summary memory | medium | conditional，用於 handoff / context recovery。 |
| episodic memory | medium | weak guidance only。 |
| project memory | medium-high | scoped to same project / repo。 |
| old execution recap | high | on-demand only。 |
| full transcript / full session replay | very high | avoid。 |

Rules：

- Prefer smallest sufficient replay scope。
- Replay depth 不得超過當前 task 需要的最小 evidence。
- 若 replay 成本高於重新讀 canonical source，優先讀 canonical source。
- Replay 不得形成 recap recursion：summary → old summary → old transcript → older summary。

### 4.5 Freshness, Scope, and Confidence Defaults

| Memory Type | Default Confidence | Freshness Handling |
| --- | --- | --- |
| transcript-derived memory | low | 必須重新驗證。 |
| episodic memory | tentative | 只能作 weak guidance。 |
| summary memory | scoped | 只恢復 session context，不證明 current truth。 |
| failure abstraction | medium | 可提示 risk，但 current source 仍需檢查。 |
| project memory | scoped | 受 repo architecture / migration / dependency changes 影響。 |
| decision memory | medium-high | status / supersession 檢查後可用。 |

每個可長期回放的 memory 應逐步補上：

```yaml
last_validated:
expires_when:
compatibility_scope:
confidence_default:
replay_allowed_as:
```

### 4.6 Memory Contamination Boundary

Memory replay 常見 contamination：

- stale architecture reuse。
- old workflow replay。
- outdated routing assumption。
- invalidated failure workaround。
- prior repo topology reuse。
- stale blocker replay。
- pseudo-active-state。

Boundary classes：

| Boundary | Meaning | Replay Rule |
| --- | --- | --- |
| `workflow-local` | 僅適用同一 workflow family。 | 可作 checklist hint，需 current source check。 |
| `domain-local` | 僅適用同一 domain / architecture family。 | 不可跨 domain 自動 replay。 |
| `project-local` | 僅適用同一 repo / project。 | repo refactor / migration 後需 revalidation。 |
| `session-global` | 會影響整個 session frame。 | 需 recap / prune / human alignment。 |

Forbidden replay：

- Replay stale blockers as active blockers。
- Replay old `.agent-goals/` state as current execution contract。
- Replay deprecated architecture frame without compatibility check。
- Replay old execution graph without revalidation。
- Replay memory-derived conclusion as canonical source。

### 4.7 Working Memory as Session Cognition Buffer

保守做法：先不新增 top-level `working-memory/`，而是強化既有 `memory/working/`。

`memory/working/` 應定義為：

- session-scoped。
- semi-stable。
- discardable。
- replay-selected。
- temporary cognition frame。

可保存：

```yaml
active_assumptions:
recent_evidence:
current_architecture_frame:
active_repo_topology:
temporary_workflow_context:
current_risk_assessment:
activated_memory_refs:
discard_after:
```

不保存：

- current blocker as durable truth。
- next action / owner / lock。
- long-term project state。
- canonical decision。
- runtime execution state。

### 4.8 Memory Compression and Promotion Pipeline

正確 promotion pipeline：

```text
memory/working/
→ compress
memory/summary/
→ select episode / project / decision / failure
memory/episodic/ | memory/project/ | memory/decision/ | memory/failure/
→ abstract / generalize
intelligence/ | knowledge/ | enforcement/ | workflow/
```

Forbidden promotion：

- Raw transcript。
- Temporary blocker。
- Active runtime assumption。
- Unstable execution graph。
- Unresolved contradiction。
- Project-secret / private evidence。

Promotion criteria：

| Condition | Required |
| --- | --- |
| reusable | yes |
| generalized | yes |
| non-project-secret | yes |
| low contamination risk | yes |
| repeated utility | preferred |
| source compatibility known | yes |

### 4.9 Routing and Runtime Boundary

Memory 可以接入 routing，但不可變成 persistent execution state。

Candidate linked updates：

- `knowledge/runtime/routing-registry.yaml`：新增 memory retrieval governance route 或擴充 `memory.operations` candidate。
- `knowledge/summaries/memory-operations.md`：更新 summary，加入 retrieval / activation / replay economics。
- `knowledge/graphs/`：建立 memory operations graph 或補強既有 graph。
- `runtime/router/activation-rules.yaml`：只有在 memory governance 成為 lazy-load rule 時才更新。
- `runtime/README.md`：必要時補一句 runtime 不保存 historical replay。

Runtime boundary：

- Memory replay 只輸出 candidate context。
- Runtime 可查 memory route / lookup metadata。
- Runtime 不保存 raw historical memory。
- 若要新增 memory activation guard，必須先完成 signal compression，避免與 cognitive-state plan 重複。

---

## 5. Suggested Implementation Phases

### Phase 0 — Boundary Confirmation

Status: draft.

Tasks:

- [ ] 確認本 plan 不重寫 `2026-05-20-1501-cognitive-state-evidence-governance.md` 的 evidence / confidence / contamination governance。
- [ ] 確認 `memory/working/` 先作為 session cognition buffer，不立即升級成 top-level `working-memory/`。
- [ ] 確認 `memory/` 與 `knowledge/` 的 durable boundary。
- [ ] 確認 `.agent-goals/` 不得 archive 成 memory，除非 compression / abstraction / promotion。
- [ ] 確認 memory retrieval 只產生 candidate context，不產生 runtime execution state。

Exit criteria:

- [ ] 分層邊界已寫入 plan。
- [ ] 與 active cognitive-state plan 的重疊已標出。
- [ ] 後續 phase 的 linked updates 已列出。

### Phase 1 — Memory Retrieval Governance Documents

Candidate files:

- `memory/retrieval-governance/README.md`
- `memory/retrieval-governance/activation-thresholds.md`
- `memory/retrieval-governance/retrieval-routing.md`
- `memory/retrieval-governance/replay-cost-governance.md`
- `memory/retrieval-governance/replay-budget.md`
- `memory/retrieval-governance/freshness-and-decay.md`
- `memory/retrieval-governance/contamination-boundary.md`
- `memory/retrieval-governance/memory-promotion-policy.md`

Tasks:

- [ ] 建立 retrieval governance index。
- [ ] 定義 activation threshold 與 trigger examples。
- [ ] 定義 replay cost classes 與 replay depth policy。
- [ ] 定義 memory freshness、scope、confidence defaults。
- [ ] 定義 contamination boundary 與 forbidden replay。
- [ ] 定義 promotion / compression / pruning policy。

### Phase 2 — Existing Memory README Updates

Candidate files:

- `memory/README.md`
- `memory/working/README.md`
- `memory/summary/README.md`
- `memory/episodic/README.md`
- `memory/project/README.md`
- `memory/failure/README.md`
- `memory/decision/README.md`

Tasks:

- [ ] 更新 `memory/README.md`：memory 是 selective replay，不是 always-loaded context。
- [ ] 更新 `memory/working/README.md`：改為 session cognition buffer，明確不保存 owner / next action / active blockers。
- [ ] 更新 `memory/summary/README.md`：summary 只作 compressed context recovery，不作 current truth。
- [ ] 更新 `memory/episodic/README.md`：episodic replay 預設 weak guidance。
- [ ] 更新 `memory/project/README.md`：加入 compatibility scope / freshness decay。
- [ ] 更新 `memory/failure/README.md`：failure memory 可提示 risk，但不能取代 enforcement rule。
- [ ] 更新 `memory/decision/README.md`：decision replay 必須檢查 status / supersession。

### Phase 3 — Knowledge / Routing Integration

Candidate files:

- `knowledge/summaries/memory-operations.md`
- `knowledge/runtime/routing-registry.yaml`
- `knowledge/graphs/`
- `knowledge/runtime/runtime-report.md`（generated if applicable）

Tasks:

- [ ] 更新 memory operations summary，加入 retrieval / activation / replay economics。
- [ ] 決定是否新增 memory retrieval route；若新增，列出 candidate sources。
- [ ] 若新增或更新 route，執行 knowledge runtime refresh。
- [ ] 視需要建立 memory graph，連結 memory type、retrieval governance、promotion target、runtime boundary。
- [ ] 確認 generated surfaces 與 source 一致。

### Phase 4 — Governance / Enforcement Boundary Updates

Candidate files:

- `governance/lifecycle/README.md`
- `governance/ai-runtime-governance/`（若需要 memory lifecycle governance）
- `enforcement/conversation-goal-ledger.md`
- `enforcement/failure-learning-system.md`

Tasks:

- [ ] 在 lifecycle governance 補上 memory promotion / cold lookup / pruning 的 durable boundary。
- [ ] 如需可執行規則，補充 `.agent-goals/` 不得長期化成 memory 的 forbidden behavior。
- [ ] 確認 failure learning promotion 與 memory promotion policy 不重複。
- [ ] 若 cognitive-state plan 已處理 contamination / confidence，這裡只連結，不重寫。

### Phase 5 — Validation Scenarios

Candidate files:

- `validation/scenarios/memory/`
- `ai-skill runtime validate`（如需新增 semantic validation）

Tasks:

- [ ] 測試 stale blocker replay 被阻止。
- [ ] 測試 old `.agent-goals/` state 不會被當成 current execution state。
- [ ] 測試 episodic memory 只能 weak guidance。
- [ ] 測試 project memory 跨 repo / architecture boundary 時需 revalidation。
- [ ] 測試 decision memory superseded 時不得被直接採用。
- [ ] 測試 full session replay 需要 explicit on-demand trigger。
- [ ] 測試 replay cost governance 阻止不必要的 context inflation。
- [ ] 測試 promotion pipeline 阻止 raw transcript / temporary blocker 進入 long-term memory。

### Phase 6 — Plan Completion Closure

Tasks:

- [ ] 確認所有 phase 完成或標 blocked。
- [ ] 執行適用 validator：ReadLints、Markdown link check、`ai-skill runtime refresh`（若修改 routing / knowledge runtime）。
- [ ] 檢查 linked updates。
- [ ] 更新 `plans/README.md` 狀態。
- [ ] 若完成，搬移至 `plans/archived/`。
- [ ] Commit / push / readback / clean status。

---

## 6. Concrete Work Breakdown

建議拆成三個可 review 的工作批次：

### Batch A — Memory Layer Boundary

Scope:

- 新增 `memory/retrieval-governance/README.md`。
- 更新 `memory/README.md` 與 `memory/working/README.md`。
- 寫清楚 `memory`、`memory/working`、`.agent-goals`、`runtime`、`knowledge` 的邊界。

Reason:

- 先解決最容易污染 execution state 的邊界問題。

### Batch B — Replay Governance and Existing Memory Types

Scope:

- 新增 activation threshold、replay cost、freshness、contamination、promotion policy 文件。
- 更新 `summary/`、`episodic/`、`project/`、`failure/`、`decision/` README。

Reason:

- 讓每種 memory type 都知道何時可 replay、何時只能 weak hint、何時必須 revalidate。

### Batch C — Routing, Validation, and Lifecycle

Scope:

- 更新 `knowledge/summaries/memory-operations.md`。
- 視需要更新 `knowledge/runtime/routing-registry.yaml` 與 graph。
- 補 validation scenarios。
- 補 lifecycle / enforcement linked updates。

Reason:

- 最後才接 routing / validation，避免文件邊界未定時把 memory 過早 runtime 化。

---

## 7. Open Questions

1. `memory/working/` 是否足夠，還是要升級為 top-level `working-memory/`？
   - Current recommendation: 先保留 `memory/working/`，改語義為 session cognition buffer。
2. Memory retrieval governance 要放在 `memory/retrieval-governance/`，還是 `governance/memory/`？
   - Current recommendation: retrieval mechanics 放 `memory/`，lifecycle / promotion gate 可由 `governance/` 引用。
3. Memory activation 是否需要 runtime guard？
   - Current recommendation: 初期不要。先由 routing / governance 文件控制，等 validation scenarios 證明 recurring failure 後再 promotion。
4. `knowledge/` 是否應索引 memory？
   - Current recommendation: 可以索引 retrieval governance 與 memory operations summary，但不要把 memory 全文 always-load。
5. Memory 是否應保存 `.agent-goals/` 完成後的 state？
   - Current recommendation: 不保存。只有經 compression / abstraction 後，才可能進 `summary/`、`project/`、`decision/`、`episodic/` 或 `failure/`。

---

## 8. Completion Definition

本 plan 完成時，系統應能做到：

- `memory/` 明確是 historical replay archive，不是 active state。
- `knowledge/` 明確是 structured navigation / reusable reference，不是 historical replay。
- `memory/working/` 明確是 session cognition buffer，不保存 owner / lock / next action。
- `.agent-goals/` 保持 active execution contract，不被長期 archive 成 memory。
- Memory replay 需要 trigger、retrieval、qualification、budget check、activation。
- Episodic memory 預設只能 weak guidance。
- Project memory 有 repo / architecture / workflow compatibility scope。
- Decision memory replay 會檢查 status / supersession。
- Summary memory 不取代 current source reading。
- Replay cost governance 能避免 full session replay 與 context inflation。
- Freshness / scope / confidence defaults 已寫入各 memory type。
- Contamination boundary 會阻止 stale blocker、old execution graph、deprecated architecture frame 直接污染 execution。
- Promotion pipeline 能把成熟 memory 抽象化到 `knowledge/`、`intelligence/`、`workflow/` 或 `enforcement/`，而不是堆積 raw context。
- Routing / summary / graph / validation surfaces 只在需要時接入，不讓 memory 成為 always-loaded context。
