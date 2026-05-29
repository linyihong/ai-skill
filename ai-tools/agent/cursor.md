# Cursor 使用說明

本檔只記錄 Cursor 與其他 agent 工具不同的地方。通用 bootstrap 與 obligation 來源一律是 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) 與 [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)。

## Thin Entry Point

Cursor 的 repo 入口是 `.cursor/rules/dependency-reading.mdc`（`alwaysApply: true`）。它必須保持 thin pointer，只指向：

1. [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)
2. [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)

不要在 Cursor rule 或本 adapter 複製 Bootstrap Receipt、Cognitive Mode、close-loop、goal ledger 或 runtime phase 的細節；需要時從 canonical source 讀取。

## Cursor 差異

- `.cursor/rules/*.mdc` 可放全域或專案規則；全域規則通常位於 `~/.cursor/rules/`，專案規則位於 `<PROJECT_ROOT>/.cursor/rules/`。
- 全域與專案 rules 會合併，不是覆蓋關係；專案 overlay 仍應保持薄，只補專案特定約束。
- `.cursor/hooks.json` 與 hooks 只應作為提醒或檢查；除非團隊明確配置，不要讓 hook 靜默修改檔案、關閉 goal 或自動 commit。
- `ai-skill init-project --tools cursor` 產生的 project hook 必須包含 `stop` close-out check，呼叫 repo-local `ai-skill hooks run stop`。Cursor stop 不能靠 exit 2 loop back；缺少 Bootstrap Receipt、compact `Cognitive:` / full `### Cognitive Mode 報告`，或必要的 `### Project Git Report` 時，hook runner 必須一次彙整缺項，輸出 `followup_message` 並 exit 0。這是 chat/session 層防漏；格式與枚舉仍只在 canonical bootstrap sources。
- 多資料夾工作區可同時打開業務專案與本 repository，讓 agent 直接讀 canonical source。

## 配置邊界

Cursor-specific 路徑、hooks、UI 行為與 workspace 操作留在本檔或 `.cursor/` 設定。跨工具規則放回 `enforcement/`，runtime contract 放回 `runtime/core-bootstrap.yaml`。

← [回到 AI 工具索引](../README.md)
