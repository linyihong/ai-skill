# Context Cost Optimization Plan

本文件是 [`next-stage-upgrade-plan.md`](../active/next-stage-upgrade-plan.md) 的**成本優化補充規劃**。它不取代原有架構分層規劃，而是針對目前最迫切的 **token 成本與 context 過重問題**，提出具體改造路徑。

## 問題摘要

現有 repository 已成功「知識 OS 化」，但產生了新的核心問題：

| 問題 | 表現 | 根本原因 |
| --- | --- | --- |
| **成本高** | 每次 session 大量 token 消耗 | Bootstrap 過大，Agent 每次重新 ingest 大量 operating rules |
| **Token 爆炸** | Max mode 容易爆 token | 深度 reasoning × 大 context 形成核爆組合 |
| **AI 漏規則** | Agent 不讀或漏讀關鍵規則 | Routing 不夠精準，沒有 lazy-load 機制 |
| **AI 失憶** | 長對話中遺失上下文 | Memory layer 不完整，無 TTL 機制 |
| **Reasoning 太慢** | Agent 花大量時間重新理解 repo | Rule graph 過重，無摘要層先行 |

## 核心診斷

```
README.md + enforcement/README.md + Default Bootstrap
    ↓
已變成「AI OS Kernel」
    ↓
問題是：Kernel 太大
```

現有 Default Bootstrap 包含 12 條規則，加上 `rule-weight`、`dependency-reading`、`linked-updates`、`failure-learning` 等形成**高密度 instruction graph**。GPT-5.5 級模型會自動做 recursive context linking、cross-reference，導致 token 快速膨脹。

## 改造原則

1. **Runtime-first，不是 documentation-first**：現有 `runtime/` 目錄偏向設計文件，需要轉成可執行的 routing 行為。
2. **Lazy-load by default**：只有當前任務需要的知識才載入，其餘 deferred。
3. **Summary layer 先行**：Agent 先讀 300-500 token summary，需要才展開全文。
4. **Context TTL**：每條 context 只活一個 task 或 session，不永久留在 context graph。
5. **Cost-aware routing**：每個知識單元標註 estimated tokens，讓 AI 自己學會省 token。

---

## Phase 1：立即省錢（第一優先）✅ 已實作

### 1.1 Bootstrap 極小化 ✅ 已實作

**現狀**：Default Bootstrap 載入 12 條完整規則，形成高密度 instruction graph。

**改造**：拆分為 Core Bootstrap + Lazy-load rules。

#### Core Bootstrap（必讀，極小集合）

```
CORE_BOOTSTRAP.md
  1. rule-weight          → 規則衝突時如何判斷優先序
  2. dependency-reading   → 依賴文件讀取鐵則
  3. conversation-goal-ledger → 對話目標閉環
```

這三條是**所有任務的最低共同需求**，無法 lazy-load。

#### Lazy-load Rules（只在特定情境載入）

| 規則 | 觸發條件 | 載入時機 |
| --- | --- | --- |
| `linked-updates` | multi-file change, architecture update | 偵測到跨文件改動時 |
| `failure-learning-system` | user 指出 agent 失誤、close-loop gap | 發生 failure 時 |
| `decision-efficiency` | 任務有多條可行路線需選擇 | 決策分歧時 |
| `tool-neutral-documentation` | 建立或修改可重用文件 | 寫文件時 |
| `governance/document-sizing` | 文件接近拆分門檻 | 文件變大時 |
| `document-todo-list` | 文件有未完成項目 | 修改文件時 |
| `goal-action-validation` | 重要工作單元需要驗證 | 執行關鍵步驟時 |
| `neutral-language` | 撰寫或審查文件用語 | 寫文件時 |
| `sanitization` | 撰寫 feedback lesson | 寫 lesson 時 |
| `authorization-scope` | 涉及授權邊界 | 分析外部系統時 |
| `cross-skill-references` | 引用其他 skill | 跨 skill 工作時 |
| `feedback-lessons` | 撰寫或 promotion lesson | 寫 lesson 時 |

**預期效益**：Token 減少 40-60%。

### 1.2 README 拆分 ✅ 已實作

**現狀**：根 `README.md` 196 行，接近 Operating System Specification，每次 session 重讀。

**改造**：拆成四份文件：

