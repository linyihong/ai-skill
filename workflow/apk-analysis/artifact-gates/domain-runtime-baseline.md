# Domain/Runtime Baseline Slice

> **Cognitive Slice**：`apk-domain-runtime-baseline`（從 [`../artifact-gates.md`](../artifact-gates.md) §3 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-domain-runtime-baseline` |
| `purpose` | 記錄環境維度 / 連線路徑 / session / 簽章 / 分頁等 runtime baseline，避免「只做 per-API shape」卻不能連線 |
| `type` | `execution` |
| `tags` | artifact-gate, domain |
| `load_when` | 需要建立 domain model / runtime baseline、SDK/client/live integration 開發前的 development readiness 檢查 |
| `do_not_load_when` | 純 stateless 端點分析、已有完整 baseline 的小修改 |
| `owner_layer` | workflow |
| `layer_justification` | 規定「baseline 要含哪些章節、何時不能只停在 skeleton」的 ordering / artifact gate；通過 workflow membership test |
| `canonical_source` | 本檔（原 `artifact-gates.md` §3 Domain/Runtime Baseline） |
| `dependencies` | `apk-api-catalog`、`apk-self-generation-audits`（若涉 live SDK） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | development-readiness gate 觸發時應載入本 slice |

## 3. Domain/Runtime Baseline

**問題：**只做「逐支 API 的 request/response shape」時，外包裝/SDK 仍可無法連線。

### 建議章節

| 章節 | 內容（去敏） |
| --- | --- |
| 環境維度 | 觀察到的 host family、path family、是否多 CDN／多 gateway、與 build／地區是否相關。 |
| 連線路徑 | App 是否走系統代理、內建 TUN、local proxy、直连；與 capture 工具相容性。 |
| Session／身分 | 列表 API 是否在未登入下可用；若否，登入／裝置／device id 與列表欄位的因果鏈。 |
| Opaque／衍生參數 | 哪些 query 由前序 response、WebView、搜尋 session 或固定 app 常數提供。 |
| 簽章與 gateway | service／hash、header 名稱集合、canonical path 規則。 |
| 分頁地面真相 | 是否有 `has_next` 類欄位；若無，記錄啟發式與反例風險。 |
| 錯誤與限流 | 影響重試的 code、冷卻、與 session 刷新關係。 |
| 重放檢查清單 | 人工或腳本重放同一列表的最小步驟。 |

### Development readiness gate

若下一步是 SDK/client/app tool/live integration 開發，baseline 不能只停在 skeleton。必須先檢查並記錄最小可跑因素：

- endpoint/path family
- route/service 對照或 adapter 策略
- session/bootstrap 依存
- opaque 參數來源與時效
- 簽章/gateway 前置
- response decrypt/unwrap 邊界
- 分頁地面真相
- 錯誤/session 恢復
- 重放檢查清單

缺任一項時，該缺口必須成為開發 blocker 或被明確 scoped out。

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
