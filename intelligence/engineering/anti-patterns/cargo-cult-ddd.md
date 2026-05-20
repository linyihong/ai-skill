# Cargo-Cult DDD

**Status**: `candidate-intelligence`

## 反模式

在沒有 business complexity、invariant density 或 bounded context 證據時，直接套用 DDD 名詞與資料夾結構。

## 訊號

- 每個 entity 都有 aggregate / repository / domain service。
- 沒有 ubiquitous language 或 context map，卻已建立 tactical pattern。
- DDD 被用來證明抽象層，而不是保護業務不變量。
- 低複雜度 CRUD 被升級成 enterprise architecture。

## 修正

回到 architecture fit analysis。若 complexity 低，改用 CRUD、vertical slice 或 simple service layer；若 complexity 中等，採 DDD Lite。