| 文件 | 用途 | 目標行數 |
| --- | --- | --- |
| `README.md` | 超短入口，只留 layout 與 quickstart | ≤ 80 行 |
| `BOOTSTRAP.md` | Runtime minimal bootstrap（Core Bootstrap） | ≤ 50 行 |
| `ARCHITECTURE.md` | 深度設計與 roadmap | 不限 |
| `GOVERNANCE.md` | 管理規範與貢獻指南 | 不限 |

### 1.3 Rule Lazy-load 機制 ✅ 已實作

**現狀**：`enforcement/README.md` 的 Default Bootstrap 要求每次載入 12 條規則。

**改造**：在 `enforcement/README.md` 中引入 **runtime activation model**：

```yaml
rule:
  id: linked-updates
  activation:
    when:
      - multi-file-change
      - architecture-update
  cost:
    estimated_tokens: 800
  load_strategy: lazy
```

Agent 不預設讀 lazy-load rules，只在符合 `activation.when` 條件時才 activate。

### 1.4 Knowledge Summary Layer ✅ 已實作

**現狀**：`knowledge/summaries/` 已有 6 個 summaries，但尚未成為 agent 的預設載入策略。

**改造**：

1. **Summary-first routing**：Agent 先讀 summary（300-500 tokens），需要才展開全文。
2. **擴充 summaries 覆蓋範圍**：每個 skill、shared rule、architecture 文件都應有對應 summary。
3. **Summary 加入 context_cost 欄位**：讓 agent 在載入前就知道成本。

#### Summary 格式升級

```markdown
## <atom-id>

| 欄位 | 值 |
| --- | --- |
| Atom ID | `skill.apk-analysis` |
| Source path | `skills/apk-analysis/SKILL.md` |
| Lifecycle | `validated` |
| Summary | 授權 APK 流量分析、動態抓包、Flutter/Dart AOT、response 解碼。 |
| When to read | 使用者要求 APK traffic/runtime/response 分析時。 |
| Do not use for | 不可取代完整 SKILL.md 或 WORKFLOW.md。 |
| Context cost | ~350 tokens |
| Estimated full cost | ~4500 tokens |
| Validation signal | Skill entrypoint links 可解析，authorization gate 已套用。 |
| Last checked | 2026-05-11 |
```

---

## Phase 2：架構升級（第二優先）✅ 已實作

### 2.1 Runtime Context Router ✅ 已實作

**現狀**：`runtime/` 目錄偏向設計文件，`knowledge/runtime/routing-registry.yaml` 已有 8 筆 records 但尚未成為真正的 runtime router。

**改造**：建立 **Runtime Context Router** 系統：

```
runtime/
  router/              ← 新增：context routing 邏輯
    README.md          ← routing 決策流程
    activation-rules.yaml  ← 何時 activate 哪些知識
    cost-budget.yaml   ← session token budget 管理
  routing/             ← 既有：設計文件
  context/             ← 既有：context pruning 設計
```

#### Activation Rules 範例

```yaml
# runtime/runtime.db
activation_rules:
  - rule_id: linked-updates
    when:
      - type: file_change
        pattern: "**/*.md"
        count: ">=2"
      - type: task_intent
        matches: ["architecture-update", "migration", "refactor"]
    load:
      strategy: lazy
      priority: P1
      estimated_tokens: 800

  - rule_id: failure-learning-system
    when:
      - type: user_signal
        matches: ["失誤", "漏讀", "忘記", "錯誤", "failure", "miss"]
      - type: validation_gap
        matches: ["close-loop", "commit", "push", "sync"]
    load:
      strategy: lazy
      priority: P1
      estimated_tokens: 1200

  - rule_id: decision-efficiency
    when:
      - type: task_complexity
        routes: ">=3"
      - type: user_signal
        matches: ["選擇", "路線", "priority", "先做哪個"]
    load:
      strategy: lazy
      priority: P2
      estimated_tokens: 600
```

### 2.2 Context Cost Metadata ✅ 已實作

**現狀**：`metadata/schema.md` 已有 `context_cost` 欄位（low/medium/high），但缺乏精確的 token 估算與 load strategy。

**改造**：在 metadata schema 中新增 `context_cost` 與 `load_strategy` 詳細欄位：

