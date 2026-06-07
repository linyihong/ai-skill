---
id: 2026-06-06-1800-sanitization-mechanical-enforcement
plan_kind: sub
status: draft
owner: linyihong
created: 2026-06-06
parent: 2026-05-31-2100-mechanical-enforcement-registry
priority: P1
required_for_completion: true
supersedes: 2026-05-31-2000-mechanical-sanitization-validator
sub_plan_reason: >
  Third sampling of the "rule-without-executor" meta-pattern (parent meta-plan
  Instance #3; sibling of 2026-05-31-2000-mechanical-sanitization-validator and
  2026-06-06-1700-workflow-activation-discovery-bridge). Empirical trigger
  reproduced 2026-06-06 in commit 214a415 — a downstream-project label
  appeared seven times across plans/active/ shared layer and was only caught
  after manual user inspection. Independent sign-off because the design
  principles (visibility-derived forbidden tokens, no allowlist, staged-content
  scan, attestation prohibited) differ from sibling sanitization-validator and
  must be governed as an independent surface to avoid silent merge of two
  conflicting executor philosophies.
---

# Sanitization: Mechanical Enforcement (Metadata-Derived)

**Status**: `draft`
Owner: framework maintainer (linyihong)
**世代**：Gen 3 runtime hardening — sanitization 第三採樣
**建立日期**：2026-06-06
**最後更新**：2026-06-06（v1 draft）
**Priority**：**P1**（升自 P2，因 supersede sibling 並承擔 `rule_classes[sanitization]` canonical executor 職責）
**Parent plan**：[`2026-05-31-2100-mechanical-enforcement-registry.md`](../archived/2026-05-31-2100-mechanical-enforcement-registry.md)
**Sibling plans**：
- [`2026-05-31-2000-mechanical-sanitization-validator.md`](2026-05-31-2000-mechanical-sanitization-validator.md) — in-progress allowlist 路線；**本 plan 建議 supersede 該 plan**（見 Decision §Relationship to sibling）。
- [`2026-06-06-1700-workflow-activation-discovery-bridge.md`](2026-06-06-1700-workflow-activation-discovery-bridge.md) — same-day sibling，同樣 third-sampling 同模式失效。

**Empirical trigger**：2026-06-06 session commit 214a415 — agent 將某 downstream project 的 label 七處寫入 `plans/active/` shared layer (該 plan 為 `2026-06-06-1700-workflow-activation-discovery-bridge` 的 v0 草稿)。`enforcement/sanitization.md` + `enforcement/reusable-guidance-boundary.md` 明文禁止「project incident 的具體 app / project 名稱」進入可重用文件，但無 mechanical executor 阻擋，agent 自律未生效，使用者手動指認後才 sanitize（commit 728282c）。

> 本 plan **不擴** cognitive mode block schema、**不改** detector ontology、**不動** workflow-activation engine。純粹建立 sanitization registry + staged-content scanner，作為 `rule_classes[sanitization]` 的 mechanical executor surface。Filename / project metadata semantic surface 議題留給條件式 follow-up。

---

## Decision Rationale

### Problem & Why Now

Parent meta-plan §Empirical Evidence 已將 sanitization 列為 instance #2；sibling `2026-05-31-2000` in-progress 處置該 gap。然而 sibling 採取的「規則先寫進 sanitization.yaml allowlist，executor 比對 allowlist」路徑需要 **長期維護 forbidden-token allowlist**，與 `enforcement/sanitization.md` 既有 prose rule 產生 dual source-of-truth。今日 commit 214a415 再度觸發同模式失效，證明**等 allowlist 路徑成熟之前，沒有任何 mechanical 防線**。

需要一個 **不依賴 allowlist** 的補強：以「**該 token 在 canonical Ai-skill repo 是否可見**」作為 forbidden 判斷基準（visibility-derived），避免雙 SOT。

### Design Principles (user-established 2026-06-06)

