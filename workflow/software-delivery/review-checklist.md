# Review Checklist（審查檢查清單）

本文件定義開發流程中的審查檢查項目。原承接 ``skills/app-development-guidance/CHECKLIST.md`` 的內容（已刪除），提取為 tool-neutral 的 workflow gates。

> **遷移狀態**：`skills/app-development-guidance/CHECKLIST.md` 已刪除。此文件為 canonical source，新內容請直接寫入此文件。

## 使用原則

1. **Keep checklists short enough to run during real development** — 檢查清單必須在實際開發中可執行。
2. **Checklist items must stay linked to implementation docs** — 檢查項目必須連結到它們要求審查者驗證的實作文件。
3. **When adding a check, update or verify implementation and control docs** — 新增檢查項目時，在同一變更中更新或驗證對應的 implementation 和 control 文件。

## 聚焦檢查清單

對於特定領域的聚焦檢查清單，使用：

- ``checklists/mobile-design-review.md``（已刪除）
- ``checklists/mobile-pr-review.md``（已刪除）
- ``checklists/mobile-release-review.md``（已刪除）
- ``checklists/api-security-review.md``（已刪除）
- ``checklists/contract-governance-review.md``（已刪除）
- ``checklists/embedded-firmware-review.md``（已刪除）

當檢查項目因新的 control 或 implementation pattern 而變更時，對應的 `controls/` 和 `implementation/` 文件必須在同一變更中更新或明確驗證。

---

## Change Intake（變更接收）

- 在程式碼之前已審查專案的企劃書、product brief、planning doc、issue、ticket、PRD、design note、BDD、API contract 或同等 artifact。
- 主要的 Product Brief 聲明已驗證或標記：目標、使用者、範圍、non-goals、假設、成功標準、限制、依賴和風險。
- 影響行為、合約、風險、所有權、測試、時程或發布關卡的未驗證 Product Brief 聲明是阻擋性問題。
- 變更已分類為新需求、bug 修復、重構 / replacement、強化或僅文件。
- 新需求在程式碼開始前更新規劃文件、BDD、受影響的合約、實作切片和測試計劃。
- Bug 修復記錄預期行為 vs 實際行為、重現/證據、受影響的 BDD scenario 或缺失 scenario、受影響的合約/錯誤和回歸測試計劃。
- 改變**可觀察**行為的 bug 修復在**同一個工作階段**中更新受影響的合約、BDD/可追溯性和整合或即時測試筆記（不僅在合併後）；延後需要明確的追蹤後續行動，而非無聲漂移。
- 重構已確認沒有行為或公開合約變更；若新入口會替代舊入口、舊腳本、API、資料流程、UI flow、runtime surface 或操作流程，已建立新舊能力 parity inventory。
- Parity inventory 逐列包含舊入口、現有能力、輸入、輸出 / 副作用、外部依賴、目標新入口、parity 狀態與測試 / fixture 證據；`deferred`、`not planned` 或 `tool-specific` 項目已說明為何不阻擋目前 phase。
- 可能影響延遲、吞吐量、資源使用、啟動、背景工作、資料庫存取、批次處理、快取、並發性或外部呼叫量的變更在程式碼之前定義效能預算和測試類型。
- 阻擋性問題在實作前已回答、有證據支持或明確排除範圍。

## Test Strategy（測試策略）

- 既有/舊有行為對受影響路徑有回歸覆蓋。
- 新或變更的行為在生產程式碼之前有 BDD scenarios。
- 新程式碼在可行時在實作之前有 failing unit、contract、property、integration 或 executable spec 覆蓋。
- 變更/新程式碼覆蓋率與總專案覆蓋率分開檢查。
- 突變測試、基於屬性的測試、不變量測試或負面案例涵蓋規則密集或安全敏感的邏輯。
- 當狀態重要時，資料庫、儲存庫、遷移或持久化行為使用 fixtures 或整合測試驗證。
- 當觀察到的成功訊號可能不同於真實狀態時，依 [`state-visibility-gap.md`](../../intelligence/engineering/execution/validation-reasoning/state-visibility-gap.md) 與 [`evidence-chain-validation.md`](../../intelligence/engineering/execution/validation-reasoning/evidence-chain-validation.md) 驗證完整 state propagation chain。
- Payment、email、external API、storage、queue、entitlement 或其他 proxy-prone path 不接受 API 200、adapter success、SMTP success 或 queue publish 作為 final proof；critical path 需要 independent observation。
- 效能敏感變更包含適合風險的 load、stress、spike、soak 或 smoke-size 效能證據。
- 效能證據報告 P95/P99 延遲、吞吐量、錯誤率和資源使用率；平均延遲不被視為足夠。
- 嵌入式或硬體支援的行為區分主機可重複測試與僅目標或硬體在迴路中的證據。
- AI 生成的程式碼接受針對規劃文件、BDD、合約、邊緣案例和安全/所有權邊界的人類審查。

