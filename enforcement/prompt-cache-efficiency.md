# Prompt Cache Efficiency

本規則用來在組裝 agent context 時維持穩定 prompt layout，提高 provider prompt cache 命中率，同時保留 required dependency reading、source-of-truth validation 與使用者當前目標的優先權。

## 核心規則

當任務涉及 bootstrap、context loading、routing、token cost optimization、prompt composition 或 provider cache 時，依下列順序組裝 context：

```text
stable prefix -> semi-stable middle -> volatile suffix
```

## 分區規則

| 區段 | 放什麼 | 不放什麼 |
| --- | --- | --- |
| Stable prefix | Core Bootstrap、固定 runtime initialization、低變動 routing policy | 時間戳、git status、tool output、使用者一次性要求 |
| Semi-stable middle | task intent 對應的 summary、route-specific rule、domain workflow | 無關 full source、跨路線的備用文件 |
| Volatile suffix | 使用者當前訊息、open files、diff、test output、live evidence | 任何會成為長期 stable prefix 的內容 |

## 行為要求

1. 不因單一任務把 volatile context 插入 stable prefix 中間。
2. 若新增常駐規則或固定入口，追加到 stable prefix 的固定區塊末端，並記錄改動原因。
3. 先用 summary / index 判斷是否需要 full source；需要時再展開。
4. `cacheable: true` 只表示 runtime / conversation 內可重用；不要把它直接解讀成 provider cache eligibility。
5. `provider_cache_candidate: true` 的 context 才可被視為 provider prompt cache 的候選內容。
6. 不得為了 cache hit 跳過 `dependency-reading.md`、validation gates、safety rules 或使用者明確要求。

## 何時啟用

啟用本規則的常見情境：

- 討論或修改 token cost optimization、context loading、prompt composition。
- 新增或調整 bootstrap、runtime initialization、routing registry、metadata schema。
- 設計 provider prompt cache、prefix cache、cached token 或 context reuse 策略。
- 長對話中需要 prune、compact、resume 或重新整理 context layout。

## 與其他規則的關係

- [`dependency-reading.md`](dependency-reading.md) 優先於本規則；需要讀的依賴不可因 cache 成本被省略。
- [`decision-efficiency.md`](decision-efficiency.md) 決定是否需要載入某 context；本規則只決定已選 context 的排列與 cache eligibility。
- [`tool-neutral-documentation.md`](tool-neutral-documentation.md) 仍要求可重用規則保持工具中立；工具或 provider 的具體設定放在 `ai-tools/` 或對應 adapter。
- [`runtime/context/prompt-cache-playbook.md`](../runtime/context/prompt-cache-playbook.md) 是本規則的 runtime layout playbook。

## 驗證

完成相關變更前，檢查：

- Stable prefix 內容與順序是否穩定。
- Volatile context 是否沒有進入 prefix。
- Metadata 是否區分 `cacheable` 與 `provider_cache_candidate`。
- Routing / summary / README 是否能導向 prompt cache playbook。
- Diff review 沒有把工具輸出、私有 evidence 或一次性狀態寫進 reusable docs。

← [回到 enforcement index](README.md)
