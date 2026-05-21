# Cognitive State & Evidence Governance

> **狀態**: completed
> **建立日期**: 2026-05-20
> **目的**: 補足 Runtime Recovery & Escalation System 尚未覆蓋的「假設、證據、信心、來源新鮮度」治理層，讓 agent 在外部 failure 發生前就能偵測 belief drift、降速或停止自動執行。

---

## 1. Problem Statement

Runtime Recovery & Escalation System 已能處理 execution 中段的 repeated failure、user contradiction、evidence conflict、source-of-truth miss、execution graph rebuild 與 recovery validation。

剩下的缺口不是 recovery procedure 不夠，而是：

```text
agent 還沒有失敗
  → 但已形成錯誤 framing / stale assumption / over-anchoring
  → 把臨時 assumption 當成 fact
  → 用低權重 evidence 宣告高權重 checkpoint 成功
  → 後續 execution graph 被 pseudo-truth 污染
```

現有 recovery 偏向「現實已打臉後」進入修復。本 plan 目標是把「belief instability」與「evidence quality mismatch」提前變成可治理 signal。

---

## 2. Scope

### In Scope

- 定義 assumption lifecycle 與 assumption ledger schema。
- 定義 evidence qualification，防止低品質或低適用性 evidence 覆蓋高品質 contradiction。
- 定義 execution confidence decay guard，讓重複 patch、stale assumption、user contradiction 會降低 autonomy。
- 定義 confidence integrity 與 claim scope governance，避免 high confidence / local evidence 被壓縮成 global conclusion。
- 定義 execution intent stability，避免 agent 的 action 漂離原始 goal 與 validation target。
- 定義 evidence lineage 與 contradiction propagation，讓底層 belief 被推翻時可降權相依 conclusion。
- 定義 belief ownership / override authority，避免多種 evidence source 同時改寫同一 belief。
- 定義 autonomy modes，把 confidence decay 轉成可執行的 autonomy downgrade。
- 定義 recovery budget，避免 recovery / rediscovery 無限循環。
- 定義 recovery exit contract，避免 recovery / rediscovery / recap 形成 cognitive recursion。
- 定義 recovery re-entry safety，避免同一 contradiction class 立刻觸發相同 recovery loop。
- 定義 human alignment checkpoint，處理 ambiguity、autonomy threshold、multiple source-of-truth conflict、missing canonical source。
- 定義 cognitive contamination boundary 與 belief garbage collection，避免 stale execution frame 污染新任務。
- 定義 governance minimality 與 cognitive cost governance，避免治理成本大於任務複雜度或形成 governance recursion。
- 定義 tier boundary、compression pass、minimal runtime principle 與 meta-stop rule，作為本 plan 的抽象化封頂機制。
- 定義 cognitive runtime state model 作為後續收斂方向，但不立即持久化為 runtime state。
- 定義 temporal confidence decay，用於 long session、multi-recovery chain 與 stale belief。
- 定義 source freshness governance，避免 stale source-of-truth 被當成 P0 evidence。
- 規劃 routing、metadata、runtime guard、validation scenario 的 linked updates。

### Out of Scope

- 不取代 `enforcement/escalation-policy.md`；該檔仍負責 real-time mismatch escalation。
- 不重寫已完成的 Runtime Recovery & Escalation System。
- 不直接新增已不存在的 `runtime/recovery/*.yaml` standalone source。
- 不把 project-specific incident、UI route、API host、token、device evidence 寫進 reusable docs。
- 不讓 assumption ledger 變成冗長日誌；只記錄會影響 execution decision 的假設。

---

## 3. Layer Responsibility

| Concern | 建議位置 | 邊界 |
| --- | --- | --- |
| Belief drift、over-anchoring、pseudo-truth 的 why | `intelligence/engineering/agent-architecture/` | 不放可執行 MUST rule。 |
| Assumption / evidence / confidence 的治理 gate | `governance/ai-runtime-governance/cognitive-state-governance.md` | 不直接寫 tool-specific 操作。 |
| Intent stability、claim scope、contradiction propagation、belief ownership、minimality、cognitive cost、tier boundary 的治理 gate | `governance/ai-runtime-governance/cognitive-state-governance.md` | 不取代 workflow success criteria；不讓 governance 遞迴膨脹。 |
| Evidence qualification、low-quality override 禁止規則 | `enforcement/evidence-hierarchy.md` | 不放 domain-specific reload set。 |
| Confidence decay、autonomy mode、recovery budget、re-entry safety、human alignment trigger | `runtime/guards/` 或 `runtime/compiler/embedded_data.rb` | 若要 machine enforcement，必須確認 compiler 會讀到 source。 |
| Domain-specific evidence priority / freshness policy | `metadata/recovery/` 或新增 `metadata/evidence/` | metadata-only 時不得假裝已編入 `runtime.db`。 |
| Failure-derived scenario | `validation/scenarios/failure-derived/` | 測 route / source selection / forbidden behavior，不測私有 project evidence。 |

---

## 4. Target State Model

目前 recovery 模型：

```text
DISCOVERY
→ EXECUTION
→ MISMATCH_DETECTION
→ ESCALATION
→ RECOVERY
→ REDISCOVERY
→ REPLAN
→ EXECUTION
```

本 plan 補上的前置治理：

```text
DISCOVERY
→ ASSUMPTION_CAPTURE
→ EVIDENCE_WEIGHTING
→ INTENT_STABILITY_CHECK
→ CONFIDENCE_MONITORING
→ AUTONOMY_MODE_SELECTION
→ GOVERNANCE_MINIMALITY_CHECK
→ EXECUTION
→ MISMATCH_DETECTION / CONFIDENCE_DECAY
→ ESCALATION / RECOVERY_EXIT / HUMAN_ALIGNMENT
```

目標是：不是等到 failure 才 recovery，而是在 assumption 不穩、evidence 不足、confidence decay 時先降低 autonomy。

---

