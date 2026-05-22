# Activation Table（Situation → Activated Rules）

本表提供人類可讀的 activation 對照，涵蓋所有常見情境。
程式化 enforcement activation 由 owner-layer executable YAML contracts 負責，投影到 `runtime/runtime.db generated_surfaces`。舊 `activation_rules` / `activation_rules_mirror` tables 已移除，不再維護 enforcement rule activation rows。

所有 enforcement rule 已有 owner-layer executable contract 時，agent 應先讀該 contract，再讀 companion Markdown；本表只保留人類導讀。YAML placement 的 canonical policy 見 [`../../governance/lifecycle/executable-contract-boundary.md`](../../governance/lifecycle/executable-contract-boundary.md)。

## 通用原則

| 載入策略 | 規則 | 觸發方式 |
|---------|------|---------|
| **Core Bootstrap**（永遠 preload） | rule-weight, dependency-reading, conversation-goal-ledger | Session 啟動時自動載入，~800 tokens |
| **Contract-backed activation** | `enforcement/*.yaml` executable contracts | 依各 contract 的 `activation` 欄位判斷，投影於 `generated_surfaces` |
| **Core bootstrap order** | `core_bootstrap_rules` | 只保存 core bootstrap 載入順序 |

## Contract-backed enforcement activations

以下規則已由 owner-layer executable YAML contract 接管 activation，不再維護 runtime lookup table，避免 runtime 與 enforcement contract 雙寫：

| Contract | Activation source |
| --- | --- |
| `authorization-scope` | [`../../enforcement/authorization-scope.yaml`](../../enforcement/authorization-scope.yaml) |
| `content-layering` | [`../../enforcement/content-layering.yaml`](../../enforcement/content-layering.yaml) |
| `cross-skill-references` | [`../../enforcement/cross-skill-references.yaml`](../../enforcement/cross-skill-references.yaml) |
| `decision-efficiency` | [`../../enforcement/decision-efficiency.yaml`](../../enforcement/decision-efficiency.yaml) |
| `document-todo-list` | [`../../enforcement/document-todo-list.yaml`](../../enforcement/document-todo-list.yaml) |
| `escalation-policy` | [`../../enforcement/escalation-policy.yaml`](../../enforcement/escalation-policy.yaml) |
| `evidence-hierarchy` | [`../../enforcement/evidence-hierarchy.yaml`](../../enforcement/evidence-hierarchy.yaml) |
| `failure-learning-system` | [`../../enforcement/failure-learning-system.yaml`](../../enforcement/failure-learning-system.yaml) |
| `feedback-lessons` | [`../../enforcement/feedback-lessons.yaml`](../../enforcement/feedback-lessons.yaml) |
| `goal-action-validation` | [`../../enforcement/goal-action-validation.yaml`](../../enforcement/goal-action-validation.yaml) |
| `linked-updates` | [`../../enforcement/linked-updates.yaml`](../../enforcement/linked-updates.yaml) |
| `neutral-language` | [`../../enforcement/neutral-language.yaml`](../../enforcement/neutral-language.yaml) |
| `prompt-cache-efficiency` | [`../../enforcement/prompt-cache-efficiency.yaml`](../../enforcement/prompt-cache-efficiency.yaml) |
| `reusable-guidance-boundary` | [`../../enforcement/reusable-guidance-boundary.yaml`](../../enforcement/reusable-guidance-boundary.yaml) |
| `sanitization` | [`../../enforcement/sanitization.yaml`](../../enforcement/sanitization.yaml) |
| `tool-neutral-documentation` | [`../../enforcement/tool-neutral-documentation.yaml`](../../enforcement/tool-neutral-documentation.yaml) |

## 任務察覺 → Workflow Discovery（必讀）

當 agent 察覺任務屬於 **需要 workflow 指揮** 的工作（命中下方 **#27** 或 registry 內任一 `route.workflow.*.activation_triggers`），**不要**直接寫碼／抓包／改 plan 就結束。依序：

