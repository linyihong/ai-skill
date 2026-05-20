# Cognitive State & Evidence Governance

> **狀態**: draft
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
- 定義 recovery budget，避免 recovery / rediscovery 無限循環。
- 定義 human alignment checkpoint，處理 ambiguity、autonomy threshold、multiple source-of-truth conflict、missing canonical source。
- 定義 cognitive contamination boundary 與 belief garbage collection，避免 stale execution frame 污染新任務。
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
| Evidence qualification、low-quality override 禁止規則 | `enforcement/evidence-hierarchy.md` | 不放 domain-specific reload set。 |
| Confidence decay、recovery budget、human alignment trigger | `runtime/guards/` 或 `runtime/compiler/embedded_data.rb` | 若要 machine enforcement，必須確認 compiler 會讀到 source。 |
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
→ CONFIDENCE_MONITORING
→ EXECUTION
→ MISMATCH_DETECTION / CONFIDENCE_DECAY
→ ESCALATION or HUMAN_ALIGNMENT
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

### 5.3 Evidence Qualification

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

### 5.4 Confidence Decay and Integrity

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

### 5.5 Claim Scope Governance

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

### 5.6 Recovery Budget

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

### 5.7 Human Alignment Checkpoint

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

### 5.8 Cognitive Contamination Policy

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

### 5.9 Belief Garbage Collection

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

### 5.10 Source Freshness and Validity Governance

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
2. **與 Escalation Policy 的邊界**：`enforcement/escalation-policy.md` 是 real-time control；本 plan 的 evidence qualification、confidence integrity 與 cognitive contamination 可成為 escalation 的上游 signal。
3. **與 Metadata Recovery 的邊界**：`metadata/recovery/` 已保存 domain reload set；若要放 domain evidence qualification，應先決定是擴充 recovery metadata，還是新增 `metadata/evidence/`，避免 metadata 目錄混雜。
4. **與 Runtime Guard 的邊界**：`runtime/guards/circuit-breaker.yaml` 已有 `hallucination_risk` 與 `mismatch_escalation`。`confidence_decay` 若新增，應避免只是把 `mismatch_escalation` 換名字；它的獨立價值是「尚未 failure 但 autonomy 應下降」。
5. **與 Runtime DB 的邊界**：若只新增 governance/enforcement docs，不需改 `runtime.db`。若新增 guard 或 embedded recovery budget，必須確認 compiler source 與 `runtime/runtime.db` 同步。

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

---

## 8. Suggested Implementation Phases

### Phase 0 — Compatibility Review

Status: draft.

Tasks:

- [ ] 確認 assumption ledger 不取代 goal ledger / dependency read ledger。
- [ ] 確認 belief / assumption 採 ephemeral cognitive governance，不進 persistent runtime state。
- [ ] 確認 evidence qualification 與 recovery metadata 的邊界。
- [ ] 確認 confidence decay 是否擴充 `circuit-breaker.yaml` 或新增 guard source。
- [ ] 確認 confidence integrity 與 claim scope governance 的 runtime / enforcement 邊界。
- [ ] 確認 cognitive contamination 的 propagation boundary。
- [ ] 確認 human alignment threshold，避免 autonomy collapse。
- [ ] 確認 belief GC 是否只在 execution frame 內執行。
- [ ] 確認 recovery budget 是否需要 runtime-state counter。
- [ ] 與使用者討論 §7 open questions。

Exit criteria:

- [ ] Layer boundary 已確認。
- [ ] Runtime source-of-truth path 已確認。
- [ ] 需要 runtime compiler / validator 的後續 phase 已標出。

### Phase 1 — Governance Translation

Candidate files:

- `governance/ai-runtime-governance/cognitive-state-governance.md`
- `governance/ai-runtime-governance/README.md`
- `intelligence/engineering/agent-architecture/`（若需要先沉澱 source intelligence）

Tasks:

- [ ] 定義 assumption lifecycle、assumption ledger、evidence qualification、confidence integrity、claim scope、cognitive contamination、belief GC 的治理 gate。
- [ ] 明確寫出與 `recovery-retry-governance.md` 的邊界。
- [ ] 更新 governance index。

