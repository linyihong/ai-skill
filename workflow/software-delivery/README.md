# Software Delivery Workflow

`workflow/software-delivery/` 負責「App 開發審查與指引的執行流程」。本目錄保存 agent 在進行 app 開發審查時可照著執行的 planning flow、review flow、handoff flow、review checklists 與 contract-first 開發流程，讓開發與審查過程可重複、可驗證。

## 何時進入此 Workflow

當工作任務需要**開發、實作、修改程式碼、進行 code review / design review / release review** 時，agent 應自行判斷並載入本 workflow。不需要 runtime 觸發——agent 知道什麼時候需要開發。

進入方式：
1. 讀取 [`execution-flow.md`](execution-flow.md) thin index 了解執行流程與 focused loading surfaces
2. 依 task intent 載入需要的 execution surface：[`intake.md`](intake.md)、[`contracts.md`](contracts.md)、[`ui-contracts.md`](ui-contracts.md)、[`test-strategy.md`](test-strategy.md)、[`validation.md`](validation.md)、[`closure.md`](closure.md)、[`surgical-changes.md`](surgical-changes.md)
3. 依流程的 Start From Evidence / Change Intake 開始；新需求、重構、parity、缺失資訊或既有專案回填時載入 [`intake.md`](intake.md)
4. 需要 artifact 規範時參考 [`artifact-gates.md`](artifact-gates.md)
5. 需要審查檢查清單時參考 [`review-checklist.md`](review-checklist.md)
6. 需要完整開發流程 overview 或 embedded / producer-consumer fallback 時參考 [`development-process.md`](development-process.md)
7. 需要前端、行動、CLI、SDK 或其他 consumer surface 的 Screen Mapping、Consumer / UI Behavior / Screen / ViewModel Contract 或 Screen Traceability 時參考 [`ui-contracts.md`](ui-contracts.md)
8. 需要 pre-build interrogation / product impact alignment / requirements cognition / BDD-lite / acceptance criteria / ambiguity resolution 時參考 [`requirements/`](requirements/README.md)
9. 需要 architecture fit analysis、DDD / CQRS / event sourcing / microservices decision 時參考 [`architecture/`](architecture/README.md)
10. 需要 Simplicity First / Surgical Changes / Think Before Coding 的行為範例時參考 [`examples/EXAMPLES.md`](examples/EXAMPLES.md)；examples 預設 suppress，僅在明確要求範例或 ambiguity 時載入
11. PR 觸動效能敏感路徑或含 AI 生成程式碼時，執行 [`perf-risk-gate.md`](perf-risk-gate.md) 的 5 步檢查（靜態 anti-pattern scan、hot-path micro-benchmark、reviewer perf checklist、pre-deploy observability gate、canary rollout）
12. 當成功訊號可能不同於真實系統狀態時，載入 validation reasoning：[`state-visibility-gap.md`](../../intelligence/engineering/execution/validation-reasoning/state-visibility-gap.md)、[`evidence-model.md`](../../intelligence/engineering/execution/validation-reasoning/evidence-model.md)、[`evidence-chain-validation.md`](../../intelligence/engineering/execution/validation-reasoning/evidence-chain-validation.md)、[`evidence-depth.md`](../../intelligence/engineering/execution/validation-reasoning/evidence-depth.md)

## Scope

本 workflow 涵蓋以下流程與審查類型：

### 開發流程

