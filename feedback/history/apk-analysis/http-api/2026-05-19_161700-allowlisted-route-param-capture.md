> 遵守 [共用規則索引](../../../../enforcement/README.md)、[authorization-scope](../../../../enforcement/authorization-scope.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - Allowlisted route parameter capture

Status: candidate

#### One-line Summary

When a live gate only needs route names or locale-like parameters, use an allowlisted Frida hook that logs only those fields instead of enabling full raw query logging.

#### Human Explanation

Gateway requests often put route `service`, locale, tokens, device identifiers, signatures, and user/session fields into the same canonical query string. Turning on full raw logging to recover one route parameter creates unnecessary secret-handling risk.

A safer probe parses the query in memory and prints only the fields needed for the current gate, such as `service`, `l`, request key names, page class, and hash verification.

#### Trigger

- A live SDK/client gate is blocked by an unknown route name, language code, page parameter, or key set.
- Existing captures have only hashes or redacted docs, and the raw value is not security-sensitive by itself.
- The full request string also contains tokens, device ids, cookies, signatures, or other private material.

#### Evidence

- Tool: Frida `_generateEhHeader` / request-normalization hook with field allowlist.
- Sanitized excerpt: hook output should look like `route=<class> keys=<names> service=<route-name> serviceHash=<hash-prefix> l=<language>` and must not include uid/token/device/cookie/eh/full query.
- Evidence path: project capture logs remain under ignored or controlled capture paths; reusable lesson stores only the generalized pattern.

#### Generalized Lesson

Prefer a route-parameter allowlist over raw request dumps:

- Parse the request material in the hook.
- Print only the specific fields needed for the decision.
- Include hash prefixes so values can be matched against sanitized API docs.
- Represent IDs and category values as presence/classes unless their raw value is explicitly in scope.
- Keep full raw query logging disabled by default.

#### Agent Action

1. Define the minimal fields required by the gate.
2. Add or use a narrow hook that extracts only those fields.
3. Validate the observed value against existing serviceHash/key-shape evidence.
4. Update project docs/tests with the derived default or keep the blocker if the hook only proves presence.

#### Applies When

- Recovering service route names, locale/server-language parameters, pagination fields, or request key order from APK traffic.
- Preparing live smoke gates for SDK/private adapters.

#### Does Not Apply When

- The raw field is itself a credential, user identifier, device identifier, signature, decrypt key, or account/session material.
- The task requires byte-for-byte replay parity; then raw material belongs only in a private, access-controlled capture.

#### Validation

- The log contains no full query string and no raw uid/token/device/cookie/eh.
- Hashes match the sanitized API catalog.
- The live gate's missing-requirements list drops only the now-derived parameter.
