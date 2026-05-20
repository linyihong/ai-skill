# Metadata Compatibility

`metadata/compatibility/` 定義 migration 進行中時，新知識分層如何保留既有 skill、enforcement rule、tool 與 script 行為。

## 相容性欄位

使用 `metadata/schema.md` 的下列欄位描述相容性：

| 欄位 | 用途 |
| --- | --- |
| `source_path` | 仍具權威性的 canonical old source path。 |
| `depends` | 必須先讀的 old entrypoints 或 enforcement rules。 |
| `related` | Candidate new layer paths 或 supporting references。 |
| `replaces` | 只有 promotion 或 deprecation approval 後才可填寫。 |
| `conflicts` | 需要解決的潛在 rule 或 entrypoint conflicts。 |
| `governance_notes` | Migration state、compatibility notes 或 deprecation requirements。 |

## 相容性狀態

| 狀態 | 意義 |
| --- | --- |
| `old-entrypoint-active` | Legacy source 仍是 active。舊 `skills/` scaffold 已退役，不應再用此狀態表示 skills source。 |
| `dual-reference` | Old entrypoint 與 new layer path 都已連結，供 discovery 使用。 |
| `new-layer-promoted` | New layer path 已受支援，但 old path 仍可解析。 |
| `deprecation-planned` | Old path 已有 replacement 與 deprecation note，但仍存在。 |
| `old-entrypoint-retired` | Old path 在 validation 與 replacement 後已移除或封存。 |

## 必要相容性備註

任何 candidate map 或 promoted atom 都要記錄：

- Old entrypoint。
- New reference path。
- Old entrypoint 是否仍 active。
- Tool-specific discovery 是否仍依賴 old path。
- 哪些 validation 證明 links 與 routing 仍可運作。

## 阻塞條件

符合下列情況時，不可 promote 或 deprecate：

- 某工具仍只能載入 old skill path，且沒有 adapter。
- 變更後尚未讀回 old entrypoint。
- `knowledge/indexes/README.md` 只指向 candidate path，卻漏掉 active source。
- Enforcement rule 指明 old source 是 canonical，且未同步更新規則。
- Link check 或 close-loop validation 失敗。

## Reference-First 預設

Compatibility metadata 應優先使用 direct canonical repository references。Tool mirrors、bundles、copied snapshots 與 local runtime paths 都是部署面，不是 source paths。
