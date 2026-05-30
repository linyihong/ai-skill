# Documentation Discipline Slice（Developer notes + Feedback tips + Backfill rules）

> **Cognitive Slice**：`apk-documentation-discipline`（從 [`../artifact-gates.md`](../artifact-gates.md) §13+§14+§15 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。
>
> 本 slice 由 [Scheme A vs B probe](../../../validation/scenarios/software-delivery/slice-load-scenario-ag-schemes-a-vs-b.yaml) 證實獨立存在的經濟性高於併入 evidence-chain：T1 標準分析任務省 ~58 行 / 次；T2 doc-discipline 審閱省 ~168 行 / 次。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-documentation-discipline` |
| `purpose` | 規範 developer guidance notes 寫法、feedback lesson 模板、與 backfill rules（分析完成後回填規則） |
| `type` | `execution` |
| `tags` | artifact-gate, documentation, backfill |
| `load_when` | 撰寫 / 審閱 developer guidance、feedback lessons、或為既有專案做文件回填 |
| `do_not_load_when` | 純分析執行中、尚未到文件撰寫階段、純 evidence capture |
| `owner_layer` | workflow |
| `layer_justification` | 規定「怎麼寫好 dev notes、怎麼寫 feedback lesson 才有 8 個必要欄位、分析後要回填到哪些 owner location」的 documentation ordering / artifact gate；通過 workflow membership test，與 evidence-chain（怎麼記錄證據）為不同 cognitive phase（怎麼寫好文件） |
| `canonical_source` | 本檔（原 `artifact-gates.md` §13 Developer Guidance Notes + §14 Feedback Lesson Writing Tips + §15 Backfill Rules） |
| `dependencies` | `apk-evidence-chain`（feedback / backfill 引用 evidence）、`apk-feature-handoff`（backfill 對象之一） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | post-analysis writeup 或 retrospective 任務應載入本 slice |

## 13. Developer Guidance Notes（可選）

若分析結果能轉成「未來開發自家 App 時可採取的設計、實作或安全做法」，可在專案分析文件加一小節：

```markdown
## Developer Guidance Notes

| Observation | Development Guidance | Owner | Validation |
| --- | --- | --- | --- |
| 已去敏觀察 | 可重用的開發建議 | client / API / backend / build / monitoring | 測試或 review 方法 |
```

這一節只寫已去敏、可泛化的開發啟發。成熟後把 App 開發 guidance 回饋到 `app-development-guidance`；本 `apk-analysis` 文件只保留分析方法、證據鏈與工具判斷。

## 14. Feedback Lesson Writing Tips

寫入 `feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md` 時，請避免只有工具名與短結論。每條技巧都應包含：

- `One-line Summary`：一句話講重點。
- `Human Explanation`：給人看的背景與誤判風險。
- `Trigger`：什麼現象會觸發這條技巧。
- `Evidence`：去敏證據或觀察。
- `Generalized Lesson`：抽象後的通用規則。
- `Agent Action`：下次 AI 要採取的具體行動。
- `Applies When` / `Does Not Apply When`：適用邊界。
- `Validation`：怎麼確認這條技巧有效。

好的 lesson 應該像這樣：

```markdown
### Proxy failure 要先拆成導流與 TLS 兩層

One-line Summary:
代理看不到明文時，先確認「有沒有進代理」，再談憑證或 pinning。

Human Explanation:
很多人看到 MITM 沒有明文就直接判斷是 pinning。更可靠的順序是先看 App 是否真的連到 proxy。如果仍直連目標 host，問題在導流或初始化時機；如果已經進 proxy 才 TLS failed，才查 CA / pinning。

Agent Action:
先檢查 CONNECT 或 connect target，不要先寫 pinning 結論。
```

## 15. Backfill Rules

每次分析完成後：

- UI Behavior 必須回填專案 UI 行為入口或 page-level map：記錄 entry path、可見 UI blocks、App sort label、tap/swipe/input 操作、API/data source 對照、截圖/UI hierarchy/live replay/hook 證據與 unknowns。若沒有 UI 證據，明確標 `needs capture`、`needs replay` 或 `Trigger confidence: low`。
- 目標 API 結論回填專案 API 文件。
- 解碼規則回填協議/解密文件。
- SDK 或 client 行為回填 BDD / tests。
- 若分析文件要用來做 app 工具、SDK、client、mock、fixture-driven implementation、contract test 或重建功能，同輪自動啟用 `app-development-guidance` 並交出 Feature Reconstruction Handoff；不要讓開發規格停留在 APK 分析文件內。
- 通用技巧回填 `feedback_history/<category>/` 或 `feedback_history/common/`（新檔），驗證後再整理到主文件或對應 `techniques/<category>/`。
- App 開發 guidance 回填 `app-development-guidance`；不要把產品開發 checklist 長期堆在 `apk-analysis`。

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