```yaml
# 新增到 metadata/schema.md 的選填欄位
context_cost:
  estimated_tokens: 1200        # 精確 token 估算
  estimated_read_time: "2-3 min"  # 人類閱讀時間
  breakdown:                     # 成本細項
    - header_and_navigation: 200
    - core_content: 700
    - examples: 200
    - validation: 100

load_strategy:
  default: lazy                  # preload | lazy | on_demand
  cacheable: true                # 是否可在 session 內 cache
  ttl:                           # context TTL
    session: 1                   # 只活一個 task
    # task: 3                   # 活 3 個 task
    # conversation: true        # 活整個對話

priority:
  criticality: high              # critical | high | medium | low
  cost_aware_routing: true       # routing 時考慮 token 成本

runtime_behavior:
  preload: false                 # 是否預設載入
  lazy_load: true                # 是否 lazy-load
  cacheable: true                # 是否可 cache
  deferrable: true               # 是否可延後載入
```

### 2.3 Skill Index ✅ 已實作

**現狀**：Agent 掃描 `skills/` 目錄尋找相關 skill，缺乏結構化索引。

**改造**：建立 `skills-index.yaml`，讓 AI 先路由 skill，而不是掃整個 repo。

```yaml
# skills-index.yaml
skills_index_version: v1
description: >
  Skill routing index。Agent 先比對 task intent 與 triggers，
  找到對應 skill 後再讀 entrypoint，不掃描整個 skills/ 目錄。

skills:
  - id: apk-analysis
    name: APK Analysis
    triggers:
      - apk
      - mitm
      - flutter
      - frida
      - traffic
      - proxy
      - reverse-engineering
    cost:
      estimated_tokens: 4500
      context_load: high
    entrypoint: skills/apk-analysis/SKILL.md
    summary: knowledge/summaries/apk-analysis-pilot.md

  - id: app-development-guidance
    name: App/API/Embedded Development Guidance
    triggers:
      - api
      - backend
      - architecture
      - mobile
      - security
      - firmware
    cost:
      estimated_tokens: 3500
      context_load: medium
    entrypoint: skills/app-development-guidance/SKILL.md
    summary: null

  - id: travel-planning
    name: Travel Planning
    triggers:
      - travel
      - itinerary
      - route
      - transportation
    cost:
      estimated_tokens: 2800
      context_load: medium
    entrypoint: skills/travel-planning/SKILL.md
    summary: null
```

### 2.4 Context TTL System ✅ 已實作

**現狀**：Context 一旦載入就永久留在 context graph，導致 Agent 越來越肥。

**改造**：建立 Context TTL 系統，定義每條 context 的生命週期：

```yaml
# runtime/runtime.db
ttl_policy_version: v1

default_ttl:
  session: 1    # 預設只活一個 task/session

rules:
  - id: shared-rule-bootstrap
    ttl:
      session: 1
    reason: Bootstrap 規則每個 session 重新載入即可

  - id: skill-workflow
    ttl:
      task: 1
    reason: Skill workflow 只在相關 task 中需要

  - id: feedback-lesson
    ttl:
      task: 1
    reason: Feedback lesson 只在相關操作時需要

  - id: architecture-roadmap
    ttl:
      session: 1
    reason: Roadmap 每個 session 確認一次即可

  - id: runtime-registry
    ttl:
      conversation: true
    reason: Routing registry 可在整個對話中 cache

  - id: knowledge-summary
    ttl:
      conversation: true
    reason: Summary 輕量且可跨 task 使用
```

---

## Phase 2.5：Provider Prompt Cache Alignment（規範層）✅ 已實作

### 2.5.1 為什麼要放在 Phase 2 與 Phase 3 之間

Phase 2 已建立 `context_cost`、`load_strategy`、`cacheable` 與 TTL 等 context metadata；Phase 3 則會處理 Semantic Retrieval、Episodic Memory 與 Runtime Orchestration。Provider prompt cache 命中率介於兩者之間：它需要 Phase 2 的 metadata 才能判斷哪些 context 穩定可重用，但不需要等 Phase 3 的 retrieval / memory 自動化完成。

因此本階段先定義 **prompt layout 與 cache eligibility 規範**，讓 agent 立刻能用人工規則提高 provider prompt cache 命中率；Phase 3 再把這些規範接進 runtime orchestration 自動執行。

**實作狀態**：規範層已建立：

