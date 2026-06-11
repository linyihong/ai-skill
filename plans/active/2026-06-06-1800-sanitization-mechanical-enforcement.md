---
id: 2026-06-06-1800-sanitization-mechanical-enforcement
plan_kind: sub
status: in-progress
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
  principles (metadata-derived forbidden tokens, no allowlist, staged-content
  scan, attestation prohibited) differ from sibling sanitization-validator and
  must be governed as an independent surface to avoid silent merge of two
  conflicting executor philosophies.
---

# Sanitization: Mechanical Enforcement (Metadata-Derived)

**Status**: `in-progress`
Owner: framework maintainer (linyihong)
**世代**：Gen 3 runtime hardening — sanitization 第三採樣
**建立日期**：2026-06-06
**最後更新**：2026-06-08（Phase 0 preflight complete）
**Priority**：**P1**（升自 P2，因 supersede sibling 並承擔 `rule_classes[sanitization]` canonical executor 職責）
**Parent plan**：[`2026-05-31-2100-mechanical-enforcement-registry.md`](../archived/2026-05-31-2100-mechanical-enforcement-registry.md)
**Sibling plans**：
- [`../archived/2026-05-31-2000-mechanical-sanitization-validator.md`](../archived/2026-05-31-2000-mechanical-sanitization-validator.md) — superseded allowlist 路線；本 plan 是 canonical metadata-derived successor（見 Decision §Relationship to sibling）。
- [`2026-06-06-1700-workflow-activation-discovery-bridge.md`](2026-06-06-1700-workflow-activation-discovery-bridge.md) — same-day sibling，同樣 third-sampling 同模式失效。

**Empirical trigger**：2026-06-06 session commit 214a415 — agent 將某 downstream project 的 label 七處寫入 `plans/active/` shared layer (該 plan 為 `2026-06-06-1700-workflow-activation-discovery-bridge` 的 v0 草稿)。`enforcement/sanitization.md` + `enforcement/reusable-guidance-boundary.md` 明文禁止「project incident 的具體 app / project 名稱」進入可重用文件，但無 mechanical executor 阻擋，agent 自律未生效，使用者手動指認後才 sanitize（commit 728282c）。

> 本 plan **不擴** cognitive mode block schema、**不改** detector ontology、**不動** workflow-activation engine。純粹建立 sanitization registry + staged-content scanner，作為 `rule_classes[sanitization]` 的 mechanical executor surface。Filename / project metadata semantic surface 議題留給條件式 follow-up。

---

## Decision Rationale

### Problem & Why Now

Parent meta-plan §Empirical Evidence 已將 sanitization 列為 instance #2；sibling `2026-05-31-2000` in-progress 處置該 gap。然而 sibling 採取的「規則先寫進 sanitization.yaml allowlist，executor 比對 allowlist」路徑需要 **長期維護 forbidden-token allowlist**，與 `enforcement/sanitization.md` 既有 prose rule 產生 dual source-of-truth。今日 commit 214a415 再度觸發同模式失效，證明**等 allowlist 路徑成熟之前，沒有任何 mechanical 防線**。

需要一個 **不依賴中央 allowlist** 的補強：以 **「該 token 是否被任一 project metadata 宣告為 private」** 作為 forbidden 判斷基準（metadata-derived），避免 framework 端與 prose `sanitization.md` 形成雙 SOT。

### Design Principles (user-established 2026-06-06)

| # | 原則 | 理由 |
|---|---|---|
| P1 | **Metadata-derived forbidden tokens** | 不掃描推導 visibility，而是 project 自己在 metadata 宣告 `visibility: private` + `private_entities: [...]`（每個 entity 含 `name` / `kind` / `match_tokens`）。`ai-skill runtime compile` 把所有 project 的 entity 宣告 projection 為 `runtime.db.derived_private_entities` (governance layer) + `runtime.db.derived_match_tokens` (execution layer)。Scanner 只查 `derived_match_tokens`，不做 absence inference。命名理由：早期用過「visibility-derived」一詞會誤導 reviewer 以為 scanner 自動推導 visibility（本 plan 明確否決該路線）；改為 metadata-derived 強調 token 是 project metadata 宣告，derivation 只發生在 case-variant 展開與 cross-project aggregation。**Entity vs token 分層 see Q8 resolution + §Phase 1A**。 |
| P2 | **Staged-content scan，不 scan commit message** | `git diff --cached` 的內容；不檢查 `commit -m` 文字。Commit message 由 cognitive mode block + 既有 commit-msg validators 治理；本 plan 不重疊。 |
| P3 | **覆蓋範圍：shared-layer classification（非 folder name classification）** | `runtime/` 與 `knowledge/runtime/` 邊界會持續模糊；用「該檔是否屬 shared / reusable layer」分類。**Topology surface 獨立於 enforcement**：新建 `runtime/repository-topology.yaml`（canonical），宣告 `shared_layer: true\|false` per source-tree subtree。Scanner 查 topology 而非硬編 folder glob。Topology 將被 sanitization / workflow activation / governance lint / dependency reading 等多個 subsystem 共用 — 不掛 `enforcement-registry.yaml` 避免它變超級桶。 |
| P4 | **Attestation 禁止** | 不接受 commit body 內 "Sanitization: yes" / "[sanitized]" 自陳。Validator 只做 verification，不做 attestation。理由：自陳是 agent 主觀宣告，與本 plan 要解決的「自律失效」同根。**此原則為 cross-cutting governance，預期被其他 obligation 重用（Dependency Read / Test Executed / Coverage Reviewed 等）。** |
| P5 | **Phase 順序：label registry → 通用 regex → LLM review (conditional)** | Phase 1 最便宜（純 token 比對）；Phase 2 涵蓋 email / phone / OS path / credential pattern 等可 regex 化的；Phase 3 為前兩階段 surface 不夠時才開（cost / determinism trade-off 留 Phase 3 自己 ADR）。 |
| P6 | **Bootstrap-safe：forbidden 由 project metadata 宣告，非「reusable absent」反推** | 「reusable layer 不存在 → forbidden」會誤殺新概念（e.g. 首次 commit `ActivationBridgeV2` 時 reusable layer 還沒有）。Forbidden 判定改為單一條件：token 出現在某 project 宣告的 `private_entities[*].match_tokens` 內。新 framework concept 不在任何 project entity 的 match_tokens 內 → 自動豁免。完全不做 visibility inference，不做 absence inference。**Schema canonical 見 §Phase 1A**（entity = governance layer，match_tokens = execution layer）。 |

