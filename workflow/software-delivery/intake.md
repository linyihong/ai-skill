# 需求接收 Slice（Start From Evidence / Change Intake / Parity / Brief 驗證 / Backfill）

> **Cognitive Slice**：`sd-intake`（跨 [`execution-flow.md`](execution-flow.md) §1/§6 + [`development-process.md`](development-process.md) 多個 intake gate 同批抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7）。

| slice 欄位 | 值 |
|---|---|
| `id` | `sd-intake` |
| `purpose` | 在程式碼工作前接收並穩定需求：從證據開始、變更分類、product brief 驗證、parity inventory、缺失資訊處置，以及既有專案文件回填 |
| `type` | `execution` |
| `tags` | requirements, parity, intake, domain-specific（backfill 條件子流程） |
| `load_when` | 接收新需求 / 變更 / bug / 重構 / 強化意圖、需求認知盤點、product brief 驗證、既有專案文件回填 |
| `do_not_load_when` | 已有明確 contract、純執行既定改動、evidence-only / 純分析任務 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「接收變更要先做什麼分類 / 過哪些 gate / 缺資訊如何處置」的 ordering / gate；通過 workflow membership test，不承載 evidence 取得方法（非 analysis），不論證長期模式（非 intelligence） |
| `canonical_source` | 本檔（原 `execution-flow.md` §1 從證據開始 + §6 回填規則；`development-process.md` §Initial Documentation Pack / §Product Brief Validation Gate / §Change Intake Gate / §Missing Information Gate / §Existing Project Documentation Backfill） |
| `dependencies` | [`requirements/README.md`](requirements/README.md)、[`requirements/pre-build-interrogation.md`](requirements/pre-build-interrogation.md)、[`templates/change-brief-template.md`](templates/change-brief-template.md)、`sd-contracts`（intake 後產出 contract） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Phase 4 Scenario A（execution-only：已有明確 contract 的小改時本 slice 應 **不** 載入，列於 forbidden_load） |

> **條件子流程（不另開 slice）**：Backfill for existing project（本檔 §既有專案文件回填 + §已實作專案的回填規則）掛 `sd-intake` 的 `tags: domain-specific,backfill`，僅在處理「已實作但文件缺失」的專案時觸發。Embedded / Hardware Product Flow 不在本 slice，掛 `sd-implementation`。

---

## 1. 從證據開始（Start From Evidence）

記錄可重複使用的觀察：

- 這個工作單元的目標是什麼、採取了什麼行動、以及如何驗證？如果驗證不可執行，引用參考來源和推理邊界。
- 觀察到了什麼行為？
- 哪一層暴露了它：客戶端程式碼、傳輸層、API 合約、儲存層、日誌、建置配置、韌體、硬體上下文、協定或執行時期行為？
- 問題是已確認、可疑還是僅為風險模式？

不要將目標特定的端點、令牌、機密、裝置 ID、原始使用者資料或專案事件細節複製到這個技能中。在將觀察提升為可重複使用的指引之前，應用 [`enforcement/reusable-guidance-boundary.md`](../../enforcement/reusable-guidance-boundary.md)。

### 變更接收（Change Intake）

在程式碼工作之前，執行變更接收：