- [`runtime/context/prompt-cache-playbook.md`](../../runtime/context/prompt-cache-playbook.md)
- [`enforcement/prompt-cache-efficiency.md`](../../enforcement/prompt-cache-efficiency.md)
- `metadata/schema.md` 的 `context_cost.provider_cache` 欄位
- `runtime/runtime.db` 與 `knowledge/runtime/routing-registry.yaml` 的 lazy-load route

Runtime 自動排序與 provider cache hit observability 仍保留給 Phase 3 Runtime Orchestration。

### 2.5.2 問題

目前 `cacheable: true` 只表示 context 可在 session / conversation 內重用，尚未區分「是否適合放進 provider prompt cache 的穩定前綴」。若每次 prompt 都把時間戳、git status、open files、tool output 或 task-specific evidence 插在固定規則前面，provider 端即使支援 prefix caching，也會因前綴 churn 而降低 cache hit。

### 2.5.3 改造目標

建立 Provider Prompt Cache Alignment playbook，明確區分：

| 類型 | 放置位置 | 範例 | 目標 |
| --- | --- | --- | --- |
| Stable prefix | Prompt 前段，順序固定 | Core Bootstrap、固定 runtime initialization、穩定 routing policy | 提高 prefix cache hit |
| Semi-stable middle | Stable prefix 之後 | task intent 對應的 summary、route-specific rules | 控制重用與任務相關性 |
| Volatile suffix | Prompt 後段 | 使用者當前要求、git status、open files、tool output、時間戳 | 避免污染穩定前綴 |

### 2.5.4 建議新增的規範

1. **Prompt cache playbook**
   - 建議位置：`runtime/context/prompt-cache-playbook.md`
   - 定義 stable prefix、semi-stable middle、volatile suffix 的排序規則。
   - 說明 Cursor / Claude 使用時，哪些內容應避免插入固定前綴。

2. **Context metadata extension**
   - 建議擴充 `metadata/schema.md`：
     - `provider_cache_candidate: true | false`
     - `prefix_stability: stable | semi_stable | volatile`
     - `cache_position: prefix | middle | suffix`
     - `churn_risk: low | medium | high`

3. **Runtime ordering rule**
   - stable context 的順序必須固定，不因 task 改變而重排。
   - volatile context 一律放在 suffix，不可插入 bootstrap / rule prefix 中間。
   - 若需要插入新的常駐規則，應追加到 stable prefix 的固定區塊末端，並記錄 prefix churn reason。

4. **Observability**
   - 在 token budget / context health 之後加入 prompt-cache 指標：
     - `stable_prefix_size`
     - `prefix_churn_count`
     - `provider_cache_candidate_tokens`
     - `volatile_prefix_violation_count`

### 2.5.5 不做事項

- 不把 provider prompt cache 視為正確性保證；它只是一個成本與延遲優化。
- 不為了 cache hit 犧牲 required dependencies、safety rules 或 source-of-truth validation。
- 不把高變動 evidence、使用者私有輸入、工具輸出或 runtime status 放入 stable prefix。
- 不要求 Phase 3 前完成自動化；本階段只先定義人工可遵守的 layout 與 metadata 規範。

### 2.5.6 與 Phase 3 的銜接

Phase 3 的 Runtime Orchestration 應讀取本階段定義的 metadata，將 context 自動排成：

```text
stable prefix → semi-stable task context → volatile suffix
```

Semantic Retrieval 負責選出相關 context；Provider Prompt Cache Alignment 負責決定 context 在 prompt 裡的位置與穩定性約束。兩者互補，不互相取代。

---

## Phase 3：真正 AI OS（第三優先）⏳ 待實作

### 3.1 Semantic Retrieval ⏳ 待實作

建立語意檢索層，讓 agent 用自然語言查詢知識，而不是依賴目錄掃描：

- 整合 SQLite / FTS runtime index（已有 prototype）
- 加入 embedding-based retrieval（future）
- Query ranking 與 context-aware scoring

### 3.2 Episodic Memory ⏳ 待實作

建立情節記憶層，讓 agent 記住「過去怎麼解決類似問題」：

- `memory/episodic/`：記錄過去 task 的關鍵決策與結果
- 每個 episodic record 包含：context、decision、outcome、token cost
- 支援 similarity-based retrieval

### 3.3 Runtime Orchestration ⏳ 待實作

