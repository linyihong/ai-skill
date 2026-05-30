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
> | Test strategy（測試策略 / BDD） | **已抽出** → [`test-strategy.md`](test-strategy.md)（含 BDD Closure / Docs-First Loop / Test Strategy Gate / Mutation / Test-First Ordering） | `sd-test-strategy` | 定義測試策略 / BDD 閉環 |
> | Implementation（執行核心） | §3 SDK 缺陷閉環、§4 同工作階段閉環 | `sd-implementation` | 實際進行程式碼變更 |
> | Validation（驗證 / 效能） | **已抽出** → [`validation.md`](validation.md)（含 Perf Gate + Validate） | `sd-validation` | 驗證變更 / 效能關卡 |
> | Closure（收尾 / 回饋） | **已抽出** → [`closure.md`](closure.md)（含 DoR / DoD / Feed Back Lessons） | `sd-closure` | 收尾、DoR/DoD 檢核、回饋可重用課程 |
> | Surgical caveats（diff 紀律） | **已抽出** → [`surgical-changes.md`](surgical-changes.md) | `sd-surgical-caveats`（`type: failure`） | 外科手術式小改、控制 diff 純度 |
>
> **Suppression 提示**：行為範例 [`examples/EXAMPLES.md`](examples/EXAMPLES.md)（`type: examples`）預設 **不載入**，僅在使用者明確要求範例或偵測到 ambiguity 時載入。evidence-only / 純分析任務不應載入本 execution-flow。
>
> > **Phase 2 進度**：已抽出 `sd-surgical-caveats`（[`surgical-changes.md`](surgical-changes.md)）、`sd-contracts`（[`contracts.md`](contracts.md)）、`sd-closure`（[`closure.md`](closure.md)）、`sd-validation`（[`validation.md`](validation.md)）、`sd-test-strategy`（[`test-strategy.md`](test-strategy.md)，§2 + §4 子節 + development-process.md §BDD Closure/§Test Strategy Gate 跨檔同批）並建立本導航。其餘 2 個 lifecycle phase（intake / implementation）的實體拆檔留待與 `development-process.md` 同批進行，見 plan Phase 2。routing-registry / execution-flow.yaml 將新 slice 納入 required source 的同步留待 Phase 3。

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

## 2. 文件優先 BDD 閉環（Docs-First BDD Closure Loop）→ 已抽出為獨立 slice

Docs-First BDD Closure Loop 連同 BDD Execution Closure、Test Strategy Gate、測試策略定義、Test-First Ordering 已抽出為 focused slice **[`test-strategy.md`](test-strategy.md)**（`sd-test-strategy`，`type: execution`，tags `artifact-gate, test, bdd`）。canonical content 在 `test-strategy.md`，此處不再保留正文以避免 dual source-of-truth。

定義測試策略 / BDD 閉環 / test-first ordering 時載入該 slice。

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

### 測試策略定義 / Test-First Ordering → 已抽出為獨立 slice

§4 的兩個子節「測試策略定義」（專案內部問題清單）與「Test-First Ordering」（framework/runtime/governance 升級強制順序）已抽出為 [`test-strategy.md`](test-strategy.md)（`sd-test-strategy`）。canonical content 在 test-strategy.md，此處不再保留正文以避免 dual source-of-truth。本節 §4「同工作階段閉環」父節（程式碼 + 持久文件的閉環紀律）保留於此，屬 `sd-implementation` 範圍。

## 5. 效能測試關卡（Performance Test Gate）→ 已抽出為獨立 slice

效能測試關卡連同 §7 驗證已抽出為 focused slice **[`validation.md`](validation.md)**（`sd-validation`，`type: execution`，tags `artifact-gate, validation, performance`）。canonical content 在 `validation.md`，此處不再保留正文以避免 dual source-of-truth。

驗證變更 / 效能關卡時載入該 slice；尚未實作完成前不需載入。

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

## 7. 驗證（Validate）→ 已抽出為獨立 slice

驗證方法清單與「舊行為仍受保護 / 新程式碼已證明」分流連同 §5 效能測試關卡已抽出為 focused slice **[`validation.md`](validation.md)**（`sd-validation`）。canonical content 在 `validation.md`，此處不再保留正文以避免 dual source-of-truth。

## 8. 回饋可重複使用的課程（Feed Back Reusable Lessons）→ 已抽出為獨立 slice

回饋可重用課程連同 Definition of Ready / Definition of Done 收尾檢核，已抽出為 focused slice **[`closure.md`](closure.md)**（`sd-closure`，`type: execution`，tags `closure, handoff, extraction-to-intelligence`）。canonical content 在 `closure.md`，此處不再保留正文以避免 dual source-of-truth。

收尾、DoR/DoD 檢核、回饋可重用課程時載入該 slice；任務中段（intake / 實作 / 驗證進行中）不需載入。

## 9. 外科手術式修改規則（Surgical Changes Rules）→ 已抽出為獨立 slice

外科手術式修改規則（9.1 只改必須改的行 / 9.2 匹配既有 code style / 9.3 不加「順便」功能 / 9.4 只清理自己的 orphan / 9.5 驗證 diff 純度）已抽出為 focused slice **[`surgical-changes.md`](surgical-changes.md)**（`sd-surgical-caveats`，`type: failure`）。

修改既有程式碼、需控制 diff 純度 / orphan 時載入該 slice；大型新功能初始實作或 evidence-only 任務不需載入。具體範例見 [`examples/EXAMPLES.md`](examples/EXAMPLES.md) §3（預設 suppress，僅在明確要求範例時載入）。