| 問題 | 必要行動 |
| --- | --- |
| 存在什麼規劃產出？ | 閱讀企劃書、product brief、規劃文件、issue、ticket、PRD、設計備註、BDD、API 合約或同等文件 |
| Product brief 本身是否已驗證？ | 根據證據或明確決策檢查目標、使用者、範圍、non-goals、假設、限制、依賴、風險和成功標準。將每個主要聲明標記為 `validated`（已驗證）、`assumption`（假設）、`open question`（開放問題）、`scoped out`（排除範圍）或 `invalidated`（無效） |
| 這是新需求還是行為變更？ | 先執行 [`requirements/`](requirements/README.md) stage：product-impact discovery、behavior-driven discovery、acceptance definition、ambiguity resolution；再更新規劃文件、BDD、合約、實作切片和測試 |
| 這是 bug 修復？ | 確認預期行為 vs 實際行為、重現/證據、受影響的 BDD 或缺失 scenario、受影響的合約/錯誤和回歸測試 |
| 這是重構？ | 先分類是純內部重構、架構重組、平台遷移、工具替換或舊系統 replacement。若會替代既有功能、腳本、API、資料流程、UI flow、runtime surface 或操作流程，必須在實作前建立新舊能力 parity inventory：舊入口、現有功能、輸入、輸出 / 副作用、外部依賴、目標新入口、parity 狀態、測試 / fixture 證據與 deferred / not planned 理由。只有純內部重構且不改變 observable behavior、public contract 或操作能力時，才可只確認沒有行為變更。 |
| 這是強化？ | 確認威脅/故障模式、擁有者層、驗證和連結的檢查清單/控制更新 |
| 這是架構決策或 domain model 變更？ | 先確認 requirements stage 已有 behavior boundary / acceptance criteria / ambiguity disposition，再執行 [`architecture/architecture-fit-analysis.md`](architecture/architecture-fit-analysis.md)，確認 CRUD / DDD Lite / Full DDD / event-driven / microservices 的 fit evidence；不得預設套用 DDD、CQRS 或 event sourcing |
| 這個變更是否影響延遲、吞吐量、資源使用、啟動、背景工作、資料庫存取、批次處理或外部呼叫量？ | 在程式碼之前定義效能預算和必要的效能測試類型。不要依賴「功能正確」作為變更可發布的證明 |
| 這個變更是否與現有文件衝突？ | 應用文件優先順序：治理/框架合約、產品計劃、BDD、合約、實作、測試。更新擁有文件，而不是僅默默修正程式碼 |

**Contract / test / implementation mismatch escalation：** 若 product brief、owner contract、BDD、測試、implementation 或使用者指正彼此衝突，不要繼續局部 patch。依 `metadata/recovery/domain-policies.yaml` 的 `software-delivery` policy 進入 recovery：重讀本 workflow、artifact gates、development process、implementation plan template、linked-updates 與 dependency-reading；寫出舊假設、反證、owner contract、行為規格或 BDD、implementation surface、validation gate 與 linked updates。未完成新 execution graph 前，不可用「測試綠了」宣稱完成，也不可把 implementation 當作唯一 source-of-truth。

如果不存在規劃產出且請求會改變行為，在實作之前建立輕量的變更簡報並詢問阻擋性問題。

如果 product brief 存在但包含影響行為、合約、風險、測試、所有權、時程或發布關卡的未驗證聲明，在實作之前將這些聲明視為阻擋項。對於純規劃答案，引用參考來源或推理邊界，而不是假裝 brief 已驗證。

> **輸出模板**：Change Intake 完成後，使用 [`templates/change-brief-template.md`](templates/change-brief-template.md) 記錄變更簡報。

### Pre-build Interrogation Gate

在 change intake 之後、implementation plan 或 framework migration 之前，讀取 [`requirements/pre-build-interrogation.md`](requirements/pre-build-interrogation.md)。若請求會變成 plan、code、workflow、governance、runtime、validation、schema、generated artifact 或 tool adapter 改動，必須先記錄：

- Goal、scope、non-goals 與 expected behavior / expected framework outcome。
- Acceptance criteria 與 validation target。
- Framework discovery：canonical source、owner layer、runtime projection、mirror / cache / generated output、compiler 與 linked updates。
- Duplication risk：是否會產生第二份 rule body、第二條 activation path、stale projection 或 ambiguous source-of-truth。
- Unknown disposition：`blocker_question`、`safe_assumption`、`scoped_out` 或 `invalidated`。

若仍有會影響 behavior、contract、runtime surface、source-of-truth、validation 或安全性的 `blocker_question`，不得產生 implementation plan；先向使用者提問或停在 planning。

### Requirements Cognition Checkpoint

在進入 architecture 或 implementation 前，若任務涉及 observable behavior，讀取 [`requirements/`](requirements/README.md)：

