# AI 工具使用說明

本目錄放各種 AI coding / agent 工具如何讀取與同步本知識庫。新增其他工具時，在這裡加新的子檔，例如 CLI agent、IDE agent、雲端 agent 或其他支援 skill/rules 的工具。

共用原則：**中央庫是真相來源**，工具端只做參照、symlink 或同步快照。

通用 shared rules、skills、templates 與根 README 應保持工具中立；工具專屬路徑、hook、UI、reload、設定與同步命令放在本目錄對應工具文件中。

若某個 skill 對某工具有必要的特殊執行策略，skill 內可用 `skills/<skill>/tool-adapters/<tool>.md` 記錄差異；本目錄仍只放該工具的全域設定、同步與操作方式。

| 工具 | 文件 | 用途 |
| --- | --- | --- |
| Claude | [claude.md](claude.md) | Claude 類工具如何明確讀取 shared rules、skill 入口、依賴文件、goal ledger 與 Ai-skill writeback 流程。 |
| Cursor | [cursor.md](cursor.md) | Cursor 如何啟用 apk-analysis、參照或同步 `.cursor`、維持中央庫一致性，並以工具中立 `.agent-goals/` 做對話目標閉環提醒。 |

← [回到根目錄](../README.md)
