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
- 定義 evidence hierarchy，防止低權重 evidence 覆蓋高權重 contradiction。
- 定義 execution confidence decay guard，讓重複 patch、stale assumption、user contradiction 會降低 autonomy。
- 定義 recovery budget，避免 recovery / rediscovery 無限循環。
- 定義 human alignment checkpoint，處理 L5、multiple source-of-truth conflict、missing canonical source。
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
| Evidence priority、low-priority override 禁止規則 | `enforcement/evidence-hierarchy.md` | 不放 domain-specific reload set。 |
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

### 5.3 Evidence Hierarchy

初版 priority：

| Evidence | Priority |
| --- | --- |
| Live UI screenshot / hierarchy | P0 |
| Foreground package / active route | P0 |
| Actual API response | P0 |
| Owner contract / canonical spec | P0 |
| Runtime DB / compiled state | P1 |
| Test output / compiler output | P1 |
| Hook success / instrumentation log | P2 |
| Search / grep result | P2 |
| Transcript memory | P3 |
| Agent inference / assumption | P4 |

核心 rule：

> Lower-priority evidence cannot override higher-priority contradiction.

例如：hook success 只能證明 hook 有觸發，不能證明 UI route 正確。

### 5.4 Confidence Decay

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
- request human alignment at L5

### 5.5 Recovery Budget

Recovery 本身也要避免 loop。

| Level | Budget |
| --- | --- |
| L1 local retry | 最多 2-3 次 |
| L2 reload local workflow | 最多 2 次 |
| L3 reload source-of-truth | 最多 2 次 |
| L4 rebuild execution graph | 最多 1 次 |
| L5 assumption collapse | 必須 human alignment |

Gate：

- 同一 failure class 超過 budget 後，不得繼續 autonomous recovery。
- 必須輸出 blocker、目前 evidence、可選路徑，交給使用者對齊。

### 5.6 Human Alignment Checkpoint

以下情況必須 request human clarification：

- multiple source-of-truth conflict
- missing canonical source
- repeated L5 escalation
- execution graph instability
- conflicting workflow ownership
- recovery budget exhausted

### 5.7 Source Freshness Governance

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
2. **與 Escalation Policy 的邊界**：`enforcement/escalation-policy.md` 是 real-time control；本 plan 的 evidence hierarchy 可成為 escalation 的上游 signal。
3. **與 Metadata Recovery 的邊界**：`metadata/recovery/` 已保存 domain reload set；若要放 domain evidence priority，應先決定是擴充 recovery metadata，還是新增 `metadata/evidence/`，避免 metadata 目錄混雜。
4. **與 Runtime Guard 的邊界**：`runtime/guards/circuit-breaker.yaml` 已有 `hallucination_risk` 與 `mismatch_escalation`。`confidence_decay` 若新增，應避免只是把 `mismatch_escalation` 換名字；它的獨立價值是「尚未 failure 但 autonomy 應下降」。
5. **與 Runtime DB 的邊界**：若只新增 governance/enforcement docs，不需改 `runtime.db`。若新增 guard 或 embedded recovery budget，必須確認 compiler source 與 `runtime/runtime.db` 同步。

---

## 7. Open Questions for Discussion

1. **Assumption ledger 放哪裡？**
   - Option A: 只作為 governance / enforcement output shape，不建 runtime state。
   - Option B: 進 `runtime-state.db`，記錄 session-level active assumptions。
   - Option C: 用 `.agent-goals/` 類似 project-local temp state，但不 commit。

2. **Evidence hierarchy 是全域 rule 還是 domain metadata？**
   - 全域 rule 適合定義 P0-P4。
   - Domain metadata 適合定義 APK / software delivery / repo maintenance 的具體 evidence priority。

3. **Confidence decay 應該放在 `runtime/guards/circuit-breaker.yaml` 還是新 guard？**
   - 放在 circuit breaker：整合成本低，但語意可能變胖。
   - 新增 guard：語意清楚，但需確認 compiler mapping。

4. **Recovery budget 是否屬於 recovery governance 或 runtime state？**
   - Governance 可先定義 budget rule。
   - Runtime enforcement 需要 counter，可能要接 `runtime-state.db`。

5. **Source freshness 是否獨立成 governance？**
   - 若只是 recovery reload 時檢查，可放在 recovery governance。
   - 若要所有 source-of-truth 都有 freshness gate，應新增 `governance/source-freshness.md` 或 `governance/ai-runtime-governance/source-freshness-governance.md`。

---

## 8. Suggested Implementation Phases

### Phase 0 — Compatibility Review

Status: draft.

Tasks:

- [ ] 確認 assumption ledger 不取代 goal ledger / dependency read ledger。
- [ ] 確認 evidence hierarchy 與 recovery metadata 的邊界。
- [ ] 確認 confidence decay 是否擴充 `circuit-breaker.yaml` 或新增 guard source。
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

- [ ] 定義 assumption lifecycle、assumption ledger、evidence hierarchy、confidence decay 的治理 gate。
- [ ] 明確寫出與 `recovery-retry-governance.md` 的邊界。
- [ ] 更新 governance index。

### Phase 2 — Enforcement Rules

Candidate files:

- `enforcement/evidence-hierarchy.md`
- `enforcement/README.md`
- `enforcement/escalation-policy.md`
- `enforcement/failure-learning-system.md`

Tasks:

- [ ] 新增 low-priority evidence 不得覆蓋 high-priority contradiction 的 MUST rule。
- [ ] 定義 assumption 被 contradiction 推翻後的 forbidden behavior。
- [ ] 將 evidence hierarchy 接到 escalation trigger。

### Phase 3 — Runtime Guard / Budget

Candidate files:

- `runtime/guards/circuit-breaker.yaml`
- `runtime/pipeline/guard-chain.yaml`
- `runtime/compiler/embedded_data.rb`（若 compiler 仍使用 embedded source）
- `runtime/runtime.db`

Tasks:

- [ ] 新增或擴充 confidence decay guard。
- [ ] 定義 recovery budget 與 L5 human alignment trigger。
- [ ] 若需要 stateful counter，評估是否接 `runtime-state.db`。
- [ ] 執行 compiler 與 runtime DB validation。

### Phase 4 — Metadata / Domain Evidence Policy

Candidate files:

- `metadata/recovery/domain-policies.yaml`
- `metadata/recovery/README.md`
- 或新增 `metadata/evidence/`

Tasks:

- [ ] 決定 domain-specific evidence policy 的位置。
- [ ] 加入 APK analysis / software delivery 的 evidence priority。
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
- 低權重 evidence 不會覆蓋高權重 contradiction。
- 重複 patch / retry 會降低 execution autonomy。
- recovery 有 budget，超過後會 human alignment。
- stale source-of-truth 不會被無條件當成 P0 evidence。
- 至少三個 failure-derived scenario 驗證 assumption drift、evidence override、confidence decay。

