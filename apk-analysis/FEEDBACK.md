# Skill 回饋與強化紀錄

這份文件用來累積新的 APK 分析技巧。當 agent 或人類在分析中發現新想法、失敗模式、工具用法或驗證規則時，先寫到這裡。確認可重用後，再整理進 `WORKFLOW.md`、`TOOLS.md` 或 `DOCUMENTATION.md`。

## 回饋原則

- **Cursor agent：** 在授權 APK 分析過程中一旦得到可重用技巧／失敗模式／驗證規則，應**主動**依「新 lesson 模板」**同一輪對話內**寫入本檔，**不要**等使用者提醒「記得回饋」；去敏與否適用條件仍須遵守。任務結束前可對照 `SKILL.md` 的 Feedback Loop 檢查是否漏寫。
- 只寫通用方法，不寫特定 App 的私有結論。
- 必須去敏。
- 不得寫入本機絕對路徑、使用者名稱、私有工作目錄、clone 位置或只能在單一機器成立的路徑；用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 這類 placeholder 代替。
- 必須說明證據。
- 必須說明適用條件與不適用條件。
- 不確定的想法要標成 `experimental`。
- 已驗證且常用的規則才 promoted 到主文件。

## 人類也能讀的寫法

每一條 lesson 都要讓沒有參與當次分析的人也看得懂。不要只寫「hook X」「改 Y」這種只有 AI 或當事人懂的短句。

建議每條 lesson 都包含三層：

- 一句話摘要：這條規則在提醒什麼。
- 人類說明：為什麼會踩坑、怎麼判斷、下一步怎麼做。
- Agent 指令：下次 AI 遇到類似情境時要採取的具體行動。

## 新 lesson 模板

