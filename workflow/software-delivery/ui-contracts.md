# UI Contract Slice（Consumer / Screen / ViewModel）

> **Cognitive Slice**：`sd-ui-contracts`（從 [`development-process.md`](development-process.md) 的 frontend / consumer 實作缺口抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7 的 workflow membership test）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-ui-contracts` |
| `purpose` | 在前端、行動、CLI、SDK 或其他 consumer 實作前，定義 screen mapping、consumer 需要的資料、screen/state 行為、UI interaction、view model derivation、accessibility expectations 與 screen traceability |
| `type` | `execution` |
| `tags` | artifact-gate, contract, frontend, consumer, ui |
| `load_when` | 前端 / 行動 / CLI / SDK / consumer surface 需要和 provider 平行實作，或 AI agent 要從需求生成 UI、state、hooks、store、tests |
| `do_not_load_when` | 無 consumer surface、純 provider 內部變更、只修不影響 UI 行為的小錯 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「實作前要產哪些 UI-facing contracts、何時阻塞 API finalization、如何驗證」的 ordering / gate；通過 workflow membership test，不承載 evidence 取得方法（非 analysis），不論證長期模式（非 intelligence） |
| `canonical_source` | 本檔 |
| `dependencies` | `sd-intake`（需求與 actor intent）、`sd-contracts`（domain / architecture / API / error contracts）、`sd-test-strategy`（BDD / contract tests） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | 前後端平行實作時，BDD scenario 能追到 screens / APIs / tables，consumer needs、screen states、view model derivation、UI behavior tests 能反查到 BDD / contracts |

## Contract Stack

UI contract 不是 API contract 的附屬品。它把「consumer 要完成什麼行為」固定下來，讓 provider 和 consumer 共同對齊，而不是讓前端被動適配已定義的 endpoint。

```text
BDD behavior
  -> Screen Mapping
  -> Domain / Architecture contracts
  -> API + Error Contract
  -> Consumer Contract
  -> UI Behavior / Screen / ViewModel contracts
  -> Frontend or consumer implementation
```

當 Screen Mapping 或 Consumer Contract 暴露出 API / table / ownership 缺口時，先修訂 Domain / Architecture / API / Error Contract，再進入平行實作。

## Frontend / Consumer Contract Gate

存在前端、行動、CLI、SDK、job runner 或其他 consumer surface 時，進入實作前至少要完成下列 gate：

| Gate | 必要產物 |
| --- | --- |
| Screen Mapping | BDD scenario 到 screens、APIs、tables / stores 的對照 |
| Consumer Contract | consumer needs、freshness、loading、empty/error behavior、permissions |
| Screen Contract | screen states、actions、permissions、events、data dependencies |
| UI Behavior Contract | interaction、state transition、validation、feedback、navigation |
| ViewModel Contract | API / domain data 到 UI display model 的 derivation |

Accessibility Contract 依 screen 風險啟用：表單、modal、dynamic update、keyboard interaction、permission denied、error recovery 或重要 feedback 都需要明確 expectation。

## 1. Screen Mapping

Screen Mapping 是 BDD-lite 後的第一個 frontend / consumer traceability seed。它不是 UI 設計稿，也不應承載 design token、layout detail 或 framework pattern；它只回答「這個行為會落到哪些 consumer surface、provider operation 與資料所有權」。

| 欄位 | 說明 |
| --- | --- |
| `scenario` | 對應 BDD scenario / feature behavior |
| `screens` | 受影響的 stable screen、route、flow、CLI command、SDK method 或 job |
| `apis` | 需要的 API operation、event、RPC、GraphQL field、local adapter 或 provider function |
| `tables_or_stores` | 涉及的 database table、collection、local store、cache key、queue 或 durable state |
| `owner_context` | 資料 / invariant 所屬 bounded context 或 provider owner |
| `validation_target` | 對應的 BDD、contract、UI behavior、integration 或 fixture proof |

最小形狀：

```yaml
scenario: Watch Episode
screens:
  - home
  - detail
  - player
apis:
  - GET /videos
  - GET /episodes
  - POST /watch-history
tables_or_stores:
  - videos
  - episodes
  - watch_history
validation_target:
  - bdd: watch-episode
  - contract: watch-history-api
  - ui: player-screen-behavior
