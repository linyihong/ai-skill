# root.bootstrap.ai-skill

| 欄位 | 值 |
| --- | --- |
| Atom ID | `root.bootstrap.ai-skill` |
| Source path | [`../../README.md`](../../README.md), [`../../CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md), [`../../enforcement/README.md`](../../enforcement/README.md) |
| Lifecycle | `validated` |
| Summary | Ai-skill 工作的 bootstrap 入口。Root README 定義 OS layout 與 cost-aware 啟動流程；CORE_BOOTSTRAP.md 定義 3 條核心規則（~800 tokens）；enforcement README 定義 Runtime Activation Model 與 lazy-load rules。 |
| When to read | 新 session、接手長對話、使用者要求繼續 Ai-skill 升級、或需要確認 source-of-truth / close-loop 流程時。 |
| Do not use for | 不可取代完整 shared rule 全文；需要實作或修改時仍要讀相關規則全文。 |
| Context cost | ~800 tokens（Core Bootstrap 3 rules） |
| Estimated full cost | ~5000 tokens（若載入所有 lazy-load rules） |
| Validation signal | Root README、CORE_BOOTSTRAP、enforcement README links 可解析；`git status --short --branch` 已檢查；required bootstrap set 已讀。 |
| Last checked | 2026-05-12 |

## Checklist

- 先讀 `CORE_BOOTSTRAP.md`（3 條核心規則，~800 tokens）。
- 再讀 root `README.md`（OS layout）。
- 查詢 `skills-index.yaml` 找到對應 skill。
- 檢查 `runtime/router/activation-rules.yaml` 決定哪些 lazy-load rules 需要 activate。
- 先讀 `knowledge/summaries/` 對應 summary（300-500 tokens）。
- 需要時才展開完整 source。
- Repository 變更必須 diff review、linked updates、commit、push、readback、clean status。
