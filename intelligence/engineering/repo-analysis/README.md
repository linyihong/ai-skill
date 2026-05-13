# Repository Analysis Intelligence

## 核心

從已實作完成的 repository 中恢復開發文件、建立追溯性、以及管理文件契約的原則。這些 atoms 描述如何從程式碼反向恢復 domain intent、如何建立需求到測試的追溯鏈、以及如何處理文件衝突。

## 目前 atoms

| Atom | 描述 |
|------|------|
| [`documentation-backfill-heuristic.md`](documentation-backfill-heuristic.md) | 文件回填：從已實作程式碼系統化恢復缺失的開發文件 |
| [`traceability-heuristic.md`](traceability-heuristic.md) | 文件追溯性：建立需求、實作、測試之間的雙向追溯 |

## 與其他層的關係

- `analysis/repo/` 提供具體的分析方法（documentation-backfill、traceability-gate、contract-governance），本層提供背後的原則與 why。
- `workflow/repo-analysis/` 提供執行流程，本層提供選擇方法的決策邏輯。
- `skills/app-development-guidance/process/README.md` 是原始來源，已不再作為 active entrypoint。