- Product-impact discovery：Impact Map、Customer Journey Map、cross-check decision。
- Behavior-driven discovery：actor intent、behavior boundary、shared language。
- Acceptance definition：acceptance criteria、validation target、regression scope。
- Ambiguity resolution：將不確定項標成 `assumption`、`open question`、`scoped out` 或 `invalidated`。

沒有 validation target 的 acceptance criteria 不能作為完成宣告基線；requirement contradiction 或 stale acceptance criteria 需要先重建 source-of-truth。

### 重構 / Replacement Parity Gate

當變更目標是重構、遷移、改寫、替換舊工具、拆分架構、搬移 runtime surface 或建立新入口取代舊入口時，不能只寫新設計。實作前必須先產出 parity inventory，讓 reviewer 能逐項確認舊能力沒有遺漏。

最低欄位：

| 欄位 | 必填內容 |
| --- | --- |
| 舊入口 | 舊 API、command、script、UI flow、資料表、job、hook、runtime surface 或文件入口。 |
| 現有能力 | 舊入口目前支援的行為、flags、輸入、輸出、錯誤模式與邊界條件。 |
| 副作用 | 寫檔、寫 DB、發送請求、commit / push、生成 artifact、同步 mirror、修改使用者設定或其他狀態變更。 |
| 外部依賴 | runtime、shell、binary、服務、權限、平台假設、環境變數與 credentials boundary。 |
| 新入口 | 對應的新 API、command、module、adapter 或 replacement surface。 |
| Parity 狀態 | `covered`、`wrapper first`、`native target`、`deferred`、`not planned` 或 `tool-specific`，並說明原因。 |
| 驗證證據 | BDD scenario、contract test、fixture、golden output、migration assertion、manual review checklist 或明確的 blocker。 |

Blocking rule：任何舊入口若狀態為 `deferred`、`not planned` 或 `tool-specific`，必須寫明為何不阻擋目前 release / phase；任何會產生副作用的舊入口，必須有 dry-run、fake-root、fixture 或等效隔離測試。缺少 parity inventory 時，不得開始 replacement implementation，也不得宣稱新功能已覆蓋舊功能。

---

## Initial Documentation Pack（初始文件包）

當這個 workflow 被用於新功能或新專案時，agent 應協助產出第一版草稿，或針對以下文件提問缺失資訊：

| 文件 | 目的 | 如果缺失，詢問 |
| --- | --- | --- |
| Product Brief | 目標、使用者、範圍、non-goals、假設 | 這是為誰做的？解決什麼問題？明確排除什麼？ |
| Product Impact Alignment | Impact Map、Customer Journey Map、cross-check decision | 目標是什麼？誰的哪段旅程有痛點？預期改變什麼行為？功能投資是否對準痛點？ |
| Requirements Cognition Notes | actor intent、behavior boundary、ambiguity disposition | 哪些行為是明確需求？哪些只是推論、假設或 open question？ |
| BDD Behavior | 使用者/系統行為場景、acceptance criteria、validation target | 關鍵 happy path 和 failure path 是什麼？什麼證據能證明它完成？ |
| Bounded Context Map | 模組/領域拆分與所有權 | 哪些概念會一起變動？哪些邊界不應洩漏？ |
| Domain Model Contract | 核心業務物件與 invariants | 什麼必須永遠為真？允許哪些狀態轉換？ |
| Architecture Contract | 層級、依賴關係、所有權、runtime 限制 | 哪一層擁有資料、side effects、安全性、持久化、外部呼叫？ |
| API / Interface Contract | Requests、responses、events、commands、public methods | 誰消費這個 contract？相容性如何測試？版本如何管理？ |
| Error Handling Contract | Error types、recovery、user messaging、logging | 哪些錯誤可 retry、可由使用者修正、致命、或安全敏感？ |
| Hardware / Firmware Contract | Datasheet/protocol truth、electrical interface、pin/context injection、driver/service/application 邊界、target 限制 | 哪些硬體事實是固定的？每個 board 注入什麼？host/target 測試如何進行？ |
| Test Plan | Unit、BDD、contract、integration tests | 什麼證明行為、invariants 和 integration contract？ |

