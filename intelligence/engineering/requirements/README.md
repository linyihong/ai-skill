# Requirements Cognition Intelligence

`requirements/` 保存需求理解與行為建模的工程判斷智慧。這一層不把 BDD 縮成 Gherkin 或 Cucumber，而是把 BDD 視為 requirement correctness 的 cognition system。

## 核心

BDD = behavior alignment and requirement cognition system。

本層處理：

- ambiguity detection
- product impact alignment
- actor intent modeling
- scenario framing
- acceptance boundaries
- behavior scope governance
- requirement traceability
- stale acceptance criteria

## 目前入口

| 子目錄 | 用途 |
| --- | --- |
| [`product-alignment/`](product-alignment/README.md) | 用 Impact Map × Customer Journey Map 檢查 business goal、target actor、journey pain 與 feature investment 是否對齊。 |
| [`behavior-modeling/`](behavior-modeling/README.md) | 建立 shared language、behavior boundary、scenario framing 與 acceptance boundary。 |
| [`specification-quality/`](specification-quality/README.md) | 檢查 requirement 是否可驗證、可追溯、沒有自行擴張。 |
| [`validation-thinking/`](validation-thinking/README.md) | 將 acceptance criteria 接到 validation target 與 proof acquisition。 |

## 與其他層的關係

- `workflow/software-delivery/requirements/` 負責實際 delivery flow。
- `governance/ai-runtime-governance/software-delivery-governance.md` 負責把 requirement / behavior gate 轉成 software-delivery governance。
- `intelligence/engineering/architecture/domain-modeling/` 接收 requirements cognition 產生的 shared language 與 behavior boundary，再建立 domain boundary / invariant。
- `runtime/` 不理解 BDD syntax，只接收 requirement contradiction、missing validation target、stale acceptance criteria 等壓縮信號。
