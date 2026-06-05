---
id: 2026-06-05-0200-cursor-enforcement-hook-adapter
plan_kind: main
status: draft
owner: linyihong
created: 2026-06-05
parent: null
---

# Cursor Enforcement Hook Adapter（讓 PreToolUse gate 在 Cursor 也機械生效）

**Status**: `draft`
Owner: framework maintainer (linyihong)

## Problem

Ai-skill 的機械 enforcement（bootstrap receipt gate、Phase 5 workflow
`primary_source` gate）目前**只在 Claude Code 生效**，因為它依賴 Claude Code 的
`PreToolUse` hook（`exit≠0` 攔截非-Read 工具）。`init-project --tools cursor`
目前只 wire 了 `sessionStart`（提醒）+ `stop`（close-out，exit 0 + followup，
不能擋）。所以 Cursor 端的 Ai-skill 約束是 **chat/session 層防漏**，沒有
**pre-tool 層攔截** —— detector 偵測到 active_route 後，Cursor 不會強制 agent
先讀 workflow primary_source。

引擎核心（detector / RuntimeContext / router）是 tool-agnostic 的 CLI，Cursor
terminal 直接可用；缺的只是**把 gate 接到 Cursor 的 pre-tool hook**。

## Key Finding（2026-06-05，已查 Cursor hooks 文件）

Cursor **有** pre-execution hooks，且與 Claude Code 的契約高度相似：

| 面向 | Claude Code `PreToolUse` | Cursor `preToolUse` |
|---|---|---|
| 觸發 | 任何工具呼叫前 | 任何工具呼叫前（all tool types） |
| transcript | `transcript_path`（stdin JSON） | `transcript_path`（stdin JSON，可能為 null） |
| 擋的方式 | `exit 2` | `exit 2` **或** `{"permission":"deny","user_message","agent_message"}` |
| fail-open | 其他 exit code 放行 | 其他 exit code 放行（明文 fail-open default） |

Cursor 另有更細的 `beforeShellExecution`（含 `command`）、`beforeMCPExecution`
（含 `tool_name`+`tool_input`）、`beforeReadFile`（含 `file_path`）。

**重大含意**：`workflowPrimarySourceGate` / `finishPreToolUse` 核心邏輯
**可直接重用**；只差 (a) payload 解析（Cursor 欄位）、(b) **block 回應格式**、
(c) `init-project` wiring。比原先預估輕很多。

### ⚠️ 連帶發現 — **CONFIRMED 2026-06-05（high severity 既有 bug）**

我們的 Go hook block 時回 `ExitValidationFailed = 30`，**不是 2**。已查證：

**證據 1 — Claude Code 官方文件**（code.claude.com/docs/en/hooks，PreToolUse）：
> Exit 2 is the **only** exit code that blocks a PreToolUse tool call. Exit 1 or
> other non-zero codes are treated as **non-blocking errors** and do **not**
> prevent the tool from running.

要 block 只有兩條路：`exit 2`，或 `exit 0` + stdout
`{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"..."}}`。

**證據 2 — 我們的 code**：
- bootstrap gate 兩條 block 路徑（`BLOCK_RECEIPT_WITHOUT_READS` /
  `BLOCK_NO_RECEIPT`）與 Phase 5 `finishPreToolUse` 全部：寫 **stderr** +
  `return ExitValidationFailed`（=30）。
- `grep -c permissionDecision hooks.go` = **0**（完全沒有 deny JSON 路徑）。
- `main.go`：`os.Exit(app.Run(...))` —— exit code 直接傳給 Claude Code，無 remap。

**結論（confirmed）**：PreToolUse 回 exit 30 → Claude Code 視為 non-blocking
error → **工具照常執行**。**bootstrap receipt gate 與 Phase 5 workflow gate
從來沒有在 Claude Code 機械 block**，一直是 **behavioral-only**（stderr 以
`<hook> hook error` notice 出現在 transcript，agent 看到後配合）。這與
`ai-tools/agent/CLAUDE.md` L3「以 exit 2 攔截」的描述、以及「mechanical L3 gate」
的定位**矛盾**。

