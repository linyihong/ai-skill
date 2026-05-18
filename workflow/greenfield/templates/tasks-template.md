# Tasks: <feature name>

**Input**: <link to plan-template.md>
**Prerequisites**: spec.md (required), plan.md (required)

## Format
`[ID] [P?] [Story] Description — <file path>`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Maps to user story (US1, US2, US3)

---

## Phase 1: Setup (Shared Infrastructure)
**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan
- [ ] T002 [P] Initialize <language> project with <framework> dependencies
- [ ] T003 [P] Configure linting and formatting tools

---

## Phase 2: Foundational (Blocking Prerequisites)
**Purpose**: Core infrastructure that MUST be complete before ANY user story

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T004 Setup database schema and migrations framework
- [ ] T005 [P] Implement authentication/authorization framework
- [ ] T006 [P] Setup API routing and middleware structure
- [ ] T007 Create base models/entities that all stories depend on
- [ ] T008 Configure error handling and logging infrastructure

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 — <title> (Priority: P1) 🎯 MVP
**Goal**: <brief description>
**Independent Test**: <how to verify this story works on its own>

### Tests (OPTIONAL — only if tests requested) ⚠️
- [ ] T010 [P] [US1] Contract test — tests/contract/test_<name>.py
- [ ] T011 [P] [US1] Integration test — tests/integration/test_<name>.py

### Implementation
- [ ] T012 [P] [US1] Create <Entity> model — src/models/<entity>.py
- [ ] T013 [US1] Implement <Service> — src/services/<service>.py
- [ ] T014 [US1] Implement endpoint — src/api/<endpoint>.py
- [ ] T015 [US1] Add validation and error handling
- [ ] T016 [US1] Add logging for user story 1 operations

**Checkpoint**: User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 — <title> (Priority: P2)
**Goal**: <brief description>
**Independent Test**: <how to verify this story works on its own>

### Tests (OPTIONAL — only if tests requested) ⚠️
- [ ] T020 [P] [US2] Contract test — tests/contract/test_<name>.py
- [ ] T021 [P] [US2] Integration test — tests/integration/test_<name>.py

### Implementation
- [ ] T022 [P] [US2] Create <Entity> model — src/models/<entity>.py
- [ ] T023 [US2] Implement <Service> — src/services/<service>.py
- [ ] T024 [US2] Implement endpoint — src/api/<endpoint>.py
- [ ] T025 [US2] Integrate with User Story 1 components (if needed)

**Checkpoint**: User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 — <title> (Priority: P3)
**Goal**: <brief description>
**Independent Test**: <how to verify this story works on its own>

### Tests (OPTIONAL — only if tests requested) ⚠️
- [ ] T030 [P] [US3] Contract test — tests/contract/test_<name>.py
- [ ] T031 [P] [US3] Integration test — tests/integration/test_<name>.py

### Implementation
- [ ] T032 [P] [US3] Create <Entity> model — src/models/<entity>.py
- [ ] T033 [US3] Implement <Service> — src/services/<service>.py
- [ ] T034 [US3] Implement endpoint — src/api/<endpoint>.py

**Checkpoint**: All user stories should now be independently functional

---

## Phase N: Polish & Cross-Cutting Concerns
**Purpose**: Improvements that affect multiple user stories

- [ ] TXXX [P] Documentation updates — docs/
- [ ] TXXX Code cleanup and refactoring
- [ ] TXXX Performance optimization across all stories
- [ ] TXXX [P] Additional unit tests (if requested) — tests/unit/
- [ ] TXXX Security hardening
- [ ] TXXX Run validation checks

---

## Dependencies & Execution Order

### Phase Dependencies
- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies
- **User Story 1 (P1)**: Can start after Foundational — No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational — May integrate with US1
- **User Story 3 (P3)**: Can start after Foundational — May integrate with US1/US2

### Within Each User Story
- Tests (if included) MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration

### Parallel Opportunities
- All tasks marked [P] can run in parallel
- Once Foundational phase completes, all user stories can start in parallel
- Different user stories can be worked on in parallel by different team members

---

## Implementation Strategy

### MVP First (User Story 1 Only)
1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery
1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo

### Parallel Team Strategy
With multiple developers:
1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes
- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
