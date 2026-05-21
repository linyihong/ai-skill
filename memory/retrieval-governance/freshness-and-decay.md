# Freshness And Decay

Memory 的 confidence 會隨時間、architecture change、workflow change、source supersession 與 context boundary 改變。Memory freshness 不等於 current truth。

## Defaults

| Memory type | Default confidence | Freshness handling |
| --- | --- | --- |
| Transcript-derived memory | low | 必須重新驗證。 |
| Episodic memory | tentative | 只能作 weak guidance。 |
| Summary memory | scoped | 只恢復 session context，不證明 current truth。 |
| Failure abstraction | medium | 可提示 risk，但 current source 仍需檢查。 |
| Project memory | scoped | 受 repo architecture、migration、dependency changes 影響。 |
| Decision memory | medium-high | Status / supersession 檢查後可用。 |

## Suggested Metadata

長期可 replay 的 memory 應逐步補齊：

```yaml
last_validated:
expires_when:
compatibility_scope:
confidence_default:
replay_allowed_as:
superseded_by:
```

## Decay Triggers

- Repo architecture refactor。
- Workflow 或 enforcement rule 更新。
- Runtime / generated surface migration。
- User goal 改變。
- Decision status 從 `accepted` 變成 `superseded` / `deprecated`。
- Memory 來源是 compacted summary 而非 current source。

## Claim Scope

Memory 可支持的 claim scope 不得超過它的 qualification scope。若 memory 只證明過去曾發生，不能聲稱現在仍成立；若 memory 只屬於某 project，不能推成跨 project rule。
