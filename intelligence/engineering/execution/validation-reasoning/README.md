# Validation Reasoning

判斷測試、fixture、manual evidence、contract check 是否足以支持 correctness claim。

## 條目

| 文件 | 用途 |
| --- | --- |
| [`state-visibility-gap.md`](state-visibility-gap.md) | Failure pattern：觀察到的狀態與真實系統狀態不一致。 |
| [`evidence-model.md`](evidence-model.md) | Evidence model：證據類型的 confidence、scope 與 domain preference。 |
| [`evidence-depth.md`](evidence-depth.md) | Decision rule：依風險選擇最低 evidence depth。 |
| [`evidence-chain-validation.md`](evidence-chain-validation.md) | Reasoning pattern：沿狀態傳播鏈驗證到最終可觀察結果。 |
| [`evidence-collapse-point.md`](evidence-collapse-point.md) | Diagnostic reasoning：找出 evidence chain 第一個失去一致性的節點。 |
| [`identity-coupled-side-effect-validation.md`](identity-coupled-side-effect-validation.md) | 特化案例：身份狀態與 side effect 同時存在時的 evidence chain。 |
