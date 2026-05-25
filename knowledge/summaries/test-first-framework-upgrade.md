## intelligence.test-first-framework-upgrade

| 欄位 | 值 |
| --- | --- |
| Atom ID | `intelligence.test-first-framework-upgrade` |
| Source path | `intelligence/engineering/development/test-first-framework-upgrade.md` |
| Lifecycle | `candidate` |
| Summary | Framework / runtime / governance 升級時 validation scenarios 必須寫在 runtime 實作之前，scenarios 是 acceptance contract 不是事後 verification。流程：列期望可觀察行為 → 寫 YAML scenarios → 驗證目前 fail → 開始實作 → 完成 = scenarios pass。與通用 TDD heuristic 互補（設計回饋 vs 順序原則）。Doc-only trial / typo / spike 可豁免；runtime.db schema / enforcement rule / compiler 改動不可豁免。 |
| When to read | 升級涉及 framework/runtime/governance/workflow/validation/scenario/metadata/compiler 改動；任務含 Phase X 實作 + acceptance criteria；跨層改動；高 blast radius；既有測試覆蓋不足。 |
| Do not use for | Doc-only trial（無 runtime 行為可測）；bug fix / hotfix（已有測試覆蓋）；typo / wording 修正；探索性 spike（throwaway prototype）。不取代通用 TDD design feedback（見 heuristics/test-driven-heuristic.md）。 |
| Context cost | ~310 tokens |
| Estimated full cost | ~2400 tokens |
| Validation signal | Plan Phase N 對應 scenarios 在 validation/scenarios/ 已存在；git log 顯示 scenarios commit hash < 實作 commit hash；commit message 含 fail-first 註記與「now passing」聲明；scenarios detection_command 輸出 empty/pass。 |
| Last checked | 2026-05-22 |
