> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-22 — Test-First for Framework Upgrade

Status: candidate

#### One-line Summary

Framework / runtime / governance 升級時 validation scenarios 必須寫在 runtime 實作之前；scenarios 是 acceptance contract 不是事後 verification。

#### Human Explanation

Cognitive Modes plan 在 2026-05-22 session 內把 6 個 validation scenarios（commit `ef305bf`）寫在 Phase 1 runtime 實作之前。使用者觀察到這種順序「先實作測試 再執行」比「先實作再補測試」更能增加框架穩定性。

實證：scenarios 在 Phase 1 實作前先跑，揭露 11 個既有 failure pattern 結構漂移（`failure-pattern-template-consistency-v1`）+ 1 個 plan 缺 section（`plan-runtime-execution-path-v1`），在實作前就先發現並修補（A+B commit `3a38e49`、plan tracking fix `2ca5b4f`）。若 scenarios 寫在實作後，這些 gaps 可能要等 Phase 3+ 才會被偵測到，修補成本更高。

#### Trigger

- 升級涉及 framework / runtime / governance / workflow / validation / scenario / metadata / compiler 改動
- 任務含「Phase X 實作」性質且有 acceptance criteria
- 跨層改動（≥ 2 個 layer）
- 高 blast radius 改動（影響多個 active workflows）
- 既有測試覆蓋不足以保護新改動

#### Evidence

- Tool: Ai-skill repo session（2026-05-21~22）
- Sanitized excerpt: commit chain
  - `ef305bf` 寫 6 scenarios（Cognitive Modes Phase 1 實作前）
  - `3a38e49` A+B 補完，scenario 揭露 11 個 patterns 結構漂移在實作前修補
  - `2ca5b4f` plan completion tracking 修補（scenario 揭露的 gap）
  - `f499397` Phase 0 Preflight 完成，未進 Phase 1 實作
- Evidence path: 本 Ai-skill repo git history

#### Generalized Lesson

升級進入 Phase N 實作前流程：

```
1. 識別 Phase N 期望可觀察行為（檔案、runtime.db、agent action）
2. 寫 validation/scenarios/<domain>/<id>-v1.yaml（每個行為 1 scenario）
3. 驗證 scenarios 目前 fail（fail-by-absence 不是 fail-by-error）
4. 開始 Phase N 實作
5. 反覆跑 scenarios，逐步 fail → pass
6. Phase N 完成 = 所有對應 scenarios pass
7. Commit message 含 scenarios commit hash 引用 + 「now passing」聲明
```

豁免：doc-only trial / bug fix / typo / 探索性 spike。
不可豁免：runtime.db schema、enforcement rule、blocking gate、compiler / generated_surfaces。

#### Agent Action

執行 framework / runtime / governance 升級時：

1. 進 Phase N 實作前先評估「是否該套 test-first」
2. 若是 → 寫 scenarios → 驗證 fail → 才開始實作
3. Commit message 區分 scenarios commit 與實作 commit
4. Phase 完成時 commit message 含「scenarios pre-written: <hash>, now passing」
5. 若無對應 scenarios 即進實作，先停下補上

#### Goal / Action / Validation

- Goal: framework 升級的 acceptance criteria 固化為機器可測 contract，先於實作
- Action: scenarios 寫入 → 驗 fail → 實作 → 驗 pass → commit message 含對照
- Validation or reference source: git log 順序（scenarios commit hash < 實作 commit hash）；commit message 含 fail-first 註記；scenarios `detection_command` 輸出 empty/pass

#### Applies When

- Framework / runtime / governance / compiler / metadata / validation 升級
- Plan 含 Phase X 實作且有 acceptance criteria
- 跨層改動
- 高 blast radius

#### Does Not Apply When

- Doc-only trial
- Bug fix / hotfix（已有測試覆蓋）
- Typo / wording 修正
- 探索性 spike（throwaway prototype）

#### Validation

- 對既有 plan 跑 commit 順序檢查能驗證 test-first 是否遵守
- Phase N 結束時 scenarios 對應 `detection_command` 全部 pass

#### Promotion Target

- ✅ `intelligence/engineering/development/test-first-framework-upgrade.md`（本次新增）
- ✅ `knowledge/summaries/test-first-framework-upgrade.md`（summary card）
- ✅ `validation/scenarios/failure-derived/test-first-for-framework-upgrades-v1.yaml`（強制 scenario）
- ✅ `governance/lifecycle/system-upgrade-governance.md` §3 規則 9 — 強制順序

#### Required Linked Updates

- `intelligence/engineering/development/README.md`（加入索引）
- `knowledge/summaries/README.md`（加入 summary）
- `plans/README.md` plan template 加 test-first ordering 說明
- Step 6（Intelligence Extraction）：done(executed) — atom 直接 promote
- Step 7（Failure Learning）：not_applicable — 為正向原則沉澱，非 failure pattern
