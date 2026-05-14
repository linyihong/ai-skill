# Development Guidance Intelligence

`intelligence/engineering/development/` 存放與「開發指引」相關的工程智慧。這些 atoms 描述如何將安全觀察轉譯為開發者行動、如何維護文件契約、以及如何確保行為變更被正確記錄。

## 目前 atoms

| Atom | 描述 |
|------|------|
| [`docs-first-bdd-closure.md`](docs-first-bdd-closure.md) | 文件優先 BDD 閉環：observable behavior 變更前必須先更新管轄契約 |
| [`risk-translation-heuristic.md`](risk-translation-heuristic.md) | 風險轉譯：將攻擊者視角的觀察轉譯為開發者行動 |
| [`contract-governance-heuristic.md`](contract-governance-heuristic.md) | 契約治理：文件衝突時的優先順序與解決方法 |

## 與其他層的關係

- `workflow/software-delivery/` 提供執行流程，本層提供背後的原則與 why。
- `analysis/development-guidance/` 提供具體的分析方法，本層提供選擇方法的決策邏輯。
- `skills/app-development-guidance/` 是原始來源，已刪除。內容已由新分層承接。
- `feedback/history/development-guidance/` 存放本領域的 feedback lessons。
