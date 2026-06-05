# PreToolUse Block Wrong Exit Code（hook 想擋卻用了不會擋的 exit code）

Status: validated（2026-06-05 confirmed；既有 bootstrap + Phase 5 gate 受影響）
Class: `meta-governance-gap` / `enforcement-self-leak`

> 這是 [`rule-without-executor`](rule-without-executor.md) meta-pattern 的一個
> 特化：executor **存在**，但它對 host 用了**錯誤的 block 協定**，於是 host 收到
> 後**放行** —— 規則寫了、executor 也寫了，卻沒有真的 enforce。enforcement
> subsystem 自己中了它在抓的 pattern。

## Trigger

- PreToolUse / 任何 pre-execution hook 想擋工具，卻用 `return 1` / `return 30` /
  其他「非 2 的 non-zero exit code」表示 block。
- hook 把 block 訊息寫進 **stderr** 後回非零，以為 host 會擋。
- 文件宣稱「以 exit 2 攔截 / mechanical gate」，但 code 實際回的是別的 code。
- 「為什麼 agent 還是跑了那個工具？」反覆出現但被當成偶發。

## Failure Mode

Host（Claude Code / Cursor）對 PreToolUse 的 block 協定是**精確**的：

**block 協定還是 per-event 的** —— 同一 host 不同 event 用不同 JSON 欄位：

| Host / event | 真正會擋 | 其他 non-zero |
|---|---|---|
| Claude PreToolUse | `exit 2`，或 `exit 0` + `hookSpecificOutput.permissionDecision="deny"` | **non-blocking → 工具照常執行** |
| Claude Stop | `exit 2`，或 `exit 0` + top-level `{"decision":"block","reason":...}` | **non-blocking → Claude 照常 stop** |
| Cursor preToolUse | `exit 2`，或 `{"permission":"deny"}` | **fail-open → 工具照常執行** |
| Cursor stop | `{"followup_message":...}` + exit 0 | （stop 無 exit-2 loop；缺項用 followup） |

⚠️ 用對 event 但用錯**欄位**也會失效（如 Stop 誤用 PreToolUse 的
`permissionDecision` 而非 `decision:block`）。本 repo 同時中過 PreToolUse 與
Claude Stop 兩處（皆回 `exit 30`）。

若 executor 回 `exit 30`（如本 repo 曾用的 `ExitValidationFailed`）：host 把它當
**non-blocking error**，stderr 只當成 `<hook> hook error` notice 顯示，**工具不會
被攔**。結果：

1. **名實不符**：文件/治理標榜「mechanical L3 gate」，實際只是 behavioral
   notice（agent 看得到 stderr → 合作型 agent 配合，但**非機械硬擋**）。
2. **靜默**：因為合作型 agent 通常會配合，看起來「有效」，leak 被長期掩蓋。
3. **可被繞**：不配合 / 被 jailbreak 的 agent 可直接無視。

## Risk

- 所有標榜 mechanical 的 PreToolUse gate（bootstrap receipt、workflow
  primary_source）其實沒在硬擋。
- 跨 host 移植時加倍危險：同一個錯 code 在 Claude 與 Cursor 都 fail-open。
- 治理 coverage 報告把這類 rule_class 標為 `mechanical`，給出**虛假的安全感**。

## Required Agent Action

### 寫 / 改任何 pre-execution hook 的 block 路徑時

1. 用 host 真正認得的 block 協定，**不要**用任意 non-zero exit code：
   - 採 transport-agnostic 的 `hookDecision{Deny, Reason}`，由 adapter render：
     - Claude：`exit 0` + `hookSpecificOutput.permissionDecision="deny"`（帶 `permissionDecisionReason`）。
     - Cursor：`exit 0` + `{"permission":"deny",...}`（或 `exit 2`）。
2. block 的 agent-facing 理由放進 **decision JSON 的 reason 欄位**，不要只放 stderr
   （exit 0 時 stderr 不一定餵回 agent）。
3. 加會驗「真的回 deny」的測試：assert stdout 的 `permissionDecision=="deny"`，
   **不要只 assert exit code**（exit code 對也可能沒擋）。

### Review 時

- 看到 hook 回 `exit 1/30/其他` 當 block → 紅旗。
- 看到「mechanical gate」字樣 → 確認對應 executor 真的用 deny 協定且有 deny-JSON 測試。

## Fix Landed（2026-06-05）

`scripts/ai-skill-cli/internal/app/hooks.go`：新增 `hookDecision` +
`renderClaudePreToolUseDecision`（`exit 0` + deny JSON）。bootstrap receipt gate
兩條 block 路徑與 Phase 5 `finishPreToolUse` 全部改用它（原本回
`ExitValidationFailed=30`）。測試：`TestRenderClaudePreToolUseDecision_{Deny,Allow}`、
`TestPreToolUseHookBlocksReceiptWithoutReads`（改 assert deny JSON）。`ai-tools/
agent/CLAUDE.md` L3 描述同步修正。Cursor adapter（plan 2026-06-05-0200）將重用同
一 `hookDecision` 抽象。

**仍待修（同 plan Phase 2）**：Claude **Stop** hook 的 block path 同樣回 `exit 30`
（`validateStopHookFinalTexts` / `blockStopHookMissingAssistantText` 的 Claude
path），須改為 `exit 0` + `{"decision":"block",...}`（Stop 用 `decision:block`，
非 PreToolUse 的 `permissionDecision`）。Cursor stop（followup_message）已正確。

## Cross-References

- [`rule-without-executor`](rule-without-executor.md) — 上層 meta-pattern
- [`scripts/ai-skill-cli/internal/app/hooks.go`](../../scripts/ai-skill-cli/internal/app/hooks.go) — `renderClaudePreToolUseDecision`
- [`plans/active/2026-06-05-0200-cursor-enforcement-hook-adapter.md`](../../plans/active/2026-06-05-0200-cursor-enforcement-hook-adapter.md) — §連帶發現（confirmed）+ HookDecision 抽象
- Claude Code hooks docs：PreToolUse exit-code / permissionDecision 契約

← [Back to failure-patterns index](README.md)
