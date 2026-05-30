# Mechanical Sanitization Validator

**Status**: `draft`
**世代**：Gen 3 runtime hardening（systemic gap remediation, 4th instance）
**建立日期**：2026-05-31
**最後更新**：2026-05-31（initial draft）
**Sibling plan**：[`2026-05-31-1900-workflow-activation-engine.md`](2026-05-31-1900-workflow-activation-engine.md) — same pattern (rule-without-executor)
**Empirical trigger**：2026-05-31 session — agent 在寫 `workflow-activation-engine` plan v1-v4 期間，向 canonical Ai-skill repo 多次 Write/Edit canonical 文件，內容夾帶 project incident details（specific filename、user 對話片段、領域 artifact 字串）。`enforcement/sanitization.md` + `enforcement/reusable-guidance-boundary.md` 規則明擺著，使用者三次追問才暴露 gap，最終 v5 patch 才人工抹除。

> 本 plan 不修個案 leak，而是補齊 **Mechanical Sanitization Validator** —— Ai-skill 第 4 個被識別出的「規則存在但無 mechanical executor」systemic gap。

---

## Decision Rationale

### Empirical Evidence（gap 體檢）

`scripts/ai-skill-cli/internal/app/hooks.go` grep 結果：

| Validator 類型 | 是否存在 mechanical hook |
|---|---|
| `gate.bootstrap.receipt_present`（Receipt + 必讀檔讀取 log） | ✅ 存在 |
| Commit-msg 19 validators（cognitive mode block / plan-status-sync / token budget 等） | ✅ 存在 |
| **Sanitization on Write/Edit to canonical paths** | ❌ **不存在** |
| **Project incident detection in reusable files** | ❌ **不存在** |
| **Absolute path / username pattern check** | ❌ **不存在** |

`enforcement/sanitization.md` 明確列必須移除的內容類型（line 4-12），但**沒有任何 PreToolUse / commit-msg hook 機械強制**。完全靠 agent 行為自律。

### Failure Pattern Classification（同 sibling plan）

| Gap 編號 | 規則檔 | 強制層 | 失誤次數（此 session） |
|---|---|---|---|
| #1 | workflow detector activation（routing-registry.yaml + dependency-reading.md §Workflow 編排） | behavioral | 1（travel-planning workflow 未觸發） |
| #2 | capability discovery activation（capability-discovery-philosophy.md） | behavioral | 1（detector miss 後未 fallback） |
| #3 | **sanitization on Write to canonical paths** | **behavioral** | **5 處 leak（v1-v4 plan 累積）** |

#1 + #3 同 systemic pattern。本 plan 處置 #3。

### 為什麼這個 gap 比 workflow detector 更危險

- workflow detector 失誤導致 review 品質不佳（可事後重做）
- sanitization 失誤導致 project incident **進入 commit history**（git 不可逆，public repo 後無法收回）
- 後者觸及 **P0 safety/privacy 邊界**，依 `enforcement/rule-weight.md` P0 規則：「不得因方便、速度或工具限制而繞過」

### Decision

建立 **Mechanical Sanitization Validator** 作為 PreToolUse + commit-msg 雙層 hook：

```
Agent calls Write/Edit on file under canonical repo paths
        ↓
PreToolUse: sanitization-prefilter
   ├─ scan content for: absolute paths, usernames, project incident keywords,
   │                    secrets, raw tokens, private hosts
   ├─ if hit → reject + suggest sanitized rewrite
   └─ pass → continue
        ↓
File written
        ↓
Commit-msg: sanitization-postcheck
   ├─ scan staged diff against same rules
   ├─ if hit → reject commit
   └─ pass → commit accepted
```

### Design Principles