```markdown
## YYYY-MM-DD - [short title]

Status: candidate | validated | deprecated | promoted

### One-line Summary

用一句人話說明這條 lesson。

### Human Explanation

給人看的背景說明：為什麼重要、常見誤判是什麼、實務上怎麼判斷。

### Trigger

遇到什麼現象或問題？

### Evidence

- Tool:
- Sanitized excerpt:
- Evidence path:

### Generalized Lesson

可重用的規則是什麼？

### Agent Action

下次 agent 看到類似情境時，應該先做什麼、不要做什麼？

### Applies When

- 條件 1
- 條件 2

### Does Not Apply When

- 條件 1

### Validation

如何確認這條 lesson 是對的？

### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`
- `DOCUMENTATION.md`
- `SKILL.md`
```

## 待沉澱候選

### 2026-04-30 - Proxy failure 要先拆成導流與 TLS 兩層

Status: promoted

#### One-line Summary

代理看不到明文時，先確認「有沒有進代理」，再談憑證或 pinning。

#### Human Explanation

很多人看到 Proxyman / Burp / mitmproxy 沒有明文，就直接判斷是 certificate pinning。這常常太早下結論。更可靠的判斷順序是先看 App 的連線目標是否已經變成 proxy；如果還是直連目標 host，問題是導流或初始化時機，不是 TLS。只有已經進 proxy 且 handshake 失敗時，才應該查 CA trust、network security config、custom trust 或 pinning。

#### Trigger

MITM 工具沒有看到明文，或顯示 `SSL Handshake Failed`。

#### Evidence

在授權 APK 分析中，曾觀察到兩種完全不同的 failure：

- 流量沒有進 proxy，裝置仍直接連目標 host `:443`。
- 流量已進 proxy，但 TLS handshake 因 CA trust / pinning 失敗。

#### Generalized Lesson

不要把「代理工具看不到明文」直接等同於 pinning。先看是否有 CONNECT / connect target 到 proxy；只有導流成功後，才進入 CA / pinning 排查。

#### Agent Action

下次遇到 MITM 失敗時，先要求或執行導流驗證：檢查 proxy 是否收到 CONNECT，或用 connect trace 觀察目標是否為 `<proxy-host>:<proxy-port>`。不要先寫 pinning 結論。

#### Promotion Target

已整理到 `WORKFLOW.md` 與 `TOOLS.md`。

### 2026-04-30 - 冷啟動比 attach 更適合抓初始化與代理設定

Status: promoted

#### One-line Summary

要影響 App 的網路 client 初始化，通常要用冷啟動 `spawn`，不要等 App 跑起來後才 attach。

#### Human Explanation

許多 App 會在啟動早期建立 HTTP client、讀 proxy 設定、初始化 token、載入 domain 或建立 session。等 App 已經跑起來後再 attach，常常只能看到後半段行為；即使成功設定環境變數，也可能太晚，原本的 client 已經決定 direct connection 或已經緩存連線設定。

#### Trigger

- attach 已啟動 PID 後看不到初始化 API。
- proxy env 設定成功，但核心流量仍直連。
- 需要抓 token refresh、domain selection、module/bootstrap、first request。

#### Evidence

- 已啟動後 attach：proxy env / hook 可能載入成功，但 connect target 仍是目標 host。
- 冷啟動 spawn：在 client 初始化前注入，connect target 才可能變成 proxy。

#### Generalized Lesson

初始化相關分析優先使用 cold-start spawn。attach 適合補抓已進入頁面後的互動流程，但不適合判斷啟動期網路設定。

#### Agent Action

當使用者要抓「App 一開啟就發生的網路流程」或「讓 Proxyman 能連得到」時，建議流程應包含：force-stop app、Frida spawn、在 resume 前載入 hook / proxy env、驗證 connect target。

#### Promotion Target

已整理到 `WORKFLOW.md`；可視需要補工具腳本範本。

### 2026-04-30 - 高語意 hook 優先於 socket bytes

Status: promoted

#### One-line Summary

能 hook request/response 物件，就不要先從 socket bytes 開始拼。

#### Human Explanation

socket、TLS read/write、`send`/`recv` 事件很多，容易卡 App，也需要自己重組 HTTP、解壓縮、切分 body。高語意 hook 例如 request options、response interceptor、decrypt function，通常事件少、內容接近業務語意，更適合建立 API 文件與測試 fixture。

#### Trigger

低層 socket / TLS hook 事件量大，容易造成 App 卡頓，也需要自行重組 HTTP。

#### Generalized Lesson

優先找 request options、response interceptor、decode/decrypt function。只有在高語意點找不到或需要補證據時，才降到 socket / TLS 層。

#### Agent Action

看到 socket hook 卡頓、ANR、輸出爆量時，停止擴大低層 hook，改回靜態搜尋 request builder / interceptor / decoder，或縮小 hook 條件。

#### Promotion Target

已整理到 `WORKFLOW.md`。

### 2026-04-30 - 動態 hook 結果要離線化

Status: promoted

#### One-line Summary

Frida 是拿樣本的手段，不應成為長期測試唯一依賴。

#### Human Explanation

動態 hook 很適合第一次看清楚 request 和 decrypted response，但它依賴裝置、App 版本、時機與 hook 穩定性。真正能讓後續 SDK 或工具穩定演進的是：把 raw response、解密輸入輸出、schema mapping 做成 fixture，並用離線 decoder 或 contract test 驗證。

#### Trigger

每次都依賴 Frida 才能驗證 API 或解密，會讓後續 SDK / client 測試不穩定。

#### Generalized Lesson

把 raw wrapper、encrypted payload、decrypted payload、schema mapping 做成 fixture 或 contract test。動態 hook 是取得樣本的手段，不應是長期驗證唯一來源。

#### Agent Action

每次拿到新解密結果時，詢問或建立：sanitized raw fixture、decoded fixture、schema 文件、最小 decoder test。不要只回報 hook log。

#### Promotion Target

已整理到 `WORKFLOW.md` 與 `DOCUMENTATION.md`。

### 2026-04-30 - Frida 只有 banner 時先做最小 hook 健康檢查

Status: validated

#### One-line Summary

Frida 沒輸出不一定是 hook 點錯，可能是 client、權限、sandbox、attach 時機或 App 沒觸發流程。

#### Human Explanation

實務上常看到 Frida 啟動後只有 banner，沒有任何自訂 log。這時不要立刻重寫完整 hook。先拆問題：Frida client 能不能列 process？最小 `console.log` hook 能不能載入？目標 PID 對不對？App 是否真的觸發該流程？執行環境是否限制了 Frida？這樣可以避免把工具連線問題誤判成分析結論。

#### Trigger

- Frida log 只有 banner。
- 完整 hook 沒有 `[INIT]` 或自訂 log。
- App 操作時 pcap 有流量，但 hook 沒任何事件。

#### Evidence

授權分析中曾遇過完整 hook 無輸出，但最小 attach 可輸出；調整執行環境或 attach/spawn 方式後，完整 hook 才正常。

#### Generalized Lesson

先用最小 hook 驗證 Frida client/server/目標 PID/權限，再檢查 hook offset 或業務流程。

#### Agent Action

遇到 Frida 無輸出時，先建議最小測試：列 process、attach PID、輸出 `HOOK_LOADED`、再逐步加 hook。不要直接擴大 hook 範圍。

#### Promotion Target

- `TOOLS.md`
- `WORKFLOW.md`

### 2026-04-30 - Session refresh 要還原 App 的真實登入/裝置流程

Status: validated

#### One-line Summary

Token 過期不一定有 refresh-token；要看 App 實際怎麼重新取得 session。

#### Human Explanation

有些 App 收到 invalid token 後，不會走標準 refresh-token API，而是清掉舊 token，回到啟動或裝置登入流程重新拿 session。若只拿舊 token 重算簽章，簽章可能正確但 token 仍無效。分析 session 問題時，要同時看 response interceptor、token store、device identity、login request builder 與 signing path。

#### Trigger

- API 回 no token、token expired、invalid token。
- 重新簽 request 仍失敗。
- App 重啟後又能成功。

#### Evidence

授權分析中曾確認：舊 session 失效後，需要按 App 的裝置登入流程取得新 token，而不是只重用舊 token 或單純重算 request signature。

#### Generalized Lesson

Session refresh 要從 App 內的 token invalidation、device identity、login body、request signing 與 token storage 一起還原。不要假設一定有 OAuth-style refresh token。

#### Agent Action

遇到 token/session 問題時，要求檢查：response interceptor 對錯誤碼的處理、token 存放位置、device id 來源、login endpoint/body、簽章 canonical path、成功後 token 寫回位置。

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`

