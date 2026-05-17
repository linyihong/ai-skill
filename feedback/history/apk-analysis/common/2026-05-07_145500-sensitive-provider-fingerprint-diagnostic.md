> Follow [enforcement rules](../../../../enforcement/README.md), [dependency-reading](../../../../enforcement/dependency-reading.md), [neutral-language](../../../../enforcement/neutral-language.md), [goal-action-validation](../../../../enforcement/goal-action-validation.md), and [feedback-lessons](../../../../enforcement/feedback-lessons.md). This lesson records only generalized guidance.
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-07 - Sensitive Provider Fingerprint Diagnostic

Status: candidate

#### One-line Summary

When a sensitive provider value must be compared across states, use a short disabled-by-default fingerprint or equality diagnostic instead of logging raw values.

#### Human Explanation

Session, language, signing, device, token, and gateway provider values can be required for replay analysis but unsafe to print. A temporary diagnostic can prove whether the provider changes, remains stable, or matches an outgoing request by logging only length/hash or equality booleans. This keeps the analysis useful without turning docs or logs into a source of secrets.

#### Trigger

- A provider has been identified but raw values should not be exposed.
- The question is whether a setting/action changes the provider, or whether two internal values match.
- Existing shape-only logs are insufficient because they only prove key presence.

#### Evidence

- Tool: Frida Java/Dart hooks, storage key hooks, request key-shape hooks.
- Sanitized excerpt: a setting switch wrote a provider key twice; the diagnostic logged distinct redacted fingerprints and subsequent signed requests still contained the provider-backed query key.
- Evidence path: project-specific fingerprints and logs stay in `<PROJECT_ROOT>/capture/`; public docs summarize only changed/same and evidence paths.

#### Generalized Lesson

Use value fingerprints only as a bounded diagnostic, not a default logging mode. Prefer length/hash for "changed vs unchanged" and equality-only comparisons for "provider equals outgoing request value." Do not publish raw values; avoid publishing hashes when the domain is tiny and easily brute-forced unless the project explicitly treats capture logs as private.

#### Agent Action

Add a disabled-by-default feature flag for the diagnostic. Run it for one short validation window, record only redacted fingerprints or equality booleans, restore app state if changed, then disable the flag before committing. In project docs, cite the private capture path and summarize the result without raw values.

#### Goal / Action / Validation

- Goal: Validate sensitive provider lifetime or state change without exposing replay material.
- Action: Log length/hash or equality only under a clearly named experimental flag.
- Validation or reference source: A valid closure includes distinct/same result, subsequent key-shape evidence, restored default hook settings, and no raw provider value in committed docs.

#### Applies When

- Values are session-, token-, language-, signing-, device-, account-, gateway-, or request-provider related.
- Replay analysis needs confidence about stability or equality.
- The diagnostic can be scoped to a single key, function, or short capture window.

#### Does Not Apply When

- Raw value disclosure is explicitly authorized and safe for the target repo.
- The value domain is so small that even hashes should not leave private capture logs.
- A schema/key-presence proof already answers the question.

#### Validation

Confirm the default hook/script has the diagnostic disabled after the run, and run syntax/lint checks on edited scripts.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Checked `feedback_history/common/README.md` and root `feedback_history/README.md`; both should index this lesson.
- No immediate promotion: keep as `candidate` until repeated across another sensitive provider class.
