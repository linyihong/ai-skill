# 開發流程

`workflow/software-delivery/development-process.md` 是從 `skills/app-development-guidance/process/README.md`（已刪除）提取的開發流程。本文件定義 contract-first 開發流程：在團隊開始 coding 前先釐清行為與領域語言，然後讓每個實作面向對著同一份 contract 工作。

> **遷移狀態**：`skills/app-development-guidance/process/README.md` 已刪除。此文件為 canonical source，新內容請直接寫入此文件。

## Default Flow（預設流程）

| 步驟 | 產出 | 說明 |
| --- | --- | --- |
| 1. 企劃書 / product brief | 目標、使用者、範圍、non-goals、限制、驗證證據 | 明確記錄假設；在將企劃視為實作輸入前先驗證 product intent；標記未知項目而非憑空創造確定性 |
| 2. AI 分析 + 模組拆分 | Bounded Context map、模組所有權、整合點 | 依領域職責拆分，而非僅依 UI 頁面或資料庫表格 |
| 3. BDD 行為 | Feature files 或 scenario tables | 用領域語言描述使用者/系統行為 |
| 4. Domain Model Contract | Entities、value objects、commands、events、invariants | 這是核心 contract；定義必須永遠為真的事項 |
| 5. Architecture Contract | 邊界、依賴關係、資料所有權、runtime/deployment 形狀 | 定義哪些層可以依賴哪些層、決策落在哪裡 |
| 6. API Contract | OpenAPI/GraphQL/schema/events/RPC/message contracts | 當多個 surface 或服務整合時至關重要；在實作前先寫好 request/response/error 形狀 |
| 7. Error Handling Contract | Error taxonomy、retry 規則、驗證錯誤、使用者可見訊息、logging | 在實作前設計 failure 行為；錯誤也是 contract 的一部分 |
| 8. 平行實作規劃 | 每個 context 和 contract 的工作切片 | 只有當共享 contracts 足夠穩定時，前端/後端才能同時開始 |
| 9. 後端 / 服務實作 | Behavior + domain + API contract 實作 | 如果沒有後端，用 local service、library、worker 或 platform 實作取代 |
| 10. 前端 / 客戶端實作 | Mock API、schema-first client、UI 行為 | 如果沒有前端，用 CLI、SDK、mobile screen、job runner 或 consumer integration 取代 |
| 11. 自動化測試 | Unit、BDD、API contract、schema tests | 測試應證明 domain invariants 和 contract 相容性 |
| 12. 效能測試計畫 | Load、stress、spike 或 soak 範圍、metric budget、runner、證據位置 | 當 latency、throughput、資源使用率、concurrency、啟動時間、背景任務、資料庫存取、外部呼叫量、caching 或 batching 可能改變時需要 |
| 13. 整合測試 | End-to-end 或 component integration 證據 | 驗證真實 adapters、auth/session、錯誤路徑、效能敏感路徑和跨 context 流程 |

> **輸出模板**：Default Flow 完成後，使用 [`templates/implementation-plan-template.md`](templates/implementation-plan-template.md) 記錄實作計畫。

## Required Contracts（必要合約）

不要假設每個專案都有前端/後端拆分。選擇符合架構的 contracts：

| 應用形狀 | 優先 contract |
| --- | --- |
| 前端 + 後端 | API Contract、Domain Model Contract、BDD scenarios |
| 僅後端/API | Domain Model Contract、API Contract、contract tests、與 consumer 的 integration tests |
| 僅前端應用 | UI behavior contract、local state/domain contract、mocked API/schema contract |
| 行動應用 | Screen/flow behavior、local storage/session contract、有遠端服務時加上 API Contract |
| CLI / 桌面 / 工具 | Command contract、input/output schema、domain model、fixture-based tests |
| Library / SDK | Public API contract、type/schema contract、examples、compatibility tests |
| 事件驅動 / worker | Event schema、command/event contract、idempotency 和 retry behavior |
| Embedded / firmware / 硬體產品 | Datasheet 或 protocol contract、hardware context contract、driver/service/application 邊界、BDD、host fixtures、hardware-in-loop checks |
| 靜態分析 / IDE extension / 開發者工具 | Rule catalog、diagnostic 或 command contract、pure kernel/adapter 邊界、fixture pairs、editor/CLI integration tests |

