# Metadata

`metadata/` 負責「知識控制系統」。本層保存 Knowledge Atom schema、ranking、confidence、compatibility 與 runtime control metadata，讓知識能被選擇、壓縮、排序、衝突仲裁與治理。

## 目前入口

- [`schema.md`](schema.md)：Knowledge Atom metadata schema v1，包含 required / optional fields、controlled values、YAML template 與 example atom。
- [`rules/`](rules/README.md)：metadata 子規則索引，連到 ranking、confidence 與 compatibility。
- [`ranking/`](ranking/README.md)：用 priority、status、confidence、context cost 與 depends/conflicts 排定讀取順序。
- [`confidence/`](confidence/README.md)：定義 low / medium / high 信心與 lifecycle state 的關係。
- [`compatibility/`](compatibility/README.md)：記錄 old entrypoint 與 new layer path 的相容狀態。

## 放什麼

- Knowledge Atom schema 與 required / optional fields。
- Ranking、priority、confidence、stability、complexity 與 context cost 規則。
- Compatibility、model suitability、depends、related 與 conflicts metadata。
- Runtime loading、promotion、cleanup 與 graph construction 所需的控制資料。
- SQLite / FTS runtime index 的欄位來源，例如 `source_path`、`tags`、`priority`、`confidence`、`context_cost`、`when_to_read` 與 `validation`。

## 不放什麼

- Knowledge Atom 正文、index 或 graph 內容；放到 `knowledge/`。
- 可執行 rule weight 與 dependency reading policy；放到 `shared-rules/`。
- Model capability profile 本身；放到 `models/`。
- Tool-specific metadata storage 或 UI；放到 `ai-tools/`。

## 與既有層的關係

- `shared-rules/rule-weight.md` 仍是目前規則衝突的可執行 policy。
- `runtime/` 會使用本層資料做 dynamic loading 與 routing。
- `knowledge/` 會引用本層 schema 讓 atoms 可搜尋、可排序、可治理。
- `knowledge/runtime/sqlite/` 應從本層 schema 取得索引欄位，不另創不相容欄位語意。
- `governance/` 使用 metadata 判斷 lifecycle、cleanup 與 validation。

## 第一批候選遷移來源

- `plans/active/next-stage-upgrade-plan.md` 的 Knowledge Atom metadata 欄位草案
- `metadata/schema.md`
- `shared-rules/rule-weight.md`、`decision-efficiency.md` 中可抽象為 metadata 欄位的概念
