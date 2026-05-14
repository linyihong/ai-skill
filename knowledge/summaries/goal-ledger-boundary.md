# governance.goal-ledger-boundary

| 欄位 | 值 |
| --- | --- |
| Atom ID | `governance.goal-ledger-boundary` |
| Source path | [`../../enforcement/conversation-goal-ledger.md`](../../enforcement/conversation-goal-ledger.md), [`../../enforcement/content-layering.md`](../../enforcement/content-layering.md), [`../../governance/lifecycle/README.md`](../../governance/lifecycle/README.md) |
| Lifecycle | `validated` |
| Summary | `.agent-goals/` 只保存 active conversation goals；長期 roadmap、phase、migration、promotion、deprecation 與治理狀態必須落到 durable planning 文件。 |
| When to read | 建立、完成或刪除 `.agent-goals/` entry，或判斷長期目標應放在 roadmap / governance / layer README 還是 temporary ledger 時。 |
| Do not use for | 不可用 summary 直接刪除 goal；刪除前仍要確認 completion criteria、validation、child goals 與 durable follow-up gate。 |
| Validation signal | `conversation-goal-ledger.md`、`content-layering.md`、`governance/lifecycle/README.md` 均描述 active vs durable 邊界；`.agent-goals/README.md` 不保留 completed archive。 |
| Last checked | 2026-05-11 |

## Checklist

- Active work 才放 `.agent-goals/`。
- 長期狀態放 durable planning 文件。
- 若完成後仍有 next phase 或治理狀態，先回寫 durable 文件。
- 完成條件與 validation 成立後，刪除 active goal 並刷新 `.agent-goals/README.md`。
