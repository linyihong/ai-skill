# Runtime Routing Philosophy

## Purpose

`governance/lifecycle/routing-philosophy.md`（本檔）與 `runtime/router/` 定義 agent 如何為任務選擇要載入的 Ai-skill knowledge。這是 routing design layer，不是 executable policy。

## Routing Pipeline

```text
task intent
  -> knowledge/indexes/README.md
  -> knowledge/runtime/routing-registry.yaml
  -> metadata/schema.md fields
  -> metadata/ranking + confidence + compatibility
  -> models/profiles + compression strategy
  -> source-of-truth gate
  -> selected primary source
  -> validation and readback
```

## Step 1: Identify Task Intent

載入深層 context 前，先把任務分類成小型 intent：

- Bootstrap / session takeover.
- Skill execution.
- Skill update or promotion.
- Knowledge index / metadata / runtime work.
- Governance, validation, or close-loop work.
- Tool adapter or compatibility work.

## Step 2: Use The Knowledge Index

讀取 `knowledge/indexes/README.md`，找到最接近的 `Task intent` row。若需要結構化路由資料，再讀 `knowledge/runtime/routing-registry.yaml`。

- 先載入 `Primary source`。
- 只有任務需要時才載入 `Related sources`。
- 若沒有符合的 row，fallback 到 root `README.md`、`enforcement/README.md` 與相關 layer README。

### Workflow 編排（activation #27 + registry-first）

當任務命中 activation **#27** 或任一 `route.workflow.*` 的 **`activation_triggers`** 時，在 Step 2 之後 **必須**：

1. 比對 [`routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) 內所有 `route.workflow.*` 的 `activation_triggers`。
2. 若多條命中，讀 [`workflow/workflow-routing.md`](../../workflow/workflow-routing.md) §常見歧義（例如 **software-delivery** vs **apk-analysis**）。
3. 載入選定 route 的 `primary_source` 與 `required_dependencies`。

未完成本 Pipeline **不得**開始寫可觀察產品行為變更。新增 workflow 時只擴充 registry，不擴充 activation-table 專向列。

## Step 3: Apply Metadata Ranking

使用 `metadata/ranking/README.md` 決定讀取順序：

1. Required enforcement rules 與 safety / source / validation gates。
2. 目前 source-of-truth entrypoint。
3. Validated 或 stable atoms。
4. Candidate maps 與 summaries。
5. Optional background references。

## Step 4: Apply Compatibility Gate

使用 `metadata/compatibility/README.md` 與 `governance/lifecycle/README.md`：

- 若舊 `skills/`（已 deprecated）或 `enforcement/` source 仍 active，它優先於 candidate new-layer content。
- 若 new layer path 只是 candidate map，它可引導 discovery，但不可覆蓋 behavior。
- 若已有 promoted atom，依賴它作 replacement 前，先確認 old links 仍可解析。

## Step 5: Validate The Route

依 routed source 行動前，確認：

- Source 是 canonical，不是 tool mirror。
- Required dependencies 已讀，或已標 not applicable。
- Migration 進行中時，old entrypoint 仍被保留。
- 選定 source 提供清楚 validation signal。
- Final response 或 commit 能說明載入了什麼、延後了什麼。

## 與既有文件的關係

- [`runtime/router/`](../../runtime/router/) — Activation 與 router 資料（`activation-table.md`、`activation-rules.yaml`）
- [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md) — Task intent routing table
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) — Machine-readable routing records
- [`metadata/ranking/README.md`](../../metadata/ranking/README.md) — Metadata ranking rules
- [`metadata/confidence/README.md`](../../metadata/confidence/README.md) — Confidence levels
- [`metadata/compatibility/README.md`](../../metadata/compatibility/README.md) — Compatibility gates
- [`models/profiles/README.md`](../../models/profiles/README.md) — Model profiles
- [`governance/lifecycle/README.md`](../../governance/lifecycle/README.md) — Lifecycle states
