#!/usr/bin/env bash
# init-new-project.sh — 新專案初始化腳本
#
# 在目標專案目錄中建立所有 AI 工具的設定檔，指向 Ai-skill 知識庫。
# 一次設定 Roo Code + Cursor + Claude Code，讓新專案立即擁有：
#   - 語言強制規則（跟隨使用者語言）
#   - 知識更新流程 Checkpoint
#   - Core Bootstrap 啟動流程
#   - 專案 durable Markdown 預設走 workflow/documentation/（無須使用者先說關鍵字）
#
# 用法：
#   ./scripts/init-new-project.sh /path/to/new-project
#   ./scripts/init-new-project.sh /path/to/new-project --dry-run    # 預覽不寫入
#   ./scripts/init-new-project.sh /path/to/new-project --tools roo  # 只設定 Roo Code
#
# 安全條件：
#   - 目標目錄必須存在
#   - 不會覆蓋已有檔案（除非 --force）
#   - 所有路徑使用 <AI_SKILL_REPO> 占位符，腳本自動填入實際絕對路徑

set -euo pipefail

# ─── 設定 ────────────────────────────────────────────────────────────────
AI_SKILL_REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DRY_RUN=false
FORCE=false
TOOLS="roo,cursor,claude"  # 預設全部

# ─── 顏色 ─────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

info()  { echo -e "${CYAN}[INFO]${NC} $*"; }
ok()    { echo -e "${GREEN}[OK]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
err()   { echo -e "${RED}[ERROR]${NC} $*" >&2; }

# ─── 參數解析 ────────────────────────────────────────────────────────────
usage() {
  cat <<EOF
用法: $(basename "$0") <target-project-dir> [選項]

在目標專案目錄中建立 AI 工具設定檔，指向 Ai-skill 知識庫。

參數:
  <target-project-dir>   目標專案目錄（必須已存在）

選項:
  --dry-run              只預覽要建立的檔案，不實際寫入
  --force                覆蓋已有檔案（預設會跳過）
  --tools <list>         要設定的工具，逗號分隔（roo,cursor,claude）
                         預設：roo,cursor,claude
  -h, --help             顯示此說明

範例:
  $(basename "$0") ~/projects/my-new-app
  $(basename "$0") ~/projects/my-new-app --dry-run
  $(basename "$0") ~/projects/my-new-app --tools roo,cursor
EOF
  exit 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help) usage ;;
    --dry-run) DRY_RUN=true; shift ;;
    --force) FORCE=true; shift ;;
    --tools)
      shift
      TOOLS="$1"
      shift
      ;;
    -*)
      err "未知選項: $1"
      usage
      ;;
    *)
      TARGET_DIR="$1"
      shift
      ;;
  esac
done

# ─── 驗證 ─────────────────────────────────────────────────────────────────
if [[ -z "${TARGET_DIR:-}" ]]; then
  err "請指定目標專案目錄"
  usage
fi

TARGET_DIR="$(cd "$TARGET_DIR" 2>/dev/null && pwd)" || {
  err "目錄不存在: $TARGET_DIR"
  exit 1
}

if [[ "$TARGET_DIR" == "$AI_SKILL_REPO" ]]; then
  err "目標目錄不能是 Ai-skill 本體"
  exit 1
fi

# 解析工具列表
IFS=',' read -ra TOOL_LIST <<< "$TOOLS"
declare -A TOOL_MAP
for t in "${TOOL_LIST[@]}"; do
  t="$(echo "$t" | xargs)"  # trim
  case "$t" in
    roo|cursor|claude) TOOL_MAP["$t"]=1 ;;
    *) warn "忽略未知工具: $t（支援: roo, cursor, claude）" ;;
  esac
done

