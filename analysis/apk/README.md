# APK Analysis Methods

`analysis/apk/` is the candidate home for reusable APK observation and decomposition methods. During the pilot, the active skill remains `skills/apk-analysis/`; this directory provides a reference-first analysis layer without moving existing files.

## 目錄結構

```
analysis/apk/
├── README.md                       # 本文件
├── traffic-triage.md               # 流量分流與路線選擇
├── tools-and-failures.md           # 工具選擇、失敗判讀、命令模板
├── techniques/                     # 保留中：舊 technique 混合層（逐步拆分中）
└── workflows/                      # 操作流程（HOW TO DO）
```

## Scope

This layer owns:

- Traffic and runtime path triage.
- Evidence-first route selection across pcap, MITM, Java hooks, native hooks, Flutter / Dart AOT, local proxy, media, and offline decoding.
- Pattern extraction from dynamic captures into reusable analysis methods.
- Analysis category routing before workflow or engineering conclusions are written.
- **HOW TO DO** execution knowledge: workflow, command, setup, tracing, hook steps, dump methods (`workflows/`).

This layer does not own:

- Tool skill trigger text; keep that in `skills/apk-analysis/SKILL.md`.
- End-to-end agent execution flow; use `workflow/apk-analysis/`.
- **HOW TO THINK** decision intelligence: heuristics, anti-patterns, failure learning, signal detection; use `intelligence/engineering/apk-analysis/`.
- Target-specific API hosts, endpoints, tokens, raw samples, or live run evidence; keep those in project docs.

## Current Source References

| Topic | Current source | Pilot target status |
| --- | --- | --- |
| Common traffic/runtime triage | `../../skills/apk-analysis/WORKFLOW.md` | ✅ extracted to `traffic-triage.md` |
| Tool choice and failure interpretation | `../../skills/apk-analysis/TOOLS.md` | ✅ extracted to `tools-and-failures.md` |
| Media verification tools | `../../skills/apk-analysis/TOOLS.md` | ✅ extracted to `tools-and-failures.md` |
| Automation script safety boundary | `../../skills/apk-analysis/TOOLS.md` | ✅ extracted to `tools-and-failures.md` |
| Flutter / Dart AOT methods | `../../skills/apk-analysis/techniques/flutter-dart-aot/` | ⏳ decomposing: workflow → `workflows/`, intelligence → `intelligence/engineering/apk-analysis/` |
| HTTP API documentation methods | `../../skills/apk-analysis/techniques/http-api/` | keep skill-local, reference here |
| Local proxy / loopback methods | `../../skills/apk-analysis/techniques/local-proxy/` | keep skill-local, reference here |
| Media / HLS methods | `../../skills/apk-analysis/techniques/media-hls/` | keep skill-local, reference here |

## Read Order

1. Start with `../../skills/apk-analysis/SKILL.md` when a tool needs the active skill trigger.
2. Use this file to understand the analysis-layer boundary.
3. Read `traffic-triage.md` for traffic/runtime path triage.
4. Read `tools-and-failures.md` for tool choice, failure interpretation, and command templates.
5. Read `workflows/` for HOW TO DO execution steps after evidence identifies a route.
6. Read `intelligence/engineering/apk-analysis/` for HOW TO THINK decision guidance.

## Migration Notes

- Traffic triage, tool choice, failure interpretation, media verification, and safety boundary have been extracted from `skills/apk-analysis/` into this directory.
- The original skill files still contain the authoritative content; this directory provides a reference-first view.
- When a technique is extracted, create metadata with `../../metadata/schema.md` and update `../../knowledge/indexes/README.md`.
- Preserve links from `skills/apk-analysis/` until the new path has been validated in real use.
- `workflows/` is the new home for HOW TO DO execution knowledge. Old `techniques/` will be gradually decomposed into `workflows/` + `intelligence/` + `techniques-archive/`.
