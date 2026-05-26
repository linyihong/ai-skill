# Gemini CLI 使用說明

本檔說明 Gemini CLI 特有的配置與操作注意事項。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## Gemini CLI 配置實作

### 自動載入入口：`GEMINI.md`

Gemini CLI 啟動時會優先讀取根目錄的 `GEMINI.md`。本庫的 `GEMINI.md` 已實作一行指向 `CORE_BOOTSTRAP.md`，Gemini 啟動後會自動依啟動流程載入核心規則與 OS layout。

**設定一次 repo 即可**：只要 clone 本 repo，Gemini CLI 啟動時自動讀 `GEMINI.md` → `CORE_BOOTSTRAP.md`，不需要每次手動指定。

### 工具配置：`.gemini/`

`.gemini/` 目錄儲存 Gemini CLI 的 session 狀態、臨時檔案與工具特定設定。詳細內容由 Gemini CLI 自動管理，人類貢獻者通常不需要手動修改。

## 全域設定 vs 專案設定

Gemini CLI 的設定層級如下：

| 層級 | 範圍 | 設定位置 | 說明 |
|------|------|----------|------|
| 層級 A：全域 | 所有專案 | `~/.gemini/GEMINI.md` | 跨專案的全域偏好與啟動引導 |
| 層級 B：專案 | 單一專案 | `<PROJECT_ROOT>/GEMINI.md` | 只對該專案生效，優先級高於全域 |

### 建議策略

```
全域 ~/.gemini/GEMINI.md（層級 A）
  ├── 指向 Ai-skill 的 CORE_BOOTSTRAP.md（絕對路徑）
  └── 語言偏好與跨專案個人規則

專案 GEMINI.md（層級 B，必選）
  ├── 指向 Ai-skill 的 CORE_BOOTSTRAP.md（絕對路徑）
  └── 專案特定的啟動引導
```

## Gemini CLI 特有能力與操作

Gemini CLI 具備強大的外部工具與子代理能力，操作時需注意以下規則：

### 1. 外部檢索（Google Search & Web Fetch）

- **工具**：`google_web_search` 與 `web_fetch`
- **使用時機**：需要查詢最新文件、修復不熟悉的編譯錯誤或進行廣泛研究時。
- **規則**：檢索結果應作為證據（evidence），並依 [`enforcement/evidence-hierarchy.md`](../../enforcement/evidence-hierarchy.md) 處理。

### 2. 子代理協調（Sub-Agents）

- **工具**：`invoke_agent`
- **使用時機**：執行耗時的批次任務、高流量輸出指令或 Speculative 研究。
- **規則**：委派任務時必須提供完整的 Context，並在子代理返回後將結果整合進當前 Session 狀態。

### 3. Shell 指令安全

- **工具**：`run_shell_command`
- **規則**：執行會修改檔案系統或系統狀態的指令前，**必須**簡要說明目的。不得使用 `ask_user` 請求執行權限，應依賴工具內建的確認機制。

## Gemini CLI 與對話目標閉環

工具中立規則見 [`enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md)。Gemini CLI 會自動讀取 `.agent-goals/`，並依據 `CORE_BOOTSTRAP.md` 的義務進行目標追蹤與驗證。

**Goal ledger 操作流程已由 runtime 管理**，請參考：
- [`runtime/runtime.db`](../../runtime/runtime.db) — `phase_machine` / `obligation_ledger` / `blocking_gates`
- `ai-skill goals` CLI helper

## 語言偏好設定

Gemini CLI 的語言行為由 `GEMINI.md` 中的 Custom Instructions 控制：

1. **在 `GEMINI.md` 中設定**：
```text
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.
```

2. **語言一致性強制規則**：所有輸出（包含 tool calls 說明、技術分析、表格欄位、commit message）都必須與使用者當前語言一致。

## 驗證

使用 Gemini CLI 完成任務後，最後要求它回報：
- 讀了哪些規則與依賴。
- 認知模式（Cognitive Mode）報告是否符合 v2 格式。
- 目標是否完成，證據是否充足。
- 驗證方法：diff review、build check、readback。

← [回到 AI 工具索引](../README.md)
