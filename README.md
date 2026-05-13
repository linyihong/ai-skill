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

## 📂 OS Layout

| 層級 | 路徑 | 用途 |
| --- | --- | --- |
| 🎯 **Core Bootstrap** | [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) | 最小必讀啟動集合（3 rules, ~800 tokens） |
| 📐 **Architecture** | [`architecture/`](architecture/) | Roadmap、升級規劃、成本優化 |
| ⚙️ **Shared Rules** | [`shared-rules/`](shared-rules/README.md) | 共用作業規則（含 lazy-load activation model） |
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
| **Claude Code** | [`ai-tools/agent/claude.md`](ai-tools/agent/claude.md) |
| **Cursor** | [`ai-tools/agent/cursor.md`](ai-tools/agent/cursor.md) |

## 📖 Key Documents

| Document | Purpose |
| --- | --- |
| [`CORE_BOOTSTRAP.md`](CORE_BOOTSTRAP.md) | Minimal bootstrap (3 rules, ~800 tokens) |
| [`skills-index.yaml`](skills-index.yaml) | Skill routing index |
| [`architecture/next-stage-upgrade-plan.md`](architecture/next-stage-upgrade-plan.md) | Full architecture upgrade plan |
| [`architecture/context-cost-optimization-plan.md`](architecture/context-cost-optimization-plan.md) | Token cost optimization plan |
| [`runtime/router/activation-rules.yaml`](runtime/router/activation-rules.yaml) | Lazy-load activation rules |
| [`runtime/context/ttl-policy.yaml`](runtime/context/ttl-policy.yaml) | Context TTL policy |

## 📌 Quickstart for New Projects

```text
Use the AI-native Knowledge Operating System.

Canonical repository:
<AI_SKILL_REPO>

First read:
<AI_SKILL_REPO>/CORE_BOOTSTRAP.md
<AI_SKILL_REPO>/README.md

Then query:
<AI_SKILL_REPO>/skills-index.yaml

For multi-step tasks, initialize:
<PROJECT_ROOT>/.agent-goals/
```

## 📋 Rules

- **Reference-first**: Agent 直接讀本 repository，不依賴 tool mirror。
- **Path convention**: 使用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>` 占位符，不寫入本機絕對路徑。
- **Close-loop**: 修改後必須 diff review、linked updates、commit、push、readback、clean status。
- **Cost-aware**: 優先讀 summary（300-500 tokens），需要才展開全文。
