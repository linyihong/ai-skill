# 開發流程

`workflow/software-delivery/development-process.md` 是從 `skills/app-development-guidance/process/README.md`（已刪除）提取的開發流程。本文件定義 contract-first 開發流程：在團隊開始 coding 前先釐清行為與領域語言，然後讓每個實作面向對著同一份 contract 工作。

> **遷移狀態**：`skills/app-development-guidance/process/README.md` 已刪除。此文件為 canonical source，新內容請直接寫入此文件。

Contract / BDD / artifact / performance 的治理 gate 見 [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md)；本檔保留 development process 的詳細流程與專案 artifact 形狀。

## Default Flow（預設流程）

| 步驟 | 產出 | 說明 |
| --- | --- | --- |
| 1. 企劃書 / product brief | 目標、使用者、範圍、non-goals、限制、驗證證據 | 明確記錄假設；在將企劃視為實作輸入前先驗證 product intent；標記未知項目而非憑空創造確定性 |
| 2. Product impact alignment | Impact Map、Customer Journey Map、cross-check decision | 先驗證 Why / Who / How / What 是否對準 journey pain，避免 AI 快速產出偏離 impact 的功能 |
| 3. Requirements cognition | Actor intent、behavior boundary、acceptance criteria、ambiguity disposition | 先用 [`requirements/`](requirements/README.md) 穩定需求與行為，不把 BDD 縮成 Gherkin，也不讓 agent 自行補需求 |
| 4. BDD-lite 行為合約 | Scenario tables、behavior contract、validation target、traceability | 用領域語言描述使用者/系統行為，並明確連到 proof target |
| 5. Screen Mapping | Scenario -> screens -> APIs -> tables / stores | 將 BDD scenario 落到畫面、provider operation 與資料所有權，作為後續 traceability seed；細節見 [`ui-contracts.md`](ui-contracts.md) |
| 6. Domain architecture cognition | Bounded Context map、invariant、consistency boundary、architecture fit | 從行為邊界推導 domain boundary；依複雜度選 CRUD / DDD Lite / Full DDD，不預設 DDD |
| 7. Domain Model Contract | Entities、value objects、commands、events、invariants | 這是核心 contract；定義必須永遠為真的事項 |
| 8. Architecture Contract | 邊界、依賴關係、資料所有權、runtime/deployment 形狀 | 定義哪些層可以依賴哪些層、決策落在哪裡 |
| 9. API Contract | OpenAPI/GraphQL/schema/events/RPC/message contracts | 當多個 surface 或服務整合時至關重要；先建立 provider operation 形狀供 consumer 對齊與修訂 |
| 10. Consumer Contract | Consumer needs、freshness、loading、empty/error behavior、permissions | 定義前端/行動/CLI/SDK/job 等 consumer 為了完成行為需要什麼；若暴露 API 缺口，回修 API / Error Contract |
| 11. Error Handling Contract | Error taxonomy、retry 規則、驗證錯誤、使用者可見訊息、logging | 在實作前設計 failure 行為；錯誤也是 contract 的一部分，並需映射到 consumer / UI error behavior |
| 12. Frontend / Consumer Contract Gate | Screen Contract、UI Behavior Contract、ViewModel Contract、Accessibility Contract | 當存在 UI 或 consumer surface 時載入 [`ui-contracts.md`](ui-contracts.md)，固定 AI agent 生成 UI 前需要的 screen、interaction 與 display semantics |
| 13. 平行實作規劃 | 每個 context、provider、consumer 和 contract 的工作切片 | 只有當共享 contracts 足夠穩定時，前端/後端或 provider/consumer 才能同時開始 |
| 14. 後端 / 服務實作 | Behavior + domain + API contract 實作 | 如果沒有後端，用 local service、library、worker 或 platform 實作取代 |
| 15. 前端 / 客戶端實作 | Mock API、schema-first client、UI behavior、screen contract、view model mapper | 如果沒有前端，用 CLI、SDK、mobile screen、job runner 或 consumer integration 取代 |
| 16. 自動化測試 | Unit、BDD、API/consumer/UI contract、schema tests | 測試應證明 behavior correctness、domain invariants、provider/consumer 相容性和 UI state transitions |
| 17. 效能測試計畫 | Load、stress、spike 或 soak 範圍、metric budget、runner、證據位置 | 當 latency、throughput、資源使用率、concurrency、啟動時間、背景任務、資料庫存取、外部呼叫量、caching 或 batching 可能改變時需要 |
| 18. 整合測試 | End-to-end 或 component integration 證據 | 驗證真實 adapters、auth/session、錯誤路徑、效能敏感路徑、screen flows 和跨 context 流程 |

