# AI-native Knowledge Operating System

AI 知識作業系統 — 讓 agent 用最小 token 成本找到正確知識。

## 🚀 Quickstart

```text
1. Read CORE_BOOTSTRAP.md  (3 rules, ~800 tokens)
2. Read this README         (OS layout)
3. Query skills-index.yaml  (find relevant skill)
4. Check activation-rules   (load lazy rules if needed)
5. Read knowledge summary   (300-500 tokens, then expand if needed)
```

詳細啟動流程：[`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md)

## 維護本 repository（人類貢獻者）

`skills-index.yaml` 與 `knowledge/indexes/` 等索引主要給 **agent 依任務路由**；若你要**修改 Ai-skill 本庫**（PR、治理、驗證指令），請從 **[`governance/contributing.md`](governance/contributing.md)** 進入（內含與 [`scripts/README.md`](scripts/README.md)、[`governance/validation/README.md`](governance/validation/README.md) 的連結）。GitHub 慣例入口：[`CONTRIBUTING.md`](CONTRIBUTING.md)。

## 📂 OS Layout

| 層級 | 路徑 | 用途 |
| --- | --- | --- |
| 🎯 **Core Bootstrap** | [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) | 最小必讀啟動集合（3 rules, ~800 tokens） |
| 📐 **Architecture** | [`architecture/`](architecture/) | Roadmap、升級規劃、成本優化 |
| ⚙️ **Shared Rules** | [`enforcement/`](enforcement/README.md) | 共用作業規則（含 lazy-load activation model） |
| 🧠 **Skills** | [`skills/`](skills/README.md) | 可重用 agent 能力模組 |
| 🗺️ **Skill Index** | [`skills-index.yaml`](skills-index.yaml) | 結構化 skill routing index |
| 🔧 **Tool Adapters** | [`ai-tools/`](ai-tools/README.md) | Claude Code、Cursor 等工具配置 |
| 🔄 **Runtime** | [`runtime/`](runtime/README.md) | Context routing、activation、TTL |
| 🧭 **Knowledge** | [`knowledge/`](knowledge/README.md) | Indexes、summaries、graphs、runtime surfaces |
| 📊 **Metadata** | [`metadata/`](metadata/README.md) | Knowledge Atom schema、ranking、confidence |
| 🧪 **Analysis** | [`analysis/`](analysis/README.md) | 觀察、拆解、pattern extraction |
| 💡 **Intelligence** | [`intelligence/`](intelligence/README.md) | Engineering decision、trade-off、anti-pattern |
| 🔄 **Workflow** | [`workflow/`](workflow/README.md) | Planning、decomposition、execution flow |
| 💾 **Memory** | [`memory/`](memory/README.md) | Episodic、project、failure memory |
| 📝 **Feedback** | [`feedback/`](feedback/README.md) | Lesson extraction、promotion、feedback loop |
| 🤖 **Models** | [`models/`](models/README.md) | Model capability profile、routing、compression |
| 🏛️ **Governance** | [`governance/`](governance/README.md) | Lifecycle、cleanup、validation |
| 📜 **Scripts** | [`scripts/`](scripts/README.md) | Close-loop automation、runtime refresh |

## 🧭 Agent 作業流程

```
Session Start
  │
  ├─ 1. Read CORE_BOOTSTRAP.md (3 rules, ~800 tokens)
  │
  ├─ 2. Read README.md (OS layout)
  │
  ├─ 3. Query skills-index.yaml → find skill by task intent
  │
  ├─ 4. Check activation-rules.yaml → load lazy rules if triggered
  │
  ├─ 5. Read skill summary (300-500 tokens)
  │
  ├─ 6. Expand to full source only if needed
  │
  └─ 7. Use .agent-goals/ for multi-step tasks
```

## 🛠️ AI Tools

| Tool | Config |
| --- | --- |
| **Roo Code** | [`ai-tools/agent/roo.md`](ai-tools/agent/roo.md) |
| **Cursor** | [`ai-tools/agent/cursor.md`](ai-tools/agent/cursor.md) |
| **Claude Code** | [`ai-tools/agent/claude.md`](ai-tools/agent/claude.md) |

## 📖 Key Documents

| Document | Purpose |
| --- | --- |
| [`governance/contributing.md`](governance/contributing.md) | 人類維護入口：驗證指令、PR gate、文件索引 |
| [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) | Minimal bootstrap (3 rules, ~800 tokens) |
| [`skills-index.yaml`](skills-index.yaml) | Skill routing index |
| [`plans/archived/2026-05-11-next-stage-upgrade-plan.md`](plans/archived/2026-05-11-next-stage-upgrade-plan.md) | Full architecture upgrade plan |
| [`plans/archived/2026-05-12-context-cost-optimization.md`](plans/archived/2026-05-12-context-cost-optimization.md) | Token cost optimization plan |
| [`runtime/router/activation-rules.yaml`](runtime/router/activation-rules.yaml) | Lazy-load activation rules |
| [`runtime/context/ttl-policy.yaml`](runtime/context/ttl-policy.yaml) | Context TTL policy |

## 📌 新專案快速啟用

### 給人類（開新專案的人）

如果你開了一個**全新的專案**，想讓它使用此知識庫，執行以下命令一次設定所有 AI 工具：

```bash
# 從 Ai-skill repo 目錄執行
./scripts/init-new-project.sh /path/to/your/new-project

# 範例
./scripts/init-new-project.sh ~/projects/my-new-app
```

這會在目標專案中建立：

| 工具 | 產出 | 效果 |
|------|------|------|
| **Roo Code** | `.roomodes` | 5 個 mode，含語言規則 + 知識更新 checkpoint |
| **Cursor** | `.cursor/rules/ai-skill-bootstrap.mdc` | alwaysApply 規則，自動載入 |
| **Claude Code** | `CLAUDE.md` | 自動載入 Core Bootstrap |
| **通用** | `.agent-goals/` | 對話目標帳本 |

詳細說明：[`ai-tools/new-project-onboarding.md`](ai-tools/new-project-onboarding.md)

### 給 AI agent（session 啟動時）

啟動流程已在 [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) 定義（含新專案自動偵測），依序執行即可。

## 📋 Rules

- **Reference-first**: Agent 直接讀本 repository，不依賴 tool mirror。
- **Path convention**: 使用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>` 占位符，不寫入本機絕對路徑。
- **Close-loop**: 修改後必須 diff review、linked updates、commit、push、readback、clean status。
- **Cost-aware**: 優先讀 summary（300-500 tokens），需要才展開全文。
