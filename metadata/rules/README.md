# Metadata Rules

`metadata/rules/` 是 metadata 操作規則索引，說明 `metadata/schema.md` 的欄位如何被 routing、governance、summaries，以及未來 graph / runtime surfaces 使用。

## 規則集合

| 規則 | 用途 |
| --- | --- |
| [`ranking/`](../ranking/README.md) | 多個 atom 或 source 都相關時，決定先讀哪一個。 |
| [`confidence/`](../confidence/README.md) | 描述證據強度，以及 atom 何時可從 candidate 移到 validated 或 stable。 |
| [`compatibility/`](../compatibility/README.md) | 新分層演進時，保留舊 skill entrypoints 與工具相容性。 |

## 邊界

- Metadata rules 不覆蓋可執行 shared rules。
- Metadata 與 shared rules 衝突時，依 `shared-rules/rule-weight.md`。
- Metadata 可以降低 context loading cost，但不能跳過 required dependency reading。
- 舊 `skills/` source files 在 lifecycle promotion gates 通過前，仍是 source of truth。
- Metadata 文件正文預設使用繁體中文；英文保留給欄位名、enum、路徑、YAML key 與必要專有名詞。
