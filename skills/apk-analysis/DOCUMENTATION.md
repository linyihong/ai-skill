# APK 分析文件寫法

分析文件的目標是讓人和 AI 都能重現推理，而不是只留下最後答案。

## 文件分層

建議分成三層：

| 層級 | 內容 | 是否放入 skill |
| --- | --- | --- |
| 方法論 | 流量路徑判斷、工具選擇、hook 策略、去敏規則。 | 可以，放在本資料夾。 |
| 專案結論 | 目標 App 的 API、host、endpoint、schema、媒體規則。 | 不放 skill，放專案 API / reference docs。 |
| 原始證據 | pcap、MITM export、Frida log、raw response、decrypted fixture。 | 不放 skill；放 gitignored 或專案指定位置，文件只引用去敏摘要。 |

## Target APK 行為邊界

專案內由 `apk-analysis` 產出的分析文件，只能描述「本次授權分析的目標 APK」已觀察或可重現的行為。這條規則適用於 UI 地圖、API 文件、動態抓包記錄、靜態分析記錄、cache / db schema、HLS / 媒體流程與所有 project-level target docs。

可以寫：

- 目標 APK 的可見 UI、入口路徑、點擊 / 滑動 / 播放 / 評論等 App 操作。
- 目標 APK 實際發出的 request chain、endpoint、query/body/header 名稱、response schema 與欄位語意。
- 來自靜態分析、hook、pcap、MITM、adb、sqlite、cache、fixture 或 replay 的證據。
- 已確認的播放鏈、本機資料行為、錯誤碼與缺口。

不要寫：

- 其他作品、clone、自研 App、外部產品或未分析 APK 的行為 / 架構 / UI / API 設計。
- 要如何實作、重構、封裝 SDK、規劃資料模型、設計產品或安排工程 TODO。
- 沒有目標 APK 證據的 fallback、相似產品慣例或推測流程。

若 APK findings 要轉成開發建議，請放在清楚分離的 handoff / guidance / plan / SDK 文件，並標示為 derived guidance；APK 行為文件只保留證據、觀察、欄位語意與 `needs capture` / `needs replay` / `needs static confirmation` 等缺口狀態。

## 功能重建交接規範

若分析目標是讓後續 agent 能用 [`app-development-guidance`](../app-development-guidance/) 重新做出同等功能，專案分析文件不能只列 endpoint。它必須把 UI 行為、資料模型、API 合約、狀態轉移、錯誤處理與驗證證據串成可交接規格。

若使用者要用分析文件做出 app 相關工具、SDK、client、mock API、fixture-driven implementation、contract test 或重建功能，agent 必須自動讀取並套用 [`app-development-guidance/SKILL.md`](../app-development-guidance/SKILL.md)。`apk-analysis` 只負責提供去敏證據、UI/API attribution、schema、fixture 與 confidence；開發文件、BDD、Domain Model Contract、API / Interface Contract、Error Handling Contract、implementation slices、tests 與 blocker questions 由 `app-development-guidance` 接手。

### Feature handoff finish gate

當某個具名 feature / tab / module 已經被分析到「核心 UI 操作與主要 API flow 可說明」的程度時，必須在同一輪對話補齊或更新 project-level feature handoff / architecture 文件。這是 completion gate，不是可選整理。

觸發條件包含任一項：

- 核心 flows 已從 `Candidate` 升到 `Confirmed`。
- agent 已能回答此功能的 entry path、主要 UI 區塊、API request keys / response schema、狀態與缺口。
- 使用者問「有沒有 API 文件」、「能不能重建」、「架構是什麼」、「繼續分析某功能」且本輪已形成穩定理解。

要求：

- 不要只更新 `schemas.md`、`endpoint-correlation.md`、hook 腳本、chat 總結或 page map；要有一份 feature-level 文件把 UI 行為、domain concepts、API contracts、state/error、data lifecycle、validation evidence、unknowns 串起來。
- 若資料仍不足，也要建立 skeleton，將缺口標成 `needs capture` / `candidate` / `low confidence`，不要等完全確認才寫。
- 若專案已有 `docs/feature-handoff-<feature>.md`、`docs/features/<feature>.md` 或類似文件，更新既有文件；沒有就新建一份。
- 完成回答前，回查該 feature 文件是否已連回 page map、operation map、API summary/schema/correlation 文件。

