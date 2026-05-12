# APK Analysis Artifact Gates

本文件定義 APK 分析過程中必須產出的 artifacts 與完成門檻。分析方法見 [`analysis/apk/`](../../analysis/apk/)；模板見本目錄的 templates。

## 1. UI Architecture Map

當使用者要求分析某個具名 page、tab、module 或功能入口，且已建立 UI/API 對照時，必須在專案 docs 產出或更新 page-level UI 架構地圖。

### 最小章節

1. Entry path：從 cold start 或主要 navigation 到此頁的路徑。
2. UI block map：tab、卡片、feed、詳情、評論、廣告或未知區塊。
3. Scroll / pagination：是否可滑動、如何取樣、`pageNumber` / cursor / `hasNext`。
4. Detail formats：至少列已見內容格式與缺口。
5. API chain：入口、列表、詳情、tag、評論、媒體或下載 API。
6. SDK/client mapping：對外欄位、resource order、fallback 與 raw JSON。
7. Validation evidence：final output、test、截圖、UI hierarchy 或去敏 API docs。
8. Follow-up template：下一個類似頁面要照什麼步驟補。

### 規則

- Screenshot 要去敏。
- 先記主要 tabs/screens 即可；高價值流程才補完整操作截圖。
- 每個重要 screen 要有可引用的 route id 與 screen id。
- Reachability recipe 要能被人工照做，也能被 automation script 改寫。
- UI 架構地圖只把 app 內頁面列入 screen inventory；外部跳轉要寫在 route recipe 的 `Destination scope` / `External transition`。
- 每個 screen 要標記是否可滑動；滑動頁面只保存代表性 top/mid/bottom。
- 每個 clickable entry 要記 target、selector/resource-id/content-desc 或座標來源。
- Capture window 要能對齊 pcap/MITM/Frida log 的時間戳。
- API 關聯要寫 `Source`（hook / pcap timing / MITM / replay）。
- 若某個 API 是 startup/preload/background sync，要在 `Notes` 標明。

## 2. API Catalog

當分析目標包含「整理 API 列表、SDK/client 對照、mock API、contract test、功能重建」時，專案文件應建立一組可維護的 API Catalog。

### 建議結構

```text
docs/API.md                         # API / host / traffic family 總入口
docs/API/<group>/README.md          # 依 path 第一段、domain、feature 或 protocol family 分組
docs/API/<group>/<operation>.md      # 單支 API 詳細文件
docs/API/coverage.md                # 已觀測、已 replay、已解密、待補與不適用清單
docs/API/supplement/<topic>.md      # HLS、media、local bridge、SDK call order 等跨 API 主題
```

### Catalog minimum

| Artifact | 必填內容 |
| --- | --- |
| API 總入口 | 已知 host/base、traffic family、response wrapper、auth/session/header 共用規則、解密/解碼入口、覆蓋率文件、UI map 入口、SDK/client 入口。 |
| 分組索引 | 分組依據、每支 API 的 method/path、request 摘要、response 用途、目前用途/結論、詳細文件連結。 |
| 單支 API 文件 | Method/path、host/base、auth/session、headers、query/path/body、response wrapper、inner payload、item schema、error/empty behavior、pagination/cache、field meaning、sensitivity、evidence、validation。 |
| Coverage / gap matrix | 靜態枚舉、動態觀測、MITM、pcap、hook、replay、decrypted fixture、contract test、UI binding、缺參、未觸發、暫不測、scope out。 |
| UI/API 對照 | UI map、route id、operation id、trigger confidence、capture window、startup/preload/background 標記。 |
| SDK/client 欄位用途 | SDK/client 實際讀取或轉換的欄位、相容性範圍、raw JSON 保留策略、fixture/test 對照。 |
| Cross-flow docs | 播放鏈、media chain、login/session、local bridge、vendor/service split、分頁與排序等跨多支 API 的流程。 |
| Sanitization | 哪些值已遮蔽、哪些 raw evidence 留在受控位置、哪些文件可 commit。 |

### Catalog completion gate

當使用者問「API 列表是否完整」、「能不能做 SDK/client/mock」時，完成回覆前要檢查：

- API 總入口是否連到分組、coverage/gap、UI map、解碼/共用 wrapper、SDK/client 文件。
- 已觀測 API 是否都落到分組索引。
- 高價值 API 是否有單支詳細文件。
- 每支 API 是否有 request、response、field meaning、evidence、validation/open questions。
- UI trigger 若未確認，是否標 `UI path: unknown` / `Trigger confidence: low`。

