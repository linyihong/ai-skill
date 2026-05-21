# Multi-model Handoff

Multi-model handoff 定義何時使用 explicit model、subagent 或 specialized runner。這是 optional 機制，且不得造成 unsafe parallel editing。

## Allowed Uses

只有下列條件全都成立時，才使用 model 或 subagent handoff：

1. Task complexity 超過 current execution profile。
2. Handoff 可以 read-only，或與 shared-file mutation 隔離。
3. Target model 或 subagent 確實 available。
4. Handoff 有明確 validation target。
5. Parent agent 仍負責 final source edits、commits 與 user-facing claims。

## Do Not Handoff

下列情況不得 handoff：

- 任務 small、direct，或已經 source-backed。
- Handoff 只會增加 context cost，卻不提高 confidence。
- Subagent 無法存取 required source 或 tools。
- 多個 agents 會修改同一批 files、shared state、git history、migrations、release steps 或 rules。
- 使用者要求 specific model，但該 model unavailable，且沒有 approved fallback。

## Handoff Packet

每次 handoff 必須包含：

```text
Goal:
Source paths:
Assumptions:
已讀 evidence:
Validation target:
禁止 actions:
Expected return:
```

## Return Contract

Parent agent 必須把 subagent output 當作 evidence，而不是 automatic truth。採取行動前需確認：

- Source paths 仍存在。
- Claims 與 source scope 匹配。
- Suggested edits 不違反 owner boundaries。
- Validation 可重現，或明確標記為 not run。
