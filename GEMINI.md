# Gemini CLI Bootstrap Entry

> **IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
>
> 本檔為 **thin tool-entry pointer**。所有核心規則與啟動義務的 canonical 來源在 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md)。

## 啟動序列

1. **載入核心規則**：讀取 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) (3 rules, ~800 tokens)。
2. **認證狀態**：輸出 **Bootstrap Receipt**，確認 phase 為 `phase.bootstrap`。
3. **報告模式**：輸出 **Cognitive Mode 報告 (v2)**，記錄 discovery signals。
4. **追蹤目標**：讀取 [`.agent-goals/README.md`](.agent-goals/README.md) 同步當前任務。

## 本專案強制回覆規則

本檔不是完整規則，只是 bootstrap pointer。Gemini 必須在處理任務前完成：

1. 使用本 repo root 作為 Ai-skill repo。
2. 讀本 repo 的 `CORE_BOOTSTRAP.md`。
3. 讀本 repo 的 `runtime/core-bootstrap.yaml`，以 YAML canonical contract 為準。
4. 每個 user-facing response 結尾都必須依 active per-turn obligations 附上 Cognitive Mode reporting。

不能只停在本檔摘要；必須 dereference pointer 到 canonical files。

## 特殊能力規範

Gemini CLI 具備外部檢索與子代理能力，詳細規範見 [`ai-tools/agent/gemini-cli.md`](ai-tools/agent/gemini-cli.md)。

## 語言一致性

Language Preference: Default to English, but always match the user's language in conversation.
語言一致性強制規則：所有輸出（含分析、表格與 commit message）必須與使用者語言一致。

---
修改規則見 [`ai-tools/agent/gemini-cli.md`](ai-tools/agent/gemini-cli.md)。
