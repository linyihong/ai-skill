# Pre-build Interrogation

## 目的

Pre-build Interrogation 是 plan 或 implementation 前的需求拷問 gate。它把使用者請求先轉成可回答、可驗證、可拒絕的問題，避免 agent 在需求、框架邊界或 source-of-truth 尚未釐清時直接產生 implementation plan。

此流程承接 mattpocock/skills 的 `/grill-me` 精神，但在 Ai-skill 中保持 tool-neutral：它不是 slash command，而是 workflow gate 與 executable contract。

## 觸發時機

符合任一條件時，先執行本 gate：

- 使用者要求 `plan`、`implement`、`build`、`refactor`、`migration`、`architecture`、`workflow`、`governance` 或 framework 改動。
- 請求會改變 observable behavior、public contract、runtime surface、workflow、validation gate、schema、tool adapter 或 generated artifact。
- 任務描述含糊，只描述想法或方向，尚未說明成功標準、非目標、驗證方式或風險。
- 改動可能建立第二份規則、第二個 activation path、runtime projection、mirror、cache、generated report 或 source-of-truth。
- 使用者指出先前產生了雙寫、漏驗證、錯 source、規則衝突或 framework drift。

## 問題分類

### 1. Intent

確認真正目標：

- 使用者要解決什麼問題？
- 成功後誰受益，或哪個 agent failure 會被防止？
- 這是新能力、修 bug、重構、治理補強、runtime migration，還是 validation 補洞？
- 若只做最小改動，必須達成哪個結果？

### 2. Scope

劃定邊界：

- 哪些行為、文件、API、workflow、runtime table、generated surface 或工具 adapter 在範圍內？
- 哪些明確不做？
- 是否需要保留相容路徑？相容期限或例外是什麼？
- 是否有專案特定 evidence 不能進 reusable docs？

### 3. Behavior And Acceptance

把需求轉成可驗證結果：

- 目前行為是什麼？期望行為是什麼？
- 哪些 edge cases、negative paths 或 failure modes 必須處理？
- 用什麼證據宣稱完成：test、runtime validate、scenario、diff review、manual review、SQLite query、link check？
- 哪些 acceptance criteria 沒有 validation target，不能當完成基線？

### 4. Framework Discovery

當改動 Ai-skill framework、governance、runtime、workflow、metadata、validation 或 tool adapter 時，必問：

- Canonical source 在哪一層？是 owner Markdown/YAML、SQLite canonical document、compiler source、generated report，還是 tool adapter？
- 改動涉及的關鍵詞彙是否已在 [`knowledge/glossary/ai-skill.md`](../../../knowledge/glossary/ai-skill.md) 有 canonical entry？若有，必須引用 owner-layer 定義；若無但屬 framework / runtime / cognitive vocabulary，應評估補進 glossary 而非 inline redefine。詞義衝突依 [`knowledge/glossary/README.md`](../../../knowledge/glossary/README.md) §Vocabulary Resolution Priority 解析。
- Runtime row 是 source 還是 projection？
- 是否存在舊入口、mirror、cache、generated output、compatibility table 或 archived plan 仍在描述同一件事？
- 修改後哪些 README、routing registry、inventory、validation scenario、compiler、runtime DB 或 tool config 需要同步？
- 如果不移除舊路徑，未來 agent 會不會讀到兩份互斥規則？

### 5. Duplication Risk

下列任一項為 blocking，除非 plan 明確處理：

- 同一條 rule 同時存在兩份 executable semantics。
- Markdown、YAML、runtime table 或 generated surface 各自維護不同 activation / gates / final report。
- SQLite 或 generated report 被當成 source 修改，owner-layer source 沒有同步。
- 新舊 entrypoint 同時 active，但沒有 deprecation、routing priority 或 parity explanation。
- Validation 只證明新 surface 存在，沒有證明舊 duplicate 已刪除或降級。

### 6. Unknowns

把不確定項分類：

| 類型 | 處理 |
| --- | --- |
| `blocker_question` | 先問使用者；不得進入受影響 implementation。 |
| `safe_assumption` | 可繼續，但必須記錄假設、風險與驗證方法。 |
| `scoped_out` | 明確寫入非目標與不阻擋理由。 |
| `invalidated` | 修正 plan；不得沿用舊假設。 |

## Ready Gate

進入 implementation plan 前，至少要能回答：

1. Goal、scope、non-goals 與 expected behavior 已明確。
2. Acceptance criteria 都有 validation target。
3. Framework 改動已識別 canonical source、projection、mirror、generated artifact 與 linked updates。
4. Duplicate source-of-truth risk 已排除，或 plan 明確包含移除 / deprecation / compatibility strategy。
5. Open questions 已分類；blocking questions 已問完或任務停在 planning。

## 輸出格式

在 plan、change brief 或回覆中記錄：

```markdown
## Pre-build Interrogation
- Goal:
- Scope:
- Non-goals:
- Acceptance / validation target:
- Framework discovery:
- Duplication risk:
- Open questions:
- Assumptions:
- Decision: proceed | ask_user | revise_plan | blocked
```

## 與其他 requirements stage 的關係

Pre-build Interrogation 是進入 `product-impact-discovery/`、`behavior-driven-discovery/`、`acceptance-definition/` 與 `ambiguity-resolution/` 前的 gate。它不取代 BDD-lite 或 product alignment；它先確認 agent 是否知道該問什麼、哪些問題會阻擋 plan、哪些可以變成明確假設。
