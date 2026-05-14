# IDE 生態系統知識

本目錄記錄 IDE 生態系統的工程智慧，包含 VS Code、JetBrains 等編輯器的內部機制與可重複使用的操作知識。

## 子目錄

| 檔案 | 描述 |
|------|------|
| [`vscode-extension-global-state.md`](vscode-extension-global-state.md) | VS Code Extension 全域設定的 SQLite 儲存機制、資料庫結構、查詢/修改方法與注意事項。 |

## 與其他層的關係

- `ai-tools/agent/` 中的工具文件（如 `roo.md`）記錄特定工具的設定細節（如 Roo Code 的 key 與欄位），本層提供通用的底層知識（VS Code Extension 的 SQLite 儲存機制）。
- `enforcement/failure-patterns/language-preference-drift.md` 記錄語言偏好漂移的失效模式，本層提供修改語言偏好的具體技術方法。
- `intelligence/engineering/` 記錄軟體工程經驗法則與架構決策智慧，與本層同為 `intelligence/` 下的獨立子目錄，兩者互補（engineering 偏軟體開發，ide 偏 IDE 生態系統）。

---

← [回到 Intelligence 索引](../README.md)
