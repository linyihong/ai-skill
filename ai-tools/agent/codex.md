# Codex Agent 規則

Codex 遵守 `ai-tools/agent/` 與其他 adapter 相同的 AI tool rule model。本文件是 Codex 專屬 canonical adapter；root [`AGENTS.md`](../../AGENTS.md) 只作為 bootstrap pointer。

## 啟動流程

Codex 必須依序 bootstrap：

1. 讀取 [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md)。
2. 讀取 [`README.md`](../../README.md)，了解 OS layout。
3. 讀取本文件，取得 Codex 專屬 adapter 行為。
4. 使用 [`runtime/runtime.db`](../../runtime/runtime.db) 作為 runtime source-of-truth。
5. 優先載入 `knowledge/summaries/`，只有 routing 需要時才展開全文。

## Runtime Source-Of-Truth

[`runtime/runtime.db`](../../runtime/runtime.db) 是 runtime configuration 的 canonical source：

- phase state：`phase_machine` / `phases`
- obligations：`obligation_ledger` / `obligations`
- blocking gates：`blocking_gates` / `gates`
- output governance：`language_policy`、`output_rules`、`governance_gates`
- canonical runtime documents：`runtime_config_documents`
- projection tables 與 generated executable surfaces

不要提交 `runtime/**/*.yaml` mirror。若 runtime document 必須改變，應透過 runtime compiler 或 approved runtime update flow 更新 SQLite canonical copy。

## Executable Contract Boundary

Codex 必須分清楚 ownership 與 runtime execution：

- governance、enforcement、workflow、ai-tools、metadata/rules 的 YAML 留在 owner layer。
- YAML 只有在 `runtime_projection.enabled: true` 時才成為 executable runtime surface。
- compiler 會將這些 contract 投影到 `runtime.db.generated_surfaces`。
- Markdown 是人類可讀說明；YAML 是 activation contract。

未來更新時使用這個 boundary：

| Source type | Owner location | Runtime behavior |
| --- | --- | --- |
| Governance principle | `governance/` Markdown | 僅供 reference |
| Governance executable flow | `governance/**/*.yaml` | enabled 時投影 |
| Enforcement prose rule | `enforcement/*.md` | 人類可讀規則正文 |
| Enforcement activation contract | `metadata/rules/*.yaml` 或 companion YAML | enabled 時投影 |
| Workflow procedure | `workflow/**/*.yaml` | enabled 時投影 |
| AI tool onboarding / adapter flow | `ai-tools/**/*.yaml` | enabled 時投影 |
| Runtime internal config | `runtime/runtime.db` | SQLite canonical only |

## Decision Promotion

不要把每個 decision 都 promotion 成 ADR，必須依內容選 target：

| Decision content | Target |
| --- | --- |
| 可執行規則或 cross-agent policy | `enforcement/` |
| reasoning heuristic、signal、tradeoff、anti-pattern | `intelligence/` |
| 操作流程或 repeatable flow | `workflow/` |
| runtime gate、activation、obligation、policy surface | `runtime/runtime.db` |
| 架構級不可逆 decision | `constitution/ADR-*` |
| session-scoped replay decision | `memory/decision/` |
| project-specific decision | `<PROJECT_ROOT>/docs/decisions/` |

正式 ADR layer 是 `constitution/`，不是 `decisions/`。Runtime decision-recording 在 SQLite 中使用 logical path `runtime/constitution/decision-recording.yaml`。

## Codex 編輯規則

修改本 repo 時：

- 優先遵守既有 repo pattern，不另建平行結構。
- 搜尋優先使用 `rg` / `rg --files`。
- 手動修改檔案使用 `apply_patch`。
- 不要 revert 使用者變更，除非使用者明確要求。
- root `AGENTS.md` 只能是 thin adapter bootstrap；可重用 AI-tool 規則放在 `ai-tools/`。
- 若變更會影響 execution，必須跑 update flow：compile、refresh、validate；使用者要求時再 commit / push。

## 驗證

提交 runtime 或 rule 變更前，執行：

```powershell
scripts\ai-skill-cli\bin\ai-skill-windows-amd64.exe runtime compile --repo . --native-compiler --json
scripts\ai-skill-cli\bin\ai-skill-windows-amd64.exe runtime refresh --repo . --json
scripts\ai-skill-cli\bin\ai-skill-windows-amd64.exe runtime validate --repo . --json
```

如果 Go compiler code 或 binaries 有變更，也執行：

```powershell
cd scripts\ai-skill-cli
go test ./...
```
