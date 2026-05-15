# Generated Surfaces

本目錄存放由 [`runtime/compiler/compiler-engine.rb`](../compiler/compiler-engine.rb) 從 canonical prose source 編譯產生的 YAML 檔案。

## 設計原則

1. **唯讀**：本目錄的檔案由 compiler 自動產生，不應手動編輯。
2. **檔頭標註**：每個 YAML 檔頭包含 `generated_from`、`generated_at`、`compiler_version`、`status`。
3. **狀態標籤**：
   - `synced`：與 prose source 一致
   - `stale`：prose source 已修改但未重新編譯
   - `orphan`：prose source 已不存在
4. **Pre-commit Hook 保護**：commit 前檢查所有 generated YAML 的 status 是否為 synced。

## 目前狀態

本目錄尚無 generated YAML。Compiler 建置完成後，執行以下命令產生：

```bash
ruby runtime/compiler/compiler-engine.rb
```

## 預計產生的檔案

| 檔案 | 來源 | 狀態 |
|------|------|------|
| `workflow-apk-analysis-phases.yaml` | `workflow/apk-analysis/execution-flow.md` | 待產生 |
| `workflow-apk-analysis-artifacts.yaml` | `workflow/apk-analysis/artifact-gates.md` | 待產生 |
| `workflow-software-delivery-phases.yaml` | `workflow/software-delivery/execution-flow.md` | 待產生 |
| `workflow-software-delivery-artifacts.yaml` | `workflow/software-delivery/artifact-gates.md` | 待產生 |
| `workflow-travel-planning-phases.yaml` | `workflow/travel-planning/execution-flow.md` | 待產生 |
| `workflow-travel-planning-artifacts.yaml` | `workflow/travel-planning/artifact-gates.md` | 待產生 |
| `workflow-documentation-phases.yaml` | `workflow/documentation/execution-flow.md` | 待產生 |
| `transaction-machine.yaml` | `enforcement/dependency-reading.md` | 待產生 |
| `goal-action-gates.yaml` | `enforcement/goal-action-validation.md` | 待產生 |
| `failure-recovery.yaml` | `enforcement/failure-learning-system.md` | 待產生 |
