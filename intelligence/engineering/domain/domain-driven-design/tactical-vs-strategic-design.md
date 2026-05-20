# Tactical vs Strategic Design

**Status**: `candidate-intelligence`

## 判斷原則

Strategic design 先決定 bounded context、語言、subdomain 與整合方式；tactical design 才決定 aggregate、entity、value object、repository、domain service。

## 風險

只做 tactical pattern 會讓系統看起來像 DDD，但沒有業務邊界。這通常導致 aggregate explosion、repository overuse 與 premature CQRS。

## 使用順序

```text
business capability / subdomain
→ bounded context / ubiquitous language
→ invariant / lifecycle
→ aggregate or service boundary
→ persistence / event / integration pattern
```

## 驗證

若無法說明 context boundary 與 business capability，不應直接引入 aggregate、CQRS 或 event sourcing。
