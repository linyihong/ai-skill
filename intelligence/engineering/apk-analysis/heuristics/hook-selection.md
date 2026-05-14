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
| Frida constructor chain 顯示 `PBC.ctor`（PaddedBlockCipher）但 mode（CBC vs SIC vs GCM）無法從 block count 區分 | 執行 live proxy test：分別用不同 mode 加密相同明文，比對 HTTP response status code | `processBlock` count 為 43（ambiguous：CBC=42+1, GCM=1+42, CTR+PKCS7=42+1）；`GCMBlockCipher.ctor` 可能屬於不同 encryption group |
| Frida 顯示自訂 AES 實作使用 8-byte 金鑰（非標準 16/24/32 bytes） | 不要嘗試猜測金鑰衍生方式；直接反組譯 Dart AOT 函式或新增 Frida hook 追蹤金鑰擴展階段 | 標準 Java `javax.crypto.Cipher` 或 BC `AESEngine` 拒絕 8-byte 金鑰；所有標準填充/雜湊方式（repeat、zero-pad、MD5、SHA-256、concat、XOR）都無法匹配 Dart 輸出 |

## 不建議的做法

- 不要從 global Dart runtime helpers（如 `LinkedHashMap._set`）開始 hook
- 不要 broad hook 多個 Dart runtime function
- 不要將 callsite `BL` addresses 當作 function hook entry points

## Token 影響

中等。此 atom 在 Flutter app 分析 session 中 lazy-load，約 200-300 tokens。
