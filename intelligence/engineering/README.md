# Engineering Intelligence

## 子目錄

| 子目錄 | 描述 |
|--------|------|
| [`heuristics/`](heuristics/README.md) | 通用軟體工程經驗法則 |
| [`architecture/`](architecture/README.md) | 架構決策智慧 |
| [`tradeoffs/`](tradeoffs/README.md) | 技術權衡 |
| [`failure/`](failure/README.md) | 失效模式 |
| [`domain/`](domain/README.md) | 領域模型智慧 |
| [`anti-patterns/`](anti-patterns/README.md) | 反模式 |
| [`distributed-systems/`](distributed-systems/README.md) | 分散式系統 |
| [`agent-architecture/`](agent-architecture/README.md) | AI Agent 自身運作的智慧（context collapse、rule overload、task routing、attention budgeting、failure recovery、cognitive boundaries、pilot-first validation、failure-to-scenario closure、linked-updates completeness、decomposition strategy、stateless validation） |
| [`development/`](development/README.md) | 開發指引工程智慧（risk translation、docs-first BDD、contract governance） |
| [`analytical-reasoning/`](analytical-reasoning/README.md) | 分析推理智慧（APK analysis、repo analysis、evidence-first routing、heuristics、signals、failure patterns） |
| [`language-specific/`](language-specific/README.md) | 語言特定知識（Java、JavaScript、Dart 等語言的 failure patterns、techniques） |

## 與其他層的關係

- `analysis/`（根目錄）提供具體分析方法，`analytical-reasoning/` 提供背後的原則與 why。
- `workflow/` 提供執行流程，本層提供選擇流程的決策邏輯。
- `skills/` 是原始來源，已不再作為 active entrypoint。
- `agent-architecture/` 研究 AI Agent 自身的認知行為模式，與 `enforcement/failure-patterns/` 互補（後者記錄具體失效模式，前者解釋為什麼會發生）。
- IDE 生態系統知識（VS Code Extension 全域設定儲存機制等）已移至 [`intelligence/ide/`](../ide/README.md)，與本層同為 `intelligence/` 下的獨立子目錄。
