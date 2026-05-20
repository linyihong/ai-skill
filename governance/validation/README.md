# Knowledge Validation Gates

`governance/validation/` 定義新 AI-native 分層變更的 validation gates。它補充 `enforcement/` 的可執行規則，但不取代 enforcement rules。

**治理與 enforcement 的分工（語言為例）：**可重用文件的語言一致性、低爭議用語等，在概念上屬於知識治理；**可執行條文**寫在 [`enforcement/neutral-language.md`](../../enforcement/neutral-language.md)，本目錄則以 **Neutral language** gate 規定「何時必須依該檔自檢」，把治理目標落實到 PR 閉環。

## 必要 Gates

| Gate | 必要檢查 | 適用時機 |
| --- | --- | --- |
| Source boundary | 確認編輯的是 canonical repository paths，不是 tool mirrors 或 runtime copies。 | 每次 Ai-skill writeback。 |
| Retired entrypoint | 確認退役入口已標明替代 source，且 active routing 指向 `workflow/`、`analysis/`、`intelligence/`、`enforcement/`、`ai-tools/` 或 `scripts/` 的現行入口。 | Candidate maps、atom promotions、migrations。 |
| Linked updates | 依 [`linked-update-governance.md`](../ai-runtime-governance/linked-update-governance.md) 檢查受影響 README、roadmap、index、metadata、generated surface 與 source entry files。 | 任何新 layer path 或 routing change。 |
| Metadata | 引入或 promotion Knowledge Atom 時，確認 `metadata/schema.md` 欄位已存在。 | Candidate atom、validated atom、promoted atom。 |
| Navigation | 應可發現新路徑時，確認 `knowledge/indexes/README.md` 可 route 到該路徑。 | Routing surfaces 與 promoted reference paths。 |
| Generated refresh | Source 變更時，確認 summaries、graphs、registry records 是否需要 refresh、revalidate 或 downgrade。 | Source-of-truth 文件、metadata、routing registry、summaries、graphs 變更。 |
| Runtime report generation | 執行 `ruby scripts/generate-knowledge-runtime-report.rb --write` 產生 deterministic runtime report。 | Routing registry、refresh policy、summaries 或 graphs 變更。 |
| Model context report generation | 執行 `ruby scripts/generate-model-context-report.rb --write` 產生 model-aware context loading report。 | Routing registry model 欄位、model profiles 或 compression strategy 變更。 |
| Model checklist generation | 執行 `ruby scripts/generate-model-checklists.rb --write` 產生 per-model context-loading checklist。 | Routing registry model 欄位、required dependencies 或 model docs 變更。 |
| SQLite runtime index boundary | 執行 `ruby scripts/validate-runtime-sqlite-index.rb`，確認 SQLite / FTS 只作 generated lookup cache，DB 可重建、被 git ignore、source checksum 未 stale，且不取代 canonical source。 | SQLite generator、query helper、feedback lesson index 或 runtime lookup cache 變更。 |
| Runtime refresh orchestration | 執行 `ruby scripts/refresh-knowledge-runtime.rb` 一鍵重建 reports / SQLite index 並跑 validators。 | 多個 generated runtime surfaces 可能同時 stale 時。 |
| Knowledge runtime helper | 執行 `ruby scripts/validate-knowledge-runtime.rb` 檢查 generated surfaces。 | Routing registry、refresh policy、summaries 或 graphs 變更。 |
| Link check | 解析 touched docs 的 Markdown links。 | Documentation changes。 |
| Neutral language | 依 [`enforcement/neutral-language.md`](../../enforcement/neutral-language.md) 檢查可重用文件：繁體中文正文、英文限路徑／指令／符號與必要專有名詞、標題與摘要中性化。 | 變更 `enforcement/`、`workflow/**/README.md`、`analysis/**/README.md`、`intelligence/**/README.md`、根 `README.md`、根 `CONTRIBUTING.md`、`governance/` 下可重用說明、模板或 onboarding 類 Markdown。 |
| Lints | 執行 touched files 可用的 lints。 | Documentation 或 code changes。 |
| Diff review | 檢查 secrets、private hostnames、raw evidence、local absolute paths 與 unrelated changes。 | Commit 前。 |
| Close-loop dry run | 執行 `./scripts/ai-skill-close-loop.sh` 確認 dirty path grouping。 | Commit 前。 |
| Commit / push / readback | Commit、push、讀回 changed entries，並確認 `git status --short --branch` clean。 | Ai-skill repository updates。 |

## Migration Validation Checklist

任何從舊 `skills/` content 移到新分層的工作，都使用此 checklist：

```text
Goal:
- What user-visible or runtime outcome does this change support?

Source:
- Old source path:
- New candidate/promoted path:
- Source-of-truth state:

Linked updates:
- Layer README:
- Knowledge index:
- Metadata:
- Roadmap:
- Old entrypoint:

Validation:
- Lints:
- Markdown links:
- Diff review:
- Close-loop dry run:
- Commit/push/readback:
- Clean status:
```

## Generated Refresh Checklist

當 source-of-truth 文件、metadata、summary、graph 或 routing registry 變更時，檢查：

```text
Changed source:
- Path:
- Change type: source / metadata / summary / graph / registry / roadmap

Generated surfaces:
- Summaries affected:
- Graph records affected:
- Registry records affected:
- Runtime docs affected:

Decision:
- Refresh now:
- Revalidate only:
- Downgrade confidence:
- No update needed because:

Validation:
- Runtime refresh orchestrator:
- Runtime report regenerated:
- Model context report regenerated:
- Model checklist regenerated:
- SQLite runtime index regenerated or marked not implemented:
- YAML parse:
- Knowledge runtime helper:
- Markdown links:
- Source paths still canonical:
- Old entrypoints still reachable:
- Roadmap / durable status updated:
```

## Pass / Block Rules

- 若 old entrypoints 失效，變更 blocked。
- 若 promoted atom 缺 metadata，變更仍只能是 `candidate-map` 或 `candidate-atom`。
- 若 links 失敗，commit 前必須修正。
- 若 validation 無法執行，記錄 blocker，且不可將 lifecycle state 標為 promoted。
- 若 generated summary / graph / registry 可能 stale，必須 refresh、revalidate 或 downgrade confidence；不可假裝仍 current。
- 若變更是 reference-only 且未使用 tool mirror，tool sync 不適用。

## 與 Enforcement Rules 的關係

- Dependency reading、canonical writeback、commit / push / readback 與 clean status 仍由 `enforcement/dependency-reading.md` 管理。
- Linked update requirements 仍由 `enforcement/linked-updates.md` 管理。
- Rule priority 仍由 `enforcement/rule-weight.md` 管理。
- 本檔提供新 knowledge system 的 architecture-layer validation shape。