## 5. Proposed Components

### 5.1 Assumption Lifecycle

建議狀態：

```text
assumption
→ tentative_belief
→ validated_belief
→ contradicted_belief
→ suspended_belief
→ deprecated_belief
```

Gate：

- 未驗證 assumption 不得寫成 conclusion。
- 低信心 assumption 不得驅動高風險 action。
- assumption 被 contradiction 推翻後，必須 suspend，不得繼續作為 execution graph 的依據。

### 5.2 Assumption Ledger

只記錄會影響 execution decision 的假設。

| Field | Meaning |
| --- | --- |
| `assumption` | 目前假設。 |
| `source` | 來源：user、doc、search、tool output、inference、memory。 |
| `confidence` | low / medium / high。 |
| `validated_by` | 哪個 evidence 支持。 |
| `contradicted_by` | 哪個 evidence 推翻。 |
| `expires_when` | 何時必須重新驗證。 |
| `execution_dependency` | 哪些操作依賴此 assumption。 |

### 5.3 Evidence Lineage and Contradiction Propagation

重要 belief 需要知道它從哪一串 evidence / inference 長出來，否則 contradiction 發生時無法知道哪些 conclusion 要一起降權。

Lineage fields：

| Field | Meaning |
| --- | --- |
| `belief` | 被使用於 execution graph 的 belief / conclusion。 |
| `derived_from` | 原始 evidence 或上游 belief。 |
| `transformed_by` | 哪個 inference、summary、classification 或 routing decision 產生此 belief。 |
| `validated_by` | 哪個 external evidence 支持此 belief。 |
| `superseded_by` | 哪個新 evidence / belief 取代它。 |
| `dependents` | 哪些 checkpoints、claims 或 execution graph nodes 依賴此 belief。 |

Contradiction cascade：

```text
contradicted evidence
→ dependent beliefs degraded
→ dependent checkpoints invalidated
→ dependent execution graph marked unstable
→ autonomy mode reevaluated
```

Rule：

- 底層 assumption 被推翻時，相依 checkpoint / claim / capture evidence 必須降權或重新驗證。
- Inference 不能失去 lineage 後偽裝成 canonical source。
- Contradiction propagation 是 cognition governance，不應只靠 post-mortem failure learning。

### 5.4 Belief Ownership and Authority Override

Lineage 說明 belief 從哪裡來；ownership 決定誰有權更新或推翻它。沒有 ownership，contradiction propagation 會被 stale workflow、automation log、user statement、runtime evidence 混亂改寫。

初版 ownership：

| Belief type | Authority owner |
| --- | --- |
| UI route / screen state | live runtime observation（screenshot、hierarchy、foreground package） |
| API contract | owner contract / canonical API spec |
| Repo architecture | current repo evidence（actual files、README、routing registry） |
| Execution success | validation evidence / artifact gates |
| Workflow intent | user goal / goal ledger / workflow primary source |
| Runtime state | `runtime/runtime.db` compiled state plus source timestamp |

Override rules：

- Lower-owner evidence may trigger review, but cannot directly overwrite higher-owner belief.
- User statement can override workflow intent, but live runtime evidence still owns UI / API reality.
- Automation log can support a runtime event claim, but cannot own UI route or feature success.
- When owners conflict, enter evidence-conflict handling rather than silently merging beliefs.

### 5.5 Evidence Qualification

Evidence 不只需要 priority。`authority` 只能描述誰比較權威，不能單獨判斷 truth quality。初版 evidence model 至少包含五軸：

| Axis | Meaning |
| --- | --- |
| `authority` | evidence source 的權威性，例如 canonical contract、live observation、tool log、agent inference。 |
| `freshness` | evidence 多新，是否接近當前 execution 時間。 |
| `validity` | evidence 是否仍適用；新文件可能已被 hotfix 推翻，舊 contract 也可能仍有效。 |
| `scope` | evidence 適用範圍，例如 current device、app version、repo、workflow、domain。 |
| `observability` | evidence 是否是 direct observation，或只是 derived / indirect signal。 |

Authority 初版可沿用 P0-P4，但必須與其他軸一起解讀：

| Evidence | Authority | Typical qualification |
| --- | --- | --- |
| Live UI screenshot / hierarchy | P0 | freshness=live, observability=direct, scope=current device |
| Foreground package / active route | P0 | freshness=live, observability=direct, scope=current device |
| Actual API response | P0 | freshness=live/current, observability=direct, scope=request context |
| Owner contract / canonical spec | P0 | authority=high, freshness/validity/scope must be checked |
| Runtime DB / compiled state | P1 | valid only if source and compiler timestamp align |
| Test output / compiler output | P1 | scope limited to tested command / fixture / environment |
| Hook success / instrumentation log | P2 | proves hook/log event, not UI or feature success |
| Search / grep result | P2 | proves textual presence, not runtime behavior |
| Transcript memory | P3 | must be revalidated before driving execution |
| Agent inference / assumption | P4 | cannot override external evidence |

核心 rule：

> Low-quality or low-scope evidence cannot justify a higher-scope claim or override a higher-quality contradiction.

例如：hook success 只能證明 hook 有觸發，不能證明 UI route 正確；fresh UI screenshot 與 stale UI map 衝突時，不是單純 P0 vs P0，而要比較 validity、scope 與 observability。

### 5.6 Confidence Decay, Integrity, and Temporal Decay

Trigger：

- repeated retry
- repeated patch
- evidence conflict
- stale assumption
- user contradiction
- source-of-truth miss
- validation without external evidence

Effect：

- reduce execution autonomy
- require source reload
- require higher-priority evidence
- suspend automation
- enter recovery L3+
- request human alignment when ambiguity or risk crosses threshold

Confidence governance 不能只看 confidence 高低，也要看 confidence 是否與 evidence quality 匹配。

