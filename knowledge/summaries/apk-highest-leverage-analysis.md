# intelligence.apk-highest-leverage-analysis

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.apk-highest-leverage-analysis` |
| Source path | [`../../intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md`](../../intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md) |
| Lifecycle | `candidate` |
| Summary | APK 分析 checkpoint 應先界定未知，再依 time-to-evidence、語意距離、安全性與 validation clarity 選擇最高收益路線。 |
| When to read | APK 分析卡住、可選多條 route，或需要判斷 UI、API replay、hook、pcap、MITM、static xref 哪條先做時。 |
| Do not use for | 不可取代 `workflow/apk-analysis/execution-flow.md`、授權確認或完整 dependency reads。 |
| Validation signal | Intelligence atom 連回 feedback lesson 與 `execution-flow.md` 的選擇主線步驟；runtime registry 有對應 route。 |
| Last checked | 2026-05-11 |

## Checklist

- 用一句話寫出目前 unknown。
- 比較 2-4 條 evidence-supported routes。
- 選擇 evidence-to-cost ratio 最高的主線。
- 記錄 fallback route 與切換條件。
- Attribution 重要時，回到 UI / operation evidence 對照。
