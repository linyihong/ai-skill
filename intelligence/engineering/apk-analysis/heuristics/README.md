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

（pilot 階段逐步建立）
