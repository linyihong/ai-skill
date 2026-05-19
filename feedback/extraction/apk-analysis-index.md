# APK Analysis Feedback History Index

本索引列出舊 `skills/apk-analysis/feedback_history/`（已刪除）下所有 lessons，依其 `Promotion Target` 分類到對應的目標層。此索引讓 `feedback/` 層可發現 lessons，並追蹤哪些 lessons 已被提取到目標層。

> ⚠️ **Lesson 已搬遷**：所有 lesson 實體已移至 [`feedback/history/apk-analysis/`](../history/apk-analysis/README.md)。本索引保留作為 extraction 追蹤用途，來源路徑 `skills/apk-analysis/feedback_history/` 已刪除。
>
> 搬遷日期：2026-05-13

## 索引說明

| 欄位 | 說明 |
| --- | --- |
| **目標層** | lesson 的 Promotion Target 對應到的新架構層 |
| **來源檔案** | 舊 `skills/apk-analysis/feedback_history/`（已刪除）下的原始 lesson 檔案 |
| **Status** | lesson 的成熟度狀態（promoted / validated / candidate） |
| **提取狀態** | 是否已提取到目標層（✅ 已提取 / ⬜ 未提取 / 🔄 部分提取） |
| **目標檔案** | 提取後的目標檔案路徑 |

---

## 1. 目標層：`workflow/apk-analysis/`

這類 lessons 的 Promotion Target 包含 `WORKFLOW.md`，適合提取到 `workflow/apk-analysis/`。

### common/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-04-30_120000-proxy-failure-要先拆成導流與-tls-兩層.md` | promoted | Proxy failure 要先拆成導流與 TLS 兩層 | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-04-30_120001-冷啟動比-attach-更適合抓初始化與代理設定.md` | promoted | 冷啟動比 attach 更適合抓初始化與代理設定 | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-04-30_120002-高語意-hook-優先於-socket-bytes.md` | promoted | 高語意 hook 優先於 socket bytes | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-04-30_120003-動態-hook-結果要離線化.md` | promoted | 動態 hook 結果要離線化 | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-04-30_120004-frida-只有-banner-時先做最小-hook-健康檢查.md` | validated | Frida 只有 banner 時先做最小 hook 健康檢查 | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-04-30_120005-session-refresh-要還原-app-的真實登入-裝置流程.md` | validated | Session refresh 要還原 App 的真實登入/裝置流程 | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-04-30_120006-登入限流要避免-tight-loop-優先-session-reuse.md` | validated | 登入限流要避免 tight-loop，優先 session reuse | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-01_101500-doh-dns-query-param-side-channel-from-okhttp-log.md` | validated | DoH 的 `dns=` 參數可作為 MITM 業務 host 空白時的側信道 | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-01_112900-proxy-config-vs-business-route.md` | validated | Proxy config is not business route proof | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_090100-frida-17-jailed-android-gadget-gate.md` | validated | Frida 17 / jailed Android Gadget gate | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_092200-frida-server-version-alignment.md` | validated | Frida server version alignment before attach debugging | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_094200-frida-heavy-window-ui-control-comparison.md` | validated | Frida-heavy window UI control comparison | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_101000-in-app-search-suggestion-before-ime-submit.md` | validated | In-app search suggestions before IME submit | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_105026-state-reset-baseline-feature-capture.md` | promoted | State reset baseline before feature capture | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_110749-scrollable-tab-strip-coverage.md` | promoted | Scrollable tab strip coverage | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_111558-foreground-package-validation.md` | promoted | Foreground package validation | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_122900-checkpoint-replay-runner.md` | promoted | Checkpoint replay runner | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_130300-feature-context-validation.md` | promoted | Feature context validation | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_133600-post-reset-window-split.md` | promoted | Post-reset window split | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-07_124100-ui-fast-path-bounded-scroll.md` | candidate | UI fast path and bounded scroll | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-07_131000-highest-leverage-analysis-path.md` | candidate | Highest leverage analysis path | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-07_143800-stateful-setting-probe-restore.md` | candidate | Stateful setting probe restore | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-07_144300-ui-route-backfill-after-runtime-probe.md` | candidate | UI route backfill after runtime probe | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-07_145500-sensitive-provider-fingerprint-diagnostic.md` | candidate | Sensitive provider fingerprint diagnostic | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-07_161400-sdk-live-self-generation-audit.md` | candidate | SDK live self-generation audit | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-07_171900-provider-read-negative-matrix.md` | candidate | Provider read negative matrix | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-11_125615-per-round-feedback-checkpoint.md` | promoted | Per-round feedback checkpoint | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-08_135700-agent-goals-before-transcripts.md` | candidate | Agent goals before transcripts when resuming | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-08_141430-identity-material-self-generation-audit.md` | candidate | Identity material self-generation audit | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-19_135900-android-cache-reset-before-image-hook-capture.md` | candidate | Android cache reset before image hook capture | ⬜ | `workflow/apk-analysis/execution-flow.md`, `analysis/apk/workflows/frida-hook-flow.md` |

