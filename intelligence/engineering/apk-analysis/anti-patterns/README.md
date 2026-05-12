# APK Analysis Anti-Patterns

`intelligence/engineering/apk-analysis/anti-patterns/` 存放 APK 分析過程中常見的反模式與錯誤做法。

## Scope

本目錄負責：

- Hook timing 錯誤（過早 hook、relocation 未完成）
- 代理設定錯誤（導流與 TLS 混淆）
- API 文件過早細節化（未先確認完整 catalog）
- 分析路徑選擇錯誤（在不該深入的點浪費時間）

## 與其他層的關係

- `anti-patterns/`（根目錄）存放跨領域通用反模式，本目錄只放 apk-analysis specific
- `intelligence/engineering/apk-analysis/failure/` 記錄具體失敗案例，本目錄記錄可預防的錯誤模式

## 目前 atoms

（pilot 階段逐步建立）
