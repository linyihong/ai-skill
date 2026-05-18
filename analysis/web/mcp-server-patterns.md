# MCP Server Design Pattern for Web Scraping

## 問題

AI agent 需要從網頁提取資料時，傳統做法是：

1. 讓 agent 自己寫爬蟲程式碼 → 耗時、容易出錯
2. 提供預先爬取的資料 → 靜態、無法應對動態需求
3. 使用通用 HTTP tool → 無法處理 JS 渲染、anti-bot

MCP（Model Context Protocol）Server 提供了一個更好的方式：**將網頁提取能力封裝為 agent 可直接呼叫的工具**。

## Scrapling MCP Server 的 3-Layer 架構

Scrapling 實作了 10 個 MCP 工具，分為 3 層：

```
┌─────────────────────────────────────────────┐
│              Session Management              │
│  open_session / close_session / list_sessions│
├─────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Basic   │  │ Dynamic  │  │ Stealth  │   │
│  │  Layer   │  │  Layer   │  │  Layer   │   │
│  ├──────────┤  ├──────────┤  ├──────────┤   │
│  │ get      │  │ fetch    │  │ stealthy │   │
│  │ bulk_get │  │ bulk_fetch│  │ _fetch   │   │
│  └──────────┘  └──────────┘  │ bulk_    │   │
│                               │ stealthy │   │
│                               │ _fetch   │   │
│                               └──────────┘   │
├─────────────────────────────────────────────┤
│              Screenshot                      │
│  screenshot (returns ImageContent)           │
└─────────────────────────────────────────────┘
```

### Layer 1: Basic HTTP（`get`, `bulk_get`）

- 使用標準 HTTP client
- 適用於靜態 HTML 頁面
- 最快、最輕量
- `bulk_get` 支援並發請求

### Layer 2: Dynamic Browser（`fetch`, `bulk_fetch`）

- 使用 headless browser（Playwright）
- 適用於需要 JavaScript 渲染的頁面
- 支援等待特定元素、網路閒置
- `bulk_fetch` 支援多頁面並發

### Layer 3: Stealth Browser（`stealthy_fetch`, `bulk_stealthy_fetch`）

- 使用 anti-detection 技術
- 適用於有 anti-bot 保護的網站（Cloudflare、Akamai 等）
- 模擬真實瀏覽器行為（指紋、headers、行為模式）
- `bulk_stealthy_fetch` 支援並發

### Session Management（`open_session`, `close_session`, `list_sessions`）

- 跨請求維持 session（cookies、headers、瀏覽器上下文）
- 登入狀態保持
- 多 session 並存（不同帳號、不同配置）

### Screenshot（`screenshot`）

- 返回 `ImageContent`（非 base64）
- 適用於視覺驗證、除錯
- 可與其他工具搭配使用

## 關鍵設計決策

### 1. CSS Selector Precision Extraction

Scrapling MCP Server 允許 agent 在請求時指定 CSS selector，**只提取需要的元素**而非整個頁面：

```
優點：
- 大幅減少 token 消耗（只傳送相關內容給 LLM）
- 提高回應速度
- 降低雜訊干擾
```

### 2. Prompt Injection Protection

MCP Server 自動 sanitize 回應內容：
- 移除 CSS-hidden elements
- 移除 aria-hidden 元素
- 移除零寬字元（zero-width characters）
- 移除 HTML comments
- 移除 template tags

### 3. 工具選擇策略

| 條件 | 推薦工具 |
|------|---------|
| 靜態 HTML，無 anti-bot | `get` |
| 需要 JS 渲染 | `fetch` |
| 有 anti-bot 保護 | `stealthy_fetch` |
| 大量頁面，同類型 | `bulk_get` / `bulk_fetch` / `bulk_stealthy_fetch` |
| 需要登入 | `open_session` → 對應 fetch tool |
| 視覺驗證 | `screenshot` |

## 適用場景

| 場景 | 適用性 | 說明 |
|------|--------|------|
| AI Agent 網頁提取 | ✅ 最佳 | 直接提供工具，無需寫程式碼 |
| 快速原型開發 | ✅ 適合 | 互動式 shell + MCP tools |
| CI/CD Pipeline | ✅ 適合 | Docker 部署，CLI 模式 |
| 大量批次爬取 | ⚠️ 可搭配 | MCP 適合 agent 驅動，大量爬取建議用 Spider |
| 即時資料監控 | ✅ 適合 | Session 管理 + 定期請求 |

## 與既有知識的關係

- [`analysis/web/README.md`](../README.md) — Web Scraping 分析方法的範圍
- [`analysis/web/sources-and-tools.md`](sources-and-tools.md) — 工具選擇參考
- [`intelligence/web-scraping/adaptive-parsing.md`](../../intelligence/web-scraping/adaptive-parsing.md) — Adaptive Parsing 可與 MCP Server 搭配
- [`enforcement/sanitization.md`](../../enforcement/sanitization.md) — Prompt Injection Protection 規則

## 參考

- [Scrapling MCP Server Guide](https://github.com/D4Vinci/Scrapling/tree/main/docs/ai/mcp-server.md)
- [MCP Specification](https://modelcontextprotocol.io/)
