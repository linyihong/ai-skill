# Framework Pattern Extraction

Heuristic: extract the stable design pattern from a mature framework before copying its implementation surface.

## Rule

When using Bootstrap or another mature UI framework as evidence:

- extract the transferable contract: container, grid, spacing scale, breakpoints, display utilities, overflow handling, or accessibility convention;
- decide whether the project should adopt the framework, emulate the pattern locally, or keep the existing system;
- avoid mixing framework classes into a mature local design system unless ownership, bundle cost, and style precedence are clear;
- record the generalized primitive in reusable knowledge, not the project-specific CSS patch.

## Smells

- Adding a full framework to fix one layout bug.
- Copying class names without adopting the framework's reset, grid assumptions, or accessibility conventions.
- Keeping separate local constants that disagree with the framework-like contract.
