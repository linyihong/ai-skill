> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/frida-hook-flow.md`](../../../../analysis/apk/workflows/frida-hook-flow.md)

### 2026-05-07 - AOT Hook Crash Static Boundary Fallback

Status: candidate

#### One-line Summary

Flutter / Dart AOT runtime hooks that already produced useful boundary evidence but then crash the app should trigger a static `call_edges` / ASM fallback plus a smaller follow-up hook, not repeated broad runtime attempts.

#### Human Explanation

High-semantic Dart AOT hooks are valuable, but some offsets or return/register summarizers can destabilize Flutter/libapp after the needed request/sign/decrypt boundary has already been observed. Treating the crash as "the flow failed" wastes time and can pollute UI/API attribution. The better move is to preserve the sanitized boundary evidence, mark the runtime window as hook-risk, and finish the missing call-chain facts through static `call_edges`, ASM, string refs, and a narrower validation hook.

#### Trigger

- A Dart AOT hook captures request key sets, signer/decrypt entry/leave shapes, or schema hashes, then the app crashes before the UI operation completes.
- Broad hooks over signer/decrypt/parser paths slow or destabilize startup, retry screens, or Flutter runtime internals.
- The remaining question is "what functions and libraries form this boundary?", not "did the user-facing route load visually?"

#### Evidence

- Tool: Frida native-offset hook, `unflutter` ASM, `call_edges.jsonl`.
- Sanitized excerpt: a lightweight signer/decrypt hook can emit request key names, `serviceHash`, and decrypted `ret/data/msg` JSON shape before a Flutter/libapp crash; static call edges then confirm the signer/decrypt library chain.
- Evidence path: keep crash logs and raw captures in `<PROJECT_ROOT>/capture/`; write only sanitized boundary conclusions in project docs.

#### Generalized Lesson

Once a crashing AOT hook has already answered the boundary question, downgrade it from primary evidence to diagnostic evidence and complete the explanation statically. Runtime crash does not invalidate already-captured sanitized boundary facts, but it does invalidate using that window as UI-success proof.

#### Agent Action

1. Preserve the sanitized events that occurred before the crash.
2. Mark the capture as `boundary evidence + hook-risk`, not a failed endpoint flow.
3. Use `call_edges.jsonl`, ASM, string refs, and function names to identify the downstream library chain.
4. If dynamic validation is still needed, write a smaller hook that only logs the exact entry/exit or hash needed.
5. Do not keep rerunning the broad hook unless the static fallback cannot answer the question.

#### Goal / Action / Validation

- Goal: close signer/decrypt/request boundary questions without destabilizing the app or overclaiming UI success.
- Action: switch from broad runtime hooks to static call-chain analysis once the runtime hook becomes crash-prone.
- Validation or reference source: static call edges agree with observed runtime event ordering and documented hook offsets; follow-up runtime hooks are narrower and no longer required to prove UI completion.

#### Applies When

- Flutter / Dart AOT apps analyzed with native offset hooks.
- The hook output is sanitized and already contains request/decrypt boundary facts.
- The app crash happens after evidence collection, or under known heavy-hook conditions.

#### Does Not Apply When

- The hook crashed before producing any relevant event.
- The question requires user-visible UI completion, purchase/write action confirmation, or server acceptance.
- Static metadata is missing or function attribution is too ambiguous to support the conclusion.

#### Validation

Confirm at least two of the following:

- Runtime events show the expected function offset, key set, length/hash, or schema shape.
- Static `call_edges` show the expected callee chain.
- ASM strings align with the observed phase, such as signer metadata, AES/decrypt errors, query mapping, or timestamp generation.
- A smaller follow-up hook or no-hook control validates the route separately.

#### Promotion Target

- `WORKFLOW.md`
- `techniques/flutter-dart-aot/`

#### Required Linked Updates

- Update the category README and root feedback index.
- Project docs should keep app-specific capture IDs, service hashes, and route conclusions outside this reusable lesson.
