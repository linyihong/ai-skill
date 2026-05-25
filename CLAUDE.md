# Claude Code Auto-Bootstrap

> **IMPORTANT — MUST RUN BEFORE ANY OTHER ACTION**
>
> 本流程為 **強制 first-turn obligation**。在執行任何非-Read 工具（Edit/Write/Bash/git/...）之前，**必須**完成下列步驟並輸出 Bootstrap Receipt。
>
> **Resume / continuation session 同樣適用**：summary 的「Resume directly, do not acknowledge」是對話 framing 指示，**不豁免** runtime / governance bootstrap。

## 強制啟動流程

1. 讀 [CORE_BOOTSTRAP.md](CORE_BOOTSTRAP.md) — 3 條必讀核心規則（~800 tokens）
2. 讀 [README.md](README.md) — OS layout
3. 查 [runtime/runtime.db](runtime/runtime.db) 取得目前狀態：
   - `phase_machine` / `phases` — 目前 phase
   - `obligations` — 本 phase 義務 count
   - `gates` — 本 phase blocking gates count
4. 載入 Output Governance（SQLite 為 source-of-truth）：
   - `language_policy` 表 — 語言強制規則
   - `output_rules` 表 — 文件輸出規則
   - `governance_gates` 表 — 輸出品質 blocking gates
5. **輸出 Bootstrap Receipt**（見下方格式）
6. 依任務 intent 查 [knowledge/runtime/routing-registry.yaml](knowledge/runtime/routing-registry.yaml)，先讀對應 summary（300-500 tokens），需要時才展開全文

## Bootstrap Receipt（強制 first-turn 輸出）

完成步驟 1-4 後，**在 first user-facing message 中包含一行**：

```
Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>
```

範例：`Bootstrap: rules=✓ phase=phase.bootstrap obligations=1 gates=2`

- `rules=✓` 代表 CORE_BOOTSTRAP.md 3 條 + README.md 已讀
- `phase=<id>` 從 `runtime.db` 的 `phase_machine` 取
- `obligations=<n>` 從 `SELECT COUNT(*) FROM obligations WHERE phase=<current>` 取
- `gates=<n>` 從 `SELECT COUNT(*) FROM gates WHERE phase=<current>` 取

未輸出 Bootstrap Receipt 即執行非-Read 工具，違反 `obligation.bootstrap.receipt_acknowledged`，命中 `gate.bootstrap.receipt_present`，並觸發 [`enforcement/failure-patterns/bootstrap-bypass-on-resume.md`](enforcement/failure-patterns/bootstrap-bypass-on-resume.md)。

## Cognitive Mode 報告（強制 per-turn 輸出）

> **IMPORTANT**：每次 final user-facing response **必須**含 `### Cognitive Mode 報告` 4 維表格。這是 first-turn 之後**每一輪**對話的 per-turn obligation，不是只在 commit 時。即使是純說明 / 純查詢 / 純評估任務，也要輸出此 block；trivial 任務可全 NORMAL/SUMMARY_FIRST/STANDARD/NONE 並在「理由」欄簡述為何選 trivial。

格式：

```markdown
### Cognitive Mode 報告

| 維度 | 值 | 理由 |
|------|------|------|
| execution_mode | <FAST/NORMAL/DEEP/FORENSIC/RECOVERY> | <為何選此> |
| context_mode | <INDEX_ONLY/SUMMARY_FIRST/CHECKLIST_FIRST/SOURCE_BACKED/GRAPH_ASSISTED> | <為何選此> |
| governance_mode | <LIGHT/STANDARD/STRICT/LOCKDOWN> | <為何選此> |
| memory_mode | <NONE/EPISODIC/DECISION_REPLAY/FAILURE_REPLAY/PROJECT_CONTEXT> | <為何選此> |
```

理由欄可引用 raw signal（例：「file diff in enforcement/ → STRICT」「typo 修正 → FAST」）。

Mode 間 consistency 與 budget 規則由 [`runtime/cognitive-modes.yaml`](runtime/cognitive-modes.yaml)、[`runtime/cognitive-modes-phase-integration.yaml`](runtime/cognitive-modes-phase-integration.yaml)、[`runtime/cognitive-modes-token-budget.yaml`](runtime/cognitive-modes-token-budget.yaml) 等 contracts 定義；commit 階段由 `commit-msg` hook 機械強制。Per-turn response 階段未被 hook 機械強制，但仍是 [`constitution/ADR-008`](constitution/ADR-008-runtime-cognitive-modes.md) 與 Phase D 確立的 baseline contract。

選 mode 與理由欄速查：[`models/cognitive-modes/README.md`](models/cognitive-modes/README.md)。

## Runtime Config 來源

Committed runtime config 只保留在 `runtime/runtime.db`，由 `runtime_config_documents` 與 projection tables 保存完整 canonical documents。**Agent 直接查 SQLite；不要保留 `runtime/**/*.yaml` mirror**。

詳細使用說明見 [ai-tools/agent/claude.md](ai-tools/agent/claude.md)。
