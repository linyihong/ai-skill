# Knowledge Graphs

`knowledge/graphs/` 描述 Knowledge Atoms、source files、skills、shared rules 與 runtime routing surfaces 之間的關係。目前本目錄保存 graph record 格式與第一批 candidate records，不生成完整自動 graph。

## 目前 graph records

| Graph record | 用途 | 狀態 |
| --- | --- | --- |
| [`source-boundary.yaml`](source-boundary.yaml) | 連接 active goal / durable roadmap 邊界、content layering 與 governance lifecycle。 | `candidate` |
| [`metadata-navigation.yaml`](metadata-navigation.yaml) | 連接 metadata schema、metadata 子規則、knowledge index、runtime registry 與 summaries。 | `candidate` |
| [`apk-analysis-pilot.yaml`](apk-analysis-pilot.yaml) | 連接 `skills/apk-analysis/` 舊入口與 analysis / workflow / intelligence 候選目的地。 | `candidate` |

## Graph 目的

Graphs 協助 agent 理解：

- 必讀 dependencies。
- Related sources。
- Conflicts。
- Replacement 與 deprecation paths。
- 舊 skills 到新分層的 promotion flow。

## Edge Types

未來 graph records 使用下列 edge labels：

| Edge | 意義 |
| --- | --- |
| `depends_on` | 使用此 atom 前必須先讀 target source。 |
| `related_to` | Target source 可能有幫助，但不是必讀。 |
| `conflicts_with` | Source 可能衝突，需要 rule-weight 或 governance resolution。 |
| `replaces` | Promotion 後，新 atom 取代舊 source。 |
| `preserves_entrypoint` | 新分層 path 保留舊 source 可達性。 |
| `promotes_from` | Atom 從舊 skill / shared rule 抽取或 promotion 而來。 |
| `routes_to` | Index 或 runtime routing 指向 target source。 |

## Graph Record 格式

```yaml
id:
source:
edges:
  - type:
    target:
    reason:
    validation:
status: candidate
```

## 相容性規則

- 使用 canonical repository-relative paths 或 atom IDs。
- 不把 tool mirror paths 建模為 canonical sources。
- 若 graph 使用 `replaces`，lifecycle state 必須已是 promoted 或 deprecated。
- Candidate maps 應使用 `preserves_entrypoint`，不要使用 `replaces`。

## 新增規則

- 新 graph record 必須能解析所有 source / target path。
- Graph record 不可包含 secrets、project incident evidence、本機絕對路徑或 tool mirror source。
- 若 source 改動，graph record 需要 revalidate 或降級 confidence。
- Graph 只描述關係；可執行規則仍以 `shared-rules/` 與 active source-of-truth 文件為準。
- Source、summary、registry 或 lifecycle state 改動時，依 [`../runtime/refresh-policy.yaml`](../runtime/refresh-policy.yaml) 判斷是否 refresh、revalidate 或 downgrade。
