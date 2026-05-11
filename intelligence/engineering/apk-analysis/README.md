# APK Analysis Engineering Intelligence

`intelligence/engineering/apk-analysis/` is the candidate home for stable engineering judgment extracted from APK analysis work. During the pilot, lesson history remains in `skills/apk-analysis/feedback_history/`; this directory only maps what should eventually be promoted.

## Scope

This layer owns:

- Stable engineering lessons from repeated APK analysis work.
- Anti-patterns and trade-offs, such as when not to keep adding broad hooks or when API shape is insufficient for live SDK work.
- Reusable decision guidance that is broader than a single technique file.
- Promotion targets for validated lessons that should influence future workflows or runtime routing.

This layer does not own:

- Raw project evidence, hosts, endpoints, tokens, device IDs, or private run logs.
- Step-by-step capture workflow; use `workflow/apk-analysis/`.
- Traffic/runtime observation methods; use `analysis/apk/`.
- Skill-specific lesson archive; keep history in `skills/apk-analysis/feedback_history/` until promotion.

## Candidate Intelligence Areas

| Area | Current source | Future shape |
| --- | --- | --- |
| Evidence-first route selection | `../../../skills/apk-analysis/WORKFLOW.md`, `../../../skills/apk-analysis/feedback_history/common/` | decision guidance for high-leverage analysis routing |
| API catalog and runtime baseline readiness | `../../../skills/apk-analysis/DOCUMENTATION.md`, `../../../skills/apk-analysis/feedback_history/http-api/`, `../../../skills/apk-analysis/feedback_history/common/` | engineering guidance for live SDK / client readiness |
| Flutter / Dart AOT anti-patterns | `../../../skills/apk-analysis/techniques/flutter-dart-aot/`, `../../../skills/apk-analysis/feedback_history/flutter-dart-aot/` | stable failure patterns and hook-selection trade-offs |
| Local proxy and routing ambiguity | `../../../skills/apk-analysis/techniques/local-proxy/`, `../../../skills/apk-analysis/feedback_history/local-proxy/` | guidance for separating routing, proxy, TLS, and attribution |
| Media chain completeness | `../../../skills/apk-analysis/techniques/media-hls/`, `../../../skills/apk-analysis/feedback_history/media-hls/` | control-plane / data-plane and validation heuristics |

## Promotion Rule

A lesson can be promoted here when:

1. It is generalized and sanitized.
2. It has validation evidence or repeated use.
3. It affects engineering judgment, not only one tool command.
4. The original `feedback_history/` entry remains discoverable.
5. `knowledge/indexes/README.md` and metadata are updated if the promoted atom becomes a routing target.

## Compatibility Notes

- Do not delete or rewrite existing feedback lessons during the pilot.
- Use this directory to point agents toward reusable judgment; use `skills/apk-analysis/feedback_history/` for historical lesson records.
- If an insight becomes a cross-skill or all-repo prevention gate, promote it to `shared-rules/` or `shared-rules/failure-patterns/` instead of keeping it only here.
