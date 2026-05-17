# Generated Surfaces (Legacy — Migrated to SQLite)

> **✅ Migration Complete**: Compiler v1.1.0 now outputs to [`runtime.db`](../runtime.db) (SQLite).
> The legacy YAML files have been removed. All agents should query `runtime.db` directly via SQLite.

本目錄先前存放由 [`runtime/compiler/compiler-engine.rb`](../compiler/compiler-engine.rb) 從 canonical prose source 編譯產生的 YAML 檔案。
Compiler v1.1.0 已將輸出目標從 YAML 遷移至 SQLite（[`runtime.db`](../runtime.db)），
所有 legacy YAML 檔案已於 2026-05-17 刪除。

## 設計原則

1. **唯讀**：`runtime.db` 由 compiler 自動產生，不應手動編輯。
2. **範圍限定**：`runtime.db` 只存放**系統層**（`workflow/`、`enforcement/`、`governance/`、`plans/`）的 generated surfaces。
   **領域層**（`analysis/`、`intelligence/`、`feedback/`）的 generated surface 應放在各自的 source 目錄下。

## 查詢方式

**新開發請直接使用 SQLite**：

```bash
# 查詢 phase 定義
sqlite3 runtime/runtime.db "SELECT id, name FROM phases;"

# 查詢 obligation 狀態
sqlite3 runtime/runtime.db "SELECT id, phase, severity FROM obligations WHERE phase = 'checkpoint';"

# 查詢 blocking gates
sqlite3 runtime/runtime.db "SELECT id, name, severity FROM gates WHERE phase = 'execution';"

# 查詢 generated surfaces（取代 legacy YAML）
sqlite3 runtime/runtime.db "SELECT type, source, updated_at FROM generated_surfaces;"

# 查詢 runtime config（已編譯至專屬表格）
sqlite3 runtime/runtime.db "SELECT model_name FROM runtime_budget;"
sqlite3 runtime/runtime.db "SELECT ttl_type FROM context_ttl_policy;"
sqlite3 runtime/runtime.db "SELECT guard_name FROM circuit_breaker;"
sqlite3 runtime/runtime.db "SELECT phase_id, content FROM phase_machine;"
sqlite3 runtime/runtime.db "SELECT obligation_id FROM obligation_ledger;"
sqlite3 runtime/runtime.db "SELECT gate_id FROM blocking_gates;"
```

完整表格清單請見 [`../README.md`](../README.md) 的 Databases 章節。
