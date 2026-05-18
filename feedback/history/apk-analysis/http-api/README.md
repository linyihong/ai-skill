# HTTP API Analysis — Feedback Lessons

| 檔名 | Status | 標題 | 一句話摘要 |
|------|--------|------|-----------|
| `2026-05-18_142100-sha256-hash-verify-python-not-shell.md` | validated | SHA256 驗證用 Python 勿用 shell | 避免 subprocess 引號／編碼導致 service hash 驗證錯誤 |
| `2026-05-15_172519-query-string-length-infer-unknown-param.md` | candidate | Query String 長度分析推算未知參數值 | 從 Frida capture 的 query string 總長度和 keys 推算未知參數（如 service name）的長度 |
| `2026-05-15_172600-same-service-hash-multi-api-type-routing.md` | candidate | 同一 Service Hash 對應多個 API | 同一個 service hash 可透過 type 參數路由到不同 API，query string 的 keys 和長度完全不同 |
| `2026-05-15_173000-category-api-returns-numeric-id-not-display-name.md` | candidate | Category API 返回 Numeric ID | LIST API 的 category 參數需要 numeric ID 而非 display name，須從 categories API 取得 ID mapping |
