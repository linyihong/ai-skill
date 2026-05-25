# Claude Code Bootstrap Entry（thin pointer）

> **IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
>
> 本檔為 **thin tool-entry pointer**。所有 obligation、format、enum、example 的 canonical 來源在 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md)。Session 啟動 first turn 必須讀 CORE_BOOTSTRAP.md 並遵守其中**所有** obligations（含 Bootstrap Receipt、Cognitive Mode 報告 per-turn block、Resume clause、knowledge-update-flow 等）。
>
> Summary 的「Resume directly」是對話 framing，**不豁免** runtime / governance bootstrap。Resume / continuation session 同樣須走完 bootstrap。

## 啟動序列

1. 讀 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) — 必讀規則 + Bootstrap Receipt + Cognitive Mode 報告 + 全部 per-session / per-turn obligations
2. 讀 [`README.md`](README.md) — OS layout
3. 查 [`runtime/runtime.db`](runtime/runtime.db) — 目前 phase / obligations / gates

Claude Code 工具專屬注意事項：[`ai-tools/agent/claude.md`](ai-tools/agent/claude.md)。

## 修改本檔的規則

本檔是 **entry pointer**，不是 canonical content。修改前先讀 [`runtime/bootstrap-entry-points.yaml`](runtime/bootstrap-entry-points.yaml) §thinness_rule。新 obligation 加到 CORE_BOOTSTRAP.md（cross-tool）或 `ai-tools/agent/claude.md`（Claude-specific）— 不加到本檔。Commit-msg hook 會擋下違反 thinness 的修改。
