# API Catalog Slice（API 文件結構 + 單支 API 詳細需求 + 模板）

> **Cognitive Slice**：`apk-api-catalog`（從 [`../artifact-gates.md`](../artifact-gates.md) §2+§11 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-api-catalog` |
| `purpose` | 紀錄 API endpoint、request/response shape、authentication 細節；維護 catalog 完整性 gate |
| `type` | `execution` |
| `tags` | artifact-gate, api |
| `load_when` | 整理 API 列表 / SDK/client 對照 / mock API / contract test / 功能重建 |
| `do_not_load_when` | 不涉及 API 文件化的純行為觀察、純 UI 探索 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「API 文件要含哪些欄位、catalog 何時算完整」的 ordering / artifact gate；通過 workflow membership test |
| `canonical_source` | 本檔（原 `artifact-gates.md` §2 API Catalog + §11 API Catalog Detail Requirements） |
| `dependencies` | `apk-ui-architecture-map`（UI binding 章節相互引用）、`apk-evidence-chain`（evidence 證據鏈）、`apk-sanitization` |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Scenario AG-A（API documentation 任務應 PRIMARY 載入本 slice） |

## 2. API Catalog

當分析目標包含「整理 API 列表、SDK/client 對照、mock API、contract test、功能重建」時，專案文件應建立一組可維護的 API Catalog。

### 建議結構

```text
docs/API.md                         # API / host / traffic family 總入口
docs/API/<group>/README.md          # 依 path 第一段、domain、feature 或 protocol family 分組
docs/API/<group>/<operation>.md      # 單支 API 詳細文件
docs/API/coverage.md                # 已觀測、已 replay、已解密、待補與不適用清單
docs/API/supplement/<topic>.md      # HLS、media、local bridge、SDK call order 等跨 API 主題
```

### Catalog minimum

| Artifact | 必填內容 |
| --- | --- |
| API 總入口 | 已知 host/base、traffic family、response wrapper、auth/session/header 共用規則、解密/解碼入口、覆蓋率文件、UI map 入口、SDK/client 入口。 |
| 分組索引 | 分組依據、每支 API 的 method/path、request 摘要、response 用途、目前用途/結論、詳細文件連結。 |
| 單支 API 文件 | Method/path、host/base、auth/session、headers、query/path/body、response wrapper、inner payload、item schema、error/empty behavior、pagination/cache、field meaning、sensitivity、evidence、validation。 |
| Coverage / gap matrix | 靜態枚舉、動態觀測、MITM、pcap、hook、replay、decrypted fixture、contract test、UI binding、缺參、未觸發、暫不測、scope out。 |
| UI/API 對照 | UI map、route id、operation id、trigger confidence、capture window、startup/preload/background 標記。 |
| SDK/client 欄位用途 | SDK/client 實際讀取或轉換的欄位、相容性範圍、raw JSON 保留策略、fixture/test 對照。 |
| Cross-flow docs | 播放鏈、media chain、login/session、local bridge、vendor/service split、分頁與排序等跨多支 API 的流程。 |
| Sanitization | 哪些值已遮蔽、哪些 raw evidence 留在受控位置、哪些文件可 commit。 |

### Catalog completion gate

當使用者問「API 列表是否完整」、「能不能做 SDK/client/mock」時，完成回覆前要檢查：

- API 總入口是否連到分組、coverage/gap、UI map、解碼/共用 wrapper、SDK/client 文件。
- 已觀測 API 是否都落到分組索引。
- 高價值 API 是否有單支詳細文件。
- 每支 API 是否有 request、response、field meaning、evidence、validation/open questions。
- UI trigger 若未確認，是否標 `UI path: unknown` / `Trigger confidence: low`。

## 11. API Catalog Detail Requirements

單支 API 文件至少要能回答：

| Area | Required detail |
| --- | --- |
| Identity | Method、host/base、path shape、operation id、分組、狀態：candidate / observed / replayed / decoded / validated / deprecated / out of scope。 |
| Request | headers、path/query/body 欄位、型別/shape、用途、必填/選填、來源、敏感性、是否參與 signing/encryption。 |
| Response | raw wrapper、decrypted/inner payload、list item schema、欄位型別、nullable/optional、欄位語意、derived-from、下游 API key。 |
| Behavior | capability、UI trigger、startup/preload/background 判斷、state impact、empty/error behavior、pagination/cache/sort semantics。 |
| Evidence | hook/MITM/pcap/replay/fixture/screenshot/UI hierarchy/automation script 的去敏引用。 |
| Validation | replay result、decoder fixture、schema assertion、SDK/client test、contract test、manual evidence，或明確標 `needs capture` / `needs replay`。 |
| Open questions | 缺少樣本、低信心 field meaning、未驗證 edge case、需要使用者或更多操作證據的 blocker。 |

### Field meaning rule

Schema 不只是型別表。欄位要盡量寫出用途：

- 哪些欄位會成為下一支 API 的 request key。
- 哪些欄位控制 UI 顯示、分頁停止、排序、播放、下載、收藏、權限或錯誤狀態。
- 哪些欄位只是樣本中出現但用途未知，必須標 `meaning unknown` / `candidate`。
- 哪些欄位是 SDK/client 已使用欄位，變動會破壞相容性。

不要把推測寫成確認規格。若只有少量 live/replay 樣本，使用 `candidate`、`sample only`、`needs more samples` 或 `low confidence`。

### API / Schema Document Template

```markdown
## Endpoint Name

