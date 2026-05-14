# Runtime Context Router

`runtime/router/` 負責 **context routing 決策邏輯**。Agent 在 session 啟動後，透過本層決定哪些知識需要載入、哪些可以 deferred。

## 路由流程

```
Session Start
  │
  ├─ 1. Load CORE_BOOTSTRAP.md（3 rules, ~800 tokens）
  │
  ├─ 2. Read README.md（OS layout, ~80 lines）
  │
  ├─ 3. Query skills-index.yaml → match task intent → find skill
  │
  ├─ 4. Run activation-engine.rb → load lazy rules if triggered
  │
  ├─ 5. Read knowledge summary（300-500 tokens）
  │
  ├─ 6. Expand to full source only if needed
  │
  └─ 7. Apply TTL policy → prune context at task/session boundary
```

## 路由決策表

| 階段 | 輸入 | 輸出 | Token 成本 |
| --- | --- | --- | --- |
| 1. Bootstrap | Session start | Core rules (3) | ~800 |
| 2. Layout | Task intent | OS map | ~200 |
| 3. Skill routing | Task intent keywords | Skill ID + entrypoint | ~100 |
| 4. Rule activation | Task context | Lazy rules list | ~200 |
| 5. Summary | Skill ID | 300-500 token summary | ~400 |
| 6. Full source | Summary match | Full document | variable |

## 路由檔案

| 檔案 | 用途 |
| --- | --- |
| [`activation-rules.yaml`](activation-rules.yaml) | 定義 lazy-load rules 的觸發條件與優先權（v2, validated） |
| [`activation-engine.rb`](activation-engine.rb) | 程式化 activation 判斷引擎，讀取 activation-rules.yaml 並輸出 activate list |
| [`activation-table.md`](activation-table.md) | 人類可讀的 Situation → Activated Rules 對照表 |
| [`cost-budget.yaml`](cost-budget.yaml)（未來） | Session token budget 管理 |

## Activation Engine 使用說明

[`activation-engine.rb`](activation-engine.rb) 是程式化的 activation 判斷工具，接受 task intent、file change、user signal 等輸入，輸出應該 activate 的 rule 列表。

### 基本用法

```bash
# 顯示所有規則的 activation 狀態（無輸入時顯示提示）
ruby runtime/router/activation-engine.rb

# 指定 task intent
ruby runtime/router/activation-engine.rb --intent migration

# 指定 file changes
ruby runtime/router/activation-engine.rb --file-changed enforcement/rule-weight.md

# 指定 user signal
ruby runtime/router/activation-engine.rb --signal 連動

# 複合條件
ruby runtime/router/activation-engine.rb --intent migration --file-changed "**/*.md" --file-changed enforcement/linked-updates.md
```

### 進階用法

```bash
# Dry-run 模式（顯示每條規則的判斷邏輯）
ruby runtime/router/activation-engine.rb --intent migration --dry-run

# 列出所有已知的 intent、signal、file pattern
ruby runtime/router/activation-engine.rb --list-known

# 多路線決策情境
ruby runtime/router/activation-engine.rb --routes 5 --signal 選擇

# 安全分析情境
ruby runtime/router/activation-engine.rb --intent security-analysis --signal authorization
```

### 支援的輸入類型

| 參數 | 說明 | 範例 |
|------|------|------|
| `--intent` / `-i` | Task intent（可多次指定） | `--intent migration --intent refactor` |
| `--signal` / `-s` | User signal（可多次指定） | `--signal 錯誤 --signal 失誤` |
| `--file-changed` / `-f` | 已修改的文件路徑（可多次指定） | `-f enforcement/rule-weight.md -f README.md` |
| `--has-todo` | 含 TODO 標記的文件（可多次指定） | `--has-todo README.md` |
| `--routes` / `-r` | 可行路線數量 | `--routes 5` |
| `--oversized` | 過大文件（可多次指定） | `--oversized large-file.md` |
| `--tool` / `-t` | 活躍工具名稱 | `--tool cursor` |
| `--dry-run` / `-n` | 顯示判斷邏輯但不輸出結果 | `--dry-run` |
| `--list-known` / `-l` | 列出所有已知條件值 | `--list-known` |

### Agent 整合建議

在 session 啟動流程的第 4 步，Agent 應執行：

```bash
ruby runtime/router/activation-engine.rb \
  --intent "$TASK_INTENT" \
  --signal "$USER_SIGNAL" \
  --file-changed "$CHANGED_FILES"
```

輸出中的 `Activated Lazy-load Rules` 列表即為需要載入的規則。

## 與既有層的關係

- `knowledge/runtime/routing-registry.yaml`：machine-readable routing records（atom → source → summary → cost）
- `knowledge/indexes/README.md`：human-readable task intent routing table
- `skills-index.yaml`：skill-level routing index（triggers → entrypoint → summary）
- `runtime/context/ttl-policy.yaml`：context 生命週期管理
- `metadata/rules/*.yaml`：每個 enforcement rule 的 metadata（含 activation_conditions、always_apply 等）
- `knowledge/graphs/rules/*.yaml`：rule dependency graph（activation 順序參考）
