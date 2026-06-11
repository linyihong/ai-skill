---
id: 2026-06-11-1100-plan-archival-link-integrity
plan_kind: main
status: ready
owner: linyihong
created: 2026-06-11
updated: 2026-06-11
priority: P2
required_for_completion: false
---

# Plan Archival Link Integrity

**Status**: `ready`
Owner: framework maintainer (linyihong)
**世代**：Gen 3 Runtime Hardening → Reference Integrity → Plan Archival Link Integrity
**建立日期**：2026-06-11

> 把「plan archive 時 relative-link 斷裂」這個 2026-06-11 親身踩到的失誤模式機械化，作為 **Reference Integrity Engine** 家族的新成員。

## Decision Rationale

### 家族定位

本 check 不屬於 Plan Governance / Archival Audit 的擴張，而是 Reference Integrity 家族的新成員：

```
Gen 3 Runtime Hardening
└─ Reference Integrity
   ├─ runtime_index_freshness        （source ↔ index drift）
   ├─ topology references            （layer 移動後引用 drift）
   ├─ parent-child references        （plan tree frontmatter）
   └─ plan_archival_link_integrity   ← 本 plan
```

`validatePlanArchivalAudit` 管的是 **Workflow Completeness**（archive 流程跑完沒：checkbox 勾完、audit 完成）；link 是否解析得到屬於 **Reference Integrity**。混在同一 validator 會讓 archival audit 變垃圾桶。

### Problem & Why Now

2026-06-11 archive `2026-06-06-1800-sanitization-mechanical-enforcement`（`active/ → archived/`）時，手動發現 **兩側 relative-link 都斷**：

1. **被移動檔案自身的 outbound link**：原本假設自己在 `plans/active/`，move 到 `plans/archived/` 後 `../archived/`、same-dir、`../active/` 全部要重算。
2. **Repo 內其他檔案的 inbound link**：8 個 active 檔案（README、2 runtime yaml `source_plan`、2 failure-pattern、2 metadata、topology-migration）指向 `plans/active/<id>`，move 後變 stale。

現有 `validatePlanArchivalAudit` 只檢查 unchecked `- [ ]`，**完全不看 link**。全靠手動 grep + 逐一修，正是「rule 無 executor」風險。

### The failure this prevents

「archive 一個 plan → 連結默默斷掉 → 半年後有人點到 404 / 工具 resolve 失敗」。本 check 在 archive commit 當下就 surface，並提供 `suggested_replacement` 讓修復成本接近零。

## Scope

**In scope**
- 偵測 staged 內 plan 檔案的 `active/ ↔ archived/` rename，透過 `git diff --cached --find-renames -M90`（high similarity threshold；archive 通常 ≥90% 內容不變）。
- **解析失敗檢查**：repo 內所有 staged + 受影響檔案的 markdown link `](relpath)`，從 **link 所在檔案的當前位置** resolve；target 不存在 → finding。
- **Stale textual path mention**：bare 文字提及舊路徑（非 markdown link 語法）→ warning。
- **suggested_replacement payload**：archive 場景 old→new 雙端已知，每個 finding 附 `suggested_replacement` 欄位（不自動修，但下游 IDE / auto-fix 可吃）。

**Out of scope**
- 自動修連結（finding 帶 suggestion 即可；auto-fix 列未來）。
- 非 plan 檔案的一般 link-rot（更大題目，另議；本 validator 只在 plan archive 事件觸發）。
- 絕對 URL / 跨 repo link。

## Phase Plan

### Phase 0 — Design decisions（已 resolved）

- [x] **D1 = 新 validator**：建立 `validatePlanArchivalLinkIntegrity` + 新 obligation `obligation.commit.plan_archival_link_integrity` + `enforcement-registry.yaml` 入 `reference_integrity` rule_class。原因：archival audit 是 Workflow Completeness，link 是 Reference Integrity；分類乾淨優於合併。
- [x] **D2 = severity 軸改為 resolvability，不分 direction**：
  - markdown link syntax `](path)` 且解析失敗 → **block**（無論 inbound/outbound；客觀錯誤）
  - bare textual path 提及（prose / frontmatter 註解 / 歷史描述）→ **warning**（避免誤殺 provenance 文字）
  - opt-out: `[skip-plan-archival-link-integrity]`
