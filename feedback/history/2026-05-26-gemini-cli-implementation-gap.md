# Feedback Lesson: 2026-05-26 Gemini CLI Implementation Gap

**ID**: `2026-05-26-gemini-cli-implementation-gap`
**類別**: `governance-failure`, `implementation-gap`
**嚴重程度**: `High`

## 事件描述

在執行 Gemini CLI 的 Agent Onboarding 時，Agent（Gemini CLI）僅完成了 Markdown 文件與 YAML 契約的撰寫，卻漏掉了對 Go CLI 源碼（`scripts/ai-skill-cli/`）的同步修改。這導致雖然文件宣稱支援 Gemini 的專案初始化，但實際工具鏈卻無法執行。

## 根因分析

1.  **認知孤島 (Cognitive Silos)**：Agent 將 `ai-tools/` (知識定義層) 與 `scripts/` (代碼實作層) 視為獨立的上下文，在執行 Onboarding 任務時未能觸發跨層級的依賴檢查。
2.  **檢核點失效 (Broken Gate)**：`agent-onboarding.yaml` 中的 CLI 同步步驟被視為「建議性」而非「阻斷性」，且 `runtime validate` 缺乏對代碼實作完整性的自動化檢驗。
3.  **影子實作 (Shadow Implementation)**：Agent 產生了「文檔即完成」的偏見，忽略了必須透過實作（Implementation）才能閉環的義務。

## 修正措施

1.  **代碼補強**：修改 `scripts/ai-skill-cli/internal/app/init_project.go` 支援 `gemini` 選項。
2.  **規則補強**：在 `enforcement/linked-updates.yaml` 中明確加入「新工具適配器必須同步 CLI 實作」的觸發條件與 Gate。
3.  **本回饋紀錄**：將此案例納入歷史紀錄以供後續 Agent 檢索。

## 給未來 Agent 的建議

- **不要只讀 Markdown**：當看到 `ai-tools/agent/` 下有新工具時，立即檢查 `scripts/ai-skill-cli/` 是否已同步。
- **實作優先**：在宣稱 Onboarding 完成前，必須嘗試編譯 CLI 並執行 `--dry-run` 驗證新工具。
- **檢查 Linked Updates**：嚴格遵守 [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) 的跨目錄更新要求。
