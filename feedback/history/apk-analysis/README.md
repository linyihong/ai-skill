# APK Analysis Feedback History

## 分類

| 分類 | 數量 | 說明 |
|------|------|------|
| [`common/`](common/) | 57 | 跨分類或通用 lesson（工具選擇、流程、UI、分析策略等） |
| [`flutter-dart-aot/`](flutter-dart-aot/) | 23 | Flutter/Dart AOT 相關 hook 與分析 |
| [`http-api/`](http-api/) | 25 | HTTP API 分析、文件化、UI 操作流程 |
| [`local-proxy/`](local-proxy/) | 8 | Local proxy 設定、診斷、hook |
| [`media-hls/`](media-hls/) | 3 | Media/HLS 串流分析 |
| [`dynamic-capture/`](dynamic-capture/) | 2 | 動態捕獲相關 |

## Recent (2026-06-22)

| Slug | Category |
|------|----------|
| `common/2026-06-22_120000-okhttp-r8-obfuscated-request-realinterceptorchain-proceed` | OkHttp R8 hook overload |
| `local-proxy/2026-06-22_120100-mitm-cdn-visible-primary-api-invisible-pinning-tier` | MITM CDN vs API pinning tier |
| `http-api/2026-06-22_120200-static-ktor-strings-not-dynamic-business-api-client` | Retrofit vs Ktor dynamic client |
| `http-api/2026-06-22_130000-in-session-api-host-differs-from-cold-start-primary` | Cold vs in-session API host failover |
| `common/2026-06-22_130100-r8-obfuscated-okhttp-response-needs-converter-hook` | Obfuscated Response / converter hook |
| `common/2026-06-22_131000-retrofit-gson-fromjson-hook-api-response-plaintext` | Gson.fromJson API JSON capture |
| `http-api/2026-06-22_141500-static-list-endpoint-zero-hit-check-detail-embedded-catalog` | List path 0 hit → check embedded catalog |
| `http-api/2026-06-22_141600-static-waterfall-path-zero-hit-check-shelf-layout-gating` | Waterfall path 0 hit → shelf layout gating |
| `http-api/2026-06-22_141700-custom-request-signatures-block-standalone-sdk-until-re-or-relay` | Custom sign blocks standalone SDK |
| `http-api/2026-06-22_141800-okhttp-interceptor-sorted-map-canonical-sign-pattern` | Interceptor sorted-map canonical for sign |
| `common/2026-06-22_141900-crypto-util-name-sha256-does-not-imply-plain-sha256` | sha256* name ≠ plain SHA256 verify gate |
| `common/2026-06-22_142000-jni-dynamic-registernatives-no-standard-java-com-export` | Dynamic JNI / no Java_com_* export |
| `http-api/2026-06-22_142100-encrypted-request-time-double-native-call-second-value-in-sign-map` | Encrypted-time header double native call |
| `common/2026-06-22_142200-frida-enumerate-loaded-classes-causes-script-load-timeout` | Avoid enumerateLoadedClasses timeout |
| `common/2026-06-22_142300-hmac-sha256-per-mode-keys-in-native-rodata` | Plain SHA256 fail → HMAC + per-mode rodata keys |
| `common/2026-06-22_142400-system-loadlibrary-name-overrides-static-so-assumption` | loadLibrary name vs wrong protection .so |
| `http-api/2026-06-22_142500-partial-offline-sign-hmac-solved-requesttime-session-remain-relay` | Hybrid SDK when sign offline, encrypted-time native |
| `common/2026-06-22_142600-frida-python-attach-java-bridge-use-cli-subprocess-for-rpc` | Python attach Java undefined → Frida CLI RPC |
| `http-api/2026-06-22_142700-opaque-api-blob-native-decrypt-secondary-compression-json` | Opaque blob → decryptStr + zlib → JSON |
| `http-api/2026-06-22_142800-native-encrypted-time-getter-may-cache-first-rotates-consecutive-same` | Encrypted-time getter: 1st rotates, 2nd+ same in burst |
| `http-api/2026-06-22_142900-wire-json-field-names-differ-from-gson-bean-paths` | Wire JSON keys ≠ Gson bean field names |
| `http-api/2026-06-22_143000-encrypt-plaintext-may-be-millis-plus-fingerprint-not-surface-seed` | Encrypt plaintext may be millis + fingerprint, not surface seed |
| `http-api/2026-06-22_143100-apk-signing-cert-sha1-may-be-stable-encrypt-plaintext-suffix` | Signing cert SHA1 as stable colon-hex suffix |
| `http-api/2026-06-22_143200-custom-f3aes-label-may-still-be-standard-aes128-cbc-with-rodata-key-iv` | Custom crypto label may still be standard AES after key hook |
| `http-api/2026-06-22_143300-native-decrypt-mode-may-reuse-encrypt-mode-key-material` | Decrypt mode N may share encrypt mode N keys |
| `http-api/2026-06-22_143400-login-may-bootstrap-with-sentinel-uid-session-headers` | Login may bootstrap with sentinel uid/session |

## 來源

所有 lesson 原位於 `skills/apk-analysis/feedback_history/`，已於 2026-05-13 搬遷至此，舊目錄已刪除。
