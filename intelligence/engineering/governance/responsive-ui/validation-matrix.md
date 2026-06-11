# Responsive Validation Matrix

Use this matrix to decide which evidence is required before claiming responsive behavior is valid.

## Minimum Shape

```yaml
responsive_validation:
  required_contexts:
    - desktop
    - mobile
  add_when_relevant:
    - narrow_mobile
    - tablet
    - landscape
    - safe_area
    - dynamic_resize
  evidence:
    - viewport_metrics
    - document_overflow_metrics
    - shell_bounding_box
    - fixed_surface_bounding_box
    - primary_content_bounding_box
    - screenshot_or_dom_snapshot
```

## Rules

- A single desktop capture cannot prove responsive validity.
- At least one wide and one constrained context are required for a responsive completion claim.
- Use render-context vocabulary from [`render-contexts/`](../../render-contexts/README.md).
- Add `safe_area`, `landscape`, or `dynamic_resize` when the defect report or UI contract depends on them.
- Record intended horizontal scrollers separately from unintended document overflow.

## Evidence Ownership

- Workflow decides when to collect evidence.
- Validation execution collects and evaluates evidence.
- Responsive UI Governance defines the expected evidence shape and failure authority.
