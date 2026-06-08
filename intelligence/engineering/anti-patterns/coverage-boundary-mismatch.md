# Coverage Boundary Mismatch

**Status**: `candidate-intelligence`

## 定義

Coverage boundary mismatch 是 validator 推論的 coverage 範圍比 source generator 實際產生的 inventory 更大。

典型失敗形狀：

```text
source generator succeeds
runtime validation succeeds
freshness / completeness validator fails
```

失敗不是因為 source stale，而是 validator 用了自己的 coverage 推論，沒有使用同一份 canonical source inventory。

## 常見症狀

- Generator 只產生明確 source rows，但 validator 依目錄、檔名或 sibling artifact 推論更多檔案應被覆蓋。
- Refresh 成功且 generated artifact 已更新，但 commit-time validator 仍要求 generator 從未承諾的 row。
- 使用 opt-out 才能提交，原因不是資料沒刷新，而是 validator 和 builder 的 coverage contract 不一致。

## 判斷規則

如果 validator 的 blocking 條件依賴「應該被覆蓋」的推論，而不是 generator 輸出的 source inventory，就有 coverage boundary mismatch 風險。

Source-of-truth 應該是：

```text
canonical source inventory
→ freshness / checksum / completeness validation
```

不是：

```text
directory or naming heuristic
→ inferred inventory
→ blocking validation
```

## 修正方式

1. 讓 validator 讀 generator 產出的 canonical source inventory。
2. 對 inventory 內的 source 做 freshness / checksum / completeness 檢查。
3. 對 inventory 外的 artifact 不做 blocking freshness claim，除非先把 generator coverage contract 擴大。
4. 若需要擴大 coverage，先改 generator，再改 validator；不要讓 validator 先行推論。

## 相關知識

- [`validation-proxy-trap.md`](validation-proxy-trap.md)
- [`mock-completeness-illusion.md`](mock-completeness-illusion.md)
