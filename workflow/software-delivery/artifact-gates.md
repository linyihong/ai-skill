# App Development Guidance Artifact Gates（開發指引產出規範）

本文件定義開發指引文件的產出規範與品質門檻。承接 ``skills/app-development-guidance/DOCUMENTATION.md``（已刪除）的內容，提取為 tool-neutral 的 artifact gates。

> **遷移狀態**：此文件為新分層的 canonical source，`skills/app-development-guidance/DOCUMENTATION.md` 已刪除。新內容請直接寫入此文件。

Artifact completeness 與 same-session closure 的治理 gate 見 [`software-delivery-governance.md`](../../governance/ai-runtime-governance/software-delivery-governance.md)；本檔保留 reusable note 與 artifact 品質門檻。

## 1. 可重複使用的筆記結構（Reusable Note Structure）

當將 APK 分析課程、實作觀察、嵌入式/韌體觀察、硬體產品課程或設計審查轉換為可重複使用的開發指引時，使用此結構：

```markdown
### 簡短標題

狀態：candidate | validated | promoted | deprecated | experimental

#### 觀察到的風險（Observed Risk）

觀察到的模式，不含目標特定機密。

#### 開發後果（Development Consequence）

為什麼這在建立我們自己的應用程式、API、SDK、韌體或硬體支援產品時很重要。

#### 建議實作 / 控制（Recommended Implementation / Control）

要實作什麼，以及哪一層擁有它。

#### 驗證（Validation）

如何證明控制存在或安全地失敗。

#### 限制（Limits）

這不能解決什麼。
```

## 2. 內容分類（保持分離）

| 內容 | 放置位置 |
| --- | --- |
| 跨平台安全原則 | `controls/` |
| 平台、應用類型、嵌入式、韌體或硬體產品實作細節 | `platforms/` |
| 語言/執行時期特定的陷阱 | `languages/` |
| 具體的可建置實作模式 | `implementation/` |
| 可重複的審查流程 | `checklists/` |
| 尚未提升的可重複使用開發原則 | 匹配的 `feedback_history/<category>/`，或跨領域的 `feedback_history/common/` |
| APK 分析方法或 hook 技術 | `analysis/apk/`、`workflow/apk-analysis/` |
| 專案特定的板子接線、校準日誌、韌體 dump、bench 測量、裝置識別碼或目標硬體結論 | 專案儲存庫文件 |
| 產品特定的 API 主機、端點、schema 或令牌細節 | 專案儲存庫文件 |
| 原始供應商文件、帳戶特定條款、憑證、沙箱/正式主機、私人 webhook 負載或真實客戶資料 | 專案儲存庫文件，附清理和存取控制 |
| 生成的客戶端、SDK、fixtures 和提供者/消費者合約檢查 | `implementation/` 和專案儲存庫 |
| Screen Mapping、Consumer Contract、UI Behavior Contract、Screen Contract、Frontend ViewModel Contract、Accessibility Contract、Screen Traceability | [`ui-contracts.md`](ui-contracts.md)、templates 和專案儲存庫的 planning / contract artifacts |
| UI governance evidence：governance domain、render context、collection method、validation mechanism、evidence class、severity、project-local design-system policy、responsive / visual / AI review scope | [`ui-governance.md`](ui-governance.md)、[`templates/ui-governance-evidence-template.md`](templates/ui-governance-evidence-template.md) 和專案 review / validation artifacts |
| Journey Specification / Journey Validation evidence：project-defined critical journey、BDD scenario ref、side-effect chain、expected outcomes、observable evidence、validation result | [`test-strategy.md`](test-strategy.md)、[`validation.md`](validation.md) 和專案 `tests/bdd` / validation artifacts |
| Product Brief 驗證、Impact Map × Customer Journey Map 對齊、文件優先順序、可追溯性和 BDD 閉環流程 | `process/`、templates 和 checklists |
| 重構、遷移、replacement 或新入口替代舊入口的新舊能力 parity inventory | 專案規劃文件、implementation plan 或專屬 parity inventory；若是可重用流程缺口，回饋到 `workflow/software-delivery/` |
| 效能預算、負載/壓力/尖峰/浸泡策略、CI smoke 檢查和發布證據 | `process/`、`CHECKLIST.md`、templates 和專案儲存庫的測試或發布筆記 |
| 共享的清理或回饋規則 | `enforcement/` |
| 僅限本機的暫記筆記、憑證或暫時性流程產出 | 僅限專案儲存庫：**gitignored** 路徑搭配**中性**目錄命名；透過環境變數和可選的未追蹤檔案進行配置；保持追蹤的 README 不含檔案系統導覽和內部代號 |

