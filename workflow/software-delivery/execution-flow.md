# App Development Guidance Execution Flow（開發指引執行流程）

本文件定義從分析觀察轉換為開發指引的執行流程。承接 ``skills/app-development-guidance/WORKFLOW.md``（已刪除）的內容，提取為 tool-neutral 的執行步驟。

> **遷移狀態**：此文件為新分層的 canonical source，`skills/app-development-guidance/WORKFLOW.md` 已刪除。新內容請直接寫入此文件。

Software-delivery 的 AI runtime gate 見 [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md)；本 workflow 保留執行順序、模板選擇與交付流程。

> ## Cognitive Slice 導航（thin index）
>
> 本檔是 software-delivery **execution-order lifecycle** 的 canonical 入口。Slice taxonomy 見 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7。依 task intent 載入對應段落 / slice，不需整份載入：
>
> | 認知階段（cognitive phase） | 本檔對應段落 | slice id | load_when |
> |---|---|---|---|
> | Intake（需求接收 / parity） | §1（Start From Evidence、Change Intake、Pre-build Interrogation、Requirements Cognition、Parity Gate）、§6 Backfill | `sd-intake` | 接收新需求 / 變更 / 重構意圖 |
> | Test strategy（測試策略 / BDD） | §2 Docs-First BDD、§4 測試策略定義、Test-First Ordering | `sd-test-strategy` | 定義測試策略 / BDD 閉環 |
> | Implementation（執行核心） | §3 SDK 缺陷閉環、§4 同工作階段閉環 | `sd-implementation` | 實際進行程式碼變更 |
> | Validation（驗證 / 效能） | §5 Perf Gate、§7 Validate | `sd-validation` | 驗證變更 / 效能關卡 |
> | Closure（收尾 / 回饋） | §8 Feed Back Reusable Lessons | `sd-closure` | 收尾、回饋可重用課程 |
> | Surgical caveats（diff 紀律） | **已抽出** → [`surgical-changes.md`](surgical-changes.md) | `sd-surgical-caveats`（`type: failure`） | 外科手術式小改、控制 diff 純度 |
>
> **Suppression 提示**：行為範例 [`examples/EXAMPLES.md`](examples/EXAMPLES.md)（`type: examples`）預設 **不載入**，僅在使用者明確要求範例或偵測到 ambiguity 時載入。evidence-only / 純分析任務不應載入本 execution-flow。
>
> > **Phase 2 進度**：本輪（execution-flow.md 先行）已抽出 `sd-surgical-caveats`（`surgical-changes.md`）並建立本導航。其餘 5 個 lifecycle phase 的實體拆檔刻意與 `development-process.md` 同批進行（兩檔的 intake / contracts / test-strategy / closure 內容重疊，分批拆會造成 dual source-of-truth），見 plan Phase 2。routing-registry / execution-flow.yaml 將 `surgical-changes.md` 納入 required source 的同步留待 Phase 3。

## 1. 從證據開始（Start From Evidence）

記錄可重複使用的觀察：

- 這個工作單元的目標是什麼、採取了什麼行動、以及如何驗證？如果驗證不可執行，引用參考來源和推理邊界。
- 觀察到了什麼行為？
- 哪一層暴露了它：客戶端程式碼、傳輸層、API 合約、儲存層、日誌、建置配置、韌體、硬體上下文、協定或執行時期行為？
- 問題是已確認、可疑還是僅為風險模式？