| Integrity | Meaning |
| --- | --- |
| `aligned` | confidence matches evidence qualification and claim scope. |
| `inflated` | confidence is higher than evidence quality supports. |
| `degraded` | confidence should drop because evidence conflicts or assumptions are stale. |
| `unsupported` | claim has no external validation. |

Temporal confidence decay applies only when session length, recovery depth, or stale belief risk justifies it. It should not force every short task through time-based scoring.

Temporal decay triggers：

- long-running session without recap。
- multi-recovery chain。
- stale repo / UI / workflow memory reused after task boundary。
- evidence validity window expired。
- context compaction that preserved conclusions but not supporting evidence。

Rule：

- Confidence should decay over time unless refreshed by valid evidence.
- Temporal decay lowers autonomy mode or requires revalidation; it should not by itself overwrite source-of-truth.
- Temporal decay is disabled for trivial tasks unless evidence instability appears.

### 5.7 Claim Scope Governance

很多 drift 不是 evidence 錯，而是 conclusion scope 超過 evidence coverage。

| Field | Meaning |
| --- | --- |
| `claim` | Agent 想宣告的結論。 |
| `claim_scope` | 結論適用範圍，例如 local hook、single screen、single API、whole feature、whole workflow。 |
| `supporting_evidence` | 支持該 claim 的 evidence qualification。 |
| `uncovered_areas` | evidence 沒覆蓋但 claim 暗示已覆蓋的範圍。 |
| `confidence_integrity` | aligned / inflated / degraded / unsupported。 |

Rule：

- Local evidence cannot justify global execution claims.
- 單一 hook success 不得宣告 feature success。
- 單一 API pass 不得宣告完整 workflow 正常。
- 單一 grep result 不得宣告 implementation complete。

### 5.8 Execution Intent Stability

Intent drift 發生時，assumption 可能仍正確、evidence 也可能有效，但 agent 的 current action 已經不再服務原始 goal。

Intent chain：

```text
goal
→ current execution target
→ active subtask
→ current action
→ validation target
```

Intent drift signals：

- 子任務開始自我擴張。
- patch 開始脫離原 validation target。
- 長時間沒有回到 success criteria。
- current action 與 original execution graph disconnected。
- 新增大量 side quest。
- recovery 本身變成主任務。

Required action：

- recap original goal。
- compare current action vs execution graph。
- prune unrelated subtask。
- require replan if drift is too large。

### 5.9 Autonomy Modes

Confidence decay 的結果應該是 autonomy downgrade，而不只是「信心變低」。

| Mode | Allowed actions | Forbidden actions |
| --- | --- | --- |
| `FULL_AUTONOMY` | patch、retry、modify execution graph、run validation | none beyond normal gates |
| `LIMITED_AUTONOMY` | small scoped edits、source reload、targeted validation | broad refactor、history rewrite、unbounded automation |
| `VALIDATION_REQUIRED` | collect evidence、run validation、compare source-of-truth | patch before validation、claim success |
| `HUMAN_ALIGNMENT_REQUIRED` | summarize options、ask blocker question | autonomous execution beyond agreed next step |
| `READ_ONLY_MODE` | discovery、read docs、inspect evidence | write files、run production action、commit/push |

Mode transition signals：

- confidence integrity becomes `inflated` or `unsupported` → `VALIDATION_REQUIRED`。
- ownership ambiguity or missing canonical authority → `HUMAN_ALIGNMENT_REQUIRED`。
- high cognitive contamination or session-global frame instability → `READ_ONLY_MODE` until rediscovery。
- evidence and validation reestablish alignment → may return to `LIMITED_AUTONOMY` or `FULL_AUTONOMY`。

### 5.10 Cognitive Runtime State Model

Current plan keeps belief execution-scoped and does not persist cognitive state. Still, the scattered mechanisms should have a conceptual consolidation target so the system does not grow into independent, overlapping gates.

Candidate cognitive states：

| State | Meaning |
| --- | --- |
| `STABLE` | Evidence, intent, and autonomy are aligned. |
| `UNCERTAIN` | Assumption or evidence quality is insufficient; validation required before broad action. |
| `DEGRADED` | Confidence integrity or temporal decay requires autonomy downgrade. |
| `CONTAMINATED` | Prior frame / route / checklist may be polluting current execution. |
| `MISALIGNED` | Current action no longer serves original goal / validation target. |
| `RECOVERY` | Active recovery is rebuilding source-of-truth and execution graph. |
| `VALIDATION_REQUIRED` | Action may continue only through evidence gathering / validation. |
| `ALIGNMENT_REQUIRED` | Human alignment is required before further autonomous execution. |
| `READ_ONLY` | Only discovery/read/inspection actions are allowed. |

Boundary：

- This state model is conceptual in this plan.
- It must not become persistent `runtime-state.db` cognition until invalidation propagation, cognitive GC, stale belief reconciliation, and lineage rules are validated.
- If promoted later, it should unify autonomy modes, recovery exit, contamination, and confidence decay instead of adding another independent layer.

### 5.11 Recovery Budget

Recovery 本身也要避免 loop。

| Level | Budget |
| --- | --- |
| L1 local retry | 最多 2-3 次 |
| L2 reload local workflow | 最多 2 次 |
| L3 reload source-of-truth | 最多 2 次 |
| L4 rebuild execution graph | 最多 1 次 |
| L5 assumption collapse | 必須停止 autonomous recovery，評估 human alignment threshold |

Gate：

- 同一 failure class 超過 budget 後，不得繼續 autonomous recovery。
- 必須輸出 blocker、目前 evidence、可選路徑，交給使用者對齊。

### 5.12 Recovery Exit Contract and Re-entry Safety

Recovery 不能只定義 entry；也必須定義何時可以退出，否則容易形成 recovery → rediscovery → recap 的 cognitive recursion。

Recovery exit criteria：

- old assumption invalidated or downgraded。
- new source-of-truth loaded or explicitly marked `source_missing` / `not_applicable`。
- execution graph rebuilt。
- contradiction propagated to dependent beliefs / checkpoints。
- validation evidence reacquired。
- autonomy mode reevaluated。

