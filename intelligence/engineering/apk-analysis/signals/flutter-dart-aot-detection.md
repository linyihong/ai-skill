# Flutter/Dart AOT Detection Signals（Flutter/Dart AOT 辨識信號）

## 問題

如何判斷一個 APK 是否使用 Flutter/Dart AOT？需要哪些信號來確認分析路線？

## 判斷信號

### 主要信號（高可信度）

| 信號 | 檢查方式 | 可信度 |
|------|---------|-------|
| `lib/<arch>/libflutter.so` 存在 | `ls lib/<arch>/libflutter.so` | 高 |
| `lib/<arch>/libapp.so` 存在 | `ls lib/<arch>/libapp.so` | 高 |
| APK 包含 Dart snapshot | `blutter libapp.so` 可識別 snapshot | 高 |

### 次要信號（中等可信度）

| 信號 | 檢查方式 | 可信度 |
|------|---------|-------|
| Java/OkHttp hooks miss 但 pcap 有網路活動 | Frida Java hook 無輸出 vs tcpdump 有流量 | 中 |
| 靜態字串包含 `Dio`、`Interceptor` | `grep -r "Dio" libapp.so` | 中 |
| 靜態字串包含 Dart package paths | `grep -r "package:" libapp.so` | 中 |
| 靜態字串包含 `HttpClient`（dart:io） | `grep -r "HttpClient" libapp.so` | 中 |
| 靜態字串包含加密相關名稱 | `grep -r "encrypt\|decrypt\|AES" libapp.so` | 中低 |

### 進階信號（Frida 動態檢測，中高可信度）

| 信號 | 檢查方式 | 可信度 |
|------|---------|-------|
| Frida constructor chain 顯示 `AES.ctor`、`Key.ctor`、`IV.ctor` | Hook Dart `encrypt` package constructors（需先識別 libapp.so offset） | 中高 |
| Frida 顯示 `PBC.ctor`（PaddedBlockCipher） | Hook `PaddedBlockCipher` constructor | 中高 |
| Frida 顯示 `CBCBlockCipher.ctor` | Hook `CBCBlockCipher` constructor | 中高 |
| Frida 顯示 `GCMBlockCipher.ctor` | Hook `GCMBlockCipher` constructor（注意：可能屬於不同 encryption group） | 中高 |
| Frida 顯示 `processBlock` 被呼叫 43 次 | Hook `CBCBlockCipher.processBlock` / `GCMBlockCipher.processBlock` | 中（需 live proxy test 確認 mode） |

### 排除信號

| 信號 | 意義 |
|------|------|
| 只有 `libflutter.so` 但無 `libapp.so` | 可能是 Flutter debug build 或非 AOT 模式 |
| Java OkHttp hooks 可攔截所有流量 | 非 Dart AOT HTTP（可能是原生或 WebView） |
| 僅有 WebView 流量 | 非 Flutter/Dart AOT |

## 判斷流程

```
libflutter.so + libapp.so 存在？
    ├── 是 → Dart AOT 分析路線
    │       ├── blutter 可識別 snapshot → 使用 Dart Frida hook
    │       └── blutter crash → 切換到 unflutter 靜態分析
    └── 否 → 檢查其他信號
            ├── Java hooks miss + pcap 有流量 → 可能 Dart AOT（檢查字串）
            └── Java hooks 可攔截 → 非 Dart AOT
```

## 相關 atoms

- `intelligence/engineering/apk-analysis/heuristics/hook-selection.md`
- `analysis/apk/traffic-triage.md`

## Token 影響

低。此 atom 在分析初期 lazy-load，約 100-150 tokens。