每個重要功能至少要留下：

| 面向 | 必填內容 | 交給 app-development-guidance 時的用途 |
| --- | --- | --- |
| Feature / Capability | 功能名稱、使用者目標、入口 screen、非目標或未知限制。 | 形成 Product Brief、BDD scenario 與 bounded context。 |
| UI Behavior | screen id、route id、operation id、前置狀態、tap/swipe/input 步驟、可見結果。 | 形成行為規格、consumer flow、E2E 測試。 |
| Domain Concepts | 從 UI 文案、response fields、狀態碼推得的 entity、value object、state、command、event；信心等級。 | 形成 Domain Model Contract；低信心項目標為 open question。 |
| API / Interface Contract | method/path shape、headers、query/body、response wrapper、inner payload、auth/session、pagination、cache、idempotency。 | 形成 API / Interface Contract 與 mock/fixture。 |
| State And Error Handling | loading/empty/error/success 狀態、錯誤碼、重試、登入過期、權限不足、限流、離線或快取行為。 | 形成 Error Handling Contract 與測試案例。 |
| Data Lifecycle | 欄位來源、derived-from、local cache/storage、刷新時機、敏感性、保留/過期行為。 | 形成 architecture/data ownership 與 storage guidance。 |
| Validation Evidence | pcap/MITM/hook/replay/fixture/screenshot/UI hierarchy/automation script 的去敏引用。 | 支援 contract test、integration test 與風險判斷。 |
| Unknowns / Assumptions | 未觸發流程、低信心 mapping、缺少樣本、未驗證 edge case。 | 形成開發前提問清單，不讓 agent 編造缺失。 |

可直接放在專案分析文件的交接段落：

```markdown
## Feature Reconstruction Handoff

| Field | Value |
| --- | --- |
| Feature ID | `<feature-id>` |
| Capability | |
| User goal | |
| Entry screens | `screen.id` / `route.id` |
| Primary operations | `operation.id` |
| Confidence | high / medium / low |
| Open questions | |

### Behavior Scenarios

| Scenario | Given | When | Then | Evidence |
| --- | --- | --- | --- | --- |

### Domain Model Candidates

| Concept | Type | Evidence source | Invariants / state | Confidence | Notes |
| --- | --- | --- | --- | --- | --- |

### API / Interface Contracts

| Operation | Method / Path or Interface | Request shape | Response shape | Auth/session | Errors | Fixture / evidence |
| --- | --- | --- | --- | --- | --- | --- |

### State And Error Handling

| State / Error | Trigger | User-visible behavior | API evidence | Retry / recovery | Logging / redaction |
| --- | --- | --- | --- | --- | --- |

### Data Lifecycle

| Field / Data | Source | Stored locally | Refresh / expiry | Sensitive | Derived from | Notes |
| --- | --- | --- | --- | --- | --- | --- |

### App Development Guidance Bridge

- Suggested bounded contexts:
- Draft BDD scenarios:
- Draft Domain Model Contract:
- Draft Architecture / ownership notes:
- Draft API / Interface Contract:
- Draft Error Handling Contract:
- Required fixtures or tests:
- Open questions for implementation:
- app-development-guidance activated: yes/no
```

規則：

- 不確定的 domain 名稱、欄位意義、狀態轉移要標 `candidate` 或 `low confidence`，不要改寫成確定事實。
- 若同一 endpoint 支援多個功能或多個 screen，要在交接文件保留多個 operation mapping。
- 若只有 API、沒有 UI 來源，仍要寫 `UI path: unknown`、`Trigger confidence: low`，並說明可用哪些操作補證。
- 若只有 UI、沒有 API 明文，仍要寫可見行為與未知 API 合約，不要假造 request/response。
- 解析 UI hierarchy / accessibility XML 時，`text`、`content-desc`、可點擊節點描述可能直接含 raw 標題、搜尋詞、聊天室文字、留言或個資；可重用或公開文件只寫 selector、bounds、通用 block 名、tab label、schema key 與 evidence path。若必須保留原文，只能放在專案受控原始 capture，不得搬進 skill lesson 或 shared docs。
- 交給 `app-development-guidance` 的內容必須是已去敏、可泛化的功能/合約描述；target-specific host、token、帳號資料與 raw response 留在專案受控位置。
- 若 output 目標包含 app 工具、SDK、client、mock、fixture-driven implementation、contract test 或重建功能，不能只輸出 APK 分析文件；必須同輪啟用 `app-development-guidance` 產出或更新開發文件需求與 blocker questions。

