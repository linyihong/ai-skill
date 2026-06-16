# Candidate Schema（shape only）

`candidate` 是一筆「這個事件可能值得某些 plan 看一下」的觀察。**只定 shape，不含 criteria 內容。**

Canonical 設計：[`plans/active/2026-06-16-1131-evidence-candidate-system.md`](../../../plans/active/2026-06-16-1131-evidence-candidate-system.md)。

## Shape（`candidate-v1`）

```yaml
candidate:
  id:                      # 唯一 id，例 C-14
  source:                  # MUST reference 原始 artifact，MUST NOT 指向另一個 candidate
    repo:                  # 本 repo 或下游 consuming repo
    artifact:              # commit / ADR / test / doc / issue
    commit:                # commit ref（可選，視 artifact 而定）
  matched_plans: []        # 可能相關的 plan id（平面列表）
  criteria_hits: []        # 命中哪幾條 plan 的 evidence-rule criterion（id 列表）
                           # 禁止 confidence 數字；未來分級才用 support_level: weak|medium|strong
  status:                  # create | accepted | discarded | expired
  expire_reason:           # 僅 status==expired：plan_closed | artifact_stale | exceeded_retention
```

## Lifecycle

```text
create ──┬─ accept    看過、採納 → 寫回對應 plan evidence
         ├─ discard   看過、不採納（= reject）
         └─ expire    沒處理但過期（≠ discard；不算 reviewed）
```

## 約束

- `source` 只能指向原始 artifact，不可指向另一個 candidate（防觀察層自我繁殖）。
- 不輸出 confidence 數字；只列 `criteria_hits`。
- candidate 存於 `../inbox/`（gitignored），accept 後 evidence 寫回 plan，不在此層保留永久 committed 狀態。
