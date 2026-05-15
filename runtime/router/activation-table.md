# Activation Table（Situation → Activated Rules）

本表提供人類可讀的 activation 對照，涵蓋所有常見情境。
程式化判斷請使用 `activation-engine.rb --dry-run`。

## 通用原則

| 載入策略 | 規則 | 觸發方式 |
|---------|------|---------|
| **Core Bootstrap**（永遠 preload） | rule-weight, dependency-reading, conversation-goal-ledger | Session 啟動時自動載入，~800 tokens |
| **Lazy-load**（依條件 activate） | 其餘 12 條規則 | 依下方情境表判斷 |

## 情境對照表

| # | 情境 | 觸發條件 | Activated Rules | 說明 |
|---|------|---------|----------------|------|
| 1 | **Repository 重構／搬遷** | task_intent: migration, refactor, rename, restructure | linked-updates, content-layering | 改目錄結構時需要連動更新 + 內容分層判斷 |
| 2 | **多文件修改** | file_change: `**/*.md` count>=2 | linked-updates | 同時改多個 .md 文件時檢查連動 |
| 3 | **Agent 重複錯誤** | user_signal: 失誤, 錯誤, failure, miss | failure-learning-system | 發現 agent 反覆犯同樣錯誤時啟動學習系統 |
| 4 | **Commit/Push 前驗證** | validation_gap: commit, push, sync, dirty | failure-learning-system | 閉環不完整或 dirty state 時檢查失效模式 |
| 5 | **Debug/Troubleshoot** | task_intent: debug, troubleshoot, fix-error | failure-learning-system | 除錯任務自動啟動失效學習 |
| 6 | **多路線決策** | task_complexity: routes>=3 | decision-efficiency | 超過 3 條可行路線時需要決策效率規則 |
| 7 | **路線選擇討論** | user_signal: 選擇, 路線, priority, 比較 | decision-efficiency | 使用者詢問「先做哪個」時啟動 |
| 8 | **撰寫文件** | task_intent: write-documentation, create-template, update-readme | tool-neutral-documentation, neutral-language | 寫文件時同時需要工具中立 + 用語規範 |
| 9 | **修改 enforcement 規則** | file_change: `enforcement/**` | tool-neutral-documentation, neutral-language | 改 enforcement 規則時檢查工具中立與用語 |
| 10 | **修改 skill 文件** | file_change: `skills/**/SKILL.md` | tool-neutral-documentation | skill 定義文件需要工具中立 |
| 11 | **更新/審閱文件** | task_intent: update-document, complete-document, review-document | document-todo-list | 文件操作時檢查 TODO 完整性 |
| 12 | **文件含 TODO** | file_has_todo: `**/*.md` | document-todo-list | 偵測到文件內有 TODO 標記時啟動 |
| 13 | **重要變更** | task_intent: critical-change, destructive-action, production-deploy, security-review | goal-action-validation | 高風險操作需要目標-執行-驗證流程 |
| 14 | **使用者要求驗證** | user_signal: 驗證, validate, 確認, confirm | goal-action-validation | 使用者明確要求驗證時啟動 |
| 15 | **寫 Feedback Lesson** | task_intent: write-feedback, create-lesson, write-feedback-history | sanitization, feedback-lessons | 寫 feedback 時同時需要去敏 + 條目規則 |
| 16 | **修改 feedback_history** | file_change: `**/feedback_history/**` | sanitization, feedback-lessons | 修改 feedback 目錄時啟動去敏與條目規則 |
| 17 | **安全分析** | task_intent: security-analysis, penetration-test, vulnerability-scan | authorization-scope | 安全相關任務需要授權範圍檢查 |
| 18 | **授權邊界討論** | user_signal: 授權, authorization, scope, 邊界 | authorization-scope | 使用者提及授權/範圍時啟動 |
| 19 | **跨 skill 整合** | task_intent: cross-skill, multi-skill, skill-integration | cross-skill-references | 跨 skill 工作時需要引用規範 |
| 20 | **多 skill 文件修改** | file_change: `skills/**` count>=2 | cross-skill-references | 同時改多個 skill 時檢查跨引用 |
| 21 | **Promote lesson** | task_intent: promote-lesson, review-feedback | feedback-lessons | 升級或審閱 feedback 條目時啟動 |
| 22 | **內容重組** | task_intent: content-organization | content-layering | 內容組織重構時需要分層規則 |
| 23 | **多 README 修改** | file_change: `**/README.md` count>=3 | content-layering | 同時改 3+ 個 README 時檢查內容分層 |
| 24 | **泛化為共用規則** | task_intent: promote-to-shared-rule | reusable-guidance-boundary | 將專案證據泛化為共用規則時需要邊界檢查 |
| 25 | **抽象化討論** | user_signal: 泛化, 可重用, reusable, 抽象化 | reusable-guidance-boundary | 使用者提及抽象化/可重用時啟動 |
| 26 | **翻譯文件** | task_intent: translate | neutral-language | 翻譯時需要中性用語規範 |
| 27 | **App/API/Embedded 開發** | task_intent: 開發, 實作, implement, coding, build, bug-fix, feature, code-review, design-review | workflow.software-delivery | 需要開發指引時載入 software-delivery workflow |
| 28 | **使用者要求開發** | user_signal: 幫我寫, 幫我做, 幫我開發, 幫我實作, 要怎麼做 | workflow.software-delivery | 使用者明確要求開發協助時載入 workflow |

## 複合情境範例

| 複合情境 | 組合條件 | Activated Rules |
|---------|---------|----------------|
| Repository 重構 + 多文件修改 | migration + file_change count>=2 | linked-updates, content-layering |
| Debug + 重複錯誤 | debug + user_signal: 錯誤 | failure-learning-system |
| 寫文件 + 修改 enforcement | write-documentation + file_change enforcement/** | tool-neutral-documentation, neutral-language |
| 安全分析 + 跨 skill | security-analysis + cross-skill | authorization-scope, cross-skill-references |
| 寫 feedback + 泛化 | write-feedback + promote-to-shared-rule | sanitization, feedback-lessons, reusable-guidance-boundary |

## 優先權參考

| Priority | 規則 | 意義 |
|---------|------|------|
| P0 | authorization-scope | 安全相關，必須立即載入 |
| P1 | linked-updates, failure-learning-system, sanitization | 高優先，情境符合時應載入 |
| P2 | 其餘 8 條規則 | 一般優先，情境符合時載入 |

## 驗證

- 每個情境至少對應一個 activation condition 類型
- 所有 lazy-load rules 都出現在至少一個情境中
- Core Bootstrap 3 條規則不在此表（永遠 preload）