| Decision | Rationale |
|---|---|
| **雙層強制**（PreToolUse + commit-msg） | PreToolUse 在寫入前擋下、減少返工；commit-msg 是 last line of defense（即使 PreToolUse 被 bypass，commit 仍會擋）。同 sibling plan 的 detector 也是雙層思維。 |
| **規則來自 canonical sanitization.md，不重新發明** | hooks.go 從 `enforcement/sanitization.md` 抽 banned patterns 表，sanitization rule 改 → validator 自動同步。 |
| **不檢查 project-local files** | `<PROJECT_ROOT>/.agent-goals/` 等本來就 project-local 的檔案不該被 sanitize 規則綁，否則 false positive 爆炸。validator 只掃 **canonical repo paths**（Ai-skill repo 內 `enforcement/`、`workflow/`、`analysis/`、`intelligence/`、`governance/`、`runtime/`、`plans/`、`feedback/extraction/` 等）。 |
| **黑名單而非白名單** | 列舉 forbidden patterns（絕對路徑 regex、known username patterns、placeholder 缺失），不是「所有可重用文件必須通過某內容白名單」—— 後者太嚴格會擋掉合理寫作。 |
| **PreToolUse reject 必須附建議改寫**，不能只 reject | agent 收到「<absolute-path-detected> at line N: replace with <placeholder>」之類具體訊息才能修正，否則迴圈試錯。 |

### Why Not Roll Into Sibling Plan

`workflow-activation-engine` 已涵蓋 4 個 phase 與 8 個 implementation step。把 sanitization validator 塞進去會：
- 模糊 plan scope（activation engine ≠ sanitization）
- 拖延 sibling plan close-out
- 違反 `governance/document-sizing.md`（單一 plan 應單一目的）

本 plan 獨立追蹤、可獨立 archive、與 sibling plan 並行進度。

---

## Architecture Compatibility Preflight

| 欄位 | 內容 |
|---|---|
| Candidate files | `scripts/ai-skill-cli/internal/app/hooks.go`（加 PreToolUse + commit-msg validators）、`enforcement/sanitization.md`（加 `machine_readable_patterns:` 區段供 validator 讀取）、`enforcement/sanitization.yaml`（新建，machine-readable spec）、新建 `enforcement/failure-patterns/sanitization-leak-on-canonical-write.md` |
| Source-of-truth | `enforcement/sanitization.md`（companion）+ `enforcement/sanitization.yaml`（canonical machine-readable，新建） |
| Compiler / generated surfaces | `runtime.db` 加 `sanitization_patterns` projection；`ai-skill runtime compile + refresh` 流程不變 |
| Layer responsibility | Validator 屬 runtime layer；Pattern spec 屬 enforcement layer |
| 衝突 | 無。本 plan 補的是 missing executor，不改 sanitization 規則本身 |
| `runtime.db` 影響 | 新增 1 個 projection；compile pipeline 加新 rule |

---

## Phase Plan

### Phase 0 — Preflight

#### Phase 0.0 — Open Questions 核對

- [ ] §Open Questions 全部標記處置

#### Phase 0.1 — Architecture Compatibility Preflight

- [ ] 確認 `enforcement/sanitization.md` 可拆 yaml sibling（同 `runtime/core-bootstrap.md` ↔ `runtime/core-bootstrap.yaml` pattern）
- [ ] 確認 `hooks.go` PreToolUse pipeline 可加 file-content scan validator
- [ ] 確認 canonical repo path 列表可從一個中央 source（建議 `enforcement/canonical-paths.yaml`）取得，validator 對該清單匹配時才 fire

### Phase 1 — Pattern Spec 形式化

把 `enforcement/sanitization.md` 內人類可讀的規則抽成 machine-readable YAML：

```yaml
# enforcement/sanitization.yaml
schema_version: 1
banned_patterns:
  absolute_path:
    - regex: '^[A-Z]:\\'                  # Windows absolute path
    - regex: '^/(Users|home)/[^/]+/'      # macOS/Linux user home
    - regex: '^/tmp/[^<]'                 # raw /tmp paths
  username:
    - regex: '\\Users\\[a-z0-9]+\\'       # Windows user dir
    - regex: '/Users/[a-z0-9]+/'          # macOS user dir
  secrets_token:
    - regex: 'Bearer [A-Za-z0-9_-]{20,}'
    - regex: 'AKIA[0-9A-Z]{16}'           # AWS access key
    - regex: 'ghp_[A-Za-z0-9]{36}'        # GitHub PAT
    # ... (extend)
  private_host:
    - regex: 'https?://\d+\.\d+\.\d+\.\d+'  # raw IPs
    - regex: '\.internal\.'
canonical_paths:
  - 'enforcement/'
  - 'workflow/'
  - 'analysis/'
  - 'intelligence/'
  - 'governance/'
  - 'runtime/'
  - 'plans/'
  - 'feedback/extraction/'
  - 'metadata/'
  - 'constitution/'
  - 'architecture/'
not_canonical:
  # validator 不掃這些 path（project-local 或專案 incident docs allowed）
  - '.agent-goals/'
  - 'feedback/history/'   # project-specific lessons 反而需要 evidence
  - 'plans/archived/'     # 已完成 plan 不再修改
project_incident_signals:
  # heuristic：當以下 keyword 集中出現於 canonical reusable doc 時 → warn
  # （這不是禁用列表，而是「太具體可能是 project incident」訊號）
  description: |
    Reusable canonical docs 通常用 generic example。當 single domain 的
    specific noun 在同一段落 cluster 出現（e.g., 5+ travel-specific nouns），
    暗示作者把 project incident 寫進來了，應提示重寫為抽象範例。
  threshold: 5  # 同一段落內同領域 specific noun 數量
```