## 3. Domain/Runtime Baseline

**問題：**只做「逐支 API 的 request/response shape」時，外包裝/SDK 仍可無法連線。

### 建議章節

| 章節 | 內容（去敏） |
| --- | --- |
| 環境維度 | 觀察到的 host family、path family、是否多 CDN／多 gateway、與 build／地區是否相關。 |
| 連線路徑 | App 是否走系統代理、內建 TUN、local proxy、直连；與 capture 工具相容性。 |
| Session／身分 | 列表 API 是否在未登入下可用；若否，登入／裝置／device id 與列表欄位的因果鏈。 |
| Opaque／衍生參數 | 哪些 query 由前序 response、WebView、搜尋 session 或固定 app 常數提供。 |
| 簽章與 gateway | service／hash、header 名稱集合、canonical path 規則。 |
| 分頁地面真相 | 是否有 `has_next` 類欄位；若無，記錄啟發式與反例風險。 |
| 錯誤與限流 | 影響重試的 code、冷卻、與 session 刷新關係。 |
| 重放檢查清單 | 人工或腳本重放同一列表的最小步驟。 |

### Development readiness gate

若下一步是 SDK/client/app tool/live integration 開發，baseline 不能只停在 skeleton。必須先檢查並記錄最小可跑因素：

- endpoint/path family
- route/service 對照或 adapter 策略
- session/bootstrap 依存
- opaque 參數來源與時效
- 簽章/gateway 前置
- response decrypt/unwrap 邊界
- 分頁地面真相
- 錯誤/session 恢復
- 重放檢查清單

缺任一項時，該缺口必須成為開發 blocker 或被明確 scoped out。

## 4. Feature Reconstruction Handoff

若分析目標是讓後續 agent 能用 `app-development-guidance` 重新做出同等功能，專案分析文件不能只列 endpoint。

### 最低表格

| 面向 | 必填內容 |
| --- | --- |
| Feature / Capability | 功能名稱、使用者目標、入口 screen、非目標或未知限制。 |
| UI Behavior | screen id、route id、operation id、前置狀態、tap/swipe/input 步驟、可見結果。 |
| Domain Concepts | 從 UI 文案、response fields、狀態碼推得的 entity、value object、state、command、event。 |
| API / Interface Contract | method/path shape、headers、query/body、response wrapper、inner payload、auth/session、pagination、cache、idempotency。 |
| State And Error Handling | loading/empty/error/success 狀態、錯誤碼、重試、登入過期、權限不足、限流、離線或快取行為。 |
| Data Lifecycle | 欄位來源、derived-from、local cache/storage、刷新時機、敏感性、保留/過期行為。 |
| Validation Evidence | pcap/MITM/hook/replay/fixture/screenshot/UI hierarchy/automation script 的去敏引用。 |
| Unknowns / Assumptions | 未觸發流程、低信心 mapping、缺少樣本、未驗證 edge case。 |

### Feature handoff finish gate

當某個具名 feature/tab/module 已被分析到「核心 UI 操作與主要 API flow 可說明」的程度時，必須在同一輪補齊或更新 project-level feature handoff 文件。

觸發條件包含任一項：

- 核心 flows 已從 `Candidate` 升到 `Confirmed`。
- agent 已能回答此功能的 entry path、主要 UI 區塊、API request keys / response schema、狀態與缺口。
- 使用者問「有沒有 API 文件」、「能不能重建」、「架構是什麼」。

## 5. 單次分析筆記模板

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

## 6. 證據鏈要求

好文件不只寫「成功」，還要寫為什麼相信它成功：

- pcap 證明對外 TLS host 存在。
- proxy CONNECT 證明導流成功。
- hook log 證明 request object 在 TLS 前可見。
- decrypt hook 或離線 decoder 證明 inner JSON 正確。
- fixture / test 證明規則可重跑。

## 7. 失敗也要記錄

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

## 8. SDK Live Self-Generation Audit

當使用者的目標是「像某些既有 SDK 一樣，除了呼叫方合法提供的 **授權身分材料** 外，其餘 host、路由、簽章、session、decrypt 都能由 SDK 自行生成」時，必須在專案 baseline 或 SDK-readiness 文件加一張 **runtime factor classification** 表：

