> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../feedback/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-20 - WebView MITM After Direct Load

Status: candidate

#### One-line Summary

當 Android WebView H5 首頁一開 proxy 就空白，但後續 XHR host 才是 schema 目標時，先直連載入前端，再開 MITM + 短窗 SSL proceed 抓後續 API。

#### Human Explanation

WebView H5 分析不一定要把首頁、靜態資源與業務 API 都放進同一個 MITM window。某些 H5 首頁在代理或自簽 CA 環境下會白屏，只剩退出按鈕；但使用者真正需要的 schema 可能是前端載入完成後，滑動列表或切換 tab 時才發出的 XHR。此時持續調代理設定或反覆重載首頁，會把時間花在非目標 host 上，還可能讓 capture 狀態變髒。

更穩定的切法是分階段：先關 Android global proxy，讓 H5 shell 和靜態資源在原始 TLS 條件下正常載入；確認目標列表或頁面已渲染後，再開 proxy / MITM 和 opt-in WebView SSL proceed，只用安全 UI 動作觸發後續 API。這樣可以避免首頁 MITM 白屏，同時仍取得目標 XHR 的 response wrapper 與 payload shape。

#### Trigger

- WebView H5 首頁在 MITM proxy 下空白，或只顯示退出 / loading。
- 目標 API 是 H5 載入後的 XHR，而不是首頁 HTML / JS。
- 已有 request metadata，但需要 runtime response schema。
- 測試範圍允許短窗 Android proxy 與 WebView SSL proceed。

#### Evidence

- Tool: Android global proxy, mitmdump, Frida WebView `onReceivedSslError` opt-in proceed, `shouldInterceptRequest`.
- Sanitized excerpt: 先開 proxy 時 H5 shell 只顯示空白；先直連載入列表，再開 proxy 後，列表滑動 / tab 切換觸發目標 API，MITM 成功取得 HTTP 200 JSON wrapper 與 decoded payload shape。
- Evidence path: 專案 raw MITM、完整 URL、token、使用者資料與 response body 留在 gitignored evidence；本 lesson 只保留 generalized capture order。

#### Generalized Lesson

1. **拆開 shell load 與 API capture**：H5 首頁 / bundle 載入可以直連，目標 XHR 再進 MITM window。
2. **先確認 UI 到達目標狀態**：只有列表或目標頁已正常渲染後，才開代理觸發後續 API。
3. **用最小 UI 動作觸發**：只做滑動、分頁、tab 切換等非消費動作；不要因抓 schema 而點進付費或 destructive flow。
4. **SSL proceed 仍是短窗 opt-in**：只對授權 response/schema capture 開啟，結束後清 Android global proxy。
5. **白屏不是 schema 失敗**：首頁白屏只代表 proxy 時序或 TLS 條件不適合載入 shell，不等於目標 XHR 不可抓。

#### Agent Action

1. 若 H5 在 proxy 下白屏，先清 Android global proxy，重新載入 H5 shell。
2. 等目標列表 / 頁面可見後，啟動 MITM、設定 Android global proxy，並 attach Frida SSL proceed hook。
3. 只做 scope 內安全 UI 操作觸發 XHR。
4. 抽取 wrapper keys / payload field types；不要把 token、完整 URL、圖片 URL 或原始文字值寫入 tracked docs。
5. capture 結束後清 proxy，並把「先直連載入、再代理 XHR」寫回專案 checklist。

#### Goal / Action / Validation

- Goal: 避免 WebView 首頁 MITM 白屏阻擋後續業務 API schema capture。
- Action: 將 H5 shell loading 與 XHR response capture 分成兩個時序，並把代理窗口縮到目標 API 觸發階段。
- Validation or reference source: MITM 檔中有目標 XHR 的 HTTP 200 wrapper；decoded payload 只記欄位 / 型別；Android proxy 已恢復。

#### Applies When

- Android WebView / Flutter InAppWebView 內嵌第三方 H5。
- H5 shell 由靜態站提供，業務 API 由另一 host 提供。
- 目標 API 可由非破壞性 UI 動作重新觸發。

#### Does Not Apply When

- 目標 response 只在首頁初次載入瞬間出現，且無法重觸發。
- 需要分析的正是首頁 HTML / bundle response。
- 測試範圍不允許代理切換或 TLS 行為 hook。

#### Validation

- 專案 checklist 記錄兩階段 capture order。
- raw capture 留在 gitignored evidence，tracked docs 只保留 schema。
- 結束後 Android global proxy 為清除狀態。

#### Promotion Target

- `analysis/apk/workflows/frida-hook-flow.md` for WebView response-capture sequencing.
- `workflow/apk-analysis/execution-flow.md` for adding "direct-load then MITM XHR" as a downgrade path after proxy-induced white screen.

#### Required Linked Updates

- Project H5 capture checklist should mention whether shell must be loaded without proxy before API MITM.
- Project API docs should distinguish runtime response schema from SDK external replay readiness.
- `feedback/history/apk-analysis/README.md` category count should be updated with this lesson.
