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
