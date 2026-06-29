# App Development Guidance Execution Flow（開發指引執行流程）

本文件定義從分析觀察轉換為開發指引的執行流程。承接 ``skills/app-development-guidance/WORKFLOW.md``（已刪除）的內容，提取為 tool-neutral 的執行步驟。

> **遷移狀態**：此文件為新分層的 canonical source，`skills/app-development-guidance/WORKFLOW.md` 已刪除。新內容請直接寫入此文件。

Software-delivery 的 AI runtime gate 見 [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md)；本 workflow 保留執行順序、模板選擇與交付流程。

> **Incident path (2026-06)**: 對未知 UI incident，software-delivery 從 execution system 擴展為 **learning system** 決策鏈：delivery produces evidence → evidence decides change → change produces knowledge。Discover → Observe → Classify → Select Layer → Execute → Ship → Retrospective。見 [`incident-observation.md`](incident-observation.md)、[`ui-incident-governance-workflow.md`](ui-incident-governance-workflow.md)、[`change-retrospective.md`](change-retrospective.md)。

> ## Cognitive Slice 導航（thin index）
>
> 本檔是 software-delivery **execution-order lifecycle** 的 canonical 入口。Slice taxonomy 見 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7。依 task intent 載入對應段落 / slice，不需整份載入：
>
> | 認知階段（cognitive phase） | 本檔對應段落 | slice id | load_when |
> |---|---|---|---|
> | Intake（需求接收 / parity） | **已抽出** → [`intake.md`](intake.md)（含 Start From Evidence、Change Intake、Pre-build Interrogation、Requirements Cognition、Parity Gate、Product Brief Validation、Missing Information、Backfill） | `sd-intake` | 接收新需求 / 變更 / 重構意圖 |
> | Incident observe（UI incident observable） | **已抽出** → [`incident-observation.md`](incident-observation.md) | `sd-incident-observation` | 未知 UI incident；需 incident card；禁止 implementation-first |
> | Incident classify + layer（UI incident 治理） | **已抽出** → [`ui-incident-governance-workflow.md`](ui-incident-governance-workflow.md) + [`layer-ownership-matrix.md`](layer-ownership-matrix.md) | `sd-ui-incident-governance` | Navigation / Continuation / Recovery 未決；選 primary modification layer |
> | Change retrospective（Ship 後學習出口） | **已抽出** → [`change-retrospective.md`](change-retrospective.md) | `sd-change-retrospective` | Incident 路徑 Ship 後；promotion 三選一；禁止 direct canonical |
> | Test strategy（測試策略 / BDD） | **已抽出** → [`test-strategy.md`](test-strategy.md)（含 BDD Closure / Journey Specification / Docs-First Loop / Test Strategy Gate / Mutation / Test-First Ordering） | `sd-test-strategy` | 定義測試策略 / BDD 閉環 / Journey Specification |
> | UI contracts（Screen Mapping / Consumer / Screen / ViewModel） | **已抽出** → [`ui-contracts.md`](ui-contracts.md)（含 Screen Mapping、Consumer Contract、UI Behavior Contract、Screen Contract、Frontend ViewModel Contract、Accessibility Contract、Screen Traceability） | `sd-ui-contracts` | 前端、行動、CLI、SDK 或其他 consumer surface 需要平行實作或 AI 生成 UI / state / tests |
> | UI governance（UI compliance / design system / evidence） | **已抽出** → [`ui-governance.md`](ui-governance.md)（含 governance domain、render context、validation mechanism、evidence class、severity policy、project-local design-system boundary） | `sd-ui-governance` | UI / consumer surface 需要 design-system enforcement、accessibility evidence、responsive evidence、behavior pattern checks、visual baseline review、AI visual review scoping 或 UI compliance completion claim |
> | Implementation（執行核心） | [`implementation/`](implementation/README.md)（[`execution-modes.md`](implementation/execution-modes.md)）+ [`execution-flow.md`](execution-flow.md) §3 SDK 缺陷閉環、§4 同工作階段閉環 | `sd-implementation` | 實際進行程式碼變更；feature 被 structure block 時載入 execution-modes |
> | Validation（驗證 / 效能） | **已抽出** → [`validation.md`](validation.md)（含 Perf Gate + Validate + Journey Validation） | `sd-validation` | 驗證變更 / 效能關卡 / Journey Validation |
> | Closure（收尾 / 回饋） | **已抽出** → [`closure.md`](closure.md)（含 DoR / DoD / Feed Back Lessons） | `sd-closure` | 收尾、DoR/DoD 檢核、回饋可重用課程 |
> | Surgical caveats（diff 紀律） | **已抽出** → [`surgical-changes.md`](surgical-changes.md) | `sd-surgical-caveats`（`type: failure`） | 外科手術式小改、控制 diff 純度 |
>
> **Suppression 提示**：行為範例 [`examples/EXAMPLES.md`](examples/EXAMPLES.md)（`type: examples`）預設 **不載入**，僅在使用者明確要求範例或偵測到 ambiguity 時載入。evidence-only / 純分析任務不應載入本 execution-flow。
>
> > **Phase 2 進度**：已抽出 `sd-surgical-caveats`（[`surgical-changes.md`](surgical-changes.md)）、`sd-contracts`（[`contracts.md`](contracts.md)）、`sd-ui-contracts`（[`ui-contracts.md`](ui-contracts.md)）、`sd-ui-governance`（[`ui-governance.md`](ui-governance.md)）、`sd-closure`（[`closure.md`](closure.md)）、`sd-validation`（[`validation.md`](validation.md)）、`sd-test-strategy`（[`test-strategy.md`](test-strategy.md)，§2 + §4 子節 + development-process.md §BDD Closure/§Test Strategy Gate 跨檔同批）、`sd-intake`（[`intake.md`](intake.md)，§1 + §6 Backfill + development-process.md §Initial Doc Pack/§Product Brief Validation/§Change Intake/§Missing Information/§Existing Project Backfill 跨檔同批）並建立本導航。**`sd-implementation`**：execution mode 正文已落地 [`implementation/execution-modes.md`](implementation/execution-modes.md)（2026-06-29 Phase 1）；SDK / 同工作階段閉環仍留本檔 §3–§4。

