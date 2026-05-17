# Distributed Runtime Philosophy

## Purpose

為 multi-agent、background jobs、async tasks、delegated execution 等分散式場景提供 **state consistency** 與 **coordination** 機制。

## Design Principles

1. **Transaction Runtime 為基礎**：distributed transaction 共用 `runtime/transactions/transaction-machine.yaml` 的 state machine，擴充 distributed lock 與 lease 機制。
2. **Phase Reconciliation 為核心**：多 agent 間的 phase inconsistency 透過 `runtime/recovery/phase-reconciliation.yaml` 的 reconciliation 程序解決。
3. **Generated Surface 版本控制**：`runtime/runtime.db` 的 `generated_surfaces` 表包含 `updated_at` 時間戳，agent 執行前檢查版本是否最新。
4. **非同步優先**：支援 background job 與 async task 的 lifecycle，不強制所有操作為同步。

## 與既有文件的關係

- [`runtime/distributed/README.md`](../../runtime/distributed/README.md) — Runtime navigation entry point
- [`runtime/distributed/distributed-locks.yaml`](../../runtime/distributed/distributed-locks.yaml) — Distributed lock / lease 機制
- [`runtime/distributed/multi-agent-coordination.yaml`](../../runtime/distributed/multi-agent-coordination.yaml) — Multi-agent 協調規則
- [`runtime/distributed/async-job-lifecycle.yaml`](../../runtime/distributed/async-job-lifecycle.yaml) — Background job / async task lifecycle
- [`runtime/transactions/transaction-machine.yaml`](../../runtime/transactions/transaction-machine.yaml) — Distributed transaction 共用 state machine
- [`runtime/recovery/phase-reconciliation.yaml`](../../runtime/recovery/phase-reconciliation.yaml) — Phase reconciliation 解決多 agent phase inconsistency
