# Engineering Intelligence

## 子目錄

| 子目錄 | 描述 |
|--------|------|
| [`requirements/`](requirements/README.md) | 需求與行為建模智慧（BDD-lite、ambiguity、acceptance、traceability、validation target） |
| [`architecture/`](architecture/README.md) | 架構決策智慧（domain modeling、architectural fit、delivery alignment、boundary/consistency/coupling） |
| [`execution/`](execution/README.md) | 軟體交付執行認知（delivery cognition、task decomposition、validation reasoning、recovery thinking） |
| [`governance/`](governance/README.md) | 工程治理判定模型（contract、validation shape、failure authority、severity、review checklist） |
| [`economics/`](economics/README.md) | 工程成本模型與 timing trade-off（例如 responsive cost curve） |
| [`ui/`](ui/README.md) | UI 工程智慧（responsive design 原理、framework-neutral layout 判斷） |
| [`render-contexts/`](render-contexts/README.md) | UI render context 共享詞彙（narrow_mobile、safe_area、dynamic_resize 等） |
| [`tradeoffs/`](tradeoffs/README.md) | 技術與 delivery 取捨 |
| [`failure/`](failure/README.md) | 失效模式 |
| [`domain/`](domain/README.md) | 業務領域智慧的 legacy / cross-domain entry；DDD tactical modeling 已移至 architecture cognition |
| [`anti-patterns/`](anti-patterns/README.md) | 反模式 |
| [`distributed-systems/`](distributed-systems/README.md) | 分散式系統 |
| [`philosophy/`](philosophy/README.md) | 原始工程哲學與認知框架（Musk Five-Step、first principles、automation last） |
| [`agent-architecture/`](agent-architecture/README.md) | AI Agent 自身運作的智慧（context collapse、rule overload、task routing、attention budgeting、failure recovery、cognitive boundaries、pilot-first validation、failure-to-scenario closure、linked-updates completeness、decomposition strategy、stateless validation） |
| [`development/`](development/README.md) | 開發指引工程智慧（risk translation、docs-first BDD、contract governance） |
| [`analytical-reasoning/`](analytical-reasoning/README.md) | 分析推理智慧（APK analysis、repo analysis、evidence-first routing、heuristics、signals、failure patterns） |
| [`apk-analysis/`](apk-analysis/README.md) | APK 分析領域專屬工程智慧（hook strategy、混合架構安全層映射、領域專用 heuristics） |
| [`language-specific/`](language-specific/README.md) | 語言特定知識（Java、JavaScript、Dart 等語言的 failure patterns、techniques） |
| [`ai-augmented-delivery/`](ai-augmented-delivery/README.md) | AI codegen 工具大幅進入開發流程後的跨工具工程判斷（generation-validation rate parity 等 meta-tool-risk 原則） |

## 與其他層的關係

- `analysis/`（根目錄）提供具體分析方法，`analytical-reasoning/` 提供背後的原則與 why。
- `workflow/` 提供執行流程，本層提供選擇流程、治理判定與成本模型的決策邏輯。
- Requirements cognition 先穩定 behavior / acceptance，再讓 architecture cognition 建立 domain boundary / invariant。
- 舊 `skills/` scaffold 已退役；執行入口在 `workflow/`，分析方法在 `analysis/`，判斷智慧在本層。
- `agent-architecture/` 研究 AI Agent 自身的認知行為模式，與 `enforcement/failure-patterns/` 互補（後者記錄具體失效模式，前者解釋為什麼會發生）。
- IDE 生態系統知識（VS Code Extension 全域設定儲存機制等）已移至 [`intelligence/ide/`](../ide/README.md)，與本層同為 `intelligence/` 下的獨立子目錄。
