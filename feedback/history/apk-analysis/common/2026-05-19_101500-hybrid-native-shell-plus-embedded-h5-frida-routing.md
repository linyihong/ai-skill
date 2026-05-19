> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - Hybrid Feature: Native API Shell Plus Embedded H5 Body

Status: candidate

#### One-line Summary

當功能同時有「少量原生 REST path」與「嵌入式 H5 主體」時，必須分兩條抓包／實作路線；不要把猜測的 `/feature/list` 當成與 Dio 主線同一類問題。

#### Human Explanation

Flutter App 常見「殼層原生、內容 H5」：原生只打餘額、配置、兌換等少數 path；列表、聊天、角色選擇在 `flutter_inappwebview` 內，由獨立 H5 域名承載。若 agent 只在 Dart Dio 或 CloudFront API 上 sweep `list` path，會得到大量 404，並誤判「沒有列表 API」。正確做法是：原生 path 用 live probe / Dio hook；H5 入口 URL 用 Java `WebView.loadUrl` 或 blutter 定位的「URL 指派點」；H5 內業務 XHR 再對第三方 host 單獨抓。

#### Trigger

- 反編譯只看到 2–4 條固定原生 path（如 balance、currencies、exchange），沒有 feed/list。
- Live path sweep 對 `/api/app/<feature>/list` 等候選皆 404。
- 存在 `*_h5_page.dart`、`InAppWebView`、`jumpToNext`、domain 配置鍵（如 `*H5`）等線索。
- 使用者期望「分頁列表」但原生 API 文件對不上。

#### Evidence

- Tool: blutter ASM、`pp.txt` 字串池、Java Frida `WebView.loadUrl`、可選極簡 Dart hook（URL 指派點）。
- Sanitized excerpt: 原生僅 `GET/POST /aimate/*` 三條；`WebView.loadUrl` 載入 `https://<third-party-h5-host>/?aId=<id>&tt=<JWT>&uu=<uid>&un=<name>&au=<api-base>`；`au` 常指向主站 CloudFront API 基底。
- Evidence path: 專案內 `mr.HS/docs/API列表/<feature>/`、`<PROJECT_ROOT>/exports/frida/*-h5-*.log`；lesson 不含 raw token。

#### Generalized Lesson

1. **分流**：原生殼 → Dio / Java SDK；H5 入口 → WebView 或單點 Dart hook；H5 內 API → MITM 對 H5 host 或延長 `shouldInterceptRequest`。
2. **勿混線**：Dio `RequestOptions` hook 看不到 WebView 內 XHR；Proxyman 對 CloudFront 核心 API 常不可靠，但對 H5 類流量較有機會。
3. **SDK 邊界**：Java SDK 可實作原生殼 API + **可選** `h5EntryUrl()`（組 query：`tt`、`uu`、`un`、`au`、`aId`）；**不可**假裝有原生 list endpoint，除非 live 證明。
4. **驗證順序**：先確認 H5 開頁 URL 形狀，再抓 H5 內 path；不要先寫 `/feature/list` facade。

#### Agent Action

1. 列出反編譯原生 path；對 list 候選做 live sweep，404 則停止猜 path。
2. 搜 `*h5*`、`InAppWebView`、`WebUri`、`jumpToNext`、domain map 鍵 `*H5`。
3. 跑 **極簡** Frida：`WebView.loadUrl` +（可選）blutter 定位的 H5 URL store offset；**不要**在啟動期 hook 全域 `Uri::parse` 或全量 Response/decrypt。
4. 將 H5 host、query 參數語意寫入專案 API 文件；可重用技巧留在 Ai-skill。
5. 企劃／SDK：分 Phase——原生 facade → `h5EntryUrl` → H5 API client（待 H5 流量證據）。

#### Goal / Action / Validation

- Goal: 正確切分混合功能的三層（原生殼 / H5 入口 / H5 API），避免錯誤 sweep 與 hook 導致崩潰。
- Action: 依 traffic triage 選兩條線；用 WebView 證據還原入口 URL 模板。
- Validation: 至少一條 `[JAVA_WEBVIEW_LOAD]` 或 `[H5_OPEN_URL]` 含預期 query 鍵；原生 path 與 H5 host 不同；list sweep 404 與 H5 證據可並存解釋。

#### Applies When

- Flutter + `flutter_inappwebview`（或系統 WebView）承載業務主體。
- 原生 API 明顯只是錢包／配置／兌換殼層。
- 需要還原「聊天／角色列表」但 Dio log 無 list path。

#### Does Not Apply When

- 功能完全原生（feed 在 Dio，如固定 `/aiplaza/list`）。
- 純外部瀏覽器打開，無 in-app WebView。
- H5 僅靜態展示、無後續 XHR（入口 URL 即全部）。

#### Validation

- 反編譯 path 清單與 live sweep 一致。
- Frida 極簡腳本 60–120s 內 App 不崩潰且能捕到開頁 URL（需使用者在錄製窗操作進入 H5）。
- `au`（或同類參數）與已知 API base 一致時，可支持 SDK 組 URL 設計。

#### Promotion Target

- `analysis/apk/traffic-triage.md`（混合功能小節）
- `analysis/apk/workflows/frida-hook-flow.md`（嵌入式 H5 步驟）

#### Required Linked Updates

- 見 [`2026-05-19_101500-embedded-h5-entry-url-minimal-frida-hooks.md`](../flutter-dart-aot/2026-05-19_101500-embedded-h5-entry-url-minimal-frida-hooks.md)（hook 細節）。
- 專案 plan／SDK 由 `<PROJECT_ROOT>/docs/plans/` 維護 incident 與 Phase，不寫入本 lesson 全文。
