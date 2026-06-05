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

- [x] 對**當前 Cursor 版本（3.4.17，已實測）** pin `preToolUse` 的 payload + 工具
      識別欄位：**`preToolUse` 帶單一 `tool_name`**（Q1 RESOLVED）。證據來自安裝版
      `workbench.desktop.main.js`（packages/hooks）：matcher-target fn
      `case"preToolUse":case"postToolUse":return e.tool_name`。共用 payload 欄位：
      `hook_event_name` / `cursor_version` / `transcript_path`（subagent 時可能缺）/
      `session_id` / `workspace_roots` / `user_email`。
- [x] 實測 block 協定：**`exit 2` 與 JSON `{"permission":"deny"}` 兩者都擋**（Q4
      RESOLVED，採 JSON）。Cursor `tsn()` 把 exit-2 轉成
      `{permission:"deny",user_message:<reason>}`；hook response 由 **stdout 讀為
      JSON**（非法 JSON → error；exit 124 = timeout）。native 欄位
      `{permission:"deny"|"allow"|"ask", user_message, agent_message, updated_input}`。
      Claude-style `hookSpecificOutput.permissionDecision` 有 compat 轉換
      (`enableClaudeNestedHookSpecificOutputCompatibility`) 但**預設關閉** (`??!1`)
      → adapter 必須用 **native `{permission}`**，不可依賴 compat。
- [x] `beforeReadFile` 是否單獨放行（Q3 RESOLVED）：blockable 事件集
      `kCp=[beforeShellExecution, beforeMCPExecution, beforeReadFile, beforeTabFileRead,
      subagentStart, preToolUse]`。`beforeReadFile` 是**獨立事件**（deny→`CFn` throw）。
      設計結論：**只 wire `preToolUse`→gate，不 wire `beforeReadFile`**（reads 不被
      gate 攔）；gate 內對 read-type `tool_name`（read_file / list_dir / grep /
      glob_file_search / codebase_search）一律 allow，與 Claude 端「Read 放行」對齊。
- [x] **實測 Claude Code 是否真的因 `exit 30` 而 block** — **CONFIRMED 2026-06-05：不擋**。
      官方文件明文「only exit 2 blocks PreToolUse；其他非零 = non-blocking」；
      我們的 code 回 exit 30 + stderr（無 deny JSON）→ 既有 bootstrap + Phase 5
      gate 一直是 behavioral-only，非機械 block。PreToolUse(2ce1189) + Stop(a47e2ee)
      皆已修為 exit 0 + decision JSON。
- [x] 產出：pinned API note + exit-code 真相表（見下方 §Phase 0 Findings）。**Gate
      PASS → 可進 Phase 1。**

#### Phase 0 Findings — PINNED（Cursor 3.4.17，static-from-installed-bundle）

來源：`/Applications/Cursor.app/.../workbench.desktop.main.js`（packages/hooks，
即實際 shipping code，非僅 web docs）。

| 項目 | Cursor 3.4.17 真相 |
|---|---|
| preToolUse 工具識別 | 單一 `tool_name` 欄位 |
| 共用 payload | `hook_event_name`,`cursor_version`,`transcript_path`(可缺),`session_id`,`workspace_roots`,`user_email` |
| block via exit 2 | ✓（`tsn()`→`{permission:"deny",user_message:reason}`） |
| block via JSON | ✓ stdout `{"permission":"deny",user_message,agent_message}` |
| allow | stdout `{"permission":"allow"}`（空 stdout 行為待 Phase 2 確認；非法 JSON=error） |
| Claude compat (`permissionDecision`) | 存在但**預設 OFF** → 不可依賴 |
| blockable events (`kCp`) | beforeShellExecution, beforeMCPExecution, beforeReadFile, beforeTabFileRead, subagentStart, **preToolUse** |
| Read 放行策略 | 不 wire beforeReadFile；preToolUse gate 對 read-type tool_name allow |
| stop 缺項回報 | followup_message（且 normalizer 也接受 `{decision:block,reason}`→轉 followup） |

**唯一 static-residual（不阻擋 Phase 1）**：agent 的 `read_file` 工具究竟走
`preToolUse(tool_name=read_file)` 還是只走 `beforeReadFile`，靜態無法 100% 確認。
但 gate 設計在兩種情形下都安全（只 wire preToolUse + read-allowlist；beforeReadFile
不 wire）。建議 Phase 4 e2e 用真實 Cursor session 回放一次 payload 雙重確認。

