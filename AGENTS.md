# Codex Adapter Bootstrap

本文件只是 repo-level Codex 自動載入入口，必須遵守 `ai-tools/` 裡的 canonical AI tool 規則。

啟動時，Codex 必須依序讀取：

1. [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) - 最小核心規則。
2. [`README.md`](README.md) - OS layout 入口。
3. [`ai-tools/agent/codex.md`](ai-tools/agent/codex.md) - Codex 專屬 adapter 規則。
4. [`runtime/runtime.db`](runtime/runtime.db) - runtime phase、obligation、gate、output governance 與 lazy-load routing 的 SQLite source-of-truth。

Runtime config canonical source：

- committed runtime config 只存在 `runtime/runtime.db`。
- canonical runtime documents 存在 `runtime_config_documents` 與 projection tables。
- 不要建立或提交 `runtime/**/*.yaml` mirror。
- governance、enforcement、workflow、ai-tools 或 metadata 擁有的 executable contract 必須留在 owner layer；只有 YAML 設定 `runtime_projection.enabled: true` 時才投影到 `runtime.db`。

不要在本文件維護獨立 Codex 規則；請更新 [`ai-tools/agent/codex.md`](ai-tools/agent/codex.md) 與相關 `ai-tools/` 文件。
