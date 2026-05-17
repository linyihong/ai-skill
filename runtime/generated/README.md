# Generated Surfaces (Legacy — Migrating to SQLite)

> **⚠️ Migration in Progress**: Compiler v1.1.0 now outputs to [`runtime.db`](../runtime.db) (SQLite).
> These YAML files are preserved for backward compatibility but are no longer updated by the compiler.
> New agents should query `runtime.db` directly via SQLite.

本目錄存放由 [`runtime/compiler/compiler-engine.rb`](../compiler/compiler-engine.rb) 從 canonical prose source 編譯產生的 YAML 檔案（legacy）。

## 設計原則

1. **唯讀**：本目錄的檔案由 compiler 自動產生，不應手動編輯。
2. **檔頭標註**：每個 YAML 檔頭包含 `generated_from`、`generated_at`、`compiler_version`、`status`。
3. **狀態標籤**：
   - `synced`：與 prose source 一致
   - `stale`：prose source 已修改但未重新編譯
   - `orphan`：prose source 已不存在
4. **範圍限定**：本目錄只存放**系統層**（`workflow/`、`enforcement/`、`governance/`、`plans/`）的 generated surfaces。
   **領域層**（`analysis/`、`intelligence/`、`feedback/`）的 generated YAML 應放在各自的 source 目錄下，
   不應集中到本目錄。

## 目前狀態

Compiler v1.1.0 已將輸出目標從 YAML 遷移至 SQLite（[`runtime.db`](../runtime.db)）。
本目錄的 YAML 檔案保留作為向後相容，但 compiler 不再更新它們。

**新開發請直接使用 SQLite**：

```bash
# 查詢 phase 定義
sqlite3 runtime/runtime.db "SELECT id, name FROM phases;"

# 查詢 obligation 狀態
sqlite3 runtime/runtime.db "SELECT id, phase, severity FROM obligations WHERE phase = 'checkpoint';"

# 查詢 blocking gates
sqlite3 runtime/runtime.db "SELECT id, name, severity FROM gates WHERE phase = 'execution';"
```

## Legacy YAML 檔案

| 檔案 | 來源 | 狀態 |
|------|------|------|
| `workflow-apk-analysis-phases.yaml` | `workflow/apk-analysis/execution-flow.md` | 🗄️ Legacy (in SQLite) |
| `workflow-apk-analysis-artifacts.yaml` | `workflow/apk-analysis/artifact-gates.md` | 🗄️ Legacy (in SQLite) |
| `workflow-software-delivery-phases.yaml` | `workflow/software-delivery/execution-flow.md` | 🗄️ Legacy (in SQLite) |
| `workflow-software-delivery-artifacts.yaml` | `workflow/software-delivery/artifact-gates.md` | 🗄️ Legacy (in SQLite) |
| `workflow-travel-planning-phases.yaml` | `workflow/travel-planning/execution-flow.md` | 🗄️ Legacy (in SQLite) |
| `workflow-travel-planning-artifacts.yaml` | `workflow/travel-planning/artifact-gates.md` | 🗄️ Legacy (in SQLite) |
| `workflow-documentation-phases.yaml` | `workflow/documentation/execution-flow.md` | 🗄️ Legacy (in SQLite) |
| `transaction-machine.yaml` | `enforcement/dependency-reading.md` | 🗄️ Legacy (in SQLite) |
| `goal-action-gates.yaml` | `enforcement/goal-action-validation.md` | 🗄️ Legacy (in SQLite) |
| `failure-recovery.yaml` | `enforcement/failure-learning-system.md` | 🗄️ Legacy (in SQLite) |
| `language-policy.yaml` | `enforcement/neutral-language.md` | 🗄️ Legacy (in SQLite) |
| `sanitization-rules.yaml` | `enforcement/sanitization.md` | 🗄️ Legacy (in SQLite) |
| `tool-neutrality-rules.yaml` | `enforcement/tool-neutral-documentation.md` | 🗄️ Legacy (in SQLite) |
| `knowledge-update-phases.yaml` | `governance/lifecycle/knowledge-update-flow.md` | 🗄️ Legacy (in SQLite) |
| `classification-rules.yaml` | `governance/lifecycle/knowledge-update-flow.md` + `intelligence/engineering/README.md` | 🗄️ Legacy (in SQLite) |
| `plans-index.yaml` | `plans/active/*.md`（聚合） | 🗄️ Legacy (in SQLite) |
