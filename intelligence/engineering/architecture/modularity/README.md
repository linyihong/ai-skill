# Modularity

Modularity intelligence 保存 module boundary、feature slice、package boundary 與 modular monolith 的判斷。

## 目前入口

| 文件 | 用途 |
| --- | --- |
| [`module-count-discipline.md`](module-count-discipline.md) | 模組數量紀律：N 對應健康做法、新增 module 的五道判斷、build system 中立原則（Maven / Gradle / workspace / Cargo / Bazel / 等）。處理「畫邊界後的數量管理」，與本目錄「畫邊界本身」互補。 |

## 與其他智慧的關係

- 畫邊界後的「跨邊界耦合」→ [`../coupling-tradeoffs/`](../coupling-tradeoffs/README.md)
- 模組規模到一定程度且有獨立 scale / deploy 需求 → [`../modular-monolith-vs-microservices.md`](../modular-monolith-vs-microservices.md)
- 多廠商整合的 N 模組特化場景 → [`../vendor-integration-architecture.md`](../vendor-integration-architecture.md)
