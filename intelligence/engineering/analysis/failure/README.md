# Analysis Failure Intelligence

`intelligence/engineering/analysis/failure/` 存放分析過程中具體的失敗模式與診斷知識，主要源自 APK 分析領域。

## Scope

本目錄負責：

- Frida spawn race condition 與 relocation timing 失敗
- JIT mismatch 與 AOT 快照載入失敗
- 代理 TLS 握手失敗的診斷
- 媒體串流解碼失敗的判讀

## 與其他層的關係

- `shared-rules/failure-learning-system.md` 定義通用 failure learning 框架，本目錄存放分析領域的具體 failure patterns
- `intelligence/engineering/analysis/anti-patterns/` 記錄可預防的錯誤模式，本目錄記錄已發生的失敗與診斷方式

## 目前 atoms

| Atom | 說明 | 來源 |
|------|------|------|
| [`frida-spawn-race.md`](frida-spawn-race.md) | Frida spawn race condition — Frida 在 spawn 模式下的 race condition 診斷與緩解策略 | `skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除） |
| [`javascript-bitwise-64bit-truncation.md`](javascript-bitwise-64bit-truncation.md) | JavaScript 位元運算子截斷 64-bit 指標 — 所有 `tryReadDartString` 回傳 `undecoded` 或 `hexdump` 報 access violation 時的診斷入口，解法指向 heuristics | `feedback/history/apk-analysis/flutter-dart-aot/` |
| [`processBlock-count-ambiguity.md`](processBlock-count-ambiguity.md) | processBlock 呼叫次數歧義 — Frida 攔截到 43 次 `processBlock` 呼叫時無法單獨從 block count 區分 CBC/GCM/CTR，需 live proxy test 確認 | `feedback/history/apk-analysis/flutter-dart-aot/` |
| [`custom-dart-aes-8byte-key-not-reproducible.md`](custom-dart-aes-8byte-key-not-reproducible.md) | Dart 自訂 AES 8-byte 金鑰無法以標準 AES 重現 — Frida 顯示 8-byte 金鑰時，標準 Java/BC AES 拒絕該金鑰長度，所有標準填充/雜湊方式都無法匹配 Dart 輸出 | `feedback/history/apk-analysis/flutter-dart-aot/` |