```

## 2. Consumer Contract

Consumer Contract 回答「consumer 為了完成行為需要什麼」，不是「provider 已經提供什麼」。適用於 web UI、mobile screen、CLI、SDK、job runner、test harness 或任何 provider 的使用方。

| 欄位 | 說明 |
| --- | --- |
| `consumer` | consumer 名稱，例如 screen、flow、CLI command、SDK method、integration job |
| `intent` | 對應的 actor intent / BDD behavior |
| `needs` | user-visible behavior 需要的資料、命令、事件或能力 |
| `freshness` | 資料新鮮度、refresh cadence、cache / stale policy |
| `loading` | skeleton、spinner、placeholder、progress、optimistic state |
| `empty_state` | 無資料時的行為與文案責任 |
| `error_behavior` | retry、fallback、inline error、toast、blocking state、recoverability |
| `permissions` | 可見性、可操作性、denied state |
| `observability` | 需要追蹤的 user action、error、latency 或 business event |

Consumer Contract 是 API finalization gate：若 `needs` 無法從現有 API / schema / event 取得，不能要求 consumer 手抄、推測或在 UI 層拼出不穩定行為。

## 3. UI Behavior Contract

UI Behavior Contract 回答「使用者操作後畫面和狀態如何變化」。它和 BDD 對齊，但更貼近 screen / component / client state。

| 欄位 | 說明 |
| --- | --- |
| `scenario` | 對應 BDD scenario 或 feature behavior |
| `given_state` | 初始資料、權限、session、feature flag、device / viewport 條件 |
| `action` | 使用者或系統觸發的操作 |
| `state_transition` | loading、dirty、validating、submitting、success、error 等狀態轉換 |
| `validation` | required fields、format、cross-field、async validation、server validation mapping |
| `feedback` | toast、inline message、banner、disabled reason、focus movement |
| `navigation` | route、modal、drawer、back behavior、deep link、redirect |
| `side_effects` | refresh list、invalidate cache、emit event、update local store |
| `test_target` | component test、interaction test、BDD/E2E、contract fixture |

## 4. Screen Contract

Screen Contract 是 AI agent 生成 UI 最穩定的輸入。它固定 screen 的 states、actions、permissions、events 與可觀察結果。

| 欄位 | 說明 |
| --- | --- |
| `screen` | stable screen / route / flow id |
| `intent` | 此 screen 支援的使用者任務 |
| `states` | loading、empty、error、success、partial、offline、permission denied |
| `actions` | create、edit、delete、submit、filter、sort、refresh、retry、navigate |
| `permissions` | role / capability 對 visibility、enabled state、navigation 的影響 |
| `events` | screen emits / consumes 的 domain 或 UI events |
| `data_dependencies` | 對應 Consumer Contract needs 與 API operation / event |
| `design_system` | 使用的 component primitives、layout constraints、responsive rules |
| `accessibility` | keyboard flow、focus order、ARIA / semantic roles、screen reader feedback |

## 5. Frontend ViewModel Contract

ViewModel Contract 固定 API / domain data 到 UI 顯示模型的 derivation。這是 AI agent 產生前端時最容易漂移的一層。

| 欄位 | 說明 |
| --- | --- |
| `source` | API response、domain entity、event、local storage、derived state |
| `view_model` | UI 使用的 stable model 名稱 |
| `fields` | 顯示欄位、格式化、fallback、單位、locale、排序 key |
| `derived_rules` | display name、badge、status label、permission flag、computed totals |
| `nullability` | missing / null / unknown / redacted 的顯示行為 |
| `error_mapping` | API / validation errors 到 field / form / page error 的映射 |
| `test_fixtures` | source fixture -> expected view model fixture |

不要讓 component 直接散落 API 欄位拼接規則。若 UI 需要 `displayName`、`isVip`、`badge`、`canEdit` 等 derived fields，先寫在 ViewModel Contract，再實作 mapper / selector / hook / presenter。

## 6. Accessibility Contract

當 screen 有 form、modal、dynamic update、keyboard interaction、permission denied、error recovery 或重要 feedback 時，需要 accessibility expectations。

| 欄位 | 說明 |
| --- | --- |
| `keyboard` | tab order、shortcut、escape / enter behavior、focus trap |
| `focus` | open / close / validation / navigation 後焦點落點 |
| `semantics` | landmark、role、label、description、live region |
| `contrast_motion` | contrast、reduced motion、loading / animation fallback |
| `assistive_feedback` | screen reader 如何得知 success、error、async progress |
| `validation` | 可用 accessibility lint、manual checklist、interaction test 或 design review 證明 |

## 7. Screen Traceability

Screen Traceability 是 Screen Mapping 的維護規則。每個 user-visible scenario 至少應能沿著下列方向追溯：

```text
BDD scenario
  -> Screen / route / command
  -> Consumer needs
  -> API / event / provider operation
  -> table / store / durable state
  -> validation target
```

| Trace link | 目的 |
| --- | --- |
| BDD -> Screen | 顯示使用者行為在哪些 consumer surface 上發生 |
| Screen -> Consumer Contract | 顯示畫面或 command 需要哪些資料、能力、freshness 和 failure behavior |
| Consumer need -> API / event | 顯示 provider 是否支援該需求，並暴露 API 缺口 |
| API / event -> table / store | 顯示資料所有權、持久化位置與 consistency boundary |
| Screen action -> UI Behavior / ViewModel | 顯示 interaction、state transition 與 display derivation |
| Trace link -> validation target | 顯示哪個 test、fixture、contract check 或 manual evidence 證明該行為 |

若某段 trace 尚未決定，使用 `TBD`、`out of scope`、`manual-only` 或 `not provider-backed` 明確標記，不要讓 agent 在實作時自行補齊。

## Contract-First Rules

- Screen Mapping 應在 BDD-lite 後建立，並在 API / table ownership 設計前暴露跨層缺口。
- Consumer Contract 應在 API Contract finalization 前完成或同步審查。
- UI Behavior Contract 不取代 BDD；它把 BDD 落到 screen / client state / interaction。
- Screen Contract 不應綁死 React、Vue、SwiftUI、Flutter 或特定框架；framework choice 屬 Architecture / Implementation decision。
- ViewModel Contract 擁有 display derivation；component 不應直接散落 API-to-display 規則。
- Accessibility Contract 是 UI behavior 的一部分，不是實作完成後的美化項目。
- Screen Traceability 應保持 BDD -> Screen -> Consumer -> API -> table/store -> validation target 可追溯。
- 若 UI contract 變更，同一批次更新 API / error contract、mocks、fixtures、generated clients、BDD / UI tests 或明確記錄 deferred scope。

## Minimum Validation

至少選一種驗證：

- Screen mapping / traceability review，證明 BDD scenario 能追到 screen、API、table/store 與 validation target。
- Consumer contract test 或 generated client compile check，證明 API / schema 支援 consumer needs。
- Source fixture -> ViewModel fixture 測試，證明 display derivation 穩定。
- UI behavior / component interaction test，證明 state transition、validation、feedback、navigation。
- BDD / E2E test，證明 screen 行為符合 actor intent。
- Accessibility lint、manual checklist 或 design review evidence，證明 keyboard / focus / semantics expectations。
