# 新專案初始化流程

當你開一個**全新的專案**，想要讓它使用 Ai-skill 知識庫時，依以下流程設定。

## 流程總覽

```
開新專案 → 執行 `ai-skill init-project` → 在專案中開啟 AI 工具 → 開始開發
```

Go CLI command `ai-skill init-project` 會一次設定所有支援的 AI 工具：

| 工具 | 產出檔案 | 效果 |
|------|---------|------|
| **Roo Code** | `.roomodes` | 5 個 mode 都含語言規則 + 知識更新 checkpoint |
| **Cursor** | `.cursor/rules/ai-skill-bootstrap.mdc` + `.cursor/hooks.json` | alwaysApply 規則，session 啟動時自動載入 |
| **Claude Code** | `CLAUDE.md` | 自動載入 Core Bootstrap 啟動流程 |

## 快速開始

```bash
# 1. 開新專案
mkdir -p ~/projects/my-new-app

# 2. 執行初始化腳本（從 Ai-skill repo 目錄）
cd /path/to/ai-skill
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app

# 3. 在目標專案中開啟 AI 工具，開始開發
cd ~/projects/my-new-app
code .   # 或 cursor . 或 claude .
```

## 進階用法

### 只設定特定工具

```bash
# 只設定 Roo Code
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --tools roo

# 只設定 Cursor + Claude Code
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --tools cursor,claude
```

### 預覽模式（不實際寫入）

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --dry-run
```

### 覆蓋已有設定檔

```bash
# 如果目標專案已有 .roomodes 或 CLAUDE.md，用 --force 覆蓋
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --force
```

## 各工具設定說明

### Roo Code（`.roomodes`）

建立 5 個 mode（code / architect / ask / debug / orchestrator），每個 mode 的 `customInstructions` 包含：

1. **Core Bootstrap 啟動**：指向 `<AI_SKILL_REPO>/CORE_BOOTSTRAP.md`
2. **語言強制規則**：無預設語言，強制跟隨使用者語言
3. **知識更新流程 Checkpoint**：每輪結束前檢查是否有新知識要回饋
4. **專案 durable Markdown 預設**：在本專案新增或大幅改寫 `docs/`、`wiki/`、`README.md`、ADR、runbook 等長期保存的 Markdown 時，動筆前先讀 `<AI_SKILL_REPO>/workflow/documentation/README.md` 與 `execution-flow.md`（無須使用者先說「docs/wiki/token」等關鍵字）

> **注意**：`.roomodes` 是專案層級設定，會**完全覆蓋**全域 Custom Instructions（不會合併）。如果之前有設定全域 Custom Instructions，`.roomodes` 中的內容會取代它。

### Cursor（`.cursor/rules/ai-skill-bootstrap.mdc`）

建立 alwaysApply 規則，包含：

1. **啟動流程**：指向 `CORE_BOOTSTRAP.md`
2. **專案 durable Markdown 預設**：寫入 `docs/`、`wiki/`、`README.md`、ADR、runbook 等長期保存的 Markdown 前，先讀 `workflow/documentation/README.md` 與 `execution-flow.md`（無須使用者先說關鍵字）
3. **語言強制規則**
4. **知識更新流程 Checkpoint**

同時建立 `.cursor/hooks.json`，在 session 啟動時提醒載入 Ai-skill。

### Claude Code（`CLAUDE.md`）

建立自動載入檔案，包含：

1. **啟動流程**：指向 `CORE_BOOTSTRAP.md`
2. **專案 durable Markdown 預設**：寫入 `docs/`、`wiki/`、`README.md`、ADR、runbook 等長期保存的 Markdown 前，先讀 `workflow/documentation/README.md` 與 `execution-flow.md`（無須使用者先說關鍵字）
3. **語言強制規則**
4. **知識更新流程 Checkpoint**

## 發現機制：新專案如何知道 Ai-skill？

新專案預設**不知道** Ai-skill 存在。以下是三種讓它知道的層級：

### 層級 A：全域設定（最強，所有專案自動生效）

| 工具 | 方式 | 狀態 |
|------|------|------|
| **Roo Code** | 全域 Custom Instructions（VS Code SQLite `state.vscdb`）→ 指向 `CORE_BOOTSTRAP.md` | 需關閉 VS Code 後執行 `ai-skill roo set-global-custom-instructions` |
| **Cursor** | User Rules（`~/.cursor/rules/`）→ alwaysApply `CORE_BOOTSTRAP.md` | 手動設定一次，所有專案生效 |
| **Claude Code** | 無全域機制 | 只能靠專案 `CLAUDE.md` |

**設定方式**：

```bash
# Roo Code 全域設定（需先關閉 VS Code）
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 roo set-global-custom-instructions

