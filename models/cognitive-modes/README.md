# Cognitive Modes

> ⚠️ **Status**: `documentation-contract trial`（Phase D of plan）
>
> 本層**僅為文件契約**，無 runtime 程式碼、無 `runtime.db` 投影、無 compiler 整合。Agent **手動套用**此 contract 並在 final report 列 Cognitive Mode 報告。執行記錄留在 git history。
>
> 完整實作（Phase 1-5：YAML executable contract、runtime.db `cognitive_modes` 表、discovery heuristics、subsystem 整合、adaptive runtime）見 [`plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md`](../../plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md)。
>
> **Rollback**：本層為 doc-only trial，撤回直接 `git revert` 對應 commit 即可；無 runtime state / 無 schema 變更 / 無 migration。

## 用途

提供 4 維 cognitive mode primitive 的**手動套用契約**，作為將來 runtime 實作的設計驗證階段。每個任務 agent 依任務性質 resolve 4 個 mode，並在 final report 回報實際使用值。

## 4 維 Mode Primitives

### execution_mode

| 值 | 適合任務 |
|------|------|
| `FAST` | typo、簡單問答、單行修改、低風險 |
| `NORMAL` | 標準流程；單檔修改 + validation |
| `DEEP` | 跨層分析、migration、promotion、複雜 debugging |
| `FORENSIC` | Incident analysis、audit、lineage tracing |
| `RECOVERY` | Failure recovery、阻擋寫入、conservative actions |

### context_mode

對應既有 [`models/compression/`](../compression/README.md) 5 級，UPPERCASE 為 runtime primitive、lowercase 為 alias：

| 值（primitive） | Alias（既有 compression） | 載入內容 |
|------|------|------|
| `INDEX_ONLY` | `index-only` | indexes + routing-registry |
| `SUMMARY_FIRST` | `summary-first` | + knowledge/summaries/<atom>.md |
| `CHECKLIST_FIRST` | `checklist-first` | + validation checklist |
| `SOURCE_BACKED` | `source-backed` | + primary source + required dependencies |
| `GRAPH_ASSISTED` | `graph-assisted` | + graph records + related sources |

### governance_mode

| 值 | 啟用 gate set | 額外 action |
|------|------|------|
| `LIGHT` | [sanitization] | — |
| `STANDARD` | [sanitization, language_policy, output_rules] | — |
| `STRICT` | [sanitization, language_policy, output_rules, linked_updates, runtime_surfaces, tool_neutral] | — |
| `LOCKDOWN` | STRICT 全集 | `block_file_writes_until_human_approval` |

### memory_mode

| 值 | 啟用 memory 子層 | 與 retrieval-governance 關係 |
|------|------|------|
| `NONE` | 不查 memory | threshold 不適用 |
| `EPISODIC` | `memory/episodic/` | AND `retrieval-governance.episodic.threshold` |
| `DECISION_REPLAY` | `memory/decision/` | AND `retrieval-governance.decision.threshold` |
| `FAILURE_REPLAY` | `memory/failure/` | AND `retrieval-governance.failure.threshold` |
| `PROJECT_CONTEXT` | `memory/project/` | AND `retrieval-governance.project.threshold` |

## Discovery（手動套用版）

Agent 任務開始時，依**raw signals**快速 resolve 4 個 mode。signal 表完整版見 [plan §Phase 2.1](../../plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md)。

### Mapping 速查（常見組合）

| 任務類型 | execution | context | governance | memory |
|------|------|------|------|------|
| Typo / wording 修正 | FAST | CHECKLIST_FIRST | LIGHT | NONE |
| 單檔修改（notes / memory/working） | NORMAL | SUMMARY_FIRST | LIGHT | NONE |
| 單檔修改（一般 reusable） | NORMAL | SUMMARY_FIRST | STANDARD | NONE |
| 新增 intelligence / failure pattern | DEEP | SOURCE_BACKED | STRICT | NONE |
| 改 enforcement / governance | DEEP | SOURCE_BACKED | STRICT | DECISION_REPLAY |
| 跨層 promotion / migration | DEEP | GRAPH_ASSISTED | STRICT | DECISION_REPLAY |
| Incident audit / lineage | FORENSIC | GRAPH_ASSISTED | STRICT | FAILURE_REPLAY |
| 連續 failure ≥ 2 | RECOVERY | SOURCE_BACKED | LOCKDOWN | FAILURE_REPLAY |
| Long session 跨專案 context | depends | depends | depends | PROJECT_CONTEXT |

## Final Report Cognitive Mode 區塊範本

Agent 完成任務時，final report **必須**含此區塊：

```markdown
### Cognitive Mode 報告

| 維度 | 值 | 理由 |
|------|------|------|
| execution_mode | <FAST/NORMAL/DEEP/FORENSIC/RECOVERY> | <為何選此> |
| context_mode | <INDEX_ONLY/SUMMARY_FIRST/CHECKLIST_FIRST/SOURCE_BACKED/GRAPH_ASSISTED> | <為何選此> |
| governance_mode | <LIGHT/STANDARD/STRICT/LOCKDOWN> | <為何選此> |
| memory_mode | <NONE/EPISODIC/DECISION_REPLAY/FAILURE_REPLAY/PROJECT_CONTEXT> | <為何選此> |
```

理由欄位可引用 raw signal（例：「file diff in enforcement/ → STRICT」「typo 修正 → FAST」）。

