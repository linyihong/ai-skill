# Web Scraping 工具與來源

## 工具比較

### 提取工具

| 工具 | 類型 | JS 渲染 | Anti-bot Bypass | 效能 | 適用場景 |
|------|------|---------|-----------------|------|---------|
| **Scrapling** | Python framework | ✅ (DynamicFetcher) | ✅ (StealthyFetcher) | ⭐⭐⭐⭐⭐ | 進階爬蟲、AI 整合、MCP |
| **Playwright** | Browser automation | ✅ | ❌ (需自行處理) | ⭐⭐⭐⭐ | 動態內容、測試自動化 |
| **httpx** | HTTP client | ❌ | ❌ | ⭐⭐⭐⭐⭐ | 簡單靜態頁面 |
| **requests** | HTTP client | ❌ | ❌ | ⭐⭐⭐⭐⭐ | 最簡單的靜態頁面 |
| **Selenium** | Browser automation | ✅ | ⚠️ (可搭配 undetected-chromedriver) | ⭐⭐ | 傳統瀏覽器自動化 |
| **curl / wget** | CLI tools | ❌ | ❌ | ⭐⭐⭐⭐⭐ | 快速測試、腳本 |

### HTML 解析工具

| 工具 | 速度 | Selector 支援 | 自適應解析 | 適用場景 |
|------|------|---------------|-----------|---------|
| **Scrapling** | ⭐⭐⭐⭐⭐ (767x BS4) | CSS + XPath + Text + Regex | ✅ Find Similar Elements | 高效能解析 |
| **lxml** | ⭐⭐⭐⭐⭐ | CSS + XPath | ❌ | 標準解析需求 |
| **BeautifulSoup4** | ⭐⭐ | CSS + 部分 XPath | ❌ | 簡單解析、學習曲線低 |
| **Parsel** | ⭐⭐⭐⭐ | CSS + XPath + Regex | ❌ | Scrapy 生態系 |
| **pyquery** | ⭐⭐⭐⭐ | CSS (jQuery-like) | ❌ | jQuery 使用者 |

### 瀏覽器自動化

| 工具 | 速度 | Stealth | Session 管理 | 並發 |
|------|------|---------|-------------|------|
| **Scrapling DynamicFetcher** | ⭐⭐⭐⭐ | ❌ | ✅ | ✅ |
| **Scrapling StealthyFetcher** | ⭐⭐⭐ | ✅ | ✅ | ✅ |
| **Playwright** | ⭐⭐⭐⭐ | ❌ (需 stealth 套件) | ✅ | ✅ |
| **Puppeteer** | ⭐⭐⭐⭐ | ❌ (需 stealth 套件) | ✅ | ✅ |
| **Selenium** | ⭐⭐ | ⚠️ | ✅ | ⚠️ |

## 選擇原則

### 1. 從最簡單的開始

```
靜態 HTML + 無 anti-bot → httpx + Scrapling 解析
    ↓ 需要 JS 渲染
DynamicFetcher (Scrapling) / Playwright
    ↓ 有 anti-bot
StealthyFetcher (Scrapling)
    ↓ 需要 AI agent 整合
Scrapling MCP Server
```

### 2. 根據資料量選擇

- **少量（< 100 頁）**：單次請求，sequential
- **中量（100-10000 頁）**：Session + 並發控制
- **大量（> 10000 頁）**：Spider（Scrapling Spider / Scrapy）+ 分散式架構

### 3. 根據反爬強度選擇

| 反爬等級 | 特徵 | 推薦工具 |
|---------|------|---------|
| **無** | 直接回傳 HTML | httpx / requests |
| **低** | User-Agent 檢查 | Scrapling Fetcher + 自訂 headers |
| **中** | Rate limiting、Cookie 驗證 | Scrapling Session + retry |
| **高** | Cloudflare、reCAPTCHA、JS challenge | Scrapling StealthyFetcher |
| **極高** | 行為分析、指紋辨識 | StealthyFetcher + proxy rotation + 人類行為模擬 |

### 4. 根據整合需求選擇

- **AI Agent 整合**：Scrapling MCP Server（10 tools, 3 layers）
- **CI/CD Pipeline**：Scrapling CLI + Docker
- **互動式開發**：Scrapling Interactive Shell
- **分散式爬蟲**：Scrapling Spider / Scrapy + Scrapling

## 參考來源

- [Scrapling GitHub](https://github.com/D4Vinci/Scrapling) — 主要參考框架
- [Scrapling MCP Server 文件](https://github.com/D4Vinci/Scrapling/tree/main/docs/ai) — MCP 整合指南
- [Scrapling Parsing 文件](https://github.com/D4Vinci/Scrapling/tree/main/docs/parsing) — 解析與選擇器文件
- [Playwright 文件](https://playwright.dev/docs/intro) — 瀏覽器自動化
- [HTTPX 文件](https://www.python-httpx.org/) — 非同步 HTTP client
