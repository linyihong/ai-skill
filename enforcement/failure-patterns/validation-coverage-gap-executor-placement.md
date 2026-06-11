# Validation Coverage Gap — Executor Placement

Status: validated
Class: `validation-coverage-gap` / `executor-placement`

## Trigger

當下列事實鏈成立時使用此 pattern：

1. **規則有了** — canonical rule yaml / 規範說明書存在
2. **Validator 有了** — lint / check function 已實作可執行
3. **Validator placement 錯了** — 只在某個次要 phase（如 `runtime compile`、`runtime validate`）跑，沒有進 commit transaction
4. **Side-channel 旁路存在** — workflow `paths:` filter、partial trigger、特定 commit 類型不過該 phase
5. **結果**：違反規則的狀態 silently 進入 main，下一個 trigger 該 phase 的 commit 才暴露

具體訊號：

- Commit A 修改 plan / metadata，只動 `plans/` 或 `governance/` → workflow paths filter 不命中 → CI 不跑 → broken state land
- Commit B 動了會 trigger CI 的 path → 跑同一套 validator → 發現 commit A 留下的破壞 → 看起來像 B 的 fault

## Failure Mode

這個 pattern 跟 `rule-without-executor.md` 是**同家族但不同變種**：

| Pattern | Rule | Executor | Placement | 失效原因 |
|---|---|---|---|---|
| `rule-without-executor` | 有 | 無 / 待寫 | — | 規則無法強制 |
| **本檔（executor-placement）** | 有 | **有** | **錯位** | 規則可強制，但被旁路 silently bypass |

執行點選錯會讓「有 validator」假象變得比「無 validator」更危險：
- Reviewer 看到 lint pass 就放心 merge
- 實際 broken 狀態 silently propagate 直到下個剛好命中的 commit
- 暴露點離 root cause commit 距離越遠，debug 越貴

### 2026-06-06 採樣

- **Validator**：`pending_implementation_child_plan_validity` (`scripts/ai-skill-cli/internal/app/enforcement_registry_lint.go::lintPendingImplementationChildPlanValidity`)
- **Rule**：`enforcement-registry.yaml` `rule_classes[].child_plan` 路徑必須 resolve 到 existing `plans/active/*.md`
- **Placement bug**：lint 只在 `ai-skill runtime compile` 跑（runtime YAML projection / SQLite refresh pipeline 的一部分）
- **Bypass path**：plan-rename commit（`4dde4de` + `94364f7`）只動 `plans/`，沒進 workflow `paths:` filter 的 `scripts/**`，CI 不跑 runtime compile，stale `child_plan` reference silently land 至 main
- **暴露 commit**：`4e9c326` (workflow-activation-discovery-bridge Phase A) 改了 `scripts/**`，下一輪 CI 跑 runtime compile，lint 觸發 — 但已離 root cause 兩個 commit
- **修補**：`93fe950 fix(registry): repoint sanitization child_plan after rename`

## Why It Recurs（為什麼不是「小心一點」就能解）

這不是個人疏忽問題。Repo 演化模式必然產生此 gap：

1. **lint 一開始寫在它最自然的所在地** — 例如 `runtime compile` 才會 load registry YAML，所以 lint 寫在那裡
2. **執行頻率與 trigger 條件耦合到 workflow paths:** — 為了省 CI 成本（plan-only commit 不該跑 Go test）
3. **新的違反規則的 commit shape 出現**（plan-only rename），第一次發生時就吃下整個 coverage gap
4. **「同樣 validator 已存在」反而誤導 reviewer** — 不會去再想是否需要在更早 phase 攔截

## Family Resemblance（與其他 plan / pattern 的關係）

本 pattern 與下列 plan 是同一家族「**Rule exists → Validation exists → Wrong execution point → Failure**」：

- [`plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md`](../../plans/active/2026-06-06-1700-workflow-activation-discovery-bridge.md) — detector miss 時 advisory 機制（補 mechanical floor at PreToolUse hook phase）
- [`plans/archived/2026-06-06-1800-sanitization-mechanical-enforcement.md`](../../plans/archived/2026-06-06-1800-sanitization-mechanical-enforcement.md) — sanitization 從 prose / advisory 升 mechanical at commit phase