**影響範圍**：所有 PreToolUse block（bootstrap + Phase 5）。Cursor adapter 也
必須回 exit 2 / deny JSON 才會擋。屬獨立 high-severity 既有修正，建議優先於
Cursor wiring。修法：PreToolUse block 路徑由 `return ExitValidationFailed(30)`
改為（a）`return 2`，或（b）`exit 0` + emit `permissionDecision:"deny"` JSON
（後者可帶 reason，且與 Cursor 的 `{permission:deny}` 同構，建議採用）。

> **✅ Claude 端 hotfix 已 landed（2026-06-05，user 選 option B）**：新增
> transport-agnostic `hookDecision{Deny,Reason}` + `renderClaudePreToolUseDecision`
> （`exit 0` + `hookSpecificOutput.permissionDecision="deny"`）。bootstrap receipt
> gate 兩條 block 路徑 + Phase 5 `finishPreToolUse` 全部改用。tests：
> `TestRenderClaudePreToolUseDecision_{Deny,Allow}` +
> `TestPreToolUseHookBlocksReceiptWithoutReads`（改 assert deny JSON）。failure
> pattern `pretooluse-block-wrong-exit-code` + CLAUDE.md L3 已更新。**本 plan 的
> Cursor adapter（Phase 1-3）將直接重用 `hookDecision` 抽象**，只需加 Cursor
> render（`{permission:deny}`）+ payload 解析 + init-project wiring。Q2 即此 hotfix。
>
### ⚠️ 連帶發現 2 — **Stop hook 同類問題（CONFIRMED 2026-06-05）**

查證結果：**Claude Stop path 有完全相同的 bug**。

**證據 1 — Claude 官方文件**（Stop hook）：
> exit 2 = blocks the stop（forces Claude to continue）；**其他非零 = non-blocking
> error → Claude 照常 stop**。JSON 阻擋用 top-level `{"decision":"block","reason":...}`
> （或 `{"continue":false,"stopReason":...}`）。

**證據 2 — code**（`hooks.go`）：`validateStopHookFinalTexts` /
`blockStopHookMissingAssistantText` 在 **Claude path（`cursorStop=false`）** 寫
stderr + `return ExitValidationFailed`（=30）。**Cursor path 已正確**（`writeCursorStopFollowup`
= `followup_message` + exit 0）。

**結論**：Claude 端的 Stop close-out 檢查（Bootstrap Receipt + Cognitive Mode +
Project Git Report）**從未機械 block**，一直 behavioral-only。影響 Cognitive Mode
close-out 與 dirty-repo report 的強制力。**Cursor 端不受影響**（followup 路徑正確）。

**修法（納入本 plan，與 PreToolUse 同抽象）**：Claude Stop block 由 exit 30 改為
`exit 0` + `{"decision":"block","reason":...}`。注意 Stop 用 **`decision:block`**（不是
PreToolUse 的 `permissionDecision:deny`）→ `hookDecision` 的 render 需 per-event：
PreToolUse→permissionDecision、Stop→decision:block、Cursor PreToolUse→permission、
Cursor Stop→followup_message（即 user 講的 Capability ≠ Transport / adapter 分離）。
loop 安全：`runStopHook` 已檢查 `stop_hook_active` 防無限 block，故 decision:block 安全。
**僅 Claude path 需改；Cursor stop 不動。**

## Scope

**IS for**：把既有 PreToolUse gate 邏輯（bootstrap receipt + workflow
primary_source）接到 Cursor 的 `preToolUse` hook，使其在 Cursor 也機械生效；
並修正 block 的 exit code 契約（2 not 30）使 Claude + Cursor 都真正 block。

**NOT for**：重寫 detector / RuntimeContext / router（已完成且 tool-agnostic）；
其他工具（Roo / Codex / Gemini / Copilot）的 enforcement 接入（各自 hook 能力
不同，另案）；Discovery graph-traversal subsystem。

