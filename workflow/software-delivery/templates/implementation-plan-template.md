# Implementation Plan: <feature>

## Architecture Compatibility Preflight
- [ ] Candidate files still exist or are marked `not applicable` / `source missing`.
- [ ] Source-of-truth and generated surfaces are identified.
- [ ] Layer responsibility is consistent with current architecture.
- [ ] Compiler / validator impact is known before implementation starts.
- [ ] Open questions reconciled: each is marked `resolved` (with preflight evidence) / `still-open` / `deferred`, and resolved ones are written back to the plan — not left answered only in working notes.

## Pre-build Interrogation
- **Goal**: <what user/system outcome this plan must achieve>
- **Scope**: <affected behavior, files, workflow, runtime surface, tool adapter, or generated artifact>
- **Non-goals**: <explicitly out of scope>
- **Acceptance / validation target**: <test, runtime validate, scenario, query, review, link check>
- **Framework discovery**: <canonical source, owner layer, projection, mirror/cache/generated output, linked updates>
- **Duplication risk**: <none | remove duplicate | deprecate old path | explicit precedence>
- **Open questions**: <blocker_question | safe_assumption | scoped_out | invalidated>
- **Decision**: <proceed | ask_user | revise_plan | blocked>

## Refactor / Replacement Parity Preflight
- [ ] Not applicable, or parity inventory exists before implementation.
- [ ] Legacy surfaces, inputs, outputs / side effects and external dependencies are mapped to target surfaces.
- [ ] Deferred / not planned / tool-specific items include a non-blocking reason.
- [ ] High-risk side effects have dry-run, fake-root, fixture or equivalent validation evidence.

## Task Breakdown
### Task 1: <title>
- **File(s)**: <path>
- **Description**: <what to do>
- **Acceptance**: <how to verify>

### Task 2: <title>
- **File(s)**: <path>
- **Description**: <what to do>
- **Acceptance**: <how to verify>

## Dependencies
- <prerequisite task or external dependency>

## Risk Assessment
| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| <risk> | <H/M/L> | <H/M/L> | <plan> |

## Traceability
- **Upstream**: <link to change-brief | contract | bdd-scenario>
- **Downstream**: → Review Report
