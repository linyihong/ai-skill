> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - Embedded H5 Entry URL: Minimal Frida Hooks (Not Uri::parse Flood)

Status: candidate

#### One-line Summary

抓取 Flutter 嵌入式 H5 **開頁 URL** 時，優先 `Java WebView.loadUrl` + blutter 單點「URL 指派」hook；禁用啟動期全域 `Uri::parse` 與全量 Dio Response/decrypt hook，否則易 stack overflow 崩潰。

#### Human Explanation

`flutter_inappwebview` 最終仍會走到 Android `WebView.loadUrl`，在 Java 層可直接看到完整 HTTPS URL（含 token query）。Dart 層可在 blutter 偽源碼找到 `jumpToNext` → `AiH5Page` 建構後、將 URL 寫入 page 欄位的那一條 `StoreField`（單次、低頻）。若同時掛 `Uri::parse` 與 `getTaggedObjectValue` 於每次 parse，啟動期會觸發上萬次而崩潰。實務上 **H5 minimal 模式**：只載入 libapp 的單點 hook + Java WebView；**不** attach 全量 `blutter_frida.js` Dio 鏈。

#### Trigger

- 需要 H5 入口 URL，但 Dio log 只有殼層 path。
- 全量 blutter Frida（Dio + Response + `aesDecryptEx`）導致 `DartWorker` SIGSEGV / stack overflow。
- `Uri::parse` hook 後 log 充滿 `Unhandle class id` 隨即 crash。

#### Evidence

- Tool: Frida 17.x、`blutter` offset、`android.webkit.WebView.loadUrl`。
- Sanitized excerpt: `__FRIDA_H5_MINIMAL__` 僅 hook URL assign offset + Java WebView；90s 內 `[JAVA_WEBVIEW_LOAD] https://<h5-host>/?aId=...&tt=...&uu=...&un=...&au=...`；全量 hook 約 15s crash。
- Evidence path: `<PROJECT_ROOT>/scripts/frida-capture-*-h5-phone.sh`、`exports/frida/mrhs-*-h5-*.log`。

#### Generalized Lesson

| 做法 | 風險 | 產出 |
| --- | --- | --- |
| Java `WebView.loadUrl` + filter（`http`、`tt=`、feature 關鍵字） | 低 | 完整開頁 URL |
| blutter 單點 URL store（`jumpToNext` 後 Assign 至 H5 page） | 低～中 | 與 Java 互證 |
| 全域 `Uri::parse` + object decode | **高（崩潰）** | 不建議 |
| 全量 Dio + Response + decrypt | **高（崩潰/ANR）** | 僅用於 CloudFront 核心 API，非 H5 入口 |

**Attach 時機**：優先 attach **已啟動且過 splash** 的 PID；避免 force-stop 後立即全量 hook。

**錄製窗**：腳本 sleep 60–120s 內，使用者必須手動進入 H5（agent 無法代點）。

#### Agent Action

1. 從 blutter ASM 找 `*h5_page*`、`jumpToNext`、`_interpolate`、`replaceAll("/api/app")` 理解 query 語意（靜態）。
2. 部署 minimal config：`__FRIDA_H5_MINIMAL__` + `WebView.loadUrl`；腳本輸出 `exports/frida/*-h5-*.log`。
3. 錄製期間提示使用者點「進入聊天／H5」；結束後 `grep WEBVIEW_LOAD|H5_OPEN_URL`。
4. H5 內列表 API：第二輪再加 `shouldInterceptRequest` 或對 `<h5-host>` MITM，**不要**回到 Uri::parse flood。

#### Goal / Action / Validation

- Goal: 穩定取得嵌入式 H5 開頁 URL，不犧牲 App 穩定性。
- Action: minimal hook set + 已啟動 attach + 使用者操作窗。
- Validation: log 含至少一條完整 HTTPS URL；App 未在 hook 後 30s 內崩潰；URL query 含 session token 與 api-base 參數（名稱依 App 而異）。

#### Applies When

- `flutter_inappwebview` 或 `WebView` 載入第三方 H5。
- blutter 已產出 `libapp.so` offset。
- 入口 URL 由 Dart 組裝後傳入 WebView。

#### Does Not Apply When

- H5 在系統瀏覽器開啟（非 in-app WebView）。
- 需 H5 **內** XHR 明文（本 lesson 只覆蓋開頁 URL；內部 API 需 MITM 或 intercept 第二階段）。
- 無 blutter offset 且 Java hook 也無輸出（可能非 WebView 通道）。

#### Validation

- 對照 blutter 中 query 片段字串（如 `&tt=`、`&uu=`）與捕獲 URL 一致。
- Java-only 與 Dart 單點 hook 若皆命中，URL 應一致。
- 全量 hook 對照實驗：確認 minimal 模式不 crash（可選）。

#### Promotion Target

- `analysis/apk/workflows/frida-hook-flow.md`
- 交叉引用 [`../common/2026-05-19_101500-hybrid-native-shell-plus-embedded-h5-frida-routing.md`](../common/2026-05-19_101500-hybrid-native-shell-plus-embedded-h5-frida-routing.md)

#### Required Linked Updates

- `traffic-triage.md` 混合功能小節（由 common lesson 觸發）。