建立 runtime orchestration layer，自動管理 context loading、pruning 與 TTL：

- `runtime/orchestration/`：context lifecycle management
- 自動根據 task intent 選擇 model profile
- 自動 prune 過期 context

### 3.4 Multi-model Routing ⏳ 待實作

根據 task 複雜度自動選擇模型：

| Task 類型 | 建議模型 | Context 策略 |
| --- | --- | --- |
| Bootstrap | small/fast | Core Bootstrap only |
| Routing | medium | Summary + index |
| Architecture reasoning | max | Full source + graph |
| Implementation | medium | Skill + workflow |
| Autocomplete | cheap | Index + checklist |

---

## Phase 4：Runtime Quality & Safety（已實作）

以下為根據外部 review 建議，優先於 Phase 3 實作的 **Runtime Quality & Safety** 層。這些元件直接控制 token 用量、context 健康度與系統穩定性。

### 4.1 Token Budget System

**檔案**：[`runtime/runtime.db`](../runtime/runtime.db)

**目的**：為每個 session 設定 token 預算上限，防止 token 爆炸。

**核心機制**：

| 元件 | 值 | 說明 |
| --- | --- | --- |
| Default max_tokens | 120,000 | 跨模型預設上限 |
| Warning threshold | 70% | 觸發 prune 建議 |
| Hard stop threshold | 90% | 強制 halt + prune |
| Per-layer budget | bootstrap/skill_index/activation_rules/summaries/full_source/tool_output/conversation | 每層分配 token 配額 |

**預期效益**：Token 使用量可預測，不再因深度 reasoning 爆 token。

### 4.2 Context Health Score

**檔案**：[`runtime/runtime.db`](../runtime/runtime.db)

**目的**：量化 context 健康度，在 context 品質惡化前主動介入。

**4 個維度**：

| 維度 | 權重 | 說明 |
| --- | --- | --- |
| Relevance | 0.35 | 當前 context 與 task 的相關性 |
| Duplication | 0.20 | 重複內容比例 |
| Staleness | 0.25 | 過期 context 比例 |
| Conflict | 0.20 | 規則衝突比例 |

**複合分數公式**：`relevance × 0.35 + (1 - duplication) × 0.20 + (1 - staleness) × 0.25 + (1 - conflict) × 0.20`

**門檻**：healthy ≥ 0.75, warning ≥ 0.50, critical < 0.50

### 4.3 Circuit Breaker

**檔案**：[`runtime/runtime.db`](../runtime/runtime.db)

**目的**：防止 agent 陷入無限迴圈、工具爆炸、context 失控。

**5 個 Guards**：

| Guard | 門檻 | 行為 |
| --- | --- | --- |
| Recursive depth | max 4 | 超過則 halt |
| Tool calls | 20/task, 100/session, 15/5min | 超過則 warn + suggest decomposition |
| Context growth | 30%/task, 80%/session | 超過則 force prune |
| Hallucination risk | 4 factors, threshold 0.7 | 超過則 halt + suggest source read |
| Conflict rules | 偵測到衝突 | warn + suggest resolution |

### 4.4 Context Pollution Detection

**檔案**：[`runtime/runtime.db`](../runtime/runtime.db)

**目的**：偵測 context 污染訊號，在 context 品質嚴重惡化前自動歸檔。

**5 個訊號**：

| 訊號 | 門檻 | 說明 |
| --- | --- | --- |
| Conversation length | 50 turns | 對話過長 |
| Repetitive edits | 5 edits | 同一檔案反覆修改 |
| Module count | 20 modules | 載入過多模組 |
| Cross-reference depth | 5 layers | 依賴鏈過深 |
| Token utilization | 85% | 接近 token 上限 |

**Critical 時**：halt agent + auto-archive session 到 `memory/working/session-archive-{timestamp}.md`

### 4.5 Tool Metadata & Lazy Activation

**檔案**：[`tools/metadata/README.md`](../tools/metadata/README.md)、[`tools/routing/README.md`](../tools/routing/README.md)

**目的**：為每個工具標註 cost、risk、activation strategy，實現工具層級的 lazy loading 與爆炸偵測。

**Tool Metadata 欄位**：

