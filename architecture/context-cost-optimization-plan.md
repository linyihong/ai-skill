# Context Cost Optimization Plan

本文件是 [`next-stage-upgrade-plan.md`](next-stage-upgrade-plan.md) 的**成本優化補充規劃**。它不取代原有架構分層規劃，而是針對目前最迫切的 **token 成本與 context 過重問題**，提出具體改造路徑。

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
README.md + shared-rules/README.md + Default Bootstrap
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

## Phase 1：立即省錢（第一優先）

### 1.1 Bootstrap 極小化

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
| `document-sizing` | 文件接近拆分門檻 | 文件變大時 |
| `document-todo-list` | 文件有未完成項目 | 修改文件時 |
| `goal-action-validation` | 重要工作單元需要驗證 | 執行關鍵步驟時 |
| `neutral-language` | 撰寫或審查文件用語 | 寫文件時 |
| `sanitization` | 撰寫 feedback lesson | 寫 lesson 時 |
| `authorization-scope` | 涉及授權邊界 | 分析外部系統時 |
| `cross-skill-references` | 引用其他 skill | 跨 skill 工作時 |
| `feedback-lessons` | 撰寫或 promotion lesson | 寫 lesson 時 |

**預期效益**：Token 減少 40-60%。

### 1.2 README 拆分

**現狀**：根 `README.md` 196 行，接近 Operating System Specification，每次 session 重讀。

**改造**：拆成四份文件：

| 文件 | 用途 | 目標行數 |
| --- | --- | --- |
| `README.md` | 超短入口，只留 layout 與 quickstart | ≤ 80 行 |
| `BOOTSTRAP.md` | Runtime minimal bootstrap（Core Bootstrap） | ≤ 50 行 |
| `ARCHITECTURE.md` | 深度設計與 roadmap | 不限 |
| `GOVERNANCE.md` | 管理規範與貢獻指南 | 不限 |

### 1.3 Rule Lazy-load 機制

**現狀**：`shared-rules/README.md` 的 Default Bootstrap 要求每次載入 12 條規則。

**改造**：在 `shared-rules/README.md` 中引入 **runtime activation model**：

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

### 1.4 Knowledge Summary Layer

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

## Phase 2：架構升級（第二優先）

### 2.1 Runtime Context Router

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
# runtime/router/activation-rules.yaml
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

### 2.2 Context Cost Metadata

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

### 2.3 Skill Index

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

### 2.4 Context TTL System

**現狀**：Context 一旦載入就永久留在 context graph，導致 Agent 越來越肥。

**改造**：建立 Context TTL 系統，定義每條 context 的生命週期：

```yaml
# runtime/context/ttl-policy.yaml
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

## Phase 3：真正 AI OS（第三優先）

### 3.1 Semantic Retrieval

建立語意檢索層，讓 agent 用自然語言查詢知識，而不是依賴目錄掃描：

- 整合 SQLite / FTS runtime index（已有 prototype）
- 加入 embedding-based retrieval（future）
- Query ranking 與 context-aware scoring

### 3.2 Episodic Memory

建立情節記憶層，讓 agent 記住「過去怎麼解決類似問題」：

- `memory/episodic/`：記錄過去 task 的關鍵決策與結果
- 每個 episodic record 包含：context、decision、outcome、token cost
- 支援 similarity-based retrieval

### 3.3 Runtime Orchestration

建立 runtime orchestration layer，自動管理 context loading、pruning 與 TTL：

- `runtime/orchestration/`：context lifecycle management
- 自動根據 task intent 選擇 model profile
- 自動 prune 過期 context

### 3.4 Multi-model Routing

根據 task 複雜度自動選擇模型：

| Task 類型 | 建議模型 | Context 策略 |
| --- | --- | --- |
| Bootstrap | small/fast | Core Bootstrap only |
| Routing | medium | Summary + index |
| Architecture reasoning | max | Full source + graph |
| Implementation | medium | Skill + workflow |
| Autocomplete | cheap | Index + checklist |

---

## 遷移路徑

### Step 1：建立 Core Bootstrap（立即）

1. 建立 [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) — 只含 rule-weight、dependency-reading、goal-ledger
2. 更新 `shared-rules/README.md` — 加入 runtime activation model
3. 更新根 `README.md` — 縮短為超短入口

### Step 2：建立 Skills Index（立即）

1. 建立 [`skills-index.yaml`](../skills-index.yaml)
2. 更新 `knowledge/indexes/README.md` 引用 skills index

### Step 3：擴充 Summary Layer（1-2 週）

1. 為每個 shared rule 建立 summary
2. 為每個 skill 建立 summary
3. 為每個 architecture 文件建立 summary
4. 在 summary 中加入 context_cost 與 estimated_tokens

### Step 4：建立 Runtime Router（2-3 週）

1. 建立 `runtime/router/activation-rules.yaml`
2. 建立 `runtime/router/cost-budget.yaml`
3. 建立 `runtime/context/ttl-policy.yaml`
4. 更新 `knowledge/runtime/routing-registry.yaml` 加入 cost metadata

### Step 5：全面導入 Context Cost Metadata（3-4 週）

1. 更新 `metadata/schema.md` 加入 context_cost、load_strategy、ttl 欄位
2. 為所有 Knowledge Atoms 補上 cost metadata
3. 更新 validation scripts 檢查 cost metadata

### Step 6：建立 Semantic Retrieval（長期）

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

## 與現有架構的關係

本規劃完全相容於 [`next-stage-upgrade-plan.md`](next-stage-upgrade-plan.md) 的 10 層架構：

| 本規劃元件 | 對應現有層 | 關係 |
| --- | --- | --- |
| Core Bootstrap | `shared-rules/` | 縮小現有 Default Bootstrap |
| Skills Index | `knowledge/indexes/` | 擴充現有 routing table |
| Summary Layer | `knowledge/summaries/` | 擴充現有 summaries |
| Runtime Router | `runtime/`, `knowledge/runtime/` | 將設計文件轉為可執行 routing |
| Context Cost Metadata | `metadata/schema.md` | 擴充現有 schema |
| Context TTL | `runtime/context/` | 新增子層 |
| Semantic Retrieval | `knowledge/runtime/sqlite/` | 深化現有 prototype |

## 風險與緩解

| 風險 | 緩解 |
| --- | --- |
| Lazy-load 導致漏讀關鍵規則 | Core Bootstrap 保留最低必要規則；activation rules 明確定義觸發條件 |
| Summary 與 source 不同步 | Refresh policy 已定義 revalidate/downgrade 條件；validation scripts 可檢查 |
| TTL 過短導致重複載入 | TTL 可設定 conversation-level cache；summary 層輕量可頻繁載入 |
| 遷移期間相容性 | 舊 entrypoint 維持 active；新路徑逐步取代 |
| Agent 不遵守 activation rules | Activation rules 寫入 Core Bootstrap，確保每個 session 都讀到 |
