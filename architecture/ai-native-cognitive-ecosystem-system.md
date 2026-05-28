# AI-native Cognitive Ecosystem System

**世代**：第四代（**vision / next-gen target，尚未 current**）
**前一代**：[`ai-native-cognitive-execution-system.md`](ai-native-cognitive-execution-system.md)（current, Gen 3）

本文件是**下一代願景文件**，不是 current canonical 入口。Gen 3 仍是 current；本檔記錄「要成為什麼樣的系統」與「在哪些 criteria 達標後才能 graduate 為 current」。

> 升級治理見 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md)。第四代升級的 graduation 條件由本檔 §Threshold Criteria 定義，並需符合該 governance 文件 §1 世代演化規則（命名變更 + 重新評估每個 ADR 的延伸狀態）。

---

## 系統世代演化

| 世代 | 名稱 | 文件 | 狀態 |
|------|------|------|------|
| 1 | Skill Repository | （未集中文件化）| historical |
| 2 | AI-native Knowledge Operating System | [`ai-native-knowledge-operating-system.md`](ai-native-knowledge-operating-system.md) | historical |
| 3 | AI-native Cognitive Execution System | [`ai-native-cognitive-execution-system.md`](ai-native-cognitive-execution-system.md) | **current** |
| 4 | **AI-native Cognitive Ecosystem System** | 本檔 | **vision** |

---

## 為什麼是「Ecosystem」而不是「Execution」

Gen 3 的核心問題是「**怎麼讓 AI 正確執行**」：

- workflow / routing / governance / runtime / validation / memory / models / tools / cognitive mode
- 它是一個 **Execution Runtime**（cognitive OS），解決 correctness、stability、token cost、failure recovery

Gen 4 的核心問題是「**整個 AI 生態如何協作與演化**」：

- Signal → Cognitive Economics → Knowledge Activation → Multi-Agent Ecology → Adaptive Runtime → Memory Evolution → Intelligence Evolution
- 它是一個 **Cognitive Ecosystem**，解決 evolution、adaptation、telemetry、fitness、decay、specialization

兩代並非取代關係 — Gen 3 的 execution runtime 仍是 Gen 4 的底層，但**重心**從「正確性」轉移到「演化效率」。Gen 3 的 runtime 在 Gen 4 中會退化成 ecosystem 中的 execution engine 子層，類比：

| 階段 | 類比 |
|---|---|
| Gen 3 | Kubernetes Runtime |
| Gen 4 | 整個 Cloud Ecosystem（含 marketplace、autoscaler、scheduler、telemetry、policy controller） |

---

## 核心 Pivot：從「Workflow」到「Cognitive Resource Management」

Gen 3 的 mental model 中心是 **workflow**：給定任務，找出該執行的 workflow / route / contract，跑完並驗證。

Gen 4 的 mental model 中心是 **cognitive resource management**：「有限 cognition 如何最有效使用」。Workflow 變成 resource 之一，與 memory、knowledge、intelligence、tool、model 並列被 economics layer 統籌調度。

延伸結果：

- `runtime/` 退化成 execution engine 子層
- 新層 `ecosystem/`（或等效）承載 cognition / economics / intelligence-evolution / discovery / orchestration / telemetry / adaptation / governance-of-evolution

Candidate Gen 4 layer skeleton（仍為 vision，待 plan 落地）：

```text
ecosystem/
├── cognition/       # cognitive resource accounting, mode adaptation
├── economics/       # token / latency / recursion / compression / model-fit costs
├── intelligence/    # intelligence fitness, decay, promotion
├── discovery/       # adaptive discovery, archaeology, gap detection
├── orchestration/   # multi-agent ecology, role specialization
├── telemetry/       # signal capture, observability, runtime emit
├── adaptation/      # pressure-based mode selection, rationale capture
├── evolution/       # workflow / heuristic / knowledge lifecycle, decay
└── governance/      # governance-of-evolution (meta-governance)

runtime/
└── execution-engine/  # the former Gen 3 runtime, now a subsystem
```

---

## Threshold Criteria（graduation 必達標準）