## Performance Test Strategy（效能測試策略）

- 當使用者體驗、營運成本、可靠性、容量或外部依賴負載可能改變時，效能測試是發布關卡的一部分。
- 負載測試涵蓋預期的穩定需求，並將結果與同意的延遲、吞吐量、錯誤率和資源預算進行比較。
- 壓力測試識別飽和行為，並確認系統可預測地降級而非無聲失敗。
- 尖峰測試涵蓋突然的流量、作業、佇列、重試、快取或外部呼叫爆量。
- 浸泡測試涵蓋長時間運行的記憶體、連線、快取、檔案控制代碼、佇列、資料庫或資源漂移。
- CI/CD 至少包含關鍵路徑的小型效能 smoke 檢查，當完整套件對每次提交來說太昂貴時。
- 效能基準已版本化或記錄，以便審查者可以判斷變更是更快、更慢還是仍在預算內。

## Product To Contract Flow（產品到合約流程）

- Product brief 命名目標、使用者、範圍、non-goals、假設和限制。
- Product brief 聲明在驅動 BDD 或實作之前有證據、明確決策、驗證計劃或 `open question` 狀態。
- 成功標準可由 BDD、測試、指標、demo、發布檢查清單或手動證據證明。
- 效能成功標準在相關時包含明確預算：P95/P99 延遲、吞吐量、錯誤率和資源上限。
- 假設有負責人、驗證計劃和如果為假的影響。
- Bounded Contexts 或模組按領域職責和整合邊界拆分。
- 關鍵行為在實作之前撰寫為 BDD scenarios。
- Domain Model Contract 定義 entities、value objects、commands、events 和 invariants。
- Architecture Contract 定義依賴方向、資料所有權、runtime boundaries 和允許的整合。
- API、event、command 或 public interface contract 在平行實作之前定義。
- Error Handling Contract 定義 error taxonomy、retry rules、user messaging、logging 和 redaction。
- 嵌入式產品定義 datasheet/protocol truth、hardware context、driver/service/application 所有權、target 限制和 bring-up validation。
- 如果沒有前端/後端拆分，仍然命名 producer 和 consumer 角色。
- Mock APIs、fixtures、schemas 或 stubs 從最新 contract 生成。
- Unit、BDD、contract 和 integration test 責任在建置工作開始前分配。
- 缺失的行為、領域、API/介面、錯誤處理、安全性、儲存、所有權或測試要求在開發繼續前作為阻擋性問題提出。

## Existing Project Documentation Backfill（既有專案文件回填）

- 既有專案文件已盤點並標記為 `exists`、`partial`、`missing` 或 `unknown`。
- 缺失的 Product Brief 欄位僅從證據重建；不可取得的原始意圖標記為 `unknown` 或 `open question`。
- 回填的 Product Brief 聲明除非有 UI、API、程式碼、測試、日誌、fixtures、使用者決策或其他證據支持，否則不被視為已驗證。
- BDD behavior 對已實作的關鍵 happy paths、failure paths、權限、空狀態、edge cases 和跨 context 流程是完整的。
- BDD scenarios 引用來自 UI 行為、API 行為、程式碼路徑、測試、日誌、fixtures 或手動驗證的證據。
- Domain Model、Architecture、API / Interface 和 Error Handling Contracts 從觀察到的行為和實作證據回填。
- 每個關鍵 BDD scenario 對應到既有測試覆蓋率或必要的測試差距。
- 任何無法從證據回填且影響行為或合約的差距在實作繼續前提出。
- 文件優先順序已定義，以便 agent 在文件不一致時知道哪個 artifact 優先。
- Stable IDs 將 product/rule/operation/command/diagnostic 條目連結到 BDD、code refs、fixtures 和 tests。
- BDD scenarios 標記為 `automated`、`fixture-backed`、`manual-evidence`、`pending-runner` 或 `not-automatable`。
- 已取消、延後、僅流程、noop、僅手動和排除範圍的項目已明確標記。

