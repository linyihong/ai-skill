# Metadata

`metadata/` 負責「知識控制系統」。本層保存 Knowledge Atom schema、ranking、confidence、compatibility、enforcement rule metadata 與 runtime control metadata，讓知識能被選擇、壓縮、排序、衝突仲裁與治理。

## 目前入口

- [`schema.md`](schema.md)：Knowledge Atom metadata schema v1，包含 required / optional fields、controlled values、YAML template 與 example atom。
- [`rules/`](rules/README.md)：metadata 子規則索引，連到 ranking、confidence、compatibility 與 enforcement rule metadata。
- [`ranking/`](ranking/README.md)：用 priority、status、confidence、context cost 與 depends/conflicts 排定讀取順序。
- [`confidence/`](confidence/README.md)：定義 low / medium / high 信心與 lifecycle state 的關係。
- [`compatibility/`](compatibility/README.md)：記錄 old entrypoint 與 new layer path 的相容狀態。
- [`recovery/`](recovery/README.md)：定義 mismatch escalation 後的 domain-specific reload set、L1-L5 metadata 與 recovery policy。
- [`architecture/`](architecture/README.md)：定義 architecture fit、DDD adoption、overengineering 與 bounded context 的 metadata-only heuristics。

## 放什麼

- Knowledge Atom schema 與 required / optional fields。
- Ranking、priority、confidence、stability、complexity 與 context cost 規則。
- Compatibility、model suitability、depends、related 與 conflicts metadata。
- **Enforcement Rule metadata**：`metadata/rules/enforcement-rule-spec.md` 定義 spec，`metadata/rules/*.yaml` 為各 rule 的 metadata 實例。
- Runtime loading、promotion、cleanup 與 graph construction 所需的控制資料。
- Recovery escalation 後的 domain-specific source-of-truth reload policy。
- Architecture selection 的 fit signal、DDD adoption signal、overengineering signal 與 bounded context heuristics（metadata-only）。
- SQLite / FTS runtime index 的欄位來源，例如 `source_path`、`tags`、`priority`、`confidence`、`context_cost`、`when_to_read` 與 `validation`。

## 不放什麼

- Knowledge Atom 正文、index 或 graph 內容；放到 `knowledge/`。
- 可執行 rule weight 與 dependency reading policy；放到 `enforcement/`。
- Model capability profile 本身；放到 `models/`。
- Tool-specific metadata storage 或 UI；放到 `ai-tools/`。

## 誰會參考這裡（Inbound References）

- [`route.runtime.context-loading`](../knowledge/runtime/routing-registry.yaml:157) — required_dependencies 引用 `metadata/ranking/README.md`、`metadata/confidence/README.md`、`metadata/compatibility/README.md`
- [`route.metadata.knowledge-atom-schema`](../knowledge/runtime/routing-registry.yaml:186) — required_dependencies 引用 `metadata/rules/README.md`、`metadata/ranking/README.md`、`metadata/confidence/README.md`、`metadata/compatibility/README.md`
- [`route.models.model-aware-routing`](../knowledge/runtime/routing-registry.yaml:316) — required_dependencies 引用 `metadata/ranking/README.md`
- [`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) — 多條 records 的 metadata 欄位使用 metadata schema
- **Enforcement Rule metadata** 將被 `runtime/router/activation-rules.yaml` 與 `knowledge/graphs/rules/` 使用。

## 與既有層的關係

- `enforcement/rule-weight.md` 仍是目前規則衝突的可執行 policy。
- `runtime/` 會使用本層資料做 dynamic loading 與 routing，包含 enforcement rule 的 activation_conditions。
- `knowledge/` 會引用本層 schema 讓 atoms 可搜尋、可排序、可治理。
- `knowledge/runtime/sqlite/` 應從本層 schema 取得索引欄位，不另創不相容欄位語意。
- `governance/` 使用 metadata 判斷 lifecycle、cleanup 與 validation。
- `metadata/rules/enforcement-rule-spec.md` 為 enforcement rule 提供結構化 metadata，與 `enforcement/` 正文互補。

## 第一批候選遷移來源

- `plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md` 的 Knowledge Atom metadata 欄位草案
- `metadata/schema.md`
- `enforcement/rule-weight.md`、`decision-efficiency.md` 中可抽象為 metadata 欄位的概念