若 exit criteria 未滿足：

- 不得宣告 recovery complete。
- 不得恢復 `FULL_AUTONOMY`。
- 必須輸出 remaining blocker 與 next minimal safe action。

Recovery re-entry safety：

- Same contradiction class cannot immediately re-trigger the identical recovery path without new evidence.
- Identical recovery path after exit must show strategy change or new source-of-truth.
- If re-entry occurs twice for the same class, downgrade autonomy and require human alignment or blocker output.
- Re-entry cooldown does not suppress new high-risk evidence; it only blocks repetitive recovery loops.

### 5.13 Human Alignment Checkpoint

Human alignment 不是 escalation level，而是 ambiguity / autonomy boundary。以下情況必須 request human clarification：

- multiple source-of-truth conflict
- missing canonical source
- policy ambiguity
- ownership ambiguity
- execution graph instability
- conflicting workflow ownership
- recovery budget exhausted
- execution risk exceeds autonomy threshold

同時需要 alignment threshold，避免 safe but unusable runtime：

| Factor | Effect |
| --- | --- |
| execution risk | increase alignment need |
| ambiguity severity | increase alignment need |
| recovery budget exhaustion | increase alignment need |
| evidence instability | increase alignment need |
| interruption cost | decrease alignment need unless risk is high |

### 5.14 Cognitive Contamination Policy

Cognitive contamination 是 stale execution frame、route belief、repo architecture 或 checklist 被錯套到新 task。

Triggers：

- stale execution memory reused
- prior route reused without validation
- prior repo architecture reused across repo boundary
- cached assumptions crossing domain boundary
- old checklist reused after contradiction

Required action：

- invalidate prior execution frame
- require fresh discovery
- mark reused assumptions as tentative
- require new external evidence

Contamination 需要 propagation boundary，避免 everything contaminates everything：

| Scope | Meaning |
| --- | --- |
| `workflow-local` | 只影響單一 workflow execution frame。 |
| `domain-local` | 影響同一 domain 的 route / checklist / evidence policy。 |
| `session-global` | 影響整個 session，必須 prune / recap / human alignment。 |

### 5.15 Belief Garbage Collection

即使 belief 不 persistent，execution frame 仍可能累積 invalidated assumptions、stale evidence、superseded route belief。

Lifecycle：

```text
contradicted_belief
→ suspended_belief
→ deprecated_belief
→ garbage_collectable
```

Gate：

- Recovery / replan 後，必須移除或降權被取代的 assumption。
- Deprecated belief 不得再作為 execution dependency。
- Long recovery frame 需要 recap / prune，避免 cognitive bloat。

### 5.16 Governance Minimality and Cognitive Cost Governance

這一層的最大風險是 governance recursion 與 cognitive runtime inflation：治理鏈本身開始 drift，或治理成本高於任務本身，導致 agent 安全但不可用。

Governance minimality principle：

> Governance must remain proportional to execution risk, ambiguity, autonomy level, evidence instability, and blast radius. Do not escalate cognitive governance depth beyond what is required to restore safe execution alignment.

Approximate cost classes：

| Mechanism | Cost |
| --- | --- |
| confidence decay signal | low |
| evidence qualification | medium |
| contradiction propagation | medium |
| human alignment | medium |
| rediscovery | high |
| recovery | high |
| full cognitive GC / session-global contamination reset | high |

Gate：

- 先選最小足夠治理深度，再決定是否升級到 recovery / rediscovery / human alignment。
- 若治理成本超過 task complexity，優先降級為 lighter guard。
- 小型任務不應啟動完整 recovery / contradiction cascade，除非有 high-risk evidence conflict。
- Runtime 應偏好最小足夠治理：signal → validate → downgrade autonomy；不要每次都跑完整 governance chain。
- Cost governance 不可覆蓋 P0 safety / source-of-truth / security gates。
- Governance chain 若本身開始產生 intent drift，必須 recap original goal and prune governance subtask。

### 5.17 Tier Boundary and Meta-Stop Rule

本 plan 需要層級封頂，避免每個 governance failure 再生成新的 governance layer。

Tier boundary：

| Tier | Purpose | Blocking rule |
| --- | --- | --- |
| Tier 0 | Safety / source-of-truth | Always blocking when violated. |
| Tier 1 | Evidence correctness | Blocking when claim or action depends on invalid evidence. |
| Tier 2 | Recovery stabilization | Blocking while recovery exit criteria are unmet. |
| Tier 3 | Cognitive optimization | Must not block Tier 0-2 execution unless promoted by concrete failure. |
| Tier 4 | Meta-governance | Documentation / governance by default; runtime enforcement requires proven recurring failure. |

Minimal runtime principle：

> Runtime only enforces the minimum viable cognitive safety signals required to prevent unsafe or invalid autonomous execution.

Meta-stop rule：

> If a governance mechanism primarily exists to govern other governance mechanisms, it should default to documentation/governance layer unless a concrete runtime failure requires enforcement.

Compression pass requirement：

- Normalize overlapping signals before adding new guards.
- Prefer one `execution_reliability_degradation` signal family over many near-duplicate triggers.
- Separate runtime primitives from governance concepts before promotion.
- Do not allow Tier 3+ optimization to delay Tier 0-2 safe execution.
- Next implementation step should be runtime reduction, not new abstraction: compress the current signal set into 3-5 runtime primitives before adding lifecycle, state, governance layer, or metadata category.

Signal family normalization candidate：

```text
execution_reliability_degradation
  ├── confidence decay
  ├── stale belief
  ├── cognitive contamination
  ├── contradiction accumulation
  └── intent instability
```

Runtime primitive boundary：

