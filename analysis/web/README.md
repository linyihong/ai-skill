# Web Scraping Analysis

## 範圍與邊界

Web Scraping Analysis 負責**分析網頁內容提取的需求與可行性**，涵蓋：

- **目標網站評估**：判斷網站結構（靜態 HTML / SPA / SSR）、是否需要 JavaScript 渲染、是否有 anti-bot 保護
- **工具選擇分析**：根據需求選擇最適合的提取工具（HTTP client / headless browser / stealth browser）
- **提取策略設計**：決定使用 CSS selector / XPath / 自適應解析 / AI 輔助提取
- **風險評估**：識別 rate limiting、IP blocking、CAPTCHA、legal/ToS 限制
- **資料品質分析**：評估提取結果的完整性、一致性、時效性

### 不屬於此領域的項目

- 實際執行提取的 workflow（屬於 `workflow/` 層）
- 提取技巧的工程知識（屬於 `intelligence/web-scraping/`）
- 特定網站的提取腳本（屬於專案本身，非知識庫）

## 何時進入此分析領域

當任務涉及以下情境時，應進入 `analysis/web/`：

1. **需要從網頁提取資料** — 爬取產品資訊、新聞、社群內容等
2. **網站有 anti-bot 保護** — Cloudflare、Akamai、reCAPTCHA 等
3. **需要動態渲染內容** — SPA（React/Vue/Angular）、無限滾動、AJAX 加載
4. **需要登入或 session 管理** — 需要 cookie、token、session 維持
5. **需要大量並發請求** — 批次爬取、多頁面爬蟲
6. **需要結構化輸出** — 從非結構化 HTML 提取為 JSON/CSV
7. **需要監控或定期檢查** — 網站變更檢測、價格監控

## 與其他 Analysis 領域的關係

| 領域 | 關係 |
|------|------|
| [`analysis/repo/`](../repo/README.md) | 分析 repo 結構與程式碼；`web/` 分析網頁內容，兩者互補 |
| [`analysis/apk/`](../apk/README.md) | APK 分析可能涉及 API endpoint 提取，需要 `web/` 輔助 |
| [`analysis/issue/`](../issue/README.md) | Issue 分析可能涉及外部連結內容提取 |
| [`analysis/production/`](../production/README.md) | 生產環境監控可能涉及網頁狀態檢查 |
| [`analysis/travel/`](../travel/README.md) | 旅遊規劃可能涉及景點/住宿資訊爬取 |

## 分析流程

```
1. 需求釐清 → 2. 網站評估 → 3. 工具選擇 → 4. 策略設計 → 5. 風險評估 → 6. 資料品質規劃
```

### 1. 需求釐清
- 需要提取哪些資料？
- 資料量級？一次性還是定期？
- 輸出格式要求？

### 2. 網站評估
- 網站使用什麼技術？（查看 `X-Powered-By`、`<meta>`、JS bundle 分析）
- 是否需要 JavaScript 渲染？
- 是否有 anti-bot 保護？（Cloudflare、reCAPTCHA、rate limiting）
- 是否有 API endpoint 可直接使用？

### 3. 工具選擇
- **HTTP client**（requests/httpx）：靜態 HTML，無 anti-bot
- **Headless browser**（Playwright/Puppeteer）：需要 JS 渲染
- **Stealth browser**（Scrapling StealthyFetcher）：有 anti-bot 保護
- **MCP Server**（Scrapling MCP）：AI agent 整合場景
- 詳細策略對照：[`multi-strategy-routing.md`](multi-strategy-routing.md) — 提取/解析/Session/並發策略選擇
- MCP 架構參考：[`mcp-server-patterns.md`](mcp-server-patterns.md) — 3-Layer 工具選擇策略

### 4. 策略設計
- CSS selector / XPath 定位
- 自適應解析（Find Similar Elements）
- 分頁處理（pagination / infinite scroll）
- 並發控制（rate limiting、retry、backoff）
- 策略組合範例：[`multi-strategy-routing.md`](multi-strategy-routing.md) — HTTP Direct / Dynamic Browser / Stealth Browser 組合

### 5. 風險評估
- 法律風險（robots.txt、ToS）
- 技術風險（IP blocking、CAPTCHA）
- 資料風險（完整性、時效性）

### 6. 資料品質規劃
- 驗證機制（schema validation、cross-reference）
- 異常處理（missing fields、format changes）
- 監控與警報

## 參考工具與技術

| 工具 | 適用場景 | 優點 | 限制 |
|------|---------|------|------|
| **Scrapling** | 進階爬蟲、anti-bot bypass | 自適應解析、MCP Server、高效能 | Python only |
| **Playwright** | 動態內容、瀏覽器自動化 | 跨瀏覽器、多語言支援 | 資源消耗較高 |
| **httpx / requests** | 簡單靜態頁面 | 輕量、快速 | 無法處理 JS 渲染 |
| **BeautifulSoup / lxml** | HTML 解析 | 成熟穩定 | 非提取工具，需搭配 HTTP client |
| **Selenium** | 傳統瀏覽器自動化 | 廣泛支援 | 速度慢、資源消耗高 |

## 誰會參考這裡（Inbound References）

- `intelligence/web-scraping/` — 提取技巧的工程知識
- `workflow/` — 實際執行提取的流程
- `knowledge/indexes/README.md` — 任務路由索引
- `knowledge/graphs/` — 知識圖譜

## 遷移狀態

- [x] 初始建立（2026-05-18）
- [ ] 與 `intelligence/web-scraping/` 的邊界已驗證
- [ ] 與 `workflow/` 的邊界已驗證
