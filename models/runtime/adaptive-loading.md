# Adaptive Loading

Adaptive loading defines when to move from compressed context to full source.

## Loading Flow

```text
index
→ summary / checklist
→ primary source
→ required dependencies
→ graph-assisted related sources
→ validation
```

## Escalation

Escalate immediately when:

- A file will be edited.
- Source-of-truth is ambiguous.
- Generated surface must be regenerated.
- Memory replay or model output conflicts with current source.
- The user corrects the framing.

## Stop Rule

Stop loading additional context once the validation target can be answered safely. Extra context is not a substitute for validation.
