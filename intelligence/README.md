# Intelligence

`intelligence/` 負責「沉澱工程智慧與領域知識」。本層不是百科知識（什麼是 Redis、什麼是 CQRS），而是**經過經驗抽象化後的工程智慧**——AI 的「專家腦內模型」。

## 與其他層的差異

| 層 | 偏 | 範例 |
|---|----|------|
| `knowledge/` | **事實** | Redis supports pub/sub |
| `skills/` | **執行流程** | How to debug Redis latency |
| `intelligence/` | **判斷力與經驗法則** | If Redis latency spikes suddenly, check connection lifecycle before scaling |

## 核心內容

- **Heuristics（經驗法則）** — 資深工程師直覺
- **Tradeoffs（取捨）** — 「沒有銀彈」的理解
- **Pattern Recognition（模式辨識）** — 可重複的設計與反設計模式
- **Failure Recognition（災難辨識）** — 抽象化後的失敗模式
- **Decision Intelligence（決策智慧）** — 架構與技術選擇的判斷力
- **Contextual Thinking（情境思考）** — 何時適用、何時不適用的邊界條件

## 結構

```text
intelligence/
  engineering/
    architecture/              # 架構思考模式（非教學）
    domain/                    # DDD / 業務模型智慧
    failure/                   # 工程災難智慧（抽象化失敗模式）
    heuristics/                # 經驗法則（intelligence 核心）
    anti-patterns/             # 常見錯誤設計
    tradeoffs/                 # 技術取捨智慧
    distributed-systems/       # 分散式系統生存經驗
    agent-architecture/        # AI Agent 自身運作智慧
    analytical-reasoning/      # 分析推理智慧（APK analysis、repo analysis、heuristics、signals、failure patterns）
    development/               # 開發指引智慧（app development guidance、controls、review checklists）
  ide/                         # IDE 生態系統知識
  business/                    # 商業決策智慧
  travel/                      # 特定領域智慧（Personal Domain Intelligence）
```

## 目前入口

- [`engineering/analytical-reasoning/`](engineering/analytical-reasoning/README.md)：分析推理智慧（APK analysis、repo analysis、evidence-first routing、heuristics、signals、failure patterns）。
  - [`highest-leverage-analysis-path.md`](engineering/analytical-reasoning/highest-leverage-analysis-path.md)：第一個實際 promoted candidate intelligence atom。
  - [`evidence-first-routing.md`](engineering/analytical-reasoning/evidence-first-routing.md)：證據驅動路線選擇（`validated-intelligence`）。
  - [`live-readiness-gates.md`](engineering/analytical-reasoning/live-readiness-gates.md)：SDK/client readiness gates（`validated-intelligence`）。

## 放什麼

- 工程決策原則、trade-off 與架構 lesson。
- 可跨專案重用的 domain knowledge。
- 失效模式、anti-pattern 與改善策略的抽象結論。
- 從分析證據萃取出的穩定判斷。
- 商業與領域的經驗法則。

## 不放什麼

- 百科式知識或技術介紹；放到 `knowledge/`。
- 觀察與拆解的原始方法；放到 `analysis/`。
- 逐步執行流程、review flow 或 task orchestration；放到 `workflow/`。
- 對話暫存 goal、目前 owner 或 next action；放到 `.agent-goals/`。
- 可執行 policy 與 close-loop gate；放到 `enforcement/`。

## 誰會參考這裡（Inbound References）

- [`route.intelligence.engineering.heuristics`](../knowledge/runtime/routing-registry.yaml:371) — required_dependencies 引用 `intelligence/engineering/README.md`
- [`route.intelligence.engineering.agent-architecture`](../knowledge/runtime/routing-registry.yaml:621) — required_dependencies 引用 `intelligence/engineering/README.md`
- [`route.intelligence.apk-analysis.atoms`](../knowledge/runtime/routing-registry.yaml:548) — primary_source 為 `intelligence/engineering/analytical-reasoning/README.md`
- [`route.intelligence.apk-highest-leverage-path`](../knowledge/runtime/routing-registry.yaml:249) — primary_source 為 `intelligence/engineering/analytical-reasoning/highest-leverage-analysis-path.md`
- [`route.workflow.apk-analysis`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `intelligence/engineering/analytical-reasoning/README.md`
- [`route.workflow.software-delivery`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `intelligence/engineering/development/README.md`
- [`route.workflow.travel-planning`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `intelligence/travel/README.md`

## 與既有層的關係

- `skills/` 已逐步遷移至新分層。舊 `skills/app-development-guidance/` 已刪除。
- `workflow/` 應 reference 本層，而不是內嵌大量知識。
- `feedback/` 可把新 lesson promotion 到本層。
- `governance/` 定義本層知識的 lifecycle、清理與 validation。

## 遷移狀態

- `skills/app-development-guidance/implementation/` — ✅ 已遷移至 `metadata/development-guidance/implementation-catalog.md`，舊目錄已刪除
- `skills/app-development-guidance/controls/` — ✅ 已遷移至 `metadata/development-guidance/controls-catalog.md`，舊目錄已刪除
- `feedback/history/` 中已成熟且跨專案可重用的 lesson — 持續進行中
- `enforcement/failure-patterns/` 中偏工程判斷的 pattern 摘要 — 持續進行中
