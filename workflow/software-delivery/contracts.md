# Contract 治理 Slice（Required Contracts / Governance / Traceability / Contract-First）

> **Cognitive Slice**：`sd-contracts`（從 [`development-process.md`](development-process.md) 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-contracts` |
| `purpose` | 建立與治理 contract（domain / architecture / API / error / UI / consumer / 硬體）、文件優先順序衝突處理與雙向可追溯性 |
| `type` | `execution` |
| `tags` | artifact-gate, contract, traceability |
| `load_when` | 需建立 / 治理 contract 與可追溯性、處理文件不一致衝突 |
| `do_not_load_when` | 無 contract 異動的小改、evidence-only / 純分析任務 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「要產哪些 contract、衝突時哪個 artifact 優先、如何追溯」的 ordering / gate；通過 workflow membership test，不承載 evidence 取得方法（非 analysis），不論證長期模式（非 intelligence） |
| `canonical_source` | 本檔（原 `development-process.md` §Required Contracts / Contract Governance Gate / Traceability Gate / Contract-First Rules） |
| `dependencies` | `sd-intake`（先完成 change intake / brief 驗證）、`sd-ui-contracts`（有 consumer surface 時）、[`templates/contract-template.md`](templates/contract-template.md) |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Phase 4 Scenario A（execution-only：無 contract 異動時本 slice 應 **不** 載入） |

## Required Contracts（必要合約）

不要假設每個專案都有前端/後端拆分。選擇符合架構的 contracts：

| 應用形狀 | 優先 contract |
| --- | --- |
| 僅後端/API | Domain Model Contract、API Contract、contract tests、與 consumer 的 integration tests |
| 前端 + 後端 | API Contract、Domain Model Contract、BDD scenarios、Screen Mapping、Consumer Contract、UI Behavior / Screen / ViewModel Contract |
| 僅前端應用 | Screen Mapping、UI Behavior Contract、Screen Contract、Frontend ViewModel Contract、local state/domain contract、mocked API/schema contract |
| 行動應用 | Screen/flow behavior、ViewModel Contract、local storage/session contract、有遠端服務時加上 Consumer / API Contract |
| CLI / 桌面 / 工具 | Command contract、input/output schema、consumer needs、domain model、fixture-based tests |
| Library / SDK | Public API contract、type/schema contract、consumer usage contract、examples、compatibility tests |
| 事件驅動 / worker | Event schema、command/event contract、idempotency 和 retry behavior |
| Embedded / firmware / 硬體產品 | Datasheet 或 protocol contract、hardware context contract、driver/service/application 邊界、BDD、host fixtures、hardware-in-loop checks |
| 靜態分析 / IDE extension / 開發者工具 | Rule catalog、diagnostic 或 command contract、pure kernel/adapter 邊界、fixture pairs、editor/CLI integration tests |

## Contract Governance Gate（合約治理關卡）

每個有多份文件的專案必須定義當文件不一致時哪個 artifact 優先。除非專案有更強的本地規則，否則使用這個預設優先順序：

1. Governance / framework contract：repository 層級的 invariants、必要更新規則、依賴方向、命名、build/run 限制
2. Product plan / accepted brief：product intent、範圍、non-goals、已取消的需求、業務語言
3. BDD behavior：可觀察的使用者/系統行為和 acceptance criteria
4. Domain、architecture、API/interface、consumer/UI、error handling、hardware 或 command contracts
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
| Screen / UI action -> Screen Mapping / Consumer / API / ViewModel Contract | 顯示 UI behavior、資料需求、provider operation、table/store ownership 與 display derivation 如何被支援 |

Stable IDs 可以是 feature IDs、rule IDs、operation IDs、route names、command names、diagnostic codes、event names 或 scenario tags。如果一個行為被有意識地記錄但未實作，標記為 `TBD`、`noop`、`not enforceable by tool`、`manual-only` 或 `out of scope`，並附上原因和負責人。

對 `.feature` / Gherkin 文件，traceability 不是附加說明。每個 scenario / scenario outline 都必須至少有一個 test ref 與一個 code / contract / schema / adapter ref；尚未實作時，test ref 指向 `pending-runner` / `todo` 測試位置，code ref 指向預期 owner path。缺少任一方向時，該 BDD artifact 只能算草稿，不能作為完成的 acceptance contract。

## Contract-First Rules（合約優先規則）

- BDD 描述行為；它不應鎖定 framework 或資料庫選擇
- Domain Model Contract 擁有 invariants、業務詞彙和狀態轉換
- Architecture Contract 擁有依賴方向、runtime boundaries、資料所有權和允許的整合路徑
- API Contract 擁有整合形狀：request、response、error、auth/session、versioning 和 compatibility
- Screen Mapping 擁有 BDD scenario 到 screens、APIs、tables / stores 與 validation target 的對照；應在 BDD-lite 後建立，並在 provider/consumer 平行實作前保持可追溯
- Consumer Contract 擁有 consumer needs、freshness、loading、empty/error behavior、permissions 和 observability；應在 API finalization 前完成或同步審查
- UI Behavior / Screen / ViewModel Contract 擁有 screen states、actions、validation、feedback、navigation、accessibility 和 API/domain data 到 UI display model 的 derivation；細節見 [`ui-contracts.md`](ui-contracts.md)
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
- 如果 Consumer / UI contract 變更，在同一次變更中更新 API / error contract、mocks、fixtures、view model mapper、screen behavior tests 或明確記錄 deferred scope
- 如果第三方整合變更，更新 sanitized integration docs、fixtures、live-test gates 和 secret/redaction notes，但不要將私人 vendor 或帳戶細節複製到可重複使用的 guidance 中
- 對於已實作的專案，BDD 成為必要的行為恢復文件。Product Brief 可能包含 unknowns，但 BDD 必須從可觀察的產品行為和實作證據填寫
- 任何改變行為、contracts、所有權、錯誤處理、儲存、安全性或 tests 的缺失資訊都會阻塞開發，直到被回答或明確排除範圍
