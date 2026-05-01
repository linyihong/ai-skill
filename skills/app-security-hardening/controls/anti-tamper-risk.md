# Anti-Tamper And Risk Signals

Use this for root, jailbreak, emulator, hook, tamper, and automation signals.

## Core Guidance

- Treat root, hook, emulator, and tamper detection as risk signals, not sole authorization decisions.
- Critical operations still require backend authorization and abuse controls.
- Design responses that tolerate false positives and false negatives.
- Prefer layered signals: device posture, account behavior, velocity, replay patterns, and transaction context.
- Do not store static bypass targets or permanent secrets in the app.

## Validation Ideas

- Test rooted/emulated/hooked environments and verify product behavior matches risk policy.
- Confirm backend can act on risk signals without trusting the client blindly.
- Review false-positive handling for legitimate users.
- Verify monitoring captures suspicious signal combinations.

## Common Overclaims

- Anti-tamper does not make client code trustworthy.
- Root detection can be bypassed.
- Blocking all suspicious devices may harm legitimate users and still miss automated abuse.
