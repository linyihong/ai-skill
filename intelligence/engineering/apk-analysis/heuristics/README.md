# APK Analysis Heuristics

`intelligence/engineering/apk-analysis/heuristics/` 存放 APK 分析領域專屬的啟發式判斷規則，這些規則尚未抽象到跨領域通用的 `analytical-reasoning/` 層級。

## Scope

本目錄負責：

- Flutter+Java 混合架構安全層映射啟發式
- 其他 APK 分析領域專屬的啟發式規則（未來擴充）

## 與其他層的關係

- `intelligence/engineering/analytical-reasoning/heuristics/` 存放跨領域通用的啟發式規則
- `feedback/history/apk-analysis/` 是這些啟發式的原始來源（raw experience）
- `analysis/apk/` 提供 instance-oriented 的操作記錄與技術文件

## 目前 atoms

| Atom | 說明 | 來源 |
|------|------|------|
| [`flutter-java-hybrid-security-layer-mapping.md`](flutter-java-hybrid-security-layer-mapping.md) | Flutter+Java 混合架構安全層映射啟發式 — 系統性識別框架、安全層（TLS 閘道 → 標頭驗證 → 請求簽名 → 回應加密）、API 簽名格式分類、呼叫圖建立的決策表 | [`feedback/history/apk-analysis/common/2026-05-18_172200-flutter-java-hybrid-security-architecture-overview.md`](../../../../feedback/history/apk-analysis/common/2026-05-18_172200-flutter-java-hybrid-security-architecture-overview.md) |
