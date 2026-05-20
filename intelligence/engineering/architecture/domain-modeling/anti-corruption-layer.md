# Anti-Corruption Layer

**Status**: `candidate-intelligence`

## 判斷原則

Anti-corruption layer（ACL）保護內部 domain model 不被外部模型、供應商 API、legacy schema 或其他 bounded context 的語言污染。

## 適用訊號

- 外部模型命名與內部業務語言衝突。
- 外部狀態轉換不能直接套用到內部 invariant。
- 第三方 API 變更會擴散到核心 domain code。
- 需要同時支援多個外部 provider。

## 不適用訊號

- 只是簡單 DTO mapping。
- 系統沒有穩定內部 domain model。
- 外部模型本身就是 source of truth 且沒有語意轉換。
- 為了「乾淨架構」而包每個 dependency。

## 最小做法

先在 adapter boundary 做 translation table、狀態 mapping 與錯誤語意 mapping。只有污染風險持續出現時才擴大 ACL。
