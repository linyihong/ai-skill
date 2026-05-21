## analysis.dual-token-audit

| 欄位 | 值 |
| --- | --- |
| Atom ID | `analysis.dual-token-audit` |
| Source path | `analysis/security/dual-token-audit.md` |
| Lifecycle | `candidate` |
| Summary | 系統內同時存在兩套以上 token 機制（JWT + JWE、HMAC + 對稱加密、平台 token + 廠商回調 token）時的審計方法。觀察點：代碼結構、key 管理、algorithm 宣告、token 流向、validation 一致性、replay 防護、log 外洩、error 訊息。Audit 五步：列 token universe → 畫 flow → key/alg matrix → 接縫盤點 → failure mode 對照。 |
| When to read | 代碼庫同時存在多套 token 工具類；系統整合多家外部廠商各家簽章不同；token validation 邏輯散落多處；audit log 顯示同一使用者多種 token format；廠商回調與內部 API 共用同一 token 解析入口。 |
| Do not use for | 不可取代具體修補執行流程（屬 workflow）；不可用於記錄具體 incident finding（屬專案 incident 文件）；不可包含真實 token / key / endpoint 證據。 |
| Context cost | ~310 tokens |
| Estimated full cost | ~2400 tokens |
| Validation signal | 能回答 token universe、algorithm/key matrix、接縫風險評估、verifier 不一致清單、log/error 洩漏狀態、replay 防護覆蓋率與 gap。 |
| Last checked | 2026-05-21 |