- **Requirements Stage**：Pre-build interrogation + product impact alignment + BDD-lite / requirements cognition，包含需求拷問、framework source-of-truth discovery、duplication risk check、Impact Map × Customer Journey Map、behavior-driven discovery、acceptance definition、ambiguity resolution、traceability 與 validation target。
- **Architecture Stage**：domain architecture cognition，包含 DDD fit、bounded context discovery、consistency boundary design、architecture escalation。
- **Contract-First Development Process**：從企劃書到實作的完整開發流程，包含 Default Flow、Required Contracts、Product Brief Validation Gate、Change Intake Gate、Contract Governance Gate、Traceability Gate、BDD Execution Closure、Test Strategy Gate、Embedded/Hardware Flow、Missing Information Gate、Existing Project Documentation Backfill 等。
- **UI / Consumer Contract Process**：在 provider/consumer 平行實作前建立 Screen Mapping、Consumer Contract、UI Behavior Contract、Screen Contract、Frontend ViewModel Contract、Accessibility Contract 與 Screen Traceability，避免 AI agent 只依 API shape 生成語意脫節的前端。
- **Evidence-Oriented Validation**：當 API/adapter/UI 成功訊號不足以證明 persisted、external、identity-specific 或 user-observable state 時，依 engineering validation reasoning 建立 evidence chain、選擇 evidence depth，必要時要求 live system proof 與 independent observation。
- **Refactor / Replacement Parity**：當新入口、平台遷移、工具改寫或架構重組要取代舊能力時，先建立新舊能力 parity inventory，逐項列出舊入口、現有能力、副作用、外部依賴、新入口、parity 狀態與測試證據。

### 審查類型

- **Design Review**：在實作前審查設計文件、API contract 與架構決策。
- **Architecture Review**：審查 architecture fit、DDD adoption、bounded context、overengineering 與 minimality。
- **Code Review**：審查實作程式碼的正確性、安全性與效能。
- **Release Review**：審查 release 前的準備狀態（build、signing、provisioning）。
- **Security Review**：審查授權、認證、資料儲存與傳輸安全。
- **Contract Governance Review**：審查 API contract 的相容性與版本管理。
- **Embedded Firmware Review**：審查韌體/硬體相關實作。

## 核心原則

1. **Review 是預防不是懲罰**。目的是在問題進入 production 前發現，不是追究責任。
2. **Checklist 是輔助不是取代**。Checklist 確保基本項目不被遺漏，但 reviewer 仍需使用工程判斷。
3. **Review 結果必須 actionable**。每個 finding 應包含：問題描述、風險等級、建議修復方式。
4. **Review 記錄應可追溯**。每個 review 的 finding、decision 與 resolution 應可追溯到對應的 commit 或 ticket。
5. **Simplicity First（簡潔優先）**：從最簡單的實作開始。不要預先加入抽象層、Strategy pattern、或 speculative features。當需求證明需要複雜度時再重構。參見 [`examples/EXAMPLES.md`](examples/EXAMPLES.md) §2。
6. **Surgical Changes（外科手術式修改）**：只改解決問題所需的行。匹配既有 code style，不要順便 refactor 不相關的 code。參見 [`surgical-changes.md`](surgical-changes.md)（`sd-surgical-caveats` slice）和 [`examples/EXAMPLES.md`](examples/EXAMPLES.md) §3。
7. **Parity Before Replacement（替換前先對照）**：重構、遷移或 replacement 若會替代既有行為、入口、腳本、API、資料流程或操作能力，先盤點舊能力到新能力的對照與驗證證據，再開始實作。

## 與既有層的關係

- `workflow/software-delivery/` 是 App 開發指引執行流程的主要入口。所有 agent 應優先參考本目錄的內容。
- `analysis/development-guidance/` 提供安全控制、實作模式、平台指引、語言陷阱的 catalog 參考，被本 workflow 引用。
- `analysis/repo/` 可被本 workflow 引用來分析 repository 結構。
- `intelligence/` 可被本 workflow 引用來輔助工程判斷。
- `workflow/software-delivery/requirements/` 提供 requirements cognition 的執行流程；其 source intelligence 來自 `intelligence/engineering/requirements/`。
- `workflow/software-delivery/architecture/` 提供 architecture fit 的執行流程；其 source intelligence 來自 `intelligence/engineering/architecture/architectural-fit/` 與 `intelligence/engineering/architecture/domain-modeling/`。
- `governance/ai-runtime-governance/software-delivery-governance.md` 定義 requirements / behavior / contract delivery gates。
- `governance/ai-runtime-governance/software-delivery-architecture-governance.md` 定義 software-delivery architecture governance gate，但不把 DDD promotion 成 runtime invariant。
- `feedback/history/development-guidance/` 儲存開發指引的具體課程記錄。
- `skills/app-development-guidance/` 是原始 skill 目錄，已刪除。所有內容已遷移至本層。

