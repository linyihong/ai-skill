# Stale Acceptance Criteria

**Status**: `candidate-intelligence`

## 判斷原則

Acceptance criteria 是 validation baseline。若 product intent、domain invariant、API contract 或 implementation truth 改變，acceptance criteria 可能過期。

## 訊號

- Test 綠，但 product / BDD 描述舊行為。
- Bug fix 改變了預期行為，卻沒有更新 scenario。
- Domain invariant 已調整，acceptance criteria 仍使用舊狀態。

## 行動

標記 stale baseline，重跑 requirement alignment。這可作為 runtime-lite 候選 signal，但不代表 runtime 理解 BDD syntax。
