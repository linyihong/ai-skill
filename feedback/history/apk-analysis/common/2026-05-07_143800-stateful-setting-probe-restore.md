> Follow [shared rules](../../../../enforcement/README.md), [dependency-reading](../../../../enforcement/dependency-reading.md), [neutral-language](../../../../enforcement/neutral-language.md), [goal-action-validation](../../../../enforcement/goal-action-validation.md), and [feedback-lessons](../../../../enforcement/feedback-lessons.md). This lesson records only generalized guidance.
# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md)

### 2026-05-07 - Stateful Setting Probe Restore

Status: candidate

#### One-line Summary

When a dynamic APK probe changes app settings, restore the original state before ending and record both the probe evidence and restoration evidence.

#### Human Explanation

Some high-value validations require changing in-app state such as language, locale, currency, notification switches, theme, or feature toggles. These probes can answer a provider/lifetime question quickly, but they also leave the device or account in a changed state if not restored. The restoration is part of the evidence workflow, not a courtesy step.

#### Trigger

- A probe taps a settings page or changes persisted app state.
- The changed setting can affect later captures, API shape, UI labels, cache behavior, or user/account state.
- The task is investigative rather than intentionally changing product configuration.

#### Evidence

- Tool: UIAutomator screenshots/XML, Frida entry hooks, value-redacted storage/key-shape logs.
- Sanitized excerpt: a language-setting probe hit the high-semantic save function, wrote the provider setting, and subsequent requests used the expected key shape; the app was then switched back to the original language and a final screenshot was captured.
- Evidence path: project-specific screenshots, XML, and logs stay in `<PROJECT_ROOT>/capture/` and are summarized in `<PROJECT_ROOT>/docs/`.

#### Generalized Lesson

Treat state restoration as a validation step whenever a probe changes app settings. Capture the changed-state proof, then restore the original setting and capture a restoration proof. Keep raw secrets and request values redacted; document timing, key names, selected-index shapes, and UI state instead.

#### Agent Action

Before running a stateful setting probe, identify the likely restoration path. After the probe, restore the original state unless the user explicitly wants the new state to remain. Update docs with both the dynamic evidence and the restore evidence.

#### Goal / Action / Validation

- Goal: Prevent one investigation from contaminating later captures or leaving the test account/device in an unexpected state.
- Action: Change the setting only for the shortest useful window, capture the result, restore it, and capture the restored state.
- Validation or reference source: A valid record includes the setting-change event, the target behavior after the change, and a restored-state screenshot/XML/log marker.

#### Applies When

- Language, locale, currency, theme, feature flags, permissions, account mode, notification switches, or cache-affecting settings are changed.
- The setting can influence subsequent API requests or UI routes.
- The app/device is reused for more analysis.

#### Does Not Apply When

- The user explicitly requested a persistent configuration change.
- The app state is disposable and will be reset immediately after the probe.
- Restoration is impossible; in that case document the blocker and current state.

#### Validation

Confirm the final visible state or stored key shape matches the pre-probe setting, and note any residual risk if only UI evidence is available.

#### Promotion Target

- `WORKFLOW.md`
- `techniques/common/`

#### Required Linked Updates

- Checked `feedback_history/common/README.md` and root `feedback_history/README.md`; both should index this lesson.
- No immediate promotion: keep as `candidate` until repeated across another app setting probe.
