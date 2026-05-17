> 遵守 [共用規則索引](../../../../enforcement/README.md)、[neutral-language](../../../../enforcement/neutral-language.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Extracted — See [`analysis/apk/tools-and-failures.md`](../../../../analysis/apk/tools-and-failures.md) (去敏規則)

### 2026-05-06 - UI hierarchy content-desc redaction

Status: validated

#### One-line Summary

When parsing UIAutomator XML, treat `text` and `content-desc` as potentially raw user/content data and sanitize before writing docs.

#### Human Explanation

UI hierarchy dumps are useful for route proof, selectors, tab labels, and clickable bounds. They can also contain full visible card titles, search suggestions, comments, chat messages, profile names, and other target-specific or personal text. If an agent copies parsed XML summaries directly into reusable docs or public project docs, it can leak raw content even when API logs are otherwise sanitized.

#### Trigger

Use this lesson when:

- parsing `uiautomator dump` XML;
- summarizing clickable nodes from screenshots/hierarchy;
- documenting search pages, feed cards, comments, chat, profile lists, or media results;
- a command/script prints node `text` or `content-desc` values.

#### Evidence

- Tool: UI hierarchy parsing during read-only APK feature analysis.
- Sanitized excerpt:
  - the XML summary listed search suggestions, ranking rows, and result-card descriptions directly from `content-desc`;
  - those values were not needed for reusable analysis;
  - the durable docs only needed selector geometry, generic block names, tab labels, schema keys, and evidence paths.
- Evidence path: project `capture/*.xml` files only; raw UI text remains in local capture and is not copied into reusable skill docs.

#### Generalized Lesson

Accessibility text is evidence, but it is also content. For reusable APK-analysis documents, extract structure rather than values: screen blocks, selected state, clickable bounds, tab order, operation path, and schema/API mapping. Store raw XML under project-controlled capture locations only.

#### Agent Action

Before writing UI hierarchy findings:

1. Decide whether each `text` / `content-desc` value is a stable UI label or raw content.
2. Keep stable navigation labels only when needed for route proof.
3. Replace result titles, query text, comments, chat messages, profile names, and media captions with generic block names.
4. Write evidence paths and selector/bounds instead of raw values.
5. If using a helper script, make its default output redacted or review the output before copying any text into docs.

#### Applies When

- The analysis uses screenshots, XML hierarchy, UIAutomator, OCR, or accessibility-node dumps.
- Visible text can come from user-generated content, search results, comments, chat, or media cards.
- Findings will be copied into project docs, reusable lessons, or enforcement rules.

#### Does Not Apply When

- The value is a stable app navigation label needed to identify a route, tab, or button.
- The raw value is intentionally retained in private local evidence under the project capture folder and not promoted to reusable docs.

#### Validation

The lesson is validated when:

- docs contain selector/bounds/generic block names rather than raw result/content text;
- evidence paths point to the local capture for traceability;
- reusable skill docs do not include target-specific titles, queries, comments, names, or media text.

#### Promotion Target

- `DOCUMENTATION.md`

#### Required Linked Updates

- Promoted into `DOCUMENTATION.md` feature handoff rules.
- Updated `feedback_history/common/README.md`.