| Concept | Runtime necessity |
| --- | --- |
| Evidence qualification | yes |
| Confidence integrity | yes |
| Claim scope | runtime-lite; enforce only when claim affects validation or success declaration |
| Contradiction propagation | partial; only for dependent execution claims/checkpoints |
| Cognitive contamination | mostly governance conceptual until validated by scenarios |
| Cognitive cost governance | governance layer by default; runtime only as minimality guard |
| Meta-stop rule | governance only unless recurring runtime recursion failure is validated |
| Tier boundary | governance/compiler gate; not a per-action runtime guard |
| Cognitive runtime state machine | conceptual consolidation target, not persistent state |

### 5.18 Source Freshness and Validity Governance

`source-of-truth` 也可能 stale，因此需要 freshness policy。

| Source | Freshness requirement |
| --- | --- |
| UI map | 必須對應目前 app / release version。 |
| API contract | 必須有 owner 或 last validation。 |
| Workflow | 必須標示 compatibility scope 或 last validation。 |
| Architecture doc | 必須標示適用版本或 migration 狀態。 |
| Runtime DB | 必須可回溯 source 與 compiler timestamp。 |

Gate：

- stale source 不能直接作為 P0 evidence。
- canonical source 與 live evidence 衝突時，進 evidence-conflict recovery，而不是盲目相信文件。

---

## 6. Architecture Compatibility Preflight

| Field | Content |
| --- | --- |
| Trigger | 使用者要求把 Cognitive State / Evidence Governance 補入 plan，並檢查是否與現有框架衝突。 |
| Checked sources | `plans/README.md`、`plans/archived/2026-05-20-1039-runtime-recovery-escalation-system.md`、`governance/ai-runtime-governance/README.md`、`governance/ai-runtime-governance/recovery-retry-governance.md`、`runtime/README.md`、`metadata/recovery/README.md`、`runtime/guards/circuit-breaker.yaml`。 |
| Conflicts | 不應新增 standalone `runtime/recovery/*.yaml`；runtime recovery source 目前在 `runtime/compiler/embedded_data.rb` 並編入 `runtime/runtime.db`。`metadata/recovery/` 目前是 metadata-only，若新增 domain evidence policy，不可宣稱 runtime enforced。`circuit-breaker.yaml` 已有 `hallucination_risk` 與 `mismatch_escalation`，confidence decay 應作為新 guard 或擴充 guard，避免重複語意。 |
| Decision | proceed as new active plan；先建立治理與 enforcement 邊界，再決定哪些 component promotion 到 runtime guard / metadata / validation。 |
| Validation | Plan readback、diff review、linked update check；若後續新增 runtime guard 或 embedded source，必須執行 runtime compiler、runtime DB validation 與 knowledge runtime refresh。 |

### Preflight Findings

1. **與 Recovery Retry Governance 的邊界**：`recovery-retry-governance.md` 處理 retry / escalation / source reload / recovery validation；本 plan 處理 failure 發生前的 assumption / evidence / confidence。兩者相鄰但不重疊。
2. **與 Escalation Policy 的邊界**：`enforcement/escalation-policy.md` 是 real-time control；本 plan 的 evidence qualification、belief ownership、confidence integrity 與 cognitive contamination 可成為 escalation 的上游 signal。
3. **與 Metadata Recovery 的邊界**：`metadata/recovery/` 已保存 domain reload set；若要放 domain evidence qualification，應先決定是擴充 recovery metadata，還是新增 `metadata/evidence/`，避免 metadata 目錄混雜。
4. **與 Runtime Guard 的邊界**：`runtime/guards/circuit-breaker.yaml` 已有 `hallucination_risk` 與 `mismatch_escalation`。`confidence_decay` 若新增，應避免只是把 `mismatch_escalation` 換名字；它的獨立價值是「尚未 failure 但 autonomy 應下降」。`autonomy_modes` 可作為 guard action target，而不是新的哲學層。Cognitive runtime state model 先保持 conceptual，不直接新增 persistent runtime state。
5. **與 Workflow / Goal Ledger 的邊界**：intent stability 不取代 goal ledger；它只檢查 current action 是否仍服務原 goal / validation target。
6. **與 Runtime DB 的邊界**：若只新增 governance/enforcement docs，不需改 `runtime.db`。若新增 guard、autonomy modes、recovery exit criteria 或 embedded recovery budget，必須確認 compiler source 與 `runtime/runtime.db` 同步。

---

## 7. Open Questions for Discussion

A. **Belief / assumption 是否應保持 ephemeral，而非 persistent runtime state？**
   - Current recommendation: adopt lightweight ephemeral cognitive governance first.
   - Belief exists in execution frame / recovery frame / validation frame.
   - Do not store active beliefs in `runtime-state.db` until belief invalidation propagation、cognitive GC、stale belief reconciliation 與 cross-session lineage 都成熟。

B. **Evidence model 是否應從 priority hierarchy 升級成 qualification model？**
   - Current recommendation: yes.
   - Authority P0-P4 is only one axis.
   - Evidence should include authority、freshness、validity、scope、observability。

C. **Confidence governance 是否需要 integrity model？**
   - Current recommendation: yes.
   - Detect high confidence + low evidence quality as `inflated` or `unsupported` confidence.
   - Confidence decay should not be only failure counter.

D. **Human alignment trigger 是否應與 escalation level 解耦？**
   - Current recommendation: yes.
   - Human alignment is an ambiguity / autonomy boundary, not L5-only behavior.

E. **是否需要 cognitive contamination governance？**
   - Current recommendation: yes.
   - Prevent stale route、repo architecture、workflow checklist or prior execution frame from crossing task / domain boundaries without validation.

F. **Evidence qualification 是否需要 validity axis？**
   - Current recommendation: yes.
   - Freshness means how new evidence is; validity means whether it still applies.
   - A fresh document can be invalid after hotfix; an old contract can still be valid.

G. **是否需要 claim scope governance？**
   - Current recommendation: yes.
   - Prevent local evidence from being compressed into global conclusion.
   - Evidence coverage must be at least as broad as claim scope.

