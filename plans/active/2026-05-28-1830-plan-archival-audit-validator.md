# Plan Archival Audit Validator

**Status**: `in-progress`
**建立日期**：2026-05-28
**最後更新**：2026-05-29
**Parent**: [`plans/active/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md`](2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md) §Phase 7 follow-up

---

## Decision Rationale

### Problem & Why Now

Parent plan §Phase 7 加了 manual "Plan Completion Audit" 步驟（`grep -nE '^- \[ \]'` 必須回空或明文交代）。Manual 步驟容易漏掉 — parent plan 自己在 Phase 3→4 過渡時就漏掉一個 deferred scenario，靠 user 點出才補救（commit `37065ea`）。Manual 規則需要 mechanical enforcement 才不會 drift。

本 plan 補上第 19 個 commit-msg validator `validatePlanArchivalAudit`，把 Phase 7 規則機械化：當 commit 把 `plans/active/<plan>.md` 移到 `plans/archived/<plan>.md` 時，scan 該檔內容；若仍有 `- [ ]` line 且 commit body 沒明文交代，block。

### Decision

新增第 19 個 commit-msg validator：

- **觸發**：staged diff 含 `plans/active/<plan>.md` 刪除 + `plans/archived/<plan>.md` 新增（或 rename detection）
- **檢查**：scan archived 版本文字，每個 `- [ ]` line 必須在 commit body 找到對應交代（regex 或關鍵字 match：deferred / non-goal / scope reduced / handover / 延後 / 拆分）
- **block default**；opt-out `[skip-plan-archival-audit]` 給特殊 cases（例如 emergency archive）
- 註冊 `obligation.commit.plan_archival_audit` 在 `runtime/core-bootstrap.yaml`
- 加入 `runtime/cli-modification-policy.yaml` 作 `gate.plan_archival_audit_required`

### Alternatives Considered

- **A. 不做，靠 manual rule**：reject — parent plan 自己已證明 manual 不夠。
- **B. 改寫 `validatePlanCheckboxSync` 讓它在 archival 時更嚴**：reject — 兩個 trigger 完全不同（一個是 reference + code work，另一個是 file move），混在一起會難維護。
- **C. 做成 pre-commit 而非 commit-msg validator**：reject — commit-msg hook 可看 body，能驗證 opt-out / 交代文字。

### Why Not an ADR Yet

只是補單一 validator，沒改架構。未來若 plan archival 規則大幅變更（例如 plan 改用 YAML schema），再評估 ADR 升級。

### ADR Promotion Criteria（completed 時驗證）

- [ ] `validatePlanArchivalAudit` Go validator 上線（呼叫 `ScanCheckboxesInFile`，不內聯 loop）
- [ ] `ai-skill scan-checkboxes <file> [--format plain|json] [--exit-code]` CLI subcommand 可用
- [ ] `obligation.commit.plan_archival_audit` 進 `runtime/core-bootstrap.yaml`
- [ ] `gate.plan_archival_audit_required` 進 `runtime/cli-modification-policy.yaml`
- [ ] Fixture tests 全綠（happy / unjustified-[]-block / justified-[]-pass / opt-out / 純內容修改不觸發）
- [ ] Unit tests 全綠（`ScanCheckboxes` pure content / file / all-checked / empty）
- [ ] CLI docs（command-contract）更新（validator entry + subcommand entry）
- [ ] Self-test：本 plan archive 時實際走一遍 validator（dogfood，與 gen3 plan 同模式）

---

## Runtime Execution Path

### Doc-only Trial 聲明 + Graduation

**目前狀態（2026-05-28）**：Plan 為 **draft**，全部 Phase 0–2 為 `[ ]`，de facto doc-only。

| Graduation Phase | Acceptance signal |
|---|---|
| Phase 2 完成 | 第 19 個 commit-msg validator active；fixture tests 全綠 |

### Runtime owner

- Validator：`scripts/ai-skill-cli/internal/app/hooks.go`
- Obligation：`runtime/core-bootstrap.yaml` §per_commit_obligations
- Policy gate：`runtime/cli-modification-policy.yaml`

### Trigger flow

```
event_or_signal:
  - file_diff: plans/active/*.md deleted AND plans/archived/*.md added with same basename
  - OR git rename detection plans/active/X.md -> plans/archived/X.md

detector:
  - commit-msg validatePlanArchivalAudit

loaded source / contract:
  - archived plan file content
  - commit body text

runtime action / blocking gate:
  - block when archived version contains `- [ ]` line not justified in commit body
  - opt-out: [skip-plan-archival-audit] trailer

observable evidence:
  - validation/scenarios/failure-derived/plan-archival-unresolved-checkbox-v1.yaml
  - Go fixture tests in scripts/ai-skill-cli/internal/app/plan_archival_audit_test.go
```