| # | 原則 | 理由 |
|---|---|---|
| P1 | **Metadata-derived forbidden tokens** | 不掃描推導 visibility，而是 project 自己在 metadata 宣告 `visibility: private` + `private_tokens: [...]`。`ai-skill runtime compile` 把所有 project 的宣告 project metadata projection 為 `runtime.db.derived_forbidden_tokens`。Scanner 只查 projection，不做 absence inference。命名理由：「visibility-derived」會暗示 scanner 自動推導，本 plan 明確否決該路線；token 是 metadata-declared，derivation 只發生在 case-variant 展開與 cross-project aggregation。 |
| P2 | **Staged-content scan，不 scan commit message** | `git diff --cached` 的內容；不檢查 `commit -m` 文字。Commit message 由 cognitive mode block + 既有 commit-msg validators 治理；本 plan 不重疊。 |
| P3 | **覆蓋範圍：shared-layer classification（非 folder name classification）** | `runtime/` 與 `knowledge/runtime/` 邊界會持續模糊；用「該檔是否屬 shared / reusable layer」分類。**Topology surface 獨立於 enforcement**：新建 `runtime/repository-topology.yaml`（canonical），宣告 `shared_layer: true\|false` per source-tree subtree。Scanner 查 topology 而非硬編 folder glob。Topology 將被 sanitization / workflow activation / governance lint / dependency reading 等多個 subsystem 共用 — 不掛 `enforcement-registry.yaml` 避免它變超級桶。 |
| P4 | **Attestation 禁止** | 不接受 commit body 內 "Sanitization: yes" / "[sanitized]" 自陳。Validator 只做 verification，不做 attestation。理由：自陳是 agent 主觀宣告，與本 plan 要解決的「自律失效」同根。**此原則為 cross-cutting governance，預期被其他 obligation 重用（Dependency Read / Test Executed / Coverage Reviewed 等）。** |
| P5 | **Phase 順序：label registry → 通用 regex → LLM review (conditional)** | Phase 1 最便宜（純 token 比對）；Phase 2 涵蓋 email / phone / OS path / credential pattern 等可 regex 化的；Phase 3 為前兩階段 surface 不夠時才開（cost / determinism trade-off 留 Phase 3 自己 ADR）。 |
| P6 | **Bootstrap-safe：forbidden 由 project metadata 宣告，非「reusable absent」反推** | 「reusable layer 不存在 → forbidden」會誤殺新概念（e.g. 首次 commit `ActivationBridgeV2` 時 reusable layer 還沒有）。Forbidden 判定改為單一條件：token 出現在某 project 宣告的 `private_tokens` 內。新 framework concept 不在任何 project private_tokens 內 → 自動豁免。完全不做 visibility inference，不做 absence inference。 |

### Decision

建立 **Sanitization Mechanical Enforcement**（rule_class `sanitization` 的 canonical executor，**supersede** sibling `2026-05-31-2000` allowlist 路線 — 見 §Relationship to sibling）：

```
git pre-commit hook
  │
  ├─ Scope filter: shared-layer classification from runtime/repository-topology.yaml
  │
  ├─ Phase 1: metadata-derived forbidden token scan
  │     derived_forbidden_tokens = compile-time projection from
  │       project metadata `private_tokens: [...]` across all known projects
  │       (case-variant expansion: CamelCase / kebab-case / SCREAMING_SNAKE)
  │     for each staged file content in shared layer:
  │       for each token in derived_forbidden_tokens:
  │         if token literal/case-variant present:
  │           emit finding (file, line, token, owning_project, suggested_placeholder)
  │
  ├─ Phase 2: generic regex pattern scan
  │     email | phone | OS absolute path (Windows / POSIX) | credential pattern
  │
  ├─ Phase 3 (deferred): LLM review hook
  │     conditional on Phase 1+2 miss rate; separate ADR
  │
  └─ Exit: any finding → block commit; zero → pass
```

**Source-of-truth (deterministic, not inferred)**：
- **Forbidden tokens 由 project metadata 宣告**：每個 project (含 `.agent-goals/` project-local 目錄、downstream consumer overlay) 在自身 metadata 內宣告 `private_tokens: [...]`。`ai-skill runtime compile` 把所有 known project 的 `private_tokens` projection 為 `runtime.db.derived_forbidden_tokens`（含 case-variant expansion + cross-project aggregation）。Scanner 比對該 table，**不做 absence inference、不做 visibility inference**。新 framework concept（如 `ActivationBridgeV2`）因不在任何 project private_tokens 內，bootstrap-safe 不誤殺。
- **Shared-layer classification**：canonical 來源是 **新建 `runtime/repository-topology.yaml`**（subtree → `shared_layer: true|false` map）。Topology 是 cross-subsystem surface（預期 sanitization / workflow activation / governance lint / dependency reading 共用），刻意不掛 `enforcement-registry.yaml` 避免該 registry 變超級桶。Scanner 從 topology projection 載入，硬編 folder glob 禁止。
- **Phase 2 regex patterns**：canonical `runtime/sanitization-patterns.yaml`，與 `enforcement/sanitization.md` companion，single SOT。

