# Cache & Replay Pattern for Web Scraping

## 問題

在開發和迭代網頁提取邏輯時，每次執行都實際發送 HTTP 請求會導致：

- **開發速度慢**：等待網路回應（尤其是有 anti-bot 保護的網站）
- **被 ban 風險**：頻繁請求觸發 rate limiting 或 IP blocking
- **成本增加**：消耗 proxy 頻寬、API 配額
- **除錯困難**：無法重現相同的回應來測試解析邏輯

## Scrapling 的解決方案：Development Mode

Scrapling 提供了 **Development Mode**，在第一次執行時快取回應，後續執行直接使用快取：

```
第一次執行（Cache Miss）
  ┌─────────┐     ┌──────────┐     ┌─────────┐
  │ Request │────▶│ Fetch    │────▶│ Cache   │
  │         │     │ (Network)│     │ Store   │
  └─────────┘     └──────────┘     └─────────┘

後續執行（Cache Hit）
  ┌─────────┐     ┌──────────┐     ┌─────────┐
  │ Request │────▶│ Cache    │────▶│ Replay  │
  │         │     │ (Local)  │     │         │
  └─────────┘     └──────────┘     └─────────┘
```

### 使用方式

```python
# 啟用 Development Mode
import scrapling
scrapling.set_development_mode(True)

# 第一次執行：實際發送請求並快取
page = scrapling.get('https://example.com')

# 第二次執行：直接從快取讀取（不會發送實際請求）
page = scrapling.get('https://example.com')  # 使用快取
```

## 通用 Cache & Replay 模式

### 1. Response Cache

快取整個 HTTP 回應，適用於：
- 頁面內容不常變動
- 需要快速迭代解析邏輯
- 多個解析器共用相同來源

```python
class ResponseCache:
    def __init__(self, ttl: int = 3600):
        self.cache = {}
        self.ttl = ttl
    
    def get(self, url: str) -> Optional[Response]:
        if url in self.cache:
            entry = self.cache[url]
            if time.time() - entry['timestamp'] < self.ttl:
                return entry['response']
        return None
    
    def set(self, url: str, response: Response):
        self.cache[url] = {
            'response': response,
            'timestamp': time.time()
        }
```

### 2. Selector Cache

快取解析結果（而非 raw HTML），適用於：
- 解析邏輯複雜、耗時
- 同一頁面需要多次解析不同欄位
- 解析邏輯不變，但頁面可能更新

```python
class SelectorCache:
    def __init__(self):
        self.cache = {}
    
    def get_selector_result(self, page_hash: str, selector: str):
        key = f"{page_hash}:{selector}"
        return self.cache.get(key)
    
    def set_selector_result(self, page_hash: str, selector: str, result):
        key = f"{page_hash}:{selector}"
        self.cache[key] = result
```

### 3. Session Cache

快取 session 狀態（cookies、headers），適用於：
- 需要登入的網站
- session 建立成本高（CAPTCHA、2FA）
- 多個請求共用相同 session

```python
class SessionCache:
    def __init__(self):
        self.sessions = {}
    
    def get_or_create_session(self, session_id: str):
        if session_id not in self.sessions:
            self.sessions[session_id] = create_new_session()
        return self.sessions[session_id]
```

## 快取策略選擇

| 策略 | 粒度 | 適用場景 | 失效處理 |
|------|------|---------|---------|
| **Response Cache** | URL-level | 頁面內容穩定 | TTL 過期後重新請求 |
| **Selector Cache** | Selector-level | 解析邏輯複雜 | 頁面變更時清除 |
| **Session Cache** | Session-level | 登入狀態維持 | Session 過期時重新建立 |
| **Hybrid Cache** | 多層 | 複雜爬蟲 | 分層失效 |

## 開發流程中的 Cache 使用

```
開發階段
  ├─ 1. 第一次執行：實際請求 → 快取回應
  ├─ 2. 撰寫/修改解析邏輯
  ├─ 3. 重新執行：使用快取 → 快速驗證
  ├─ 4. 解析邏輯正確後：清除快取 → 用真實回應驗證
  └─ 5. 部署：關閉 Development Mode

優點：
- 步驟 2-3 迭代速度提升 10-100x
- 避免開發期間被 ban
- 可重現的除錯環境
```

## 注意事項

- **快取一致性**：確保快取內容與實際回應一致（尤其是頁面更新時）
- **快取大小**：設定合理的 TTL 和最大快取數量，避免記憶體爆炸
- **敏感資料**：不要在快取中儲存敏感資訊（token、密碼）
- **開發/生產分離**：Development Mode 僅用於開發，生產環境應關閉

## 與既有知識的關係

- [`analysis/web/README.md`](../../analysis/web/README.md) — Web Scraping 分析方法的範圍
- [`analysis/web/multi-strategy-routing.md`](../../analysis/web/multi-strategy-routing.md) — Cache 作為並發策略的一部分
- [`analysis/web/mcp-server-patterns.md`](../../analysis/web/mcp-server-patterns.md) — MCP Server 可搭配 Cache 使用
