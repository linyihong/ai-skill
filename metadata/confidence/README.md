# Metadata Confidence

`metadata/confidence/` 定義 Knowledge Atoms 與 routing surfaces 的證據強度標記方式。

## 信心值

| 值 | 意義 | 可用情境 |
| --- | --- | --- |
| `low` | 合理但尚未在使用中驗證。 | Candidate maps、早期 summaries、尚未測試的 atom proposals。 |
| `medium` | 已有 source review、link validation 或一次成功使用支撐。 | Navigation rows、candidate atoms、pilot maps。 |
| `high` | 已重複使用或 review，且有清楚 validation evidence。 | Validated atoms、stable routing、promoted guidance。 |

## 與狀態的關係

| Lifecycle status | 最低 confidence | 備註 |
| --- | --- | --- |
| `temporary` | `low` | 短期或 project-local；不要當作 durable knowledge 索引。 |
| `candidate` | `low` | 可被 routing 使用，但必須標明 candidate。 |
| `validated` | `medium` | 需要真實使用、review 或明確 validation record。 |
| `stable` | `high` | 需要重複使用或強 review evidence。 |
| `deprecated` | any | 必須附 replacement 或 reason。 |

## 證據信號

符合下列情況時，confidence 可以提高：

- Markdown links 與 source paths 可解析。
- Atom 已在完成的任務中使用。
- Reviewer 或 user 接受此 guidance。
- Test、fixture、lint、link check 或 close-loop validation 通過。
- 舊 source-of-truth entrypoints 更新後，此 atom 仍保持對齊。

符合下列情況時，confidence 應保持低：

- Atom 只是 planning guess。
- 舊 skill 仍在變動，且尚無 synchronization rule。
- Validation 缺失或 blocked。
- Guidance 依賴尚未泛化的 project-specific evidence。

## 降級條件

符合下列情況時，降級或標示 stale：

- 舊 `skills/` source 已變更，而 atom 尚未重新檢查。
- Links 失效。
- 與 `shared-rules/` 或 source-of-truth skill behavior 出現衝突。
- Candidate path 被誤認為 replacement path。

## 驗證

每個 `confidence: high` 的 atom 應包含或連到：

- Source path。
- Validation method。
- Last known lifecycle state。
- Compatibility 或 deprecation note。
