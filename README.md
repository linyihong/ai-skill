# Ai-skill

這個 repository 是這台機器上的共用 AI skill 知識庫。

未來不同專案需要 APK 分析、抓包、Frida、Proxyman、Dart AOT、解密或分析方法沉澱時，優先從這裡讀 skill，並把新的可重用技巧回饋回這裡。

## 現有 Skills

| Skill | 用途 |
| --- | --- |
| `apk-analysis/` | 授權 APK 流量分析、動態抓包、Flutter/Dart AOT、response 解碼、文件化與技巧回饋。 |

## 使用方式

在任一專案中，可以要求 AI 先讀。請把 `<AI_SKILL_REPO>` 替換成你本機 clone 這個 repository 的位置：

```text
<AI_SKILL_REPO>/apk-analysis/SKILL.md
<AI_SKILL_REPO>/apk-analysis/RUNBOOK.md
```

建議提示：

```text
請先閱讀 <AI_SKILL_REPO>/apk-analysis/SKILL.md 和 RUNBOOK.md。
之後依 WORKFLOW.md / TOOLS.md / DOCUMENTATION.md 分析這個 APK。
如果發現新的可重用技巧，請回饋到 <AI_SKILL_REPO>/apk-analysis/FEEDBACK.md。
回饋內容要讓人類也能看懂，不只給 AI 看。
完成後請在 <AI_SKILL_REPO> 建立 commit 並 push。
```

## 回饋規則

可以回饋：

- 新的分析流程。
- 新的工具使用技巧。
- 新的失敗判讀方式。
- 新的去敏與 fixture 沉澱方式。
- 通用的媒體、解密、session、proxy、Frida 經驗。

不要回饋：

- 特定產品的完整 host / endpoint。
- token、secret、device id、帳號或個資。
- 未去敏 raw response。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置；請改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>`。
- 只對單一產品有效、沒有泛化價值的結論。

## Git 規則

這個 repository 有自己的 git remote。更新 skill 後：

```bash
cd <AI_SKILL_REPO>
git status
git add .
git commit -m "Update apk analysis skill"
git push
```

不要把專案私有資料、抓包原始檔或未去敏樣本 commit 到這裡。
