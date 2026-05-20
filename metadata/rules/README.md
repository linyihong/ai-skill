# Metadata Rules

`metadata/rules/` 是 metadata 操作規則索引，說明 `metadata/schema.md` 的欄位如何被 routing、governance、summaries，以及未來 graph / runtime surfaces 使用。

## 規則集合

| 規則 | 用途 |
| --- | --- |
| [`ranking/`](../ranking/README.md) | 多個 atom 或 source 都相關時，決定先讀哪一個。 |
| [`confidence/`](../confidence/README.md) | 描述證據強度，以及 atom 何時可從 candidate 移到 validated 或 stable。 |
| [`compatibility/`](../compatibility/README.md) | 新分層演進時，保留舊 skill entrypoints 與工具相容性。 |
| [`enforcement-rule-spec.md`](enforcement-rule-spec.md) | Enforcement Rule 專屬 metadata spec，繼承 Knowledge Atom schema 並新增 activation_conditions、always_apply、scope 等欄位。 |
| [`rule-weight.yaml`](rule-weight.yaml) | 規則權重與衝突優先序（Core Bootstrap，P0） |
| [`dependency-reading.yaml`](dependency-reading.yaml) | 依賴文件讀取鐵則（Core Bootstrap，P0） |
| [`conversation-goal-ledger.yaml`](conversation-goal-ledger.yaml) | 對話目標帳本（Core Bootstrap，P0） |
| [`linked-updates.yaml`](linked-updates.yaml) | 連動更新（Lazy-load，P1） |
| [`escalation-policy.yaml`](escalation-policy.yaml) | Escalation policy（Lazy-load，P1） |
| [`failure-learning-system.yaml`](failure-learning-system.yaml) | 失效學習系統（Lazy-load，P1） |
| [`sanitization.yaml`](sanitization.yaml) | 去敏與占位符（Lazy-load，P1） |
| [`authorization-scope.yaml`](authorization-scope.yaml) | 授權與範圍（Lazy-load，P0） |
| [`feedback-lessons.yaml`](feedback-lessons.yaml) | Feedback 與技巧條目（Lazy-load，P2） |
| [`goal-action-validation.yaml`](goal-action-validation.yaml) | 目標、執行、驗證共同流程（Lazy-load，P2） |
| [`tool-neutral-documentation.yaml`](tool-neutral-documentation.yaml) | 工具中立文件（Lazy-load，P2） |
| [`document-todo-list.yaml`](document-todo-list.yaml) | 文件 TODO 清單（Lazy-load，P2） |
| [`neutral-language.yaml`](neutral-language.yaml) | 中性與低爭議文件用語（Lazy-load，P2） |
| [`cross-skill-references.yaml`](cross-skill-references.yaml) | Cross-skill references（Lazy-load，P2） |
| [`reusable-guidance-boundary.yaml`](reusable-guidance-boundary.yaml) | 可重用規則與專案證據邊界（Lazy-load，P2） |
| [`content-layering.yaml`](content-layering.yaml) | 內容分層（Lazy-load，P2） |
| [`decision-efficiency.yaml`](decision-efficiency.yaml) | 決策效率（Lazy-load，P2） |
| [`prompt-cache-efficiency.yaml`](prompt-cache-efficiency.yaml) | Prompt cache efficiency（Lazy-load，P2） |
| [`conflict-matrix.yaml`](conflict-matrix.yaml) | Rule 衝突配對矩陣，定義已知衝突與解決方式（12 個配對） |

## 邊界

- Metadata rules 不覆蓋可執行 enforcement rules。
- Metadata 與 enforcement rules 衝突時，依 `enforcement/rule-weight.md`。
- Metadata 可以降低 context loading cost，但不能跳過 required dependency reading。
- 舊 `skills/` scaffold 已退役；active source-of-truth 由 promoted workflow / analysis / intelligence / enforcement source 承接。
- Metadata 文件正文預設使用繁體中文；英文保留給欄位名、enum、路徑、YAML key 與必要專有名詞。
