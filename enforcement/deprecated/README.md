# Deprecated Enforcement Rules

本目錄存放已 deprecated 的 enforcement rules。Deprecation 流程遵循 [`governance/lifecycle/README.md`](../../governance/lifecycle/README.md) 的通用 lifecycle 規則，並加入 enforcement rule 專屬的補充說明。

## Deprecation 流程

Enforcement rule 的 deprecation 分為 4 個階段：

```
標記（Mark）→ 公告（Announce）→ 緩衝期（Buffer）→ 搬移（Move）
```

### 階段 1：標記（Mark）

在原始 rule 檔案與 metadata 中標記為 deprecated：

1. **更新 metadata**：在 `metadata/rules/<rule-name>.yaml` 中設定：
   ```yaml
   status: deprecated
   replaces: enforcement.<new-rule-id>  # 指向取代它的新 rule
   deprecation_date: "2026-05-14"
   removal_date: "2026-08-14"  # 預設 3 個月緩衝期
   ```

2. **更新原始 rule 檔案**：在檔案開頭加入 deprecation notice：
   ```markdown
   # [Deprecated] 原規則名稱

   > ⚠️ 本規則已 deprecated，將於 YYYY-MM-DD 後移除。
   > 請參閱 [`enforcement/<new-rule>.md`](../<new-rule>.md) 取代。
   ```

3. **更新 activation-rules.yaml**：將該 rule 從 `rules` 列表中移除，或標註 `status: deprecated`。

### 階段 2：公告（Announce）

確保所有相關方知道此 deprecation：

1. **更新 `enforcement/README.md`**：在規則索引中標記該規則為 deprecated。
2. **更新 `knowledge/graphs/rules/`**：在 graph 記錄中標記該 rule 的 `status` 為 `deprecated`。
3. **更新 `runtime/router/activation-rules.yaml`**：移除或註解掉 deprecated rule。
4. **檢查 linked updates**：執行 `enforcement/linked-updates.md` 的連動更新檢查。

### 階段 3：緩衝期（Buffer）

預設緩衝期為 **3 個月**，從標記日開始計算。在緩衝期內：

- Deprecated rule 仍可被讀取和使用。
- Agent 應優先使用取代它的新 rule。
- 所有指向舊 rule 的連結應逐步更新為新 rule。
- 若發現遺漏的依賴，延長緩衝期並更新 `removal_date`。

### 階段 4：搬移（Move）

緩衝期結束後：

1. **移動檔案**：將原始 rule 從 `enforcement/` 移至 `enforcement/deprecated/`。
2. **保留 redirect notice**：在原位置建立一個簡短檔案，內容為：
   ```markdown
   # 本檔案已搬遷

   此規則已 deprecated 並搬移至 [`enforcement/deprecated/<rule-name>.md`](deprecated/<rule-name>.md)。
   請使用 [`enforcement/<new-rule>.md`](../<new-rule>.md) 取代。
   ```
3. **清理 metadata**：更新 `metadata/rules/<rule-name>.yaml` 中的 `source_path` 指向新位置。
4. **更新所有索引**：`enforcement/README.md`、`metadata/rules/README.md`、`knowledge/graphs/rules/` 等。

## 狀態轉換圖

```
active ──→ deprecated ──→ removed
              │
              ├── 仍可讀取
              ├── 優先使用取代規則
              └── 3 個月緩衝期
```

| 狀態 | 意義 | 檔案位置 | Metadata status |
|------|------|---------|----------------|
| `active` | 當前有效規則 | `enforcement/<rule>.md` | `validated` 或 `stable` |
| `deprecated` | 即將移除，有替代方案 | `enforcement/<rule>.md`（標記 notice） | `deprecated` |
| `removed` | 已搬移至 deprecated/ | `enforcement/deprecated/<rule>.md` | `deprecated`（source_path 更新） |

## 何時觸發 Deprecation

符合以下任一條件時，應啟動 deprecation 流程：

1. **規則被合併**：兩條規則合併為一條，舊規則不再需要。
2. **規則被取代**：新規則提供了更完整或更正確的 guidance。
3. **規則不再適用**：技術或流程變更使該規則不再相關。
4. **規則被提升**：規則內容被提升到更高層（如 `governance/` 或 `intelligence/`），不再需要 enforcement 層的副本。

## 何時不應 Deprecate

- 規則仍在 active use，且無明確替代方案。
- 規則被多個外部檔案引用，但尚未完成 linked updates。
- 規則的取代規則尚未通過驗證（至少需經過 1 次成功使用）。

## 清理政策（Cleanup Policy）

| 時機 | 動作 |
|------|------|
| 標記 deprecation 時 | 更新 metadata、原始檔案 notice、activation-rules.yaml |
| 緩衝期內 | 逐步更新所有外部引用指向新規則 |
| 搬移至 deprecated/ 時 | 移動檔案、建立 redirect notice、更新所有索引 |
| 最終移除時 | 刪除 redirect notice、確認無任何外部引用指向舊路徑 |

## 與既有文件的關係

- [`governance/lifecycle/README.md`](../../governance/lifecycle/README.md)：通用 lifecycle 規則（本流程是其 enforcement rule 專屬補充）
- [`metadata/schema.md`](../../metadata/schema.md)：`status: deprecated` 已在 schema 中定義
- [`metadata/rules/enforcement-rule-spec.md`](../../metadata/rules/enforcement-rule-spec.md)：Enforcement Rule 專屬 metadata spec
- [`enforcement/linked-updates.md`](../linked-updates.md)：連動更新檢查
- [`knowledge/graphs/rules/README.md`](../../knowledge/graphs/rules/README.md)：Rule dependency graph（deprecation 時需更新）