這些文件可以從輕量的 Markdown 草稿開始。如果專案很小，放在一個 planning file 中；如果長大，拆成一個資料夾加上 `README.md` 和聚焦的子文件。

## Product Brief Validation Gate（Product Brief 驗證關卡）

企劃書 / Product Brief 不會因為存在就自動可信。在使用它作為 BDD、contracts、估算、實作切片或測試的來源之前，先把它當作獨立的 artifact 來驗證。

### Product Impact Alignment Check

對新產品、新功能或高成本 feature investment，先用 Impact Map × Customer Journey Map 做 product alignment：

| 檢查 | 必要問題 |
| --- | --- |
| Impact Map | Why / Who / How / What 是否清楚？目標是否能追到 actor behavior change？ |
| Customer Journey Map | 需求是否對準具體 journey stage、pain point、emotional low 或 blocker？ |
| Cross-check | Impact Map 的 Who / How / What 是否真的出現在 journey 中，且資源投資對準最高摩擦點？ |
| Decision | `proceed`、`revise`、`reject` 或 `ask_user`；阻擋性缺口不得直接進 BDD / implementation |

若 business goal、target actor、journey pain 或 feature investment 不一致，先修正 product brief 或標記 `open question`。不要用完整的 BDD scenarios 掩蓋 product direction mismatch。

| Brief 項目 | 驗證問題 | 可接受的證據 |
| --- | --- | --- |
| Goal / problem | 問題是真實的、具體的、且與使用者/系統結果相關嗎？ | 使用者請求、利害關係人決策、支援 ticket、觀察到的工作流程、metric、分析發現、或明確假設 |
| Users / actors | Actors 有命名且對應到權限、角色、裝置、系統或外部服務嗎？ | 現有帳號/角色、UI/API 行為、領域文件、組織決策、或 open question |
| Scope | 能判斷現在要建什麼嗎？ | BDD scenario list、module/context map、已接受的功能列表、API/interface list |
| Non-goals | 排除的行為夠明確，能防止意外實作嗎？ | 已取消/延後/排除範圍的表格、issue decision、stakeholder 回答 |
| Assumptions | 假設是可測試的、有時效的、或標記為風險的嗎？ | 證據連結、驗證計畫、負責人、到期/審查日期 |
| Success criteria | 測試、審查、metric、效能預算或 demo 能證明它有效嗎？ | BDD acceptance criteria、contract tests、P95/P99 latency budget、throughput target、error-rate budget、release checklist、analytics/telemetry query、manual evidence |
| Constraints | 法律、安全、隱私、平台、硬體、預算、時程、相容性和營運限制有列出嗎？ | Policy、platform docs、architecture contract、risk review、hardware/vendor docs |
| Dependencies | 外部服務、vendor、團隊、generated clients、遷移、資料或硬體依賴有識別嗎？ | Integration contract、API docs、schema、vendor excerpt、migration plan、owner confirmation |
| Risks | 濫用、失敗、安全、隱私、replay、資料遺失和營運風險有列出 controls 或 blockers 嗎？ | Threat model、hardening note、controls/checklists、open blocker questions |

如果任何 brief 項目影響行為、domain invariants、API/interface 形狀、錯誤處理、安全性、儲存、所有權、測試、時程或 release gate 且無法驗證，它就是 blocker。在使用者確認、取得證據或明確排除該項目之前，不應繼續開發。

對於已實作優先的專案，根據可觀察的證據驗證回填的 Product Brief。無法恢復的 Product intent 可以保持 `unknown`，但每個已實作的行為仍需要 BDD、contract 和 test 證據。

對每個 major brief claim 使用以下狀態：