## 單次分析筆記模板

```markdown
# [APK / 功能] 分析紀錄

## Scope

- APK:
- Version:
- Package:
- Device / emulator:
- Authorization:
- Goal:

## Environment

- OS:
- adb:
- Frida:
- Proxy tool:
- Static tools:

## Hypotheses

| Hypothesis | Test | Result |
| --- | --- | --- |
| localhost bridge | lo pcap | |
| system proxy / MITM | proxy capture | |
| Java HTTP stack | Java hook | |
| Flutter / native | connect backtrace / AOT strings | |

## Evidence

| Evidence | Path / excerpt | Interpretation |
| --- | --- | --- |
| pcap | `<path>` | |
| hook log | `<path>` | |
| static search | `<path or command>` | |
| screenshot / UI hierarchy | `<path>` | |

## Findings

- Finding 1.
- Finding 2.

## Feature Reconstruction Handoff

- Feature ID:
- Capability:
- User goal:
- Entry screens:
- Primary operations:
- Candidate domain concepts:
- API / interface contracts:
- State and error handling:
- Data lifecycle:
- Fixtures / validation:
- Open questions for app-development-guidance:

## Unknowns

- Unknown 1.

## Next Steps

1. Next validation.
2. Next fixture or test.

## Sanitization

- Tokens redacted:
- Device identifiers redacted:
- User data removed:
```

## UI 架構地圖模板

用 screenshot、UI hierarchy 與可重放操作，把 App 的可見架構寫成地圖。這份文件放專案分析文件，不放 reusable skill，skill 只保留模板與方法。若截圖太多會拖慢裝置或干擾 hook/pcap 時，先做輕量盤點，等核心 API 解完後再補關鍵 UI 綁定。

### Page-level artifact rule

當使用者要求分析某個具名 page、tab、module 或功能入口，且本次工作已建立 UI/API 對照時，必須在專案 docs 產出或更新一份 page-level 輕量 UI 架構地圖，例如：

```text
docs/UI架構地圖/<page-name>.md
```

不要只把結論寫進 endpoint 文件、工具使用說明或 chat。endpoint 文件回答「API 長什麼樣」，page map 回答「使用者怎麼到這頁、每個 UI 區塊用哪條 API、哪些格式和分頁已驗證、SDK/測試要怎麼覆蓋」。

最小章節：

1. Entry path：從 cold start 或主要 navigation 到此頁的路徑。
2. UI block map：tab、卡片、feed、詳情、評論、廣告或未知區塊。
3. Scroll / pagination：是否可滑動、如何取樣、`pageNumber` / cursor / `hasNext`。
4. Detail formats：至少列已見內容格式與缺口。
5. API chain：入口、列表、詳情、tag、評論、媒體或下載 API。
6. SDK/client mapping：對外欄位、resource order、fallback 與 raw JSON。
7. Validation evidence：final output、test、截圖、UI hierarchy 或去敏 API docs。
8. Follow-up template：下一個類似頁面要照什麼步驟補。

