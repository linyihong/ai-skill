# AI-native Cognitive Execution System Overview

| 欄位 | 值 |
| --- | --- |
| kind | overview |
| audience | human |
| stability | active |
| routing | leaf |

AI-native Cognitive Execution System 是一套 AI-native 認知執行框架，用可攜、可版本化、可演化的結構累積工程知識，而不被單一 Agent、模型或工具綁架。

這份 overview 給第一次接觸本 repository 的人類讀者。Agent bootstrap、runtime obligations、receipt 格式、per-turn rules 與 SQLite table semantics 不在本文維護；需要執行時請回到 [`../CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) 與 [`../runtime/core-bootstrap.yaml`](../runtime/core-bootstrap.yaml)。

## 它是什麼

AI-native Cognitive Execution System 把工程經驗整理成 agent 可讀、可路由、可驗證的 source repository。它不是把所有內容塞進一個 prompt，而是把不同類型的知識放在不同層：

- `workflow/`：可執行流程與 checklist。
- `analysis/`：觀察、拆解與分析方法。
- `intelligence/`：可重用的工程判斷、heuristics、trade-offs 和 anti-patterns。
- `enforcement/`：所有 agent 都必須遵守的共用規則。
- `runtime/`：machine-readable contracts、state 和 gates。
- `knowledge/`：routing、summaries 和 graphs，讓 agent 不必每次讀完整 repository。

它的核心不是「提供一個 agent」，而是讓工程知識能被多個 agent 穩定執行。

## 解決什麼問題

現代 AI 開發 workflow 常見三個問題：

1. **知識被工具綁住**：團隊規則散在 IDE rules、聊天記憶、prompt snippets、agent-specific settings 裡，換工具時很難遷移。
2. **經驗無法累積**：一次 review、debug 或 incident 學到的東西，沒有被整理成下一個 agent 會自動遵守的結構。
3. **執行不可驗證**：就算 prompt 說了「請遵守規則」，也很難知道 agent 是否真的讀了正確 source、走了正確 workflow、完成了必要驗證。

本系統用 repository、routing、contracts、gates 和 validation loop 把這些問題收斂成可維護的知識結構。

## 核心價值

### Knowledge Ownership

團隊的工程知識應該屬於團隊，而不是某個模型供應商、hosted memory 或 IDE extension。把知識放在 repository 裡，代表它可以被 review、version control、diff、rollback、fork 和遷移。

### Agent Portability

不同 agent 的能力與操作介面會變，但它們都可以讀 source files、遵守 routing、執行 workflow、產出 validation evidence。這讓同一套工程知識能接到 Cursor、Claude Code、Codex、Roo Code 或未來 runtime。

### Human Experience Compounding

工程團隊每天都在累積經驗：哪些架構選型容易過度設計、哪些驗證容易漏、哪些 agent 失誤會重複發生。AI-native Cognitive Execution System 讓這些經驗不只停在人的記憶裡，而是變成下一次 agent 執行時可用的 cognitive structure。

## 和 prompt engineering 的差異

Prompt engineering 通常是把指令塞進一次對話或一段固定提示詞。它適合短期引導，但不擅長管理大型、長期、會演化的工程知識。

本系統把 prompt 視為輸入界面之一，而不是知識的唯一容器。真正的 source-of-truth 放在可維護的文件、contract、workflow、index 和 runtime state 中。

## 和 MCP / hosted memory 的差異

MCP 解決的是 agent 如何連接工具、資料與外部能力。Hosted memory 解決的是某個 agent runtime 如何保存記憶。

AI-native Cognitive Execution System 解決的是另一個問題：工程知識本身應該如何組織、版本化、路由、驗證，以及如何避免被單一 runtime 鎖住。它可以搭配 MCP 或 hosted memory，但不把長期可重用知識交給它們當唯一 source。

## 和 LangGraph / agent framework 的差異

LangGraph 或其他 agent framework 主要定義 agent 如何執行狀態圖、工具調用或多步驟流程。它們是 runtime / orchestration layer。

本系統關注的是 cognitive layer：哪些工程知識要被載入、什麼 source 是 canonical、哪些行為必須被 enforcement、失敗如何回饋成可重用規則，以及完成時需要什麼 validation evidence。

換句話說，agent framework 可以執行流程；AI-native Cognitive Execution System 定義工程知識如何變成可執行、可驗證、可遷移的流程。

## 和一般 workflow automation 的差異

一般 automation 偏向 deterministic steps：跑測試、發佈、產生報告、同步檔案。

本系統處理的是 agent cognition：在資訊不完整、需求模糊、文件分層、規則衝突或驗證不確定時，agent 應該讀什麼、如何判斷、何時停下、如何回報 evidence。它也可以調用 automation，但核心是讓 AI 的工程判斷有可追溯結構。

## 下一步

- 想把本系統接到新專案：讀 [`../ai-tools/new-project-onboarding.md`](../ai-tools/new-project-onboarding.md)。
- 想理解目前世代架構：讀 [`../architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md)。
- 想維護本 repository：讀 [`../governance/contributing.md`](../governance/contributing.md)。
- 你是 AI agent：從 [`../CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md) 進入。