### Relationship to sibling 2026-05-31-2000 (supersede)

兩條路徑的核心是 **explicit registry vs derived registry** 的哲學分歧，不是 implementation 差異。同時保留必然導致「allowlist 說可以、visibility 說不可以」的 conflicting verdict 案例，rule_class 治理失靈。

**裁決：本 plan supersede sibling `2026-05-31-2000`**。

理由：
1. Sibling 路線需要每加一個 downstream project 手動更新 allowlist → 與 P1 visibility-derived 原則衝突且 dual SOT
2. 2026-06-06 incident 證明 allowlist cold-start window 存在 systemic risk
3. P6 bootstrap-safe 條件已解決原本傾向 sibling 的「visibility 推論誤殺新概念」風險
4. Registry rule_class 一個 class 一個 canonical executor，避免 cross-executor verdict drift

執行：Phase 0 完成後 sibling plan archive + status=superseded，本 plan 升 P1 + `required_for_completion: true`，parent meta-plan §Child plans 同步更新。

### Why Not Just Wait for Sibling Plan

| 角度 | Sibling `2026-05-31-2000` (allowlist-based) | This plan (visibility-derived) |
|---|---|---|
| Source-of-truth | `enforcement/sanitization.yaml` allowlist | 跟著 canonical repo 可見性走，不另建 allowlist |
| Maintenance cost | 每加一個新 project 要更新 allowlist | 零 — 新 project 進入 `.agent-goals/` 自動被 derive |
| Drift risk vs `sanitization.md` prose | 中高（dual SOT） | 低（無 dual SOT） |
| 對 today's incident (214a415) | 若 allowlist 尚未含該 label → miss | 命中（label 在 project-local 出現、在 reusable 缺席） |
| 與 sibling 共存？ | 兩個 executor 同 rule_class | 視 Q1 裁決 |

兩條路徑各自處理不同失誤模式。今日 incident 證明 allowlist 路徑有 cold-start window，本 plan 補的就是 cold-start window 的 mechanical 防線。

---

## Architecture Compatibility Preflight

| 欄位 | 內容 |
|---|---|
| Candidate files | 新建 `runtime/sanitization-patterns.yaml`（canonical Phase 2 regex）；新建 `scripts/ai-skill-cli/internal/app/sanitization_scan.go`（scanner）；新建 validator entry in `scripts/ai-skill-cli/internal/app/hooks.go` pre-commit dispatcher；新建 `enforcement/sanitization-mechanical.md`（companion，philosophy + 與 prose `sanitization.md` 邊界）；`enforcement/enforcement-registry.yaml` 更新 `rule_classes[sanitization]` executors block（或新增 second executor entry，依 Q1） |
| Source-of-truth | `enforcement/sanitization.md` 仍是 prose rule canonical；`runtime/sanitization-patterns.yaml` 是 Phase 2 regex canonical；visibility-derived Phase 1 沒有 forbidden-token canonical（依存可見性查詢） |
| Compiler / generated surfaces | `runtime/sanitization-patterns.yaml` 經 `ai-skill runtime compile` projection 到 `runtime.db.generated_surfaces`；scanner 從 runtime.db 載入 patterns |
| Layer responsibility | enforcement-mechanical 屬 enforcement layer (rule_class executor)；scanner 屬 runtime layer；pre-commit hook 整合屬 ai-skill-cli layer |
| 與現行架構衝突 | **與 sibling `2026-05-31-2000` 可能衝突**；Q1 必須先裁決 |
| `runtime.db` 影響 | 新增 `sanitization_patterns` projection table；無 schema migration risk（純讀） |

---

## Scope Seal

**In scope**：
- Pre-commit `git diff --cached` content scanner，scope filter = `plans/** workflow/** enforcement/** governance/** knowledge/**`
- Phase 1: visibility-derived project label scan
- Phase 2: generic regex（email / phone / OS absolute path / credential pattern）
- Block-on-finding semantics；無 attestation surface
- `rule_classes[sanitization]` registry entry 更新

