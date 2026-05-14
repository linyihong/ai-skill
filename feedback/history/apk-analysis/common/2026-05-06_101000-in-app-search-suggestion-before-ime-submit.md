> 遵守 [共用規則索引](../../../../enforcement/README.md)、[neutral-language](../../../../enforcement/neutral-language.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Extracted — See [`workflow/apk-analysis/execution-flow.md`](../../../../workflow/apk-analysis/execution-flow.md) (Section 1: Capture Window 詳細規則)

### 2026-05-06 - In-app search suggestions before IME submit

Status: validated

#### One-line Summary

When an app search box can trigger system search, permission UI, or another app, validate search flows first through in-app suggestion chips and result-category tabs.

#### Human Explanation

Search pages often combine an app-owned search UI with platform text input behavior. Automated `input text` plus a submit tap can accidentally activate an IME action, system search, or permission screen, which contaminates UI/API attribution. If the search page exposes in-app suggestions, hot terms, or result category tabs, those controls can exercise the app search logic without leaving the app.

#### Trigger

Use this when authorized dynamic analysis shows:

- tapping a search submit button or IME action leaves the target app;
- the search page has in-app suggestion chips, hot list rows, or category tabs;
- API evidence is needed for search without collecting raw query text.

#### Evidence

- Tool: adb UIAutomator hierarchy, screenshots, low-overhead Frida Dart AOT hooks.
- Sanitized excerpt:
  - IME-style search submit opened an external/system UI and was excluded from app evidence;
  - tapping an in-app suggestion stayed inside the target app;
  - switching the app-owned result category tab triggered the feature-specific search logic and API request/decrypt hooks.
- Evidence path: `<PROJECT_ROOT>/capture/` UI hierarchy and Frida logs only; reusable docs store selector/bounds, request key names, schema shapes, and hashes, not raw query text or private content.

#### Generalized Lesson

For search-flow attribution, prefer this order:

1. open the app search page;
2. capture the `EditText`, search button, suggestion chips, and result category tabs;
3. tap an app-owned suggestion or category tab before using IME submit;
4. verify the package remains the target app;
5. correlate feature-specific search logic hooks with request keys and decrypted response shape.

Only use typed input / IME submit after the in-app path is understood or when the app has no suggestion/category controls.

#### Agent Action

Next time this symptom appears:

1. Mark the external/system transition as excluded app evidence.
2. Re-run search through in-app suggestions or result tabs.
3. Log only sanitized query key names and service hashes.
4. Document whether the sampled query produced visible results, empty state, or recommendations.

#### Applies When

- The app exposes in-app search suggestions, hot words, or category tabs.
- UI automation can identify stable bounds/selectors for those controls.
- The analysis scope allows read-only search navigation.

#### Does Not Apply When

- The app has no in-app suggestion/result controls.
- The task specifically requires testing typed query submission or IME behavior.
- Tapping suggestions performs a write action, purchase, follow, like, or other non-read-only behavior.

#### Validation

The lesson is validated when:

- package remains the target app after the in-app search action;
- feature-specific search hooks or request keys fire;
- decrypted response schema can be correlated with the UI action window;
- raw query/result content is not copied into reusable documentation.

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`

#### Required Linked Updates

- Updated `feedback_history/common/README.md`.
- Project-specific selectors, bounds, service hashes, and schema shapes belong in project docs.
