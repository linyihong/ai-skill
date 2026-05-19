# Flutter + Java Hybrid Security Layer Mapping（混合架構安全層映射啟發式）

## 問題

分析 Flutter / Dart AOT + Java/Kotlin bridge 的 Android App 時，安全邏輯可能分散在 Dart runtime、Java SDK、本機 proxy、gateway 與 response decoder。若只從單一層 hook，容易誤判簽名、加密、TLS gate 或 API 呼叫圖。

## 原則

- Flutter AOT 的 Dart 程式碼與 Java/Kotlin bridge 是不同執行層，hook 策略必須分層對應。
- 多層安全防護通常以「TLS / gateway → HTTP header validation → request signing → response unwrap/decrypt」堆疊出現。
- 不同 API family 可能使用不同 signing format、identity material 或 response decrypt path，不應假設全 App 共用一套。
- Base64、固定長度前綴或 wrapper field 不代表 plaintext；必須用 decoder / decrypt boundary 驗證。

## 決策表

### 1. 框架與執行層識別

| 情境 | 判斷 | 下一步 |
| --- | --- | --- |
| APK 同時包含 Flutter native libraries 與大量 Java/Kotlin classes | 混合架構 | 分別建立 Dart AOT、Java bridge、network gateway 的觀察點 |
| Java HTTP hook miss，但 pcap 顯示業務流量 | HTTP path 可能在 Dart / native 層 | 轉向 Dart AOT HttpClient、socket 或高語意 request object |
| Java/Kotlin 層出現 local server / proxy / bridge 類別 | 可能有本機轉發層 | 先確認 loopback / handler evidence，再 hook handler |
| Flutter WebView 或 H5 bridge 存在 | API path 可能分成 native shell 與 embedded web | 將 native facade、H5 entry URL、H5 XHR 分線記錄 |

### 2. 安全層識別

| 信號 | 可能層級 | 驗證方式 |
| --- | --- | --- |
| App-owned request 成功，但外部 client 相同 payload 失敗 | TLS / gateway / fingerprint gate | 比對 TLS path、gateway response class 與 App-owned client 行為 |
| Request header 或 body 有 opaque signature field | Request signing layer | Hook signing input、canonicalization 與 output boundary，值本身留在專案 evidence |
| Response wrapper 看似 encoded / encrypted | Response unwrap/decrypt layer | Hook decoder return value或建立 raw-to-decoded fixture |
| 同一 App 不同 API family 使用不同 request shape | API family-specific signing | 按 route family 建立 signing matrix，不套用單一路徑 |

### 3. API family mapping

| 情境 | 建議分類 |
| --- | --- |
| 需要 session / identity material | Authenticated API family |
| 不需要 session，但需要 install/device/bootstrap material | Guest/bootstrap API family |
| 只有 public read parameters，但仍經過 gateway/signing | Read-only signed API family |
| 進入 WebView / H5 後由另一個 host 或 runtime 發出 | Embedded web API family |

## 不建議的做法

- 不要假設所有 API 共用同一套 signature、header 或 decrypt path。
- 不要只從 Java 層 hook；Dart AOT 的 HTTP 請求可能完全繞過 Java HTTP stack。
- 不要只看 Base64 / prefix 就判斷 plaintext 或 encryption mode。
- 不要把具體 header 名、service id、host、port、payload、token 或 raw key 寫進 reusable atom。

## 驗證

- 已建立分層 map：Dart AOT、Java/Kotlin bridge、local proxy / gateway、response decoder。
- 每個 API family 都有自己的 signing / identity / decrypt boundary 判斷，不共用未驗證假設。
- 對外部 client 或 SDK 實作只回填 mode、boundary、normalization 與 pass/fail，不回填 raw secret 或 project-only values。
- 若需要實作 live client，先套用 [`../live-readiness-gates.md`](../live-readiness-gates.md)。

## 相關 atoms

- [`hook-selection.md`](hook-selection.md)
- [`local-proxy-routing-diagnosis.md`](local-proxy-routing-diagnosis.md)
- [`dart-encrypt-package-mode-detection.md`](dart-encrypt-package-mode-detection.md)
- [`../signals/flutter-dart-aot-detection.md`](../signals/flutter-dart-aot-detection.md)

## Token 影響

中等。此 atom 在分析 Flutter + Java/Kotlin 混合架構 App 時 lazy-load，約 250-350 tokens。
