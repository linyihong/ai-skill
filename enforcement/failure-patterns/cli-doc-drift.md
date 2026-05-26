# CLI Doc Drift（Go script 改了但 command-contract.md 沒同步）

Status: validated
Class: `source-of-truth-duplication` / `governance-drift`

## Trigger

當 commit 在 `scripts/ai-skill-cli/internal/app/*.go` 新增 / 改寫 / 重命名 CLI subcommand dispatch（`case "run <X>":`、`case "<X>":` for runtime subcommand）或 hook handler（`runXxxHook`、`buildRuntimeXxxResult`）但**沒同 commit stage** `scripts/ai-skill-cli/docs/command-contract.md`，使用此 pattern。

具體訊號：

- Staged Go file 含 `+case "run <X>":` diff line
- Staged Go file 含 `+func runXxxHook` 或 `+func buildRuntimeXxxResult`
- Staged 集合**不**含 `scripts/ai-skill-cli/docs/command-contract.md`

## Failure Mode

CLI implementation 跑得比 doc 快，導致：

1. Doc-vs-impl drift — 使用者讀 doc 看不到實際可用 subcommand
2. Migration map 失效 — `legacy-to-go-migration-map.md` 沒列新 commands
3. Onboarding 漏洞 — 新 agent 不知有 `runtime obligations` 等指令

本 session 真實案例：`hooks run commit-msg` / `hooks run pre-push` / `runtime obligations` 三個 subcommand 上線後沒在 command-contract.md，多日 drift 才被使用者發現。

## Required Agent Action

修改 `scripts/ai-skill-cli/internal/app/*.go` 涉及 subcommand dispatch / hook handler 時：

1. 同 commit stage `scripts/ai-skill-cli/docs/command-contract.md`
2. 更新 §初始命令範圍 table + 該 command 的 per-command section
3. 若是純內部 refactor（function body 改動不涉 subcommand 表面）→ 用 `[skip-cli-doc-sync]` 並在 commit message 寫明 reason

## Prevention Gate

Commit-msg validator: `validateCLIDocSync` in `scripts/ai-skill-cli/internal/app/hooks.go`。

Detection logic: git diff `--cached` 含 `+case "run `、`+case "obligations"`、`+func runCommitMsgHook` 等 pattern + command-contract.md 沒 staged → block exit 30。

Canonical contract: [`runtime/cli-modification-policy.yaml`](../../runtime/cli-modification-policy.yaml) §gate.cli.command_contract_synced。

## Validation

符合下列條件即此 pattern 已被防止：

- 每個含 CLI dispatch 變更的 commit 都同時 stage command-contract.md
- 或標示明確 `[skip-cli-doc-sync]` 並有 reason

## Source

- 2026-05-26 session：3 個 subcommand 上線後多日漏 doc；使用者提醒才補 commit `2b106e9`。Phase 6 of [`bootstrap-contract-yaml-migration`](../../plans/archived/2026-05-25-2200-bootstrap-contract-yaml-migration.md) 加入 `validateCLIDocSync` 機械強制。

## Related

- [`runtime/cli-modification-policy.yaml`](../../runtime/cli-modification-policy.yaml) — canonical rule
- [`workflow/software-delivery/execution-flow.md`](../../workflow/software-delivery/execution-flow.md) — parent workflow
- [`markdown-yaml-sync-drift.md`](markdown-yaml-sync-drift.md) — 同類 doc-impl drift

## Linked Validation Scenarios

- `cli-doc-sync-enforcement-v1`

← [Back to failure patterns](README.md)
