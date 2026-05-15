# Distributed Runtime

## 目的

為 multi-agent、background jobs、async tasks、delegated execution 等分散式場景提供 **state consistency** 與 **coordination** 機制。

## 設計原則

1. **Transaction Runtime 為基礎**：distributed transaction 共用 `runtime/transactions/transaction-machine.yaml` 的 state machine，擴充 distributed lock 與 lease 機制。
2. **Phase Reconciliation 為核心**：多 agent 間的 phase inconsistency 透過 `runtime/recovery/phase-reconciliation.yaml` 的 reconciliation 程序解決。
3. **Generated Surface 版本控制**：`runtime/generated/*.yaml` 的 header 包含 `generated_at` 時間戳，agent 執行前檢查版本是否最新。
4. **非同步優先**：支援 background job 與 async task 的 lifecycle，不強制所有操作為同步。

## 檔案結構

```
runtime/distributed/
├── README.md                          # 設計原則（本檔）
├── distributed-locks.yaml             # Distributed lock / lease 機制
├── multi-agent-coordination.yaml      # Multi-agent 協調規則
└── async-job-lifecycle.yaml           # Background job / async task lifecycle
```

## 與既有層的關係

| 元件 | 關係 |
|------|------|
| `runtime/transactions/transaction-machine.yaml` | Distributed transaction 共用 state machine |
| `runtime/recovery/phase-reconciliation.yaml` | Phase reconciliation 解決多 agent phase inconsistency |
| `runtime/generated/README.md` | Generated surface 版本控制防止 stale state |
| `runtime/gates/blocking-gates.yaml` | Distributed gates 確保跨 agent 一致性 |
| `runtime/obligations/obligation-ledger.yaml` | Distributed obligation 追蹤 |

## 使用方式

Agent 在 multi-agent 或 async 場景下：

1. 讀取 `distributed-locks.yaml` 取得 distributed lock
2. 讀取 `multi-agent-coordination.yaml` 確認協調規則
3. 執行操作前檢查 generated surface 版本
4. 操作完成後釋放 lock 並更新 phase
