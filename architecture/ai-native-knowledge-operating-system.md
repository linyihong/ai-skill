# AI-native Knowledge Operating System

本文件定義 AI-native Knowledge Operating System 的 repository-level architecture direction。它是 roadmap 與 boundary document，不是 executable shared rule。可執行政策仍放在 `enforcement/`；工具專屬設定仍放在 `ai-tools/`。

下一階段完整升級規劃見 [`next-stage-upgrade-plan.md`](../plans/archived/2026-05-11-next-stage-upgrade-plan.md)。

## 目的

AI-native Knowledge Operating System 是讓 agents 能可靠載入知識、追蹤 goals、改善 reusable guidance、驗證變更，並關閉 git writeback loops 的操作層；它不把任何單一 AI tool 或 local mirror 當成 source of truth。

預設模型是 **reference-first**：

1. Agents 直接讀 canonical `<AI_SKILL_REPO>`。
2. `enforcement/` 提供 operating rules。
3. `skills/` 提供 capability modules。
4. `ai-tools/` 與 tool config 提供特定工具 adapters。
5. `.agent-goals/` 記錄 project-local temporary execution state。
6. Scripts 支援 validation、goal tracking、commit/push/readback 與 optional tool sync。

`symlink`、`bundle` 與 `copy snapshot` flows 仍是 compatibility layers，用於需要 native scanning 或 offline snapshots 的工具。它們不是 default architecture，也不能取代直接讀 canonical repository。

## 層級

| 層級 | 位置 | 職責 |
| --- | --- | --- |
| Source of truth | `<AI_SKILL_REPO>` git repository | Canonical rules、skills、templates、scripts 與 architecture docs。 |
| Operating rules | `enforcement/` | Cross-skill policy、dependency reading、linked updates、goal ledger rules、failure learning、validation 與 documentation boundaries。 |
| Capability modules | `skills/` | Domain-specific workflows、checklists、documentation templates、techniques 與 feedback lessons。 |
| Tool adapters | `ai-tools/`、tool config、optional skill adapters | Tool-specific paths、hooks、UI behavior、sync strategy 與 troubleshooting。 |
| Temporary execution state | `<PROJECT_ROOT>/.agent-goals/` | Active user goals、owner/lock decisions、open work、next actions 與 completion validation。 |
| Close-loop automation | `scripts/` | Goal ledger operations、grouped commits、push、readback support 與 optional tool sync 的保守 helpers。 |

## Reference-first 載入

`reference-first` 表示 agent 應以路徑載入中央 repo，而不是依賴每個專案內複製的 skill package。

最小啟動形狀：

```text
<AI_SKILL_REPO>/enforcement/README.md
<AI_SKILL_REPO>/skills/<skill-name>/SKILL.md
```

完成 shared-rule bootstrap 後，agent 依 dependency routing 只讀目前任務需要的 skill workflows、tool adapters、feedback lessons、templates 或 architecture docs。

這讓更新流程保持簡單：

- 更新 canonical repo。
- 驗證 linked updates。
- Commit 並 push。
- Push 後讀回更新入口。
- 讓未來 sessions 直接讀 canonical repo，而不是刷新 copied packages。

## 相容層

只有 active tool 無法可靠 reference central repo，或需要 native scan behavior 時，才使用 compatibility layers。

| Strategy | 使用時機 | Close-loop rule |
| --- | --- | --- |
| Reference-first | 一般 agent work 的預設。 | 不需要 tool mirror sync；確認 canonical repo 是最新且可讀。 |
| Symlink / bundle bridge | 工具需要 native skill discovery，但 local paths 可指回 `<AI_SKILL_REPO>`。 | 只有此 bridge 被有意使用或使用者明確要求時才 sync。 |
| Copy snapshot | 工具不能讀 central repo、不能用 symlink，或需要 offline snapshot。 | 記錄 source commit/date，並有意識地 refresh；不可把 copy 當成 source。 |

Tool-specific details 屬於 `ai-tools/`。Generic rules 應使用「configured tool sync」或「optional tool sync」，不要把單一工具寫成 default。

## Current Compatibility Inventory

本 inventory 記錄目前仍有意提到 native scan、symlink、bundle 或 copy snapshot 行為的 compatibility surfaces。

| Surface | 目前角色 | 保留條件 | 移除或 deprecation 信號 |
| --- | --- | --- | --- |
| `ai-tools/agent/cursor.md` | Cursor adapter，說明 reference-first loading、Core Bootstrap 與語言偏好設定。 | Cursor 仍是 active tool adapter。 | 保留作為 tool adapter。 |
| `scripts/sync-cursor-bundle.sh` | `~/.cursor/bundles` 與 `~/.cursor/skills` 的 optional symlink/bundle bridge。 | 任何 local setup 仍需要 Cursor native scan 或 bundle paths。 | 當 native scan workflows 有 documented reference-first 或 symlink-free replacement 時 deprecate。 |
| `scripts/git-hooks/post-commit` with `AI_SKILL_SYNC_CURSOR_BUNDLE=1` | 明確 opt in 的 post-commit mirror refresh。 | Users 想在 commits 後自動 refresh local mirror。 | `sync-cursor-bundle.sh` deprecated 後，或沒有 active setup export 該 opt-in flag 時移除。 |
| `dependency-reading.md`、`linked-updates.md`、`failure-learning-system.md` 與 `failure-patterns/source-mirror-write-drift.md` 中的 source/mirror guardrails | 防止 `.cursor`、`~/.cursor`、bundles 或 generated copies 被誤認為 source updates。 | 任何 mirror、runtime copy 或 tool deployment surface 仍存在。 | 即使 bundle scripts 移除，只要 tool adapter 仍可能建立 runtime copies，就保留。 |