| Classification | 意義 | 可開始 live SDK self-generation? |
| --- | --- | --- |
| `sdk-generatable` | 可由公開 SDK 程式、常數、演算法、穩定 public config 或已去敏規則自行生成；不需要私有 runtime bridge。 | 是 |
| `identity-material-bound` | 需要授權方提供或初始化的身分材料，例如 device id、install id、授權帳號、session seed、合法裝置初始化結果。 | 可，若這是唯一剩餘未知或唯一使用者提供項 |
| `private-adapter-required` | 需要 raw service、私有 host 選擇、簽章 key、decrypt key、in-app bridge、未公開 provider，或只能靠 app runtime 生成。 | 否 |
| `unknown` | 還不知道來源、時效、錯誤行為或是否可重建。 | 否 |
| `scoped-out` | 不屬於本 SDK live scope（例如 media download、write actions）。 | 不阻塞該 scope，但必須明寫 |

建議至少列：

| Runtime factor | 必問問題 |
| --- | --- |
| Base endpoint / host | SDK 能否從固定 fallback、public config、DNS/config API 自行選擇？還是必須 private host table / app storage？ |
| Route binding / service | raw route id/service 是否可由 SDK deterministic 生成？ |
| Authorized identity material | 是否只剩 device id / install id / 授權帳號 / session seed 等呼叫方身分材料需要注入？每個 key group 能否由 SDK/tool 自行生成或初始化？ |
| Session/bootstrap | guest/device login 是否可由 SDK 生成？是否仍需要 app-only token、captcha、human login 或私有 WebView state？ |
| Opaque query/header | 每個 opaque 欄位是 app 常數、locale、device/session 派生，還是 response/session 私有值？ |
| Signing/gateway | canonicalization、排序、hash/HMAC/AES、timestamp/random 來源是否可重建且已 fixture 驗證？ |
| Response decrypt/unwrap | SDK 是否能自行把 wire response 解成 JSON？key/IV/KDF 是否還在 private app helper？ |
| Error/session recovery | token 過期、bad signature、bad device、bad opaque 的 code 與 refresh/writeback 是否已 live matrix 驗證？ |
| Pagination/data truth | 是否已知如何終止分頁、辨識空資料、避免把錯誤 envelope 當空列表？ |
| Media/download | 若 scope 包含媒體，signed URL、key unwrap、decrypt、package 是否可重建；不含媒體時標 scoped-out。 |

完成後給出 verdict：

```text
Live SDK self-generation verdict:
- ready except authorized identity material: yes/no
- remaining non-device blockers: <factor list>
- allowed next work: live SDK implementation / private adapter only / offline parser only
```

**判讀規則：**只要仍有 `private-adapter-required` 或 `unknown` 的 base host、route service、signing、decrypt、session bootstrap、opaque provider，就不能說「只剩授權身分材料」。

## 9. Authorized Identity Material Self-Generation Audit

當 runtime 需要 device / install / account / session seed / vendor attestation / server-issued session 類材料時，必須逐 key group 回答「能否自生成」與「怎麼生成」。

最低表格：

| Field | 必填內容 |
| --- | --- |
| Key group / surface | 欄位名稱群、storage key 名、request key 名或 provider function boundary；只寫 name/shape，不寫 raw value。 |
| Role in live access | 它是 app/build constant、device/install material、account material、guest/session seed、vendor attestation、server-issued session，還是其它 runtime factor。 |
| Self-generation verdict | `sdk-generatable` / `caller-provided` / `server-issued` / `trusted-bridge` / `private-adapter-required` / `unknown` / `scoped-out`。 |
| Generation recipe or provider boundary | 若 `sdk-generatable`，寫 sanitized recipe：inputs、algorithm family、canonical order、storage key name、refresh trigger、validation fixture；若不能，寫由 caller、server response、trusted bridge 或 private adapter 提供。 |
| Lifecycle and reset behavior | first install、cold start、guest/login、preserved session、logout、`clear app data`、reinstall、token expiry 時如何建立、重用、更新或清除。 |
| Cooldown / risk controls | 是否會觸發 rate limit、device health、attestation check、account lock、captcha/human step；只寫 status/error class，不寫可濫用細節。 |
| Error / negative matrix | missing、empty、stale、bad-fixed、bad-signature、expired-session 等情況的 wrapper/UI/recovery class；若未驗證，標 `pending`。 |
| Validation evidence | 去敏 hook summary、static provider trace、fixture、replay parity、unit/contract test；不得保存 raw token、device id、account、vendor payload、signature 或 host。 |

