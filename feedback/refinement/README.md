# Workflow Refinement

`feedback/refinement/` 定義「流程精煉」的系統設計。本目錄保存如何從實作經驗、feedback lesson 與 replay 結果中，持續改進 `workflow/` 的可執行流程，讓 workflow 隨著使用經驗演化而不需要每次都從頭設計。

## 核心責任

- Workflow 的持續改進流程（從使用經驗到流程更新）。
- Refinement 的觸發條件（何時應該修改 workflow）。
- Refinement 的版本管理（如何追蹤 workflow 的變更歷史）。
- Refinement 的驗證（如何確認修改後的 workflow 確實更好）。
- 與 `workflow/` 各子目錄的對應關係。

## 核心原則

1. **Refinement 是演化不是重寫**。Workflow 應隨著使用經驗逐步調整，而不是每次發現問題就重新設計。
2. **Refinement 必須有證據**。每個 refinement 應基於至少一次實際使用經驗或 feedback lesson，而不是猜測。
3. **Refinement 保持向後相容**。修改 workflow 時應確保既有 reference 仍可運作，或提供遷移路徑。
4. **Refinement 是選擇性的**。不是每個 feedback 都需要修改 workflow。一次性問題或特殊情境不應觸發 refinement。

## Refinement 觸發條件

| 條件 | 說明 | 優先級 |
| --- | --- | --- |
| Workflow 步驟導致 agent 卡住或繞路 | 表示流程有缺口或模糊地帶。 | high |
| 同一 workflow 被多次繞過或覆寫 | 表示流程不符合實際需求。 | high |
| Feedback lesson 明確建議 workflow 修改 | Lesson 已泛化且經過驗證。 | medium |
| Replay 結果指出 workflow 缺口 | Replay 發現可預防的流程問題。 | medium |
| 新工具或新技術引入 | Workflow 可能需要更新以利用新能力。 | low |
| 定期維護（每季或每半年） | 預防性 refinement。 | low |

## Refinement 流程

```
1. 識別 refinement 候選
   ├─ 來自 replay 結果的 workflow_gap
   ├─ 來自 feedback_history 中標記為 workflow-related 的 lesson
   ├─ 來自 agent 在執行 workflow 時的繞路或覆寫記錄
   └─ 來自 code review 或流程審查的建議

2. 分析問題類型
   ├─ 步驟遺漏：workflow 缺少某個必要步驟。
   ├─ 步驟模糊：workflow 的某個步驟描述不夠清楚。
   ├─ 步驟順序錯誤：workflow 的步驟順序不符合實際需求。
   ├─ 步驟冗餘：workflow 包含不必要的步驟。
   └─ 邊界條件缺失：workflow 未涵蓋特殊情境。

3. 設計修改方案
   ├─ 最小修改原則：只修改有問題的部分。
   ├─ 保持格式一致：使用 workflow/ 的既有格式。
   ├─ 更新 decision point 或 gate 條件（如適用）。
   └─ 加入修改註記（修改日期、原因、來源）。

4. 執行修改
   ├─ 修改 workflow/ 對應文件。
   ├─ 更新 workflow/ 子目錄的 README（如需要）。
   ├─ 更新 knowledge/indexes/README.md（如 routing 變更）。
   └─ 更新 knowledge/runtime/routing-registry.yaml（如需要）。

5. 驗證修改
   ├─ 確認修改後的 workflow 可被 agent 正確解讀。
   ├─ 確認既有 reference 仍可運作。
   ├─ 確認修改有對應的 feedback lesson 或 replay 記錄。
   └─ 執行 validation gates（governance/validation/README.md）。
```

## Refinement 記錄格式

每次 refinement 應記錄：

```yaml
refinement_id: <YYYY-MM-DD>-<slug>
target: <workflow/ 路徑>
trigger: <step_missing | step_vague | step_order | step_redundant | boundary_missing>
source:
  - <feedback_history 路徑或 replay_id>
change_summary: |
  簡述修改內容與原因。
before_after:
  - 修改前：<簡述>
  - 修改後：<簡述>
validation:
  - 既有 reference 仍可運作：<yes | no>
  - Agent 可正確解讀：<yes | no>
  - Linked updates 完成：<yes | no>
```

## Refinement 的版本管理

Workflow 的 refinement 使用以下方式追蹤變更：

1. **Git log**：每次 refinement 作為獨立 commit，commit message 包含 `refinement_id`。
2. **文件內註記**：在 workflow 文件底部或修改處加入修改記錄區塊：

```markdown
## 修改記錄

| 日期 | 修改內容 | 原因 | 來源 |
|------|---------|------|------|
| 2026-05-12 | 新增步驟 3.5：檢查 proxy 連線 | 多次遇到 proxy 失敗未先檢查 | feedback/replay/2026-05-12-proxy-check |
```

3. **Refinement 記錄**：每次 refinement 的詳細記錄保存在 `feedback/refinement/` 目錄（可選，僅重大 refinement 需要）。

## 與其他層的關係

- `feedback/replay/`：Replay 發現的 workflow_gap 是 refinement 的主要輸入。
- `feedback/extraction/`：Refinement 可能發現 intelligence 缺口，進入 extraction 流程。
- `feedback/promotion/`：重大的 workflow refinement 可能需要通過 promotion pipeline。
- `workflow/`：Refinement 的最終目的地。各子目錄的 workflow 文件由此持續改進。
- `governance/validation/`：Refinement 完成後需通過 validation gates。
- `knowledge/indexes/README.md`：Refinement 可能影響 routing，需更新 index。
- `knowledge/runtime/routing-registry.yaml`：Refinement 可能影響 registry 記錄。
