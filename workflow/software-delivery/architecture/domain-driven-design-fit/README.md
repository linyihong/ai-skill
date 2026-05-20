# Domain-Driven Design Fit

## 目的

判斷需求與 domain complexity 是否需要 DDD Lite 或 Full DDD，而不是預設使用 DDD。

## 步驟

1. 讀 requirements stage 的 behavior boundary。
2. 評估 domain complexity、invariant density、language instability、integration pressure、lifecycle。
3. 對照 `metadata/architecture/architecture-fit-matrix.yaml`。
4. 輸出 CRUD / DDD Lite / Full DDD recommendation。
