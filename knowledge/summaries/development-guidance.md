## workflow.software-delivery

| 欄位 | 值 |
| --- | --- |
| Atom ID | `workflow.software-delivery` |
| Source path | `workflow/software-delivery/execution-flow.md` |
| Lifecycle | `validated` |
| Summary | 將授權 App/API/Embedded/Firmware 觀察轉成開發 guidance、實作模式、控制項、檢查清單。software-delivery 入口已切成 focused loading surfaces：`intake.md`（需求/變更/parity/backfill）、`contracts.md`（contract governance/traceability）、`test-strategy.md`（BDD/test strategy）、`validation.md`（validation/performance gate）、`closure.md`（DoR/DoD/lesson feedback）、`surgical-changes.md`（diff purity）。原 `skills/app-development-guidance/` 已刪除，所有內容已遷移至新分層。提供 5 個標準化輸出模板（change-brief / contract / bdd-scenario / implementation-plan / review-report），位於 `workflow/software-delivery/templates/`。另提供 Greenfield 標準化流程（`workflow/greenfield/`）與 Slash Command 模式（`ai-tools/slash-commands.md`）。 |
| When to read | 使用者要求 plan、需求釐清、API/backend/mobile/embedded/security 開發 guidance、實作模式、控制項或檢查清單時；先讀 `execution-flow.md` thin index，再依 task intent 載入對應 focused surface。 |
| Do not use for | 不可取代完整的 `workflow/software-delivery/`、`metadata/development-guidance/` 或 `analysis/development-guidance/` 目錄。不可用於未授權的系統分析。不要預設載入 `examples/EXAMPLES.md`；只有使用者明確要求範例或偵測到 ambiguity 時才載入。 |
| Context cost | ~500 tokens |
| Estimated full cost | ~3500 tokens |
| Validation signal | Pre-build interrogation、requirements、contract、BDD、implementation plan 與 validation target 可追溯；focused loading surfaces 可解析並掛在 `route.workflow.software-delivery`；新分層 entrypoint links 可解析，metadata catalogs 結構完整。舊 `skills/app-development-guidance/` 已刪除。 |
| Last checked | 2026-05-30 |
