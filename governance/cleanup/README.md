# Duplicate Cleanup & Splitting

`governance/cleanup/` 定義知識重複偵測、拆分規則與所有權邊界。本層不儲存實際的 duplicate records，而是提供判斷流程與治理規則，讓 agent 在遷移、promotion 或日常維護時能系統化處理重複內容。

## 核心原則

1. **Duplicate 不等於錯誤**。在遷移過渡期，舊 `skills/` 與新分層可能暫時持有相似內容。只要舊 entrypoint 仍 active，重複是預期狀態。
2. **Cleanup 只在 promotion 或 deprecation 階段執行**。`candidate-map` 或 `candidate-atom` 階段不清理舊 source，避免破壞尚未完成的遷移。
3. **Splitting 優先於 dedup**。當一個文件包含多個責任（例如 workflow + intelligence 混雜），應先拆分再清理重複。
4. **Ownership boundary 決定誰是 canonical**。重複內容的 canonical source 由 lifecycle state 與 layer 責任邊界決定，不靠內容相似度。

## Duplicate 類型

| 類型 | 說明 | 處理策略 |
| --- | --- | --- |
| **Cross-layer duplicate** | 同一知識同時存在於 `skills/` 與新分層（如 `analysis/`、`workflow/`、`intelligence/`）。 | 保留舊 entrypoint 直到 promotion 完成；promotion 後在舊 source 加 deprecation note 指向新路徑。 |
| **Intra-layer duplicate** | 同一層內兩個文件描述相同概念（如兩個 `intelligence/` atoms 都討論 retry strategy）。 | 合併為單一 atom，更新所有指向舊路徑的 index / summary / graph / registry。 |
| **Cross-skill duplicate** | 兩個不同 skill 包含相似的操作說明或知識（如 `apk-analysis` 與 `app-development-guidance` 都提到 proxy 設定）。 | 抽取共用部分到 `enforcement/` 或 `intelligence/`，兩邊 skill 保留 reference。 |
| **Source vs summary drift** | `knowledge/summaries/` 的 summary 與 `skills/` 或新分層的 canonical source 不一致。 | 執行 generated refresh checklist，決定 refresh / revalidate / downgrade。 |
| **Tool mirror duplicate** | `ai-tools/` 文件複製了 `enforcement/` 或 `skills/` 的可執行規則。 | 移除 tool mirror 中的重複規則，改為 reference 指向 canonical source。 |

## 偵測流程

```
1. 發現疑似重複
   ├─ 同一知識出現在多個路徑？
   ├─ 兩個文件的標題 / 核心原則高度相似？
   └─ Summary 與 source 內容不一致？

2. 判斷類型
   ├─ Cross-layer → 檢查 lifecycle state
   ├─ Intra-layer → 檢查 ownership boundary
   ├─ Cross-skill → 檢查是否可抽取共用規則
   ├─ Source vs summary → 執行 generated refresh
   └─ Tool mirror → 確認是否為 reference-only

3. 決定行動
   ├─ 保留（過渡期預期重複）
   ├─ 合併（intra-layer 或可抽取共用規則）
   ├─ 更新（summary 過時）
   ├─ 移除 mirror（改為 reference）
   └─ 標記 deprecation（promotion 完成後）
```

## Splitting 規則

當一個文件包含多個層的責任時，應拆分而非保留混合文件：

| 混合情況 | 拆分目標 |
| --- | --- |
| Workflow 步驟 + 判斷智慧 | Workflow → `workflow/`，Intelligence → `intelligence/` |
| 分析步驟 + 工具命令 | Analysis → `analysis/`，Tools → `skills/<name>/TOOLS.md` |
| 操作流程 + 產出規範 | Workflow → `workflow/`，Artifact gates → `workflow/<domain>/artifact-gates.md` |
| 領域知識 + 執行腳本 | Knowledge → `knowledge/` 或 `intelligence/`，Script → `scripts/` |

### Splitting 門檻

符合以下任一條件時應考慮拆分：

- 單一文件超過 200 行且包含兩個以上層的內容。
- 文件的 README 或目錄描述提到「本文件同時包含 X 與 Y」。
- Agent 在讀取時只用到文件的前半部，後半部從未被 routing 命中。
- 文件被多個 layer 的 index 或 registry 同時指向。

### Splitting 流程

```
1. 識別文件中每個區塊的責任層（analysis / workflow / intelligence / tools / knowledge）。
2. 為每個責任層建立新文件在對應目錄。
3. 保留原文件作為 entrypoint，內容改為 reference-first 格式（連結到各拆分目標）。
4. 更新 index、summary、graph、registry 與 layer README。
5. 標記原文件為 `candidate-map` 狀態。
6. 經過 promotion gates 後，原文件可標為 `deprecated`。
```

## Ownership Boundary

當兩個層都聲明同一知識的所有權時，使用以下規則決定 canonical source：

| 知識類型 | Canonical layer | 其他層的行為 |
| --- | --- | --- |
| 可執行操作步驟 | `workflow/` | `intelligence/` 只能引用，不複製步驟 |
| 判斷原則與 heuristics | `intelligence/` | `workflow/` 可引用判斷原則，不複製原理 |
| 分析方法與工具用法 | `analysis/` | `workflow/` 可引用分析步驟，不複製方法細節 |
| 事實性知識 | `knowledge/` | 各層可引用，不複製事實正文 |
| 可執行規則 | `enforcement/` | 各層可引用，不複製規則正文 |
| 工具設定 | `ai-tools/` | 不複製中央庫內容，只 reference |

## 清理執行流程

當 cleanup 被觸發（promotion、deprecation、或定期維護）：

```
1. 列出受影響的所有路徑（舊 source、新 layer、index、summary、graph、registry）。
2. 對每個路徑判斷：
   ├─ 保留（仍在過渡期）
   ├─ 合併（intra-layer duplicate）
   ├─ 更新（summary / graph / registry stale）
   ├─ 加 deprecation note（promotion 完成）
   └─ 刪除（僅限已 promotion + deprecation 完成 + 無外部 reference）
3. 執行 linked updates（index、summary、graph、registry、layer README、roadmap）。
4. 執行 validation gates（governance/validation/README.md）。
5. 執行 close-loop（commit / push / readback / clean status）。
```

## 與其他層的關係

- `governance/lifecycle/README.md`：cleanup 的觸發時機由 lifecycle state 決定（promoted / deprecated 才可清理）。
- `governance/validation/README.md`：cleanup 完成後需通過 validation gates。
- `metadata/`：ranking、confidence、compatibility 可輔助判斷哪個 source 是 canonical。
- `knowledge/graphs/`：graph edges 可輔助發現跨層 duplicate（多個 source 指向同一 target）。
- `knowledge/indexes/README.md`：cleanup 後需更新 index 移除或重新指向已變更的路徑。
- `enforcement/dependency-reading.md`：cleanup 過程中的文件讀取仍受 dependency reading 約束。
- `enforcement/linked-updates.md`：cleanup 後的 linked updates 需符合 linked update 規則。
