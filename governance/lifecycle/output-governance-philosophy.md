# Output Governance Philosophy

## Purpose

將語言偏好、文件輸出規則從分散的 prose 檔案升級為 **declarative YAML**，讓 runtime 可以直接檢查、強制執行、並在 validation phase 自動驗證輸出品質。

## Design Principles

1. **Canonical Source 仍在 prose**：`enforcement/neutral-language.md`、`enforcement/sanitization.md`、`enforcement/tool-neutral-documentation.md` 仍為原始規則定義。本目錄的 YAML 為 compiled version，由 compiler 同步更新。
2. **Phase-aware 檢查**：Output governance gates 掛在 `validation` 與 `finalize` phase，確保每輪輸出都經過語言/格式/去敏檢查。
3. **工具中立**：語言政策定義核心規則，各工具的具體設定方式（Roo Code SQLite、Cursor `.cursor/rules/`、Claude `CLAUDE.md`）留在 `ai-tools/agent/*.md`。
4. **Compiler 整合**：compiler 在編譯 generated YAML 時同時檢查 output rules。

## 與既有文件的關係

- [`runtime/output-governance/`](../../runtime/output-governance/) — Runtime navigation entry point (data files: `language-policy.yaml`, `output-rules.yaml`, `governance-gates.yaml`)
- [`runtime/output-governance/language-policy.yaml`](../../runtime/output-governance/language-policy.yaml) — 語言強制規則
- [`runtime/output-governance/output-rules.yaml`](../../runtime/output-governance/output-rules.yaml) — 文件輸出規則
- [`runtime/output-governance/governance-gates.yaml`](../../runtime/output-governance/governance-gates.yaml) — Output governance blocking gates
- [`runtime/runtime.db`](../../runtime/runtime.db) — `phase_machine` / `blocking_gates` / `governance_gates` compiled runtime surface
- [`runtime/compiler/compiler-rules.yaml`](../../runtime/compiler/compiler-rules.yaml) — validation / finalize phase 與 blocking gates 的 embedded source
- [`enforcement/neutral-language.md`](../../enforcement/neutral-language.md) — 語言規則的 prose source
- [`enforcement/sanitization.md`](../../enforcement/sanitization.md) — 去敏規則的 prose source
- [`enforcement/tool-neutral-documentation.md`](../../enforcement/tool-neutral-documentation.md) — 工具中立性規則的 prose source
