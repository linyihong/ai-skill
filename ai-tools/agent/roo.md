# Roo Code 使用說明

本檔說明 Roo Code（VS Code extension）與其他工具的差異。通用配置原則見 [`ai-tools/README.md`](../README.md)；知識庫入口見 [`README.md`](../../README.md)；啟動流程見 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。

## Roo Code 與其他工具的差異

| 特性 | Claude Code | Cursor | Roo Code |
|------|------------|--------|----------|
| 執行環境 | CLI terminal | 獨立 IDE | VS Code extension |
| 自動載入入口 | `CLAUDE.md` | `.cursor/rules/*.mdc` | Custom instructions（VS Code settings 或 `.roomodes`） |
| 工具配置 | `.claude/settings.json` | `.cursor/hooks.json` | VS Code settings（`roo-code.*`） |
| Modes | 單一模式 | 單一模式 | 多 modes（code/architect/ask/debug 等） |
| File restrictions | 無 | 無 | 每個 mode 可限制可讀寫的 file patterns |

## Roo Code 配置實作

### 自動載入入口：Custom Instructions

Roo Code 沒有像 `CLAUDE.md` 那樣的自動載入機制。你需要手動設定 **Custom Instructions**。

Custom Instructions 有**兩個層級**，各有不同的用途：

#### 層級 A：全域 Custom Instructions（所有專案共用）

在 Roo Code Extension 設定面板中設定的 Custom Instructions，**每個專案都會載入**。

**設定位置**：VS Code → Roo Code Extension 圖示 → 設定（⚙️）→ **Custom Instructions** 欄位

**建議內容**（貼入該欄位）：

```text
你是一個運行在 Roo Code（VS Code AI extension）的 AI agent。

開始工作前，請依 <AI_SKILL_REPO>/CORE_BOOTSTRAP.md 的啟動流程載入核心規則與 OS layout。

Language Preference: Default to English, but always match the user's language in conversation. If the user writes in Chinese, respond in Chinese. If the user writes in Japanese, respond in Japanese. If the user switches languages, follow their switch.

語言一致性強制規則：所有輸出（包含 attempt_completion 結果、技術分析、表格欄位、章節標題、commit message）都必須與使用者當前語言一致。如果使用者使用中文，所有內容（包括技術關鍵詞、程式碼註解、分析報告）都必須使用中文。在 attempt_completion 前必須先確認語言一致性。
```

**適用時機**：
- 你希望**所有專案**都自動啟用 Ai-skill 知識庫系統
- 新專案不需要額外設定，直接開啟 Roo Code 即可
- 全域設定是**軟性預設值**，專案可以透過 `.roomodes` 覆蓋

##### 自動寫入方式（AI agent 專用）

如果 AI agent 有檔案系統寫入權限，可以直接修改 VS Code 的 `state.vscdb`（SQLite 資料庫），無需使用者手動操作。

**⚠️ 重要限制**：VS Code 的 extension host 會主動管理 `state.vscdb`。如果 VS Code 正在執行，直接寫入 SQLite 後，VS Code 在下一次狀態變更（如切換 mode、變更 API 設定）時會用自己的記憶體狀態覆寫資料庫。因此**必須先關閉 VS Code**再執行寫入。

**資料庫路徑**：
```
~/Library/Application Support/Code/User/globalStorage/state.vscdb
```

**Table / Key 結構**：
- Table: `ItemTable`（key-value store）
- Key: `RooVeterinaryInc.roo-cline`
- Value: JSON blob，包含所有 Roo Code 全域設定

**寫入欄位**：在 JSON blob 中加入 `customInstructions` 欄位

**建議使用專用腳本**（內含 VS Code 執行中檢查）：

```bash
# 1. 先關閉 VS Code（Cmd+Q）
# 2. 執行腳本
python3 scripts/set-roo-global-custom-instructions.py
# 3. 重新開啟 VS Code
```

腳本路徑：[`scripts/set-roo-global-custom-instructions.py`](scripts/set-roo-global-custom-instructions.py)

**手動 Python 寫入範例**（供參考）：