H. **Cognitive contamination 是否需要 propagation boundary？**
   - Current recommendation: yes.
   - Use workflow-local / domain-local / session-global to avoid marking everything contaminated by everything.

I. **Human alignment 是否需要 alignment cost / threshold governance？**
   - Current recommendation: yes.
   - Avoid safe but unusable runtime where agent asks human too often.

J. **Execution-scoped belief 是否需要 garbage collection lifecycle？**
   - Current recommendation: yes.
   - Contradicted or deprecated beliefs should be pruned from execution dependencies after recovery / replan.

K. **是否需要 execution intent stability gate？**
   - Current recommendation: yes.
   - Intent drift is separate from belief drift: evidence and assumptions can be valid while the current action no longer serves the original goal.

L. **Evidence lineage 是否需要支援 contradiction propagation？**
   - Current recommendation: yes.
   - Beliefs should keep `derived_from` / `validated_by` / `dependents` so contradicted assumptions can degrade dependent claims and checkpoints.

M. **Autonomy mode 是否應成為 runtime guard 的 action target？**
   - Current recommendation: yes.
   - Confidence decay should downgrade autonomy mode instead of only lowering an abstract confidence score.

N. **Recovery exit criteria 是否要成為 blocking contract？**
   - Current recommendation: yes.
   - Recovery must not exit until assumptions, source reload, graph rebuild, contradiction propagation, validation evidence, and autonomy mode are reconciled.

O. **是否需要 cognitive cost governance？**
   - Current recommendation: yes.
   - Prevent governance overgrowth where runtime spends more effort governing itself than completing the task.

P. **是否需要 governance minimality principle？**
   - Current recommendation: yes.
   - Governance depth should stay proportional to risk、ambiguity、autonomy、evidence instability and blast radius.

Q. **Belief 是否需要 ownership / override authority？**
   - Current recommendation: yes.
   - Belief owner decides who can update or override UI route、API contract、repo architecture、execution success、workflow intent and runtime state beliefs.

R. **Recovery exit 後是否需要 re-entry safety / cooldown？**
   - Current recommendation: yes.
   - Same contradiction class should not immediately re-trigger the identical recovery path without new evidence or strategy change.

S. **Cognitive runtime state machine 是否應先作 conceptual consolidation？**
   - Current recommendation: yes.
   - Use it to unify scattered governance semantics, but do not persist it as runtime state yet.

T. **Temporal confidence decay 是否只在 long-running / stale-context 情境啟用？**
   - Current recommendation: yes.
   - Apply it to long sessions, multi-recovery chains, and stale belief reuse; avoid burdening trivial tasks.

U. **是否需要 tier boundary 來封頂抽象化？**
   - Current recommendation: yes.
   - Tier 3 cognitive optimization and Tier 4 meta-governance must not block Tier 0-2 execution unless a concrete runtime failure promotes them.

V. **是否需要 compression pass 先 normalize signals？**
   - Current recommendation: yes.
   - Confidence decay、intent instability、contradiction propagation、contamination should be evaluated as variants of execution reliability degradation before adding separate guards.
   - Next step is runtime reduction: compress 12+ conceptual signals into 3-5 runtime primitives before adding new lifecycle/state/layer.

W. **哪些 concept 是 runtime primitive，哪些只是 governance concept？**
   - Current recommendation: separate before implementation.
   - Runtime should enforce only minimum viable cognitive safety signals.

X. **是否需要 meta-stop rule？**
   - Current recommendation: yes.
   - Governance mechanisms that mainly govern other governance mechanisms default to governance docs, not runtime enforcement.

---

## 8. Suggested Implementation Phases

### Phase 0 — Compatibility Review

Status: completed (2026-05-21).

Tasks:

- [x] 確認 assumption ledger 不取代 goal ledger / dependency read ledger。
- [x] 確認 belief / assumption 採 ephemeral cognitive governance，不進 persistent runtime state。
- [x] 確認 evidence qualification 與 recovery metadata 的邊界。
- [x] 確認 belief ownership / override authority 的初版 owner table。
- [x] 確認 confidence decay 是否擴充 `circuit-breaker.yaml` 或新增 guard source。
- [x] 確認 temporal confidence decay 只在 long session / stale context 啟用。
- [x] 確認 confidence integrity 與 claim scope governance 的 runtime / enforcement 邊界。
- [x] 確認 execution intent stability 不取代 goal ledger / workflow success criteria。
- [x] 確認 evidence lineage 與 contradiction propagation 的最小可行 scope。
- [x] 確認 autonomy modes 是否作為 runtime guard action target。
- [x] 確認 recovery exit contract 是否屬於 governance、runtime recovery config 或 validation scenario。
- [x] 確認 recovery re-entry safety / cooldown 的最小規則。
- [x] 確認 cognitive runtime state model 先作 conceptual consolidation，不 persistent。
- [x] 確認 cognitive contamination 的 propagation boundary。
- [x] 確認 human alignment threshold，避免 autonomy collapse。
- [x] 確認 belief GC 是否只在 execution frame 內執行。
- [x] 確認 governance minimality / cognitive cost governance，避免小任務觸發過重治理或 governance recursion。
- [x] 確認 tier boundary，限制 Tier 3+ cognitive optimization / meta-governance 不阻塞 Tier 0-2 execution。
- [x] 確認 compression pass，把重疊 signals normalise 為較少的 runtime primitives。
- [x] 確認 runtime reduction：先把 12+ conceptual signals 壓成 3-5 個 runtime primitives，不再新增 lifecycle / state / governance layer。
- [x] 確認 minimal runtime principle 與 meta-stop rule。
- [x] 確認 recovery budget 是否需要 runtime-state counter。
- [x] 依使用者「可以開始」指示，將 §7 open questions 收斂為 Phase 0 initial decisions；後續 phase 若遇到互斥選項再回問使用者。

Exit criteria:

- [x] Layer boundary 已確認。
- [x] Runtime source-of-truth path 已確認。
- [x] 需要 runtime compiler / validator 的後續 phase 已標出。


