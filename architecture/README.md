# Architecture

`architecture/` 存放**每世代系統的 canonical 設計入口**。每一次世代升級（系統命名變更）需在本層建立新文件作為當前世代 navigation，舊世代文件保留為 historical。

## 用途

- 定義當前世代的 OS canonical 入口（navigation 文件，指向 executable contracts 與 philosophy 的真實 source-of-truth）
- 保留歷史世代文件，提供回溯能力
- 不包含執行計畫（執行計畫請見 [`plans/`](../plans/README.md)）
- 不包含工程架構判斷（請見 [`intelligence/engineering/architecture/`](../intelligence/engineering/architecture/README.md)）

## 目前文件

| 文件 | 世代 | 狀態 | 說明 |
|------|------|------|------|
| [`ai-native-cognitive-ecosystem-system.md`](ai-native-cognitive-ecosystem-system.md) | 4 | **vision** | 下一代願景文件：從 execution runtime 轉為 cognitive ecosystem；含 10 條 graduation threshold criteria 與當前狀態評估。**尚未 current**，Gen 3 仍是當前世代。 |
| [`ai-native-cognitive-execution-system.md`](ai-native-cognitive-execution-system.md) | 3 | **current** | 當前世代 canonical 入口：系統定位、canonical 入口表、核心機制、演化里程碑 |
| [`ai-native-knowledge-operating-system.md`](ai-native-knowledge-operating-system.md) | 2 | historical | 第二代設計文件（`skills/` 為 capability 層、Default Bootstrap 12 條） |

## 與 `intelligence/engineering/architecture/` 的邊界

| 層 | 範疇 |
|----|------|
| `architecture/`（本層） | **OS / 知識庫**架構：repo 怎麼組織、啟動、契約、世代演化 |
| [`intelligence/engineering/architecture/`](../intelligence/engineering/architecture/README.md) | **工程**架構判斷：domain modeling、modularity、coupling、選型 trade-off |

混淆會導致 reusable lesson 放錯層；新增內容前先用上表分流。

## 誰會參考這裡（Inbound References）

變更本層內容時，需要一併檢查以下依賴方：

| 來源 | 關係 |
|------|------|
| [`route.architecture.permanent-docs`](../knowledge/runtime/routing-registry.yaml) | Routing registry record，agent 依此找到 architecture/ |
| [`README.md`](../README.md) | 根目錄 OS Layout 表格列出 architecture/ 層 |
| [`constitution/`](../constitution/README.md) | ADR 可能引用 architecture 架構原則 |
| [`plans/`](../plans/README.md) | 執行計畫可能引用 architecture 設計原則 |
| [`governance/lifecycle/system-upgrade-governance.md`](../governance/lifecycle/system-upgrade-governance.md) | 升級治理規則要求世代升級時更新本層 |

## 與既有層的關係

- [`plans/`](../plans/README.md)：執行計畫存放處，完成後歸檔至 `plans/archived/`；當前世代文件引用已 archived 的關鍵升級計畫
- [`constitution/`](../constitution/README.md)：架構決策記錄（ADR），記錄為什麼做出某個架構選擇
- [`governance/lifecycle/`](../governance/lifecycle/README.md)：升級治理、各層哲學；當前世代的核心機制散落於此
- [`README.md`](../README.md)：根目錄 README 已列出 architecture 作為 OS Layout 的一層
