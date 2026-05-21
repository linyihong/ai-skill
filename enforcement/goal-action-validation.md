# 目標、執行、驗證共同流程

本規則適用於本庫所有 enforcement rules、skills、templates、feedback lessons，以及由 skill 產出的專案文件。目的不是要求暴露模型內部推理，而是讓每個對外可見的工作單元都有清楚的目標、實際做了什麼、以及如何確認沒有偏離。

## 核心規則

每個重要工作單元都要能回答三件事：

| 欄位 | 必填內容 |
| --- | --- |
| 目標 | 這一步要解決什麼問題、回答什麼問題、產出什麼文件或完成什麼變更。 |
| 執行 | 實際做了哪些檢查、讀了哪些文件、改了哪些檔案、跑了哪些命令或採用了哪些判斷流程。 |
| 驗證 | 如何確認結果可信：測試、lint、link check、diff review、fixture、證據交叉比對、人工檢查、或明確的參考來源。 |

若該工作單元是純分析、規劃、判斷或問答，沒有可執行測試，也不能把驗證欄留空。改用：

- **參考來源**：讀過的文件、程式碼、規格、證據、使用者提供的條件、既有 enforcement rule 或 skill 條文。
- **推論邊界**：哪些是已確認、哪些是推測、哪些仍是 open question。
- **反查方式**：使用者或下一個 agent 可以去哪裡驗證這個結論。

## 使用時機

| 情境 | 要求 |
| --- | --- |
| 修改文件或程式 | 說明目標、實際修改範圍、驗證方式。 |
| 產出分析結論 | 說明問題、檢查來源、結論依據與未知項。 |
| 制定計畫 | 說明計畫目標、依據、風險、完成後如何驗證。 |
| 回答「是否需要改」 | 說明判斷目標、檢查了哪些文件/證據、為什麼足夠或不足。 |
| 新增 feedback lesson 或 enforcement rule | 說明要沉澱的目標、推廣到哪些文件、連動更新與驗證。 |
| 使用者指出閉環不完整 | 說明漏掉的原因、已補強的 shared/skill 規則、連動更新清單，以及如何驗證不會再只修局部。 |

## 建議輸出形狀

簡短任務可以用一句話：

```text
目標：確認 X 是否缺失。執行：讀取 A/B/C。驗證：與規則 D 對照，缺口為 Y。
```

較複雜任務用表格：

| 工作單元 | 目標 | 執行 | 驗證 / 參考來源 |
| --- | --- | --- | --- |
| `<step>` | | | |

不需要每個微小動作都列表，但下列項目必須有目標/執行/驗證或參考來源：

- 重要結論。
- 文件或程式修改。
- 風險判斷。
- 缺口判斷。
- 「不需要修改」的決策。
- 更新流程完成聲明。

## 驗證方式範例

| 工作類型 | 驗證方式 |
| --- | --- |
| Markdown / skill 文件 | Markdown link check、lints、索引反查、diff review、required linked updates 檢查。 |
| API / SDK / tool guidance | Contract test、fixture、generated client compile check、sample run、BDD traceability。 |
| APK 分析文件 | UI path / operation id 交叉比對、pcap / hook / replay / fixture 引用、confidence label、open questions。 |
| app-development-guidance 文件 | BDD、Domain/API/Error Contract、test strategy、blocker question 是否齊全。 |
| 純問答或判斷 | 引用文件路徑、規則條文、使用者限制、已知證據與未知項。 |

## 防呆規則

- 不要只寫「已完成」「看起來可以」「不需要改」而沒有驗證或參考來源。
- 不要把未跑過的測試寫成已驗證；要寫「未執行」與原因。
- 不要把參考來源寫成空泛的「依經驗」；至少列出文件、程式碼、規格、證據或使用者條件。
- 不要為了滿足格式而產生冗長流水帳；重點是每個可見結論都能被反查。
- 若驗證失敗，應回到目標與執行重新修正，直到驗證通過、被明確標為 blocker，或被使用者決定 scope out。
- 若 reusable guidance 更新被指出混入 project incident 或漏做 linked updates，不可只修文字；必須依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) 分析原因並補強規則層。

### 驗證 Gate 參考