if [[ ${#TOOL_MAP[@]} -eq 0 ]]; then
  err "沒有有效的工具可設定"
  exit 1
fi

info "目標專案: $TARGET_DIR"
info "Ai-skill 路徑: $AI_SKILL_REPO"
info "要設定的工具: ${!TOOL_MAP[*]}"
[[ "$DRY_RUN" == true ]] && warn "DRY RUN 模式 — 不會寫入任何檔案"
echo ""

# ─── 輔助函數 ────────────────────────────────────────────────────────────

# 寫入檔案，支援 dry-run 與 force 檢查
write_file() {
  local filepath="$1"
  local content="$2"
  local desc="$3"

  if [[ -f "$filepath" && "$FORCE" != true ]]; then
    warn "跳過（已存在）: $filepath"
    warn "  使用 --force 覆蓋"
    return 0
  fi

  if [[ "$DRY_RUN" == true ]]; then
    info "[DRY RUN] 將建立 $desc: $filepath"
    return 0
  fi

  mkdir -p "$(dirname "$filepath")"
  echo "$content" > "$filepath"
  ok "已建立 $desc: $filepath"
}

# ─── 1. Roo Code 設定 ────────────────────────────────────────────────────
if [[ -n "${TOOL_MAP[roo]:-}" ]]; then
  info "── Roo Code ──"

  ROOMODES_FILE="$TARGET_DIR/.roomodes"
  ROOMODES_CONTENT=$(cat <<ROOMODES_EOF
{
  "customModes": [
    {
      "slug": "code",
      "name": "💻 Code",
      "roleDefinition": "You are Roo, a highly skilled software engineer with extensive knowledge in many programming languages, frameworks, design patterns, and best practices.",
      "customInstructions": "依 ${AI_SKILL_REPO}/CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。\n\n## 語言強制規則（最高優先級）\n\n### 核心原則\n- **強制跟隨使用者語言**：使用者用什麼語言提問，所有輸出就必須用什麼語言。\n- **無預設語言**：不存在「Default to English」。使用者第一次提問的語言即為本次對話語言。\n- **禁止自行切換**：即使分析技術內容、寫 commit message、產生表格，也必須與使用者當前語言一致。\n- **attempt_completion 前必須檢查**：確認結果的語言與使用者最後一次提問的語言一致。\n\n### 防漂移機制\n- 如果使用者用中文，所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解）都必須使用中文。\n- 如果使用者用日文，所有輸出都必須使用日文。\n- 如果使用者切換語言，跟隨其切換。\n- **禁止出現「Default to English」思維**：不存在預設語言，只有「使用者當前語言」。\n\n## 知識更新流程 Checkpoint\n\n每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：\n\n1. 讀取 ${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md 了解完整流程。\n2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？\n3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：\n   - Step 1-2：觸發檢查 + 分類知識類型\n   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）\n   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 ${AI_SKILL_REPO}/enforcement/sanitization.md 去敏）\n   - Step 5：更新目標層\n   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning\n   - Step 8：執行 Linked Updates\n   - Step 9：更新 Runtime Surfaces\n   - Step 10：驗證（diff review、去敏檢查、link check）\n   - Step 11：Commit / Push / Readback（關閉 writeback transaction）\n4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。\n\n## 專案 durable Markdown 預設（Ai-skill）\n\n當任務為在本專案新增或大幅改寫長期保存的 Markdown（例如 docs/、wiki/、README.md、ADR、runbook 等，以本 repo 慣例為準）時，動筆前先讀 ${AI_SKILL_REPO}/workflow/documentation/README.md，再讀 ${AI_SKILL_REPO}/workflow/documentation/execution-flow.md。拆分與篇幅見 ${AI_SKILL_REPO}/governance/document-sizing.md。無須使用者先說「docs/wiki/token」等關鍵字；此為上述文件類型之預設流程。",
      "groups": ["read", "edit", "command", "mcp"]
    },
    {
      "slug": "architect",
      "name": "🏗️ Architect",
      "roleDefinition": "You are Roo, an expert software architect specializing in system design, technical planning, and breaking down complex problems.",
      "customInstructions": "依 ${AI_SKILL_REPO}/CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。\n\n## 語言強制規則（最高優先級）\n\n### 核心原則\n- **強制跟隨使用者語言**：使用者用什麼語言提問，所有輸出就必須用什麼語言。\n- **無預設語言**：不存在「Default to English」。使用者第一次提問的語言即為本次對話語言。\n- **禁止自行切換**：即使分析技術內容、寫 commit message、產生表格，也必須與使用者當前語言一致。\n- **attempt_completion 前必須檢查**：確認結果的語言與使用者最後一次提問的語言一致。\n\n### 防漂移機制\n- 如果使用者用中文，所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解）都必須使用中文。\n- 如果使用者用日文，所有輸出都必須使用日文。\n- 如果使用者切換語言，跟隨其切換。\n- **禁止出現「Default to English」思維**：不存在預設語言，只有「使用者當前語言」。\n\n## 知識更新流程 Checkpoint\n\n每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：\n\n1. 讀取 ${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md 了解完整流程。\n2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？\n3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：\n   - Step 1-2：觸發檢查 + 分類知識類型\n   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）\n   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 ${AI_SKILL_REPO}/enforcement/sanitization.md 去敏）\n   - Step 5：更新目標層\n   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning\n   - Step 8：執行 Linked Updates\n   - Step 9：更新 Runtime Surfaces\n   - Step 10：驗證（diff review、去敏檢查、link check）\n   - Step 11：Commit / Push / Readback（關閉 writeback transaction）\n4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。\n\n## 專案 durable Markdown 預設（Ai-skill）\n\n當任務為在本專案新增或大幅改寫長期保存的 Markdown（例如 docs/、wiki/、README.md、ADR、runbook 等，以本 repo 慣例為準）時，動筆前先讀 ${AI_SKILL_REPO}/workflow/documentation/README.md，再讀 ${AI_SKILL_REPO}/workflow/documentation/execution-flow.md。拆分與篇幅見 ${AI_SKILL_REPO}/governance/document-sizing.md。無須使用者先說「docs/wiki/token」等關鍵字；此為上述文件類型之預設流程。",
      "groups": ["read"],
      "fileRestrictions": ["**/*.md"]
    },
    {
      "slug": "ask",
      "name": "❓ Ask",
      "roleDefinition": "You are Roo, a knowledgeable technical advisor who provides clear explanations, documentation, and answers to technical questions.",
      "customInstructions": "依 ${AI_SKILL_REPO}/CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。\n\n## 語言強制規則（最高優先級）\n\n### 核心原則\n- **強制跟隨使用者語言**：使用者用什麼語言提問，所有輸出就必須用什麼語言。\n- **無預設語言**：不存在「Default to English」。使用者第一次提問的語言即為本次對話語言。\n- **禁止自行切換**：即使分析技術內容、寫 commit message、產生表格，也必須與使用者當前語言一致。\n- **attempt_completion 前必須檢查**：確認結果的語言與使用者最後一次提問的語言一致。\n\n### 防漂移機制\n- 如果使用者用中文，所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解）都必須使用中文。\n- 如果使用者用日文，所有輸出都必須使用日文。\n- 如果使用者切換語言，跟隨其切換。\n- **禁止出現「Default to English」思維**：不存在預設語言，只有「使用者當前語言」。\n\n## 知識更新流程 Checkpoint\n\n每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：\n\n1. 讀取 ${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md 了解完整流程。\n2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？\n3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：\n   - Step 1-2：觸發檢查 + 分類知識類型\n   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）\n   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 ${AI_SKILL_REPO}/enforcement/sanitization.md 去敏）\n   - Step 5：更新目標層\n   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning\n   - Step 8：執行 Linked Updates\n   - Step 9：更新 Runtime Surfaces\n   - Step 10：驗證（diff review、去敏檢查、link check）\n   - Step 11：Commit / Push / Readback（關閉 writeback transaction）\n4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。\n\n## 專案 durable Markdown 預設（Ai-skill）\n\n當任務為在本專案新增或大幅改寫長期保存的 Markdown（例如 docs/、wiki/、README.md、ADR、runbook 等，以本 repo 慣例為準）時，動筆前先讀 ${AI_SKILL_REPO}/workflow/documentation/README.md，再讀 ${AI_SKILL_REPO}/workflow/documentation/execution-flow.md。拆分與篇幅見 ${AI_SKILL_REPO}/governance/document-sizing.md。無須使用者先說「docs/wiki/token」等關鍵字；此為上述文件類型之預設流程。",
      "groups": ["read"]
    },
    {
      "slug": "debug",
      "name": "🪲 Debug",
      "roleDefinition": "You are Roo, a systematic debugger specializing in troubleshooting issues, investigating errors, and diagnosing problems.",
      "customInstructions": "依 ${AI_SKILL_REPO}/CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。\n\n## 語言強制規則（最高優先級）\n\n### 核心原則\n- **強制跟隨使用者語言**：使用者用什麼語言提問，所有輸出就必須用什麼語言。\n- **無預設語言**：不存在「Default to English」。使用者第一次提問的語言即為本次對話語言。\n- **禁止自行切換**：即使分析技術內容、寫 commit message、產生表格，也必須與使用者當前語言一致。\n- **attempt_completion 前必須檢查**：確認結果的語言與使用者最後一次提問的語言一致。\n\n### 防漂移機制\n- 如果使用者用中文，所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解）都必須使用中文。\n- 如果使用者用日文，所有輸出都必須使用日文。\n- 如果使用者切換語言，跟隨其切換。\n- **禁止出現「Default to English」思維**：不存在預設語言，只有「使用者當前語言」。\n\n## 知識更新流程 Checkpoint\n\n每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：\n\n1. 讀取 ${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md 了解完整流程。\n2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？\n3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：\n   - Step 1-2：觸發檢查 + 分類知識類型\n   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）\n   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 ${AI_SKILL_REPO}/enforcement/sanitization.md 去敏）\n   - Step 5：更新目標層\n   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning\n   - Step 8：執行 Linked Updates\n   - Step 9：更新 Runtime Surfaces\n   - Step 10：驗證（diff review、去敏檢查、link check）\n   - Step 11：Commit / Push / Readback（關閉 writeback transaction）\n4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。\n\n## 專案 durable Markdown 預設（Ai-skill）\n\n當任務為在本專案新增或大幅改寫長期保存的 Markdown（例如 docs/、wiki/、README.md、ADR、runbook 等，以本 repo 慣例為準）時，動筆前先讀 ${AI_SKILL_REPO}/workflow/documentation/README.md，再讀 ${AI_SKILL_REPO}/workflow/documentation/execution-flow.md。拆分與篇幅見 ${AI_SKILL_REPO}/governance/document-sizing.md。無須使用者先說「docs/wiki/token」等關鍵字；此為上述文件類型之預設流程。",
      "groups": ["read", "edit", "command", "mcp"]
    },
    {
      "slug": "orchestrator",
      "name": "🪃 Orchestrator",
      "roleDefinition": "You are Roo, a skilled orchestrator coordinating complex, multi-step projects that require coordination across different specialties.",
      "customInstructions": "依 ${AI_SKILL_REPO}/CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。\n\n## 語言強制規則（最高優先級）\n\n### 核心原則\n- **強制跟隨使用者語言**：使用者用什麼語言提問，所有輸出就必須用什麼語言。\n- **無預設語言**：不存在「Default to English」。使用者第一次提問的語言即為本次對話語言。\n- **禁止自行切換**：即使分析技術內容、寫 commit message、產生表格，也必須與使用者當前語言一致。\n- **attempt_completion 前必須檢查**：確認結果的語言與使用者最後一次提問的語言一致。\n\n### 防漂移機制\n- 如果使用者用中文，所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解）都必須使用中文。\n- 如果使用者用日文，所有輸出都必須使用日文。\n- 如果使用者切換語言，跟隨其切換。\n- **禁止出現「Default to English」思維**：不存在預設語言，只有「使用者當前語言」。\n\n## 知識更新流程 Checkpoint\n\n每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：\n\n1. 讀取 ${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md 了解完整流程。\n2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？\n3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：\n   - Step 1-2：觸發檢查 + 分類知識類型\n   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）\n   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 ${AI_SKILL_REPO}/enforcement/sanitization.md 去敏）\n   - Step 5：更新目標層\n   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning\n   - Step 8：執行 Linked Updates\n   - Step 9：更新 Runtime Surfaces\n   - Step 10：驗證（diff review、去敏檢查、link check）\n   - Step 11：Commit / Push / Readback（關閉 writeback transaction）\n4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。\n\n## 專案 durable Markdown 預設（Ai-skill）\n\n當任務為在本專案新增或大幅改寫長期保存的 Markdown（例如 docs/、wiki/、README.md、ADR、runbook 等，以本 repo 慣例為準）時，動筆前先讀 ${AI_SKILL_REPO}/workflow/documentation/README.md，再讀 ${AI_SKILL_REPO}/workflow/documentation/execution-flow.md。拆分與篇幅見 ${AI_SKILL_REPO}/governance/document-sizing.md。無須使用者先說「docs/wiki/token」等關鍵字；此為上述文件類型之預設流程。",
      "groups": ["read", "edit", "command", "mcp"]
    }
  ]
}
ROOMODES_EOF
)
  write_file "$ROOMODES_FILE" "$ROOMODES_CONTENT" "Roo Code 設定檔 (.roomodes)"
fi

# ─── 2. Cursor 設定 ──────────────────────────────────────────────────────
if [[ -n "${TOOL_MAP[cursor]:-}" ]]; then
  info "── Cursor ──"

  CURSOR_RULES_DIR="$TARGET_DIR/.cursor/rules"
  CURSOR_FILE="$CURSOR_RULES_DIR/ai-skill-bootstrap.mdc"
  CURSOR_CONTENT=$(cat <<CURSOR_EOF
---
description: Ai-skill 知識庫啟動流程
globs: 
alwaysApply: true
---

# Ai-skill Bootstrap

本專案使用 [Ai-skill 知識庫](${AI_SKILL_REPO}/) 作為 AI 輔助開發的知識基礎。

## 啟動流程

每個 session 啟動時，依以下順序載入：

1. 讀取 Core Bootstrap：[\`${AI_SKILL_REPO}/CORE_BOOTSTRAP.md\`](${AI_SKILL_REPO}/CORE_BOOTSTRAP.md)
2. 了解 OS layout：[\`${AI_SKILL_REPO}/README.md\`](${AI_SKILL_REPO}/README.md)
3. 依任務 intent 查詢 routing-registry：[\`${AI_SKILL_REPO}/knowledge/runtime/routing-registry.yaml\`](${AI_SKILL_REPO}/knowledge/runtime/routing-registry.yaml)
4. 載入 lazy-load rules（依 activation conditions）

## 專案 durable Markdown（預設 workflow）

當任務為在本專案**新增或大幅改寫**長期保存的 Markdown（例如 \`docs/\`、\`wiki/\`、根目錄或子目錄 \`README.md\`、ADR、runbook 等，以本 repo 慣例為準）時，**動筆前**先讀：

1. [\`${AI_SKILL_REPO}/workflow/documentation/README.md\`](${AI_SKILL_REPO}/workflow/documentation/README.md)
2. 再依需要讀 [\`${AI_SKILL_REPO}/workflow/documentation/execution-flow.md\`](${AI_SKILL_REPO}/workflow/documentation/execution-flow.md)

拆分與篇幅：[\`${AI_SKILL_REPO}/governance/document-sizing.md\`](${AI_SKILL_REPO}/governance/document-sizing.md)。

**無須**使用者先說「docs / wiki / token」等關鍵字；凡寫入上述類型之專案文件，預設即遵守此流程。

## 語言強制規則（最高優先級）

### 核心原則
- **強制跟隨使用者語言**：使用者用什麼語言提問，所有輸出就必須用什麼語言。
- **無預設語言**：不存在「Default to English」。使用者第一次提問的語言即為本次對話語言。
- **禁止自行切換**：即使分析技術內容、寫 commit message、產生表格，也必須與使用者當前語言一致。

### 防漂移機制
- 如果使用者用中文，所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解）都必須使用中文。
- 如果使用者用日文，所有輸出都必須使用日文。
- 如果使用者切換語言，跟隨其切換。
- **禁止出現「Default to English」思維**：不存在預設語言，只有「使用者當前語言」。

## 知識更新流程 Checkpoint

每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：

1. 讀取 [\`${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md\`](${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md) 了解完整流程。
2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？
3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：
   - Step 1-2：觸發檢查 + 分類知識類型
   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）
   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 sanitization.md 去敏）
   - Step 5：更新目標層
   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning
   - Step 8：執行 Linked Updates
   - Step 9：更新 Runtime Surfaces
   - Step 10：驗證（diff review、去敏檢查、link check）
   - Step 11：Commit / Push / Readback（關閉 writeback transaction）
4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。
CURSOR_EOF
)
  write_file "$CURSOR_FILE" "$CURSOR_CONTENT" "Cursor 規則檔 (.cursor/rules/ai-skill-bootstrap.mdc)"

  # Cursor hooks（可選）
  CURSOR_HOOKS_DIR="$TARGET_DIR/.cursor"
  CURSOR_HOOKS_FILE="$CURSOR_HOOKS_DIR/hooks.json"
  CURSOR_HOOKS_CONTENT=$(cat <<HOOKS_EOF
{
  "sessionStart": [
    {
      "description": "提醒載入 Ai-skill 知識庫",
      "command": "echo '提示：本專案使用 Ai-skill 知識庫，請確認已載入 CORE_BOOTSTRAP.md'"
    }
  ]
}
HOOKS_EOF
)
  write_file "$CURSOR_HOOKS_FILE" "$CURSOR_HOOKS_CONTENT" "Cursor hooks 設定 (.cursor/hooks.json)"
fi

# ─── 3. Claude Code 設定 ─────────────────────────────────────────────────
if [[ -n "${TOOL_MAP[claude]:-}" ]]; then
  info "── Claude Code ──"

  CLAUDE_FILE="$TARGET_DIR/CLAUDE.md"
  CLAUDE_CONTENT=$(cat <<CLAUDE_EOF
# Claude Code Auto-Bootstrap

本專案使用 [Ai-skill 知識庫](${AI_SKILL_REPO}/) 作為 AI 輔助開發的知識基礎。

## 啟動流程

每個 session 啟動時，依以下順序載入：

1. 讀取 Core Bootstrap：[\`${AI_SKILL_REPO}/CORE_BOOTSTRAP.md\`](${AI_SKILL_REPO}/CORE_BOOTSTRAP.md)
2. 了解 OS layout：[\`${AI_SKILL_REPO}/README.md\`](${AI_SKILL_REPO}/README.md)
3. 依任務 intent 查詢 routing-registry：[\`${AI_SKILL_REPO}/knowledge/runtime/routing-registry.yaml\`](${AI_SKILL_REPO}/knowledge/runtime/routing-registry.yaml)
4. 載入 lazy-load rules（依 activation conditions）

## 專案 durable Markdown（預設 workflow）

當任務為在本專案**新增或大幅改寫**長期保存的 Markdown（例如 \`docs/\`、\`wiki/\`、根目錄或子目錄 \`README.md\`、ADR、runbook 等，以本 repo 慣例為準）時，**動筆前**先讀：

1. [\`${AI_SKILL_REPO}/workflow/documentation/README.md\`](${AI_SKILL_REPO}/workflow/documentation/README.md)
2. 再依需要讀 [\`${AI_SKILL_REPO}/workflow/documentation/execution-flow.md\`](${AI_SKILL_REPO}/workflow/documentation/execution-flow.md)

拆分與篇幅：[\`${AI_SKILL_REPO}/governance/document-sizing.md\`](${AI_SKILL_REPO}/governance/document-sizing.md)。

**無須**使用者先說「docs / wiki / token」等關鍵字；凡寫入上述類型之專案文件，預設即遵守此流程。

## 語言強制規則（最高優先級）

### 核心原則
- **強制跟隨使用者語言**：使用者用什麼語言提問，所有輸出就必須用什麼語言。
- **無預設語言**：不存在「Default to English」。使用者第一次提問的語言即為本次對話語言。
- **禁止自行切換**：即使分析技術內容、寫 commit message、產生表格，也必須與使用者當前語言一致。

### 防漂移機制
- 如果使用者用中文，所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解）都必須使用中文。
- 如果使用者用日文，所有輸出都必須使用日文。
- 如果使用者切換語言，跟隨其切換。
- **禁止出現「Default to English」思維**：不存在預設語言，只有「使用者當前語言」。

## 知識更新流程 Checkpoint

每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：

1. 讀取 [\`${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md\`](${AI_SKILL_REPO}/governance/lifecycle/knowledge-update-flow.md) 了解完整流程。
2. 自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？
3. 若是，依 knowledge-update-flow.md 的 11 個步驟執行：
   - Step 1-2：觸發檢查 + 分類知識類型
   - Step 3：決定 Promotion Target（intelligence / workflow / analysis / shared-rules / runtime / memory）
   - Step 4：寫入 feedback/history/<domain>/<category>/ lesson（寫入前依 sanitization.md 去敏）
   - Step 5：更新目標層
   - Step 6-7：選擇性執行 Intelligence Extraction 或 Failure Learning
   - Step 8：執行 Linked Updates
   - Step 9：更新 Runtime Surfaces
   - Step 10：驗證（diff review、去敏檢查、link check）
   - Step 11：Commit / Push / Readback（關閉 writeback transaction）
4. 若否，簡短說明本輪只有 project-specific evidence 或尚未達可泛化標準。
CLAUDE_EOF
)
  write_file "$CLAUDE_FILE" "$CLAUDE_CONTENT" "Claude Code 設定檔 (CLAUDE.md)"
fi

# ─── 4. 對話目標帳本目錄（選用） ─────────────────────────────────────────
if [[ "$DRY_RUN" != true ]]; then
  AGENT_GOALS_DIR="$TARGET_DIR/.agent-goals"
  if [[ ! -d "$AGENT_GOALS_DIR" ]]; then
    mkdir -p "$AGENT_GOALS_DIR"
    # 建立初始 README.md
    cat > "$AGENT_GOALS_DIR/README.md" <<GOALS_EOF
# Agent Goals

本目錄由 [Ai-skill 對話目標帳本](${AI_SKILL_REPO}/enforcement/conversation-goal-ledger.md) 管理。
用 \`scripts/agent-goals.sh\` 操作。

## 目前目標

（尚無 active goal）
GOALS_EOF
    ok "已建立對話目標目錄: $AGENT_GOALS_DIR"
  fi
fi

# ─── 完成 ─────────────────────────────────────────────────────────────────
echo ""
if [[ "$DRY_RUN" == true ]]; then
  info "DRY RUN 完成。使用以下命令實際執行："
  echo "  $(basename "$0") $TARGET_DIR"
else
  ok "初始化完成！"
  echo ""
  echo "設定的工具：${!TOOL_MAP[*]}"
  echo "目標專案：$TARGET_DIR"
  echo "Ai-skill：$AI_SKILL_REPO"
  echo ""
  echo "下一步："
  echo "  1. 在目標專案中開啟對應的 AI 工具"
  echo "  2. 工具會自動載入 CORE_BOOTSTRAP.md 啟動流程"
  echo "  3. 開始開發！"
  echo ""
  echo "如果 Ai-skill 路徑變更，重新執行此腳本即可更新所有路徑。"
fi
