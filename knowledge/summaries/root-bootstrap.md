# root.bootstrap.ai-skill

| 欄位 | 值 |
| --- | --- |
| Atom ID | `root.bootstrap.ai-skill` |
| Source path | [`../../README.md`](../../README.md), [`../../shared-rules/README.md`](../../shared-rules/README.md) |
| Lifecycle | `validated` |
| Summary | Ai-skill 工作的 bootstrap 入口。Root README 定義 repository layout 與 reference-first 工作流；shared-rules README 定義 Default Bootstrap 與依任務補讀規則。 |
| When to read | 新 session、接手長對話、使用者要求繼續 Ai-skill 升級、或需要確認 source-of-truth / close-loop 流程時。 |
| Do not use for | 不可取代完整 Default Bootstrap；需要實作或修改時仍要讀相關 shared rule 全文。 |
| Validation signal | Root README 與 shared-rules README links 可解析；`git status --short --branch` 已檢查；required bootstrap set 已讀。 |
| Last checked | 2026-05-11 |

## Checklist

- 先讀 root `README.md`。
- 再讀 `shared-rules/README.md` 的 Default Bootstrap。
- 依任務讀 skill / layer / tool-specific dependencies。
- Repository 變更必須 diff review、linked updates、commit、push、readback、clean status。
