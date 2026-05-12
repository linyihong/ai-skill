# Engineering Heuristics

放**經驗法則**。這是 intelligence 核心之一。

## 核心

資深工程師直覺。

## 目前 atoms

| Atom | 原則 | 狀態 |
|------|------|------|
| [`premature-optimization.md`](premature-optimization.md) | If performance issue is not measured, optimization is likely harmful. | `candidate-intelligence` |
| [`abstraction-threshold.md`](abstraction-threshold.md) | If abstraction removes more clarity than duplication, do not abstract. | `candidate-intelligence` |
| [`retry-smell.md`](retry-smell.md) | More than 3 retries often indicates architectural failure, not transient instability. | `candidate-intelligence` |
| [`single-responsibility-heuristic.md`](single-responsibility-heuristic.md) | If you can't describe what a module does without using "and", it has too many responsibilities. | `candidate-intelligence` |
| [`test-driven-heuristic.md`](test-driven-heuristic.md) | If writing a test for a function is difficult, the function's design is likely wrong. | `candidate-intelligence` |

## 與其他層的關係

- 具體執行步驟 → `workflow/`
- 可執行的規則 → `shared-rules/`
