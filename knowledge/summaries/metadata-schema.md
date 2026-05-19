# metadata.schema.knowledge-atom

| 欄位 | 值 |
| --- | --- |
| Atom ID | `metadata.schema.knowledge-atom` |
| Source path | [`../../metadata/schema.md`](../../metadata/schema.md) |
| Lifecycle | `validated` |
| Summary | Knowledge Atom metadata schema v1，定義 atom 的必填欄位、選填欄位、受控值、YAML 範本、驗證規則與 provider prompt cache hints。 |
| When to read | 建立或評估 Knowledge Atom、summary、graph record、runtime registry record，或需要判斷 ranking / confidence / compatibility / provider cache 欄位時。 |
| Do not use for | 不可用 metadata 覆蓋可執行 enforcement rules；規則衝突仍依 `enforcement/rule-weight.md`。 |
| Validation signal | Schema 欄位可套用到 routing registry 與第一批 summary atoms；Markdown links 可解析。 |
| Last checked | 2026-05-11 |

## Checklist

- 每個 atom 至少要有 `id`、`title`、`type`、`layer`、`source_path`、`summary`、`domains`、`tags`、`status`、`priority`、`confidence`、`stability`、`context_cost`、`when_to_read`、`validation`。
- `source_path` 必須指向 canonical repository path。
- `context_cost.cacheable` 表示 runtime / conversation 內可重用；provider prompt cache eligibility 需使用 `context_cost.provider_cache.provider_cache_candidate`。
- `stable` 需要真實使用或 review 的 validation evidence。
