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
| Frida constructor chain 顯示 `PBC.ctor`（PaddedBlockCipher）但 mode（CBC vs SIC vs GCM）無法從 block count 區分 | 依 [`dart-encrypt-package-mode-detection.md`](dart-encrypt-package-mode-detection.md) 分組驗證 mode、padding 與 output length | `processBlock` count 可能同時符合多種 mode；`GCMBlockCipher.ctor` 可能屬於不同 encryption group |
| Frida 顯示短 key/IV（非標準 16/24/32 bytes） | 不要直接宣稱自訂 AES；先追蹤 zero-pad、truncate、copy 或 Uint8List normalization，再用 fixture / live proxy test 驗證 | Hook 到的參數可能是 normalization 前的 material，不一定是 cipher 最終輸入 |
| 初始化階段的 hook 在 attach 模式下從未觸發 | 改用 spawn 模式（`frida -U -f <package> -l script.js`） | 同一 hook 在 spawn 模式下正常觸發，attach 模式下不觸發 |
| Dart AOT `padRight` 回傳包含 null bytes（`\0`）的字串 | 不要假設使用標準 Dart space padding；hex dump 確認填充字元 | `padRight` 回傳字串的 hex dump 顯示 `\0` bytes 而非 `0x20` |
| Dart AOT `substring(0, N)` 在短字串上不 crash | 不要假設 `substring(0, N)` 代表字串長度 >= N；檢查實際字串長度 | `substring(0, 32)` 在 16-char 字串上回傳原字串，無異常 |
| 需要捕獲 vtable dispatch 的完整流程 | hook dispatch 指令（`UBFX` + `LDR` + `BLR`）而非猜測的目標函數 | Dispatch hook 成功捕獲所有可能的目標函數 |

## 不建議的做法

- 不要從 global Dart runtime helpers（如 `LinkedHashMap._set`）開始 hook
- 不要 broad hook 多個 Dart runtime function
- 不要將 callsite `BL` addresses 當作 function hook entry points

## Token 影響

中等。此 atom 在 Flutter app 分析 session 中 lazy-load，約 200-300 tokens。
