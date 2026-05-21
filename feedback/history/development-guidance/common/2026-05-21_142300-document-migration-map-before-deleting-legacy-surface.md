> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-21 - Document Migration Map Before Deleting Legacy Surface

Status: validated

#### One-line Summary

Legacy script deletion must be preceded by a developer-facing migration map that names the new Go owner, source-of-truth, side effects, and validation evidence.

#### Human Explanation

Updating scattered docs after an implementation is not enough for a migration. Other developers need a single handoff map that explains where each old capability moved, what file they should edit next time, and which tests prove the old surface can stay deleted.

#### Trigger

An agent is about to delete, port, or declare closure for an old shell, Ruby, Python, or tool adapter surface.

#### Evidence

- Tool: repository documentation and validation review.
- Sanitized excerpt: A migration removed legacy runtime surfaces but initially lacked a single developer map from old entries to new Go CLI commands and implementation owners.
- Evidence path: `<AI_SKILL_REPO>/scripts/ai-skill-cli/docs/legacy-to-go-migration-map.md`

#### Generalized Lesson

Before deleting a legacy surface, write the migration map first. The map must identify the old surface, new command or package owner, source-of-truth file, side effects, and validation gate. Only then should implementation, deletion, parity docs, and disposition docs be finalized.

#### Agent Action

For migration work:

1. Create or update the migration map before deleting the old surface.
2. Link the map from parity, command contract, and disposition docs.
3. Add a BDD or fixture gate that prevents future deletion without the map.
4. Run validation after docs and code agree.

#### Goal / Action / Validation

- Goal: Preserve developer handoff clarity during migration closure.
- Action: Add a migration map and gate it from CLI runtime docs.
- Validation or reference source: Diff review confirms links from CLI docs; runtime validation and Go tests confirm the migrated behavior remains covered.

#### Applies When

- A legacy entrypoint is deleted or converted to a Go-native command.
- A script has side effects such as writing files, mutating Git, updating runtime DB, syncing tool settings, or touching generated artifacts.
- The replacement is expected to be maintained by other developers after the current session.

#### Does Not Apply When

- The change is a local typo fix or refactor with no entrypoint migration.
- The legacy surface remains fully retained and no replacement is being declared.

#### Validation

Check that the migration map exists, the old surface has a row, and linked docs point to it before deletion is committed.

#### Promotion Target

- `scripts/ai-skill-cli/docs/legacy-to-go-migration-map.md`
- `scripts/ai-skill-cli/docs/bdd-scenarios.md`
- `scripts/ai-skill-cli/docs/test-fixture-plan.md`

#### Required Linked Updates

- Updated CLI docs index, command contract, parity inventory, legacy disposition, BDD scenarios, and fixture plan.
- No category README existed under `feedback/history/development-guidance/common/`; domain README count is updated instead.
