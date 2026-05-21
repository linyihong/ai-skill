> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-21 — Dual-Token Security Audit（雙簽章/雙加密 token 並存的審計方法）

Status: candidate

#### One-line Summary

系統同時存在兩套以上 token 機制（JWT + JWE、HMAC + 對稱加密、平台 token + 廠商回調 token）時，每一道接縫都是簽章誤用、replay、key 混用、降級攻擊的高風險點；本 lesson 沉澱可重用審計路線。

#### Human Explanation

雙 token 系統的危險不在單一 token 本身，而在「接縫」：同一個 filter chain 處理多種 token 時的判別邏輯、廠商回調進入內部 API 時的轉換點、內部服務間呼叫的二次包裝。每個接縫都可能因為 alg confusion、key reuse、replay window 不足、log 洩漏、error oracle 而被利用。需要系統化方法蒐證並評估，不能只看單點代碼。

#### Trigger

- 代碼庫同時存在 `jwt` / `jwe` / `jws` / `signing` / `encrypt` 等不同 token 工具類
- 系統整合多家外部廠商，各家簽章規範不同
- Token validation 邏輯散落多處（middleware / controller / service 各有一份）
- Audit log 顯示同一使用者出現多種 token format
- 廠商回調與內部 API 共用同一個 token 解析入口

#### Evidence

- Tool: 後端 repo 結構與依賴觀察
- Sanitized excerpt: 觀察到一個多廠商後端 repo 共用工具層內同時存在 `util/jwt` 與 `util/jwe` 兩套子目錄，廠商整合層各自有自己的簽章流程；尚未實際進入代碼審計，僅憑結構訊號識別風險
- Evidence path: 證據留在 `<PROJECT_ROOT>/` 原始 review note，不複製到本庫

#### Generalized Lesson

審計五步驟（依優先順序）：

1. 列出 token universe（種類、發行者、消費者、用途、lifetime）
2. 畫 token flow（每個 endpoint accepts which token, validates with which key）
3. 建立 key/algorithm matrix（檢查共用、混用、互通可能）
4. 接縫盤點（共用 filter / 廠商回調入內部 / 服務間二次包裝）
5. Failure mode 對照（alg confusion、key reuse、replay window、log leak、error oracle）

八個觀察點：代碼結構、key 管理、algorithm 宣告、token 流向、validation 一致性、replay 防護、log 外洩、error 訊息。

#### Agent Action

在後端 repo 看到雙 token 訊號時：

1. `grep -r "jwt\|jwe\|sign\|verify"` 列出 token 相關 util、filter、interceptor
2. 對每一個 token 種類執行五步驟審計
3. 每個接縫獨立 review，不要假設「同 filter chain 必然安全」
4. 若發現 red flag（alg: none、key reuse、log 洩漏完整 token）立即提報

#### Goal / Action / Validation

- Goal: 系統化辨識 token 接縫風險
- Action: 五步驟審計 + 八觀察點檢查；產出 token universe 表、key/alg matrix、接縫風險表
- Validation or reference source: 能回答「token 種類、algorithm 與 key 是否獨立、接縫數量與風險、verifier 不一致清單、log/error 洩漏狀態、replay 防護覆蓋率與 gap」

#### Applies When

- 後端系統含使用者驗證 + 外部廠商回調
- 系統涉及金流、合規、博弈、支付等高敏感業務

#### Does Not Apply When

- 只有單一 token 機制且發行/消費端單一
- 純內部服務無外部回調

#### Validation

- 對既有後端 repo 跑五步驟能產出可審查的 token universe + 接縫風險表
- 已知 alg confusion / key reuse / replay 等 CVE 模式可在 audit 中被識別

#### Promotion Target

- ✅ `analysis/security/dual-token-audit.md`（已於 `d5ec684` 寫入，並建立新分類 `analysis/security/`）
- ✅ `knowledge/summaries/dual-token-audit.md`（已於 `d5ec684` 寫入）

#### Required Linked Updates

- ✅ `analysis/security/README.md` 新建（已於 `d5ec684` 寫入）
- ✅ `analysis/README.md` 加入 `security/` 入口（已於 `d5ec684` 更新）
- ✅ `knowledge/summaries/README.md` 索引（已於 `d5ec684` 更新）
- Step 6（Intelligence Extraction）不適用：analysis atom 已直接寫入
- Step 7（Failure Learning）不適用於本 lesson 內容；本輪流程失誤已由既有 `enforcement/failure-patterns/knowledge-update-flow-bypassed-by-sub-pipeline.md` 覆蓋
