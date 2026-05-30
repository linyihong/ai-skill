# Evidence Chain Slice（單次分析筆記模板 + 證據鏈要求 + 失敗紀錄）

> **Cognitive Slice**：`apk-evidence-chain`（從 [`../artifact-gates.md`](../artifact-gates.md) §5+§6+§7 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-evidence-chain` |
| `purpose` | 提供單次分析筆記模板、證據鏈完整性要求、失敗也要記錄的格式 |
| `type` | `execution` |
| `tags` | artifact-gate, evidence |
| `load_when` | 記錄分析筆記 / 證據鏈 / 失敗 capture、需要 evidence template |
| `do_not_load_when` | 純 reference 查閱、無新 analysis 進行中 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「證據要怎麼記錄、要含哪些層次（pcap → CONNECT → hook → decrypt → fixture）」的 ordering / artifact gate；通過 workflow membership test |
| `canonical_source` | 本檔（原 `artifact-gates.md` §5 單次分析筆記模板 + §6 證據鏈要求 + §7 失敗也要記錄） |
| `dependencies` | `apk-sanitization`（evidence 必先去敏）、`apk-documentation-discipline`（寫作品質） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Scenario B + C 顯示本 slice 為 evidence-only / mixed 任務必載入 |

## 5. 單次分析筆記模板

```markdown
# [APK / 功能] 分析紀錄

## Scope
- APK:
- Version:
- Package:
- Device / emulator:
- Authorization:
- Goal:

## Environment
- OS:
- adb:
- Frida:
- Proxy tool:
- Static tools:

## Hypotheses
| Hypothesis | Test | Result |
| --- | --- | --- |
| localhost bridge | lo pcap | |
| system proxy / MITM | proxy capture | |
| Java HTTP stack | Java hook | |
| Flutter / native | connect backtrace / AOT strings | |

## Evidence
| Evidence | Path / excerpt | Interpretation |
| --- | --- | --- |
| pcap | `<path>` | |
| hook log | `<path>` | |
| static search | `<path or command>` | |
| screenshot / UI hierarchy | `<path>` | |

## Findings
- Finding 1.
- Finding 2.

## Feature Reconstruction Handoff
- Feature ID:
- Capability:
- User goal:
- Entry screens:
- Primary operations:
- Candidate domain concepts:
- API / interface contracts:
- State and error handling:
- Data lifecycle:
- Fixtures / validation:
- Open questions for app-development-guidance:

## Unknowns
- Unknown 1.

## Next Steps
1. Next validation.
2. Next fixture or test.

## Sanitization
- Tokens redacted:
- Device identifiers redacted:
- User data removed:
```

## 6. 證據鏈要求

好文件不只寫「成功」，還要寫為什麼相信它成功：

- pcap 證明對外 TLS host 存在。
- proxy CONNECT 證明導流成功。
- hook log 證明 request object 在 TLS 前可見。
- decrypt hook 或離線 decoder 證明 inner JSON 正確。
- fixture / test 證明規則可重跑。

## 7. 失敗也要記錄

失敗紀錄應包含：

- 嘗試了什麼。
- 期望看到什麼。
- 實際看到什麼。
- 排除了什麼假設。
- 是否要重試，或是否停止投入。

例：

```text
Java OkHttp hook installed successfully, but no target host/path appeared while pcap showed TLS traffic to the API host. This rules out the Java OkHttp path for the tested flow and shifts the next step to native/Flutter analysis.
```

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
