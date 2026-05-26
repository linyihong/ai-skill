# AGENTS.md — Generic Agent Bootstrap Entry

> **本檔為 thin generic agent entry**。適用任何遵循 `AGENTS.md` 慣例的 AI agent（Codex、Cursor partial、Aider、Cline、其他 AGENTS.md-aware tools）。Canonical bootstrap obligations 在 [`runtime/core-bootstrap.yaml`](runtime/core-bootstrap.yaml) + [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) companion。

## 啟動序列

1. 讀 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) 與 [`runtime/runtime.db`](runtime/runtime.db) — 必讀規則 + 所有 obligations（Bootstrap Receipt、Cognitive Mode 報告 per-turn block 等）
2. 讀 [`README.md`](README.md) — OS layout
3. 讀 [`ai-tools/README.md`](ai-tools/README.md) — **routing hub**：選擇你的 AI 工具對應的 adapter（claude / cursor / roo / codex / future tools）
4. 依 routing hub 指示讀對應的 tool adapter 取得 tool-specific 操作注意事項

## 修改本檔的規則

本檔是 entry pointer，不是 canonical content。修改前先讀 [`runtime/bootstrap-entry-points.yaml`](runtime/bootstrap-entry-points.yaml) §thinness_rule。新 obligation 加到 [`runtime/core-bootstrap.yaml`](runtime/core-bootstrap.yaml) + CORE_BOOTSTRAP.md（cross-tool）— 不加到本檔。Commit-msg hook 會擋下違反 thinness 的修改。

本檔不 link 到單一工具的 adapter，避免把 generic entry 鎖死到特定工具。Routing 統一經 `ai-tools/README.md`。
