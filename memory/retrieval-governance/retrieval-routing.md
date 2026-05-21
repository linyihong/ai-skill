# Retrieval Routing

Retrieval routing 將目前 trigger signal 對應到候選 memory type。Routing 只產生 candidate context；不產生 truth，也不產生 runtime execution state。

## Routing Table

| Signal | Candidate memory | Qualification |
| --- | --- | --- |
| Repeated failure class | `memory/failure/`、`memory/episodic/` | Failure class 或 execution graph 必須相似；current source 仍需檢查。 |
| Same repo / project | `memory/project/` | 確認 repo、architecture boundary、branch / migration 狀態仍相容。 |
| Architecture decision recall | `memory/decision/` | 檢查 `accepted`、`superseded`、`deprecated` status。 |
| Context compaction recovery | `memory/summary/`、`memory/working/` archive | 只 replay 最小摘要，不 replay full transcript。 |
| Stale assumption suspicion | `memory/failure/`、`memory/episodic/` | 只作 weak hint，必須重新驗證 current source。 |
| Workflow family match | `memory/episodic/`、`memory/project/` | 確認 workflow scope 與 domain boundary。 |
| User asks to continue prior session | `memory/summary/`、active plan、latest git state | Summary 只恢復 context；source truth 仍是 current files。 |

## Routing Order

1. 先讀 current source / active plan / git state。
2. 若 current source 不足，再查 memory index 或 memory type README。
3. 只讀與 trigger 相符的最小 memory。
4. 將 qualified memory 放入 `memory/working/` frame。
5. 使用前重新驗證 current source。

## Discard Rules

立刻 discard memory candidate：

- Scope 不匹配。
- Status 已 deprecated / superseded。
- Freshness 已過期且無 revalidation path。
- 包含 project-secret / private evidence。
- 會造成 current goal drift。