| Field | Value |
| --- | --- |
| Method | `GET` / `POST` |
| Path | `/path` |
| Auth | Required / Optional |
| Source | pcap / MITM / hook / replay |
| UI path | `Tab > Screen > Action` |
| Operation ID | `open-home` / `open-detail` |
| Trigger confidence | high / medium / low |
| Capability / feature | user-visible function this API supports |
| Domain concept candidates | entity/value object/state names inferred from evidence |
| State impact | creates / reads / updates / deletes / refreshes / paginates / authenticates |

### HTTP Request Headers

| Header | Type / Shape | Meaning | Required | Source | Sensitive | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `Authorization` | bearer / custom / none | session auth | yes/no | token provider | yes | value redacted |
| `User-Agent` | string shape only | client identity | yes/no | app/runtime | no | |

### Request Query / Path Parameters

| Field | Type / Shape | Meaning | Required | Source | Sensitive | Notes |
| --- | --- | --- | --- | --- | --- | --- |

### Request Body

| Field | Type / Shape | Meaning | Required | Source | Sensitive | Notes |
| --- | --- | --- | --- | --- | --- | --- |

### Response Wrapper

| Field | Type / Shape | Meaning | Required / Optional | Notes |
| --- | --- | --- | --- | --- |

### Decrypted / Inner Payload

| Field | Type / Shape | Meaning | Required / Optional | Source / Derived From | Notes |
| --- | --- | --- | --- | --- | --- |

### Response Headers

| Header | Type / Shape | Meaning | Notes |
| --- | --- | --- | --- |

### Evidence

- Sanitized log:
- Fixture:
- UI path:
- Screenshot / UI evidence:

### Validation

- Replay:
- Contract test:
- Manual verification:

### Reconstruction Notes

- BDD scenario candidate:
- Domain Model Contract candidates:
- API / Interface Contract notes:
- Error Handling Contract notes:
- Fixtures needed for rebuild:
- Open questions:
```

API 文件要求：

- 分析完 API 後要回填專案文件；不要只把 endpoint 留在暫存 log。
- HTTP/HTTPS API 必須記錄可見的 headers、request、response；看不到的部分要寫明是 MITM 不可見、hook 未到位、加密包裹、或尚未驗證。
- 每個 request/response 字段都要逐欄位分析 type/shape、meaning、required/optional、source/derived-from、敏感性與備註。
- 每個高價值 API 都要標明支援哪個 capability、對應 operation id、可能的 domain concept、狀態影響、錯誤/空狀態與 fixture。
- Header 名稱、path shape、query key、schema 可以保留；header value、token、cookie、device id、個資與可重放 URL 必須去敏。
- 截圖可用來輔助說明 UI path、tab、screen 與操作，但不能取代 HTTP header/request/response 的字段分析。

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
