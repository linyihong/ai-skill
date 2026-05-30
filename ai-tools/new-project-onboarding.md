# 新專案初始化流程

當你開一個**全新的專案**，想要讓它使用 Ai-skill 知識庫時，依以下流程設定。

## 流程總覽

```
開新專案 → 執行 `ai-skill init-project` → 在專案中開啟 AI 工具 → 開始開發
```

Go CLI command `ai-skill init-project` 會一次設定所有支援的 AI 工具：

| 工具 | 產出檔案 | 效果 |
|------|---------|------|
| **Roo Code** | `.roomodes` | thin bootstrap pointer，載入 bootstrap canonical contract |
| **Cursor** | `.cursor/rules/ai-skill-bootstrap.mdc` + `.cursor/hooks.json` | thin alwaysApply pointer；session 啟動提示 bootstrap，stop hook 以 `followup_message` loop back 檢查 Bootstrap Receipt 與 final Cognitive Mode 報告 |
| **Claude Code** | `CLAUDE.md` + `.claude/settings.json` | thin bootstrap pointer + Go hook runner，載入 bootstrap canonical contract 並強制 final Cognitive Mode / nested Git report |
| **Gemini CLI** | `GEMINI.md` | thin bootstrap pointer，載入 bootstrap canonical contract |
| **Codex / generic agents** | `AGENTS.md` | thin generic entry，透過 `ai-tools/README.md` 路由到 active adapter |
| **GitHub Copilot** | `.github/copilot-instructions.md` + `.github/instructions/ai-skill.instructions.md` + `.copilot/` guided startup package | compatibility thin pointer；新 session 用 guided prompt 導流到 canonical source，hard enforcement 由 hooks / CI / runtime validate 負責 |

初始化檔案不得寫入本機 Ai-skill repo 的絕對路徑；project-local bootstrap files 使用 `<AI_SKILL_REPO>` placeholder，Claude hook wiring 使用 `AI_SKILL_REPO` 環境變數（未設定時才嘗試 `$HOME/Documents/Ai-skill` 等不含使用者名稱的 fallback）。

## 快速開始

```bash
# 建議先設定一次，避免 project files 寫入任何本機絕對路徑
export AI_SKILL_REPO=/path/to/ai-skill

# 1. 開新專案
mkdir -p ~/projects/my-new-app

# 2. 執行初始化腳本（從 Ai-skill repo 目錄）
cd /path/to/ai-skill
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app

# 3. 在目標專案中開啟 AI 工具，開始開發
cd ~/projects/my-new-app
code .   # 或 cursor . 或 claude .
```

## 本機環境變數

`AI_SKILL_REPO` 是 project-local bootstrap files 連回 Ai-skill knowledge repo 的唯一約定。它是每台機器自己的設定，不應寫入會 commit 的專案檔或 `.env`。

`init-project` 會自動建立：

- `.ai-skill/local.env`：本機設定，內容包含 `AI_SKILL_REPO=<目前 Ai-skill repo>`，權限為 `0600`
- `.ai-skill/.gitignore`：只忽略 `local.env`，避免本機路徑被 commit，同時允許提交 `.ai-skill/project/` project overlay
- `.ai-skill/project/README.md`：project-local overlay index；專案特有規則應放在 `.ai-skill/project/rules/`，再由工具入口薄薄指向

Claude hooks 會先讀 process environment；如果沒有 `AI_SKILL_REPO`，會 source `.ai-skill/local.env`，所以剛初始化後不用重開 shell 也能連回 Ai-skill。

```bash
# macOS / Linux / WSL
export AI_SKILL_REPO=/path/to/ai-skill
```

```powershell
# Windows PowerShell
[Environment]::SetEnvironmentVariable("AI_SKILL_REPO", "C:\path\to\Ai-skill", "User")
```

長期建議做法：macOS / Linux / WSL 放在 shell profile（例如 `~/.zshrc`、`~/.bashrc`）；Windows 放在 User environment variable。Project files 只保留 `<AI_SKILL_REPO>` placeholder，讓同一份 repo 可以在不同 OS / 使用者路徑上共用。

## 進階用法

### 只設定特定工具