> **輸出模板**：Default Flow 完成後，使用 [`templates/implementation-plan-template.md`](templates/implementation-plan-template.md) 記錄實作計畫。

## Required Contracts（必要合約）→ 已抽出為獨立 slice

> **Cognitive Slice 重構（Phase 2）**：本節連同 Contract Governance Gate、Traceability Gate、Contract-First Rules 已抽出為 [`contracts.md`](contracts.md)（slice `sd-contracts`，type `execution`，tags `artifact-gate, contract, traceability`）。canonical content 在 `contracts.md`，此處不再保留正文以避免 dual source-of-truth。

## Initial Documentation Pack（初始文件包）→ 已抽出為獨立 slice

> **Cognitive Slice 重構（Phase 2）**：Initial Documentation Pack、Product Brief Validation Gate（含 Product Impact Alignment Check）、Change Intake Gate（含 Refactor / Replacement Parity Inventory）、Missing Information Gate、Existing Project Documentation Backfill 已連同 execution-flow.md §1 / §6 跨檔同批抽出為 [`intake.md`](intake.md)（slice `sd-intake`，type `execution`，tags `requirements, parity, intake`）。需求接收 / brief 驗證 / parity / 缺失資訊處置 / 既有專案回填的 canonical content 在 `intake.md`，此處不再保留正文以避免 dual source-of-truth。

## Product Brief Validation Gate（Product Brief 驗證關卡）→ 已抽出為獨立 slice

見 [`intake.md`](intake.md) §Product Brief Validation Gate（`sd-intake`）。

## Change Intake Gate（變更接收關卡）→ 已抽出為獨立 slice

見 [`intake.md`](intake.md) §Change Intake Gate（含 Refactor / Replacement Parity Inventory，`sd-intake`）。

## Contract Governance Gate / Traceability Gate → 已抽出為獨立 slice

> **Cognitive Slice 重構（Phase 2）**：Contract Governance Gate 與 Traceability Gate 已抽出為 [`contracts.md`](contracts.md)（slice `sd-contracts`）。文件優先順序、衝突分類、雙向可追溯性的 canonical content 在 `contracts.md`，此處不再保留正文以避免 dual source-of-truth。

## BDD Execution Closure / Test Strategy Gate → 已抽出為獨立 slice

> **Cognitive Slice 重構（Phase 2）**：BDD Execution Closure 與 Test Strategy Gate（含 Mutation Testing / Test Effectiveness Check）已連同 execution-flow.md §2 Docs-First BDD Closure Loop + §4 子節「測試策略定義」+「Test-First Ordering」抽出為 [`test-strategy.md`](test-strategy.md)（slice `sd-test-strategy`，type `execution`，tags `artifact-gate, test, bdd`）。canonical content 在 `test-strategy.md`，此處不再保留正文以避免 dual source-of-truth。
>
> **perf 測試類型表**（load / stress / spike / soak）原於本節重複，現已歸 canonical 至 [`validation.md`](validation.md) §1 效能測試關卡（`sd-validation`）；`test-strategy.md` 引用 perf 風險作策略選型考量，不複製選型表。

## Embedded / Hardware Product Flow（嵌入式/硬體產品流程）

當專案涉及 firmware、sensors、boards、UART/I2C/SPI/BLE/CAN/GPIO、RTOS tasks、hardware bring-up 或 host/target validation 時使用這個流程：

