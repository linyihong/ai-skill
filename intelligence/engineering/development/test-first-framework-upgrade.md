# Test-First for Framework Upgrade（框架升級先寫測試）

**Status**: `candidate-intelligence`
**Source**: 通用軟體交付經驗 + Ai-skill 2026-05-22 Cognitive Modes plan 實證

## 原則

**Framework / runtime / governance 升級時，validation scenarios 必須寫在 runtime 實作之前；scenarios 是 acceptance contract 不是事後 verification。**

新 framework 進入 Phase 1 runtime 實作前，先：
1. 列出該升級**期望的可觀察行為**
2. 把這些行為寫成 `validation/scenarios/<domain>/<scenario-id>.yaml`
3. 確認 scenarios 目前**會 fail**（因為實作還沒做）
4. 才開始 runtime 實作
5. 實作完成後，scenarios 通過 = acceptance criteria 達成

## 為什麼這比「先實作再寫測試」穩定

| 維度 | 先測試再實作（test-first） | 先實作再寫測試（test-after） |
|------|------|------|
| Acceptance criteria 清晰度 | 高 — scenarios 是契約 | 低 — 實作後才回頭定義「應該測什麼」 |
| 過度設計風險 | 低 — 只實作 scenarios 覆蓋的 | 高 — 容易加 scenarios 沒要求的功能 |
| Scope creep 防護 | 強 — 新 scope 需先加 scenarios | 弱 — 邊做邊加實作 |
| Regression 安全 | 強 — scenarios 早期偵測 break | 弱 — 實作完才測 |
| 重構自信 | 高 — scenarios 守護 | 低 — 不知道 break 什麼 |
| Doc-runtime drift 防護 | 強 — scenarios 是強制 runtime check | 弱 — 文件可能與實作脫節 |
| 跨 phase 整合 | 順 — scenarios 為各 phase 邊界 | 亂 — phase 後才補測試 |

對 framework / governance upgrade 特別重要：**升級錯誤的爆炸半徑大**，scenarios 在實作前固化契約能避免「實作了不該實作的、忘了實作該實作的」。

## 與既有 TDD heuristic 的差異

| 維度 | [`intelligence/engineering/heuristics/test-driven-heuristic.md`](../heuristics/test-driven-heuristic.md) | 本 atom |
|------|------|------|
| 範圍 | 通用 unit test 設計 | Framework / runtime / governance 升級流程 |
| 關注 | 測試難寫 → 設計有問題的 design feedback | 寫 scenarios 在實作前的 ordering principle |
| 對象 | 函數 / class | 升級 plan 的 Phase ordering |
| 主張 | 「測試難寫代表設計需重構」 | 「框架升級 scenarios 必須先於實作」 |
| 互補 | 設計回饋 | 流程順序 |

兩者**互補不衝突**：通用 TDD 處理「測試難不難寫」，本 atom 處理「測試何時寫」。

## 訊號（何時必須套用）

- 升級涉及 `framework / runtime / governance / workflow / validation / scenario / metadata / compiler / generated artifact` 改動
- 任務含「Phase X 實作」性質且有 acceptance criteria
- 跨層改動（≥ 2 個 layer）
- 高 blast radius 改動（影響多個 active workflows）
- 既有測試覆蓋不足以保護新改動

## 操作流程

```
升級進入 Phase N 實作前
  ↓
1. 識別 Phase N 的期望可觀察行為
   - 哪些檔案會被建立 / 修改？
   - runtime.db 哪些 tables / surfaces 會被觸動？
   - 哪些 agent action 會啟用 / 阻擋？
  ↓
2. 寫對應 validation/scenarios/<domain>/<id>-v1.yaml
   - 每個期望行為 = 1 scenario
   - 含 given / when / then + detection_command + pass_criteria
  ↓
3. 跑 scenario 確認**目前 fail**
   - 若 pass → 已實作過 or scenario 太寬鬆，要修
   - 若 fail → 確認失敗訊息符合預期（fail-by-absence 不是 fail-by-error）
  ↓
4. 開始 Phase N 實作
  ↓
5. 實作期間反覆跑 scenarios，逐步從 fail → pass
  ↓
6. Phase N 完成 = 所有對應 scenarios pass
  ↓
7. Commit message 含「scenarios pre-written: <commit hash>」+「now passing」
```

## 何時可以不適用

- **Doc-only trial（如 Cognitive Modes Phase D）**：純 documentation contract，無 runtime 行為可測；改用 manual application checks（agent final report 套 contract）
- **Bug fix / hotfix**：已有測試覆蓋的修補；補測試但不阻斷修補
- **Typo / wording 修正**：無 runtime 行為變更
- **探索性 spike**：明確 throwaway 性質的 prototype

## 何時必須套用（不可豁免）

- 任何進入 `runtime/runtime.db` 的 schema / table 改動
- 任何新增 / 修改 `governance/` 的 enforcement rule
- 任何修改 `enforcement/` 的 blocking gate / activation rule
- 任何修改 compiler / generated_surfaces 投影邏輯
- 任何 framework 命名變更或世代升級

## 與其他智慧的關係

- [`test-driven-heuristic.md`](../heuristics/test-driven-heuristic.md)：互補（設計 vs 順序）
- [`docs-first-bdd-closure.md`](docs-first-bdd-closure.md)：類比思想（observable behavior 前先更新契約 ↔ 實作前先寫 scenarios）
- [`migration-feature-bundling.md`](../anti-patterns/migration-feature-bundling.md)：相關「先驗證再加新功能」邏輯
- [`plan-first-decision-promotion.md`](../architecture/plan-first-decision-promotion.md)：相關「先驗證再升級到 ADR」邏輯
- [`governance/lifecycle/system-upgrade-governance.md`](../../../governance/lifecycle/system-upgrade-governance.md) §3 規則 8：強制 plan 含 validation scenarios（本 atom 進一步指定**順序**：scenarios 先於實作）

## 驗證

| 檢查 | 通過條件 |
|------|------|
| Plan Phase N 對應 scenarios 存在 | 進 Phase N 實作前 `validation/scenarios/` 已有對應 YAML |
| Scenarios commit 早於實作 commit | git log 顯示 scenarios commit hash < 實作 commit hash |
| Scenarios fail-first 已驗證 | Commit message / plan 註記「scenarios pre-written, initially failing」 |
| 實作完成後全部 scenarios pass | `validation/scenarios/...` 對應 scenario `detection_command` 輸出 empty / pass |
| Phase 完成 commit message | 含「now passing」+ scenarios commit hash 引用 |

## 案例證據

**Ai-skill Cognitive Modes plan**（2026-05-22）：
- `ef305bf` 寫入 6 個 validation scenarios（含 `plan-runtime-execution-path-v1`、`failure-pattern-template-consistency-v1` 等）
- 對應 Phase 1 runtime 實作**仍在 Phase 0 等待啟動**（commit `f499397`）
- 已驗證 scenarios 揭露 11 個既有 failure pattern 結構漂移 + 1 個 plan 缺章節 — 在實作前發現問題，避免實作後才補

## Token Impact

實作前寫 scenarios 成本 ~3-5k tokens（依複雜度），但避免：
- 實作走偏：估 10-30k tokens 的 rework
- Phase 後測試補登：估 5-15k tokens 的回溯整理
- Production-like 失敗：跨 session 修復成本不可估

長期 ROI 高。

---

← [回到 engineering/development/](README.md)