### Generated surfaces (with named consumer)

| Generated surface key | Named consumer | Consumer 類型 |
|---|---|---|
| `obligation.commit.plan_archival_audit` | `validatePlanArchivalAudit` (hooks.go) + commit-msg dispatcher | Go validator + dispatcher |
| `gate.plan_archival_audit_required` | cli-modification-policy.yaml + validateCLIDocSync | policy gate + CLI doc sync |

### Validation scenarios

- `plan-archival-unresolved-checkbox-v1`：plan archived with `- [ ]` lines and no body justification → block
- `plan-archival-justified-checkbox-v1`：plan archived with `- [ ]` but body cites "deferred → <handover>" → pass
- `plan-archival-opt-out-v1`：plan archived with `[skip-plan-archival-audit]` → bypass

---

## Open Questions

| # | Question | 影響範圍 |
|---|---|---|
| 1 | 「明文交代」的 keyword set 要多寬？只認 deferred / non-goal / scope-reduced / handover / 延後 / 拆分，還是更寬？ | Phase 1 / 2 |
| 2 | 是否要求每個 `- [ ]` 與 body 的 justification 一一對應，還是只要 body 有 ANY justification keyword 就放行？ | Phase 1 |
| 3 | 純 prose / refactor 修改 plan（不 archive）會誤觸嗎？ | Phase 0 trigger 設計 |

---

## 完成條件

- [ ] Phase 0–2 全部達成
- [ ] ADR Promotion Criteria 全綠
- [ ] `validatePlanArchivalAudit` 為第 19 個 commit-msg validator
- [ ] `ai-skill scan-checkboxes` CLI subcommand 可跨專案使用（任意 Markdown checklist）
- [ ] 對應 scenarios + Go fixture tests（validator）+ unit tests（scan utility）
- [ ] 自身 archive 時通過 validator（dogfood）

---

## Phase 0 — Pre-Build Interrogation

| 欄位 | 內容 |
|---|---|
| Trigger | Parent plan §Phase 7 follow-up；user 點出 deferred-scenario drift（commit 37065ea） |
| Checked sources | `plans/active/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md` §Phase 7 / `scripts/ai-skill-cli/internal/app/hooks.go`（plan_checkbox_sync 參考）/ `runtime/core-bootstrap.yaml` / `runtime/cli-modification-policy.yaml` |
| Goal | 把 Phase 7 manual completion audit 機械化為 commit-msg validator |
| Scope | 一個新 validator + obligation + gate + 3 scenarios + fixture tests + CLI doc |
| Non-goals | 不改 plans/README.md template；不擴大到 active plans 的 in-progress drift（那是 `validatePlanCheckboxSync` 的責任）；不做 plan content semantic validation |
| Acceptance | validator active；3 scenarios green；本 plan dogfood archive 通過 |
| Duplication risk | 與 `validatePlanCheckboxSync` 互補：前者在 archival moment 嚴查；後者在 in-progress commit 嚴查 plan 推進 |
| Open questions | 見 §Open Questions |
| Decision | proceed |

---

## Phase 1 — Test-First Validation

### Tasks

- [x] 新增 `validation/scenarios/failure-derived/plan-archival-unresolved-checkbox-v1.yaml`
- [x] 新增 `validation/scenarios/failure-derived/plan-archival-justified-checkbox-v1.yaml`
- [x] 新增 `validation/scenarios/failure-derived/plan-archival-opt-out-v1.yaml`
- [x] 每個 scenario 綁定 Go fixture test（同 parent plan §Phase 1 schema lock 規則）

### Phase 1 完成條件

- [x] 3 scenarios 符合 `validation/scenario.schema.json`
- [x] Runtime validate 通過

---

## Phase 2 — Shared Utility + Validator Implementation + Self-Dogfood Archive

### Design: `scan-checkboxes` 三層架構

Phase 2 實作分三層，checkbox 掃描邏輯集中在 `scan_checkboxes.go`，validators 與 CLI 共用同一實作：

```
Layer 1 — Pure utility（零 shell 依賴，純 Go strings）
  ScanCheckboxes(content string) CheckboxScanResult
    → 掃描任意字串，回傳 UncheckedLines / CheckedLines

Layer 2 — File helper
  ScanCheckboxesInFile(path string) (CheckboxScanResult, error)
    → 讀檔 + 呼叫 Layer 1，供 hooks.go validator 呼叫

Layer 3 — CLI subcommand
  ai-skill scan-checkboxes <file> [--format plain|json] [--exit-code]
    → 呼叫 Layer 2，輸出結果
    → --exit-code: 有 unchecked 時 exit 1（可用於任何專案的 CI / pre-push hook）
```

