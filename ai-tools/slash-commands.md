# Slash Commands

## 格式
`/skill.<command> [arguments]`

Agent 在收到 `/skill.<command>` 時，解析 command 名稱，對應到 [`routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) 的 route，載入對應 workflow 與模板，根據 arguments 或 conversation context 填寫模板後輸出結構化文件。

## 指令列表

### Greenfield Workflow
對應 [`workflow/greenfield/`](../workflow/greenfield/README.md) 的 4 階段流程。

| 指令 | 功能 | 對應階段 | 產出模板 | 範例 |
|------|------|---------|---------|------|
| `/skill.spec` | 建立 Feature Specification | Specify | [`spec-template.md`](../workflow/greenfield/templates/spec-template.md) | `/skill.spec 建立使用者登入功能` |
| `/skill.plan` | 建立 Implementation Plan | Plan | [`plan-template.md`](../workflow/greenfield/templates/plan-template.md) | `/skill.plan`（需先有 spec） |
| `/skill.tasks` | 建立 Task Breakdown | Tasks | [`tasks-template.md`](../workflow/greenfield/templates/tasks-template.md) | `/skill.tasks`（需先有 plan） |

### Software Delivery Workflow
對應 [`workflow/software-delivery/`](../workflow/software-delivery/README.md) 的變更管理流程。

| 指令 | 功能 | 對應階段 | 產出模板 | 範例 |
|------|------|---------|---------|------|
| `/skill.brief` | 建立 Change Brief | Change Intake | [`change-brief-template.md`](../workflow/software-delivery/templates/change-brief-template.md) | `/skill.brief 修復登入頁面 500 錯誤` |
| `/skill.contract` | 建立 Contract | Contract Governance | [`contract-template.md`](../workflow/software-delivery/templates/contract-template.md) | `/skill.contract`（需先有 change brief） |
| `/skill.bdd` | 建立 BDD Scenarios | BDD Closure | [`bdd-scenario-template.md`](../workflow/software-delivery/templates/bdd-scenario-template.md) | `/skill.bdd`（需先有 contract） |

## 實作方式

Agent 收到 `/skill.<command>` 時的處理流程：

1. **解析指令**：從 `/skill.<command>` 中提取 command 名稱
2. **路由查詢**：對應到 [`routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) 的 `slash_commands` 欄位
3. **載入 Workflow**：讀取對應 workflow 的 `execution-flow.md`
4. **載入模板**：讀取對應的模板檔案
5. **填寫模板**：根據 arguments 或 conversation context 填寫模板內容
6. **輸出**：輸出結構化文件，並記錄 traceability

## 路由對應

```yaml
# 定義在 routing-registry.yaml 中
route.workflow.greenfield:
  slash_commands:
    - /skill.spec
    - /skill.plan
    - /skill.tasks

route.workflow.software-delivery:
  slash_commands:
    - /skill.brief
    - /skill.contract
    - /skill.bdd
```

## 注意事項

- 部分指令有前置依賴（如 `/skill.plan` 需先有 spec），Agent 應檢查後提示使用者
- 指令可帶 arguments 作為上下文輸入，也可不帶（從 conversation history 推斷）
- 產出文件應包含 Traceability 欄位，記錄上下游關係
