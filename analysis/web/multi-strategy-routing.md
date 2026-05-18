# Multi-strategy Routing for Web Scraping

## 問題

不同的網頁提取任務需要不同的策略組合。單一策略無法應對所有場景：

- 靜態頁面用 HTTP client 最快
- 動態頁面需要 headless browser
- 有 anti-bot 保護的頁面需要 stealth 技術
- 大量頁面需要並發控制
- 需要登入的頁面需要 session 管理

## 策略分類

### 1. 提取策略（Fetch Strategy）

| 策略 | 工具 | 適用場景 | 速度 | 隱蔽性 |
|------|------|---------|------|--------|
| **HTTP Direct** | Scrapling Fetcher / httpx | 靜態 HTML，無 anti-bot | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **Dynamic Browser** | Scrapling DynamicFetcher / Playwright | SPA、JS 渲染、AJAX | ⭐⭐⭐ | ⭐⭐⭐ |
| **Stealth Browser** | Scrapling StealthyFetcher | Cloudflare、Akamai、reCAPTCHA | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **API Direct** | 直接呼叫網站 API | 網站有公開或未公開 API | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

### 2. 解析策略（Parse Strategy）

| 策略 | 工具 | 適用場景 | 穩定性 |
|------|------|---------|--------|
| **CSS Selector** | Scrapling `.css()` | 結構穩定的頁面 | ⭐⭐⭐ |
| **XPath** | Scrapling `.xpath()` | 複雜 DOM 導航 | ⭐⭐⭐ |
| **Text Content** | Scrapling `.find_by_text()` | 文字比對 | ⭐⭐⭐⭐ |
| **Regex** | Scrapling `.find_by_regex()` | 模式比對（email、電話等） | ⭐⭐⭐⭐⭐ |
| **Adaptive Parsing** | Scrapling `.find_similar()` | 結構會變動的頁面 | ⭐⭐⭐⭐ |
| **AI-assisted** | 將提取結果傳給 LLM | 非結構化內容、需要理解語意 | ⭐⭐ |

### 3. Session 策略（Session Strategy）

| 策略 | 適用場景 | 資源消耗 |
|------|---------|---------|
| **Stateless** | 一次性請求，無需登入 | 低 |
| **Single Session** | 需要登入、cookie 維持 | 中 |
| **Multi-session** | 多帳號、多配置並存 | 高 |
| **Session Pool** | 大量並發請求，需要輪換 | 高 |

### 4. 並發策略（Concurrency Strategy）

| 策略 | 適用場景 | 風險 |
|------|---------|------|
| **Sequential** | 少量頁面，無需加速 | 低 |
| **Bounded Parallel** | 中量頁面，可控並發數 | 中 |
| **Full Parallel** | 大量頁面，最快速度 | 高（可能被 ban） |
| **Rate-limited** | 需要遵守 robots.txt 或 rate limit | 低 |

## 路由決策流程

```
輸入：目標 URL + 任務需求
  │
  ├─ 1. 網站評估
  │    ├─ 靜態 HTML？ → HTTP Direct
  │    ├─ 需要 JS 渲染？ → Dynamic Browser
  │    └─ 有 anti-bot？ → Stealth Browser
  │
  ├─ 2. 解析策略選擇
  │    ├─ 結構穩定？ → CSS Selector / XPath
  │    ├─ 結構變動？ → Adaptive Parsing
  │    └─ 非結構化？ → AI-assisted
  │
  ├─ 3. Session 策略
  │    ├─ 需要登入？ → Session
  │    └─ 一次性？ → Stateless
  │
  └─ 4. 並發策略
       ├─ 少量 → Sequential
       ├─ 中量 → Bounded Parallel
       └─ 大量 → Rate-limited / Session Pool
```

## 策略組合範例

### 範例 1：爬取新聞網站

```python
# 評估：靜態 HTML，無 anti-bot，結構穩定
# 策略組合：HTTP Direct + CSS Selector + Stateless + Sequential
page = scrapling.get('https://news.example.com')
titles = page.css('.article-title').text()
```

### 範例 2：爬取 SPA 電商網站

```python
# 評估：React SPA，需要 JS 渲染，有 rate limiting
# 策略組合：Dynamic Browser + Adaptive Parsing + Session + Rate-limited
async with scrapling.DynamicFetcher() as fetcher:
    page = await fetcher.fetch('https://shop.example.com/products')
    products = page.find('.product-card').find_similar()
```

### 範例 3：爬取有 Cloudflare 保護的網站

```python
# 評估：Cloudflare 保護，需要登入，大量頁面
# 策略組合：Stealth Browser + XPath + Multi-session + Bounded Parallel
async with scrapling.StealthyFetcher() as fetcher:
    page = await fetcher.stealthy_fetch('https://protected.example.com/data')
    items = page.xpath('//div[@class="data-item"]')
```

### 範例 4：MCP Server 自動路由

```python
# Scrapling MCP Server 會根據 URL 自動選擇策略
# Agent 只需指定目標和需要的資料
"""
使用 Scrapling MCP 工具：
- get() → 自動選擇 HTTP Direct
- fetch() → 自動選擇 Dynamic Browser
- stealthy_fetch() → 自動選擇 Stealth Browser
"""
```

## 策略切換條件

當以下情況發生時，應考慮切換策略：

| 觸發條件 | 目前策略 | 切換至 |
|---------|---------|--------|
| HTTP Direct 回傳狀態碼 403/503 | HTTP Direct | Stealth Browser |
| Dynamic Browser 被檢測為 bot | Dynamic Browser | Stealth Browser |
| 解析結果為空或異常 | 任何解析策略 | Adaptive Parsing |
| Rate limiting 觸發 | 任何並發策略 | Rate-limited |
| Session 過期 | 任何 Session 策略 | 重新建立 Session |

## 與既有知識的關係

- [`analysis/web/README.md`](../README.md) — Web Scraping 分析方法的範圍
- [`analysis/web/sources-and-tools.md`](sources-and-tools.md) — 工具選擇參考
- [`intelligence/web-scraping/adaptive-parsing.md`](../../intelligence/web-scraping/adaptive-parsing.md) — Adaptive Parsing 作為解析策略之一
- [`analysis/web/mcp-server-patterns.md`](mcp-server-patterns.md) — MCP Server 的 3-Layer 架構對應提取策略
- [`enforcement/sanitization.md`](../../enforcement/sanitization.md) — Prompt Injection Protection
