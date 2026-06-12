# Coverage-Evidence Dangling Reference（coverage_evidence 引用不存在的 scenario）

Status: validated
Class: `governance-drift` / `metadata-drift`

## Trigger

當下列事實鏈成立時使用此 pattern：

1. **某 rule_class 宣告 `coverage_evidence`** — `validation_scenarios[]` 或
   `regression_scenarios[]` 列出 scenario yaml 路徑
2. **路徑指向的檔案從未建立**（或被 rename 後沒同步）— 引用是 aspirational filename，
   不是 real evidence
3. **沒有 compile-time 檢查** — coverage_evidence 看起來像「有證據」，實際是裝飾品
4. **下游消費者誤信** — coverage report / governance dashboard / promotion gate 把
   dangling reference 當成已驗證覆蓋

具體訊號：

- `git grep <scenario-name>.yaml` 只在 `enforcement-registry.yaml` 命中，檔案本身不存在
- coverage_evidence 引用的 scenario 路徑 `test -f` 失敗
- rule_class verification_status 標 `scenario_exists: ✓` 但檔案缺席

## Failure Mode

這個 pattern 與 [`rule-without-executor.md`](rule-without-executor.md) /
[`validation-coverage-gap-executor-placement.md`](validation-coverage-gap-executor-placement.md)
同家族（Layer 2.5 meta-governance 的「宣稱 vs 實作」落差），但變種不同：

| Pattern | 宣稱 | 實作落差 |
|---|---|---|
| `rule-without-executor` | 規則存在 | executor 沒寫 |
| `validation-coverage-gap-executor-placement` | executor 存在 | placement 錯位被旁路 |
| **本檔（coverage-evidence-dangling-reference）** | **coverage_evidence 指向 scenario** | **scenario 檔案不存在** |

「有 coverage_evidence」的假象比「沒有」更危險：reviewer 看到 `validation_scenarios: [...]`
就假設覆蓋成立，不會去 `test -f` 每個路徑。dangling reference 可以存活數月而無人察覺。

## 高階 Pattern：Declared Evidence ≠ Existing Evidence

本檔不是 coverage_evidence 專屬的 path-lint，而是一個更高階治理 pattern 的 instance：

> **任何「宣告自己有證據 / 有成員 / 有來源」的結構，其宣告與實存可以分岔。**
> Declared Evidence ≠ Existing Evidence。

`dangling_coverage_ref` 驗的不是「YAML 格式對不對」，而是「**這個宣稱是否屬實**」——
它是 **Evidence Integrity check**，不是 syntax check。任何 declared-membership 結構都會遇到同型失效：

| 結構類型 | 宣告 | 可能的 dangling 失效 |
|---|---|---|
| registry / catalog | rule_class → executor / scenario | symbol 或 scenario 檔不存在 |
| inventory | 列出「應有」資產清單 | 清單項目從未 land（over-promise） |
| index / manifest | path → resource | rename / move 後路徑斷鏈 |
| coverage / evidence map | claim → proof artifact | proof 不存在或不可追溯 |

通用防呆原則：**任何 declared reference 都必須有一個 mechanical resolver**（`test -f` /
symbol grep / API probe），且 resolver 必須在 reference 進 main 前執行。沒有 resolver 的
declared-membership 結構，本質上是「希望」而非「事實」。這與 `rule-without-executor`
（規則無 executor）是同一條公理在不同層的展開：**宣稱必須可機械驗證，否則治理只是裝飾**。

### 2026-06-12 採樣

- **Rule**：`enforcement-registry.yaml` `bootstrap_integrity.coverage_evidence.validation_scenarios`
- **Dangling refs**：`bootstrap-receipt-required-reads-gate-v1.yaml` +
  `bootstrap-bypass-on-resume-v1.yaml` — Phase 2 inventory（`1b523fd`）宣告的檔名，
  從未 land；real evidence 以不同檔名存在
- **Validator**：`scripts/ai-skill-cli/internal/app/scenario_lint.go::LintValidationScenarios`
  的 `dangling_coverage_ref`（FAIL），wire 進 `ai-skill runtime compile`
- **修補**：補齊兩個 scenario（對應 read-log gate 與 bypass-on-resume failure pattern 的
  real regression evidence），使引用 resolve；非僅 repoint metadata