```bash
# 只設定 Roo Code
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --tools roo

# 只設定 Cursor + Claude Code
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --tools cursor,claude

# 只設定 GitHub Copilot compatibility adapter
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project ~/projects/my-new-app --tools copilot
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

建立 tool-local modes，每個 mode 的 `customInstructions` 只保留 bootstrap pointer：

1. **Bootstrap companion**：指向 `<AI_SKILL_REPO>/CORE_BOOTSTRAP.md`
2. **Canonical contract**：指向 `<AI_SKILL_REPO>/runtime/core-bootstrap.yaml`

語言、Cognitive Mode、knowledge-update 與 close-loop obligations 由 bootstrap contract 載入，不複製到 `.roomodes`。

> **注意**：`.roomodes` 是專案層級設定，會**完全覆蓋**全域 Custom Instructions（不會合併）。如果之前有設定全域 Custom Instructions，`.roomodes` 中的內容會取代它。

### Cursor（`.cursor/rules/ai-skill-bootstrap.mdc`）

建立 thin alwaysApply pointer，包含：

1. **Bootstrap companion**：指向 `CORE_BOOTSTRAP.md`
2. **Canonical contract**：指向 `runtime/core-bootstrap.yaml`

Bootstrap obligations、format、enum、examples、required reads 與 runtime DB 載入不複製在 Cursor rule；由 `runtime/core-bootstrap.yaml` 管理。

若專案有自己的檢查規則，不要把規則本體放進 `.cursor/rules`。將規則本體放在 `.ai-skill/project/rules/`，`.cursor/rules/*.mdc` 只保留 thin pointer，這樣 Claude Code、AGENTS.md-aware tools、Gemini CLI、Roo Code 也能讀到同一份 project overlay。

同時建立 `.cursor/hooks.json`：

- `sessionStart`：提示載入 Ai-skill canonical bootstrap。
- `stop`：呼叫 repo-local `ai-skill hooks run stop`；若對話缺 Bootstrap Receipt、final response 缺 compact `Cognitive:` / full `### Cognitive Mode 報告`，或需要但未附 `### Project Git Report`，hook runner 一次彙整缺項，輸出 `followup_message` 並 exit 0 讓 Cursor loop back。Cognitive 格式與枚舉仍只在 canonical bootstrap sources，不複製到 Cursor rule。

### Claude Code（`CLAUDE.md` + `.claude/settings.json`）

建立 thin bootstrap pointer，包含：

1. **Bootstrap companion**：指向 `CORE_BOOTSTRAP.md`
2. **Canonical contract**：指向 `runtime/core-bootstrap.yaml`

工具入口不複製 obligations；Claude Code 依 bootstrap contract 載入後續規則。`.claude/settings.json` 只保存 hook wiring，呼叫 Ai-skill repo-local Go binary：

- `SessionStart`：注入 Bootstrap Receipt context。
- `UserPromptSubmit`：注入 final close-out Cognitive Mode reminder，並掃描 project root 底下的 dirty nested Git repos。
- `Stop`：檢查 final response 是否含 Cognitive Mode；若有 dirty root / nested repo，要求合併 `### Project Git Report`。

`.claude/settings.json` 不寫本機絕對路徑。Hook command 會讀 `AI_SKILL_REPO`；未設定時只嘗試 `$HOME/Documents/Ai-skill` / `$HOME/Ai-skill` / `$PWD/../Ai-skill` 這類 portable fallback，找不到則 fail-open 並提示設定環境變數。

### GitHub Copilot（`.github/` instructions + `.copilot/` guided startup）

建立 Copilot project-wide custom instructions、scoped instructions 與 guided startup package，全部保持 thin pointer：

1. **Project-wide entry**：`.github/copilot-instructions.md` 指向 `<AI_SKILL_REPO>/CORE_BOOTSTRAP.md`、`<AI_SKILL_REPO>/runtime/core-bootstrap.yaml` 與 Copilot adapter。
2. **Scoped entry**：`.github/instructions/ai-skill.instructions.md` 使用 `applyTo: "**"`，讓 Copilot 在可載入 instructions files 的 surface 上回到同一組 canonical source。
3. **Guided startup**：`.copilot/bootstrap-prompt.md` 與 `ai-skill copilot start --project <project>` 產生新 session 第一則訊息；`.copilot/start-copilot.sh` 只是 temporary thin wrapper，呼叫 repo-local Go CLI 並寫明 deletion condition。
4. **Enforcement boundary**：Copilot instructions 只導流；若 Copilot 無法可靠強制 runtime gate，hard enforcement 仍交給 repository hooks、CI 與 `ai-skill runtime validate`。

> **限制**：部分 Copilot / VS Code agent 功能會依 workspace project detector 啟用，可能只支援特定 language 或 framework project。對文件型或 knowledge repo，Copilot 應視為 compatibility adapter，不視為 primary governed runtime。

## 發現機制：新專案如何知道 Ai-skill？

新專案預設**不知道** Ai-skill 存在。以下是三種讓它知道的層級：

### 層級 A：全域設定（最強，所有專案自動生效）

| 工具 | 方式 | 狀態 |
|------|------|------|
| **Roo Code** | 全域 Custom Instructions（VS Code SQLite `state.vscdb`）→ 指向 `CORE_BOOTSTRAP.md` | 需關閉 VS Code 後執行 `ai-skill roo set-global-custom-instructions` |
| **Cursor** | User Rules（`~/.cursor/rules/`）→ alwaysApply `CORE_BOOTSTRAP.md` | 手動設定一次，所有專案生效 |
| **Claude Code** | 無全域機制 | 只能靠專案 `CLAUDE.md` |
| **GitHub Copilot** | VS Code / Copilot custom instructions 或 user instructions locations | 可作為輔助，但本 repo 只保證 project-level thin pointer |

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
如果沒有，使用此全域規則載入 Core Bootstrap companion 與 canonical contract：

[`<AI_SKILL_REPO>/CORE_BOOTSTRAP.md`](<AI_SKILL_REPO>/CORE_BOOTSTRAP.md)
[`<AI_SKILL_REPO>/runtime/core-bootstrap.yaml`](<AI_SKILL_REPO>/runtime/core-bootstrap.yaml)
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
| `.github/copilot-instructions.md` | `.github/` | 有 `CORE_BOOTSTRAP.md` 參考 |
| `.github/instructions/ai-skill.instructions.md` | `.github/instructions/` | 有 scoped thin pointer |
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

# 檢查 GitHub Copilot 設定
ls -la /path/to/project/.github/copilot-instructions.md
ls -la /path/to/project/.github/instructions/ai-skill.instructions.md
ls -la /path/to/project/.copilot/bootstrap-prompt.md
/path/to/ai-skill/scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 copilot start --project /path/to/project

# 檢查對話目標目錄
ls -la /path/to/project/.agent-goals/README.md
```

如果所有檔案都存在，設定完成。

## 路徑變更處理

如果 Ai-skill repo 移動位置，重新執行初始化腳本即可更新所有路徑：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 init-project --project /path/to/project --force
```

初始化檔只保留 `<AI_SKILL_REPO>` / `AI_SKILL_REPO` 參照；若 Ai-skill repo 移動位置，更新本機環境變數即可，不需要把新絕對路徑寫進 project files。

## 與既有文件的關係

- [`scripts/ai-skill-cli/`](../scripts/ai-skill-cli/README.md) — `ai-skill init-project` 的 Go CLI source
- [`ai-skill roo set-global-custom-instructions`](../ai-skill roo set-global-custom-instructions) — Roo Code 全域設定寫入腳本
- `ai-skill sync-cursor-bundle` — Cursor bundle 同步的 Go CLI 入口；目前只有 dry-run planner，write mode 待實作
- [`ai-tools/agent-onboarding.md`](agent-onboarding.md) — 新 AI 工具 onboarding（不是新專案）
- [`ai-tools/agent/roo.md`](agent/roo.md) — Roo Code 使用說明
- [`ai-tools/agent/cursor.md`](agent/cursor.md) — Cursor 使用說明
- [`ai-tools/agent/claude.md`](agent/claude.md) — Claude Code 使用說明
- [`ai-tools/agent/gemini-cli.md`](agent/gemini-cli.md) — Gemini CLI 使用說明
- [`ai-tools/agent/copilot.md`](agent/copilot.md) — GitHub Copilot compatibility adapter 使用說明
- [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) — 最小啟動集合

← [回到 AI 工具索引](README.md)
# Runtime Projection

本 new project onboarding 是可執行流程，companion YAML 為 [`new-project-onboarding.yaml`](new-project-onboarding.yaml)。

更新 `ai-skill init-project`、project-level bootstrap files、工具清單 or `ai-tools/agent/<tool>.md` 並影響新專案初始化時，必須同步檢查本文件與 YAML activation contract。

**重要義務**：當新增支援的 AI 工具時，**必須同時更新根目錄的 `README.md`** 中的「Works With」區塊，確保系統能力透明。

目前 `ai-skill init-project` 支援的工具：

| 工具 | Project bootstrap file | Canonical adapter |
| --- | --- | --- |
| Roo Code | `.roomodes` | [`agent/roo.md`](agent/roo.md) |
| Cursor | `.cursor/rules/ai-skill-bootstrap.mdc`、`.cursor/hooks.json` | [`agent/cursor.md`](agent/cursor.md) |
| Claude Code | `CLAUDE.md`、`.claude/settings.json` | [`agent/claude.md`](agent/claude.md) |
| Gemini CLI | `GEMINI.md` | [`agent/gemini-cli.md`](agent/gemini-cli.md) |
| Codex | `AGENTS.md` | [`agent/codex.md`](agent/codex.md) |
| GitHub Copilot | `.github/copilot-instructions.md`、`.github/instructions/ai-skill.instructions.md`、`.copilot/` guided startup package | [`agent/copilot.md`](agent/copilot.md) |
