# Model Governance

`models/governance/` 定義 model-aware execution 的安全邊界。核心原則：可以調整 behavior，不可未經證實宣稱實際 model selection。

## 入口

- [`model-selection-governance.md`](model-selection-governance.md)：provider-neutral selection boundary。
- [`hallucination-boundaries.md`](hallucination-boundaries.md)：low-confidence output 的 claim scope 與 source gate。
- [`context-budget-governance.md`](context-budget-governance.md)：長任務與 compression budget。
- [`model-confidence-governance.md`](model-confidence-governance.md)：capability confidence 與 downgrade rules。

## Layer Boundary

Tool-specific model names、selector UI、availability 與 model list 屬於 `ai-tools/`。`models/` 只保存 reusable behavior contract、capability dimensions 與 workflow adaptation。
