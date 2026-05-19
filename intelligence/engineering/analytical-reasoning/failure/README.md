# Analysis Failure Intelligence

`intelligence/engineering/analytical-reasoning/failure/` 存放分析過程中具體的失敗模式與診斷知識，主要源自 APK 分析領域。

## Scope

本目錄負責「看到某種失敗症狀時如何診斷」：

- Frida spawn race condition 與 relocation timing 失敗
- JIT mismatch 與 AOT 快照載入失敗
- 代理 TLS 握手失敗的診斷
- 媒體串流解碼失敗的判讀
- 加密模式、參數 normalization 或 hook 讀值造成的誤判診斷

## 與其他層的關係

- `enforcement/failure-learning-system.md` 定義通用 failure learning 框架，本目錄存放分析領域的具體 failure patterns
- `intelligence/engineering/analytical-reasoning/anti-patterns/` 記錄可預防的錯誤模式，本目錄記錄已發生的失敗與診斷方式
- 如果 failure 入口已經有穩定預防規則，解法應指向 `heuristics/`，不要在 failure atom 重複大量 procedure。

## 目前 atoms

| Atom | 說明 | 來源 |
|------|------|------|
| [`frida-spawn-race.md`](frida-spawn-race.md) | Frida spawn race condition — Frida 在 spawn 模式下的 race condition 診斷與緩解策略 | `skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除） |
| [`javascript-bitwise-64bit-truncation.md`](javascript-bitwise-64bit-truncation.md) | JavaScript 位元運算子截斷 64-bit 指標 — 所有 `tryReadDartString` 回傳 `undecoded` 或 `hexdump` 報 access violation 時的診斷入口，解法指向 heuristics | `feedback/history/apk-analysis/flutter-dart-aot/` |
| [`processBlock-count-ambiguity.md`](processBlock-count-ambiguity.md) | processBlock 呼叫次數歧義 — Frida 攔截到 43 次 `processBlock` 呼叫時無法單獨從 block count 區分 CBC/GCM/CTR，需 live proxy test 確認 | `feedback/history/apk-analysis/flutter-dart-aot/` |
| [`custom-dart-aes-8byte-key-not-reproducible.md`](custom-dart-aes-8byte-key-not-reproducible.md) | Dart AES 短參數誤判 — Frida 顯示短 key/IV 時，先診斷是否為 normalization 前 material，解法指向 mode detection heuristic | `feedback/history/apk-analysis/flutter-dart-aot/` |