**Out of scope (本 plan 不處理)**：
- Cognitive mode block schema 擴張（NONE）
- Detector ontology / route_type 改動（NONE）
- Workflow Activation Engine 任何變動（NONE）
- Commit message 文字檢查（既有 cognitive validators 已治理）
- `scripts/**` / `runtime/**` / `validation/**` 等 implementation surface（Phase 1 v2 才評估）
- LLM-based semantic review（Phase 3，獨立 ADR）
- 自陳 attestation marker（P4 明確禁止）
- 修個案 leak（今日 commit 214a415 已由 user 手動修正於 728282c）

---

## Phase Plan

### Phase 0 — Preflight & Sibling Supersede

- [ ] User 認可 §Relationship to sibling 的 supersede 裁決
- [ ] Sibling `2026-05-31-2000` 移入 `plans/archived/` + frontmatter `status: superseded` + `superseded_by: 2026-06-06-1800-sanitization-mechanical-enforcement`
- [ ] 本 plan frontmatter 升 `priority: P1` + `required_for_completion: true`
- [ ] Parent meta-plan §Child plans 列表同步更新（移除 sibling 進入 archived 引用，加入本 plan）
- [ ] Phase 0 不啟動任何 implementation

### Phase 1 — Metadata-Derived Forbidden Token Scanner

- [ ] 定義 project metadata schema：`<PROJECT_ROOT>/.ai-skill-project.yaml` (或 overlay 既有 metadata 機制)，欄位：
  ```yaml
  project:
    id: <project-slug>
    visibility: private | public
    private_tokens:
      - <ProjectName>
      - <project-name>
      - <PROJECT_NAME>
      # case variants 自動 derive (CamelCase / kebab / SCREAMING_SNAKE)
  ```
- [ ] 新建 `runtime/repository-topology.yaml`（canonical cross-subsystem topology surface，**不掛 enforcement-registry**）：
  ```yaml
  # runtime/repository-topology.yaml
  schema_version: 1
  shared_layer_classification:
    - subtree: plans/
      shared: true
    - subtree: workflow/
      shared: true
    - subtree: enforcement/
      shared: true
    - subtree: governance/
      shared: true
    - subtree: knowledge/
      shared: true
    - subtree: runtime/
      shared: true   # reusable framework runtime
    - subtree: scripts/
      shared: false  # implementation
    - subtree: validation/
      shared: true   # reusable scenarios
    - subtree: .agent-goals/
      shared: false  # project-local
  expected_consumers:
    - sanitization (this plan)
    - workflow_activation (path-signal context)
    - governance_lint (future)
    - dependency_reading (future)
  ```
- [ ] Topology projection target：`runtime.db.repository_topology`（query via `ai-skill runtime ...` CLI）
- [ ] `ai-skill runtime compile` 階段 projection：所有 known project metadata 的 `private_tokens` → `runtime.db.derived_forbidden_tokens` table (含 case variants expansion)
- [ ] 實作 scanner core：staged file 落在 shared-layer subtree → 比對 derived_forbidden_tokens → emit finding
- [ ] **Bootstrap-safe guard**：scanner 不做「reusable layer 是否包含 token」inference；新 framework concept 只要不在任何 project private_tokens 內，自動通過
- [ ] False-positive guardrail：`enforcement/sanitization.md` 自身的 example 段落 self-reference exception
- [ ] Unit tests：
  - fail case：commit 214a415 reconstruction（七處 leak 必須全部 emit finding）
  - pass case：commit 728282c 後內容（零 finding）
  - bootstrap-safety case：首次提交一個全新 framework concept token，確認 0 finding

### Phase 2 — Generic Regex Patterns (inherited from superseded plan)

