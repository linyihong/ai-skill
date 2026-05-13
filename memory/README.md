# Memory

`memory/` 負責「長期記憶」。本層保存可重用、可回放、可治理的歷史脈絡設計與記憶分類，不保存專案私有 raw evidence 或 active conversation state。

## 目前入口

- [`working/`](memory/working/README.md) — Session-local 工作記憶（可丟棄）
- [`summary/`](memory/summary/README.md) — 壓縮 session 歷史（≤500 tokens）
- [`decision/`](memory/decision/README.md) — 輕量 ADR（immutable, numbered）
- [`episodic/`](memory/episodic/README.md) — 情境記憶（跨 session 經驗 recall）
- [`project/`](memory/project/README.md) — 專案記憶（跨 session 專案脈絡）
- [`failure/`](memory/failure/README.md) — 失效記憶（抽象化失效模式）

## 放什麼

- Long-term memory、episodic memory 與 experience replay 的設計。
- 可重用 historical context 的分類、使用條件與 validation 方法。
- Failure memory 與 project memory 的抽象化邊界。
- 記憶如何被 `runtime/` 與 `metadata/` 載入的策略。

## 不放什麼

- Active goal、owner、lock、next action；放到 `.agent-goals/`。
- 專案 incident raw logs、tokens、host、private evidence；留在業務專案。
- Feedback lesson 的 promotion workflow；放到 `feedback/`。
- 可執行 shared policy；放到 `shared-rules/`。

## 與既有層的關係

- `shared-rules/failure-learning-system.md` 仍定義 failure learning 的可執行流程。
- `feedback/history/` 保存 skill-specific lesson，成熟後可抽象成 memory 或 intelligence。
- `knowledge/` 管導航與 atom；本層管記憶類型、回放與保存邊界。
- `governance/` 管記憶 lifecycle、deprecation 與清理。
- 冷資料查找由 `knowledge/runtime/sqlite/` 這類 generated lookup cache 處理；本層只定義哪些 historical context 值得長期保留與回放。

## 第一批候選遷移來源

- `shared-rules/failure-learning-system.md` 中的 storage 與 promotion 概念
- `shared-rules/failure-patterns/`
- `feedback/history/` 中可抽象成長期記憶類型的經驗