### 2026-04-30 - 登入限流要避免 tight-loop，優先 session reuse

Status: validated

#### One-line Summary

遇到 login too frequently，不要盲目旋轉單一參數；先重用 session 並記錄風控維度。

#### Human Explanation

登入限流可能不是單一 request 欄位造成，而是伺服器用 device、User-Agent、IP、時間窗、App fingerprint、帳號狀態等多維度計算。一直換 device id 或 tight-loop login 可能讓問題更嚴重。穩定做法是同一輪測試重用 session，記錄每次登入嘗試時間與參數，必要時用 device/session pool。

#### Trigger

- API 回 login too frequently。
- 多個測試每次都重新登入。
- 改某個 device/body 欄位後結果不穩定。

#### Evidence

授權測試中曾做過參數 probe，無法把限流穩定歸因於單一欄位；session reuse 明顯降低重複登入風險。

#### Generalized Lesson

登入流程測試要有節流與重用策略。不要為每個測試方法重新登入；不要在沒有證據時假設旋轉單一 device 欄位即可繞過限制。

#### Agent Action

設計 live integration 或 runner 時，優先共用 session/context，記錄 login attempt metadata。遇到限流時先停止 tight-loop，再分析時間窗與風控維度。

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

### 2026-04-30 - 媒體播放鏈要分控制面、金鑰與資料面

Status: validated

#### One-line Summary

影片不是只有一個 URL；HLS 需要分開記錄 playlist、key、segment、解密與合併驗證。

#### Human Explanation

媒體播放流程常分成 API 詳情、HLS playlist、AES key、`.ts` segments、CDN signed URL。只抓到詳情 API 或只保存 m3u8，通常不足以證明影片可播放。可重現分析需要記錄 playlist 如何取得、key URI 如何解析、segment URL 是否短效、是否需要 AES 解密、最後是否能 remux / 播放。

#### Trigger

- App 可播放，但重放單一 URL 失敗。
- m3u8 有 key 或 signed segment query。
- segments 下載後不能直接播放。

#### Evidence

授權分析中曾確認 HLS playlist、key 與 segments 是不同請求；segments 需要依 playlist/key/IV 規則處理後才可合併成可播放檔。

#### Generalized Lesson

媒體分析文件要拆控制面與資料面，並用工具驗證 final media，而不是只看副檔名。

#### Agent Action

遇到 HLS 時，要求保存 sanitized playlist metadata、key 長度/來源、segment count、解密規則、final output 驗證結果。若要測播放，使用 `ffprobe` / `ffmpeg` 驗證容器與 frame/duration。

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

### 2026-04-30 - APK metadata：`aapt` 不在 PATH 時走 SDK build-tools；launcher 用 `resolve-activity`

Status: validated

#### One-line Summary

`apkanalyzer` 或環境找不到 build-tools 時，改用 `$ANDROID_HOME/build-tools/<version>/aapt`；badging 若沒有 launcher 行，用 `cmd package resolve-activity` 取得 `am start -n` 所需的完整 component。