## 與其他 models/ 入口的關係

| 既有 | 與 cognitive-modes 的關係 |
|------|------|
| `models/profiles/` | 粗略 small/large/specialized profile 是 cognitive-mode 組合的 backward-compat label；profile 仍可作為 reference doc |
| `models/compression/` | 5 級 compression 對應 `context_mode` 5 級（lowercase alias）；compression 文件是 `context_mode` 的 implementation strategy reference |
| `models/capabilities/` | 5 個 capability dimension 可細化 mode 選擇（例：low context stability → context 偏向 SOURCE_BACKED） |
| `models/routing/` | task-routing / autonomy-routing 是 mode discovery 的 raw signal 之一 |
| `models/governance/` | model selection / hallucination / context budget governance 與 governance_mode compose |
| `models/runtime/` | 將於 Phase 1+ 接收 runtime primitive |
| `models/workflow-adaptation/` | workflow adaptation shape 與 mode 組合 compose |

## Phase D Trial 範圍

本層的「手動套用」trial：

- ✅ 每次 final report 列 Cognitive Mode（已於 commit `db9b515` 之後的 commits 開始實施）
- ⏳ 累積 ≥ 5-10 個任務的 mode 報告，驗證設計
- ⏳ 蒐集「mode 不夠用」「組合誤判」「raw signal 不足」等實證
- ⏳ Phase D completed 後決定：
  - 通過 → 進 Phase 0（Pre-Build Interrogation）+ Phase 1（runtime YAML contract）
  - 不通過 → 修 plan 設計、重啟 Phase D，或撤回整個 plan

## Rollback

| 動作 | 操作 | 實證狀態 |
|------|------|---------|
| 完全撤回 Phase D（單一 commit） | `git revert <Phase D commit>` — 移除本 README 與 plan §Phase D 段落；無 runtime state 變更 | ✅ Clean **僅在後續無修改本 README 或 plan 時**成立 |
| 完全撤回 Phase D + 後續所有修改 | 反向順序逐 commit revert：先 revert 最新，再 revert Phase D commit | ✅ T5 演練驗證；`--abort` 可隨時取消 |
| 暫停手動套用 | 在 plan §Phase D 加 `paused` 標記，agent final report 不再列 Cognitive Mode | ✅ 純文件變更，無風險 |
| 修改 mode 定義 | 編輯本 README + plan 對應 mode 描述；下次任務套用新定義 | ✅ |
| 升級到 Phase 1 runtime 實作 | 本 README 仍保留為 source-of-truth doc；Phase 1 建立 `runtime/cognitive-modes.yaml` 引用本 README | ⏳ Phase 1 開始時驗證 |

### T5 Rollback Dry-Run 實證（2026-05-22）

依 plan §Phase D §T5 設計演練：

```bash
git revert --no-commit 9df20ae   # 嘗試直接 revert Phase D 啟動 commit
# 結果：CONFLICT — models/cognitive-modes/README.md 在 HEAD 被後續 T1 (f98d6e4) 修改
# 衝突類型：UD（deleted in revert target / modified in HEAD）

git revert --abort               # 取消 revert
# 結果：✅ 工作樹回到原狀，無殘留
```

**發現**：當 Phase D 啟動 commit 之後有任何 commit 修改 `models/cognitive-modes/README.md` 或 plan §Phase D 段落，**單一 commit `git revert` 會 conflict**。Conflict 本身可解（接受 HEAD 版或先 revert 後續 commits），但「100% safe rollback」的承諾應改為「依累積修改情況的多步 rollback」。

**推導出的真實 rollback 路徑**（截至 2026-05-22）：

```bash
# 反向順序撤回所有 Phase D 相關 commits：
git revert df37b1a   # T2: README 索引修補
git revert f98d6e4   # T1: stale 描述修正
git revert 9df20ae   # Phase D 啟動

# 或一次三個 revert commit（會產生 3 個 inverse commits）：
git revert df37b1a f98d6e4 9df20ae
```

**`--abort` 安全性**：✅ 已驗證 — 任何 conflict 狀態下 `git revert --abort` 都能清乾淨回到原狀。

**對未來 commits 的建議**：
- 若 commit 改動 `models/cognitive-modes/` 或 plan §Phase D，commit message 列出 rollback 順序
- 若想保留「單一 commit 一鍵 rollback」能力，避免後續 commit 修改 Phase D 相關檔案 — 或在新 commit 中明確標註「此 commit 後 Phase D rollback 需多步」

## 不放什麼

- Runtime 程式碼（屬 Phase 1+，放 `scripts/ai-skill-cli/internal/`）
- Executable YAML contract（屬 Phase 1+，放 `runtime/cognitive-modes.yaml`）
- Runtime state schema（屬 Phase 1+，加進 `runtime.db`）
- Discovery heuristics 機器化規則（屬 Phase 2+，放 `runtime/cognitive-modes-discovery.yaml`）
- Token budget gate 邏輯（屬 Phase 4+）
- Adaptive triggers 機器化（屬 Phase 5+）

## 與 plan 的對應

完整 5 phase 實作藍圖、Open Questions 全部 resolved 內容、ADR Promotion Criteria 見：

→ [`plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md`](../../plans/active/2026-05-22-1629-runtime-cognitive-modes-system.md)

← [回到 models/](../README.md)