# Cursor 全域設定
mkdir -p ~/.cursor/rules
# 在 ~/.cursor/rules/ 下建立 ai-skill-global.mdc，內容參考下方
```

**Cursor 全域規則範本**（`~/.cursor/rules/ai-skill-global.mdc`）：

```markdown
---
description: Ai-skill 知識庫全域啟動
globs: 
alwaysApply: true
---

# Ai-skill Global Bootstrap

每個專案啟動時，檢查是否有專案層級的 `.cursor/rules/ai-skill-bootstrap.mdc`。
如果沒有，使用此全域規則載入 Core Bootstrap：

[`<AI_SKILL_REPO>/CORE_BOOTSTRAP.md`](<AI_SKILL_REPO>/CORE_BOOTSTRAP.md)
```

### 層級 B：專案初始化 CLI（中等，開新專案時跑一次）

就是 `ai-skill init-project`。

每次開新專案時執行一次，所有工具設定一次到位。

### 層級 C：文件引導（被動，需要有人告訴你）

如果你拿到一個已經設定好的專案，或想知道某個專案是否使用 Ai-skill，可以檢查：

| 檢查項目 | 位置 | 表示有設定 |
|---------|------|-----------|
| `.roomodes` | 專案根目錄 | 有 `CORE_BOOTSTRAP.md` 參考 |
| `.cursor/rules/ai-skill-bootstrap.mdc` | `.cursor/rules/` | 檔案存在 |
| `CLAUDE.md` | 專案根目錄 | 有 `CORE_BOOTSTRAP.md` 參考 |
| `.agent-goals/` | 專案根目錄 | 目錄存在 |

## 驗證設定是否正確

執行初始化後，檢查目標專案：

```bash
# 檢查 Roo Code 設定
ls -la /path/to/project/.roomodes

# 檢查 Cursor 設定
ls -la /path/to/project/.cursor/rules/ai-skill-bootstrap.mdc

# 檢查 Claude Code 設定
ls -la /path/to/project/CLAUDE.md

# 檢查對話目標目錄
ls -la /path/to/project/.agent-goals/README.md
```

如果所有檔案都存在，設定完成。

## 路徑變更處理

如果 Ai-skill repo 移動位置，重新執行初始化腳本即可更新所有路徑：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project /path/to/project --force
```

這會用新的絕對路徑覆蓋所有設定檔。

## 與既有文件的關係

- [`scripts/ai-skill-cli/`](../scripts/ai-skill-cli/README.md) — `ai-skill init-project` 的 Go CLI source
- [`ai-skill roo set-global-custom-instructions`](../ai-skill roo set-global-custom-instructions) — Roo Code 全域設定寫入腳本
- `ai-skill sync-cursor-bundle` — Cursor bundle 同步的 Go CLI 入口；目前只有 dry-run planner，write mode 待實作
- [`ai-tools/agent-onboarding.md`](agent-onboarding.md) — 新 AI 工具 onboarding（不是新專案）
- [`ai-tools/agent/roo.md`](agent/roo.md) — Roo Code 使用說明
- [`ai-tools/agent/cursor.md`](agent/cursor.md) — Cursor 使用說明
- [`ai-tools/agent/claude.md`](agent/claude.md) — Claude Code 使用說明
- [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) — 最小啟動集合

← [回到 AI 工具索引](README.md)