- [x] **D4 = textual finding 分 category + opt-in provenance marker**：避免「半年後 warning 一直不修，其實是故意保留歷史」的雜訊。
  - default category = `stale_textual_reference`（warning）
  - 若同行或上一行有 `<!-- archival-provenance -->` marker → category = `historical_provenance_reference`，severity = `info`（不進 warning 列表）
  - 機械區分意圖比 NLP 猜 phrasing 可靠；責任落在寫 provenance 的人
  - Finding payload 帶 `category` 欄位，下游工具可分流
- [x] **D3 = `git diff --cached --find-renames -M90`**：吃 Git 算好的 rename intent，比自己重建對照少 edge case；threshold 90% 對 archive 場景合適（pure move，極少改動）。

### Parser Strategy（Design Note）

本 validator **刻意**只支援 framework plan 文件實際使用的 markdown link 子集，不引入 markdown AST 依賴（goldmark 等）。理由：domain 是 **Plan Archive Event** + `plans/` 子樹，不是 universal markdown lint。引入完整 AST 對不存在的構造付成本。

實作為 `extractMarkdownLinks()` 小型 state machine（char-by-char 掃描），不是散落的 regex，回傳 `Link{Target, Line, Column}` struct。

**Supported（會被解析並驗證）**
- inline link `[text](path)`
- inline link with title `[text](path "title")`
- escaped parens in path `[text](../a\(b\).md)`
- code-fence exclusion（``` 與 ~~~ block 內整段跳過）
- 相對路徑（絕對 URL / `mailto:` / `#anchor` 過濾）

**Not supported（**ignored**，不部分解析）**
- reference-style link `[text][ref]` + `[ref]: path`
- HTML anchor `<a href="...">`
- autolink `<https://...>`

**Contract**: 不支援構造 **整段忽略**，不做部分解析。誤解析比漏報危險。

### Phase 1 — Implementation

- [ ] `scripts/ai-skill-cli/internal/app/markdown_links.go`（新檔）：實作 `extractMarkdownLinks(content []byte) []Link` bounded parser（state machine，40-80 行）
- [ ] `scripts/ai-skill-cli/internal/app/hooks.go` 新增 `validatePlanArchivalLinkIntegrity`
- [ ] 偵測 staged plan rename（`active/ ↔ archived/`）：跑 `git diff --cached --find-renames -M90 --name-status` 取所有 `R*` 條目，過濾 plan 路徑
- [ ] **建立整批 rename map（必須在掃描前完成）**：multi-archive in same commit 時，A、B 同時 archive 且互相引用，每個檔的 resolve 都要看完整 rename map，不能逐檔處理
- [ ] **Markdown link parsing**：使用 `extractMarkdownLinks()`（bounded parser，**非** regex、**非** markdown AST），取得 `Link{Target, Line, Column}`；避免 prose 中的路徑字串被誤判為 link
- [ ] **解析**：對每個 link，從 link 所在檔案的 **新位置**（若該檔本身被 rename）或 **當前位置** resolve 相對路徑；target 不存在 → finding
- [ ] **Bare textual path scan**：對被 rename 檔案的舊路徑（`plans/active/<id>`）做 plain-text 搜尋，命中且不在 markdown link node 內 → finding。檢查命中行（與上一行）是否含 `<!-- archival-provenance -->`：有 → category `historical_provenance_reference` / severity `info`；無 → category `stale_textual_reference` / severity `warning`
- [ ] **suggested_replacement**：finding payload 帶 `{old_path, new_path, suggested_replacement, category}`，old/new 從 rename map 反查

### Phase 2 — Tests

