# Ai-skill Glossary

> Canonical shared-language source for Ai-skill framework vocabulary. 第一批只收 framework / runtime semantics / cognitive vocabulary / architecture contracts。Domain glossary（`software-delivery.md` / `runtime.md` / `validation.md`）延後。
>
> 上游 plan：[`plans/active/2026-05-25-1000-context-language-glossary-system.md`](../../plans/active/2026-05-25-1000-context-language-glossary-system.md)
>
> Validator 機械強制：`ai-skill glossary validate`（[command contract](../../scripts/ai-skill-cli/docs/command-contract.md#ai-skill-glossary-validate)）。

## 用途與邊界

| 是 | 不是 |
| --- | --- |
| Ai-skill framework / runtime / cognitive / architecture 共享語彙的 canonical 定義 | 業務領域百科 |
| 每個詞條有單一 `owner-layer`、單一 canonical meaning | 跨 owner 重複定義同一詞 |
| 其他 layer（workflow / intelligence / memory / plan / ADR / runtime YAML）引用 owner entry | 在 workflow / intelligence 內 inline redefine |
| Markdown canonical（human readable + git diffable + PR friendly） | 第二份 YAML / SQLite source-of-truth |
| Phase 5 SQLite projection 是 derived index | 把 projection 當 canonical |

## Vocabulary Resolution Priority

當多個 source 對同一詞有不同說法時，agent 必須依下列順序解析：

1. Active project canonical docs（僅限該專案語境）
2. `knowledge/glossary/`
3. Accepted ADRs / `constitution/`
4. Workflow docs
5. Intelligence heuristics
6. Memory replay（永遠不能 canonical）

限制：

- Project docs 只能覆蓋 project-local usage，不得改寫 Ai-skill framework term。
- Memory replay 永遠不能覆蓋 glossary / ADR / workflow current source。
- 舊 ADR 若與 active runtime docs 衝突，必須檢查 Framework Generation / Vocabulary Evolution section。

## Semantic Owner Domains

每個 entry 必須宣告 `owner-layer`，限定為下列其中之一：

| Owner | 語義範疇 | Status |
| --- | --- | --- |
| `runtime-cognition` | Runtime cognitive mode core、deterministic control plane | canonical |
| `semantic-routing` | Routing registry、route activation、graph traversal | canonical |
| `workflow-orchestration` | Workflow phase machine、execution flow、gate sequencing | canonical |
| `validation-governance` | Validation scenarios、governance gates、enforcement rules | canonical |
| `memory-replay` | Memory layer semantics、replay boundary、staleness policy | canonical |
| `runtime-projection` | SQLite projection、generated surfaces、derived index | canonical |
| `architecture-contracts` | Owner-layer YAML contracts、ADR-level architecture boundary | canonical |
| `ecosystem-adaptation` | Economics / pressure / adaptation / telemetry signals | candidate（為 [economics plan](../../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) 預留） |
| `runtime-economics` | Runtime-owned cost / budget / compression contracts | candidate（僅在 economics plan Phase 0 確認 `runtime/economics/` 路徑時啟用；否則降為 deprecated 並由 `ecosystem-adaptation` 承載） |

Owner 擁有該詞的 canonical meaning；其他 layer 只能引用、alias 或標記 local usage，不得重新定義。同一詞在不同 context 有不同 meaning 時，必須拆 bounded context 或改名。

## Entry Status Set

```
canonical       — 正式生效
candidate       — 已收錄但 owner / scope 未最終確認
deprecated      — 已退役，引用前必須改用 replaced_by 指向的 term
superseded      — 被 replaced_by 取代，保留歷史 anchor
alias-only      — 僅限既有 deprecated entry 過渡使用；新 entry 禁用（用 aliases: 欄位）
experimental    — Gen 3 演化期實驗條目
project-local   — 專案層級條目（不在 ai-skill.md，僅供 templates / project memory 引用）
```

## Entry Physical Shape

每個 entry 由兩部分組成：

1. `## <snake_case_term>` H2 heading（提供 human anchor + GitHub link）
2. 緊接一個 ```` ```yaml ```` code block（提供 machine-readable schema）

H2 文字必須與 YAML block 的 `term:` 完全相同。

### Required fields

- `term` — snake_case canonical name
- `status` — 上述 status set 之一
- `meaning` — 一段話的 canonical definition
- `affects` — string array，影響的 layer / workflow / contract paths
- `owner-layer` — 上述 owner domains 之一

### Optional fields

- `aliases` — string array，已知的別名（kebab-case 自動由 validator 視為 alias，不需列出）
- `anti-meaning` — disambiguator，「這個詞不是另一個聽起來像它的詞」（純人類閱讀，validator 不引用）
- `excludes` — string array，「這個詞不涵蓋語義鄰居 X 的職責」；每個字串必須為現有 entry 的 `term`（validator 強制）
- `related-terms` — list of `{type, target}`，定義 directed relation
- `introduced-by` — 來源 plan / ADR，格式 `plans/<path>` 或 `constitution/ADR-XXX.md`
- `deprecated-by` — 取代來源，同 introduced-by 格式

### Forbidden patterns

不得在 entry body 出現：

- Project-specific hosts、paths、class / test names、sample IDs、incident evidence
- Temporary implementation detail / current runtime state
- Raw historical transcript / memory replay
- Executable contract semantics（屬於 workflow / enforcement / runtime source）
- Commit SHA、issue 編號、PR URL（會 rot）

## Naming Convention

- `term:` canonical 寫法為 **snake_case**（例：`context_mode`、`cognitive_cost`）
- kebab-case（例：`context-mode`）由 validator 自動列為 alias，不需手動寫入 `aliases:`
- Markdown H2 heading 使用 snake_case 與 `term:` 完全一致

## Alias 表達方式

Alias 不建獨立 entry。所有 alias 寫在 canonical term 的 `aliases:` 欄位（flat string array）。Validator 強制：

- `aliases:` 中的字串不得出現為其他 entry 的 `term:`
- Alias chain 不得形成 cycle
- `status: alias-only` 在新 schema 下禁用（已由 `aliases:` 取代）；僅保留 enum 供既有 deprecated entry 過渡

## Relation Lifecycle

`related-terms:` 內每個 relation 必須宣告 `type:` 與 `target:`。Allowed types：

```yaml
allowed_relation_types:
  - alias_of
  - related_to
  - conflicts_with
  - owned_by
  - used_by
  - deprecated_by
  - replaced_by
  - derived_from
  - aggregates
```

`derived_from` 與 `aggregates` 為 [economics plan](../../plans/active/2026-05-27-1557-tool-runtime-signal-economics-integration.md) §Split cost model 預留（`cognitive_cost` aggregates `thinking_cost` + `context_cost` + `execution_cost` + `knowledge_cost`）。

### Symmetry Classification（machine-readable）

```yaml
relation_symmetry:
  alias_of:        { symmetry: asymmetric, reverse_field: aliases }
  related_to:      { symmetry: symmetric,  reverse_field: related_to }
  conflicts_with:  { symmetry: symmetric,  reverse_field: conflicts_with }
  owned_by:        { symmetry: asymmetric, reverse_name: owns }
  used_by:         { symmetry: asymmetric, reverse_name: used_in }
  deprecated_by:   { symmetry: asymmetric, reverse_name: deprecates }
  replaced_by:     { symmetry: asymmetric, reverse_name: replaces }
  derived_from:    { symmetry: asymmetric, reverse_name: derives }
  aggregates:      { symmetry: asymmetric, reverse_name: aggregated_into }
```

Validator 強制：

- **Symmetric relation**（`related_to` / `conflicts_with`）必須在兩端 entry 雙向出現
- **Asymmetric relation** 僅在 source entry 寫一次；反向關係由 Phase 5 SQLite projection 自動 derive

## 反向關係查找（人類路徑）

Phase 5 projection 上線前後皆可用的 fallback：

1. **CLI（Phase 5 後）**：
   ```bash
   ai-skill glossary inspect <term>
   ```
   輸出 term 的 outgoing relations + 由 SQLite 即時 derive 的 incoming relations。

2. **SQL（Phase 5 後）**：
   Projection DB 位置：`knowledge/runtime/sqlite/runtime-index.sqlite`。
   例如：
   ```sql
   SELECT source_term FROM glossary_relations
   WHERE target_term = 'compression' AND relation_type = 'alias_of';
   ```

3. **Markdown 肉眼（Phase 5 前唯一可用）**：
   ```bash
   grep -rn "target: compression" knowledge/glossary/
   # 找誰指向 compression
   grep -rn "^## compression$" knowledge/glossary/
   # 直接看 compression 自己的 entry
   ```
   去 target term 的 entry 看 `related-terms` 與 `aliases` 欄位。

## Usage Index Source Types

Phase 5 `glossary_usage` 表的 `source_type` 允許值：

```
workflow / validation / runtime / knowledge / adr / plan / memory
```

Phase 1 usage index 只收 declared references（`related-terms` / `affects` / `excludes`）。Repo-wide textual scan 只能產生 candidate relations，不能直接 canonicalize（短詞如 `mode` / `context` / `runtime` 會造成 false semantic relationships）。

## Drift Detection Categories

Validator 與 Phase 5 projection 必須能偵測：

- **Duplicate meaning** — 兩個 entries 不同 term 但 meaning 高度重疊
- **Conflicting ownership** — 同一概念在兩個 owner-layer 各有 canonical entry
- **Alias loop** — `A.aliases includes B`、`B.aliases includes A`
- **Deprecated term resurrection** — 已 `deprecated` / `superseded` term 重新出現為 `canonical`
- **Near-duplicate concept fork** — 兩個 entries 同 owner-layer 但 meaning 不同（owner boundary 內部分裂）

## Entry Example（worked）

````markdown
## context_mode

```yaml
term: context_mode
status: canonical
owner-layer: runtime-cognition
meaning: >
  Runtime control plane 對 context expansion strategy 的 deterministic
  enum。決定 agent 在 task entry 時以何種深度載入 source / summary /
  graph / checklist。是 cognitive mode 6 維 vector 的其中一維。
affects:
  - runtime/cognitive-modes.yaml
  - runtime/cognitive-modes-discovery.yaml
  - knowledge/runtime/routing-registry.yaml
aliases:
  - ctx_mode
anti-meaning: >
  不是 LLM context window 的 token budget（那是 token_budget_pressure 或
  context_cost）；不是 IDE / editor 的 "context menu"。
excludes:
  - discovery_mode
related-terms:
  - { type: related_to,  target: execution_mode }
  - { type: related_to,  target: governance_mode }
  - { type: used_by,     target: cognitive_cost }
introduced-by: plans/archived/2026-05-22-1629-runtime-cognitive-modes-system.md
```
````

`excludes: [discovery_mode]` 表示「決定載入哪些檔案是 `discovery_mode` 的職責，不在本詞範圍」。validator 會檢查 `discovery_mode` entry 確實存在。

## Semantic Governance（暫存於本檔）

未來拆出 `governance/semantic/` 的 promotion triggers（任一發生即考慮拆層）：

- 累積 ≥ 3 次 drift incidents（duplicate meaning / conflicting ownership / alias explosion）
- Deprecation lifecycle 需要跨 release 階段（目前只有單階段 `deprecated`）
- Semantic migration 工具化（多筆 entries 批次改 owner / status）
- Owner domain set 出現結構性 conflict 或需動態擴張

未達 trigger 前，semantic governance rules 保留在本 README。

## Validator Surface

`ai-skill glossary validate` 檢查（完整契約見 [command-contract.md](../../scripts/ai-skill-cli/docs/command-contract.md#ai-skill-glossary-validate)）：

1. H2 heading + YAML block 配對；H2 文字 = `term:`
2. Required fields 完整
3. `status` / `owner-layer` / `relation type` 屬於 allowed enum
4. `term:` 為 snake_case
5. `aliases:` 規則（不可為他 term、無 cycle、新 entry 禁 `alias-only`）
6. `introduced-by` / `deprecated-by` 形狀為 `plans/<path>` 或 `constitution/ADR-XXX.md`
7. `excludes:` 字串為現有 term
8. Symmetric relation 雙向出現

Exit `30` (`validation_failed`) 並列出每個 violation；通過 exit `0` 並回報計數。

## 與其他規則的關係

- 詞條被 workflow / intelligence / runtime YAML 引用時，引用方不得 inline redefine — 必須連回本目錄
- Memory 引用詞條只能作為 replay 標籤，不得覆寫 canonical meaning
- 新增 entry 必須通過 validator；schema 變更必須先更新本 README、再改 validator、再加 entry