本規則的「驗證」概念已升級為 declarative blocking gates，現行 source 是 [`runtime/runtime.db`](../runtime/runtime.db) / [`runtime/compiler/embedded_data.rb`](../runtime/compiler/embedded_data.rb)。每個 phase 都有對應的 blocking gates（critical/high/medium severity），定義了 phase transition 的必要條件。

Agent 在執行驗證時應：

1. 查詢目前 phase 的 blocking gates（`runtime/runtime.db`，必要時讀 `runtime/compiler/embedded_data.rb`）
2. 依 severity 順序通過各 gate（critical → high → medium）
3. 若 gate 未通過 → 進入 recovery（`runtime/runtime.db` 的 recovery tables / `runtime/compiler/embedded_data.rb`）
4. 若所有 gate 通過 → 允許 phase transition

本節的 prose 規則（目標/執行/驗證三欄）仍適用於工作單元的內部品質檢查，但 phase transition 的 blocking gate 檢查應以 `runtime/runtime.db` / `runtime/compiler/embedded_data.rb` 為 authoritative reference。

### State-based Enforcement（狀態化強制規則）

下列驗證規則已對應到 runtime state machine 的 phase/gate 定義：

```yaml
# State-based enforcement mapping for goal-action-validation.md
# 這些規則已由 runtime state machine 管理，agent 不應再以 prose 方式逐條檢查。
state_based_enforcement:
  version: v1
  status: active
  owner_layer: enforcement/goal-action-validation
  description: >
    將目標/執行/驗證三欄規則對應到 runtime state machine 的 blocking gates。
    Agent 應優先查詢 runtime/runtime.db，必要時讀 runtime/compiler/embedded_data.rb，
    而非依賴本節的 prose 摘要或已移除的 standalone YAML。

  # 目標定義 → 由 runtime.db / embedded_data.rb 的 gate.execution.goal_defined 管理
  - rule: goal_defined
    phase: phase.execution
    gate: gate.execution.goal_defined
    severity: critical
    description: "本輪的執行目標已明確定義"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.execution.goal_defined"

  # 執行範圍 → 由 runtime.db / embedded_data.rb 的 gate.execution.scope_clear 管理
  - rule: scope_clear
    phase: phase.execution
    gate: gate.execution.scope_clear
    severity: high
    description: "執行範圍已明確，無 scope creep 風險"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.execution.scope_clear"

  # 驗證完整性 → 由 runtime.db / embedded_data.rb 的 gate.validation.all_obligations_met 管理
  - rule: validation_complete
    phase: phase.validation
    gate: gate.validation.all_obligations_met
    severity: critical
    description: "本輪所有 obligation 已完成，包含目標/執行/驗證三欄"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.validation.all_obligations_met"

  # 連動更新 → 由 runtime.db / embedded_data.rb 的 gate.validation.linked_updates_complete 管理
  - rule: linked_updates_complete
    phase: phase.validation
    gate: gate.validation.linked_updates_complete
    severity: critical
    description: "所有連動更新已執行"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.validation.linked_updates_complete"

  # 產出完整性 → 由 runtime.db / embedded_data.rb 的 gate.validation.artifacts_complete 管理
  - rule: artifacts_complete
    phase: phase.validation
    gate: gate.validation.artifacts_complete
    severity: high
    description: "本輪產出的 artifacts 完整且符合規範"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.validation.artifacts_complete"

  # 禁止動作 → 由 runtime.db / embedded_data.rb 的 gate.validation.no_forbidden_actions_used 管理
  - rule: no_forbidden_actions
    phase: phase.validation
    gate: gate.validation.no_forbidden_actions_used
    severity: critical
    description: "本輪未使用任何 forbidden actions"
    runtime_ref: runtime/runtime.db
    runtime_section: "gate.validation.no_forbidden_actions_used"
```

## 與其他規則的關係

- 去敏與敏感資料處理依 [`sanitization.md`](sanitization.md)。
- 文件用語依 [`neutral-language.md`](neutral-language.md)。
- 文件放置位置依 [`content-layering.md`](content-layering.md)。
- 可重用規則與專案證據邊界依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md)。
- 連動更新依 [`linked-updates.md`](linked-updates.md)。

← [回到共用規則索引](README.md)
