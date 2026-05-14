# Engineering Heuristics

放**經驗法則**。這是 intelligence 核心之一。

## 核心

資深工程師直覺。

## 目前 atoms

| Atom | 原則 | 狀態 | 來源 |
|------|------|------|------|
| [`premature-optimization.md`](premature-optimization.md) | If performance issue is not measured, optimization is likely harmful. | `candidate-intelligence` | 通用軟體工程經驗 |
| [`abstraction-threshold.md`](abstraction-threshold.md) | If abstraction removes more clarity than duplication, do not abstract. | `candidate-intelligence` | 通用軟體工程經驗 |
| [`retry-smell.md`](retry-smell.md) | More than 3 retries often indicates architectural failure, not transient instability. | `candidate-intelligence` | 通用軟體工程經驗 |
| [`single-responsibility-heuristic.md`](single-responsibility-heuristic.md) | If you can't describe what a module does without using "and", it has too many responsibilities. | `candidate-intelligence` | 通用軟體工程經驗 |
| [`test-driven-heuristic.md`](test-driven-heuristic.md) | If writing a test for a function is difficult, the function's design is likely wrong. | `candidate-intelligence` | 通用軟體工程經驗 |
| [`field-confidence-judgment.md`](field-confidence-judgment.md) | 在不確定時使用明確標記（confirmed/candidate/needs capture/meaning unknown），不要 invent 不存在的資訊。 | `candidate-intelligence` | 從 `intelligence/engineering/analytical-reasoning/heuristics/api-documentation-completeness.md` 提取的跨領域通用部分 |
| [`magic-bytes-reference.md`](magic-bytes-reference.md) | 檔案 extension 不可信；magic bytes 是判斷容器類型的最可靠方式。 | `candidate-intelligence` | 從 `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md` 提取的跨領域通用部分 |
| [`document-priority-hierarchy.md`](document-priority-hierarchy.md) | 文件衝突時更新管轄文件，而不是默默地只修正其中一個來源。 | `candidate-intelligence` | 從 `intelligence/engineering/development/contract-governance-heuristic.md` 提取的跨領域通用部分 |

## 與其他層的關係

- 具體執行步驟 → `workflow/`
- 可執行的規則 → `enforcement/`
