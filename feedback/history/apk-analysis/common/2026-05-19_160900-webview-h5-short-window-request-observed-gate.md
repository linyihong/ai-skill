> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../feedback/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - WebView H5 Short-Window Capture And Request-Observed Gate

Status: candidate

#### One-line Summary

嵌入式 H5 的 API path 抓取優先用短窗口 `WebViewClient.shouldInterceptRequest` 取得 request metadata；只看到 path 時只能標 `request-observed`，不能直接進 SDK client 實作。

#### Human Explanation

WebView / H5 分析容易有兩個誤判：第一，看到沒有命中就把 capture 時間越拉越長；第二，抓到 host/path 後就把它當成 SDK-ready API。更穩定的做法是短窗口重放使用者操作，用 `WebView.loadUrl` 與實際掛上的 `WebViewClient` 子類觀察 request metadata。若需要 response body，再改走 MITM 或安全的 response hook。不要預設注入 `evaluateJavascript` / JS fetch-XHR probe，舊 WebView 或特定裝置上可能因跨 callback 保存 WebView object reference 造成 JNI object lifecycle crash。

#### Trigger

- 功能主體在 `WebView` / `flutter_inappwebview` / H5 裡，而 native / Dart API 只看到少量殼層 request。
- `WebView.loadUrl` 已能捕到 H5 入口，但 H5 內 API 尚未文件化。
- Agent 想延長 capture window、啟用重型 JS probe，或準備根據單次 path observation 做 SDK facade。

#### Evidence

- Tool: Frida Java hook around `WebView.loadUrl`, `WebView.setWebViewClient`, subclass `shouldInterceptRequest`, `onPageFinished`, and optional `WebChromeClient.onConsoleMessage`.
- Sanitized excerpt: A 60-second WebView capture can reveal H5 entry URL, H5 API host, CORS preflight, opaque `data` query key, and static resource requests without recording token values or user identifiers.
- Evidence path: project-specific raw logs and exact host/path evidence stay in `<PROJECT_ROOT>` docs / gitignored evidence; this lesson only keeps generalized method and gates.

#### Generalized Lesson

1. **Keep capture windows short first**: use 30-60 seconds with a concrete operation recipe before increasing duration. Longer windows add noise and do not fix a wrong UI path or unstable hook.
2. **Hook the actual WebViewClient class**: hook `WebView.setWebViewClient` to discover and instrument the concrete subclass, not only `android.webkit.WebViewClient`.
3. **Prefer request metadata before JS injection**: `shouldInterceptRequest` gives method, URL, main-frame / gesture flags, and safe header names. This is often enough to establish candidate endpoints.
4. **Disable JS probe by default**: `evaluateJavascript` / JS `fetch` / `XMLHttpRequest` probe can be useful, but should be opt-in short-window experimentation because WebView object references may be unsafe across callbacks.
5. **Use a request-observed gate**: path + query keys + header names prove only request observation. SDK implementation needs response wrapper/schema, opaque parameter generation, auth semantics, and error behavior.

#### Agent Action

1. Before running capture, define a 60-second operation recipe: enter H5, trigger one list/detail/chat action, scroll once.
2. Run a light Frida script: `WebView.loadUrl`, `setWebViewClient`, `shouldInterceptRequest`, and `onPageFinished`; keep JS probe off unless explicitly testing it.
3. After capture, summarize only sanitized host family, path shape, method, query key names, header names, and trigger confidence.
4. Mark docs as `request-observed` until response body/status and opaque parameter generation are known.
5. If SDK work is requested, block on response/schema and generation rules, or explicitly scope SDK to H5 entry URL / browser automation only.

#### Goal / Action / Validation

- Goal: Prevent unstable WebView H5 captures and premature SDK client work from path-only evidence.
- Action: Use short-window WebView request metadata capture, document candidate endpoints, and apply a request-observed readiness gate.
- Validation or reference source: A valid capture has `WebView.loadUrl` and at least one `shouldInterceptRequest` record for the H5 host; SDK readiness remains false until response schema and opaque query generation are documented.

#### Applies When

- Embedded H5 / WebView carries feature content and native APIs are only a shell.
- The goal is API discovery or SDK readiness assessment.
- Request metadata can be captured without raw response bodies.

#### Does Not Apply When

- The feature is fully native and already covered by Dart / Java HTTP hooks.
- The only needed artifact is the H5 entry URL, not H5 internal APIs.
- The team has an approved MITM setup that already captures request and response safely.

#### Validation

- Script syntax passes before capture.
- Capture completes without app crash under the default light hook profile.
- Tracked docs contain no raw token, cookie, user identifier, full replay URL, or private payload fragment.
- API docs clearly distinguish `request-observed` from `schema-validated` or `sdk-ready`.

#### Promotion Target

- `analysis/apk/workflows/frida-hook-flow.md` for WebView short-window hook strategy.
- `workflow/apk-analysis/execution-flow.md` for short-window replay and request-observed gates if repeated across more H5 captures.
- `workflow/software-delivery/execution-flow.md` only if the SDK-readiness gate needs a general development handoff rule.

#### Required Linked Updates

- Project evidence stays in `<PROJECT_ROOT>` API docs / gitignored capture logs; this lesson intentionally omits concrete hosts, endpoints, tokens, user ids, and payload values.
- `feedback/history/apk-analysis/README.md` category count should be updated with this lesson.
- Existing per-round feedback checkpoint lesson should be checked because this lesson was initially identified before being written; that gap is a feedback-close-loop failure, not a new APK technique.
