# Rule Without Executor（規則寫了但 runtime 沒接執行器）

Status: validated（2026-05-31 session 連續 5 instance 暴露此 meta-pattern）
Class: `meta-governance-gap` / `framework-self-audit-miss`

## Trigger

Knowledge layer（`enforcement/`、`governance/`、`workflow/`、`runtime/*.yaml`、`knowledge/runtime/routing-registry.yaml` 等）宣告一條規則 / obligation / activation trigger，但 Runtime layer（`scripts/ai-skill-cli/internal/app/hooks.go`、`detector.go`、runtime.db state machine）沒有對應的 executor 真正 enforce 它。

具體觸發訊號：

- 新增或修改 `enforcement/*.yaml`、`runtime/*.yaml`、`governance/**/*.yaml`、`routing-registry.yaml` 但 commit 內**沒有**對應的 hooks.go / detector.go 修改
- 規則文字使用 "MUST" / "必須" / "blocking" 等強制語但沒有對應 PreToolUse / PreCommit / commit-msg / Stop hook
- Knowledge layer 的 `activation_triggers` / `obligation` / `gate` 在 runtime SQLite tables 找不到對應 row
- 規則被使用者反覆追問同一個問題（"X 為什麼沒擋住？"）但每次 fix 都只 patch 個案
- 任務 review 時必須說「這條規則是 behavioral」但沒有 explicit `coverage: behavioral_only` 宣告

## Failure Mode

「Knowledge Layer 有規則，Runtime Layer 沒執行器」 → 規則只活在文件裡，agent 在實際執行時無任何 forcing function。導致：

1. **Silent leak**：rule 寫好之後沒人記得補 executor；下一個違反此規則的場景照樣發生
2. **半年延遲偵測**：要等使用者明確追問才會發現「啊我以為這條會擋住」
3. **每次只 fix 個案**：發現 leak 後通常只補當下這條 rule 的 executor，不檢查其他 rule 是否也 leak
4. **虛假的 mechanical 安全感**：governance dashboard / coverage report 顯示「規則齊全」但實際上沒 enforcement
5. **Meta-pattern 無法被偵測**：因為框架自身沒有「rule ↔ executor binding」結構，當下無人能回答「目前有幾條規則沒 executor」

## Risk

- **規則信任崩塌**：當使用者發現 N 條看似 enforced 的規則其實是空殼，會質疑整個 framework
- **Failure pattern 重演**：剛寫進 repo 的規則立刻 leak，self-inconsistency
- **Maintenance cost 倍增**：每個 instance 都要重複「使用者追問 → patch → 寫 failure pattern」流程
- **Layer responsibility 模糊**：分不清「這條應該是 mechanical / behavioral / not_mechanizable」，最後變「先放著」永遠 leak

## 2026-05-31 Session 暴露的 5 個 Inaugural Instances

這 5 個 instance 共享同一個 meta-pattern。詳細 binding 見 [`enforcement-registry.yaml`](../enforcement-registry.yaml)。

| # | Knowledge Source | Runtime Executor 缺失 | 發現方式 | 目前 registry 狀態 |
|---|---|---|---|---|
| 1 | `runtime/core-bootstrap.yaml` `obligation.bootstrap.receipt` | PreToolUse gate enforcing Receipt before non-Read tools | 使用者半年內反覆追問 | ✓ 已補（read-log gate landed） |
| 2 | `knowledge/runtime/routing-registry.yaml` `activation_triggers` | `detector.go`（沒有 Go executor 讀 activation_triggers） | 使用者追問 3 次 | pending（child plan 1900 active） |
| 3 | `governance/lifecycle/capability-discovery-philosophy.md` | Discovery fallback 機制 | 同一 session | behavioral_only（與 workflow_activation Phase 6.1 綁定 sunset） |
| 4 | `enforcement/sanitization.md` + `reusable-guidance-boundary.md` | PreToolUse Write hook + commit-msg validator | 使用者追問 5 次 | pending（child plan 2000 active） |
| 5 | `knowledge/runtime/routing-registry.yaml`（intelligence routes） | intelligence-class routes 的 activation semantics | 同一 session | research_required（research_questions 列具體未解問題） |

完整 regression scenario 鎖定見 [`validation/scenarios/enforcement/2026-05-31-regression-five-instances-v1.yaml`](../../validation/scenarios/enforcement/2026-05-31-regression-five-instances-v1.yaml)。

## Required Agent Action

### 修改 rule yaml 時

每次新增 / 修改 `enforcement/*.yaml`、`runtime/*.yaml`、`governance/**/*.yaml`、`knowledge/runtime/routing-registry.yaml` 時，必須**同一 commit** 處理 enforcement-registry binding：

