# Model Routing

`models/routing/` 定義 model-aware execution strategy。它不控制 provider model selection、Cursor Auto、工具 UI 狀態，也不預設進入 runtime enforcement。

## 邊界

| 層級 | 責任 |
| --- | --- |
| `models/routing/` | 依 task class、cognitive state、autonomy mode、context budget 與 tool capability 選擇 execution strategy。 |
| `models/profiles/` | 保存 `small`、`large`、`specialized` 等粗略 context-loading profile。 |
| `models/capabilities/` | 描述用來修正 profile 假設的 capability dimensions。 |
| `models/compression/` | 選擇 context compression level 與 escalation path。 |
| `ai-tools/` | 記錄工具特定 model selector 行為與可用的 explicit model controls。 |
| `workflow/` | 定義任務執行形狀與成功條件。 |
| `runtime/` | 只有在 validation 與 runtime reduction 後，才接收 minimal routing primitives。 |

## Strategy Primitives

| Primitive | 意義 |
| --- | --- |
| `execution-heavy` | 有邊界的實作或文件修改，搭配一般 validation。 |
| `validation-heavy` | 修改前先 reload source、收集 evidence，並縮小 claim scope。 |
| `source-backed` | 任何實質變更前先讀 primary source 與 required dependencies。 |
| `rediscovery-only` | 忽略 stale route、memory 或 checklist，重新建立 execution frame。 |
| `goal-realignment` | 繼續前回到 user goal、`.agent-goals` 或 workflow success criteria。 |
| `recovery-specialized` | 走 recovery / escalation workflow，避免 unrelated optimization。 |
| `validation-only` | 只執行 checks 與 source-of-truth 比對，不提前宣稱完成。 |
| `human-facing-summary` | 為使用者整理 options、blockers 與 assumptions，等待決策。 |
| `inspection-only` | 只讀取、搜尋、分析，不寫檔、不 commit、不執行 production actions。 |

## 必要 Routing Steps

1. 用 [`task-routing.md`](task-routing.md) 分類 task shape。
2. 用 [`autonomy-routing.md`](autonomy-routing.md) 檢查 cognitive state 與 autonomy mode。
3. 若 explicit model selection 不可用，用 [`fallback-routing.md`](fallback-routing.md) 做 behavior-only adaptation。
4. 委派 subagent 或 explicit model run 前，先套用 [`multi-model-handoff.md`](multi-model-handoff.md)。
5. 從 [`../compression/README.md`](../compression/README.md) 選擇 compression。

## 禁止主張

- 除非 tool-specific source 證明，否則不得宣稱 provider model、Cursor Auto 或 main chat model 已切換。
- 不得把使用者要求的 model 靜默替換成另一個 model。
- 不得讓 model strategy 覆蓋 source-of-truth、safety、evidence hierarchy 或 user goal。
- Validation scenarios 證明 runtime cost 合理前，不得把 routing primitives promoted 到 `runtime/`。

## Validation Signal

Model-aware routing decision 必須記錄：

```text
Task class:
Cognitive state:
Autonomy mode:
使用的 capability dimensions:
Compression level:
Strategy:
Tool capability:
Fallback behavior:
Validation target:
```
