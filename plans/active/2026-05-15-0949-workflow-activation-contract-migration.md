# Workflow Activation Contract — Migration Plan

## 問題

目前的 [`activation-rules.yaml`](../../runtime/router/activation-rules.yaml) 同時包含兩種不同性質的規則：

| 類型 | 範例 | 該誰管 |
|------|------|--------|
| **Enforcement Rule** | `enforcement.linked-updates`, `enforcement.failure-learning-system` | Runtime — 跨 workflow 的基礎設施邏輯 |
| **Workflow Activation** | `workflow.software-delivery` | Workflow — 特定 workflow 的領域知識 |

這違反了 [`runtime/README.md`](../../runtime/README.md:17) 的邊界原則：「Skill-specific workflow 全文放到 `workflow/`」。

## 目標

1. Runtime 只提供「進入點時機」（entry point timing），不維護任何 workflow 的觸發詞彙表
2. 每個 workflow 自帶 `activation-contract.yaml`，自己決定要不要被觸發
3. 新增 workflow 不需修改 runtime

---

## 1. activation-contract.yaml Schema

每個 workflow 目錄下（如 `workflow/software-delivery/activation-contract.yaml`）放置此檔案：

```yaml
# Workflow Activation Contract
# 定義此 workflow 在什麼條件下應被觸發。
# Runtime 的 activation-engine.rb 會掃描所有 workflow/*/activation-contract.yaml
# 建立 Entry Point 索引，並在收到 intent event 時 broadcast 給符合條件的 workflow。

activation_contract_version: v1
workflow_id: software-delivery
description: App/API/Embedded 開發指引 workflow

triggers:
  task_intent:
    description: "當任務意圖包含以下關鍵字時觸發"
    keywords:
      - 開發
      - 實作
      - implement
      - 寫程式
      - coding
      - 建置
      - build
      - 修 bug
      - bug-fix
      - 新功能
      - feature
      - 開發指引
      - development-guidance
      - contract-first
      - bdd
      - code-review
      - design-review
      - release-review

  user_signal:
    description: "當使用者訊息包含以下關鍵字時觸發"
    keywords:
      - 開發
      - 寫 code
      - 實作
      - 要怎麼做
      - 幫我寫
      - 幫我做
      - 幫我開發
      - 幫我實作

  file_context:
    description: "當工作上下文包含以下檔案類型時觸發"
    patterns:
      - "**/*.swift"
      - "**/*.kt"
      - "**/*.java"
      - "**/*.dart"
      - "**/*.ts"
      - "**/*.tsx"
      - "**/*.js"
      - "**/*.jsx"
      - "**/*.py"
      - "**/*.rb"
      - "**/*.go"
      - "**/*.rs"
      - "**/pubspec.yaml"
      - "**/build.gradle*"
      - "**/Package.swift"
      - "**/Cargo.toml"
      - "**/go.mod"

load:
  strategy: lazy
  priority: P1
  estimated_tokens: 3500
  source: workflow/software-delivery/execution-flow.md
  summary: knowledge/summaries/development-guidance.md

# ── 以下為選擇性中繼資料 ──

routing:
  registry_ref: route.workflow.software-delivery
  required_dependencies:
    - workflow/software-delivery/artifact-gates.md
    - enforcement/README.md
  candidate_sources:
    - analysis/development-guidance/README.md
    - workflow/software-delivery/README.md
    - intelligence/engineering/development/README.md
    - knowledge/summaries/development-guidance.md
    - knowledge/graphs/workflow-software-delivery.yaml

compatibility:
  models:
    preferred: specialized
    fallback: default
  context_cost: 3500
```

### Schema 欄位說明

| 欄位 | 必要 | 說明 |
|------|------|------|
| `activation_contract_version` | 是 | Schema 版本，用於向後相容 |
| `workflow_id` | 是 | 對應 workflow 目錄名稱 |
| `description` | 是 | 人類可讀的說明 |
| `triggers.task_intent` | 否 | 任務意圖關鍵字觸發 |
| `triggers.user_signal` | 否 | 使用者訊息關鍵字觸發 |
| `triggers.file_context` | 否 | 工作上下文檔案類型觸發 |
| `load` | 是 | 載入策略（strategy, priority, tokens, source, summary） |
| `routing` | 否 | 對應 routing-registry.yaml 的參考 |
| `compatibility` | 否 | Model profile 與 context cost 資訊 |

