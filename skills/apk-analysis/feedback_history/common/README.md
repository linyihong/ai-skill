# Common Feedback Lessons

Cross-cutting APK analysis lessons that apply before a runtime/API family is known, or that affect the shared workflow across categories.

| File | Status | Topic | Notes |
| --- | --- | --- | --- |
| `2026-04-30_120000-proxy-failure-要先拆成導流與-tls-兩層.md` | promoted | Proxy failure 要先拆成導流與 TLS 兩層 | Common MITM/proxy diagnosis. |
| `2026-04-30_120001-冷啟動比-attach-更適合抓初始化與代理設定.md` | promoted | 冷啟動比 attach 更適合抓初始化與代理設定 | Common startup/injection timing. |
| `2026-04-30_120002-高語意-hook-優先於-socket-bytes.md` | promoted | 高語意 hook 優先於 socket bytes | Common hook selection principle. |
| `2026-04-30_120003-動態-hook-結果要離線化.md` | promoted | 動態 hook 結果要離線化 | Common fixture/offline validation rule. |
| `2026-04-30_120004-frida-只有-banner-時先做最小-hook-健康檢查.md` | validated | Frida 只有 banner 時先做最小 hook 健康檢查 | Common Frida troubleshooting. |
| `2026-04-30_120005-session-refresh-要還原-app-的真實登入-裝置流程.md` | validated | Session refresh 要還原 App 的真實登入/裝置流程 | Common session/token analysis. |
| `2026-04-30_120006-登入限流要避免-tight-loop-優先-session-reuse.md` | validated | 登入限流要避免 tight-loop，優先 session reuse | Common live-test safety. |
| `2026-04-30_120008-aapt-sdk-build-tools-resolve-activity.md` | validated | APK metadata：`aapt` 不在 PATH 時走 SDK build-tools | Common APK metadata/setup. |
| `2026-05-01_101500-doh-dns-query-param-side-channel-from-okhttp-log.md` | validated | DoH 的 `dns=` 參數可作為 MITM 業務 host 空白時的側信道 | Common network side-channel triage. |
| `2026-05-01_112900-proxy-config-vs-business-route.md` | validated | Proxy config is not business route proof | Common proxy route proof rule. |
| `2026-05-01_174100-technique-category-routing.md` | promoted | Technique category routing | Repo organization/category routing rule. |
| `2026-05-01_175500-workflow-as-routing-not-technique-dump.md` | promoted | Workflow as routing, not technique dump | Keep `WORKFLOW.md` common and move category details into `techniques/<category>/`. |
| `2026-05-06_090100-frida-17-jailed-android-gadget-gate.md` | validated | Frida 17 / jailed Android Gadget gate | Distinguish CLI compatibility and injection transport failures from hook correctness. |
| `2026-05-06_092200-frida-server-version-alignment.md` | validated | Frida server version alignment before attach debugging | If `frida-ps` works but attach closes globally, align local/device server versions before changing hooks. |
| `2026-05-06_094200-frida-heavy-window-ui-control-comparison.md` | validated | Frida-heavy window UI control comparison | When Frida-heavy captures trigger retry/spinner states, compare no-Frida and attach-after-load controls before blaming backend timeout. |
| `2026-05-06_101000-in-app-search-suggestion-before-ime-submit.md` | validated | In-app search suggestions before IME submit | Use app-owned suggestion chips/result tabs to validate search flows before typed/IME submit that may leave the app. |
| `2026-05-06_102100-feature-handoff-finish-gate.md` | promoted | Feature handoff finish gate | Named-feature APK analysis is not complete until project-level feature handoff/architecture docs are created or updated. |
| `2026-05-06_102800-ui-hierarchy-content-desc-redaction.md` | validated | UI hierarchy content-desc redaction | Treat UIAutomator `text` / `content-desc` as potentially raw content; document structure and evidence paths instead of values. |
| `2026-05-06_104000-detailed-api-list-finish-gate.md` | promoted | Detailed API list finish gate | Confirmed API flows need per-API request/response docs, not only schema or correlation summaries. |
| `2026-05-06_105026-state-reset-baseline-feature-capture.md` | promoted | State reset baseline before feature capture | App-start-to-feature API analysis should record reset/cache/session baseline so startup, navigation, and feature APIs are not conflated. |
| `2026-05-06_110749-scrollable-tab-strip-coverage.md` | promoted | Scrollable tab strip coverage | Tab/category/result strips need horizontal swipe coverage before claiming complete feature API coverage. |
| `2026-05-06_111558-foreground-package-validation.md` | promoted | Foreground package validation | UI evidence must prove the foreground package is the target app before API/UI alignment claims. |
| `2026-05-06_113302-dart-aot-offset-from-asm-address.md` | promoted | Dart AOT offset from ASM address | Use the first asm-line function address, not the output filename suffix, as the Frida hook offset. |
| `2026-05-06_122900-checkpoint-replay-runner.md` | promoted | Checkpoint replay runner | Stable UI/API recapture should use checkpointed replay scripts with package validation and target stops. |
| `2026-05-06_130300-feature-context-validation.md` | promoted | Feature context validation | Package-correct UI evidence can still be the wrong in-app screen; add target feature anchors to checkpoint guards. |
| `2026-05-06_133600-post-reset-window-split.md` | promoted | Post-reset window split | Split reset/session recovery and feature-checkpoint attach windows when launch-time hooks or first-run UI can pollute attribution. |
| `2026-05-07_124100-ui-fast-path-bounded-scroll.md` | candidate | UI fast path and bounded scroll | Compare app-owned shortcut controls against generic tap paths, and use bounded scroll convergence when scroll extent is unavailable. |
| `2026-05-07_131000-highest-leverage-analysis-path.md` | candidate | Highest leverage analysis path | Choose the analysis route with the best time-to-evidence ratio instead of continuing a fixed technique sequence. |