| 欄位 | 說明 |
| --- | --- |
| cost.avg_input_tokens | 平均輸入 token |
| cost.avg_output_tokens | 平均輸出 token |
| cost.risk | recursive_expansion / side_effects |
| contexts | 工具使用的 context 類型 |
| activation | preload / lazy / on_demand |
| compression | 支援的 compression level |

**Tool Explosion Detection**：recursive_search、repetitive_read、tool_chain_too_long、output_too_large

### 4.6 Tool Output Compression

**檔案**：[`tools/compression/README.md`](../tools/compression/README.md)

**目的**：根據 context 健康度動態壓縮工具輸出，減少 token 消耗。

**4 個 Compression Levels**：

| Level | 壓縮率 | 適用情境 |
| --- | --- | --- |
| raw | 1.0x | Context 健康度高 |
| summary | 0.2-0.3x | Context 健康度 warning |
| structured | 0.1-0.2x | Context 健康度 critical |
| minimal | 0.05-0.1x | Token budget 接近 hard stop |

**Per-output-type 策略**：stack trace（top 5 frames）、JSON（relevant fields only）、Git diff（file list + top 3 lines）、Log（ERROR/WARN only）、Search results（top 5）、File content（summary-first）

### 4.7 Memory Architecture（子層拆分）

**檔案**：[`memory/working/README.md`](../memory/working/README.md)、[`memory/summary/README.md`](../memory/summary/README.md)、[`memory/decision/README.md`](../memory/decision/README.md)

**目的**：將單一 memory 層拆分為三個子層，各自有不同的生命週期與 token 策略。

| 子層 | 生命週期 | Token 策略 | 用途 |
| --- | --- | --- | --- |
| Working | Session-local, discardable | 不保留 | 當前 session 的工作記憶 |
| Summary | Compressed, ≤500 tokens | 低 | 壓縮的 session 歷史 |
| Decision | Immutable, numbered | 極低 | 輕量 ADR（架構決策記錄） |

### 4.8 Decision System（ADR）

**檔案**：[`decisions/README.md`](../decisions/README.md)

**目的**：建立 Architecture Decision Records（ADR）系統，記錄關鍵架構決策，避免重複討論。

**ADR Lifecycle**：proposed → accepted → deprecated → superseded

**ADR 命名**：`ADR-{number}-{short-title}.md`

### 4.9 Anti-patterns

**檔案**：[`anti-patterns/README.md`](../anti-patterns/README.md) + 5 個 anti-pattern 文件

**目的**：記錄已知失效模式，讓 agent 快速辨識並避免。

| Anti-pattern | Severity | 說明 |
| --- | --- | --- |
| Context Explosion | critical | Context 無限制增長 |
| Recursive Tool Loop | critical | 工具反覆呼叫無進展 |
| Hallucination Loop | critical | 無 canonical source 時過度推理 |
| Stale Summary | high | Summary 與 source 不同步 |
| Skill Pollution | high | 不相關 skill 浪費 token |

### 4.10 Skills Metadata v2

**檔案**：[`skills-index.yaml`](../skills-index.yaml)（version: v2）

**目的**：為每個 skill 加入 weight、domains、dependencies、conflicts、priority.runtime，支援 runtime relevance scoring 與 conflict detection。

**新增欄位**：

| 欄位 | 說明 |
| --- | --- |
| weight | 0.0-1.0，relevance scoring 權重 |
| domains | 領域標籤，用於 domain-based routing |
| dependencies | 依賴文件路徑 |
| conflicts | 衝突 skill ID |
| priority.runtime | high / medium / low，runtime 載入優先序 |

---

## 遷移路徑

### Step 1：建立 Core Bootstrap（立即）✅ 已實作

1. 建立 [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) — 只含 rule-weight、dependency-reading、goal-ledger
2. 更新 `enforcement/README.md` — 加入 runtime activation model
3. 更新根 `README.md` — 縮短為超短入口

### Step 2：建立 Skills Index（立即）✅ 已實作

1. 建立 [`skills-index.yaml`](../skills-index.yaml)
2. 更新 `knowledge/indexes/README.md` 引用 skills index

### Step 3：擴充 Summary Layer（1-2 週）✅ 已實作

1. 為每個 shared rule 建立 summary
2. 為每個 skill 建立 summary
3. 為每個 architecture 文件建立 summary
4. 在 summary 中加入 context_cost 與 estimated_tokens

### Step 4：建立 Runtime Router（2-3 週）✅ 已實作

