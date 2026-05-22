> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-22 — Migration + Feature Bundling 反模式與 stakeholder 風險翻譯

Status: candidate

#### One-line Summary

大型系統改版（rewrite / migration / platform 升級）時，**新版必須先達成 parity 才能加新功能**；同時混入新功能會讓驗證失去 ground truth，bug 來源無法定位；對 stakeholder 的有效翻譯模式是「失望總比絕望好」。

#### Human Explanation

這是 Martin Fowler「Refactoring 不可同時改 behavior」法則的放大版。Migration 本質上是一場大規模 refactor — 從舊版搬到新版時，必須保持外部可觀察行為等價，才能用既有的 BDD / 回歸測試當 ground truth 驗證搬遷成功。一旦在新版同時加入舊版沒有的功能，差異來源就變成多元：「這個 bug 是搬遷錯了？還是新功能引入？還是兩者互動？」團隊會花大量時間在排除假設，回歸成本指數增長，時程黑洞接連發生。

工程上的論證難以說服 PM / 業務 / 客戶，因為他們聽到的是「不能加功能」「使用者會失望」。需要一句具象、有情緒對比的翻譯：**「失望總比絕望好」** — 把抽象的驗證風險翻成 stakeholder 聽得懂的後果對比。失望是短期感受，絕望是長期系統崩塌。

#### Trigger

- Migration / rewrite / platform 升級的 plan 同時列「移植 module X」+「新增 feature Y」
- 客戶 / PM 在改版 kickoff 提出「既然要重寫，順便加 OOO 吧」
- Migration 沒有獨立的 **parity gate**（新版功能完全等於舊版）作為 milestone
- 新版上線後出現行為差異，團隊無法判斷該不該是 bug，因為 spec 同時被 migration + new feature 改寫，沒有 ground truth 可對照

#### Evidence

- Tool: 真實 stakeholder 對話（顧問與客戶就改版範圍討論）
- Sanitized excerpt:「大型系統改版切忌同時在新版調整現有功能！改了功能你如何驗證原本商業邏輯的正確性？」「如果不加一些新功能的話，使用者會很失望耶！」「失望總比絕望好」
- Evidence path: stakeholder 對話本身為證據，已抽象化，不涉及任何專案、品牌或機密

#### Generalized Lesson

**反模式**：把 migration（refactor）與 new feature（behavior change）綁進同一階段交付。

**正確路徑**：
```
Phase 1: Parity Migration
   ↓ Gate: 舊版 BDD / regression 必須在新版 100% 通過
Phase 2: 新功能（基於穩定新版）
```

**Trade-off 軸向**：

| 軸 | Parity-first | Bundled-with-features |
|------|------|------|
| 客戶短期滿意度 | 低 | 高 |
| 工程驗證成本 | 線性、可定位 | 指數、不可定位 |
| 上線風險 | 可控（差異 = 搬遷 bug） | 不可控（差異來源不明） |
| 時程預測性 | 高 | 低（兩個 unknown 疊加） |
| 回滾可行性 | 高 | 低 |

**Stakeholder 翻譯模式**：當客戶以「使用者失望」反對 parity-first 時，回應「失望總比絕望好」。這是把工程語「驗證失去 ground truth」翻成具象情緒對比，留下「絕望」的錨點讓對方記住後果嚴重性。

#### Agent Action

在 migration / rewrite / platform 升級的分析或 plan review 中：

1. 檢查 plan 是否含 parity gate（新版 = 舊版功能等價）作為獨立 milestone
2. 若 plan 同時列「移植」+「新增功能」，標記為 anti-pattern candidate
3. 提出 Phase 1 (parity) → Phase 2 (new feature) 的拆分建議
4. 若團隊 / 客戶反對，準備「失望總比絕望好」這類具象風險翻譯
5. 若新版上線後出現「不確定該不該是 bug」的差異，溯源檢查是否本反模式造成

#### Goal / Action / Validation

- Goal: Migration 的驗證有 ground truth；bug 來源可定位
- Action: 建立 parity gate；scope-lock 新功能到 Phase 2；對 stakeholder 提出具象風險翻譯
- Validation or reference source: 舊版 BDD / regression suite 在新版 100% 通過；新版任一行為差異能明確指向「搬遷 bug」「環境差異」或「明確記錄的接受差異」；無「不確定該不該是 bug」的情況

#### Applies When

- 大型系統 rewrite（換語言、換框架、換 platform）
- Service migration（monolith → microservices、on-prem → cloud）
- DB engine 替換（含 schema 演進但業務邏輯保持）
- Legacy modernization、tech debt 清算
- 任何「外部可觀察行為應保持等價」的搬遷專案

#### Does Not Apply When

- 純新建（greenfield）專案，無舊版可比對
- 完全砍掉舊功能重新設計使用者體驗（明確聲明 break compatibility）
- 小規模 patch / hotfix，無「搬遷」性質

#### Validation

- 可在實際 migration plan 中找到「parity gate」這個 milestone
- Migration phase 結束時，舊版測試套件在新版通過率 = 100%（或明列已知例外）
- 出現行為差異時，能在 24 小時內定位來源是「搬遷錯」「環境差」還是「已知接受差異」

#### Promotion Target

- `intelligence/engineering/anti-patterns/migration-feature-bundling.md`（本次新增）
- `analysis/development-guidance/risk-translation.md`（加入 §5 stakeholder 翻譯案例）
- `knowledge/summaries/migration-feature-bundling.md`（summary card）

#### Required Linked Updates

- `intelligence/engineering/anti-patterns/README.md`（加入索引條目）
- `knowledge/summaries/README.md`（加入 summary 條目）
- Step 6（Intelligence Extraction）：done(executed) — 本 lesson 同時 promote 為 intelligence anti-pattern atom
- Step 7（Failure Learning）：not_applicable — 本 lesson 來自 user-supplied 工程原則與 stakeholder 對話，非 agent failure 補救
