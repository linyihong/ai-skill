# Analysis Domain Discovery Gap（分析領域發現缺口）

## Linked Validation Scenarios

- `validate_failure_pattern_validator_coverage` — 檢查每個 failure pattern 的 Linked Validation Scenarios 是否為空
- `validate_new_category_registration` — 驗證 agent 在遇到新領域時是否會檢查 `analysis/README.md` 並建立對應入口

## 症狀

當 agent 分析一個外部 library 或工具時，只考慮將其中的知識放入 `intelligence/`（工程知識），而忽略了該工具可能代表一個全新的**分析領域**，應該在 `analysis/` 下建立對應的入口。

## 具體案例

分析 [Scrapling](https://github.com/D4Vinci/Scrapling)（Python Web Scraping 框架）時：
- ❌ 只想到 Adaptive Parsing → `intelligence/web-scraping/`
- ❌ 只想到 MCP Server → `intelligence/engineering/`
- ✅ 應該先想到：Web Scraping 是一個分析領域 → `analysis/web/`

## 根本原因

1. **Discovery checkpoint 的 `search_sources` 缺少 `analysis/`**：`phase.checkpoint` 的 discovery_targets 只搜尋 `workflow`、`intelligence`、`validation_rules`、`governance`，沒有 `analysis` 類型。所以 agent 不會被引導去檢查 `analysis/` 下是否有對應領域。

2. **`knowledge/indexes/README.md` 缺少對應路由**：即使 discovery 搜尋了 indexes，也沒有「Web Scraping 分析」這條路由。

3. **Agent 的思考捷徑**：看到「library/tool」直接對應到「engineering knowledge」，跳過了「這是否是一個新的分析領域」的檢查。

## 預防方式

### 1. Discovery 層（系統性）

在 `phase.checkpoint` 的 discovery_targets 中加入 `analysis` 類型，搜尋 `analysis/README.md`。

### 2. Indexes 層（系統性）

在 `knowledge/indexes/README.md` 中加入新分析領域的路由。

### 3. Agent 思考流程（個人）

當分析一個外部 library/tool 時，強制執行以下檢查順序：

```
Step 1: 這個 library/tool 代表哪個「領域」？
  → 不是「這是什麼技術」，而是「這是用來分析/觀察/拆解什麼的？」
  → 例如：Scrapling = Web Scraping（分析網頁的領域）

Step 2: 這個領域在 analysis/ 下已經有入口了嗎？
  → 檢查 analysis/README.md 的「目前入口」列表
  → 如果沒有 → 先建立 analysis/<domain>/README.md

Step 3: 從這個 library/tool 可以萃取出哪些「工程知識」？
  → 放到 intelligence/ 下
  → 這些是從分析中學到的可重用 pattern
```

### 檢查清單

- [ ] 這個外部資源代表一個「分析/觀察/拆解」的領域嗎？
- [ ] `analysis/README.md` 的「目前入口」列表包含這個領域嗎？
- [ ] 如果沒有，先建立 `analysis/<domain>/README.md` 再萃取 intelligence
- [ ] Discovery checkpoint 的 search_sources 包含 `analysis/README.md` 嗎？

## 檢測

- 當 agent 只把外部知識放到 `intelligence/` 而沒有檢查 `analysis/` 時，觸發警告
- 定期檢查 `analysis/README.md` 的「目前入口」列表 vs `intelligence/` 下的領域知識是否對應

## 恢復

如果已經犯了這個錯（如本案例）：
1. 在 `analysis/` 下建立對應的領域入口
2. 更新 `analysis/README.md` 的「目前入口」列表
3. 將原本放在 `intelligence/` 的知識標註為「從 analysis/<domain> 萃取」
4. 更新 discovery checkpoint 的 search_sources
5. 更新 `knowledge/indexes/README.md` 的路由

## 修復狀態（2026-05-18）

本 failure pattern 已透過以下變更修復：

| # | 變更 | 檔案 | 狀態 |
|---|------|------|------|
| 1 | Discovery checkpoint 加入 `analysis` 類型 | [`runtime/compiler/embedded_data.rb`](../../runtime/compiler/embedded_data.rb:674) | ✅ 已新增 `analysis` discovery_target，搜尋 `analysis/README.md` |
| 2 | 路由索引加入 Web Scraping analysis | [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md:54) | ✅ 已新增「執行 Web Scraping 分析」路由 |
| 3 | Analysis 入口加入 web/ | [`analysis/README.md`](../../analysis/README.md:12) | ✅ 已新增 `web/` 到「目前入口」列表 |
| 4 | Failure pattern 記錄 | 本檔案 | ✅ 已記錄修復狀態 |

下次遇到類似情境（分析外部 library/tool）時，discovery checkpoint 會自動引導 agent 檢查 `analysis/README.md`，避免再次遺漏。
