# 失效學習系統

本規則把重複 agent 錯誤轉成 durable prevention。它連接 failure capture、classification、reusable pattern records、linked updates 與 validation，讓同類失效模式較不容易重演。

當使用者指出 agent mistake、close-loop gap、wrong source/mirror update、missed dependency、incomplete validation、forgotten goal、unsafe parallel work，或任何應成為本系統一部分的重複行為時，使用本規則。

## 核心規則

發現 failure 時，不要只修當下檔案。必須跑 failure learning loop：

1. **Capture**：在目前工作脈絡記錄發生什麼、在哪裡被發現、造成什麼 user-facing risk。
2. **Classify**：用下方 taxonomy 分類失效模式。
3. **Contain**：繼續廣泛工作前先控制當前風險。
4. **Promote**：把可重用 lesson 放到正確 durable location。**寫入前必須依 [`sanitization.md`](sanitization.md) 去敏**：不得包含本機絕對路徑（改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 占位符）、使用者名稱、私有工作目錄、clone 位置、secrets、raw tokens、私人 host、個資或 project-specific evidence。
5. **Strengthen**：補強原本可防止它的 rule、workflow、checklist、tool adapter 或 validation gate。
6. **Validate**：確認未來 agent 找得到並能套用這個 prevention。

目標不是歸檔責任，而是把已觀察到的 failure 轉成有明確 trigger 與 validation method 的 reusable guardrail。

### Recovery 與 Escalation 參考

當 failure 或 mismatch 正在發生時，先依 [`escalation-policy.md`](escalation-policy.md) 判斷是否需要停止局部 patch、重讀 source-of-truth、重建 execution graph。Escalation 是 real-time control；本檔是 post-mortem learning。

若 failure 導致 phase transition 被 blocking gate 阻斷，修復流程由 runtime recovery 系統管理。部分 legacy YAML 已 embedded 到 [`runtime/compiler/embedded_data.rb`](../runtime/compiler/embedded_data.rb)，並編譯進 [`runtime/runtime.db`](../runtime/runtime.db)；不要引用不存在的 `runtime/recovery/*.yaml` 作為 standalone source，除非另有 source restoration migration。

本節的 failure learning loop 專注於「將 failure 轉為 reusable prevention」（capture → classify → promote → strengthen → validate），而 escalation / runtime recovery 專注於「failure 發生時的即時停止、重讀與狀態恢復」。兩者互補：

| 面向 | Failure Learning System | Escalation / Runtime Recovery |
|------|------------------------|-------------------------------|
| 目標 | 防止同類 failure 重演 | 停止錯路線並修復目前 failure 造成的損害 |
| 產出 | failure pattern、feedback lesson、validation scenario | recovery output、source reload、new execution graph、runtime recovery state |
| 時機 | failure 被控制後（post-mortem） | failure / mismatch 正在發生時（real-time） |
| 範圍 | cross-session、cross-agent | 目前 session、目前 phase / workflow |

## Failure Taxonomy

| Class | 意義 | 常見 prevention |
| --- | --- | --- |
| `source-mirror-drift` | Agent 更新了 local tool mirror、project `.cursor`、runtime copy 或 generated bundle，而不是 canonical source repository。 | 要求 canonical repo check、source-first edit，再做 tool sync。 |
| `dependency-miss` | Agent 修改或使用 rule/skill 時沒有讀 required linked dependencies。 | 補強 dependency read ledger 與 linked updates。 |
| `goal-ledger-miss` | Multi-step 或可恢復 user goal 沒有正確記錄、更新、拆分、暫停或完成。 | 繼續前更新 `.agent-goals/`，並連到 todos/plans。 |
| `validation-gap` | Agent 未做 diff review、lints、tests、link check、source check、sync、push、readback 或 clean status 就宣稱完成。 | 加入具體 validation gate，並回報實際跑了什麼。 |
| `scope-drift` | Agent 把無關變更、project incident details 或 local absolute paths 混進 reusable docs。 | 套用 reusable guidance boundary 與 sanitization。 |
| `handoff-gap` | Agent 留下不清楚的 next actions、blockers、owner/lock state 或 remaining decisions。 | 更新 goal ledger、Document TODO 或 handoff notes。 |
| `tool-strategy-gap` | Reusable rule 假設單一工具行為，而沒有隔離 tool-specific execution。 | 將工具細節移到 `ai-tools/` 或 skill adapter。 |
| `parallelization-risk` | 多個 agents 可能獨立編輯 shared state、git history、migrations、release steps 或 rules。 | 標記 `non-parallelizable` 或 `single-owner`，遇到衝突 lock 就停止。 |

