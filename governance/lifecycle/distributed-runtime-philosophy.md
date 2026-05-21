# Distributed Runtime Philosophy

## Purpose

為 multi-agent、background jobs、async tasks、delegated execution 等分散式場景提供 **state consistency** 與 **coordination** 機制。

## Design Principles

1. **Transaction Runtime 為基礎**：distributed transaction 共用 `runtime/runtime.db` 的 `transaction_states` / `transaction_transitions` state machine（source：`runtime/runtime.db`），擴充 distributed lock 與 lease 機制。
2. **Phase Reconciliation 為核心**：多 agent 間的 phase inconsistency 透過 `runtime/runtime.db` 的 `phase_reconciliation` / `recovery_strategies`（source：`runtime/runtime.db`）解決。
3. **Generated Surface 版本控制**：`runtime/runtime.db` 的 `generated_surfaces` 表包含 `updated_at` 時間戳，agent 執行前檢查版本是否最新。
4. **非同步優先**：支援 background job 與 async task 的 lifecycle，不強制所有操作為同步。

## 與既有文件的關係

- [`runtime/distributed/`](../../runtime/distributed/) — Runtime navigation entry point (data files: `distributed-locks.yaml`, `multi-agent-coordination.yaml`, `async-job-lifecycle.yaml`)
- [`runtime/runtime.db`](../../runtime/runtime.db) — Distributed lock / lease 機制
- [`runtime/runtime.db`](../../runtime/runtime.db) — Multi-agent 協調規則
- [`runtime/runtime.db`](../../runtime/runtime.db) — Background job / async task lifecycle
- [`runtime/runtime.db`](../../runtime/runtime.db) — Distributed transaction 與 phase reconciliation 的 compiled runtime surface
- [`runtime/runtime.db`](../../runtime/runtime.db) — transaction / recovery state machine 的 source