### Phase 1 — Payload adapter（Cursor → 既有 gate 輸入）

- [x] `runPreToolUseHook` 加 host 辨識：`detectPreToolUseHost`（payload 有
      `cursor_version` → hostCursor；否則 hostClaude）。**不**用 `hook_event_name`
      分流（Claude payload 也帶該欄位）。
- [x] 正規化沿用既有 `tool_name` / `transcript_path`（Cursor preToolUse 帶單一
      `tool_name`，Q1）。`preToolUseReadAllowed(host,tool)`：Cursor read-allowlist
      （read_file/list_dir/grep/glob_file_search/codebase_search）放行；Claude 維持
      只放行 `Read`（不擴大既有語意）。
- [x] `transcript_path` 為空/不存在時 fail-open（沿用既有 `ALLOW_NO_TRANSCRIPT`）。
- [x] 重用 `workflowPrimarySourceGate` + bootstrap gate（**零邏輯改動**；只在 block
      site 換 `renderPreToolUseDecision(host,...)`，`finishPreToolUse` 加 host 參數）。

### Phase 2 — Response-format adapter + exit-code 契約修正

- [x] **Claude PreToolUse**：已 landed（commit 2ce1189）—— `hookDecision` +
      `renderClaudePreToolUseDecision`（`exit 0` + `permissionDecision:deny`）。
- [x] **Claude Stop（連帶發現 2，本 plan 新增）**：已 landed —— `validateStopHookFinalTexts` /
      `blockStopHookMissingAssistantText`（+ dead 的 `validateStopHookFinalText` singular）
      的 Claude path 由 `ExitValidationFailed(30)` 改為 `exit 0` +
      `{"decision":"block","reason":<合併缺項訊息>}`，經新增的 `renderClaudeStopDecision`。
      Stop 用 `decision:block`（非 permissionDecision）→ `hookDecision` render per-event。
      **Cursor stop 不動**（followup 已正確）。測試：`TestRenderClaudeStopDecision_{Deny,Allow}`
      + `TestRunStopHookBlocksClaude{PayloadWithoutCognitive,OkOnlyPayload}`（assert
      top-level `decision:block` JSON，非只 exit code）；Cursor stop 回 followup 不變；
      `stop_hook_active` re-entry 仍 allow（防 loop）。binary 待 releasebuild 重建。
- [x] Cursor preToolUse block/allow：`renderCursorPreToolUseDecision`（block = stdout
      `{"permission":"deny",user_message,agent_message}`；allow = 空 stdout，與 Claude
      renderer parity，亦為 shipped Cursor stop hook 已證行為）。`renderPreToolUseDecision`
      per-host dispatch。Claude `permissionDecision` compat 預設 OFF，故用 native。
      Tests：`TestRenderCursorPreToolUseDecision_{Deny,Allow}`、`TestDetectPreToolUseHost`、
      `TestPreToolUseReadAllowed`、`TestRunPreToolUseHookCursor{BlocksWithoutReceipt,
      AllowsReadTool,FailOpenWithoutTranscript}`。
- [x] 回歸：`go test ./...` 全綠（既有 Claude PreToolUse / stop 測試不受影響；改動只在
      block-site render + 新增 host 參數）。binary 待 releasebuild 重建。

### Phase 3 — init-project Cursor wiring

- [x] `initProjectCursorHooksContent()` 加 `preToolUse` entry → 呼叫 `hooks run
      pre-tool-use --repo <ROOT>`，沿用 `sh -c` uname wrapper。**關鍵差異**：preToolUse
      wrapper 在 `AI_SKILL_REPO` 無法解析時 **exit 0（fail-open）**，不像 stop 的
      exit 2 + failClosed —— 避免缺 repo 時擋掉所有工具。無 `matcher` 欄位（Cursor 不需）。
- [x] 保留既有 `sessionStart` + `stop`（三 hook：sessionStart/preToolUse/stop）。
- [x] `cursor.md` + `cursor.yaml` 同步：說明 Cursor 現有 pre-tool 機械 gate（native
      `{permission:deny}`、read-allowlist、fail-open-on-missing-repo）；改寫舊「hooks
      只是提醒」框架。cursor.yaml runtime_projection → `runtime compile + refresh` 已跑。
      init_project_test 加 preToolUse wiring + fail-open 斷言。

