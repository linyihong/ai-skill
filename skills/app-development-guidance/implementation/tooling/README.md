# Tooling, CLI, IDE Extension, And Static Analysis Patterns

Use this guidance for developer tools: IDE extensions, linters, static analyzers, CLIs, code generators, import/export tools, migration tools, and internal automation.

## Core Shape

| Layer | Owns | Must not own |
| --- | --- | --- |
| Rule catalog / command catalog | Stable IDs, order, titles, diagnostic or command codes, implementation status. | Product-specific secrets or environment assumptions. |
| Pure kernel | Parsing, validation, transformations, domain rules, command execution over explicit inputs. | Editor APIs, filesystem policy, network calls, UI state, process globals unless injected. |
| Adapter | VS Code/IDE diagnostics, CLI argv/stdout, filesystem discovery, debounce, workspace policy, generated output location. | Hidden rule semantics that bypass the kernel. |
| Fixtures | Input/output pairs, invalid examples, parser failures, edge cases, snapshots. | Private project data unless sanitized. |
| Integration shell | Editor/CLI/E2E tests that prove registration, diagnostics, commands, packaging, or distribution. | Re-testing every pure kernel branch. |

## Development Rules

- Give every rule, command, diagnostic, generator, or migration a stable ID.
- Link ID -> BDD or behavior spec -> implementation refs -> fixtures/tests.
- Keep core logic pure and fixture-testable before adapter/UI tests.
- Keep adapter policy explicit: file selection, debounce, workspace probing, path filters, environment variables, and output handling.
- Use separate diagnostic/error sources or codes for parser/tool failures versus product/domain rule failures.
- If parsing or a third-party tool fails, define degradation behavior: continue string-level checks, emit parse diagnostic, skip rule, or block command.
- Mark process-only, not-tool-enforceable, unimplemented, or intentionally noop catalog entries explicitly.
- Keep helper UX templates, generated snippets, and validators in sync so users can satisfy the tool.

## BDD And Test Closure

| Scenario type | Preferred validation |
| --- | --- |
| Pure rule or command behavior | Fixture-backed unit tests. |
| Parser degradation | Fixture tests for invalid syntax, partial input, unsupported dialect, or malformed files. |
| Adapter behavior | Integration test for IDE diagnostics, CLI output, generated files, exit codes, or package activation. |
| Workspace policy | Tests or documented manual evidence for path filters, project markers, config discovery, and opt-in rules. |
| Process-only rule | Manual review checklist or release gate; mark `not enforceable by tool`. |

## Distribution Gate

- Build/package command is documented.
- Runtime or host compatibility is declared.
- Generated bundles are traceable to source.
- Sample project or fixture proves basic activation.
- Release notes distinguish enforced rules from documented-only guidance.
- Diagnostics, logs, and telemetry avoid leaking source code or private project data beyond the user's environment.

## Required Linked Updates

Follow [`../../../../shared-rules/linked-updates.md`](../../../../shared-rules/linked-updates.md). Tooling implementation changes must update or verify [`../../process/`](../../process/), [`../../CHECKLIST.md`](../../CHECKLIST.md), templates if they add traceability fields, and any relevant language/platform notes.
