# 文件 TODO 清單

可重用文件可以在前段放一個短 TODO 清單，讓人和 agent 不必讀完整份文件，就能立即看出哪些內容尚未完成。

本規則處理 document-local unfinished work，並補充 [`conversation-goal-ledger.md`](conversation-goal-ledger.md)；後者追蹤跨文件、跨 session 的 conversation-level goals。

## 何時新增

當文件有下列任一情況時，加入 TODO 區塊：

- 未完成章節。
- 缺少 examples、tables、templates、links 或 references。
- 仍需驗證的 claims。
- 已知需要補強的弱點。
- 影響文件可用性的 follow-up work。
- 文件完成前必須解決的 open questions。

若文件已完成且沒有 document-local open work，可省略 TODO 區塊，或用簡短文字寫 `No open document TODOs`。

## 放置位置

TODO 區塊要放在詳細正文之前、容易看見的位置：

1. YAML/frontmatter（若有）。
2. 標題。
3. 簡短 purpose/overview。
4. `## Document TODO` 或 `## TODO`。
5. 主體內容。

不要把 document TODOs 藏在文件最後。目的在於長對話或 handoff 後能立即定位。

## 模板

使用精簡表格：

```markdown
## Document TODO

| Priority | Status | TODO | Link | Owner / Goal |
| --- | --- | --- | --- | --- |
| P1 | pending | Add validation checklist for package sync | [Validation](#validation) | `.agent-goals/goals/P1-example.md` |
```

建議欄位：

| 欄位 | 意義 |
| --- | --- |
| `Priority` | `P0`、`P1`、`P2` 或 `P3`，對齊 goal ledger priority vocabulary。 |
| `Status` | `pending`、`in_progress`、`blocked`、`needs-validation`、`done` 或 `cancelled`。 |
| `TODO` | 具體未完成工作，不是模糊提醒。 |
| `Link` | 同文件相關章節 anchor，或相關檔案/章節。 |
| `Owner / Goal` | 可選 owner、goal file、issue 或 todo ID，用來說明誰應決定或完成。 |

## Linking Rules

每個 TODO 都應指向可行動的位置：

- 工作屬於某章節時，連到同文件 heading。
- 缺口在另一份文件時，連到該文件。
- TODO 屬於 conversation-level active goal 時，連到 `.agent-goals/goals/<goal-id>.md`。
- 若 issue、planning document 或 checklist item 是 source of truth，連到該項目。

若暫時沒有有用連結，寫 `needs anchor`，並在關閉 TODO 前建立相關章節。

## 與 Goal Ledger 的關係

Document TODOs 與 goal ledger entries 應互相補強：

- Document TODO 是單一文件內的局部缺口。
- Goal ledger entry 追蹤跨文件、工具或 session 的 user-facing objective。
- 如果 document TODO 變成 user-facing objective 或跨多檔，建立或連到 goal ledger entry。
- 若 goal ledger entry 依賴某文件章節，從 goal 連回 document TODO 或 heading。

Goal 完成時，要更新或移除相關 document TODO。若 document TODO 仍未處理，除非 goal completion criteria 明確排除該 TODO，否則不要刪除 linked goal。

## 維護

- TODO 表要短。大型 task breakdown 移到 `.agent-goals/`、issue tracker 或 planning document。
- 只有 linked work 真正完成並驗證後，才移除或標成 `done`。
- 不清楚的工作優先標成 `blocked` 或 `needs-validation`，不要直接刪除。
- Review 時，宣稱文件完成前先檢查 TODO 區塊。

← [Back to enforcement index](README.md)
