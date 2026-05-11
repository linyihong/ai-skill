# APK Analysis Methods

`analysis/apk/` is the candidate home for reusable APK observation and decomposition methods. During the pilot, the active skill remains `skills/apk-analysis/`; this directory provides a reference-first analysis layer without moving existing files.

## Scope

This layer owns:

- Traffic and runtime path triage.
- Evidence-first route selection across pcap, MITM, Java hooks, native hooks, Flutter / Dart AOT, local proxy, media, and offline decoding.
- Pattern extraction from dynamic captures into reusable analysis methods.
- Analysis category routing before workflow or engineering conclusions are written.

This layer does not own:

- Tool skill trigger text; keep that in `skills/apk-analysis/SKILL.md`.
- End-to-end agent execution flow; use `workflow/apk-analysis/`.
- Stable engineering trade-offs, anti-patterns, and reusable decision lessons; use `intelligence/engineering/apk-analysis/`.
- Target-specific API hosts, endpoints, tokens, raw samples, or live run evidence; keep those in project docs.

## Current Source References

| Topic | Current source | Pilot target status |
| --- | --- | --- |
| Common traffic/runtime triage | `../../skills/apk-analysis/WORKFLOW.md` | candidate for extraction |
| Tool choice and failure interpretation | `../../skills/apk-analysis/TOOLS.md` | candidate for extraction |
| Flutter / Dart AOT methods | `../../skills/apk-analysis/techniques/flutter-dart-aot/` | keep skill-local, reference here |
| HTTP API documentation methods | `../../skills/apk-analysis/techniques/http-api/` | keep skill-local, reference here |
| Local proxy / loopback methods | `../../skills/apk-analysis/techniques/local-proxy/` | keep skill-local, reference here |
| Media / HLS methods | `../../skills/apk-analysis/techniques/media-hls/` | keep skill-local, reference here |

## Read Order

1. Start with `../../skills/apk-analysis/SKILL.md` when a tool needs the active skill trigger.
2. Use this file to understand the analysis-layer boundary.
3. Read `../../skills/apk-analysis/WORKFLOW.md` for the current decision tree.
4. Read only the matching `../../skills/apk-analysis/techniques/<category>/` folder after evidence identifies a route.

## Migration Notes

- Do not duplicate technique content here until a specific Knowledge Atom is selected.
- When a technique is extracted, create metadata with `../../metadata/schema.md` and update `../../knowledge/indexes/README.md`.
- Preserve links from `skills/apk-analysis/` until the new path has been validated in real use.
