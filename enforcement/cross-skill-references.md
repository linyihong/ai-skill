# Cross-workflow references（跨 workflow 引用）

當一個 workflow 需要另一個 workflow 的 output contract、documentation format、validation checklist 或 implementation guidance 時，使用 cross-workflow reference。

## 原則

Workflow 可以引用其他 workflow，但不能複製另一個 workflow 的完整規則。Referring workflow 應說明何時讀取 target workflow、期待什麼輸出，以及目前 workflow 保留哪些 ownership boundary。

## 何時引用另一個 workflow

符合下列情況時，加入明確 cross-workflow reference：

- 某 workflow 產生的 artifact 需要由另一個 skill 消費。
- 某 workflow 需要另一個 workflow 的 template、checklist 或 contract format。
- 某 finding 應轉成另一個領域的 guidance。
- 使用者要求從一個 workflow handoff 到另一個 workflow。

例：`apk-analysis` 負責授權 APK evidence、traffic attribution、schema recovery 與 sanitized feature reconstruction handoff。`app-development-guidance` 負責把 handoff 轉成 BDD、Domain Model Contract、API / Interface Contract、Error Handling Contract、implementation slices、checks 與 tests。

## 必要引用格式

新增 cross-workflow reference 時，包含：

| 欄位 | 要求 |
| --- | --- |
| Target workflow | 連到 `workflow/<domain>/` 或要讀取的特定檔案/資料夾。 |
| Trigger | 明確寫出 agent 何時應讀 target workflow。 |
| Expected input/output | 命名 workflow 之間傳遞的 artifact，例如 handoff、checklist、contract、fixture 或 implementation guidance。 |
| Ownership boundary | 說明哪些仍屬於目前 workflow，哪些屬於 target workflow。 |
| Sanitization boundary | 說明 target-specific data、secrets、raw evidence 或 product conclusions 是否必須留在 project docs。 |
| Linked updates | 更新雙方 skill 入口；若 target skill 已涵蓋 handoff，說明無需更新的理由。 |

## 不要做

- 不要把另一個 workflow 的完整流程貼進目前 workflow。
- 不要讓每個 workflow 預設讀取所有其他 workflow。
- 不要建立循環 `always read` 鏈。
- 不要因為兩個 skills 互相引用，就把 target-specific evidence 搬進 reusable skill。
- 如果 requested output 必須 cross-workflow handoff，不要把它描述成 optional。

## 良好模式

```markdown
Use `workflow/target-domain/` when <specific trigger>. The current workflow owns <current boundary>; the target workflow owns <target boundary>. Pass <artifact name> with <required fields>. Keep <sensitive or target-specific data> in project docs.
```

## 連動更新

新增或修改 cross-workflow reference 時：

- 更新 referring workflow 的入口（`workflow/<domain>/execution-flow.md`），以及相關 `README.md`。
- 若 target workflow 需要辨識 incoming handoff，更新 target workflow 的入口。
- 如果該關係變成 recurring repo-wide rule，更新 [`linked-updates.md`](linked-updates.md)。
- 變更位於 `enforcement/`、`workflow/`、`analysis/` 或 `intelligence/` 時，依 configured tool sync 處理；具體工具命令放在 `ai-tools/`。

← [回到共用規則索引](README.md)
