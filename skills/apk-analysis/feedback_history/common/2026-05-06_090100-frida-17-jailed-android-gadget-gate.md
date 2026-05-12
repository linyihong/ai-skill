# Extracted — See [`analysis/apk/tools-and-failures.md`](../../../../analysis/apk/tools-and-failures.md) (Frida 健康檢查) and [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-06 - Frida 17 / jailed Android Gadget gate

Status: validated

#### One-line Summary

Frida 17 may reject legacy `--no-pause`, and jailed Android spawn/attach can fail before any hook logic runs unless Gadget or a reachable frida-server is available.

#### Human Explanation

When a capture runner upgrades to Frida 17, `frida -U -f <pkg> -l hook.js --no-pause` can fail immediately with `unrecognized arguments: --no-pause`; Frida 17 uses `--pause` as the opt-in paused mode and resumes by default. After fixing the CLI flag, a non-root / jailed Android environment may still fail with `need Gadget to attach on jailed Android` or `unable to connect to remote frida-server: closed`. In that state, changing Dart offsets or UI actions will not create logs, because injection did not happen.

#### Trigger

- Frida exits before hook initialization with `unrecognized arguments: --no-pause`.
- Frida spawn reports `need Gadget to attach on jailed Android`.
- Frontmost or process-name attach reports frida-server connection closed.
- UI automation succeeds, but Frida log contains only CLI banner or injection failure.

#### Evidence

During an authorized Flutter/Dart AOT capture, the first run failed on `--no-pause`; the corrected Frida 17 command then failed at spawn with a jailed Android Gadget requirement, and frontmost attach failed because the remote frida-server connection closed. UIAutomator evidence remained valid, but same-window Dart hook evidence could not be produced.

#### Generalized Lesson

Separate "Frida CLI compatibility", "device injection transport", and "hook correctness". Only debug offsets after confirming the current Frida command line is accepted and the device can actually spawn/attach.

#### Agent Action

For Frida 17+, avoid adding `--no-pause`; use default resume behavior unless `--pause` is explicitly needed. On jailed Android, first verify Gadget/frida-server availability with a minimal attach before running full hooks. If injection is unavailable, continue UI/API documentation with explicit confidence labels and do not claim same-window hook evidence.

#### Promotion Target

- `TOOLS.md`
- `WORKFLOW.md`