判斷規則：

- `sdk-generatable` 需要可重跑的生成 recipe 或測試，不能只因為值看起來像 UUID、hash、locale、random 就宣稱可生成。
- `caller-provided` 仍需定義 lifecycle、reset/cooldown、health/error 行為；否則是 `unknown` 或 `private-adapter-required`。
- `server-issued` 可以由 SDK 建模 storage/refresh boundary，但 raw material 來自授權 server response。
- `trusted-bridge` / `private-adapter-required` 可以讓 private live smoke 成立，但不能支撐「standalone self-generating SDK」宣稱。
- 若任何 live-required identity key group 的 generation recipe、provider boundary、reset/cooldown 或 error matrix 是 `unknown`，live-facing development 只能繼續在 private adapter / bridge scope。

**Finish gate：** 若本輪目標包含「可程式化拉取真實資料」「接 SDK transport」「寫 integration test」之一，而專案尚無 baseline 或僅有 API 條目：必須在同一工作單建立 **skeleton baseline**，並把 **open** 項寫成可驗證問題。若本輪要開始開發 SDK、client、app tool、live integration，baseline 不能只停在 skeleton。若 live flow 需要 device/install/account/session/vendor/server-issued material，必須先補 authorized identity material self-generation audit。

## 10. UI Architecture Map Template

用 screenshot、UI hierarchy 與可重放操作，把 App 的可見架構寫成地圖。這份文件放專案分析文件，不放 reusable skill，skill 只保留模板與方法。

### Page-level artifact rule

當使用者要求分析某個具名 page、tab、module 或功能入口，且本次工作已建立 UI/API 對照時，必須在專案 docs 產出或更新一份 page-level 輕量 UI 架構地圖。

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

### Interaction Inventory

| Interaction ID | Screen ID | Type | Target / Gesture | Selector or coordinate source | Expected result | API capture needed |
| --- | --- | --- | --- | --- | --- | --- |
| `home-scroll-mid` | `home.feed` | swipe | vertical swipe up once | screenshot coordinates / hierarchy bounds | feed mid-page visible | yes/no |
| `open-detail` | `home.feed` | tap | first item card | visible label / bounds / coordinates | `item.detail` | yes |

### Operation To API Matrix

| Operation ID | Route ID | UI path / action | Automation script | Binding phase | Capture window | Method / Path | Source | Response shape | Confidence | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `open-home` | `launch-authenticated` | cold start -> Home | `<script or manual>` | initial map | `<start-end>` | `GET /<path>` | hook / pcap / MITM | top-level keys only | medium | may include preload/cache |
| `open-detail` | `open-detail` | `Home > item tap` | `scripts/ui/open-detail.sh` | after API decoded | `<start-end>` | `POST /<path>` | hook | schema-only summary | high | |

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
- UI 架構地圖只把 app 內頁面列入 screen inventory；外部跳轉要寫在 route recipe 的 `Destination scope` / `External transition`。
- 若跳到系統設定、瀏覽器、支付、分享、第三方 App、外部 intent 或不可控 Web 流程，記錄觸發點、外部目的地類型、是否需要人工接手，以及同窗 API capture window。
- 若分析主線不是 UI 地圖（例如 provider、session、signing、pagination、storage），但動態操作新發現了 app 內 route、settings screen、global menu、dialog、tab 或可重放操作，完成前仍要回填或建立對應 `docs/UI架構地圖/<route-or-area>.md`。
- 每個 screen 要標記是否可滑動；滑動頁面只保存代表性 top/mid/bottom 或關鍵分頁，不做無限制全量截圖。
- 每個 clickable entry 要記 target、selector/resource-id/content-desc 或座標來源，以及預期跳轉/操作結果。
- Automation script 只記可重放操作與時間窗；不要把帳密、token、付款、刪除、發文、下單等高風險動作寫成無保護腳本。
- Capture window 要能對齊 pcap/MITM/Frida log 的時間戳或 sequence id。
- API 關聯要寫 `Source`，例如 hook、pcap timing、MITM、replay；只靠 screenshot 不足以證明 API 來源。
- 若某個 API 是 startup/preload/background sync，要在 `Notes` 標明，避免誤判為當前點擊觸發。
- 若採 API-first，先在 API 文件標 `UI path: unknown` / `Trigger confidence: low`，等核心 API 穩定後再回填 UI binding。

