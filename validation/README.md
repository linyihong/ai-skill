# AI Decision Contract Testing

本目錄定義 **Behavior Routing Validation** — 測試 AI 的 decision path、rule obedience、routing stability，而不是 function output。

## 核心原則

- **Stateless**：每個 scenario 是無狀態的，沒有前文提示、conversation memory、context 殘留
- **Deterministic**：相同 scenario 應產出相同決策路徑
- **Traceable**：AI 必須輸出結構化 trace，記錄 signals、heuristics、rejected routes、final route

治理 gate 見 [`governance/ai-runtime-governance/validation-scenario-governance.md`](../governance/ai-runtime-governance/validation-scenario-governance.md)；本檔只保存 scenario、trace 與 evaluation 的執行格式。

## 目錄結構

```
validation/
├── evidence-types/      # L3 evidence_type catalog（what was proven；gate requires 只用 type token）
├── scenarios/           # 測試情境（YAML）
│   ├── apk-analysis/    # APK 分析領域
│   ├── app-dev/         # App 開發領域
│   ├── software-delivery/ # Requirements cognition / behavior validation / delivery correctness
│   ├── architecture/    # 架構選型與 DDD adoption 決策
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

## Evidence Types Catalog

L3 **Validation Capability** produces Evidence（artifact + proof shape）。`evidence_type` 只回答「證明了什麼」；`collection_method` 與 `artifact_shape` 不得作 gate token。

- Catalog: [`evidence-types/README.md`](evidence-types/README.md)
- Types: `source_contract`, `user_visible`, `navigation`, `state_persistence`, `media_playback`, `temporal_behavior`
- Gate `requires:` 只列 `evidence:<type>`；trace chain：**gate → claim → artifact**
- OQ-5 **reject inheritance** — 用各 type 檔內 `supported_collection_methods` / `supported_artifact_shapes` 對照，不建 subtype 樹
- Gate vocabulary: [`workflow/software-delivery/validation/evidence-gate-vocabulary.md`](../workflow/software-delivery/validation/evidence-gate-vocabulary.md)
- Authority table: [`workflow/software-delivery/validation/authority-decision-table.md`](../workflow/software-delivery/validation/authority-decision-table.md)
- Failure catalog: [`workflow/software-delivery/validation/failure-evolution-catalog.md`](../workflow/software-delivery/validation/failure-evolution-catalog.md)
- Experience runtime (cross-cutting): [`workflow/cross-cutting/experience-runtime/README.md`](../workflow/cross-cutting/experience-runtime/README.md)

## Software Delivery Scenarios

- `software-delivery/evidence-type-projection-break-v1.yaml` — L2 behavior 已寫但 L3 evidence_type 缺失時不得宣稱 UX complete；`browser_review` 不得作 pass/fail token。
- `software-delivery/requirement-contradiction.yaml` — requirement / BDD / tests 衝突時不得直接 implementation。
- `software-delivery/product-impact-misalignment.yaml` — Impact Map 與 Customer Journey 不一致時不得直接產生 BDD 或 implementation plan。
- `software-delivery/missing-validation-target.yaml` — acceptance criteria 缺 validation target 時不得宣稱 ready。
- `software-delivery/stale-acceptance-criteria.yaml` — domain invariant 改變後舊 acceptance baseline 必須刷新。
- `software-delivery/behavior-scope-overclaim.yaml` — single scenario pass 不得宣稱 global feature correctness。
- `software-delivery/mutation-testing-effectiveness.yaml` — high coverage 不得直接等同測試有效；高風險邏輯用 targeted mutants / negative checks 驗證測試能抓錯。

## Cognitive State / Evidence Scenarios

- `failure-derived/projection-break-missing-browser-evidence-v1.yaml` — BDD/contract pass but preview enforcement on wrong DOM; require evidence envelope + playerStage scoping before L3 closure.
- `failure-derived/assumption-as-fact-v1.yaml` — unvalidated assumption 不得作為 execution fact。
- `failure-derived/hook-log-overrides-ui-contradiction-v1.yaml` — hook/log success 不得覆蓋 live UI / contract contradiction。
- `failure-derived/repeated-patch-autonomy-downgrade-v1.yaml` — repeated patch 無新證據時必須降 autonomy 或進 recovery。
- `failure-derived/local-evidence-global-claim-v1.yaml` — local evidence 不得支撐 global completion claim。
- `failure-derived/stale-frame-cross-domain-contamination-v1.yaml` — stale execution frame 跨 domain reuse 必須觸發 rediscovery。
- `failure-derived/intent-stability-drift-v1.yaml` — current action 脫離 original goal / validation target 時必須 realign。
- `failure-derived/contradiction-propagates-checkpoint-v1.yaml` — upstream assumption 被推翻時 downstream checkpoint 必須 invalidated。
- `failure-derived/recovery-exit-without-criteria-v1.yaml` — recovery exit criteria 未滿足時不得恢復 autonomy。
- `failure-derived/recovery-reentry-no-new-evidence-v1.yaml` — 同 contradiction class 無新 evidence 時不得重跑相同 recovery。
- `failure-derived/governance-minimality-small-task-v1.yaml` — 小型低風險任務不應啟動過重 governance chain。
- `failure-derived/tier3-does-not-block-tier0-tier2-v1.yaml` — Tier 3 cognitive optimization 不得阻塞 Tier 0-2。
- `failure-derived/meta-governance-no-runtime-promotion-v1.yaml` — 無具體 validated failure 時 meta-governance 不 promotion 到 runtime。
- `failure-derived/metadata-yaml-not-executable-contract-v1.yaml` — metadata YAML 不得被誤判為 executable contract；必須檢查 contract fields 與 runtime projection data。
- `failure-derived/markdown-yaml-contract-bypass-v1.yaml` — companion YAML contract 存在時不得只讀 Markdown；必須讀 YAML gates 並驗證 runtime projection。
- `failure-derived/executable-contract-close-loop-miss-v1.yaml` — 新增 executable contract 後必須 compile / refresh / validate / projection assertion / commit / push / readback。

## Runtime Close-Out Scenarios

- `runtime/feedback-report-required-v1.yaml` — final close-out 不得缺 Feedback / Learning Report。
- `runtime/feedback-report-schema-v1.yaml` — Feedback / Learning Report 只做 presence / schema / enum / field-combination 機械檢查，不做語義判斷。
- `runtime/non-local-repo-feedback-none-allowed-v1.yaml` — non-local repo 仍可在沒有 reusable learning 時回報 `feedback_decision: NONE`。
- `failure-derived/feedback-needed-but-not-reported-v1.yaml` — user correction / reusable runtime gap 後不得漏報 `feedback_decision: NEEDED`。

## Sanitization Scenarios

- `runtime/sanitization-metadata-derived-fail-v1.yaml` — shared-layer 內容含 project metadata 宣告的 private token（214a415 reconstruction）必須 block。
- `runtime/sanitization-metadata-derived-pass-v1.yaml` — 同 token 在 project-local（`shared_layer:false`）放行；未宣告的 framework concept bootstrap-safe 放行。
- `runtime/sanitization-placeholder-allowed-v1.yaml` — `<PROJECT_ROOT>` / `<USER>` 等 placeholder 形不得被 generic pattern 誤擋。

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

| # | ID | 來源 Atom / 架構變更 | 測試目標 |
|---|-----|---------------------|---------|
| 1 | `flutter-aot-hooking-v1` | `hook-selection.md` + `flutter-dart-aot-detection.md` | Flutter AOT 分析路線選擇 |
| 2 | `local-proxy-vs-pinning-v1` | `local-proxy-routing-diagnosis.md` + `local-proxy-detection.md` | Local proxy vs TLS pinning 判斷 |
| 3 | `early-hook-prevention-v1` | `early-hook-instability.md` + `frida-spawn-race.md` | 過早 hook 的預防行為 |
| 4 | `new-category-registration-v1` | Go-native `ai-skill runtime validate` covers directory structure checks | 新增 intelligence/analysis/workflow 類別時是否正確建立 README.md、更新 routing registry |
| 5 | `cargo-cult-ddd` | DDD integration plan + architecture metadata | 低 complexity 專案不得預設 full DDD |
| 6 | `architecture-fit-mismatch` | architecture selection governance | 高 complexity 專案不得被壓成 CRUD-only |
| 7 | `overengineering-detection` | overengineering metadata / workflow | 架構複雜度超過 business complexity 時需 simplification review |
| 8 | `bounded-context-collapse` | bounded context heuristics | 多 domain 不得被錯誤混成單一 global model |
| 9 | `aggregate-explosion` | aggregate boundary heuristics | aggregate boundaries 過度切分時需 invariant review |

## 執行方式

```bash
# 目前可執行的 Go-native runtime validation
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime validate