## Contract Governance（合約治理）

- Governance/framework contract、product plan、BDD、domain/API/interface/error contracts、implementation 和 tests 有明確的優先順序。
- 當這些表面存在時，存在最小文件同步矩陣，涵蓋 API、permission、database、UI flow、generated client、vendor integration、CLI command、diagnostic rule 和 release setting 變更。
- OpenAPI/schema/API contract 變更重新生成 typed clients、SDKs、mocks、fixtures 或 schema packages。
- Vendor integration docs 分開 raw vendor sources 與 sanitized integration excerpts、fixtures、live-test gates 和 secret handling。
- Tooling/extension rule catalogs 將 stable IDs 映射到 diagnostics/commands、fixtures、tests 和明確的 non-enforceable entries。

## UI Governance Review（UI 治理審查）

- UI compliance claim 已連到 [`ui-governance.md`](ui-governance.md)，且沒有把 visual diff、screenshot 或 AI review 當成 governance domain。
- Governance domain、collection method、validation mechanism、evidence class、severity 已分開記錄；需要 artifact 時使用 [`templates/ui-governance-evidence-template.md`](templates/ui-governance-evidence-template.md)。
- Browser Review、screenshot capture、DOM snapshot、accessibility tree 或 human observation 被記錄為 evidence acquisition，不被當成 validator 或 compliance proof。
- Design-system claim 引用 project-local token / primitive / component policy；本 workflow 不要求全域 token scale。
- Accessibility 或 behavior pattern blocker 只基於客觀 expectation；主觀 visual taste 和 AI visual review 預設是 warning / research，除非有專案 opt-in、客觀 rubric、deterministic capture 和 review policy。
- Completion / DoD 沒有隱藏 unresolved UI governance blocker；deferred UI compliance scope 有 owner、理由和 follow-up。

## Reusable Guidance Boundary（可重複使用指引邊界）

