# Product Impact Discovery

## 目的

在 behavior-driven discovery 之前，先確認 product brief、business goal、target actor 與 customer journey 是否對齊。這一步防止 AI 快速產出功能，但沒有對準真正 user pain 或 business impact。

## 步驟

1. 建立 Impact Map：Why / Who / How / What。
2. 建立或讀取 Customer Journey Map：stage、trigger、pain point、emotional low、blocker、evidence。
3. 交叉驗證 Impact Map 與 Journey Map：Who、timing、pain、investment 是否對齊。
4. 將每個主要 claim 標記為 `validated`、`assumption`、`open question`、`scoped out` 或 `invalidated`。
5. 輸出 decision：`proceed`、`revise`、`reject` 或 `ask_user`。

## 觸發時機

- 新產品、新功能或高成本 feature investment。
- Product brief 有 goal / feature list，但缺少 user journey evidence。
- AI 產出很多功能，但不知道哪個真的有 impact。
- BDD scenario 看似完整，但不確定是否對準真實使用者痛點。

## 參考

- `intelligence/engineering/requirements/product-alignment/`
- `workflow/software-delivery/templates/product-impact-alignment-template.md`
