# Plans（計畫目錄）

## 目錄規則

| 子目錄 | 用途 | 生命週期 |
|--------|------|---------|
| [`active/`](active/) | 進行中或待審閱的計畫（draft / in-progress） | 完成後搬移至 `archived/` |
| [`archived/`](archived/) | 已執行完成的計畫（執行結果記錄） | 永久保留，作為決策記錄 |

## 原則

1. **`active/` 只放尚未開始或正在執行的計畫** — 一旦計畫執行完畢，立即搬移至 `archived/`
2. **`archived/` 的計畫不刪除** — 作為歷史決策記錄，可供日後查閱
3. **計畫檔案命名規則**：`<slug>.md`，slug 需能反映計畫核心目標
4. **每個計畫必須在檔頭標註狀態**：`draft` / `in-progress` / `completed`
5. **計畫完成後，若從中提煉出可重用的系統經驗，應建立對應的 intelligence atom**

## 目前狀態

| 檔案 | 狀態 | 說明 |
|------|------|------|
| [`active/cognitive-boundary-system.md`](active/cognitive-boundary-system.md) | ⏳ draft | Cognitive Boundary System 整合計畫，待審閱後開始實作 |
| [`archived/technique-intelligence-pilot.md`](archived/technique-intelligence-pilot.md) | ✅ completed | Phase 28：Technique → Intelligence Pilot（flutter-dart-aot） |
| [`archived/skill-specific-extraction.md`](archived/skill-specific-extraction.md) | ✅ completed | Phase 33：Skill-Specific Intelligence Extraction |
| [`archived/ai-decision-contract-testing.md`](archived/ai-decision-contract-testing.md) | ✅ completed | AI Decision Contract Testing 框架設計與實作 |

## 與其他層的關係

- [`architecture/next-stage-upgrade-plan.md`](../architecture/next-stage-upgrade-plan.md) — 全局升級路線圖，`plans/` 中的計畫是路線圖的具體執行計畫
- [`governance/lifecycle/README.md`](../governance/lifecycle/README.md) — Skills Deprecation Timeline 等生命週期規則
- [`intelligence/engineering/agent-architecture/`](../intelligence/engineering/agent-architecture/) — 從已完成計畫中提煉的系統經驗結晶
