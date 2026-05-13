# Skill Classification Boundary Confusion（回饋分類邊界混淆）

Status: candidate
Class: `scope-drift`

## Trigger

當以下情況同時發生時，agent 可能誤判 feedback lesson 的放置位置：

1. Agent 需要為多個 skill 建立 feedback lessons（例如同時處理 `app-development-guidance` 和 `apk-analysis` 的回饋）
2. 其中一個 skill 已遷移至新分層（`workflow/<name>/`、`analysis/<name>/`、`intelligence/<name>/`），另一個尚未遷移
3. Agent 沒有先檢查每個 skill 的 `SKILL.md` 是否有「新分層路徑（優先讀取）」章節
4. Agent 混淆了「lesson 的技術主題」與「lesson 所屬的 skill 邊界」

## Failure Mode

Agent 把某個 skill 的 analysis technique lesson 放到另一個不相關 skill 的 feedback_history 中，因為：

- **錯誤假設**：認為「lesson 的技術主題」決定放置位置（例如 TLS hook 技巧 → `app-development-guidance`）
- **正確規則**：lesson 的放置位置由「該 lesson 屬於哪個 skill 的 scope」決定（例如 TLS hook 技巧 → `apk-analysis`，因為這是 APK 分析技術）
- **忽略遷移狀態**：沒有檢查 skill 是否已遷移至新分層，導致 lesson 被放到舊 `skills/<name>/feedback_history/` 而非新分層路徑

### 具體錯誤模式

| 錯誤 | 正確 |
|------|------|
| 把 APK 分析技術（非標準 TLS、dart:io HttpClient）放到 `app-development-guidance` 的 feedback_history | 這些是 APK 分析技術，應放 `apk-analysis` 對應的分層路徑 |
| 把開發指引（anti-bot gateway、SDK 設計模式）放到 `apk-analysis` 的 feedback_history | 這些是開發指引，應放 `app-development-guidance` 對應的分層路徑 |
| 已遷移 skill 的 lesson 仍放 `skills/<name>/feedback_history/` | 已遷移 skill 的 lesson 應放新分層路徑（`workflow/`、`analysis/`、`intelligence/`） |

## Risk

- Lesson 放在錯誤的 skill 中，後續 agent 在該 skill 工作時讀不到相關技巧
- 跨 skill 的 lesson 分散在各處，降低可發現性
- 使用者需要手動糾正，浪費時間
- 已遷移 skill 的舊路徑 lesson 不會被新分層的 agent 讀取到

## Required Agent Action

在建立 feedback lesson 之前：

1. **檢查每個 skill 的 `SKILL.md`**：確認是否有「新分層路徑（優先讀取）」章節
   - 有 → 已遷移，lesson 放新分層路徑
   - 無 → 尚未遷移，lesson 放舊 `skills/<name>/feedback_history/`
2. **確認 lesson 的 skill 歸屬**：這個技巧屬於哪個 skill 的 scope？
   - APK 分析技術（Frida hook、TLS capture、proxy 架構）→ `apk-analysis`
   - 開發指引（SDK 設計、API contract、BDD、測試策略）→ `app-development-guidance`
3. **確認 lesson 的層級歸屬**（僅已遷移 skill）：
   - 執行流程相關 → `workflow/<name>/feedback_history/`
   - 分析方法相關 → `analysis/<name>/feedback_history/`
   - 工程智慧相關 → `intelligence/<name>/feedback_history/`
4. **如果對分類不確定**，先讀取 [`feedback-lessons.md`](../feedback-lessons.md) 的「條目放哪裡」和「判斷流程」章節

## Prevention Gate

在建立 feedback lesson 之前，必須能回答：

| Check | Required answer |
| --- | --- |
| Skill 遷移狀態 | 該 skill 的 `SKILL.md` 是否有「新分層路徑（優先讀取）」？ |
| Lesson 的 skill 歸屬 | 這個技巧屬於哪個 skill 的 scope？（apk-analysis / app-development-guidance / 其他） |
| Lesson 的層級歸屬 | 屬於 workflow / analysis / intelligence 哪一層？（僅已遷移 skill） |
| 路徑確認 | 最終路徑是否符合 `feedback-lessons.md` 的規則？ |

## 驗證

- 所有新 lesson 都放在正確的 skill 路徑下
- 已遷移 skill 的 lesson 不在舊 `skills/<name>/feedback_history/` 中
- 未遷移 skill 的 lesson 不在新分層路徑中
- `feedback-lessons.md` 的「判斷流程」章節有明確的錯誤案例

## Linked Rules

- [`../feedback-lessons.md`](../feedback-lessons.md)（條目放哪裡、判斷流程）
- [`../failure-learning-system.md`](../failure-learning-system.md)（失效學習系統）
- [`../scope-drift` 類別](../failure-learning-system.md#failure-taxonomy)

## Linked Validation Scenarios

- `validate_intelligence_classification_boundary` — 檢查 `intelligence/README.md` 的結構圖與實際目錄是否一致，防止新目錄加入時沒有正確分類

## Error Examples

### ❌ 錯誤範例 1：APK 分析技術放到 app-development-guidance

**情境**：從 TATA 專案學到「非標準 TLS bypasses SSL hooks」和「dart:io HttpClient bypasses Java hooks」

**錯誤行為**：把這兩個 lesson 放到 `intelligence/engineering/app-development-guidance/feedback_history/common/`

**為什麼錯**：這兩個是 APK 分析技術（Frida hook 策略），不是開發指引。`app-development-guidance` 的 scope 是 SDK 設計、API contract、BDD、測試策略。

**正確位置**：
- `feedback/history/apk-analysis/local-proxy/`（非標準 TLS — APK 分析技術）
- `feedback/history/app-development-guidance/common/`（dart:io HttpClient — 這個 lesson 同時涉及開發指引，因為它影響 SDK 的 HTTP client 選擇）

### ❌ 錯誤範例 2：已遷移 skill 的 lesson 放到舊路徑

**情境**：`apk-analysis` 已遷移至新分層（`SKILL.md` 有「新分層路徑（優先讀取）」），但 agent 仍把 lesson 放到 `skills/apk-analysis/feedback_history/`

**錯誤行為**：使用舊 `skills/<name>/feedback_history/` 路徑

**為什麼錯**：所有 lesson 的統一目標路徑是 `feedback/history/<domain>/`，不再按 lesson 性質分散到各層。

**正確位置**：`feedback/history/apk-analysis/local-proxy/`（APK 分析技術）

### ✅ 正確範例：判斷 lesson 歸屬

**Lesson 主題**：「External JVM SDK 無法直接呼叫受 PerimeterX 保護的 API，因為 TLS fingerprint 驗證會阻擋非 App 的 TLS stack」

**判斷流程**：
1. 這是一個開發指引（SDK 設計限制），不是 APK 分析技術
2. `app-development-guidance` 已遷移 → 使用新分層路徑
3. 這是工程智慧（理解 anti-bot 的行為邊界）→ `intelligence/engineering/app-development-guidance/feedback_history/`
4. ✅ 最終位置：`intelligence/engineering/app-development-guidance/feedback_history/common/2026-05-13_094500-anti-bot-gateway-blocks-external-sdk.md`

← [Back to failure patterns](README.md)
