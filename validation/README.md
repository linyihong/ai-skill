# AI Decision Contract Testing

本目錄定義 **Behavior Routing Validation** — 測試 AI 的 decision path、rule obedience、routing stability，而不是 function output。

## 核心原則

- **Stateless**：每個 scenario 是無狀態的，沒有前文提示、conversation memory、context 殘留
- **Deterministic**：相同 scenario 應產出相同決策路徑
- **Traceable**：AI 必須輸出結構化 trace，記錄 signals、heuristics、rejected routes、final route

## 目錄結構

```
validation/
├── scenarios/           # 測試情境（YAML）
│   ├── apk-analysis/    # APK 分析領域
│   ├── app-dev/         # App 開發領域
│   └── travel/          # 旅遊規劃領域
├── rules/               # 規則定義（YAML）
│   ├── heuristics/      # Heuristic obedience rules
│   ├── routing/         # Routing stability rules
│   └── failure/         # Failure recovery rules
├── expected-routes/     # 預期決策路徑（YAML）
├── forbidden-routes/    # 禁止決策路徑（YAML）
├── traces/              # AI 執行 trace 記錄
└── evaluations/         # 評估結果
```

## 與 `governance/validation/` 的分工

| 層級 | `governance/validation/` | `validation/` |
|------|--------------------------|---------------|
| 測試對象 | 文件、metadata、routing registry | AI 的 decision path |
| 測試方式 | Link check、lint、diff review | Scenario-based stateless test |
| 驗證目標 | 文件完整性、路徑正確性 | Rule obedience、routing stability |
| 執行時機 | Commit 前 | 架構變更後、模型升級後 |

## Scenario 格式

每個 scenario 定義在 `scenarios/<domain>/<id>.yaml`：

```yaml
id: <unique-id>
domain: <domain>
type: routing-decision | heuristic-obedience | failure-recovery
priority: P1 | P2 | P3

given:
  # 情境條件（無狀態，不含前文提示）
  app_type: <type>
  artifacts:
    - <path>
  signals:
    - <signal>
  constraints:
    - <constraint>

when:
  action: <decision-action>

then:
  expected_route:
    - <step-1>
    - <step-2>
  expected_heuristics:
    - <heuristic-id>
  forbidden_routes:
    - <forbidden-route>
  expected_final_route: <route>
  expected_intelligence:
    - <path>
  unexpected_intelligence:
    - <path>
```

## Trace 格式

AI 執行 scenario 後輸出結構化 trace 到 `traces/<domain>/<id>-<date>.yaml`：

```yaml
scenario: <id>
timestamp: <ISO-8601>
model: <model-name>

trace:
  signals_detected:
    - <signal>
  heuristics_used:
    - <heuristic-id>
  rejected_routes:
    - <route>: <reason>
  final_route: <route>
  intelligence_loaded:
    - <path>
  decision_rationale: <text>
```

## Evaluation 格式

評估結果寫入 `evaluations/<domain>/<id>-<date>.yaml`：

```yaml
scenario: <id>
timestamp: <ISO-8601>
model: <model-name>

result:
  passed: <bool>
  route_correctness:
    expected: <int>
    actual: <int>
    violations:
      - step: <n>
        expected: <step>
        actual: <step>
  heuristic_obedience:
    expected_used: <int>
    actual_used: <int>
    missing:
      - <heuristic-id>
    violations:
      - <behavior>
  forbidden_routes_used:
    - <route>
  intelligence_usage:
    expected: <int>
    actual: <int>
    missing:
      - <path>
  summary:
    passed_checks: <n>/<total>
    failed_checks:
      - <check>
    critical_failures:
      - <failure>
```

## 首批 Scenario

