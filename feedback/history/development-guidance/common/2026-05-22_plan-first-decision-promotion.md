> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-22 — Plan-First Decision Promotion（提案在 plan，憲法在驗證後）

Status: candidate

#### One-line Summary

架構決策的提案、討論、alternatives 評估在 `plans/active/<plan>.md` §Decision Rationale 完成；只有 plan completed 且通過 ADR Promotion Criteria 後才升級為 accepted ADR。憲法層只放已驗證決策，不放 proposed/draft，避免「廢棄憲法」累積與狀態爆炸。

#### Human Explanation

2026-05-21～22 session 內為 Runtime Cognitive Modes 提案直接建立 `constitution/ADR-008-runtime-cognitive-modes.md` (Status: proposed) + 對應 plan (Status: draft) 平行存在。使用者指出此模式有問題：constitution 應該是「已驗證決策集合」，proposed ADR 違反此定位。Plan 已是天然提案容器，平行 ADR 造成雙寫、狀態爆炸、「廢棄憲法」累積風險。

修正方向：plan-first / ADR-after-completion。提案階段全在 plan §Decision Rationale，accepted ADR 在 plan completed + Promotion Criteria 通過後才建立。

#### Trigger

- 想為新架構/流程/跨層改動寫 ADR
- 想標記某個 ADR 為 `proposed` / `draft` / `pending`
- 已有對應 `plans/active/<plan>.md` 但仍想平行建 ADR
- 使用者問「這個要不要寫成 ADR」
- 看到既有 constitution/ 內有非 accepted status 的 ADR

#### Evidence

- Tool: 對話與 commit history
- Sanitized excerpt: ADR-008 (proposed) + plan (draft) 平行存在 → 使用者指出「以後計畫的東西是在 plan 寫好，完成後才把東西寫到 constitution」→ 撤回 ADR-008 + 整併入 plan §Decision Rationale + 確立 §No-Proposed-ADR Rule
- Evidence path: commit history `90e87fa` (建立 ADR-008 proposed) → `db9b515` (撤回 ADR-008 + 確立規則)

#### Generalized Lesson

完整提案 → 驗證 → 憲法的 promotion flow：

```
有新架構提案
  ↓
寫 plan（draft）含 §Decision Rationale 6 子章節：
  - Problem & Why Now
  - Decision
  - Alternatives Considered
  - Why Not an ADR Yet
  - ADR Promotion Criteria
  - Consequences
  ↓
討論、調整、進 in-progress
  ↓
Phase 0 pre-build interrogation → Phase 1-N
  ↓
plan completed
  ↓
評估 5 條 ADR Promotion Criteria：
  - foundational + cross-session + cross-project + expensive-to-reverse + explains-why
  - 全中 → ADR (accepted)；不全中 → 改用更輕 promotion target
```

#### Agent Action

寫架構決策時：

1. **不直接建 proposed ADR**
2. 改寫對應 plan 的 §Decision Rationale section
3. Plan 進行中持續更新 §Decision Rationale（決策可動，未升憲法）
4. Plan completed 時逐項對 §ADR Promotion Criteria 評估
5. 通過才建 ADR（直接 accepted）；不通過則改用更輕 promotion target（per ADR-007 路由表）

看到既有 proposed ADR 時：
1. 撤回 ADR 檔
2. 內容整併入對應 plan
3. constitution/README 移除該列
4. 補 failure pattern 記錄

#### Goal / Action / Validation

- Goal: constitution 純度（全 accepted），plan 為提案容器
- Action: 採 plan-first 流程；plan completed 後才評估 ADR promotion
- Validation or reference source: constitution/ 內所有 ADR Status ∈ {accepted, deprecated, superseded}；plan 含完整 §Decision Rationale；新 ADR 建立時引用 completed plan 為 evidence

#### Applies When

- Ai-skill 自身架構決策治理
- 採用 ADR-007 promotion pipeline 的專案
- Plan-driven 開發節奏

#### Does Not Apply When

- 業務專案的 ADR 治理（各專案自行評估是否採用此模式）
- 純 spike/exploratory 開發（無 plan 結構）
- 不採用 ADR pipeline 的決策記錄方式

#### Validation

- 對既有 constitution/ 跑 status 檢查，無 proposed/draft/pending
- 對既有 plans/active/ 跑 §Decision Rationale completeness 檢查
- 新 ADR 建立時可在 commit message 找到對應 completed plan 引用

#### Promotion Target

- ✅ `intelligence/engineering/architecture/plan-first-decision-promotion.md`（已於本次 commit 寫入）
- ✅ `knowledge/summaries/plan-first-decision-promotion.md`（summary card）
- ✅ Source ADR: [`constitution/ADR-007`](../../../constitution/ADR-007-constitution-and-decision-promotion-boundary.md)
- ✅ Canonical rule: [`governance/lifecycle/decision-promotion-pipeline.md`](../../../governance/lifecycle/decision-promotion-pipeline.md) §No-Proposed-ADR Rule
- ✅ Failure pattern: [`enforcement/failure-patterns/premature-adr-promotion.md`](../../../enforcement/failure-patterns/premature-adr-promotion.md)

#### Required Linked Updates

- `intelligence/engineering/architecture/README.md`（加入索引條目）
- `knowledge/summaries/README.md`（加入 summary 條目）
- Step 6（Intelligence Extraction）：done(executed) — atom 已 promote
- Step 7（Failure Learning）：not_applicable — 失效已沉澱於 `premature-adr-promotion.md`
