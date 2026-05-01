> 遵守 [共用規則索引](../../../shared-rules/README.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-01 - Dart AOT Callsite Hooks Are Not Function Hooks

Status: validated

#### One-line Summary

Dart AOT disassembly 裡的 `BL` callsite PC 不能當成穩定 Frida `Interceptor.attach()` 目標；優先 hook function entry，必要時才用 Stalker 或更低階 instrumentation。

#### Human Explanation

`unflutter` 或反組譯輸出的 call edge 會列出 caller 內部 `BL` 指令位址。這些位址很適合做靜態 xref 與流程理解，但它們不是函式入口；Frida `Interceptor.attach()` 在 Android/ARM64 上通常要求可攔截的 function boundary 或可安全替換的目標。把多個 Dart AOT callsite 直接拿去 attach，可能全部被拒絕，或在熱路徑上造成 App 初始化後立即退出。

另一個相近陷阱是直接 hook 全域 Dart collection helper，例如 `LinkedHashMap._set`，再在每次呼叫中解 Dart object。這類 helper 會在啟動與 UI/runtime 路徑高頻命中，容易造成過大 overhead 或不穩定，且會產生大量非目標噪音。

#### Trigger

- 靜態 call edge 已指出某個 Dart interceptor 會呼叫 `Map._set`、`replaceAll` 或其他 helper。
- 想在「某個 caller 呼叫 helper 的瞬間」讀 register 來抓 header/value。
- Frida 對 `base + from_pc` 報 `unable to intercept function ...`，或全域 helper hook 後 App 很快結束。

#### Evidence

- Tool: `unflutter` function/call edge output plus Frida native offset hook.
- Sanitized excerpt: repeated `unable to intercept function at <base+callsite>; please file a bug`; global Dart `LinkedHashMap._set` hook installed but app terminated before useful target events.
- Evidence path: project-private Frida logs under `<PROJECT_ROOT>/capture/` and static call edges under `<PROJECT_ROOT>/work/unflutter_out/`.

#### Generalized Lesson

Use Dart AOT callsites as navigation hints, not as default hook anchors. For dynamic capture, start with function-entry offsets from `functions.jsonl` and decode arguments/returns there. If the needed value exists only at an internal callsite, prefer:

- hook the callee function only when it is narrow enough and not a global hot helper;
- hook the semantic Java/native boundary that receives the final value;
- use Stalker or instruction-level instrumentation as an explicit experiment, with short windows and tight filters;
- document that a callsite attach failed instead of treating it as proof the callsite was not executed.

#### Agent Action

When a Dart AOT call edge points at a promising helper, first check whether the target is a narrow app-owned function or a global runtime/collection/string helper. Hook app-owned function entries first. Avoid defaulting to global helper hooks; if attempted, keep the capture short, sanitize output, and be ready to remove the hook if the process exits or becomes noisy.

#### Applies When

- Flutter/Dart AOT APK analysis uses `unflutter`, `blutter`, IDA, Ghidra, or similar output with function PCs and call edges.
- The analyst wants to correlate signed/encrypted headers, Map writes, string replacements, or response decode helpers.
- Frida native hooks are being placed at `libapp.so` offsets.

#### Does Not Apply When

- The address is a real exported/native function entry or a narrow app-owned Dart function PC from `functions.jsonl`.
- The toolchain supports safe instruction-level tracing for the chosen address and the analyst intentionally uses it.
- A high-level Java/Netty/OkHttp/Dio hook already exposes the final value with enough context.

#### Validation

- Function-entry hooks at app-owned Dart PCs remain stable and produce events.
- Direct callsite attach either fails with Frida's `unable to intercept function` error or produces no reliable target events.
- Global helper hooks are shown to be noisy or destabilizing compared with semantic function/handler hooks.

#### Promotion Target

- `WORKFLOW.md`
- `TOOLS.md`

#### Required Linked Updates

- Add a concise warning to the Flutter/Dart AOT workflow and Frida/tooling failure patterns. No app-specific endpoints, hosts, or header values are included here.
