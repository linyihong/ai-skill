# 目錄結構治理（Directory Structure Governance）

## 為什麼需要這個流程

目錄結構是知識 OS 的骨架。一個好的目錄結構讓 agent 能直覺找到正確的檔案，一個不好的目錄結構會導致：

- **路徑混淆**：兩個名字相似的目錄讓 agent 無法判斷該放哪裡
- **慣性命名**：沿用舊 skill 的名稱，而不是反映內容的本質
- **邊界模糊**：目錄名稱無法清楚表達「這裡面放什麼」和「這裡面不放什麼」
- **重構成本累積**：每次重構都要更新上百個外部引用

本文件定義一套 checkpoint 流程，讓 agent 在建立或重構目錄時能自我檢查。

## Checkpoint：目錄結構合理性檢查

在以下時機執行此 checkpoint：

1. **建立新目錄時** — 任何新目錄建立前
2. **重構既有目錄時** — 重新命名、合併、拆分目錄前
3. **發現路徑混淆時** — 當 agent 或使用者無法直覺判斷某個檔案該放哪裡
4. **每輪知識更新 checkpoint 中** — 作為可選的結構健康檢查

### Step 1：名稱衝突檢查

檢查新目錄名稱是否與既有目錄（尤其是不同層的目錄）產生混淆：

```
問題：新目錄名稱是否與以下任一衝突？
├── 同層其他目錄（如 engineering/ 下的兄弟目錄）
├── 不同層但同名目錄（如根目錄 analysis/ vs engineering/analysis/）
├── 語意相近但邊界不同的目錄（如 intelligence/engineering/heuristics/ vs enforcement/failure-patterns/）
└── 已棄用或正在遷移的舊目錄
```

**判斷原則**：
- 如果兩個目錄名稱相同或極相似，但它們的內容邊界（放什麼/不放什麼）不同，**名稱需要區分**。
- 好的名稱應該讓 agent 在不讀 README 的情況下，也能大致猜出內容類型。
- 名稱應該反映「內容的本質」（intelligence / method / workflow / rule），而不是「來源 skill 的名稱」。

### Step 2：邊界清晰度檢查

檢查目錄的「放什麼 / 不放什麼」邊界是否清晰：

```
問題：以下邊界是否清晰？
├── 與上層目錄的邊界（這個目錄的內容為什麼不能直接放在上層？）
├── 與兄弟目錄的邊界（這個目錄和隔壁目錄的差別是什麼？）
├── 與同名跨層目錄的邊界（如果有的話）
└── 與 enforcement/ 的邊界（工程判斷 vs 可執行 policy）
```

**判斷原則**：
- 如果無法用一句話說明「這個目錄和那個目錄的差別」，邊界就不夠清晰。
- 如果邊界需要依賴「歷史原因」或「因為以前就是這樣」來解釋，需要重新思考。
- 如果兩個目錄的「放什麼」有大量重疊，考慮合併或重新劃分。

### Step 3：慣性命名檢查

檢查目錄名稱是否受到「舊 skill 名稱」或「既有習慣」的影響：

```
問題：這個命名是否有以下慣性偏誤？
├── 直接沿用舊 skill 名稱（apk-analysis → analysis，但內容本質不是 analysis 方法）
├── 使用過於通用的名稱（analysis、common、misc、other）
├── 名稱包含實作細節而非抽象概念（frida-hook 而非 dynamic-instrumentation）
├── 名稱反映「來源」而非「用途」（from-apk-skill 而非 analysis-intelligence）
└── 名稱與實際內容不符（目錄叫 analysis 但內容是 intelligence）
```

**判斷原則**：
- 目錄名稱應該反映「這個目錄的內容是什麼」，而不是「這個內容從哪裡來」。
- 如果一個目錄的內容本質上屬於另一層（例如 intelligence 放到 analysis 目錄下），名稱應該反映真實層級。
- 當 rename 時，優先考慮「內容本質」而非「最短名稱」。

### Step 4：路徑深度與扁平化檢查

```
問題：這個目錄的位置是否合理？
├── 是否放在正確的層級？（intelligence/ 下的 engineering/ 是對的，但 analysis/ 呢？）
├── 路徑深度是否超過 4 層？（intelligence/engineering/analytical-reasoning/heuristics/ 是 4 層）
├── 是否可以在不破壞既有引用的情況下移動？
└── 是否有其他更直覺的位置？
```