| 狀態 | 意義 | 必要行動 |
| --- | --- | --- |
| `validated` | 有證據或明確的使用者/利害關係人決策支持 | 連結證據或決策 |
| `assumption` | 合理但未經證明 | 加上負責人、驗證計畫、以及如果為假的影響 |
| `open question` | 在實作繼續前需要答案 | 提問並阻止受影響的工作 |
| `scoped out` | 明確不屬於當前工作範圍 | 記錄 non-goal 並防止意外實作 |
| `invalidated` | 證據與 brief 矛盾 | 在 code 之前修訂 brief、BDD、contracts 和 tests |

## Change Intake Gate（變更接收關卡）

在任何由這個 workflow 驅動的 code 變更之前，檢查專案的企劃書、product brief、planning docs、issue、ticket、PRD、design note、BDD、API contract 或同等專案 artifact。在實作前分類請求：

| 變更類型 | code 之前需要 |
| --- | --- |
| 新需求 / 功能 / 行為變更 | 先執行 requirements stage：product-impact discovery、behavior-driven discovery、acceptance definition、ambiguity resolution；再更新或建立 planning docs：Product Brief 或 change brief、BDD scenarios、受影響的 Domain Model Contract、Architecture Contract、API / Interface Contract、Error Handling Contract、實作切片和 tests。在 blocker questions 解決前不要開始 code |
| Bug 修復 | 確認預期行為 vs 實際行為、重現步驟或證據、受影響的 BDD scenario 或缺失 scenario、受影響的 contract/error handling、以及 regression test plan。如果修復改變了預期行為或 public contract，也視為新需求 |
| Refactor / replacement / 內部清理 | 先判斷是否取代既有功能、入口、腳本、API、資料流程、runtime surface 或操作流程。若是 replacement / migration，code 前必須建立新舊能力 parity inventory，列出舊入口、現有能力、輸入、輸出 / 副作用、外部依賴、目標新入口、parity 狀態與測試 / fixture 證據。只有純內部清理且沒有行為或 public contract 變更時，才可只記錄無行為變更；若行為、資料所有權、API、錯誤處理、安全性、儲存或 tests 改變，重新分類為新需求或 bug |
| 安全 / 強化變更 | 確認威脅或 failure mode、owner layer、必要 control、驗證方法、以及行為/API/contracts/checklists 是否需要改變 |

如果沒有 planning artifact，在實作前建立輕量的 change brief。如果請求是新需求，缺失的 planning docs 是 blockers；向使用者提問並在寫 code 前填寫 BDD/contracts。

> **輸出模板**：Change Intake 完成後，使用 [`templates/change-brief-template.md`](templates/change-brief-template.md) 記錄變更簡報。

### Refactor / Replacement Parity Inventory

當 refactor 實際上會替代舊能力時，parity inventory 是 code 前 artifact，不是 implementation 後補充。它可以放在 change brief、implementation plan 或專案專屬 inventory，但必須能被 reviewer 逐列檢查：

| 欄位 | 用途 |
| --- | --- |
| 舊入口 / 舊能力 | 定義不能遺漏的既有 API、command、script、UI flow、job、hook、資料流程或 runtime surface。 |
| 輸入 / 輸出 / 副作用 | 捕捉 flags、payload、生成物、寫入、同步、網路呼叫、commit / push 或使用者可見狀態變更。 |
| 外部依賴 | 捕捉 binary、shell、服務、平台、權限、環境變數、credential boundary 或 generated surface。 |
| 新入口 / 對照狀態 | 標記 `covered`、`wrapper first`、`native target`、`deferred`、`not planned` 或 `tool-specific`。 |
| 驗證證據 | 連到 BDD、contract test、fixture、golden output、dry-run、fake-root、manual review checklist 或 blocker。 |

任何 `deferred`、`not planned` 或 `tool-specific` 項目都要寫明不阻擋目前 phase 的理由；任何有副作用的舊能力都需要隔離測試或 dry-run 證據。

### Task Scope / Ownership Awareness Check

Before implementation or commit, separate task scope from ownership awareness. They often appear together, but they answer different questions:

