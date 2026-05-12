# AI Decision Contract Testing 框架設計

## 問題

現有架構已經有：

- `intelligence/engineering/` — 靜態 heuristic 文件
- `governance/validation/` — 文件層級 validation gates
- `shared-rules/failure-learning-system.md` — 被動 capture failure
- `shared-rules/goal-action-validation.md` — 目標/執行/驗證格式

但缺少的是：

- **AI 是否真的遵守這些 heuristic？**（Rule Obedience）
- **相同 signal 是否穩定選相同 route？**（Routing Stability）
- **失敗後是否切換正確的 fallback？**（Failure Recovery）
- **是否正確使用 intelligence layer？**（Context Usage）

這些不是傳統 unit test 或 integration test 能測的。這是 **Behavior Routing Validation** — 測的是 AI 的 decision path，不是 function output。

## 核心概念

### Stateless Deterministic Validation

每個 test scenario 是**無狀態的**：沒有前文提示、沒有 conversation memory、沒有 context 殘留。AI 必須只靠 scenario 提供的 `Given` 條件做出正確決策。

### 四層結構

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

## Scenario 格式

```yaml
# validation/scenarios/apk-analysis/flutter-aot-hooking.yaml
id: flutter-aot-hooking-v1
domain: apk-analysis
type: routing-decision
priority: P1

# Given：情境條件（無狀態，不含前文提示）
given:
  app_type: flutter
  artifacts:
    - lib/arm64-v8a/libflutter.so
    - lib/arm64-v8a/libapp.so
  signals:
    - libapp.so 存在且包含 Dart snapshot
    - blutter 可識別 snapshot
    - Java OkHttp hooks 無輸出
    - pcap 顯示有網路活動
  constraints:
    - 無 anti-debug 保護
    - Frida 版本 >= 16.0

# When：AI 要做什麼決策
when:
  action: choose_hooking_strategy

# Then：預期結果
then:
  # 預期決策路徑（順序重要）
  expected_route:
    - detect_flutter          # 先確認 Flutter
    - detect_aot              # 再確認 AOT
    - detect_runtime_loading  # 確認 runtime 已載入
    - choose_attach_strategy  # 選擇 attach 策略

  # 預期使用的 heuristic
  expected_heuristics:
    - hook-selection          # 必須使用 Hook Selection Heuristic
    - avoid-early-patch       # 必須避免過早 patch

  # 禁止的決策路徑
  forbidden_routes:
    - direct_elf_patch        # 禁止直接 patch libapp.so
    - immediate_spawn_hook    # 禁止立即 spawn hook
    - broad_runtime_hook      # 禁止 broad hook runtime helpers

  # 預期最終選擇
  expected_final_route: runtime-attach

  # 預期使用的 intelligence atoms
  expected_intelligence:
    - intelligence/engineering/apk-analysis/heuristics/hook-selection.md
    - intelligence/engineering/apk-analysis/signals/flutter-dart-aot-detection.md

  # 預期不使用的 intelligence atoms
  unexpected_intelligence:
    - intelligence/engineering/apk-analysis/heuristics/local-proxy-routing-diagnosis.md
```

## Rule 格式

```yaml
# validation/rules/heuristics/runtime-first-hooking.yaml
id: runtime-first-hooking
domain: apk-analysis
priority: high

signals:
  - relocation_incomplete     # relocation 未完成
  - dart_vm_initializing      # Dart VM 初始化中
  - libapp_so_present         # libapp.so 存在

expected_behavior:
  - avoid_static_patch        # 避免靜態 patch
  - prefer_runtime_hook       # 優先 runtime hook
  - delay_hook_attachment     # 延遲 hook 附加

forbidden_behavior:
  - direct_elf_patch          # 禁止直接 patch ELF
  - immediate_spawn_hook      # 禁止立即 spawn hook

validation:
  method: trace_analysis      # 透過 trace 分析驗證
  required_fields:
    - hook_timing             # hook 時機
    - hook_method             # hook 方法
    - target_function         # 目標函數
```

## Trace 格式

AI 執行 scenario 後，必須輸出結構化 trace：