若 failure 不符合現有 class，先檢查能否用既有 class 清楚描述；只有必要時才新增 class。

## Storage Rules

| 內容 | Durable location |
| --- | --- |
| 當前未完成修復、owner、lock、next action | `<PROJECT_ROOT>/.agent-goals/` |
| 跨 skill 可重用 failure pattern | `enforcement/failure-patterns/` |
| Skill-specific technique 或 failure lesson | `feedback/history/<domain>/` |
| Tool-specific reminder、hook、prompt 或 UI detail | `ai-tools/<tool>.md` 或 tool config |
| Project incident evidence、raw logs、exact private paths、hosts、tokens | Project docs、issue tracker 或 private evidence，不進 reusable docs |

不要把 secrets、real tokens、raw private data 或 local absolute paths 寫進 failure patterns。使用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<tool-mirror>`、`<runtime-copy>` 等 placeholders。

## Failure Pattern Record

當同一 failure 可能跨 projects、agents、tools 或 skills 重演時，建立或更新 reusable pattern。

建議檔案：

```text
enforcement/failure-patterns/<short-slug>.md
```

建議格式：

```markdown
# <Pattern Title>

Status: candidate | validated | promoted | deprecated
Class: <taxonomy class>

## Trigger
When should an agent suspect this failure?

## Failure Mode
What goes wrong in generalized terms?

## Risk
What user-facing, repo, validation, or handoff risk does it create?

## Required Agent Action
What must the agent do next time?

## Prevention Gate
What check would have stopped the mistake?

## 驗證
How to confirm the prevention worked?

