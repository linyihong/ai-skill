# Vendor Integration Architecture（多廠商整合架構）

**Status**: `candidate-intelligence`
**Source**: 通用後端架構經驗（適用於支付聚合、社群登入、IM、雲廠商 SDK、博弈聚合、廣告聯播、地圖等場景）

## 原則

**整合超過 3 個外部廠商時，先在 vendor SPI、Strategy、Compile-time submodule、Out-of-process service 之間選型，不要預設「每個廠商一個編譯期模組」。** 模組數量隨廠商數線性成長的架構，會在 N 達兩位數時開始崩壞，三位數時失控。

## 為什麼

1. **編譯時間爆炸**：N 個 module 對應 N 套相依、N 個 build target。`mvn install` / `gradle build` 時間隨 N 線性甚至超線性成長。
2. **IDE 載入崩潰**：IntelliJ / VS Code 對 100+ Maven module 的 indexing 與 navigation 顯著退化。
3. **升級風險擴散**：升級共用框架（Spring Boot、序列化器）必須觸及每個 vendor 模組，PR diff 龐大、review 困難。
4. **架構不一致**：不同廠商由不同團隊／不同時期實作，導致同一 repo 內有 N 種整合風格（有的 4 層拆分、有的平鋪、有的只是 test 殼），認知負擔高。
5. **熱換與動態啟用困難**：編譯期綁定的廠商無法在不重啟服務的情況下加入、停用或灰度發布。

## 整合策略選項

| 策略 | 形式 | 何時適合 |
|------|------|------|
| **A. Adapter / Strategy（單模組內多實作）** | 一個 vendor interface，N 個 class 實作，共享同 Maven module | N ≤ 5，廠商行為高度相似，無需熱換 |
| **B. Compile-time submodule per vendor** | 每個廠商一個 Maven/Gradle 子模組 | 廠商 SDK 相依嚴重衝突，需要 classpath 隔離，且 N 可控（< 10） |
| **C. Plugin / SPI（runtime 載入）** | 共用 interface 在 core，廠商實作打包成 JAR/插件，啟動時用 `ServiceLoader` / OSGi / 自製 plugin loader 載入 | N 中大（10~100），廠商頻繁增減，需要動態啟用 |
| **D. Out-of-process（廠商各自服務）** | 每家廠商獨立 process / container，主服務透過 RPC/gRPC/HTTP 呼叫 | 廠商 SDK 不穩定（記憶體洩漏、阻塞）、需要獨立 scale 或合規隔離 |
| **E. Hybrid (Tier strategy)** | 主流廠商用 A，長尾廠商用 C，特殊隔離需求用 D | N 三位數以上、廠商重要性差異大 |

## 軸向（選型 trade-off）

| 軸 | 偏 A/B（compile-time） | 偏 C/D（runtime） |
|----|------|------|
| 廠商數量 N | 小 | 大 |
| 新增廠商頻率 | 低 | 高 |
| 廠商 SDK 相依衝突 | 不嚴重 | 嚴重 |
| 熱換 / 動態啟用需求 | 無 | 有 |
| 廠商行為相似度 | 高（同 interface） | 任意 |
| Ops 成熟度 | 任意 | 高（須有插件管理、隔離） |
| 法規 / 合規隔離 | 任意 | 需要時必須 D |

## 何時必須跳出「每個廠商一個編譯模組」

當任一條成立：

- N ≥ 10 且預期繼續增長
- IDE 載入時間 ≥ 30 秒或卡頓明顯
- 升級共用框架需要動 ≥ 5 個 vendor 模組
- 新增廠商的平均工程成本 ≥ 1 週（含建置、CI、deploy 流程）
- 廠商需要灰度／A/B、停用、降級切換而不能重啟服務

→ 應遷移到 **C（SPI）** 或 **E（混合）**。

## 決策流程

```text
要整合新廠商前先問：

1. 目前已有幾個廠商？
   < 5 → 用 A（Adapter/Strategy）
   5-10 → 評估行為相似度與 SDK 相依
   > 10 → 強烈考慮 C 或 E

2. 廠商 SDK 之間有相依衝突嗎（同一函式庫不同版本）？
   是 → B 或 C（classpath 隔離）
   否 → A 可接受

3. 需要不重啟動態啟停？
   是 → C 或 D
   否 → A 或 B

4. 廠商穩定性差或合規要求隔離？
   是 → D
   否 → 上述任一
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「Maven module 就是天然的隔離邊界，每家廠商一個」 | Module 邊界用來描述**穩定的程式碼組織**，不該隨外部廠商數量線性膨脹 |
| 「插件架構太複雜，先 compile-time 一陣子」 | 廠商整合通常只增不減；一旦 N ≥ 10，重構成 SPI 的成本遠高於一開始就採 SPI |
| 「共用一個 cross-core 模組就算插件化了」 | 共用 base class 不等於 plugin。Plugin 的關鍵是**獨立 lifecycle、獨立載入、獨立啟停**，不是繼承 |
| 「為了清楚把每個產品 SKU/遊戲也拆模組」 | SKU/遊戲是**資料**，不是**程式碼分群**。應該是 config / catalog，不是 Maven module |

## 演化路徑

從 N 小到 N 大的健康演化：

```text
1. 第 1-2 家廠商：直接寫，不抽象（避免過早抽象）
2. 第 3 家：建立 vendor SPI interface，重構前兩家為 A
3. 第 5-10 家：評估是否需要 classpath 隔離，必要時部分轉 B
4. 第 10+ 家：建立 plugin 載入機制，新廠商走 C
5. 第 50+ 家或穩定性差的廠商：考慮 D（獨立服務）
6. 廠商重要性差距明顯：採 E（混合分層）
```

## 與其他智慧的關係

- [`modular-monolith-vs-microservices.md`](modular-monolith-vs-microservices.md)：當廠商整合走到 D（out-of-process），就是 microservices 拆分決策的一個 trigger。
- [`modularity/`](modularity/README.md)：本檔處理「外部整合的模組策略」，與模組化單體的內部 module 策略互補。
- [`coupling-tradeoffs/`](coupling-tradeoffs/README.md)：選 A/B/C/D 的本質是廠商耦合度與 lifecycle 對齊的選擇。
- [`../anti-patterns/`](../anti-patterns/README.md) 的相關反模式：把資料分群當程式碼分群（例如把 SKU/遊戲拆成 Maven module）。

## 驗證

| 檢查 | 通過條件 |
|------|------|
| 廠商數對應策略 | 廠商數 vs 策略選擇符合上表軸向 |
| 模組數合理性 | 廠商相關 Maven/Gradle module 數 ≤ 廠商數 × 4（base / biz / config / data 等） |
| 編譯時間 | 全量 build < 5 分鐘（CI），增量 build < 30 秒 |
| 新增廠商成本 | 從接到需求到 PR 可 merge < 3 個工作日 |
| 升級框架影響 | 升級 Spring Boot / 序列化器等共用相依不需要每個廠商模組 reconfigure |

## Token Impact

選錯策略的代價遠超過初期的「多想一下」。N=100 的編譯期模組架構，重構回 SPI 通常需要 1-3 個工程師月。早期投入 1-2 週設計 SPI，回報是長期的線性成本控制。

---

← [回到 engineering/architecture/](README.md)