| 層級 | Contract | 說明 |
| --- | --- | --- |
| Datasheet / vendor spec | Electrical interface、protocol bytes、timing、default parameters、valid ranges、errata | 將 vendor docs 視為 protocol truth；分別記錄觀察到的偏差 |
| Protocol Parsing Contract | Frame format、state machine、length/checksum rules、command/ACK/report shapes、fixtures、invalid frames | 將 byte-level parsing 與產品意義分離 |
| BDD Behavior | 使用者/系統行為、裝置狀態、設定/配置流程、故障處理、target events | BDD 使用領域術語，而非 raw registers 或 UART calls |
| Domain Model Contract | Pure DTOs、units、ranges、invariants、timestamps、validity windows | 將 HAL/RTOS types 保持在 domain objects 之外 |
| Hardware Context Contract | Board-specific pins、UART/I2C/SPI bus、baud/rates、buffers、interrupts、power modes、injected configuration | Board 變更應改變 context/config，而非 protocol/domain logic |
| Embedded Architecture Contract | Driver/service/application layering、task/ISR boundaries、queues、ownership、concurrency、lifecycle、error escalation | Drivers 處理 bytes；services 解析；applications 決定產品行為 |
| Public API / Interface Contract | Context lifecycle、callbacks/subscriptions、commands、errors、consumer ownership、multi-device rules | 避免平行的第二個 context API，除非 contracts 先被修訂 |
| Test Plan | Host unit tests、protocol fixtures、negative cases、property/invariant tests、simulator/mocks、hardware-in-loop、bring-up log evidence | 區分 host-repeatable proof 和 bench-only evidence |

在 firmware code 之前：

1. 讀取 datasheet/protocol spec 和專案 contracts
2. 確認 hardware context 是每個 board 可注入的，而非寫死為唯一 truth source
3. 撰寫或更新 BDD 和 protocol/domain/API contracts
4. 為 protocol parsing 和 negative cases 加上 host-side fixtures
5. 只在無法在 host 上證明的證據才定義 target 或 hardware-in-loop validation
6. 記錄 bring-up evidence：board revision、wiring、pins、bus settings、firmware version、logs 和已知偏差

## Missing Information Gate（缺失資訊關卡）→ 已抽出為獨立 slice

見 [`intake.md`](intake.md) §Missing Information Gate（`sd-intake`）。canonical content 在 `intake.md`，此處不再保留正文以避免 dual source-of-truth。

## Existing Project Documentation Backfill（既有專案文件回填）→ 已抽出為獨立 slice

> **Cognitive Slice 重構（Phase 2）**：既有專案文件回填已連同 execution-flow.md §6 作為 backfill 條件子流程抽出至 [`intake.md`](intake.md) §Backfill（`sd-intake`，`tags: domain-specific,backfill`）。canonical content 在 `intake.md`，此處不再保留正文以避免 dual source-of-truth。僅在處理「已實作但文件缺失」的專案時載入。

## Contract-First Rules（合約優先規則）→ 已抽出為獨立 slice

> **Cognitive Slice 重構（Phase 2）**：Contract-First Rules 已抽出為 [`contracts.md`](contracts.md)（slice `sd-contracts`）。contract ownership 與 contract-first 紀律的 canonical content 在 `contracts.md`，此處不再保留正文以避免 dual source-of-truth。

## When Frontend And Backend Do Not Both Exist（當前後端不都存在時）

將「前端」和「後端」替換為 producer/consumer 角色：

| 原始角色 | 通用角色 |
| --- | --- |
| 前端 | Consumer：UI、CLI、SDK、job、mobile screen、test harness |
| 後端 | Provider：API、domain service、library function、local adapter、worker |

流程仍然適用：

1. 定義行為
2. 定義 domain invariants
3. 定義 provider/consumer contract
4. 各自對著 mock、fixture 或 schema 建構
5. 用 contract 和 integration tests 證明相容性

## Minimum Definition Of Ready / Minimum Definition Of Done → 已抽出為獨立 slice

> **Cognitive Slice 重構（Phase 2）**：Minimum Definition Of Ready 與 Minimum Definition Of Done 已連同 execution-flow.md §8 Feed Back Reusable Lessons 抽出為 [`closure.md`](closure.md)（slice `sd-closure`，type `execution`，tags `closure, handoff, extraction-to-intelligence`）。DoR / DoD 收尾檢核的 canonical content 在 `closure.md`，此處不再保留正文以避免 dual source-of-truth。