## Initial Documentation Pack（初始文件包）

當這個 workflow 被用於新功能或新專案時，agent 應協助產出第一版草稿，或針對以下文件提問缺失資訊：

| 文件 | 目的 | 如果缺失，詢問 |
| --- | --- | --- |
| Product Brief | 目標、使用者、範圍、non-goals、假設 | 這是為誰做的？解決什麼問題？明確排除什麼？ |
| Bounded Context Map | 模組/領域拆分與所有權 | 哪些概念會一起變動？哪些邊界不應洩漏？ |
| BDD Behavior | 使用者/系統行為場景 | 關鍵 happy path 和 failure path 是什麼？ |
| Domain Model Contract | 核心業務物件與 invariants | 什麼必須永遠為真？允許哪些狀態轉換？ |
| Architecture Contract | 層級、依賴關係、所有權、runtime 限制 | 哪一層擁有資料、side effects、安全性、持久化、外部呼叫？ |
| API / Interface Contract | Requests、responses、events、commands、public methods | 誰消費這個 contract？相容性如何測試？版本如何管理？ |
| Error Handling Contract | Error types、recovery、user messaging、logging | 哪些錯誤可 retry、可由使用者修正、致命、或安全敏感？ |
| Hardware / Firmware Contract | Datasheet/protocol truth、electrical interface、pin/context injection、driver/service/application 邊界、target 限制 | 哪些硬體事實是固定的？每個 board 注入什麼？host/target 測試如何進行？ |
| Test Plan | Unit、BDD、contract、integration tests | 什麼證明行為、invariants 和 integration contract？ |

這些文件可以從輕量的 Markdown 草稿開始。如果專案很小，放在一個 planning file 中；如果長大，拆成一個資料夾加上 `README.md` 和聚焦的子文件。

## Product Brief Validation Gate（Product Brief 驗證關卡）

企劃書 / Product Brief 不會因為存在就自動可信。在使用它作為 BDD、contracts、估算、實作切片或測試的來源之前，先把它當作獨立的 artifact 來驗證。

| Brief 項目 | 驗證問題 | 可接受的證據 |
| --- | --- | --- |
| Goal / problem | 問題是真實的、具體的、且與使用者/系統結果相關嗎？ | 使用者請求、利害關係人決策、支援 ticket、觀察到的工作流程、metric、分析發現、或明確假設 |
| Users / actors | Actors 有命名且對應到權限、角色、裝置、系統或外部服務嗎？ | 現有帳號/角色、UI/API 行為、領域文件、組織決策、或 open question |
| Scope | 能判斷現在要建什麼嗎？ | BDD scenario list、module/context map、已接受的功能列表、API/interface list |
| Non-goals | 排除的行為夠明確，能防止意外實作嗎？ | 已取消/延後/排除範圍的表格、issue decision、stakeholder 回答 |
| Assumptions | 假設是可測試的、有時效的、或標記為風險的嗎？ | 證據連結、驗證計畫、負責人、到期/審查日期 |
| Success criteria | 測試、審查、metric、效能預算或 demo 能證明它有效嗎？ | BDD acceptance criteria、contract tests、P95/P99 latency budget、throughput target、error-rate budget、release checklist、analytics/telemetry query、manual evidence |
| Constraints | 法律、安全、隱私、平台、硬體、預算、時程、相容性和營運限制有列出嗎？ | Policy、platform docs、architecture contract、risk review、hardware/vendor docs |
| Dependencies | 外部服務、vendor、團隊、generated clients、遷移、資料或硬體依賴有識別嗎？ | Integration contract、API docs、schema、vendor excerpt、migration plan、owner confirmation |
| Risks | 濫用、失敗、安全、隱私、replay、資料遺失和營運風險有列出 controls 或 blockers 嗎？ | Threat model、hardening note、controls/checklists、open blocker questions |

