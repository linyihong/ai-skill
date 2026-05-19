# Analysis Heuristics

`intelligence/engineering/analytical-reasoning/heuristics/` 存放分析過程中的啟發式判斷規則，主要源自 APK 分析領域。

## Scope

本目錄負責「未來遇到類似情況時應如何判斷」：

- Hook 策略選擇啟發式（何時用 Frida、何時靜態分析、何時用 Dart-level hook）
- API 文件完整性判斷啟發式
- 代理導流診斷啟發式
- 媒體串流鏈完整性判斷啟發式
- UI 操作穩定性啟發式（何時用 bounded scroll、operation script、API-first replay）
- UI-to-API 歸因啟發式（前景 package 驗證、feature context 驗證、操作時間窗對齊）
- 加密模式、padding、參數來源與 hook timing 的決策表
- Flutter + Java/Kotlin 混合架構的安全層映射啟發式

## 與其他層的關係

- `analysis/apk/workflows/` 提供操作步驟，本目錄提供「何時該用哪個步驟」的判斷
- `intelligence/engineering/analytical-reasoning/evidence-first-routing.md` 決定分析路線，本目錄決定路線內的技術選擇
- 如果內容只是某次失敗症狀與診斷入口，放 `failure/`；如果是可重用的預防規則與決策表，放本目錄。

## 目前 atoms

| Atom | 說明 | 來源 | 跨領域推廣 |
|------|------|------|-----------|
| [`hook-selection.md`](hook-selection.md) | Hook 策略選擇啟發式 — 根據 signal（Flutter/Dart AOT、Java OkHttp、Socket）選擇 hook 策略的決策表 | `skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除） | — |
| [`api-documentation-completeness.md`](api-documentation-completeness.md) | API 文件完整性啟發式 — 何時開始 API Catalog、何時完成、Field Confidence 判斷的決策表 | `skills/apk-analysis/techniques/http-api/`（已刪除） | Field Confidence 判斷表已提取到 [`intelligence/engineering/heuristics/field-confidence-judgment.md`](../../heuristics/field-confidence-judgment.md) |
| [`local-proxy-routing-diagnosis.md`](local-proxy-routing-diagnosis.md) | Local Proxy vs TLS Pinning 診斷 — 根據 signal 判斷流量走 local proxy 還是 TLS pinning 的決策表 | `skills/apk-analysis/techniques/local-proxy/`（已刪除） | — |
| [`javascript-bitwise-64bit-truncation.md`](javascript-bitwise-64bit-truncation.md) | JavaScript 位元運算子截斷 64-bit 指標啟發式 — 在 Frida hook 中永遠使用 `NativePointer` 方法或字串操作做 64-bit 指標運算，避免 JS 位元運算子的 32-bit 截斷 | `feedback/history/apk-analysis/flutter-dart-aot/` | — |
| [`ui-operation-stability.md`](ui-operation-stability.md) | UI 操作穩定性啟發式 — 決定何時該用 bounded scroll、operation script、API-first replay 的決策表 | `feedback/history/apk-analysis/http-api/` + `feedback/history/apk-analysis/common/` | — |
| [`ui-to-api-attribution.md`](ui-to-api-attribution.md) | UI-to-API 歸因啟發式 — 確保 UI 操作能正確對應到 API 請求，避免 attribution 錯誤的決策表 | `feedback/history/apk-analysis/http-api/` + `feedback/history/apk-analysis/common/` | — |
| [`dart-aot-trivial-field-accessor-detection.md`](dart-aot-trivial-field-accessor-detection.md) | Dart AOT trivial field accessor 檢測 — 短函數（≤12 bytes）幾乎是 field accessor 而非業務 getter | `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_081500-iv-f5447c-trivial-field-accessor-not-iv-getter.md` | — |
| [`dart-aot-padright-substring-dispatch-hooking.md`](dart-aot-padright-substring-dispatch-hooking.md) | Dart AOT 字串操作假設與 dispatch hooking — `padRight` 可能用 null bytes、`substring` 越界不 crash、dispatch hooking 技巧 | `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_073700-dart-aot-padright-null-bytes-substring-noop.md` | — |
| [`dart-encrypt-package-mode-detection.md`](dart-encrypt-package-mode-detection.md) | Dart `encrypt` 套件 AES 模式檢測 — 如何區分 CBC vs CTR/SIC vs GCM，包含 live proxy test 確認方法 | `feedback/history/apk-analysis/flutter-dart-aot/2026-05-13_232600-dart-encrypt-package-aes-mode-cbc-not-ctr.md` | — |
| [`frida-spawn-vs-attach-init-timing.md`](frida-spawn-vs-attach-init-timing.md) | Frida spawn vs attach 初始化時機 — 初始化函數必須用 spawn 模式；Frida JS 無 Buffer API | `feedback/history/apk-analysis/common/2026-05-14_073700-frida-spawn-vs-attach-init-timing-no-buffer.md` + `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_081500-initialize-hook-does-fire-in-spawn.md` | — |
| [`dart-aot-pool-constant-pp-peep.md`](dart-aot-pool-constant-pp-peep.md) | Dart AOT pool constant (PP_peep) 啟發式 — `string_refs.jsonl` 中 `kind:"PP_peep"` 的條目可直接揭露硬編碼字串常數（如 AES key、secret），無需反組譯 | `feedback/history/apk-analysis/flutter-dart-aot/2026-05-15_094700-dart-aot-pool-constant-pp-peep-hardcoded-key.md` | — |
