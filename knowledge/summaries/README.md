# Knowledge Summaries

`knowledge/summaries/` will hold compact summaries of Knowledge Atoms and source-of-truth documents. During the current phase, this directory defines summary format only; it does not replace old skills or shared rules.

## Summary Purpose

Summaries help agents:

- Decide whether a source is relevant before reading the full file.
- Reduce context loading cost.
- Preserve source-of-truth links.
- Support small-model or checklist-first routing.

## Summary Format

Use this shape for future summaries:

| Field | Required | Purpose |
| --- | --- | --- |
| `Atom ID` | yes | Metadata ID from `metadata/schema.md`. |
| `Source path` | yes | Canonical repository-relative source path. |
| `Lifecycle` | yes | `candidate`, `validated`, `stable`, or `deprecated`. |
| `Summary` | yes | One or two sentences describing the source. |
| `When to read` | yes | Trigger condition for loading the full source. |
| `Do not use for` | yes | Boundaries and non-goals. |
| `Validation signal` | yes | How to confirm the summary is still aligned with source. |
| `Last checked` | optional | Date or commit if a summary becomes stable. |

## Example

```markdown
## knowledge.indexes.task-routing

| Field | Value |
| --- | --- |
| Source path | `knowledge/indexes/README.md` |
| Lifecycle | `candidate` |
| Summary | Routes task intents to canonical primary sources and related references. |
| When to read | Use before loading deep skill or shared-rule context. |
| Do not use for | Replacing required dependency reading or old skill entrypoints. |
| Validation signal | Links resolve and routing rows still point to canonical sources. |
```

## Rules

- A summary must link to its source path.
- A summary must not contain secrets, raw evidence, private hosts, tokens, or local absolute paths.
- A summary cannot promote a candidate path into a replacement path.
- If the source changes materially, revalidate or downgrade the summary confidence.
