# Render Context Library

Render contexts are reusable names for UI constraints. They are not device models.

## Contexts

| Context | Use When |
| --- | --- |
| [`narrow-mobile.md`](narrow-mobile.md) | Very constrained phone-like width or height can affect layout survival. |
| [`mobile.md`](mobile.md) | Standard single-column touch layout is in scope. |
| [`large-mobile.md`](large-mobile.md) | Wider phone or compact phablet layout may reveal density assumptions. |
| [`tablet.md`](tablet.md) | Medium viewport can support more density but still touch-oriented constraints. |
| [`desktop.md`](desktop.md) | Wide viewport or pointer-oriented layout is in scope. |
| [`landscape.md`](landscape.md) | Height or aspect ratio changes can affect content reachability. |
| [`safe-area.md`](safe-area.md) | System UI, notch, home indicator, or webview insets affect usable space. |
| [`dynamic-resize.md`](dynamic-resize.md) | Browser chrome, orientation, devtools emulation, split view, or reload changes dimensions. |

## Rule

Use render-context IDs in contracts and validation plans. Device SKUs may be test fixtures, but they should not be the governance vocabulary.