達成下列 ≥ 7/10 criteria 且其他 3 個有 active plan + 明確 entry condition，才考慮 graduate 本檔為 `current`。每條 criteria 必須是 **machine-verifiable**（generated surface / SQLite query / scenario evidence），不接受純散文宣稱。

### A. Cognitive Economics — 活的成本模型

- 至少 6 個 cost dimensions（token / latency / recursion / compression / model-fit / memory-loading / workflow-validation）有 executable contract，且 projected to `runtime.db generated_surfaces`
- Cognitive Mode discovery 消費這些 cost signals 做 mode selection
- 有 split cost model（thinking / context / execution / knowledge）取代單一 `cognitive_cost` summary
- **Acceptance signal**：`SELECT * FROM cognitive_cost_split` 在 runtime SQLite index 可查；commit-msg validator 比對 declared cost vs evidence

### B. Knowledge Evolution & Decay — 知識有 lifecycle

- 每個 knowledge atom 有 usage telemetry（last used、hit count、failure-after-load count）
- Stale knowledge 有 decay scoring；達閾值自動 deprecate proposal
- Promotion pipeline 不是手動 review；有 fitness-based promotion candidate emitter
- **Acceptance signal**：`SELECT * FROM atom_fitness ORDER BY decay_score DESC` 可查；validator 阻擋引用 decayed atom

### C. Discovery Evolution — 主動發現

- Discovery 不再只是 static routing；有 heuristic discovery、archaeology mode、knowledge-gap detection 三種模式 active
- 新 `discovery_mode` enum 已 wired 到 cognitive-modes-discovery
- Gap detection 可在 plan execution 時主動建議「缺少 X 類型的 knowledge atom」
- **Acceptance signal**：`discovery_mode` 為 6 維 cognitive vector 的第七維；有 ≥ 1 個 scenario 證明 gap detection 主動發 promotion candidate

### D. Ecosystem Signals — 跨層訊號完整

- Cognitive-modes-discovery signals ≥ 25 條（目前 14 條），涵蓋 tool-derived、model-derived、memory-derived、workflow-derived、knowledge-derived 五大類
- Signal 不只觸發 mode，也餵入 adaptation layer
- **Acceptance signal**：`runtime/cognitive-modes-discovery.yaml` 內 `signal_type:` enum 含 `tool_usage`、`model_capability`、`memory_freshness`、`workflow_pressure`、`knowledge_activation` 五類

### E. Adaptive Runtime — 壓力驅動的 mode 選擇

- Mode selection 不只看 preset signal threshold，看實際 observed pressure（context expansion rate、recursion depth、retry cost）
- 有 feedback loop contract（即使 v1 是 static）允許「上次 DEEP 浪費了 → 下次同樣 signal 改 NORMAL」
- **Acceptance signal**：`ecosystem.adaptation.contract` 在 `runtime.db generated_surfaces` 存在；scenario 證明同 signal 在不同 pressure 下產生不同 mode

### F. Memory Economics & Decay — 記憶的價值會衰減

- Memory entry 有 staleness 計分（時間、引用次數、互斥 evidence）
- Replay budget 不只是 token 預算，是 cost-of-replay vs expected-value
- Project memory 與 episodic memory 有不同 decay 曲線
- **Acceptance signal**：`SELECT memory_id, decay_score FROM memory_lifecycle WHERE decay_score > threshold` 可查；validator 阻擋復活已 decayed memory

### G. Intelligence Evolution — 判斷準則的適應

- Intelligence atom 有 fitness score（activation 後 outcome 是否成功）
- 多次 failure-after-activation 觸發 deprecation proposal
- Cross-intelligence synthesis：當兩個 intelligence 常被一起 activate，emit candidate merged atom
- **Acceptance signal**：`intelligence_mode` enum 已 wired；fitness telemetry 在 runtime reports 可見

### H. Multi-Agent Ecology — 不再假設單一 agent

- 同一 repo 可由多個 agent role 並行使用（reviewer / implementer / explorer / debugger），每 role 有 cognitive profile
- Multi-agent lock / coordination contract 存在
- `.agent-goals/` 升級為多 agent owner 對齊面
- **Acceptance signal**：governance YAML 內有 `multi_agent_coordination.contract`；至少一個 scenario 驗證 role conflict resolution

