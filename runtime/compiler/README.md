# Runtime Compiler

將 canonical prose source 編譯為 `runtime/generated/*.yaml` 的編譯層。

## 為什麼需要 Compiler

- **Prose 是人讀的**：workflow、enforcement、governance 文件包含 judgment、heuristics、troubleshooting strategy，這些需要 human 才能理解。
- **YAML 是機器讀的**：runtime 需要 machine-readable 的 state（current_phase、allowed_actions、blocking_gates）才能控制執行流程。
- **Compiler 是橋樑**：將 prose 中的 deterministic state 提取出來，轉換為 YAML，同時保持 prose 作為 canonical source。

## 核心原則

1. **Deterministic Only**：只編譯 execution-critical state（current_phase、allowed_actions、blocking_gates、required_artifacts、open_obligations、transaction_state）。Heuristics、judgment、troubleshooting 永遠留在 prose。
2. **Prose is Canonical**：所有修改應在 prose source 進行，再透過 compiler 更新 generated YAML。不應手動編輯 generated YAML。
3. **Idempotent**：相同 prose source → 相同 generated YAML。重複執行不改變結果。
4. **Sync on Commit**：pre-commit hook 檢查 prose 與 YAML 是否一致，不一致則 block commit。

## 使用方式

```bash
# 編譯所有 source
ruby runtime/compiler/compiler-engine.rb

# 只檢查是否需要編譯（exit code 0 = up to date, 1 = stale）
ruby runtime/compiler/compiler-engine.rb --check

# 顯示預期變更（不實際編譯）
ruby runtime/compiler/compiler-engine.rb --diff
```

## Source-Target Mapping

| Source | Target | 提取內容 |
|--------|--------|----------|
| `workflow/*/execution-flow.md` | `runtime/generated/workflow-{domain}-phases.yaml` | phase definitions、allowed/forbidden actions、blocking gates |
| `workflow/*/artifact-gates.md` | `runtime/generated/workflow-{domain}-artifacts.yaml` | required artifacts、verification criteria |
| `enforcement/dependency-reading.md` | `runtime/generated/transaction-machine.yaml` | transaction states、rules |
| `enforcement/goal-action-validation.md` | `runtime/generated/goal-action-gates.yaml` | validation gates、criteria |
| `enforcement/failure-learning-system.md` | `runtime/generated/failure-recovery.yaml` | failure patterns、recovery strategies |

## 檔案結構

```
runtime/compiler/
├── README.md              # 本檔
├── compiler-rules.yaml    # 編譯規則（scope、mapping、rules、workflow）
└── compiler-engine.rb     # 編譯引擎（Ruby CLI）

runtime/generated/
├── README.md              # Generated surfaces 概覽
└── .gitkeep               # 佔位檔
```

## 與其他層的關係

- `runtime/phases/phase-machine.yaml` — 手寫的 phase machine（P0），compiler 未來可從 workflow prose 自動產生
- `runtime/gates/blocking-gates.yaml` — 手寫的 blocking gates（P0），compiler 未來可從 enforcement prose 自動產生
- `enforcement/dependency-reading.md` — prose source，compiler 從中提取 transaction state machine
- `scripts/` — compiler 可被 pre-commit hook 或 close-loop script 呼叫
