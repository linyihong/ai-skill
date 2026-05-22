# Module Count Discipline（模組數量紀律）

**Status**: `candidate-intelligence`
**Source**: 通用軟體架構經驗（適用於 Maven module / Gradle subproject / npm workspace / Cargo workspace / Go multi-module / Bazel package / .NET project / pnpm / Lerna / Nx workspace / poetry workspace）

## 原則

**Repo 內模組數量隨業務需求增長，但「畫一個模組」是有成本的工程行為，不是免費的資料分類。** 當 N（模組總數）成長到讓 build time、IDE 載入、升級成本、refactor 成本任一指標線性以上增長時，必須主動約束新增模組的條件，並評估合併與替代結構（plugin / SPI / config / catalog）。

本檔是「畫邊界後的數量管理」，與 [`README.md`](README.md) 的「畫邊界本身」互補。

## 為什麼

1. **Build 成本**：每個模組是一份相依宣告、一份產出、一份測試 lifecycle。N 個模組 = N 倍 build graph，並行度只能緩解一部分。
2. **IDE 與 tooling 成本**：對 100+ 模組的 indexing、navigation、refactor 工具退化明顯，從幾秒到幾十秒甚至卡頓。
3. **升級擴散成本**：升級共用框架、序列化器、編譯器版本需要動每一個模組的設定檔；PR diff 龐大、review 困難、漏改風險高。
4. **認知負擔**：每個模組都是新進工程師要建立 mental model 的單位；不必要的模組增加 onboarding 時間。
5. **規範漂移**：N 個模組由不同時期、不同團隊建立，容易產生 N 種風格（4 層拆分、平鋪、test-only 殼），同 repo 內失去一致性。

## 訊號（何時該警覺）

| 指標 | 警戒線 | 動作 |
|------|------|------|
| 全量 build 時間（CI） | > 5 分鐘 | 評估並行度、incremental、merge 候選 |
| 增量 build 時間（本機） | > 30 秒 | 評估 build cache、module 粒度 |
| IDE 載入時間 | > 30 秒或明顯卡頓 | 評估 lazy load、module 合併、SPI |
| 全 repo 模組數 N | ≥ 30 | 開始建立 module 新增門檻 |
| N ≥ 50 | 強烈考慮 plugin / SPI / catalog 替代方案 |
| N ≥ 100 | 結構味道明顯，多半是把「資料」誤分群為「程式碼」 |
| 升級共用相依的影響範圍 | 需動 ≥ 5 個模組 | 評估共用 platform module |
| 新增一個模組到 PR merge 的 lead time | ≥ 1 工作日 | 流程過重；考慮 plugin loader |
| 同 repo 內不同模組整合風格 | ≥ 2 種風格 | 缺乏 module 模板與審查 |

## 健康演化路徑（依 N 對應做法）

```
N ≤ 5    → 直接寫，不抽象（避免過早抽象）
N 5-10   → 建立 module 模板與審查；定義依賴方向（禁止循環）
N 10-30  → 建立明確 module 新增門檻；共用配置抽到 platform module
N 30-50  → 引入 build cache、incremental build；evaluation 是否該轉 plugin/SPI
N 50-100 → 強制 plugin / SPI / catalog 替代；新增 module 需架構審查
N ≥ 100  → 結構味道警報；多半需要重構或拆 repo（monorepo split）
```

## Trade-off 軸向

| 軸 | 偏「多模組」 | 偏「少模組」 |
|----|------|------|
| 物理隔離強度 | 強（編譯期相依隔離） | 弱（package boundary 靠約定） |
| 重構成本 | 高（跨模組 refactor 需多 PR） | 低（單模組內可全域 refactor） |
| 並行開發友善度 | 高（小單位 merge 衝突少） | 低（單模組多人改易衝突） |
| 新人理解成本 | 高（多模組 mental model） | 低（單模組單入口） |
| Build / IDE 成本 | 高 | 低 |
| 升級擴散範圍 | 大（每模組各動） | 小（單點修改） |
| Hot-swap / 動態啟用 | 不支援（compile-time bound） | 不支援（要轉 plugin） |

## 何時新增模組是對的

| 場景 | 理由 |
|------|------|
| 物理隔離需求 | 廠商 SDK 相依衝突需 classpath 隔離 |
| 獨立 lifecycle | 該模組有獨立 release cadence、版本、owner |
| 強制依賴方向 | 用 module 邊界執行 layering（如 domain 不能 import infrastructure） |
| 重用單元 | 該模組會被多個 consumer 使用（含外部） |
| 測試隔離 | 該模組有大量重型測試需獨立執行 |