# 修改 validation/ 或 knowledge/ 後，先 refresh generated surfaces 再 validate
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime refresh
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime validate
```

Decision-contract scenario runner 尚未實作；不要引用已移除或不存在的 Ruby runner。新增 scenario 時，先用 Go-native `runtime validate` 檢查 registry / knowledge runtime invariants，scenario-specific runner 需等 `ai-skill` CLI 增加對應命令後再寫入文件。

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

| # | ID | 來源 Failure Pattern / 架構變更 | 測試目標 |
|---|-----|-------------------------------|---------|
| 1 | `entrypoint-drift-v1` | `enforcement/failure-patterns/entrypoint-positioning-drift.md` | 更新 entrypoint 時是否同時更新 title、opening、indexes |
| 2 | `source-mirror-write-v1` | `enforcement/failure-patterns/source-mirror-write-drift.md` | 編輯前是否先確認 canonical source vs mirror |
| 3 | `shared-rules-architecture-drift-v1` | `enforcement/failure-patterns/shared-rules-architecture-drift.md` | 架構重構後是否同步檢查 enforcement/ 路徑參考 |
| 4 | `feedback-history-consolidation-v1` | `feedback/history/` 目錄整合（2026-05-13） | 新 feedback lesson 是否正確使用 feedback/history/<domain>/ 而非舊 skills/ 路徑 |
| 5 | `runtime-recovery-navigation-mismatch` | Runtime recovery escalation system（2026-05-20） | APK navigation / UI evidence mismatch 時是否停止 capture patch 並進 recovery |
| 6 | `runtime-recovery-user-contradiction` | Runtime recovery escalation system（2026-05-20） | 使用者指出 route/source miss 時是否停止舊 execution graph |
| 7 | `runtime-recovery-source-miss` | Runtime recovery escalation system（2026-05-20） | source-of-truth miss 時是否補 required reload set 而非只靠局部測試 |

## 誰會參考這裡（Inbound References）

- [`route.validation.ai-decision-contract`](../knowledge/runtime/routing-registry.yaml:583) — primary_source 為 `validation/README.md`
- [`route.evaluations.scenario-results`](../knowledge/runtime/routing-registry.yaml:748) — required_dependencies 引用 `validation/README.md`
- [`route.traces.decision-traces`](../knowledge/runtime/routing-registry.yaml:802) — required_dependencies 引用 `validation/README.md`
- [`enforcement/failure-learning-system.md`](../enforcement/failure-learning-system.md) — 定義 Failure → Scenario 閉環流程
- [`enforcement/failure-patterns/`](../enforcement/failure-patterns/) — 每個 failure pattern 可對應 failure-derived scenario

## 與其他層的關係

- `intelligence/engineering/` — 每個 heuristic atom 可對應 1+ 個 scenario
- `governance/validation/` — 文件層級 validation gates，與本層互補
- `enforcement/failure-learning-system.md` — 被動 capture failure，與本層主動驗證互補；failure 可 promotion 為 scenario
- `enforcement/failure-patterns/` — 每個 validated failure pattern 可對應 1+ 個 failure-derived scenario
- `plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md` — Durable Roadmap Goals 記錄本層狀態
