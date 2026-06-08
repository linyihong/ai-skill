# 去敏、占位符與 Prompt Injection Protection

## 去敏與占位符

可重用文件（含 `feedback_history` lesson）**不要**包含：

- 完整 `Authorization` / token、session cookie、可識別特定使用者的實體裝置識別。
- AES/HMAC/簽名密鑰（**除非**教學用、合成、可公開測試向量）。
- 未去敏的 raw request/response。
- 本機**真實**絕對路徑、使用者帳號名稱、私用工作目錄、git clone 實體路徑。
- Project incident 的具體 app / project 名稱、endpoint、host、payload fragment、sample ID、class/test 名稱、live run 結果或環境 quirks；這些依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) 留在專案文件。

一律改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等占位符（見本庫 [README.md](../README.md)）。

Mechanical pre-commit enforcement is documented in [`sanitization-mechanical.md`](sanitization-mechanical.md). The mechanical surfaces are [`../runtime/repository-topology.yaml`](../runtime/repository-topology.yaml) and [`../runtime/sanitization-patterns.yaml`](../runtime/sanitization-patterns.yaml).

## Prompt Injection Protection（網頁內容去敏）

當從網頁提取內容傳遞給 AI 模型時，**必須** sanitize 以下類型的內容以防止 prompt injection：

### 必須移除的內容類型

| 類型 | 說明 | 範例 |
|------|------|------|
| **CSS-hidden elements** | 使用 `display:none`、`visibility:hidden`、`hidden` attribute 隱藏的元素 | 隱藏的廣告、追蹤像素 |
| **aria-hidden elements** | 標記為 `aria-hidden="true"` 的元素 | 無障礙輔助元素 |
| **Zero-width characters** | Unicode 零寬字元（`\u200B`、`\u200C`、`\u200D`、`\uFEFF` 等） | 隱藏的文字注入 |
| **HTML comments** | `<!-- ... -->` 中的內容 | 開發者註解可能包含隱藏指令 |
| **Template tags** | `{{ }}`、`{% %}`、`<%= %>` 等模板語法 | 可能被誤解為指令 |
| **Script tags** | `<script>` 中的 JavaScript 程式碼 | 可能包含惡意程式碼 |
| **Event handlers** | `onclick`、`onload`、`onerror` 等 inline event handler | 可能被用於 XSS |
| **Meta refresh/redirect** | `<meta http-equiv="refresh">` 中的 URL | 可能導向惡意網站 |
| **Base64-encoded content** | 隱藏在 data URI 或 base64 字串中的內容 | 編碼後的隱藏指令 |

### 實作建議

```python
# Scrapling MCP Server 的自動 sanitization 範例
def sanitize_web_content(html: str) -> str:
    # 1. 移除 CSS-hidden elements
    # 2. 移除 aria-hidden elements
    # 3. 移除 zero-width characters
    # 4. 移除 HTML comments
    # 5. 移除 template tags
    # 6. 保留結構化內容（表格、列表、標題）
    return sanitized_html
```

### 適用場景

- **MCP Server 回應**：Scrapling MCP Server 已內建此保護
- **Web Scraping 結果**：使用 Scrapling 提取的內容自動 sanitize
- **Browser automation 輸出**：Playwright/Puppeteer 提取的內容需手動 sanitize
- **HTTP client 回應**：requests/httpx 提取的 raw HTML 需手動 sanitize

### 不適用場景

- 內部產生的文件（非網頁來源）
- 已確認安全的受信任來源
- 不需要傳遞給 AI 模型的內容（如直接寫入資料庫）

← [回到共用規則索引](README.md)
