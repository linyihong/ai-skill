# Hook Selection Heuristic（Hook 策略選擇啟發式）

## 問題

Flutter/Dart AOT 分析中，何時該用 Java-level Frida hook、何時該用 Dart-level hook、何時該用靜態分析？

## 原則

- Dart AOT 編譯的程式碼無法由 Java-level Frida hook 攔截
- `dart:io` HttpClient 請求完全繞過 Java OkHttp 層
- 靜態分析 libapp.so 可還原字串常數與部分控制流
- Dart AOT callsite `BL` addresses 是 navigation hints，不是 function hook entry points

## 決策表

| 情境 | 建議做法 | 判斷信號 |
|------|---------|---------|
| App 使用 Flutter 且無 AOT 混淆 | Dart Frida hook（`-dart` 參數） | `libapp.so` 存在且包含 Dart snapshot |
| App 使用 Flutter 且有 AOT 混淆 | 靜態分析 libapp.so + 字串還原 | `libapp.so` 中字串被 XOR/加密 |
| App 使用 `dart:io` HttpClient | Dart-level Frida hook | 流量未出現在 Java OkHttp 層 |
| App 使用 `http` package | Dart-level Frida hook | 同上 |
| Java helper/plugin hooks miss 但 proxy 顯示加密 headers | 先 inspect Dart AOT interceptors | Java hooks 無流量但 pcap 有網路活動 |
| blutter 識別 snapshot 但 crash | 保留失敗證據，切換到 unflutter | blutter 輸出 crash log |

## 不建議的做法

- 不要從 global Dart runtime helpers（如 `LinkedHashMap._set`）開始 hook
- 不要 broad hook 多個 Dart runtime function
- 不要將 callsite `BL` addresses 當作 function hook entry points

## Token 影響

中等。此 atom 在 Flutter app 分析 session 中 lazy-load，約 200-300 tokens。
