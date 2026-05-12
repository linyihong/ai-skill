# Review Checklists

本文件定義開發流程中的審查檢查清單。承接 [`skills/app-development-guidance/checklists/`](../../skills/app-development-guidance/checklists/) 的內容，提取為 tool-neutral 的 workflow gates。

> **相容性規則**：`skills/app-development-guidance/checklists/` 仍為 active skill entrypoint。本文件為 reference target，兩者應保持同步。

## 檢查清單類型

| 清單 | 使用時機 | 原始來源 |
|------|----------|----------|
| **Mobile Design Review** | Before implementing a mobile feature or security-sensitive flow | [`checklists/mobile-design-review.md`](../../skills/app-development-guidance/checklists/mobile-design-review.md) |
| **Mobile PR Review** | During code review | [`checklists/mobile-pr-review.md`](../../skills/app-development-guidance/checklists/mobile-pr-review.md) |
| **Mobile Release Review** | Before shipping a mobile release | [`checklists/mobile-release-review.md`](../../skills/app-development-guidance/checklists/mobile-release-review.md) |
| **API Security Review** | When mobile/web clients depend on API security properties | [`checklists/api-security-review.md`](../../skills/app-development-guidance/checklists/api-security-review.md) |
| **Contract Governance Review** | When multiple planning, BDD, contract, generated, and test docs must stay traceable | [`checklists/contract-governance-review.md`](../../skills/app-development-guidance/checklists/contract-governance-review.md) |
| **Embedded Firmware Review** | When firmware, sensors, boards, protocols, or hardware-in-loop validation are involved | [`checklists/embedded-firmware-review.md`](../../skills/app-development-guidance/checklists/embedded-firmware-review.md) |

## 使用原則

1. **Keep checklists short enough to run during real development** — 檢查清單必須在實際開發中可執行。
2. **Checklist items must stay linked to implementation docs** — 檢查項目必須連結到它們要求審查者驗證的實作文件。
3. **When adding a check, update or verify implementation and control docs** — 新增檢查項目時，在同一變更中更新或驗證對應的 implementation 和 control 文件。

## 與其他層的關係

- `workflow/app-development-guidance/execution-flow.md` 提供執行流程，本文件提供流程中的審查門檻。
- `analysis/app-development-guidance/controls-catalog.md` 提供檢查清單引用的控制原則。
- `skills/app-development-guidance/checklists/` 是原始來源，仍為 active entrypoint。