#### Phase 0 Decision Ledger (2026-05-21)

| Field | Decision |
| --- | --- |
| Trigger | Start implementation for this active plan before Phase 1 governance translation. |
| Checked sources | `plans/README.md`, `governance/ai-runtime-governance/README.md`, `governance/ai-runtime-governance/recovery-retry-governance.md`, `runtime/README.md`, `metadata/recovery/README.md`, `runtime/guards/circuit-breaker.yaml`, `runtime/pipeline/guard-chain.yaml`, `enforcement/README.md`, `enforcement/dependency-reading.md`, `enforcement/linked-updates.md`, `enforcement/content-layering.md`, `enforcement/rule-weight.md`, `enforcement/conversation-goal-ledger.md`. |
| Conflicts | No blocking conflict. Current architecture supports a governance-first implementation. `runtime/recovery/*.yaml` remains removed; recovery runtime source is `runtime/compiler/embedded_data.rb` and compiled `runtime/runtime.db`. `metadata/recovery/` is metadata-only. Existing `mismatch_escalation` and `hallucination_risk` guards overlap with confidence decay, so Phase 3 must normalize signals before adding runtime guard keys. |
| Decision | Proceed. Phase 1 should create the governance translation first. Phase 2 may promote only crisp MUST / forbidden behavior into enforcement. Phase 3 may touch runtime guard sources only after compression to 3-5 runtime primitives. |
| Validation | Use diff review, linked-update check, `ai-skill runtime refresh`, `go test ./...` for CLI guard safety when runtime scripts are involved, commit / push / readback, and clean `git status --short --branch`. |

#### Phase 0 Architecture Decisions

| Topic | Phase 0 decision | Follow-up phase |
| --- | --- | --- |
| Assumption ledger boundary | Keep assumption ledger execution-scoped and lightweight. It does not replace `.agent-goals/`, the dependency-read ledger, or durable planning docs. | Phase 1 governance wording; Phase 5 scenarios. |
| Belief persistence | Do not persist belief / assumption state in `runtime-state.db` in this plan. Treat cognitive state as conceptual consolidation until invalidation and GC rules are validated. | Phase 1; possible future runtime-state plan. |
| Evidence qualification vs recovery metadata | Generic evidence qualification belongs in governance / enforcement. Domain-specific reload sets stay in `metadata/recovery/` unless a later compiler target is explicitly added. | Phase 2 and Phase 4. |
| Belief ownership | Use the plan's owner table as initial governance. Lower-owner evidence may trigger review but cannot silently overwrite higher-owner belief. | Phase 1 and Phase 2. |
| Confidence decay runtime source | Do not add a new runtime guard yet. First compress confidence decay, intent instability, stale belief, contradiction accumulation, and contamination into `execution_reliability_degradation` primitives. | Phase 3. |
| Temporal confidence decay | Enable only for long sessions, multi-recovery chains, expired evidence windows, stale memory reuse, or context compaction without evidence. Trivial tasks remain exempt. | Phase 1 and Phase 3. |
| Claim scope and confidence integrity | Governance owns the model. Enforcement may add a narrow rule: local evidence cannot justify global success claims. Runtime stays runtime-lite unless validation scenarios prove repeated failure. | Phase 1, Phase 2, Phase 5. |
| Intent stability | Intent stability checks current action against original goal / validation target. It does not replace goal ledger or workflow success criteria. | Phase 1 and Phase 2. |
| Lineage and contradiction propagation | Minimum viable scope is dependent execution claims and checkpoints, not every thought or note. | Phase 1 and Phase 5. |
| Autonomy modes | Autonomy modes are valid runtime guard action targets, but the allowed / forbidden action table should be defined in governance first. | Phase 1 and Phase 3. |
| Recovery exit and re-entry safety | Treat recovery exit as a blocking governance contract; runtime promotion requires compiler-source updates. Re-entry cooldown should be minimal and evidence-based. | Phase 1, Phase 3, Phase 5. |
| Human alignment threshold | Human alignment is ambiguity / autonomy boundary, not only L5 escalation. Threshold must consider risk and interruption cost to avoid autonomy collapse. | Phase 1 and Phase 3. |
| Cognitive contamination boundary | Use workflow-local / domain-local / session-global boundaries. Default action is rediscovery and assumption downgrade, not full session reset. | Phase 1 and Phase 5. |
| Belief GC | Keep GC inside the execution or recovery frame; do not create persistent garbage collection state. | Phase 1. |
| Governance minimality / tier boundary | Tier 3+ optimization and Tier 4 meta-governance do not block Tier 0-2 execution unless concrete failures promote them. | Phase 1 and Phase 5. |
| Recovery budget | Do not add a `runtime-state.db` counter yet. Start with governance thresholds and validation scenarios; evaluate stateful counters only if repeated recovery loops remain unbounded. | Phase 3 and future runtime-state plan. |

### Phase 1 — Governance Translation

Candidate files:

- `governance/ai-runtime-governance/cognitive-state-governance.md`
- `governance/ai-runtime-governance/README.md`
- `intelligence/engineering/agent-architecture/`（若需要先沉澱 source intelligence）

Tasks:

- [x] 定義 assumption lifecycle、assumption ledger、belief ownership、evidence lineage、contradiction propagation、evidence qualification、confidence integrity、temporal decay、claim scope、intent stability、cognitive state model、cognitive contamination、belief GC、governance minimality / cognitive cost、tier boundary、meta-stop rule 的治理 gate。
- [x] 明確寫出與 `recovery-retry-governance.md` 的邊界。
- [x] 更新 governance index。

### Phase 2 — Enforcement Rules

Candidate files:

- `enforcement/evidence-hierarchy.md`
- `enforcement/README.md`
- `enforcement/escalation-policy.md`
- `enforcement/failure-learning-system.md`

Tasks:

