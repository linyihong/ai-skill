# Existing Project Documentation Backfill

從 ``skills/app-development-guidance/process/README.md``（已刪除）提取。當分析一個已實作完成的 repository 時，使用本方法系統化地恢復缺失的開發文件。

## 適用時機

- 第一次接觸一個已實作完成的 repository。
- 需要為 legacy codebase 補上開發文件。
- 需要從現有程式碼反向恢復 domain intent 與架構決策。

## 文件恢復規則

| 缺失文件 | 恢復方法 |
| --- | --- |
| **Product Brief** | 只恢復證據支持的內容：可見的目標、使用者/角色、範圍、非目標、假設、限制。若原始意圖不可得，標記為 `unknown` / `open question`；不要發明商業理由。 |
| **Bounded Context Map** | 從程式碼所有權、runtime boundaries、資料庫表格、API groups、UI areas、queues、SDK/public APIs、deployment units 推斷模組邊界。 |
| **BDD Behavior** | **必須完成。** 從已實作產品的 UI、API behavior、tests、logs 恢復 critical happy paths、failure paths、permissions、empty states、edge cases、cross-context flows。 |
| **Domain Model Contract** | 從 code、schemas、storage、UI states、tests 推斷 entities、value objects、commands、events、invariants、state transitions；不確定的詞彙標記為 candidate。 |
| **Architecture Contract** | 記錄實際的 dependency direction、data ownership、side-effect boundaries、integrations、runtime/deployment shape、已知違規。 |
| **API / Interface Contract** | 提取實際的 request/response schemas、public methods、events、commands、auth/session behavior、versioning、compatibility、fixtures、consumers。 |
| **Error Handling Contract** | 恢復觀察到的 error taxonomy、retry rules、user messages、logging/redaction behavior、security-sensitive failures、gaps。 |
| **Test Plan** | 將現有 tests 對應到 behavior/contracts，列出未覆蓋的 BDD scenarios、invariants、contracts、integration paths 所需的 tests。 |

## Pipeline Artifact 恢復

對已實作完成的專案，同時恢復交付管線：

| Pipeline artifact | 恢復方法 |
| --- | --- |
| **Plan index / product radar** | 將 source product docs、PDFs、tickets、screenshots、legacy notes 對應到 modules、controllers、screens、commands、packages。標記已取消或已取代的需求。 |
| **Contract taxonomy** | 列出哪些 documents 管轄 build/run、HTTP/API shape、auth/tenant/session、persistence、domain layering、frontend/backend integration、third-party integration、testing、documentation sync。 |
| **Minimum doc sync matrix** | 對每種 change type，列出最少需更新的 docs/tests：API、permission、database、UI flow、generated client、vendor integration、CLI command、diagnostic rule、release setting。 |
| **OpenAPI / schema / generated client** | 驗證 generated consumer code 來自 source contract，不是手抄 endpoints 或 DTOs。 |
| **Vendor / third-party integration** | 分離 raw vendor docs 與 sanitized integration excerpts、request/response contracts、fixture examples、live-test gates、secret handling。 |
| **Tooling / extension rule catalog** | 對應 catalog order、rule IDs、diagnostics/commands、fixtures、tests；明確標記 process-only 或 non-enforceable rules。 |

## 恢復順序

1. **Inventory** 現有 docs、source folders、tests、schemas、API specs、fixtures、release notes、observed behavior。
2. 建立 **documentation gap table**，狀態：`exists`、`partial`、`missing`、`unknown`。
3. **先恢復 BDD Behavior**（當 product brief 缺失時），因為已實作的行為是最強的可用真相來源。
4. 從完成的 behavior 與 implementation evidence 恢復 **Domain Model、Architecture、API/Interface、Error Handling Contracts**。
5. 將 **unknown product intent** 與 **observed behavior** 分開標記。Unknown intent 不阻礙 BDD 完成。
6. 若 BDD 無法從可用證據完成，**停止並要求** missing behavior、screen/API examples、logs、test cases 或 user decisions。
7. 對任何缺乏 coverage 的 critical BDD scenario，**新增 tests 或 test TODOs**。

## 與其他層的關係

- `workflow/repo-analysis/` 引用本方法作為分析步驟的具體實作。
- `intelligence/engineering/architecture/` 承接從 repo 分析中萃取的架構判斷。
- `intelligence/engineering/domain/` 承接從 repo 分析中萃取的領域模型理解。
- `skills/app-development-guidance/process/README.md` 是原始來源，已刪除。內容已由本文件承接。

## 遷移狀態

- 本文件是 `skills/app-development-guidance/process/README.md`（已刪除）的 reference target，已取代舊入口。
- 新內容請直接寫入此文件。
- 恢復的 Product Brief 欄位若為 `unknown`，必須明確標記，不得發明商業理由。
