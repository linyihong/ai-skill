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

| Source | Target (SQLite) | 提取內容 |
|--------|----------------|----------|
| `workflow/*/execution-flow.md` | `runtime.db → generated_surfaces (type='workflow_phases')` | phase definitions、allowed/forbidden actions、blocking gates |
| `workflow/*/artifact-gates.md` | `runtime.db → generated_surfaces (type='workflow_artifacts')` | required artifacts、verification criteria |
| `enforcement/dependency-reading.md` | `runtime.db → transaction_states` | transaction states、rules |
| `enforcement/goal-action-validation.md` | `runtime.db → gates (type='validation_gate')` | validation gates、criteria |
| `enforcement/failure-learning-system.md` | `runtime.db → generated_surfaces (type='failure_recovery')` | failure patterns、recovery strategies |

## 與既有文件的關係

- [`runtime/README.md`](../../runtime/README.md) — Runtime navigation entry point
- [`scripts/ai-skill-cli/internal/app/runtime_compiler.go`](../../scripts/ai-skill-cli/internal/app/runtime_compiler.go) — Compiler implementation
- [`runtime/runtime.db`](../../runtime/runtime.db) — Source-target mapping rules
- [`runtime/runtime.db`](../../runtime/runtime.db) — Compiled SQLite surfaces (replaces `runtime/generated/`)
## Executable YAML Contract Projection

流程、gate、activation、blocking condition、required evidence 或 failure action 若需要 agent 直接執行，source 層必須提供 companion YAML contract。Contract 不搬進 `runtime/`；它留在 owner layer，並用以下欄位 opt in：

```yaml
runtime_projection:
  enabled: true
  target_key: governance.example.contract
  surface: generated_surfaces
```

Compiler 只會投影 `runtime_projection.enabled: true` 的 contract，避免一般 metadata、graph、validation YAML 變成 runtime 噪音。投影目標是 `runtime.db` 的 `generated_surfaces`；runtime internal config 則仍由 `runtime_config_documents` / projection tables 管理。
