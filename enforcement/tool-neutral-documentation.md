# 工具中立文件

可重用規則、skills、templates、索引與 feedback lessons 預設應保持工具中立。工具專屬設定、路徑、hook、命令、UI 標籤與同步細節，應放在 `ai-tools/` 或該工具自己的設定檔。

## 核心規則

撰寫可重用文件時：

1. 先描述可攜的行為、決策規則、workflow 或 artifact。
2. 使用通用詞彙，例如 `agent`、`AI tool`、`tool-specific adapter`、`local tool mirror`、`project tool config` 與 `<PROJECT_ROOT>`。
3. 除非規則確實只適用於單一工具，否則不要讓特定工具聽起來像必要條件。
4. 工具行為不同時，使用 Strategy-style adapter：核心 skill/enforcement rule 保留共同契約，工具專屬執行差異獨立放置。
5. 工具層級設定放在對應的 `ai-tools/<tool>.md`。
6. skill-specific 工具執行差異只在確實屬於該 skill 時，放入小型 skill-local adapter 文件。
7. 使用者需要具體設定步驟時，從可重用文件連到 `ai-tools/` 或 skill-local adapter。

## 允許工具專屬內容的位置

工具名稱、路徑與 UI 操作可出現在：

| 位置 | 允許內容 |
| --- | --- |
| `ai-tools/<tool>.md` | 工具設定、同步路徑、UI 步驟、hooks、settings 與 troubleshooting。 |
| `tools/adapters/<tool>.md` | 單一工具的 skill-specific 執行差異；核心 workflow 仍保持工具中立。 |
| 工具設定檔（如 `.cursor/rules/*.mdc`、`.roomodes`、`CLAUDE.md`） | 該工具載入的規則。 |
| 工具專屬 scripts 或 script docs | 只屬於該工具的命令；必要時從通用 docs 連過去。 |
| Project-local tool files | 若可安全提交到該專案，可放專案專屬 adapter 設定。 |

## 可重用文件中應避免

除非章節明確討論工具整合，否則下列位置避免工具專屬措辭：

- Root `README.md`。
- `enforcement/README.md` 索引摘要。
- `workflow/<domain>/`、`analysis/<domain>/`、`intelligence/<domain>/` 下的所有文件。
- Skill templates 與 `skills/ADDING_SKILLS.md`。
- Feedback lessons 與可重用 checklists。

範例：

| 較不通用的寫法 | 優先寫法 |
| --- | --- |
| `Cursor agent entry point` | `Agent entry point` |
| `Reload Cursor` | `Reload or refresh the active tool if it caches skills/rules` |
| `copy to .cursor/skills` | `deploy to the active tool's skill/rule location; see ai-tools/` |
| generic rule 中寫 `run sync-cursor-bundle.sh` | `run the configured tool sync; Cursor details live in ai-tools/cursor.md` |

## 與既有工具文件的關係

可重用文件可以把工具當作範例提及，但不可把工具寫成預設要求；除非該文件本來就在該工具的文件區域內。工具專屬設定與路徑應集中在 `ai-tools/agent/` 中各工具文件。

## Strategy-style Tool Adapters

當某個 skill 在不同 AI tools 中有實際執行差異時，使用此模式：

```text
<domain>/
  README.md                 # tool-neutral overview
  execution-flow.md         # tool-neutral workflow
tools/
  adapters/
    README.md               # index of supported adapters
    <tool>.md               # only the execution differences for that tool
```

核心 workflow 文件像 strategy interface：

- Trigger conditions。
- Inputs and outputs。
- Required evidence and validation。
- Safety、sanitization 與 handoff rules。
- Tool-neutral workflow 與 terminology。

每個 tool adapter 像一個 strategy implementation：

- 使用哪些工具事件、命令、hooks、prompts 或 settings。
- 工具可以自動化什麼，哪些仍需人工處理。
- 工具專屬 failure modes 與 validation。
- 連回它實作的核心 workflow 步驟。

不要把完整核心 workflow 複製到每個 adapter。若 tool-specific adapter 需要重述共同行為，應把共同內容移回核心 skill/enforcement rule，再由 adapter 連回去。

依 scope 選擇放置位置：

| Scope | 放置位置 |
| --- | --- |
| Tool-wide setup、sync、global hooks、UI、settings | `ai-tools/<tool>.md` |
| 單一工具的 skill-specific 執行差異 | `tools/adapters/<tool>.md` |
| Project-specific tool config | Project docs 或 project tool config |
| Reusable cross-tool policy | `enforcement/` |

## Review Checklist

完成可重用文件變更前，檢查：

- Root 或 skill-level 文字是否依賴單一 IDE 或 agent 產品？
- `.cursor/` 或 `~/.cursor/` 等工具專屬路徑是否只出現在 tool docs、tool config 或明確工具專屬 scripts？
- Generic rule 是否先使用「configured tool sync」，並連到 `ai-tools/` 取得具體命令？
- 如果 skill 需要 tool-specific 行為，是否隔離在 `tool-adapters/<tool>.md`，並連回核心 workflow？
- 新增的 skill 或 enforcement rule 是否誤複製工具專屬設定章節，而不是連到工具文件？

← [Back to enforcement index](README.md)