1. **強制**執行 **[`governance/lifecycle/routing-philosophy.md`](../../governance/lifecycle/routing-philosophy.md)** Routing Pipeline（Step 1–5）：`task intent` → [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md) → [`routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) → `primary_source`。
2. **比對** registry 中所有 `route.workflow.*` 的 **`activation_triggers`**（task_intents、user_signals、file_change_globs）；列出命中者。
3. 若多條命中 → 讀 **[`workflow/workflow-routing.md`](../../workflow/workflow-routing.md)** §常見歧義 裁決；否則直接採唯一命中 route。
4. **進入**該 route 的 `README.md` → `execution-flow.md`，並載入該 route 的 **`required_dependencies`**（`dependency-reading`：必須讀 `primary_source`）。
5. 若該 route 含 **project_overlays**，在 **已進入 workflow 之後** 再載入 `<PROJECT_ROOT>` 專案 yaml。

| 常見任務 | registry `id`（觸發條件見各 route 的 `activation_triggers`） |
| --- | --- |
| SDK／plan／契約／實作／BDD | `route.workflow.software-delivery` |
| APK 逆向／Frida／協議分析 | `route.workflow.apk-analysis` |
| 純 agent 友善文件、零行為變更 | `route.workflow.documentation-ai-native` |
| 從零新專案 | `route.workflow.greenfield` |
| 旅遊行程 | `route.workflow.travel-planning` |

**新增 workflow 時**：只擴充 `routing-registry.yaml` 一筆 `route.workflow.*`（含 `activation_triggers`），**不必**在本表新增專向列。

## 情境對照表

| # | 情境 | 觸發條件 | Activated Rules | 說明 |
|---|------|---------|----------------|------|
| 1 | **Agent 重複錯誤** | user_signal: 失誤, 錯誤, failure, miss | `failure-learning-system.yaml` | 發現 agent 反覆犯同樣錯誤時啟動學習系統 |
| 2 | **Commit/Push 前驗證** | validation_gap: commit, push, sync, dirty | `failure-learning-system.yaml` | 閉環不完整或 dirty state 時檢查失效模式 |
| 3 | **Debug/Troubleshoot** | task_intent: debug, troubleshoot, fix-error | `failure-learning-system.yaml` | 除錯任務自動啟動失效學習 |
| 4 | **多路線決策** | task_complexity: routes>=3 | `decision-efficiency.yaml` | 超過 3 條可行路線時需要決策效率規則 |
| 5 | **路線選擇討論** | user_signal: 選擇, 路線, priority, 比較 | `decision-efficiency.yaml` | 使用者詢問「先做哪個」時啟動 |
| 6 | **撰寫文件** | task_intent: write-documentation, create-template, update-readme | `tool-neutral-documentation.yaml`, `neutral-language.yaml` | 寫文件時同時需要工具中立 + 用語規範 |
| 7 | **修改 enforcement 規則** | file_change: `enforcement/**` | `tool-neutral-documentation.yaml`, `neutral-language.yaml` | 改 enforcement 規則時檢查工具中立與用語 |
| 8 | **修改 workflow 文件** | file_change: `workflow/**/execution-flow.md` | `tool-neutral-documentation.yaml` | workflow 定義文件需要工具中立 |
| 9 | **更新/審閱文件** | task_intent: update-document, complete-document, review-document | `document-todo-list.yaml` | 文件操作時檢查 TODO 完整性 |
| 10 | **文件含 TODO** | file_has_todo: `**/*.md` | `document-todo-list.yaml` | 偵測到文件內有 TODO 標記時啟動 |
| 11 | **跨 workflow 整合** | task_intent: cross-workflow, multi-workflow, workflow-integration | `cross-skill-references.yaml` | 跨 workflow 工作時需要引用規範 |
| 12 | **多 workflow 文件修改** | file_change: `workflow/**` count>=2 | `cross-skill-references.yaml` | 同時改多個 workflow 時檢查跨引用 |
| 13 | **Promote lesson** | task_intent: promote-lesson, review-feedback | `feedback-lessons.yaml` | 升級或審閱 feedback 條目時啟動 |
| 14 | **內容重組** | task_intent: content-organization | `content-layering.yaml` | 內容組織重構時需要分層規則 |
| 15 | **多 README 修改** | file_change: `**/README.md` count>=3 | `content-layering.yaml` | 同時改 3+ 個 README 時檢查內容分層 |
| 16 | **泛化為共用規則** | task_intent: promote-to-enforcement-rule | `reusable-guidance-boundary.yaml` | 將專案證據泛化為共用規則時需要邊界檢查 |
| 17 | **抽象化討論** | user_signal: 泛化, 可重用, reusable, 抽象化 | `reusable-guidance-boundary.yaml` | 使用者提及抽象化/可重用時啟動 |
| 18 | **翻譯文件** | task_intent: translate | `neutral-language.yaml` | 翻譯時需要中性用語規範 |
| 19 | **Workflow 編排任務** | 任一 [`routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) 內 `route.workflow.*` 的 `activation_triggers` 命中；或 user_signal: workflow, 走 workflow；task_intent: workflow-orchestration | 見上方 §Workflow Discovery；各 route 的 `required_dependencies` 由 registry 定義 | **Blocking**：未完成 discovery 不得執行可觀察產品行為變更；具體 route 與 lazy rules **不在此表重複** |

