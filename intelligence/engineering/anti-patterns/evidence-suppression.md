# Evidence Suppression

Anti-pattern: hiding a failure signal instead of resolving or explicitly accepting the underlying failure.

## Pattern

```text
failure signal appears
↓
implementation hides, catches, retries, filters, or masks the signal
↓
validation observes a quieter system
↓
underlying failure remains
```

Evidence suppression is dangerous because it converts a diagnosable failure into missing evidence.

## Common Variants

| Variant | Suppresses | Better Response |
| --- | --- | --- |
| `css_overflow_hidden` | Horizontal or vertical layout overflow. | Identify the overflowing element, source layout contract, and intended scroller boundary. |
| `exception_swallowing` | Runtime error path. | Classify the error, log sufficient context, or handle with explicit fallback. |
| `retry_until_success` | Intermittent failure rate. | Record retry count, terminal failure, and root cause pressure. |
| `warning_filtering` | Static or runtime warning evidence. | Fix, scope, or document the warning with owner and expiration. |
| `test_assertion_weakening` | Regression signal. | Adjust the contract only when behavior changed intentionally and evidence supports it. |

## Rule

- Hiding a signal is allowed only when the hidden condition is intentionally out of scope and the scope is documented.
- Suppression must name the owner, reason, and alternative evidence path.
- Suppression must not be used to support a completion claim.

## UI Example

`overflow: hidden` on a page shell can be valid for a component with intentional clipping. It becomes evidence suppression when it hides document overflow, final-content clipping, fixed-surface overlap, or scroll-root failure without proving the source layout is correct.

For layout authority failures, see [`../ui/layout-authority-governance/`](../ui/layout-authority-governance/README.md).
