> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-08 - Component traceability marker depth

Status: validated

#### One-line Summary

新增或提升共用 UI component 時，BDD / contract 不應只驗 component 名稱或檔案存在，還要驗 feature refs 與最小實作語意 marker。

#### Human Explanation

在文件先行或 BDD 驅動的 UI 開發中，component inventory 很容易退化成「名稱清單」。名稱清單能防止 component 從企劃中消失，但不能證明該 component 的實作真的承擔了預期責任，也不能保證 feature 文件能追到具體檔案。

較可靠的 traceability 應同時覆蓋三層：

1. Feature / contract 明確列出 component 名稱與具體 code path。
2. BDD 測試確認 code path 存在，且相關 component index 有登記。
3. BDD 或 contract test 檢查該 component 的最小語意 marker，例如輸入 ViewModel、主要狀態、互動元件、關閉 / 開啟控制、空態或權限提示。

#### Trigger

- 新增 `src/components/*`、`app/**/_components/*` 或把 route-local UI 上提為共用 component。
- 使用者或 reviewer 問「為什麼 feature / test 沒有明確指到這個 component」。
- BDD 只檢查 component inventory 名稱，沒有讀取 component 實作檔。

#### Evidence

- Tool: agent-assisted software delivery
- Sanitized excerpt: A shared UI component was implemented and indexed, but the feature refs only pointed to a broad components directory; BDD passed because it checked names and existence, not implementation semantics.
- Evidence path: Project-specific evidence remains in `<PROJECT_ROOT>` feature files and BDD tests; this lesson only records the generalized traceability rule.

#### Generalized Lesson

Component traceability is not complete until a future agent can answer:

| Layer | Required signal |
| --- | --- |
| Feature / contract | Component name and concrete source path are referenced. |
| Component index | Shared component location and purpose are documented. |
| BDD / contract test | Test reads the actual component file or a focused fixture and checks semantic markers. |
| Implementation | Component exposes the expected props, state, accessibility label, interaction primitive, or fallback behavior. |

#### Agent Action

When adding or promoting a shared UI component:

1. Update the feature / contract with the component name and exact source path.
2. Update the component index or nearest owner README.
3. Add a BDD / contract assertion that reads the concrete component file.
4. Assert at least one marker for each core responsibility, not just `"ComponentName"`.
5. Run the relevant BDD / contract tests before claiming traceability is complete.

#### Goal / Action / Validation

- Goal: Prevent component inventory tests from passing when implementation traceability is shallow.
- Action: Require feature refs, component index, and implementation semantic markers for shared UI components.
- Validation or reference source: Relevant BDD / contract test reads the component source file and fails if core markers disappear.

#### Applies When

- UI component is shared across routes or domains.
- A route-local component is promoted to a common component folder.
- Feature inventory names component-level obligations.
- Component behavior includes user-visible state, accessibility, permissions, or external contract fields.

#### Does Not Apply When

- A component is private to one route and already covered by a route-level integration test with adequate semantic checks.
- The change is a pure style-only edit that does not alter component responsibility, props, state, or traceability.

#### Validation

The prevention worked when:

- Removing the concrete feature ref to the component path fails a BDD / traceability test.
- Removing a key implementation marker (for example the expected ViewModel prop, dialog primitive, QR/link display, empty-state marker, or auth prompt) fails a focused test.
- The component index still records ownership and common-scope placement.

#### Promotion Target

- `workflow/software-delivery/artifact-gates.md`
- `workflow/software-delivery/execution-flow.md`
- `enforcement/failure-patterns/`

#### Required Linked Updates

- Updated `feedback/history/development-guidance/README.md` to include this lesson count.
- Added a cross-skill failure pattern for shallow component traceability validation.
- Checked reusable guidance boundary: this lesson contains generalized rules only; project-specific file names and live evidence remain outside Ai-skill reusable docs.
