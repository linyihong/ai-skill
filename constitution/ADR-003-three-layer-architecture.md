# ADR-003: Three-Layer Architecture（Knowledge / Skills / Intelligence）

## Status

**Accepted**

## Framework Generation

- **世代分類**：Gen 2 / Gen 3 transition（**詞彙演化注意**，見下方 Vocabulary Evolution section）
- **當前世代文件**：[`architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md)
- **適用狀態**：「三層分離」的**核心精神**（事實 / 流程 / 判斷）在 Gen 3 仍成立，但**具體層級已演化**：`skills/` 已 deprecated，內容遷移到 `workflow/` 與 `analysis/`。本 ADR 正文保持 immutable，演化補在下方 Vocabulary Evolution。

## Context

本 repository 的內容可以分為三種本質不同的類型：

1. **Knowledge（事實）**：Redis 支援 pub/sub、HTTP 200 意義、CQRS 定義。
2. **Skills（流程）**：APK 分析執行步驟、Code review checklist、旅行規劃流程。
3. **Intelligence（判斷）**：何時該用 CQRS、如何判斷 Redis 連線池問題、效能與可維護性的取捨。

最初所有內容混在 `skills/` 下，導致：
- Skill 文件同時包含流程步驟、事實查詢表、經驗判斷。
- 跨 skill 重用困難（同樣的 intelligence 分散在不同 skill 中）。
- Agent 載入 skill 時被迫載入所有內容，無法只載入需要的部分。

需要一個架構讓這三種內容可以獨立管理、獨立演化、獨立載入。

## Decision

採用 **Three-Layer Architecture**，將內容分為三個平行層：

```text
knowledge/    ← 事實層（What is）
skills/       ← 流程層（How to）
intelligence/ ← 判斷層（When / Why / Trade-off）
```

### 各層核心責任

| 層 | 核心問題 | 變動頻率 | Lifecycle | Token 策略 |
| --- | --- | --- | --- | --- |
| `knowledge/` | What is X? | 低 | stable → deprecated | 可壓縮為 summary |
| `skills/` | How to do X? | 中 | active → updated → archived | 保留執行步驟 |
| `intelligence/` | When / Why / Trade-off? | 中高 | candidate → validated → promoted → updated | 保留決策邏輯 |

### 跨層關係

- **Skills 引用 Knowledge**：執行流程中需要查詢事實（如 HTTP 狀態碼）。
- **Skills 引用 Intelligence**：執行流程中需要判斷（如選擇分析路線）。
- **Knowledge 不引用 Skills**：事實獨立於使用方式。
- **Intelligence 可引用 Knowledge**：經驗判斷基於事實。
- **Intelligence 可引用 Skills**：經驗判斷可建議特定執行流程。

### 與既有 `analysis/`、`workflow/`、`feedback/` 的關係

- `analysis/` 是分析方法（屬於 knowledge 面向的事實性方法論）。
- `workflow/` 是可執行流程（屬於 skills 面向）。
- `feedback/` 是演化引擎（驅動三層的內容更新與 promotion）。

## Consequences

### 正面

- **關注點分離**：每層只負責一種本質，內容更純淨。
- **精準載入**：agent 可以只載入需要的層，減少 token 浪費。
- **獨立演化**：knowledge 可穩定不變，intelligence 可持續更新，skills 可隨工具改變而調整。
- **跨 skill 重用**：intelligence atom 可被多個 skill 引用。

### 負面

- **初期遷移成本**：需要將現有 `skills/` 內容拆解到三層。
- **跨層連結維護**：需要 governance 確保跨層引用不中斷。
- **學習曲線**：新貢獻者需要理解三層的差異才能正確放置內容。

## Alternatives Considered

- **兩層架構（knowledge + skills）**：將 intelligence 放在 knowledge 下。但 intelligence 的 lifecycle 與使用方式與 knowledge 差異太大。不採用。
- **單層架構（全部放在 skills 下）**：維持現狀。但無法解決跨 skill 重用與精準載入問題。不採用。
- **四層架構（加上 experience 層）**：將 feedback/replay 獨立。但 replay 本質上是 intelligence 的原料，不是獨立層。不採用。

## Vocabulary Evolution

> 本 section 不修改上方 immutable 正文，僅標註 Gen 2 → Gen 3 的詞彙演化。決策的**核心精神**（事實 / 流程 / 判斷三層分離）保留；具體層級與路徑名稱已演化。

| 原文用詞 | Gen 3 對應 | 說明 |
|---------|----------|------|
| `skills/`（流程層） | `workflow/` + `analysis/` | Skills 層在 Gen 3 已 deprecated；原 skills 內容拆解為「執行流程 → workflow/」與「分析方法 → analysis/」 |
| `skills/ADDING_SKILLS.md` | （已移除） | 對應的 capability 入口已分散到 `workflow/<domain>/README.md` 與 `analysis/<domain>/README.md` |
| 「Skills 引用 Knowledge / Intelligence」 | 「Workflow / Analysis 引用 Knowledge / Intelligence」 | 跨層引用關係保留，但發起層改為 workflow + analysis |
| skill lifecycle（active → updated → archived） | workflow / analysis 自有 lifecycle | 由 `governance/lifecycle/` 統一管理 |

**核心精神保留**：
- 事實（knowledge）/ 流程（workflow + analysis）/ 判斷（intelligence）的三向分離仍是 Gen 3 架構基礎
- 跨層引用方向不變（流程 → 事實 + 判斷）
- 各層獨立 lifecycle 不變

**為什麼不 supersede 本 ADR**：核心決策（三層分離）仍成立，僅標籤變化。Supersession 保留給「決策本身被推翻」的情況；參考 [ADR-007](ADR-007-constitution-and-decision-promotion-boundary.md) 的 promotion boundary 原則。

## Related

- [`intelligence/README.md`](../intelligence/README.md) — intelligence 層定義
- [`knowledge/README.md`](../knowledge/README.md) — knowledge 層定義
- [`workflow/README.md`](../workflow/README.md) — Gen 3 workflow 層（取代部分 skills/ 內容）
- [`analysis/README.md`](../analysis/README.md) — Gen 3 analysis 層（取代部分 skills/ 內容）
- [`plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md`](../plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md) — 整體架構規劃
- [`feedback/README.md`](../feedback/README.md) — 演化引擎
