# Frida Spawn Race Condition（Frida Spawn 競爭條件）

## 問題

使用 `frida -f <app>` spawn 模式時，Frida 附加時機與 Dart VM 初始化之間存在 race condition，導致 hook 失效或 app crash。

## 症狀

| 症狀 | 發生時機 | 頻率 |
|------|---------|------|
| Frida 報 `Process terminated` | Spawn 後立即 | 高 |
| Hook 從未觸發但 app 正常運行 | Spawn 後 1-3 秒 | 中 |
| App crash 在 Dart VM 初始化階段 | Spawn 後 0.5-2 秒 | 中高 |
| Frida 報 `unable to find target process` | Spawn 後立即 | 低 |

## 診斷方式

```bash
# 使用 --pause 參數暫停 app 啟動
frida -f <app.bundle.id> --pause -l script.js

# 手動 resume
%resume
```

如果 `--pause` 模式正常但普通 spawn 失敗，確認是 race condition。

## 緩解方式

1. **使用 `--pause`**：讓 app 暫停直到 hook script 就緒
2. **延遲 hook**：在 script 中加入 `setTimeout(() => { ... }, 2000)` 延遲附加
3. **改用 attach 模式**：先啟動 app，再用 `frida -n <app>` attach
4. **檢查 Frida 版本**：較舊版本（<16.0）的 spawn race 更常見

## 相關 atoms

- `intelligence/engineering/analytical-reasoning/anti-patterns/early-hook-instability.md`
- `intelligence/engineering/analytical-reasoning/heuristics/hook-selection.md`

## Token 影響

低。此 atom 在 spawn 失敗時 lazy-load，約 100-150 tokens。
