# GitHub Copilot 使用說明

本檔只記錄 GitHub Copilot 與其他 agent 工具不同的地方。通用 bootstrap 與 obligation 來源一律是 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) 與 [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)。

## Thin Entry Point

Copilot 的 repo-wide 入口是 `.github/copilot-instructions.md`。它必須保持 thin pointer，只指向：

1. [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)
2. [`runtime/core-bootstrap.yaml`](../../runtime/core-bootstrap.yaml)
3. 本 adapter 與 companion contract [`copilot.yaml`](copilot.yaml)

Scoped instructions 可放在 `.github/instructions/*.instructions.md`，但仍只能導流到 canonical bootstrap / enforcement source，不可複製中央規則正文。

## Copilot 差異

- Copilot custom instructions 是提示與上下文 surface，不是可靠 enforcement boundary；行為強制仍依 repository hooks、CI、`ai-skill runtime validate` 與 runtime gates。
- 某些 Copilot / VS Code agent 功能會依 workspace project detector 啟用，可能只支援特定 language 或 framework project。Ai-skill 這類 knowledge repo 應視為 compatibility adapter，不視為 primary governed runtime。
- `.github/copilot-instructions.md` 為 project-wide instructions；`.github/instructions/*.instructions.md` 可用 `applyTo` frontmatter 做 scoped pointer。
- Copilot 新 session 若沒有可靠讀取 bootstrap，使用 `ai-skill copilot start --project <project>` 或 project-local `.copilot/start-copilot.sh` 產生第一則 guided bootstrap prompt，貼給 Copilot 後再開始任務；不得把「只是列檔 / read-only / 說明原因」視為可跳過 bootstrap 的例外。
- `.github/prompts/*.prompt.md` 若未來加入，只能作為手動任務模板，不可取代 bootstrap contract。

## 配置邊界

Copilot-specific 路徑、instructions、prompts、agent UI 與 project detector 限制留在本檔或 `.github/` 設定。跨工具規則放回 `enforcement/`，runtime contract 放回 `runtime/core-bootstrap.yaml`。

← [回到 AI 工具索引](../README.md)