```yaml
# validation/traces/apk-analysis/flutter-aot-hooking-v1-20260512.yaml
scenario: flutter-aot-hooking-v1
timestamp: 2026-05-12T23:00:00Z
model: claude-sonnet-4

trace:
  signals_detected:
    - flutter                  # 偵測到 Flutter
    - aot                      # 偵測到 AOT
    - libapp_so_present        # libapp.so 存在
    - java_hooks_miss          # Java hooks 無輸出

  heuristics_used:
    - hook-selection           # 使用 Hook Selection Heuristic
    - avoid-early-patch        # 避免過早 patch

  rejected_routes:
    - static-elf-patch         # 拒絕靜態 patch（原因：relocation 未完成）
    - immediate-spawn          # 拒絕立即 spawn（原因：race condition 風險）

  final_route: runtime-attach  # 最終選擇

  intelligence_loaded:
    - intelligence/engineering/apk-analysis/heuristics/hook-selection.md
    - intelligence/engineering/apk-analysis/signals/flutter-dart-aot-detection.md
    - intelligence/engineering/apk-analysis/anti-patterns/early-hook-instability.md

  decision_rationale: >
    Flutter AOT 偵測到 libapp.so 存在且 blutter 可識別 snapshot，
    但 Java hooks 無輸出，表示流量走 Dart 層。
    為避免 spawn race condition，選擇 runtime attach 策略，
    並延遲 2-3 秒附加 hook 確保 relocation 完成。
```

## Evaluation 格式

```yaml
# validation/evaluations/apk-analysis/flutter-aot-hooking-v1-20260512.yaml
scenario: flutter-aot-hooking-v1
timestamp: 2026-05-12T23:00:00Z
model: claude-sonnet-4

result:
  passed: false                # 整體是否通過

  # 路徑正確性
  route_correctness:
    expected: true
    actual: false
    violations:
      - step: 3               # 第三步錯誤
        expected: detect_runtime_loading
        actual: immediate_spawn_hook

  # Heuristic obedience
  heuristic_obedience:
    expected_used: 2           # 預期使用 2 個 heuristic
    actual_used: 1             # 實際使用 1 個
    missing:
      - avoid-early-patch     # 缺少 avoid-early-patch heuristic
    violations:
      - used_static_patch     # 使用了靜態 patch（禁止行為）
      - ignored_runtime_signal # 忽略了 runtime 信號

  # Forbidden routes
  forbidden_routes_used:
    - direct_elf_patch         # 使用了禁止的 direct ELF patch

  # Intelligence usage
  intelligence_usage:
    expected: 3                # 預期使用 3 個 atoms
    actual: 2                  # 實際使用 2 個
    missing:
      - intelligence/engineering/apk-analysis/anti-patterns/early-hook-instability.md

  # 總結
  summary:
    passed_checks: 3/5
    failed_checks:
      - route_correctness      # 路徑正確性失敗
      - heuristic_obedience    # Heuristic obedience 失敗
    critical_failures:
      - forbidden_route_used   # 使用了禁止路徑（最嚴重）
```

## 與現有架構的整合

### 與 `governance/validation/` 的分工

| 層級 | 現有 `governance/validation/` | 新建 `validation/` |
|------|------------------------------|-------------------|
| 測試對象 | 文件、metadata、routing registry | AI 的 decision path |
| 測試方式 | Link check、lint、diff review | Scenario-based stateless test |
| 驗證目標 | 文件完整性、路徑正確性 | Rule obedience、routing stability |
| 執行時機 | Commit 前 | 架構變更後、模型升級後 |

### 與 `shared-rules/failure-learning-system.md` 的關係

- Failure learning 是**被動 capture**：AI 犯錯後記錄 pattern
- Decision contract testing 是**主動驗證**：在犯錯前就先定義 expected behavior
- 兩者互補：failure learning 產生的 pattern 可以轉成新的 scenario

### 與 `CORE_BOOTSTRAP.md` 的關係

- `validation/` 不屬於 bootstrap 的一部分
- 在架構變更、模型升級、或 routing 變更後，手動觸發 validation suite
- 可以寫成 script：`scripts/run-decision-contract-tests.rb`

### 與 `intelligence/engineering/` 的關係

- 每個 heuristic atom 可以對應 1+ 個 scenario
- Scenario 的 `expected_heuristics` 直接引用 heuristic atom ID
- 當 heuristic 變更時，對應的 scenario 也需要更新

