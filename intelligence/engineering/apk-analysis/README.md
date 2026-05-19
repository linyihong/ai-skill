# APK Analysis Engineering Intelligence

`intelligence/engineering/apk-analysis/` 存放 APK 分析領域專屬的工程判斷。這些 atoms 比 `analysis/apk/` 更偏向「如何判斷與取捨」，但尚未抽象到跨領域通用的 `analytical-reasoning/` 層級。

## 子目錄

| 子目錄 | 描述 |
|--------|------|
| [`heuristics/`](heuristics/README.md) | APK 分析領域專屬啟發式，例如 Flutter + Java 混合架構安全層映射 |

## 與其他層的關係

- `analysis/apk/` 提供具體分析方法與 workflow 操作細節。
- `workflow/apk-analysis/` 提供端到端執行流程。
- `intelligence/engineering/analytical-reasoning/` 保留跨領域分析推理智慧；APK-only atoms 先放在本目錄。
- `feedback/history/apk-analysis/` 是 lesson 原始來源，promotion 後仍需保留可追溯連結。