### Decision

建立 **Sanitization Mechanical Enforcement**（rule_class `sanitization` 的 canonical executor，**supersede** sibling `2026-05-31-2000` allowlist 路線 — 見 §Relationship to sibling）：

```
git pre-commit hook
  │
  ├─ Scope filter: shared-layer classification from runtime/repository-topology.yaml
  │
  ├─ Phase 1: metadata-derived forbidden token scan
  │     derived_private_entities  = compile-time projection of
  │       project metadata `private_entities: [...]` (governance layer:
  │       entity name + kind + owning project; NOT case-expanded)
  │     derived_match_tokens      = case-variant expansion of every
  │       entity's match_tokens (execution layer: scanner direct query)
  │     for each staged file content in shared layer:
  │       for each token in derived_match_tokens:
  │         if token literal/case-variant present:
  │           emit finding (file, line, matched_token, entity_name,
  │                         owning_project, suggested_placeholder)
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
- **Forbidden tokens 由 project metadata 宣告**：每個 project (含 `.agent-goals/` project-local 目錄、downstream consumer overlay) 在自身 metadata 內宣告 `private_entities: [...]`（governance layer：entity 帶 `name` / `kind` / `match_tokens`）。`ai-skill runtime compile` projection 為兩張 table：`runtime.db.derived_private_entities`（governance query — entity name + kind + owning project）+ `runtime.db.derived_match_tokens`（execution layer：每 entity 的 match_tokens × case-variant expansion × cross-project aggregation 完全展開）。Scanner 比對 `derived_match_tokens`，**不做 absence inference、不做 visibility inference**。新 framework concept（如 `ActivationBridgeV2`）因不在任何 project entity 的 match_tokens 內，bootstrap-safe 不誤殺。
- **Shared-layer classification**：canonical 來源是 **新建 `runtime/repository-topology.yaml`**（subtree → `shared_layer: true|false` map）。Topology 是 cross-subsystem surface（預期 sanitization / workflow activation / governance lint / dependency reading 共用），刻意不掛 `enforcement-registry.yaml` 避免該 registry 變超級桶。Scanner 從 topology projection 載入，硬編 folder glob 禁止。
- **Phase 2 regex patterns**：canonical `runtime/sanitization-patterns.yaml`，與 `enforcement/sanitization.md` companion，single SOT。

### Relationship to sibling 2026-05-31-2000 (supersede)

兩條路徑的核心是 **explicit registry vs derived registry** 的哲學分歧，不是 implementation 差異。同時保留必然導致「allowlist 說可以、visibility 說不可以」的 conflicting verdict 案例，rule_class 治理失靈。

**裁決：本 plan supersede sibling `2026-05-31-2000`**。

理由：
1. Sibling 路線需要每加一個 downstream project 手動更新 allowlist → 與 P1 metadata-derived 原則衝突且 dual SOT
2. 2026-06-06 incident 證明 allowlist cold-start window 存在 systemic risk
3. P6 bootstrap-safe 條件已解決原本傾向 sibling 的「visibility 推論誤殺新概念」風險
4. Registry rule_class 一個 class 一個 canonical executor，避免 cross-executor verdict drift

執行：Phase 0 完成後 sibling plan archive + status=superseded，本 plan 升 P1 + `required_for_completion: true`，parent meta-plan §Child plans 同步更新。

### Why Not Just Wait for Sibling Plan

| 角度 | Sibling `2026-05-31-2000` (allowlist-based) | This plan (metadata-derived) |
|---|---|---|
| Source-of-truth | `enforcement/sanitization.yaml` 中央 allowlist | 各 project metadata `private_entities[*].match_tokens` 宣告，projection 為 `derived_private_entities` (governance) + `derived_match_tokens` (execution) |
| Maintenance cost | 中央 allowlist 必須隨每個新 project 手動更新 | 只有 per-project metadata 維護（owner 在 project 端宣告 `private_entities`）；framework 端零中央維護、無 cross-project allowlist drift |
| Drift risk vs `sanitization.md` prose | 中高（dual SOT） | 低（無 dual SOT） |
| 對 today's incident (214a415) | 若 allowlist 尚未含該 label → miss | 命中（label 在 project-local 出現、在 reusable 缺席） |
| 與 sibling 共存？ | 兩個 executor 同 rule_class | 視 Q1 裁決 |

兩條路徑各自處理不同失誤模式。今日 incident 證明 allowlist 路徑有 cold-start window，本 plan 補的就是 cold-start window 的 mechanical 防線。

---

## Architecture Compatibility Preflight

| 欄位 | 內容 |
|---|---|
| Candidate files | 新建 `runtime/sanitization-patterns.yaml`（canonical Phase 2 regex）；新建 `scripts/ai-skill-cli/internal/app/sanitization_scan.go`（scanner）；新建 validator entry in `scripts/ai-skill-cli/internal/app/hooks.go` pre-commit dispatcher；新建 `enforcement/sanitization-mechanical.md`（companion，philosophy + 與 prose `sanitization.md` 邊界）；`enforcement/enforcement-registry.yaml` 更新 `rule_classes[sanitization]` executors block（或新增 second executor entry，依 Q1） |
| Source-of-truth | `enforcement/sanitization.md` 仍是 prose rule canonical；`runtime/sanitization-patterns.yaml` 是 Phase 2 regex canonical；Phase 1 forbidden tokens 由各 project metadata `private_entities[*].match_tokens` 宣告，projection 到 `runtime.db.derived_private_entities` (governance) + `runtime.db.derived_match_tokens` (execution)（framework 端無中央 allowlist） |
| Compiler / generated surfaces | `runtime/sanitization-patterns.yaml` 經 `ai-skill runtime compile` projection 到 `runtime.db.generated_surfaces`；scanner 從 runtime.db 載入 patterns |
| Layer responsibility | enforcement-mechanical 屬 enforcement layer (rule_class executor)；scanner 屬 runtime layer；pre-commit hook 整合屬 ai-skill-cli layer |
| 與現行架構衝突 | Q1 已裁決 supersede；sibling `2026-05-31-2000` 已在 `plans/archived/` 且 frontmatter `status: superseded` |
| `runtime.db` 影響 | 新增 `sanitization_patterns` projection table；無 schema migration risk（純讀） |

---

## Scope Seal

**In scope**：
- Pre-commit `git diff --cached` content scanner，scope filter = `plans/** workflow/** enforcement/** governance/** knowledge/**`
- Phase 1: metadata-derived project label scan
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

- [x] User 認可 §Relationship to sibling 的 supersede 裁決
- [x] Sibling `2026-05-31-2000` 移入 `plans/archived/` + frontmatter `status: superseded` + `superseded_by: 2026-06-06-1800-sanitization-mechanical-enforcement`
- [x] 本 plan frontmatter 升 `priority: P1` + `required_for_completion: true`
- [x] Live registry / index references synchronized; archived parent meta-plan body left historical to avoid re-opening an archived main plan with active required sub-plans
- [x] Phase 0 不啟動任何 implementation

### Phase 1 — Metadata-Derived Forbidden Token Scanner

> **Phase 1 拆分（2026-06-08 review）**：原始 Phase 1 把 4 個獨立 surface 塞在同一 phase（metadata schema / topology / projection / scanner）。Dependency 隱藏在文字下，commit 也會過大。拆成 1A/1B/1C/1D，依賴顯式：
>
> ```
> 1A (Metadata Schema) ─┐
>                        ├─→ 1C (Projection) ─→ 1D (Scanner)
> 1B (Topology Surface) ─┘
> ```
>
> 1A 與 1B 可平行 PR；1C 依賴前兩個但本身輕；1D 才是肉。每 phase 一個 focused commit。

#### Phase 1A — Project Metadata Schema

**Governance / execution layer 分層**（resolved Q8 2026-06-08）：治理層 declare **entities**（被保護的對象 — codename / customer / product 等 identity）；scanner 比對的是 **match tokens**（entity 的 alias、case variants）。Schema 直接做出分層，projection 也拆 governance table（`derived_private_entities`）與 execution table（`derived_match_tokens`）。

**Phase 1A 落地 commit**（2026-06-08）：產出 `metadata/project/ai-skill-project-schema.yaml` (canonical) + `example-ai-skill-project.yaml` + `README.md` + `migration-notes.md` + `scripts/ai-skill-cli/internal/app/project_metadata.go` (canonical parser) + `project_metadata_test.go` (17/17 pass)。**嚴格未動** `sanitization_scan.go`（legacy reader 保留，Phase 1C 才動 projection、Phase 1D 才動 scanner pointer）。詳見 [`metadata/project/migration-notes.md`](../../metadata/project/migration-notes.md)。

- [x] 定義 project metadata schema：`<PROJECT_ROOT>/.ai-skill-project.yaml`（Phase 1A implements direct metadata scan; overlay metadata 維持 future-compatible）：

  ```yaml
  project:
    id: <project-slug>
    visibility: private | public
    private_entities:
      - name: <canonical entity name>            # 治理層 ID（人讀；governance refers to this）
        kind: codename | client | product | individual | other
        match_tokens:                             # 實作層 matching surface（scanner 比對這個）
          - "ProjectFoo"
          - "project-foo"
          - "PROJECT_FOO"
        case_variants: auto                       # 或 manual: ["..."] 完全列舉
        # 未來可擴：severity / expires_at / owner / related_entity
  ```

  關鍵 design properties：
  - **Entity ≠ Tokens 1:1** — 一個 entity 可多個 alias
  - **Tokens 可 collision** — 兩個 entity 的 match_tokens 重疊 → governance lint 可偵測
  - **Future-proof** — 加 entity-level metadata 不破 schema

- [x] schema validation test：parser 接受合法 schema、拒絕缺 `kind` / 缺 `match_tokens` / `case_variants` 既非 `auto` 也非 list 等變形
- [x] sample `.ai-skill-project.yaml` documentation（README in plan body example）

#### Phase 1B — Repository Topology Surface

**設計決策**（resolved 2026-06-08）：採 `schema_version: 2`，每 subtree 必須有 `owner` + `purpose`。**移除** flat `expected_consumers` list（容易 stale → 又一個 stale-reference 失效模式）；改寫死 `consumer_tracking.strategy: code_reference` 把治理決策 freeze，避免欄位 6 個月後復活。

**Phase 1B 落地 commit**（2026-06-09）：產出 `runtime/repository-topology-migration.md`（schema spec + migration trajectory）+ `scripts/ai-skill-cli/internal/app/repository_topology.go` (canonical loader supports v1 read / v2 read / v2 write) + `repository_topology_test.go` (14/14 pass, including live `runtime/repository-topology.yaml` v1 parses regression guard)。**嚴格未動** `runtime/repository-topology.yaml`（仍 v1 in production）、`runtime_compiler.go` line 339 projection rule、`sanitization_scan.go::repositoryTopologyRow`。YAML v1→v2 在線升級由 Phase 1C 跟 projection rewrite 同 commit 做（rationale: compiler line 339 hard-codes v1 field names；單獨升 YAML 會破 projection）。詳見 [`runtime/repository-topology-migration.md`](../../runtime/repository-topology-migration.md)。

- [x] 新建 `runtime/repository-topology.yaml`（canonical cross-subsystem topology surface，**不掛 enforcement-registry**）：

  ```yaml
  # runtime/repository-topology.yaml
  schema_version: 2
  consumer_tracking:
    strategy: code_reference   # consumer list NOT manually maintained; derive via
                                # grep of code that loads this file. See
                                # governance/lifecycle/topology-consumers.md
                                # (TBD; created on first dispute).
    rationale: |
      Manual consumer lists go stale (new subsystem reads topology but no one
      updates the list). Code-reference derivation is the durable governance
      decision. This block is intentionally permanent — do not re-add a
      manual `expected_consumers:` field.

  subtrees:
    - path: plans/
      shared_layer: true
      owner: framework-maintainer
      purpose: "Plan tracking; referenced by enforcement-registry child_plan"
    - path: workflow/
      shared_layer: true
      owner: framework-maintainer
      purpose: "Cross-skill workflow contracts and execution flows"
    - path: enforcement/
      shared_layer: true
      owner: framework-maintainer
      purpose: "Reusable enforcement rules and failure patterns"
    - path: governance/
      shared_layer: true
      owner: framework-maintainer
      purpose: "Lifecycle / promotion / ADR governance surfaces"
    - path: knowledge/
      shared_layer: true
      owner: knowledge-maintainer
      purpose: "Atom summaries, indexes, graphs, routing registry"
    - path: runtime/
      shared_layer: true
      owner: framework-maintainer
      purpose: "Compiled / projected runtime surfaces consumed by Go executor"
    - path: validation/
      shared_layer: true
      owner: framework-maintainer
      purpose: "Reusable validation scenarios and fixtures"
    - path: scripts/
      shared_layer: false
      owner: tooling-maintainer
      purpose: "CLI / runtime implementation; not consumed as reusable knowledge"
    - path: .agent-goals/
      shared_layer: false
      owner: project-local
      purpose: "Per-project ephemeral goal ledger; excluded from git"
  ```

- [x] **Backward-compat**：v1 schema (`shared_layer_classification:` 配 `subtree: / shared:`) 仍可解析作為 read-only fallback；compiler 寫入時一律 v2。Transition window 一個 minor release（之後可移除 v1 reader）。
- [x] Topology projection target：`runtime.db.repository_topology`（query via `ai-skill runtime topology` CLI）
- [x] Topology loader unit test：v1 read、v2 read、v2 write round-trip、缺 `owner`/`purpose` reject

#### Phase 1C — Runtime Compile Projection

**Phase 1C 拆分為 1C₁ + 1C₂**（mirroring 1A/1B discipline success）：

- **Phase 1C₁** — Topology projection migration（landed 2026-06-09，commit landing this section）
  - `runtime/repository-topology.yaml` 升級 v1 → v2 in-place（11 subtrees with owner+purpose；`expected_consumers` 移除；`consumer_tracking.strategy: code_reference` frozen governance block）
  - `runtime_compiler.go` line 339 tuple 移除；新 `repository_topology_compile.go::compileRepositoryTopology` 接 Phase 1B `LoadRepositoryTopology`
  - 投影 JSON 內容 dual-shape — 同時帶 v1 keys (`subtree`, `shared`) + v2 keys (`path`, `shared_layer`, `owner`, `purpose`)；legacy `sanitization_scan.go::repositoryTopologyRow` 不破
  - 2 compiler tests pass：`TestCompileRepositoryTopology_WritesBackwardCompatJSON` + `TestCompileRepositoryTopology_LegacyScannerCompat`
  - `TestLoadRepositoryTopology_LiveFileParses` 更新斷言 v2 + 驗 owner/purpose 完整
  - 完整 internal/app 測試套 (110s) 全 pass，無 regression
  - 詳見 [`runtime/repository-topology-migration.md`](../../runtime/repository-topology-migration.md)

- **Phase 1C₂** — Project metadata projection（additive，landed 2026-06-10）
  - 新 `derived_private_entities`（governance：name / kind / owning_project / placeholder）+ `derived_match_tokens`（execution：token / canonical_token / entity_name / kind / owning_project / placeholder）兩張表，**additive** 不取代既有 `derived_forbidden_tokens`（Phase 1D 才退役 legacy）
  - 新 `scripts/ai-skill-cli/internal/app/project_metadata_compile.go`：`compileProjectMetadataDerived`（walk + parse via `LoadProjectMetadata`）+ pure `projectMetadataDerivedRows` + `expandEntityMatchTokens`
  - case-variant expansion 在 projection 階段（`auto` → `sanitizationTokenVariants`；`explicit` → verbatim 並抑制 auto）；parser 維持純讀
  - cross-entity token collision 保留 entity identity（PK 含 entity_name），供 Phase 1D finding 報 entity 而非僅 token
  - transition tolerance：新-schema validation 失敗的 project file 跳過（legacy 投影仍覆蓋），I/O 錯誤照常 propagate
  - source-file 註冊交給 legacy 投影（避免 `runtime_source_files` source_path PK 衝突）
  - `project_metadata_compile_test.go`（5/5 pass：auto expansion / explicit suppresses auto / cross-entity collision / public skip / DB integration + bootstrap-safety）
  - 嚴格未動 `sanitization_scan.go`；live `ai-skill runtime compile` 驗證新表存在且 0 row（bootstrap-safe — 尚無 project 以新 schema 宣告 `private_entities`）

- [x] **1C₁** topology projection migration（landed this commit）
- [x] **1C₂** project metadata projection（additive，landed 2026-06-10）
- [x] Phase 1B/1A discipline 保持：1C₁/1C₂ 不動 `sanitization_scan.go`，legacy reader/scanner 完全靠 backward-compat JSON 持續運作
- [x] Compiler unit test（topology + 1C₂ project metadata projection 各自 unit test 全 pass）

#### Phase 1D — Scanner Implementation

- [x] 實作 scanner core：staged file 路徑落在 shared-layer subtree（query `runtime.db.repository_topology`）→ 比對 `runtime.db.derived_match_tokens` → emit finding
- [x] **Bootstrap-safe guard**：scanner 不做「reusable layer 是否包含 token」inference；新 framework concept 只要不在任何 project `private_entities` 內，自動通過
- [x] False-positive guardrail：`enforcement/sanitization.md` 自身的 example 段落 self-reference exception（Phase 1 以 metadata-derived scope 達成：未被 project metadata 宣告的 synthetic examples 不會觸發；真實 private token 仍會被擋）
- [x] Unit tests：
  - fail case：synthetic reconstruction of private-token leak in shared-layer plan content
  - pass case：同一 private token 允許在 project-local `.agent-goals/` content（topology shared=false）
  - bootstrap-safety case：首次提交一個全新 framework concept token，確認 0 finding
  - entity vs token granularity case：兩個 entity 的 match_tokens 部分重疊 → finding 報 entity name（不只 token），便於 governance debug

#### Phase 1D — Shape-Aware Skip Remediation（bound 2026-06-10，源自 1C₂ self-review）

> **背景（1C₂ self-review Finding 1）**：`compileProjectMetadataDerived`（`project_metadata_compile.go`）目前對 new-schema validation 失敗的 `.ai-skill-project.yaml` 採 **silent `continue`**，rationale 是「transition 期間 legacy flat-shape 檔案本來就會 validation fail，legacy `derived_forbidden_tokens` 投影仍覆蓋」。此 tolerance **僅在 legacy 表還活著時成立**。
>
> **為何綁 1D 而非更早**：1C₂ 階段 silent-skip 完全無害（scanner 仍吃 legacy，0 個 `.ai-skill-project.yaml` 存在）。真正咬人的時間點正是 **1D**——scanner pointer 改指 `derived_match_tokens` **且 legacy 表退役**的那一刻：「被 legacy 覆蓋」的豁免理由蒸發，此後任何 skip = 該 project 零保護、且靜默無聲，恰好是本 plan 要消滅的「rule 無 executor / 自律靜默失效」class。**硬截止點：1D 完成前必須處置；不得拖到 Phase 4 gate 打開（那時壞檔直接靜默放行）。**
>
> **正確修法（非單純 stderr warning）**：純 warning 無法區分「合法舊格式」vs「新格式打錯字」。1D 必須做 **shape-aware** 分流：
> - 偵測檔案使用的是 new-schema 鍵（`private_entities` 物件陣列）還是 legacy flat 形（`private_tokens` / string-list `private_entities`）
> - **legacy flat 形** → 容忍跳過（與當前行為一致）
> - **new-schema 但 validation fail** → **hard error**（compile 失敗），不得 silent `continue`

- [x] **1D-remediation**：`compileProjectMetadataDerived` 改 shape-aware tolerance（legacy-flat 容忍／new-schema-malformed hard-fail），與 scanner pointer 遷移、legacy `derived_forbidden_tokens` 退役同批
- [x] Unit test：new-schema-malformed file（缺 `kind`）→ compile hard error；legacy-flat file → 容忍跳過 0 error
- [x] 移除 `project_metadata_compile.go` 內 `IsValidationError(err) { continue }` 的無條件 silent-skip，改判 shape 後分流

**Phase 1D 落地 commit**（2026-06-10）：scanner pointer migration + legacy retirement + shape-aware remediation 同批完成：

- **Scanner migration**：`sanitization_scan.go` 的 `loadDerivedForbiddenTokens` → `loadDerivedMatchTokens`（query `derived_match_tokens`，含 `entity_name`/`kind`）；`derivedForbiddenToken` struct → `derivedMatchToken`；finding 字串改為 `contains "<token>" (entity "<name>" from <project>); use <placeholder>`，達成 entity-granularity 要求。
- **Legacy retirement**：移除 `compileDerivedForbiddenTokens` / `readProjectMetadata` / legacy flat structs / `insertProjectMetadataSourceFile`；`runtime_compiler.go` 移除 `derived_forbidden_tokens` CREATE TABLE + projection call；`runtime.go::nativeRuntimeRequiredTables` 改列 `derived_private_entities` + `derived_match_tokens`。`compileProjectMetadataDerived` 成為唯一 project-metadata 投影，並接回 `runtime_source_files` 註冊（legacy 退役後無 PK 衝突）。
- **Shape-aware remediation**：新增 `classifyProjectMetadataShape`（probe `yaml.Node`：`private_tokens` 或 scalar `private_entities` → legacy-flat；object `private_entities` → new-schema）。legacy-flat → stderr warning + 跳過；new-schema 但 validation fail → **hard error**（關閉 1C₂ self-review Finding 1 的 silent-skip 缺口）。
- **驗證**：`go test ./...`（全模組）+ `go vet` clean；新增 `TestCompileProjectMetadataDerived_LegacyFlatShapeTolerated` + `_NewSchemaMalformedHardFails` + scanner entity-name 斷言；`ai-skill runtime compile + refresh + validate` 全 pass，live `runtime.db` 確認 `derived_forbidden_tokens` 已消失、兩張新表存在。
- **未動**：`enforcement-registry.yaml` `coverage` 仍 `pending_implementation`（promote 到 `mechanical` + 加 blocking per-commit obligation 屬 Phase 4，本批不碰）。

#### Phase 1D review — Finding A（hard-fail blast radius；綁 Failure Authority invariant）

> **觀察（2026-06-10 Phase 1D self-review）**：shape-aware hard-fail 正確關閉了 silent-skip，但 `discoverProjectMetadataFiles` 掃全 repo（只跳 `.git`/`node_modules`/`vendor`），因此一個 malformed `.ai-skill-project.yaml` **不論落在哪一層都會 hard-fail compile** —— 包括 git-ignored、project-local 的 `.agent-goals/…`。等於把「silent under-protection」換成「對 ephemeral local 檔案 over-blocking」。
>
> **不在本 plan 直接實作修正**。這不是 sanitization 個案，而是一個 cross-cutting governance 問題「**誰有資格阻塞 compile？**」（與 Discovery fail-open、Runtime Index source-row scope 同類）。已抽為 incubator 第三 family：[`governance/lifecycle/governance-pattern-library-draft.md`](../../governance/lifecycle/governance-pattern-library-draft.md) §"Parallel observation: Failure Authority family"。
>
> **Invariant（observation-stage）**：只有 compile-authoritative source（shared-layer / tracked / runtime-index `sources` row）可阻塞 compile；non-authoritative source（`shared_layer:false` / `owner:project-local` / untracked）只能 warn。分類靠 topology v2 —— 若 `shared_layer:true|false` 最終都 compile fail，topology 的 path-classification 治理價值在最該發揮的時刻被削弱。
>
> **落地方向（dependency inversion + Standing，2026-06-10 確認）**：核心原則是 **Standing**（Validity ≠ Authority）——「檔案 100% invalid 不代表它有資格阻塞 compile」。不從 Finding A 推導 classifier，而從 invariant 推導，且中間多一層 **Authority Classification Contract**（subject-based，非 path-based）。Build order = (1) family observation ✅ → (2) **Authority Classification Contract** ✅（`ClassifyFailureAuthority(subject)`，subject kind = discovery-provider / runtime-index-row / metadata-file / generated-surface，非 `isCompileAuthoritative(path)`）→ (3) classifier 實作（contract 第一個 implementation）→ (4) Finding A 作為 **Executor #1**（metadata-file subject caller）。如此 Discovery / Runtime Index / 未來 surface 共用同一 contract。**不採** scanner 內 `.agent-goals/` 特例 hack。詳見 [`governance/lifecycle/governance-pattern-library-draft.md`](../../governance/lifecycle/governance-pattern-library-draft.md) §"Authority Classification Contract" + §"Dependency inversion"。

- [x] **Finding A 落地（Failure Authority Executor #1，2026-06-10）**：`compileProjectMetadataDerived` 改 call `ClassifyFailureAuthority(metadata-file subject)` —— 由 `repository-topology.yaml`（直接 load，因 projection 早於 topology compile）longest-prefix 解析 `shared_layer`/`owner`，authoritative → hard-fail、non-authoritative（`shared_layer:false` / `owner:project-local` / `.agent-goals/`）→ warn+skip。topology miss → SharedUnknown → authoritative（fail-safe）。新增 `SubtreeForPath` + `metadataFileAuthoritySubject` + 2 compile authority-scoping tests + SubtreeForPath test。Family sample #3 ❌→✅，N=4。

### Phase 2 — Generic Regex Patterns (inherited from superseded plan)

- [x] 新建 `runtime/sanitization-patterns.yaml`（canonical）+ companion section in `enforcement/sanitization-mechanical.md`
- [x] Pattern set：
  - Email：`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
  - Phone：國際 + 區域格式（保留 conservative，避免誤抓 plan timestamp）
  - OS absolute path：`/Users/[^/\s]+` / `/home/[^/\s]+` / `C:\\Users\\[^\\\s]+` / `\\\\\\\\?\\\\` UNC
  - Credential pattern：高熵字串 + 已知 prefix（`sk-` / `ghp_` / `xoxb-` 等）
- [x] Pattern allowlist override：placeholder forms `<USER>` / `<PROJECT_ROOT>` / `<AI_SKILL_REPO>` 必須被 pattern 排除
- [x] `ai-skill runtime compile` projection → `runtime.db.sanitization_patterns`
- [x] Unit tests：每 pattern ≥ 1 fail + 1 pass + 1 placeholder pass

### Phase 2.5 — Incident-Score Heuristic (warn only, inherited from superseded plan)

補強 Phase 1 metadata-derived 抓不到的「無 private token 但明顯是 project incident」案例（連續的 domain-specific noun cluster、quoted user text、具體 dated filename 引用等）。**warn-only，不 block**，理由：heuristic 性質與 P4 attestation-prohibited 的 deterministic 精神有張力，降為 advisory 守住 deterministic-block 邊界。

- [x] 將 sibling plan §Phase 1 `incident_score` weighted-signal schema 整段搬入 `runtime/sanitization-patterns.yaml` `incident_score:` 區段（filename_pattern / quoted_user_text / artifact_string / domain_noun_cluster，weights 5 / 5 / 3 / 1，warn_if_total_score_gte: 7）
- [x] Scanner 在 Phase 1+2 全部 pass 時才執行 incident-score；任一 phase block → 跳過（避免 noise）
- [x] Output 為 `warning` 等級，commit 不 block；finding 寫入 stderr 供 agent / user 評估
- [x] Unit tests：sibling plan §incident_score `examples` 三個 case（leak 命中、registry 內合法描述放行、archived/feedback 路徑跳過）必須全部 pass
- [x] Acceptance：incident-score 不得對 Phase 1 已 derived_forbidden 的 token 重複 finding（避免雙報）

### Phase 3 — LLM Review Hook (Conditional, Deferred)

- [ ] Trigger condition：Phase 1+2 三個月真實使用 miss rate 量測；若 > X% (X 待 ADR 定) → 啟動 Phase 3 ADR
- [ ] Phase 3 自己一個 sub-plan（本 plan 不展開）

### Phase 4 — Registry & Bootstrap Integration

> **Phase 4 拆 step**（2026-06-11）：step 1（可逆，已完成）= failure-pattern doc + 3 validation scenarios；step 2/3（promotion，待執行）= registry `coverage→mechanical` + 不可 opt-out blocking obligation + hooks 註冊。Step 1 先行因 registry Status Transition Matrix 要求 `scenario_exists`/`regression_exists` 才能合法升 mechanical。

- [x] 新建 `enforcement/failure-patterns/sanitization-leak-on-canonical-write.md`（empirical evidence: 214a415 + sibling plan v1-v4 累積 leak；含 executor binding + Failure Authority 關係）
- [ ] 更新 `enforcement/enforcement-registry.yaml` `rule_classes[sanitization]`：
  - 新增 executor entry: `file: scripts/ai-skill-cli/internal/app/sanitization_scan.go`, `symbol: scanStagedContentForSanitization`, `hook_phase: pre-commit`, `block_or_warn: block`
  - 若 Q1 裁決 supersede → 移除舊 allowlist executor entry
  - `coverage` 視 promotion criteria 是否成立決定（pending_implementation → mechanical）
- [ ] 更新 `runtime/core-bootstrap.yaml` `per_commit_obligations` 新增 `obligation.commit.sanitization_visibility_scan`（**無 opt_out_marker**，per P4 attestation-prohibited + Q5 reject）
- [ ] `commitMsgValidatorRegistry` 在 hooks.go 註冊新 obligation
- [ ] Companion `enforcement/sanitization-mechanical.md` cross-link to prose `sanitization.md`
- [ ] Re-dry-run `ai-skill runtime compile` + enforcement coverage report
- [x] Validation scenarios（+ `validation/README.md` 索引）：
  - `validation/scenarios/runtime/sanitization-metadata-derived-pass-v1.yaml`
  - `validation/scenarios/runtime/sanitization-metadata-derived-fail-v1.yaml`（commit 214a415 reconstruction）
  - `validation/scenarios/runtime/sanitization-placeholder-allowed-v1.yaml`
- [ ] Owner-grouped commit + push + readback

---

## Open Questions

| # | Question | Pri | Status | 處置 |
|---|---|---|---|---|
| Q1 | ~~與 sibling 關係~~ | — | **resolved 2026-06-06 / verified 2026-06-08** | supersede。Sibling already archived with `status: superseded`; this plan is the active successor. |
| Q2 | Visibility metadata 載入機制：scan `<PROJECT_ROOT>/.ai-skill-project.yaml`，還是嵌入既有 routing/overlay metadata？ | P1 | open | Phase 1 ADR；不影響本 plan 接受度。 |
| Q3 | Phase 3 LLM review 的 miss-rate threshold 是多少才 trigger？ | P2 | open | Phase 3 sub-plan 自己決定。 |
| Q4 | ~~Shared-layer classification 進 enforcement-registry？~~ | — | **resolved 2026-06-06: 否** | 改建獨立 `runtime/repository-topology.yaml` cross-subsystem surface。Rationale: topology 將被 sanitization / workflow activation / governance lint / dependency reading 共用，掛 enforcement-registry 會讓該 registry 變超級桶（同 parent meta-plan R6 警告的反模式）。 |
| Q5 | ~~`[skip-sanitization-scan]` opt-out marker？~~ | — | **resolved 2026-06-06: reject** | 與 P4 attestation-prohibited 衝突。若未來需要 emergency override，須另開 ADR 設計 admin-override surface（含 owner / reason / time-boxed expiry），不採 commit-message marker。 |
| Q6 | False-positive 處理機制：suggested_placeholder 是否自動 patch staged content？ | P2 | open | v1 不自動 patch；自動 patch 列 v2 評估。 |
| Q7 | P4 attestation-prohibited 原則升級為 cross-cutting governance：何時抽出獨立 enforcement rule (`enforcement/verification-not-declaration.md`)，覆蓋 Dependency Read / Test Executed / Coverage Reviewed 等？ | P2 | open | 本 plan 內留 inline；累積 ≥3 個受益 obligation 後 promote 為獨立 rule，列入 parent meta-plan tracking。 |
| Q8 | Metadata 欄位是否改名 `private_entities` 而非 `private_tokens`？ | P1 (架構級) | **resolved 2026-06-08** | 不是 naming 問題，是 governance ≠ execution layer boundary 問題。決議：(a) schema 用 `private_entities:`（治理層 — entity 是 governance refers to 的對象，有 name / kind 身份）；(b) entity 含 `match_tokens:`（執行層 — scanner 比對的字串 surface）；(c) projection 也拆兩張 table — `derived_private_entities`（治理 query）+ `derived_match_tokens`（scanner direct query）。完整 schema 見 §Phase 1A。 |

---

## Acceptance

- Phase 1 blocks staged shared-layer content containing project-declared private tokens from `.ai-skill-project.yaml`.
- Phase 1 remains bootstrap-safe: no project metadata means no private-token findings.
- Phase 2 regex and Phase 2.5 incident-score heuristics remain deferred until their checklist items are explicitly executed.
- `rule_classes[sanitization]` has one active implementation path: this metadata-derived plan.

---

## Validation

| 欄位 | 內容 |
|---|---|
| Trigger | 2026-06-06 session 確認本 plan 已 commit + push（per `enforcement/dependency-reading.md` writeback transaction） |
| Required set | `enforcement/sanitization.md` / `enforcement/reusable-guidance-boundary.md` / `plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md` / `knowledge/runtime/routing-registry.yaml` (sanitization slots) / `scripts/ai-skill-cli/internal/app/hooks.go` (dispatcher) |
| Read | ✅ 全部 |
| Not applicable | — |
| Deferred / blocked | Phase 1 implementation pending; Phase 0 does not start scanner implementation |
| Validation | Phase 0 preflight verified sibling supersede, active successor status, live index/failure-pattern references, sanitization source rules, and registry child_plan pointer |

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
| 任何 allowlist-style `private_project_names:` / `allowed_tokens:` 設計 | 被 metadata-derived `private_entities[*].match_tokens` 取代 |
| `[skip-sanitization]` opt-out marker | P4 + Q5 reject |
| `canonical_paths` / `not_canonical` 硬列清單 | 由 `shared_layer_classification` table 取代 |
| PreToolUse Write/Edit warning hook | commit-diff-only 立場下不需要（避免 Edit `old_string` 誤判） |

## Handoff Notes

- 本 plan 與 sibling `2026-05-31-2000` 同 coverage class；Q1 已裁決 supersede，後續 implementation 只走本 plan，避免 registry rule_class 出現 conflicting executor
- Empirical trigger commit 214a415 + 修正 commit 728282c 是 Phase 1 unit test 的 ground truth dataset
- Phase 4 Registry integration 必須與 sibling 協調 — `enforcement-registry.yaml` 同檔，避免 commit conflict
- 若 user 將本 plan 升 P1，需同步 archive 或降級 sibling，並更新 parent meta-plan §Child plans 列表
