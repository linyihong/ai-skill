# Analysis Engineering Intelligence

`intelligence/engineering/analytical-reasoning/` 存放與「分析」相關的工程智慧，涵蓋 APK 分析、repository 分析等領域。這些 atoms 描述如何選擇分析路線、辨識技術特徵、避免常見錯誤，以及從分析結果中提煉可重用原則。

## 目錄結構

```
intelligence/engineering/analytical-reasoning/
├── README.md                           # 本文件
├── evidence-first-routing.md           # 以證據驅動分析路線選擇
├── highest-leverage-analysis-path.md   # 最高槓桿分析路徑
├── live-readiness-gates.md             # 實作就緒門檻
├── documentation-backfill-heuristic.md # 文件回填經驗法則（from repo-analysis）
├── traceability-heuristic.md           # 文件追溯性經驗法則（from repo-analysis）
├── heuristics/                         # 跨語言啟發式判斷規則（不包含語言特定知識）
├── signals/                            # 技術特徵辨識信號
├── failure/                            # 語言/框架特定的失敗模式與診斷
└── anti-patterns/                      # 可預防的錯誤模式（跨語言通用）
```

## 目前 atoms

| Atom | 狀態 | 來源 |
| --- | --- | --- |
| [`evidence-first-routing.md`](evidence-first-routing.md) | `validated` | `feedback/history/apk-analysis/common/` |
| [`highest-leverage-analysis-path.md`](highest-leverage-analysis-path.md) | `candidate-intelligence` | `feedback/history/apk-analysis/common/` |
| [`live-readiness-gates.md`](live-readiness-gates.md) | `validated` | `feedback/history/apk-analysis/common/` |
| [`documentation-backfill-heuristic.md`](documentation-backfill-heuristic.md) | `candidate-intelligence` | `analysis/repo/documentation-backfill.md` |
| [`traceability-heuristic.md`](traceability-heuristic.md) | `candidate-intelligence` | `analysis/repo/traceability-gate.md` |
| [`heuristics/`](heuristics/) | `pilot` | Technique decomposition + UI operation intelligence extraction |
| [`signals/`](signals/) | `pilot` | Technique decomposition |
| [`failure/`](failure/) | `pilot` | Technique decomposition + feedback history |
| [`anti-patterns/`](anti-patterns/) | `pilot` | Technique decomposition |

## Scope

本層負責：

- 從重複分析工作中抽出的穩定工程判斷。
- Anti-pattern 與 trade-off，例如何時不應繼續加 broad hooks，或何時 API shape 不足以支援 live SDK work。
- 比單一 technique file 更高層的 reusable decision guidance。
- 會影響未來 workflow 或 runtime routing 的 validated lesson promotion targets。
- **HOW TO THINK** 決策智慧：heuristics、anti-patterns、failure learning、signal detection。

本層不負責：

- Raw project evidence、hosts、endpoints、tokens、device IDs 或 private run logs。
- Step-by-step capture workflow；使用 `workflow/apk-analysis/` 或 `workflow/repo-analysis/`。
- Traffic/runtime observation methods；使用 `analysis/apk/`。
- **HOW TO DO** 操作流程、command、setup；使用 `analysis/apk/workflows/`。
- Skill-specific lesson archive；在 promotion 前保留於 `feedback/history/`。

## Atom 分類規則

新增或整理 atom 時，用以下問題判斷：

| 主要回答 | 分類 |
| --- | --- |
| 看到什麼 signal 代表哪種技術路線？ | `signals/` |
| 這個症狀發生時如何診斷 root cause？ | `failure/` |
| 未來應採用哪個判斷規則或取捨？ | `heuristics/` |
| 哪種做法應避免，因為它常導致錯誤或不穩？ | `anti-patterns/` |
| 只是操作順序、命令模板或 capture procedure | `analysis/apk/workflows/` 或 `workflow/<domain>/` |
| 只是某次 App 的 raw evidence / case transcript | 業務專案 evidence；去敏後才可成為 `feedback/history/` lesson |

## 與其他層的關係

- `analysis/apk/` 提供 APK 分析的具體操作方法，本層提供背後的原則與 why。
- `analysis/repo/` 提供 repository 分析的具體操作方法，本層提供背後的原則與 why。
- `workflow/apk-analysis/` 提供執行流程，本層提供選擇流程的決策邏輯。
- `intelligence/engineering/heuristics/` 存放跨領域通用的啟發式規則（如 field-confidence-judgment、magic-bytes-reference），本層存放分析領域專用的啟發式。

## Promotion Rule

A lesson can be promoted here when:

1. It is generalized and sanitized.
2. It has validation evidence or repeated use.
3. It affects engineering judgment, not only one tool command.
4. The original `feedback/history/` entry remains discoverable.
5. `knowledge/indexes/README.md` and metadata are updated if the promoted atom becomes a routing target.

## Compatibility Notes

- Do not delete or rewrite existing feedback lessons during the pilot.
- Use this directory to point agents toward reusable judgment; use `feedback/history/apk-analysis/` for historical lesson records.
- If an insight becomes a cross-skill or all-repo prevention gate, promote it to `enforcement/` or `enforcement/failure-patterns/` instead of keeping it only here.