如果任何 brief 項目影響行為、domain invariants、API/interface 形狀、錯誤處理、安全性、儲存、所有權、測試、時程或 release gate 且無法驗證，它就是 blocker。在使用者確認、取得證據或明確排除該項目之前，不應繼續開發。

對於已實作優先的專案，根據可觀察的證據驗證回填的 Product Brief。無法恢復的 Product intent 可以保持 `unknown`，但每個已實作的行為仍需要 BDD、contract 和 test 證據。

對每個 major brief claim 使用以下狀態：

| 狀態 | 意義 | 必要行動 |
| --- | --- | --- |
| `validated` | 有證據或明確的使用者/利害關係人決策支持 | 連結證據或決策 |
| `assumption` | 合理但未經證明 | 加上負責人、驗證計畫、以及如果為假的影響 |
| `open question` | 在實作繼續前需要答案 | 提問並阻止受影響的工作 |
| `scoped out` | 明確不屬於當前工作範圍 | 記錄 non-goal 並防止意外實作 |
| `invalidated` | 證據與 brief 矛盾 | 在 code 之前修訂 brief、BDD、contracts 和 tests |

## Change Intake Gate（變更接收關卡）

在任何由這個 workflow 驅動的 code 變更之前，檢查專案的企劃書、product brief、planning docs、issue、ticket、PRD、design note、BDD、API contract 或同等專案 artifact。在實作前分類請求：

| 變更類型 | code 之前需要 |
| --- | --- |
| 新需求 / 功能 / 行為變更 | 先更新或建立 planning docs：Product Brief 或 change brief、BDD scenarios、受影響的 Domain Model Contract、Architecture Contract、API / Interface Contract、Error Handling Contract、實作切片和 tests。在 blocker questions 解決前不要開始 code |
| Bug 修復 | 確認預期行為 vs 實際行為、重現步驟或證據、受影響的 BDD scenario 或缺失 scenario、受影響的 contract/error handling、以及 regression test plan。如果修復改變了預期行為或 public contract，也視為新需求 |
| Refactor / 內部清理 | 確認沒有行為或 public contract 變更。如果行為、資料所有權、API、錯誤處理、安全性、儲存或 tests 改變，重新分類為新需求或 bug |
| 安全 / 強化變更 | 確認威脅或 failure mode、owner layer、必要 control、驗證方法、以及行為/API/contracts/checklists 是否需要改變 |

如果沒有 planning artifact，在實作前建立輕量的 change brief。如果請求是新需求，缺失的 planning docs 是 blockers；向使用者提問並在寫 code 前填寫 BDD/contracts。

> **輸出模板**：Change Intake 完成後，使用 [`templates/change-brief-template.md`](templates/change-brief-template.md) 記錄變更簡報。

## Contract Governance Gate（合約治理關卡）

每個有多份文件的專案必須定義當文件不一致時哪個 artifact 優先。除非專案有更強的本地規則，否則使用這個預設優先順序：

1. Governance / framework contract：repository 層級的 invariants、必要更新規則、依賴方向、命名、build/run 限制
2. Product plan / accepted brief：product intent、範圍、non-goals、已取消的需求、業務語言
3. BDD behavior：可觀察的使用者/系統行為和 acceptance criteria
4. Domain、architecture、API/interface、error handling、hardware 或 command contracts
5. 實作和 generated clients
6. Tests、fixtures 和 examples

如果較低層級發現較高層級是錯的，不要默默地「修正」code。將衝突分類為以下之一：

| 衝突類型 | 必要行動 |
| --- | --- |
| Product intent 已變更 | 更新 product brief 或 plan，然後更新 BDD/contracts/tests |
| BDD 缺失或過時 | 從證據回填或修訂 BDD，並連結受影響的 tests |
| Contract 過時 | 在同一次變更中更新 contract 和所有 consumers、mocks、generated clients、fixtures 和 tests |
| 實作 bug | 保持 docs 穩定，加上或更新 regression tests，然後修復 code |
| Test 或 fixture 過時 | 將 tests/fixtures 更新到當前 contract 並引用來源 |

明確記錄已取消、延後、排除範圍和無法由工具強制執行的項目。不要讓它們成為未來 agent 可能重新引入的隱形空缺。