---

## 2. activation-engine.rb 改造

### 現狀

```ruby
def load_activation_rules(path = ACTIVATION_RULES_PATH)
  # 只讀 activation-rules.yaml
  # 直接比對 task_intent / user_signal
end
```

### 目標架構

```ruby
# 新的 Entry Point 掃描機制
WORKFLOW_DIR = 'workflow/*/activation-contract.yaml'

def scan_workflow_entry_points
  # 掃描所有 workflow/*/activation-contract.yaml
  # 建立 entry_point_index = { workflow_id => contract }
  # 回傳索引供 broadcast 使用
end

def evaluate_rules(rules, inputs, entry_points)
  # 1. 先評估 enforcement rules（現有邏輯）
  # 2. 如果有 intent event，broadcast 給所有 entry points
  # 3. 每個 entry point 自己比對 triggers
  # 4. 回傳符合條件的 workflow
end
```

### 關鍵設計決策

| 決策 | 選擇 | 理由 |
|------|------|------|
| 掃描時機 | Session 啟動時掃描一次 + cache | 避免每次 event 都掃描磁碟 |
| Cache 失效 | 偵測 `workflow/*/activation-contract.yaml` 的 mtime 變化 | 支援 runtime 中新增 workflow |
| Broadcast 實作 | 順序比對（非平行），依 priority 排序 | 確保可預測的行為 |
| 衝突處理 | 多個 workflow 符合時，選 priority 最高的 | 若同 priority，依 workflow_id 字母序 |

---

## 3. activation-rules.yaml 精簡

### 移除

```yaml
# 整個「Workflow Activation Rules」區塊移除
- rule_id: workflow.software-delivery
  ...
```

### 保留

所有 `enforcement.*` 規則保持不變——它們是跨 workflow 的基礎設施邏輯，本來就該在 runtime。

---

## 4. activation-table.md 更新

Row 27-28 的內容改為：

```markdown
| 27 | **App/API/Embedded 開發** | task_intent 含開發相關詞彙 | workflow.software-delivery | 透過 activation-contract.yaml 觸發 |
| 28 | **使用者要求開發** | user_signal 含開發相關詞彙 | workflow.software-delivery | 透過 activation-contract.yaml 觸發 |
```

並在表格下方新增說明：

> Row 27-28 的觸發條件定義在 `workflow/software-delivery/activation-contract.yaml`，而非 `activation-rules.yaml`。Runtime 只負責掃描 entry point 並 broadcast event。

---

## 5. routing-registry.yaml 精簡

[`route.workflow.software-delivery`](../../knowledge/runtime/routing-registry.yaml:592) 的 `triggers` 欄位（30+ 關鍵字）應移除，改為引用 `activation-contract.yaml`：

```yaml
- id: route.workflow.software-delivery
  task_intent: 產出 app/API/embedded development guidance
  # triggers 欄位移除，改為：
  activation_contract: workflow/software-delivery/activation-contract.yaml
  primary_source: workflow/software-delivery/execution-flow.md
  ...
```

Registry 只保留 route metadata（priority, confidence, context_cost, model profile），觸發邏輯完全交給 `activation-contract.yaml`。

---

## 6. Validation Scenarios 更新

### 新增 scenario

`validation/scenarios/app-dev/activation-contract-scan-v1.yaml` — 測試 runtime 是否能正確掃描 `workflow/*/activation-contract.yaml` 並建立 entry point 索引。

### 修改現有 scenario

- [`workflow-activation-v1.yaml`](../../validation/scenarios/app-dev/workflow-activation-v1.yaml) — 更新 `expected_route.steps` 改為描述 broadcast + self-matching 流程
- [`compiler-output-v1.yaml`](../../validation/scenarios/app-dev/compiler-output-v1.yaml) — 新增 `activation-contract.yaml` 到 compiler 的 source_target_mapping（選擇性）

