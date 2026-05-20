# Event Storming

**Status**: `candidate-intelligence`

## 判斷原則

Event storming 是發現業務流程、語言、命令、事件、policy 與 boundary 的 planning method，不是 runtime workflow，也不是必做儀式。

## 適用訊號

- 流程跨多個角色或團隊。
- 需求文件有大量狀態轉換但責任不清。
- 團隊對「發生了什麼」沒有共同語言。
- 需要找出 bounded context 或 domain events。

## 不適用訊號

- 小型 CRUD feature。
- 單人短期 prototype。
- 已有清楚 BDD 與 domain contract。

## 產出

Event storming 的可重用產出應是：事件語言、command、policy、hot spot、context boundary 與 open question；不是大量便利貼的轉錄。
