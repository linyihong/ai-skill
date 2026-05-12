# Contract Governance Gate

從 [`skills/app-development-guidance/process/README.md`](../../skills/app-development-guidance/process/README.md) 提取。當分析一個 repository 並需要定義文件間的優先順序與衝突處理規則時，使用本方法。

## 適用時機

- Repository 有多份開發文件且需要定義哪份文件優先。
- 文件之間存在衝突，需要系統化的解決方案。
- 需要建立文件變更的治理規則。

## 預設文件優先順序

除非專案有更強的 local rule，否則使用以下優先順序：

| 優先級 | 文件類型 | 說明 |
| --- | --- | --- |
| 1 | Governance / framework contract | Repository-wide invariants、required update rules、dependency direction、naming、build/run constraints |
| 2 | Product plan / accepted brief | Product intent、scope、non-goals、canceled requirements、business language |
| 3 | BDD behavior | Observable user/system behavior、acceptance criteria |
| 4 | Domain, architecture, API/interface, error handling, hardware, or command contracts | 各層契約 |
| 5 | Implementation and generated clients | 實作與產生的客戶端 |
| 6 | Tests, fixtures, and examples | 測試與範例 |

## 衝突處理

若較低層發現較高層有誤，不要默默在程式碼中「修正」。分類衝突類型並執行對應行動：

| 衝突類型 | 必要行動 |
| --- | --- |
| **Product intent changed** | 更新 product brief 或 plan，然後更新 BDD/contracts/tests |
| **BDD missing or stale** | 從 evidence 回填或修訂 BDD，連結受影響的 tests |
| **Contract stale** | 在同一次變更中更新 contract 與所有 consumers、mocks、generated clients、fixtures、tests |
| **Implementation bug** | 保持 docs 穩定，新增或更新 regression tests，然後修正程式碼 |
| **Test or fixture stale** | 更新 tests/fixtures 到當前 contract，引用來源 |

## 取消與延後行為的記錄

明確記錄以下項目，不要讓它們成為未來 agent 可能重新引入的隱形空缺：

- Canceled（已取消）
- Deferred（已延後）
- Out-of-scope（排除在範圍外）
- Not tool-enforceable（工具無法強制執行）

## 與其他層的關係

- `workflow/repo-analysis/` 引用本方法作為分析步驟的具體實作。
- `analysis/repo/documentation-backfill.md` 提供文件恢復的完整流程。
- `analysis/repo/traceability-gate.md` 提供追溯連結的建立方法。
- `skills/app-development-guidance/process/README.md` 是原始來源，仍為 active entrypoint。
