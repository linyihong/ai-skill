> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Frida Python create_script attach may lack Java bridge — use CLI subprocess for Java RPC

Status: candidate

#### One-line Summary

Attach 模式用 Python `frida` 的 `session.create_script().load()` 時，腳本頂層可能 **`Java is not defined`**；同一 PID 用 **`frida -U -p <pid> -l script.js`** 卻正常。Java 層 RPC（`Java.use` / `Java.choose`）應改 **Frida CLI 子程序** 取 JSON，而非假設 Python binding 與 CLI 等價。

#### Human Explanation

Frida 17 attach + Python API 與 CLI `-l` 的 runtime 初始化順序不同，導致 `Java.perform` 在 script load 時失敗。症狀：`ReferenceError: 'Java' is not defined` at line 1；`rpc.exports` 永遠註冊不了。解法：(1) host 用 `subprocess` 呼叫 `frida -l sign_rpc.js -e "console.log(JSON.stringify(...))"`；(2) script 內 `rpc.exports` 放在 `Java.perform` 內；(3) singleton 用 `Java.choose` 而非 static 假設。驗證：CLI 與 Python 子程序輸出一致後再接離線 signer。

#### Trigger

- `frida.get_usb_device().attach(pid)` + `create_script` with `Java.perform` at load
- Error: `Java is not defined` in script message
- Same script via `frida -U -p PID -l` works

#### Evidence

- Tool: Frida 17.9.6 attach; Python 3.12 `frida` pip vs CLI
- Sanitized excerpt: Python attach fails; CLI returns RPC JSON; hybrid POST succeeds via CLI wrapper
- Evidence path: `<PROJECT_ROOT>/scripts/sign/frida_cli_rpc.py`

#### Generalized Lesson

```text
Java RPC from host:
  1. Try Python frida module — if Java bridge works, keep it
  2. If "Java is not defined" on script load in attach mode:
     → subprocess frida CLI -l <rpc.js> -e "console.log(JSON.stringify(rpc.exports.foo()))"
  3. Do not assume pip frida ≡ CLI for JVM apps
  4. Register rpc.exports inside Java.perform; resolve singletons with Java.choose
```

#### Agent Action

1. Project 提供 `frida_cli_rpc.py` 薄封裝；`FRIDA` env 指向 venv binary。
2. Ai-skill 只寫分支，不寫 package 名。
3. 與 `142000`、`142500` 並用。

#### Goal / Action / Validation

- Goal: host 離線 signer 能穩定取 live native 欄位。
- Action: CLI subprocess JSON bridge。
- Validation: hybrid API POST 200 with offline sign + CLI requestTime。

#### Applies When

- Attach (not spawn) to running JVM app
- Python frida script needs Java.use / Java.choose

#### Does Not Apply When

- Python create_script Java bridge works on target Frida version
- Pure native hook (no Java)

#### Validation

- Documented Python failure + CLI success + E2E sign verify

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §Frida host relay

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
