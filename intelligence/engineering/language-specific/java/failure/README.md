# Java Failure Patterns

`intelligence/engineering/language-specific/java/failure/` stores Java-specific failure patterns and diagnostic knowledge.

## Scope

This directory is responsible for:

- Java standard library pitfalls (e.g., `String.trim()` removing tab characters)
- Java runtime behavior that differs from other languages
- Java-specific debugging techniques

## Relationship to Other Layers

- `intelligence/engineering/analytical-reasoning/failure/` stores cross-language or analysis-technique failures; this directory stores failures specific to the Java language
- `enforcement/failure-learning-system.md` defines the generic failure learning framework

## Current Atoms

| Atom | Description | Source |
|------|-------------|--------|
| [`java-tsv-trim-split-trailing-empty.md`](java-tsv-trim-split-trailing-empty.md) | Java TSV `trim()` before `split("\t", -1)` destroys trailing empty fields — `String.trim()` removes tab (0x09), silently dropping trailing empty TSV fields | `intelligence/engineering/analytical-reasoning/failure/` (migrated) |