| # | ID | 來源 Atom | 測試目標 |
|---|-----|-----------|---------|
| 1 | `flutter-aot-hooking-v1` | `hook-selection.md` + `flutter-dart-aot-detection.md` | Flutter AOT 分析路線選擇 |
| 2 | `local-proxy-vs-pinning-v1` | `local-proxy-routing-diagnosis.md` + `local-proxy-detection.md` | Local proxy vs TLS pinning 判斷 |
| 3 | `early-hook-prevention-v1` | `early-hook-instability.md` + `frida-spawn-race.md` | 過早 hook 的預防行為 |

## 執行方式

```bash
# 執行完整 validation suite（未來實作）
ruby scripts/run-decision-contract-tests.rb

# 執行特定 domain
ruby scripts/run-decision-contract-tests.rb --domain apk-analysis

# 執行特定 scenario
ruby scripts/run-decision-contract-tests.rb --scenario flutter-aot-hooking-v1
```

## Failure → Scenario 閉環

當 AI 系統在執行中發生 routing 錯誤、heuristic 誤用、forbidden route 被選中、或任何可被 scenario 捕捉的行為錯誤時，該錯誤經驗必須被轉化為新的 validation scenario。

### 閉環流程

```
AI 執行錯誤
    │
    ▼
1. Capture ─── 記錄錯誤：什麼 signals 被誤判、什麼 forbidden route 被選中
    │
    ▼
2. Classify ── 判斷這是新的 failure pattern，還是既有 scenario 的變體
    │
    ▼
3. Create ──── 建立新的 scenario YAML 到 scenarios/failure-derived/
    │
    ▼
4. Link ────── 在 failure pattern 或 feedback lesson 中標註對應的 scenario ID
    │
    ▼
5. Validate ── 執行新 scenario 確認 trace 可產出、evaluation 可比對
```

### 判斷是否要建立 scenario

| 條件 | 應建立 scenario | 不應建立 scenario |
|------|----------------|-------------------|
| 錯誤類型 | Routing 錯誤、heuristic 誤用、forbidden route 被選中 | Token 耗盡、工具執行失敗、模型 hallucination |
| 可重現性 | 相同 signals 應產出相同決策 | 隨機性或環境依賴高 |
| Prevention 價值 | 未來可能再次發生 | 一次性事件 |
| 與既有 scenario 的關係 | 新的 signal 組合或新的 forbidden route | 既有 scenario 已涵蓋 |

### Failure-Derived Scenarios

`scenarios/failure-derived/` 存放從實際 AI 執行錯誤提煉的 scenario。這些 scenario 的格式與一般 scenario 相同，但多了 `failure_source` 欄位指向原始 failure 記錄：

```yaml
id: <id>
domain: <domain>
type: <type>
priority: <priority>
failure_source:           # 新增欄位：指向原始 failure
  pattern: <path-to-failure-pattern>
  lesson: <path-to-feedback-lesson>
  date: <ISO-8601>
  description: <簡述原始錯誤>

given:
  # ... 與一般 scenario 相同
when:
  # ...
then:
  # ...
```

### 首批 Failure-Derived Scenarios

| # | ID | 來源 Failure Pattern | 測試目標 |
|---|-----|---------------------|---------|
| 1 | `entrypoint-drift-v1` | `shared-rules/failure-patterns/entrypoint-positioning-drift.md` | 更新 entrypoint 時是否同時更新 title、opening、indexes |
| 2 | `source-mirror-write-v1` | `shared-rules/failure-patterns/source-mirror-write-drift.md` | 編輯前是否先確認 canonical source vs mirror |

## 與其他層的關係

- `intelligence/engineering/` — 每個 heuristic atom 可對應 1+ 個 scenario
- `governance/validation/` — 文件層級 validation gates，與本層互補
- `shared-rules/failure-learning-system.md` — 被動 capture failure，與本層主動驗證互補；failure 可 promotion 為 scenario
- `shared-rules/failure-patterns/` — 每個 validated failure pattern 可對應 1+ 個 failure-derived scenario
- `architecture/next-stage-upgrade-plan.md` — Durable Roadmap Goals 記錄本層狀態