> **輸出模板**：Contract Governance 完成後，使用 [`templates/contract-template.md`](templates/contract-template.md) 記錄合約。

## Traceability Gate（可追溯性關卡）

當專案是先實作後補文件時，要求雙向追溯：

| 連結 | 目的 |
| --- | --- |
| Product 或 rule ID -> BDD | 顯示哪個行為證明該需求 |
| BDD -> code refs | 顯示行為在哪裡實作 |
| BDD -> test refs | 顯示行為如何被驗證，或還有什麼 gap |
| Contract operation / command / diagnostic -> fixture | 顯示 provider/consumer 相容性和 edge cases |
| Generated client 或 SDK method -> API/OpenAPI/source contract | 防止手抄 endpoint 和 drift |

Stable IDs 可以是 feature IDs、rule IDs、operation IDs、route names、command names、diagnostic codes、event names 或 scenario tags。如果一個行為被有意識地記錄但未實作，標記為 `TBD`、`noop`、`not enforceable by tool`、`manual-only` 或 `out of scope`，並附上原因和負責人。

## BDD Execution Closure（BDD 執行閉環）

Narrative BDD 在回填期間是可接受的，但不能被當作完成的測試覆蓋率。對每個關鍵 scenario，記錄以下狀態之一：

| 狀態 | 意義 | 必要下一步 |
| --- | --- | --- |
| `automated` | Scenario 由 unit、contract、API、integration、E2E、fixture 或 runner test 覆蓋 | 連結 test path/name |
| `fixture-backed` | Scenario 由 checked-in input/output fixtures 證明但沒有完整 runner | 連結 fixture 和 assertion 負責人 |
| `manual-evidence` | Scenario 需要手動、UI、bench 或外部服務證據 | 記錄執行步驟、證據和限制 |
| `pending-runner` | Gherkin 存在但沒有 runner/step definition 被接上 | 加上 runner 選擇或對應到可執行的 test type |
| `not-automatable` | Tooling 無法直接強制執行 | 說明手動審查或 release checklist 項目 |

BDD closure 不要求每個 scenario 都使用 Cucumber-style runner。它要求每個關鍵 scenario 有明確的驗證路徑，且沒有模糊的「已記錄但未測試」狀態。

> **輸出模板**：BDD Execution Closure 完成後，使用 [`templates/bdd-scenario-template.md`](templates/bdd-scenario-template.md) 記錄行為場景。

## Test Strategy Gate（測試策略關卡）

在實作前區分「保護舊行為」和「驗證新 code」。高總體覆蓋率可以證明舊行為受到保護，但不能證明新產生或新撰寫的 code 是正確的。

| 目標 | 目的 | 必要驗證 |
| --- | --- | --- |
| 既有 / 舊有行為 | 防止 regression 並保護已知 contracts | 執行覆蓋受影響行為的既有 unit、BDD、contract、integration 和 regression tests |
| 新需求或新 code | 證明新行為正確、安全且與 docs 一致 | 先寫或更新 BDD，然後在可行時在 production code 前加上 failing tests 或 executable specs。將 changed/new-code 覆蓋率與整體專案覆蓋率分開追蹤 |
| AI 生成的 code | 補償合理但錯誤的 code | 要求 BDD scenario、unit/contract tests，以及專注於 intent、edge cases 和安全/所有權邊界的人類審查 |
| 業務規則 / 演算法 | 捕捉通過範例但規則錯誤的情況 | 加上 property-based tests、invariant tests 或 table-driven edge cases |
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

對於效能測試，選擇能證明 release 風險的最小測試：

| 測試類型 | 必要時機 |
| --- | --- |
| Load test | 預期流量、job volume 或 batch size 已知且必須保持在預算內 |
| Stress test | 極限、飽和點或 degradation 行為未知 |
| Spike test | 突然的 burst、queue pressure、retries、cache misses 或 external-call fan-out 是合理的 |
| Soak test | 長時間運行的 memory、connection、cache、file-handle、queue 或 database drift 是合理的 |

不要只接受平均 latency 作為效能證據。記錄 percentile latency、throughput、error rate 和資源使用率，以及環境和資料集大小。

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

