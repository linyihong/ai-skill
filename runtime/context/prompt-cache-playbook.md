# Provider Prompt Cache Alignment

本 playbook 定義 agent 組裝 context 時的 prompt layout。目標是提高 provider 端 prefix cache 命中率，同時不犧牲 required dependency reading、source-of-truth validation 或使用者當前目標。

## 核心原則

Provider prompt cache 依賴穩定前綴。每次任務都會變的內容若插入前綴中間，會讓後續穩定內容失去 cache reuse 機會。

將 context 分成三段：

```text
stable prefix -> semi-stable middle -> volatile suffix
```

## Context 分區

| 區段 | 放置內容 | 穩定性要求 | 範例 |
| --- | --- | --- | --- |
| `stable prefix` | 所有任務都會用到、順序固定、低變動的規則與 runtime 入口 | 不因任務重排；新增項目只追加到固定區塊末端 | Core Bootstrap、root layout、固定 runtime initialization、stable routing policy |
| `semi-stable middle` | 本任務需要、但仍可跨相似任務重用的 context | 依 task intent 選擇；同一路線內順序固定 | knowledge summary、route-specific enforcement rule、model compression checklist |
| `volatile suffix` | 高變動或一次性內容 | 永遠放後段；不得插入 stable prefix 中間 | 使用者當前要求、open files、git status、tool output、時間戳、live evidence |

## Stable Prefix 規則

1. Stable prefix 的順序必須固定。若新增穩定項目，追加到該區塊末端，不在中間插入。
2. Stable prefix 只放低變動、跨任務共用、可重讀也不改變語意的 context。
3. 不把 task-specific evidence、工具輸出、時間戳、使用者臨時偏好或工作樹狀態放入 stable prefix。
4. 若 required dependency reading 要求讀取某檔案，即使該檔不適合 cache，也必須讀取；prompt cache 不能覆蓋正確性規則。

## Semi-Stable Middle 規則

Semi-stable middle 用於 task-specific routing 後的可重用 context。它可以依任務增減，但同一路線內應保持固定順序：

1. 先放 route summary，再放 canonical source 或 enforcement rule。
2. 先放 index / summary，再放 full source。
3. 先放 cross-task rule，再放 domain-specific workflow / intelligence。
4. 當 summary 已足夠回答低風險問題，不展開 full source。

## Volatile Suffix 規則

Volatile suffix 放所有高變動輸入：

- 使用者當前訊息與最新優先序。
- IDE 附加狀態，例如 open files、recent files、lint diagnostics。
- git status、diff、test output、tool output。
- live evidence、專案私有資料、一次性分析片段。
- 時間戳、session-local note、臨時 todo。

若 volatile context 需要影響後續任務，應沉澱成 summary、goal ledger、feedback lesson 或 durable plan，再由對應 route 載入；不要直接把原始 volatile payload 升格成 stable prefix。

## Metadata 對應

`metadata/schema.md` 中的 provider cache hints 用來描述 context 適合放在哪一段：

| 欄位 | 用途 |
| --- | --- |
| `provider_cache_candidate` | 此 atom 是否適合 provider prefix cache reuse。 |
| `prefix_stability` | `stable`、`semi_stable` 或 `volatile`。 |
| `cache_position` | 建議位置：`prefix`、`middle` 或 `suffix`。 |
| `churn_risk` | 前綴變動風險：`low`、`medium`、`high`。 |

`cacheable: true` 仍表示 runtime / conversation 內可重用；`provider_cache_candidate: true` 才表示適合納入 provider prompt cache 的穩定布局。兩者不可混用。

## Observability

第一版先以人工檢查與 runtime metadata 記錄下列信號：

| 指標 | 用途 |
| --- | --- |
| `stable_prefix_size` | 估算穩定前綴 token 大小，避免前綴過重。 |
| `prefix_churn_count` | 記錄 stable prefix 在同一類任務中被改動的次數。 |
| `provider_cache_candidate_tokens` | 估算可被 provider cache 重用的 token。 |
| `volatile_prefix_violation_count` | 偵測高變動內容被放入 prefix 的次數。 |

## 檢查清單

- Stable prefix 是否只包含低變動、跨任務共用 context？
- 使用者當前要求、git status、tool output 是否都在 suffix？
- 新增常駐規則是否追加到固定區塊末端，而不是插入中間？
- `cacheable` 是否沒有被誤解成 provider cache eligibility？
- 是否沒有為了提高 cache hit 跳過 required dependencies 或 validation gates？

← [回到 context/](README.md)
