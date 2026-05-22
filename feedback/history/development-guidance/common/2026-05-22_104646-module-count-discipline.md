> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-22 — Module Count Discipline（模組數量紀律）

Status: candidate

#### One-line Summary

Repo 內模組（build system module / workspace package / subproject）數量隨業務需求增長，但「畫一個模組」有工程成本；N 失控時 build time、IDE 載入、升級擴散、refactor 成本會線性以上增長。本 lesson 沉澱跨 build system 通用的模組數量管理原則。

#### Human Explanation

模組劃分原本是「物理隔離」工具，但常被誤用為「資料分類」工具。當 SKU、商品、遊戲、廠商目錄被拆成編譯期模組時，每一筆資料變動都觸發整套 module 機制（build、相依、版本、IDE 索引），代價遠超過用 config / catalog 表達相同分類。

當 N 達到兩位數開始有 build / IDE 退化訊號，三位數則幾乎必然是「資料當程式碼分群」的反模式。每個 build system（Maven、Gradle、npm workspace、Cargo workspace、Go multi-module、Bazel、.NET solution、pnpm / Lerna / Nx / poetry workspace）都有相同的數量爆炸風險，原則跨工具中立。

#### Trigger

- 觀察到 repo 模組數 N ≥ 30 且仍在增長
- CI 全量 build > 5 分鐘 或本機增量 > 30 秒
- IDE 載入 / indexing 卡頓
- 升級共用框架（編譯器、序列化器、ORM）需要動 ≥ 5 個模組
- 看到「每個廠商 / SKU / 商品 / 遊戲 / 地區一個 module」的結構
- 同 repo 內模組整合風格不一致（有的多層拆、有的平鋪、有的 test-only 殼）

#### Evidence

- Tool: Repo 結構觀察、build system 配置檔閱讀
- Sanitized excerpt: 多廠商後端整合 repo 觀察到對某一廠商展開三位數量級「每個資料項一個 module」的結構，跨廠商整合風格不一致；屬於本 lesson 的典型案例
- Evidence path: 證據留在 `<PROJECT_ROOT>/` 原始 review note，不複製到本庫

#### Generalized Lesson

模組數量是工程成本，不是免費分類工具。N 對應健康做法：

| N | 健康做法 |
|---|------|
| ≤ 5 | 直接寫，不抽象 |
| 5-10 | Module 模板 + 依賴方向審查 |
| 10-30 | 明確新增門檻、共用 platform module |
| 30-50 | Build cache、incremental、評估 plugin/SPI |
| 50-100 | 強制 plugin / SPI / catalog 替代 |
| ≥ 100 | 結構味道警報，多半要重構或拆 repo |

**新增 module 的判斷流程**（依序問）：
1. 有獨立 lifecycle / version / release cadence？
2. 有 classpath 相依衝突？
3. 會被外部使用？
4. 加進去後 N 會超警戒？
5. 本質是「資料」還是「程式碼」？

任一答「是資料」→ 用 config / catalog / 資料表，不要用編譯期模組。

#### Agent Action

在後端 / monorepo / 多模組專案的 review 中：

1. 列出 module 數 N 與分布（Maven / Gradle / workspace / etc.）
2. 對照 N 警戒線給出評估
3. 找出「資料當程式碼分群」訊號（廠商 / SKU / 商品等以 module 為單位）
4. 若存在，提出 plugin / SPI / catalog 替代方案
5. 引導建立 module 新增門檻文件，避免 N 持續失控

#### Goal / Action / Validation

- Goal: N 在可承受範圍，build / IDE / 升級成本線性可控
- Action: 建立新增 module 的判斷門檻；既有資料分群型 module 評估替代結構
- Validation or reference source: Build 時間在警戒線內；升級共用相依影響範圍 ≤ N × 10%；無「廠商 / SKU / 商品」為單位的編譯期 module

#### Applies When

- 任何使用編譯期 module 結構的 build system（Maven、Gradle、npm workspace、Cargo workspace、Go multi-module、Bazel、.NET solution、pnpm、Lerna、Nx、poetry workspace 等）
- Repo 預期長期維護且模組數會增長

#### Does Not Apply When

- 純單模組專案（無 sub-module 結構）
- 短期 PoC、demo 專案
- 純 monorepo 但各專案完全獨立（不互相 import）

#### Validation

- 找得到實際案例的 N 與健康做法對應
- 對既有 repo 跑 N 計算與警戒線比對，能定位問題模組

#### Promotion Target

- `intelligence/engineering/architecture/modularity/module-count-discipline.md`（本次新增）
- `knowledge/summaries/module-count-discipline.md`（summary card）

#### Required Linked Updates

- `intelligence/engineering/architecture/modularity/README.md`（加入索引條目）
- `knowledge/summaries/README.md`（加入 summary 條目）
- Step 6（Intelligence Extraction）：done(executed) — 本 lesson 同 commit promote 為 intelligence atom
- Step 7（Failure Learning）：not_applicable — 來源為通用工程原則與專案觀察，非 agent failure 補救
