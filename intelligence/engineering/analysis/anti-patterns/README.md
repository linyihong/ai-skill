# Analysis Anti-Patterns

`intelligence/engineering/analysis/anti-patterns/` 存放分析過程中常見的反模式與錯誤做法，主要源自 APK 分析領域。

## Scope

本目錄負責：

- Hook timing 錯誤（過早 hook、relocation 未完成）
- 代理設定錯誤（導流與 TLS 混淆）
- API 文件過早細節化（未先確認完整 catalog）
- 分析路徑選擇錯誤（在不該深入的點浪費時間）

## 與其他層的關係

- `anti-patterns/`（根目錄）存放跨領域通用反模式，本目錄只放分析領域 specific
- `intelligence/engineering/analysis/failure/` 記錄具體失敗案例，本目錄記錄可預防的錯誤模式

## 目前 atoms

| Atom | 說明 | 來源 |
|------|------|------|
| [`early-hook-instability.md`](early-hook-instability.md) | 過早 hook 導致不穩定 — 在 app 初始化完成前 hook 導致 crash 的症狀表與診斷方法 | `skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除） |