1. 開啟 [`enforcement/enforcement-registry.yaml`](../enforcement-registry.yaml)
2. 找到對應的 `rule_class` entry；若不存在則新建
3. 宣告 `coverage` 屬於 6-value 哪一種：
   - `mechanical`（有 executor）→ 列 `executors[].symbol`，必須真實存在於 hooks.go / detector.go
   - `pending_implementation`（child plan 在跑）→ 列 `child_plan` + `target_promotion`
   - `behavioral_only`（故意不機械化）→ 列 `rationale` + `sunset_decision.{revisit_when, success_criteria}`（雙必填）
   - `not_mechanizable`（永遠不該機械化）→ 列 `rationale` + `objective_validation_impossible_because`
   - `research_required`（還不知怎麼機械化）→ 列 `rationale` + `research_questions` ≥ 1 + `estimated_unblock_timeline`
   - `deprecated`（移除中）→ 列 `replaced_by` 或 `removal_date`
4. 若是 `mechanical` 且有 instance set（如 workflow_activation 57 routes），填 `coverage_evidence`

### 修改 hooks.go / detector.go 時

每次新增 Go executor 必須同一 commit 在 registry 對應 rule_class 補 `executors[]` entry。若 Go function 是 helper（不在 `binding_required_for` 三類），加進 `internal_helper_allowlist`。

### Review PR 時

- 檢查 commit 是否同時動了 rule yaml 和 registry
- 檢查 registry 的 `coverage` 與 commit 描述一致
- 檢查 behavioral_only 是否補了 sunset_decision 雙欄位
- demotion（mechanical → behavioral_only）必須附 ADR

## Mechanical Enforcement（Phase 3 後）

Phase 3 lint 實作完成後，下列 violation 在 `ai-skill runtime compile` 階段自動 fail：

| Violation | Lint Error | Hard Block |
|---|---|---|
| 新增 rule yaml 但 registry 沒對應 entry | `orphan_rule` | ✓ |
| 新增 hooks.go executor 但 registry 沒 binding（且不在 allowlist） | `orphan_executor` | ✓ |
| `coverage: mechanical` 但 executor symbol 不存在 | `missing_executor_symbol` | ✓ |
| `coverage: behavioral_only` 缺 sunset_decision.revisit_when 或 success_criteria | `behavioral_only_incomplete_sunset` | ✓ |
| `coverage: not_mechanizable` 缺 objective_validation_impossible_because | `not_mechanizable_missing_rationale` | ✓ |
| `coverage: deprecated` 缺 replaced_by 與 removal_date | `deprecated_missing_disposal` | ✓ |
| `coverage: deprecated` 過 removal_date 仍存在 | `deprecated_past_removal_date` | ✓ |

完整 lint contract 見 Phase 7 scenarios：
- [`registry-lint-orphan-rule-v1.yaml`](../../validation/scenarios/enforcement/registry-lint-orphan-rule-v1.yaml)
- [`registry-lint-missing-executor-v1.yaml`](../../validation/scenarios/enforcement/registry-lint-missing-executor-v1.yaml)
- [`registry-lint-behavioral-without-rationale-v1.yaml`](../../validation/scenarios/enforcement/registry-lint-behavioral-without-rationale-v1.yaml)

## 為什麼這個 pattern 比一般 failure pattern 更重要

一般 failure pattern 處理「某個具體規則被違反」；本 pattern 處理「framework 自身缺結構性檢查」。沒有 Layer 2.5 binding，每個未來 instance 都會獨立發生、獨立被使用者追問、獨立 patch。

這條 pattern 的存在本身就是治理閉環的最後一塊 —— 它讓 framework 能夠回答「目前有幾條規則沒 executor」這個 meta-question，並讓答案隨時可機械驗證。

## Cross-References

- [`enforcement/enforcement-registry.yaml`](../enforcement-registry.yaml) — canonical binding 表
- [`enforcement/enforcement-registry.md`](../enforcement-registry.md) — companion philosophy（Layer 2.5 framing）
- [`plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md`](../../plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md) — 建立 registry 的 plan
- Child plan：[`2026-05-31-1900-workflow-activation-engine.md`](../../plans/archived/2026-05-31-1900-workflow-activation-engine.md) — sanitizes routing-registry rules with detector.go (archived 2026-06-05)
- Superseded child plan：[`2026-05-31-2000-mechanical-sanitization-validator.md`](../../plans/archived/2026-05-31-2000-mechanical-sanitization-validator.md) — allowlist-based sanitization validator route
- Child plan（completed 2026-06-11）：[`2026-06-06-1800-sanitization-mechanical-enforcement.md`](../../plans/archived/2026-06-06-1800-sanitization-mechanical-enforcement.md) — metadata-derived sanitization，已 coverage=mechanical
- Related pattern：[`bootstrap-bypass-on-resume`](bootstrap-bypass-on-resume.md) — instance #1 of this meta-pattern（已 fix）
- Related pattern：[`markdown-yaml-sync-drift`](markdown-yaml-sync-drift.md) — adjacent meta-pattern（companion drift）

← [Back to failure-patterns index](README.md)