## Linked Rules
- <enforcement rule / skill / tool doc links>
```

Pattern records 要短。若 pattern 變長，將 examples 拆成較小 pattern files，並讓 `failure-patterns/README.md` 保持索引功能。

## Promotion Decision

分類 failure 後，選擇能防止重犯的最小 durable promotion target：

| Failure scope | Promotion target |
| --- | --- |
| 只影響單一 active conversation | `.agent-goals/` progress 或 handoff note |
| 單一 reusable document 有局部 open gap | 該文件前段的 Document TODO |
| Cross-document 或 cross-agent workflow failure | `enforcement/failure-patterns/` 加上相關 enforcement rule |
| Skill-specific repeated mistake | 該 skill 的 `feedback_history/`，成熟後再推進 workflow/checklist |
| Tool-specific execution failure | `ai-tools/<tool>.md`、tool config 或 skill tool adapter |
| **架構重構後 enforcement 未同步** | **`enforcement/failure-patterns/shared-rules-architecture-drift.md`** + 執行 `governance/lifecycle/intelligence-extraction-pipeline.md` Step 7a |
| **AI 系統面執行錯誤**（routing 錯誤、heuristic 誤用、forbidden route 被選中） | **`validation/scenarios/failure-derived/`** — 建立 stateless scenario，未來可自動驗證同類錯誤是否重演 |

### Failure → Validation Scenario 條件

Scenario promotion 的治理 gate 見 [`governance/ai-runtime-governance/validation-scenario-governance.md`](../governance/ai-runtime-governance/validation-scenario-governance.md)；本節保留 failure learning loop 需要的最小條件。

當 failure 符合以下所有條件時，應 promotion 為 validation scenario：

1. **可 stateless 重現**：不需要前文提示或 conversation memory，給定相同 signals 應產出相同決策
2. **有明確的 expected/forbidden route**：可以定義「應該走哪條路」和「不該走哪條路」
3. **有 prevention 價值**：未來模型升級、架構變更或 routing 調整後，可能再次發生
4. **不是一次性事件**：同類 signals 組合可能在真實任務中再次出現

建立 scenario 後，在 failure pattern 或 feedback lesson 中標註對應的 scenario ID：

```markdown
## Validation Scenario
- [`validation/scenarios/failure-derived/<id>.yaml`](../../validation/scenarios/failure-derived/<id>.yaml)
```

不要把 project incident 直接推進 reusable docs。必須先泛化 cause、trigger、required action 與 validation。

## Source And Mirror Failures

`source-mirror-drift` 是本 repository 的高優先失效類別。

當使用者要求更新 rules、skills、feedback lessons、templates 或 OS guidance 時，agent 必須：

1. 定位 canonical `<AI_SKILL_REPO>` git root。
2. 在該 repo 確認 `git status --short --branch`。
3. 先編輯 canonical source files。
4. 將工具部署 / mirror 路徑（如 `~/.cursor/`、`~/.claude/`、專案本機設定檔、generated bundles 與 project-local mirrors）視為 deployment/runtime surfaces。具體工具部署路徑見 [`ai-tools/agent/`](../ai-tools/agent/) 中各工具文件。
5. Source repo change 完成後才同步 mirrors。
6. Commit、push、read back，並確認 clean status 後才宣稱完成。

Reference-first tool setup 可減少 duplicate copies，但不取代 source check。若沒有明確 gate，agent 仍可能寫錯地方。

## 驗證

關閉 failure-learning update 前，確認：

- Immediate issue 已被控制，或明確記錄仍 open。
- Failure class 已命名。
- Durable location 正確。
- Prevention gate 已寫在未來 agent 會讀到的位置。
- Linked updates 已檢查。
- 若 canonical repository 有變更，tool sync、commit、push、readback 與 clean status 已完成；reference-first 時 tool sync 可標 not applicable。

### State-based Enforcement（狀態化強制規則）

下列 failure-learning close-loop 檢查已對應到 runtime state machine 的 phase/gate 定義：

```yaml
# State-based enforcement mapping for failure-learning-system.md
# 這些規則已由 runtime state machine 管理，agent 不應再以 prose 方式逐條檢查。
state_based_enforcement:
  version: v1
  status: active
  owner_layer: enforcement/failure-learning-system
  description: >
    將 failure-learning close-loop 檢查對應到 runtime state machine 的 phase/gate。
    Agent 應優先查詢 runtime/runtime.db，必要時讀 runtime/compiler/embedded_data.rb，
    而非依賴本節的 prose 摘要或已移除的 standalone YAML。

  # Immediate issue control → 由 runtime/runtime.db 的 phase.recovery 管理
  - rule: immediate_issue_control
    phase: phase.recovery
    description: "Immediate issue 已被控制，或明確記錄仍 open"
    runtime_ref: runtime/runtime.db
    runtime_section: "phase.recovery.entry_conditions"

  # Failure classification → 由 runtime/runtime.db 的 obligation.recovery.diagnose_root_cause 管理
  - rule: failure_classification
    phase: phase.recovery
    obligation: obligation.recovery.diagnose_root_cause
    description: "Failure class 已命名，root cause 已診斷"
    runtime_ref: runtime/runtime.db
    runtime_section: "obligation.recovery.diagnose_root_cause"

  # Durable location → 由 runtime/runtime.db 的 phase.execution.allowed_actions 管理
  - rule: durable_location
    phase: phase.execution
    description: "Durable location 正確（feedback_history/、enforcement/failure-patterns/ 等）"
    runtime_ref: runtime/runtime.db
    runtime_section: "phase.execution.allowed_actions"

  # Prevention gate → 由 runtime/runtime.db 的 gate.recovery.repair_strategy_defined 管理
  - rule: prevention_gate
    phase: phase.validation
    gate: gate.validation.artifacts_complete
    description: "Prevention gate 已寫在未來 agent 會讀到的位置"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.validation.artifacts_complete"

  # Linked updates → 由 runtime/runtime.db 的 gate.validation.linked_updates_complete 管理
  - rule: linked_updates
    phase: phase.validation
    gate: gate.validation.linked_updates_complete
    description: "Linked updates 已檢查"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.validation.linked_updates_complete"

  # Writeback close-loop → 由 runtime/runtime.db 的 verified → closed 管理
  - rule: writeback_close_loop
    phase: phase.readback
    description: |
      若 canonical repository 有變更：
      tool sync → commit → push → readback → clean status 已完成。
      Reference-first 時 tool sync 可標 not applicable。
    runtime_ref: runtime/runtime.db
    runtime_section: "state.verified → state.closed"
```

← [Back to enforcement index](README.md)