### flutter-dart-aot/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `flutter-dart-aot/2026-05-01_133900-dart-aot-interceptor-strings-after-java-helper-miss.md` | validated | Dart AOT interceptor strings after Java helper miss | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_140900-unflutter-aot-offset-hook-after-blutter-crash.md` | promoted | Unflutter AOT offset hook after blutter crash | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_142000-exhaustive-java-okhttp-hooks-may-still-miss-flutter-business-http.md` | validated | 廣覆蓋 Java OkHttp 仍無業務 host 時應轉 Dart／native／pcap | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_145948-dart-aot-callsite-hooks-are-not-function-hooks.md` | validated | Dart AOT Callsite Hooks Are Not Function Hooks | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_151551-schema-only-jsondecode-hook.md` | validated | Schema-only jsonDecode Hook | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_153437-sequence-jsondecode-before-api-response.md` | validated | Sequence jsonDecode Before Calling It API Response | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_155200-dart-compressed-response-fields.md` | candidate | Dart AOT compressed response fields | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_164741-dart-inline-onebyte-string-smi-length.md` | validated | Dart inline one-byte string Smi length | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-07_135100-entry-boundary-before-unstable-aot-hooks.md` | candidate | Entry boundary before unstable AOT hooks | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-07_135900-dio-requestoptions-uri-keyset-probe.md` | candidate | Dio RequestOptions URI keyset probe | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-07_141200-dart-aot-lazy-static-slot-trace.md` | candidate | Dart AOT lazy static slot trace | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-07_142600-dart-aot-async-provider-return-shape.md` | candidate | Dart AOT async provider return shape | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-07_151200-aot-hook-crash-static-boundary-fallback.md` | candidate | AOT hook crash static boundary fallback | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |

### http-api/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `http-api/2026-05-01_171500-json-shape-before-query-shape.md` | validated | JSON Shape Before Query Shape | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md` | validated | UI Architecture Map from screenshots | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-01_173800-api-field-documentation-after-analysis.md` | validated | API field documentation after analysis | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md` | validated | UI automation operation scripts for API capture | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_183700-scrollable-clickable-screen-mapping.md` | validated | Scrollable clickable screen mapping | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_184218-playable-list-item-still-needs-detail-api.md` | validated | Playable list item still needs detail API | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_184700-screen-reachability-operation-recipes.md` | validated | Screen reachability operation recipes | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_184900-in-app-route-map-external-transitions.md` | validated | In-app route map external transitions | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_185921-scroll-depth-before-api-conclusion.md` | validated | Scroll depth before API conclusion | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_190838-richtext-html-images-are-api-resources.md` | validated | Richtext HTML images are API resources | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_192012-infinite-scroll-needs-pagination-proof.md` | validated | Infinite scroll needs pagination proof | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_195200-feature-reconstruction-handoff.md` | validated | Feature reconstruction handoff | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_202400-page-analysis-requires-ui-map-artifact.md` | validated | Page analysis requires UI map artifact | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_205500-ui-behavior-backfill-finish-gate.md` | validated | UI behavior backfill finish gate | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-06_081400-auto-app-guidance-for-sdk-tools.md` | validated | Auto app guidance for SDK tools | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-06_104300-api-catalog-finish-gate.md` | promoted | API Catalog finish gate | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-07_125600-api-first-pagination-replay.md` | candidate | API-first pagination replay | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-11_095600-ui-api-top-items-parity.md` | candidate | UI/API top items parity before API blame | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-11_113100-post-selection-gesture-for-lazy-api.md` | candidate | Post-selection gesture for lazy API | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-11_113200-read-only-argument-override-preserve-app-path.md` | candidate | Read-only argument override preserve app path | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-11_135000-redacted-sample-targeting-classifier.md` | candidate | Redacted sample targeting classifier | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-11_135700-articles-first-live-adapter-smoke.md` | candidate | Articles-first live adapter smoke | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `dynamic-capture/2026-05-19_143500-avoid-single-cached-target-for-decrypt-capture.md` | candidate | Avoid single cached target for decrypt capture | ⬜ | `workflow/apk-analysis/execution-flow.md`, `analysis/apk/workflows/frida-hook-flow.md` |

