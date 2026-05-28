# Copilot agent adapter (thin pointer)

用途
- 為 Copilot 提供專案啟動時的 thin adapter 指標（不要複製中央規則）。

必讀（不重複）
- CORE_BOOTSTRAP.md
- runtime/core-bootstrap.yaml

使用說明
- adapter contract 在 ai-tools/agent/copilot.yaml
- Copilot 啟動時會讀取此檔，請勿在此放置絕對路徑或機敏資訊；使用 <AI_SKILL_REPO> 或 <repo-root> 佔位。

建議流程
- 將可執行或啟動腳本放在 scripts/ 或記錄 start_command 到 copilot.yaml
- 若需本機路徑，使用 `.ai-skill/local.env`（被 .gitignore 忽略）作為機器本地橋接

如需更多細節，參考 ai-tools/new-project-onboarding.yaml 的 render_project_bootstrap_files 與 apply_tool_specific_adapters 步驟。
