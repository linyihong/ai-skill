# Repository Analysis

`analysis/repo/` 負責「Repository 分析與理解方法」。本目錄保存如何觀察、拆解與理解一個陌生 repository 的分析框架，讓 agent 能快速建立 repository 的心智模型，而不需要完整讀完所有原始碼。

## 核心責任

- Repository 結構觀察（目錄佈局、語言分佈、建置系統）。
- Dependency graph 分析（套件依賴、內部模組依賴）。
- Codebase 心智模型建立（entrypoint、核心抽象、資料流）。
- 技術債與品質評估（測試覆蓋、lint 狀態、dead code）。
- 遷移或重構前的影響範圍分析。
- **已實作專案的文件恢復**：從現有程式碼反向恢復缺失的開發文件。
- **文件追溯性建立**：建立需求、實作、測試之間的雙向連結。
- **契約治理**：定義文件優先順序與衝突處理規則。

## 分析方法

### 1. 靜態結構分析

```
1. 讀取 repository 根目錄結構。
2. 識別建置系統（Makefile、Cargo.toml、package.json、build.gradle 等）。
3. 識別主要語言與框架。
4. 識別測試目錄與測試框架。
5. 識別 CI/CD 配置（.github/、.gitlab-ci.yml、Jenkinsfile 等）。
6. 識別文件目錄（docs/、wiki/、README.md 等）。
```

### 2. 依賴分析

```
1. 讀取套件管理檔案（Cargo.toml、package.json、Gemfile、requirements.txt 等）。
2. 分類依賴類型（runtime / dev / build / optional）。
3. 識別內部模組之間的依賴關係。
4. 標記 circular dependency 風險。
5. 評估外部依賴的維護狀態（版本新鮮度、替代方案）。
```

### 3. Entrypoint 與核心流程

```
1. 找到應用程式 entrypoint（main.rs、main.go、index.js、Application.java 等）。
2. 追蹤請求生命週期（request → middleware → handler → response）。
3. 識別核心抽象（trait、interface、abstract class）。
4. 識別資料流（input → processing → storage → output）。
5. 識別錯誤處理策略（panic、Result、exception、error code）。
```

### 4. 技術債評估

```
1. 檢查測試覆蓋率（有無測試、測試類型、測試品質）。
2. 檢查 lint / format 配置與執行結果。
3. 檢查 dead code 可能性（未使用的 import、export、function）。
4. 檢查安全性問題（硬編碼 secret、SQL injection 風險、授權缺失）。
5. 檢查建置與部署流程（Dockerfile、deployment script、migration）。
```

### 5. 已實作專案文件恢復

參見 [`documentation-backfill.md`](documentation-backfill.md)。從現有程式碼、tests、schemas、API specs、fixtures 系統化恢復缺失的開發文件。

### 6. 文件追溯性建立

參見 [`traceability-gate.md`](traceability-gate.md)。建立 Product/rule ID → BDD → code → test 的雙向追溯連結。

### 7. 契約治理

參見 [`contract-governance.md`](contract-governance.md)。定義文件優先順序與衝突處理規則。

## 已提取內容

| 來源 | 目標 | 內容 |
| --- | --- | --- |
| `skills/app-development-guidance/process/README.md` §Existing Project Documentation Backfill | [`documentation-backfill.md`](documentation-backfill.md) | 8 種文件恢復規則、6 種 pipeline artifact 恢復方法、7 步恢復順序 |
| `skills/app-development-guidance/process/README.md` §Traceability Gate | [`traceability-gate.md`](traceability-gate.md) | 5 種追溯連結、stable ID 類型、未實作行為標記 |
| `skills/app-development-guidance/process/README.md` §Contract Governance Gate | [`contract-governance.md`](contract-governance.md) | 6 級文件優先順序、5 種衝突處理規則、取消/延後記錄方法 |

## 與其他層的關係

- `workflow/repo-analysis/` 可引用本層的分析步驟，但不複製分析方法細節。
- `intelligence/engineering/architecture/` 承接從 repo 分析中萃取的架構判斷。
- `intelligence/engineering/domain/` 承接從 repo 分析中萃取的領域模型理解。
- `skills/app-development-guidance/process/README.md` 是原始來源，已不再作為 active entrypoint。

## 產出格式

每次 repo 分析應產出：

- **Repository 心智模型摘要**（≤300 tokens）：語言、框架、核心抽象、資料流。
- **依賴圖重點**（≤200 tokens）：關鍵外部依賴、內部模組依賴、circular dependency 風險。
- **技術債清單**（≤200 tokens）：需要關注的品質問題與優先順序。
- **影響範圍分析**（≤200 tokens）：如果修改 X，哪些檔案會受影響。
- **文件恢復狀態**（≤200 tokens）：已恢復的文件與仍缺失的文件。
- **追溯性矩陣**（≤200 tokens）：需求 → BDD → code → test 的對應狀態。
