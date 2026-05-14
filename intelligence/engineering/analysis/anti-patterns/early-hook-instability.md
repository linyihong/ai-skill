# Early Hook Instability（過早 Hook 導致不穩定）

## 問題

在 Flutter/Dart AOT 分析中，過早附加 hook（特別是在 app 啟動初期）會導致 app 不穩定、crash、或 hook 失效。

## 原因

- Dart AOT snapshot 載入需要時間，relocation 未完成前 hook 會指向錯誤位址
- Global Dart runtime helpers（如 `LinkedHashMap._set`）是 hot path，broad hook 會 destabilize app
- Spawn 模式中，Frida 附加時機與 Dart VM 初始化有 race condition

## 症狀

| 症狀 | 可能原因 | 診斷方式 |
|------|---------|---------|
| App 在啟動後數秒內 crash | Hook 附加在 relocation 完成前 | 延遲 hook 附加時間（+2-3 秒） |
| Hook 從未觸發 | Function PC 錯誤或 snapshot 未載入 | 確認 libapp.so base address 是否正確 |
| App 變慢或 ANR | Broad hook on hot runtime helpers | 減少 hook 數量，只 hook app-owned functions |
| Frida 報 `Process terminated` | Spawn race condition | 改用 `-f` 搭配 `--pause` 參數 |

## 預防方式

1. **延遲 hook**：在 app 啟動後 2-3 秒再附加 hook
2. **精準 hook**：只 hook app-owned functions，不要 broad hook runtime helpers
3. **使用 `--pause`**：在 spawn 模式中使用 `frida -f <app> --pause` 讓 app 暫停直到 hook 就緒
4. **驗證 PC**：確認 function PC 是 libapp.so 的 offset，不是 callsite `BL` address

## 相關 atoms

- `intelligence/engineering/analysis/failure/frida-spawn-race.md`
- `intelligence/engineering/analysis/heuristics/hook-selection.md`

## Token 影響

低。此 atom 在 hook 設定階段 lazy-load，約 150-200 tokens。