產出：
- [ ] `enforcement/sanitization.yaml` 新建
- [ ] `enforcement/sanitization.md` 加 companion 章節指向 yaml + `machine_readable_patterns` 連結
- [ ] `runtime/runtime.db` compile pipeline 加 `sanitization_patterns` projection

### Phase 2 — PreToolUse Validator

新建 `hooks.go` validator `validateSanitizationOnWrite`：

簽名 + 行為：

```go
func validateSanitizationOnWrite(toolName string, params ToolParams) HookResult {
    if toolName != "Write" && toolName != "Edit" { return Pass }

    filePath := params.GetString("file_path")
    if !isCanonicalRepoPath(filePath) { return Pass }  // project-local 不掃

    content := params.GetString("content")  // Write 用
    if toolName == "Edit" {
        content = params.GetString("new_string")
    }

    hits := scanBannedPatterns(content, sanitizationPatterns)
    if len(hits) > 0 {
        return Reject(formatRejectionWithSuggestions(hits))
    }

    return Pass
}
```

Reject 訊息範例：
```
Sanitization gate rejected this write to <canonical-path>:

Line 7: matched absolute_path pattern `^[A-Z]:\\`
  Found: "C:\Users\actual-user\..."
  Suggest: replace with <PROJECT_ROOT> or <WORKSPACE> placeholder

Line 42: matched project_incident_signal (5 travel-specific nouns in same paragraph)
  Found: paragraph containing [<5 listed terms>]
  Suggest: abstract to <domain-keyword> placeholders, move concrete
           evidence to feedback/history/<domain>/ if it's a reusable lesson

See enforcement/sanitization.md for full rules.
```

產出：
- [ ] `hooks.go` validator + unit tests
- [ ] Pattern compile/cache（避免每次 Write 重 parse yaml）
- [ ] Hook registration in PreToolUse dispatcher

### Phase 3 — Commit-msg Validator（last line of defense）

新建 `validateSanitizationOnCommit`：

- 掃 staged diff（所有 added / modified lines）
- 對每行套用 `sanitization.yaml` patterns
- hit → reject commit + 列出 hit 位置

範例 opt-out：`[skip-sanitization]` trailer（極少用，例如 `feedback/history/` 例外，但實際上該目錄已在 `not_canonical` 列表）。

產出：
- [ ] `validateSanitizationOnCommit` + unit tests
- [ ] 加入 commit-msg validator registry（成為第 20 個 validator）
- [ ] 更新 `runtime/core-bootstrap.yaml` `per_commit_obligations` 加 `obligation.commit.sanitization_diff`

### Phase 4 — Failure Pattern + Documentation

- [ ] 新建 `enforcement/failure-patterns/sanitization-leak-on-canonical-write.md`
  - 記錄 2026-05-31 session 為 inaugural case
  - 同 sibling plan pattern：rule existed, behavioral enforcement only, no mechanical executor
- [ ] 更新 `enforcement/sanitization.md` companion 加章節「Mechanical enforcement reference」指向 validator + yaml + failure pattern

### Phase 5 — Validation Scenarios

新建 scenarios：
- `validation/scenarios/sanitization/pretooluse-rejects-absolute-path-v1.yaml`
- `validation/scenarios/sanitization/commit-rejects-leaked-username-v1.yaml`
- `validation/scenarios/sanitization/project-incident-signal-warns-v1.yaml`
- `validation/scenarios/sanitization/non-canonical-path-passes-through-v1.yaml`（確認 `.agent-goals/` etc. 不被擋）
- `validation/scenarios/sanitization/2026-05-31-regression-v1.yaml`（回放本 session 的 5 處 leak，全 reject）