### local-proxy/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `local-proxy/2026-04-30_120009-內建-sing-box-tun-類通道可能繞過-wi-fi-系統代理.md` | candidate | 內建 sing-box／TUN 類通道可能繞過 Wi‑Fi 系統代理 | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |
| `local-proxy/2026-04-30_120010-本機-loopback-proxyserver-轉發會讓-wi-fi-http-mitm-看不到業務-connect.md` | candidate | 本機 loopback「ProxyServer」轉發會讓 Wi‑Fi HTTP MITM 看不到業務 CONNECT | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |
| `local-proxy/2026-05-01_114300-local-proxy-handler-uri-hook.md` | validated | Hook local proxy handler URI, not just OkHttp | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |
| `local-proxy/2026-05-01_131000-cast-netty-request-for-handler-route.md` | validated | Cast Netty request interfaces for handler routes | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |
| `local-proxy/2026-05-01_132400-netty-aggregated-request-tostring-headers.md` | validated | Netty aggregated request toString can expose headers | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |

### media-hls/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `media-hls/2026-04-30_120007-媒體播放鏈要分控制面-金鑰與資料面.md` | validated | 媒體播放鏈要分控制面、金鑰與資料面 | ✅ | `analysis/apk/workflows/media-hls-analysis-flow.md` |
| `media-hls/2026-05-06_121650-hls-playlist-key-may-be-wrapped.md` | promoted | HLS playlist key may be wrapped | ✅ | `analysis/apk/workflows/media-hls-analysis-flow.md` |
| `media-hls/2026-05-11_093300-media-prefix-matrix-cdn-stale.md` | candidate | Media prefix matrix and stale CDN classification | ✅ | `analysis/apk/workflows/media-hls-analysis-flow.md` |

---

## 2. 目標層：`intelligence/engineering/analytical-reasoning/`

這類 lessons 包含可泛化的工程判斷、啟發式或信號檢測，適合提取到 `intelligence/engineering/analytical-reasoning/`。

### 適合 heuristics/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-04-30_120002-高語意-hook-優先於-socket-bytes.md` | promoted | 高語意 hook 優先於 socket bytes | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md` |
| `common/2026-05-07_131000-highest-leverage-analysis-path.md` | candidate | Highest leverage analysis path | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/` |
| `flutter-dart-aot/2026-05-01_145948-dart-aot-callsite-hooks-are-not-function-hooks.md` | validated | Dart AOT Callsite Hooks Are Not Function Hooks | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/` |
| `flutter-dart-aot/2026-05-07_135100-entry-boundary-before-unstable-aot-hooks.md` | candidate | Entry boundary before unstable AOT hooks | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/` |
| `http-api/2026-05-01_171650-ui-architecture-map-from-screenshots.md` | validated | UI Architecture Map from screenshots | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-operation-stability.md` |
| `http-api/2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md` | validated | UI automation operation scripts for API capture | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-operation-stability.md` |
| `http-api/2026-05-05_183700-scrollable-clickable-screen-mapping.md` | validated | Scrollable clickable screen mapping | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-operation-stability.md` |
| `common/2026-05-07_124100-ui-fast-path-bounded-scroll.md` | candidate | UI fast path and bounded scroll | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-operation-stability.md` |
| `http-api/2026-05-05_184700-screen-reachability-operation-recipes.md` | validated | Screen reachability operation recipes | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-to-api-attribution.md` |
| `common/2026-05-07_144300-ui-route-backfill-after-runtime-probe.md` | candidate | UI route backfill after runtime probe | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-to-api-attribution.md` |
| `common/2026-05-06_111558-foreground-package-validation.md` | promoted | Foreground package validation | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-to-api-attribution.md` |
| `common/2026-05-06_130300-feature-context-validation.md` | promoted | Feature context validation | ✅ | `intelligence/engineering/analytical-reasoning/heuristics/ui-to-api-attribution.md` |