- **Task scope validation** asks whether the changed surfaces still belong to the user-approved task.
- **Ownership awareness** asks whether the changed surfaces cross repo, module, security, platform, generated artifact, or team ownership boundaries that need explicit review.

Use this shape in change brief, implementation plan, or review notes when a change starts expanding:

```yaml
change_boundary_review:
  declared_task_scope:
    - <surface or behavior the task is allowed to change>
  changed_surfaces:
    - <actual files/modules/repos/workflows touched>
  task_scope_result: in_scope | expanded_with_approval | overreach | blocked
  ownership_boundaries:
    crossed:
      - <owner boundary, repo, module, generated artifact, shared library, security/payment/platform area>
    review_status: not_needed | owner_reviewed | owner_missing | blocked
  required_action:
    - narrow_diff
    - update_scope
    - ask_owner
    - split_change
    - document_deferred_scope
```

Examples:

- Editing login, registration, profile, and settings for a "fix Login page" request is task-scope overreach even if one team owns all screens.
- Editing shared auth, security framework, or payment adapters for a "fix Login API" request may be ownership-awareness failure even if the changes are technically related.

Do not use same-owner as proof of scope, and do not use task relatedness as ownership permission.

## Missing Information Gate（缺失資訊關卡）

在開發規劃或實作繼續之前，缺失資訊必須被明確處理：

| 缺失項目類型 | 必要行動 |
| --- | --- |
| 可從證據恢復 | 回填並引用證據來源 |
| Product intent 無法恢復 | 標記為 `unknown` / `open question`，向使用者提問，不要憑空創造 intent |
| 影響 BDD behavior、domain invariants、API/interface shape、error handling、security、storage 或 tests | 視為 blocker：在繼續實作前向使用者提問或要求證據 |
| 不改變行為或 contracts 的 nice-to-have 上下文 | 記錄為非阻塞的 open question，並說明為什麼不阻塞 |

不要在未解決的 blockers 下繼續開發。Agent 必須將缺失項目列為問題，等待答案或證據，然後在繼續前更新文件。

---

## Backfill（既有專案文件回填，條件子流程：`tags: domain-specific,backfill`）

> 僅在 workflow 用於「已完全或大部分實作但文件缺失」的專案時觸發。新專案不需載入本段。

### 已實作專案的回填規則（Backfill Rules for Implemented Projects）

如果專案已實作且文件缺失，在提出新的指引之前先進行文件差距審計：

| 文件 | 回填要求 |
| --- | --- |
| Product Brief | 僅回填證據支援的目標、使用者、範圍、限制和假設；將不可取得的意圖標記為 `unknown` 或 `open question` |
| BDD 行為 | **必須完成。** 從可觀察的 UI、API、程式碼、測試、日誌、fixtures 和手動驗證中完成 |
| 合約 | 從已實作的行為和證據回填領域模型、架構、API/介面、錯誤處理和測試計劃 |
| 嵌入式/硬體證據 | 從程式碼、日誌、接線筆記和測試中回填 datasheet/協定參考、硬體上下文、驅動程式/服務/應用程式邊界、主機 fixtures 和啟動證據 |
| 可追溯性 | 將產品/規則 ID 連結到 BDD，BDD 連結到程式碼引用，BDD 連結到測試，API/命令/診斷合約連結到 fixtures，以及生成的客戶端連結到來源合約 |

不要讓缺失的 product brief 阻擋已實作產品的 BDD 回填。

對於先實作優先的專案，也要恢復交付管線：來源產品文件或計劃雷達、文件優先順序、最小文件同步矩陣、OpenAPI/schema/codegen 流程、供應商整合摘錄，以及明確的已取消/延後/排除範圍的決策。

影響行為、領域不變量、API/介面形狀、錯誤處理、安全性、儲存、所有權或測試的缺失資訊是阻擋項。向使用者提問或要求證據，用答案更新文件，然後才繼續開發規劃或實作。非阻擋性的未知項必須標記為什麼它們不改變行為或合約。