Acceptance：五 scenario 全 PASS。

### Phase 6 — Close-out

- [ ] phase done
- [ ] `git status` clean
- [ ] `git push` 完成、`git log origin/main..HEAD` empty
- [ ] 讀回 sanitization.md / sanitization.yaml / failure pattern
- [ ] Archive 本 plan

---

## Open Questions

| # | Question | 處置 |
|---|---|---|
| Q1 | `not_canonical` 列表是否該 include `analysis/<domain>/sources-and-tools.md`？該檔可能合理引用外部 service 名稱 | still-open — 建議 default 掃，例外加 `[allow-domain-names]` marker |
| Q2 | PreToolUse reject 訊息語言：英文或中文？sibling plan workflow detector 同問題 | still-open — 建議 reject 訊息結構化（key=value），訊息文字依 locale |
| Q3 | `project_incident_signals` heuristic threshold（5 nouns/段落）是否太鬆 / 太緊？ | still-open — 建議 Phase 5 validation scenario 跑回放校準 |
| Q4 | hooks.go validator 共用 RuntimeContext from sibling plan 嗎？ | resolved → 不依賴。本 plan 的 validator 純基於 file content 與 path，不需要 workflow context |
| Q5 | 已 committed leak（v1-v4 in commit history）是否要 git history rewrite？ | still-open — 建議**不**做（git rewrite 風險高），但記在 failure pattern 作為 "irreversible incident" |

---

## Validation Plan

- [ ] Phase 1 yaml schema 經 review
- [ ] Phase 2 PreToolUse rejection 訊息實際對 agent 有幫助（reject + actionable suggestion）
- [ ] Phase 3 commit-msg validator 不誤殺正當 commit（false positive < 1%）
- [ ] Phase 4 failure pattern entry 同模式 actionable
- [ ] Phase 5 regression scenario 五個全 PASS
- [ ] Phase 6 close-loop

---

## Dependency Read Ledger

| 欄位 | 內容 |
|---|---|
| Trigger | 2026-05-31 session 使用者指出 sibling plan v1-v4 寫作期間 sanitization gate 未觸發、project incident 洩漏 |
| Required set | `enforcement/sanitization.md`、`enforcement/reusable-guidance-boundary.md`、`enforcement/rule-weight.md`、`enforcement/dependency-reading.md`、`runtime/core-bootstrap.yaml`、`scripts/ai-skill-cli/internal/app/hooks.go`（部分 grep 確認 validator 不存在）、sibling plan `2026-05-31-1900-workflow-activation-engine.md` |
| Read | 以上 |
| Not applicable | 無 |
| Deferred | Phase 2-3 實作細節 source（hooks.go 完整檔案、commit-msg dispatcher 細節）—— Phase 0.1 進入 implementation 前再補讀 |
| Validation | Architecture Compatibility Preflight 已列；Phase 0.1 unlock 前驗證 |

---

## Source

2026-05-31 session：
- 使用者連續追問五次，依序暴露：(1) sqlite3 vs ai-skill CLI 認知偏差，(2) `route.workflow.travel-planning` activation gap，(3) Discovery vs Detector 混淆，(4) intelligence 預設 advisory 風險，(5) **sanitization 自我觸發失敗導致 v1-v4 plan 寫作期間 project incident 洩漏**。
- 第五次追問識別出本 plan 處理的 systemic gap。

本 plan 是 sibling `workflow-activation-engine` plan 同模式問題（rule-without-executor）的第二個顯式 case，獨立追蹤是因為 scope 完全不同（sanitization vs routing）。

## Companion References

- `enforcement/sanitization.md` —— canonical 規則（companion markdown）
- `enforcement/reusable-guidance-boundary.md` —— project incident vs reusable rule 邊界
- `enforcement/rule-weight.md` §P0 —— sanitization 屬 P0 safety/privacy 邊界
- `enforcement/failure-patterns/bootstrap-bypass-on-resume.md` —— validator 採同 PreToolUse pattern 範例
- Sibling plan：`plans/active/2026-05-31-1900-workflow-activation-engine.md` —— 同 systemic gap pattern