## Phases

### Phase 0 — Verify & Pin Cursor hook API + exit-code 契約（gate）

- [ ] 對**當前 Cursor 版本**實測，pin `preToolUse` 的精確 stdin payload（特別是
      工具識別欄位：是否有 `tool_name`？generic preToolUse vs
      beforeShellExecution 的差異）與 `cursor_version`。
- [ ] 實測 block：`exit 2` 與 `{"permission":"deny"}` 哪個在 Cursor 確實擋下；
      `beforeReadFile` 是否需要單獨放行（gate 對 Read 一律 allow）。
- [x] **實測 Claude Code 是否真的因 `exit 30` 而 block** — **CONFIRMED 2026-06-05：不擋**。
      官方文件明文「only exit 2 blocks PreToolUse；其他非零 = non-blocking」；
      我們的 code 回 exit 30 + stder（無 deny JSON）→ 既有 bootstrap + Phase 5
      gate 一直是 behavioral-only，非機械 block。詳見上方 §連帶發現。Phase 2 一併
      修成 exit 2 / permissionDecision:deny。**仍需做的**：在真實 Claude session
      回放一次以雙重確認（文件 + code 已足以下結論，此為保險）。
- [ ] 產出：pinned API note + exit-code 真相表 → 決定 Phase 2 回應策略。**此 gate
      未過不得進 Phase 2**（避免照錯誤假設寫 adapter）。

### Phase 1 — Payload adapter（Cursor → 既有 gate 輸入）

- [ ] `runPreToolUseHook` 加 caller/payload 辨識：Cursor payload 含
      `hook_event_name` / `cursor_version`，Claude 沒有 → 以此分流。
- [ ] 正規化成既有內部表示（toolName, transcriptPath）。Cursor `preToolUse`
      若無單一 `tool_name`，定義「哪些算 non-Read 工具」的映射（含 Read 放行）。
- [ ] `transcript_path` 為 null 時 fail-open（與現行一致）。
- [ ] 重用 `workflowPrimarySourceGate` + bootstrap gate（**零邏輯改動**）。

### Phase 2 — Response-format adapter + exit-code 契約修正

- [x] **Claude PreToolUse**：已 landed（commit 2ce1189）—— `hookDecision` +
      `renderClaudePreToolUseDecision`（`exit 0` + `permissionDecision:deny`）。
- [ ] **Claude Stop（連帶發現 2，本 plan 新增）**：`validateStopHookFinalTexts` /
      `blockStopHookMissingAssistantText` 的 Claude path 由 `ExitValidationFailed(30)`
      改為 `exit 0` + `{"decision":"block","reason":<合併缺項訊息>}`。Stop 用
      `decision:block`（非 permissionDecision）→ 擴 `hookDecision` render 成
      per-event（`renderClaudeStopDecision`）。**Cursor stop 不動**（followup 已正確）。
      測試：assert Claude stop block 回 `decision:block` JSON（非只 exit code）；
      Cursor stop 回 followup 不變；`stop_hook_active` re-entry 仍 allow（防 loop）。
- [ ] Cursor preToolUse block/allow：`{"permission":"deny"|"allow"}`（依 Phase 0 實測）。
- [ ] 回歸：確認改動不破壞既有 commit-msg / 其他 hook（它們不經此路徑）；
      `TestRunStopHookBlocksCursorPayloadWithoutCognitive` 等既有 stop 測試對齊新契約。

### Phase 3 — init-project Cursor wiring

- [ ] `initProjectCursorHooksContent()` 加 `preToolUse`（或 beforeShellExecution
      + beforeMCPExecution，依 Phase 0）entry → 呼叫 `ai-skill hooks run
      pre-tool-use --repo <ROOT>`，沿用既有 `sh -c` uname 跨平台 wrapper。
- [ ] 保留既有 `sessionStart` + `stop`。
- [ ] `cursor.md` 更新：說明 Cursor 現在有 pre-tool 機械 gate（不再只是提醒）。

### Phase 4 — Tests + e2e

