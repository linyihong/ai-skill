# Tool Metadata

`tools/metadata/` 負責定義每個 AI tool 的 metadata schema，讓 runtime 可以根據 tool 的成本、風險與適用情境做 routing 與 lazy activation。

## Tool Metadata Schema

```yaml
tool:
  id: repo-search
  name: Repository Search
  description: 在 repository 中搜尋檔案或內容

  cost:
    avg_input_tokens: 800     # 平均每次呼叫的 input tokens
    avg_output_tokens: 1500   # 平均每次呼叫的 output tokens
    risk:
      recursive_expansion: true   # 是否可能引發遞迴搜尋
      side_effects: false         # 是否有副作用（寫入/修改）

  contexts:
    - debug
    - coding
    - analysis

  activation:
    strategy: lazy             # preload | lazy | on_demand
    priority: P2

  compression:
    supported: true            # 是否支援 output compression
    default_level: summary     # raw | summary | structured
```

## Tool 清單

| Tool ID | avg_input | avg_output | recursive | contexts | Activation |
| --- | --- | --- | --- | --- | --- |
| `repo-search` | 800 | 1500 | true | debug, coding, analysis | lazy (P2) |
| `read-file` | 200 | 2000 | false | all | preload (P0) |
| `write-to-file` | 500 | 100 | false | all | lazy (P1) |
| `execute-command` | 300 | 3000 | true | dev, test, deploy | lazy (P1) |
| `apply-diff` | 400 | 200 | false | coding, review | lazy (P1) |
| `search-files` | 600 | 1200 | true | debug, analysis | lazy (P2) |
| `ask-followup` | 200 | 500 | false | all | on_demand (P3) |

## 與既有層的關係

- `ai-tools/`：工具配置與同步細節（Claude Code、Cursor 等）
- `runtime/runtime.db`：tool call 次數限制與 explosion detection
- `tools/routing/`：tool lazy activation 決策邏輯
- `tools/compression/`：tool output compression 策略