## 3. 可重複使用指引邊界（Reusable Guidance Boundary）

本節應用 [`enforcement/reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md) 中的全域規則。

技能文件應描述可重複使用的原因、故障模式、決策規則和驗證方法。不要將專案事件透過複製其應用程式名稱、模組名稱、端點路徑、範例負載、類別名稱、本機路徑、主機或即時資料特異性提升到技能中。

如果一個事件教導了有用的課程，將其拆分：

- **技能：** 通用規則，例如「針對即時服務報告的 SDK bug 應透過 SDK 公開表面重現，然後用行為規格和回歸測試鎖定。」
- **專案文件：** 具體的重現目標、受影響功能、樣本 ID、即時環境筆記、BDD 檔案名稱、測試類別名稱和執行結果。

從 `templates/README.md` 開始選擇可複製的文件形狀。從 product brief 開始或驗證時使用 `templates/initial-development-docs.md`，可重複使用指引使用 `templates/hardening-note.md`，快速功能審查使用 `templates/threat-model-lite.md`。

## 4. 必要連結更新聲明（Required Linked Update Statement）

每個影響多個資料夾的可重複使用筆記必須遵循 [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) 並包含簡短的連結更新聲明：

```markdown
#### 必要連結更新（Required Linked Updates）

- `controls/...`：已更新或已檢查，因為 ...
- `implementation/...`：已更新或已檢查，因為 ...
- `checklists/...`：已更新或已檢查，因為 ...
```

如果不需要連結更新，說明原因。這使得相關文件必須保持同步而非可選的後續行動變得明確。

## 5. 良好指引標準（Good Guidance Criteria）

良好的開發指引是：

- 工程師可執行的。
- 明確說明擁有者層。
- 足夠清楚以轉化為程式碼、配置、測試或審查檢查清單項目。
- 可測試的。
- 誠實面對殘留風險。
- 已清理且不含目標特定細節。
- 對 replacement / migration / refactor 類變更，能反查舊入口到新入口的能力覆蓋率、明確 deferred 項目與測試證據。
- 對 consumer / UI 變更，能反查 BDD scenario、screen mapping、consumer needs、screen states、view model derivation、accessibility expectations、API / table ownership 與對應 contract / fixture / test。
- 對 UI compliance claim，能反查 UI governance domain、render context、collection method、validation mechanism、evidence class、severity、project-local design-system policy 與 responsive / visual / AI review scope；不要把 browser review、visual diff 或 AI review 當成 governance domain 本身。
- 對 critical journey claim，能反查 BDD-owned Journey Specification、criticality reason、side-effect chain、expected outcomes、observable evidence 與 validation result；不要把 Journey 當成 `validation_domain` 或把 API success 當成 outcome proof。
- 在需要效能證據時明確，包括指標、預算、環境、執行器和發布關卡。

## 5.1 UI Governance Evidence Shape

當 artifact 聲稱 UI compliance、design-system compliance、accessibility compliance、behavior pattern coverage、visual baseline 或 AI visual review 時，使用 focused UI governance evidence，而不是把 enforcement 細節塞回 UI contract。

最小欄位：

- **Owner layer**：`workflow` evidence；project-specific raw screenshots / scan outputs 留在專案 artifact。
- **Governance domain**：Contract / Design System / Accessibility / Responsive / Behavior / Closure / not_applicable。
- **Render context**：desktop / tablet / mobile / narrow_mobile / not_applicable，或 project-declared equivalent；responsive claim 至少需要 desktop + mobile（或同等寬 / 窄 context）evidence。
- **Collection method**：contract_readback / static_analysis / runtime_trace / browser_review / human_observation / not_applicable。
- **Validation mechanism**：deterministic / screenshot_diff / ai_review / manual_review / not_applicable。
- **Evidence class**：contract / runtime / accessibility_scan / visual_diff / screenshot / responsive_capture / ai_review / human_review / not_applicable。
- **Severity**：block_candidate / warn / research / not_applicable，且 AI visual review 預設不升級為 hard block。
- **Project-local policy**：design token 或 component primitive policy 的專案來源；本 workflow 不定義全域 token scale。
- **Linked artifacts**：UI contract、test strategy、validation result、review decision 或 deferred scope。

## 5.2 Journey Validation Evidence Shape

當 artifact 聲稱 critical user journey、membership / entitlement / identity flow、payment-like path、irreversible action 或其他 project-defined critical journey 已完成時，使用 Journey Validation evidence shape。

最小欄位：

- **Owner layer**：BDD owns Journey Specification；validation owns Journey Execution / Evidence Evaluation。
- **Journey name**：project-defined name，不使用 framework canonical journey list。
- **Criticality**：critical / optional；critical 需列出 revenue、identity、entitlement、security、irreversible_action 或 project-specific equivalent。
- **BDD reference**：`tests/bdd` scenario、feature、manual evidence spec 或 pending runner state。
- **Side-effect chain**：使用者 action 後預期發生的 state transition chain。
- **Expected outcomes**：真實狀態或產品結果，例如 `membership_active`、`playback_allowed`。
- **Observable evidence**：證明 expected outcomes 的 artifact / readback path，例如 DB readback、profile badge、protected resource access、event record。
- **Validation result**：pass / fail / blocked，並記錄缺失 evidence 或 deferred scope。

API 200、adapter success、mock pass、screen-level UI pass 或單一 screenshot 不能單獨完成 Journey Validation。

## 5.3 Governance Invariant Evidence Shapes

當 artifact 聲稱 runtime capability、authority-coupled side effect、configuration applied/readback、或 operational transaction 已完成時，使用 evidence shape，而不是新增 incident-specific checklist。

### Runtime Capability

最小欄位：

- **Capability**：要證明的 runtime capability，例如 browser API、platform permission、filesystem access、container mount、service discovery。
- **Runtime family**：browser / platform / filesystem / container / orchestration。
- **Claim scope**：這個 capability 支撐哪個 feature、workflow 或 release claim。
- **Capability readback**：feature detection、permission check、runtime probe、contract readback 或等效證據。
- **Fallback evidence**：capability absent / denied / unavailable 時的 validated behavior。
- **Result**：pass / fail / blocked。

### Authority-Coupled Side Effect

最小欄位：

- **Business truth**：被宣稱完成的產品或業務狀態。
- **Authority event**：有資格宣告 business truth 成立的事件或狀態。
- **Observable proxy**：click、API 200、adapter success、log、local counter 等 proxy signal。
- **Durable / external evidence**：DB、event record、provider/gateway confirmation、object/read model 或等效 readback。
- **User / business observable readback**：使用者或業務結果可觀察證據。
- **Rejected proxy-only signals**：不可作為 final proof 的訊號。

### Configuration Readback

最小欄位：

- **Desired state**：輸入設定或期望設定值。
- **Applied state**：部署、同步、生成或套用步驟。
- **Readback state**：實際 runtime/deployed state 的 readback path。
- **Validation evidence**：config input、runtime readback、user/API observable state。
- **Mismatch handling**：desired/applied/readback 不一致時的 fail / blocked / rollback decision。

### Operational Transaction Closure

最小欄位：

- **Operation**：deploy、migration、backfill、cache rebuild、data import、batch job 或 project-defined operation。
- **Transaction state**：started、partial、interrupted、resumed、completed、verified。
- **Final state readback**：runtime state、data count、version、health、business effect 或等效證據。
- **Evidence captured**：start record、completion record、final-state verification。
- **Residual unknowns**：若 final state 無法驗證，必須縮小完成宣告或標記 blocked。

## 6. 避免（Avoid）

- 將未發布的工作流程映射到暗示性目錄名稱、列出開發者機器路徑或重複內部調查故事的追蹤 Markdown——這些屬於 gitignore 下的僅本機筆記，而非預設分支敘述。
- 使用「混淆」而不說明它保護什麼和不保護什麼。
- 僅基於功能測試或平均延遲就說「效能沒問題」，而沒有 P95/P99、吞吐量、錯誤率、資源、基準或環境上下文。
- 沒有輪換計劃或威脅模型就「添加固定」。
- 將「檢測 root」作為硬性授權決策。
- 將「在應用程式中隱藏機密」作為持久的安全控制。
- 將第三方 APK 的原始發現複製到可重複使用的文件中。
- 只描述新架構或新命令，卻沒有列出舊能力、新舊對照、副作用與 parity 驗證。
