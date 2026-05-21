# 貢獻與驗證入口（人類維護者）

本頁給**修改 Ai-skill 本 repository** 的維護者與 reviewer：把「要跑什麼、去哪讀權威文件」收成單一路徑。
`knowledge/indexes/` 等索引主要服務 **agent 依 task intent 路由**；本頁服務 **人類改 repo 時的 PR / 閉環流程**。

## 開始前

- 可執行政策與 writeback 閉環：[`enforcement/dependency-reading.md`](../enforcement/dependency-reading.md)、[`enforcement/linked-updates.md`](../enforcement/linked-updates.md)。
- 可重用 Markdown 的**語言與用語**：變更 `enforcement/`、`workflow/`、`analysis/`、`intelligence/`、`governance/`、根 `README.md`、根 `CONTRIBUTING.md`、模板或 onboarding 類文件時，依 [`enforcement/neutral-language.md`](../enforcement/neutral-language.md) 檢查（繁體中文正文；英文限路徑、指令、環境變數、程式符號與必要專有名詞）。
- Agent 啟動與最小上下文：[`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md)（與維護 PR 非強制，但有助對齊語彙）。

## 依變更類型該做什麼

| 你改的是 | 必讀 / 必跑（摘要） |
| --- | --- |
| `knowledge/`、`validation/` 或會影響 runtime surface | 在提交前執行 `ai-skill runtime refresh`；細節見 [`scripts/README.md`](../scripts/README.md#knowledge-runtime-validation)。 |
| `enforcement/`、`workflow/`、`analysis/`、`intelligence/`、根 `README`、模板、同步腳本 | 依 [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) 做連動更新或明列「已檢查，無需更新」。 |
| 僅文件、連結、排版 | 對 touched docs 做 Markdown link check；並依 [`enforcement/neutral-language.md`](../enforcement/neutral-language.md) 做語言與低爭議用語檢查。大改見 [`governance/validation/README.md`](validation/README.md) 的 Link check / Lints / Diff review。 |
| 新分層 / migration / 架構重構 | 依 [`governance/validation/README.md`](validation/README.md) 的 **Migration Validation Checklist** 與 [`enforcement/linked-updates.md`](../enforcement/linked-updates.md) 架構重構列。 |

## 常用指令（canonical）

在 `<AI_SKILL_REPO>` 根目錄：

```bash
# 改 knowledge / validation 後：重建 reports、SQLite index 並跑 validators
ai-skill runtime refresh

# 提交前檢查 dirty 分組（預設 dry-run）
ai-skill close-loop --dry-run
```

選用：若已設定 `git config core.hooksPath scripts/git-hooks`，在 staged 檔案觸及 `knowledge/`、`validation/` 等時，`pre-commit` 會跑 `ai-skill runtime validate`（見 [`scripts/git-hooks/pre-commit`](../scripts/git-hooks/pre-commit)）。

## 與「新專案接線」的差別

- **在別的專案掛載本知識庫**（Roo / Cursor / Claude 等）：見 [`ai-tools/new-project-onboarding.md`](../ai-tools/new-project-onboarding.md)。
- **改本 repo 的內容與治理**：以本頁與 [`governance/validation/README.md`](validation/README.md) 為準；工具專屬 sync、hook 細節見 [`ai-tools/README.md`](../ai-tools/README.md) 與各工具子文件。

## 權威索引（避免重複維護）

| 主題 | 文件 |
| --- | --- |
| 腳本與 runtime refresh 全文 | [`scripts/README.md`](../scripts/README.md) |
| Validation gates 與 checklist | [`governance/validation/README.md`](validation/README.md) |
| 中性用語與語言一致性 | [`enforcement/neutral-language.md`](../enforcement/neutral-language.md) |
| 文件拆分門檻 | [`governance/document-sizing.md`](document-sizing.md) |
| OS 目錄總覽 | [`README.md`](../README.md) |

← [回到 governance 索引](README.md)