### Phase 2 — Enforcement Rules

Candidate files:

- `enforcement/evidence-hierarchy.md`
- `enforcement/README.md`
- `enforcement/escalation-policy.md`
- `enforcement/failure-learning-system.md`

Tasks:

- [ ] 新增 low-quality / low-scope evidence 不得覆蓋 high-quality contradiction 的 MUST rule。
- [ ] 定義 assumption 被 contradiction 推翻後的 forbidden behavior。
- [ ] 定義 local evidence 不得支持 global claim 的 forbidden behavior。
- [ ] 將 evidence qualification、confidence integrity、cognitive contamination 接到 escalation trigger。

### Phase 3 — Runtime Guard / Budget

Candidate files:

- `runtime/guards/circuit-breaker.yaml`
- `runtime/pipeline/guard-chain.yaml`
- `runtime/compiler/embedded_data.rb`（若 compiler 仍使用 embedded source）
- `runtime/runtime.db`

Tasks:

- [ ] 新增或擴充 confidence decay guard。
- [ ] 定義 recovery budget、alignment threshold 與 contamination boundary trigger。
- [ ] 若需要 stateful counter，評估是否接 `runtime-state.db`。
- [ ] 執行 compiler 與 runtime DB validation。

### Phase 4 — Metadata / Domain Evidence Policy

Candidate files:

- `metadata/recovery/domain-policies.yaml`
- `metadata/recovery/README.md`
- 或新增 `metadata/evidence/`

Tasks:

- [ ] 決定 domain-specific evidence policy 的位置。
- [ ] 加入 APK analysis / software delivery 的 evidence qualification。
- [ ] 標明 metadata-only 或 runtime-enforced boundary。

### Phase 5 — Validation Scenarios

Candidate files:

- `validation/scenarios/failure-derived/cognitive-state-assumption-drift.yaml`
- `validation/scenarios/failure-derived/evidence-hierarchy-override.yaml`
- `validation/scenarios/failure-derived/confidence-decay-repeated-patch.yaml`
- `scripts/validate-knowledge-runtime.rb`（若需 semantic validator）

Tasks:

- [ ] 測試 assumption 被當成 fact 的失效模式。
- [ ] 測試 hook/log success 不能覆蓋 UI / contract contradiction。
- [ ] 測試 repeated patch 後必須降 autonomy 或進 recovery。
- [ ] 測試 local evidence 被壓縮成 global claim 時會被阻擋。
- [ ] 測試 stale execution frame 跨 domain reuse 會觸發 cognitive contamination。

### Phase 6 — Plan Completion Closure

Tasks:

- [ ] 確認所有 phase 完成或標 blocked。
- [ ] 執行適用 validator。
- [ ] 檢查 linked updates。
- [ ] 更新 `plans/README.md` 狀態。
- [ ] 若完成，搬移至 `plans/archived/`。
- [ ] Commit / push / readback / clean status。

---

## 9. Completion Definition

本 plan 完成時，系統應能做到：

- 未驗證 assumption 不會沉默升級成 fact。
- Belief / assumption 預設保持 execution-scoped，不會形成 hallucination cache。
- Evidence qualification 同時檢查 authority、freshness、validity、scope、observability。
- 低品質或低適用性 evidence 不會覆蓋高品質 contradiction。
- Local evidence 不會被壓縮成 global conclusion。
- 重複 patch / retry 會降低 execution autonomy。
- Human alignment 由 ambiguity / autonomy threshold 觸發，不與 L5 綁死。
- Recovery 有 budget，超過後會停止 autonomous recovery。
- Cognitive contamination 有 workflow-local / domain-local / session-global boundary。
- Execution-scoped belief 有 GC lifecycle，避免 recovery / replan 後 cognitive bloat。
- stale source-of-truth 不會被無條件當成 P0 evidence。
- 至少五個 failure-derived scenario 驗證 assumption drift、evidence qualification override、confidence integrity、claim scope overcompression、cognitive contamination。