- [ ] fail/markdown link broken：archive A，A 內含 `[parent](../active/sibling.md)`（move 後相對路徑錯）→ block
- [ ] fail/inbound markdown link broken：另一 active 檔含 `[source](plans/active/<moved-id>.md)` → block
- [ ] warn/stale textual mention：另一檔 prose 寫 `Archived from plans/active/<moved-id>.md`（非 link 語法，無 provenance marker）→ warning，category `stale_textual_reference`
- [ ] info/historical provenance：同上 prose，但同行/上一行有 `<!-- archival-provenance -->` → severity `info`，category `historical_provenance_reference`，不進 warning 列表
- [ ] pass/clean archive：move 且所有 inbound/outbound markdown link 都已 retarget → 0 finding
- [ ] pass/bare id provenance：純歷史 prose 提及 bare id（無路徑）→ 不誤報
- [ ] pass/escaped parens in path：plan 含 `[text](../a\(b\).md)`，target 存在 → bounded parser 正確解析跳脫括號，0 finding（驗證 state machine 相對 regex 的主要價值）
- [ ] **fail/multi-archive cross-reference**：同一 commit 內 A、B 都 archive，A 內有 `[B](../active/B.md)` 但未更新為 `../archived/B.md` → block（驗證 rename map 整批建立邏輯）
- [ ] pass/multi-archive cross-reference resolved：同上但 A 已更新為 `B.md`（same-dir archived）→ 0 finding

### Phase 3 — Registry & Bootstrap Integration

- [ ] `enforcement/enforcement-registry.yaml`：新增 rule_class entry（`reference_integrity` 家族）+ executor 指向 validator + coverage = `mechanical`
- [ ] `runtime/core-bootstrap.yaml` §per_commit_obligations：新增 `obligation.commit.plan_archival_link_integrity` + opt_out_marker + contract_source 指本 plan
- [ ] `enforcement/failure-patterns/plan-archival-link-drift.md`（empirical: 2026-06-11 sanitization archive commit 3f7c4b4 手動修 8 inbound + 3 outbound link）
- [ ] validation scenarios（scenarios/ 對應 test list）
- [ ] `ai-skill enforcement coverage` 確認 bucket 從 `none` → `mechanical`
- [ ] commit / push / readback

## Acceptance

- Archiving a plan with a broken markdown link (任一方向) emits a **block** finding with `suggested_replacement` payload.
- Archiving a plan with a stale **textual** path mention (non-link prose) emits a **warning** finding.
- Clean archive (all markdown links retargeted) passes with zero findings.
- Bare-id provenance mentions（無路徑語法）do not false-positive.
- Multi-archive in same commit: cross-references between simultaneously-archived plans are correctly resolved against the batch rename map.
- Unsupported markdown constructs (reference-style links, HTML anchors, autolinks) are **ignored, not partially interpreted**. Validator intentionally supports only the markdown subset used by framework plan documents.

## Future Extensibility

本 plan **不**抽象共用元件。先實作具體 executor，跑過一輪 telemetry 後再決定是否抽 `pkg/reference/` 共用層（`BuildRenameMap()` / `ResolveReference()` / `SuggestReplacement()`），給其他 rename-shaped 場景重用：

```
未來潛在共用 engine（不進本 plan）
└─ ReferenceRewriteEngine
   ├─ PlanArchivalLinkIntegrity       ← 本 plan，先具體後抽象
   ├─ TopologyMigrationIntegrity      (potential)
   └─ MetadataRelocationIntegrity     (potential)
```

過早抽象成本高於收益；以本 executor 的實際 finding 分佈為證據再決定。

## Validation

| 欄位 | 內容 |
|---|---|
| Trigger | 2026-06-11 sanitization plan archive 親身踩到 inbound+outbound link 斷裂 |
| Empirical evidence | commit 3f7c4b4（手動修 8 inbound + 3 outbound link） |
| Required set | `scripts/ai-skill-cli/internal/app/hooks.go`（新 validator）/ `scripts/ai-skill-cli/internal/app/markdown_links.go`（custom bounded link parser）/ `enforcement/enforcement-registry.yaml` / `runtime/core-bootstrap.yaml` §per_commit_obligations |
| Deferred | auto-fix；非 plan link-rot；跨 repo link |