## 何時新增模組是錯的

| 場景 | 為什麼錯 | 正確做法 |
|------|------|------|
| 為「整齊」拆 | Package / namespace 已能表達結構 | 用 package boundary |
| 把資料分群當程式碼分群 | SKU、商品、遊戲、廠商目錄等是 data | 用 config / catalog / 資料表 |
| 為「以後可能拆服務」預拆 | YAGNI；多半永遠不拆 | 等真正需要再拆 |
| 同團隊小功能 | 沒有 lifecycle / owner 差異 | 用 feature folder |
| 把 test fixture 變 module | Test 應該跟著被測模組 | 用 test source set |

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「每個 feature 一個 module 才乾淨」 | Feature 是業務分群，常用 package 表達即可；只有需要獨立 lifecycle 才升級為 module |
| 「先拆模組，以後合併也很容易」 | 合併比拆分痛 — 跨 module 的相依、import path、版本綁定都要回退 |
| 「N 多沒關係，CI 機器加速就好」 | CI 成本只是 N 失控的其中一個症狀；IDE、認知、refactor 成本不會因 CI 加速解決 |
| 「把廠商 / SKU / 商品拆模組好做版控」 | 那是 data 該做的事；用 config / migration / catalog，不要用編譯期模組 |
| 「拆得越細越好 review」 | Review 友善的單位是 PR diff，不是 module 數量 |

## 決策流程

```text
要新增一個 module 前先問：

1. 這個東西有獨立 lifecycle / version / release cadence 嗎？
   是 → 評估 module
   否 → 用 package / feature folder

2. 它的相依與既有模組有 classpath 衝突嗎？
   是 → 評估 module 隔離
   否 → 同上

3. 它會被外部（其他 repo / 服務）使用嗎？
   是 → 評估 module 作為 publishable artifact
   否 → 同上

4. 加進去後全 repo N 會超過警戒線嗎？
   是 → 評估 plugin / SPI / catalog 等替代結構
   否 → 可加，但設定 module 模板與檢查

5. 它本質是「資料」（廠商、商品、SKU、遊戲、地區）嗎？
   是 → 用 config / catalog / 資料表，不是 module
   否 → 同上
```

## 與其他智慧的關係

- [`README.md`](README.md)：本目錄是「module boundary、feature slice、package boundary、modular monolith 的判斷」；本檔是其數量管理面，互補不重複。
- [`../modular-monolith-vs-microservices.md`](../modular-monolith-vs-microservices.md)：當 N 達到一定規模並有獨立 scale / deploy 需求時，下一步可能轉微服務；但模組爆炸本身不是拆服務的理由。
- [`../vendor-integration-architecture.md`](../vendor-integration-architecture.md)：本檔的特定場景特化 — 多廠商整合的 N 模組問題；該檔 N ≥ 10 跳出 compile-time module per vendor 的建議是本檔原則的具體實作。
- [`../../anti-patterns/migration-feature-bundling.md`](../../anti-patterns/migration-feature-bundling.md)：Migration 時若把搬遷與新功能 +「順便重新拆模組」一起做，會同時觸發三個 unknown 疊加。
- [`../coupling-tradeoffs/`](../coupling-tradeoffs/README.md)：模組邊界本質是 coupling 與 cohesion 的取捨。

## 驗證

| 檢查 | 通過條件 |
|------|------|
| 模組新增門檻 | Repo 有書面的 module 新增條件（lifecycle / classpath / external use / data-vs-code 判別） |
| Build 時間 | 全量 CI < 5 分鐘、增量本機 < 30 秒（或有明確改善 plan） |
| 升級擴散 | 升級共用相依時，需動的模組數 ≤ N × 10%（共用 platform module 設計成功的訊號） |
| 一致性 | 同 repo 內模組整合風格 ≤ 2 種，且有對應模板 |
| 資料 / 程式碼分群 | 沒有以「廠商 / SKU / 商品 / 遊戲 / 地區」為單位的編譯期模組 |

## Token Impact

「module 拆得越細越好 review」是常見錯覺，會在 N 達到 30-50 後反咬。提前建立 module 新增門檻的工程成本約 1-2 天；事後重構回合理 N 通常需要 1-3 工程師月。早期 discipline 是高槓桿投資。

---

← [回到 modularity/](README.md)