Inventory rule：compatibility surfaces 可以保留，但每一項都必須說明存在原因、使用時機，並維持 `reference-first` 是預設。

## 遷移 Roadmap

### Phase 1: Reference-first default

- 保持 `enforcement/README.md` 作為 bootstrap index。
- Tool docs 必須清楚說明 `reference-first` 不需要 copying 或 bundle sync。
- Close-loop automation 保持保守：canonical repository changes 必須 commit/push/readback；tool sync 是條件式。
- 保留 compatibility scripts，供仍需要 symlink、bundle 或 copy snapshot workflows 的 users 使用。

### Phase 2: Compatibility Inventory

- 當 native scan、bundle、symlink 或 copy snapshot references 改變時，維護上方 current compatibility inventory。
- 辨識仍依賴 native tool scanning 或 copied skill directories 的 active workflows。
- 在相關 `ai-tools/` 文件記錄每個剩餘 compatibility case。
- 工具允許時，優先使用 symlink 或 reference strategies，而不是 copy snapshots。
- 當 stale docs 讓 copy 或 bundle sync 看似 default path 時，更新文件。

### Phase 3: Deprecation Readiness

所有條件成立時，copy 與 bundle flows 才可 deprecated：

- Active tool docs 指向 `reference-first` 作為 normal path。
- 沒有 active project workflow 在日常使用中需要 copied skill packages。
- 剩餘 native-scan needs 已由 symlink/reference adapters 或明確例外處理。
- Close-loop validation 不再依賴 mirror refresh checks，除非使用者要求 compatibility tests。
- 移除 scripts 或 docs 前，已有 deprecation note 與 migration path。

### Phase 3 Deprecation Checklist

移除或封存 copy、bundle 或 native-scan compatibility paths 前，使用此 checklist：

| Gate | 檢查 | 證據 |
| --- | --- | --- |
| Scope | 辨識要 deprecated 的確切 surface：copy snapshot docs、bundle sync script、post-commit hook、native-scan setup 或 tool mirror path。 | 編輯前列出 affected files 與 owner group。 |
| Search | 搜尋 root docs、`enforcement/`、`ai-tools/`、`skills/`、scripts、tool rules 與 failure patterns 中是否仍有 required/default usage。 | 搜尋結果顯示沒有 required/default usage；或剩餘引用都明確標為 compatibility-only。 |
| Replacement | 為每個剩餘 user workflow 記錄 reference-first、symlink 或 manual snapshot replacement。 | `ai-tools/` 或 architecture docs 指向 replacement。 |
| Source boundary | 確認 replacement 不把 `.cursor`、`~/.cursor`、bundles、generated files 或 snapshots 當成 canonical source。 | Source/mirror guardrails 仍從 `dependency-reading.md` 與 failure patterns 連到。 |
| Script behavior | 若移除 script，同步更新 script docs、hooks、close-loop automation 與 env var examples。 | `scripts/README.md` 與相關 tool docs 不再把 removed commands 當成 active paths。 |
| Validation | Commit 前跑 link checks、stale wording search、diff review 與 close-loop dry run。 | Validation output 或 final response 說明哪些檢查通過。 |
| Rollback | 保留足夠 history 或 migration notes，讓仍需 manual compatibility path 的 users 可回退。 | Deprecation note 或 release note 說明替代方案。 |

## Copy And Bundle Sync Removal Criteria

不要只因為 `reference-first` 存在就移除 compatibility scripts。只有在下列條件完成後，才移除或 archive：

1. 搜尋確認沒有 required workflow 預設要求 agents 執行它們。
2. Tool-specific docs 記錄 replacement strategy。
3. 依賴 native scanning 的 users 有 symlink/reference path，或明確承認 manual snapshot path。
4. `enforcement/dependency-reading.md` 與 `enforcement/linked-updates.md` 仍說明如何處理必要 tool sync，但不把它 universal 化。
5. Close-loop process 已驗證不會發生 accidental source/mirror drift。

## 與既有文件的關係

- `enforcement/` 仍是 executable policy layer。
- `skills/` 仍是 capability layer。
- `ai-tools/` 仍是 tool adapter layer。
- `scripts/` 仍是 helper automation，不是 architecture 本身。
- `.agent-goals/` 仍是 temporary project state，完成後刪除。

當 repository 改變 agent 如何載入 rules、skills 如何被 discover、goal state 如何追蹤，或 source/mirror boundaries 如何 enforcement 時，更新本 architecture document。
