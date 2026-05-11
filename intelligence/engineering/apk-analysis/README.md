# APK Analysis Engineering Intelligence

`intelligence/engineering/apk-analysis/` 是 `apk-analysis` pilot 中用來保存穩定工程判斷的候選位置。Pilot 期間，lesson history 仍保留在 `skills/apk-analysis/feedback_history/`；本目錄只保存已抽象成 engineering intelligence 的判斷與 promotion map。

## 目前 intelligence atoms

| Atom | Status | Source |
| --- | --- | --- |
| [`highest-leverage-analysis-path.md`](highest-leverage-analysis-path.md) | `candidate-intelligence` | `skills/apk-analysis/feedback_history/common/2026-05-07_131000-highest-leverage-analysis-path.md` |

## Scope

本層負責：

- 從重複 APK 分析工作中抽出的穩定工程判斷。
- Anti-pattern 與 trade-off，例如何時不應繼續加 broad hooks，或何時 API shape 不足以支援 live SDK work。
- 比單一 technique file 更高層的 reusable decision guidance。
- 會影響未來 workflow 或 runtime routing 的 validated lesson promotion targets。

本層不負責：

- Raw project evidence、hosts、endpoints、tokens、device IDs 或 private run logs。
- Step-by-step capture workflow；使用 `workflow/apk-analysis/`。
- Traffic/runtime observation methods；使用 `analysis/apk/`。
- Skill-specific lesson archive；在 promotion 前保留於 `skills/apk-analysis/feedback_history/`。

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