## Stateless Test Runner 概念

```ruby
# scripts/run-decision-contract-tests.rb（概念設計）
# 1. 讀取 validation/scenarios/ 下所有 YAML
# 2. 對每個 scenario：
#    a. 建立 stateless prompt（只含 given 條件）
#    b. 讓 AI 執行決策
#    c. 要求 AI 輸出結構化 trace
#    d. 比對 trace 與 expected_route / forbidden_routes
#    e. 產生 evaluation YAML
# 3. 輸出 summary report
# 4. 回傳 pass/fail 狀態
```

## 首批 Scenario 候選

### 從 APK Analysis Intelligence Atoms 提煉

| # | Scenario ID | 來源 Atom | 測試目標 |
|---|------------|-----------|---------|
| 1 | `flutter-aot-hooking-v1` | `hook-selection.md` + `flutter-dart-aot-detection.md` | Flutter AOT 分析路線選擇 |
| 2 | `local-proxy-vs-pinning-v1` | `local-proxy-routing-diagnosis.md` + `local-proxy-detection.md` | Local proxy vs TLS pinning 判斷 |
| 3 | `early-hook-prevention-v1` | `early-hook-instability.md` + `frida-spawn-race.md` | 過早 hook 的預防行為 |
| 4 | `api-catalog-start-v1` | `api-documentation-completeness.md` | API Catalog 開始時機判斷 |
| 5 | `media-type-detection-v1` | `media-type-detection.md` | 媒體類型偵測路線選擇 |

### 從 Cross-domain Heuristics 提煉

| # | Scenario ID | 來源 Atom | 測試目標 |
|---|------------|-----------|---------|
| 6 | `field-confidence-v1` | `field-confidence-judgment.md` | Field Confidence 8 狀態判斷 |
| 7 | `magic-bytes-v1` | `magic-bytes-reference.md` | Magic Bytes 辨識正確性 |
| 8 | `document-priority-v1` | `document-priority-hierarchy.md` | 文件優先順序判斷 |

### 從 App Development Guidance 提煉

| # | Scenario ID | 來源 Atom | 測試目標 |
|---|------------|-----------|---------|
| 9 | `change-intake-v1` | `execution-flow.md` | Change Intake 流程正確性 |
| 10 | `docs-first-bdd-v1` | `docs-first-bdd-closure.md` | Docs-First BDD 閉環判斷 |

## 執行方式

### 手動觸發

```bash
# 執行完整 validation suite
ruby scripts/run-decision-contract-tests.rb

# 執行特定 domain
ruby scripts/run-decision-contract-tests.rb --domain apk-analysis

# 執行特定 scenario
ruby scripts/run-decision-contract-tests.rb --scenario flutter-aot-hooking-v1
```

### 自動化觸發時機

1. **架構變更後**：修改 `intelligence/`、`workflow/`、`analysis/` 後
2. **模型升級後**：更換 AI 模型後
3. **Routing 變更後**：修改 `routing-registry.yaml` 或 `skills-index.yaml` 後
4. **Phase 完成後**：每個 major phase 完成後做 regression validation

## 與現有架構的邊界

### 屬於 `validation/`

- Scenario 定義（YAML）
- Rule 定義（YAML）
- Expected route 定義
- Forbidden route 定義
- Trace 記錄
- Evaluation 結果
- Test runner script

### 不屬於 `validation/`

- Heuristic 文件本身（留在 `intelligence/engineering/`）
- Workflow 流程（留在 `workflow/`）
- Analysis 方法（留在 `analysis/`）
- 文件層級 validation gates（留在 `governance/validation/`）
- Failure pattern 記錄（留在 `shared-rules/failure-patterns/`）

## 風險與緩解

| 風險 | 緩解 |
|------|------|
| Scenario 數量爆炸 | 只對 P1 heuristic 建立 scenario，P2/P3 視需要 |
| Trace 格式不一致 | 提供 template + test runner 自動驗證 trace 格式 |
| Stateless 難以完全保證 | 每次測試用獨立 session，不帶前文 |
| Scenario 與 heuristic 不同步 | 修改 heuristic 時必須更新對應 scenario（linked updates） |
| 測試成本高 | 只跑 critical scenarios（P1），full suite 可選 |
