# Web Scraping Analysis

`analysis/web/` 負責 Web Scraping 的分析方法。本層保存網站結構觀察、反爬機制判讀、動態內容分析、MCP 工具設計評估等可重用分析方法。

> ⚠️ **Stub**: 此目錄為初始建立，內容將由 [`plans/active/2026-05-18-scrapling-knowledge-integration-plan.md`](../../plans/active/2026-05-18-scrapling-knowledge-integration-plan.md) 逐步補完。

## 分析範圍

- 網站結構分析（DOM tree、API endpoints、authentication flow）
- 反爬機制判讀（rate limiting、fingerprinting、challenge pages）
- 動態內容分析（JavaScript rendering、WebSocket、XHR）
- MCP 工具設計評估（tool granularity、session management、stealth strategy）

## 工具與資源

- [`sources-and-tools.md`](sources-and-tools.md)：Web Scraping 工具比較與選擇原則

## 與既有層的關係

- `intelligence/web-scraping/` 承接從分析結果萃取出的工程智慧（如 adaptive parsing）
- `intelligence/engineering/mcp-server-patterns.md` 承接 MCP 工具設計模式
