# Docs-First BDD Closure Heuristic（文件優先 BDD 閉環經驗法則）

**Status**: `candidate-intelligence`
**Source**: [`workflow/software-delivery/execution-flow.md`](../../../workflow/software-delivery/execution-flow.md), `skills/app-development-guidance/WORKFLOW.md`（已刪除）

## 原則

**If observable behavior changes, the owning contract must be updated before code.**

可觀察行為變更前，必須先更新管轄該行為的契約文件。

## 為什麼

1. **程式碼是最不可靠的真相來源** — 程式碼只顯示「當前的實作」，不顯示「為什麼這樣做」或「原本預期做什麼」。
2. **文件與程式碼不同步時，下一次修改會基於錯誤的前提** — 這是技術債累積最快的方式。
3. **BDD 場景是人類可讀的行為規格** — 它們是 product intent 與 executable test 之間的橋樑，缺少 BDD 等同於缺少可驗證的 acceptance criteria。
4. **先更新契約再寫程式碼，強制開發者思考「我要改變什麼行為」而非「我要改哪一行」**。

## 何時適用

- 任何會改變 observable behavior 的變更（新功能、行為修改、bug fix 影響 public contract）。
- 在已實作完成的 repository 中恢復開發文件時。
- 當 frontend 與 backend 由不同團隊或 agent 維護時。
- 當 repository 已有 BDD 或 contract testing 基礎設施時。

## 何時不適用

- 純 refactor（不改變 observable behavior 或 public contract）。
- 僅修改內部實作細節（private method rename、extract helper）。
- 專案處於 prototyping 階段且明確同意不維護文件。
- 變更僅影響 tooling 或 build 流程，不影響 runtime behavior。

## 決策流程

```text
有 observable behavior 變更？
  ├── 否 → 純 refactor，不需更新 BDD/contracts
  └── 是 → 哪份文件管轄這個行為？
        ├── Governance / framework contract
        ├── Product plan / brief
        ├── BDD behavior spec
        ├── Domain / architecture / API contract
        ├── Implementation
        └── Tests / fixtures
        按優先順序更新管轄文件，再實作程式碼
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「先寫 code，再補 doc」 | 先更新契約再寫 code；補 doc 通常不會發生 |
| 「BDD 太花時間，直接用 unit test 代替」 | BDD 描述行為意圖，unit test 驗證實作細節；兩者不同層次 |
| 「契約衝突時，以程式碼為準」 | 程式碼可能是 bug；應分類衝突類型並更新管轄文件 |
| 「只改 frontend，backend 契約不用動」 | 如果 frontend 行為依賴 backend contract，兩邊都要同步 |

## Token Impact

避免「先寫 code 再補 doc」的習慣，後續補 doc 的成本通常是即時更新的 3-5 倍，且經常被跳過。

---

← [回到 engineering/app-development-guidance/](README.md)
