## memory.operations

| 欄位 | 值 |
| --- | --- |
| Atom ID | `memory.operations` |
| Source path | `memory/README.md` |
| Lifecycle | `candidate` |
| Summary | Memory 是 selective replay system：working buffer、summary、episodic、project、failure、decision 與 retrieval-governance。Replay 需要 trigger、qualification、budget、freshness/scope check 與 current source revalidation。 |
| When to read | 需要恢復 session context、查 decision status、參考同 project memory、處理 repeated failure、或判斷 memory 是否可 replay / promotion 時。 |
| Do not use for | 不可取代 canonical source、knowledge navigation、`.agent-goals` active contract 或 runtime execution state。不可 full session replay 除非 explicit on-demand。 |
| Context cost | ~350 tokens |
| Estimated full cost | ~2200 tokens |
| Validation signal | memory README、retrieval-governance、各 subtype README 可解析；validation/scenarios/memory 覆蓋 stale blocker、old goal state、weak guidance、supersession 與 replay budget cases。 |
| Last checked | 2026-05-21 |