```python
import json, sqlite3, os, subprocess

CUSTOM_INSTRUCTIONS = """你是一個運行在 Roo Code（VS Code AI extension）的 AI agent。

開始工作前，請依 <AI_SKILL_REPO>/CORE_BOOTSTRAP.md 的啟動流程載入核心規則與 OS layout。

Language Preference: Default to English, but always match the user's language in conversation. If the user writes in Chinese, respond in Chinese. If the user writes in Japanese, respond in Japanese. If the user switches languages, follow their switch.

語言一致性強制規則：所有輸出（包含 attempt_completion 結果、技術分析、表格欄位、章節標題、commit message）都必須與使用者當前語言一致。如果使用者使用中文，所有內容（包括技術關鍵詞、程式碼註解、分析報告）都必須使用中文。在 attempt_completion 前必須先確認語言一致性。"""

DB_PATH = os.path.expanduser(
    "~/Library/Application Support/Code/User/globalStorage/state.vscdb"
)

# 步驟 1：檢查 VS Code 是否正在執行
result = subprocess.run(["pgrep", "-f", "Visual Studio Code"],
                        capture_output=True, text=True, timeout=5)
if result.returncode == 0 and len(result.stdout.strip()) > 0:
    print("⚠️  VS Code 正在執行中！請先關閉 VS Code 再執行。")
    exit(1)

# 步驟 2：讀取現有 JSON
conn = sqlite3.connect(DB_PATH)
cursor = conn.cursor()
cursor.execute("SELECT value FROM ItemTable WHERE key = 'RooVeterinaryInc.roo-cline'")
row = cursor.fetchone()
data = json.loads(row[0])

# 步驟 3：寫入 customInstructions
data["customInstructions"] = CUSTOM_INSTRUCTIONS
new_value = json.dumps(data, ensure_ascii=False)
cursor.execute(
    "UPDATE ItemTable SET value = ? WHERE key = 'RooVeterinaryInc.roo-cline'",
    (new_value,)
)
conn.commit()

# 步驟 4：強制 WAL checkpoint
conn.execute("PRAGMA wal_checkpoint(TRUNCATE);")
conn.close()

print("✅ 寫入成功！請重新開啟 VS Code。")
```

**注意事項**：
1. **必須先關閉 VS Code**，否則寫入會被 VS Code 覆寫
2. `state.vscdb` 使用 WAL（Write-Ahead Log）模式，寫入後必須執行 `PRAGMA wal_checkpoint(TRUNCATE);` 才能確保持久化
3. 修改後需**重新開啟 VS Code** 才會生效
4. 如果 Ai-skill 路徑變更，需同步更新 `CUSTOM_INSTRUCTIONS` 中的絕對路徑
5. 此方法也適用於修改其他 Roo Code 全域設定（如 `language` 欄位）

#### 層級 B：專案 `.roomodes`（單一專案專用）

在專案根目錄建立 `.roomodes` 檔案，Roo Code 會自動讀取。此設定**只對該專案生效**。

**適用時機**：
- 該專案需要**不同的 mode 定義**（如不同的 file restrictions、role definition）
- 該專案需要**覆蓋**全域 Custom Instructions 中的某些設定
- 該專案不在本機，無法使用絕對路徑指向 Ai-skill

**注意**：`.roomodes` 中的 `customInstructions` 會**完全覆蓋**全域 Custom Instructions，不會合併。所以如果用了 `.roomodes`，需要把全域的設定也複製進去。

#### 建議策略

```
全域 Custom Instructions（層級 A）
  ├── 指向 Ai-skill 的 CORE_BOOTSTRAP.md（絕對路徑）
  ├── 語言偏好設定
  └── 語言一致性強制規則

專案 .roomodes（層級 B，可選）
  ├── 只在需要自訂 mode 定義時建立
  ├── 必須包含層級 A 的所有內容（因為會完全覆蓋）
  └── 加上該專案特有的 mode 設定
```

#### 如果 Ai-skill 路徑變更

