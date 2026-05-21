# Enforcement Layer Enhancement（enforcement/ 後續強化計畫）

Status: `completed`
Created: 2026-05-14
Depends on: [`shared-rules-to-enforcement-migration.md`](shared-rules-to-enforcement-migration.md)（搬遷完成後才能開始）
Related: [`metadata/schema.md`](../../metadata/schema.md), [`runtime/runtime.db`](../../runtime/runtime.db), [`enforcement/rule-weight.md`](../../enforcement/rule-weight.md), [`governance/lifecycle/README.md`](../../governance/lifecycle/README.md)

---

## 背景

在 [`enforcement/ → enforcement/` 搬遷](shared-rules-to-enforcement-migration.md)完成後，enforcement layer 需要進一步強化以發揮完整潛力。本計畫涵蓋 5 個強化方向，對應外部建議中「尚未做到」的部分。

---

## 方向 A：Enforcement Rule Metadata Spec

### 目標

建立 enforcement rule 專屬的 metadata spec，讓每條 rule 有 machine-readable 的結構化描述。

### 現狀對照

| 現有資產 | 缺口 |
|---------|------|
| `metadata/schema.md` 有通用 Knowledge Atom schema | 無 enforcement-specific spec |
| schema 涵蓋 8/9 建議欄位 | 無 `owner` 欄位、無 `activation_conditions` 專屬欄位 |
| schema 是 YAML template | 無實際的 rule metadata 實例 |

### 實作項目

- [ ] 建立 `metadata/rules/enforcement-rule-spec.md`：
  - 定義 Enforcement Rule 專屬的 metadata 欄位（繼承 Knowledge Atom schema 的必填欄位）
  - 新增欄位：`activation_conditions`（結構化觸發條件）、`always_apply`（boolean）、`scope`（適用範圍列舉）、`deprecated_by`（被哪條 rule 取代）
  - 提供 YAML template 與範例
- [ ] 為每條 enforcement rule 建立對應的 `.yaml` metadata 檔案：
  - `metadata/rules/rule-weight.yaml`
  - `metadata/rules/dependency-reading.yaml`
  - `metadata/rules/conversation-goal-ledger.yaml`
  - `metadata/rules/linked-updates.yaml`
  - `metadata/rules/failure-learning-system.yaml`
  - `metadata/rules/sanitization.yaml`
  - `metadata/rules/feedback-lessons.yaml`
  - `metadata/rules/goal-action-validation.yaml`
  - `metadata/rules/authorization-scope.yaml`
  - `metadata/rules/tool-neutral-documentation.yaml`
  - `metadata/rules/document-sizing.yaml`
  - `metadata/rules/document-todo-list.yaml`
  - `metadata/rules/neutral-language.yaml`
  - `metadata/rules/cross-skill-references.md`
  - `metadata/rules/reusable-guidance-boundary.md`
  - `metadata/rules/content-layering.md`
  - `metadata/rules/decision-efficiency.md`
- [ ] 更新 `metadata/rules/README.md` 加入 enforcement rule spec 的索引
- [ ] 更新 `metadata/README.md` 加入 enforcement rule metadata 的說明
- [ ] 驗證：每個 `.yaml` 檔案的 `source_path` 指向 `enforcement/` 下的對應檔案

---

## 方向 B：Enforcement Rule Dependency Graph

### 目標

將 `linked-updates.md` 中的文字依賴關係轉換為 machine-readable graph，存放在 `knowledge/graphs/rules/`。

### 現狀對照

| 現有資產 | 缺口 |
|---------|------|
| `linked-updates.md` 有完整的文字依賴表（51 行） | 無 machine-readable graph |
| `knowledge/graphs/` 有 graph infrastructure（edge types、record 格式） | 無 `knowledge/graphs/rules/` 目錄 |
| 15 個既有 graph records | 無任何 rule dependency graph |

### 實作項目

