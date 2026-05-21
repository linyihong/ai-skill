# Change Brief: <title>

## Metadata
- **Change Type**: <feature|bugfix|refactor|replacement|perf|docs|chore>
- **Priority**: <p0|p1|p2|p3>
- **Evidence Source**: <issue|incident|product-brief|customer-feedback>
- **Date**: <YYYY-MM-DD>

## Evidence Summary
<Link to source evidence, max 3 bullet points>

## Product Impact Alignment
- **Impact / journey artifact**: <link to product-impact-alignment-template.md if applicable>
- **Decision**: proceed | revise | reject | ask_user | not_applicable
- **Blocking mismatch**: <goal / actor / timing / pain / investment mismatch, if any>

## Scope
### In Scope
- <item>

### Out of Scope
- <item>

## Refactor / Replacement Parity
- **Applies**: <yes|no>
- **Inventory Artifact**: <link to parity inventory or not_applicable>
- **Legacy Surfaces Covered**: <old APIs / commands / scripts / UI flows / jobs / hooks / runtime surfaces>
- **Deferred / Not Planned**: <items and why they do not block this phase>
- **Validation Evidence**: <BDD / contract test / fixture / golden output / dry-run / manual review>

## Blocker Assessment
- [ ] No blocker — proceed to Requirements / Contract phase
- [ ] Blocker identified: <description>

## Traceability
- **Spec Source**: <link to spec-template.md if from Greenfield workflow>
- **Downstream**: → Product Impact Alignment → Contract → BDD Scenarios → Implementation Plan → Review Report
- **Linked Artifacts**:
  - **Product Impact Alignment**: <link after creation>
  - **Contract**: <link after creation>
  - **BDD Scenarios**: <link after creation>
  - **Implementation Plan**: <link after creation>