## Missing Information Gate（缺失資訊關卡）

在開發規劃或實作繼續之前，缺失資訊必須被明確處理：

| 缺失項目類型 | 必要行動 |
| --- | --- |
| 可從證據恢復 | 回填並引用證據來源 |
| Product intent 無法恢復 | 標記為 `unknown` / `open question`，向使用者提問，不要憑空創造 intent |
| 影響 BDD behavior、domain invariants、API/interface shape、error handling、security、storage 或 tests | 視為 blocker：在繼續實作前向使用者提問或要求證據 |
| 不改變行為或 contracts 的 nice-to-have 上下文 | 記錄為非阻塞的 open question，並說明為什麼不阻塞 |

不要在未解決的 blockers 下繼續開發。Agent 必須將缺失項目列為問題，等待答案或證據，然後在繼續前更新文件。

## Existing Project Documentation Backfill（既有專案文件回填）

當這個 workflow 被用於已經完全或大部分實作的專案時，先審查既有文件並回填任何缺失的開發文件。不要因為實作已經存在就跳過流程。

| 缺失文件 | 回填規則 |
| --- | --- |
| Product Brief | 只重建證據支持的內容：可見的目標、使用者/actors、範圍、non-goals、假設和限制。如果原始 intent 不可取得，將欄位標記為 `unknown` 或 `open question`；不要憑空創造業務理由 |
| Bounded Context Map | 從 code ownership、runtime boundaries、database tables、API groups、UI areas、queues、SDK/public APIs 和 deployment units 推斷模組 |
| BDD Behavior | **必須完成。** 從已實作的產品、tests、UI、API 行為和 logs 重建關鍵 happy paths、failure paths、權限、空狀態、edge cases 和跨 context 流程 |
| Domain Model Contract | 從 code、schemas、storage、UI states 和 tests 推斷 entities、value objects、commands、events、invariants 和 state transitions；將不確定的詞彙標記為 candidate |
| Architecture Contract | 記錄實際的 dependency direction、data ownership、side-effect boundaries、integrations、runtime/deployment shape 和已知違規 |
| API / Interface Contract | 提取實際的 request/response schemas、public methods、events、commands、auth/session behavior、versioning、compatibility、fixtures 和 consumers |
| Error Handling Contract | 回填觀察到的 error taxonomy、retry rules、user messages、logging/redaction behavior、security-sensitive failures 和 gaps |
| Test Plan | 將既有 tests 對應到 behavior/contracts，並列出未覆蓋的 BDD scenarios、invariants、contracts 和 integration paths 所需的 tests |

對於已實作優先的專案，也要恢復 delivery pipeline：

| Pipeline artifact | 回填規則 |
| --- | --- |
| Plan index / product radar | 將來源 product docs、PDFs、tickets、screenshots 或 legacy notes 對應到 modules、controllers、screens、commands 或 packages。標記已取消或已取代的需求 |
| Contract taxonomy | 列出哪些文件管轄 build/run、HTTP/API shape、auth/tenant/session、persistence、domain layering、frontend/backend integration、third-party integration、testing 和 documentation sync |
| Minimum doc sync matrix | 對每個變更類型，說明最少要更新的 docs/tests：API、permission、database、UI flow、generated client、vendor integration、CLI command、diagnostic rule、release setting |
| OpenAPI / schema / generated client | 驗證 generated consumer code 來自 source contract，而非手抄 endpoints 或 DTOs |
| Vendor / third-party integration | 區分 raw vendor docs 與 sanitized integration excerpts、request/response contracts、fixture examples、live-test gates 和 secret handling |
| Tooling / extension rule catalog | 對應 catalog order、rule IDs、diagnostics/commands、fixtures 和 tests；明確標記 process-only 或 non-enforceable rules |

既有專案的回填順序：

