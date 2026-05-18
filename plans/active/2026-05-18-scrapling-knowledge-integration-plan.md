# Scrapling Knowledge Integration Plan

## 背景

分析 [Scrapling](https://github.com/D4Vinci/Scrapling)（Python Web Scraping 框架）後，發現多個我們系統目前缺乏的知識點與設計模式。本計畫定義哪些要加入、加入哪裡、以及執行順序。

## 提案總覽

| # | 知識點 | 類型 | 優先級 | 預計文件數 | 預計工時 |
|---|--------|------|--------|-----------|---------|
| 1 | Web Scraping 分析方法（analysis 入口） | `analysis/web/` | 🔴 P0 | 1 分析入口 + 1 README | 1 session |
| 2 | Adaptive Parsing（自適應解析） | `intelligence/` | 🔴 P0 | 1 知識文件 | 1 session |
| 3 | MCP Server Design Pattern | `intelligence/` | 🔴 P0 | 1 知識文件 | 1 session |
| 4 | Prompt Injection Protection 強化 | `enforcement/` | 🟡 P1 | 修改 1 現有文件 | 0.5 session |
| 5 | Multi-session / Multi-strategy Routing | `intelligence/` | 🟡 P1 | 1 知識文件 | 1 session |
| 6 | Cache & Replay Pattern | `intelligence/` | 🟢 P2 | 1 知識文件 | 0.5 session |
| 7 | Anti-bot Bypass Techniques | `intelligence/` | 🟢 P2 | 1 知識文件 | 0.5 session |

## 詳細規劃

### P0：必須加入

#### 1. Web Scraping 分析方法（analysis/web/）

**位置**：`analysis/web/README.md` + `analysis/web/sources-and-tools.md`

**理由**：`analysis/` 目前有 `apk/`、`development-guidance/`、`repo/`、`travel/`、`production/`、`issue/`，但沒有 `web/`。Scrapling 的知識本質上是「如何分析與爬取網頁」，這是一個獨立的分析領域，應該與其他 analysis 入口平級。

**`analysis/web/README.md` 內容**：
- Web Scraping 分析方法的範圍與邊界
- 何時進入此分析領域（需要從網頁提取資料、網站有 anti-bot 保護、需要動態渲染內容等）
- 與其他 analysis 領域的關係（例如 `repo/` 分析 repo 結構，`web/` 分析網頁內容）
- 可參考的工具與技術（Scrapling、Playwright、HTTP clients 等）

**`analysis/web/sources-and-tools.md` 內容**：
- 網頁分析的資料來源類型（靜態 HTML、SPA、API-driven sites）
- 工具選擇策略（何時用 HTTP request、何時用 browser、何時用 stealth mode）
- Anti-bot 保護的識別與應對
- 與既有 `analysis/travel/sources-and-tools.md` 類比設計

**不包含**：
- 不複製 Scrapling 的完整文件
- 不包含實作細節（那是實作層的事）

#### 2. Adaptive Parsing（自適應解析）

**位置**：`intelligence/web-scraping/adaptive-parsing.md`

**內容**：
- 問題：網站改版後 selector 失效
- Scrapling 的解法：DOM tree depth + tag name hierarchy + fuzzy attribute matching
- 我們可以怎麼用：在 analysis 領域中，當目標網站/APP 改版時，agent 可以自動重新定位元素
- 適用場景：任何需要定期爬取或分析的任務
- 與既有知識的關係：可參考 `intelligence/engineering/` 下的分析模式

**不包含**：
- 不實作實際的 fuzzy matching 演算法（那是 library 層的事）
- 不複製 Scrapling 的程式碼

#### 3. MCP Server Design Pattern

**位置**：`intelligence/engineering/mcp-server-patterns.md`

**內容**：
- Scrapling MCP Server 的架構：10 個工具分三層（basic / dynamic / stealth）
- 關鍵設計決策：
  - CSS selector 精準提取 → 減少 token 消耗
  - Session persistence → 避免重複啟動 browser
  - Bulk operations → 平行處理多 URL
  - Screenshot as ImageContent → 模型可直接看到圖片
- 我們可以怎麼用：如果未來需要讓 AI agent 與外部工具互動
- 與既有文件的關係：可參考 `ai-tools/` 下的工具使用說明

**不包含**：
- 不實作 MCP server（那是實作層的事）
- 不複製 Scrapling 的 MCP 工具程式碼

### P1：建議加入

#### 4. Prompt Injection Protection 強化

**位置**：修改 `enforcement/sanitization.md`

**新增內容**：
- CSS-hidden elements sanitization（`display:none`, `visibility:hidden`, `opacity:0`）
- Zero-width unicode characters sanitization
- Template tags sanitization
- HTML comments sanitization
- `aria-hidden` elements sanitization

#### 5. Multi-session / Multi-strategy Routing

**位置**：`intelligence/engineering/multi-strategy-routing.md`

**內容**：
- Scrapling 的做法：同一個 spider 內混用 HTTP session + stealth browser session
- 抽象模式：根據任務特徵路由到不同的執行策略
- 我們可以怎麼用：在 workflow 中，不同類型的子任務用不同的執行策略
- 與既有文件的關係：可參考 `workflow/software-delivery/` 的流程設計

### P2：可選加入

#### 6. Cache & Replay Pattern

**位置**：`intelligence/engineering/cache-replay-pattern.md`

**內容**：
- Development mode 的概念：第一次執行時 cache，後續 replay
- 適用場景：重複性分析任務、測試、開發迭代
- 與既有文件的關係：可參考 `enforcement/failure-learning-system.md`

#### 7. Anti-bot Bypass Techniques

**位置**：`intelligence/web-scraping/anti-bot-bypass.md`

**內容**：
- TLS fingerprint impersonation
- Cloudflare Turnstile 自動化解法
- WebRTC blocking
- Canvas noise / fingerprint spoofing
- DNS-over-HTTPS 防止 DNS leak

## 不加入的項目

以下 Scrapling 的功能**不適合**加入我們的系統：

| 功能 | 原因 |
|------|------|
| CLI extract 指令 | 我們的系統不是爬蟲工具 |
| Interactive Web Scraping Shell | 同上 |
| Spider framework（CrawlSpider / SitemapSpider） | 太特定於爬蟲領域 |
| Proxy rotation 實作細節 | 太特定於爬蟲領域 |
| Browser automation（Playwright） | 太特定於爬蟲領域 |

## 執行順序

```
Phase 1 (P0): analysis/web/ + Adaptive Parsing + MCP Server Design Pattern
  → 建立 analysis/web/README.md + analysis/web/sources-and-tools.md
  → 建立 intelligence/web-scraping/adaptive-parsing.md
  → 建立 intelligence/engineering/mcp-server-patterns.md
  → 更新 analysis/README.md（加入 web/ 入口）
  → 更新 knowledge/graphs/
  → 更新 knowledge/runtime/routing-registry.yaml（加入 route.analysis.web）
  → 執行 compiler + validator

Phase 2 (P1): Prompt Injection Protection 強化 + Multi-strategy Routing
  → 修改 enforcement/sanitization.md
  → 建立 intelligence/engineering/multi-strategy-routing.md
  → 更新 knowledge/graphs/
  → 執行 compiler + validator

Phase 3 (P2): Cache & Replay + Anti-bot Bypass
  → 建立 intelligence/engineering/cache-replay-pattern.md
  → 建立 intelligence/web-scraping/anti-bot-bypass.md
  → 更新 knowledge/graphs/
  → 執行 compiler + validator
```

## 完成條件

- [ ] Phase 1 完成（analysis/web/ 建立 + 2 份知識文件 + graphs/routing 更新 + compiler/validator 通過）
- [ ] Phase 2 完成（1 份修改 + 1 份建立 + graphs 更新 + compiler/validator 通過）
- [ ] Phase 3 完成（2 份知識文件建立 + graphs 更新 + compiler/validator 通過）
- [ ] 所有新文件都經過 linked-updates 檢查
- [ ] Commit + Push

## 與既有文件的關係

- 新 `analysis/web/` 遵循 `analysis/README.md` 的規範
- 新知識文件放在 `intelligence/` 下，遵循 `intelligence/README.md` 的規範
- 修改 `enforcement/sanitization.md` 需遵循 `enforcement/` 的格式
- 所有新文件需在 `knowledge/graphs/` 中有對應的 graph 記錄
- `analysis/README.md` 需加入 `web/` 入口
- `knowledge/runtime/routing-registry.yaml` 需加入 `route.analysis.web`
- 需執行 compiler 更新 `runtime/compiler/embedded_data.rb`
- 需執行 validator 確保一致性