```markdown
## App Architecture Map

### Capture Strategy

| Field | Value |
| --- | --- |
| Mode | lightweight overview / API-first then bind / full operation map |
| Capture budget | main tabs only / key flows only / exhaustive |
| Automation mode | none / adb single-flow / uiautomator / manual replay |
| Reason | avoid device lag / core API unknown / documentation completeness |
| Deferred binding | endpoints or screens to revisit later |

### Navigation Summary

| Area | Visible label | Entry point | Screenshot | Notes |
| --- | --- | --- | --- | --- |
| bottom tab | Home | cold start / bottom nav | `<screenshot-path>` | |
| bottom tab | Search | bottom nav | `<screenshot-path>` | |
| drawer/menu | Profile | avatar/menu tap | `<screenshot-path>` | |

### Screen Inventory

| Screen ID | Canonical route | Screenshot | Scrollable | Clickable entries | Key visible elements | State / Preconditions |
| --- | --- | --- | --- | --- | --- | --- |
| `home.feed` | `launch-authenticated -> home.feed` | `<screenshot-path>` | vertical list: top/mid/bottom sampled | item card, banner, tab buttons | feed list, banner | logged in |
| `item.detail` | `launch-authenticated -> home.feed -> open-detail` | `<screenshot-path>` | no / yes | play button, favorite, related item | title, action buttons | item available |

### Screen Reachability / Operation Recipe

| Route ID | Target screen | Destination scope | Start state | Step order | Step type | Target / Gesture | Selector or coordinate source | Expected result | External transition | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `launch-authenticated` | `home.feed` | in-app | logged-in cold start | 1 | launch | open package | package name | `home.feed` visible | none | may trigger preload APIs |
| `open-detail` | `item.detail` | in-app | `home.feed` top | 1 | tap | first item card | visible title / hierarchy bounds | `item.detail` visible | none | item must be available |
| `reach-feed-bottom` | `home.feed` bottom sample | in-app | `home.feed` top | 1 | swipe | vertical swipe up | screenshot coordinates | feed mid visible | none | bounded scroll |
| `open-external-help` | none | external | `profile.menu` | 1 | tap | Help / Support | label / hierarchy bounds | browser or external WebView opens | browser/custom tab; stop app screen mapping | document trigger only |

### Interaction Inventory

| Interaction ID | Screen ID | Type | Target / Gesture | Selector or coordinate source | Expected result | API capture needed |
| --- | --- | --- | --- | --- | --- | --- |
| `home-scroll-mid` | `home.feed` | swipe | vertical swipe up once | screenshot coordinates / hierarchy bounds | feed mid-page visible | yes/no |
| `open-detail` | `home.feed` | tap | first item card | visible label / bounds / coordinates | `item.detail` | yes |

### Operation To API Matrix

| Operation ID | Route ID | UI path / action | Automation script | Binding phase | Capture window | Method / Path | Source | Response shape | Confidence | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `open-home` | `launch-authenticated` | cold start -> Home | `<script or manual>` | initial map | `<start-end>` | `GET /<path>` | hook / pcap / MITM | top-level keys only | medium | may include preload/cache |
| `open-detail` | `open-detail` | `Home > item tap` | `scripts/ui/open-detail.sh` (`tap`) | after API decoded | `<start-end>` | `POST /<path>` | hook | schema-only summary | high | |
| `scroll-feed` | `reach-feed-bottom` | `Home > swipe feed` | `scripts/ui/scroll-feed.sh` (`swipe`) | UI binding | `<start-end>` | `GET /<path>` | hook / pcap | list page shape | medium | may be pagination/preload |

### Unknown / Untested Navigation

- Screen or tab not yet opened:
- Operation that produced no network:
- API seen without confirmed UI trigger:
- Binding deferred because screenshots/UI traversal were too heavy:
```

文件要求：

- Screenshot 要去敏；不要保留帳號、頭像、電話、email、訂單、私訊或個資。
- 先記主要 tabs/screens 即可；只有高價值流程或需要 attribution 的 API 才補完整操作截圖。
- 每個重要 screen 都要有可引用的 route id；route id 記「怎麼到頁面」，screen id 記「頁面是什麼」。
- Reachability recipe 要能被人工照做，也能被 automation script 改寫成 tap/swipe/launch 步驟。
- UI 架構地圖只把 app 內頁面列入 screen inventory；外部跳轉要寫在 route recipe 的 `Destination scope` / `External transition`，不要繼續當成本 app screen 展開。
- 若跳到系統設定、瀏覽器、支付、分享、第三方 App、外部 intent 或不可控 Web 流程，記錄觸發點、外部目的地類型、是否需要人工接手，以及同窗 API capture window。
- 每個 screen 要標記是否可滑動；滑動頁面只保存代表性 top/mid/bottom 或關鍵分頁，不做無限制全量截圖。
- 每個 clickable entry 要記 target、selector/resource-id/content-desc 或座標來源，以及預期跳轉/操作結果。
- Automation script 只記可重放操作與時間窗；不要把帳密、token、付款、刪除、發文、下單等高風險動作寫成無保護腳本。
- Capture window 要能對齊 pcap/MITM/Frida log 的時間戳或 sequence id。
- API 關聯要寫 `Source`，例如 hook、pcap timing、MITM、replay；只靠 screenshot 不足以證明 API 來源。
- 若某個 API 是 startup/preload/background sync，要在 `Notes` 標明，避免誤判為當前點擊觸發。
- 若採 API-first，先在 API 文件標 `UI path: unknown` / `Trigger confidence: low`，等核心 API 穩定後再回填 UI binding。

