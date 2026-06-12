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
