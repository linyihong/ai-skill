# Dependency Graph Maintenance

`governance/dependency/` 定義知識依賴圖的維護規則。本層不儲存實際的 graph edges（那些在 `knowledge/graphs/`），而是定義 graph 何時需要更新、如何驗證一致性、以及依賴變更時的連動更新流程。

## 核心原則

1. **Graph 是輔助發現工具，不是 source of truth**。`knowledge/graphs/` 的 YAML records 是人工維護的 edge list，不是自動生成的完整 graph。Agent 應將 graph 視為候選清單，需要高信心判斷時仍讀 canonical source。
2. **依賴變更必須更新 graph**。當 source file 被拆分、合併、promotion 或 deprecation 時，對應的 graph edges 必須同步更新或標記 stale。
3. **Graph 一致性由 validation 確保**。`scripts/validate-knowledge-runtime.rb` 與 `scripts/query-knowledge-graph.rb` 可輔助檢查 graph edges 是否指向存在的路徑。
4. **不要為了 graph 而建立 graph**。只有當依賴關係對 agent 的 routing 或閱讀順序有實際影響時，才需要記錄 edge。

## 何時需要更新 Graph

| 事件 | 必要行動 |
| --- | --- |
| 新 Knowledge Atom 建立 | 在對應 graph YAML 或新 graph record 中加入 edges（depends_on / related_to / promotes_from）。 |
| Source file 被拆分 | 更新原 graph record 的 edges，指向新的拆分目標。 |
| Source file 被合併 | 移除指向舊路徑的 edges，加入指向合併目標的 edges。 |
| Promotion 完成 | 加入 `replaces` edge 從新 atom 指向舊 source；加入 `preserves_entrypoint` edge。 |
| Deprecation 啟動 | 加入 `deprecates` edge 或更新舊 source 的 graph record 狀態。 |
| 新 dependency 被發現 | 在現有 graph record 中加入新的 `depends_on` 或 `related_to` edge。 |
| 既有 dependency 不再相關 | 移除或標記 edge 為 `stale`，不要直接刪除（保留變更歷史）。 |

## Graph Record 維護流程

```
1. 識別變更類型（新 atom / 拆分 / 合併 / promotion / deprecation / 新 dependency）。
2. 找出受影響的 graph records：
   ├─ 直接包含變更 source 的 graph record。
   └─ 透過 depends_on / related_to 間接連結的 graph record。
3. 對每個受影響的 graph record：
   ├─ 新增必要的 edges。
   ├─ 移除不再有效的 edges（或標記 stale）。
   └─ 更新 validation 欄位（檢查日期、檢查者）。
4. 執行 graph validation：
   ├─ 所有 target paths 存在且可讀。
   ├─ 沒有 dangling edges（指向不存在的路徑）。
   ├─ 沒有 duplicate edges（同一對 source-target 重複）。
   └─ Edge types 使用 controlled vocabulary（見下方）。
5. 執行 linked updates：
   ├─ 更新 `knowledge/indexes/README.md` 路由。
   ├─ 更新 `knowledge/summaries/` 對應 summary。
   ├─ 更新 `knowledge/runtime/routing-registry.yaml`。
   └─ 更新 layer README 的「目前入口」章節。
6. 執行 validation gates（governance/validation/README.md）。
7. 執行 close-loop（commit / push / readback / clean status）。
```

## Edge Type Controlled Vocabulary

所有 graph edges 必須使用以下 types，不可自創：

| Edge | 意義 | 使用時機 |
| --- | --- | --- |
| `depends_on` | 使用此 atom 前必須先讀 target source。 | 當 target 包含 prerequisite 知識。 |
| `related_to` | Target source 可能有幫助，但不是必讀。 | 當 target 是補充或延伸閱讀。 |
| `conflicts_with` | Source 可能衝突，需要 rule-weight 或 governance resolution。 | 當兩個 source 對同一主題有不同建議。 |
| `replaces` | Promotion 後，新 atom 取代舊 source。 | Promotion 完成時，從新 atom 指向舊 source。 |
| `preserves_entrypoint` | 新分層 path 保留舊 source 可達性。 | 新 atom 建立時，確保舊 entrypoint 仍可到達。 |
| `promotes_from` | Atom 從舊 skill / shared rule 抽取或 promotion 而來。 | Candidate-atom 或 validated-atom 階段，記錄來源。 |
| `routes_to` | Index 或 runtime routing 指向 target source。 | 當 index / registry 新增 routing 條目時。 |
| `stale` | Edge 不再準確，但保留作為變更歷史。 | 當 dependency 不再相關但不確定是否完全移除。 |

## 依賴變更的連動更新

當一個 source 的依賴關係變更時，以下路徑可能需要連動更新：

| 變更類型 | 受影響路徑 |
| --- | --- |
| 新增 depends_on | `knowledge/indexes/README.md`（確保 routing 順序正確）、`knowledge/summaries/`（更新 summary 的 prerequisite 欄位）、`knowledge/runtime/routing-registry.yaml`（更新 dependencies 欄位）。 |
| 移除 depends_on | 同上，移除對應 reference。 |
| 新增 replaces | `governance/lifecycle/README.md`（檢查 promotion gates）、舊 source 加 deprecation note、`knowledge/indexes/README.md`（更新路由優先順序）。 |
| 新增 conflicts_with | `shared-rules/rule-weight.md`（確認衝突解決規則）、`metadata/ranking/README.md`（更新 ranking 優先順序）。 |
| 新增 preserves_entrypoint | 舊 source README 或 entrypoint 文件加 reference note。 |

## Graph Validation

每次 graph record 變更後，至少驗證：

1. **Path existence**：所有 `source` 與 `target` 路徑存在於 repository。
2. **Edge type validity**：只使用 controlled vocabulary 中的 types。
3. **No duplicate edges**：同一對 `(source, type, target)` 不重複出現。
4. **No dangling references**：沒有指向已刪除或已 deprecation 完成的路徑（除非 edge type 是 `stale`）。
5. **Reciprocal consistency**：如果 A `depends_on` B，B 的 graph record 不需要反向 edge，但 agent 應能從 B 找到 A（透過 index 或 summary）。
6. **Validation script**：執行 `ruby scripts/validate-knowledge-runtime.rb` 檢查 generated surfaces 一致性。

## 與其他層的關係

- `knowledge/graphs/README.md`：graph record 格式、edge types、查詢方式。本層定義維護規則，graph 層儲存實際資料。
- `governance/lifecycle/README.md`：lifecycle state 決定 graph edge 的有效性（`promoted` 狀態才能加 `replaces` edge）。
- `governance/validation/README.md`：graph 變更後需通過 validation gates。
- `governance/cleanup/README.md`：cleanup 可能導致 graph edges 需要更新或移除。
- `knowledge/indexes/README.md`：index 與 graph 應一致，graph 的 depends_on 應反映在 index 的 routing 順序。
- `knowledge/runtime/routing-registry.yaml`：registry 的 dependencies 欄位應與 graph 的 depends_on edges 一致。
- `metadata/schema.md`：metadata 的 depends / conflicts 欄位應與 graph edges 一致。
- `shared-rules/linked-updates.md`：graph 變更後的 linked updates 需符合 linked update 規則。
