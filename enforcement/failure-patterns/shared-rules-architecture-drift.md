# Enforcement Architecture Drift（enforcement 與架構不同步）

Status: validated
Class: `dependency-miss` / `validation-gap`

## Trigger

當 repository 發生架構重構（目錄重組、分層新增、路徑變更、命名變更），且 agent 更新了主要檔案（workflow、intelligence、analysis）但**沒有同步檢查 enforcement/ 中的路徑參考是否過期**時，使用此 pattern。

典型觸發情境：

- 新增或重組 `workflow/<domain>/`、`analysis/<domain>/`、`intelligence/<domain>/` 等新分層
- 將舊 `skills/<skill-name>/` 內容提取到新分層
- 修改 `skills/` 結構後，未檢查 `enforcement/` 中的範例路徑、表格、模板是否仍指向舊路徑
- 只更新了「被提取的檔案本身」，但沒有更新「描述如何提取的規則」

## Failure Mode

Agent 專注在「搬移內容」或「重組目錄」的主要任務上，卻忽略了 `enforcement/` 中可能包含指向舊結構的參考。這些參考包括：

1. **範例路徑**：`governance/document-sizing.md` 中的目錄結構範例
2. **表格內容**：`linked-updates.md` 中的連動關係表
3. **模板**：`feedback-lessons.md` 中的 Promotion Target 模板
4. **索引**：`enforcement/README.md` 中的 lazy-load 表格
5. **教學文件**：`skills/ADDING_SKILLS.md` 中的目錄結構建議
6. **規則正文**：`content-layering.md`、`tool-neutral-documentation.md` 中的路徑描述
7. **流程規則**：`decision-efficiency.md` 中的 Context Loading 步驟
8. **引用規則**：`cross-skill-references.md` 中的引用格式

## Root Cause Analysis

### 為什麼 agent 會漏掉 enforcement 同步？

| 原因 | 說明 | 對應預防 |
|------|------|---------|
| **任務焦點偏誤** | Agent 將「架構重構」視為「搬移檔案」任務，而非「全庫一致性更新」任務。完成主要搬移後，認為任務已完成，沒有意識到 enforcement 也需要更新 | 在 task plan 中加入 enforcement 檢查項 |
| **enforcement 不在預設讀取範圍** | `dependency-reading.md` 的 Default Bootstrap Boundary 沒有要求架構重構時讀取 enforcement | 在 pipeline Step 7a 中明確定義觸發條件 |
| **無自動化檢查** | 沒有 grep 命令或 validation gate 在架構重構完成後自動檢查 enforcement 是否過期 | 在 pipeline Step 7a 中加入 grep 驗證命令 |
| **linked-updates.md 表格不完整** | 原有連動關係表沒有「架構重構 → enforcement 同步」這一項 | 在 linked-updates.md 中加入此連動關係 |
| **無 failure pattern 預防** | 這是第一次發生此類錯誤，之前沒有對應的 failure pattern 可以觸發 prevention gate | 建立本 pattern |

### 為什麼這次（enforcement 整合）又發生？

這次的「enforcement 整合」任務本身就是要更新 enforcement，但：

1. **初始分析時只識別了 9 個受影響檔案**，但沒有意識到「這個錯誤本身也應該被記錄為 failure pattern」
2. **沒有在 task plan 中加入「記錄 failure pattern」的步驟**——只專注在「修好檔案」，沒有做「為什麼會需要修」的反思
3. **沒有建立 validation scenario**——即使修好了，未來同類錯誤仍可能重演，因為沒有 stateless 驗證

## Risk

- 即使新分層檔案已完整建立，`enforcement/` 仍指向舊路徑，導致未來 agent 讀到過期資訊
- 新加入的開發者或 agent 會依照 enforcement 的範例建立舊結構檔案
- 架構重構的「完整性」被破壞：使用者看到新分層已建立，但規則仍說「放 skills/」
- 每次架構變更都需要人工檢查 enforcement，增加維護成本
- 若未及時修正，累積的 drift 會讓 enforcement 失去作為 source of truth 的信任

## Required Agent Action

進行任何架構重構（目錄重組、分層新增、路徑變更）時，**必須**執行 [`governance/lifecycle/intelligence-extraction-pipeline.md`](../../governance/lifecycle/intelligence-extraction-pipeline.md) 的 **Step 7a：Enforcement 同步檢查**，該步驟已定義完整的檢查範圍、更新原則與驗證命令。

此外，**必須**在 task plan 中加入以下步驟：

```
□ 架構重構完成後：
   1. 執行 intelligence-extraction-pipeline.md Step 7a（enforcement 同步檢查）
   2. 判斷本次 failure 是否需要記錄為 failure pattern
   3. 若需要，建立 failure pattern 並加入 failure-patterns/README.md 索引
   4. 判斷是否符合 validation scenario 條件（stateless、有明確 expected/forbidden route、有 prevention 價值）
   5. 若符合，建立 validation/scenarios/failure-derived/<id>.yaml
   6. 更新 failure-learning-system.md 的 Promotion Decision（如需要）
```

## Prevention Gate

進行任何架構重構前，在 task plan 或 `.agent-goals/` 中加入以下檢查項：

```
□ 架構重構完成後，已執行 intelligence-extraction-pipeline.md Step 7a
   - 已搜尋 enforcement/ 中所有指向舊結構的參考
   - 已分類受影響的參考類型
   - 已為每個受影響檔案新增新分層路徑 + 保留舊路徑向後相容
   - 已執行 grep 驗證無遺漏
   - 已依 linked-updates.md 完成連動更新
□ 已判斷是否需要記錄為 failure pattern
□ 已判斷是否需要建立 validation scenario
```

## Validation Scenario

- [`validation/scenarios/failure-derived/shared-rules-architecture-drift-v1.yaml`](../../validation/scenarios/failure-derived/shared-rules-architecture-drift-v1.yaml)

## 驗證

1. 對 `enforcement/` 執行 `grep -rn "skills/<old-path>"` — 應無未預期的舊路徑參考
2. 對 `enforcement/` 執行 `grep -rn "workflow/<domain>/"` — 應出現新分層路徑
3. 從 `enforcement/README.md` 開始逐檔讀取，確認每個檔案的路徑參考都是最新的
4. 確認 failure pattern 本身有記錄 root cause 和 validation scenario
5. Commit/push/readback 後確認 `git status --short --branch` 乾淨

## Linked Validation Scenarios

- `validate_directory_structure` — 檢查 `enforcement/` 中各目錄的 README 是否列出所有子檔案，防止架構重構後 enforcement 路徑未同步
- `validate_failure_pattern_validator_coverage` — 檢查每個 failure pattern 的 Linked Validation Scenarios 是否為空，防止架構重構後未記錄 failure pattern

## Linked Rules

- [`../../governance/lifecycle/intelligence-extraction-pipeline.md`](../../governance/lifecycle/intelligence-extraction-pipeline.md)（Step 7a）
- [`../failure-learning-system.md`](../failure-learning-system.md)
- [`../linked-updates.md`](../linked-updates.md)
- [`../content-layering.md`](../content-layering.md)
- [`../tool-neutral-documentation.md`](../tool-neutral-documentation.md)
- [`../governance/document-sizing.md`](../governance/document-sizing.md)
- [`../decision-efficiency.md`](../decision-efficiency.md)
- [`../cross-skill-references.md`](../cross-skill-references.md)
- [`../feedback-lessons.md`](../feedback-lessons.md)
- [`../../skills/ADDING_SKILLS.md`](../../skills/ADDING_SKILLS.md)
