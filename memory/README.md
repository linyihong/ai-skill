# Memory

`memory/` 負責「長期記憶」。本層保存可重用、可回放、可治理的歷史脈絡設計與記憶分類，不保存專案私有 raw evidence 或 active conversation state。

Memory 是 selective replay system，不是 always-loaded context。每次 replay 前都必須通過 retrieval trigger、qualification、freshness / scope check 與 replay budget；memory-derived conclusion 不得取代 canonical source。

## 目前入口

- [`retrieval-governance/`](retrieval-governance/README.md) — selective retrieval、activation、replay budget、freshness、contamination 與 promotion policy
- [`working/`](working/README.md) — Session cognition buffer（可丟棄）
- [`summary/`](summary/README.md) — 壓縮 session 歷史（≤500 tokens）
- [`decision/`](decision/README.md) — 輕量 ADR（immutable, numbered）
- [`episodic/`](episodic/README.md) — 情境記憶（跨 session 經驗 recall）
- [`project/`](project/README.md) — 專案記憶（跨 session 專案脈絡）
- [`failure/`](failure/README.md) — 失效記憶（抽象化失效模式）

## 放什麼

- Long-term memory、episodic memory 與 experience replay 的設計。
- 可重用 historical context 的分類、使用條件與 validation 方法。
- Failure memory 與 project memory 的抽象化邊界。
- 記憶如何被 selective retrieval、qualification 與 activation 的策略。

## 不放什麼

- Active goal、owner、lock、next action；放到 `.agent-goals/`。
- 專案 incident raw logs、tokens、host、private evidence；留在業務專案。
- Feedback lesson 的 promotion workflow；放到 `feedback/`。
- 可執行 shared policy；放到 `enforcement/`。
- Canonical source、current truth 或 runtime execution state。

## 誰會參考這裡（Inbound References）

- [`route.decisions.adr`](../knowledge/runtime/routing-registry.yaml:696) — candidate_sources 引用 `memory/decision/README.md`
- [`enforcement/failure-learning-system.md`](../enforcement/failure-learning-system.md) — 定義 failure memory 的 storage 與 promotion 規則
- [`plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md`](../plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md) — 引用 memory/ 的設計概念
- [`plans/active/2026-05-20-1745-memory-retrieval-activation-governance.md`](../plans/active/2026-05-20-1745-memory-retrieval-activation-governance.md) — 定義 selective cognitive replay system 的執行計畫

## 與既有層的關係

- `enforcement/failure-learning-system.md` 仍定義 failure learning 的可執行流程。
- `feedback/history/` 保存 skill-specific lesson，成熟後可抽象成 memory 或 intelligence。
- `knowledge/` 管導航與 atom；本層管記憶類型、回放與保存邊界。若內容不依賴特定 incident 也成立，應 promotion 到 `knowledge/`、`intelligence/`、`workflow/` 或 `governance/`。
- `governance/` 管記憶 lifecycle、deprecation 與清理。
- 冷資料查找由 `knowledge/runtime/sqlite/` 這類 generated lookup cache 處理；本層只定義哪些 historical context 值得長期保留與回放。
- `runtime/` 可查 memory route / lookup metadata，但不得保存 raw historical memory 或 active execution contract。

## 第一批候選遷移來源

- `enforcement/failure-learning-system.md` 中的 storage 與 promotion 概念
- `enforcement/failure-patterns/`
- `feedback/history/` 中可抽象成長期記憶類型的經驗
