# Plan-First Decision Promotion（提案在 plan，憲法在驗證後）

**Status**: `candidate-intelligence`
**Source**: 通用架構治理經驗（generalised from ADR/RFC 治理模式）

## 原則

**架構決策的提案、討論與 alternatives 評估在 `plans/active/<plan>.md` 的 `§Decision Rationale` section 完成；只有 plan completed 且通過 `ADR Promotion Criteria` 後，才升級為 `accepted` ADR 寫入 constitution。憲法層只放已驗證的決策，不放 proposed/draft 決策。**

## 為什麼

1. **避免「廢棄憲法」累積**：proposed ADR 是「未驗證的憲法」，與 constitution 的「已驗證跨 session 決策集合」定位矛盾。失敗的 proposed ADR 會留在 constitution 變成噪音。
2. **狀態爆炸**：proposed → accepted → 反悔 → superseded → withdrawn 的狀態組合難維護。每個 reviewer 都要先檢查 status 才能判斷 weight。
3. **平行維護成本**：ADR 與 plan 同時記錄 Decision/Alternatives 內容，雙寫易漂移。
4. **Build-then-bless 邏輯**：「實作驗證後才升憲法」比「先寫憲法再實作」對工程師而言更直覺；後者容易紙上談兵。
5. **Plan 已是天然容器**：plan 既有 phases、phases 既有完成條件，把 Decision Rationale 加進 plan 不增加結構複雜度。

## 訊號（何時該套用此原則）

- 想為新架構/流程/跨層改動寫 ADR
- 想標記某個 ADR 為 `proposed` / `draft` / `pending`
- 已有對應 `plans/active/<plan>.md` 但仍想平行建 ADR
- 使用者問「這個要不要寫成 ADR」
- 看到既有 constitution/ 內有非 accepted status 的 ADR

## 操作流程

```
有新架構提案
  ↓
寫 plan（status: draft）
  必含 §Decision Rationale section：
    - Problem & Why Now
    - Decision
    - Alternatives Considered
    - Why Not an ADR Yet
    - ADR Promotion Criteria（completed 時驗證）
    - Consequences（預期）
  ↓
討論 Open Questions、調整 plan
  ↓
plan → status: in-progress
  ↓
Phase 0 pre-build interrogation
  ↓
Phase 1-N 執行
  ↓
plan → status: completed
  ↓
評估 §ADR Promotion Criteria：
  ✓ foundational + cross-session + cross-project + expensive-to-reverse + explains-why 全中？
  ✓ Plan 結果證實 decision 可行？
  ✓ Open Questions 全解？
  ✓ 沒有更輕的 promotion target 適用（per ADR-007）？
  ✓ 系統真實使用此 contract，具體 evidence 達標？
  ↓
  ├─ 通過 → 建立 ADR（直接 accepted），引用 completed plan
  └─ 不通過 → archive plan，改用更輕 promotion target（runtime gate / enforcement / intelligence）
```

## ADR Promotion Criteria（plan completed 時驗證）

升級為 accepted ADR 的必要條件（per ADR-007 §ADR Boundary）：

| 條件 | 翻譯 |
|------|------|
| foundational | 基礎性架構決策 |
| cross_session_and_cross_project | 跨 session、跨 project 都成立 |
| expected_to_remain_stable | 預期穩定，不會頻繁修改 |
| expensive_to_reverse | 反悔代價高（動結構/命名/外部 API） |
| explains_why_system_is_shaped_this_way | 解釋系統為何長這樣（why 而非 how） |

**5 條必須全部符合**。任一不符合 → 改用其他 promotion target。

## 替代 Promotion Targets（per ADR-007）

依決策內容選對應層，避免 ADR 過度膨脹：

| 決策內容 | 應放 |
|------|------|
| 可執行規則 / cross-agent policy | `enforcement/` |
| 推理 heuristic / tradeoff / signal / anti-pattern | `intelligence/` |
| 操作流程 / repeatable workflow | `workflow/` |
| Runtime gate / activation / phase / obligation | `runtime/runtime.db` |
| Session 範疇 replay 決策 | `memory/decision/` |
| 專案專屬決策 | `<PROJECT_ROOT>/docs/decisions/` |
| 架構級不可逆 + 5 條全中 | `constitution/ADR-*`（**直接 accepted**） |

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「跨 session 就該寫 ADR」 | 跨 session 是必要但非充分；還要 cross-project + foundational + expensive-to-reverse |
| 「先寫 proposed ADR，討論後改 accepted」 | 提案階段不該進憲法；plan 才是提案容器 |
| 「重要決策都寫 ADR 才正式」 | ADR 的「正式」反而是負擔（immutable、需 supersession）；輕量決策放 memory/decision/ 或更輕層 |
| 「memory/decision/ 永遠是 ADR 預備區」 | 不對 — 多數 memory/decision 永遠停在那層，不需升級 |
| 「升級世代要先寫 proposed ADR 規劃」 | 改為 plan 內 §Decision Rationale；ADR 在 plan completed 後寫 |

## 與其他智慧的關係

- [ADR-007](../../../constitution/ADR-007-constitution-and-decision-promotion-boundary.md)：本原則的 source ADR — constitution 命名與 promotion target 路由
- [`governance/lifecycle/decision-promotion-pipeline.md`](../../../governance/lifecycle/decision-promotion-pipeline.md) §No-Proposed-ADR Rule：本原則的 canonical executable rule（2026-05-22 確立）
- [`enforcement/failure-patterns/premature-adr-promotion.md`](../../../enforcement/failure-patterns/premature-adr-promotion.md)：違反本原則的 failure pattern
- [`plans/README.md`](../../../plans/README.md) §Plan 模板必填章節：強制 plan 含 §Decision Rationale
- [`migration-feature-bundling.md`](../anti-patterns/migration-feature-bundling.md)：相關「build-then-bless」邏輯（先 parity 後 feature）
- [`vendor-integration-architecture.md`](vendor-integration-architecture.md)：類似「先驗證再升級」邏輯（先 N 小，N 大後才升 SPI）

## 驗證

| 檢查 | 通過條件 |
|------|------|
| Constitution status | `constitution/` 內所有 ADR Status 為 `accepted` / `deprecated` / `superseded`（無 `proposed` / `draft` / `pending`） |
| Plan §Decision Rationale | 涉及架構/流程/跨層 plan 含完整 §Decision Rationale 6 子章節 |
| ADR creation timing | 新 ADR 建立時對應 plan 已 `completed`；commit message 引用 completed plan 為 evidence |
| 雙寫消除 | Plan 與 ADR 不重複維護同份 Decision Rationale；ADR 寫入後 plan 可保留為 implementation evidence |
| Promotion target 多元化 | 過去 N 個 plan completed，並非全部升級為 ADR — 應依內容路由至更輕 layer |

## Token Impact

對 reviewer：constitution 純度高（讀 7 個 ADR 全部都是 active canonical 決策），無需先 filter status。
對作者：Plan 模板強制 §Decision Rationale，等效於 ADR 的 forcing function；不需平行維護 ADR 草稿。
對 agent：discovery 時可直接信任 constitution；提案階段查 `plans/active/`。

## 邊界

| 適用 | 不適用 |
|------|------|
| Ai-skill 自身架構決策治理 | 業務專案的 ADR 治理（採用前各專案自行評估） |
| Plan-driven 開發節奏 | 純探索性 / spike-driven 開發（無 plan 結構） |
| 採用 ADR-007 promotion pipeline | 不採用 ADR pipeline 的決策記錄方式 |

---

← [回到 engineering/architecture/](README.md)
