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

## Worked Example：網路速度優化

來自 2025 infographic「想讓網路變快，不要靠感覺」，作為 wish → task list 翻譯的具體範例。展示同一個 wish 在 agent 兩種處理路徑下的差異。

### User wish

> 「想讓網路快一點。」

### Wish 信號判斷

- 動詞「想」「快」→ 模糊
- 「快」未量化（多快？哪個指標？）
- 無 success criteria
- 無 baseline 引用
- → 4/4 wish 信號命中，**必須翻譯為 task list 才能執行**

### Agent 翻譯後的 task list（measure → process → verify）

**Measure（baseline）**

| 指標 | 工具 | 目的 |
|---|---|---|
| 下載 / 上傳速度 | `speedtest-cli` | 量化目前頻寬 |
| DNS 解析時間 | `dig` / `nslookup` | 抓 DNS layer 延遲 |
| MTU / 封包遺失 | `ping -s` / `mtr` | 抓 L2-L3 健康度 |
| Wi-Fi 訊號 / 干擾 | OS 工具 | 抓無線層問題 |

→ 完成後手上有 4 個 baseline 數字，且使用者已確認「目標」（例如下載 > 100Mbps、ping < 30ms）。

**Process（單變數變更，最多 3 個假設）**

依 baseline 找到 3 個最可能的瓶頸（不要 8 個全部一次改 —— 那是 shotgun debugging）：

1. 移除舊的 / 衝突的 network profiles
2. 關閉吃頻寬的背景程式
3. 優化 mDNS（若有 .local 解析慢的訊號）

每個變更**分別執行**，每次只動一個變數。

**Verify（after baseline + 對照）**

跑回 measure 階段同樣的 4 個指標，逐項對照 before / after。

### 反例：Agent 直接執行 wish

```
User: 「想讓網路快一點」
Agent 直接動手：
  → 重設 DNS
  → 重啟 Wi-Fi
  → 升級驅動
  → 移除 profile
  → 關閉背景程式
  → 改 MTU
  → 重開 router
  → 跑 speedtest
Result: 「我做完了，你試試看」
```

問題：
- 沒 baseline → 不知道「快多少才算修好」
- 沒單變數變更 → 不知道哪個動作真的有效
- 沒對照 → 下次再慢無法從同樣的指標開始查
- → 命中 [`shotgun-debugging`](../anti-patterns/shotgun-debugging.md)

### 為什麼這個 example 適合放這裡

1. **跨領域**：infographic 本身是 IT 維護領域，但 wish-to-task-list translation 是 agent 通用 pattern，這個 example 證明 atom 不只適用 software-delivery
2. **明顯的 wish 訊號**：4 個 wish 信號全中，是教學範例的好材料
3. **measure / process / verify 三段清楚**：infographic 把流程顯式畫出來，可直接搬到本 atom 不需再加工

如果未來 repo 真的需要服務多種 troubleshooting 任務（網路 / 效能 / 系統設定 / 硬體調校），再把本段抽到 `workflow/troubleshooting/measure-process-verify.md` + 多個 worked examples。目前以 incremental 為主，不開新 workflow domain。

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
