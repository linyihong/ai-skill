# Domain-Driven Design Intelligence

本目錄保存 DDD 的業務模型判斷智慧。它不是教學百科，也不是 runtime rule；用途是幫 agent 在軟體交付時判斷何時需要 DDD、需要多少 DDD，以及何時應該降級為更簡單的架構。

## 核心原則

DDD 是 selectable architecture strategy，不是 universal doctrine。

使用 DDD 前先問：

| 問題 | 判斷 |
| --- | --- |
| 業務規則是否密集？ | 若只有 CRUD，不需要 full DDD。 |
| 不變量是否關鍵且跨流程？ | 若是，才需要 aggregate / bounded context 深化。 |
| 語言是否不穩定？ | 若詞彙變動頻繁，先建立 ubiquitous language。 |
| 系統是否長期演化？ | 短期 MVP 優先 DDD Lite 或 simple service layer。 |
| 整合壓力是否高？ | 高外部模型污染風險才需要 anti-corruption layer。 |

## 目前條目

| 文件 | 用途 |
| --- | --- |
| [`bounded-context.md`](bounded-context.md) | 判斷業務語言與模型邊界。 |
| [`aggregate-design.md`](aggregate-design.md) | 判斷 aggregate 是否保護真正的不變量。 |
| [`ubiquitous-language.md`](ubiquitous-language.md) | 建立跨產品、工程與測試一致的語言。 |
| [`domain-events.md`](domain-events.md) | 判斷 domain event 是否代表業務事實。 |
| [`anti-corruption-layer.md`](anti-corruption-layer.md) | 防止外部模型污染內部領域模型。 |
| [`domain-services.md`](domain-services.md) | 判斷何時需要 domain service。 |
| [`repository-pattern.md`](repository-pattern.md) | 判斷 repository 是否代表 aggregate persistence boundary。 |
| [`tactical-vs-strategic-design.md`](tactical-vs-strategic-design.md) | 區分 tactical pattern 與 strategic modeling。 |
| [`ddd-lite-vs-full-ddd.md`](ddd-lite-vs-full-ddd.md) | 在 DDD Lite 與 Full DDD 間選擇。 |
| [`event-storming.md`](event-storming.md) | 使用 event storming 發現流程與語言邊界。 |
| [`overengineering-risks.md`](overengineering-risks.md) | 偵測 DDD 過度套用。 |
| [`architecture-fit-signals.md`](architecture-fit-signals.md) | 將 DDD adoption signal 接到 architecture fit analysis。 |

## 與其他層的關係

- `intelligence/engineering/architecture/architecture-selection/` 負責選架構，不把 DDD 當預設。
- `workflow/software-delivery/architecture/` 負責實際執行 architecture fit analysis。
- `governance/ai-runtime-governance/software-delivery-architecture-governance.md` 負責把 architecture minimality 與 overengineering review 轉成治理 gate。
- `metadata/architecture/` 保存可機讀 heuristics，但預設 metadata-only。
- `runtime/` 不保存 DDD tactical rules；只有壓縮後的 reliability signal 才可能另案 promotion。

## 驗證

使用本目錄後，agent 應能說明：為什麼選 CRUD / DDD Lite / Full DDD、哪些 DDD pattern 被拒絕、以及拒絕是否基於 business complexity 而非個人偏好。
