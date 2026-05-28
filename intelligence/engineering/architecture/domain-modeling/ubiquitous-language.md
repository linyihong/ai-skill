# Ubiquitous Language

**Status**: `candidate-intelligence`

## 判斷原則

Ubiquitous language 是 DDD Lite 也應保留的最小能力。它讓 product brief、BDD、domain model、API contract、測試與實作使用同一組可驗證語言。

## 適用訊號

- 同一需求在產品、設計、後端、前端、測試中有不同名稱。
- Bug 來自詞彙誤解，而不是程式錯誤。
- BDD scenario 與 domain model 對不上。
- 文件中的名詞無法對應到 code 或 test。

## 行動

1. 建立 domain glossary，只收會影響行為的詞。
2. 每個詞記錄狀態、責任 context、允許操作與反例。
3. 在 BDD / contract / implementation plan 中使用同一詞彙。
4. 若詞在不同 context 意義不同，回到 bounded context analysis。

## 避免

不要把 glossary 變成百科。只記錄會影響決策、資料狀態或驗證的語言。

## Ai-skill Framework Glossary 邊界

本文件談的是 **domain glossary**（產品 / 業務語彙），與 Ai-skill 自身的 framework / runtime / cognitive vocabulary 是兩個邊界：

| 類型 | Canonical location | 範例 |
| --- | --- | --- |
| Ai-skill framework / runtime / cognitive / architecture | [`knowledge/glossary/ai-skill.md`](../../../../../knowledge/glossary/README.md) | `context_mode`、`cognitive_cost`、`generated_surface`、`owner_layer` |
| Project / business domain | Project 自有 docs（or `<PROJECT_ROOT>/memory/project/context-language.md` replay aid） | 業務模型、產品功能、領域概念 |

Project domain glossary 不得改寫 Ai-skill framework term；Ai-skill framework glossary 不收業務詞。兩者皆遵守同一原則：每個詞單一 owner、不在多處 inline redefine。
