# Feature Reconstruction Handoff Slice

> **Cognitive Slice**：`apk-feature-handoff`（從 [`../artifact-gates.md`](../artifact-gates.md) §4 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-feature-handoff` |
| `purpose` | 把分析結果整理成可讓後續 agent 用 `app-development-guidance` 重建功能的 handoff 文件 |
| `type` | `execution` |
| `tags` | artifact-gate, handoff |
| `load_when` | 產出 feature reconstruction handoff、使用者問「能不能重建」/「架構是什麼」/「有沒有 API 文件」 |
| `do_not_load_when` | 探索期、尚未要 handoff、單純 API capture |
| `owner_layer` | workflow |
| `layer_justification` | 規定「handoff 文件要含哪 8 個面向、什麼時候必須補齊」的 ordering / artifact gate；通過 workflow membership test |
| `canonical_source` | 本檔（原 `artifact-gates.md` §4 Feature Reconstruction Handoff） |
| `dependencies` | `apk-ui-architecture-map`、`apk-api-catalog`、`apk-domain-runtime-baseline`、`apk-evidence-chain` |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | feature handoff finish gate 任一觸發條件成立時應載入本 slice |

## 4. Feature Reconstruction Handoff

若分析目標是讓後續 agent 能用 `app-development-guidance` 重新做出同等功能，專案分析文件不能只列 endpoint。

### 最低表格

| 面向 | 必填內容 |
| --- | --- |
| Feature / Capability | 功能名稱、使用者目標、入口 screen、非目標或未知限制。 |
| UI Behavior | screen id、route id、operation id、前置狀態、tap/swipe/input 步驟、可見結果。 |
| Domain Concepts | 從 UI 文案、response fields、狀態碼推得的 entity、value object、state、command、event。 |
| API / Interface Contract | method/path shape、headers、query/body、response wrapper、inner payload、auth/session、pagination、cache、idempotency。 |
| State And Error Handling | loading/empty/error/success 狀態、錯誤碼、重試、登入過期、權限不足、限流、離線或快取行為。 |
| Data Lifecycle | 欄位來源、derived-from、local cache/storage、刷新時機、敏感性、保留/過期行為。 |
| Validation Evidence | pcap/MITM/hook/replay/fixture/screenshot/UI hierarchy/automation script 的去敏引用。 |
| Unknowns / Assumptions | 未觸發流程、低信心 mapping、缺少樣本、未驗證 edge case。 |

### Feature handoff finish gate

當某個具名 feature/tab/module 已被分析到「核心 UI 操作與主要 API flow 可說明」的程度時，必須在同一輪補齊或更新 project-level feature handoff 文件。

觸發條件包含任一項：

- 核心 flows 已從 `Candidate` 升到 `Confirmed`。
- agent 已能回答此功能的 entry path、主要 UI 區塊、API request keys / response schema、狀態與缺口。
- 使用者問「有沒有 API 文件」、「能不能重建」、「架構是什麼」。

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
