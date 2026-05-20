# Bounded Context Discovery

## 目的

從 shared language、behavior scope、domain invariants 和 ownership 找出 bounded context。

## 步驟

1. 檢查同一詞是否在不同 workflow 中含義不同。
2. 檢查 invariant / lifecycle 是否不同。
3. 檢查 external model 是否污染 internal language。
4. 決定 module boundary、bounded context 或 anti-corruption layer。