- [x] 新增 low-quality / low-scope evidence 不得覆蓋 high-quality contradiction 的 MUST rule。
- [x] 定義 assumption 被 contradiction 推翻後的 forbidden behavior。
- [x] 定義 local evidence 不得支持 global claim 的 forbidden behavior。
- [x] 定義 current action 不得長期脫離 original goal / validation target 的 forbidden behavior。
- [x] 將 evidence qualification、belief ownership conflict、confidence integrity、intent drift、cognitive contamination 接到 escalation trigger。

### Phase 3 — Runtime Guard / Budget

Candidate files:

- `runtime/guards/circuit-breaker.yaml`
- `runtime/pipeline/guard-chain.yaml`
- `runtime/compiler/embedded_data.rb`（若 compiler 仍使用 embedded source）
- `runtime/runtime.db`

Tasks:

- [x] 新增或擴充 confidence decay guard。
- [x] 定義 autonomy modes、recovery budget、recovery exit criteria、re-entry safety、alignment threshold 與 contamination boundary trigger。
- [x] 在新增 guard 前執行 signal normalization / compression pass。
- [x] 僅 promotion 最小 runtime primitives；Tier 3+ governance concept 預設不 runtime 化。
- [x] 若需要 stateful counter，評估是否接 `runtime-state.db`。
- [x] 執行 compiler 與 runtime DB validation。

### Phase 4 — Metadata / Domain Evidence Policy

Candidate files:

- `metadata/recovery/domain-policies.yaml`
- `metadata/recovery/README.md`
- 或新增 `metadata/evidence/`

Tasks:

- [x] 決定 domain-specific evidence policy 的位置。
- [x] 加入 APK analysis / software delivery 的 evidence qualification。
- [x] 標明 metadata-only 或 runtime-enforced boundary。

### Phase 5 — Validation Scenarios

Candidate files:

- `validation/scenarios/failure-derived/cognitive-state-assumption-drift.yaml`
- `validation/scenarios/failure-derived/evidence-hierarchy-override.yaml`
- `validation/scenarios/failure-derived/confidence-decay-repeated-patch.yaml`
- `ai-skill runtime validate`（若需 semantic validator）

Tasks:

- [x] 測試 assumption 被當成 fact 的失效模式。
- [x] 測試 hook/log success 不能覆蓋 UI / contract contradiction。
- [x] 測試 repeated patch 後必須降 autonomy 或進 recovery。
- [x] 測試 local evidence 被壓縮成 global claim 時會被阻擋。
- [x] 測試 stale execution frame 跨 domain reuse 會觸發 cognitive contamination。
- [x] 測試 current action 脫離 original goal 時會觸發 intent stability gate。
- [x] 測試底層 assumption 被推翻時相依 checkpoint 會被 invalidated。
- [x] 測試 recovery exit criteria 未滿足時不得恢復 autonomy。
- [x] 測試相同 contradiction class 在沒有新 evidence 時不得立即重跑相同 recovery。
- [x] 測試 governance minimality 會避免小任務啟動過重 governance chain。
- [x] 測試 Tier 3+ cognitive optimization 不會阻塞 Tier 0-2 safe execution。
- [x] 測試 meta-governance mechanism 不會在無 concrete failure 時 promotion 到 runtime。

### Phase 6 — Plan Completion Closure

Tasks:

- [x] 確認所有 phase 完成或標 blocked。
- [x] 執行適用 validator。
- [x] 檢查 linked updates。
- [x] 更新 `plans/README.md` 狀態。
- [x] 若完成，搬移至 `plans/archived/`。
- [x] Commit / push / readback / clean status。

---

## 9. Completion Definition

本 plan 完成時，系統應能做到：

- 未驗證 assumption 不會沉默升級成 fact。
- Belief / assumption 預設保持 execution-scoped，不會形成 hallucination cache。
- Evidence qualification 同時檢查 authority、freshness、validity、scope、observability。
- 低品質或低適用性 evidence 不會覆蓋高品質 contradiction。
- Local evidence 不會被壓縮成 global conclusion。
- Execution intent stability 會阻止 side quest 取代原 goal。
- Evidence lineage 可支援 contradiction propagation。
- Belief ownership / override authority 會決定誰能更新或推翻特定 belief。
- 重複 patch / retry 會降低 execution autonomy。
- Temporal confidence decay 只在 long-running / stale-context 情境觸發，不拖慢短任務。
- Autonomy modes 能把 confidence decay 轉成具體 allowed / forbidden actions。
- Human alignment 由 ambiguity / autonomy threshold 觸發，不與 L5 綁死。
- Recovery 有 budget，超過後會停止 autonomous recovery。
- Recovery 有 exit contract，未滿足前不得宣告完成或恢復 full autonomy。
- Recovery re-entry safety 會阻止同類 contradiction 無新 evidence 時立即重跑相同 recovery。
- Cognitive runtime state model 作為 conceptual consolidation target，不立即持久化為 runtime state。
- Cognitive contamination 有 workflow-local / domain-local / session-global boundary。
- Execution-scoped belief 有 GC lifecycle，避免 recovery / replan 後 cognitive bloat。
- Governance minimality / cognitive cost governance 會防止小任務啟動過重治理鏈或 governance recursion。
- Tier boundary 會封頂抽象化，Tier 3+ 不會無條件阻塞 Tier 0-2 execution。
- Compression pass 會先合併重疊 signals，再決定是否新增 runtime guard。
- Minimal runtime principle 會限制 runtime 只 enforcement 最小必要 cognitive safety primitives。
- Meta-stop rule 會阻止 governance-of-governance 預設進 runtime。
- stale source-of-truth 不會被無條件當成 P0 evidence。
- 至少十二個 failure-derived scenario 驗證 assumption drift、evidence qualification override、belief ownership conflict、confidence integrity、claim scope overcompression、intent drift、contradiction propagation、recovery exit、re-entry safety、cognitive contamination、governance recursion、Tier 3+ over-blocking。

