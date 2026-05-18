# Anti-bot Bypass Techniques

## 問題

現代網站使用多種 anti-bot 技術來防止自動化爬取：

- **Cloudflare**：JS challenge、CAPTCHA、瀏覽器指紋
- **Akamai**：行為分析、請求指紋
- **reCAPTCHA**：使用者行為驗證
- **Rate Limiting**：IP 為單位的請求頻率限制
- **Browser Fingerprinting**：WebGL、Canvas、Font、Audio 指紋
- **TLS Fingerprinting**：JA3/JA3S 指紋
- **HTTP Headers 驗證**：順序、大小寫、自訂 header

## Scrapling 的解決方案

Scrapling 提供了 **StealthyFetcher**，整合了多種 anti-detection 技術：

### 1. 瀏覽器指紋偽裝

```python
# StealthyFetcher 自動處理：
# - 真實的 User-Agent（模擬真實瀏覽器）
# - 正確的 Accept-Language / Accept-Encoding
# - WebGL、Canvas、Font 指紋
# - 瀏覽器外觀（navigator、screen 等）
# - TLS handshake 指紋

async with scrapling.StealthyFetcher() as fetcher:
    page = await fetcher.stealthy_fetch('https://protected.example.com')
```

### 2. 人類行為模擬

```python
# 模擬真實使用者行為：
# - 隨機的滑鼠移動軌跡
# - 自然的滾動模式
# - 隨機的按鍵間隔
# - 頁面載入等待時間

# 可自訂行為參數
page = await fetcher.stealthy_fetch(
    'https://protected.example.com',
    wait_until='networkidle',  # 等待網路閒置
    humanize=True,             # 啟用人類行為模擬
    random_delay=(1, 3)       # 隨機延遲 1-3 秒
)
```

### 3. Session 維持

```python
# 跨請求維持相同的瀏覽器上下文
session = await fetcher.open_session()
page1 = await session.stealthy_fetch('https://protected.example.com/login')
# 登入後，session 自動維持 cookies 和 token
page2 = await session.stealthy_fetch('https://protected.example.com/dashboard')
```

## 通用 Anti-bot Bypass 策略

### 1. 代理輪換（Proxy Rotation）

```python
# 使用代理池輪換 IP
proxies = [
    'http://proxy1.example.com:8080',
    'http://proxy2.example.com:8080',
    'http://proxy3.example.com:8080',
]

async with scrapling.StealthyFetcher() as fetcher:
    for url in urls:
        proxy = random.choice(proxies)
        page = await fetcher.stealthy_fetch(url, proxy=proxy)
```

### 2. 請求間隔控制

```python
# 隨機化請求間隔，避免被檢測為 bot
import random
import asyncio

async def scrape_with_delay(urls, fetcher):
    for url in urls:
        delay = random.uniform(2, 5)  # 2-5 秒隨機間隔
        await asyncio.sleep(delay)
        page = await fetcher.stealthy_fetch(url)
```

### 3. 請求頭偽裝

```python
# 自訂請求頭以模擬真實瀏覽器
headers = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ...',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'Accept-Language': 'en-US,en;q=0.5',
    'Accept-Encoding': 'gzip, deflate, br',
    'Connection': 'keep-alive',
    'Upgrade-Insecure-Requests': '1',
    'Sec-Fetch-Dest': 'document',
    'Sec-Fetch-Mode': 'navigate',
    'Sec-Fetch-Site': 'none',
    'Sec-Fetch-User': '?1',
}
```

### 4. CAPTCHA 處理

```python
# CAPTCHA 出現時的處理策略：
# 1. 降低請求頻率
# 2. 更換代理 IP
# 3. 使用 CAPTCHA 解決服務（如 2Captcha、Anti-Captcha）
# 4. 等待一段時間後重試

async def handle_captcha(fetcher, url, max_retries=3):
    for attempt in range(max_retries):
        page = await fetcher.stealthy_fetch(url)
        if not page.contains_captcha():
            return page
        # CAPTCHA 出現，更換代理並等待
        await change_proxy()
        await asyncio.sleep(60 * (attempt + 1))  # 遞增等待
    raise Exception("CAPTCHA 無法繞過")
```

## Anti-bot 等級與對應策略

| 等級 | 特徵 | 基本策略 | 進階策略 |
|------|------|---------|---------|
| **L0 - 無保護** | 直接回傳 HTML | HTTP Direct | 不需要 |
| **L1 - 基本檢查** | User-Agent、Referer 檢查 | 自訂 headers | 不需要 |
| **L2 - Rate Limiting** | 短時間內多次請求被限 | 請求間隔 + retry | 代理輪換 |
| **L3 - JS Challenge** | Cloudflare 5 秒盾 | StealthyFetcher | 瀏覽器指紋偽裝 |
| **L4 - CAPTCHA** | reCAPTCHA、hCaptcha | StealthyFetcher + 代理輪換 | CAPTCHA 解決服務 |
| **L5 - 行為分析** | 滑鼠軌跡、滾動模式 | 人類行為模擬 | 真實瀏覽器錄製回放 |
| **L6 - 指紋辨識** | TLS、WebGL、Canvas 指紋 | StealthyFetcher | 專用指紋偽裝工具 |

## 注意事項

- **法律合規**：遵守網站的 `robots.txt` 和 Terms of Service
- **道德使用**：不要對網站造成過度負載
- **IP 管理**：使用代理池時注意 IP 品質和地理分佈
- **監控與警報**：監控成功率、錯誤率，及時發現被 ban 的情況
- **優雅降級**：當 bypass 失敗時，有備用方案（如使用官方 API）

## 與既有知識的關係

- [`analysis/web/README.md`](../../analysis/web/README.md) — Web Scraping 分析方法的範圍
- [`analysis/web/sources-and-tools.md`](../../analysis/web/sources-and-tools.md) — 工具選擇參考（StealthyFetcher）
- [`analysis/web/multi-strategy-routing.md`](../../analysis/web/multi-strategy-routing.md) — 提取策略中的 Stealth Browser
- [`intelligence/engineering/cache-replay-pattern.md`](../engineering/cache-replay-pattern.md) — 開發階段使用 Cache 避免被 ban
- [`analysis/web/mcp-server-patterns.md`](../../analysis/web/mcp-server-patterns.md) — MCP Server 的 Stealth Layer