1. 建立 `runtime/runtime.db`
2. 建立 `runtime/router/cost-budget.yaml`
3. 建立 `runtime/runtime.db`
4. 更新 `knowledge/runtime/routing-registry.yaml` 加入 cost metadata

### Step 5：全面導入 Context Cost Metadata（3-4 週）✅ 已實作

1. 更新 `metadata/schema.md` 加入 context_cost、load_strategy、ttl 欄位
2. 為所有 Knowledge Atoms 補上 cost metadata
3. 更新 validation scripts 檢查 cost metadata

### Step 6：建立 Runtime Quality & Safety（已實作）

1. 建立 Token Budget System → [`runtime/runtime.db`](../runtime/runtime.db)
2. 建立 Context Health Score → [`runtime/runtime.db`](../runtime/runtime.db)
3. 建立 Circuit Breaker → [`runtime/runtime.db`](../runtime/runtime.db)
4. 建立 Context Pollution Detection → [`runtime/runtime.db`](../runtime/runtime.db)
5. 建立 Tool Metadata & Lazy Activation → [`tools/metadata/README.md`](../tools/metadata/README.md)、[`tools/routing/README.md`](../tools/routing/README.md)
6. 建立 Tool Output Compression → [`tools/compression/README.md`](../tools/compression/README.md)
7. 建立 Memory Architecture 子層 → [`memory/working/README.md`](../memory/working/README.md)、[`memory/summary/README.md`](../memory/summary/README.md)、[`memory/decision/README.md`](../memory/decision/README.md)
8. 建立 Decision System（ADR）→ [`decisions/README.md`](../decisions/README.md)
9. 建立 Anti-patterns → [`anti-patterns/README.md`](../anti-patterns/README.md)
10. 升級 Skills Metadata v2 → [`skills-index.yaml`](../skills-index.yaml)（weight、domains、dependencies、conflicts、priority.runtime）

### Step 7：建立 Runtime Pipeline（已實作）

將所有 Runtime Quality & Safety 元件串接成可執行的 orchestration flow：

1. 建立 Pipeline 概覽 → [`runtime/README.md`](../runtime/README.md) — 元件接線圖、跨階段通訊表（10 個觸發事件）
2. 建立 Session Lifecycle → [`runtime/runtime.db`](../runtime/runtime.db) — 4 階段（bootstrap → routing → execution → close-loop），每階段有 token budget、guard chain、進入/離開條件
3. 建立 Progressive Context Expansion → [`runtime/runtime.db`](../runtime/runtime.db) — 4 層級（summary → module summary → detailed source → raw source），每層有 cache policy、entry/exit conditions
4. 建立 Guard Chain → [`runtime/runtime.db`](../runtime/runtime.db) — 每 stage 的 guard 執行順序（ordered by severity）、檢查頻率（per_tool_call / per_task / per_edit）、中斷行為
5. 建立 Skill Relevance Engine → [`runtime/runtime.db`](../runtime/runtime.db) — 3 維度 scoring（trigger_match 0.5 + domain_match 0.3 + weight 0.2）、threshold 0.5、conflict penalty ×0.5

### Step 8：建立 Feedback Promotion Pipeline（已實作）

將 feedback lesson 從 `skills/*/feedback_history/` 的原始觀察，透過機器可讀的 scoring、workflow 與 lifecycle automation，推進到 `workflow/`、`intelligence/`、`enforcement/`、`memory/` 或 runtime surfaces：

1. 建立 Promotion Pipeline 概覽 → [`feedback/pipeline/README.md`](../feedback/pipeline/README.md) — pipeline 架構圖、與既有層的關係、使用方式
2. 建立 Promotion Engine → [`feedback/pipeline/promotion-engine.yaml`](../feedback/pipeline/promotion-engine.yaml) — 5 維度 scoring（impact 0.30 + maturity 0.25 + frequency 0.20 + freshness 0.15 + urgency 0.10）、threshold 0.7 immediate / 0.5 backlog、5 種 promotion target decisions、3 個 scoring examples
3. 建立 Promotion Workflow → [`feedback/pipeline/promotion-workflow.yaml`](../feedback/pipeline/promotion-workflow.yaml) — 5 階段 workflow（assess → prepare → write → update → validate）、每階段有 entry/exit conditions、steps、output
4. 建立 Lifecycle Automation → [`feedback/pipeline/lifecycle-automation.yaml`](../feedback/pipeline/lifecycle-automation.yaml) — 4 種 automation（auto-archive cold 180 days、auto-downgrade stale 90 days、periodic promotion check weekly、cold data threshold monitor）、完整 state machine（new → experimental → candidate → validated → promoted → archived）

