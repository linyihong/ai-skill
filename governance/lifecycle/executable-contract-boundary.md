# Executable Contract Boundary

本文件是 Ai-skill YAML 放置政策的 canonical 規章。它定義哪些 YAML 是 owner-layer executable contract、哪些 YAML 只是 metadata 或 runtime config，以及影響執行的 contract 如何進入 `runtime/runtime.db`。

## 核心規則

Source ownership 留在擁有該概念的 layer。Runtime execution surface 只投影到 `runtime/runtime.db`。

```text
owner-layer Markdown / YAML contract
  -> runtime compiler
  -> runtime.db generated_surfaces / projection tables
```

不要因為 governance、enforcement 或 workflow source 會影響執行，就把它搬進 `runtime/`。`runtime/` 仍然是 runtime engine 與 SQLite registry boundary。

## 責任地圖

| 檔案或 table | 責任 | 不負責 |
| --- | --- | --- |
| `governance/lifecycle/executable-contract-boundary.md` | 人類可讀的 YAML placement policy 與 source-of-truth boundary。 | 維護逐檔 inventory 狀態。 |
| `governance/lifecycle/executable-contract-boundary.yaml` | 給 agent 執行的 placement / projection gates。 | 取代人類可讀 policy。 |
| `metadata/executable-contract-schema.md` | executable contract 欄位 schema。 | 決定 YAML 應放在哪個 layer。 |
| `governance/lifecycle/executable-contract-inventory.yaml` | Inventory state：`contract_exists`、`contract_required`、`markdown_only`、`not_applicable`。 | 重新定義 placement policy。 |
| `runtime/runtime.db generated_surfaces` | owner-layer executable contracts 與 deterministic generated surfaces 的 projection。 | 擁有 governance、enforcement、workflow 或 ai-tools source content。 |
| `runtime/runtime.db generated_surfaces` 中的 executable contracts | Contract-backed activation 與 execution surface。 | 取代 owner-layer source 或維護第二份 rule body。 |
| `runtime/runtime.db runtime_config_documents` | runtime-owned config documents 的 canonical copy。 | 擁有非 runtime 的 governance、enforcement、workflow 或 ai-tools contracts。 |

## 何時需要 YAML

當 agent 必須把文件當成 workflow 或 gate 執行時，該文件需要 YAML contract。

需要 YAML 的訊號：

- Ordered steps
- Trigger 或 activation conditions
- Required reads 或 dependencies
- `depends_on` relationships
- Exit conditions
- Blocking gates
- Required evidence
- Failure actions
- Final status / report requirements

如果文件只說明哲學、背景、tradeoff 或設計理由，保持 Markdown-only；除非後續 workflow 從中抽出可執行 gates。

## 放置規則

| Source 類型 | YAML contract 位置 | Runtime projection |
| --- | --- | --- |
| Governance lifecycle flow | `governance/**/*.yaml` | 影響執行時必須投影 |
| Enforcement policy contract | `enforcement/**/*.yaml` | 影響執行時必須投影 |
| Workflow execution flow | `workflow/**/*.yaml` | 影響執行時必須投影 |
| AI tool adapter flow | `ai-tools/**/*.yaml` | 影響執行時必須投影 |
| Rule metadata | `metadata/rules/*.yaml` | 除非明確 promotion 成完整 executable contract，否則不投影 |
| Runtime internal config | `runtime.db` canonical documents | 已由 runtime 擁有 |
| Philosophy / rationale / ADR | Markdown only | 不需要 |

`metadata/rules/*.yaml` 預設是 metadata。只有在明確例外情況下，且包含完整 executable contract schema、execution-bearing fields 與 `runtime_projection.enabled: true`，才可以承載 executable contract。

## Runtime Projection 規則

會影響 agent 執行的 YAML contract 必須包含：

```yaml
runtime_projection:
  enabled: true
  target_key: governance.example.contract
  surface: generated_surfaces
```

Compiler 只投影明確設定 `runtime_projection.enabled: true` 的 contract。這避免一般 metadata、graph 與 validation YAML 變成 runtime noise。

Projection 不會轉移 ownership。若 `enforcement/authorization-scope.yaml` 被投影到 `generated_surfaces`，source 仍然是 `enforcement/authorization-scope.yaml`；`runtime.db` 裡的 row 只是 compiled runtime surface。

## Contract-backed Activation 規則

Enforcement、governance、workflow 與 ai-tools 的 activation 由 owner-layer executable YAML contract 定義，並投影到 `runtime/runtime.db generated_surfaces`。Runtime 不再維護 `activation_rules` / `activation_rules_mirror` tables 作為 enforcement lazy-load lookup。

修改 executable contract placement、runtime projection 或 framework source-of-truth 前，先執行 [`../../workflow/software-delivery/requirements/pre-build-interrogation.md`](../../workflow/software-delivery/requirements/pre-build-interrogation.md)。若無法說明 canonical owner、projection boundary、duplicate surface removal / deprecation 與 validation target，不得開始 migration。

當 owner-layer executable contract 已存在且包含 `activation` 時，activation path 採 contract-first：

```text
contract activation
  -> owner-layer executable YAML contract
  -> companion Markdown for rationale and maintenance context
```

Markdown 透過 contract 的 `source_markdown` 或 `required_sources` 引用；不得另建 runtime lookup table 維護同一份 activation semantics。

## Schema 規則

新的 executable contract 應遵守 [`../../metadata/executable-contract-schema.md`](../../metadata/executable-contract-schema.md)。Metadata YAML 不是 executable contract；除非它定義 `contract_type`、`blocking_level`、`activation`、execution-bearing fields，且設定 `runtime_projection.enabled: true`。

## Agent 規則

當 Markdown 文件要求某流程必須作為 workflow 執行時，agent 必須先載入 companion YAML contract，再使用 Markdown 作為解釋與維護脈絡。

當 runtime surface 與 owner-layer executable contract 同時描述同一個 activation 時，owner-layer executable contract 對 execution semantics 具有權威性。Runtime row 是 projection，不是第二份 source。

如果文件有 executable signals 但沒有 YAML contract，agent 必須把它視為 linked-update gap，建立 contract 或記錄為何不適用。

## Contract Inventory

Active inventory 位於 [`executable-contract-inventory.yaml`](executable-contract-inventory.yaml)。它把目前文件標記為：

- `contract_exists`：YAML 已存在。
- `contract_required`：下一步需要 YAML contract。
- `markdown_only`：刻意保持非 executable。
- `not_applicable`：template、example、deprecated stub 或非 owner source。

[`executable-contract-boundary.yaml`](executable-contract-boundary.yaml) 保存 executable boundary gates。[`executable-contract-inventory.yaml`](executable-contract-inventory.yaml) 保存目前 inventory decisions；需要更新清單時，應更新 inventory 文件，不要擴寫 boundary contract 的 seed inventory。

## 相關文件

- [`knowledge-update-flow.yaml`](knowledge-update-flow.yaml)
- [`executable-contract-inventory.yaml`](executable-contract-inventory.yaml)
- [`compiler-philosophy.md`](compiler-philosophy.md)
- [`../../runtime/README.md`](../../runtime/README.md)
- [`../../scripts/ai-skill-cli/internal/app/runtime_compiler.go`](../../scripts/ai-skill-cli/internal/app/runtime_compiler.go)
