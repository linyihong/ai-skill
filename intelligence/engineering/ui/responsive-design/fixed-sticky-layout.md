# Fixed And Sticky Layout

Principle: fixed and sticky surfaces must align to the same layout contract as the content they control.

## Rule

- Bottom bars, sticky headers, floating actions, and fixed docks should derive width and inline position from the page/container contract.
- If content is capped by `max-width`, fixed surfaces should be capped by the same source of truth.
- Reserve enough scroll or safe-area space so fixed surfaces do not cover primary content.
- Prefer CSS positioning tied to viewport/container rules before adding JavaScript measurement.

## Smells

- Fixed footer width is computed separately from the content shell.
- Sticky controls drift after viewport resize or browser zoom.
- A fixed surface overlaps the final list item, CTA, or form submit button.
- JavaScript `visualViewport` measurement is used to paper over duplicated CSS width constants.
