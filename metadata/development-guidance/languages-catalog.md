# Language-Specific Pitfalls Catalog

本文件列出各語言/runtime 特有的陷阱與注意事項。原承接 ``skills/app-development-guidance/languages/`` 的內容（已刪除），提取為 tool-neutral 的分析參考。

> **遷移狀態**：`skills/app-development-guidance/languages/` 已刪除。此文件為 canonical source，新內容請直接寫入此文件。

## 語言分類

| 語言 | 範圍 | 原始來源 |
|------|------|----------|
| **Dart** | Dart and Flutter-specific concerns | ``languages/dart.md``（已刪除） |
| **Kotlin/Java** | Kotlin/Java Android-specific code patterns | ``languages/kotlin-java.md``（已刪除） |
| **Swift** | Swift/iOS-specific code patterns | ``languages/swift.md``（已刪除） |
| **TypeScript** | TypeScript frontend/backend client code concerns | ``languages/typescript.md``（已刪除） |

## 使用原則

1. **Use this directory only for language-specific or framework-runtime-specific pitfalls** — 如果 lesson 是關於 API design、token lifecycle、logging、storage 或 release controls，把原則放在 `controls/`。
2. **Concrete how-to steps belong in `implementation/`** — 語言陷阱記錄「要注意什麼」，實作步驟在 implementation。
3. **Language guidance changes may require implementation updates** — 當語言特定指引改變工程師實作控制的方式時，在同一變更中更新或驗證對應的 implementation 檔案。

## 與其他層的關係

- `analysis/development-guidance/controls-catalog.md` 提供跨平台控制原則。
- `analysis/development-guidance/implementation-catalog.md` 提供可建置實作模式。
- `skills/app-development-guidance/languages/` 是原始來源，已刪除。內容已由本文件承接。
