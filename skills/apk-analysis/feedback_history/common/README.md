# Common Feedback Lessons

Cross-cutting APK analysis lessons that apply before a runtime/API family is known, or that affect the shared workflow across categories.

| File | Status | Topic | Notes |
| --- | --- | --- | --- |
| `../2026-04-30_120000-proxy-failure-要先拆成導流與-tls-兩層.md` | promoted | Proxy failure 要先拆成導流與 TLS 兩層 | Common MITM/proxy diagnosis. |
| `../2026-04-30_120001-冷啟動比-attach-更適合抓初始化與代理設定.md` | promoted | 冷啟動比 attach 更適合抓初始化與代理設定 | Common startup/injection timing. |
| `../2026-04-30_120002-高語意-hook-優先於-socket-bytes.md` | promoted | 高語意 hook 優先於 socket bytes | Common hook selection principle. |
| `../2026-04-30_120003-動態-hook-結果要離線化.md` | promoted | 動態 hook 結果要離線化 | Common fixture/offline validation rule. |
| `../2026-04-30_120004-frida-只有-banner-時先做最小-hook-健康檢查.md` | validated | Frida 只有 banner 時先做最小 hook 健康檢查 | Common Frida troubleshooting. |
| `../2026-04-30_120005-session-refresh-要還原-app-的真實登入-裝置流程.md` | validated | Session refresh 要還原 App 的真實登入/裝置流程 | Common session/token analysis. |
| `../2026-04-30_120006-登入限流要避免-tight-loop-優先-session-reuse.md` | validated | 登入限流要避免 tight-loop，優先 session reuse | Common live-test safety. |
| `../2026-04-30_120008-aapt-sdk-build-tools-resolve-activity.md` | validated | APK metadata：`aapt` 不在 PATH 時走 SDK build-tools | Common APK metadata/setup. |
| `../2026-05-01_101500-doh-dns-query-param-side-channel-from-okhttp-log.md` | validated | DoH 的 `dns=` 參數可作為 MITM 業務 host 空白時的側信道 | Common network side-channel triage. |
| `../2026-05-01_112900-proxy-config-vs-business-route.md` | validated | Proxy config is not business route proof | Common proxy route proof rule. |
| `../2026-05-01_174100-technique-category-routing.md` | promoted | Technique category routing | Repo organization/category routing rule. |
