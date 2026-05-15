# Language-Specific Intelligence

`intelligence/engineering/language-specific/` stores knowledge that is specific to a particular programming language. This is a separate dimension from `analytical-reasoning/` (which stores cross-language analytical techniques).

## Directory Structure

```
language-specific/
├── README.md           # This file
├── java/
│   ├── README.md
│   └── failure/        # Java-specific failure patterns
└── ...                 # Other languages as needed
```

## Scope

This directory is responsible for:

- Language-specific failure patterns (e.g., Java `String.trim()` behavior, JavaScript bitwise operator truncation)
- Language-specific techniques and idioms
- Language runtime quirks and pitfalls

This directory does NOT replace:

- `analytical-reasoning/failure/` — cross-language or analysis-technique failures belong there
- `analytical-reasoning/heuristics/` — cross-language heuristics belong there
- `enforcement/failure-patterns/` — cross-skill failure patterns belong there

## When to Add a New Language

When a language accumulates at least 1 validated failure pattern or technique, create a subdirectory:

```
language-specific/<language>/
├── README.md
└── failure/        # Failure patterns for this language
└── techniques/     # (future) Techniques for this language
```

## Current Languages

| Language | Directory | Atoms |
|----------|-----------|-------|
| Java | [`java/`](java/README.md) | 1 failure pattern |