不要將目標特定的端點、令牌、機密、裝置 ID、原始使用者資料或專案事件細節複製到這個技能中。在將觀察提升為可重複使用的指引之前，應用 [`enforcement/reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md)。

### 變更接收（Change Intake）

在程式碼工作之前，執行變更接收：

| 問題 | 必要行動 |
| --- | --- |
| 存在什麼規劃產出？ | 閱讀企劃書、product brief、規劃文件、issue、ticket、PRD、設計備註、BDD、API 合約或同等文件 |
| Product brief 本身是否已驗證？ | 根據證據或明確決策檢查目標、使用者、範圍、non-goals、假設、限制、依賴、風險和成功標準。將每個主要聲明標記為 `validated`（已驗證）、`assumption`（假設）、`open question`（開放問題）、`scoped out`（排除範圍）或 `invalidated`（無效） |
| 這是新需求還是行為變更？ | 先執行 [`requirements/`](requirements/README.md) stage：product-impact discovery、behavior-driven discovery、acceptance definition、ambiguity resolution；再更新規劃文件、BDD、合約、實作切片和測試 |
| 這是 bug 修復？ | 確認預期行為 vs 實際行為、重現/證據、受影響的 BDD 或缺失 scenario、受影響的合約/錯誤和回歸測試 |
| 這是重構？ | 先分類是純內部重構、架構重組、平台遷移、工具替換或舊系統 replacement。若會替代既有功能、腳本、API、資料流程、UI flow、runtime surface 或操作流程，必須在實作前建立新舊能力 parity inventory：舊入口、現有功能、輸入、輸出 / 副作用、外部依賴、目標新入口、parity 狀態、測試 / fixture 證據與 deferred / not planned 理由。只有純內部重構且不改變 observable behavior、public contract 或操作能力時，才可只確認沒有行為變更。 |
| 這是強化？ | 確認威脅/故障模式、擁有者層、驗證和連結的檢查清單/控制更新 |
| 這是架構決策或 domain model 變更？ | 先確認 requirements stage 已有 behavior boundary / acceptance criteria / ambiguity disposition，再執行 [`architecture/architecture-fit-analysis.md`](architecture/architecture-fit-analysis.md)，確認 CRUD / DDD Lite / Full DDD / event-driven / microservices 的 fit evidence；不得預設套用 DDD、CQRS 或 event sourcing |
| 這個變更是否影響延遲、吞吐量、資源使用、啟動、背景工作、資料庫存取、批次處理或外部呼叫量？ | 在程式碼之前定義效能預算和必要的效能測試類型。不要依賴「功能正確」作為變更可發布的證明 |
| 這個變更是否與現有文件衝突？ | 應用文件優先順序：治理/框架合約、產品計劃、BDD、合約、實作、測試。更新擁有文件，而不是僅默默修正程式碼 |

**Contract / test / implementation mismatch escalation：** 若 product brief、owner contract、BDD、測試、implementation 或使用者指正彼此衝突，不要繼續局部 patch。依 `metadata/recovery/domain-policies.yaml` 的 `software-delivery` policy 進入 recovery：重讀本 workflow、artifact gates、development process、implementation plan template、linked-updates 與 dependency-reading；寫出舊假設、反證、owner contract、行為規格或 BDD、implementation surface、validation gate 與 linked updates。未完成新 execution graph 前，不可用「測試綠了」宣稱完成，也不可把 implementation 當作唯一 source-of-truth。

如果不存在規劃產出且請求會改變行為，在實作之前建立輕量的變更簡報並詢問阻擋性問題。

如果 product brief 存在但包含影響行為、合約、風險、測試、所有權、時程或發布關卡的未驗證聲明，在實作之前將這些聲明視為阻擋項。對於純規劃答案，引用參考來源或推理邊界，而不是假裝 brief 已驗證。

> **輸出模板**：Change Intake 完成後，使用 [`templates/change-brief-template.md`](templates/change-brief-template.md) 記錄變更簡報。

### Pre-build Interrogation Gate

在 change intake 之後、implementation plan 或 framework migration 之前，讀取 [`requirements/pre-build-interrogation.md`](requirements/pre-build-interrogation.md)。若請求會變成 plan、code、workflow、governance、runtime、validation、schema、generated artifact 或 tool adapter 改動，必須先記錄：

- Goal、scope、non-goals 與 expected behavior / expected framework outcome。
- Acceptance criteria 與 validation target。
- Framework discovery：canonical source、owner layer、runtime projection、mirror / cache / generated output、compiler 與 linked updates。
- Duplication risk：是否會產生第二份 rule body、第二條 activation path、stale projection 或 ambiguous source-of-truth。
- Unknown disposition：`blocker_question`、`safe_assumption`、`scoped_out` 或 `invalidated`。

若仍有會影響 behavior、contract、runtime surface、source-of-truth、validation 或安全性的 `blocker_question`，不得產生 implementation plan；先向使用者提問或停在 planning。

### Requirements Cognition Checkpoint

在進入 architecture 或 implementation 前，若任務涉及 observable behavior，讀取 [`requirements/`](requirements/README.md)：

- Product-impact discovery：Impact Map、Customer Journey Map、cross-check decision。
- Behavior-driven discovery：actor intent、behavior boundary、shared language。
- Acceptance definition：acceptance criteria、validation target、regression scope。
- Ambiguity resolution：將不確定項標成 `assumption`、`open question`、`scoped out` 或 `invalidated`。

沒有 validation target 的 acceptance criteria 不能作為完成宣告基線；requirement contradiction 或 stale acceptance criteria 需要先重建 source-of-truth。

### 重構 / Replacement Parity Gate

當變更目標是重構、遷移、改寫、替換舊工具、拆分架構、搬移 runtime surface 或建立新入口取代舊入口時，不能只寫新設計。實作前必須先產出 parity inventory，讓 reviewer 能逐項確認舊能力沒有遺漏。

最低欄位：

| 欄位 | 必填內容 |
| --- | --- |
| 舊入口 | 舊 API、command、script、UI flow、資料表、job、hook、runtime surface 或文件入口。 |
| 現有能力 | 舊入口目前支援的行為、flags、輸入、輸出、錯誤模式與邊界條件。 |
| 副作用 | 寫檔、寫 DB、發送請求、commit / push、生成 artifact、同步 mirror、修改使用者設定或其他狀態變更。 |
| 外部依賴 | runtime、shell、binary、服務、權限、平台假設、環境變數與 credentials boundary。 |
| 新入口 | 對應的新 API、command、module、adapter 或 replacement surface。 |
| Parity 狀態 | `covered`、`wrapper first`、`native target`、`deferred`、`not planned` 或 `tool-specific`，並說明原因。 |
| 驗證證據 | BDD scenario、contract test、fixture、golden output、migration assertion、manual review checklist 或明確的 blocker。 |

Blocking rule：任何舊入口若狀態為 `deferred`、`not planned` 或 `tool-specific`，必須寫明為何不阻擋目前 release / phase；任何會產生副作用的舊入口，必須有 dry-run、fake-root、fixture 或等效隔離測試。缺少 parity inventory 時，不得開始 replacement implementation，也不得宣稱新功能已覆蓋舊功能。

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

> **輸出模板**：BDD Closure Loop 完成後，使用 [`templates/bdd-scenario-template.md`](templates/bdd-scenario-template.md) 記錄行為場景、acceptance criteria、validation target 與 traceability。

## 3. SDK 缺陷閉環（即時重現 + BDD 回歸）

當報告涉及**針對生產環境或供應商服務的 SDK 可觀察行為**時使用此流程。目標：**使用整合者使用的相同程式碼路徑進行驗證**，然後**在文件和測試中鎖定通用行為**——而不僅是聊天結論或一次性除錯筆記。

| 步驟 | 行動 |
| --- | --- |
| 1 | **重現**：使用應用程式的即時或整合測試工具，僅呼叫**支援的 SDK 公開表面**；除非專案明確將其定義為測試中的表面，否則不要注入捕獲的授權資料或手動建構的 HTTP |
| 2 | **記錄**：在專案的整合測試筆記或同等可追溯性產出中記錄通過/失敗、不穩定性、環境限制和已清理的失敗類別。捕獲原因類別，而非私人主機名稱、令牌、使用者資料或一次性樣本 ID |
| 3 | 如果行為錯誤或必須記錄為合約：更新專案的人類可讀行為規格、鏡像可執行行為測試，並添加或擴展 fixture 支援和/或即時回歸覆蓋 |
| 4 | 應用**同工作階段閉環**：擁有合約和連結更新與程式碼/文件在同一批次中移動 |

穩定的組合、映射、解碼和錯誤語義通常屬於**fixture 支援的行為測試**；遠端服務 availability 和邊緣保護可能保持**僅即時**，並記錄不穩定性和備用驗證。

## 4. 同工作階段閉環（程式碼 + 持久文件）

實作工作通常在**測試通過**後停滯。這使得合約、BDD 和整合筆記描述過時的行為——特別是當任務被框架化為「修復 bug」或「調整 SDK」而未提及文件時。

| 根本原因 | 緩解措施 |
| --- | --- |
| 任務被框架化為僅程式碼 | 仍然分類變更；如果**可觀察**行為改變了，連結的規格必須在**同一批次**中與生產程式碼一起移動 |
| 使用者訊息省略了「更新文件」 | 技能仍然適用：改變語義的 bug 修復在治理目的上**是**行為變更 |
| 完成定義 = 綠色 CI | 擴展 DoD：擁有合約/BDD 行已更新，或在變更簡報中記錄明確的**範圍化**文件債務 ticket ID |
| 專案連結更新僅存在於儲存庫 README 中 | 在觸及所列表面時遵循該矩陣；agent 在編輯這些套件時應讀取它 |

**在標記完成之前：** 驗證持久文件符合新的執行時期真相——至少是適用的專案架構/API/錯誤/測試合約、受影響流程的 Gherkin 或可執行規格，以及當公開語義改變時的即時測試或整合筆記。

> **輸出模板**：Same-Session Closure 完成後，使用 [`templates/implementation-plan-template.md`](templates/implementation-plan-template.md) 記錄實作計畫。

### 測試策略定義

然後定義測試策略：

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

### Test-First Ordering（Framework / Runtime / Governance 升級強制）

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

## 5. 效能測試關卡（Performance Test Gate）

當變更可能影響回應時間、吞吐量、資源使用、啟動工作、背景處理、資料庫存取、外部 API 扇出、快取、批次處理或並發性時，功能正確性是不夠的。當使用者體驗、成本、可靠性或營運容量依賴於它時，將效能視為發布合約的一部分。

| 測試類型 | 使用時機 | 證明 |
| --- | --- | --- |
| 負載測試 | 預期流量或正常批次量已知 | 系統在正常需求下保持在延遲、吞吐量、錯誤率和資源預算內 |
| 壓力測試 | 容量限制或擴展行為未知 | 系統可預測地降級，並在生產之前暴露第一個瓶頸 |
| 尖峰測試 | 流量可能突然跳升、佇列可能爆量、或 AI 生成的變更改變了呼叫量 | 自動擴展、佇列、速率限制、快取和重試行為能承受突然的需求變化 |
| 浸泡測試 | 記憶體、連線、快取、檔案控制代碼、佇列或資料庫漂移可能隨時間出現 | 長時間運行的行為保持穩定，不會洩漏資源或逐漸降級 |

最低指標：

- 延遲：使用者可見或合約可見操作的 P95 和 P99；平均值僅為支援性上下文。
- 吞吐量：相關表面的每秒/分鐘請求、作業、訊息或操作數。
- 錯誤率：超時、5xx、重試耗盡、佇列失敗或領域特定的失敗預算。
- 資源使用率：相關時的 CPU、記憶體、磁碟、網路、資料庫連線、佇列深度、執行緒/任務計數和外部呼叫量。

CI/CD 可以從小的 smoke 級別效能檢查開始。較大的負載、壓力、尖峰或浸泡套件可以夜間運行、預發布或按需運行，但其觸發條件、擁有者、預算和證據位置必須記錄。

## 6. 已實作專案的回填規則（Backfill Rules for Implemented Projects）

如果專案已實作且文件缺失，在提出新的指引之前先進行文件差距審計：

| 文件 | 回填要求 |
| --- | --- |
| Product Brief | 僅回填證據支援的目標、使用者、範圍、限制和假設；將不可取得的意圖標記為 `unknown` 或 `open question` |
| BDD 行為 | **必須完成。** 從可觀察的 UI、API、程式碼、測試、日誌、fixtures 和手動驗證中完成 |
| 合約 | 從已實作的行為和證據回填領域模型、架構、API/介面、錯誤處理和測試計劃 |
| 嵌入式/硬體證據 | 從程式碼、日誌、接線筆記和測試中回填 datasheet/協定參考、硬體上下文、驅動程式/服務/應用程式邊界、主機 fixtures 和啟動證據 |
| 可追溯性 | 將產品/規則 ID 連結到 BDD，BDD 連結到程式碼引用，BDD 連結到測試，API/命令/診斷合約連結到 fixtures，以及生成的客戶端連結到來源合約 |

不要讓缺失的 product brief 阻擋已實作產品的 BDD 回填。

對於先實作優先的專案，也要恢復交付管線：來源產品文件或計劃雷達、文件優先順序、最小文件同步矩陣、OpenAPI/schema/codegen 流程、供應商整合摘錄，以及明確的已取消/延後/排除範圍的決策。

影響行為、領域不變量、API/介面形狀、錯誤處理、安全性、儲存、所有權或測試的缺失資訊是阻擋項。向使用者提問或要求證據，用答案更新文件，然後才繼續開發規劃或實作。非阻擋性的未知項必須標記為什麼它們不改變行為或合約。

對於嵌入式或硬體支援的產品，缺失的 datasheet/協定真相、電氣介面、引腳/匯流排映射、硬體上下文所有權、時間/並發限制、安全行為、fixture 來源或目標驗證方法也是阻擋項，除非明確排除範圍。

## 7. 驗證（Validate）

使用至少一種驗證方法：

- 單元或整合測試。
- 發布檢查清單項目。
- 靜態掃描或建置斷言。
- 附證據的手動審查。
- 執行時期或後端遙測查詢。
- 嵌入式/硬體行為的主機端 fixture 測試、模擬器測試、bench 日誌或硬體在迴路中運行。
- 提供者/消費者合約測試、生成的客戶端編譯檢查、fixture 對、診斷快照或閘控即時整合測試。

在驗證實作之前，確認沒有影響行為、合約、錯誤處理、安全性、儲存、所有權或測試的未解決阻擋性問題。

驗證應區分「舊行為仍受保護」與「新程式碼已證明」。優先使用 BDD/TDD 加上變更程式碼測試；當範例單獨無法證明規則時，添加突變、基於屬性、合約、資料庫支援、生成的客戶端、fixture 支援、主機端 fixture 或硬體在迴路中的測試。

> **輸出模板**：Validate 完成後，使用 [`templates/review-report-template.md`](templates/review-report-template.md) 記錄審查報告。

## 8. 回饋可重複使用的課程（Feed Back Reusable Lessons）

如果一個課程超越了一個產品：

1. 在匹配的 `feedback_history/<category>/` 或跨領域的 `feedback_history/common/` 下添加一個檔案。
2. 連結共享規則而不是複製它們。
3. 將已驗證的指引提升到結構化資料夾、檢查清單或此工作流程中。

如果課程來自 APK 分析，將分析方法保留在 `analysis/apk/` 或 `workflow/apk-analysis/` 中，將開發行動保留在此處。

## 9. 外科手術式修改規則（Surgical Changes Rules）→ 已抽出為獨立 slice

外科手術式修改規則（9.1 只改必須改的行 / 9.2 匹配既有 code style / 9.3 不加「順便」功能 / 9.4 只清理自己的 orphan / 9.5 驗證 diff 純度）已抽出為 focused slice **[`surgical-changes.md`](surgical-changes.md)**（`sd-surgical-caveats`，`type: failure`）。

修改既有程式碼、需控制 diff 純度 / orphan 時載入該 slice；大型新功能初始實作或 evidence-only 任務不需載入。具體範例見 [`examples/EXAMPLES.md`](examples/EXAMPLES.md) §3（預設 suppress，僅在明確要求範例時載入）。
