# Runtime Routing

`runtime/routing/` 定義 agent 如何為任務選擇要載入的 Ai-skill knowledge。這是 routing design layer，不是 executable policy。

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

若任務對應 active `.agent-goals/` entry，以該 goal 作為目前 user-visible objective。不要讓 stale goal 覆蓋最新 user request。

## Step 2: Use The Knowledge Index

讀取 `../../knowledge/indexes/README.md`，找到最接近的 `Task intent` row。若需要結構化路由資料，再讀 `../../knowledge/runtime/routing-registry.yaml`。

- 先載入 `Primary source`。
- 只有任務需要時才載入 `Related sources`。
- 若沒有符合的 row，fallback 到 root `README.md`、`enforcement/README.md` 與相關 layer README。
- 若 row 指向 candidate path，也要載入舊 source-of-truth entrypoint。

## Step 3: Apply Metadata Ranking

使用 `../../metadata/ranking/README.md` 決定讀取順序：

1. Required shared rules 與 safety / source / validation gates。
2. 目前 source-of-truth entrypoint。
3. Validated 或 stable atoms。
4. Candidate maps 與 summaries。
5. Optional background references。

使用 `../../metadata/confidence/README.md`，避免把 low-confidence candidates 當成 stable behavior。

## Step 3.5: Apply Model Profile

使用 `../../models/profiles/README.md` 與 `../../models/compression/README.md` 決定 context loading 深度：

- `small` profile 優先使用 index、registry、summary、checklist，但不能跳過 required dependencies。
- `large` profile 預設使用 source-backed loading，適合跨層規劃、規則更新、migration 與 close-loop 任務。
- `specialized` profile 在 domain / tool / data format 任務中使用，但仍需遵守 shared-rule bootstrap 與 source-of-truth gate。

若任務要修改 canonical source、commit / push / readback、處理 conflict、promotion 或 deprecation，必須升級到 `source-backed` 或 `graph-assisted`。

## Step 4: Apply Compatibility Gate

使用 `../../metadata/compatibility/README.md` 與 `../../governance/lifecycle/README.md`：

- 若舊 `skills/`（已於 2026-05-13 標記為 deprecated）或 `enforcement/` source 仍 active，它優先於 candidate new-layer content。
- 若 new layer path 只是 candidate map，它可引導 discovery，但不可覆蓋 behavior。
- 若已有 promoted atom，依賴它作 replacement 前，先確認 old links 仍可解析。

## Step 5: Validate The Route

依 routed source 行動前，確認：

- Source 是 canonical，不是 tool mirror。
- Required dependencies 已讀，或已標 not applicable。
- Migration 進行中時，old entrypoint 仍被保留。
- 選定 source 提供清楚 validation signal。
- Final response 或 commit 能說明載入了什麼、延後了什麼。

## Runtime Output Shape

重要 routing decision 可簡短回報：

```text
Task intent:
Primary source:
Related sources loaded:
Candidate sources deferred:
Source-of-truth gate:
Validation signal:
```

## Stop Conditions

符合下列情況時停止載入更多 sources：

- Primary source 已回答目前決策。
- 更多 context 只會重複同一 source-of-truth。
- Candidate path 會與 active skill entrypoint 衝突。
- Required validation 失敗。
- User 改變任務優先序。