## 遷移狀態

| 來源 | 目標 | 狀態 |
|------|------|------|
| `skills/app-development-guidance/WORKFLOW.md` | [`execution-flow.md`](execution-flow.md)、[`artifact-gates.md`](artifact-gates.md)、[`analysis/development-guidance/risk-translation.md`](../../analysis/development-guidance/risk-translation.md) | ✅ 已遷移，舊目錄已刪除 |
| `skills/app-development-guidance/process/` | [`development-process.md`](development-process.md) | ✅ 已遷移，舊目錄已刪除 |
| `skills/app-development-guidance/checklists/` | [`review-checklist.md`](review-checklist.md) | ✅ 已遷移，舊目錄已刪除 |

## 已提取內容

| 檔案 | 來源 | 說明 |
|------|------|------|
| [`execution-flow.md`](execution-flow.md) | `WORKFLOW.md` §1, §5-8（已刪除） | Start From Evidence、Change Intake、BDD Closure Loop、SDK Defect Closure、Same-Session Closure、Performance Gate、Backfill Rules、Validate |
| [`execution-flow.yaml`](execution-flow.yaml) | `execution-flow.md` | Software delivery execution executable contract：change intake、requirements、BDD closure、parity、performance、validation gates |
| [`intake.md`](intake.md) | `execution-flow.md` §1/§6 + `development-process.md` intake gates | Focused execution surface：需求接收、Change Intake、Pre-build Interrogation、Requirements Cognition、Parity Gate、Product Brief Validation、Missing Information、Backfill |
| [`contracts.md`](contracts.md) | `development-process.md` contract gates | Focused execution surface：Required Contracts、Contract Governance、Traceability、Contract-First Rules |
| [`ui-contracts.md`](ui-contracts.md) | `development-process.md` frontend / consumer contract gap | Focused execution surface：Screen Mapping、Consumer Contract、UI Behavior Contract、Screen Contract、Frontend ViewModel Contract、Accessibility Contract、Screen Traceability |
| [`test-strategy.md`](test-strategy.md) | `execution-flow.md` §2/§4 子節 + `development-process.md` BDD/Test Strategy gates | Focused execution surface：Docs-first BDD closure、test strategy、mutation testing、test-first ordering |
| [`validation.md`](validation.md) | `execution-flow.md` §5/§7 | Focused execution surface：validation、performance gate、old/new behavior proof、completion evidence |
| [`closure.md`](closure.md) | `execution-flow.md` §8 + `development-process.md` DoR/DoD | Focused execution surface：Definition of Ready/Done、handoff、close-loop、reusable lesson feedback |
| [`surgical-changes.md`](surgical-changes.md) | `execution-flow.md` §9 | Focused failure surface：surgical change discipline、diff purity、orphan cleanup boundary |
| [`requirements/pre-build-interrogation.md`](requirements/pre-build-interrogation.md) | mattpocock/skills `/grill-me` pattern + Ai-skill framework failure learning | Plan / implementation 前的需求拷問、framework discovery 與 source-of-truth duplication gate |
| [`artifact-gates.md`](artifact-gates.md) | `DOCUMENTATION.md`（已刪除） | Reusable Note Structure、Content Classification、Guidance Boundary、Linked Update Statement、Good Guidance Criteria |
| [`artifact-gates.yaml`](artifact-gates.yaml) | `artifact-gates.md` | Software delivery artifact executable contract：artifact shape、owner layer、sanitization、linked updates、quality gates |
| [`analysis/development-guidance/risk-translation.md`](../../analysis/development-guidance/risk-translation.md) | `WORKFLOW.md` §2-5（已刪除） | Risk Translation Table、Owner Layer Selection、Control Definition、Guidance Classification、Linked Updates |
| [`review-checklist.md`](review-checklist.md) | `skills/app-development-guidance/checklists/`（已刪除） | 6 種審查 checklist 的 catalog（Mobile Design Review、Mobile PR Review、Mobile Release Review、API Security Review、Contract Governance Review、Embedded Firmware Review） |
| [`development-process.md`](development-process.md) | `skills/app-development-guidance/process/README.md`（已刪除） | Contract-first 開發流程：Default Flow、Required Contracts、Product Brief Validation Gate、Change Intake Gate、Contract Governance Gate、Traceability Gate、BDD Execution Closure、Test Strategy Gate、Embedded/Hardware Flow、Missing Information Gate、Existing Project Documentation Backfill、Contract-First Rules、Definition of Ready/Done |