全域 Custom Instructions 中使用的是絕對路徑（如 `<AI_SKILL_REPO>/`）。如果：
- **Ai-skill 移動位置** → 更新全域 Custom Instructions 中的路徑
- **在其他電腦使用** → 修改為對應的絕對路徑
- **使用相對路徑** → 只能在 Ai-skill repo 內生效，不適合全域設定

### Modes 設定

Roo Code 支援多種 modes，每個 mode 可以有不同的 custom instructions 與 file restrictions。
你可以在 VS Code 的 `settings.json` 中設定：

```json
{
  "roo-code.modes": [
    {
      "slug": "code",
      "name": "💻 Code",
      "roleDefinition": "You are Roo, a highly skilled software engineer...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read", "edit", "command", "mcp"]
    },
    {
      "slug": "architect",
      "name": "🏗️ Architect",
      "roleDefinition": "You are Roo, an expert software architect...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read"],
      "fileRestrictions": ["**/*.md"]
    },
    {
      "slug": "ask",
      "name": "❓ Ask",
      "roleDefinition": "You are Roo, a knowledgeable technical advisor...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read"]
    },
    {
      "slug": "debug",
      "name": "🪲 Debug",
      "roleDefinition": "You are Roo, a systematic debugger...",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。",
      "groups": ["read", "edit", "command", "mcp"]
    }
  ]
}
```

或者，你也可以在專案根目錄建立 `.roomodes` 檔案（Roo Code 會自動讀取）：

```json
{
  "customModes": [
    {
      "slug": "ai-skill-maintainer",
      "name": "Ai-skill Maintainer",
      "roleDefinition": "你負責維護 Ai-skill 知識庫的架構、規則與技能。",
      "customInstructions": "依 CORE_BOOTSTRAP.md 啟動流程載入核心規則與 OS layout。\n\nLanguage Preference: Default to English, but always match the user's language in conversation. If the user writes in Chinese, respond in Chinese. If the user writes in Japanese, respond in Japanese. If the user switches languages, follow their switch.",
      "groups": ["read", "edit", "command", "mcp"]
    }
  ]
}
```

> **注意**：`.roomodes` 中的 `customInstructions` 會**覆蓋** Roo Code Extension 設定面板中的全域 Custom Instructions。如果你已經在全域設定中寫了固定的語言偏好，`.roomodes` 的設定會優先。

### 建議的 `.roomodes` 設定

由於本知識庫有多個操作面向，建議建立以下自訂 modes：

| Mode | Slug | 用途 | Groups |
|------|------|------|--------|
| 🏗️ Architect | `architect` | 規劃架構、設計 pipeline | read |
| 💻 Code | `code` | 實作元件、寫規則、建立檔案 | read, edit, command, mcp |
| ❓ Ask | `ask` | 查詢知識、解釋架構 | read |
| 🪲 Debug | `debug` | 診斷問題、追蹤失效模式 | read, edit, command, mcp |
| 🪃 Orchestrator | `orchestrator` | 協調多步驟任務 | read, edit, command, mcp |
| 📝 Skill Writer | `skill-writer` | 撰寫 skill、feedback lesson | read, edit |
| 🧹 Governance | `governance` | 維護 lifecycle、validation | read, edit |

## 語言偏好設定（重要）

Roo Code 的語言偏好設定涉及**兩個層級**，需要分別處理才能完整解決語言漂移問題。

### 問題

如果 Custom Instructions 中寫了固定的語言偏好（例如 `You should always speak and think in the "English" (en) language`），agent 會**無視使用者實際使用的語言**，強制用該語言回應。這是因為 Custom Instructions 在 system prompt 中的優先級高於對話上下文。

此外，Roo Code 的全域設定中有一個 `language` 欄位（預設 `"en"`），此欄位會被傳入 system prompt 建構過程，進一步強化英文偏好。

### 解決方案：兩個層級

#### 層級 1：Custom Instructions（`.roomodes`）

在 `.roomodes` 的每個 mode 的 `customInstructions` 中，將語言偏好設為**軟性預設值**，加上「跟隨使用者語言」的例外：

```text
Language Preference: Default to English, but always match the user's language in conversation.
If the user writes in Chinese, respond in Chinese.
If the user writes in Japanese, respond in Japanese.
If the user switches languages, follow their switch.
```