- 應用 [`reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md) 中的全域規則。
- 可重複使用的技能指引說明通用故障模式、決策規則、擁有者層和驗證方法。
- 專案名稱、本機路徑、主機、端點字串、負載片段、樣本 ID、類別名稱、即時資料特異性和執行結果保留在專案儲存庫中，而非可重複使用的技能中。
- 從事件衍生的課程被拆分：技能中的通用方法；專案文件中的具體重現證據和 BDD/test 檔案名稱。
- 編輯此技能後，搜尋變更的技能資料夾中的專案特定字串，並移除或重新定位任何匹配項。

## Embedded / Hardware Product Review（嵌入式/硬體產品審查）

- Datasheet、vendor protocol、errata 和觀察到的偏差與產品行為分開記錄。
- Hardware context 記錄板子修訂、引腳、匯流排/UART/I2C/SPI/BLE/CAN 設定、緩衝區、時間和電源假設。
- 板子特定的接線和引腳選擇透過 context/config 注入，而非硬編碼為唯一的生產路徑。
- Driver、service、domain 和 application layers 有明確的所有權；raw bytes/registers 不會洩漏到產品行為程式碼中。
- BDD 使用領域行為和裝置狀態，而非 raw UART/register 操作。
- Protocol fixtures 包含正面範例、無效長度/形狀、重新同步或 checksum 失敗，以及 spec 中的邊界值。
- 主機可重複測試在依賴 bench-only 檢查之前涵蓋 parsing、domain invariants、command/API contracts 和 error mapping。
- 硬體在迴路中或手動啟動記錄韌體版本、板子/接線、測試命令、日誌、測量證據和已知限制。
- 安全相關行為定義 fail-safe state、timeout、debounce/cooldown、retry 和 recovery rules。
- 發布關卡涵蓋目標建置、燒錄路徑、配置預設值、機密、除錯日誌、校準/版本備註和回滾或復原路徑。

## API And Transport（API 與傳輸）

- 敏感流程僅使用 HTTPS；明文流量在發布建置中禁用。
- 對高風險應用程式考慮憑證固定，並附有輪換和事件計劃。
- 後端授權不信任僅客戶端的標誌、角色、價格、餘額或功能閘門。
- 對重放敏感的請求具有伺服器端時間戳、nonce、冪等性或風險檢查。
- 請求簽名（如果使用）簽署正確的欄位，且不依賴靜態客戶端機密。
- 錯誤回應避免洩漏堆疊追蹤、內部主機名稱或敏感業務規則。
- 生成的客戶端和 SDK 來自當前 API/schema contract；消費者不手抄 routes、DTOs 或 response envelopes。
- 第三方 webhooks/callbacks 驗證簽名、時間戳、body binding、冪等性和重放行為。

## Auth, Tokens, And Sessions（驗證、令牌與工作階段）

- 存取令牌有範圍、時限且可撤銷。
- 重新整理流程有速率限制且綁定到帳戶/工作階段上下文。
- 登出、密碼變更和風險事件使相關工作階段失效。
- 令牌不會寫入除錯日誌、崩潰報告、分析或螢幕截圖。
- 工作階段識別碼不用作長期裝置識別碼。
- 依身份、權限、tenant、ownership、會員或 feature flag 改變的 UI/API/SSR 行為，以真實身份材料走產品路徑驗證 protected resource、persisted state 與 user-observable state 一致。
- 忘記密碼、寄信、播放權限、收藏/追蹤、會員狀態或 entitlement 變更若產生 side effect，完成證據包含 live system proof 與必要的 DB/external/readback 驗證。

## Local Storage（本地儲存）

- 機密僅在需要時儲存，並在可能的情況下使用平台支援的安全儲存。
- 快取檔案、SQLite、共享偏好設定和下載的媒體會審查敏感內容。
- 備份和螢幕截圖根據產品風險進行配置。
- 離線資料有到期日、加密計劃或明確的業務理由。

## Flutter And Android Build（Flutter 與 Android 建置）

- 發布建置禁用除錯標誌、測試端點、開發選單和詳細的網路日誌。
- 在與崩潰符號化和支援需求相容的地方啟用混淆/最小化。
- 除非診斷需要，否則移除原生符號。
- Flutter platform channels 在沒有伺服器端或 OS 級別檢查的情況下不暴露特權操作。
- 第三方 SDK 會審查權限、遙測和資料收集。

## Logging And Telemetry（日誌與遙測）

- 日誌清理令牌、cookies、授權標頭、裝置 ID 和個人資料。
- 崩潰報告和分析事件避免原始 request/response 負載。
- 除錯日誌與發布建置隔離。
- 安全相關的故障是可觀察的，而不暴露機密。

## Anti-Tamper And Risk Signals（反篡改與風險訊號）

- Root/jailbreak/emulator/hook 檢測被視為風險訊號，而非唯一的存取控制決策。
- 伺服器端風險評分可以容忍誤報和漏報。
- 關鍵操作仍然需要後端授權和濫用控制。
- 應用程式避免儲存成為永久繞過目標的靜態機密。

## Release Gate（發布關卡）

- 審查者可以指向測試、建置檢查或記錄的證據，證明每個必要控制存在。
- 若 release claim 涉及 state visibility gap，審查者可以指向 evidence chain 中每個必要 segment 的證據，且 claim scope 不超過 evidence scope。
- 沒有影響行為、合約、錯誤、安全性、儲存、所有權或測試的未解決阻擋性問題。
- 已知的殘留風險記錄在專案儲存庫中。
- 可重複使用的課程僅在清理後才提升到此技能中。

---

## 與其他層的關係

- `workflow/software-delivery/execution-flow.md` 提供執行流程，本文件提供流程中的審查門檻。
- `analysis/development-guidance/controls-catalog.md` 提供檢查清單引用的控制原則。
- `skills/app-development-guidance/CHECKLIST.md` 是原始來源，已刪除。內容已由本文件承接。
