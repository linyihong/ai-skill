# Dual-Token Security Audit（雙簽章 / 雙加密 token 並存的審計方法）

**Status**: `candidate-analysis`

## 目的

系統內同時存在兩套以上 token 機制（典型如 **JWT + JWE**、**HMAC 簽章 + 對稱加密**、**平台自簽 token + 廠商回調 token**）時，每一道接縫都是**簽章誤用、replay、key 混用、降級攻擊**的高風險點。本方法提供可重用的觀察路線，不是修補步驟。

## 何時觸發此分析

- 代碼庫同時存在 `jwt`、`jwe`、`jws`、`signing`、`encrypt` 等不同 token 工具類。
- 系統整合多家外部廠商（支付、博弈、IM、SSO），各家簽章規範不同。
- 看到 token validation 邏輯散落在多個位置（middleware / controller / service 各有一份）。
- Audit log 顯示同一使用者出現多種 token format。
- 廠商回調與內部 API 共用同一個 token 解析入口。

## 觀察點（依優先順序）

| 層 | 觀察什麼 | 工具 |
|----|------|------|
| **1. 代碼結構** | Token 相關 util、filter、interceptor 的位置與數量；是否有單一 entry 或多入口 | `grep -r "jwt\|jwe\|sign\|verify"`、依賴圖 |
| **2. Key 管理** | 簽章 key、加密 key、廠商 secret 的儲存位置；是否與 token 類別對應；rotation 機制 | 配置檔、Vault/KMS、env var |
| **3. Algorithm 宣告** | JWT 是否接受 `alg: none`；是否硬編碼 algorithm 防止 confusion attack；JWE 的 `enc` 與 `alg` 是否同樣硬編碼 | 看 verifier 代碼 |
| **4. Token 流向** | 哪些 endpoint 接受哪一種 token；廠商回調 vs 使用者 session vs 內部 RPC 是否走同一驗證 | API 文件 + middleware 順序 |
| **5. Validation 一致性** | 同一種 token 在不同地方驗證時，是否有的檢查 `exp` 有的不檢查、有的檢查 `aud` 有的不檢查 | 比對所有 verifier |
| **6. Replay 防護** | 是否有 `jti` / nonce / timestamp window；廠商回調是否有 idempotency key | 看 callback handler |
| **7. Log 外洩** | Log 中是否完整輸出 token / signature / decrypted payload | log config、live log 抽樣 |
| **8. Error 訊息** | Token 驗證失敗的錯誤訊息是否揭露原因（"invalid signature" vs "expired" vs "malformed"）導致 oracle attack | error response 內容 |

## Audit Steps

```text
1. 列出 token universe
   找出所有 token 類別（JWT / JWE / HMAC / opaque / vendor-specific）
   每一種列出：發行者、消費者、用途、lifetime

2. 畫 token flow 圖
   每一道 endpoint 標註：accepts which token type, validates with which key, falls back to ?

3. Key/algorithm matrix
   每一種 token × 每一個 key → 是否有共用、混用、互通的可能

4. 接縫盤點
   雙 token 系統最危險的是「接縫」：
   - 同一個 filter chain 處理多種 token 時的判別邏輯
   - 廠商回調進入內部 API 時的轉換點
   - 內部服務間呼叫的二次包裝
   每個接縫獨立 review

5. Failure mode 對照
   對每個接縫，問：
   - 攻擊者能否用 token A 假裝 token B？
   - 簽章降級（alg confusion）是否被阻擋？
   - Key reuse 是否存在（同一個 key 同時用於 sign 與 encrypt）？
   - Replay window 是多久？跨 token 種類是否累積防護？
```

## Failure Signals（紅旗）

| Signal | 可能問題 |
|--------|------|
| Token verifier 接受 `alg: none` 或從 token header 讀 alg 不做 whitelist | Algorithm confusion |
| JWT 的 secret 與 JWE 的 key 是同一字串 | Key reuse 漏洞 |
| 廠商回調 endpoint 與使用者 session endpoint 共用同一個 token middleware | Cross-domain token 互通 |
| Log 出現 `Authorization: Bearer eyJ...` 完整字串 | Token leak |
| 驗證失敗訊息透露具體欄位（`invalid signature` vs `unknown key id`） | Oracle attack |
| 沒有 `jti` 或 replay window 的 callback endpoint | Replay attack |
| Refresh token 與 access token 用相同 audience/claims 結構 | Token confusion |
| Token validation 在多處實作（filter / interceptor / controller manual） | 不一致風險，必有一處漏 |

## 證據蒐集規範

審計過程中產出的證據（key fragments、token samples、internal endpoints、vendor secrets）**絕對不可寫進本檔或任何 reusable doc**。依 [`enforcement/sanitization.md`](../../enforcement/sanitization.md)：

- Token 範例僅用 RFC 教科書範例或自製 dummy。
- 廠商名稱抽象為 `Vendor A / B / C`。
- 內部 endpoint 抽象為 `/api/<resource>/<action>`。
- 真實 finding 留在專案 incident 資料夾，不進本庫。

## 驗證

完成 dual-token audit 後應能回答：

| 問題 | 期望答案形式 |
|------|------|
| 系統有幾種 token？ | 列表 + 用途 + 發行/消費端 |
| 每種 token 的 algorithm 與 key 是否獨立？ | Matrix |
| 接縫有幾個？每個接縫的 risk 評估？ | 表格（接縫 / 風險 / 證據 / 建議） |
| 哪些 verifier 不一致？ | 具體差異清單 |
| Log / error 是否會洩漏 token 內容？ | yes/no + 證據位置 |
| Replay 防護是否覆蓋所有 token 種類？ | 覆蓋率 + gap |

若任一問題無法回答，audit 不完整。

## 與其他智慧的關係

- 安全架構決策（如「為什麼這個系統需要 JWE 而不只是 JWT」）→ [`intelligence/engineering/architecture/system-boundaries/`](../../intelligence/engineering/architecture/system-boundaries/README.md)
- 多廠商整合導致的多 token 並存 → [`intelligence/engineering/architecture/vendor-integration-architecture.md`](../../intelligence/engineering/architecture/vendor-integration-architecture.md)
- Secret 不可進 log / commit → [`enforcement/sanitization.md`](../../enforcement/sanitization.md)

---

← [回到 analysis/security/](README.md)