**判斷原則**：
- 路徑深度建議不超過 4 層（`層/領域/類別/檔案`）。
- 如果一個目錄的內容可以被多個上層共用，考慮提升一層。
- 如果一個目錄只有一個子目錄或檔案，考慮合併到上層。

### Step 5：全域引用影響評估

在執行任何目錄重構前，評估影響範圍：

```
影響評估：
├── 有多少外部檔案引用這個路徑？（grep -rl "old/path"）
├── 這些引用分布在哪些層？（architecture/、knowledge/、runtime/、skills/ 等）
├── 是否有 binary 或 generated 檔案包含這個路徑？（如 SQLite index）
├── 是否有 tool adapter 或設定檔依賴這個路徑？
└── 重構後是否需要更新 validator 或 runtime surface？
```

**判斷原則**：
- 如果外部引用超過 50 個，考慮分階段更新或使用 symlink 過渡。
- 如果 SQLite 或 generated index 包含舊路徑，需要在重構後重建。
- 如果 tool adapter 依賴舊路徑，需要更新 adapter 或加入相容性處理。

## 常見反模式

### 反模式 1：同名跨層目錄

```
analysis/                    # 根目錄：分析方法（how to observe）
intelligence/engineering/
  analysis/                  # intelligence 層：分析工程智慧（why & principles）
```

**問題**：兩個目錄名稱完全相同，但內容本質不同（一個是 method，一個是 intelligence）。Agent 無法從路徑名稱判斷該去哪個。

**解決方案**：
- 改名為反映真實內容的名稱，例如 `intelligence/engineering/analytical-reasoning/`
- 或在 `intelligence/engineering/` 下使用更能區分的前綴

### 反模式 2：技能名稱直接映射

```
# 舊結構（skill-specific）
intelligence/engineering/apk-analysis/

# 新結構（直接改名但沒反映本質）
intelligence/engineering/analytical-reasoning/
```

**問題**：只是把 `apk-analysis` 縮短為 `analysis`，但沒有思考「這個目錄的內容本質是什麼」。如果內容本質是 intelligence（工程智慧），就不該用 analysis（分析方法）這個詞。

**解決方案**：
- 先分析內容本質：這些檔案是「分析方法」還是「從分析中萃取的智慧」？
- 如果是後者，名稱應該用 intelligence 相關的詞彙

### 反模式 3：過度拆分

```
intelligence/engineering/
  apk-analysis/
    heuristics/
    signals/
    failure/
    anti-patterns/
  app-development-guidance/
```

**問題**：以 skill 為單位拆分 intelligence，而不是以「知識類型」為單位。這導致：
- 每個 skill 都重複 heuristics/signals/failure/anti-patterns 子目錄
- 跨 skill 的通用 intelligence 沒有位置
- 新 skill 加入時需要建立整套目錄結構

**解決方案**：
- 以知識類型為單位組織（heuristics/、signals/、failure/、anti-patterns/ 是全局分類）
- 在檔案層級用 frontmatter 或 naming convention 標記 domain
- 跨 domain 的 intelligence 放在通用目錄

## 與既有流程的關係

| 流程 | 關係 |
|------|------|
| [`knowledge-update-flow.md`](knowledge-update-flow.md) | 本 checkpoint 可作為 Step 1（觸發檢查）的可選子檢查 |
| [`dependency-reading.md`](../../enforcement/dependency-reading.md) | Step 5 的影響評估需要遵守 dependency reading 規則 |
| [`linked-updates.md`](../../enforcement/linked-updates.md) | 目錄重構後的 linked updates 需要遵守此規則 |
| [`content-layering.md`](../../enforcement/content-layering.md) | 目錄結構需要與內容分層原則一致 |
| [`intelligence/README.md`](../../intelligence/README.md) | 定義了 intelligence 層的放什麼/不放什麼 |
| [`analysis/README.md`](../../analysis/README.md) | 定義了 analysis 層的放什麼/不放什麼 |

## 驗證

1. 新目錄建立後，檢查是否有同名或語意衝突的既有目錄。
2. 新目錄的 README 必須包含「放什麼」和「不放什麼」。
3. 重構完成後，執行 `grep -rl "old/path"` 確認無殘留引用。
4. 如果 SQLite runtime index 存在，重建 index。
5. 更新 `knowledge/runtime/routing-registry.yaml` 和 `knowledge/indexes/README.md`。