## 1. 從證據開始（Start From Evidence）→ 已抽出為獨立 slice

從證據開始、變更接收（Change Intake）、Pre-build Interrogation Gate、Requirements Cognition Checkpoint、重構 / Replacement Parity Gate 已連同 development-process.md 的 Initial Documentation Pack / Product Brief Validation Gate / Change Intake Gate / Missing Information Gate / 既有專案文件回填，跨檔同批抽出為 focused slice **[`intake.md`](intake.md)**（`sd-intake`，`type: execution`，tags `requirements, parity, intake`）。canonical content 在 `intake.md`，此處不再保留正文以避免 dual source-of-truth。

接收新需求 / 變更 / bug / 重構意圖、需求認知盤點時載入該 slice；已有明確 contract 的純執行改動或 evidence-only 任務不需載入。

## 2. 文件優先 BDD 閉環（Docs-First BDD Closure Loop）→ 已抽出為獨立 slice

Docs-First BDD Closure Loop 連同 BDD Execution Closure、Journey Specification、Test Strategy Gate、測試策略定義、Test-First Ordering 已抽出為 focused slice **[`test-strategy.md`](test-strategy.md)**（`sd-test-strategy`，`type: execution`，tags `artifact-gate, test, bdd`）。canonical content 在 `test-strategy.md`，此處不再保留正文以避免 dual source-of-truth。

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

效能測試關卡連同 §7 驗證與 Journey Validation 已抽出為 focused slice **[`validation.md`](validation.md)**（`sd-validation`，`type: execution`，tags `artifact-gate, validation, performance`）。canonical content 在 `validation.md`，此處不再保留正文以避免 dual source-of-truth。

驗證變更 / 效能關卡時載入該 slice；尚未實作完成前不需載入。

## 6. 已實作專案的回填規則（Backfill Rules for Implemented Projects）→ 已抽出為獨立 slice

已實作專案的回填規則已連同 development-process.md §Existing Project Documentation Backfill，作為 backfill 條件子流程抽出至 [`intake.md`](intake.md) §Backfill（`sd-intake`，`tags: domain-specific,backfill`）。canonical content 在 `intake.md`，此處不再保留正文以避免 dual source-of-truth。僅在處理「已實作但文件缺失」的專案時載入。

## 7. 驗證（Validate）→ 已抽出為獨立 slice

驗證方法清單與「舊行為仍受保護 / 新程式碼已證明」分流連同 §5 效能測試關卡已抽出為 focused slice **[`validation.md`](validation.md)**（`sd-validation`）。canonical content 在 `validation.md`，此處不再保留正文以避免 dual source-of-truth。

## 8. 回饋可重複使用的課程（Feed Back Reusable Lessons）→ 已抽出為獨立 slice

回饋可重用課程連同 Definition of Ready / Definition of Done 收尾檢核，已抽出為 focused slice **[`closure.md`](closure.md)**（`sd-closure`，`type: execution`，tags `closure, handoff, extraction-to-intelligence`）。canonical content 在 `closure.md`，此處不再保留正文以避免 dual source-of-truth。

收尾、DoR/DoD 檢核、回饋可重用課程時載入該 slice；任務中段（intake / 實作 / 驗證進行中）不需載入。

## 9. 外科手術式修改規則（Surgical Changes Rules）→ 已抽出為獨立 slice

外科手術式修改規則（9.1 只改必須改的行 / 9.2 匹配既有 code style / 9.3 不加「順便」功能 / 9.4 只清理自己的 orphan / 9.5 驗證 diff 純度）已抽出為 focused slice **[`surgical-changes.md`](surgical-changes.md)**（`sd-surgical-caveats`，`type: failure`）。

修改既有程式碼、需控制 diff 純度 / orphan 時載入該 slice；大型新功能初始實作或 evidence-only 任務不需載入。具體範例見 [`examples/EXAMPLES.md`](examples/EXAMPLES.md) §3（預設 suppress，僅在明確要求範例時載入）。
