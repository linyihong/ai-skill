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

## 成熟度階梯（L1–L5）

Gen 演進對應 5 階成熟度。本系統目前位於 **L3.5 → L4 過渡期**：

| 階段 | 名稱 | 本系統位置 | 主要特徵 |
|---|---|---|---|
| L1 | Prompt / Skill System | 已超過 | 散落 prompt + skill 模板，無 runtime |
| L2 | Agent Runtime System | 已超過 | 有 agent loop + tool calling，無 governance |
| L3 | **Cognitive Execution System** | **完成大半（Gen 3 current）** | runtime governance、cognitive mode、failure-derived evolution、knowledge/workflow 分離 |
| L4 | **Cognitive Ecosystem System** | **過渡中（Gen 4 vision）** | knowledge activation runtime、cognitive economics engine、adaptive discovery、intelligence fitness、cross-agent ecology、cognitive telemetry |
| L5 | Adaptive Cognitive Civilization | 很後期 | 自演化、自治理、自分化的多 agent 生態 |

跨越 L3 → L4 的判定見下方 §The Final Test。

---

## The Final Test（單句判定 Gen 4 是否達成）

> **系統開始能自己決定**：該讀什麼、不該讀什麼、該啟動什麼 cognition、哪些知識太昂貴、哪些 workflow 已過時、哪些 memory 值得留下、哪些 heuristics 應 promotion、哪些 governance 過重 — 而不是由人類定義固定 workflow。

Gen 3 是 **human knows what to load**（routing-registry 列哪個 task intent 載入哪個 source）；Gen 4 是 **runtime discovers what must load**（economics + telemetry + activation contract 決定）。10+1 條 threshold criteria 是這句話的 machine-verifiable 拆解。

---

## 已完成的東西（Gen 3 substantial achievements）

避免 vision document 變成「全是 gap、看不到累積」。本系統在 Gen 3 已建立 4 個 substantial 基礎，**超過大多數 agent framework 的 baseline**：

| 已完成 | 內容 | 證據 |
|---|---|---|
| **Runtime Governance** | routing / validation / forbidden routes / heuristics / failure recovery / recovery escalation / governance gates / runtime validate / generated surfaces | `runtime/runtime.db` phase_machine + obligations + gates + 14 commit-msg validators |
| **Cognitive State Management** | cognitive mode（6 維 vector）/ compression / token budget / memory boundary / context loading / model profiles / discovery signals 14 條 | `runtime/cognitive-modes*.yaml` + ADR-008 + commit-msg cognitive contract block |
| **Failure-derived Evolution** | Failure → Scenario → Validation → Trace → Governance → Runtime Prevention 完整 pipeline | `validation/scenarios/failure-derived/*.yaml`（35+ scenarios）+ `enforcement/failure-patterns/` + `feedback/history/` + ADR-004 |
| **Knowledge / Workflow / Intelligence 分離** | knowledge atom（事實導航）/ workflow（執行流程）/ intelligence（判斷準則）/ runtime（contract）/ governance（約束）/ memory（演化）/ models（cognition economics）/ tools（execution capability）各層職責清楚 | ADR-002 / ADR-003 + 8 個 top-level layer 各自 README + routing-registry |

也就是說，本系統的 Gen 3 基礎已經是 **Cognitive Architecture 等級**，不是 skill repo 等級。Gen 4 是在此基礎之上**加 evolution / economics / ecology / telemetry**，不是推翻重做。

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
├── suppression/     # negative activation — what NOT to load / activate
├── ecology-boundary/  # what belongs in ecosystem vs runtime vs workflow vs docs
└── governance/      # governance-of-evolution (meta-governance)

