# 失效模式 Patterns

本目錄存放跨 skill 可重用的 agent failure patterns。每個 pattern 記錄泛化後的 failure mode、trigger、required action、prevention gate 與 validation method。

當 [`failure-learning-system.md`](../failure-learning-system.md) 要求 promote 或查找 reusable failure pattern 時，先讀本索引。

| 模式 | 類別 | 狀態 | 摘要 |
| --- | --- | --- | --- |
| [Correction loop bypass](correction-loop-bypass.md) | `validation-gap` | validated | 防止 agent 在使用者指出修正不完整時，只修當下文字，卻漏掉 `.agent-goals`、failure learning、linked updates、validation、commit/push/readback。 |
| [Entrypoint positioning drift](entrypoint-positioning-drift.md) | `validation-gap` | validated | 防止 agent 在命名或架構變更後，只更新次要連結或段落，卻留下 root title、opening paragraph 或主要入口 framing 過期。 |
| [Shared-rules architecture drift](shared-rules-architecture-drift.md) | `dependency-miss` / `validation-gap` | validated | 防止 agent 在架構重構後，只更新主要檔案（workflow、intelligence、analysis）卻漏掉 enforcement/ 中的路徑參考同步。 |
| [Skill-local feedback bypass](skill-local-feedback-bypass.md) | `dependency-miss` / `validation-gap` | validated | 防止 agent 只補單一 skill 的 feedback lesson，卻沒有讀取全庫 failure-learning system 並沉澱 cross-skill prevention gate。 |
| [Source / mirror write drift](source-mirror-write-drift.md) | `source-mirror-drift` | validated | 防止 agent 更新 project-local tool mirrors 或 runtime copies，而不是 canonical source repo。 |
| [Tool config design without rule check](tool-config-design-without-rule-check.md) | `tool-strategy-gap` | candidate | 防止 agent 設計新工具配置時漏讀 `ai-tools/<tool>.md` 的現有規則，導致重複或邊界混淆。 |
| [Language preference drift](language-preference-drift.md) | `configuration-gap` / `instruction-conflict` | validated | 防止 agent 的 Custom Instructions 中設定了固定語言偏好，導致無視使用者實際使用的語言，強制用英文回應。 |
| [Failure-to-validator closure](failure-to-validator-closure.md) | `validation-gap` / `process-gap` | validated | 防止 agent 修復錯誤後，沒有把錯誤模式抽象化為可重複檢測的 validator 測試案例。 |
| [Framework duplication without interrogation](framework-duplication-without-interrogation.md) | `source-of-truth-duplication` / `requirements-cognition-gap` | validated | 防止 agent 修改 framework / runtime / governance 時，未先做需求拷問與 source-of-truth discovery，留下雙寫 rule、activation path、projection 或 generated surface。 |
| [Refactor parity feedback miss](refactor-parity-feedback-miss.md) | `validation-gap` / `dependency-miss` | candidate | 防止 agent 只在單一計畫補新舊功能對照，卻沒有把 replacement / refactor parity gate 回饋到 software-delivery workflow。 |
| [Commit/push before writeback transaction close](commit-before-validation-skip.md) | `validation-gap` | candidate | 防止 agent 在 commit/push 前跳過 Ai-skill writeback transaction 關閉條件（依 dependency-reading.md §Writeback Transaction Guard）。 |
| [Mandatory step blocker bypass](mandatory-step-blocker-bypass.md) | `process-gap` / `validation-gap` | validated | 防止 agent 遇到強制步驟的環境阻斷（如工具未安裝）時自行判斷「環境限制」而靜默跳過，應立即停止並通知使用者。 |
| [Knowledge-update-flow bypassed by sub-pipeline](knowledge-update-flow-bypassed-by-sub-pipeline.md) | `process-gap` / `validation-gap` | validated | 防止 agent 在新增 intelligence atom / failure pattern / scenario 時，只讀 sub-pipeline 文件（如 intelligence-extraction-pipeline）並把它當完整流程，跳過 master `knowledge-update-flow.md` 的 Step 4/7/9/11。 |
| [Premature ADR promotion](premature-adr-promotion.md) | `process-gap` / `governance-drift` | validated | 防止 agent 在 plan completed 前就建立 proposed/draft ADR，違反 constitution/ 只放 accepted ADRs 的定位；架構提案階段應在 `plans/active/<plan>.md` §Decision Rationale 完成，依新規則（2026-05-22）只在 plan completed 且通過 ADR Promotion Criteria 後才建立 accepted ADR。 |
| [Analysis domain discovery gap](analysis-domain-discovery-gap.md) | `discovery-gap` / `routing-miss` | validated | 防止 agent 分析外部 library / 工具時只放入 `intelligence/`，忽略該工具可能代表全新分析領域；應先檢查 `analysis/` 是否有對應入口。Discovery checkpoint `search_sources` 需含 `analysis/`。 |
| [Template drift](template-drift.md) | `template-inconsistency` / `governance-drift` | validated | 防止 agent 建立 reusable layer 文件時憑直覺寫 section header，造成 canonical template 漂移；應先 list_files + Read 同類既有檔，依 canonical section list 寫入。對應 scenario：`failure-pattern-template-consistency-v1`。 |
| [Skill classification boundary confusion](skill-classification-boundary-confusion.md) | `scope-drift` | candidate | 防止 agent 在處理多 skill feedback lessons 時，把 A skill 的 analysis technique 放到 B skill 的 feedback_history；放置位置由 skill scope 決定，不是 lesson 技術主題。 |
| [Cognitive mode resolution bypass](cognitive-mode-resolution-bypass.md) | `process-gap` / `governance-drift` | candidate | 防止 agent 跳過 cognitive mode 解析或 final response Cognitive close-out，導致壓縮策略錯誤、governance gate 未激活、memory isolation 失效與 chat/session 報告漏證據。 |
| [Inflated cognitive mode reporting](inflated-cognitive-mode-reporting.md) | `governance-drift` / `validation-gap` | validated | 防止 agent 膨脹或自造 Cognitive Contract mode / cost / activation signal，讓報告失去 capability semantics。 |
| [Bootstrap bypass on resume](bootstrap-bypass-on-resume.md) | `process-gap` / `governance-drift` | validated | 防止 agent 從 conversation summary 被喚起時，把 "Resume directly" 對話 framing 當成 runtime/governance bootstrap 豁免，跳過 CORE_BOOTSTRAP.md / runtime.db 查詢與 Bootstrap Receipt 輸出。 |
| [CLI doc drift](cli-doc-drift.md) | `source-of-truth-duplication` / `governance-drift` | validated | 防止 Go 改了 CLI subcommand / hook handler 但 `scripts/ai-skill-cli/docs/command-contract.md` 沒同步。對應 validator: `validateCLIDocSync`；canonical rule: `runtime/cli-modification-policy.yaml`。 |
| [Runtime YAML unprojected](runtime-yaml-unprojected.md) | `source-of-truth-duplication` / `governance-drift` | validated | 防止 `runtime/*.yaml` 沒設 `runtime_projection.enabled: true` 或缺 `target_key`，compiler silent skip 導致規則寫了不生效。對應 validator: `validateRuntimeYamlProjects`；plan 例外路徑: §Deferred Runtime Projection。 |
| [Markdown / YAML sync drift](markdown-yaml-sync-drift.md) | `source-of-truth-duplication` / `governance-drift` | validated | 防止改 canonical doc markdown 但沒同步改 sibling YAML companion。對應 validator: `validateMarkdownYamlSync`（sibling-pair only；cross-path mapping 列入 Phase 7 backlog）。 |
| [Bootstrap YAML bypass](bootstrap-yaml-bypass.md) | `governance-drift` / `source-of-truth-duplication` | validated | 防止 agent / hook 跳過 `generated_surfaces[runtime.core_bootstrap.contract]` query 直接讀 `CORE_BOOTSTRAP.md` prose，導致 obligation 列表落後 YAML schema 變更。對應 prevention: Phase 6 per-obligation dispatcher refactor + enhanced Bootstrap Receipt（含 `Active per-turn obligations:` 行）。 |
| [Intelligence layer bypass via tool adapter](intelligence-layer-bypass-via-tool-adapter.md) | `process-gap` / `knowledge-routing-miss` | validated | 防止 agent 因任務主題「關於某工具」而把跨工具可重用的設計洞見直接寫進 `ai-tools/<tool>.md`（P3），繞過 knowledge-update-flow Step 1 觸發，intelligence 層從未建立。Tool adapter 應只含工具專屬細節，設計原理必須先在 intelligence atom 再引用。 |
| [Shell script added without Go migration](shell-script-added-without-go-migration.md) | `process-gap` / `platform-governance-miss` | validated | 防止 agent 新增 `.sh` 腳本而非實作 Go CLI，違反跨平台治理政策。對應 validator: `validateNoNewShellScripts`；opt-out: `[skip-go-migration]`；canonical rule: `runtime/cli-modification-policy.yaml §gate.cli.no_new_shell_scripts`。 |
| [AI codegen passes CI / fails production](ai-codegen-passes-ci-fails-production.md) | `validation-gap` / `meta-tool-risk` | candidate | 防止 AI 生成程式碼通過 unit/integration test 卻在 production 失敗（外部觀察：43% 仍需 prod debug、88% 需 2–3 redeploy）。4 個 perf anti-pattern（loop DB query、unbounded collection、no timeout、string SQL）；對應 workflow: `workflow/software-delivery/perf-risk-gate.md`；對應 scenario: `validation/scenarios/software-delivery/ai-codegen-perf-risk-checklist.yaml`；opt-out: `[skip-perf-bench]`。 |
| [Rule without executor](rule-without-executor.md) | `meta-governance-gap` / `framework-self-audit-miss` | validated | Meta-pattern：Knowledge layer 寫了規則但 Runtime layer 沒對應 executor，規則只活在文件裡。2026-05-31 session 連續 5 個 instance 暴露（bootstrap bypass / workflow activation / capability discovery / sanitization / intelligence classification）。對應結構性修正：`enforcement/enforcement-registry.yaml` Layer 2.5 binding + Phase 3 compile-time lint（orphan_rule / missing_executor_symbol / behavioral_without_sunset）；regression scenario: `validation/scenarios/enforcement/2026-05-31-regression-five-instances-v1.yaml`。 |

## 維護

- 不要把 project-specific evidence 放進本目錄。
- 當 failure mode 可能跨 projects、tools、skills 或 agents 重演時，新增 pattern。
- 若 pattern 變成 skill-specific，把 lesson 移到該 skill 的 `feedback_history/`；只有 cross-skill trigger 仍有價值時，才從這裡連回。
- 若 pattern 變長，拆出獨立 examples，不要膨脹索引。

← [Back to enforcement index](../README.md)
