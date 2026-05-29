# 外科手術式修改規則（Surgical Changes Rules）

> **Cognitive Slice**：`sd-surgical-caveats`（從 [`execution-flow.md`](execution-flow.md) §9 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-surgical-caveats` |
| `purpose` | 進行外科手術式小改時，控制 diff 純度、避免 scope creep 與 orphan 污染 |
| `type` | `failure`（diff-purity / scope-discipline caveat，非 execution-order） |
| `tags` | caveat, surgical, diff-purity |
| `load_when` | 進行外科手術式小改、需控制 diff 純度 / orphan、修改既有程式碼 |
| `do_not_load_when` | 大型新功能初始實作、純分析 / evidence-only 任務 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「修改既有程式碼時做什麼 / 不做什麼」的 procedure 與紀律 gate，通過 workflow membership test；不承載 evidence 取得方法（非 analysis），也不論證長期模式（非 intelligence） |
| `canonical_source` | 本檔（原 `execution-flow.md` §9.1–9.5） |
| `dependencies` | [`examples/EXAMPLES.md`](examples/EXAMPLES.md) §3（具體範例，預設 suppress） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Phase 4 Scenario C（mixed：debug 失敗 deployment，需 surgical caveat） |

當修改既有程式碼時，遵循以下規則以最小化 diff 和避免引入不相關的變更。參見 [`examples/EXAMPLES.md`](examples/EXAMPLES.md) §3 的具體範例。

## 9.1 只改必須改的行

- 只修改解決問題所需的行。不要順便 refactor 不相關的函式、區塊或檔案。
- 如果發現不相關的 code smell，記錄為獨立 issue 或 TODO，不要在本次變更中一併修改。
- 例外：如果不相關的 code 會直接導致本次變更無法正確測試或驗證，則可一併修改，但必須在 commit message 中明確說明。

## 9.2 匹配既有 code style

- 不要改變既有程式碼的風格（quote style、indentation、命名慣例、type hints、docstring 格式）。
- 如果既有 code 使用單引號，新 code 也使用單引號。如果既有 code 沒有 type hints，不要加 type hints。
- 如果既有 code 有 docstring，新 code 也加 docstring，且格式一致。如果既有 code 沒有 docstring，不要加 docstring。
- 不要重新格式化既有 code 的 whitespace、換行或括號位置。

## 9.3 不要順便加「順便」的功能

- 不要加「既然來了就順便」的 validation、error handling、logging、caching 或 notification。
- 不要加 speculative features（「以後可能會用到」的參數、選項、抽象層）。
- 不要加「最佳實踐」風格的改進（Strategy pattern、Factory、DI container），除非需求明確要求。

## 9.4 只清理自己的 orphan

- 如果本次變更產生了未使用的 import、變數或函式，清理它們。
- 不要清理本次變更之前就存在的 unused code，除非它直接阻擋編譯或測試。

## 9.5 驗證 diff 純度

在標記完成之前，檢查 diff：

```
git diff --stat          # 確認只有預期的檔案被修改
git diff                 # 逐行檢查是否有不相關的變更
```

如果 diff 包含超過解決問題所需的變更，還原不相關的部分。
