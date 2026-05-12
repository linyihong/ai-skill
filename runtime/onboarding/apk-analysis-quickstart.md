# APK Analysis Quick Start（APK 分析快速入門）

本文件提取 [`skills/apk-analysis/SKILL.md`](../../skills/apk-analysis/SKILL.md) 中 Quick Start 的操作步驟（非 routing 部分），作為 `runtime/onboarding/` 層的執行指引。

> **相容性規則**：`skills/apk-analysis/SKILL.md` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## 前置步驟（Routing）

開始操作前，先完成以下 routing 判斷：

1. **確認範圍與授權**：APK、版本、裝置/模擬器、架構、允許的分析行為。不收集或發布 credentials、完整 token、私鑰、個資、無關第三方流量。
2. **分離方法與目標事實**：可重用技巧放 skill folder；target-specific endpoints/hosts/schemas/findings 放專案 API/reference docs。
3. **網路路徑分類**：檢查 localhost traffic → whole-device pcap → system proxy/MITM → Java HTTP hooks → Flutter Dart AOT/native paths。
4. **路由到 technique category**：證據指向特定技術才讀對應 techniques 文件。

## 步驟 5：建立 UI 架構地圖

當裝置/app 可操作時，建立 UI 架構地圖：

### 5.1 輕量開始

- 僅拍攝足夠的去敏截圖/UI hierarchy，以理解主要 tab、drawer、routes、key screens。
- 標記哪些 screen 可 scroll、哪些可見元素是 clickable entry point。

### 5.2 記錄操作路徑

- 記錄如何到達每個重要 screen：entry state、tap/swipe steps、expected destination、reusable operation id。
- 當使用者要求分析特定 page/tab/module 時，建立或更新 project-level page map artifact（例如 `docs/UI架構地圖/<page>.md`），而非僅留在 API docs、tool docs 或對話中。

### 5.3 UI Behavior 回填

在回報分析完成前，強制執行 UI Behavior 回填：

- 更新專案的 UI Behavior entry/index 與相關 page-level map。
- 記錄：observed App actions、visible sort labels、tap/swipe steps、data source mapping、evidence、unknowns。
- 若 UI behavior 未被 capture，在專案文件中明確標記 `needs capture` / `Trigger confidence: low`，而非省略該章節。

### 5.4 操作地圖範圍

- 操作地圖僅限於 in-app pages。
- 若步驟開啟其他 app、system screen、browser 或 external intent，記錄該轉換而非視為 app screen。

### 5.5 可選：自動化腳本

- 對關鍵流程，可選擇建立小型可 replay 的 app-operation script，使 API capture 可重複執行且 timing 穩定。

### 5.6 順序調整

- 若截圖或裝置控制使 app 變慢，調整順序：先解決 core API/decode，再將重要 API 綁回 UI actions。
- 僅對需要 API attribution 的流程記錄精確 UI path 與 action window。

## 步驟 7：將動態結果轉換為持久化資產

分析完成後，將動態 hook 結果轉換為以下持久化資產：

### 7.1 必要資產

| 資產類型 | 說明 |
|---------|------|
| UI 架構地圖 + Operation-to-API Matrix | 頁面層級的操作與 API 對應表 |
| 去敏 HTTP/API docs | headers、request fields、response fields、per-field meaning/type notes |
| Domain / Runtime Baseline（實取數據門檻） | 見下方詳細說明 |
| Feature Reconstruction Handoff | capability、behavior scenarios、domain concepts、API/interface contracts、state/error handling、data lifecycle、fixtures、open questions |
| 去敏 request/response samples | 替換敏感資訊後的樣本 |
| Offline decoders or fixtures | 可離線執行的解碼器或測試 fixture |
| API/schema docs | 完整的 API 列表與 schema 文件 |
| Contract tests | 當專案有 SDK 或 client implementation 時 |

### 7.2 Domain / Runtime Baseline 詳細要求

API 條目與 schema **不足以**讓 SDK、client 或回放工具自動連線。必須在**專案**內另有一份 **domain／執行環境基線**，包含：

- host/path family 與環境維度（去敏占位）
- TLS／代理／內嵌通道假設
- **session／登入／裝置身分**對列表資料的依存
- 必填 **opaque query**（例如 session 標量）從哪些操作或 response 衍生
- 簽章或 gateway 前置條件（不留 secret）
- 分頁**地面真相**（明確旗標 vs 長度推斷）
- 限流／錯誤恢復
- `service`/`serviceHash` 與公開端點的對應表（值可 placeholder）

以表格**連回** API Catalog 條目與 UI operation id。可新建 `docs/domain-baseline.md`、`api/domain-environment.md` 或附在現有 inventory。**不放進**可重用 skill 正文。

### 7.3 Development Readiness Gate

若下一步是 SDK/client/app/tool development 且必須與真實服務通訊、replay feature 或執行 live integration，**不可**僅從 API shape 開始實作。先驗證專案文件包含最小可跑因素：

- endpoint/path family、route/service mapping 或 placeholder strategy
- session/bootstrap dependency
- opaque parameter source and lifetime
- Authorized identity material self-generation audit
- signing/gateway prerequisites
- response decrypt/unwrap boundary
- pagination truth
- error/session recovery
- replay checklist

缺失因素必須寫為 blocker 或明確 scope out。僅 offline parser、fixture、mock 或 docs-only 工作可從 skeleton baseline 開始。

### 7.4 Authorized Identity Material Self-Generation Audit

當 device、install、account、session seed、vendor attestation 或 server-issued session material 是 live access 的必要條件時，分析必須回答兩個問題：

1. **SDK/tool 能否在沒有 target app runtime 的情況下產生或初始化該 material？**
2. **如果可以，去敏的產生 recipe、lifecycle、reset/cooldown、validation matrix 是什麼？**

若答案為 no 或 unknown，命名 provider boundary（`caller-provided`、`server-issued`、`trusted-bridge`、`private-adapter-required`、`unknown`）並視為 live-development blocker，除非明確 scope 到 private adapter。

### 7.5 API Catalog / API List Finish Gate

當任務需要 API reference、SDK/client input、mock API、contract test 或 rebuildable feature 時，建立或更新 project-level API Catalog：

- total API entry
- grouped indexes
- per-API detail files
- coverage/gap status
- UI/API mapping
- SDK/client field usage
- evidence
- validation
- open questions

不可將 confirmed APIs 僅留在 schema catalogs、correlation tables、hook logs 或對話中。

### 7.6 Feature Handoff Finish Gate

當 named feature/module 已足夠理解到可標記 core flows 為 `Confirmed` 或回答 implementation-style 問題時，在同一 session 建立或更新 project-level feature handoff / architecture document。不可將理解僅留在 API tables、hook logs、對話或 page maps 中。

## 與其他層的關係

- `workflow/apk-analysis/execution-flow.md` 提供分析執行流程，本文件提供快速入門的操作步驟。
- `workflow/apk-analysis/artifact-gates.md` 提供產出規範與品質門檻。
- `intelligence/engineering/apk-analysis/` 提供決策啟發式（如 hook 選擇）。
- `skills/apk-analysis/SKILL.md` 是原始來源，仍為 active entrypoint。
