# Adaptive Parsing（自適應解析）

## 問題

傳統的網頁解析依賴固定的 CSS selector 或 XPath，但：

- 網頁結構經常變動（A/B testing、改版、個人化內容）
- 不同頁面可能有不同的 DOM 結構（列表頁 vs 詳情頁）
- 動態內容（JavaScript 渲染）導致 selector 失效
- 維護大量 selector 成本高

## Scrapling 的解決方案：Find Similar Elements

Scrapling 實作了 **Find Similar Elements** 演算法，能夠在 DOM 結構變動時自動找到相似的目標元素。

### 演算法步驟（3 層過濾）

```
Step 1: DOM Tree Depth 過濾
  └─ 找出所有與目標元素相同 DOM 深度的元素
  
Step 2: Tag Name 層級過濾
  └─ 比對 tag name、parent tag name、grandparent tag name
  
Step 3: Fuzzy Attribute Matching
  └─ 使用相似度閾值（default: 0.2）比對屬性
```

#### Step 1：DOM Tree Depth

找出所有與目標元素在 DOM 樹中**相同深度**的元素。這是最高效的過濾條件，因為相同深度的元素通常屬於同一層級的內容。

#### Step 2：Tag Name Hierarchy

進一步過濾，比對：
- 元素本身的 tag name（如 `div`、`span`、`a`）
- 父層級的 tag name
- 祖父層級的 tag name

這確保找到的元素具有相似的 HTML 結構上下文。

#### Step 3：Fuzzy Attribute Matching

對剩餘元素進行模糊屬性比對：
- 使用 `similarity_threshold` 參數（default: `0.2`）
- 可忽略特定屬性（`ignore_attributes`，default: `href`, `src`）
- 可選擇比對文字內容（`match_text: true`）

### API 使用

```python
# 基本用法 — 找到與第一個 .quote 元素相似的所有元素
page.find('.quote').find_similar()

# 自訂相似度閾值
page.find('.quote').find_similar(similarity_threshold=0.5)

# 忽略特定屬性
page.find('.quote').find_similar(ignore_attributes=['href', 'src', 'data-id'])

# 比對文字內容
page.find('.quote').find_similar(match_text=True)
```

### 與其他選擇方法的搭配

```python
# CSS Selector — 精確定位
page.css('.quote .text')

# XPath — 結構化導航
page.xpath('//div[@class="quote"]/span[@class="text"]')

# Text Content — 文字比對
page.find_by_text('Hello')

# Regex — 模式比對
page.find_by_regex(r'\b\w+@\w+\.\w+\b')

# Find Similar — 自適應解析（當上述方法因結構變動失效時）
page.find('.quote').find_similar()
```

## 適用場景

| 場景 | 傳統方法 | Adaptive Parsing |
|------|---------|-----------------|
| 固定結構網站 | ✅ CSS selector | ✅（但不需要） |
| 結構會變動的網站 | ❌ selector 常失效 | ✅ 自動適應 |
| 多頁面爬蟲（不同模板） | ❌ 需要多組 selector | ✅ 單一錨點即可 |
| A/B testing 頁面 | ❌ 無法預測結構 | ✅ 模糊比對 |
| 動態內容（JS 渲染） | ⚠️ 需等待渲染完成 | ✅ 渲染後自適應 |

## 限制

- **相似度閾值敏感**：太低會抓到不相關元素，太高會漏掉
- **依賴初始錨點**：需要一個有效的初始元素作為比對基準
- **效能開銷**：模糊比對比精確 selector 慢（但 Scrapling 仍比 BS4 快 767x）
- **不適用於**：結構完全不同的頁面（需要完全不同的解析策略）

## 與既有知識的關係

- [`analysis/web/README.md`](../../analysis/web/README.md) — Web Scraping 分析方法的範圍
- [`analysis/web/sources-and-tools.md`](../../analysis/web/sources-and-tools.md) — 工具選擇參考
- [`intelligence/engineering/mcp-server-patterns.md`](../engineering/mcp-server-patterns.md) — MCP Server 整合 Adaptive Parsing
