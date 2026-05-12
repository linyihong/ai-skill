# APK Analysis Execution Flow

本文件定義授權 APK 分析的 agent 執行流程。分析方法見 [`analysis/apk/`](../../analysis/apk/)；工程智慧見 [`intelligence/engineering/apk-analysis/`](../../intelligence/engineering/apk-analysis/)。

## 1. 開始前確認

記錄：

- APK 名稱與版本。
- package name。
- CPU 架構：`arm64-v8a`、`armeabi-v7a` 等。
- 裝置或 emulator。
- 是否 root。
- 是否允許 Frida / MITM / pcap / decompile。
- 是否有測試帳號與授權範圍。

禁止：

- 保存完整 token、cookie、真實 device id、私密 key。
- 把第三方或非目標流量當成分析資料。
- 把單一 App 的私有 host / secret 寫進 reusable skill。

### Reset baseline / 起始狀態

當目標是「從 App 開始到某個具名功能的完整 API 流程」時，開始 capture 前先決定並記錄起始狀態：

| Reset level | 用途 | 注意事項 |
| --- | --- | --- |
| `force-stop only` | 保留帳號/session，只重新啟動 App 與網路 client。 | 適合避免 login/rate limit，同時觀察冷啟動與導航流程。 |
| `clear cache` | 減少資源/cache 對列表或詳情的干擾。 | 不一定清掉 DB/session；需記錄是否仍有本機資料。 |
| `clear app data` | 還原 first-run / session recovery / onboarding 狀態。 | 可能移除測試 session、觸發登入或限流；需使用授權測試帳號並記錄邊界。 |
| `reinstall` | 驗證安裝後首輪 bootstrap / permission / migration。 | 成本最高；不要在不需要 first-run 行為時使用。 |

每個 reset-to-feature capture 應拆成可反查的 window：

1. Reset / state preparation：`force-stop`、可選 data/cache clear、權限、proxy、Frida/MITM/pcap 狀態。
2. Cold start / bootstrap：launch、公告/ onboarding / login / session recovery、startup/background APIs。
3. Navigation：從 launcher/home 到目標 page/tab/module 的 tap/swipe/input。
4. Feature operations：列表、分類/filter、搜尋、分頁/scroll、詳情、評論/媒體/action 等。
5. Documentation closure：page map、operation map、API list、schema/correlation、feature handoff、unknowns。

## 2. Quick Start（預設執行順序）

1. 確認 scope 與 authorization。
2. 分離 method 與 target facts：
   - 可重用技術放 skill 目錄。
   - target-specific endpoints、hosts、schemas、findings 放專案 API/reference docs。
3. 從 network path triage 開始（見 [`analysis/apk/traffic-triage.md`](../../analysis/apk/traffic-triage.md)）。
4. 證據指向特定技術類別後，才讀對應的 `analysis/apk/` 方法。
5. 建立 UI architecture map（若可操作裝置）。
6. 優先使用高語意 hook（request object > raw socket, response decoder > TLS bytes）。
7. 將動態結果轉換為 durable assets：
   - UI architecture map + operation-to-API matrix。
   - Redacted HTTP/API docs。
   - Domain/runtime baseline。
   - Feature reconstruction handoff。
   - Redacted request/response samples。
   - Offline decoders or fixtures。
   - API/schema docs。
   - Contract tests（若專案有 SDK/client implementation）。
8. **Automatic skill feedback**：每次學到新 reusable technique 時，在同一輪寫入 `feedback_history/`。

## 3. 分析結束定義

一次分析可以收斂時，應具備：

- 清楚知道核心流量走哪個 stack。
- 若使用者要求完整 app-start-to-feature 流程，已記錄 reset/cache/session baseline，並把 windows 分開回填。
- 有 request metadata 或已證明拿不到的原因。
- 有 response outer shape。
- 若有加密，有解碼點或下一步定位計畫。
- 有去敏樣本或 fixture。
- 有文件回填位置。
- **若下游要程式化取數／接 SDK／寫 integration：** 專案內已有或可指向的 Domain/runtime baseline。缺則視為收口不完整。
- **若要開始開發 live-facing SDK/client/app tool：** baseline 必須通過 development readiness gate。
- 有新的 reusable lesson，或使用者/reviewer 提出可泛化改進時，已在 `feedback_history/` 新增對應檔案。

## 4. Required Output Style

回報分析進度時，包含：

- Goal, action, and validation or reference source。
- What was tested。
- What evidence was observed。
- What was ruled out。
- What remains unknown。
- The next lowest-risk action。

記錄新發現時，包含：

- Trigger or UI path。
- Tool and command summary。
- Evidence file path or sanitized excerpt。
- Feature/capability mapping and operation id。
- Page-level UI map path（若任務針對具名 page/tab/module）。
- Domain/runtime baseline 更新點。
- Generalized lesson。
- Follow-up validation。

## 5. Safety and Sanitization

Never write raw secrets into reusable skill docs：

- Full Authorization tokens。
- Session cookies。
- Device identifiers that belong to a real user/device。
- AES/HMAC secrets unless they are synthetic examples。
- Private hostnames not meant for the reusable method guide。
- Personal user data。

Use placeholders：

```text
<package-name>
<device-serial>
<api-host>
<proxy-host>:<proxy-port>
<token-redacted>
<secret-redacted>
```

## 6. Feedback Loop

若分析發現新的 reusable idea，或使用者/reviewer 建議可泛化的改進：

1. 在 `feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md` 建立 dated lesson（同輪，除非缺證據）。
2. 泛化，使其不限於單一 APK。
3. 加入 evidence 與 validation criteria。
4. 驗證後 promotion 到 `WORKFLOW.md`、`TOOLS.md`、`DOCUMENTATION.md` 或 `techniques/<category>/`。

**Root-cause check when feedback did not trigger：** 判斷 trigger 是否太隱晦、idea 被誤判為 project-only、validation uncertainty 阻擋了 candidate lesson、或 writeback transaction 未開啟。強化對應 trigger/checklist 後完成 sync/commit/push。

**Agent checklist before ending an APK-analysis task：** 是否有新的 generalized lesson？若有 → `feedback_history/<category>/` 或 `feedback_history/common/` 有新檔（minimum）；optional promotion 到主文件。

## 7. 文件分層

| 層級 | 內容 | 存放位置 |
| --- | --- | --- |
| 方法論 | 流量路徑判斷、工具選擇、hook 策略、去敏規則。 | `analysis/apk/` |
| 專案結論 | 目標 App 的 API、host、endpoint、schema、媒體規則。 | 專案 API/reference docs |
| 原始證據 | pcap、MITM export、Frida log、raw response、decrypted fixture。 | 專案受控位置（gitignored） |

## 8. 回填規則

每次分析完成後：

- UI Behavior 回填專案 UI 行為入口或 page-level map。
- 目標 API 結論回填專案 API 文件。
- 解碼規則回填協議/解密文件。
- SDK/client 行為回填 BDD/tests。
- 若分析文件要用於 app 工具/SDK/client/mock/contract test，同輪啟用 `app-development-guidance` 並交出 Feature Reconstruction Handoff。
- 通用技巧回填 `feedback_history/<category>/` 或 `feedback_history/common/`。
- App 開發 guidance 回填 `app-development-guidance/`。

---

← [回到 workflow/apk-analysis/](README.md)