runtime/
└── execution-engine/  # the former Gen 3 runtime, now a subsystem
```

`suppression/` 與 `ecology-boundary/` 是 2026-05-28 第三輪 review 增補：

- **suppression/** 防止「多 activate 即省 token」的反向錯覺。真正省成本不是 activate 對的東西，而是 **不 activate** 不該開的東西。`governance-minimality-small-task-v1` 與 `tier3-does-not-block-tier0-tier2-v1` 的精神在此層機械化。
- **ecology-boundary/** 防止 ecosystem inflation — 不是所有跨層問題都該進 ecosystem。明確定義「什麼進 ecosystem、什麼留 runtime / workflow / docs」，machine-verifiable 邊界。

---

## Cognitive Kernel Boundary（OS 類比）

第三輪 review 點出本系統實質已接近 **AI OS kernel architecture**。為避免 ecosystem 邊界失控，明列 OS 類比：

| 本系統 | OS kernel 角色 |
|---|---|
| `runtime/execution-engine/` | kernel core |
| `runtime/cognitive-modes*.yaml` | scheduler policy |
| `tools/` | device drivers |
| `memory/` | virtual memory |
| `governance/` | security policy |
| `workflow/` | process / task |
| `ecosystem/telemetry/` | observability subsystem |
| `ecosystem/activation/` | loader / dynamic linker |
| `ecosystem/economics/` | resource allocator |
| `ecosystem/discovery/` | dynamic linker（symbol resolution）|
| `ecosystem/suppression/` | OOM killer / cgroup limits |
| `ecosystem/ecology-boundary/` | kernel API contract（什麼算 kernel call）|

**Kernel Boundary Rule（防 ecosystem inflation）**：

一個 concept 進入 `ecosystem/` 必須同時滿足：

1. **Cross-layer**：跨 ≥ 2 個 source-of-truth layer（models / tools / memory / workflow / intelligence / knowledge）
2. **Resource interaction**：影響 cognition cost / activation / decay / pressure
3. **Machine-verifiable**：有 SQLite query / generated surface / scenario / validator 可驗
4. **Non-duplicating**：不重新定義既有 source layer 的 canonical 概念

不滿足 4 條的 concept 留在原 layer。Vision 詞彙不算進入 ecosystem，必須降為 `ecology-boundary/` 內的 candidate term 並走 graduation。

---

## Threshold Criteria（graduation 必達標準）

達成下列 **≥ 10/13 criteria** 且其餘 3 個有 active plan + 明確 entry condition，才考慮 graduate 本檔為 `current`。每條 criteria 必須是 **machine-verifiable**（generated surface / SQLite query / scenario evidence），不接受純散文宣稱。

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

### C. Discovery Evolution — 主動發現（6 種模式）

Discovery 不再只是 static routing；6 種模式至少 ≥ 4 種 active：

| Mode | 意思 |
|---|---|
| `architecture_archaeology` | 逆向探勘既有 system / repo，重建未文件化結構 |
| `bounded_context_discovery` | 從 code / docs 自動 mapping domain bounded context |
| `workflow_synthesis` | 從既有 workflow 組合新 workflow（不寫死）|
| `knowledge_gap_detection` | 偵測「應該有但目前沒有」的 atom |
| `heuristic_discovery` | 從 failure pattern 抽出新 heuristic candidate |
| `ecosystem_telemetry` | 從 runtime signals 推斷需要的 discovery |

要求：

- `discovery_mode` enum wired 到 `cognitive-modes-discovery.yaml`（成為 cognitive vector 的第七維）
- Gap detection 可在 plan execution 時主動發 promotion candidate
- **Acceptance signal**：≥ 1 個 scenario 證明 gap detection 主動發 candidate；`discovery_mode` 出現在 commit-msg cognitive contract block

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

### I. Telemetry Layer — 系統能觀察自己（6 類 telemetry）

不只是 commit-msg Cognitive Mode 報告 — 是 system-wide observability。6 類 telemetry 至少 ≥ 4 類 active：

| Telemetry | 內容 |
|---|---|
| `activation_frequency` | 哪些 intelligence / workflow / atom 常被啟用、哪些幾乎不用 |
| `token_burn_hotspots` | 哪些 route / workflow / tool 最貴 |
| `governance_friction` | 哪些 gate 常擋人、opt-out 用得多 |
| `recovery_loops` | 哪些 failure 反覆出現 |
| `stale_intelligence` | 哪些知識 / heuristic 過期或衰退 |
| `context_explosion` | 哪些 route 或載入策略容易炸 token |

要求：

- Runtime emit telemetry：每次 mode 選擇 / signal 觸發 / cost actual / route activation 都有紀錄
- 有 queryable surface（SQLite table / time-series file / runtime report）
- **Acceptance signal**：`ecosystem.telemetry.contract` 在 generated surfaces；至少一個 SQL 證明可查 mode selection 歷史 + token burn hotspots

### J. Cognitive Resource Management — 真的有預算分配

- 每個 task 有 cognitive budget，runtime 在 thinking / context / execution / knowledge 之間實際分配
- Budget exceedance 觸發 compression / escalation / abort
- 預算消耗有 evidence trail
- **Acceptance signal**：`ecosystem.cognitive_budget_policy.contract` projected；commit-msg validator 比對 declared cost class vs actual budget consumption

### K. Knowledge Activation Graph — 多輸入聲明式啟用契約

Gen 3 是「human knows what to load」（plan 作者寫死 candidate_sources）；Gen 4 是 **runtime discovers what must load**。但 activation 不是單純 signal → load 的 1:1 mapping，而是 **graph convergence**：多個輸入共同決定啟用集合。

```yaml
activation:
  inputs:                      # 多輸入收斂
    signals: [architecture_complexity_high, vendor_count_large]
    pressure: context_explosion_low
    economics: budget_available_high
    context_shape: cross_domain_synthesis
    failure_history: no_recent_archaeology_failure
    role: implementer  # 不是 reviewer / explorer
    memory_freshness: project_context_stale
  activate:
    - intelligence: vendor-integration-architecture
    - intelligence: bounded-context-analysis
    - governance: architecture-fit-governance
  economics:
    estimated_cost: HIGH
    why: cross-domain synthesis + bounded context discovery
