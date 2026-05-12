# Language-Specific Pitfalls Catalog

本文件列出各語言/runtime 特有的陷阱與注意事項。承接 [`skills/app-development-guidance/languages/`](../../skills/app-development-guidance/languages/) 的內容，提取為 tool-neutral 的分析參考。

> **相容性規則**：`skills/app-development-guidance/languages/` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## 語言分類

| 語言 | 範圍 | 原始來源 |
|------|------|----------|
| **Dart** | Dart and Flutter-specific concerns | [`languages/dart.md`](../../skills/app-development-guidance/languages/dart.md) |
| **Kotlin/Java** | Kotlin/Java Android-specific code patterns | [`languages/kotlin-java.md`](../../skills/app-development-guidance/languages/kotlin-java.md) |
| **Swift** | Swift/iOS-specific code patterns | [`languages/swift.md`](../../skills/app-development-guidance/languages/swift.md) |
| **TypeScript** | TypeScript frontend/backend client code concerns | [`languages/typescript.md`](../../skills/app-development-guidance/languages/typescript.md) |

## 使用原則

1. **Use this directory only for language-specific or framework-runtime-specific pitfalls** — 如果 lesson 是關於 API design、token lifecycle、logging、storage 或 release controls，把原則放在 `controls/`。
2. **Concrete how-to steps belong in `implementation/`** — 語言陷阱記錄「要注意什麼」，實作步驟在 implementation。
3. **Language guidance changes may require implementation updates** — 當語言特定指引改變工程師實作控制的方式時，在同一變更中更新或驗證對應的 implementation 檔案。

## 與其他層的關係

- `analysis/app-development-guidance/controls-catalog.md` 提供跨平台控制原則。
- `analysis/app-development-guidance/implementation-catalog.md` 提供可建置實作模式。
- `skills/app-development-guidance/languages/` 是原始來源，仍為 active entrypoint。