## 建議 Workflow 流程

### Design Review Flow

```
1. 確認審查範圍（新功能 / 架構變更 / API 變更）。
2. 讀取設計文件或 RFC。
3. 檢查設計是否涵蓋：
   ├─ Pre-build interrogation：goal、scope、non-goals、acceptance、framework source-of-truth、duplication risk。
   ├─ Requirements cognition：actor intent、behavior boundary、acceptance criteria、ambiguity disposition。
   ├─ API contract（request/response 格式、錯誤處理）。
   ├─ Domain / architecture fit：bounded context、invariant、consistency boundary、CRUD / DDD Lite / Full DDD decision。
   ├─ 資料模型與儲存方案。
   ├─ 安全考量（授權、認證、資料保護）。
   └─ Validation target 與測試策略（BDD-lite、單元測試、合約測試、整合測試、E2E 測試）。
4. 對每個項目給出 verdict：approve / approve-with-comments / changes-requested。
5. 記錄所有 finding 與 decision。
```

### Code Review Flow

```
1. 確認 review 範圍（PR diff / 特定檔案 / 完整功能）。
2. 讀取 diff，檢查：
   ├─ 邏輯正確性。
   ├─ 邊界條件處理。
   ├─ 錯誤處理與復原。
   ├─ 安全性（注入、授權缺失、敏感資料暴露）。
   ├─ 效能（不必要的迴圈、過度查詢、記憶體洩漏）。
   └─ 程式碼風格與可讀性。
3. 對每個 finding 標記 severity：blocker / major / minor / nit。
4. 提供具體的修改建議（不只是「這有問題」，而是「建議改成 X」）。
5. 確認所有 blocker 與 major finding 被解決後 approve。
```

### Release Review Flow

```
1. 確認 release 版本號與 changelog。
2. 檢查 build 狀態（CI/CD pipeline 是否通過）。
3. 檢查 signing 與 provisioning（iOS code signing、Android keystore）。
4. 檢查相依套件版本（無已知漏洞的版本）。
5. 檢查 release note 是否完整（新功能、修復、已知問題）。
6. 確認 rollback 計畫。
7. 給出 release verdict：go / go-with-caveats / block。
```

## 產出格式

本 workflow 各階段使用標準化輸出模板，確保產出格式一致、可追溯、可被後續階段自動消費：

| 模板 | 對應階段 | 用途 |
|------|---------|------|
| [`templates/product-impact-alignment-template.md`](templates/product-impact-alignment-template.md) | Product Impact Discovery | 記錄 Impact Map、Customer Journey Map、cross-check decision 與進入 BDD 前的缺口 |
| [`templates/change-brief-template.md`](templates/change-brief-template.md) | Change Intake | 記錄變更類型、證據、範圍、blocker 評估 |
| [`templates/contract-template.md`](templates/contract-template.md) | Contract Governance | 記錄 domain model、架構決策、API / error / consumer / UI 合約 |
| [`templates/bdd-scenario-template.md`](templates/bdd-scenario-template.md) | BDD Closure Loop | 記錄 requirement link、behavior boundary、Given/When/Then、acceptance criteria、validation target、regression scope |
| [`templates/implementation-plan-template.md`](templates/implementation-plan-template.md) | Implementation | 記錄任務拆解、檔案路徑、驗收條件、風險評估 |
| [`templates/review-report-template.md`](templates/review-report-template.md) | Review（6 種類型） | 記錄 finding、decision、reviewed artifacts |

每次 review 應產出：

- **Review 摘要**（≤200 tokens）：審查類型、範圍、verdict。
- **Finding 清單**（每個 finding ≤100 tokens）：問題描述、severity、建議修復方式。
- **Decision 記錄**（≤100 tokens）：最終決定、決定理據、相關連結。