差異：parent meta-plan [`2026-05-31-2100-mechanical-enforcement-registry`](../../plans/archived/2026-05-31-2100-mechanical-enforcement-registry.md) 處理的是「executor 缺」，本 pattern 處理的是「executor placement 錯」。前者要寫新 validator，後者只要把既有 validator 接到正確 phase。

## Reusable Application Surfaces

此 pattern 不限於 registry — 任何「lint 只在週期性 refresh phase 跑」的場域都適用：

| 場域 | Validator | 目前 placement | 應加 placement |
|---|---|---|---|
| `enforcement-registry.yaml` child_plan validity | `lintPendingImplementationChildPlanValidity` | `runtime compile` | commit transaction（**本 session 採樣**）|
| `enforcement-registry.yaml` orphan_rule | `lintOrphanRules` | `runtime compile` | commit transaction（同類風險）|
| `runtime-index.sqlite` source checksum | `nativeRuntimeIndexChecksumsCheck` | `runtime validate` + commit-msg `validateRuntimeIndexFreshness` | **已雙 placement，本 pattern 範例** |
| `enforcement-registry.yaml` source_files 引用的檔案存在 | （潛在）| 無 | 評估 |
| topology consistency（shared layer / project layer 邊界）| 待建 | — | commit transaction |

`runtime_index_freshness` 已經是「先 runtime validate single placement → 加 commit-msg validator dual placement」的 precedent（commit `c5874a8`）。本 pattern 把該方法論抽象成可重複套用的家族。

## Required Agent Action

### 修改 / 新增 validator 時

審查 validator placement 必須回答：

1. **此 validator 守護的 invariant 是什麼？** — 例如「registry 引用解析」、「checksum 同步」、「topology 邊界」
2. **違反此 invariant 可能在哪種 commit shape 中產生？** — 哪些 staged file 組合會破壞 invariant？
3. **目前 placement 是否覆蓋所有 trigger commit shapes？**
4. **若 (3) 為否，是否要 dual-placement（runtime + commit transaction）或 placement 上移？**

### 修改 workflow `paths:` filter 時

不可只把 paths 拓寬作為「修 coverage gap」的辦法 — 那是把 validator 變成 CI-only，commit 仍然 silently 進 main，只是 CI 之後才報。Coverage gap 的正確修法是：

- 把 validator 接到 commit transaction（pre-commit / commit-msg hook adapter），讓 commit 本身被擋
- 必要時 workflow paths 也擴，作為 defense in depth（**不**作為主要修法）

### 觀察到 silent state drift 時

當 reviewer 發現 main 上有應被 lint 攔住但已 land 的狀態：

1. 不要只修當下的狀態（治標）
2. 找出該 lint 的當前 placement 與 trigger 條件
3. 識別 bypass commit shape
4. 補 placement 或 raise 至本 pattern 應用層級

## Prevention Gate

**現有機械防護**：

- `enforcement_registry_lint.go::nativeEnforcementRegistryLint` 在 `runtime compile` pipeline 跑 — **只覆蓋 trigger compile 的 commit**
- `validateRuntimeIndexFreshness` 在 `commit-msg` adapter 跑 — 覆蓋所有 commit，是 dual-placement precedent

**缺口**：

- enforcement-registry.yaml 本身的 lint findings（orphan_rule / pending_implementation_child_plan_validity）尚無 commit transaction 層強制
- → 由 [`Commit-Time Registry Reference Consistency`](TBD-plan-link) plan 處理（spawn task chip 2026-06-08）

## Validation

符合下列任一條件即此 pattern 已被防止：

- Validator 至少有一個 placement 覆蓋 commit transaction（不只 runtime compile）
- 違反 invariant 的最小 commit shape（如 plan-only edit）被該 placement 攔截
- 暴露 commit 與 root cause commit 之間距離縮為 0

## Source

- 2026-06-06 session: 採樣 #1（registry child_plan stale reference after sanitization plan rename）
- 廣義家族（plan-move → cross-reference rot）先前已採樣 2 次但 shape 不同：
  - `5b3e089 chore(plans): archive mechanical-enforcement-registry plan + relink 22 cross-references`
  - `4fd626d chore(plans): archive workflow-activation-engine + fix link rot`
- 後續 plan：commit-time validator wire-up（spawn task chip）

← [Back to failure-patterns index](README.md)
