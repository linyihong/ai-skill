# Hardening Implementation Patterns

Use this file for short, sanitized implementation patterns before they grow into their own file.

## Server-Owned Authorization

```text
client sends user intent
server loads authenticated subject
server checks object ownership / role / policy
server computes business result
client renders server-approved result
```

Do not let the client send authoritative role, price, balance, ownership, or feature-gate values.

## Redaction Before Export

```text
event created
redaction function removes forbidden fields
safe event sent to logs / analytics / crash reporter
```

Redaction should happen before data leaves the app or server process.

## Required Linked Updates

When adding a pattern here, either promote it to a focused implementation file or update the relevant:

- `controls/`
- `platforms/`
- `languages/`
- `checklists/`
