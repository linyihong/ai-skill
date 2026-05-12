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

---

← [回到 workflow/apk-analysis/](README.md)