### 適合 anti-patterns/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-04-30_120000-proxy-failure-要先拆成導流與-tls-兩層.md` | promoted | Proxy failure 要先拆成導流與 TLS 兩層 | ✅ | `intelligence/engineering/analytical-reasoning/anti-patterns/` |
| `common/2026-05-01_112900-proxy-config-vs-business-route.md` | validated | Proxy config is not business route proof | ✅ | `intelligence/engineering/analytical-reasoning/anti-patterns/` |

### 適合 failure/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-04-30_120004-frida-只有-banner-時先做最小-hook-健康檢查.md` | validated | Frida 只有 banner 時先做最小 hook 健康檢查 | ✅ | `intelligence/engineering/analytical-reasoning/failure/` |
| `common/2026-05-06_090100-frida-17-jailed-android-gadget-gate.md` | validated | Frida 17 / jailed Android Gadget gate | ✅ | `intelligence/engineering/analytical-reasoning/failure/` |
| `common/2026-05-06_092200-frida-server-version-alignment.md` | validated | Frida server version alignment before attach debugging | ✅ | `intelligence/engineering/analytical-reasoning/failure/` |

### 適合 signals/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-01_101500-doh-dns-query-param-side-channel-from-okhttp-log.md` | validated | DoH 的 `dns=` 參數可作為 MITM 業務 host 空白時的側信道 | ✅ | `intelligence/engineering/analytical-reasoning/signals/` |
| `local-proxy/2026-04-30_120009-內建-sing-box-tun-類通道可能繞過-wi-fi-系統代理.md` | candidate | 內建 sing-box／TUN 類通道可能繞過 Wi‑Fi 系統代理 | ✅ | `intelligence/engineering/analytical-reasoning/signals/local-proxy-detection.md` |
| `local-proxy/2026-04-30_120010-本機-loopback-proxyserver-轉發會讓-wi-fi-http-mitm-看不到業務-connect.md` | candidate | 本機 loopback「ProxyServer」轉發會讓 Wi‑Fi HTTP MITM 看不到業務 CONNECT | ✅ | `intelligence/engineering/analytical-reasoning/signals/local-proxy-detection.md` |

---

## 3. 目標層：`analysis/apk/workflows/`

這類 lessons 包含具體的操作流程、步驟或命令，適合提取到 `analysis/apk/workflows/`。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-04-30_120008-aapt-sdk-build-tools-resolve-activity.md` | validated | APK metadata：`aapt` 不在 PATH 時走 SDK build-tools | ✅ | `analysis/apk/tools-and-failures.md` |
| `common/2026-05-06_113302-dart-aot-offset-from-asm-address.md` | promoted | Dart AOT offset from ASM address | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_140900-unflutter-aot-offset-hook-after-blutter-crash.md` | promoted | Unflutter AOT offset hook after blutter crash | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_151551-schema-only-jsondecode-hook.md` | validated | Schema-only jsonDecode Hook | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `flutter-dart-aot/2026-05-01_164741-dart-inline-onebyte-string-smi-length.md` | validated | Dart inline one-byte string Smi length | ✅ | `analysis/apk/workflows/frida-hook-flow.md` |
| `http-api/2026-05-06_104300-api-catalog-finish-gate.md` | promoted | API Catalog finish gate | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `http-api/2026-05-05_183400-ui-automation-operation-scripts-for-api-capture.md` | validated | UI automation operation scripts for API capture | ✅ | `analysis/apk/workflows/http-api-documentation-flow.md` |
| `local-proxy/2026-05-01_114300-local-proxy-handler-uri-hook.md` | validated | Hook local proxy handler URI, not just OkHttp | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |
| `local-proxy/2026-05-01_131000-cast-netty-request-for-handler-route.md` | validated | Cast Netty request interfaces for handler routes | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |
| `local-proxy/2026-05-01_132400-netty-aggregated-request-tostring-headers.md` | validated | Netty aggregated request toString can expose headers | ✅ | `analysis/apk/workflows/local-proxy-hook-flow.md` |
| `media-hls/2026-04-30_120007-媒體播放鏈要分控制面-金鑰與資料面.md` | validated | 媒體播放鏈要分控制面、金鑰與資料面 | ✅ | `analysis/apk/workflows/media-hls-analysis-flow.md` |
| `media-hls/2026-05-06_121650-hls-playlist-key-may-be-wrapped.md` | promoted | HLS playlist key may be wrapped | ✅ | `analysis/apk/workflows/media-hls-analysis-flow.md` |

---

## 4. 目標層：`analysis/apk/techniques/` 或 `analysis/apk/tools-and-failures.md`

這類 lessons 包含工具使用技巧或技術說明。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-04-30_120008-aapt-sdk-build-tools-resolve-activity.md` | validated | APK metadata：`aapt` 不在 PATH 時走 SDK build-tools | ✅ | `analysis/apk/tools-and-failures.md` |
| `common/2026-05-06_102800-ui-hierarchy-content-desc-redaction.md` | validated | UI hierarchy content-desc redaction | ✅ | `analysis/apk/tools-and-failures.md` |
| `common/2026-05-06_154500-domain-runtime-baseline-finish-gate.md` | validated | Domain/runtime baseline for real data | ✅ | `analysis/apk/tools-and-failures.md` |

