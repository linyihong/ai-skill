# BDD Scenario: <feature/behavior>

## Requirement Link
- **Source**: <product brief / issue / contract link>
- **Actor / system role**: <actor>
- **Behavior boundary**: <in scope / out of scope>
- **Ambiguity disposition**: validated | assumption | open question | scoped out | invalidated

## Scenario: <title>
<!-- If this is rendered as a .feature file, keep these as Gherkin comments above the Scenario. -->
- **Test ref**: <tests/path>::<test/scenario/checklist name>
- **Code ref**: <implementation / contract / schema / adapter path>

**Given** <precondition>
**When** <action>
**Then** <expected outcome>

## Scenario: <title>
- **Test ref**: <tests/path>::<test/scenario/checklist name>
- **Code ref**: <implementation / contract / schema / adapter path>

**Given** <precondition>
**When** <action>
**Then** <expected outcome>

## Edge Cases
- <edge case 1>
- <edge case 2>

## Acceptance Criteria
- <observable acceptance criterion>
- <negative or failure criterion>

## Validation Target
- **Proof type**: automated | fixture-backed | manual-evidence | pending-runner | not-automatable
- **Test / fixture / checklist**: <path or owner>
- **Code / contract owner path**: <implementation / contract / schema / adapter path>
- **Limitations**: <what this does not prove>

## Regression Scope
- [ ] Existing tests affected: <list>
- [ ] New tests required: <count>
- [ ] Test data/fixtures needed: <description>

## Traceability
- **Upstream**: <link to product brief / behavior contract / contract-template.md>
- **BDD -> tests**: <test refs, fixture refs, checklist refs, or pending-runner todo>
- **BDD -> code / contracts**: <implementation, API, schema, adapter, UI, command, data migration refs>
- **Downstream**: → Implementation Plan → Review Report
