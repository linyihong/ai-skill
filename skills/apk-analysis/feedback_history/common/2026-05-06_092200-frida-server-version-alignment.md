> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[neutral-language](../../../../shared-rules/neutral-language.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Extracted — See [`analysis/apk/tools-and-failures.md`](../../../../analysis/apk/tools-and-failures.md) (Frida 健康檢查)

### 2026-05-06 - Frida server version alignment before attach debugging

Status: validated

#### One-line Summary

When `frida-ps` can list Android processes but every attach fails with `remote frida-server: closed`, verify and align the device `frida-server` version before changing hooks or target offsets.

#### Human Explanation

Frida process listing can succeed even when the server-side transport is not healthy enough for attach/spawn. In that state it is easy to misread attach failures as target app protection, bad hook code, or invalid native offsets. A fast version check and minimal attach to both the target process and a benign process separates transport health from hook correctness.

#### Trigger

Use this when an authorized Android APK analysis session shows:

- `frida-ps -U` lists device processes.
- `frida -U -p <pid>` or `frida -U -n <name>` fails with `unable to connect to remote frida-server: closed`.
- The same error occurs for a benign process, not just the target APK.

#### Evidence

- Tool: Frida CLI + Android shell.
- Sanitized excerpt:
  - local Frida tools reported one patch version;
  - device `/data/local/tmp/frida-server --version` reported a different patch version;
  - minimal attach failed before replacing the server;
  - after pushing a matching `frida-server` for the device ABI and starting it as root, minimal attach printed `HOOK_LOADED` for both target and benign processes.
- Evidence path: project-local capture and terminal logs only; no reusable file stores target package names, paths, tokens, or raw traffic.

#### Generalized Lesson

Treat Frida attach as a transport health check before treating it as hook validation. If `frida-ps` works but attach closes the connection, do not immediately edit hook offsets or switch to Gadget/repackaging. First check:

1. local `frida --version`;
2. device `frida-server --version`;
3. device ABI;
4. whether a server process is actually running;
5. minimal attach to a benign process.

If versions differ, install and start a matching server build for the device ABI, then re-run minimal attach.

#### Agent Action

Next time this symptom appears:

1. Run `frida-ps -U`, `frida --version`, and device-side `frida-server --version`.
2. Test minimal attach to the target PID and one benign PID.
3. If both fail with `remote frida-server: closed`, replace/restart the device server with the matching local version before changing hook scripts.
4. Only move to Gadget/repackaging after matching server + benign attach still fails.

#### Applies When

- The device is in authorized analysis scope.
- The device permits replacing or restarting `frida-server` through root, test harness, or approved device management.
- Attach fails globally, not only for one protected process.

#### Does Not Apply When

- The device is not authorized for dynamic instrumentation.
- The target requires Gadget because there is no root/server channel.
- Only one target process fails while benign process attach works; then investigate target-specific controls or timing.

#### Validation

The lesson is validated when:

- local and device Frida versions match;
- `frida-ps -U` still lists processes;
- minimal attach prints a known marker for both a benign PID and the target PID;
- the full hook script can attach and emit initialization lines.

#### Promotion Target

- `TOOLS.md`

#### Required Linked Updates

- Updated `feedback_history/common/README.md`.
- Not promoted to `TOOLS.md` yet because the existing Frida health-check section already covers minimal attach; this lesson adds a narrower version-alignment failure mode that can be promoted after it recurs.
