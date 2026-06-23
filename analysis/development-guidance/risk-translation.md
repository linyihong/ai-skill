# 風險翻譯與控制選擇

本文件定義如何將分析觀察轉換為開發者視角的風險陳述，並選擇最適合的控制層。原承接 ``skills/app-development-guidance/WORKFLOW.md`` §2-4 的內容（已刪除）。

## 1. 翻譯為風險

以開發者用語撰寫風險：

| 觀察 | 開發者結論 |
| --- | --- |
| Request 可以用相同 body 和 token 重放 | 後端合約可能缺少 nonce、timestamp、idempotency 或風險評分 |
| Token 有效期長且範圍廣泛 | Token 洩漏時帳號接管爆炸半徑過大 |
| 敏感值出現在日誌中 | Debug 或遙測管線可能洩漏機密 |
| App 依賴隱藏的客戶端邏輯進行授權 | 伺服器可能信任攻擊者可修改的客戶端狀態 |
| 本地儲存包含可復原的機密 | 裝置入侵或備份可暴露憑證 |
| Release build 暴露符號、debug flag 或測試端點 | 逆向工程成本不必要地低 |
| 生成的客戶端與 OpenAPI/schema 不一致 | 前端、SDK 或工具消費者可能呼叫過時路由或反序列化錯誤形狀 |
| UI 直接從 API 欄位推導顯示狀態且無 ViewModel Contract | 畫面可能能運作，但 display semantics、permission、error mapping 或 derived fields 會和需求語意漂移 |
| Gherkin 存在但無可執行連結 | 團隊可能認為行為已測試，但實際上只有文件 |
| Impact Map 的目標族群與 Customer Journey 的痛點不一致 | 團隊可能快速實作功能，但沒有改善真正影響目標的使用者旅程 |
| 供應商文件直接複製到產品流程中 | 機密、無關欄位或不穩定的第三方假設可能洩漏到實作和可重用文件中 |

## 2. 選擇擁有者層

優先選擇由最強層擁有的控制：

| 層 | 適用場景 |
| --- | --- |
| 後端/API | 授權、重放防禦、速率限制、詐欺訊號、Token 輪換、審計日誌 |
| 客戶端 App | 安全儲存、安全預設、必要時的憑證固定、高風險流程的 UX 摩擦、遙測衛生 |
| 全端合約 | OpenAPI/schema 生成、型別客戶端、Screen Mapping、Consumer Contract、UI Behavior / Screen / ViewModel Contract、provider/consumer fixture、合約測試、相容性規則 |
| 工具/擴充 | 純規則引擎或命令核心、adapter 邊界、診斷/命令、fixture 支援的規則、編輯器/CLI 整合測試 |
| 第三方整合 | 清理後的供應商摘錄、憑證邊界、即時測試關卡、重放/冪等/webhook/結算行為、審計日誌 |
| 嵌入式/韌體 | 感測器/協定解析、硬體上下文注入、驅動/服務/應用程式邊界、RTOS/任務擁有權、主機 fixture、硬體迴圈檢查 |
| 建置/發布 | 混淆、符號剝離、debug flag 強制執行、相依性審查、機密掃描 |
| 監控 | 異常偵測、裝置風險訊號、濫用模式警示 |

客戶端強化可以提高成本並改善訊號品質，但不能作為授權或財務/業務完整性的唯一控制。

## 3. 定義控制

針對每個風險或有用的實作教訓，定義：

- 必要控制
- 擁有者
- 實作說明
- 驗證方法
- 殘餘風險

範例：

```text
風險：API request 可被重放
控制：伺服器驗證 timestamp、nonce、帳號/session 綁定與冪等金鑰
擁有者：後端/API
驗證：整合測試重放相同簽章 request，預期拒絕或冪等處理
殘餘風險：裝置入侵仍可竊取有效 session；需保留監控
```

## 4. 歸檔指引

在撰寫文件前先分類結果：

