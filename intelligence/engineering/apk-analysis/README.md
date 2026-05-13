# APK Analysis Engineering Intelligence

`intelligence/engineering/apk-analysis/` 是 `apk-analysis` pilot 中用來保存穩定工程判斷的候選位置。Pilot 期間，lesson history 保留在 `feedback/history/apk-analysis/`；本目錄只保存已抽象成 engineering intelligence 的判斷與 promotion map。

## 目錄結構

```
intelligence/engineering/apk-analysis/
├── README.md                       # 本文件
├── evidence-first-routing.md       # 既有：以證據驅動分析路線選擇
├── highest-leverage-analysis-path.md # 既有：最高槓桿分析路徑
├── live-readiness-gates.md         # 既有：實作就緒門檻
├── heuristics/                     # 啟發式判斷規則（何時該用哪個技術）
├── anti-patterns/                  # 可預防的錯誤模式
├── failure/                        # 具體失敗模式與診斷
└── signals/                        # 技術特徵辨識信號
```

## 目前 intelligence atoms

| Atom | 狀態 | 來源 |
| --- | --- | --- |
| [`evidence-first-routing.md`](evidence-first-routing.md) | `validated` | `skills/apk-analysis/feedback_history/common/` |
| [`highest-leverage-analysis-path.md`](highest-leverage-analysis-path.md) | `candidate-intelligence` | `skills/apk-analysis/feedback_history/common/2026-05-07_131000-highest-leverage-analysis-path.md` |
| [`live-readiness-gates.md`](live-readiness-gates.md) | `validated` | `skills/apk-analysis/feedback_history/common/` |
| [`heuristics/`](heuristics/) | `pilot` | Technique decomposition from `skills/apk-analysis/techniques/`（已刪除） |
| [`anti-patterns/`](anti-patterns/) | `pilot` | Technique decomposition from `skills/apk-analysis/techniques/`（已刪除） |
| [`failure/`](failure/) | `pilot` | Technique decomposition from `skills/apk-analysis/techniques/`（已刪除） |
| [`signals/`](signals/) | `pilot` | Technique decomposition from `skills/apk-analysis/techniques/`（已刪除） |

## Scope

本層負責：

- 從重複 APK 分析工作中抽出的穩定工程判斷。
- Anti-pattern 與 trade-off，例如何時不應繼續加 broad hooks，或何時 API shape 不足以支援 live SDK work。
- 比單一 technique file 更高層的 reusable decision guidance。
- 會影響未來 workflow 或 runtime routing 的 validated lesson promotion targets。
- **HOW TO THINK** 決策智慧：heuristics、anti-patterns、failure learning、signal detection。

本層不負責：

- Raw project evidence、hosts、endpoints、tokens、device IDs 或 private run logs。
- Step-by-step capture workflow；使用 `workflow/apk-analysis/`。
- Traffic/runtime observation methods；使用 `analysis/apk/`。
- **HOW TO DO** 操作流程、command、setup；使用 `analysis/apk/workflows/`。
- Skill-specific lesson archive；在 promotion 前保留於 `feedback/history/apk-analysis/`。

## Candidate Intelligence Areas

| Area | Current source | Future shape |
| --- | --- | --- |
| Evidence-first route selection | `../../../skills/apk-analysis/WORKFLOW.md`, `../../../feedback/history/apk-analysis/common/` | decision guidance for high-leverage analysis routing |
| API catalog and runtime baseline readiness | `../../../skills/apk-analysis/DOCUMENTATION.md`, `../../../feedback/history/apk-analysis/http-api/`, `../../../feedback/history/apk-analysis/common/` | engineering guidance for live SDK / client readiness |
| Flutter / Dart AOT anti-patterns | `../../../skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除）, `../../../feedback/history/apk-analysis/flutter-dart-aot/` | stable failure patterns and hook-selection trade-offs |
| Local proxy and routing ambiguity | `../../../skills/apk-analysis/techniques/local-proxy/`（已刪除）, `../../../feedback/history/apk-analysis/local-proxy/` | guidance for separating routing, proxy, TLS, and attribution |
| Media chain completeness | `../../../skills/apk-analysis/techniques/media-hls/`（已刪除）, `../../../feedback/history/apk-analysis/media-hls/` | control-plane / data-plane and validation heuristics |

## Promotion Rule

A lesson can be promoted here when:

1. It is generalized and sanitized.
2. It has validation evidence or repeated use.
3. It affects engineering judgment, not only one tool command.
4. The original `feedback_history/` entry remains discoverable.
5. `knowledge/indexes/README.md` and metadata are updated if the promoted atom becomes a routing target.

## Compatibility Notes

- Do not delete or rewrite existing feedback lessons during the pilot.
- Use this directory to point agents toward reusable judgment; use `feedback/history/apk-analysis/` for historical lesson records.
- If an insight becomes a cross-skill or all-repo prevention gate, promote it to `shared-rules/` or `shared-rules/failure-patterns/` instead of keeping it only here.
- Old `skills/apk-analysis/techniques/` and `analysis/apk/techniques/` have been deleted (Phase C). Intelligence atoms here were extracted from those sources.
