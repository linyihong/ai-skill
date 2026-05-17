# Compiler Philosophy

## Why Compiler

- **Prose 是人讀的**：workflow、enforcement、governance 文件包含 judgment、heuristics、troubleshooting strategy，這些需要 human 才能理解。
- **YAML 是機器讀的**：runtime 需要 machine-readable 的 state（current_phase、allowed_actions、blocking_gates）才能控制執行流程。
- **Compiler 是橋樑**：將 prose 中的 deterministic state 提取出來，轉換為 YAML，同時保持 prose 作為 canonical source。

## Core Principles

1. **Deterministic Only**：只編譯 execution-critical state（current_phase、allowed_actions、blocking_gates、required_artifacts、open_obligations、transaction_state）。Heuristics、judgment、troubleshooting 永遠留在 prose。
2. **Prose is Canonical**：所有修改應在 prose source 進行，再透過 compiler 更新 generated YAML。不應手動編輯 generated YAML。
3. **Idempotent**：相同 prose source → 相同 generated YAML。重複執行不改變結果。
4. **Sync on Commit**：pre-commit hook 檢查 prose 與 YAML 是否一致，不一致則 block commit。

## Runtime Surface

| Source | Target | 提取內容 |
|--------|--------|----------|
| `workflow/*/execution-flow.md` | `runtime/generated/workflow-{domain}-phases.yaml` | phase definitions、allowed/forbidden actions、blocking gates |
| `workflow/*/artifact-gates.md` | `runtime/generated/workflow-{domain}-artifacts.yaml` | required artifacts、verification criteria |
| `enforcement/dependency-reading.md` | `runtime/generated/transaction-machine.yaml` | transaction states、rules |
| `enforcement/goal-action-validation.md` | `runtime/generated/goal-action-gates.yaml` | validation gates、criteria |
| `enforcement/failure-learning-system.md` | `runtime/generated/failure-recovery.yaml` | failure patterns、recovery strategies |

## 與既有文件的關係

- [`runtime/compiler/README.md`](../../runtime/compiler/README.md) — Runtime navigation entry point
- [`runtime/compiler/compiler-engine.rb`](../../runtime/compiler/compiler-engine.rb) — Compiler implementation
- [`runtime/compiler/compiler-rules.yaml`](../../runtime/compiler/compiler-rules.yaml) — Source-target mapping rules
- [`runtime/generated/`](../../runtime/generated/) — Compiled YAML surfaces
