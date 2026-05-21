# Execution Cost Strategy

Execution cost strategy balances model capability, context cost, validation cost, and risk.

## Cost Signals

- Context loading size。
- Number of owner layers touched。
- Runtime / generated surface involvement。
- Evidence instability。
- Tool reliability。
- User-facing blast radius。

## Strategy

| Cost / Risk | Strategy |
| --- | --- |
| Low | checklist-first or bounded source-backed. |
| Medium | source-backed + targeted validation. |
| High | validation-heavy + staged edits. |
| Unknown | discovery / human alignment before execution. |

## Rule

Do not reduce validation cost by increasing claim risk. If cost is too high, narrow the claim or ask for alignment.