- [ ] 建立 `knowledge/graphs/rules/` 目錄
- [ ] 建立 `knowledge/graphs/rules/README.md`：
  - 說明此目錄存放 enforcement rule 之間的依賴關係
  - 引用 `linked-updates.md` 作為 source of truth
  - 列出所有 rule graph records
- [ ] 從 `linked-updates.md` 的「常見連動關係」表格提取依賴關係，建立以下 graph records：
  - `knowledge/graphs/rules/core-bootstrap.yaml`（rule-weight、dependency-reading、conversation-goal-ledger 之間的關係）
  - `knowledge/graphs/rules/linked-updates-deps.yaml`（linked-updates 依賴的所有 rules）
  - `knowledge/graphs/rules/failure-learning-deps.yaml`（failure-learning-system 依賴的所有 rules）
  - `knowledge/graphs/rules/content-layering-deps.yaml`（content-layering 被哪些 rules 依賴）
  - `knowledge/graphs/rules/full-rule-graph.yaml`（完整的 rule dependency graph，含所有 17 條 rules）
- [ ] 更新 `knowledge/graphs/README.md` 加入 rules/ 子目錄的索引
- [ ] 驗證：graph 中的每個 edge 都能對應到 `linked-updates.md` 中的一行

---

## 方向 C：Runtime Activation Engine

### 目標

讓 `runtime/runtime.db` 從「人工參考文件」變成「可被 runtime 實際使用的 activation engine」。

### 現狀對照

| 現有資產 | 缺口 |
|---------|------|
| `runtime/runtime.db` 已完成 14 條規則的結構化定義 | 未被任何 engine 讀取 |
| 每條 rule 有 `activation.when` 條件 | activation 條件是人工比對 |
| 有 `load.strategy`、`load.priority`、`load.estimated_tokens` | 無程式化的 activation 判斷邏輯 |

### 實作項目

- [ ] 建立 `runtime/router/activation-engine.rb`（或 shell script）：
  - 讀取 `activation-rules.yaml`
  - 接受 task intent、file change、user signal 等輸入
  - 輸出應該 activate 的 rule 列表
  - 支援 `--dry-run` 模式顯示判斷結果
- [ ] 建立 `runtime/router/activation-table.md`（人類可讀的 activation 對照表）：
  - 格式：Situation → Activated Rules → 說明
  - 涵蓋所有常見情境（repo refactor、reusable extraction、failure repeated、close-loop 等）
- [ ] 更新 `activation-rules.yaml`：
  - 狀態從 `candidate` 改為 `validated`
  - 所有 `source` 路徑從 `enforcement/` 改為 `enforcement/`
  - 加入 `rule_id` 與 `metadata/rules/` 中 `.yaml` 檔案的對應
- [ ] 更新 `runtime/README.md` 加入 activation engine 的使用說明
- [ ] 驗證：activation engine 能正確判斷 5 個以上測試情境

---

## 方向 D：Rule Conflict Resolution Enhancement

### 目標

將 `rule-weight.md` 的文字衝突規則升級為 machine-readable conflict matrix。

### 現狀對照

| 現有資產 | 缺口 |
|---------|------|
| `rule-weight.md` 有完整框架（6 條規則 + 8 個範例） | 無 machine-readable conflict matrix |
| P0/P1/P2/P3 權重體系已定義 | 無自動化衝突檢測 |
| 不確定時的處理流程 | 無程式化的 resolution engine |

### 實作項目

- [ ] 建立 `metadata/rules/conflict-matrix.yaml`：
  - 定義已知的 rule 衝突配對
  - 每個衝突配對包含：`rule_a`、`rule_b`、`resolution`（哪條勝出）、`rationale`（原因）、`example`（範例連結）
  - 從 `rule-weight.md` 的 8 個範例中提取衝突配對