---

## 5. 目標層：`analysis/apk/artifact-gates.md`（`DOCUMENTATION.md` 對應）

這類 lessons 的 Promotion Target 包含 `DOCUMENTATION.md`，適合提取到 `analysis/apk/` 或 `workflow/apk-analysis/artifact-gates.md`。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-06_102100-feature-handoff-finish-gate.md` | promoted | Feature handoff finish gate | ✅ | `workflow/apk-analysis/artifact-gates.md` |
| `common/2026-05-06_104000-detailed-api-list-finish-gate.md` | promoted | Detailed API list finish gate | ✅ | `workflow/apk-analysis/artifact-gates.md` |
| `common/2026-05-06_105026-state-reset-baseline-feature-capture.md` | promoted | State reset baseline before feature capture | ✅ | `workflow/apk-analysis/execution-flow.md` |
| `common/2026-05-06_154500-domain-runtime-baseline-finish-gate.md` | validated | Domain/runtime baseline for real data | ✅ | `workflow/apk-analysis/artifact-gates.md` |

---

## 6. 目標層：`enforcement/`

這類 lessons 影響全庫共用規則。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-01_174100-technique-category-routing.md` | promoted | Technique category routing | ✅ | `enforcement/cross-skill-references.md` |
| `common/2026-05-01_175500-workflow-as-routing-not-technique-dump.md` | promoted | Workflow as routing, not technique dump | ✅ | `enforcement/cross-skill-references.md` |
| `common/2026-05-11_125615-per-round-feedback-checkpoint.md` | promoted | Per-round feedback checkpoint | ✅ | `../feedback-lessons.md` |
| `common/2026-05-08_135700-agent-goals-before-transcripts.md` | candidate | Agent goals before transcripts when resuming | ✅ | `enforcement/dependency-reading.md` |

---

## 統計摘要

| 目標層 | lessons 數量 | 已提取 | 未提取 |
| --- | ---: | ---: | ---: |
| `workflow/apk-analysis/` | 29 | 29 | 0 |
| `intelligence/engineering/analytical-reasoning/` | 17 | 17 | 0 |
| `analysis/apk/workflows/` | 12 | 12 | 0 |
| `analysis/apk/tools-and-failures.md` | 3 | 3 | 0 |
| `workflow/apk-analysis/artifact-gates.md` | 4 | 4 | 0 |
| `enforcement/` | 4 | 4 | 0 |
| **總計** | **69** | **69** | **0** |

---

## 相容性說明

- ✅ `skills/apk-analysis/feedback_history/` 已於 2026-05-13 刪除，所有 lesson 已搬遷至 `feedback/history/apk-analysis/`。
- 此索引僅供 `feedback/` 層發現 lessons，不改變 lesson 的 storage location。
- 提取到目標層時，在原始 lesson 檔案開頭加入 `# Extracted — See <target path>` 標記。
- 此索引應隨 lessons 的新增或提取狀態變更而更新。