### Phase 4 — Tests + e2e

- [x] Unit：Cursor payload 解析（block / allow / fail-open / Read 放行）—— 見 Phase 1/2 tests。
- [x] Unit：response 格式（native `{permission:deny}` / allow 空 stdout）—— `TestRenderCursorPreToolUseDecision_*`。
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
| Q1 | Cursor `preToolUse` 是否帶單一 `tool_name`，還是要靠 `beforeShellExecution`/`beforeMCPExecution` 分別接？ | **RESOLVED（3.4.17）：帶單一 `tool_name`**（`case"preToolUse":return e.tool_name`）。只 wire `preToolUse` 即可，不需分接 shell/MCP。 |
| Q2 | Claude Code 在 `exit 30` 下到底擋不擋？（影響既有 enforcement 是否其實 leak） | **RESOLVED 2026-06-05：不擋（confirmed）**。官方文件 only-exit-2-blocks + 我們回 30 + 無 deny JSON → 既有 bootstrap/Phase 5 gate 實為 behavioral-only。修法 Phase 2（→ exit 2 或 permissionDecision:deny）。屬獨立 high-severity 既有 bug。 |
| Q3 | Cursor `beforeReadFile` 是否需單獨 allow，避免 gate 擋掉 agent 要遵從的 Read？ | **RESOLVED（3.4.17）：不 wire `beforeReadFile`**；只 wire `preToolUse` 並在 gate 內對 read-type `tool_name` allow。`beforeReadFile` 是獨立 blockable 事件，不接它即不會擋 read。 |
| Q4 | block 用 `exit 2` 還是 JSON `{permission:deny}`？兩者 Cursor 都支援，選哪個較穩？ | **RESOLVED（3.4.17）：兩者都擋（`tsn()` 把 exit-2 轉 permission:deny）；採 JSON `{permission:"deny",user_message,agent_message}`**（可帶 agent_message 指明要讀哪個檔；response 由 stdout 讀為 JSON）。Claude `permissionDecision` compat 預設 OFF，不可依賴。 |
| Q5 | 是否一個 adapter 同時涵蓋 Roo（VSCode extension，hook 能力不同）？ | out-of-scope，另案 |
| Q6 | Claude Stop hook 是否同樣因 exit 30 而不 block？ | **RESOLVED 2026-06-05：是（confirmed）**。官方文件 Stop 亦 only-exit-2-blocks（或 `decision:block` JSON）；code 的 Claude stop path 回 30 → 不 block，close-out 檢查實為 behavioral-only。Cursor stop（followup_message）不受影響。修法 Phase 2（Claude stop → `decision:block` JSON）。見 §連帶發現 2。 |

## Risks

- **Q2 若成真**：既有 Claude enforcement 可能一直是「顯示但不擋」。本 plan 會
  順手修正，但也意味先前「機械 gate」的安全性需重新評估（behavioral 仍在，
  因為 agent 看得到 stderr/reminder）。
- Cursor hook API 仍在演進 —— Phase 0 pin 版本，docs 註明 `cursor_version`。
- 跨工具 payload 分流若判斷錯誤，可能在某工具誤觸/漏觸 → Phase 4 e2e 雙工具覆蓋。

## Validation Plan

- [x] Phase 0 API + exit-code 真相已 pin（Cursor 3.4.17 實測 + Claude exit-30 confirmed；見 §Phase 0 Findings）
- [ ] Cursor preToolUse payload 解析正確（unit）
- [ ] block 在 Cursor 真正 deny、allow 真正放行（e2e）
- [ ] advisory/miss 在 Cursor fail-open 不誤殺（e2e）
- [ ] Claude 端 enforcement 仍正常（回歸；若修 exit code 則確認真正 block）
- [x] **Claude Stop block 真正 block（assert `decision:block` JSON，非只 exit code）**；Cursor stop followup 不變；`stop_hook_active` re-entry 仍 allow（unit 已綠：`TestRenderClaudeStopDecision_*` + `TestRunStopHookBlocksClaude*`；既有 Cursor stop / loop-guard 測試回歸通過）
- [ ] close-loop：commit / push / readback / clean
