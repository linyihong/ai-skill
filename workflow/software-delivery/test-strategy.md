# Test Strategy Slice（BDD Closure / Test Strategy / Test-First Ordering）

> **Cognitive Slice**：`sd-test-strategy`（從 [`execution-flow.md`](execution-flow.md) §2 + §4 子節 與 [`development-process.md`](development-process.md) §BDD Execution Closure / §Test Strategy Gate 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-test-strategy` |
| `purpose` | 定義測試策略：BDD 狀態閉環、docs-first BDD ordering、保護舊行為 vs 驗證新 code 的分流、mutation testing 有效性檢查、framework/runtime/governance 升級的 test-first 強制順序 |
| `type` | `execution` |
| `tags` | artifact-gate, test, bdd |
| `load_when` | 定義測試策略 / BDD 閉環 / test-first ordering |
| `do_not_load_when` | 不涉測試設計的純文件改動、evidence-only / 純分析任務 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「哪種風險用哪種測試、什麼順序寫、何時 test-first 為強制」的 ordering / gate；通過 workflow membership test，不承載 evidence 取得方法（非 analysis），不論證長期模式（非 intelligence） |
| `canonical_source` | 本檔（原 `execution-flow.md` §2 Docs-First BDD Closure Loop + §4 子節「測試策略定義」+「Test-First Ordering」+ `development-process.md` §BDD Execution Closure + §Test Strategy Gate 含 Mutation Testing） |
| `dependencies` | `sd-intake`（先 requirements / acceptance criteria）、`sd-contracts`（測試針對 contract 設計）、`sd-validation`（perf 執行關卡，本 slice cross-link 不複製） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Phase 4 Scenario A（execution-only：新需求的 BDD/test 順序）、Scenario C（mixed：debug 引用本 slice） |

> **Perf 內容邊界（與 sd-validation 的分工）**：本 slice 引用 perf 風險作為「測試策略選型考量」（風險→測試類型對應、新需求順序中的 smoke），但**不複製** perf 測試類型表與最低指標。perf **執行關卡 / 觸發條件 / 最低指標 / 完整測試類型表** canonical 在 [`validation.md`](validation.md)（`sd-validation`），請直接參照。

## 1. BDD Execution Closure（BDD 執行閉環）

Narrative BDD 在回填期間是可接受的，但不能被當作完成的測試覆蓋率。對每個關鍵 scenario，記錄以下狀態之一：

| 狀態 | 意義 | 必要下一步 |
| --- | --- | --- |
| `automated` | Scenario 由 unit、contract、API、integration、E2E、fixture 或 runner test 覆蓋 | 連結 test path/name |
| `fixture-backed` | Scenario 由 checked-in input/output fixtures 證明但沒有完整 runner | 連結 fixture 和 assertion 負責人 |
| `manual-evidence` | Scenario 需要手動、UI、bench 或外部服務證據 | 記錄執行步驟、證據和限制 |
| `pending-runner` | Gherkin 存在但沒有 runner/step definition 被接上 | 加上 runner 選擇或對應到可執行的 test type |
| `not-automatable` | Tooling 無法直接強制執行 | 說明手動審查或 release checklist 項目 |

BDD closure 不要求每個 scenario 都使用 Cucumber-style runner。它要求每個關鍵 scenario 有明確的驗證路徑，且沒有模糊的「已記錄但未測試」狀態。

> **輸出模板**：BDD Execution Closure 完成後，使用 [`templates/bdd-scenario-template.md`](templates/bdd-scenario-template.md) 記錄行為場景、acceptance criteria、validation target 與 traceability。

## 2. 文件優先 BDD 閉環（Docs-First BDD Closure Loop）

當在行為由人類可讀規格加上可執行測試管理的儲存庫中工作時，在改變可觀察行為之前保持產出同步：

| 步驟 | 行動 |
| --- | --- |
| 1 | 更新**擁有合約**（架構、API/介面、領域模型、錯誤處理、product brief 或同等文件），並在程式碼之前解決或標記開放問題 |
| 2 | 在專案指定位置添加或編輯**人類可讀的行為規格** |
| 3 | 在專案的**可執行行為測試**中鏡像相同的 scenarios；保持書面 scenarios 和可執行檢查同步 |
| 4 | 實作步驟定義、適配器、fixtures 和生產程式碼 |
| 5 | **完成定義：** 聚合的可執行行為套件或同等驗證在同一次變更中通過，且任何狀態表或可追溯性文件已更新 |

除非團隊已在可追溯的變更記錄中明確範圍化缺失的規格工作，否則不要僅使用**單元測試**而沒有行為規格對齊就合併可觀察行為或共享合約變更。

對於專案特定路徑、測試執行器和狀態表，使用應用程式儲存庫自己的治理文件；不要將這些細節複製到這個可重複使用的技能中。

## 3. Test Strategy Gate（測試策略關卡）

在實作前區分「保護舊行為」和「驗證新 code」。高總體覆蓋率可以證明舊行為受到保護，但不能證明新產生或新撰寫的 code 是正確的。

| 目標 | 目的 | 必要驗證 |
| --- | --- | --- |
| 既有 / 舊有行為 | 防止 regression 並保護已知 contracts | 執行覆蓋受影響行為的既有 unit、BDD、contract、integration 和 regression tests |
| 新需求或新 code | 證明新行為正確、安全且與 docs 一致 | 先寫或更新 BDD，然後在可行時在 production code 前加上 failing tests 或 executable specs。將 changed/new-code 覆蓋率與整體專案覆蓋率分開追蹤 |
| AI 生成的 code | 補償合理但錯誤的 code | 要求 BDD scenario、unit/contract tests，以及專注於 intent、edge cases 和安全/所有權邊界的人類審查 |
| 業務規則 / 演算法 | 捕捉通過範例但規則錯誤的情況 | 加上 property-based tests、invariant tests、targeted mutation checks 或 table-driven edge cases |
| 關鍵條件判斷 / 驗證邏輯 | 證明當邏輯錯誤時 tests 會失敗 | 在可行時加上 mutation testing，或手動測試如果 guards 被移除會失敗的 negative cases |
| 資料庫 / 持久化行為 | 保護真實的狀態轉換和遷移 | 加上 fixture-backed repository tests、migration tests 或針對代表性資料的 integration tests |
| 效能敏感行為 | 防止功能正確的 code 超出 latency、throughput、error-rate 或資源預算 | 根據風險加上 load、stress、spike 或 soak tests；追蹤 P95/P99 latency、throughput、error rate、CPU、memory、disk、network、database connections、queue depth 和 external-call volume |

新需求的建議順序：

1. BDD scenarios
2. 針對新行為的 failing unit、contract、property 或 integration tests
3. Production code
4. 關鍵規則的 mutation/negative checks
5. 當變更可能影響 latency、throughput 或資源時的效能 smoke 或 targeted scenario
6. 將 planning docs、BDD、tests 和效能證據並排進行人類審查

### Mutation Testing / Test Effectiveness Check

Mutation testing 是測試有效性檢查，不是 coverage KPI。當變更涉及 AI 生成邏輯、權限/安全/金流、domain invariant、複雜條件判斷、或 refactor 宣稱無行為變更時，使用 targeted mutant flow：

1. 描述要防止的錯誤或 invariant break。
2. 產生小型 mutants，例如 boundary、comparison、boolean、nullability、error handling 或 guard 移除。
3. 過濾 equivalent mutants，避免把語意相同的版本當成測試缺口。
4. 若 mutant survived，補 BDD / unit / property / contract / fixture test，或縮小 correctness claim。

通過標準不是「mutant 越多越好」，而是至少能殺掉代表真實風險的 mutant；若沒有合適工具，可用手動 negative check 或 code review 方式模擬。

### Performance Test 選型 → 見 sd-validation

效能測試類型（load / stress / spike / soak）的選型表與最低指標**不在此處複製**，canonical 在 [`validation.md`](validation.md) §1 效能測試關卡。當測試策略涉及效能風險時，從本 slice 的「效能敏感行為」列指向 `validation.md` 取得選型細節與 metric budget。

## 4. 測試策略定義（專案內部問題清單）

實作前，依專案具體情境定義測試策略：

| 問題 | 必要行動 |
| --- | --- |
| 哪些既有行為不能回歸？ | 為受影響的舊行為執行或添加回歸測試 |
| 引入了什麼新行為？ | 在可行時在生產程式碼之前撰寫 BDD 和失敗測試或可執行規格 |
| 總覆蓋率是否隱藏了未測試的新程式碼？ | 分別追蹤變更/新程式碼覆蓋率與整個專案覆蓋率 |
| 邏輯是否規則密集或安全敏感？ | 添加 targeted mutation checks、基於屬性的測試、不變量測試或負面案例 |
| 測試是否真的能抓到錯誤？ | 對 AI-generated logic、critical branch、domain invariant 或 refactor-no-behavior-change claim，產生小型 mutant 或手動 negative check；若 mutant survived，補 validation target 或縮小完成宣告 |
| 持久化是否重要？ | 添加 fixture 支援的資料庫/儲存庫/遷移測試或整合測試 |
| 程式碼是否由 AI 生成？ | 需要測試加上針對規劃文件、BDD、合約和邊緣案例的人類審查 |
| 這是嵌入式或硬體支援的？ | 分開主機可重複測試與僅目標或硬體在迴路中的證據；記錄板子、接線、引腳/匯流排設定、韌體版本、日誌和觀察到的偏差 |
| 這個變更是否涉及效能？ | 首先添加一個小的、可重複的效能檢查；根據風險選擇負載、壓力、尖峰或浸泡測試。追蹤 P95/P99 延遲、吞吐量、錯誤率和資源使用率，而不僅是平均延遲 |

## 5. Test-First Ordering（Framework / Runtime / Governance 升級強制）

> ⚠️ 上方策略表中「在可行時在生產程式碼之前撰寫 BDD 和失敗測試」對一般開發為**建議**；但若變更涉及 **framework / runtime / governance / workflow / validation / scenario / metadata / compiler / generated artifact** 改動，由 [`governance/lifecycle/system-upgrade-governance.md`](../../governance/lifecycle/system-upgrade-governance.md) §3 規則 9 升級為**強制**順序：

```
1. 列出 Phase N 期望可觀察行為（檔案、runtime.db、agent action）
2. 寫對應 validation/scenarios/<domain>/<id>-v1.yaml
3. 驗證 scenarios 目前 fail（fail-by-absence）
4. 才開始 Phase N 實作
5. Commit message 含「scenarios pre-written: <hash>, now passing」
```

**豁免**（須明寫理由）：doc-only trial / bug fix / typo / 探索性 spike
**不可豁免**：runtime.db schema、enforcement rule、blocking gate、compiler、`generated_surfaces`

完整原則見 [`intelligence/engineering/development/test-first-framework-upgrade.md`](../../intelligence/engineering/development/test-first-framework-upgrade.md)；
對應 validation scenario [`validation/scenarios/failure-derived/test-first-for-framework-upgrades-v1.yaml`](../../validation/scenarios/failure-derived/test-first-for-framework-upgrades-v1.yaml)。
