# Migration Seeder Anti-Patterns

**Status**: `candidate-intelligence`

## 反模式

把大量業務資料（廠商目錄、商品 catalog、遊戲/SKU 清單、權限矩陣、地區字典等）以巨型 SQL `INSERT` 包進 schema migration，使**資料 lifecycle 與 schema lifecycle 被強制綁定**。

## 訊號

- 單一 migration 檔超過 50KB，且絕大多數內容是 `INSERT` 語句。
- 檔名含 `dataSeeder`、`seed`、`fixture`、`initial_data`，但不在獨立的 seeder pipeline 中執行。
- 同一張表反覆出現「補資料」的後續 migration（A 新增 N 筆 → B 修正 5 筆 → C 補欄位 + 重灌）。
- 業務人員想新增/下架一筆資料時，工程師需要寫新 migration、PR、deploy。
- 不同環境（dev / staging / prod）需要不同資料，但同一 migration 對所有環境執行，迫使內含 `IF NOT EXISTS` / `ON CONFLICT` 補丁。
- Migration rollback 會同時 rollback schema 與資料，造成意外資料遺失風險。

## 根本原因

- 把 Flyway / Liquibase / Alembic 等 schema migration 工具當作通用資料載入器使用。
- 缺乏獨立的 seeder / fixture / reference-data pipeline。
- 把「上線時資料庫該有什麼」與「schema 該長什麼樣」誤等同為同一件事。

## 影響

| 維度 | 後果 |
|------|------|
| Schema / 資料耦合 | 修一筆資料要新 migration；rollback schema 會誤刪資料。 |
| 部署視窗佔用 | 大型 `INSERT` 在 prod 跑數十秒到數分鐘，鎖表期間阻塞流量。 |
| Code review 困難 | 數十到上百 KB 的 SQL 無人能逐行 review，異常難以發現。 |
| 環境漂移 | 補丁式 `IF`/`ON CONFLICT` 散落各 migration，環境間真實差異無法觀察。 |
| 業務速度 | 新增一筆資料的 lead time 從分鐘變成 PR → review → deploy 週期。 |
| 真實 source-of-truth 模糊 | 同一份資料同時存在 migration、admin UI、配置檔，無人知道何者權威。 |

## 替代方案

| 資料性質 | 建議 |
|------|------|
| 啟動必需的少量列舉（enum、預設角色） | 保留在 migration，但**每筆獨立、明確版本化** |
| 動態目錄（廠商、商品、遊戲清單、設施列表） | Application-level seeder（啟動時讀 YAML/JSON）或 admin UI 為 source-of-truth |
| 環境差異資料（dev fixture） | 獨立 seeder script，不進 migration |
| 大量參考資料（地區字典、稅率表） | 外部 CSV/Parquet + bulk loader（`COPY` / `LOAD DATA INFILE`），版本化資料檔案而非 SQL |
| 可變配置 | 移出 DB，用 config service / feature flag / 環境變數 |

## 決策流程

```text
要寫新 dataSeeder migration 前先問：

1. 這份資料未來會頻繁變動嗎？
   是 → 不放 migration，改用 application seeder 或 admin UI
   否 → 繼續

2. 這份資料超過 10KB / 100 筆嗎？
   是 → 改用外部資料檔（CSV/JSON）+ bulk loader
   否 → 繼續

3. 不同環境需要不同資料嗎？
   是 → 環境感知的 seeder pipeline，不在 migration
   否 → 可以放 migration，但保持小、原子化、單一 INSERT 集合
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「migration 框架支援 SQL，那資料當然也用它載」 | Schema migration 工具設計目標是 schema 演進，不是 reference data lifecycle |
| 「rollback 機制就靠 migration」 | 業務資料的 rollback 需要獨立的 audit / soft delete，不該綁在 schema 版本 |
| 「只是初始化用一次」 | 「一次」會被廠商擴充、補資料、修錯需求打破，往往演化成上百 KB 的 patch 鏈 |

## 與其他 anti-pattern 的關係

- [`generic-repository-overuse.md`](generic-repository-overuse.md)：兩者皆是把 schema 層工具當業務工具用。
- 與 `enforcement/failure-patterns/` 的 source/mirror drift：類比，業務資料在 migration 與 application config 兩處同時存在時容易漂移。

## 驗證

- `find <migrations_dir> -name "*.sql" -size +50k` 應為 0，或每一筆能解釋為何不可分。
- 沒有任一 migration 同時做 schema change + 大量 data insert。
- 業務目錄資料有明確 source-of-truth 與修改路徑文件（誰能改、改在哪、如何同步到所有環境）。

---

← [回到 engineering/anti-patterns/](README.md)
