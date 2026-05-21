# Support Matrix：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## Platform Support

| Platform | Desktop single binary | Git dependency | Runtime DB | Mobile posture |
| --- | --- | --- | --- | --- |
| Windows | supported target | external Git required | pure Go SQLite target | not applicable |
| macOS | supported target | external Git required | pure Go SQLite target | not applicable |
| Linux | supported target | external Git required | pure Go SQLite target | not applicable |
| iOS | not native arbitrary binary | only through app / remote | Browser/WASM or app-contained only | control plane / inspect UI / remote trigger |
| Android | feasibility target | Termux / app / remote | Termux or app-contained | app sandbox / remote runner |

## Desktop Baseline

Desktop support targets:

- Windows
- macOS
- Linux

Baseline assumptions:

- `ai-skill` is a single Go binary for the runtime toolchain.
- YAML, JSON, SQLite engine, runtime logic, scheduler, migration / repair logic should be bundled where feasible.
- Git is not bundled. Git remains a required external dependency for writeback, commit, push, hooks, and close-loop commands.

## iOS Boundary

iOS is not a native arbitrary binary target.

Supported evaluation routes:

| Route | Position | Notes |
| --- | --- | --- |
| App-contained runtime | possible | iOS app must bundle runtime, Git / file access, SQLite, and UI |
| Browser/WASM | possible | best candidate for inspect UI, replay UI, state validation, governance control plane |
| SSH remote runner | high feasibility | iPhone acts as control plane; runtime executes on desktop, VPS, NAS, Mac mini, or Linux host |
| Native arbitrary binary | unsupported | iOS security model does not allow general executable persistence |

## Android Boundary

Android feasibility must be evaluated separately:

- Termux may support more local execution than iOS.
- App sandbox has constraints similar in category but not identical to iOS.
- Remote runner remains the conservative option.

## Unsupported / Blocked Conditions

| Condition | Required Behavior |
| --- | --- |
| Missing Git on desktop close-loop | block and prompt install |
| iOS native binary request | reject and suggest App / Browser-WASM / SSH remote runner |
| Missing write permission | block write commands |
| Merge / rebase / cherry-pick state | block commit / push |
| Unsupported platform | return stable `unsupported_platform` exit code |
