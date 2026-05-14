# Traceability Gate

從 [`skills/app-development-guidance/process/README.md`](../../skills/app-development-guidance/process/README.md)（已刪除）提取。當分析一個已實作完成的 repository 並需要建立文件與程式碼之間的雙向追溯時，使用本方法。

## 適用時機

- 需要為已實作完成的專案建立文件追溯性。
- 需要確認每個需求都有對應的實作與測試。
- 需要從程式碼反向追溯回原始需求或規則。

## 追溯連結

| 連結 | 目的 |
| --- | --- |
| Product or rule ID → BDD | 顯示哪個 behavior 證明該需求。 |
| BDD → code refs | 顯示 behavior 在哪裡實作。 |
| BDD → test refs | 顯示 behavior 如何被驗證，或存在什麼 gap。 |
| Contract operation / command / diagnostic → fixture | 顯示 provider/consumer 相容性與 edge cases。 |
| Generated client or SDK method → API/OpenAPI/source contract | 防止手抄 endpoints 與 drift。 |

## Stable ID 類型

Stable IDs 可以是：
- Feature IDs
- Rule IDs
- Operation IDs
- Route names
- Command names
- Diagnostic codes
- Event names
- Scenario tags

## 未實作行為的處理

若某個 behavior 被刻意記錄但未實作，標記為以下之一並附上原因與 owner：

| 標記 | 含義 |
| --- | --- |
| `TBD` | 待決定 |
| `noop` | 無操作（intentionally empty） |
| `not enforceable by tool` | 工具無法強制執行 |
| `manual-only` | 僅手動驗證 |
| `out of scope` | 明確排除在範圍外 |

## 與其他層的關係

- `workflow/repo-analysis/` 引用本方法作為分析步驟的具體實作。
- `analysis/repo/documentation-backfill.md` 提供文件恢復的完整流程。
- `skills/app-development-guidance/process/README.md` 是原始來源，已刪除。內容已由本文件承接。
