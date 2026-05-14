# Documentation Backfill Heuristic（文件回填經驗法則）

**Status**: `candidate-intelligence`
**Source**: [`analysis/repo/documentation-backfill.md`](../../analysis/repo/documentation-backfill.md), [`skills/app-development-guidance/process/README.md`](../../skills/app-development-guidance/process/README.md)

## 原則

**If a repository has working code but missing documentation, recover behavior first, then contracts, then intent.**

如果一個 repository 有可運作的程式碼但缺少文件，先恢復行為，再恢復契約，最後恢復意圖。

## 為什麼

1. **已實作的行為是最強的可用真相來源** — 程式碼可能不是「正確的」，但它顯示了「實際發生的」。
2. **BDD 是恢復的錨點** — 從 observable behavior 回填 BDD 場景，再從 BDD 推斷 domain model、architecture、API contracts。
3. **Product intent 是最難恢復的** — 如果原始商業理由不可得，標記為 `unknown` 而不是發明理由。發明的理由會在未來導致錯誤的決策。
4. **Pipeline artifacts 也需要恢復** — 不只是開發文件，CI/CD、deployment、testing 管線的設計意圖也需要記錄。

## 何時適用

- 第一次接觸一個已實作完成的 repository。
- 需要為 legacy codebase 補上開發文件。
- 需要從現有程式碼反向恢復 domain intent 與架構決策。
- 團隊交接或知識轉移時。

## 何時不適用

- 專案已有完整的開發文件（只需要驗證同步狀態）。
- 專案即將被棄用且不需要長期維護。
- 只修改現有功能，不需要全面恢復文件。

## 決策流程

```text
已實作完成的 repository，缺少文件？
  ├── Step 1: Inventory 現有 docs、tests、schemas、fixtures
  ├── Step 2: 建立 documentation gap table（exists / partial / missing / unknown）
  ├── Step 3: 先恢復 BDD Behavior（從 UI、API、tests、logs）
  │     ├── Critical happy paths
  │     ├── Failure paths
  │     ├── Permissions / authorization
  │     ├── Empty states
  │     ├── Edge cases
  │     └── Cross-context flows
  ├── Step 4: 從 BDD + implementation evidence 恢復 contracts
  │     ├── Domain Model（entities、value objects、commands、events）
  │     ├── Architecture Contract（dependency direction、data ownership）
  │     ├── API / Interface Contract（request/response schemas）
  │     └── Error Handling Contract（error taxonomy、retry rules）
  ├── Step 5: 標記 unknown product intent，不發明商業理由
  ├── Step 6: 如果 BDD 無法從可用證據完成，停止並要求更多資訊
  └── Step 7: 對缺乏 coverage 的 critical BDD scenario，新增 test TODOs
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「這個 function 叫 calculatePrice，所以它的目的是計算價格」 | 記錄 observed behavior，不要發明 intent |
| 「先補 product brief，再補技術文件」 | 先補 BDD（從行為恢復），product brief 的 intent 可能無法恢復 |
| 「文件恢復是一次性工作」 | 文件恢復是 iterative 的，每次修改都會增加更多理解 |
| 「所有 missing 文件都要補到 100%」 | 優先恢復 critical paths，非關鍵路徑可以標記為 partial |

## Token Impact

避免因缺少文件導致的重複 reverse engineering。一個沒有 BDD 的 repository 每次修改都需要 15-30 分鐘理解現有行為，而有 BDD 的 repository 只需 2-5 分鐘確認場景。

---

← [回到 intelligence/engineering/analysis/](README.md)
