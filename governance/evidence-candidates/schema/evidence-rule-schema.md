# Evidence Rule Schema（shape only — 空殼）

每個 plan 宣告「什麼 artifact 算我的候選證據」。**本檔只凍 shape，不寫任何 plan 的 criteria 內容**
（criteria authoring 是 Phase 1A Step 3）。

Canonical 設計：[`plans/active/2026-06-16-1131-evidence-candidate-system.md`](../../../plans/active/2026-06-16-1131-evidence-candidate-system.md)。

## Shape（`evidence-rule-v1`）

```yaml
evidence_rule:
  collect: bool            # 是否收集本 plan 的候選
  match:
    artifact_types: []     # commit / adr / test / doc / issue ...
    criteria:
      - id:                # criterion 識別碼
        description:        # 人讀說明（candidate.criteria_hits 對應此 id）
  exclusions: []           # 明確排除的情況
```

## 邊界

- `evidence_rule` 是 plan 自己宣告的判準。rule **定義 owner = plan**（內嵌「Evidence Rule」section）；
  `../evidence-rules/*.pointer.yaml` 只 index，不複製內容。
- **`acceptance-gate` 是 sibling schema（未來）**：`notify` 是 gate 的 projection（達 gate 即「可以
  review」≠「已升級」），不獨立成 rule，避免第二個 state machine。`schema/` 採目錄形式正是為了預留
  `acceptance-gate-schema.md` 而不需重命名。