```

要求：

- 至少 ≥ 5 個 activation contract 存在於 `ecosystem/activation/`，每個 contract 至少消費 **3 種以上**輸入（signals / pressure / economics / context_shape / failure_history / role / memory_freshness）
- Runtime 依 contract 自動載入，**不需 plan 作者手動列 candidate_sources**
- Contract 內 `economics.estimated_cost` 與 actual cost 有 telemetry 比對（criterion I）；偏差累積觸發 contract refinement
- **Acceptance signal**：`SELECT * FROM activation_graph` 在 runtime SQLite 可查；scenario 證明同 signal 在不同 pressure / role / memory_freshness 下產生不同 activation set

這條與 §The Final Test 直接對應：activation graph 完成 = system「自己決定該讀什麼、該啟動什麼 cognition」的 machine-verifiable proof。**Discovery ≠ Activation**：discovery 找候選，activation graph 才真正決策。

### L. Telemetry Economics — 觀察的成本不能大於被觀察

Criterion I（Telemetry Layer）有六種 telemetry，但本身會產生 runtime exhaust。若無 economics 控制，會出現 telemetry 成本 > execution 成本 / replay 成本 > 推理成本 / trace 成本 > workflow 成本的「過度自我觀察」病態。

要求：

- Telemetry 自帶 budget：每類 telemetry 有 retention / aggregation / decay 規則
- 高頻 telemetry 自動降頻或 aggregate（例如：mode selection 每 N 次 sample 一次，不是每次全紀錄）
- Telemetry decay 與 criterion F（memory decay）共用 lifecycle 模型
- 有 SQL / scenario 證明 telemetry 自身成本 ≤ 觀察對象成本的某個 ratio（例如 10%）
- **Acceptance signal**：`ecosystem.telemetry_economics.contract` 在 generated surfaces；scenario 證明 telemetry budget exceedance 觸發 aggregation / suppression / decay

### M. Cognitive Suppression Layer — 負向 activation

Activation 的正向決策（criterion K）只解決一半。**真正省 cognition 不是 activate 對的東西，是不 activate 不該開的東西**。本層機械化「什麼時候 disable validation / 跳過 deep discovery / 略過 architecture review」。

```yaml
suppression:
  when:
    - task_size_small
    - risk_low
    - failure_history_clean
  disable:
    - workflow: heavy-governance-chain
    - intelligence: cargo-cult-ddd
    - validation: exhaustive-discovery
  why: small low-risk task does not warrant Tier-3 governance