## API / Schema 文件模板

```markdown
## Endpoint Name

| Field | Value |
| --- | --- |
| Method | `GET` / `POST` |
| Path | `/path` |
| Auth | Required / Optional |
| Source | pcap / MITM / hook / replay |
| UI path | `Tab > Screen > Action` |
| Operation ID | `open-home` / `open-detail` |
| Trigger confidence | high / medium / low |
| Capability / feature | user-visible function this API supports |
| Domain concept candidates | entity/value object/state names inferred from evidence |
| State impact | creates / reads / updates / deletes / refreshes / paginates / authenticates |

### HTTP Request Headers

| Header | Type / Shape | Meaning | Required | Source | Sensitive | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `Authorization` | bearer / custom / none | session auth | yes/no | token provider | yes | value redacted |
| `User-Agent` | string shape only | client identity | yes/no | app/runtime | no | |

### Request Query / Path Parameters

| Field | Type / Shape | Meaning | Required | Source | Sensitive | Notes |
| --- | --- | --- | --- | --- | --- | --- |

### Request Body

| Field | Type / Shape | Meaning | Required | Source | Sensitive | Notes |
| --- | --- | --- | --- | --- | --- | --- |

### Response Wrapper

| Field | Type / Shape | Meaning | Required / Optional | Notes |
| --- | --- | --- | --- | --- |

### Decrypted / Inner Payload

| Field | Type / Shape | Meaning | Required / Optional | Source / Derived From | Notes |
| --- | --- | --- | --- | --- | --- |

### Response Headers

| Header | Type / Shape | Meaning | Notes |
| --- | --- | --- | --- |

### Evidence

- Sanitized log:
- Fixture:
- UI path:
- Screenshot / UI evidence:

### Validation

- Replay:
- Contract test:
- Manual verification:

### Reconstruction Notes

- BDD scenario candidate:
- Domain Model Contract candidates:
- API / Interface Contract notes:
- Error Handling Contract notes:
- Fixtures needed for rebuild:
- Open questions:
```

API 文件要求：

- 分析完 API 後要回填專案文件；不要只把 endpoint 留在暫存 log。
- HTTP/HTTPS API 必須記錄可見的 headers、request、response；看不到的部分要寫明是 MITM 不可見、hook 未到位、加密包裹、或尚未驗證。
- 每個 request/response 字段都要逐欄位分析 type/shape、meaning、required/optional、source/derived-from、敏感性與備註。
- 每個高價值 API 都要標明支援哪個 capability、對應 operation id、可能的 domain concept、狀態影響、錯誤/空狀態與 fixture；這些欄位是後續用 `app-development-guidance` 重建功能的輸入。
- Header 名稱、path shape、query key、schema 可以保留；header value、token、cookie、device id、個資與可重放 URL 必須去敏。
- 截圖可用來輔助說明 UI path、tab、screen 與操作，但不能取代 HTTP header/request/response 的字段分析。

## 去敏規則

必須遮蔽：

- `Authorization`、cookie、session token。
- device id、install id、advertising id。
- 真實帳號、電話、email、邀請碼。
- AES/HMAC key material。
- 能直接重放付費內容或個人內容的 URL。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置。請改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等 placeholder。

可以保留：

- header 名稱。
- path shape。
- query key 名稱。
- response top-level schema。
- schema-only JSON 摘要：字串 length/hash、top-level keys、key type；不要保留 value。
- 已去敏的 fixture。
- magic bytes、容器格式、演算法步驟。

## 證據鏈要求

