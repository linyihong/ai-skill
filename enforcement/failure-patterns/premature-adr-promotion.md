# Premature ADR Promotion（在 plan completed 前就寫 ADR）

Status: validated
Class: `process-gap` / `governance-drift`

## Trigger

當 agent 為新架構提案建立 `constitution/ADR-<n>-<slug>.md` 並標 `Status: proposed`，**在對應 plan completed 前**就把 ADR 寫入 constitution/，使用此 pattern。

具體觸發訊號：

- 新增 ADR 檔，Status 為 `proposed` / `draft` / `pending`（任何非 `accepted` 狀態）
- 對應 `plans/active/<plan>.md` 仍是 `draft` 或 `in-progress`
- ADR 內含 §Open Questions 或「待審查」「待驗證」字眼
- ADR 與 plan 互相 cross-reference，但 plan 的 §ADR Promotion Criteria 未通過

## Failure Mode

把「未驗證的架構提案」放進 `constitution/`，違反 constitution 的核心定位（**accepted = 已驗證、跨 session+project、難回退的決策**）。

具體後果：

1. **「廢棄憲法」累積**：失敗或被取代的 proposed ADRs 留在 constitution/ 成為噪音
2. **狀態爆炸**：proposed → accepted → 反悔 → superseded → withdrawn 的狀態組合難維護
3. **誤導讀者**：新人讀 constitution/ 不知道哪些是 canonical 哪些是提案
4. **平行維護成本**：ADR 與 plan 同時記錄 Decision/Alternatives，內容易漂移
5. **違反 ADR-007 §Decision**：「ADR is NOT the default endpoint」— proposed ADR 是違反此原則的快速路徑

## Risk

- **Constitution 純度損失**：constitution/ 不再代表「已驗證的架構決策集合」
- **Forcing function 弱化**：plan 提案如果有 ADR 平行存在，反而傾向用 ADR 表達而忽略 plan 的 phase 設計
- **Promotion gate 模糊**：proposed → accepted 的 gate 沒有明確 criteria，容易 silently flip
- **Reviewer 認知負擔**：每個 ADR 都要先檢查 status 才能判斷 weight

## Required Agent Action

當有新架構提案時：

1. **不在 constitution/ 建 ADR 檔**（無論 status 標什麼）
2. **在 `plans/active/<plan>.md` 內加 §Decision Rationale section**，包含原本要寫進 proposed ADR 的內容：
   - Problem & Why Now
   - Decision
   - Alternatives Considered
   - **Why Not an ADR Yet**（明寫為什麼此階段不寫 ADR）
   - ADR Promotion Criteria（completed 時驗證）
   - Consequences（預期）
3. **Plan 經 Phase 0-N 執行至 completed**
4. **Completed 時評估 ADR Promotion Criteria**：
   - 若通過 → 建立 ADR（**直接 accepted**），引用 completed plan 為 evidence
   - 若不通過 → 改用更輕的 promotion target（runtime gate / enforcement / intelligence），不寫 ADR

## Prevention Gate

開始任何架構提案前，agent 必須能回答：

| Check | Required answer |
|-------|-----------------|
| Plan 是否已存在 | `plans/active/<plan>.md` 已建立，含 §Decision Rationale |
| 是否打算建 proposed ADR | **必須是「否」**；提案內容已在 plan 內 |
| §Why Not an ADR Yet | Plan 中已明寫 |
| §ADR Promotion Criteria | Plan 中已列出可驗證條件 |
| ADR 建立時機 | 僅在 plan completed + criteria pass 後 |

若回答不確定，**先讀** [`governance/lifecycle/decision-promotion-pipeline.md`](../../governance/lifecycle/decision-promotion-pipeline.md) §No-Proposed-ADR Rule。

## 驗證

符合下列條件時，此 pattern 已被防止：

- `constitution/` 中所有 ADRs 的 Status 為 `accepted` / `deprecated` / `superseded`（無 `proposed` / `draft` / `pending`）
- 新架構提案均存在於 `plans/active/<plan>.md` 而非 `constitution/`
- Plan 與 ADR 不重複維護同份 Decision Rationale
- ADR 編號連續且僅在實作驗證後分配

## Source

- 2026-05-21 ~ 2026-05-22 session：agent 為 Runtime Cognitive Modes 提案建立 `constitution/ADR-008-runtime-cognitive-modes.md` (Status: proposed) 與對應 plan (Status: draft) 平行存在。使用者指出「以後計畫的東西在 plan 寫好，完成後才寫到 constitution」即觸發本 pattern。ADR-008 已於 2026-05-22 撤回，內容整併入 plan 的 §Decision Rationale。
- Related decision: [`governance/lifecycle/decision-promotion-pipeline.md`](../../governance/lifecycle/decision-promotion-pipeline.md) §No-Proposed-ADR Rule（2026-05-22 確立）

## Related

- [`governance/lifecycle/decision-promotion-pipeline.md`](../../governance/lifecycle/decision-promotion-pipeline.md) — §No-Proposed-ADR Rule canonical 出處
- [`governance/lifecycle/decision-promotion-pipeline.yaml`](../../governance/lifecycle/decision-promotion-pipeline.yaml) — `blocking_rules.no_proposed_adr` 機器化執行
- [`constitution/README.md`](../../constitution/README.md) — ADR 生命週期已移除 `proposed` 狀態
- [`plans/README.md`](../../plans/README.md) — Plan 模板必填 §Decision Rationale section
- [`constitution/ADR-007-constitution-and-decision-promotion-boundary.md`](../../constitution/ADR-007-constitution-and-decision-promotion-boundary.md) — ADR 不是 default endpoint 原則
- [`knowledge-update-flow-bypassed-by-sub-pipeline.md`](knowledge-update-flow-bypassed-by-sub-pipeline.md) — 同性質的 process-gap 失誤

← [Back to failure patterns](README.md)
