# Feedback 與技巧條目（共用規則）

「怎麼寫回饋、檔名怎麼取、模板長怎樣」**全部**在本檔維護；各 skill **不再**另存一份 `FEEDBACK.md` 正文（`apk-analysis` 目錄下僅保留極短入口檔，指向本檔）。

**每一條 lesson 全文**統一放在 **`feedback/history/<domain>/`** 之下（對應 skill 領域名稱，如 `apk-analysis/`、`app-development-guidance/`、`travel-planning/`），**不要**再把 lesson 分散到 `workflow/<name>/feedback_history/`、`analysis/<name>/feedback_history/`、`intelligence/<name>/feedback_history/` 或 `skills/<name>/feedback_history/` 等舊路徑。

## 原則

- **不要**在每條 lesson 裡重複貼上 [sanitization.md](sanitization.md)、[dependency-reading.md](dependency-reading.md)、[neutral-language.md](neutral-language.md)、[goal-action-validation.md](goal-action-validation.md)、[authorization-scope.md](authorization-scope.md) 等全文；條目頂部用一行**引用** [README.md](README.md) 或**本檔**即可。
- **Agent 行為：** 在授權分析過程中一旦得到可重用技巧／失敗模式／驗證規則，應**主動**在同一輪對話內於 **`feedback/history/<domain>/`** 對應分類新增檔案（依下方**分類規則**、**檔名規則**與**模板**），**不要**等使用者提醒。
- **每輪回饋檢查：** 每個有實質進展的工作回合結束前、切回長時間專案工作前、提交 project-only evidence 前、或使用者說「繼續」展開下一輪前，agent 必須明確自問：本輪是否新增可重用技巧、validation rule、replay knob、hook/runner guard、錯誤模式、或閉環缺口？若是，先開啟 canonical `<AI_SKILL_REPO>` writeback transaction；若否，最終回覆可簡短說明本輪沒有新增泛化 lesson，或只留下 project-specific evidence。
- 只寫**通用方法**，不寫特定 App / project incident 的私有結論或具體證據；必須依 [reusable-guidance-boundary.md](reusable-guidance-boundary.md) 先抽象化原因、規則與驗證，再依 [sanitization.md](sanitization.md) 去敏；若 lesson 來自 agent 失誤或閉環不完整，另依 [failure-learning-system.md](failure-learning-system.md) 分類失效模式並判斷是否需要 `shared-rules/failure-patterns/`；不要只新增 skill-local feedback lesson 就宣稱已吸收 agent 錯誤，需先判斷是否也需要 cross-skill failure pattern；發現 skill/rule/template/lesson 更新時依 [dependency-reading.md](dependency-reading.md) 讀取依賴；標題、slug、摘要與正文必須依 [neutral-language.md](neutral-language.md) 使用中性低爭議用語；每個重要結論必須依 [goal-action-validation.md](goal-action-validation.md) 說明目標、執行、驗證或參考來源；必須說明證據與適用／不適用條件；不確定標 `experimental`。
- **去敏檢查：** 寫入任何 lesson 前，必須依 [`sanitization.md`](sanitization.md) 逐項檢查：不得包含本機絕對路徑（改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 占位符）、使用者名稱、私有工作目錄、clone 位置、secrets、raw tokens、私人 host、個資或 project-specific evidence。寫入完成後在 `git diff` 中再次確認去敏項目已處理。

## 條目放哪裡

| 內容 | 位置 |
| --- | --- |
| **共用政策（全庫）** | [`shared-rules/README.md`](README.md) |
| **本檔** | 命名規則、模板、索引與 Git 約定（**唯一正文**） |
| **每一條獨立 lesson（所有 skill，無論遷移狀態）** | **`feedback/history/<domain>/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`** |
| **條目總覽表**（可選） | **`feedback/history/<domain>/README.md`** 與必要的 **`feedback/history/<domain>/<category>/README.md`** |

> **重要：** `feedback/history/` 是 lesson 的唯一目標路徑。舊的 `skills/<name>/feedback_history/`、`workflow/<name>/feedback_history/`、`analysis/<name>/feedback_history/`、`intelligence/<name>/feedback_history/` 路徑已於 2026-05-13 刪除（`skills/` 下的 `feedback_history/` 已全部搬遷完畢），**新 lesson 一律寫入 `feedback/history/<domain>/`**。

### 判斷流程

