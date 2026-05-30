# Validation Slice（Performance Test Gate / Validate）

> **Cognitive Slice**：`sd-validation`（從 [`execution-flow.md`](execution-flow.md) §5 + §7 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-validation` |
| `purpose` | 變更出貨前的驗證關卡：效能測試的觸發與執行（performance test gate）、以及最低驗證方法清單（validate）；確認「舊行為仍受保護 / 新程式碼已證明」二者皆有證據 |
| `type` | `execution` |
| `tags` | artifact-gate, validation, performance |
| `load_when` | 驗證變更 / 效能關卡 |
| `do_not_load_when` | 尚未實作完成前、純 intake / contract 規劃任務 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「出貨前要過哪些 gate、用哪些驗證方法」的 ordering / gate；通過 workflow membership test，不承載 evidence 取得方法（非 analysis），不論證長期模式（非 intelligence） |
| `canonical_source` | 本檔（原 `execution-flow.md` §5 效能測試關卡 + §7 驗證） |
| `dependencies` | `sd-implementation`（實作完成才能驗證）、`sd-test-strategy`（perf 測試類型選型在 test-strategy slice，本 slice 引用） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Phase 4 Scenario A（execution-only：完成宣告前的最低驗證）、Scenario C（mixed：debug 失敗 deployment pipeline 引用本 slice） |

> **Perf 內容邊界（與 sd-test-strategy 的分工）**：本 slice 擁有 perf **執行關卡 / gate 觸發條件 / 最低指標**（即「何時必須有 perf 證據、要追哪些 metric」）。Perf **測試類型選型表**（load / stress / spike / soak 何時用哪一種）的 canonical 在 `sd-test-strategy`（development-process.md §Test Strategy Gate 內），本 slice 引用而不複製，避免 dual source-of-truth。

## 1. 效能測試關卡（Performance Test Gate）

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

## 2. 驗證（Validate）

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
