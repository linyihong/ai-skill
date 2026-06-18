# Closure Slice（Definition of Ready / Definition of Done / Feed Back Reusable Lessons）

> **Cognitive Slice**：`sd-closure`（從 [`execution-flow.md`](execution-flow.md) §8 與 [`development-process.md`](development-process.md) §DoR / §DoD 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-closure` |
| `purpose` | 收尾關卡：實作開始前的 Definition of Ready 檢核、出貨/合併前的 Definition of Done 檢核、以及把超越單一產品的課程回饋成可重用指引 |
| `type` | `execution` |
| `tags` | closure, handoff, extraction-to-intelligence |
| `load_when` | 收尾、DoR/DoD 檢核、回饋可重用課程 |
| `do_not_load_when` | 任務中段（尚在 intake / 實作 / 驗證進行中） |
| `owner_layer` | workflow |
| `layer_justification` | 規定「收尾要過哪些 gate（DoR/DoD）、完成後做什麼順序（回饋）」的 ordering / gate；通過 workflow membership test。`extraction-to-intelligence` 僅為候選標記——「Feed Back Reusable Lessons」**產生** intelligence 候選，但本身是 workflow 閉環步驟；真正的 intelligence 內容（為何某模式長期有效）不在此 slice，須另經 `evidence_refs`≥2 gate 升層 |
| `canonical_source` | 本檔（原 `execution-flow.md` §8 Feed Back Reusable Lessons + `development-process.md` §Minimum Definition Of Ready / §Minimum Definition Of Done） |
| `dependencies` | `sd-intake`（DoR 引用 brief / requirements / contracts 產出）、`sd-validation`（DoD 引用驗證證據） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Phase 4 Scenario A（execution-only：完成宣告的 DoD 檢核）；Scenario E（greenfield/SDD 真實任務的收尾銜接） |

## Minimum Definition Of Ready（最低就緒定義）

在實作開始前，功能應具備：

- 包含範圍和 non-goals 的 Product brief
- Product impact alignment：Impact Map、Customer Journey Map、cross-check decision
- Requirements cognition notes：actor intent、behavior boundary、acceptance criteria、ambiguity disposition
- 關鍵行為的 BDD-lite scenarios，且每個 critical scenario 有 validation target
- Bounded Context 或模組拆分，且已由 behavior boundary / invariant evidence 支撐
- 核心 invariants 的 Domain Model Contract
- 依賴關係、所有權和 runtime boundaries 的 Architecture Contract
- 整合用的 API、event、command 或 public interface contract
- 預期失敗和 recovery behavior 的 Error Handling Contract
- 涵蓋 unit、behavior、contract 和 integration 層級的 Test plan
- 區分既有-regression 覆蓋率與變更/new-code 驗證的 Test strategy
- 當 latency、throughput、資源使用率、concurrency、啟動時間、背景任務、資料庫存取或外部呼叫量可能改變時的效能預算和測試類型
- 沒有影響實作行為或 contracts 的未解決 blocker questions

對於已實作的專案，「ready」表示缺失文件審查已完成，且 BDD 涵蓋已實作的關鍵行為，即使原始 product intent 仍部分未知。

## Minimum Definition Of Done（最低完成定義）

在出貨或合併前：

- Domain invariants 已測試
- Contract tests 對 provider 和 consumer 都通過
- Mocks/fixtures 符合最新 contract
- Integration test 至少涵蓋關鍵 happy path 和一個重要 failure path
- 效能敏感變更已記錄 load、stress、spike、soak 或 smoke-size 效能證據，對照 agreed budget
- 殘留的 unknowns 或延後的行為已在專案 repository 中記錄

### Operational Transaction Close-Out

When a task starts an operational transaction, do not close the work from a start signal alone. This guidance applies to deploys, migrations, backfills, cache rebuilds, data imports, batch jobs, or project-defined operations that can become partial or interrupted.

Close-out evidence should record:

```yaml
operational_transaction_closeout:
  operation: deploy | migration | backfill | cache_rebuild | data_import | batch_job
  transaction_state:
    started: true
    partial: true | false
    interrupted: true | false
    resumed: true | false
    completed: true | false | unknown
    verified: true | false
  final_state_readback: <runtime state, data count, version, health, or business effect>
  closure_decision: complete | blocked | rolled_back | deferred_with_owner
  evidence:
    - start_record
    - completion_record
    - final_state_verification
```

If `completed` or `verified` is unknown, narrow the completion claim or mark the task blocked. If the same transaction-state shape proves reusable across deploy, migration, backfill, cache rebuild, import, and batch job scenarios, promote the reasoning to shared execution reasoning rather than expanding this workflow checklist.

## Feed Back Reusable Lessons（回饋可重複使用的課程）

**Incident-driven path**: after Ship, run [`change-retrospective.md`](change-retrospective.md) first — promotion 僅 `keep local` | `promote project` | `candidate canonical`；禁止 direct canonical promote。

如果一個課程超越了一個產品：

1. 在匹配的 `feedback_history/<category>/` 或跨領域的 `feedback_history/common/` 下添加一個檔案。
2. 連結共享規則而不是複製它們。
3. 將已驗證的指引提升到結構化資料夾、檢查清單或此工作流程中。

如果課程來自 APK 分析，將分析方法保留在 `analysis/apk/` 或 `workflow/apk-analysis/` 中，將開發行動保留在此處。