> **原理**：`.roomodes` 中的 `customInstructions` 會**覆蓋** Roo Code Extension 設定面板中的全域 Custom Instructions，因此這是設定語言偏好的主要方式。

#### 層級 2：全域 `language` 欄位（SQLite）

Roo Code 的全域設定儲存在 VS Code 的 **globalStorage SQLite 資料庫**中。其中的 `language` 欄位（預設 `"en"`）會影響 system prompt 的建構。

VS Code Extension 全域設定的通用查詢/修改方法見 [`intelligence/ide/vscode-extension-global-state.md`](../../intelligence/ide/vscode-extension-global-state.md)。Roo Code 專屬的 key 為 `RooVeterinaryInc.roo-cline`，需修改的欄位為 `language`（設為 `"zh-CN"`）。

### 如何預設語言偏好以跟隨文件規則

本知識庫的文件規則要求：**預設英文，但跟隨使用者語言**。要讓 Roo Code 遵守此規則，需完成以下步驟：

1. **建立 `.roomodes`**（已建立）：在專案根目錄建立 `.roomodes`，每個 mode 的 `customInstructions` 中包含上述軟性語言偏好設定。
2. **修改全域 `language` 欄位**（已修改）：將 SQLite 中的 `language` 從 `"en"` 改為 `"zh-CN"`，避免 system prompt 強制使用英文。通用修改方法見 [`intelligence/ide/vscode-extension-global-state.md`](../../intelligence/ide/vscode-extension-global-state.md)。
3. **驗證**：啟動新 session 後，用中文提問確認 agent 以中文回應。

> **注意**：如果只做步驟 1 不做步驟 2，Roo Code 的 system prompt 仍可能因 `language: "en"` 而傾向英文。兩個層級都設定才能完整解決。

### 相關資源

- VS Code Extension 全域設定修改通用指南: [`intelligence/ide/vscode-extension-global-state.md`](../../intelligence/ide/vscode-extension-global-state.md)
- Failure pattern: [`enforcement/failure-patterns/language-preference-drift.md`](../../enforcement/failure-patterns/language-preference-drift.md)
- Feedback lesson: [`feedback/history/roo-code/2026-05-13_124800-language-preference-drift.md`](../../feedback/history/roo-code/2026-05-13_124800-language-preference-drift.md)
- 設定檔: [`.roomodes`](../../.roomodes)

## Roo Code 與對話目標閉環

