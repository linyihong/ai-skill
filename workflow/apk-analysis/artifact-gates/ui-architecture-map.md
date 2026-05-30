# UI Architecture Map Slice（page-level UI 架構地圖 + 完整模板）

> **Cognitive Slice**：`apk-ui-architecture-map`（從 [`../artifact-gates.md`](../artifact-gates.md) §1+§10 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-ui-architecture-map` |
| `purpose` | 為具名 page/tab/module 建立或更新 page-level UI 架構地圖（含 operation-to-API matrix） |
| `type` | `execution` |
| `tags` | artifact-gate, ui |
| `load_when` | 建立或更新 APK UI architecture map、需要 navigation segment library / operation-to-API 對照 |
| `do_not_load_when` | 純後台 / 純網路分析無 UI 觀察需要、僅做 API catalog 不涉及 UI binding |
| `owner_layer` | workflow |
| `layer_justification` | 規定「UI map 要產出哪些章節、要過哪些 reachability/automation 規則」的 ordering / artifact gate；通過 workflow membership test |
| `canonical_source` | 本檔（原 `artifact-gates.md` §1 UI Architecture Map + §10 UI Architecture Map Template） |
| `dependencies` | `apk-api-catalog`（API binding 章節相互引用）、`analysis/apk/` evidence acquisition methods |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Scenario AG-A（execution-only：UI map 已存在的任務應 **不** 載入本 slice） |

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
- 可重用 App 操控要拆成具名 navigation segments，而不是只保存一條長腳本；每段要能從已知 checkpoint 進入、停在已知 checkpoint，供後續自由組合。
- UI 架構地圖只把 app 內頁面列入 screen inventory；外部跳轉要寫在 route recipe 的 `Destination scope` / `External transition`。
- 每個 screen 要標記是否可滑動；滑動頁面只保存代表性 top/mid/bottom。
- 每個 clickable entry 要記 target、selector/resource-id/content-desc 或座標來源。
- Capture window 要能對齊 pcap/MITM/Frida log 的時間戳。
- API 關聯要寫 `Source`（hook / pcap timing / MITM / replay）。
- 若某個 API 是 startup/preload/background sync，要在 `Notes` 標明。

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

### Navigation Segment Library（UI map 的可組合操控腳本登記處）

| Segment ID | Entry checkpoint | Action | Exit checkpoint | Script / function | Preconditions | Evidence | Reusable by |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `launch-to-home` | app not running / force-stopped | launch and wait for authenticated home | `home.feed` | `scripts/ui/segments/launch-to-home.sh` | valid test session | screenshot + hierarchy + package check | any home-based flow |
| `home-to-detail` | `home.feed` | tap first visible item card | `item.detail` | `scripts/ui/segments/home-to-detail.sh` | item list visible | screenshot + hierarchy + timestamp | detail API capture, media capture |

文件要求：UI map 是 navigation segment 的 source-of-truth；script 檔案只是執行物。每個 segment 都要在 UI map 記 entry / exit checkpoint、前置狀態、selector 或 coordinate 來源、script path、輸出 evidence 與可重用範圍。組合流程時引用 segment IDs，例如 `launch-to-home -> home-to-detail`；若某段失效，只重測該 segment，不重跑整條 route。後續 agent 想去任何頁面時，應先查 `docs/UI架構地圖/<route-or-area>.md` 的 segment library，再組合既有 segments。

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
- App 操控腳本要優先拆成可重用 navigation segments，並以 UI map 作為登記處：記錄 segment id、entry checkpoint、exit checkpoint、script path、preconditions、evidence 與可組合範例；後續 capture 應先查 UI map 組合既有 segments，避免重複從頭測試。
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

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