1. **確認 lesson 的 domain 歸屬**：這個技巧屬於哪個 skill 的 scope？
   - APK 分析技術（Frida hook、TLS capture、proxy 架構、response decoding）→ `apk-analysis`
   - 開發指引（SDK 設計、API contract、BDD、測試策略、release gate）→ `app-development-guidance`
   - 旅遊規劃（行程、交通、住宿、預算）→ `travel-planning`
   - 若不確定，先讀取該 skill 的 `SKILL.md` 確認 scope 描述。
2. **確認 domain 下的分類**：對應 `feedback/history/<domain>/` 下是否有對應分類目錄（如 `common/`、`flutter-dart-aot/`、`controls/`）。
3. **寫入 `feedback/history/<domain>/<category>/`**，若尚無對應分類目錄，**應主動建立**，而非退回舊路徑。

### 範例

- **APK 分析 lesson**：`feedback/history/apk-analysis/common/2026-05-13_094500-anti-bot-gateway-blocks-external-sdk.md`
- **開發指引 lesson**：`feedback/history/development-guidance/controls/2026-05-01_142100-client-encrypted-header-not-boundary.md`
- **旅遊規劃 lesson**：`feedback/history/travel-planning/common/2026-05-13_094500-xxx.md`

### ⚠️ 常見錯誤（必讀）

以下錯誤案例來自實際 feedback lesson 寫入失誤，新 agent 在建立 lesson 前應先閱讀，避免重蹈覆轍：

#### ❌ 錯誤 1：混淆 lesson 的 domain 歸屬

**情境**：從 TATA 專案學到「非標準 TLS bypasses SSL hooks」和「dart:io HttpClient bypasses Java hooks」

**錯誤行為**：把這兩個 lesson 放到 `feedback/history/development-guidance/common/`

**為什麼錯**：這兩個是 APK 分析技術（Frida hook 策略），不是開發指引。`app-development-guidance` 的 scope 是 SDK 設計、API contract、BDD、測試策略。

**正確做法**：
- 先確認 lesson 屬於哪個 domain：非標準 TLS hook 是 APK 分析技術 → `apk-analysis`
- 最終位置：`feedback/history/apk-analysis/local-proxy/`

#### ❌ 錯誤 2：已遷移 skill 的 lesson 放到舊路徑

**情境**：`apk-analysis` 已遷移至新分層，但 agent 仍把 lesson 放到 `skills/apk-analysis/feedback_history/` 或 `workflow/apk-analysis/feedback_history/`

**錯誤行為**：使用舊路徑（無論是 `skills/` 還是 `workflow/`、`analysis/`、`intelligence/` 底下的 `feedback_history/`）

**為什麼錯**：所有 lesson 的統一目標路徑是 `feedback/history/<domain>/`，不再按 lesson 性質分散到各層。

**正確做法**：`feedback/history/apk-analysis/local-proxy/`

#### ✅ 正確範例：完整的判斷流程

**Lesson 主題**：「External JVM SDK 無法直接呼叫受 PerimeterX 保護的 API」

**判斷流程**：
1. 這是一個開發指引（SDK 設計限制），不是 APK 分析技術 → `app-development-guidance`
2. 確認 domain 分類：`feedback/history/development-guidance/common/`
3. ✅ 最終位置：`feedback/history/development-guidance/common/2026-05-13_094500-anti-bot-gateway-blocks-external-sdk.md`

> 更多 failure pattern 請見 [`shared-rules/failure-patterns/skill-classification-boundary-confusion.md`](failure-patterns/skill-classification-boundary-confusion.md)。

成熟後可將 lesson 整理進 `workflow/<domain>/execution-flow.md`、`analysis/<domain>/` 或 `intelligence/<domain>/`（見模板中 **Promotion Target**）。

## 分類規則

當某個 domain 內部已經開始按 runtime、platform、control、technique、checklist 等方式分類時，`feedback/history/<domain>/` 也要跟著分類，避免所有 lesson 混在同一層：

