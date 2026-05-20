# Implementation Plan: <feature>

**Spec Source**: <link to spec-template.md>

## Summary
<primary requirement + technical approach>

## Technical Context
- **Language/Version**: <e.g., Python 3.11>
- **Primary Dependencies**: <e.g., FastAPI, SQLAlchemy>
- **Storage**: <e.g., PostgreSQL | Redis | files>
- **Testing**: <e.g., pytest | XCTest | cargo test>
- **Target Platform**: <e.g., Linux server | iOS 15+ | WASM>
- **Project Type**: <library | cli | web-service | mobile-app | compiler>
- **Performance Goals**: <e.g., 1000 req/s | 60 fps>
- **Constraints**: <e.g., <200ms p95 | <100MB memory>
- **Scale/Scope**: <e.g., 10k users | 1M LOC>

## Constitution Check
*GATE: Must pass before proceeding to Tasks phase.*

- [ ] Architecture consistency with existing system
- [ ] Architecture Compatibility Preflight completed per `plans/README.md`
- [ ] Dependency license compatibility
- [ ] Security & compliance requirements
- [ ] Performance targets achievable with chosen stack

## Project Structure

### Documentation
```
specs/<feature>/
├── spec.md       # Feature Specification
├── plan.md       # This file
└── tasks.md      # Task Breakdown (next phase)
```

### Source Code
```
<project structure based on project type>
```

## Complexity Tracking
> Fill ONLY if Constitution Check has violations that must be justified.

| Violation | Why Needed | Simpler Alternative Rejected |
|-----------|------------|------------------------------|
| <e.g., 4th service> | <current need> | <why 3 insufficient> |
| <e.g., Repository pattern> | <specific problem> | <why direct DB access insufficient> |
