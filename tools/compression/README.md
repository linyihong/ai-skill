# Tool Output Compression

`tools/compression/` 負責壓縮 tool output，避免 raw output 直接塞入 context 造成 token 浪費。

## 壓縮策略

| Output 類型 | 預設壓縮 | 方法 |
| --- | --- | --- |
| Stack trace | summary | 只保留關鍵錯誤訊息與 stack top 5 frames |
| JSON | structured | 只保留相關欄位，移除冗餘結構 |
| Git diff | summary | 只保留變更摘要與關鍵行 |
| Log | summary | 只保留 error/warning 層級 |
| Search results | ranked | 只保留 top 5 相關結果 |
| File content | summary-first | 先讀 summary，需要才展開 |

## 壓縮層級

```yaml
compression_levels:
  raw:
    description: 不壓縮，完整輸出
    use_case: debug、需要精確內容時
    token_ratio: 1.0

  summary:
    description: 摘要輸出，保留關鍵資訊
    use_case: 預設、快速瀏覽
    token_ratio: 0.2-0.3

  structured:
    description: 結構化輸出，只保留相關欄位
    use_case: JSON、YAML 等結構化資料
    token_ratio: 0.1-0.2

  minimal:
    description: 極簡輸出，只保留結論
    use_case: 確認狀態、快速檢查
    token_ratio: 0.05-0.1
```

## 壓縮規則

1. **Stack trace**：保留 exception type + message + top 5 frames。移除 `... N more` 與無關內部框架。
2. **JSON**：移除 null/empty 欄位。只保留 task-relevant 的 key。
3. **Git diff**：只保留變更檔案清單 + 每個檔案的前 3 行變更。完整 diff 在需要時才展開。
4. **Log**：只保留 ERROR、WARN 層級。INFO/DEBUG 在需要時才展開。
5. **Search results**：依 relevance score 排序，只保留 top 5。
6. **File content**：先讀 summary（300-500 tokens），需要才讀全文。

## 與既有層的關係

- `models/compression/README.md`：model-aware compression strategy
- `tools/metadata/`：每個 tool 的 compression 支援標記
- `tools/routing/`：tool activation 時決定 compression level