好文件不只寫「成功」，還要寫為什麼相信它成功：

- pcap 證明對外 TLS host 存在。
- proxy CONNECT 證明導流成功。
- hook log 證明 request object 在 TLS 前可見。
- decrypt hook 或離線 decoder 證明 inner JSON 正確。
- fixture / test 證明規則可重跑。

## 失敗也要記錄

失敗紀錄應包含：

- 嘗試了什麼。
- 期望看到什麼。
- 實際看到什麼。
- 排除了什麼假設。
- 是否要重試，或是否停止投入。

例：

```text
Java OkHttp hook installed successfully, but no target host/path appeared while pcap showed TLS traffic to the API host. This rules out the Java OkHttp path for the tested flow and shifts the next step to native/Flutter analysis.
```

## Developer Guidance Notes（可選）

若分析結果能轉成「未來開發自家 App 時可採取的設計、實作或安全做法」，可在專案分析文件加一小節：

```markdown
## Developer Guidance Notes

| Observation | Development Guidance | Owner | Validation |
| --- | --- | --- | --- |
| 已去敏觀察 | 可重用的開發建議 | client / API / backend / build / monitoring | 測試或 review 方法 |
```

這一節只寫已去敏、可泛化的開發啟發。成熟後把 App 開發 guidance 回饋到 [`app-development-guidance`](../app-development-guidance/)；本 `apk-analysis` skill 只保留分析方法、證據鏈與工具判斷。若使用者要用這些文件實作 app 工具、SDK 或 client，必須立即套用 `app-development-guidance`，不要在 `apk-analysis` 文件裡直接展開產品開發流程。

## 技巧回饋文件要給人讀

寫入 **`feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`**（跨分類用 `feedback_history/common/`；規則與模板見 [`../../shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md)）時，請避免只有工具名與短結論。每條技巧都應包含：

- `One-line Summary`：一句話講重點。
- `Human Explanation`：給人看的背景與誤判風險。
- `Trigger`：什麼現象會觸發這條技巧。
- `Evidence`：去敏證據或觀察。
- `Generalized Lesson`：抽象後的通用規則。
- `Agent Action`：下次 AI 要採取的具體行動。
- `Applies When` / `Does Not Apply When`：適用邊界。
- `Validation`：怎麼確認這條技巧有效。

好的 lesson 應該像這樣：

```markdown
### Proxy failure 要先拆成導流與 TLS 兩層

One-line Summary:
代理看不到明文時，先確認「有沒有進代理」，再談憑證或 pinning。

Human Explanation:
很多人看到 MITM 沒有明文就直接判斷是 pinning。更可靠的順序是先看 App 是否真的連到 proxy。如果仍直連目標 host，問題在導流或初始化時機；如果已經進 proxy 才 TLS failed，才查 CA / pinning。

Agent Action:
先檢查 CONNECT 或 connect target，不要先寫 pinning 結論。
```

## 回填規則

每次分析完成後：

- UI Behavior 必須回填專案 UI 行為入口或 page-level map（例如 `docs/UI-Behavior.md`、`docs/UI架構地圖/<page>.md` 或專案等價位置）：記錄 entry path、可見 UI blocks、App sort label、tap/swipe/input 操作、API/data source 對照、截圖/UI hierarchy/live replay/hook 證據與 unknowns。若沒有 UI 證據，明確標 `needs capture`、`needs replay` 或 `Trigger confidence: low`。
- 目標 API 結論回填專案 API 文件。
- 解碼規則回填協議/解密文件。
- SDK 或 client 行為回填 BDD / tests。
- 若分析文件要用來做 app 工具、SDK、client、mock、fixture-driven implementation、contract test 或重建功能，同輪自動啟用 [`app-development-guidance`](../app-development-guidance/) 並交出 Feature Reconstruction Handoff；不要讓開發規格停留在 APK 分析文件內。
- 通用技巧回填 **`feedback_history/<category>/`** 或 **`feedback_history/common/`**（新檔），驗證後再整理到本 skill 的主文件或對應 `techniques/<category>/`。
- App 開發 guidance 回填 [`app-development-guidance`](../app-development-guidance/)；不要把產品開發 checklist 長期堆在 `apk-analysis`。
