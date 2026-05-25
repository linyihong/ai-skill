# Template Drift（文件模板漂移）

Status: validated
Class: `template-inconsistency` / `governance-drift`

## Trigger

當 agent 建立新的 failure pattern / intelligence atom / validation scenario / workflow / analysis 等可重用 reusable layer 文件時，**沒有先讀同類既有文件的 canonical 結構**，憑直覺或記憶決定 section header 與必要欄位，導致：

具體觸發訊號：

- 新文件 section header 用中文（如 `## 驗證`）但同類既有文件多用英文（`## Validation`）
- 新文件缺必要 section（如 failure pattern 缺 `## Linked Validation Scenarios`）
- 新文件 section 順序與既有不一致
- 同類目錄下混用兩種以上 template 風格（例：5 個用英文 + 3 個用中文）
- 工具的 grep / scanner 因 section 名稱不一致而誤判完整性

## Failure Mode

模板逐步漂移，破壞下列系統性檢查：

1. **Validation scenarios 失效**：依 section header 偵測缺失的 scenario（如 `validate_failure_pattern_template_consistency`）會誤報或漏報
2. **Cross-doc 比對困難**：reviewer / agent 看不出某 section 是否相同概念
3. **Template 失去 forcing function**：新作者看不到統一範本，繼續用直覺寫
4. **跨世代漂移累積**：每次升級多一個風格，最終形成「考古學」狀態

## Risk

- **靜默違反**：drift 通常不會 break runtime，所以沒有阻斷訊號
- **規模放大**：每個漂移檔案降低後續 reviewer 的標準感
- **Validation Token 浪費**：scanner 必須維護「同義 section header」清單來補救
- **新成員學習成本**：看不到單一 canonical reference

## Required Agent Action

建立或修改 reusable layer 文件前：

1. **先 list_files 該目錄** — 看有多少同類檔案
2. **挑 1-2 個最近建立的同類檔** Read 全文 — 作為 template 參考
3. **抓出 canonical section list**：
   - failure-patterns/：`## Trigger` / `## Failure Mode` / `## Risk` / `## Required Agent Action` / `## Prevention Gate` / `## Validation` / `## Linked Rules` / `## Linked Validation Scenarios` / `## Source` / `## Related`
   - intelligence/<sub>/：依各 subdir README 內既有 atom 結構
   - validation/scenarios/：依既有 YAML schema
4. **依 canonical list 寫**，section 名稱 100% match（不翻譯、不簡寫、不省略）
5. **必要時更新 README 索引** — 加新檔到索引表

## Prevention Gate

開始建立 reusable layer 文件前，agent 必須能回答：

| Check | Required answer |
|-------|-----------------|
| 目錄已有幾個同類檔 | `ls <target-dir>/*.md \| wc -l` |
| 我參考了哪個檔當 template | 至少 1 個既有檔的路徑 |
| Canonical section list | 列出本檔將使用的 section headers |
| 是否會出現新 section name | 若是，理由是什麼？是否更新 README / scanner 規則？ |
| Index 更新 | 是否已加新檔到對應 README 索引表 |

若任一無法清楚回答，**先讀同類既有檔案再下筆**。

## Validation

符合下列條件時，此 pattern 已被防止：

- 新增的 failure pattern / intelligence atom / validation scenario 通過對應 template-consistency scenario
- 同類目錄內所有檔案的 section header 風格一致（中或英擇一，不混用）
- README 索引表與實際檔案 1:1 對應
- 跨世代升級時，舊風格檔案在升級 commit 一併對齊新風格

## Linked Rules

- [`enforcement/dependency-reading.md`](../dependency-reading.md) — 修改前讀依賴的延伸：建立前讀模板
- [`enforcement/neutral-language.md`](../neutral-language.md) — 中性低爭議與一致性
- [`enforcement/linked-updates.md`](../linked-updates.md) — 索引與相關文件同步
- [`governance/lifecycle/system-upgrade-governance.md`](../../governance/lifecycle/system-upgrade-governance.md) §3 — 升級 governance（template 屬其中一個維度）

## Linked Validation Scenarios

- `validation/scenarios/failure-derived/failure-pattern-template-consistency-v1.yaml` — failure-patterns/ 結構檢查（揭露本 pattern 的具體案例）
- `validation/scenarios/failure-derived/readme-index-consistency-v1.yaml` — README 索引與實際檔案對應

## Source

- 2026-05-22 Phase D Trial T4 (commit `7dc96d2`) — failure-patterns/ audit 發現 4 個檔結構不一致：
  - `analysis-domain-discovery-gap.md`：全中文 section（症狀/根本原因/預防方式）
  - `framework-duplication-without-interrogation.md`：`## 驗證` + 缺 `## Linked Validation Scenarios`
  - `knowledge-update-flow-bypassed-by-sub-pipeline.md`：`## Prevention` 而非 `## Prevention Gate`
  - `premature-adr-promotion.md`（agent 自寫）：`## 驗證` + 缺 `## Linked Validation Scenarios`

- 2026-05-22 Phase D Trial scenario 套用（commit `ef305bf`）— `failure-pattern-template-consistency-v1.yaml` 揭露更廣的 drift：9 個檔的 `## 驗證`、1 個檔的 `## Prevention`、1 個檔的全中文 section。本 pattern 即為這次發現的沉澱。

- 修復 commit chain：
  - `7dc96d2`（T4，部分修）→ `ef305bf`（scenario 揭露完整 drift）→ 本 commit（補完 + 沉澱 pattern）

## Related

- [`failure-pattern-template-consistency-v1`](../../validation/scenarios/failure-derived/failure-pattern-template-consistency-v1.yaml) — 直接對應 scenario
- [`entrypoint-positioning-drift.md`](entrypoint-positioning-drift.md) — 同類「漂移而非中斷」的失效

← [Back to failure patterns](README.md)