對於嵌入式或硬體支援的產品，缺失的 datasheet/協定真相、電氣介面、引腳/匯流排映射、硬體上下文所有權、時間/並發限制、安全行為、fixture 來源或目標驗證方法也是阻擋項，除非明確排除範圍。

### Existing Project Documentation Backfill（既有專案文件回填）

當這個 workflow 被用於已經完全或大部分實作的專案時，先審查既有文件並回填任何缺失的開發文件。不要因為實作已經存在就跳過流程。

| 缺失文件 | 回填規則 |
| --- | --- |
| Product Brief | 只重建證據支持的內容：可見的目標、使用者/actors、範圍、non-goals、假設和限制。如果原始 intent 不可取得，將欄位標記為 `unknown` 或 `open question`；不要憑空創造業務理由 |
| Bounded Context Map | 從 code ownership、runtime boundaries、database tables、API groups、UI areas、queues、SDK/public APIs 和 deployment units 推斷模組 |
| BDD Behavior | **必須完成。** 從已實作的產品、tests、UI、API 行為和 logs 重建關鍵 happy paths、failure paths、權限、空狀態、edge cases 和跨 context 流程 |
| Domain Model Contract | 從 code、schemas、storage、UI states 和 tests 推斷 entities、value objects、commands、events、invariants 和 state transitions；將不確定的詞彙標記為 candidate |
| Architecture Contract | 記錄實際的 dependency direction、data ownership、side-effect boundaries、integrations、runtime/deployment shape 和已知違規 |
| API / Interface Contract | 提取實際的 request/response schemas、public methods、events、commands、auth/session behavior、versioning、compatibility、fixtures 和 consumers |
| Error Handling Contract | 回填觀察到的 error taxonomy、retry rules、user messages、logging/redaction behavior、security-sensitive failures 和 gaps |
| Test Plan | 將既有 tests 對應到 behavior/contracts，並列出未覆蓋的 BDD scenarios、invariants、contracts 和 integration paths 所需的 tests |

對於已實作優先的專案，也要恢復 delivery pipeline：

| Pipeline artifact | 回填規則 |
| --- | --- |
| Plan index / product radar | 將來源 product docs、PDFs、tickets、screenshots 或 legacy notes 對應到 modules、controllers、screens、commands 或 packages。標記已取消或已取代的需求 |
| Contract taxonomy | 列出哪些文件管轄 build/run、HTTP/API shape、auth/tenant/session、persistence、domain layering、frontend/backend integration、third-party integration、testing 和 documentation sync |
| Minimum doc sync matrix | 對每個變更類型，說明最少要更新的 docs/tests：API、permission、database、UI flow、generated client、vendor integration、CLI command、diagnostic rule、release setting |
| OpenAPI / schema / generated client | 驗證 generated consumer code 來自 source contract，而非手抄 endpoints 或 DTOs |
| Vendor / third-party integration | 區分 raw vendor docs 與 sanitized integration excerpts、request/response contracts、fixture examples、live-test gates 和 secret handling |
| Tooling / extension rule catalog | 對應 catalog order、rule IDs、diagnostics/commands、fixtures 和 tests；明確標記 process-only 或 non-enforceable rules |

既有專案的回填順序：

1. 盤點既有 docs、source folders、tests、schemas、API specs、fixtures、release notes 和觀察到的行為
2. 建立 documentation gap table，狀態為：`exists`、`partial`、`missing` 或 `unknown`
3. 當 product brief 缺失時先回填 BDD Behavior，因為已實作的行為是最強的可用 truth source
4. 從已完成的行為和實作證據回填 Domain Model、Architecture、API / Interface 和 Error Handling Contracts
5. 將未知的 product intent 與觀察到的行為分開標記。未知的 intent 不阻塞 BDD 完成
6. 如果 BDD 無法從可用證據完成，停止並要求缺失的行為、screen/API 範例、logs、test cases 或使用者決策，然後再繼續開發
7. 對任何缺乏覆蓋率的關鍵 BDD scenario 加上 tests 或 test TODOs
