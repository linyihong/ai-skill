# Metadata Compatibility

`metadata/compatibility/` defines how new knowledge layers preserve existing skill, shared-rule, tool, and script behavior while migration is in progress.

## Compatibility Fields

Use these `metadata/schema.md` fields to express compatibility:

| Field | Use |
| --- | --- |
| `source_path` | Canonical old source path that remains authoritative. |
| `depends` | Required old entrypoints or shared rules. |
| `related` | Candidate new layer paths or supporting references. |
| `replaces` | Only set after promotion or deprecation approval. |
| `conflicts` | Potential rule or entrypoint conflicts that require resolution. |
| `governance_notes` | Migration state, compatibility notes, or deprecation requirements. |

## Compatibility States

| State | Meaning |
| --- | --- |
| `old-entrypoint-active` | Existing `skills/` or `shared-rules/` source remains active. |
| `dual-reference` | Old entrypoint and new layer path are both linked for discovery. |
| `new-layer-promoted` | New layer path is supported, but old path still resolves. |
| `deprecation-planned` | Old path has replacement and deprecation note, but still exists. |
| `old-entrypoint-retired` | Old path is removed or archived after validation and replacement. |

## Required Compatibility Notes

For any candidate map or promoted atom, record:

- Old entrypoint.
- New reference path.
- Whether old entrypoint remains active.
- Whether tool-specific discovery still depends on the old path.
- What validation proves links and routing still work.

## Blocking Conditions

Do not promote or deprecate when:

- A tool still loads only the old skill path and no adapter exists.
- The old entrypoint was not read after the change.
- `knowledge/indexes/README.md` points only to a candidate path and omits the active source.
- A shared rule says the old source is canonical and no rule update has been made.
- Link check or close-loop validation fails.

## Reference-First Default

Compatibility metadata should prefer direct canonical repository references. Tool mirrors, bundles, copied snapshots, and local runtime paths are deployment surfaces, not source paths.