- [ ] 更新 `rule-weight.md`（搬遷後為 `enforcement/rule-weight.md`）：
  - 在「常見範例」表格中加入「對應 conflict matrix」欄位
  - 加入 conflict matrix 的連結
- [ ] 更新 `metadata/rules/README.md` 加入 conflict matrix 的索引
- [ ] 驗證：conflict matrix 中的每個配對都能對應到 `rule-weight.md` 中的一個範例

---

## 方向 E：Enforcement Rule Deprecation Lifecycle

### 目標

建立 enforcement rule 專屬的 deprecation 流程與儲存位置。

### 現狀對照

| 現有資產 | 缺口 |
|---------|------|
| `governance/lifecycle/README.md` 有通用 lifecycle 與 deprecation 機制 | 無 enforcement rule 專屬的 deprecation 流程 |
| Skills Deprecation Timeline（Phase A/B/C/D） | 無 rule 版本的追蹤機制 |
| 刪除規則與檢查清單 | 無 `enforcement/deprecated/` 目錄 |

### 實作項目

- [ ] 建立 `enforcement/deprecated/` 目錄（在搬遷完成後）
- [ ] 建立 `enforcement/deprecated/README.md`：
  - 說明此目錄存放已 deprecated 的 enforcement rules
  - 定義 deprecation 流程：標記 → 公告 → 緩衝期 → 搬移
  - 引用 `governance/lifecycle/README.md` 的 lifecycle 規則
- [ ] 在 `metadata/schema.md` 中確認 `status: deprecated` 的使用方式：
  - 被 deprecated 的 rule 在 metadata 中標記 `status: deprecated`
  - 加入 `replaces` 欄位指向取代它的新 rule
- [ ] 更新 `governance/lifecycle/README.md`：
  - 在 lifecycle states 中加入 enforcement rule 專屬的 deprecation 說明
  - 加入 rule 版本的追蹤建議
- [ ] 驗證：deprecation 流程文件涵蓋所有必要步驟

---

## 執行順序

```
Phase 0: enforcement/ → enforcement/ 搬遷完成（必要前置條件）
    │
    ├── 方向 A (Enforcement Metadata Spec)
    │   └── 方向 B (Rule Dependency Graph) — 依賴 A 的 rule_id
    │       └── 方向 C (Activation Engine) — 依賴 A 的 metadata + B 的 graph
    │
    ├── 方向 D (Conflict Resolution) — 可與 A/B/C 並行
    │
    └── 方向 E (Deprecation Lifecycle) — 可最後做，但 infrastructure 先建立
```

### 建議優先順序

| 優先級 | 方向 | 原因 |
|--------|------|------|
| 🔴 P0 | **方向 A**（Metadata Spec） | 其他方向的基礎；沒有 rule_id 就無法建立 graph、activation、conflict matrix |
| 🔴 P1 | **方向 C**（Activation Engine） | 最有實際價值的產出；讓 activation-rules.yaml 真正可用 |
| 🟡 P2 | **方向 B**（Rule Graph） | 依賴 A 的 rule_id；可批次產生 |
| 🟡 P2 | **方向 D**（Conflict Matrix） | 可與 A 並行，但需要 rule-weight.md 的範例 |
| 🟢 P3 | **方向 E**（Deprecation） | 目前尚無 deprecated rule，可延後但 infrastructure 先建立 |

---

## 與既有文件的關係

- [`shared-rules-to-enforcement-migration.md`](shared-rules-to-enforcement-migration.md) — 本計畫的前置條件
- [`metadata/schema.md`](../../metadata/schema.md) — 方向 A 的基礎 schema
- [`runtime/runtime.db`](../../runtime/runtime.db) — 方向 C 的基礎資料
- [`enforcement/rule-weight.md`](../../enforcement/rule-weight.md) — 方向 D 的基礎規則
- [`governance/lifecycle/README.md`](../../governance/lifecycle/README.md) — 方向 E 的基礎 lifecycle
- [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) — 方向 B 的 source of truth