#### Human Explanation

許多機器只把 `adb`（platform-tools）放在 PATH，但沒有 `aapt`。`apkanalyzer` 仍會嘗試在 SDK 內解析 `aapt`，若本機 SDK 配置不完整，可能拋出「Cannot locate latest build tools」之類錯誤。此時不必急著裝 jadx：只要已安裝 Android SDK，`build-tools` 目錄裡通常已有對應版本的 `aapt`，直接呼叫即可取得 `package`、`versionName`、`permissions`、`native-code` 等盤點資訊。

另外，`aapt dump badging` 的輸出有時**沒有** `launchable-activity:` 行（多重 activity、工具版本或 manifest 複雜度都可能造成）。這不代表無法冷啟動：在**已安裝該 package 的裝置**上，`adb shell cmd package resolve-activity --brief <package>` 常能給出預設 launcher 的 `package/class`，可用於腳本化的 `am start -n`，也比依賴 `monkey` 更穩定。

#### Trigger

- `aapt` / `aapt2`：command not found。
- `apkanalyzer`：報錯無法定位 build tools / `aapt`。
- `aapt dump badging app.apk | grep launchable-activity` 無結果，但需要自動化啟動 App。

#### Evidence

- Tool：`apkanalyzer`（依賴 SDK 內部 `aapt` 解析）。
- Tool：`aapt dump badging`（部分 APK 無 `launchable-activity` 行）。
- Tool：`adb shell cmd package resolve-activity --brief <package>`（回傳形如 `pkg/component.name`）。
- Sanitized excerpt：`IllegalStateException: Cannot locate latest build tools`（`apkanalyzer`）；`resolve-activity` 成功回傳預設 activity component。

#### Generalized Lesson

1. **Metadata fallback**：PATH 無 `aapt` 時，優先嘗試 `$ANDROID_HOME/build-tools/<任意已裝版本>/aapt`（版本目錄可用 `ls` 選最新穩定版），再考慮安裝獨立工具鏈。
2. **Launcher fallback**：需要可重現的 `am start` 時，若 badging 無 launcher，可在裝置上用 `cmd package resolve-activity --brief`；靜態-only 環境則仍可用 `aapt dump xmltree AndroidManifest.xml` 搭配 intent-filter 判讀（較費工）。

#### Agent Action

- 盤點 APK 時若 `aapt` 不在 PATH，先檢查 `$ANDROID_HOME/build-tools`，用該路徑下的 `aapt dump badging`。
- 不要僅因 `apkanalyzer` 失敗就判定「無法讀 manifest metadata」。
- 撰寫冷啟動／代理測試腳本時，優先用 `resolve-activity` 或明確的 `-n pkg/activity`；`monkey` 僅作備用。
- 技能文件與 FEEDBACK 中**不要**寫入使用者本機絕對路徑；用 `$ANDROID_HOME` 與占位符描述。

#### Applies When

- 已安裝 Android SDK（含 `build-tools`），且需要快速 badging。
- 裝置上已安裝目標 App，需要 launcher component 做自動化。

#### Does Not Apply When

- 完全沒有 Android SDK／沒有 `build-tools`（需改用人類安裝 SDK、或改用其他解析器）。
- 僅有 APK 檔、無法 adb 到裝置：launcher 需靠 xmltree／反編譯推斷。

#### Validation

- `aapt ... dump badging` 輸出含 `package:` 與 `versionName`，且與裝置 `dumpsys package` 一致（同一簽署／同一版）。
- `am start -n pkg/activity` 能啟動 App 到預設桌面入口。

#### Promotion Target

- `TOOLS.md`（命令模板、常見失敗判讀）

### 2026-04-30 - 內建 sing-box／TUN 類通道可能繞過 Wi‑Fi 系統代理

Status: candidate

#### One-line Summary

當 `libapp.so`／`classes.dex` 出現 **sing-box**、**MethodChannel `.../singbox`** 或 embedded VPN 類命名時，MITM 若「完全收不到業務 host」，不一定是 App 沒連線，而可能是流量**不經** Android Wi‑Fi HTTP 代理。

#### Human Explanation

`sing-box` 等工具常用 **TUN** 或本機轉發鏈讓流量從虛擬網卡出去，與「在系統設定裡填 HTTP proxy」是不同路徑。此時 Proxyman／Charles 可能只看到少數仍走代理的流量，或只看到 DNS／其他分流結果；全機 pcap 仍可能顯示對外 TLS。若誤判為「沒有 API 流量」，會浪費時間在 CA／pinning 上。

