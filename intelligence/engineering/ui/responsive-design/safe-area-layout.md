# Safe-Area Layout

Principle: mobile UI must reserve space for system UI and physical display constraints.

## Rule

- Use safe-area environment values for fixed bottom or top surfaces when the target includes mobile webviews or modern phones.
- Add safe-area space to interactive surfaces, not just decorative backgrounds.
- Validate that primary CTAs and navigation controls remain visible and tappable when safe areas are present.
- Do not assume a rectangular viewport with no home indicator, notch, or browser chrome.

## Smells

- Fixed bottom controls use a hard-coded bottom inset.
- Final scroll content is hidden under a home indicator or tab bar.
- Safe-area padding is applied to the shell but not to the fixed surface that actually overlaps system UI.