1. 盤點既有 docs、source folders、tests、schemas、API specs、fixtures、release notes 和觀察到的行為
2. 建立 documentation gap table，狀態為：`exists`、`partial`、`missing` 或 `unknown`
3. 當 product brief 缺失時先回填 BDD Behavior，因為已實作的行為是最強的可用 truth source
4. 從已完成的行為和實作證據回填 Domain Model、Architecture、API / Interface 和 Error Handling Contracts
5. 將未知的 product intent 與觀察到的行為分開標記。未知的 intent 不阻塞 BDD 完成
6. 如果 BDD 無法從可用證據完成，停止並要求缺失的行為、screen/API 範例、logs、test cases 或使用者決策，然後再繼續開發
7. 對任何缺乏覆蓋率的關鍵 BDD scenario 加上 tests 或 test TODOs

## Contract-First Rules（合約優先規則）

- BDD 描述行為；它不應鎖定 framework 或資料庫選擇
- Domain Model Contract 擁有 invariants、業務詞彙和狀態轉換
- Architecture Contract 擁有依賴方向、runtime boundaries、資料所有權和允許的整合路徑
- API Contract 擁有整合形狀：request、response、error、auth/session、versioning 和 compatibility
- Error Handling Contract 擁有 failure taxonomy、retry policy、user messaging、logging 和 security redaction
- Contract Governance 擁有文件優先順序、衝突處理、已取消/延後範圍和最小 linked updates
- 新需求必須在 code 開始前更新 planning docs、BDD、contracts、實作切片和 tests
- Bug 修復必須在 code 開始前識別預期 vs 實際行為以及 regression test
- 新的或 AI 生成的 code 必須用針對變更行為的 tests 驗證，而非僅整體專案覆蓋率
- 當一般範例無法證明規則時，使用 mutation、property-based、contract 或 database-backed tests
- Embedded 變更必須區分 datasheet/protocol truth、hardware context、driver/service/application 所有權、host-testable logic 和 target-only evidence
- 只有當共享 contracts 已版本化到足以進行 mock、stub 或 schema-first 工作時，才能平行進行實作
- 如果 contract 變更，在同一次變更中更新 BDD、實作、mocks 和 tests，或明確記錄為什麼不更新
- 如果 API/schema contract 變更，從 source contract 重新產生 typed clients 或 SDKs；不要手抄 routes、DTOs 或 operation names
- 如果第三方整合變更，更新 sanitized integration docs、fixtures、live-test gates 和 secret/redaction notes，但不要將私人 vendor 或帳戶細節複製到可重複使用的 guidance 中
- 對於已實作的專案，BDD 成為必要的行為恢復文件。Product Brief 可能包含 unknowns，但 BDD 必須從可觀察的產品行為和實作證據填寫
- 任何改變行為、contracts、所有權、錯誤處理、儲存、安全性或 tests 的缺失資訊都會阻塞開發，直到被回答或明確排除範圍

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

## Minimum Definition Of Ready（最低就緒定義）

在實作開始前，功能應具備：

- 包含範圍和 non-goals 的 Product brief
- Bounded Context 或模組拆分
- 關鍵行為的 BDD scenarios
- 核心 invariants 的 Domain Model Contract
- 依賴關係、所有權和 runtime boundaries 的 Architecture Contract
- 整合用的 API、event、command 或 public interface contract
- 預期失敗和 recovery behavior 的 Error Handling Contract
- 涵蓋 unit、behavior、contract 和 integration 層級的 Test plan
- 區分既有-regression 覆蓋率與變更/new-code 驗證的 Test strategy
- 當 latency、throughput、資源使用率、concurrency、啟動時間、背景任務、資料庫存取或外部呼叫量可能改變時的效能預算和測試類型
- 沒有影響實作行為或 contracts 的未解決 blocker questions

對於已實作的專案，「ready」表示缺失文件審查已完成，且 BDD 涵蓋已實作的關鍵行為，即使原始 product intent 仍部分未知。

## Minimum Definition Of Done（最低完成定義）

在出貨或合併前：

- Domain invariants 已測試
- Contract tests 對 provider 和 consumer 都通過
- Mocks/fixtures 符合最新 contract
- Integration test 至少涵蓋關鍵 happy path 和一個重要 failure path
- 效能敏感變更已記錄 load、stress、spike、soak 或 smoke-size 效能證據，對照 agreed budget
- 殘留的 unknowns 或延後的行為已在專案 repository 中記錄
