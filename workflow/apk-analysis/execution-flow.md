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

### Capture Window 詳細規則

以下規則定義 capture window 的品質門檻，適用於「從 App 開始到某個具名功能的完整 API 流程」分析。

**Tab / category / filter strip coverage rule：** 當頁面有 top tabs、category chips、search result tabs、carousel-like tabs 或任何看起來可水平滑動的分類列時，要先記錄 first viewport 的可見項、hierarchy 暴露的總數（例如「第 X 个标签，共 N 个」）、左右滑動後新增項，以及每個 reachable tab 的 `captured` / `needs capture` / `no-network-observed` 狀態。只測可見 tab 不能宣稱完整 tab 面 API 覆蓋。

**Post-selection lazy-load rule：** 選中 tab、category、filter、grid label 或 chip 後若短窗口沒有新增 feature API，不要直接判定該 UI action 沒有 API。先在已驗證 target package / feature context 下補一個低風險後續 gesture（例如列表 scroll、refresh 或 bounded wait），並把結果分成 `selection-only`、`post-selection-triggered`、`no-network-after-follow-up`。這類 replay knob 或 trigger pattern 一旦可重用，必須立即走 feedback lesson 檢查，不等專案收尾。

**UI evidence package validation rule：** 每個 screenshot / hierarchy 用於 UI-to-API 對齊前，必須驗證 foreground package / activity 屬於目標 App。若 XML package 變成 launcher、browser、Google/search、settings、permission page 或其他外部 App，該 window 要標 `external` / `invalid for target UI`，automation 應中止或記錄明確轉場；Frida 仍命中目標 PID 只能證明目標進程內事件，不能自動證明是該 UI step 觸發。對重要 feature checkpoint，package 正確後還要驗證目標 feature context（例如穩定 tab label、page title、section heading、selected tab 或 route anchor）；同 package 但跑到充值、活動、WebView 或其它 module 的 window 要標 `wrong in-app screen` / `invalid for target feature`，不可當作該 feature evidence。

**Checkpoint replay runner rule：** 同一 feature/page 需要反覆測 Frida、media、tab sweep 或 reset baseline 時，將已確認路徑固化成 replay script，並為 `launch`、目標 tab、列表、詳情、媒體區等節點提供 `--target` / checkpoint 停點。每個 checkpoint 都應截圖、dump XML、驗證 target package；如果跑歪，先修 selector、fallback coordinate、wait 或 scroll，再把後續 capture 當證據。

**Post-reset window split rule：** `clear app data` / reinstall 後同時需要 session recovery 與 feature API attribution 時，優先拆成「reset + startup/session recovery」與「已驗證導航後 attach feature hooks」兩個 capture window。若 Frida-from-launch 的長窗口導致外部 App、錯頁、公告、更新、WebView 或 timing drift，不要把 feature 操作硬接在同一窗口；先用 package / feature-context guard 證明 session recovery 成功，再從目標 feature checkpoint attach 低負載 hook 取得 feature API 證據。

文件中要把 API 標為 `startup/preload`、`session-recovery`、`navigation`、`feature-triggered`、`cache-hydration` 或 `background/ambiguous`，避免把啟動期或預載 request 誤判成當前點擊觸發。

**Read-only argument override rule：** 若要驗證分頁、排序、語言、filter 這類 read-only 行為，但完整 UI 很難自然觸發下一頁或邊界值，優先選擇高語意函式參數覆寫的短窗口，並保留 App 自己的 session、signing、gateway 與 decrypt path。覆寫腳本必須預設短窗、輸出 schema/hash/key set，不輸出 raw token、raw signature、raw service 或 raw response value；文件要標 `app-owned signing/decrypt preserved`，不能把結果當成 standalone replay parity。

**Redacted sample-targeting classifier rule：** 若 UI 盲抽樣本低收益，但 decrypted response 裡有可判斷樣本可用性的欄位（例如 count、availability、status、type），可新增 disabled-by-default classifier，只輸出 value class 與 item index（例如 `zero/nonzero/missing/other`），再用 UI replay 點擊候選項。不要輸出 raw id、title、body、comment、user、URL、token、完整 count 或其他內容值；文件要標明此 classifier 只是 sample targeting aid，不是資料擷取或 standalone replay parity。

**Articles-first live adapter smoke rule：** 當 APK 分析輸出要接 SDK/private adapter 並驗證真實 read-only 資料時，先選一條核心 read route 做最小 smoke（通常是 list/page 1），只把 base endpoint、該 route binding、opaque/session provider、identity readiness、signing、decrypt/plaintext boundary 設為必填。分類、詳情、留言、媒體、next-page 等 secondary routes 應是 optional follow-up，除非當前目標明確是 full route parity；不要讓 secondary binding 缺失阻塞第一條 live proof。

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