> **Registry-first**：#27 是進入 workflow 世界的**通用閘門**。開發、APK、greenfield、documentation 等觸發條件與依賴文件寫在對應 `route.workflow.*`，不為每個 workflow 新增 activation 列。專案 overlay（如 `apk-analysis-sdk/runtime/workflow-activation.yaml`）僅在選定 **software-delivery** 後附加。

## 複合情境範例

| 複合情境 | 組合條件 | Activated Rules |
|---------|---------|----------------|
| Repository 重構 + 多文件修改 | migration + file_change count>=2 | `linked-updates.yaml` contract + `content-layering.yaml` |
| Debug + 重複錯誤 | debug + user_signal: 錯誤 | `failure-learning-system.yaml` |
| 寫文件 + 修改 enforcement | write-documentation + file_change enforcement/** | `tool-neutral-documentation.yaml`, `neutral-language.yaml` |
| 安全分析 + 跨 skill | security-analysis + cross-skill | `authorization-scope.yaml` contract + cross-skill-references |
| 寫 feedback + 泛化 | write-feedback + promote-to-enforcement-rule | `sanitization.yaml` contract + `feedback-lessons.yaml`, `reusable-guidance-boundary.yaml` |
| 改 SDK plan + 實作 | file_change docs/plans/** + task_intent implement | registry → `route.workflow.software-delivery`（含 docs-first、linked-updates） |
| Frida 抓短劇 API | task_intent frida + file_change TATA/scripts/frida/** | registry → `route.workflow.apk-analysis`（含 `authorization-scope.yaml` contract） |
| 多 route 同時命中 | 改 `*-sdk/**` 且開 Frida | §Workflow Discovery 步驟 3 + workflow-routing §歧義 → 通常 software-delivery vs apk-analysis |

## 優先權參考

| Priority | 規則 | 意義 |
|---------|------|------|
| P0 | `authorization-scope.yaml` contract | 安全相關，必須依 contract activation 立即載入 |
| P1 | `linked-updates.yaml`、`failure-learning-system.yaml`、`sanitization.yaml` | 高優先，情境符合時應載入 |
| P2 | 其餘 enforcement contracts | 一般優先，情境符合時載入 |

## 驗證

- 每個 enforcement activation 情境至少對應一個 contract 的 `activation` 欄位
- 舊 `activation_rules` / `activation_rules_mirror` tables 不再存在
- Contract-backed enforcement activation 不在 runtime lookup 維護第二份 rule body
- Core Bootstrap 3 條規則不在此表（永遠 preload）