## Why It Recurs

不是個人疏忽。Registry inventory 與 scenario authoring 是兩個動作，天然會 drift：

1. **Inventory 先寫 aspirational 清單** — 規劃階段列出「應該有」的 scenario 名稱
2. **Authoring 排在後面或被遺忘** — scenario 沒跟上，但 inventory 不會自我修正
3. **Rename / archive 斷鏈** — plan/scenario 改名後 coverage_evidence 沒同步
4. **沒有 `test -f` 機械門檻** — 全靠人工抽查，coverage_evidence 越長越沒人逐條驗

## 未來演進與觀察（Future Evolution）

本 executor 目前驗「證據是否存在 + 結構是否合格」，但 evidence governance 還有兩條
未來演進路線（**不是現在要改**，列為 observation）：

1. **Scenario Maturity Ladder（M0–M3）**：detection_command 改 WARNING（而非 FAIL）
   後，structural lint 已天然分兩層——Layer 1「scenario 是否存在」（id/given/when/then
   = FAIL）、Layer 2「scenario 是否成熟」（domain / detection_command = WARNING）。
   這預留了 maturity 階梯：`M0 只有骨架 → M1 有 detection_command → M2 有 regression
   linkage → M3 有可執行 validation`。未來可把 WARNING 升級為 per-level gate，而現有
   設計不阻礙此路。

2. **coverage_target_pct 的 Goodhart 風險**：目前 floor 只看「百分比」（<50 FAIL /
   <80 WARNING）。但 percentage 會被 game——Rule A（10 scenarios）與 Rule B
   （2 scenarios）都可宣告 100%，實際治理強度天差地遠。未來（Gen4/Gen5）應引入
   **coverage count / diversity / criticality** 維度，而非單看 percentage。將
   coverage_target_pct 當「可被 game 的代理指標」持續觀察，避免 metric 變成 target 後
   失去意義。

## Required Agent Action

- 新增或修改 `coverage_evidence.validation_scenarios[]` / `regression_scenarios[]` 時，
  確認每個路徑 `test -f` 通過；引用不存在的 scenario 不可進 main
- 若 inventory 宣告的 scenario 尚未 land：**補 scenario**（real evidence）而非
  把引用當成佔位符留著；確認是早期草稿才 repoint
- regression scenario 應 link 回 failure pattern（`failure_source:` block /
  `enforcement/failure-patterns/<x>.md` 引用 / 置於 `validation/scenarios/failure-derived/`）

## Prevention Gate

**現有機械防護**：

- `scenario_lint.go::LintValidationScenarios` 的 `dangling_coverage_ref`（FAIL）在
  `runtime compile` 攔截 coverage_evidence 引用不存在的 scenario
- `referenced_scenario_structure`（FAIL: id/given/when/then）確保 referenced scenario
  是真正 well-formed 的 BDD 證據，不只是空檔
- F19 `validation_scenario_governance` rule_class（`coverage: mechanical`）binding 此 executor

**缺口（tracked follow-up）**：

- 目前 placement 僅 `runtime compile`；commit-transaction dual-placement 列為 follow-up
  （precedent: `runtime_index_freshness`），見 child plan §Follow-up

## Validation

符合下列任一即此 pattern 已被防止：

- `runtime compile` 對任何 dangling coverage_evidence reference emit FAIL 並 block
- coverage_evidence 引用的每個 scenario 路徑 `test -f` 通過
- regression scenario 可追溯回其 failure pattern 來源

## Source

- 2026-06-12 session：F19 validation_scenario_governance promotion（child plan
  `2026-06-01-0100-validation-scenario-governance-executor`）採樣 bootstrap_integrity
  的兩個 dangling coverage_evidence reference

## Related

- [`rule-without-executor.md`](rule-without-executor.md) — 同家族：規則無 executor
- [`validation-coverage-gap-executor-placement.md`](validation-coverage-gap-executor-placement.md) — 同家族：executor placement 錯位
- [`enforcement/enforcement-registry.yaml`](../enforcement-registry.yaml) — `validation_scenario_governance` rule_class

← [Back to failure patterns](README.md)
