# Cross-skill references（跨 skill 引用）

當一個 skill 需要另一個 skill 的 output contract、documentation format、validation checklist 或 implementation guidance 時，使用 cross-skill reference。

## 原則

Skills 可以引用其他 skills，但不能複製另一個 skill 的完整規則。Referring skill 應說明何時讀取 target skill、期待什麼輸出，以及目前 skill 保留哪些 ownership boundary。

## 何時引用另一個 skill

符合下列情況時，加入明確 cross-skill reference：

- 某 workflow 產生的 artifact 需要由另一個 skill 消費。
- 某 skill 需要另一個 skill 的 template、checklist 或 contract format。
- 某 finding 應轉成另一個領域的 guidance。
- 使用者要求從一個 skill handoff 到另一個 skill。

例：`apk-analysis` 負責授權 APK evidence、traffic attribution、schema recovery 與 sanitized feature reconstruction handoff。`app-development-guidance` 負責把 handoff 轉成 BDD、Domain Model Contract、API / Interface Contract、Error Handling Contract、implementation slices、checks 與 tests。

## 必要引用格式

新增 cross-skill reference 時，包含：

| 欄位 | 要求 |
| --- | --- |
| Target skill | 連到 `../<skill-name>/` 或要讀取的特定檔案/資料夾。 |
| Trigger | 明確寫出 agent 何時應讀 target skill。 |
| Expected input/output | 命名 skills 之間傳遞的 artifact，例如 handoff、checklist、contract、fixture 或 implementation guidance。 |
| Ownership boundary | 說明哪些仍屬於目前 skill，哪些屬於 target skill。 |
| Sanitization boundary | 說明 target-specific data、secrets、raw evidence 或 product conclusions 是否必須留在 project docs。 |
| Linked updates | 更新雙方 skill 入口；若 target skill 已涵蓋 handoff，說明無需更新的理由。 |

## 不要做

- 不要把另一個 skill 的完整 workflow 貼進目前 skill。
- 不要讓每個 skill 預設讀取所有其他 skills。
- 不要建立循環 `always read` 鏈。
- 不要因為兩個 skills 互相引用，就把 target-specific evidence 搬進 reusable skill。
- 如果 requested output 必須 cross-skill handoff，不要把它描述成 optional。

## 良好模式

```markdown
Use `../target-skill/` when <specific trigger>. The current skill owns <current boundary>; the target skill owns <target boundary>. Pass <artifact name> with <required fields>. Keep <sensitive or target-specific data> in project docs.
```

## 連動更新

新增或修改 cross-skill reference 時：

- 更新 referring skill 的 `SKILL.md`，以及相關 `README.md`、`WORKFLOW.md`、`DOCUMENTATION.md` 或 technique file。
- 若 target skill 需要辨識 incoming handoff，更新 target skill 的入口。
- 如果該關係變成 recurring repo-wide rule，更新 [`linked-updates.md`](linked-updates.md)。
- 變更位於 `shared-rules/` 或 `skills/` 時，依 configured tool sync 處理；具體工具命令放在 `ai-tools/`。

← [回到共用規則索引](README.md)