- [ ] 新建 `runtime/sanitization-patterns.yaml`（canonical）+ companion section in `enforcement/sanitization-mechanical.md`
- [ ] Pattern set：
  - Email：`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
  - Phone：國際 + 區域格式（保留 conservative，避免誤抓 plan timestamp）
  - OS absolute path：`/Users/[^/\s]+` / `/home/[^/\s]+` / `C:\\Users\\[^\\\s]+` / `\\\\\\\\?\\\\` UNC
  - Credential pattern：高熵字串 + 已知 prefix（`sk-` / `ghp_` / `xoxb-` 等）
- [ ] Pattern allowlist override：placeholder forms `<USER>` / `<PROJECT_ROOT>` / `<AI_SKILL_REPO>` 必須被 pattern 排除
- [ ] `ai-skill runtime compile` projection → `runtime.db.sanitization_patterns`
- [ ] Unit tests：每 pattern ≥ 1 fail + 1 pass + 1 placeholder pass

### Phase 2.5 — Incident-Score Heuristic (warn only, inherited from superseded plan)

補強 Phase 1 visibility-derived 抓不到的「無 private token 但明顯是 project incident」案例（連續的 domain-specific noun cluster、quoted user text、具體 dated filename 引用等）。**warn-only，不 block**，理由：heuristic 性質與 P4 attestation-prohibited 的 deterministic 精神有張力，降為 advisory 守住 deterministic-block 邊界。

- [ ] 將 sibling plan §Phase 1 `incident_score` weighted-signal schema 整段搬入 `runtime/sanitization-patterns.yaml` `incident_score:` 區段（filename_pattern / quoted_user_text / artifact_string / domain_noun_cluster，weights 5 / 5 / 3 / 1，warn_if_total_score_gte: 7）
- [ ] Scanner 在 Phase 1+2 全部 pass 時才執行 incident-score；任一 phase block → 跳過（避免 noise）
- [ ] Output 為 `warning` 等級，commit 不 block；finding 寫入 stderr 供 agent / user 評估
- [ ] Unit tests：sibling plan §incident_score `examples` 三個 case（leak 命中、registry 內合法描述放行、archived/feedback 路徑跳過）必須全部 pass
- [ ] Acceptance：incident-score 不得對 Phase 1 已 derived_forbidden 的 token 重複 finding（避免雙報）

### Phase 3 — LLM Review Hook (Conditional, Deferred)

- [ ] Trigger condition：Phase 1+2 三個月真實使用 miss rate 量測；若 > X% (X 待 ADR 定) → 啟動 Phase 3 ADR
- [ ] Phase 3 自己一個 sub-plan（本 plan 不展開）

### Phase 4 — Registry & Bootstrap Integration

- [ ] 新建 `enforcement/failure-patterns/sanitization-leak-on-canonical-write.md`（inherited 規劃，empirical evidence: 214a415 + sibling plan v1-v4 累積 leak）
- [ ] 更新 `enforcement/enforcement-registry.yaml` `rule_classes[sanitization]`：
  - 新增 executor entry: `file: scripts/ai-skill-cli/internal/app/sanitization_scan.go`, `symbol: scanStagedContentForSanitization`, `hook_phase: pre-commit`, `block_or_warn: block`
  - 若 Q1 裁決 supersede → 移除舊 allowlist executor entry
  - `coverage` 視 promotion criteria 是否成立決定（pending_implementation → mechanical）
- [ ] 更新 `runtime/core-bootstrap.yaml` `per_commit_obligations` 新增 `obligation.commit.sanitization_visibility_scan`（**無 opt_out_marker**，per P4 attestation-prohibited + Q5 reject）
- [ ] `commitMsgValidatorRegistry` 在 hooks.go 註冊新 obligation
- [ ] Companion `enforcement/sanitization-mechanical.md` cross-link to prose `sanitization.md`
- [ ] Re-dry-run `ai-skill runtime compile` + enforcement coverage report
- [ ] Validation scenarios：
  - `validation/scenarios/runtime/sanitization-visibility-derived-pass-v1.yaml`
  - `validation/scenarios/runtime/sanitization-visibility-derived-fail-v1.yaml`（commit 214a415 reconstruction）
  - `validation/scenarios/runtime/sanitization-placeholder-allowed-v1.yaml`
- [ ] Owner-grouped commit + push + readback

---

## Open Questions

| # | Question | Pri | Status | 處置 |
|---|---|---|---|---|
| Q1 | ~~與 sibling 關係~~ | — | **resolved 2026-06-06** | supersede。詳見 §Relationship to sibling。 |
| Q2 | Visibility metadata 載入機制：scan `<PROJECT_ROOT>/.ai-skill-project.yaml`，還是嵌入既有 routing/overlay metadata？ | P1 | open | Phase 1 ADR；不影響本 plan 接受度。 |
| Q3 | Phase 3 LLM review 的 miss-rate threshold 是多少才 trigger？ | P2 | open | Phase 3 sub-plan 自己決定。 |
| Q4 | ~~Shared-layer classification 進 enforcement-registry？~~ | — | **resolved 2026-06-06: 否** | 改建獨立 `runtime/repository-topology.yaml` cross-subsystem surface。Rationale: topology 將被 sanitization / workflow activation / governance lint / dependency reading 共用，掛 enforcement-registry 會讓該 registry 變超級桶（同 parent meta-plan R6 警告的反模式）。 |
| Q5 | ~~`[skip-sanitization-scan]` opt-out marker？~~ | — | **resolved 2026-06-06: reject** | 與 P4 attestation-prohibited 衝突。若未來需要 emergency override，須另開 ADR 設計 admin-override surface（含 owner / reason / time-boxed expiry），不採 commit-message marker。 |
| Q6 | False-positive 處理機制：suggested_placeholder 是否自動 patch staged content？ | P2 | open | v1 不自動 patch；自動 patch 列 v2 評估。 |
| Q7 | P4 attestation-prohibited 原則升級為 cross-cutting governance：何時抽出獨立 enforcement rule (`enforcement/verification-not-declaration.md`)，覆蓋 Dependency Read / Test Executed / Coverage Reviewed 等？ | P2 | open | 本 plan 內留 inline；累積 ≥3 個受益 obligation 後 promote 為獨立 rule，列入 parent meta-plan tracking。 |

---

## Validation

| 欄位 | 內容 |
|---|---|
| Trigger | 2026-06-06 session 確認本 plan 已 commit + push（per `enforcement/dependency-reading.md` writeback transaction） |
| Required set | `enforcement/sanitization.md` / `enforcement/reusable-guidance-boundary.md` / `plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md` / `knowledge/runtime/routing-registry.yaml` (sanitization slots) / `scripts/ai-skill-cli/internal/app/hooks.go` (dispatcher) |
| Read | ✅ 全部 |
| Not applicable | — |
| Deferred / blocked | Q1 裁決前不啟動 Phase 1；sibling plan diff review 留 Q1 之後 |
| Validation | 本 plan 是 draft；驗證為 user 對 Q1 + Decision Rationale 認可後升 `status: in-progress` |

---

## Inherited from Superseded Plan

從 `plans/archived/2026-05-31-2000-mechanical-sanitization-validator` 繼承的設計遺產（哲學一致、可直接吸收）：

| 來源段落 | 吸收進本 plan | 變動 |
|---|---|---|
| `banned_patterns.filesystem_reference` / `username` / `secrets_token` / `private_host`（4 類 ~15 條 regex） | Phase 2 `runtime/sanitization-patterns.yaml` | 直接搬入；placeholder allowlist override 保留 |
| `incident_score` weighted-signal heuristic + examples | Phase 2.5（warn-only） | 從 sibling 的 block-level 降為 warn-level，與 P4 attestation-prohibited 的 deterministic-block 邊界一致 |
| Commit-diff-only 立場（PreToolUse 不 block，避免 Edit 看到舊內容誤判） | 本 plan 從一開始就採此立場（pre-commit only），無 PreToolUse 層 | 不變 |
| `failure-patterns/sanitization-leak-on-canonical-write.md` 建立規劃 | Phase 4 任務追加 | 保留為 framework learning artifact，不屬 executor |

**棄用段落（不繼承）**：

| 段落 | 棄用理由 |
|---|---|
| `sanitization.yaml` 內 `incident_score` block-level threshold | 由 Phase 2.5 warn-level 取代 |
| 任何 allowlist-style `private_project_names:` / `allowed_tokens:` 設計 | 被 visibility-derived `private_tokens` 取代 |
| `[skip-sanitization]` opt-out marker | P4 + Q5 reject |
| `canonical_paths` / `not_canonical` 硬列清單 | 由 `shared_layer_classification` table 取代 |
| PreToolUse Write/Edit warning hook | commit-diff-only 立場下不需要（避免 Edit `old_string` 誤判） |

## Handoff Notes

- 本 plan 與 sibling `2026-05-31-2000` 同 coverage class，Q1 不裁決前不可雙線 implementation（避免 registry rule_class 出現 conflicting executor）
- Empirical trigger commit 214a415 + 修正 commit 728282c 是 Phase 1 unit test 的 ground truth dataset
- Phase 4 Registry integration 必須與 sibling 協調 — `enforcement-registry.yaml` 同檔，避免 commit conflict
- 若 user 將本 plan 升 P1，需同步 archive 或降級 sibling，並更新 parent meta-plan §Child plans 列表