#### Trigger

- 靜態字串出現 `singbox`、`sing-box`、`tun`、`MethodChannel` 名稱含 `singbox`。
- Wi‑Fi 手動代理已設好，但 MITM 幾乎沒有核心 API，而裝置側仍有對外連線證據（時序、DNS、SNI、使用者可見功能正常）。

#### Evidence

- Tool：`strings` on `libapp.so` / dex；Flutter **MethodChannel** 命名（例如 `.../singbox`）。
- Sanitized excerpt：含 `singbox` channel 與 `singbox_status` 類字串；MITM 無 `CONNECT` 到 proxy 但 pcap 仍有 TLS。

#### Generalized Lesson

遇到 embedded proxy／VPN／sing-box 線索時，**代理失敗**要先拆兩層：⑴ 是否根本沒進 HTTP 代理；⑵ 進代理後 TLS 是否失敗。優先用 **whole-device pcap**（或授權環境下的路由鏡像）確認 host／SNI；再用 Frida hook **實際發請求的 client**（例如 Dart interceptor）取得明文或 signing 結果。

#### Agent Action

- 靜態掃到 `singbox` 相關 channel 時，在報告中標註「可能繞過系統 HTTP 代理」，並建議 pcap／hook 路線，不要先下結論「pinning 導致工具無效」。
- 不要在 reusable skill 中寫入特定 App 的完整通道名稱以外的業務結論；僅泛化方法。

#### Applies When

- App 內建 VPN／代理／分流（sing-box、Clash 系、自架 tun）。
- 使用者依賴 Wi‑Fi 代理 MITM 但流量走 tun。

#### Does Not Apply When

- 已確認全部 HTTPS 都對 `<proxy-host>:<proxy-port>` 發起 CONNECT（代表至少 HTTP 層有使用代理設定）。
- 純 Java `OkHttp` 且無 tun／無第三方 VPN SDK。

#### Validation

- 同一操作下：MITM 與 pcap 對照（是否有相同 SNI／時序）。
- 關閉或繞過內建通道後（僅在授權測試環境）MITM 行為是否改變。

#### Promotion Target

- `TOOLS.md`（常見失敗判讀）
- `WORKFLOW.md`（可選：代理判斷節補一句）

### 2026-04-30 - 本機 loopback「ProxyServer」轉發會讓 Wi‑Fi HTTP MITM 看不到業務 CONNECT

Status: candidate

#### One-line Summary

若流量先到 **`127.0.0.1:<local-port>`** 的 HTTP 代理，再由該代理連向真實 API host，電腦上的 Wi‑Fi HTTP MITM 通常只看到 loopback，業務 **CONNECT** 可能完全不進電腦代理。

#### Human Explanation

有些 App（含 Flutter + OkHttp 組合）會在裝置上開 **本機 proxy**，Dart／Java client 先對 loopback 發請求，再由服務端元件轉發 HTTPS。此時「系統 Wi‑Fi 手動代理」指向電腦時，**不一定**能複製到這段路徑：對 OS 而言出站連線可能仍是 **直連** API IP／走分流／經其他路由。

#### Trigger

- MITM／Proxyman 幾乎沒有業務流量。
- `adb logcat` 出現 **`ProxyServer`**／**`ProxyServerHandler`**、`Listening on port ... for target https://<api-host>`。

#### Evidence

- Tool：`adb logcat`（tag／line 依裝置而有差異）。
- Sanitized excerpt：存在「本地埠 → `https://<api-host>`」之類描述；MITM 側無對應 CONNECT。

#### Generalized Lesson

輔助抓網域時可同步搜 log：**proxy／handler／forward／localhost**。Wi‑Fi MITM 失敗時不要只猜 pinning；先確認是否存在 **loopback 中介**。

#### Agent Action

當使用者 MITM 為空但業務明顯有 HTTPS：建議 **logcat grep `ProxyServer`**／相關 tag；並提示原始 log 可能含敏感標頭，必須去敏後才寫入 repo。

#### Applies When

- App 內建本地 proxy／middleware／sing-box 類鏈路。
- OkHttp／Dart HttpClient 連 `127.0.0.1`。

#### Does Not Apply When

- 已確認業務 HTTPS **CONNECT** 目標為電腦 `<proxy-host>`。

#### Validation

- logcat 宣稱的 `<api-host>` 與 root **pcap SNI**／Frida hook URL 一致。

#### Promotion Target

- `TOOLS.md`（常見失敗判讀）