- 新 lesson 優先放到 **`feedback/history/<domain>/<category>/`**，其中 `<category>` 應對應該 domain 內既有分類名稱，例如 `flutter-dart-aot`、`http-api`、`controls`、`common`。
- 跨分類、全域適用或分類尚未明確的 lesson 放到 **`feedback/history/<domain>/common/`**。
- 若一條 lesson 會 promote 到多個分類，放在主要分類，並在 lesson 的 **Promotion Target** / **Required Linked Updates** 寫出其他同步更新位置。
- `feedback/history/<domain>/README.md` 應是總索引，列出 category folders；每個 category folder 可有自己的 `README.md` 表格。
- ✅ 既有歷史 lesson 已全部搬遷至 `feedback/history/<domain>/`，舊路徑 `skills/<name>/feedback_history/` 已於 2026-05-13 刪除。
- 不要為了分類而重複複製 lesson 內容；一條 lesson 只保留一份全文，其他地方用連結。

## 檔名規則（時間 + `<slug>`）

- 使用 **`YYYY-MM-DD_HHMMSS-<slug>.md`**：
  - **`YYYY-MM-DD`**：建立 lesson 的日期（本機）。
  - **`HHMMSS`**：**24 小時制**本機時間（6 位數字，例：`143052` = 14:30:52）。含時間可避免同日多檔碰撞、也方便依檔名排序。
- `<slug>` 建議 **短英文 kebab-case** 或 **有意義的英數縮寫**（例：`proxy-two-layer-tls`、`aapt-resolve-activity`）；中文標題可保留但不宜過長。
- **同一秒多條**：微調秒數或改 `<slug>`，勿覆寫既有檔。
- **修改既有 lesson**：在原檔 **追加修訂說明**（簡短段落）或建新檔並在舊檔頂部標 `deprecated → 見 xxx.md`；不要默默刪除歷史。

## 新 lesson 模板

複製到新檔：

- **所有 lesson（統一目標路徑）**：**`feedback/history/<domain>/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`**
- 若 domain 下尚無對應分類目錄，**應主動建立**，而非退回舊路徑。

注意：以下引用路徑以檔案位於 `feedback/history/<domain>/<category>/` 內為準；若檔案在 `feedback/history/<domain>/common/`，共用規則連結路徑相同。

```markdown
> 遵守 [共用規則索引](../../../shared-rules/README.md)、[dependency-reading](../../../shared-rules/dependency-reading.md)、[neutral-language](../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### YYYY-MM-DD - [short title]

Status: candidate | validated | deprecated | promoted | experimental

#### One-line Summary

用一句人話說明這條 lesson。

#### Human Explanation

給人看的背景說明：為什麼重要、常見誤判是什麼、實務上怎麼判斷。

#### Trigger

遇到什麼現象或問題？

#### Evidence

- Tool:
- Sanitized excerpt: 只寫可公開、不可識別單一專案或使用者的摘要；具體 project incident 證據留在專案文件。
- Evidence path: 使用 `<PROJECT_ROOT>` 等 placeholder，或引用 project docs 的相對位置；不要寫本機實體路徑。

#### Generalized Lesson

可重用的規則是什麼？

#### Agent Action

下次 agent 看到類似情境時，應該先做什麼、不要做什麼？

#### Goal / Action / Validation

- Goal:
- Action:
- Validation or reference source:

#### Applies When

- 條件 1

#### Does Not Apply When

- 條件 1

#### Validation

如何確認這條 lesson 是對的？

#### Promotion Target

- `workflow/<domain>/execution-flow.md`
- `analysis/<domain>/`
- `intelligence/<domain>/`
- `shared-rules/`（若 lesson 適合提升為全庫規則）

#### Required Linked Updates

- 依 [`linked-updates.md`](linked-updates.md) 列出必須同步更新或已檢查的相關文件；若無需連動更新，寫明原因。
- 若 lesson 來自 project incident，列出已依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) 檢查：skill 只保留 generalized lesson，具體證據留 project docs。
```

## 同步與索引

- **索引**：維護者可定期整理 **`feedback/history/<domain>/README.md`**（表格：檔名、Status、標題、一句話摘要，或 category index）；若該 domain 已有索引，agent 新增 lesson 檔後**必須**追加表格列或明確說明為何暫不更新。已分類 domain 同時更新對應 **`feedback/history/<domain>/<category>/README.md`**。
- **Git**：**`feedback/history/`** 版本控制；不要提交含機密的原始 log。
- **歷史**：✅ `skills/<name>/feedback_history/` 下的既有 lesson 已於 2026-05-13 全部搬遷至 `feedback/history/<domain>/`，舊目錄已刪除。新 lesson 一律寫入 `feedback/history/`。

← [回到共用規則索引](README.md)
