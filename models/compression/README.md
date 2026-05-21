# Model-Aware Compression

`models/compression/` 定義不同 model profiles 使用 summaries、checklists、registry 與 full source 的壓縮策略。目標是降低 context cost，同時不跳過 required dependencies 或 source-of-truth validation。

## Compression Levels

| Level | 使用內容 | 適用情境 | 不可省略 |
| --- | --- | --- | --- |
| `index-only` | `knowledge/indexes/README.md`、`knowledge/runtime/routing-registry.yaml` | 快速定位 primary source 或列出下一步。 | 若要執行變更，仍需讀 primary source。 |
| `summary-first` | Index + registry + `knowledge/summaries/<atom>.md` | 判斷是否需要深入讀 source、低風險問答、handoff 概覽。 | Source-of-truth gate、validation signal。 |
| `checklist-first` | Summary checklist + required validation rules | 小模型執行格式化、檢查或簡短 close-out。 | Required bootstrap、linked updates、diff review。 |
| `source-backed` | Primary source + required dependencies + relevant summaries | 文件修改、規則更新、migration、commit/push 任務。 | Full source 與 validation gate。 |
| `graph-assisted` | Source-backed + graph records + related sources | 跨層衝突、promotion / deprecation、dependency graph 維護。 | Conflict resolution 與 old entrypoint check。 |

## Profile Defaults

| Profile | 預設壓縮層級 | 升級條件 |
| --- | --- | --- |
| `small` | `summary-first` 或 `checklist-first` | 要修改 canonical source、遇到 conflict、缺 validation signal。 |
| `large` | `source-backed` | 需要跨 layer dependency、promotion、deprecation 或 graph reasoning。 |
| `specialized` | `source-backed` + domain workflow | 任務需要 tool adapter、domain technique、live evidence 或 project-specific artifact。 |

## Capability Adjustment

Capability dimensions 可以收緊 compression：

| Capability signal | Compression adjustment |
| --- | --- |
| Low reasoning depth | 只有 bounded low-risk tasks 才優先 checklist-first；其他情況 source-backed。 |
| Low context stability | Edits 前 reread primary source，並縮小 claim scope。 |
| Unknown tool reliability | Tool capability validated 前，避免 close-loop automation。 |
| Medium / high hallucination risk | 使用 source-backed validation 與 evidence hierarchy。 |
| Low compression resilience | 不把 generated reports 或 summaries 當作 source replacement。 |

## Escalation Rules

從壓縮內容升級到 full source 的條件：

- Summary 與 source-of-truth 可能不一致。
- 任務需要修改檔案、commit、push 或 readback。
- 任務涉及 safety、secrets、authorization、source/mirror 或 destructive actions。
- Routing registry 指向 candidate path，但 old entrypoint 仍 active。
- Validation signal 不足以支持結論。

## Output Requirement

使用壓縮策略時，agent 應記錄：

```text
Profile:
Compression level:
Primary source:
Summaries used:
Required full sources:
Deferred sources:
Escalation trigger:
Validation signal:
```

## Boundary

- 壓縮策略不能取代 `enforcement/dependency-reading.md` 的 required reads。
- 壓縮策略不能把 candidate summary 當成 replacement path。
- 壓縮策略不能把 tool mirror 或 generated output 當成 canonical source。
