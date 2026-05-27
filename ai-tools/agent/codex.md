# Codex Agent 規則

本檔只記錄 Codex 與其他 agent 工具不同的地方。通用 bootstrap 與 obligation 來源一律是 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) 與 [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)。

## Thin Entry Point

Root [`AGENTS.md`](../../AGENTS.md) 是 Codex / generic agent 的薄入口。它必須只指向：

1. [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)
2. [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)

不要在 `AGENTS.md` 或本 adapter 複製 Bootstrap Receipt、Cognitive Mode、close-loop、runtime phase 或 dependency checklist 的細節；需要時從 canonical source 讀取。

## Codex 差異

- Codex / generic AGENTS-aware tools 通常會自動讀取 repo root 的 `AGENTS.md`。
- Codex 工作時優先依使用者目前 workspace 與 repository pattern 行動；不要建立平行的工具專屬規則來源。
- 若 Codex 需要工具專屬操作說明，只保留在本檔；跨工具規則放回 `enforcement/`。
- 修改 runtime 行為時，先更新 canonical runtime contract 或 source，再依本 repo 的 runtime compile / refresh / validate 流程驗證。

## 配置邊界

`AGENTS.md` 是 entry pointer，不是規則正文。Codex-specific 注意事項留在本檔；可重用規則、workflow、decision promotion 與 runtime contract 不在 adapter 中重述。

← [回到 AI 工具索引](../README.md)
