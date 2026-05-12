# APK Analysis Heuristics

`intelligence/engineering/apk-analysis/heuristics/` 存放 APK 分析過程中的啟發式判斷規則。

## Scope

本目錄負責：

- Hook 策略選擇啟發式（何時用 Frida、何時靜態分析、何時用 Dart-level hook）
- API 文件完整性判斷啟發式
- 代理導流診斷啟發式
- 媒體串流鏈完整性判斷啟發式

## 與其他層的關係

- `analysis/apk/workflows/` 提供操作步驟，本目錄提供「何時該用哪個步驟」的判斷
- `intelligence/engineering/apk-analysis/evidence-first-routing.md` 決定分析路線，本目錄決定路線內的技術選擇

## 目前 atoms

| Atom | 說明 | 來源 | 跨領域推廣 |
|------|------|------|-----------|
| [`hook-selection.md`](hook-selection.md) | Hook 策略選擇啟發式 — 根據 signal（Flutter/Dart AOT、Java OkHttp、Socket）選擇 hook 策略的決策表 | `skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除） | — |
| [`api-documentation-completeness.md`](api-documentation-completeness.md) | API 文件完整性啟發式 — 何時開始 API Catalog、何時完成、Field Confidence 判斷的決策表 | `skills/apk-analysis/techniques/http-api/`（已刪除） | Field Confidence 判斷表已提取到 [`intelligence/engineering/heuristics/field-confidence-judgment.md`](../../heuristics/field-confidence-judgment.md) |
| [`local-proxy-routing-diagnosis.md`](local-proxy-routing-diagnosis.md) | Local Proxy vs TLS Pinning 診斷 — 根據 signal 判斷流量走 local proxy 還是 TLS pinning 的決策表 | `skills/apk-analysis/techniques/local-proxy/`（已刪除） | — |