- [ ] Unit：Cursor payload 解析（block / allow / fail-open / Read 放行）。
- [ ] Unit：response 格式（exit 2 vs JSON deny；allow）。
- [ ] e2e：模擬 Cursor `preToolUse` payload 餵 `hooks run pre-tool-use` ——
      locked route + primary_source 未讀 → deny；已讀 → allow；advisory/miss →
      fail-open。對應 Claude e2e 等價案。
- [ ] 重建 5 平台 binary。

### Phase 5 — Docs + close-out

- [ ] `command-contract.md`、`cursor.md`、`ai-tools/README.md` 同步。
- [ ] 若 Phase 0 確立 Claude exit-30 leak：開 failure-pattern
      `pretooluse-block-wrong-exit-code` 或在既有 pattern 補。
- [ ] enforcement-registry：評估是否新增 `cursor_enforcement_parity` rule_class
      或併入既有；coverage 記錄。
- [ ] close-loop：git clean / push / readback。

## Open Questions

| # | Question | 處置 |
|---|---|---|
| Q1 | Cursor `preToolUse` 是否帶單一 `tool_name`，還是要靠 `beforeShellExecution`/`beforeMCPExecution` 分別接？ | Phase 0 實測 pin |
| Q2 | Claude Code 在 `exit 30` 下到底擋不擋？（影響既有 enforcement 是否其實 leak） | **RESOLVED 2026-06-05：不擋（confirmed）**。官方文件 only-exit-2-blocks + 我們回 30 + 無 deny JSON → 既有 bootstrap/Phase 5 gate 實為 behavioral-only。修法 Phase 2（→ exit 2 或 permissionDecision:deny）。屬獨立 high-severity 既有 bug。 |
| Q3 | Cursor `beforeReadFile` 是否需單獨 allow，避免 gate 擋掉 agent 要遵從的 Read？ | Phase 0 + Phase 1 映射 |
| Q4 | block 用 `exit 2` 還是 JSON `{permission:deny}`？兩者 Cursor 都支援，選哪個較穩？ | Phase 0 實測，傾向 JSON（可帶 agent_message 指明要讀哪個檔） |
| Q5 | 是否一個 adapter 同時涵蓋 Roo（VSCode extension，hook 能力不同）？ | out-of-scope，另案 |
| Q6 | Claude Stop hook 是否同樣因 exit 30 而不 block？ | **RESOLVED 2026-06-05：是（confirmed）**。官方文件 Stop 亦 only-exit-2-blocks（或 `decision:block` JSON）；code 的 Claude stop path 回 30 → 不 block，close-out 檢查實為 behavioral-only。Cursor stop（followup_message）不受影響。修法 Phase 2（Claude stop → `decision:block` JSON）。見 §連帶發現 2。 |

## Risks

- **Q2 若成真**：既有 Claude enforcement 可能一直是「顯示但不擋」。本 plan 會
  順手修正，但也意味先前「機械 gate」的安全性需重新評估（behavioral 仍在，
  因為 agent 看得到 stderr/reminder）。
- Cursor hook API 仍在演進 —— Phase 0 pin 版本，docs 註明 `cursor_version`。
- 跨工具 payload 分流若判斷錯誤，可能在某工具誤觸/漏觸 → Phase 4 e2e 雙工具覆蓋。

## Validation Plan

- [ ] Phase 0 API + exit-code 真相已 pin（含 Claude exit-30 實測）
- [ ] Cursor preToolUse payload 解析正確（unit）
- [ ] block 在 Cursor 真正 deny、allow 真正放行（e2e）
- [ ] advisory/miss 在 Cursor fail-open 不誤殺（e2e）
- [ ] Claude 端 enforcement 仍正常（回歸；若修 exit code 則確認真正 block）
- [ ] **Claude Stop block 真正 block（assert `decision:block` JSON，非只 exit code）**；Cursor stop followup 不變；`stop_hook_active` re-entry 仍 allow
- [ ] close-loop：commit / push / readback / clean
