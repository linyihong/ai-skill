# Wish-to-Task-List Translation（將願望翻譯成任務單）

Status: candidate-intelligence
Layer: `intelligence/engineering/agent-architecture/`

## 主張

> 當 user 給 agent 的 prompt 是「願望」（vague wish: 「想讓 X 變好」「修一下這個」「讓它變快」），agent 必須在開始執行前**先翻譯**為三段式任務單（measure / process / verify），否則進入 shotgun mode、亂試一輪。

## 為什麼成立

Wish 與 task list 的差別：

| 維度 | Wish | Task list |
|---|---|---|
| 目標 | 模糊（「快一點」） | 量化（「p95 < 200ms」） |
| 步驟 | 隱含（agent 自填） | 明確（n 個 ordered step） |
| 驗證 | 缺席（「應該好了」） | 內建（before / after 對照） |
| 失敗信號 | 無法判斷（「再試試？」） | 明確（某 step 沒達標） |

Agent 接到 wish 時的失敗模式：

1. **強行視為 task list**：把 wish 拆成自己直覺的步驟，但步驟沒驗證標準 → shotgun
2. **無限拷問**：把所有不確定都問回去 → user 嫌煩，下次更模糊
3. **正確路徑（本 atom 主張）**：自填 task list 初稿 + 點出量化空白讓 user 補

## 觸發信號（wish vs task）

判斷 prompt 是否為 wish：

- 動詞含「想」「希望」「讓」「變得」「修一下」（語氣模糊）
- 形容詞無量化（「快」「順」「好」未對應數字或閾值）
- 無 success criteria（「修好」未定義）
- 無 baseline 引用（不提到目前的狀態數字）

任一條成立 → agent 應視為 wish，套用本 atom。

## Required Agent Action（wish 進來時）

### Step 1：先回顯 task list 初稿

```
我把你的需求翻成任務單：
- 目標（measure）：[嘗試量化的版本]
- 步驟（process）：[初擬 n 步]
- 驗證（verify）：[before / after 對照標準]

確認後我才開始。
```

### Step 2：標出 must-fill blanks

如果某些數字 / 標準確實必須 user 決定（例如「p95 要多少才算 OK」），不要瞎猜，明確列為 blank：

```
以下 2 點需要你決定：
- p95 latency 目標（目前 800ms，要降到多少？）
- 可接受的 perf budget（願意多花多少記憶體換速度？）
```

### Step 3：執行 + verify

只有當 task list 確認後才動手；每執行完一步即時回報 verify 結果，不集中到結尾才報告。

## 與既有 repo 內容的差別

repo 已有：
- [`intelligence/engineering/requirements/`](../requirements/README.md) → 處理「需求是什麼」的 ambiguity resolution
- [`workflow/software-delivery/requirements/`](../../../workflow/software-delivery/requirements/) → 處理 acceptance criteria 的執行流程

本 atom 補的洞：
- 範圍更廣（不只 software-delivery，任何 troubleshooting / optimization / 操作類任務）
- 焦點是「prompt 翻譯責任」而非「需求工程方法論」
- 與 cognitive contract v2 的 `validation_mode` 對應：non-NONE 即此原則的 runtime 投影

## 反模式：直接執行 wish

```
User: 「想讓網路變快」
反模式 agent：直接執行
  → 重啟 Wi-Fi
  → 改 DNS
  → 移除 network profiles
  → 升級驅動
  → 全部完成後說「應該變快了，你試試看」
問題：哪個有效？user 怎麼判斷「快」？下次再慢從哪裡查？
```

正模式：

```
User: 「想讓網路變快」
正模式 agent：
  「我把任務翻譯為：
   - measure: 先跑 speedtest-cli + ping baseline
   - process: 列 8 個常見瓶頸，逐一檢查
   - verify: 跑前後同樣的 speedtest + ping，對照數字
  
   你的『快』有具體目標嗎？例如下載 100Mbps 以上、ping < 30ms？」
```

## 與 shotgun-debugging 的關係

兩個 atom 是同一現象的兩面：

- 本 atom（agent-architecture/）：prompt 接收端，agent 該如何防止
- [`shotgun-debugging`](../anti-patterns/shotgun-debugging.md)（anti-patterns/）：執行端，反模式本身命名與識別

User 給 wish + agent 不翻譯 → shotgun debugging。

## Prevention Gate

- **Agent 自檢**：每次 user prompt 進來，先用本 atom 的 4 個 wish 信號判斷，是 wish 就先翻譯
- **PR / commit 層**：若 commit body 含「修好了」「應該 OK」但無 before / after，視為違反本原則
- **Reviewer 層**：reviewer 看到 PR 沒有明確 success criteria + baseline，要求補上
- **Runtime 層（已部分覆蓋）**：cognitive contract v2 的 `validation_mode != NONE` + 強制 `Capability summary` 是本原則的 runtime projection

## Related

- [`shotgun-debugging.md`](../anti-patterns/shotgun-debugging.md) — Agent 沒翻譯時的具體執行反模式
- [`intelligence/engineering/requirements/`](../requirements/README.md) — 鄰近：requirements cognition（功能需求方向）
- [`workflow/software-delivery/perf-risk-gate.md`](../../../workflow/software-delivery/perf-risk-gate.md) — measure → process → verify 的具體執行模板
- [`failure-recovery.md`](failure-recovery.md) — Agent 第一次 recovery attempt 最可靠，盲試會降低品質（同源觀察）

## Source

- 2026-05-27 session：使用者提供外部 infographic「Prompt 與其寫成願望，不如寫成任務單」，含 8 步網路除錯例 + 「先量測、再處理、要驗證」三原則；本 atom 將其抽象為跨工具 agent prompt-handling pattern。Status `candidate-intelligence` 至 repo 內首次套用此模板於非 software-delivery 任務後 promote 為 `validated`。
