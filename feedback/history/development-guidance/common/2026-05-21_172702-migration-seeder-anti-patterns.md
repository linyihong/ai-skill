> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-21 — Migration Seeder Anti-Patterns

Status: candidate

#### One-line Summary

把大量業務資料以巨型 `INSERT` 包進 schema migration，會強制資料 lifecycle 與 schema lifecycle 綁定，造成部署、review、業務速度與環境漂移等多重後果。

#### Human Explanation

Schema migration 工具（Flyway / Liquibase / Alembic）的設計目標是 schema 演進，但容易被誤用為通用資料載入器。當廠商目錄、商品 catalog、SKU 清單、權限矩陣等業務資料被塞進 migration 時，每一筆業務變動都需要工程師寫新 migration、PR、deploy；migration 檔案膨脹到上百 KB 無人能 review；rollback schema 會誤刪資料；不同環境必須用 `IF NOT EXISTS` / `ON CONFLICT` 補丁掩蓋差異。表面上「把資料放 migration」看起來方便，實際上是把短期方便換成長期 lifecycle 耦合。

#### Trigger

- 觀察到單一 migration 檔超過 50KB 且大多是 `INSERT`
- 檔名含 `dataSeeder` / `seed` / `fixture` / `initial_data`
- 同張表反覆出現補資料 migration
- 業務人員想新增一筆資料需要工程師寫 migration
- 不同環境需要不同資料卻共用同一 migration

#### Evidence

- Tool: 多模組後端專案 review（目錄結構觀察、`find` 檔案大小）
- Sanitized excerpt: 觀察到一個多廠商整合的後端 repo，`<migrations_dir>/` 內存在多個 vendor seeder migration，單檔 50KB 至 150KB，內容絕大多數為 `INSERT INTO <table> VALUES (...)`，且同類資料反覆出現於後續 migration
- Evidence path: 證據留在 `<PROJECT_ROOT>/` 原始 review note，不複製到本庫

#### Generalized Lesson

依資料性質選擇對應載入路徑，不一律用 schema migration：

| 資料性質 | 路徑 |
|------|------|
| 啟動必需的少量列舉 | 保留 migration，但每筆獨立 |
| 動態目錄（廠商、商品、SKU） | Application-level seeder 或 admin UI 為 source-of-truth |
| 環境差異 fixture | 獨立 seeder pipeline，不放 migration |
| 大量參考資料 | 外部 CSV/Parquet + bulk loader |
| 可變配置 | Config service / feature flag / 環境變數 |

#### Agent Action

看到後端 repo 的 schema migration 目錄時：

1. `find <migrations_dir> -name "*.sql" -size +50k` 列出大型 migration
2. 對大型 migration 檢查內容比例（`INSERT` vs `CREATE / ALTER`）
3. 若 `INSERT` 比例高，標記為 anti-pattern candidate
4. 引導使用者依資料性質選擇替代載入路徑，不要把大量業務資料塞進 schema migration

#### Goal / Action / Validation

- Goal: 業務資料的 lifecycle 與 schema lifecycle 解耦
- Action: 依資料性質選擇載入路徑；現有大型 seeder 評估抽出到 application seeder 或外部資料檔
- Validation or reference source: 大型 migration 數量 = 0 或可解釋；無 migration 同時做 schema change + 大量 data insert；業務目錄資料有文件化 source-of-truth 與修改路徑

#### Applies When

- 後端服務含 schema migration（Flyway / Liquibase / Alembic / golang-migrate 等）
- 業務目錄資料量會增長或頻繁變動
- 多環境部署（dev / staging / prod）

#### Does Not Apply When

- 純 enum / lookup table，量小且穩定
- 一次性遷移（如系統初始化）且確認後續不會擴充

#### Validation

- 對既有 repo 跑 size check，能找出實際違反案例
- 替代方案有明確的 source-of-truth 文件與修改路徑

#### Promotion Target

- ✅ `intelligence/engineering/anti-patterns/migration-seeder-anti-patterns.md`（已於 commit `d5ec684` 寫入）
- ✅ `knowledge/summaries/migration-seeder-anti-patterns.md`（已於 commit `d5ec684` 寫入）

#### Required Linked Updates

- ✅ `intelligence/engineering/anti-patterns/README.md` 索引（已於 `d5ec684` 更新）
- ✅ `knowledge/summaries/README.md` 索引（已於 `d5ec684` 更新）
- Step 6（Intelligence Extraction）不適用：本 lesson 為反向補登，intelligence atom 已在 `d5ec684` 直接寫入，不需重新 extract
- Step 7（Failure Learning）：本 lesson 本身不適用，但本輪 session 跳過 knowledge-update-flow 的失誤已由既有 `enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md` 覆蓋
