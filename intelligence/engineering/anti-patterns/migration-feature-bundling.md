# Migration + Feature Bundling（搬遷與新功能混綁反模式）

**Status**: `candidate-intelligence`
**Source**: 通用軟體交付經驗

## 反模式

把 migration（大型 refactor — rewrite / platform 升級 / monolith→microservices / cloud 遷移 / DB engine 替換）與 new feature（behavior change）綁進同一階段交付。新版同時做「搬遷舊邏輯」+「新增舊版沒有的功能」。

這是 Martin Fowler「Refactoring 不可同時改 behavior」法則的放大版。

## 訊號

- Migration plan 同時列「移植 module X」+「新增 feature Y」作為同一 phase 交付
- 客戶 / PM / 業務在改版 kickoff 提出「既然要重寫，順便加 OOO 吧」
- Plan 沒有獨立的 **parity gate**（新版功能完全等於舊版）作為 milestone
- 新版上線後出現行為差異，團隊無法判斷該不該是 bug
- Code review 出現「不確定這是搬遷預期還是新功能預期」這類對話
- 既有 BDD / 回歸測試無法直接套用於新版（因為 spec 已被新功能改寫）

## 根本原因

- 把 migration 誤當成「重寫機會」，順手把累積的功能需求一起塞進去
- 沒有意識到 migration 需要的不是「最佳新版本」，而是「外部可觀察行為等價的新版本」
- Stakeholder 壓力（「不加功能使用者會失望」）優先於工程可驗證性

## 影響：Verification Identity Crisis

當新功能與 migration 混雜，**驗證失去 ground truth**：

| 出現差異時 | Parity-first 可判斷 | Bundled 無法判斷 |
|------|------|------|
| 行為與舊版不同 | 是搬遷 bug，需修正 | 來源 = 搬遷？新功能？互動？ |
| 效能退化 | 搬遷引入，需 profiling | 多重來源疊加，難 profiling |
| 回歸測試失敗 | 對應到具體搬遷段落 | 不知道該檢查哪段代碼 |
| 客戶回報「以前可以這樣用」 | 搬遷遺漏，補回 | 可能是「故意不支援」也可能是 bug |

### 其他連帶後果

| 維度 | 後果 |
|------|------|
| 工程驗證成本 | 從線性 → 指數（每個 bug 都要先排除「是不是新功能造成的」） |
| 時程預測性 | 兩個 unknown（搬遷量、新功能量）疊加，估算誤差成倍 |
| 回滾可行性 | 從「回滾搬遷」變成「沒法回滾，因為新功能不在舊版」 |
| 團隊認知負擔 | 每個 PR 都要同時帶兩種 context（搬遷 + 新功能） |

## 正確路徑：Parity-First Migration

```
Phase 1: Parity Migration
   ↓ Gate（必須通過才能進 Phase 2）：
   - 舊版 BDD / regression suite 在新版 100% 通過
   - 已知例外明確記錄為「故意接受的差異」
   - 客戶可在新版執行所有舊版業務流程
Phase 2: 新功能（基於穩定新版）
```

### Trade-off 軸向

| 軸 | Parity-first | Bundled-with-features |
|------|------|------|
| 客戶短期滿意度 | 低（看不到新東西） | 高 |
| 工程驗證成本 | 線性、可定位 | 指數、不可定位 |
| 上線風險 | 可控（差異 = 搬遷 bug） | 不可控（差異來源不明） |
| 時程預測性 | 高 | 低 |
| 回滾可行性 | 高 | 低 |
| Stakeholder 溝通難度 | 高（需要說服） | 低（先 buy-in、後爆炸） |

## Stakeholder 翻譯：「失望總比絕望好」

當客戶 / PM 以「不加功能使用者會失望」反對 parity-first 時，需要把抽象的工程風險翻成具象情緒對比：

> 工程語：「驗證失去 ground truth，bug 來源無法定位，時程不可預測」
> ↓
> Stakeholder 語：**「失望總比絕望好」**

為什麼這話術有效：

1. 把抽象（無法驗證）轉成具象（情緒對比）
2. Stakeholder 不需要懂 BDD 也能理解
3. 留下「絕望」這個錨點，讓對方記住後果嚴重性
4. 短期失望 vs 長期絕望，時間框架對比清楚

詳見 [`analysis/development-guidance/risk-translation.md`](../../../analysis/development-guidance/risk-translation.md) §5。

## 何時不適用

- 純新建（greenfield）專案，沒有舊版可作 ground truth
- 完全砍掉舊功能重新設計使用者體驗（明確聲明 break compatibility，且已與客戶達成新版 spec 共識）
- 小規模 patch / hotfix，無搬遷性質
- 舊版已無人使用，新版實質上是「新產品」

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「既然要動，順便加 OOO 吧」 | Migration 的目標是「外部行為等價」，不是「最佳化版本」 |
| 「parity 太無聊，團隊沒動力」 | Parity gate 通過後 Phase 2 立刻開始；新功能不是被砍掉，是被延後到可驗證階段 |
| 「我們是 rewrite，不是 migration，所以可以隨便改」 | Rewrite 仍要對外部使用者保持行為等價，否則就是新產品而非 rewrite |
| 「測試覆蓋率不夠，反正怎麼樣都驗不全」 | 那就更不該再加新功能讓問題擴大；先補回歸測試或 BDD |

## 驗證

| 檢查 | 通過條件 |
|------|------|
| Migration plan 含 parity gate | Plan 中有明確 milestone「新版 = 舊版功能等價」 |
| 舊測試套件可重用 | 舊版 BDD / regression 在新版 100% 通過或明列例外 |
| Bug 可溯源 | 出現差異時 24 小時內能定位「搬遷 bug」「環境差異」或「已知接受差異」 |
| Scope-lock | 新功能延後到 Phase 2，且 Phase 1 完成前不開工 |

## 與其他智慧的關係

- [`migration-seeder-anti-patterns.md`](migration-seeder-anti-patterns.md)：兩者皆是 migration 反模式；前者談 schema migration 不該載業務資料，本檔談 system migration 不該綁新功能
- [`architecture-absolutism.md`](architecture-absolutism.md)：把單一架構當 universal default；本檔是「把單一 phase 當『修一切的機會』」的時序版
- [`analysis/development-guidance/risk-translation.md`](../../../analysis/development-guidance/risk-translation.md)：本檔 §「Stakeholder 翻譯」是該文件的具體案例
- Martin Fowler "Refactoring"：refactor + behavior change 不可同時的法則放大版

---

← [回到 engineering/anti-patterns/](README.md)