### Step 9：建立 Semantic Retrieval（長期）⏳ 待實作

1. 深化 SQLite / FTS runtime index
2. 加入 embedding-based retrieval
3. 建立 query ranking 與 context-aware scoring

---

## 預期效益

| 項目 | 目前 | 優化後 | 改善幅度 |
| --- | --- | --- | --- |
| Bootstrap token 消耗 | ~5000 tokens | ~800 tokens | **-84%** |
| 每次 session 初始載入 | 12 條規則全文 | 3 條核心規則 | **-75%** |
| Skill 發現成本 | 掃描整個 skills/ 目錄 | 查 skills-index.yaml | **-90%** |
| 知識載入決策 | 無成本資訊 | cost-aware routing | **可控** |
| Context 生命週期 | 永久保留 | TTL 管理 | **可預測** |
| Max mode 爆 token 風險 | 高 | 低（分階段載入） | **顯著降低** |
| Token 用量控制 | 無預算管理 | Token Budget + Health Score | **可預測** |
| 工具輸出浪費 | 原始輸出全量載入 | 4-level compression | **-50~95%** |
| 記憶管理 | 單一 memory 層 | 3 子層（working/summary/decision） | **精準控制** |
| 架構決策重複 | 無記錄 | ADR 系統 | **避免重複** |
| 失效模式辨識 | 被動發現 | Anti-pattern 目錄 | **主動預防** |

## 與現有架構的關係

本規劃完全相容於 [`next-stage-upgrade-plan.md`](../active/next-stage-upgrade-plan.md) 的 10 層架構：

| 本規劃元件 | 對應現有層 | 關係 |
| --- | --- | --- |
| Core Bootstrap | `enforcement/` | 縮小現有 Default Bootstrap |
| Skills Index | `knowledge/indexes/` | 擴充現有 routing table |
| Summary Layer | `knowledge/summaries/` | 擴充現有 summaries |
| Runtime Router | `runtime/`, `knowledge/runtime/` | 將設計文件轉為可執行 routing |
| Context Cost Metadata | `metadata/schema.md` | 擴充現有 schema |
| Context TTL | `runtime/context/` | 新增子層 |
| Token Budget System | `runtime/budget/` | 新增子層 |
| Context Health Score | `runtime/health/` | 新增子層 |
| Circuit Breaker / Guards | `runtime/README.md` | 新增子層 |
| Tool Metadata / Routing / Compression | `tools/` | 新增層 |
| Memory 子層 | `memory/working/`, `memory/summary/`, `memory/decision/` | 拆分現有 memory 層 |
| Decision System | `decisions/` | 新增層 |
| Anti-patterns | `anti-patterns/` | 新增層 |
| Skills Metadata v2 | `skills-index.yaml` | 升級現有 index |
| Semantic Retrieval | `knowledge/runtime/sqlite/` | 深化現有 prototype |

## 風險與緩解

| 風險 | 緩解 |
| --- | --- |
| Lazy-load 導致漏讀關鍵規則 | Core Bootstrap 保留最低必要規則；activation rules 明確定義觸發條件 |
| Summary 與 source 不同步 | Refresh policy 已定義 revalidate/downgrade 條件；validation scripts 可檢查 |
| TTL 過短導致重複載入 | TTL 可設定 conversation-level cache；summary 層輕量可頻繁載入 |
| 遷移期間相容性 | 舊 entrypoint 維持 active；新路徑逐步取代 |
| Agent 不遵守 activation rules | Activation rules 寫入 Core Bootstrap，確保每個 session 都讀到 |
| Token Budget hard stop 中斷工作 | Warning threshold 設 70%，預留緩衝空間；auto-archive 保留工作進度 |
| Circuit breaker 誤判 | 每個 guard 有獨立 threshold，可依模型/任務調整 |
| Compression 遺失關鍵資訊 | Raw level 保留完整資訊；summary/structured/minimal 僅用於非關鍵輸出 |
