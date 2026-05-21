> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-21 — Vendor Integration Architecture（多廠商整合架構選型）

Status: candidate

#### One-line Summary

整合超過 3 個外部廠商時，需在 Adapter / Compile-time submodule / Plugin SPI / Out-of-process service / Hybrid 之間選型；當廠商數量 N ≥ 10 仍維持「每個廠商一個編譯期模組」，編譯時間、IDE 載入、升級成本會在 N 達兩位數時崩壞，三位數時失控。

#### Human Explanation

支付聚合、社群登入、IM、博弈聚合、廣告聯播、地圖等場景常需整合多家外部廠商。最直覺的做法是「每個廠商一個 Maven/Gradle 子模組」，但這條路在 N 達到兩位數後快速劣化：CI build 時間線性甚至超線性增長、IDE 對 100+ 模組的 indexing 顯著退化、升級共用框架需要動每個 vendor 模組、團隊風格不一致造成同一 repo 內有多種整合風格。Plugin / SPI 的設計需要前期投入，但長期成本是線性可控的。

#### Trigger

- Vendor 相關 Maven/Gradle module 數量 ≥ 10
- 升級共用框架（Spring Boot、序列化器）必須觸及大量 vendor 模組
- IDE 載入時間 ≥ 30 秒或卡頓明顯
- 廠商需要灰度 / A/B / 停用 / 降級切換而不能重啟服務
- 廠商 SDK 之間存在相依衝突（同一函式庫不同版本）

#### Evidence

- Tool: 多廠商整合的後端專案 review（目錄結構觀察）
- Sanitized excerpt: 觀察到一個後端 repo 對某一家廠商展開了三位數量級的「每個 SKU/遊戲一個 Maven module」結構，且同 repo 內不同廠商整合風格高度不一致（有 4 層拆分、有平鋪、有 test-only 殼）
- Evidence path: 證據留在 `<PROJECT_ROOT>/` 原始 review note，不複製到本庫

#### Generalized Lesson

五種整合策略與適用條件：

| 策略 | 適用條件 |
|------|------|
| A. Adapter / Strategy | N ≤ 5，廠商行為高度相似 |
| B. Compile-time submodule | 廠商 SDK 相依嚴重衝突，需 classpath 隔離，N < 10 |
| C. Plugin / SPI | N 中大（10~100），廠商頻繁增減，需動態啟用 |
| D. Out-of-process | 廠商 SDK 不穩定、需獨立 scale 或合規隔離 |
| E. Hybrid | N 三位數以上、廠商重要性差異大 |

關鍵：**SKU/遊戲是資料不是程式碼分群**，不應拆成 Maven module；應以 config / catalog 形式存在。

#### Agent Action

評估後端 repo 的多廠商整合架構時：

1. 數出 vendor 相關 module 數量
2. 比對策略軸向（廠商數量、新增頻率、SDK 衝突、熱換需求、行為相似度）
3. 若 N ≥ 10 且仍是 compile-time submodule per vendor，明確指出架構味道並建議遷移到 SPI 或 Hybrid
4. 區分「廠商」與「SKU/遊戲/商品」：前者是程式碼分群，後者是資料

#### Goal / Action / Validation

- Goal: 廠商整合成本維持線性，避免 N 達兩位數後架構崩壞
- Action: 依軸向選擇策略；既有 compile-time module per vendor 在 N ≥ 10 時規劃遷移到 SPI
- Validation or reference source: 廠商相關 module 數 ≤ 廠商數 × 4；全量 build < 5 分鐘；新增廠商 < 3 工作日；升級共用相依不需動所有 vendor 模組

#### Applies When

- 後端需要整合多家外部廠商（payment、login、IM、gaming、ads、map、cloud SDK 等）
- 廠商列表會持續增長

#### Does Not Apply When

- 廠商數量穩定且小（N ≤ 3）
- 整合內容非常薄（只是呼叫 REST endpoint，無 SDK）

#### Validation

- 觀察 N 與架構策略對應關係是否符合軸向
- 重構案例顯示 compile-time → SPI 遷移後 build 時間與新增廠商成本下降

#### Promotion Target

- ✅ `intelligence/engineering/architecture/vendor-integration-architecture.md`（已於 `d5ec684` 寫入）
- ✅ `knowledge/summaries/vendor-integration-architecture.md`（已於 `d5ec684` 寫入）

#### Required Linked Updates

- ✅ `intelligence/engineering/architecture/README.md` 索引（已於 `d5ec684` 更新）
- ✅ `knowledge/summaries/README.md` 索引（已於 `d5ec684` 更新）
- Step 6（Intelligence Extraction）不適用：intelligence atom 已直接寫入
- Step 7（Failure Learning）不適用於本 lesson 內容；本輪流程失誤已由既有 `enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md` 覆蓋
