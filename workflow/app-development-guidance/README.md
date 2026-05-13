# App Development Guidance Workflow

`workflow/app-development-guidance/` 負責「App 開發審查與指引的執行流程」。本目錄保存 agent 在進行 app 開發審查時可照著執行的 planning flow、review flow、handoff flow、review checklists 與 contract-first 開發流程，讓開發與審查過程可重複、可驗證。

## Scope

本 workflow 涵蓋以下流程與審查類型：

### 開發流程

- **Contract-First Development Process**：從企劃書到實作的完整開發流程，包含 Default Flow、Required Contracts、Product Brief Validation Gate、Change Intake Gate、Contract Governance Gate、Traceability Gate、BDD Execution Closure、Test Strategy Gate、Embedded/Hardware Flow、Missing Information Gate、Existing Project Documentation Backfill 等。

### 審查類型

- **Design Review**：在實作前審查設計文件、API contract 與架構決策。
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

## 與既有層的關係

- `skills/app-development-guidance/` 目前仍是 active skill entrypoint；本層只承接逐步抽出的通用執行流程。
- `skills/app-development-guidance/WORKFLOW.md` 是目前的 workflow source of truth。
- `skills/app-development-guidance/process/README.md` 提供 contract-first 開發流程（已提取至本層）。
- `skills/app-development-guidance/checklists/` 提供各審查類型的 checklist 正文（已提取 catalog 至本層）。
- `skills/app-development-guidance/controls/` 提供各安全控制的評估標準（已提取 catalog 至 `analysis/app-development-guidance/`）。
- `analysis/repo/` 可被本 workflow 引用來分析 repository 結構。
- `analysis/app-development-guidance/` 提供安全控制、實作模式、平台指引、語言陷阱的 catalog 參考。
- `intelligence/` 可被本 workflow 引用來輔助工程判斷。

## 第一批候選遷移來源

- `skills/app-development-guidance/WORKFLOW.md` — ✅ 已提取（execution-flow.md, artifact-gates.md, analysis/app-development-guidance/risk-translation.md）
- `skills/app-development-guidance/process/` — ✅ 已提取（development-process.md）
- `skills/app-development-guidance/checklists/` — ✅ 已提取（[`review-checklist.md`](review-checklist.md)）

## 已提取內容

| 檔案 | 來源 | 說明 |
|------|------|------|
| [`execution-flow.md`](execution-flow.md) | `WORKFLOW.md` §1, §5-8 | Start From Evidence、Change Intake、BDD Closure Loop、SDK Defect Closure、Same-Session Closure、Performance Gate、Backfill Rules、Validate |
| [`artifact-gates.md`](artifact-gates.md) | `DOCUMENTATION.md` | Reusable Note Structure、Content Classification、Guidance Boundary、Linked Update Statement、Good Guidance Criteria |
| [`analysis/app-development-guidance/risk-translation.md`](../analysis/app-development-guidance/risk-translation.md) | `WORKFLOW.md` §2-5 | Risk Translation Table、Owner Layer Selection、Control Definition、Guidance Classification、Linked Updates |
| [`review-checklist.md`](review-checklist.md) | `skills/app-development-guidance/checklists/` | 6 種審查 checklist 的 catalog（Mobile Design Review、Mobile PR Review、Mobile Release Review、API Security Review、Contract Governance Review、Embedded Firmware Review） |
| [`development-process.md`](development-process.md) | `skills/app-development-guidance/process/README.md` | Contract-first 開發流程：Default Flow、Required Contracts、Product Brief Validation Gate、Change Intake Gate、Contract Governance Gate、Traceability Gate、BDD Execution Closure、Test Strategy Gate、Embedded/Hardware Flow、Missing Information Gate、Existing Project Documentation Backfill、Contract-First Rules、Definition of Ready/Done |

## 建議 Workflow 流程

### Design Review Flow

```
1. 確認審查範圍（新功能 / 架構變更 / API 變更）。
2. 讀取設計文件或 RFC。
3. 檢查設計是否涵蓋：
   ├─ 功能需求與非功能需求。
   ├─ API contract（request/response 格式、錯誤處理）。
   ├─ 資料模型與儲存方案。
   ├─ 安全考量（授權、認證、資料保護）。
   └─ 測試策略（單元測試、整合測試、E2E 測試）。
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

每次 review 應產出：

- **Review 摘要**（≤200 tokens）：審查類型、範圍、verdict。
- **Finding 清單**（每個 finding ≤100 tokens）：問題描述、severity、建議修復方式。
- **Decision 記錄**（≤100 tokens）：最終決定、決定理據、相關連結。
