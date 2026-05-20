# Lifecycle Evaluation

**Status**: `candidate-intelligence`

## 判斷原則

架構投資必須對齊系統壽命。短期探索需要 learning speed；長期核心系統需要模型穩定性與演化能力。

## 分級

| Lifecycle | 優先事項 | 架構偏好 |
| --- | --- | --- |
| Prototype | 快速驗證需求 | simple service layer、vertical slice |
| MVP | 快速交付且保留改造空間 | feature module、DDD Lite |
| Growth | 穩定領域語言與邊界 | modular monolith、selective DDD |
| Long-lived core | 保護不變量與 context boundary | Full DDD、ACL、event coordination |

## 升級條件

只有當維護成本、語言衝突或 invariant break 成為 recurring evidence 時，才從簡單架構升級。