### I. Telemetry Layer — 系統能觀察自己

- Runtime emit telemetry：每次 mode 選擇 / signal 觸發 / cost actual / route activation
- 有 queryable surface（SQLite table / time-series file / runtime report）
- 不只是 commit-msg Cognitive Mode 報告 — 是 system-wide observability
- **Acceptance signal**：`ecosystem.telemetry.contract` 在 generated surfaces；至少一個 SQL 證明可查 mode selection 歷史

### J. Cognitive Resource Management — 真的有預算分配

- 每個 task 有 cognitive budget，runtime 在 thinking / context / execution / knowledge 之間實際分配
- Budget exceedance 觸發 compression / escalation / abort
- 預算消耗有 evidence trail
- **Acceptance signal**：`ecosystem.cognitive_budget_policy.contract` projected；commit-msg validator 比對 declared cost class vs actual budget consumption

---

## 現況評估（2026-05-28）

針對上述 10 criteria 的當前狀態。**沒有任何一條完全達標**；多數仍在 plan 或部分實作。

| Criteria | 狀態 | 關鍵證據 / 缺口 |
|---|---|---|
| A. Cognitive Economics | ⚠️ 規劃中 | [`plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) Phase 0–10 已展開；split cost model 已決議；無 live executable contract |
| B. Knowledge Evolution & Decay | ❌ 未開始 | ADR-004 feedback-promotion-pipeline 存在；無 atom usage telemetry、無 decay scoring、無自動 deprecation proposal |
| C. Discovery Evolution | ⚠️ 詞彙就位 | `discovery_mode` 已 candidate entry in [`knowledge/glossary/ai-skill.md`](../knowledge/glossary/ai-skill.md)；無 implementation；無 gap detection |
| D. Ecosystem Signals | ⚠️ 14/25 | `runtime/cognitive-modes-discovery.yaml` 目前 14 條 signal，全屬 user_keyword / file_diff / git_status / session / phase / token_budget 類；無 tool / model / memory / workflow / knowledge 派生 signal |
| E. Adaptive Runtime | ❌ 未開始 | Mode selection 仍純 deterministic signal → mode；無 pressure-based adaptation；無 feedback loop |
| F. Memory Economics & Decay | ❌ 未開始 | `memory/retrieval-governance/` 有 selective replay；無 decay scoring、無 cost-of-replay 模型 |
| G. Intelligence Evolution | ❌ 未開始 | `intelligence_mode` 同 C，是 glossary candidate term；無 fitness score、無 cross-intelligence synthesis |
| H. Multi-Agent Ecology | ❌ 未開始 | `.agent-goals/` 有 lock 機制但仍 single-owner assumption；無 multi-agent role contract |
| I. Telemetry Layer | ⚠️ 部分 | commit-msg Cognitive Mode 報告 + `runtime.db generated_surfaces` 是局部 telemetry；無 system-wide observability surface；無 mode-selection history |
| J. Cognitive Resource Management | ⚠️ 部分 | `runtime/cognitive-modes-token-budget.yaml` 有 token 預算；無 split-cost 分配、無 budget exceedance trigger |

**綜合判斷**：本系統正站在 Gen 3 → Gen 4 的**轉折起點**，不是中段。約 2/10 criteria 達到「⚠️ 部分」，其餘 ≥ 6 條未開始。要 graduate 為 `current` 至少還需要 4–8 個大型 plan 落地。

---

## 為什麼有人認為已經在往 Gen 4 走

User 評價提到本系統已開始出現第二代特徵。誠實對照：

| 觀察 | 屬實程度 | 證據 |
|---|---|---|
| 已在想 token / compression / tool / model / activation cost | ✅ | economics plan Phase 0–10 已寫，cognitive_cost split 已決議 |
| 已開始問哪些 decision 要進 memory、哪些 failure 要 promotion | ⚠️ partial | ADR-004 + ADR-005 提供 promotion 框架；但**自動化 promotion** 尚未動工 |
| 已覺得 routing 不夠、discovery 太弱 | ✅ | `discovery_mode` 已成為 glossary candidate，economics plan Phase 9 明列 knowledge acquisition signals |
| 已想 tool-derived / model-derived / knowledge-derived signals | ✅ vision | economics plan §Cognitive Mode discovery consumes economics-derived signals only |

**結論**：方向正確，**詞彙與設計已開始進入 Gen 4 範疇**（cognitive economics、ecosystem signals、knowledge acquisition、cognitive resource management），但 **executable contract 與 runtime 機械強制仍主要停在 Gen 3 範疇**。Gen 4 的 vocabulary 領先 Gen 4 的 enforcement 約一個世代。

---

## 與其他層的邊界（Gen 4 邊界提案）

| 層 | Gen 3 角色 | Gen 4 角色提案 |
|---|---|---|
| `runtime/` | execution + mode contract + obligation ledger | 退化為 `runtime/execution-engine/`，被 `ecosystem/orchestration/` 調度 |
| `ecosystem/`（新） | — | cognition / economics / intelligence-evolution / discovery / orchestration / telemetry / adaptation / evolution / governance-of-evolution |
| `workflow/` | execution flow | 被 `ecosystem/orchestration/` 與 `ecosystem/evolution/` 雙向消費（execute + decay） |
| `memory/` | selective replay | 加入 `ecosystem/cognition/memory-economics`；entry 有 fitness |
| `intelligence/` | 判斷準則 | 加入 `ecosystem/intelligence/` evolution；atom 有 fitness、cross-synthesis |
| `knowledge/` | navigation + glossary | 加入 `ecosystem/discovery/` adaptive lookup；atom usage telemetry |
| `governance/` | governance rules | 拆 `governance/runtime/` + `ecosystem/governance/`（meta-governance / evolution rules） |

**重點**：source-of-truth layer（`models/` / `tools/` / `memory/` / `workflow/` / `intelligence/` / `knowledge/`）保留 — ecosystem 不 own 真理，它 own **interaction、economics、evolution**。這條與 [economics integration plan §Source-of-truth layers](../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) 一致。

---

## 與 Active Plans 的關係

| Plan | 對 Gen 4 的貢獻 |
|---|---|
| [`plans/active/2026-05-25-1000-context-language-glossary-system.md`](../plans/active/2026-05-25-1000-context-language-glossary-system.md) | 提供 ecosystem 共享語彙基礎；Gen 4 概念（ecosystem、pressure_model、knowledge_mode、discovery_mode、intelligence_mode）已 candidate entry |
| [`plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) | 直接攻 criteria A（cognitive economics）+ D 部分（ecosystem signals）+ J 部分（resource management）。是 Gen 4 的核心 plan，但仍 draft，未實作 |

