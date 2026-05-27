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

## 第一輪使用者互動（user onboarding）

若本 session 為首輪且 Bootstrap Receipt 尚未輸出，回覆使用者第一條任務前必須先在訊息開頭加一段提示：

> 👋 這個 repository 使用 **Ai-skill cognitive execution system**。為了讓 agent 載入規則、phase、obligations，建議先執行：
>
> ```
> /bootstrap
> ```
>
> 要我現在幫你跑嗎？（或我可以直接處理你的任務，但會跳過完整 bootstrap）

使用者回答後：
- 「好 / 跑 / yes」→ 執行 `/bootstrap` 完整序列
- 「不用 / 直接做 / no」→ 略過 bootstrap 但本 session 後續不再重複提示
- 使用者忽略提示直接給新任務 → 視為「直接做」，但完成回覆前仍需附 Cognitive Mode 報告

此提示只在 first turn 出現，避免每輪打擾。

## 修改本檔的規則

本檔是 **entry pointer**，不是 canonical content。修改前先讀 [`runtime/bootstrap-entry-points.yaml`](runtime/bootstrap-entry-points.yaml) §thinness_rule。新 obligation 加到 CORE_BOOTSTRAP.md（cross-tool）或 `ai-tools/agent/claude.md`（Claude-specific）— 不加到本檔。Commit-msg hook 會擋下違反 thinness 的修改。

## 修改本檔的規則

本檔是 **entry pointer**，不是 canonical content。修改前先讀 [`runtime/bootstrap-entry-points.yaml`](runtime/bootstrap-entry-points.yaml) §thinness_rule。新 obligation 加到 CORE_BOOTSTRAP.md（cross-tool）或 `ai-tools/agent/claude.md`（Claude-specific）— 不加到本檔。Commit-msg hook 會擋下違反 thinness 的修改。
