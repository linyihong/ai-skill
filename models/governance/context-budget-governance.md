# Context Budget Governance

Context budget governance 定義 model-aware routing 如何控制長任務的讀取深度與 replay 成本。

## Budget Order

1. Index / registry。
2. Summary / checklist。
3. Primary source。
4. Required dependencies。
5. Related graph / memory candidate。
6. Raw transcript 或 broad replay（預設禁止）。

## Escalation Triggers

- 修改 canonical source。
- Runtime / generated surface touched。
- Evidence conflict。
- User correction。
- Long context 或 compaction。
- Validation signal 不足。

## Model / Memory Alignment

Model context budget 與 memory replay budget 必須一致：若 memory replay 會造成 context inflation，優先讀 current canonical source 或 summary-first route。