---

## 維護規則

- 本檔為 **vision document**，**不是** canonical entry。Gen 3 仍為 current。
- 修改本檔需先確認 graduation criteria 仍 machine-verifiable；不接受純散文加 criteria。
- 當某條 criteria 從 ❌ → ⚠️ → ✅，需在 §現況評估表記錄 evidence link（plan / commit / contract path）。
- 達 ≥ 7/10 ✅ 且其他 3 條有 active plan + entry condition 時，依 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) §1 啟動世代升級流程，本檔 status 從 `vision` 升為 `current`，Gen 3 文件降為 `historical`，重新評估每個 ADR 在 Gen 4 的延伸狀態。
- 新加入的 Gen 4 plan 必須遵循 [`governance/lifecycle/system-upgrade-governance.yaml`](../governance/lifecycle/system-upgrade-governance.yaml) §`define_runtime_trigger_flow`，不得以 routing-only 或 projection-only 宣稱已完成 runtime integration（見 plans/README.md §Runtime Execution Path 強化條款）。

---

## Gen 3 → Gen 4 的 graduation 不變條件

無論 criteria 達標多少，下列原則跨世代維持：

1. **Reference-First**（ADR-001）— Gen 4 仍以 source repository 為 single source-of-truth
2. **Knowledge vs Intelligence 分離**（ADR-002）
3. **Owner-layer canonical / projection derived** 邊界
4. **Doc-only trial 不能聲稱已完成 runtime integration**（governance §define_runtime_trigger_flow）
5. **Single owner per term**（glossary plan §Semantic Ownership）

Gen 4 是 layer 擴張 + 重心轉移，不是推翻 Gen 3。
