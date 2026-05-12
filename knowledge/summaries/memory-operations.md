## memory.operations

| 欄位 | 值 |
| --- | --- |
| Atom ID | `memory.operations` |
| Source path | `memory/README.md` |
| Lifecycle | `candidate` |
| Summary | 長期記憶層：short-term（目前 task context）、episodic（過去 task 關鍵決策與結果）、project（專案歷史脈絡）、failure（反覆失效模式）。支援 similarity-based retrieval。 |
| When to read | 需要參考過去類似任務的決策、查詢專案歷史、或避免重複失效模式時。 |
| Do not use for | 不可取代 knowledge/ 的結構化知識。不可用於儲存可從 canonical source 重建的暫態資料。 |
| Context cost | ~250 tokens |
| Estimated full cost | ~1000 tokens |
| Validation signal | Memory layer README 可解析，episodic/project/failure 路徑存在。 |
| Last checked | 2026-05-12 |
