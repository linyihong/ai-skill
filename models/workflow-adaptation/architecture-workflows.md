# Architecture Workflows

Architecture workflow 用於 planning、tradeoff、layer responsibility、contradiction analysis 與 migration / promotion decision。

## Workflow Shape

1. 界定 current decision。
2. 讀 architecture / governance / owner README。
3. 列出 constraints、tradeoffs、source-of-truth。
4. 檢查 overreach：是否把 design doc 提前變 runtime invariant。
5. 若要 implementation，轉回 bounded execution plan。

## Required Evidence

- Current architecture source。
- Relevant governance / lifecycle rules。
- Known conflicting plans or ADRs。
- Validation or rollback path。

## Forbidden

- 從單一 memory 或 summary 推導 global architecture。
- 未檢查 existing layer boundary 就新增 layer。
- 未驗證 recurring failure 就新增 runtime guard。
- 把 tool-specific limitation 寫成 tool-neutral invariant。
