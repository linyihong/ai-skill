# evidence_type: `source_contract`

## Proves

Static alignment between specification and implementation: Gherkin scenario text, source asserts, schema fields, or documented contract strings.

## Non-goals

- Runtime browser behavior
- User-visible layout
- DB row existence alone without contract mapping

## Supported collection_methods

- `contract_readback`
- `static_analysis`

## Supported artifact_shapes

- `source_diff`
- `string_match_log`
- `schema_assertion_log`

## Proxy traps

- BDD string assert pass **without** executable validation → behavior gap, not validation complete
- Copy-pasted contract refs with no test runner → `source_contract` only

## Example claim

`preview_overlay_described_in_feature`
