# Software Delivery Governance

## Source Intelligence

source_intelligence:

- [`intelligence/engineering/requirements/README.md`](../../intelligence/engineering/requirements/README.md)
- [`intelligence/engineering/development/docs-first-bdd-closure.md`](../../intelligence/engineering/development/docs-first-bdd-closure.md)
- [`analysis/development-guidance/risk-translation.md`](../../analysis/development-guidance/risk-translation.md)
- [`workflow/software-delivery/requirements/README.md`](../../workflow/software-delivery/requirements/README.md)
- [`workflow/software-delivery/development-process.md`](../../workflow/software-delivery/development-process.md)

本文件把 requirements cognition、docs-first BDD closure、contract-first development 與 development guidance 的風險翻譯方法轉譯成 AI runtime software-delivery governance。原始 intelligence 回答「如何穩定 observable behavior、acceptance、traceability 與 validation target」；本文件定義 change intake、requirements cognition、contract precedence、BDD closure、artifact completeness、performance evidence 與 same-session documentation closure 的治理 gate。

## 觸發時機

在下列情況套用本治理：

- 使用者要求 implement、develop、SDK、API、embedded、plan、code review 或 design review。
- 變更可能影響 observable behavior、public contract、domain invariant、API/schema、error handling、storage、安全性、ownership、tests 或 performance。
- Product brief、BDD、contract、implementation 或 tests 之間出現 mismatch。
- 回填已實作專案的 BDD、contract、test plan、hardware/embedded evidence 或 traceability。

## Runtime Gate

| Gate | 通過條件 |
| --- | --- |
| Change intake | 已分類為新需求、bug、refactor、安全/強化、performance 或 planning-only，並確認 code 前需要的 artifacts。 |
| Brief validation | Product brief 的主要 claim 已標記 `validated`、`assumption`、`open question`、`scoped out` 或 `invalidated`。 |
| Requirements cognition | Observable behavior 已有 actor intent、behavior boundary、acceptance criteria、validation target 與 ambiguity disposition。 |
| Contract precedence | 衝突時先判斷 governing contract、product plan、BDD、domain/API/error/hardware contract、implementation、tests 的優先序。 |
| Docs-first BDD closure | Observable behavior 變更前，owning contract、BDD scenario、executable validation 與 implementation slice 已同步或明確 scope out。 |
| Artifact completeness | Change brief、contract、BDD scenario、implementation plan、review report 或 project-local equivalent 已產出或標記 not applicable。 |
| Test strategy | 新行為與舊行為 regression 分開驗證；contract、fixture、integration、property、mutation 或 hardware-in-loop evidence 依風險選擇。 |
| Performance evidence | 影響 latency、throughput、資源、concurrency、external-call fan-out 時，不以功能測試取代 performance budget / smoke / load / stress / spike / soak evidence。 |
| Same-session closure | Code、docs、contracts、BDD、tests、generated clients、fixtures 與 linked updates 在同一批次閉環，或留下明確 owner 與 scoped debt。 |

## 分層判斷

| 內容類型 | 目標層 |
| --- | --- |
| 為什麼 BDD 是 requirements cognition、如何處理 ambiguity / acceptance / traceability | `intelligence/engineering/requirements/` |
| 為什麼 contract 要先於 code、為什麼 BDD 是 behavior bridge | `intelligence/engineering/development/` |
| Software delivery 的 AI runtime gate 與 completion criteria | `governance/ai-runtime-governance/` |
| 實際執行順序、輸出模板選擇、工作流程步驟 | `workflow/software-delivery/` |
| 風險翻譯、owner layer selection、guidance classification | `analysis/development-guidance/` |
| Controls / implementation / platform / language catalog | `metadata/development-guidance/` |

## Workflow Mapping

- [`workflow/software-delivery/requirements/README.md`](../../workflow/software-delivery/requirements/README.md) — requirements cognition stage。
- [`workflow/software-delivery/execution-flow.md`](../../workflow/software-delivery/execution-flow.md) — software-delivery workflow entry and execution order。
- [`workflow/software-delivery/development-process.md`](../../workflow/software-delivery/development-process.md) — contract-first development process and detailed gates。
- [`workflow/software-delivery/artifact-gates.md`](../../workflow/software-delivery/artifact-gates.md) — reusable note structure and artifact quality gates。
- [`analysis/development-guidance/README.md`](../../analysis/development-guidance/README.md) — development guidance analysis methods。

## Validation Candidate

後續若要 promotion 到 `validation/`，可建立 scenario 檢查：

- Requirement / BDD / tests 互相矛盾卻繼續 implementation。
- Acceptance criteria 缺 validation target 卻宣稱 ready。
- Observable behavior change 只改 code，未更新 owning contract / BDD / tests。
- Product brief claim 未驗證卻被當成 implementation input。
- API/schema contract 改變但 generated client、fixtures 或 consumer tests 未同步。
- Performance-sensitive change 只用 unit/functional tests 宣稱可 release。
- Existing project backfill 憑空創造 product intent，而不是標記 `unknown` 或 `open question`。
