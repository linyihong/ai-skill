# Consistency Boundary Design

## 目的

根據 business invariant 設計 consistency boundary，不以 aggregate pattern 作為預設。

## 步驟

1. 列出 critical invariants。
2. 判斷哪些必須同步成立。
3. 判斷哪些可 eventual consistency / compensation。
4. 定義 validation target 與 failure behavior。