**跨專案用途**：任何專案的 Markdown checklist（release checklist、onboarding checklist、規格確認單）都能用 `ai-skill scan-checkboxes <file> --exit-code` 驗證，讓使用者方便追蹤完成狀態。

**`CheckboxScanResult` struct**：

```go
type CheckboxScanResult struct {
    File          string   // 空字串表示 raw content scan
    UncheckedLines []string // "- [ ] 任務A"
    CheckedLines   []string // "- [x] 任務B"
}

func (r CheckboxScanResult) HasUnchecked() bool { return len(r.UncheckedLines) > 0 }
```

**`validatePlanArchivalAudit` 呼叫方式**：
```go
result, _ := ScanCheckboxesInFile(archivedPath)
if result.HasUnchecked() && !bodyJustifiesUnchecked(commitBody) {
    return "plan-archival-audit: ..."
}
```

### Tasks

- [x] 新增 `scripts/ai-skill-cli/internal/app/scan_checkboxes.go`
  - `CheckboxScanResult` struct + `HasUnchecked()`
  - `ScanCheckboxes(content string) CheckboxScanResult`（Layer 1，純 Go，不呼叫任何 shell）
  - `ScanCheckboxesInFile(path string) (CheckboxScanResult, error)`（Layer 2）
  - `runScanCheckboxes(args []string, stdout, stderr io.Writer) int`（Layer 3 CLI handler）
- [x] 新增 `scripts/ai-skill-cli/internal/app/scan_checkboxes_test.go`（≥ 4 個 unit tests：pure content / file / all-checked / empty）
- [x] 在 `app.go` 加 `case "scan-checkboxes"` + `printUsage` 說明
- [x] 新增 `validatePlanArchivalAudit` 在 `hooks.go`（呼叫 `ScanCheckboxesInFile`，不內聯 loop）
- [x] 新增 `findArchivedPlans(staged []string) []string`（偵測 active→archived move，純 Go regex，不用 grep）
- [x] 註冊 `obligation.commit.plan_archival_audit` 在 `runtime/core-bootstrap.yaml` §per_commit_obligations
- [x] 加 `gate.plan_archival_audit_required` 在 `runtime/cli-modification-policy.yaml`
- [x] Fixture tests（≥ 5 個：happy / unjustified block / justified pass / opt-out / non-archival commit no-trigger）
- [x] 更新 `scripts/ai-skill-cli/docs/command-contract.md`：新增第 19 個 validator entry + `scan-checkboxes` subcommand entry
- [x] bin rebuild
- [ ] **Dogfood**：本 plan 自身 archive 時通過 validator（deferred — 本 plan 仍 active，待功能成熟後 archive）

### Phase 2 完成條件

- [x] `scan_checkboxes.go` 含三層實作，零 shell 依賴
- [x] `ai-skill scan-checkboxes <file> --exit-code` 可在任何專案使用
- [x] 第 19 個 commit-msg validator active
- [x] per_commit_obligations 含 `obligation.commit.plan_archival_audit`
- [x] cli-modification-policy 含 `gate.plan_archival_audit_required`
- [x] Unit tests（scan utility）+ fixture tests（validator）全綠
- [x] CLI docs updated（subcommand + validator）
- [ ] Dogfood self-test 通過（deferred）

---

## Stakeholder 同意項目

- [ ] 接受新增第 19 個 commit-msg validator
- [ ] 接受 block default（與 `validateRuntimeTriggerWiring` 一致）
- [ ] 接受 opt-out trailer `[skip-plan-archival-audit]`

---

## 與其他 plans 的關係

| Plan | 關係 |
|---|---|
| [`plans/active/2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md`](2026-05-28-1200-gen3-runtime-trigger-audit-and-completion.md) | **Parent**。本 plan 是 §Phase 7 follow-up；parent 自身 archive 時受本 validator 保護（若本 plan 比 parent 早 Phase 2 graduate）|
| Gen 4 audit-extension cluster（未建）| **Sibling pattern**。Plan-archival audit 與 routing-orphan audit / glossary coverage 同屬 "static analysis at commit time" 家族 |

---

## 預估規模

| Phase | 變動 | LOC |
|---|---|---|
| Phase 1 | 3 scenarios YAML | ~150 |
| Phase 2 | `scan_checkboxes.go`（utility + CLI）+ validator + obligation + gate + tests + docs | ~350 |
| **Total** | | **~500**，3–4 commits |
