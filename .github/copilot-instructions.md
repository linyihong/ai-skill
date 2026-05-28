# Copilot instructions (minimal)

請在啟動前閱讀：
- CORE_BOOTSTRAP.md
- runtime/core-bootstrap.yaml
- runtime/runtime.db

最小啟動步驟：
1. 讀取 bootstrap 合約（Copilot 在 repository 內會自動讀取 .github/copilot-instructions.md）：
   ! ai-skill runtime obligations
2. 在第一個非-讀取工具前輸出 Bootstrap Receipt，範例：
   Bootstrap: rules=✓ phase=<phase-id> obligations=<comma-separated-obligation-ids> gates=<n>

參考檔案：
- copilot.md

備註：避免在檔案中使用絕對路徑或機敏資訊；使用 <repo-root> 占位並確保必要路徑已授權（/add-dir）。