工具中立規則見 [`enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md)。Roo Code 沒有像 Cursor 那樣的 hooks 機制，但可以透過以下方式實作對話目標閉環。

### 在 Custom Instructions 中加入 goal ledger 提醒

在 `.roomodes` 的每個 mode 的 `customInstructions` 中，加入以下內容：

```text
開始工作前，若 `<PROJECT_ROOT>/.agent-goals/` 存在，先讀取確認是否有未完成的 active goal。
若 goal 標示 `single-owner` 或 `non-parallelizable`，不要和其他 agent 分工同一流程；先取得使用者確認。
完成目標後，若仍有長期 roadmap 或治理狀態，先回寫到 durable planning 文件，再刪除 active goal。
```

### 建議操作

在 Roo Code 開始可中斷、可拆解或多目標工作時，或已看到 active project 有 modified / staged / untracked files、已建立 TodoWrite、使用者說「繼續」前一個多步驟任務時：

1. 讀取 `<PROJECT_ROOT>/.agent-goals/`，確認是否已有 active / blocked / needs-validation goal，以及 priority、owner、lock、parallelization mode、plan/todo links、missing/decision/strengthen。若使用者說 agent 中斷、突然關閉、要從哪裡重做、剩下什麼或下一步是什麼，必須先讀 `.agent-goals/README.md` 與對應 active goal，再用 transcripts、terminal output、git status 交叉確認；不要把 transcript/git 當成第一真相來源。
2. 若沒有 ledger 且任務不是單一回覆即可完成，使用本庫 helper 初始化；不要因為已有 TodoWrite 就跳過 goal ledger：

   ```bash
   <AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> init
   ```

3. 建立或更新本輪主要目標：

   ```bash
   <AI_SKILL_REPO>/scripts/agent-goals.sh --project <PROJECT_ROOT> start \
     --id P1-short-goal \
     --title "Short goal title" \
     --source "User request summary" \
     --parallelization single-owner \
     --next "Next concrete action" \
     --criteria "Observable completion condition"
   ```

4. 若使用者轉移目標，先 `pause` 或 `update --status superseded` 舊 goal，再建立新的 `P1`。
5. 若有 planning 文件或 TodoWrite todo，使用 `--plan` / `--todo` 連到 goal，並讓 `.agent-goals/README.md` 的主目標表可快速跳回該 goal。
6. 若發現需要拆小目標，使用 `split` 或在 goal 檔的 `Subgoals` 區塊記錄；若發現不能分工或需單一 owner，使用 `--parallelization single-owner|non-parallelizable` 更新。
7. 在回覆完成前，只有完成條件與驗證都成立時才 `complete --validated`；條件已成立時必須同輪刪除 goal 並刷新 `.agent-goals/README.md`，不要把 `completed` row 留在 active 表。若該 goal 完成後仍有長期 roadmap、phase、migration、promotion、deprecation 或治理狀態，先回寫到 durable planning 文件，再刪除 active goal。否則保留 goal，讓下一個 agent 可接手。

## Roo Code 與知識更新流程 Checkpoint

工具中立規則見 [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md)。Roo Code 沒有 hooks 機制，但可以透過 `.roomodes` 的 `customInstructions` 加入 checkpoint 提醒。

### 在 Custom Instructions 中加入 checkpoint 提醒

在 `.roomodes` 的每個 mode 的 `customInstructions` 中，在 goal ledger 提醒之後加入以下內容：

```text
## 知識更新流程 Checkpoint

每輪工作結束前、切回長時間專案工作前、或使用者說「繼續」展開下一輪前，必須執行知識更新檢查：

1. 讀取 [`<AI_SKILL_REPO>/governance/lifecycle/knowledge-update-flow.md`] 了解完整流程。
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
```

### 與 Cursor 的差異

| 特性 | Cursor | Roo Code |
|------|--------|----------|
| 自動提醒機制 | `.cursor/hooks.json`（sessionStart / preCompact / stop） | 無 hooks，需在 Custom Instructions 中手動提醒 |
| Goal ledger 操作 | 可透過 hooks 自動檢查 | 需在每個 session 開始時手動讀取 |
| 知識更新流程 Checkpoint | 可在 `.cursor/rules/*.mdc` 中加入，可選 hooks 輔助 | 需在 `.roomodes` Custom Instructions 中加入 |
| 多 mode 支援 | 單一模式 | 多 modes，每個 mode 可獨立設定 checkpoint 提醒 |

## Roo Code 操作注意

### 1. 工作區設定

Roo Code 直接在 VS Code 中操作，所以工作區就是目前開啟的資料夾。
建議用多資料夾工作區同時開啟業務專案與本 repo：

```text
<PROJECT_ROOT>/          ← 業務專案
<AI_SKILL_REPO>/         ← 本知識庫
```

### 2. 同步與更新

由於 Roo Code 直接操作 VS Code 工作區，同步流程如下：

1. 開始工作前：`git pull` 確保本庫最新
2. 修改檔案：直接編輯本庫內的檔案
3. 完成後：`git add` → `git commit` → `git push`
4. 確認：`git log --oneline -3` + `git status`

### 3. 與其他工具的協作

如果你同時使用 Claude Code（CLI）和 Roo Code（VS Code）：

- **Claude Code** 適合 CLI 操作、批次腳本、git 操作
- **Roo Code** 適合 VS Code 內的開發、檔案編輯、即時預覽
- 兩者共用同一份 `CORE_BOOTSTRAP.md`、`skills-index.yaml` 與 `runtime/pipeline/`
- 修改規則時，只需修改本庫一份，兩個工具都會讀到最新版本

---

← [回到 AI 工具索引](../README.md)