| 如果教訓是關於... | 放在 |
| --- | --- |
| 跨堆疊適用的安全屬性 | `controls/` |
| 行動端、Web、後端、嵌入式、韌體、硬體或 OS 特定實作 | `platforms/` |
| Dart、Kotlin/Java、Swift、TypeScript 或執行時期特定陷阱 | `languages/` |
| 具體可建置的模式或操作指南 | `implementation/` |
| 重複出現的設計、PR、發布或 API 審查步驟 | `checklists/` |
| 可重用但仍在浮現中的教訓 | `feedback/history/development-guidance/<category>/` 或 `feedback/history/development-guidance/common/` |
| 可複製的文件形狀 | `templates/` |
| 已實作專案中遺失的開發文件 | `process/` 和 `templates/initial-development-docs.md` |
| 缺少需求或合約的阻礙問題 | `process/` 和當前規劃文件 |
| Product goal、target actor、journey pain 或 feature investment 的不一致 | `process/`、`templates/product-impact-alignment-template.md` 和當前規劃文件 |
| 程式碼前的變更接收 | `process/` 和當前規劃文件 |
| 新程式碼或 AI 生成程式碼的測試策略 | `process/`、`CHECKLIST.md` 和當前規劃文件 |
| 嵌入式/硬體產品流程 | `platforms/embedded/`、`implementation/embedded/`、`process/` 和硬體感知檢查清單 |
| OpenAPI/schema/codegen 或全端 provider/consumer 合約 | `implementation/backend/`、`process/`、`CHECKLIST.md` 和相關 API 檢查清單 |
| Screen mapping、consumer needs、screen states、UI behavior、view model derivation、screen traceability 或 accessibility expectations | `workflow/software-delivery/ui-contracts.md`、`implementation/`、templates 和專案規劃文件 |
| 工具、CLI、IDE 擴充、linter 或靜態分析架構 | `implementation/tooling/`、`process/` 和相關審查檢查清單 |
| 供應商或第三方 API 整合 | `implementation/backend/`、`controls/`、`checklists/` 和專案特定的清理後文件 |

優先使用資料夾間的連結，而非複製相同指引。

## 5. 執行必要連結更新

在完成變更前，遵循倉儲級規則 [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md)。如果更新影響相關資料夾，這些連結更新在**同一變更中為必要**：

| 變更區域 | 必須更新或驗證 |
| --- | --- |
| `controls/` | 相關的 `implementation/`、`platforms/`、`languages/` 和 `checklists/` 文件 |
| `implementation/` | 相關的 `controls/`、`platforms/`、`languages/` 和 `checklists/` 文件 |
| `platforms/` | 相關的 `controls/`、`implementation/`、語言說明和檢查清單 |
| `languages/` | 相關的平台和實作文件 |
| `checklists/` | 相關的控制和實作文件 |
| `templates/` | `templates/README.md`、`DOCUMENTATION.md` 以及指示使用者複製範本的文件 |
| `process/` 治理或回填規則 | `templates/`、`CHECKLIST.md`、`WORKFLOW.md` 和相關的 `implementation/` 文件 |

如果不需要連結更新，說明原因。不要將必要的連結更新稱為可選。

## 6. 對 Stakeholder 翻譯（業務語的風險陳述）

§1 把觀察翻成「開發者風險陳述」；本節再翻一層 — 把開發者風險翻成 stakeholder（PM / 業務 / 客戶）聽得懂的具象話術。工程上的論證若停在開發者語彙，往往無法說服 stakeholder 接受 scope 限制或階段延後。

### 翻譯原則

| 從 | 到 |
| --- | --- |
| 抽象工程概念（ground truth、validation、coverage） | 具象後果（看不到、不確定、爆炸、回不來） |
| 量化效益（驗證成本指數成長） | 情緒對比（短期 vs 長期、可控 vs 失控） |
| 教科書詞彙（idempotency、parity） | 日常比喻（重複按按鈕、新版跟舊版做一樣的事） |

### 案例：Migration 不要綁新功能

**情境**：大型系統改版的 kickoff，客戶要求「既然要重寫，順便加 OOO 吧」。

**開發者風險陳述**（§1 風格）：

> 同時加新功能會讓驗證失去 ground truth。新版出現的行為差異無法判斷來源是搬遷錯、新功能 bug 還是兩者互動，回歸成本指數增長。

**Stakeholder 反駁**：「不加新功能的話，使用者會很失望耶。」

**Stakeholder 翻譯**：

> **失望總比絕望好**

為什麼這話術有效：

1. **抽象 → 具象**：把「驗證失去 ground truth」翻成「失望 vs 絕望」的情緒對比
2. **不需要技術背景就能理解**
3. **錨點留在「絕望」**：對方會記住後果嚴重性（系統崩塌、無法回滾、時程黑洞）
4. **時間框架對比**：短期失望（沒新功能）vs 長期絕望（系統壞了無法定位）

完整論證請見 [`intelligence/engineering/anti-patterns/migration-feature-bundling.md`](../../intelligence/engineering/anti-patterns/migration-feature-bundling.md)。

### 寫話術時的 checklist

| 檢查 | 通過條件 |
|------|------|
| 是否具象 | 對方腦中能形成畫面或感覺，不是抽象名詞 |
| 是否避免術語 | 不含 idempotent、parity、coverage、ground truth 等開發者詞彙 |
| 是否有時間框架 | 短期 vs 長期、現在 vs 之後 |
| 是否留錨點 | 對方記住一個關鍵詞（如「絕望」「爆炸」「黑洞」） |
| 是否不誇大 | 翻譯不能把可控風險講成不可控（中性化不能改變結論） |