```

要求：

- 至少 ≥ 3 個 suppression contract 在 `ecosystem/suppression/`
- 既有 scenarios `governance-minimality-small-task-v1` 與 `tier3-does-not-block-tier0-tier2-v1` 被 suppression contract 機械引用
- 反向：governance / validation 被 suppress 的事件有 telemetry 紀錄（criterion I 連動）
- **Acceptance signal**：`SELECT * FROM suppression_events` 在 runtime SQLite 可查；scenario 證明 small low-risk task 被 suppress 後 cognitive_cost 顯著下降

---

## Migration Roadmap（Phase A / B / C）

從 Gen 3.5（目前）到 Gen 4 graduate，建議分 3 階段。每階段對應若干 criteria。

### Phase A — Strengthen Runtime（**目前位置**）

繼續 Gen 3 的深化：cognitive modes、tool signals、model economics、memory boundary、routing、discovery。

- Active：[`plans/active/2026-05-25-1000-context-language-glossary-system.md`](../plans/active/2026-05-25-1000-context-language-glossary-system.md) Phase 6 加 glossary runtime auto-detect
- Active：[`plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md`](../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) 把 economics 與 ecosystem signals 拉到 runtime
- 觸及 criteria：A（partial）/ D（partial）/ J（partial）

### Phase B — Cognitive Activation Layer（**下一個大爆炸點**）

建立 declarative activation runtime（criterion K 的實作）。signals → intelligence / workflow / governance / memory 自動 activation。這是 Gen 4 的**結構性轉折**：discovery 從 static route lookup 升級為 contract-driven activation。

- 主要 criterion：**K**（必達）+ C（discovery evolution，連帶 unlock 4+ 個 discovery 模式）+ D（signals 從 14 → 25+ 條，涵蓋 5 大類）
- 連帶要求：Phase B 一旦上線，使所有 active plans 在 §Runtime Execution Path 的 trigger chain 都能引用 activation contract，而不是手寫 candidate_sources
- 預估規模：與 glossary system + economics system 同級或更大

### Phase C — Cognitive Economics Engine（**真正進入 Gen 4 ecosystem**）

完成後系統能：評估 cognition 成本、動態調整、淘汰 knowledge、選擇 workflow、管理 memory、控制 governance friction、預測 token burn。

- 主要 criteria：**A**（完成）+ **B**（knowledge decay）+ **F**（memory economics）+ **G**（intelligence fitness）+ **I**（telemetry）
- 達 ≥ 8/11 criteria 時，graduate 為 `current`

### Phase D（可選）— Cross-Agent Ecology

達到 L5（Adaptive Cognitive Civilization）的前哨。對應 criterion **H**。可在 Gen 4 graduate 後再啟動。

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
| K. Knowledge Activation Graph | ❌ 未開始 | `knowledge/runtime/routing-registry.yaml` 是 task_intent → primary_source 的 manual map，**不是** multi-input signal → activation 的 declarative graph；無 economics / pressure / role / memory_freshness aware activation |
| L. Telemetry Economics | ❌ 未開始 | 目前 telemetry 只有 commit-msg cognitive mode block + generated reports；無 retention / aggregation / decay budget；無 self-observation cost 控制 |
| M. Cognitive Suppression Layer | ⚠️ 詞彙就位 | scenarios `governance-minimality-small-task-v1` + `tier3-does-not-block-tier0-tier2-v1` 已 enforce 部分 suppression；無 declarative suppression contract、無 suppression telemetry |

**綜合判斷**：本系統正站在 Gen 3 → Gen 4 的**轉折起點**，不是中段。**0 ✅ / 5 ⚠️ partial / 8 ❌ 未開始**（共 13 criteria）。Gen 4 詞彙領先 Gen 4 enforcement 約一個世代。要 graduate 為 `current` 至少還需要 4–8 個大型 plan 落地，重心會在 Phase B（K + C + D）、Phase C（A + B + F + G + I）、Phase B′（L + M，與 telemetry 同步）。

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
- 達 **≥ 10/13 ✅** 且其他 3 條有 active plan + entry condition 時，依 [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) §1 啟動世代升級流程，本檔 status 從 `vision` 升為 `current`，Gen 3 文件降為 `historical`，重新評估每個 ADR 在 Gen 4 的延伸狀態。
- 新加入的 Gen 4 plan 必須遵循 [`governance/lifecycle/system-upgrade-governance.yaml`](../governance/lifecycle/system-upgrade-governance.yaml) §`define_runtime_trigger_flow`，不得以 routing-only 或 projection-only 宣稱已完成 runtime integration（見 plans/README.md §Runtime Execution Path 強化條款）。

---

## Watch-Out List（Gen 4 路上會撞到的牆）

第三輪 review 識別 5 個下一階段必然遇到的反模式。寫進本檔以便未來 plan 主動規避：

| Wall | 風險敘述 | 緩解 |
|---|---|---|
| 1. Discovery confused with Activation | 把 routing-registry 當 activation graph，導致系統仍是「人決定載入什麼」 | criterion K 強制 multi-input activation contract；route ≠ activation |
| 2. Workflow inflation | workflow 越寫越大、越 rigid，最後變 BPMN；agent 失去自主性 | 核心心智模型轉為 cognitive resource management；workflow 是 resource 之一不是中心 |
| 3. Ecosystem boundary inflation | 所有跨層概念都塞 ecosystem，最後 ecosystem 變第二個 monolith | §Cognitive Kernel Boundary 4 條 rule 機械強制；`ecosystem/ecology-boundary/` 為守門 |
| 4. Telemetry explosion | telemetry 成本 > execution 成本，系統陷入 over-self-observation | criterion L（Telemetry Economics）強制 retention / aggregation / decay budget |
| 5. Positive-activation bias | 只想 activate 對的東西，忽略不 activate 不該開的東西 | criterion M（Cognitive Suppression Layer）機械化負向 activation |

每條 wall 對應一個 criterion；新 plan 在 §Decision Rationale 應 cite 自己對應緩解的 wall。

## Public Positioning（系統的真實層級）

外部觀察者已指出本系統實質不是 agent framework 而是 **AI runtime architecture / cognitive operating architecture / adaptive cognition infrastructure**。代價是：

- complexity 會比一般 framework 大數量級
- governance 會越來越難（meta-governance 比 governance 更難）
- layering 重要性遠超 features
- economics 比 features 更重要（cognition resource allocation 是核心問題）

本檔接受此定位。Gen 4 不追求成為更好的 agent framework，追求成為 **可演化的 cognitive OS**。

---

## Gen 3 → Gen 4 的 graduation 不變條件

無論 criteria 達標多少，下列原則跨世代維持：

1. **Reference-First**（ADR-001）— Gen 4 仍以 source repository 為 single source-of-truth
2. **Knowledge vs Intelligence 分離**（ADR-002）
3. **Owner-layer canonical / projection derived** 邊界
4. **Doc-only trial 不能聲稱已完成 runtime integration**（governance §define_runtime_trigger_flow）
5. **Single owner per term**（glossary plan §Semantic Ownership）

Gen 4 是 layer 擴張 + 重心轉移，不是推翻 Gen 3。
