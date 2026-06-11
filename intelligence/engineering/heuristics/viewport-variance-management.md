# Viewport Variance Management

Heuristic: design and validate against render constraints, not device names.

## Rule

When a UI bug is described with a device name, translate it into the underlying render context:

- width constraint;
- height constraint;
- orientation or aspect ratio;
- safe-area requirement;
- dynamic resize or browser chrome behavior;
- embedded container or wrapper constraint.

## Why

Device names are useful test fixtures, but poor governance vocabulary. A render-context name stays valid across emulator presets, browser zoom, split-screen, native wrappers, and future devices.

## Smells

- The only validation target is a device preset such as `small_phone` or `tablet_model`.
- A breakpoint is named after a device SKU rather than a layout constraint.
- A viewport resize bug is treated as a one-device CSS patch.

## Better Target

Use shared render contexts such as `narrow_mobile`, `mobile`, `tablet`, `desktop`, `landscape`, `safe_area`, and `dynamic_resize`.
