# Domain／執行環境基線 — 取數門檻

**Date:** 2026-05-06  
**Category:** common  
**Status:** validated (method) — apply per project  

## Observation

Teams often stop after per-API request/response documentation and schema catalogs. Implementers then cannot reliably attach transports, derive opaque parameters (`l`-like session scalars), or choose pagination semantics without guessing.

## Lesson

Maintain a **project-level domain / runtime baseline** (separate from entity-level Domain Concepts in feature handoff) that records: environment/host family placeholders, TLS/proxy path, login/device dependency for list calls, lineage of opaque query fields, signing/gateway prerequisites (no secrets), pagination ground truth, rate limits. Cross-link rows to API Catalog entries and UI operation ids.

**Finish gate:** if the outcome includes SDK/client/replay/integration, baseline must exist or be an explicit skeleton with tracked open questions in the same work unit.

## Validation

Concrete baseline shape and checklist: [`../../DOCUMENTATION.md`](../../DOCUMENTATION.md) § *Domain／執行環境基線*. Skill entry updated in [`../../SKILL.md`](../../SKILL.md) Quick Start §7.

## Applicability

Any authorized APK traffic analysis whose downstream consumes real HTTP or decrypted JSON outside the APK.
