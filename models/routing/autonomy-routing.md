# Autonomy Routing

Autonomy routing 把 cognitive state 映射到 model-aware execution behavior。State semantics 來自 [`../../governance/ai-runtime-governance/cognitive-state-governance.md`](../../governance/ai-runtime-governance/cognitive-state-governance.md)；本檔只負責轉成 execution strategy。

## Cognitive State Mapping

| Cognitive state | Strategy | 必要行為 |
| --- | --- | --- |
| `STABLE` | `execution-heavy` | 使用 bounded edits、normal validation 與 concise closeout。 |
| `UNCERTAIN` | `validation-heavy` | 先讀 primary source、收集 evidence，避免 broad patching。 |
| `DEGRADED` | `source-backed` | 降低 autonomy，使用 source-of-truth 加 validation gates。 |
| `CONTAMINATED` | `rediscovery-only` | 不重用 stale route、memory、checklist 或 prior conclusion。 |
| `MISALIGNED` | `goal-realignment` | 重新連接 current action 到 user goal、`.agent-goals` 或 workflow criteria。 |
| `RECOVERY` | `recovery-specialized` | 遵循 recovery workflow，避免 unrelated improvement work。 |
| `VALIDATION_REQUIRED` | `validation-only` | 執行 checks 並比對 source；有 evidence 前不宣稱完成。 |
| `ALIGNMENT_REQUIRED` | `human-facing-summary` | 整理 options 與 blockers，等待 user decision。 |
| `READ_ONLY` | `inspection-only` | 只分析；不寫檔、不 commit、不執行 production actions。 |

## Autonomy Downgrade Triggers

- 沒有 new evidence 卻反覆 patch。
- 用 low-scope evidence 支撐 high-scope claims。
- Source-of-truth ambiguity 或缺 canonical authority。
- Tool output 與 live observation 或 owner contract 矛盾。
- Context compaction 保留 conclusions，但缺 evidence。
- Memory 或 prior route 跨 task boundary 重用，且未 revalidation。

## Recovery To Execution

只有在下列條件成立後，才從 downgraded autonomy 回到 execution：

1. Old assumptions 已 downgraded 或 invalidated。
2. Source-of-truth 已讀取，或標記 missing / not applicable。
3. Validation target 明確。
4. Evidence scope 與 next claim 匹配。
5. Selected strategy 不再依賴 contradicted beliefs。

## Model Routing Boundary

Autonomy routing 可以選擇更嚴格的 behavior shape。除非 tool-specific adapter 確認可控制 model，否則不得宣稱已選用更強 provider model。