## 11. API Catalog Detail Requirements

單支 API 文件至少要能回答：

| Area | Required detail |
| --- | --- |
| Identity | Method、host/base、path shape、operation id、分組、狀態：candidate / observed / replayed / decoded / validated / deprecated / out of scope。 |
| Request | headers、path/query/body 欄位、型別/shape、用途、必填/選填、來源、敏感性、是否參與 signing/encryption。 |
| Response | raw wrapper、decrypted/inner payload、list item schema、欄位型別、nullable/optional、欄位語意、derived-from、下游 API key。 |
| Behavior | capability、UI trigger、startup/preload/background 判斷、state impact、empty/error behavior、pagination/cache/sort semantics。 |
| Evidence | hook/MITM/pcap/replay/fixture/screenshot/UI hierarchy/automation script 的去敏引用。 |
| Validation | replay result、decoder fixture、schema assertion、SDK/client test、contract test、manual evidence，或明確標 `needs capture` / `needs replay`。 |
| Open questions | 缺少樣本、低信心 field meaning、未驗證 edge case、需要使用者或更多操作證據的 blocker。 |

### Field meaning rule

Schema 不只是型別表。欄位要盡量寫出用途：

- 哪些欄位會成為下一支 API 的 request key。
- 哪些欄位控制 UI 顯示、分頁停止、排序、播放、下載、收藏、權限或錯誤狀態。
- 哪些欄位只是樣本中出現但用途未知，必須標 `meaning unknown` / `candidate`。
- 哪些欄位是 SDK/client 已使用欄位，變動會破壞相容性。

不要把推測寫成確認規格。若只有少量 live/replay 樣本，使用 `candidate`、`sample only`、`needs more samples` 或 `low confidence`。

### API / Schema Document Template

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
- 每個高價值 API 都要標明支援哪個 capability、對應 operation id、可能的 domain concept、狀態影響、錯誤/空狀態與 fixture。
- Header 名稱、path shape、query key、schema 可以保留；header value、token、cookie、device id、個資與可重放 URL 必須去敏。
- 截圖可用來輔助說明 UI path、tab、screen 與操作，但不能取代 HTTP header/request/response 的字段分析。

## 12. Sanitization Rules

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

## 13. Developer Guidance Notes（可選）

若分析結果能轉成「未來開發自家 App 時可採取的設計、實作或安全做法」，可在專案分析文件加一小節：

```markdown
## Developer Guidance Notes

| Observation | Development Guidance | Owner | Validation |
| --- | --- | --- | --- |
| 已去敏觀察 | 可重用的開發建議 | client / API / backend / build / monitoring | 測試或 review 方法 |
```

這一節只寫已去敏、可泛化的開發啟發。成熟後把 App 開發 guidance 回饋到 `app-development-guidance`；本 `apk-analysis` 文件只保留分析方法、證據鏈與工具判斷。

## 14. Feedback Lesson Writing Tips

寫入 `feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md` 時，請避免只有工具名與短結論。每條技巧都應包含：

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

## 15. Backfill Rules

每次分析完成後：

- UI Behavior 必須回填專案 UI 行為入口或 page-level map：記錄 entry path、可見 UI blocks、App sort label、tap/swipe/input 操作、API/data source 對照、截圖/UI hierarchy/live replay/hook 證據與 unknowns。若沒有 UI 證據，明確標 `needs capture`、`needs replay` 或 `Trigger confidence: low`。
- 目標 API 結論回填專案 API 文件。
- 解碼規則回填協議/解密文件。
- SDK 或 client 行為回填 BDD / tests。
- 若分析文件要用來做 app 工具、SDK、client、mock、fixture-driven implementation、contract test 或重建功能，同輪自動啟用 `app-development-guidance` 並交出 Feature Reconstruction Handoff；不要讓開發規格停留在 APK 分析文件內。
- 通用技巧回填 `feedback_history/<category>/` 或 `feedback_history/common/`（新檔），驗證後再整理到主文件或對應 `techniques/<category>/`。
- App 開發 guidance 回填 `app-development-guidance`；不要把產品開發 checklist 長期堆在 `apk-analysis`。

---

← [回到 workflow/apk-analysis/](README.md)
