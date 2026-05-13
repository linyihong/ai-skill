# Feedback Promotion Pipeline

`feedback/promotion/` 定義 feedback lesson 如何從 skill-local history 推進到 `workflow/`、`intelligence/`、`shared-rules/`、`memory/` 或 runtime surfaces。它是 promotion design layer，不取代 `shared-rules/feedback-lessons.md` 的可執行寫作規則。

## Source Of Truth

| Layer | Role |
| --- | --- |
| `feedback/history/<domain>/` | Lesson 全文與歷史紀錄的 source-of-truth。 |
| `shared-rules/feedback-lessons.md` | Lesson 命名、模板、索引與 agent 行為規則。 |
| `shared-rules/failure-learning-system.md` | Agent failure 或 close-loop gap 的分類與 promotion target。 |
| `governance/lifecycle/README.md` | Candidate、validated、promoted、deprecated 的 lifecycle gate。 |
| `governance/validation/README.md` | Promotion、generated refresh 與 close-loop validation gate。 |

## Promotion Targets

| Lesson type | Target | Gate |
| --- | --- | --- |
| Single-skill technique | `skills/<skill>/WORKFLOW.md`、`TOOLS.md`、`DOCUMENTATION.md` 或 `techniques/` | Lesson 已 generalized、去敏，且 skill index 已更新。 |
| Engineering judgment | `intelligence/` | 影響 trade-off、anti-pattern、route selection 或 cross-project decision。 |
| Execution flow | `workflow/` | 影響 agent 如何執行 planning、review、handoff 或 validation。 |
| Cross-skill or all-repo rule | `shared-rules/` 或 `shared-rules/failure-patterns/` | Failure class 或 prevention gate 可跨 skill 重演。 |
| Runtime navigation | `knowledge/`、`metadata/`、`runtime/` | 需要被 registry、summary、graph 或 model context report route 到。 |
| Long-term lesson memory | `memory/` | 需要保留 replay / episodic / project abstraction boundary。 |

## Promotion Checklist

1. 保留原 `feedback/history/<domain>/` lesson；不要刪除或覆寫歷史。
2. 檢查 lesson 是否只含 generalized rule、適用條件與 validation，不含 project incident raw evidence。
3. 決定最小 durable target：skill workflow、intelligence、workflow、shared-rule、memory 或 runtime surface。
4. 若 promotion 變成 runtime route，更新 `knowledge/indexes/README.md`、`knowledge/runtime/routing-registry.yaml`、summary、graph 與 generated reports。
5. 若 source lesson 或 old skill entrypoint 仍 active，在新 layer 明確標 `candidate` / `dual-reference` / `old-entrypoint-active`。
6. 執行 lints、Markdown link check、generated runtime validation、close-loop dry run、commit、push、readback 與 clean status。

## Downgrade Or Rework

將 promotion 降級或退回時：

- 保留原 lesson，新增修訂或 deprecation note，不默默刪除。
- 若 validation signal 不再成立，將 summary / graph / registry confidence 降級或標 `needs refresh`。
- 若發現混入 project-specific evidence，移回 project docs，並依 `reusable-guidance-boundary.md` 重寫 generalized rule。
- 若是 agent failure prevention，不只退回 skill lesson；同時檢查 `failure-learning-system.md` 是否需要 shared failure pattern。

## Runtime Outputs

目前 runtime surfaces 會把 promotion pipeline 當成 design route，而不是自動 promotion engine。修改 promotion docs 或已進 runtime 的 feedback-derived atom 時，需重新執行：

```bash
ruby scripts/generate-model-context-report.rb --write
ruby scripts/generate-knowledge-runtime-report.rb --write
ruby scripts/validate-knowledge-runtime.rb
```
