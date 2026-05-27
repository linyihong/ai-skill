# Claude Code Bootstrap Entry（thin pointer）

> **IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
> 本檔為 **thin tool-entry pointer**。所有 obligation / format / enum /
> example 的 canonical 來源在 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md)
> 與 [`runtime/core-bootstrap.yaml`](runtime/core-bootstrap.yaml)。
> Summary 的「Resume directly」不是豁免；resume / continuation session
> 同樣必須走 bootstrap。

## 啟動序列

1. 讀 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) — companion entry。
2. 讀 [`README.md`](README.md) — OS layout
3. 讀 [`runtime/core-bootstrap.yaml`](runtime/core-bootstrap.yaml) — canonical contract。
4. 查 [`runtime/runtime.db`](runtime/runtime.db) — phase / obligations / gates。

依 canonical contract 執行 Bootstrap Receipt、每輪 Cognitive Mode 報告與 close-loop checks。

Claude Code 工具專屬注意事項與 hook behavior 見 [`ai-tools/agent/claude.md`](ai-tools/agent/claude.md)。

## 修改本檔的規則

本檔是 entry pointer，不是 canonical content。修改前先讀 [`runtime/bootstrap-entry-points.yaml`](runtime/bootstrap-entry-points.yaml) §thinness_rule。新 obligation 加到 canonical bootstrap source；Claude-specific 內容放 [`ai-tools/agent/claude.md`](ai-tools/agent/claude.md)。