---

## 7. Migration 步驟

### Phase 1：建立 activation-contract.yaml（P0）

1. 在 `workflow/software-delivery/` 建立 `activation-contract.yaml`
2. 內容包含目前 `activation-rules.yaml` 中 `workflow.software-delivery` 的所有觸發條件
3. 加上 `file_context` 觸發條件（程式碼檔案類型）

### Phase 2：改造 activation-engine.rb（P0）

1. 新增 `scan_workflow_entry_points` 方法
2. 修改 `evaluate_rules` 支援 broadcast 模式
3. 新增 `--scan-entry-points` CLI 參數
4. 維持向後相容：如果 `activation-contract.yaml` 不存在，fallback 到 `activation-rules.yaml`

### Phase 3：精簡 activation-rules.yaml（P1）

1. 移除 `workflow.software-delivery` rule
2. 確認所有 enforcement rules 不受影響

### Phase 4：精簡 routing-registry.yaml（P1）

1. 移除 `route.workflow.software-delivery` 的 `triggers` 欄位
2. 改為 `activation_contract` 引用

### Phase 5：更新 activation-table.md（P1）

1. Row 27-28 改為標註「透過 activation-contract.yaml 觸發」
2. 新增說明段落

### Phase 6：更新 Validation Scenarios（P1）

1. 新增 `activation-contract-scan-v1.yaml`
2. 修改 `workflow-activation-v1.yaml`

### Phase 7：推廣到其他 Workflow（P2）

1. 為 `workflow/apk-analysis/` 建立 `activation-contract.yaml`
2. 為 `workflow/travel-planning/` 建立 `activation-contract.yaml`
3. 為 `workflow/documentation/` 建立 `activation-contract.yaml`
4. 為 `workflow/repo-analysis/` 建立 `activation-contract.yaml`

---

## 8. 受影響檔案清單

| 檔案 | 變更類型 | Phase |
|------|---------|-------|
| `workflow/software-delivery/activation-contract.yaml` | **新增** | P0 |
| `runtime/router/activation-engine.rb` | **修改** — 新增 scan + broadcast | P0 |
| `runtime/router/activation-rules.yaml` | **修改** — 移除 workflow.software-delivery | P1 |
| `knowledge/runtime/routing-registry.yaml` | **修改** — triggers → activation_contract | P1 |
| `runtime/router/activation-table.md` | **修改** — Row 27-28 說明 | P1 |
| `validation/scenarios/app-dev/activation-contract-scan-v1.yaml` | **新增** | P1 |
| `validation/scenarios/app-dev/workflow-activation-v1.yaml` | **修改** — 更新流程描述 | P1 |
| `workflow/apk-analysis/activation-contract.yaml` | **新增** | P2 |
| `workflow/travel-planning/activation-contract.yaml` | **新增** | P2 |
| `workflow/documentation/activation-contract.yaml` | **新增** | P2 |
| `workflow/repo-analysis/activation-contract.yaml` | **新增** | P2 |

---

## 9. 風險與注意事項

### 向後相容

- Phase 1-2 期間，`activation-rules.yaml` 的 `workflow.software-delivery` 仍然存在
- `activation-engine.rb` 應先檢查 `activation-contract.yaml`，若不存在則 fallback 到 `activation-rules.yaml`
- 這樣可以在不中斷現有功能的情況下逐步遷移

### Compiler 整合

- 可考慮將 `activation-contract.yaml` 加入 [`compiler-rules.yaml`](../../runtime/compiler/compiler-rules.yaml) 的 `source_target_mapping`
- 但 activation contract 本身已經是 YAML，不需要編譯——它直接是 runtime 可讀格式
- 如果需要 generated 版本，可以加一個 pass-through mapping

### 與 Pre-commit Hook 的關係

- 修改 `activation-contract.yaml` 不影響現有的 compiler sync check
- 但應在 pre-commit hook 中新增檢查：如果 `activation-contract.yaml` 存在，對應的 `activation-rules.yaml` 不應包含該 workflow 的 rule
